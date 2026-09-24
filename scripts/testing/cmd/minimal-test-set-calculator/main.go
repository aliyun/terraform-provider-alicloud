//nolint:all
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/aliyun/terraform-provider-alicloud/scripts/testing/internal/common"
	mapset "github.com/deckarep/golang-set"
	log "github.com/sirupsen/logrus"
)

func init() {
	customFormatter := new(log.TextFormatter)
	customFormatter.FullTimestamp = true
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	customFormatter.DisableTimestamp = false
	customFormatter.DisableColors = false
	customFormatter.ForceColors = true
	log.SetFormatter(customFormatter)
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

var (
	resourceTypeName = flag.String("resource", "", "resource name (e.g., drds_polardbx_instance)")
	outputFormat     = flag.String("format", "summary", "output format: summary, json, detail")
)

type TestCase struct {
	Name               string
	CoveredAttributes  mapset.Set
	ModifiedAttributes mapset.Set
	HasImportTest      bool
	Steps              int
	LineNumber         int
}

type MinimalTestSetResult struct {
	ResourceName          string                    `json:"resource_name"`
	TotalTestCases        int                       `json:"total_test_cases"`
	MinimalTestSet        []string                  `json:"minimal_test_set"`
	AttributeCoverage     float64                   `json:"attribute_coverage"`
	ModificationCoverage  float64                   `json:"modification_coverage"`
	TotalAttributes       int                       `json:"total_attributes"`
	CoveredAttributes     int                       `json:"covered_attributes"`
	TotalModifyAttributes int                       `json:"total_modify_attributes"`
	CoveredModifyAttr     int                       `json:"covered_modify_attributes"`
	UncoveredAttributes   []string                  `json:"uncovered_attributes,omitempty"`
	UncoveredModifyAttr   []string                  `json:"uncovered_modify_attributes,omitempty"`
	AllTestCases          map[string]TestCaseDetail `json:"all_test_cases"`
}

type TestCaseDetail struct {
	CoveredAttributes       []string `json:"covered_attributes"`
	ModifiedAttributes      []string `json:"modified_attributes"`
	HasImportTest           bool     `json:"has_import_test"`
	Steps                   int      `json:"steps"`
	LineNumber              int      `json:"line_number"`
	IsSelected              bool     `json:"is_selected"`
	NewAttributesCovered    int      `json:"new_attributes_covered"`
	NewModificationsCovered int      `json:"new_modifications_covered"`
}

func main() {
	flag.Parse()

	if *resourceTypeName == "" {
		log.Fatal("Error: -resource parameter is required")
	}

	// Remove alicloud_ prefix if present
	resName := strings.TrimPrefix(*resourceTypeName, "alicloud_")

	log.Infof("Analyzing test cases for resource: %s", resName)

	// Get schema attributes using common function
	schemaAllSet, schemaMustSet, schemaModifySet, _ := common.GetSchemaAttributes(resName, true)

	if schemaMustSet.Cardinality() == 0 {
		log.Fatal("Error: Resource not found or has no attributes")
	}

	log.Infof("Total schema attributes: %d", schemaAllSet.Cardinality())
	log.Infof("Must-set attributes: %d", schemaMustSet.Cardinality())
	log.Infof("Modifiable attributes: %d", schemaModifySet.Cardinality())

	// Parse test file using common function
	testFilePath := fmt.Sprintf("alicloud/resource_alicloud_%s_test.go", resName)
	testMustSet := mapset.NewSet()
	testModifySet := mapset.NewSet()
	testIgnoreSet := mapset.NewSet()

	testCaseInfoMap, success := common.ParseTestFile(testFilePath, &testMustSet, &testModifySet, &testIgnoreSet)

	if !success || len(testCaseInfoMap) == 0 {
		log.Fatal("Error: No test cases found or failed to parse test file")
	}

	// Convert TestCaseInfo to TestCase for minimal set calculation
	testCases := make(map[string]*TestCase)
	for name, info := range testCaseInfoMap {
		tc := &TestCase{
			Name:               name,
			CoveredAttributes:  info.CoveredAttributes,
			ModifiedAttributes: info.ModifiedAttributes,
			HasImportTest:      info.HasImportTest,
			Steps:              info.Steps,
			LineNumber:         info.LineNumber,
		}
		testCases[name] = tc
	}

	log.Infof("Found %d test cases", len(testCases))

	// Calculate minimal test set using greedy algorithm (100% coverage)
	minimalSet := calculateMinimalTestSet(testCases, schemaMustSet, schemaModifySet)

	// Calculate coverage
	result := buildResult(resName, testCases, minimalSet, schemaMustSet, schemaModifySet)

	// Output result
	outputResult(result, *outputFormat)

	// Exit with error if 100% coverage is not met
	if result.AttributeCoverage < 1.0 || result.ModificationCoverage < 1.0 {
		log.Errorf("Failed to achieve 100%% coverage!")
		os.Exit(1)
	}

	log.Infof("✓ Minimal test set calculated successfully with 100%% coverage!")
}

// calculateMinimalTestSet uses a greedy algorithm to achieve 100% coverage
func calculateMinimalTestSet(testCases map[string]*TestCase, mustSet, modifySet mapset.Set) []string {
	minimalSet := []string{}
	coveredAttrs := mapset.NewSet()
	coveredModifyAttrs := mapset.NewSet()

	remainingTests := make(map[string]*TestCase)
	for k, v := range testCases {
		remainingTests[k] = v
	}

	totalAttrs := mustSet.Cardinality()
	totalModifyAttrs := modifySet.Cardinality()

	log.Infof("Starting greedy algorithm for 100%% coverage...")
	log.Infof("Target: %d attributes and %d modifiable attributes", totalAttrs, totalModifyAttrs)

	iteration := 0
	for len(remainingTests) > 0 {
		iteration++

		// Calculate current coverage
		attrCoverage := float64(coveredAttrs.Intersect(mustSet).Cardinality()) / float64(totalAttrs)
		modifyCoverage := 0.0
		if totalModifyAttrs > 0 {
			modifyCoverage = float64(coveredModifyAttrs.Intersect(modifySet).Cardinality()) /
				float64(totalModifyAttrs)
		}

		// Check if 100% coverage is met
		if attrCoverage >= 1.0 && (totalModifyAttrs == 0 || modifyCoverage >= 1.0) {
			log.Infof("100%% coverage achieved after %d iterations", iteration)
			break
		}

		// Find the best test case to add
		bestTest := ""
		bestScore := -1.0

		for name, tc := range remainingTests {
			// Calculate incremental coverage
			newAttrs := tc.CoveredAttributes.Difference(coveredAttrs).Intersect(mustSet)
			newModifyAttrs := tc.ModifiedAttributes.Difference(coveredModifyAttrs).Intersect(modifySet)

			// Score based on new coverage (weighted equally)
			attrScore := float64(newAttrs.Cardinality())
			modifyScore := float64(newModifyAttrs.Cardinality())
			score := attrScore*0.5 + modifyScore*0.5

			// Bonus for import tests if we don't have one yet
			hasImport := false
			for _, selected := range minimalSet {
				if testCases[selected].HasImportTest {
					hasImport = true
					break
				}
			}
			if !hasImport && tc.HasImportTest {
				score += 10.0
			}

			if score > bestScore && (newAttrs.Cardinality() > 0 || newModifyAttrs.Cardinality() > 0) {
				bestScore = score
				bestTest = name
			}
		}

		if bestTest == "" {
			log.Warnf("No more tests can improve coverage")
			break
		}

		// Add the best test
		minimalSet = append(minimalSet, bestTest)
		tc := remainingTests[bestTest]
		coveredAttrs = coveredAttrs.Union(tc.CoveredAttributes)
		coveredModifyAttrs = coveredModifyAttrs.Union(tc.ModifiedAttributes)
		delete(remainingTests, bestTest)
	}

	return minimalSet
}

func buildResult(resName string, allTests map[string]*TestCase, minimalSet []string,
	mustSet, modifySet mapset.Set) *MinimalTestSetResult {

	result := &MinimalTestSetResult{
		ResourceName:   resName,
		TotalTestCases: len(allTests),
		MinimalTestSet: minimalSet,
		AllTestCases:   make(map[string]TestCaseDetail),
	}

	// Calculate coverage
	coveredAttrs := mapset.NewSet()
	coveredModifyAttrs := mapset.NewSet()

	for name, tc := range allTests {
		isSelected := false
		for _, selected := range minimalSet {
			if selected == name {
				isSelected = true
				break
			}
		}

		// Build attribute lists
		attrList := make([]string, 0, tc.CoveredAttributes.Cardinality())
		for attr := range tc.CoveredAttributes.Iter() {
			attrList = append(attrList, attr.(string))
		}
		sort.Strings(attrList)

		modifyList := make([]string, 0, tc.ModifiedAttributes.Cardinality())
		for attr := range tc.ModifiedAttributes.Iter() {
			modifyList = append(modifyList, attr.(string))
		}
		sort.Strings(modifyList)

		// Calculate new coverage if this test is added
		newAttrsCovered := tc.CoveredAttributes.Difference(coveredAttrs).Intersect(mustSet).Cardinality()
		newModifyCovered := tc.ModifiedAttributes.Difference(coveredModifyAttrs).Intersect(modifySet).Cardinality()

		detail := TestCaseDetail{
			CoveredAttributes:       attrList,
			ModifiedAttributes:      modifyList,
			HasImportTest:           tc.HasImportTest,
			Steps:                   tc.Steps,
			LineNumber:              tc.LineNumber,
			IsSelected:              isSelected,
			NewAttributesCovered:    newAttrsCovered,
			NewModificationsCovered: newModifyCovered,
		}

		result.AllTestCases[name] = detail

		if isSelected {
			coveredAttrs = coveredAttrs.Union(tc.CoveredAttributes)
			coveredModifyAttrs = coveredModifyAttrs.Union(tc.ModifiedAttributes)
		}
	}

	result.TotalAttributes = mustSet.Cardinality()
	result.CoveredAttributes = coveredAttrs.Intersect(mustSet).Cardinality()
	result.AttributeCoverage = float64(result.CoveredAttributes) / float64(result.TotalAttributes)

	result.TotalModifyAttributes = modifySet.Cardinality()
	result.CoveredModifyAttr = coveredModifyAttrs.Intersect(modifySet).Cardinality()
	if result.TotalModifyAttributes > 0 {
		result.ModificationCoverage = float64(result.CoveredModifyAttr) /
			float64(result.TotalModifyAttributes)
	} else {
		result.ModificationCoverage = 1.0
	}

	// Uncovered attributes
	uncovered := mustSet.Difference(coveredAttrs)
	for attr := range uncovered.Iter() {
		result.UncoveredAttributes = append(result.UncoveredAttributes, attr.(string))
	}
	sort.Strings(result.UncoveredAttributes)

	uncoveredModify := modifySet.Difference(coveredModifyAttrs)
	for attr := range uncoveredModify.Iter() {
		result.UncoveredModifyAttr = append(result.UncoveredModifyAttr, attr.(string))
	}
	sort.Strings(result.UncoveredModifyAttr)

	return result
}

func outputResult(result *MinimalTestSetResult, format string) {
	switch format {
	case "json":
		outputJSON(result)
	case "detail":
		outputDetail(result)
	default:
		outputSummary(result)
	}
}

func outputSummary(result *MinimalTestSetResult) {
	fmt.Printf("\n=== Minimal Test Set for %s ===\n\n", result.ResourceName)

	fmt.Printf("Total test cases: %d\n", result.TotalTestCases)
	fmt.Printf("Minimal test set: %d (%.1f%% reduction)\n",
		len(result.MinimalTestSet),
		(1.0-float64(len(result.MinimalTestSet))/float64(result.TotalTestCases))*100)

	fmt.Printf("\nCoverage:\n")
	fmt.Printf("  Attributes:     %d/%d (%.1f%%)\n",
		result.CoveredAttributes, result.TotalAttributes, result.AttributeCoverage*100)
	fmt.Printf("  Modifications:  %d/%d (%.1f%%)\n",
		result.CoveredModifyAttr, result.TotalModifyAttributes, result.ModificationCoverage*100)

	if len(result.UncoveredAttributes) > 0 {
		fmt.Printf("\nUncovered Attributes (%d):\n", len(result.UncoveredAttributes))
		for _, attr := range result.UncoveredAttributes {
			fmt.Printf("  - %s\n", attr)
		}
	}

	if len(result.UncoveredModifyAttr) > 0 {
		fmt.Printf("\nUncovered Modifications (%d):\n", len(result.UncoveredModifyAttr))
		for _, attr := range result.UncoveredModifyAttr {
			fmt.Printf("  - %s\n", attr)
		}
	}

	fmt.Printf("\nSelected Test Cases:\n")
	for i, testName := range result.MinimalTestSet {
		fmt.Printf("  %d. %s\n", i+1, testName)
	}

	fmt.Printf("\nRun command:\n")
	fmt.Printf("  make test-resource-debug RESOURCE=alicloud_%s TESTCASE=\"%s\"\n",
		result.ResourceName, strings.Join(result.MinimalTestSet, "|"))
	fmt.Println()
}

func outputJSON(result *MinimalTestSetResult) {
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal result: %v", err)
	}
	fmt.Println(string(data))
}

func outputDetail(result *MinimalTestSetResult) {
	fmt.Printf("\n=== Test Case Coverage Analysis for %s ===\n\n", result.ResourceName)

	// Sort test cases: selected first, then by new coverage contribution
	type testEntry struct {
		name   string
		detail TestCaseDetail
	}
	var entries []testEntry
	for name, detail := range result.AllTestCases {
		entries = append(entries, testEntry{name, detail})
	}
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].detail.IsSelected != entries[j].detail.IsSelected {
			return entries[i].detail.IsSelected
		}
		return entries[i].detail.NewAttributesCovered+entries[i].detail.NewModificationsCovered >
			entries[j].detail.NewAttributesCovered+entries[j].detail.NewModificationsCovered
	})

	// Print table header
	fmt.Printf("%-60s | %8s | %5s | %6s | %5s | %6s | %8s | %8s\n",
		"Test Case", "Selected", "Attrs", "Modify", "Steps", "Import", "New Attr", "New Mod")
	fmt.Println(strings.Repeat("-", 140))

	// Print each test case
	for _, entry := range entries {
		name := entry.name
		tc := entry.detail

		selected := "No"
		if tc.IsSelected {
			selected = "Yes"
		}

		importMark := "No"
		if tc.HasImportTest {
			importMark = "Yes"
		}

		fmt.Printf("%-60s | %8s | %5d | %6d | %5d | %6s | %8d | %8d\n",
			truncate(name, 60), selected,
			len(tc.CoveredAttributes), len(tc.ModifiedAttributes),
			tc.Steps, importMark,
			tc.NewAttributesCovered, tc.NewModificationsCovered)
	}

	fmt.Println()
	fmt.Printf("Summary: %d/%d test cases selected for 100%% coverage\n",
		len(result.MinimalTestSet), result.TotalTestCases)
	fmt.Println()
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
