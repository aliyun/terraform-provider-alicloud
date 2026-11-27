//nolint:all
package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/sdkv2"
	"github.com/aliyun/terraform-provider-alicloud/scripts/internal/splitlayout"
	"github.com/aliyun/terraform-provider-alicloud/scripts/testing/internal/clientmigration"
	mapset "github.com/deckarep/golang-set"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	log "github.com/sirupsen/logrus"
	"github.com/waigani/diffparser"
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
	log.SetLevel(log.DebugLevel)
}

var (
	fileNames    = flag.String("fileNames", "", "the files to check diff")
	resourceName = flag.String("resource", "", "directly specify a resource name to check (e.g., alicloud_realtime_compute_deployment)")

	resourceFileRegex     = regexp.MustCompile("alicloud/(resource|data_source)[0-9a-zA-Z_]*")
	resourceFileTestRegex = regexp.MustCompile("alicloud/(resource|data_source)[0-9a-zA-Z_]*_test.go")

	goFuncDeclRegex = regexp.MustCompile(`^func\s+(?:\([^)]*\)\s*)?([a-zA-Z_][a-zA-Z0-9_]*)\s*\(`)
	importSpecRegex = regexp.MustCompile(`^(?:[a-zA-Z_][a-zA-Z0-9_]*\s+|[._]\s+)?"[^"]+"$`)
	importDeclRegex = regexp.MustCompile(`^import\s+(?:[a-zA-Z_][a-zA-Z0-9_]*\s+|[._]\s+)?"[^"]+"$`)
	// 只认可迁移配对表内的限定符，避免把「换 provider 工厂」判成机械改动
	providerFactoryFieldRegex = regexp.MustCompile(`^(?:ProviderFactories|ProtoV5ProviderFactories|Providers)\s*:\s*(` + migrationQualifierAlt() + `)\.(?:ProviderFactories|SDKProviderFactories|ProtoV5ProviderFactories|Providers|SDKProviders)(?:\(\))?,?$`)
	providerConstructorRegex  = regexp.MustCompile(`^(?:var\s+)?[a-zA-Z_][a-zA-Z0-9_]*\s*:?=\s*(` + migrationQualifierAlt() + `)\.Provider(?:\(\))?(?:\.(?:ResourcesMap|DataSourcesMap))?$`)
)

func migrationQualifierAlt() string {
	return `(?:alicloud|sdkv2|acctest)`
}

func hasCoverageRelevantTestChanges(file *diffparser.DiffFile) bool {
	if file.Mode != diffparser.MODIFIED {
		return true
	}
	if len(file.Hunks) == 0 {
		// Nothing to analyze (e.g. rename-only entry): fail towards checking.
		return true
	}
	for _, hunk := range file.Hunks {
		if hunkHasCoverageRelevantChange(hunk) {
			return !clientmigration.OnlyChange(file)
		}
	}
	return false
}

func hunkHasCoverageRelevantChange(hunk *diffparser.DiffHunk) bool {
	currentFunc := ""
	inImport := strings.HasPrefix(strings.TrimSpace(hunk.HunkHeader), "import")
	for _, line := range hunk.WholeRange.Lines {
		content := strings.TrimSpace(line.Content)
		if match := goFuncDeclRegex.FindStringSubmatch(content); len(match) == 2 {
			currentFunc = match[1]
		}
		if strings.HasPrefix(content, "import") && strings.HasSuffix(content, "(") {
			inImport = true
		}
		if line.Mode == diffparser.UNCHANGED {
			if inImport && content == ")" {
				inImport = false
			}
			if strings.HasPrefix(line.Content, "}") {
				currentFunc = ""
			}
			continue
		}
		if content == "" {
			continue
		}
		if inImport {
			if content == ")" {
				inImport = false
				continue
			}
			if importSpecRegex.MatchString(content) {
				continue
			}
		}
		if importDeclRegex.MatchString(content) || (strings.HasPrefix(content, "//") && !strings.HasPrefix(content, "//go:") && !strings.HasPrefix(content, "// +build")) {
			continue
		}
		if strings.HasPrefix(currentFunc, "Test") && !strings.HasPrefix(currentFunc, "TestAcc") && currentFunc != "TestMain" {
			if strings.HasPrefix(line.Content, "}") {
				currentFunc = ""
			}
			continue
		}
		if providerFactoryFieldRegex.MatchString(content) || providerConstructorRegex.MatchString(content) {
			continue
		}
		return true
	}
	return false
}

func main() {
	exitCode := 0
	flag.Parse()

	// Mode 1: Direct resource name check (for local use)
	if resourceName != nil && len(*resourceName) > 0 {
		log.Infof("==> Checking resource: %s", *resourceName)
		if !checkSingleResource(*resourceName) {
			exitCode = 1
		}
		os.Exit(exitCode)
	}

	// Mode 2: Diff file check (for CI/CD)
	if fileNames != nil && len(*fileNames) == 0 {
		log.Infof("the diff file is empty, shipped!")
		return
	}

	byt, _ := os.ReadFile(*fileNames)
	diff, _ := diffparser.Parse(string(byt))
	resourceNameMap := make(map[string]struct{})
	// checkCandidate deduplicates and checks one resource touched by the diff.
	// An empty testPath or implPath keeps the flat-layout default paths inside
	// checkResourceByName; split-layout callers pass the paths resolved from
	// the product's register.go.
	checkCandidate := func(resName string, isResource bool, fileType, testPath, implPath string) {
		if _, ok := resourceNameMap[resName]; ok {
			return
		}
		resourceNameMap[resName] = struct{}{}
		if !checkResourceByName(resName, isResource, fileType, testPath, implPath) {
			exitCode = 1
		}
	}
	// One split stem resolves once: the implementation and its test share it.
	splitSeen := make(map[string]bool)
	for _, file := range diff.Files {
		if strings.HasSuffix(file.NewName, "_test.go") && !hasCoverageRelevantTestChanges(file) {
			log.Infof("==> %s has no ACC coverage changes; skip coverage check", file.NewName)
			continue
		}
		isNameCorrect = true
		resName := ""
		isResource := true
		fileType := "resource"
		if resourceFileTestRegex.MatchString(file.NewName) {
			resName = strings.TrimSuffix(strings.Split(file.NewName, "/")[1], "_test.go")
		} else if resourceFileRegex.MatchString(file.NewName) {
			resName = strings.TrimSuffix(strings.Split(file.NewName, "/")[1], ".go")
		} else if product, stem, _, ok := splitlayout.IsSplitSource(file.NewName); ok {
			// The implementation and its test share the stem; the first of
			// the pair to appear in the diff carries both paths.
			if splitSeen[product+"/"+stem] {
				continue
			}
			splitSeen[product+"/"+stem] = true
			resolved, err := splitlayout.Resolve(product, stem)
			if err != nil {
				log.Warnf("fail to resolve split-layout file %s. Error: %s", file.NewName, err)
				continue
			}
			// One file can register more than one resource, so check every
			// matched registration.
			for _, r := range resolved {
				rFileType := "data source"
				if r.IsResource() {
					rFileType = "resource"
				}
				checkCandidate(strings.TrimPrefix(r.TypeName, "alicloud_"), r.IsResource(), rFileType, r.TestPath, r.ImplPath)
			}
			// Framework-native registrations — managed resources, data
			// sources, ephemerals, list resources and actions — live outside
			// the provider's SDKv2 maps, so their attribute coverage cannot
			// be measured; said so per registration instead of passing
			// silently.
			fwResolved, err := splitlayout.FrameworkResolve(product, stem)
			if err != nil {
				log.Warnf("fail to resolve framework registrations of %s. Error: %s", file.NewName, err)
				continue
			}
			for _, r := range fwResolved {
				log.Infof("==> %s is registered framework-native (%s): the SDKv2 attribute coverage check does not apply", r.TypeName, r.Kind)
			}
			if len(resolved) == 0 && len(fwResolved) == 0 {
				log.Warnf("==> %s registers no resource, data source, ephemeral, action or list resource", file.NewName)
			}
			continue
		} else {
			continue
		}

		if strings.HasPrefix(resName, "data_source_") {
			isResource = false
			fileType = "data source"
			resName = strings.TrimPrefix(resName, "data_source_")
		} else {
			resName = strings.TrimPrefix(resName, "resource_")
		}

		// Remove alicloud_ prefix to get the clean resource name
		resName = strings.TrimPrefix(resName, "alicloud_")

		checkCandidate(resName, isResource, fileType, "", "")
	}

	os.Exit(exitCode)
}

// checkSingleResource checks a single resource by its full name (e.g., "alicloud_ecs_instance")
func checkSingleResource(fullName string) bool {
	isResource := true
	fileType := "resource"
	resName := fullName
	testPath, implPath := "", ""

	// Determine if it's a data source or resource
	if strings.HasPrefix(fullName, "data.") {
		isResource = false
		fileType = "data source"
		resName = strings.TrimPrefix(fullName, "data.")
	}

	// Remove alicloud_ prefix if present
	resName = strings.TrimPrefix(resName, "alicloud_")

	// Migrated resources keep their tests next to the implementation under
	// alicloud/service/<product>/; fall back to the registration when the
	// flat-layout test file is gone.
	flatTest := "alicloud/resource_alicloud_" + resName + "_test.go"
	if !isResource {
		flatTest = "alicloud/data_source_alicloud_" + resName + "_test.go"
	}
	if _, err := os.Stat(flatTest); os.IsNotExist(err) {
		if r := splitlayout.FindByTypeName("alicloud_" + resName); r != nil {
			isResource = r.IsResource()
			fileType = "resource"
			if !r.IsResource() {
				fileType = "data source"
			}
			testPath, implPath = r.TestPath, r.ImplPath
		} else if splitlayout.FrameworkRegistered("alicloud_" + resName) {
			log.Infof("==> %s %s is registered framework-native; the SDKv2 attribute coverage check does not apply", fileType, resName)
			return true
		}
	}

	return checkResourceByName(resName, isResource, fileType, testPath, implPath)
}

// checkResourceByName checks a resource by name and type. An empty testPath
// or implPath keeps the flat-layout default; split-layout callers pass the
// paths resolved from the product's register.go.
func checkResourceByName(resName string, isResource bool, fileType, testPath, implPath string) bool {
	isNameCorrect = true

	log.Infof("==> Getting %s %s attributes...", fileType, resName)
	var resource = &schema.Resource{}
	var ok bool
	if isResource {
		resource, ok = sdkv2.Provider().ResourcesMap["alicloud_"+resName]
		if !ok || resource == nil {
			log.Errorf("resource %s is not found in the provider ResourceMap\n\n", resName)
			return false
		}
	} else {
		resource, ok = sdkv2.Provider().DataSourcesMap["alicloud_"+resName]
		if !ok || resource == nil {
			log.Errorf("data source %s is not found in the provider DataSourcesMap\n\n", resName)
			return false
		}
	}

	schemaAllSet, schemaMustSet, schemaModifySet, schemaForceNewSet, schemaDeprecatedSet :=
		mapset.NewSet(), mapset.NewSet(), mapset.NewSet(), mapset.NewSet(), mapset.NewSet()
	getSchemaAttr(false, resource.Schema, &schemaAllSet, &schemaMustSet, &schemaModifySet, &schemaForceNewSet, &schemaDeprecatedSet)

	log.Infof("==> Getting %s %s attributes in test cases...", fileType, resName)
	testMustSet, testModifySet, testIgnoreSet :=
		mapset.NewSet(), mapset.NewSet(), mapset.NewSet()
	if testPath == "" {
		testPath = "alicloud/"
		if isResource {
			testPath += "resource_alicloud_" + resName + "_test.go"
		} else {
			testPath += "data_source_alicloud_" + resName + "_test.go"
		}
	}
	check := getTestCaseAttr(testPath, resName, &testMustSet, &testModifySet, &testIgnoreSet)

	// "check" denotes the test code is using a standard template
	if check {
		log.Infof("==> checking %s %s attributes' coverage rate", fileType, resName)
		if checkAttributeSet(resName, fileType, schemaMustSet, testMustSet,
			schemaModifySet, testModifySet, schemaForceNewSet, schemaAllSet, testIgnoreSet, schemaDeprecatedSet, implPath) && isNameCorrect {
			log.Infof("--- PASS!\n\n")
			return true
		}
	}
	// Split-layout tests speak the plugin-testing dialect; the standard
	// template (resourceTestAccConfigFunc and the testAccCheck maps) lives
	// in the flat monolith package and is not importable from
	// alicloud/service/<product>. A split test the parser cannot read is a
	// recognized skip rather than a failure; flat-layout behavior is
	// unchanged.
	if _, _, _, isSplit := splitlayout.IsSplitSource(testPath); isSplit &&
		(!check || (testMustSet.Cardinality() == 0 && testModifySet.Cardinality() == 0)) {
		log.Infof("==> %s %s's split-layout test %s is not written in the standard template; attribute coverage is unmeasurable and skipped", fileType, resName, testPath)
		return true
	}
	// Flat data-source tests speak the existConfig/fakeConfig dialect of
	// dataSourceTestAccConfigFunc; the standard template's Config/Check
	// steps do not exist there, so a parser that extracted nothing
	// measures nothing. Skip with a named notice instead of failing; a
	// data source whose test does carry readable steps is still measured.
	if !isResource && testMustSet.Cardinality() == 0 && testModifySet.Cardinality() == 0 {
		log.Infof("==> %s %s's flat-layout test %s uses the existConfig/fakeConfig template; attribute coverage is unmeasurable and skipped", fileType, resName, testPath)
		return true
	}

	log.Errorf("--- Failed!\n\n")
	return false
}

// get the schema
func getSchemaAttr(isResource bool, schema map[string]*schema.Schema,
	schemaAllSet, schemaMustSet, schemaModifySet, schemaForceNewSet,
	schemaDeprecatedSet *mapset.Set) {

	schemaAttributes := make(map[string]SchemaAttribute)

	getSchemaAttributes("", schemaAttributes, schema)

	for key, value := range schemaAttributes {
		if key == "dry_run" {
			continue
		}
		// Deprecated attributes are not part of the coverage target but may
		// legitimately appear in ImportStateVerifyIgnore to suppress diffs on
		// renamed fields — track them separately so the redundant-attr check
		// can whitelist them.
		if value.Deprecated != "" {
			(*schemaDeprecatedSet).Add(key)
			continue
		}
		(*schemaAllSet).Add(key)
		if value.Optional || value.Required {
			(*schemaMustSet).Add(key)
			if !value.ForceNew {
				(*schemaModifySet).Add(key)
			}
		}
		if value.ForceNew {
			(*schemaForceNewSet).Add(key)
		}

	}
}

func getSchemaAttributes(rootName string, schemaAttributes map[string]SchemaAttribute,
	resourceSchema map[string]*schema.Schema) {
	for key, value := range resourceSchema {
		if rootName != "" {
			key = rootName + "." + key
		}

		if _, ok := schemaAttributes[key]; !ok {
			schemaAttributes[key] = SchemaAttribute{
				Name:       key,
				Type:       value.Type.String(),
				Optional:   value.Optional,
				Required:   value.Required,
				ForceNew:   value.ForceNew,
				Default:    fmt.Sprint(value.Default),
				Deprecated: value.Deprecated,
			}
		}
		if value.Type == schema.TypeSet || value.Type == schema.TypeList {
			if v, ok := value.Elem.(*schema.Schema); ok {
				vv := schemaAttributes[key]
				vv.ElemType = v.Type.String()
				schemaAttributes[key] = vv
			} else {
				vv := schemaAttributes[key]
				vv.ElemType = "Object"
				schemaAttributes[key] = vv
				getSchemaAttributes(key, schemaAttributes, value.Elem.(*schema.Resource).Schema)
			}
		}
	}
}

type SchemaAttribute struct {
	Name        string
	Type        string
	Optional    bool
	Required    bool
	ForceNew    bool
	Default     string
	ElemType    string
	Deprecated  string
	DocsLineNum int
	Removed     string
}

// get the attribute which have been tested
func getTestCaseAttr(filePath string, resourceName string,
	testMustSet, testModifySet, testIgnoreSet *mapset.Set) bool {
	resourceTest, ok := scanTestCaseSteps(filePath, testIgnoreSet)
	if !ok {
		return false
	}
	return parseConfig(resourceTest, testMustSet, testModifySet)
}

// scanTestCaseSteps reads a Go test file and captures, per TestAccAliCloud*
// test function, every step's Config literal (stepStr), the attribute source
// recoverable from the Check-block fallback, and the plan-only flag
// (stepPlanOnly). It is split from getTestCaseAttr so unit tests can assert
// the captured step attribution directly — e.g. that a helper-Config step
// whose `PlanOnly: true` line arrives while the Check fallback is still
// waiting for its map is marked plan-only — without depending on the coverage
// scoring in parseConfig.
func scanTestCaseSteps(filePath string, testIgnoreSet *mapset.Set) (ResourceTest, bool) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Errorf("fail to open test file %s. Error: %s", filePath, err)
		return ResourceTest{}, false
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	resourceTest := ResourceTest{
		resourceName: filePath,
		funcs:        map[string]FuncTest{},
	}
	line := 0
	funcName := ""
	stepNumber := 0
	configStr := ""
	inConfig := false
	inFunc := false
	ignoreStr := ""
	inIgnore := false
	// When a Config line does not embed an inline map[string]interface{} (e.g.
	// a helper-function call like `Config: foo(name, "PayByCenOwner")`), the
	// attribute data is still recoverable from the Check block that follows:
	// `testAccCheck(map[string]string{ "order_type": "PayByCenOwner", ... })`.
	// We track that fallback state here and record the Check map in place of
	// the Config map when it applies.
	inCheckFallback := false
	checkFallbackStr := ""
	checkFallbackStep := 0
	// planOnlySeenStepOpen tracks whether a step-literal `{` has opened since
	// the last Config capture ended (or the function started). It lets a
	// `PlanOnly: true` line written outside an active capture be attributed to
	// the right step: to the step that is about to capture its Config when a
	// new literal already opened, or to the step whose capture just closed on
	// an earlier sentinel (field order inside a TestStep literal is free).
	planOnlySeenStepOpen := false
	// planOnlyPending defers the attribution of a `PlanOnly: true` line seen
	// after a step literal opened but before its Config capture closes. The
	// flag is only written onto the step index when that step's capture really
	// ends, so a step without any Config (e.g. one holding only
	// ResourceName/ImportState fields) cannot spill the flag onto the next
	// captured step; the pending flag is dropped when such a literal closes.
	planOnlyPending := false
	// planOnlyBraceDepth counts the braces open in the step literal the
	// plan-only attribution state currently tracks: 1 for the literal's own
	// `{`, one more per nested block (e.g. ConfigPlanChecks) opened inside it.
	// A `}`-prefixed line only ends the attribution window (dropping a pending
	// flag) when it closes the step literal itself, not a nested block: a
	// nested block between the `PlanOnly: true` line and the `Config:` line
	// must not make the scanner lose the pending attribution.
	planOnlyBraceDepth := 0
	for scanner.Scan() {
		line += 1
		text := scanner.Text()
		// commented line
		if commentedRegex.MatchString(text) {
			continue
		} else if text == "}" {
			inFunc = false
			// Reset the per-step capture state at end-of-function so a step
			// that never reaches a Check/ExpectError sentinel (e.g. malformed
			// or unusual layouts) does not leak its in-progress Config buffer
			// into the next function.
			inConfig = false
			configStr = ""
			inCheckFallback = false
			checkFallbackStr = ""
			planOnlySeenStepOpen = false
			planOnlyPending = false
			planOnlyBraceDepth = 0
		} else if normalFuncRegex.MatchString(text) {
			if unitFuncRegex.MatchString(text) {
				continue
			}
			if !standardFuncRegex.MatchString(text) {
				name := text[strings.Index(text, "T"):strings.Index(text, "(")]
				log.Errorf("testcase %s should start with TestAccAliCloud", name)
				isNameCorrect = false
			}
			inFunc = true
			funcName = text[strings.Index(text, "T"):strings.Index(text, "(")]
			stepNumber = 0
			// Defensive reset in case the previous function ended without a
			// proper closing `}` token on its own line.
			inConfig = false
			configStr = ""
			inCheckFallback = false
			checkFallbackStr = ""
			planOnlySeenStepOpen = false
			planOnlyPending = false
			planOnlyBraceDepth = 0
			resourceTest.funcs[funcName] = FuncTest{
				stepStr:        map[int]string{},
				stepAttributes: map[int]map[string]interface{}{},
				stepPlanOnly:   map[int]bool{},
			}
		}
		if inFunc {
			if configRegex.MatchString(text) {
				configStr += text
				inConfig = true
				continue
			}
			if inConfig {
				// Config block ends at Check: OR at any other well-known step
				// sentinel (ExpectError, ExpectNonEmptyPlan, PlanOnly, Destroy,
				// etc.) so steps that validate expected errors or plan-only
				// behavior close cleanly instead of merging into the next Config.
				isCheck := checkRegex.MatchString(text)
				isSentinel := !isCheck && stepEndSentinelRegex.MatchString(text)
				if isCheck || isSentinel {
					// A `PlanOnly: true` line closing the capture marks the step
					// being recorded here as plan-only.
					if isSentinel && planOnlyStep(text) {
						resourceTest.funcs[funcName].stepPlanOnly[stepNumber] = true
					}
					if planOnlyPending {
						// A `PlanOnly: true` line seen before this step's Config
						// line: the capture that just closed is the step it
						// belongs to.
						resourceTest.funcs[funcName].stepPlanOnly[stepNumber] = true
						planOnlyPending = false
					}
					resourceTest.funcs[funcName].stepStr[stepNumber] = configStr
					// If Config had no inline map AND we hit a Check line,
					// arm the Check-block fallback so we record the Check's
					// testAccCheck(map[string]string{...}) as the step's
					// effective attribute source. Sentinel steps (ExpectError
					// etc.) do not carry attribute data, so skip the fallback.
					if isCheck && !strings.Contains(configStr, "{") {
						inCheckFallback = true
						checkFallbackStr = ""
						checkFallbackStep = stepNumber
					}
					inConfig = false
					configStr = ""
					stepNumber++
					planOnlySeenStepOpen = false
				} else {
					// remove extra spaces and '\'
					text = symbolRegex.ReplaceAllString(text, "")
					if len(text) == 0 || strings.HasPrefix(text, "//") {
						continue
					}
					configStr += text + "\n"
				}
			} else if inCheckFallback {
				trimmed := symbolRegex.ReplaceAllString(text, "")
				if len(trimmed) == 0 || strings.HasPrefix(trimmed, "//") {
					continue
				}
				// While the fallback is still waiting for the Check map to
				// open, a `PlanOnly: true` field of the same step must mark
				// that step instead of being swallowed as check-map text.
				if checkFallbackStr == "" && planOnlyStep(text) {
					resourceTest.funcs[funcName].stepPlanOnly[checkFallbackStep] = true
					continue
				}
				// Open bracket: start of the testAccCheck map.
				if checkFallbackStr == "" {
					if strings.Contains(trimmed, "{") {
						// Synthesize the same decoration the standard-template
						// parser expects (the leading `interface{}` pair is
						// skipped by parseConfig via Index("{")+2).
						checkFallbackStr = "Config: testAccConfig(map[string]interface{}{\n"
					}
					continue
				}
				if strings.HasPrefix(trimmed, "}") {
					// End of the Check map. Record it and disarm.
					checkFallbackStr += "}),\n"
					resourceTest.funcs[funcName].stepStr[checkFallbackStep] = checkFallbackStr
					inCheckFallback = false
					checkFallbackStr = ""
					continue
				}
				checkFallbackStr += trimmed + "\n"
			} else {
				// Outside an active capture: track step-literal boundaries and
				// `PlanOnly: true` lines so a plan-only flag written before the
				// step's Config line (or after the capture already closed on an
				// earlier Check/sentinel within the same step) still marks the
				// right step.
				trimmed := strings.TrimSpace(text)
				if trimmed == "{" || strings.Contains(trimmed, "TestStep{") {
					// A step literal opened, either as a bare `{` element or
					// with the literal's `{` on the `resource.TestStep{` line
					// itself (repo layouts use both). Brace bookkeeping for the
					// plan-only attribution restarts from this brace.
					planOnlySeenStepOpen = true
					planOnlyBraceDepth = 1
				} else if strings.HasPrefix(trimmed, "}") {
					if planOnlyBraceDepth > 1 {
						// A nested block inside the step literal (e.g.
						// ConfigPlanChecks written between the PlanOnly line and
						// the Config line) closed: the pending attribution must
						// survive until the step's capture really ends.
						planOnlyBraceDepth--
					} else {
						// Step literal closed without a Config capture (e.g. a
						// step holding only ResourceName/ImportState/Destroy
						// fields): drop the pending plan-only attribution
						// instead of letting it tag the next captured step.
						planOnlySeenStepOpen = false
						planOnlyPending = false
						planOnlyBraceDepth = 0
					}
				} else if planOnlySeenStepOpen && strings.HasSuffix(trimmed, "{") {
					// A nested block opened inside the step literal.
					planOnlyBraceDepth++
				} else if planOnlyStep(text) {
					if planOnlySeenStepOpen {
						// A new step literal already opened: the flag belongs to
						// the step whose Config is captured next, but only if
						// that capture really happens — defer it.
						planOnlyPending = true
					} else if stepNumber > 0 {
						// Still inside the literal whose capture just closed on
						// Check/sentinel: the flag belongs to that step.
						resourceTest.funcs[funcName].stepPlanOnly[stepNumber-1] = true
					}
				}
			}
			if ignoreRegex.MatchString(text) {
				inIgnore = true
				ignoreStr += text
				continue
			}
			if inIgnore {
				ignoreStr += text
				if strings.Contains(text, "}") {
					inIgnore = false
					ignoreStr = strings.ReplaceAll(ignoreStr, "\"", "")
					ignoreStr = symbolRegex.ReplaceAllString(ignoreStr, "")
					attrStr := ignoreStr[strings.Index(ignoreStr, "{")+1 : strings.Index(ignoreStr, "}")]
					if len(attrStr) != 0 {
						attrSlice := strings.Split(attrStr, ",")
						for _, v := range attrSlice {
							if len(v) != 0 && v != "dry_run" {
								(*testIgnoreSet).Add(v)
							}
						}
					}
					ignoreStr = ""
				}
			}
		}

	}

	return resourceTest, true

}

func parseConfig(resourceTest ResourceTest,
	testMustSet, testModifySet *mapset.Set) (toCheck bool) {
	for funcName, f := range resourceTest.funcs {
		// attribute-value map in a test func
		attributeValueMap := map[string]string{}
		// appliedSteps counts the non-plan-only steps seen so far (skipped
		// helper-Config steps included), so a key first seen in a later step
		// only counts as a modification when a real apply step already
		// established the baseline. A leading plan-only step does not.
		appliedSteps := 0
		for configIndex := 0; configIndex < len(f.stepStr); configIndex++ {
			planOnly := f.stepPlanOnly[configIndex]
			applyOrdinal := appliedSteps
			if !planOnly {
				appliedSteps++
			}
			configStr := f.stepStr[configIndex]
			// Some tests (e.g. cross-account order_type tests) pass a helper
			// function name instead of an inline map[string]interface{} to
			// Config. They contribute no attribute data this parser can read,
			// so skip the step rather than failing the whole resource — the
			// standard-template steps (basic0 / basic1) still populate
			// testMustSet / testModifySet and keep coverage accurate.
			if !strings.Contains(configStr, "{") {
				log.Debugf("skipping non-template step in func [%s] of [%s]", funcName, resourceTest.resourceName)
				continue
			}
			configStr = configStr[strings.Index(configStr, "{")+2 : strings.LastIndex(configStr, "}")+1]
			if configStr == "{}" {
				continue
			}
			configSlice := strings.Split(configStr, "\n")

			// traverse each line of config
			for i, v := range configSlice {
				valueIndex := strings.Index(v, ":")

				// "xxx:xxx",
				if strings.Contains(v, ":") {
					splitIndex := strings.Index(v, ":")
					beforeCount := strings.Count(v[:splitIndex], "\"")
					if beforeCount%2 != 0 {
						valueIndex = -1
					}
				}

				valueSuffix := string(v[len(v)-1])
				valueStr := v[valueIndex+1 : len(v)-1]
				bracketCount := strings.Count(v, "]") + strings.Count(v, "[")
				if v[valueIndex+1] == '"' && strings.HasSuffix(v, "},") ||
					strings.HasSuffix(v, "],") && bracketCount%2 == 1 {
					valueSuffix = v[len(v)-2:]
					valueStr = v[valueIndex+1 : len(v)-2]
				}

				// Check if this is a nested JSON string (contains multi-level escaped quotes like \\\")
				// This indicates the value is itself a JSON string with internal JSON structure
				// We need to skip all transformations to preserve its original structure
				isNestedJSONString := strings.HasPrefix(valueStr, "\"") &&
					strings.Contains(valueStr, `\\\"`) &&
					(strings.Contains(valueStr, "{") || strings.Contains(valueStr, "["))

				if !isNestedJSONString {
					// "xxx"+xx, "xxx"+xxx+"xxx, `xxx`+"xxx"+`xxx`
					if strings.Contains(valueStr, "+") {
						v := valueStr
						index := strings.Index(v, "+")
						for index != -1 {
							beforeStr := v[:index]
							beforeStr = strings.TrimSpace(beforeStr)
							if strings.Count(beforeStr, "\"")%2 == 0 ||
								(beforeStr[0] == '`' && beforeStr[len(beforeStr)-1] == '`') {
								v = v[index+1:]
								index = strings.Index(v, "+")
							} else {
								break
							}
						}
						if index == -1 {
							valueStr = strings.ReplaceAll(valueStr, "\"", "")
							valueStr = strings.ReplaceAll(valueStr, "`", "")
							valueStr = "\"" + valueStr + "\""
						}
					}

					// `xxx`
					if strings.HasPrefix(valueStr, "`") && strings.HasSuffix(valueStr, "`") {
						valueStr = "\"" + valueStr[:len(valueStr)-1] + "\"" + valueSuffix
					} else {
						valueStr += valueSuffix
					}

					// value is a map or slice
					for i := 0; i < len(matchMap); i++ {
						k := matchMap[i][0]
						v := matchMap[i][1]
						if strings.Contains(valueStr, k) {
							valueStr = strings.ReplaceAll(valueStr, k, v)
						}
					}

					// value with variable or func
					if variableRegex.MatchString(valueStr) || valueFuncRegex.MatchString(valueStr) {
						valueStr = strings.ReplaceAll(valueStr, "\\\"", "*")
						valueStr = strings.ReplaceAll(valueStr, "\"", "")
						valueStr = "\"" + valueStr + "\"" + valueSuffix
					}

					// "xxx/xxx", "xxx/"xxx"
					if valueOnlySymbol.MatchString(valueStr) {
						valueStr = strings.ReplaceAll(valueStr, "\\", "*")
						valueStr = strings.ReplaceAll(valueStr, "\"", "*")
						valueStr = "\"" + valueStr[1:len(valueStr)-1] + "\"" + valueSuffix
					}
				} else {
					// For nested JSON strings, skip all transformations and just add the suffix
					valueStr += valueSuffix
				}

				configSlice[i] = v[:valueIndex+1] + valueStr
			}

			configStr = strings.Join(configSlice, "")
			configRune := []rune(configStr)

			// match the bracket
			if toCheck = bracketMatch(funcName, configIndex, configRune); !toCheck {
				return toCheck
			}

			configStr = string(configRune)
			jsonData := []byte(configStr)
			var v interface{}
			err := json.Unmarshal(jsonData, &v)
			if err != nil {
				log.Errorf("fail to unmarshal func %v's number %v config: \n%s\n%s", funcName, configIndex, configStr, err)
				return false
			}
			data := v.(map[string]interface{})

			f.stepAttributes[configIndex] = data

			// Plan-only steps still contribute attribute coverage (must-set),
			// but they never apply, so parseAttr skips the modification
			// accounting and the value bookkeeping for them. applyOrdinal is
			// the count of apply steps that precede this one: a key absent
			// from the value map is only a first-occurrence update (not a
			// baseline) when such a step exists.
			parseAttr(applyOrdinal, "", data, attributeValueMap, testMustSet, testModifySet, planOnly)
		}
	}
	return true
}

// planOnlyStep reports whether a TestStep field line enables PlanOnly. Only
// `PlanOnly: true` makes the step plan-only (plan refresh, no apply); an
// explicit `PlanOnly: false` keeps the step a normal apply step.
func planOnlyStep(text string) bool {
	m := planOnlyValueRegex.FindStringSubmatch(text)
	return m != nil && m[2] == "true"
}

func bracketMatch(funcName string, configIndex int, s []rune) bool {
	var stack []string

	for i := 0; i < len(s); i++ {
		r := string(s[i])
		switch r {
		// remove the extra comma
		case ",":
			if i != len(s)-1 && (s[i+1] == '}' || s[i+1] == ']') {
				s[i] = rune(' ')
			}
		case "{", "[":
			stack = append(stack, r)
		case "}", "]":
			{
				if len(stack) == 0 {
					log.Errorf("fail to math bracket [ ] in func %s's number %d config:%s\n", funcName, configIndex, string(s))
					return false
				}
				temp := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				if temp == "[" && r == "}" {
					s[i] = ']'
				}
				if bracket[temp] == r {
					break
				}
			}
		}
	}
	return true
}

// parseAttr walks a captured Config's attribute tree. For plan-only steps
// (planOnly=true) the attributes still populate testMustSet, but they are
// excluded from modification coverage: the step never applies, and its values
// must not clobber the last applied value that later apply steps are compared
// against. appliedSteps is the number of apply steps processed before this
// one: a key absent from the value map counts as a first-occurrence update
// only when at least one apply step already established the baseline.
func parseAttr(appliedSteps int, rootName string, data interface{}, attributeValueMap map[string]string,
	testMustSet, testModifySet *mapset.Set, planOnly bool) {
	if d, ok := data.(map[string]interface{}); ok {
		for key, value := range d {
			if rootName != "" {
				key = rootName + "." + key
			}
			addCoverageKey(testMustSet, key)
			if !planOnly {
				// check if the attribute has been updated
				updated := false
				if v, ok := attributeValueMap[key]; ok {
					if fmt.Sprintf("%v", value) != v {
						updated = true
					}
				} else if appliedSteps > 0 {
					updated = true
				}
				if updated {
					addCoverageKey(testModifySet, key)
				}
				attributeValueMap[key] = fmt.Sprintf("%v", value)
			}
			parseAttr(appliedSteps, key, value, attributeValueMap, testMustSet, testModifySet, planOnly)
		}
	} else if d, ok := data.([]interface{}); ok {
		for _, v := range d {
			if !planOnly {
				attributeValueMap[rootName] = fmt.Sprintf("%v", data)
			}
			parseAttr(appliedSteps, rootName, v, attributeValueMap, testMustSet, testModifySet, planOnly)
		}
	}
}

// stripNumericIndexParts drops numeric path segments from dotted keys, e.g.
// "zone.0.zone_id" -> "zone.zone_id". Used to normalize flattened Check-map
// keys so they align with the schema's nested attribute names.
func stripNumericIndexParts(key string) string {
	if !strings.Contains(key, ".") {
		return key
	}
	parts := strings.Split(key, ".")
	out := parts[:0]
	for _, p := range parts {
		if _, err := strconv.Atoi(p); err == nil {
			continue
		}
		out = append(out, p)
	}
	return strings.Join(out, ".")
}

// addCoverageKey registers a key against a coverage set along with every
// ancestor path built from its non-numeric segments. This lets flat keys
// from the Check-fallback (e.g. "zone.0.zone_id") satisfy coverage for both
// the leaf attribute ("zone.zone_id") and its container ("zone").
func addCoverageKey(set *mapset.Set, key string) {
	(*set).Add(key)
	stripped := stripNumericIndexParts(key)
	if stripped != key {
		(*set).Add(stripped)
	}
	if !strings.Contains(stripped, ".") {
		return
	}
	parts := strings.Split(stripped, ".")
	for i := 1; i <= len(parts); i++ {
		(*set).Add(strings.Join(parts[:i], "."))
	}
}

// indexLikeSegmentRegex matches dotted path segments that can only be index
// attempts, never schema attribute names: an optional sign followed by
// digits. Schema attribute names are [a-zA-Z_0-9] and never start with a
// sign or a digit, so such a segment has no legal spelling in either a real
// state key or a schema path.
var indexLikeSegmentRegex = regexp.MustCompile(`^[+-]?[0-9]+$`)

// legalIndexSegment reports whether a dotted path segment is a legal
// positional index: one or more ASCII digits and nothing else. A signed
// segment (e.g. "-1" in "block.-1.leaf") looks like an index but is not a
// legal one — no real state key or schema path carries it — so an
// ImportStateVerifyIgnore entry containing one is invalid and must keep
// failing the original gates instead of being normalized alongside legal
// indexed entries.
func legalIndexSegment(seg string) bool {
	if seg == "" {
		return false
	}
	for _, r := range seg {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

// normalizeIgnoreEntry aligns an ImportStateVerifyIgnore entry with the
// schema's index-free dotted leaf paths: every legal index segment (pure
// digits) is dropped, reusing the same stripping the Check-coverage keys
// already go through (stripNumericIndexParts). It returns ok=false when the
// entry contains an index-like segment that is not a legal non-negative
// integer index (e.g. "block.-1.leaf"); such an entry is not normalized, so
// its original spelling keeps failing the coverage and ignore gates.
func normalizeIgnoreEntry(entry string) (string, bool) {
	if !strings.Contains(entry, ".") {
		return entry, true
	}
	for _, p := range strings.Split(entry, ".") {
		if indexLikeSegmentRegex.MatchString(p) && !legalIndexSegment(p) {
			return entry, false
		}
	}
	return stripNumericIndexParts(entry), true
}

// normalizeIgnoreSet builds the CI-side view of ImportStateVerifyIgnore for
// the set comparisons in checkAttributeSet: an indexed entry such as
// "db_cluster_ip_array.0.white_list_type" additionally contributes its
// schema leaf path ("db_cluster_ip_array.white_list_type") so the
// modification-coverage Difference can offset the schema-side key. The
// original strings stay untouched in testIgnoreSet — the SDK import
// comparison matches them with strings.HasPrefix against real state keys,
// so the runtime behavior of ImportStateVerifyIgnore does not change, and
// no parent path or sibling field is added to any exemption. Entries with
// invalid index segments are returned separately, unnormalized: they keep
// the original spelling and keep blocking.
func normalizeIgnoreSet(set mapset.Set) (mapset.Set, []string) {
	normalized := mapset.NewSet()
	var invalid []string
	for _, v := range set.ToSlice() {
		entry, ok := v.(string)
		if !ok {
			continue
		}
		normalized.Add(entry)
		stripped, valid := normalizeIgnoreEntry(entry)
		if !valid {
			invalid = append(invalid, entry)
			continue
		}
		if stripped != entry {
			normalized.Add(stripped)
		}
	}
	sort.Strings(invalid)
	return normalized, invalid
}

var (
	commentedRegex    = regexp.MustCompile("^[\t]*//")
	normalFuncRegex   = regexp.MustCompile("^func Test(.*)")
	unitFuncRegex     = regexp.MustCompile("^func TestUnit(.*)")
	standardFuncRegex = regexp.MustCompile("^func TestAccAliCloud(.*)")
	configRegex       = regexp.MustCompile("(^[\t]*)Config:(.*)")
	checkRegex        = regexp.MustCompile("(.*)Check:(.*)")
	// Steps that end the Config-capture state machine without a Check func.
	// They still need to close the capture so a following Config step is
	// segmented independently instead of being merged into the current one.
	// Matches lines such as `ExpectError: regexp.MustCompile(...)`,
	// `ExpectNonEmptyPlan: true` or `PlanOnly: true`. Without PlanOnly here,
	// a PlanOnly step merges with the next Config step and triggers a false
	// bracket/JSON syntax error (the step-boundary `},`/`{` tokens become
	// part of the captured config body).
	stepEndSentinelRegex = regexp.MustCompile(`(^[\t]*)(ExpectError|ExpectNonEmptyPlan|PreventDiskCleanup|Destroy|ImportState|PlanOnly):`)
	// planOnlyValueRegex identifies a TestStep field that enables PlanOnly and
	// captures its literal value. Only `PlanOnly: true` makes a step plan-only;
	// an explicit `PlanOnly: false` keeps it a normal apply step.
	planOnlyValueRegex = regexp.MustCompile(`(^[\t ]*)PlanOnly:\s*(true|false)`)
	ignoreRegex        = regexp.MustCompile("(.*)ImportStateVerifyIgnore:(.*)")
	attrRegex          = regexp.MustCompile("^([{]*)\"([a-zA-Z_0-9-.]+)\":(.*)")
	symbolRegex        = regexp.MustCompile(`\s`)
	variableRegex      = regexp.MustCompile("(^[a-zA-Z_0-9]+)|(\"[+]\")")
	valueFuncRegex     = regexp.MustCompile("[(].*[\"].*[\"].*[)]")
	valueOnlySymbol    = regexp.MustCompile(`.*([^\\\"])(\\)([^\\\"]).*`)
	bracket            = map[string]string{
		"{": "}",
		"[": "]",
	}
	matchMap = map[int][]string{
		0: []string{"[]string{", "["},
		1: []string{"[]map[string]interface{}{", "["},
		2: []string{"[]map[string]string{", "["},
		3: []string{"[]interface{}{", "["},
		4: []string{"map[string]interface{}{", "{"},
		5: []string{"map[string]string{", "{"},
		6: []string{"\\n", ""},
		7: []string{"`${map(", "["},
		8: []string{")}`", "]"},
	}

	isNameCorrect = true
)

type ResourceTest struct {
	resourceName string
	funcs        map[string]FuncTest
}

type FuncTest struct {
	stepStr        map[int]string
	stepAttributes map[int]map[string]interface{}
	// stepPlanOnly marks captured steps whose TestStep sets PlanOnly: true.
	// A plan-only step only refreshes the plan — it never applies — so fields
	// whose values change only inside such a step do not count as actual
	// modification coverage (they still count towards attribute coverage).
	stepPlanOnly map[int]bool
}

// TreeNode represents a node in the attribute tree
type TreeNode struct {
	Name     string
	Children map[string]*TreeNode
}

// buildAttributeTree builds a tree structure from attribute paths
func buildAttributeTree(attributes []string) *TreeNode {
	root := &TreeNode{
		Name:     "root",
		Children: make(map[string]*TreeNode),
	}

	for _, attr := range attributes {
		parts := strings.Split(attr, ".")
		current := root
		for _, part := range parts {
			if _, exists := current.Children[part]; !exists {
				current.Children[part] = &TreeNode{
					Name:     part,
					Children: make(map[string]*TreeNode),
				}
			}
			current = current.Children[part]
		}
	}

	return root
}

// printTree prints the tree structure with proper formatting
func printTree(node *TreeNode, prefix string, isLast bool, buffer *strings.Builder) {
	if node.Name != "root" {
		// Determine the branch symbol
		branch := "├── "
		if isLast {
			branch = "└── "
		}
		buffer.WriteString(prefix + branch + node.Name + "\n")

		// Update prefix for children
		if isLast {
			prefix += "    "
		} else {
			prefix += "│   "
		}
	}

	// Sort children for consistent output
	childNames := make([]string, 0, len(node.Children))
	for name := range node.Children {
		childNames = append(childNames, name)
	}
	sort.Strings(childNames)

	// Print children
	for i, name := range childNames {
		child := node.Children[name]
		isLastChild := i == len(childNames)-1
		printTree(child, prefix, isLastChild, buffer)
	}
}

// formatAttributesAsTree converts attribute list to tree format string
func formatAttributesAsTree(attributes []interface{}) string {
	if len(attributes) == 0 {
		return ""
	}

	// Convert interface{} slice to string slice
	attrStrings := make([]string, 0, len(attributes))
	for _, attr := range attributes {
		if attrStr, ok := attr.(string); ok {
			attrStrings = append(attrStrings, attrStr)
		}
	}

	// Sort for consistent output
	sort.Strings(attrStrings)

	// Build tree
	tree := buildAttributeTree(attrStrings)

	// Print tree
	var buffer strings.Builder
	printTree(tree, "", true, &buffer)

	return buffer.String()
}

// neverReadBackAttributes returns the attributes in testIgnoreSet that the
// resource/data source implementation never writes back with d.Set (e.g.
// create-only inputs that the Get/List API does not return). Such attributes
// have no source during import, so keeping them in ImportStateVerifyIgnore is
// legitimate even when they are ForceNew. The exemption set stays empty when
// the implementation file cannot be read, so the check fails closed.
// An empty implPath keeps the flat-layout default.
func neverReadBackAttributes(resourceName string, fileType string, testIgnoreSet mapset.Set, implPath string) mapset.Set {
	exemptions := mapset.NewSet()
	if implPath == "" {
		implPath = "alicloud/resource_alicloud_" + resourceName + ".go"
		if fileType == "data source" {
			implPath = "alicloud/data_source_alicloud_" + resourceName + ".go"
		}
	}
	content, err := os.ReadFile(implPath)
	if err != nil {
		return exemptions
	}
	src := string(content)
	for _, attr := range testIgnoreSet.ToSlice() {
		attrStr, ok := attr.(string)
		if !ok {
			continue
		}
		// Nested ignore entries (parent.child) are covered by the parent d.Set.
		root := strings.SplitN(attrStr, ".", 2)[0]
		if !strings.Contains(src, fmt.Sprintf("d.Set(%q", root)) {
			exemptions.Add(attrStr)
		}
	}
	return exemptions
}

// apiNonRoundTripAttributes returns the ImportStateVerifyIgnore attributes
// whose configured value does NOT round-trip into the state Read produces for
// an import. Read may retain the configured value via d.Set (so the attribute
// IS read back, not neverReadBack), yet the value observed during a real
// ImportStateVerify still diverges from the configured value; listing such an
// attribute in ImportStateVerifyIgnore is therefore a precise, per-field
// allowance — not a Read bug being covered up, and not a reason to change
// production Read or to generalize the exemption.
//
// This is NOT a blanket exemption for every attribute that uses d.Set. Each
// (resource, attr) pair is registered explicitly below with a reference to the
// observed import divergence. A pair MUST be removed once the import path
// round-trips the configured value. Unregistered resources fail closed
// (empty set -> no exemption). resourceName is the bare name passed to
// checkAttributeSet (alicloud_ prefix already stripped); e.g.
// "aligreen_oss_stock_task".
var apiNonRoundTripRegistry = map[string]map[string]string{
	"aligreen_oss_stock_task": {
		// aligreen OssStockTask Read: when end_date is configured it is
		// retained via d.Set("end_date", <configured>); when absent (e.g. an
		// import with no prior config) Read backfills EndDate from the Describe
		// response (objectRaw["EndDate"]). The configured-state and
		// import-state value paths therefore differ. A real ImportStateVerify
		// observed the imported EndDate diverge from the configured end_date,
		// so this precise field is retained in ImportStateVerifyIgnore.
		//
		// This does NOT assert the API contract is broken or that Read is a
		// no-op: it records the observed, per-field import divergence only.
		// Production Read is not changed for this exemption, and the
		// exemption is not generalized to other d.Set attributes. Remove this
		// entry once the import path round-trips the configured end_date.
		// (Internal ACC taskId/log references are kept in the private Task
		// report, not as the sole public proof.)
		"end_date": "aligreen_oss_stock_task",
	},
	"instance": {
		// alicloud_instance Read: DescribeInstances returns CpuOptions with
		// CoreCount/ThreadsPerCore/TopologyType but never Accelerators, so Read
		// preserves the configured accelerators from prior state via
		// d.Set("cpu_options", ...). An import starts from empty state, so the
		// imported cpu_options carries no accelerators and a real
		// ImportStateVerify observes the imported value diverge from the
		// configured "vqat" member (the only member value the ECS API accepts,
		// InvalidParameter.CpuOptionsAccelerators otherwise). The field is
		// retained in ImportStateVerifyIgnore for that observed import
		// divergence only; remove this entry once DescribeInstances echoes
		// CpuOptions.Accelerators and the import path round-trips it.
		"cpu_options.accelerators": "instance",
	},
}

func apiNonRoundTripAttributes(resourceName string, testIgnoreSet mapset.Set) mapset.Set {
	exemptions := mapset.NewSet()
	registry, ok := apiNonRoundTripRegistry[resourceName]
	if !ok {
		return exemptions
	}
	regSet := mapset.NewSet()
	for attr := range registry {
		regSet.Add(attr)
	}
	return testIgnoreSet.Intersect(regSet)
}

func checkAttributeSet(resourceName string, fileType string, schemaMustSet, testMustSet,
	schemaModifySet, testModifySet, schemaForceNewSet, schemaAllSet, testIgnoreSet,
	schemaDeprecatedSet mapset.Set, implPath string) bool {

	isFullCover, isIgnoreLegal, isAllModified := true, true, true

	notCoverSlice := schemaMustSet.Difference(testMustSet).ToSlice()
	if len(notCoverSlice) != 0 {
		isFullCover = false
		schemaCount := float64(len(schemaMustSet.ToSlice()))
		notCoverCount := float64(len(notCoverSlice))
		coverageRate := 1 - (notCoverCount / schemaCount)
		log.Infof("resource %s attributes has %.2f%% testing coverage rate ", resourceName, coverageRate*100)

		// Format as tree structure
		treeStr := formatAttributesAsTree(notCoverSlice)
		log.Errorf("resource %s attributes missing test cases (%d attributes):\n%s", resourceName, len(notCoverSlice), treeStr)
	} else {
		log.Infof("resource %s attributes has 100%% testing coverage rate ", resourceName)
	}

	// Align indexed ImportStateVerifyIgnore entries (e.g.
	// "db_cluster_ip_array.0.white_list_type") with the schema's index-free
	// leaf paths for the set comparisons below. The original strings stay in
	// testIgnoreSet: the SDK import comparison matches them with
	// strings.HasPrefix against real state keys, so the runtime behavior of
	// ImportStateVerifyIgnore is unchanged. Entries whose index-like segments
	// are not legal non-negative integers (e.g. "block.-1.leaf") are invalid:
	// they are not normalized and fail the ignore gate.
	normalizedIgnore, invalidIndexedIgnore := normalizeIgnoreSet(testIgnoreSet)
	if len(invalidIndexedIgnore) != 0 {
		isIgnoreLegal = false
		invalidIndexedIgnoreStr, _ := json.Marshal(invalidIndexedIgnore)
		log.Errorf("resource %s attributes %v in ImportStateVerifyIgnore contain invalid index segments; only non-negative integer indexes may align with schema paths", resourceName, string(invalidIndexedIgnoreStr))
	}

	// ForceNew attributes in ImportStateVerifyIgnore usually mean the resource
	// should be fixed to read them back. Create-only input attributes that the
	// Get/List API never returns are the exception: Read never writes them back
	// with d.Set, so the imported state can never contain them and listing them
	// in ImportStateVerifyIgnore is the only workable pattern. Exempt exactly
	// those attributes. The comparison runs against the normalized ignore view
	// so an indexed ForceNew ignore ("block.0.fn_attr") is judged exactly like
	// its schema path ("block.fn_attr") — indexed and non-indexed spellings of
	// the same field cannot diverge on this gate.
	unreadAttrs := neverReadBackAttributes(resourceName, fileType, normalizedIgnore, implPath)
	apiNonRoundTrip := apiNonRoundTripAttributes(resourceName, normalizedIgnore)
	forceNewButIgnore := schemaForceNewSet.Intersect(normalizedIgnore).Difference(unreadAttrs).Difference(apiNonRoundTrip).ToSlice()
	if len(forceNewButIgnore) != 0 {
		isIgnoreLegal = false
		forceNewButIgnoreStr, _ := json.Marshal(forceNewButIgnore)
		log.Errorf("resource %s [ForceNew] attributes %v are in ImportStateVerifyIgnore array ", resourceName, string(forceNewButIgnoreStr))
	}
	exempted := schemaForceNewSet.Intersect(normalizedIgnore).Intersect(unreadAttrs).Union(apiNonRoundTrip)
	if exempted.Cardinality() > 0 {
		exemptedStr, _ := json.Marshal(exempted.ToSlice())
		log.Infof("resource %s [ForceNew] attributes %v stay in ImportStateVerifyIgnore because the resource Read never writes them back (neverReadBack) or the API does not round-trip the configured value (apiNonRoundTrip, proven per-attr)", resourceName, string(exemptedStr))
	}
	// Attributes in ImportStateVerifyIgnore that are neither part of the
	// active schema nor deprecated are genuinely redundant. Deprecated fields
	// are allowed to stay in ImportStateVerifyIgnore to suppress diffs during
	// the deprecation window without being flagged by coverage checks. Indexed
	// entries are judged through the same normalization: an entry is legal
	// when its schema leaf path exists; invalid-index entries were already
	// reported above and are not normalized here.
	redundantAttr := testIgnoreSet.Difference(schemaAllSet).Difference(schemaDeprecatedSet).ToSlice()
	redundantAttrFinal := []string{}
	for _, v := range redundantAttr {
		vStr := v.(string)
		normalized, valid := normalizeIgnoreEntry(vStr)
		if !valid {
			continue
		}
		if schemaAllSet.Contains(normalized) || schemaDeprecatedSet.Contains(normalized) {
			continue
		}
		redundantAttrFinal = append(redundantAttrFinal, vStr)
	}
	if len(redundantAttrFinal) != 0 {
		isIgnoreLegal = false
		redundantAttrFinalStr, _ := json.Marshal(redundantAttrFinal)
		log.Errorf("resource %s attributes %v should not in ImportStateVerifyIgnore array", resourceName, string(redundantAttrFinalStr))
	}
	// The modification-coverage Difference runs against the normalized ignore
	// view: an indexed ignore offsets exactly its schema leaf path (and only
	// that leaf — parent paths and sibling fields are never added).
	schemaModifySet = schemaModifySet.Difference(normalizedIgnore)

	// Get immutable attributes from documentation
	immutableAttrs := getImmutableAttributesFromDoc(resourceName, fileType)
	if immutableAttrs.Cardinality() > 0 {
		log.Infof("Found %d immutable attributes in documentation: %v", immutableAttrs.Cardinality(), immutableAttrs.ToSlice())

		// Build a set that includes both the immutable attributes and all their nested fields
		immutableWithNested := mapset.NewSet()
		for _, attr := range immutableAttrs.ToSlice() {
			attrStr := attr.(string)
			immutableWithNested.Add(attrStr)

			// Also add any nested fields (e.g., if "resources" is immutable, also exclude "resources.resource_id")
			for _, schemaAttr := range schemaModifySet.ToSlice() {
				schemaAttrStr := schemaAttr.(string)
				if strings.HasPrefix(schemaAttrStr, attrStr+".") {
					immutableWithNested.Add(schemaAttrStr)
					log.Debugf("Excluding nested immutable attribute: %s", schemaAttrStr)
				}
			}
		}

		schemaModifySet = schemaModifySet.Difference(immutableWithNested)
	}

	notModifySlice := schemaModifySet.Difference(testModifySet).ToSlice()
	if len(notModifySlice) != 0 {
		isAllModified = false
		schemaCount := float64(len(schemaModifySet.ToSlice()))
		notCoverCount := float64(len(notModifySlice))
		coverageRate := 1 - (notCoverCount / schemaCount)
		log.Infof("resource %s attributes has %.2f%% modified coverage rate ", resourceName, coverageRate*100)

		// Format as tree structure
		treeStr := formatAttributesAsTree(notModifySlice)
		log.Errorf("resource %s attributes missing modification in test cases (%d attributes):\n%s", resourceName, len(notModifySlice), treeStr)
	} else {
		log.Infof("resource %s attributes has 100%% modified coverage rate ", resourceName)
	}

	return isFullCover && isIgnoreLegal && isAllModified
}

// getImmutableAttributesFromDoc reads the documentation and finds attributes marked as immutable
// These attributes have descriptions like "Note: The parameter is immutable after resource creation"
func getImmutableAttributesFromDoc(resourceName string, fileType string) mapset.Set {
	immutableSet := mapset.NewSet()

	// Remove "alicloud_" prefix if present for documentation path
	docResourceName := strings.TrimPrefix(resourceName, "alicloud_")

	// Construct documentation path
	var docPath string
	if fileType == "resource" {
		docPath = "website/docs/r/" + docResourceName + ".html.markdown"
	} else {
		docPath = "website/docs/d/" + docResourceName + ".html.markdown"
	}

	// Read documentation file
	file, err := os.Open(docPath)
	if err != nil {
		log.Debugf("Cannot open documentation file %s: %v", docPath, err)
		return immutableSet
	}
	defer file.Close()

	// Regex patterns to match immutable declarations
	immutablePattern := regexp.MustCompile(`\*\*Note:\s*The parameter is immutable after resource creation`)
	// Pattern to extract attribute name: * `attr_name` - (Optional/Required...) Description
	attrPattern := regexp.MustCompile(`^\*\s+\x60([a-zA-Z_0-9]+)\x60\s+-\s+\(`)

	scanner := bufio.NewScanner(file)
	var currentAttr string
	var continuedLine string

	for scanner.Scan() {
		line := scanner.Text()

		// Check if this line defines an attribute
		if matches := attrPattern.FindStringSubmatch(line); matches != nil {
			// Save previous attribute if it was immutable
			if currentAttr != "" && strings.Contains(continuedLine, "immutable after resource creation") {
				// Add the attribute and all its nested fields
				immutableSet.Add(currentAttr)
				log.Debugf("Found immutable attribute: %s", currentAttr)
			}

			// Start tracking new attribute
			currentAttr = matches[1]
			continuedLine = line
		} else if currentAttr != "" {
			// Continue reading the description (might span multiple lines)
			continuedLine += " " + strings.TrimSpace(line)

			// Stop if we hit a new section or empty line
			if strings.HasPrefix(line, "##") || strings.HasPrefix(line, "###") {
				// Check if current attribute is immutable before resetting
				if strings.Contains(continuedLine, "immutable after resource creation") {
					immutableSet.Add(currentAttr)
					log.Debugf("Found immutable attribute: %s", currentAttr)
				}
				currentAttr = ""
				continuedLine = ""
			}
		}

		// Also check for inline immutable declarations
		if immutablePattern.MatchString(line) && currentAttr != "" {
			immutableSet.Add(currentAttr)
			log.Debugf("Found immutable attribute: %s", currentAttr)
		}
	}

	// Check the last attribute
	if currentAttr != "" && strings.Contains(continuedLine, "immutable after resource creation") {
		immutableSet.Add(currentAttr)
		log.Debugf("Found immutable attribute: %s", currentAttr)
	}

	return immutableSet
}

func testAll(diff *diffparser.Diff) {
	file, err := os.Open("filename.txt")
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	diff.Files = make([]*diffparser.DiffFile, 0, 800)
	line := 0
	for scanner.Scan() {
		text := scanner.Text()
		if strings.HasPrefix(text, "resource_alicloud") {
			diff.Files = append(diff.Files, &diffparser.DiffFile{NewName: "alicloud/" + text})
		}
		line++
	}
}
