//nolint:all
package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/scripts/internal/splitlayout"
)

// TestNewRequiredFieldDetection tests if we can detect new required fields
func TestNewRequiredFieldDetection(t *testing.T) {
	oldAttrs := map[string]map[string]interface{}{
		"field1": {
			"Name":     "field1",
			"Type":     "TypeString",
			"Optional": true,
		},
	}

	newAttrs := map[string]map[string]interface{}{
		"field1": {
			"Name":     "field1",
			"Type":     "TypeString",
			"Optional": true,
		},
		"new_required_field": {
			"Name":     "new_required_field",
			"Type":     "TypeBool",
			"Required": true,
		},
	}

	hasBreaking := IsBreakingChange(oldAttrs, newAttrs)
	if !hasBreaking {
		t.Error("Expected breaking change for new required field, but got none")
	} else {
		fmt.Println("✓ Test passed: New required field detected")
	}
}

// TestAttributeDeletion tests if we can detect attribute deletion
func TestAttributeDeletion(t *testing.T) {
	oldAttrs := map[string]map[string]interface{}{
		"field1": {
			"Name":     "field1",
			"Type":     "TypeString",
			"Optional": true,
		},
		"field2": {
			"Name":     "field2",
			"Type":     "TypeString",
			"Optional": true,
		},
	}

	newAttrs := map[string]map[string]interface{}{
		"field1": {
			"Name":     "field1",
			"Type":     "TypeString",
			"Optional": true,
		},
	}

	hasBreaking := IsBreakingChange(oldAttrs, newAttrs)
	if !hasBreaking {
		t.Error("Expected breaking change for attribute deletion, but got none")
	} else {
		fmt.Println("✓ Test passed: Attribute deletion detected")
	}
}

// TestOptionalToRequired tests if we can detect optional->required change
func TestOptionalToRequired(t *testing.T) {
	oldAttrs := map[string]map[string]interface{}{
		"field1": {
			"Name":     "field1",
			"Type":     "TypeString",
			"Optional": true,
		},
	}

	newAttrs := map[string]map[string]interface{}{
		"field1": {
			"Name":     "field1",
			"Type":     "TypeString",
			"Required": true,
		},
	}

	hasBreaking := IsBreakingChange(oldAttrs, newAttrs)
	if !hasBreaking {
		t.Error("Expected breaking change for optional->required, but got none")
	} else {
		fmt.Println("✓ Test passed: Optional->Required detected")
	}
}

// TestTypeChange tests if we can detect type changes
func TestTypeChange(t *testing.T) {
	oldAttrs := map[string]map[string]interface{}{
		"count_field": {
			"Name":     "count_field",
			"Type":     "TypeInt",
			"Optional": true,
		},
	}

	newAttrs := map[string]map[string]interface{}{
		"count_field": {
			"Name":     "count_field",
			"Type":     "TypeString",
			"Optional": true,
		},
	}

	hasBreaking := IsBreakingChange(oldAttrs, newAttrs)
	if !hasBreaking {
		t.Error("Expected breaking change for type change, but got none")
	} else {
		fmt.Println("✓ Test passed: Type change detected")
	}
}

// TestForceNewAdded tests if we can detect ForceNew being added
func TestForceNewAdded(t *testing.T) {
	oldAttrs := map[string]map[string]interface{}{
		"field1": {
			"Name":     "field1",
			"Type":     "TypeString",
			"Optional": true,
		},
	}

	newAttrs := map[string]map[string]interface{}{
		"field1": {
			"Name":     "field1",
			"Type":     "TypeString",
			"Optional": true,
			"ForceNew": true,
		},
	}

	hasBreaking := IsBreakingChange(oldAttrs, newAttrs)
	if !hasBreaking {
		t.Error("Expected breaking change for ForceNew added, but got none")
	} else {
		fmt.Println("✓ Test passed: ForceNew addition detected")
	}
}

// TestSafeChange tests if safe changes are allowed
func TestSafeChange(t *testing.T) {
	oldAttrs := map[string]map[string]interface{}{
		"field1": {
			"Name":     "field1",
			"Type":     "TypeString",
			"Required": true,
		},
	}

	newAttrs := map[string]map[string]interface{}{
		"field1": {
			"Name":     "field1",
			"Type":     "TypeString",
			"Required": true,
		},
		"new_optional_field": {
			"Name":     "new_optional_field",
			"Type":     "TypeString",
			"Optional": true,
		},
	}

	hasBreaking := IsBreakingChange(oldAttrs, newAttrs)
	if hasBreaking {
		t.Error("Expected no breaking change for new optional field, but got one")
	} else {
		fmt.Println("✓ Test passed: New optional field allowed")
	}
}

// TestRetryCodeRemoval tests if retry error code removal is detected
func TestRetryCodeRemoval(t *testing.T) {
	oldCodes := map[string]map[string]struct{}{
		"CreateInstance": {
			"Throttling":         {},
			"ServiceUnavailable": {},
			"OperationConflict":  {},
		},
	}

	newCodes := map[string]map[string]struct{}{
		"CreateInstance": {
			"Throttling":         {},
			"ServiceUnavailable": {},
		},
	}

	hasBreaking := IsRetryCodeBreaking(oldCodes, newCodes)
	if !hasBreaking {
		t.Error("Expected breaking change for retry code removal, but got none")
	} else {
		fmt.Println("✓ Test passed: Retry code removal detected")
	}
}

// TestRetryCodeAddition tests if adding retry error codes is allowed (safe change)
func TestRetryCodeAddition(t *testing.T) {
	oldCodes := map[string]map[string]struct{}{
		"CreateInstance": {
			"Throttling": {},
		},
	}

	newCodes := map[string]map[string]struct{}{
		"CreateInstance": {
			"Throttling":         {},
			"ServiceUnavailable": {},
			"OperationConflict":  {},
		},
	}

	hasBreaking := IsRetryCodeBreaking(oldCodes, newCodes)
	if hasBreaking {
		t.Error("Expected no breaking change for retry code addition, but got one")
	} else {
		fmt.Println("✓ Test passed: Retry code addition allowed")
	}
}

// TestRetryCodeCompleteRemoval tests if complete removal of IsExpectedErrors is detected
func TestRetryCodeCompleteRemoval(t *testing.T) {
	oldCodes := map[string]map[string]struct{}{
		"CreateInstance": {
			"Throttling":         {},
			"ServiceUnavailable": {},
		},
	}

	newCodes := map[string]map[string]struct{}{
		// CreateInstance completely removed
	}

	hasBreaking := IsRetryCodeBreaking(oldCodes, newCodes)
	if !hasBreaking {
		t.Error("Expected breaking change for complete retry code removal, but got none")
	} else {
		fmt.Println("✓ Test passed: Complete retry code removal detected")
	}
}

// TestRetryCodeParsingFromContent tests parsing retry codes from actual Go code
func TestRetryCodeParsingFromContent(t *testing.T) {
	content := `
package alicloud

func resourceCreate(d *schema.ResourceData) error {
	action := "CreateInstance"
	if err := client.DoAction(action); err != nil {
		if IsExpectedErrors(err, []string{"Throttling", "ServiceUnavailable"}) {
			return resource.RetryableError(err)
		}
		return err
	}
	
	action = "DescribeInstance"
	if err := client.DoAction(action); err != nil {
		if IsExpectedErrors(err, []string{"NotFound", "InvalidId"}) {
			return resource.RetryableError(err)
		}
		return err
	}
	return nil
}
`

	codes := ParseRetryErrorCodesFromContent(content)

	// Should have 2 APIs
	if len(codes) != 2 {
		t.Errorf("Expected 2 APIs, got %d", len(codes))
		return
	}

	// Check CreateInstance codes
	if createCodes, ok := codes["CreateInstance"]; ok {
		if len(createCodes) != 2 {
			t.Errorf("Expected 2 codes for CreateInstance, got %d", len(createCodes))
		}
		if _, ok := createCodes["Throttling"]; !ok {
			t.Error("Expected Throttling code for CreateInstance")
		}
		if _, ok := createCodes["ServiceUnavailable"]; !ok {
			t.Error("Expected ServiceUnavailable code for CreateInstance")
		}
	} else {
		t.Error("CreateInstance not found in parsed codes")
	}

	// Check DescribeInstance codes
	if describeCodes, ok := codes["DescribeInstance"]; ok {
		if len(describeCodes) != 2 {
			t.Errorf("Expected 2 codes for DescribeInstance, got %d", len(describeCodes))
		}
		if _, ok := describeCodes["NotFound"]; !ok {
			t.Error("Expected NotFound code for DescribeInstance")
		}
		if _, ok := describeCodes["InvalidId"]; !ok {
			t.Error("Expected InvalidId code for DescribeInstance")
		}
	} else {
		t.Error("DescribeInstance not found in parsed codes")
	}

	fmt.Println("✓ Test passed: Retry code parsing from content")
}

// generateExtremeDistanceContent creates test content with 1000 lines between action definition and IsExpectedErrors call
func generateExtremeDistanceContent() string {
	var content strings.Builder

	// Action definition at the beginning
	content.WriteString(`action := "CreateTrail"
var request map[string]interface{}
var response map[string]interface{}
query := make(map[string]interface{})
var err error
request = make(map[string]interface{})

`)

	// Generate 1000 lines of request assignments
	for i := 1; i <= 1000; i++ {
		content.WriteString(fmt.Sprintf(`if v, ok := d.GetOk("test%d"); ok {
	request["Test%d"] = v
}
`, i, i))
	}

	// Add the IsExpectedErrors call at the end
	content.WriteString(`
wait := incrementalWait(3*time.Second, 5*time.Second)
err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
	response, err = client.RpcPost("Actiontrail", "2020-07-06", action, query, request, true)
	if err != nil {
		if IsExpectedErrors(err, []string{"InsufficientBucketPolicyException"}) || NeedRetry(err) {
			wait()
			return resource.RetryableError(err)
		}
		return resource.NonRetryableError(err)
	}
	return nil
})`)

	return content.String()
}

// extractContentFromDiff extracts old and new content from unified diff lines
func extractContentFromDiff(diffContent string) (oldContent, newContent string) {
	lines := strings.Split(diffContent, "\n")
	var oldLines, newLines []string

	for _, line := range lines {
		if strings.HasPrefix(line, "@@") || strings.HasPrefix(line, "diff ") ||
			strings.HasPrefix(line, "index ") || strings.HasPrefix(line, "---") ||
			strings.HasPrefix(line, "+++") {
			continue
		}
		if strings.HasPrefix(line, "-") {
			oldLines = append(oldLines, strings.TrimPrefix(line, "-"))
		} else if strings.HasPrefix(line, "+") {
			newLines = append(newLines, strings.TrimPrefix(line, "+"))
		} else if strings.HasPrefix(line, " ") {
			oldLines = append(oldLines, strings.TrimPrefix(line, " "))
			newLines = append(newLines, strings.TrimPrefix(line, " "))
		}
	}

	return strings.Join(oldLines, "\n"), strings.Join(newLines, "\n")
}

func TestParseRetryErrorCodesFromHunk(t *testing.T) {
	tests := []struct {
		name           string
		diffContent    string
		expectedOld    map[string][]string // context -> error codes
		expectedNew    map[string][]string // context -> error codes
		expectBreaking bool
	}{
		{
			name: "remove_single_error_code",
			diffContent: `diff --git a/alicloud/resource_test.go b/alicloud/resource_test.go
index 1234567..abcdefg 100644
--- a/alicloud/resource_test.go
+++ b/alicloud/resource_test.go
@@ -10,7 +10,7 @@ func resourceTest() {
 	action := "CreateInstance"
 	conn := client.EcsConn
 	wait := incrementalWait(3*time.Second, 3*time.Second)
-	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
-		raw, err := conn.DoAction(request, action)
-		if err != nil {
-			if IsExpectedErrors(err, []string{"Throttling", "SystemBusy", "ServiceUnavailable"}) {
-				return resource.RetryableError(err)
-			}
-			return resource.NonRetryableError(err)
-		}
+	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
+		raw, err := conn.DoAction(request, action)
+		if err != nil {
+			if IsExpectedErrors(err, []string{"Throttling", "SystemBusy"}) {
+				return resource.RetryableError(err)
+			}
+			return resource.NonRetryableError(err)
+		}`,
			expectedOld: map[string][]string{
				"CreateInstance": {"Throttling", "SystemBusy", "ServiceUnavailable"},
			},
			expectedNew: map[string][]string{
				"CreateInstance": {"Throttling", "SystemBusy"},
			},
			expectBreaking: true,
		},
		{
			name: "completely_remove_expected_errors_call",
			diffContent: `diff --git a/alicloud/resource_test.go b/alicloud/resource_test.go
index 1234567..abcdefg 100644
--- a/alicloud/resource_test.go
+++ b/alicloud/resource_test.go
@@ -10,10 +10,6 @@ func resourceTest() {
 	action := "DeleteInstance"
 	conn := client.EcsConn
 	wait := incrementalWait(3*time.Second, 3*time.Second)
-	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
-		raw, err := conn.DoAction(request, action)
-		if err != nil {
-			if IsExpectedErrors(err, []string{"InvalidInstanceId.NotFound"}) {
-				return resource.RetryableError(err)
-			}
-			return resource.NonRetryableError(err)
-		}
+	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
+		raw, err := conn.DoAction(request, action)
+		if err != nil {
+			return resource.NonRetryableError(err)
+		}`,
			expectedOld: map[string][]string{
				"DeleteInstance": {"InvalidInstanceId.NotFound"},
			},
			expectedNew:    map[string][]string{},
			expectBreaking: true,
		},
		{
			name: "add_new_error_codes_non_breaking",
			diffContent: `diff --git a/alicloud/resource_test.go b/alicloud/resource_test.go
index 1234567..abcdefg 100644
--- a/alicloud/resource_test.go
+++ b/alicloud/resource_test.go
@@ -10,7 +10,7 @@ func resourceTest() {
 	action := "UpdateInstance"
 	conn := client.EcsConn
 	wait := incrementalWait(3*time.Second, 3*time.Second)
-	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
-		raw, err := conn.DoAction(request, action)
-		if err != nil {
-			if IsExpectedErrors(err, []string{"Throttling"}) {
-				return resource.RetryableError(err)
-			}
-			return resource.NonRetryableError(err)
-		}
+	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
+		raw, err := conn.DoAction(request, action)
+		if err != nil {
+			if IsExpectedErrors(err, []string{"Throttling", "SystemBusy", "ServiceUnavailable"}) {
+				return resource.RetryableError(err)
+			}
+			return resource.NonRetryableError(err)
+		}`,
			expectedOld: map[string][]string{
				"UpdateInstance": {"Throttling"},
			},
			expectedNew: map[string][]string{
				"UpdateInstance": {"Throttling", "SystemBusy", "ServiceUnavailable"},
			},
			expectBreaking: false,
		},
		{
			name: "no_changes_stable_case",
			diffContent: `diff --git a/alicloud/resource_test.go b/alicloud/resource_test.go
index 1234567..abcdefg 100644
--- a/alicloud/resource_test.go
+++ b/alicloud/resource_test.go
@@ -10,7 +10,7 @@ func resourceTest() {
 	action := "DescribeInstance"
 	conn := client.EcsConn
 	wait := incrementalWait(3*time.Second, 3*time.Second)
 	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
 		raw, err := conn.DoAction(request, action)
 		if err != nil {
 			if IsExpectedErrors(err, []string{"Throttling", "SystemBusy"}) {
 				return resource.RetryableError(err)
 			}
 			return resource.NonRetryableError(err)
 		}`,
			expectedOld: map[string][]string{
				"DescribeInstance": {"Throttling", "SystemBusy"},
			},
			expectedNew: map[string][]string{
				"DescribeInstance": {"Throttling", "SystemBusy"},
			},
			expectBreaking: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Extract old and new content from diff
			oldContent, newContent := extractContentFromDiff(tt.diffContent)

			// Parse retry codes from content
			oldRetryCodes := ParseRetryErrorCodesFromContent(oldContent)
			newRetryCodes := ParseRetryErrorCodesFromContent(newContent)

			// Verify old codes
			if !verifyExpectedCodes(t, "old", oldRetryCodes, tt.expectedOld) {
				return
			}

			// Verify new codes
			if !verifyExpectedCodes(t, "new", newRetryCodes, tt.expectedNew) {
				return
			}

			// Test breaking change detection
			isBreaking := IsRetryCodeBreaking(oldRetryCodes, newRetryCodes)
			if isBreaking != tt.expectBreaking {
				t.Errorf("Expected breaking change: %t, got: %t", tt.expectBreaking, isBreaking)
			}
		})
	}
}

func verifyExpectedCodes(t *testing.T, prefix string, actual map[string]map[string]struct{}, expected map[string][]string) bool {
	if len(actual) != len(expected) {
		t.Errorf("%s codes: expected %d contexts, got %d", prefix, len(expected), len(actual))
		return false
	}

	for expectedContext, expectedCodes := range expected {
		actualCodes, exists := actual[expectedContext]
		if !exists {
			t.Errorf("%s codes: missing context '%s'", prefix, expectedContext)
			return false
		}

		if len(actualCodes) != len(expectedCodes) {
			t.Errorf("%s codes for context '%s': expected %d codes, got %d", prefix, expectedContext, len(expectedCodes), len(actualCodes))
			return false
		}

		for _, expectedCode := range expectedCodes {
			if _, exists := actualCodes[expectedCode]; !exists {
				t.Errorf("%s codes for context '%s': missing code '%s'", prefix, expectedContext, expectedCode)
				return false
			}
		}
	}

	return true
}

func TestIsRetryCodeBreaking(t *testing.T) {
	tests := []struct {
		name           string
		oldRetryCodes  map[string]map[string]struct{}
		newRetryCodes  map[string]map[string]struct{}
		expectBreaking bool
		description    string
	}{
		{
			name: "remove_error_codes",
			oldRetryCodes: map[string]map[string]struct{}{
				"CreateInstance": {
					"Throttling":         {},
					"SystemBusy":         {},
					"ServiceUnavailable": {},
				},
			},
			newRetryCodes: map[string]map[string]struct{}{
				"CreateInstance": {
					"Throttling": {},
					"SystemBusy": {},
				},
			},
			expectBreaking: true,
			description:    "Removing ServiceUnavailable error code should be detected as breaking change",
		},
		{
			name: "completely_remove_context",
			oldRetryCodes: map[string]map[string]struct{}{
				"DeleteInstance": {
					"InvalidInstanceId.NotFound": {},
				},
			},
			newRetryCodes:  map[string]map[string]struct{}{},
			expectBreaking: true,
			description:    "Completely removing IsExpectedErrors call should be detected as breaking change",
		},
		{
			name: "add_error_codes",
			oldRetryCodes: map[string]map[string]struct{}{
				"UpdateInstance": {
					"Throttling": {},
				},
			},
			newRetryCodes: map[string]map[string]struct{}{
				"UpdateInstance": {
					"Throttling":         {},
					"SystemBusy":         {},
					"ServiceUnavailable": {},
				},
			},
			expectBreaking: false,
			description:    "Adding error codes should not be detected as breaking change",
		},
		{
			name: "no_changes",
			oldRetryCodes: map[string]map[string]struct{}{
				"DescribeInstance": {
					"Throttling": {},
					"SystemBusy": {},
				},
			},
			newRetryCodes: map[string]map[string]struct{}{
				"DescribeInstance": {
					"Throttling": {},
					"SystemBusy": {},
				},
			},
			expectBreaking: false,
			description:    "No changes should not be detected as breaking change",
		},
		{
			name:           "empty_maps",
			oldRetryCodes:  map[string]map[string]struct{}{},
			newRetryCodes:  map[string]map[string]struct{}{},
			expectBreaking: false,
			description:    "Empty maps should not be detected as breaking change",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isBreaking := IsRetryCodeBreaking(tt.oldRetryCodes, tt.newRetryCodes)
			if isBreaking != tt.expectBreaking {
				t.Errorf("%s: expected breaking change: %t, got: %t", tt.description, tt.expectBreaking, isBreaking)
			}
		})
	}
}

func TestRegexPatterns(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected map[string][]string // context -> codes
	}{
		{
			name: "standard_expected_errors_call",
			content: `action := "TestAction"
if IsExpectedErrors(err, []string{"Throttling", "SystemBusy"}) {`,
			expected: map[string][]string{
				"TestAction": {"Throttling", "SystemBusy"},
			},
		},
		{
			name: "expected_errors_call_with_spaces",
			content: `action := "TestAction2"
if IsExpectedErrors(err,  []string{ "Throttling" , "SystemBusy" }) {`,
			expected: map[string][]string{
				"TestAction2": {"Throttling", "SystemBusy"},
			},
		},
		{
			name: "single_error_code",
			content: `action := "DeleteAction"
if IsExpectedErrors(err, []string{"InvalidInstanceId.NotFound"}) {`,
			expected: map[string][]string{
				"DeleteAction": {"InvalidInstanceId.NotFound"},
			},
		},
		{
			name: "with_action_context",
			content: `action := "CreateInstance"
if IsExpectedErrors(err, []string{"Throttling"}) {`,
			expected: map[string][]string{
				"CreateInstance": {"Throttling"},
			},
		},
		{
			name:    "extreme_distance_action_context_1000_lines",
			content: generateExtremeDistanceContent(),
			expected: map[string][]string{
				"CreateTrail": {"InsufficientBucketPolicyException"},
			},
		},
		{
			name: "function_context_with_action_definition",
			content: `func resourceAliCloudEcsInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "RunInstances"
	var request map[string]interface{}
	var response map[string]interface{}
	var err error
	
	wait := incrementalWait(1*time.Second, 1*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Ecs", "2014-05-26", action, nil, request, false)
		if err != nil {
			if NeedRetry(err) || IsExpectedErrors(err, []string{"IncorrectVSwitchStatus"}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})`,
			expected: map[string][]string{
				"RunInstances": {"IncorrectVSwitchStatus"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Parse retry codes from content using the actual implementation
			parsedCodes := ParseRetryErrorCodesFromContent(tt.content)

			// Verify parsed codes match expected
			if !verifyExpectedCodes(t, "parsed", parsedCodes, tt.expected) {
				t.Errorf("Failed to parse expected codes from content")
			}
		})
	}
}

func TestExactRdsErrorClassifierCompatibility(t *testing.T) {
	old := `action := "DeleteDatabase"
if IsExpectedErrors(err, []string{"InvalidDBName.NotFound"}) { return nil }`
	updated := `action := "DeleteDatabase"
if rdsErrorHasCode(deleteErr, "InvalidDBInstanceId.NotFound", "InvalidDBName.NotFound") { return nil }
if IsExpectedErrors(deleteErr, []string{"OperationDenied.DBInstanceStatus"}) { return retry(deleteErr) }`
	before := ParseRetryErrorCodesFromContent(old)
	after := ParseRetryErrorCodesFromContent(updated)
	if !verifyExpectedCodes(t, "updated", after, map[string][]string{"DeleteDatabase": {"InvalidDBInstanceId.NotFound", "InvalidDBName.NotFound", "OperationDenied.DBInstanceStatus"}}) {
		return
	}
	if IsRetryCodeBreaking(before, after) {
		t.Fatal("changing the classifier must not appear to remove a retained code")
	}
	removed := ParseRetryErrorCodesFromContent(strings.ReplaceAll(updated, `, "InvalidDBName.NotFound"`, ""))
	if !IsRetryCodeBreaking(before, removed) {
		t.Fatal("real removal must remain detectable")
	}
}

// frameworkEphemeralSrc is a plugin-framework ephemeral shaped like
// alicloud/service/kms/ephemeral_plaintext.go: map[string]schema.Attribute
// with StringAttribute/MapAttribute entries and a RequiresReplace modifier.
const frameworkEphemeralSrc = `package kms

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (e *plaintextEphemeralResource) Schema(ctx context.Context, req ephemeral.SchemaRequest, resp *ephemeral.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"ciphertext_blob": schema.StringAttribute{
				Required: true,
			},
			"encryption_context": schema.MapAttribute{
				Optional:    true,
				ElementType: types.StringType,
			},
			"plaintext": schema.StringAttribute{
				Computed:  true,
				Sensitive: true,
			},
			"key_id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"immutable_field": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}
`

func TestParseFrameworkSchemaFromAST(t *testing.T) {
	attrs := ParseFrameworkSchemaFromAST(frameworkEphemeralSrc)

	want := map[string]struct {
		typ      string
		required bool
		optional bool
		forceNew bool
	}{
		"ciphertext_blob":    {"TypeString", true, false, false},
		"encryption_context": {"TypeMap", false, true, false},
		"plaintext":          {"TypeString", false, false, false},
		"key_id":             {"TypeString", false, false, false},
		"immutable_field":    {"TypeString", false, true, true},
	}
	if len(attrs) != len(want) {
		t.Fatalf("ParseFrameworkSchemaFromAST returned %d attributes (%v), want %d", len(attrs), attrs, len(want))
	}
	for name, w := range want {
		attr, ok := attrs[name]
		if !ok {
			t.Errorf("attribute %s missing from parsed framework schema", name)
			continue
		}
		if attr["Type"] != w.typ {
			t.Errorf("attribute %s Type = %v, want %s", name, attr["Type"], w.typ)
		}
		if required, _ := attr["Required"].(bool); required != w.required {
			t.Errorf("attribute %s Required = %v, want %v", name, attr["Required"], w.required)
		}
		if optional, _ := attr["Optional"].(bool); optional != w.optional {
			t.Errorf("attribute %s Optional = %v, want %v", name, attr["Optional"], w.optional)
		}
		_, hasForceNew := attr["ForceNew"]
		if hasForceNew != w.forceNew {
			t.Errorf("attribute %s ForceNew = %v, want %v", name, hasForceNew, w.forceNew)
		}
	}
}

func TestParseSchemasMergesSDKv2AndFramework(t *testing.T) {
	src := `package mixed

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

var sdkFields = map[string]*schema.Schema{
	"sdk_field": {
		Type:     schema.TypeString,
		Optional: true,
	},
}

func (r *resource) Schema() {
	resp.Schema = fwschema.Schema{
		Attributes: map[string]fwschema.Attribute{
			"framework_field": fwschema.StringAttribute{
				Required: true,
			},
		},
	}
}
`
	attrs := parseSchemas(src)
	if _, ok := attrs["sdk_field"]; !ok {
		t.Error("parseSchemas lost the SDKv2 map[string]*schema.Schema attribute")
	}
	fw, ok := attrs["framework_field"]
	if !ok {
		t.Fatal("parseSchemas lost the framework map[string]fwschema.Attribute attribute")
	}
	if fw["Type"] != "TypeString" || fw["Required"] != true {
		t.Errorf("framework_field = %v, want TypeString Required", fw)
	}
}

func TestFrameworkSchemaBreakingChanges(t *testing.T) {
	oldSrc := `package kms

func (e *ephemeral) Schema() {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"keep_me": schema.StringAttribute{
				Optional: true,
			},
			"remove_me": schema.StringAttribute{
				Optional: true,
			},
			"loosen_later": schema.StringAttribute{
				Optional: true,
			},
		},
	}
}
`
	newSrc := `package kms

func (e *ephemeral) Schema() {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"keep_me": schema.StringAttribute{
				Optional: true,
			},
			"loosen_later": schema.StringAttribute{
				Required: true,
			},
			"new_optional": schema.StringAttribute{
				Optional: true,
			},
		},
	}
}
`
	if !IsBreakingChange(ParseFrameworkSchemaFromAST(oldSrc), ParseFrameworkSchemaFromAST(newSrc)) {
		t.Error("Expected framework removal + optional-to-required to be detected as breaking")
	}

	// A purely additive optional attribute is safe.
	safeNew := strings.ReplaceAll(oldSrc, `"remove_me": schema.StringAttribute{
				Optional: true,
			},`, `"remove_me": schema.StringAttribute{
				Optional: true,
			},
			"added_optional": schema.StringAttribute{
				Optional: true,
			},`)
	if IsBreakingChange(ParseFrameworkSchemaFromAST(oldSrc), ParseFrameworkSchemaFromAST(safeNew)) {
		t.Error("Expected adding an optional framework attribute to be safe")
	}
}

// TestMigratedEntityNoPhantomTypeBreak proves a migrated entity compares
// cleanly against its pre-split flat-layout content: schema.TypeString and
// schema.StringAttribute normalize to the same TypeString.
func TestMigratedEntityNoPhantomTypeBreak(t *testing.T) {
	flatSrc := `package alicloud

func resourceUser() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"user_principal_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"email": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}
`
	frameworkSrc := `package ims

func (r *userResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"user_principal_name": schema.StringAttribute{
				Required: true,
			},
			"email": schema.StringAttribute{
				Optional: true,
			},
		},
	}
}
`
	if IsBreakingChange(parseSchemas(flatSrc), parseSchemas(frameworkSrc)) {
		t.Error("A migrated SDKv2->framework entity must not report phantom breaking changes")
	}
}

// TestFactoryResolutionAndBreakingComparison exercises the factory→declaring-file
// resolver and the name-agnostic breaking-change comparator against the #11402
// alicloud_havip model: two public resource aliases (alicloud_havip and
// alicloud_vpc_ha_vip) register the same real factory resourceAliCloudVpcHaVip,
// declared in resource_alicloud_vpc_ha_vip.go, while a dead, never-registered
// resourceAlicloudHavip with a different retry set lives in the path-guessable
// resource_alicloud_havip.go. The resolver must follow the registered factory
// to its real declaring file — not the filename — and the comparator must stay
// name-agnostic across a split move that renames and package-qualifies the
// constructor. These are the breaking-change-checker regressions (AC1, AC2,
// AC3, AC4, AC8); the collection-order deletion regressions (AC5, AC6, AC7)
// live in the coverage package's select tests.
func TestFactoryResolutionAndBreakingComparison(t *testing.T) {
	// Root provider.go registration: both resource aliases point at the same
	// bare factory; the data-source alias of the same public name points at
	// its own factory, so kind is a separate namespace.
	providerSrc := []byte(`package alicloud

func Provider() *schema.Provider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"alicloud_havip":      resourceAliCloudVpcHaVip(),
			"alicloud_vpc_ha_vip": resourceAliCloudVpcHaVip(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"alicloud_havip": dataSourceAliCloudHavip(),
		},
	}
}
`)

	// The real factory, in resource_alicloud_vpc_ha_vip.go. Its retry set is
	// IncorrectStatus (no %s).
	realSrc := `package alicloud

func resourceAliCloudVpcHaVip() *schema.Resource {
	return &schema.Resource{Schema: map[string]*schema.Schema{
		"name": {Type: schema.TypeString, Optional: true},
	}}
}

func resourceAliCloudVpcHaVipCreate(d *schema.ResourceData, meta interface{}) error {
	action := "WaitHavip"
	if IsExpectedErrors(err, []string{"IncorrectStatus"}) {
		return err
	}
	return nil
}
`

	// The dead, never-registered factory in resource_alicloud_havip.go — the
	// file a type-name path guess would wrongly pick. Its retry set is
	// IncorrectStatus.%s, differing from the real factory.
	deadSrc := `package alicloud

func resourceAlicloudHavip() *schema.Resource {
	return &schema.Resource{Schema: map[string]*schema.Schema{
		"name": {Type: schema.TypeString, Optional: true},
	}}
}

func resourceAlicloudHavipCreate(d *schema.ResourceData, meta interface{}) error {
	action := "WaitHavip"
	if IsExpectedErrors(err, []string{"IncorrectStatus.%s"}) {
		return err
	}
	return nil
}
`

	impls := map[string][]byte{
		"alicloud/resource_alicloud_vpc_ha_vip.go": []byte(realSrc),
		"alicloud/resource_alicloud_havip.go":      []byte(deadSrc),
	}

	// AC1: two aliases resolve to the same real factory and its real declaring
	// file, not the dead factory file whose name matches the alias. The real
	// factory compares cleanly against itself; the dead factory file, if
	// mistakenly selected, yields a false breaking change.
	t.Run("AC1 alias-shared factory resolves to real file not dead one", func(t *testing.T) {
		factory, registered, err := resolveFactory("alicloud_havip", splitlayout.KindResource, providerSrc, nil)
		if err != nil || !registered || factory != "resourceAliCloudVpcHaVip" {
			t.Fatalf("resolveFactory alicloud_havip resource: factory=%q registered=%v err=%v", factory, registered, err)
		}
		path, err := findDeclaringFile(factory, impls)
		if err != nil || path != "alicloud/resource_alicloud_vpc_ha_vip.go" {
			t.Fatalf("findDeclaringFile real factory: path=%q err=%v", path, err)
		}
		deadPath, err := findDeclaringFile("resourceAlicloudHavip", impls)
		if err != nil || deadPath != "alicloud/resource_alicloud_havip.go" {
			t.Fatalf("findDeclaringFile dead factory: path=%q err=%v", deadPath, err)
		}
		if !compareBreakingChange(realSrc, realSrc) {
			t.Fatal("real factory compared to itself must pass")
		}
		// The dead factory's retry set differs from the real implementation, so
		// comparing it against the real implementation reports a breaking change
		// — the false positive the old filename guess produced.
		if compareBreakingChange(deadSrc, realSrc) {
			t.Fatal("dead factory file must report breaking vs the real implementation (false positive under filename guess)")
		}
	})

	// AC2: a registered factory that drops its retry codes is still flagged.
	t.Run("AC2 registered factory retry removal still fails", func(t *testing.T) {
		headNoRetry := `package alicloud

func resourceAliCloudVpcHaVip() *schema.Resource {
	return &schema.Resource{Schema: map[string]*schema.Schema{
		"name": {Type: schema.TypeString, Optional: true},
	}}
}

func resourceAliCloudVpcHaVipCreate(d *schema.ResourceData, meta interface{}) error {
	action := "WaitHavip"
	return nil
}
`
		if compareBreakingChange(realSrc, headNoRetry) {
			t.Fatal("removing a registered factory's retry codes must be flagged as breaking")
		}
	})

	// AC3: the comparator is not relaxed for real schema breaking changes.
	t.Run("AC3 real schema breaking changes still fail", func(t *testing.T) {
		optional := `package alicloud

func f() *schema.Resource {
	return &schema.Resource{Schema: map[string]*schema.Schema{
		"name": {Type: schema.TypeString, Optional: true},
	}}
}
`
		required := `package alicloud

func f() *schema.Resource {
	return &schema.Resource{Schema: map[string]*schema.Schema{
		"name": {Type: schema.TypeString, Required: true},
	}}
}
`
		if compareBreakingChange(optional, required) {
			t.Fatal("optional->required must be flagged as breaking")
		}
		forceNew := `package alicloud

func f() *schema.Resource {
	return &schema.Resource{Schema: map[string]*schema.Schema{
		"name": {Type: schema.TypeString, Optional: true, ForceNew: true},
	}}
}
`
		if compareBreakingChange(optional, forceNew) {
			t.Fatal("adding ForceNew must be flagged as breaking")
		}
		extra := `package alicloud

func f() *schema.Resource {
	return &schema.Resource{Schema: map[string]*schema.Schema{
		"name":  {Type: schema.TypeString, Optional: true},
		"extra": {Type: schema.TypeString, Optional: true},
	}}
}
`
		if compareBreakingChange(extra, optional) {
			t.Fatal("attribute deletion must be flagged as breaking")
		}
	})

	// AC4: a registered implementation moved with an export rename and package
	// qualification is still selected for comparison, not exempted as a new
	// entity; the comparator is name-agnostic so it compares cleanly.
	t.Run("AC4 factory rename move is compared not exempted", func(t *testing.T) {
		headExported := `package vpc

func ResourceAliCloudVpcHaVip() *schema.Resource {
	return &schema.Resource{Schema: map[string]*schema.Schema{
		"name": {Type: schema.TypeString, Optional: true},
	}}
}

func ResourceAliCloudVpcHaVipCreate(d *schema.ResourceData, meta interface{}) error {
	action := "WaitHavip"
	if IsExpectedErrors(err, []string{"IncorrectStatus"}) {
		return err
	}
	return nil
}
`
		factory, registered, err := resolveFactory("alicloud_havip", splitlayout.KindResource, providerSrc, nil)
		if err != nil || !registered || factory != "resourceAliCloudVpcHaVip" {
			t.Fatalf("baseline factory must be registered: factory=%q registered=%v err=%v", factory, registered, err)
		}
		if !compareBreakingChange(realSrc, headExported) {
			t.Fatal("factory rename + package move must compare cleanly (name-agnostic)")
		}
	})

	// AC8: kind is a separate namespace, resource aliases share a factory, and
	// an ambiguous dual-factory registration fails closed.
	t.Run("AC8 kind separation alias share and ambiguity", func(t *testing.T) {
		resFactory, resRegistered, err := resolveFactory("alicloud_havip", splitlayout.KindResource, providerSrc, nil)
		if err != nil || !resRegistered || resFactory != "resourceAliCloudVpcHaVip" {
			t.Fatalf("resource alicloud_havip: factory=%q registered=%v err=%v", resFactory, resRegistered, err)
		}
		dsFactory, dsRegistered, err := resolveFactory("alicloud_havip", splitlayout.KindDataSource, providerSrc, nil)
		if err != nil || !dsRegistered || dsFactory != "dataSourceAliCloudHavip" {
			t.Fatalf("data source alicloud_havip: factory=%q registered=%v err=%v", dsFactory, dsRegistered, err)
		}
		aliasFactory, _, err := resolveFactory("alicloud_vpc_ha_vip", splitlayout.KindResource, providerSrc, nil)
		if err != nil || aliasFactory != "resourceAliCloudVpcHaVip" {
			t.Fatalf("alias alicloud_vpc_ha_vip: factory=%q err=%v", aliasFactory, err)
		}
		if resFactory == dsFactory {
			t.Fatal("resource and data source factories must not collide by name")
		}
		ambiguous := []byte(`package alicloud

func Provider() *schema.Provider {
	return &schema.Provider{ResourcesMap: map[string]*schema.Resource{
		"alicloud_havip":      resourceAliCloudVpcHaVip(),
		"alicloud_vpc_ha_vip": resourceAliCloudVpcHaVip(),
		"alicloud_havip":      resourceOther(),
	}}
}
`)
		if _, _, err := resolveFactory("alicloud_havip", splitlayout.KindResource, ambiguous, nil); err == nil {
			t.Fatal("ambiguous dual factory registration must fail closed")
		}
	})
}

// gitInitRepo writes files into a fresh temp git repository, commits them as
// a single baseline, and returns the repo directory and the baseline commit
// SHA. The sources need only parse with parser.SkipObjectResolution; they do
// not have to compile, since the checkers only AST-parse and text-grep.
func gitInitRepo(t *testing.T, files map[string]string) (dir, rev string) {
	t.Helper()
	dir = t.TempDir()
	for path, content := range files {
		full := filepath.Join(dir, filepath.FromSlash(path))
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatalf("mkdir %s: %v", path, err)
		}
		if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
			t.Fatalf("write %s: %v", path, err)
		}
	}
	for _, args := range [][]string{
		{"init", "-q"},
		{"config", "user.email", "checker@example.com"},
		{"config", "user.name", "checker"},
		{"-c", "commit.gpgsign=false", "add", "-A"},
		{"-c", "commit.gpgsign=false", "commit", "-q", "-m", "baseline"},
	} {
		c := exec.Command("git", args...)
		c.Dir = dir
		if out, err := c.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v: %s", args, err, out)
		}
	}
	c := exec.Command("git", "rev-parse", "HEAD")
	c.Dir = dir
	out, err := c.Output()
	if err != nil {
		t.Fatalf("git rev-parse HEAD: %v", err)
	}
	return dir, strings.TrimSpace(string(out))
}

// withGitRepo builds a temp git repo, changes the process working directory
// into it (the checker's git/show/grep run against CWD), runs fn, and restores
// the original directory. Tests do not run in parallel, so the global chdir is
// safe and always restored.
func withGitRepo(t *testing.T, files map[string]string, fn func(dir, rev string)) {
	t.Helper()
	dir, rev := gitInitRepo(t, files)
	orig, err := os.Getwd()
	if err != nil {
		t.Fatalf("os.Getwd: %v", err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("os.Chdir %s: %v", dir, err)
	}
	defer func() { _ = os.Chdir(orig) }()
	fn(dir, rev)
}

// TestLegacyImplGitFullPathFromTempRepo drives the full
// resolveLegacyImplGit -> registersAt -> resolveFactory -> grepDeclaringFiles
// -> getFileContentFromGit -> findDeclaringFile chain against a real temp git
// repository at a baseline commit. Unlike the in-memory map tests above, this
// exercises the git grep "<rev>:<path>" prefix and the registersAt
// absent-vs-error distinction through the actual git commands, not just a map
// handed to findDeclaringFile.
func TestLegacyImplGitFullPathFromTempRepo(t *testing.T) {
	providerSrc := `package alicloud

func Provider() *schema.Provider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"alicloud_havip":      resourceAliCloudVpcHaVip(),
			"alicloud_vpc_ha_vip": resourceAliCloudVpcHaVip(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"alicloud_havip": dataSourceAliCloudHavip(),
		},
	}
}
`
	realSrc := `package alicloud

func resourceAliCloudVpcHaVip() *schema.Resource {
	return &schema.Resource{Schema: map[string]*schema.Schema{
		"name": {Type: schema.TypeString, Optional: true},
	}}
}

func resourceAliCloudVpcHaVipCreate(d *schema.ResourceData, meta interface{}) error {
	action := "WaitHavip"
	if IsExpectedErrors(err, []string{"IncorrectStatus"}) {
		return err
	}
	return nil
}
`
	deadSrc := `package alicloud

func resourceAlicloudHavip() *schema.Resource {
	return &schema.Resource{Schema: map[string]*schema.Schema{
		"name": {Type: schema.TypeString, Optional: true},
	}}
}

func resourceAlicloudHavipCreate(d *schema.ResourceData, meta interface{}) error {
	action := "WaitHavip"
	if IsExpectedErrors(err, []string{"IncorrectStatus.%s"}) {
		return err
	}
	return nil
}
`
	dataSrc := `package alicloud

func dataSourceAliCloudHavip() *schema.Resource {
	return &schema.Resource{Schema: map[string]*schema.Schema{
		"id": {Type: schema.TypeString, Optional: true},
	}}
}
`
	// A service register present at the baseline so registersAt exercises the
	// present-and-readable path, not just absence. It registers nothing the
	// resolver cares about; it only needs to be readable from git.
	vpcRegister := "package vpc\n\nfunc init() {}\n"

	files := map[string]string{
		"alicloud/provider.go":                     providerSrc,
		"alicloud/resource_alicloud_vpc_ha_vip.go": realSrc,
		"alicloud/resource_alicloud_havip.go":      deadSrc,
		"alicloud/data_source_alicloud_havip.go":   dataSrc,
		"alicloud/service/vpc/register.go":         vpcRegister,
	}

	withGitRepo(t, files, func(dir, rev string) {
		// git grep with a tree revision emits "<rev>:<path>" entries;
		// grepDeclaringFiles must strip the revision prefix and yield plain
		// alicloud/ paths. Without the strip, getFileContentFromGit would run
		// "git show rev:rev:path" and exit 128.
		t.Run("git grep candidates carry no revision prefix", func(t *testing.T) {
			paths, err := grepDeclaringFiles("resourceAliCloudVpcHaVip", rev)
			if err != nil {
				t.Fatalf("grepDeclaringFiles: %v", err)
			}
			if len(paths) != 1 || paths[0] != "alicloud/resource_alicloud_vpc_ha_vip.go" {
				t.Fatalf("paths=%#v want [alicloud/resource_alicloud_vpc_ha_vip.go]", paths)
			}
			for _, p := range paths {
				if strings.HasPrefix(p, rev+":") || !strings.HasPrefix(p, "alicloud/") {
					t.Fatalf("candidate path not stripped or illegal: %q", p)
				}
			}
		})

		// AC1 via full git path: the alias resolves to the real factory's
		// declaring file, not the dead file whose name matches the alias.
		t.Run("AC alias resolves real declaring file via full git path", func(t *testing.T) {
			path, content, err := resolveLegacyImplGit("alicloud_havip", splitlayout.KindResource, rev)
			if err != nil {
				t.Fatalf("resolveLegacyImplGit: %v", err)
			}
			if path != "alicloud/resource_alicloud_vpc_ha_vip.go" {
				t.Fatalf("path=%q want alicloud/resource_alicloud_vpc_ha_vip.go", path)
			}
			if !strings.Contains(content, "func resourceAliCloudVpcHaVip") {
				t.Fatalf("content does not declare the real factory: %q", content)
			}
		})

		// AC: a genuinely new (unregistered) entity returns an empty path with
		// no error, not a guessed file.
		t.Run("AC unregistered new entity returns empty no error", func(t *testing.T) {
			path, _, err := resolveLegacyImplGit("alicloud_brand_new_type", splitlayout.KindResource, rev)
			if err != nil {
				t.Fatalf("unregistered type must not error: %v", err)
			}
			if path != "" {
				t.Fatalf("unregistered type path=%q want empty", path)
			}
		})

		// AC: kind is a separate namespace; the data-source alias resolves its
		// own factory file through the full git path.
		t.Run("AC data source kind resolves its declaring file", func(t *testing.T) {
			path, content, err := resolveLegacyImplGit("alicloud_havip", splitlayout.KindDataSource, rev)
			if err != nil {
				t.Fatalf("resolveLegacyImplGit data source: %v", err)
			}
			if path != "alicloud/data_source_alicloud_havip.go" {
				t.Fatalf("path=%q want alicloud/data_source_alicloud_havip.go", path)
			}
			if !strings.Contains(content, "func dataSourceAliCloudHavip") {
				t.Fatalf("content does not declare the data source factory: %q", content)
			}
		})

		// AC2/AC3 via resolved git baseline: the comparator is not relaxed —
		// retry removal and optional->required are still flagged.
		t.Run("AC breaking changes flagged from resolved git baseline", func(t *testing.T) {
			_, base, err := resolveLegacyImplGit("alicloud_havip", splitlayout.KindResource, rev)
			if err != nil {
				t.Fatalf("resolve baseline: %v", err)
			}
			retryRemoved := `package alicloud

func resourceAliCloudVpcHaVip() *schema.Resource {
	return &schema.Resource{Schema: map[string]*schema.Schema{
		"name": {Type: schema.TypeString, Optional: true},
	}}
}

func resourceAliCloudVpcHaVipCreate(d *schema.ResourceData, meta interface{}) error {
	action := "WaitHavip"
	return nil
}
`
			if compareBreakingChange(base, retryRemoved) {
				t.Fatal("retry removal from resolved baseline must be flagged breaking")
			}
			optionalToRequired := `package alicloud

func resourceAliCloudVpcHaVip() *schema.Resource {
	return &schema.Resource{Schema: map[string]*schema.Schema{
		"name": {Type: schema.TypeString, Required: true},
	}}
}
`
			if compareBreakingChange(base, optionalToRequired) {
				t.Fatal("optional->required from resolved baseline must be flagged breaking")
			}
		})

		// AC4 via resolved git baseline: a factory rename + package move
		// compares cleanly (name-agnostic), not exempted as a new entity.
		t.Run("AC factory rename move compares clean from git baseline", func(t *testing.T) {
			_, base, err := resolveLegacyImplGit("alicloud_havip", splitlayout.KindResource, rev)
			if err != nil {
				t.Fatalf("resolve baseline: %v", err)
			}
			renamed := `package vpc

func ResourceAliCloudVpcHaVip() *schema.Resource {
	return &schema.Resource{Schema: map[string]*schema.Schema{
		"name": {Type: schema.TypeString, Optional: true},
	}}
}

func ResourceAliCloudVpcHaVipCreate(d *schema.ResourceData, meta interface{}) error {
	action := "WaitHavip"
	if IsExpectedErrors(err, []string{"IncorrectStatus"}) {
		return err
	}
	return nil
}
`
			if !compareBreakingChange(base, renamed) {
				t.Fatal("factory rename move must compare cleanly against resolved baseline")
			}
		})

		// registersAt distinction: a register present at the revision is read
		// into the map; a register that exists on disk (Glob returns it) but is
		// absent at the revision is skipped without error. The old `continue`
		// masked both; the fix distinguishes expected absence from a read error.
		t.Run("AC registersAt present read and absent skipped no error", func(t *testing.T) {
			// Untracked register on disk but not in the baseline commit.
			newProd := filepath.Join(dir, "alicloud/service/newprod/register.go")
			if err := os.MkdirAll(filepath.Dir(newProd), 0o755); err != nil {
				t.Fatalf("mkdir newprod: %v", err)
			}
			if err := os.WriteFile(newProd, []byte("package newprod\n\nfunc init() {}\n"), 0o644); err != nil {
				t.Fatalf("write newprod: %v", err)
			}
			regs, err := registersAt(rev)
			if err != nil {
				t.Fatalf("registersAt must not error when a register is absent at rev: %v", err)
			}
			if _, ok := regs["alicloud/service/vpc/register.go"]; !ok {
				t.Fatal("present register alicloud/service/vpc/register.go must be read into the map")
			}
			if _, ok := regs["alicloud/service/newprod/register.go"]; ok {
				t.Fatal("absent-at-rev register alicloud/service/newprod/register.go must be skipped, not stored")
			}
		})
	})
}

// TestLegacyImplGitFailClosedFromTempRepo verifies the resolver fails closed
// (never guesses a file) for ambiguous, registered-but-undeclared and broken
// git state, all through the real git commands in a temp repository.
func TestLegacyImplGitFailClosedFromTempRepo(t *testing.T) {
	// A factory declared in two files: git grep returns both rev-prefixed
	// paths; findDeclaringFile fails closed (multiple declarations).
	t.Run("factory declared in two files fails closed", func(t *testing.T) {
		providerSrc := `package alicloud

func Provider() *schema.Provider {
	return &schema.Provider{ResourcesMap: map[string]*schema.Resource{
		"alicloud_dup": resourceDup(),
	}}
}
`
		dupA := `package alicloud

func resourceDup() *schema.Resource {
	return &schema.Resource{Schema: map[string]*schema.Schema{"a": {Type: schema.TypeString, Optional: true}}}
}
`
		dupB := `package alicloud

func resourceDup() *schema.Resource {
	return &schema.Resource{Schema: map[string]*schema.Schema{"b": {Type: schema.TypeString, Optional: true}}}
}
`
		files := map[string]string{
			"alicloud/provider.go":       providerSrc,
			"alicloud/resource_dup_a.go": dupA,
			"alicloud/resource_dup_b.go": dupB,
		}
		withGitRepo(t, files, func(dir, rev string) {
			_, _, err := resolveLegacyImplGit("alicloud_dup", splitlayout.KindResource, rev)
			if err == nil {
				t.Fatal("factory declared in two files must fail closed, got nil error")
			}
		})
	})

	// A registered factory with no declaring file at the baseline: git grep
	// matches nothing (exit 1) -> fail closed, not a filename guess.
	t.Run("registered but undeclared factory fails closed", func(t *testing.T) {
		providerSrc := `package alicloud

func Provider() *schema.Provider {
	return &schema.Provider{ResourcesMap: map[string]*schema.Resource{
		"alicloud_ghost": resourceGhost(),
	}}
}
`
		files := map[string]string{
			"alicloud/provider.go": providerSrc,
		}
		withGitRepo(t, files, func(dir, rev string) {
			_, _, err := resolveLegacyImplGit("alicloud_ghost", splitlayout.KindResource, rev)
			if err == nil {
				t.Fatal("registered-but-undeclared factory must fail closed, got nil error")
			}
		})
	})

	// An invalid revision is a real git read error: provider.go cannot be read,
	// so resolveLegacyImplGit fails closed at the first git access rather than
	// returning an empty (silently skipped) result.
	t.Run("bad revision fails closed at provider read", func(t *testing.T) {
		files := map[string]string{
			"alicloud/provider.go": "package alicloud\n",
		}
		withGitRepo(t, files, func(dir, rev string) {
			_, _, err := resolveLegacyImplGit("alicloud_any", splitlayout.KindResource, "no-such-revision-zzz")
			if err == nil {
				t.Fatal("bad revision must fail closed, got nil error")
			}
		})
	})

	// D2: a type registered at the baseline whose provider.go map value is not a
	// zero-argument call on a plain identifier (a helper, a variable, or a call
	// with arguments) is still registered by name. resolveFactory cannot follow
	// such a value (providerFactories omits it), so the old code waived it as a
	// genuinely-new entity and exempted it from comparison.
	// ProviderRegisteredNames records it by name, so the resolver must fail
	// closed rather than guess a factory expression or declaring file; it does
	// not widen to support arbitrary expressions.
	t.Run("registered but unresolvable constructor fails closed", func(t *testing.T) {
		// newWidgetResource(ctx) is a call with an argument: not a followable
		// zero-argument plain-identifier call, but ProviderRegisteredNames
		// records alicloud_widget by name and kind.
		providerSrc := `package alicloud

func Provider() *schema.Provider {
	return &schema.Provider{ResourcesMap: map[string]*schema.Resource{
		"alicloud_widget": newWidgetResource(ctx),
	}}
}
`
		widgetSrc := `package alicloud

func newWidgetResource(ctx context.Context) *schema.Resource {
	return &schema.Resource{Schema: map[string]*schema.Schema{
		"name": {Type: schema.TypeString, Optional: true},
	}}
}
`
		files := map[string]string{
			"alicloud/provider.go":        providerSrc,
			"alicloud/resource_widget.go": widgetSrc,
		}
		withGitRepo(t, files, func(dir, rev string) {
			// Control: providerFactories does not follow the arg-bearing call,
			// so resolveFactory reports the type as not registered by
			// construction — the口径 narrower than ProviderRegisteredNames.
			factory, registered, err := resolveFactory("alicloud_widget", splitlayout.KindResource, []byte(providerSrc), nil)
			if err != nil || registered || factory != "" {
				t.Fatalf("resolveFactory must not follow an arg-bearing factory value: factory=%q registered=%v err=%v", factory, registered, err)
			}
			// ProviderRegisteredNames nonetheless records the registration by name.
			if !registeredByName([]byte(providerSrc), "alicloud_widget", splitlayout.KindResource) {
				t.Fatal("registeredByName must report alicloud_widget as registered by name")
			}
			// resolveLegacyImplGit must fail closed: the type is registered by
			// name at the baseline, so it cannot be waived as a genuinely-new
			// entity and exempted from comparison.
			if _, _, err := resolveLegacyImplGit("alicloud_widget", splitlayout.KindResource, rev); err == nil {
				t.Fatal("registered-but-unresolvable constructor must fail closed, got nil error")
			}
		})
	})

	// D2 control: a genuinely unregistered type still returns an empty path
	// with no error (a genuinely new entity), so the ProviderRegisteredNames
	// cross-check does not over-close.
	t.Run("genuinely unregistered type still returns empty no error", func(t *testing.T) {
		providerSrc := `package alicloud

func Provider() *schema.Provider {
	return &schema.Provider{ResourcesMap: map[string]*schema.Resource{
		"alicloud_widget": newWidgetResource(ctx),
	}}
}
`
		files := map[string]string{
			"alicloud/provider.go": providerSrc,
		}
		withGitRepo(t, files, func(dir, rev string) {
			path, _, err := resolveLegacyImplGit("alicloud_truly_absent", splitlayout.KindResource, rev)
			if err != nil {
				t.Fatalf("genuinely unregistered type must not error: %v", err)
			}
			if path != "" {
				t.Fatalf("genuinely unregistered type path=%q want empty", path)
			}
		})
	})
}

// TestRegistersAtTreeExistenceDistinguishesMissingBlobFromAbsent verifies the
// D1 fix through the real git commands: registersAt must distinguish a register
// genuinely absent from the baseline tree (skip without error) from a register
// whose tree entry still exists but whose blob object is unreadable (fail
// closed). The old isFileNew short-circuit ran git cat-file -e, which exits 1
// for a missing blob object exactly as it does for an absent path, so a broken
// git state was silently skipped like an expected absence and could hide a real
// registration. pathExistsInTree inspects the tree entry (git ls-tree) instead,
// which still lists a path whose blob object is missing.
func TestRegistersAtTreeExistenceDistinguishesMissingBlobFromAbsent(t *testing.T) {
	providerSrc := `package alicloud

func Provider() *schema.Provider {
	return &schema.Provider{ResourcesMap: map[string]*schema.Resource{
		"alicloud_vpc_ha_vip": resourceAliCloudVpcHaVip(),
	}}
}
`
	realSrc := `package alicloud

func resourceAliCloudVpcHaVip() *schema.Resource {
	return &schema.Resource{Schema: map[string]*schema.Schema{
		"name": {Type: schema.TypeString, Optional: true},
	}}
}
`
	vpcRegister := "package vpc\n\nfunc init() {}\n"
	files := map[string]string{
		"alicloud/provider.go":                     providerSrc,
		"alicloud/resource_alicloud_vpc_ha_vip.go": realSrc,
		"alicloud/service/vpc/register.go":         vpcRegister,
	}
	withGitRepo(t, files, func(dir, rev string) {
		registerPath := "alicloud/service/vpc/register.go"

		// Sanity: the committed register is present in the tree and readable.
		exists, err := pathExistsInTree(registerPath, rev)
		if err != nil || !exists {
			t.Fatalf("pathExistsInTree present register: exists=%v err=%v", exists, err)
		}
		regs, err := registersAt(rev)
		if err != nil {
			t.Fatalf("registersAt before corruption must succeed: %v", err)
		}
		if _, ok := regs[registerPath]; !ok {
			t.Fatalf("present register %s must be read into the map", registerPath)
		}

		// A truly absent path (present on disk but not committed) is not in the
		// tree: pathExistsInTree returns (false, nil), and registersAt skips it
		// without error.
		absentPath := filepath.Join(dir, "alicloud/service/newprod/register.go")
		if err := os.MkdirAll(filepath.Dir(absentPath), 0o755); err != nil {
			t.Fatalf("mkdir newprod: %v", err)
		}
		if err := os.WriteFile(absentPath, []byte("package newprod\n\nfunc init() {}\n"), 0o644); err != nil {
			t.Fatalf("write newprod: %v", err)
		}
		absentExists, err := pathExistsInTree("alicloud/service/newprod/register.go", rev)
		if err != nil {
			t.Fatalf("pathExistsInTree absent register must not error: %v", err)
		}
		if absentExists {
			t.Fatal("pathExistsInTree must report an uncommitted register as absent from the tree")
		}
		if _, err := registersAt(rev); err != nil {
			t.Fatalf("registersAt must skip a tree-absent register without error: %v", err)
		}

		// Corrupt the git state: remove the blob object for the committed
		// register while leaving its tree entry in place. ls-tree still lists
		// the path (it reads the tree object, not the blob), but cat-file -e
		// and git show fail on the missing blob.
		blobSHA := gitRevParseBlob(t, dir, rev+":"+registerPath)
		removeLooseObject(t, dir, blobSHA)

		// pathExistsInTree still sees the tree entry — the path is in the tree.
		exists, err = pathExistsInTree(registerPath, rev)
		if err != nil || !exists {
			t.Fatalf("pathExistsInTree must still see the tree entry after blob removal: exists=%v err=%v", exists, err)
		}
		// The old isFileNew (git cat-file -e) reports the missing-blob path as
		// new — documenting the flaw the fix replaces.
		if !isFileNew(registerPath, rev) {
			t.Fatal("isFileNew (cat-file -e) must report a missing-blob path as new; the flaw pathExistsInTree replaces")
		}

		// registersAt must now fail closed: the tree lists the register but its
		// blob cannot be read, which is a broken git state, not an expected
		// absence. The old code skipped it silently via isFileNew.
		if _, err := registersAt(rev); err == nil {
			t.Fatal("registersAt must fail closed when a register blob is missing but its tree entry remains")
		}
		// resolveLegacyImplGit propagates the fail-closed error rather than
		// returning an empty (silently skipped) result that would waive the
		// type as a genuinely-new entity.
		if _, _, err := resolveLegacyImplGit("alicloud_vpc_ha_vip", splitlayout.KindResource, rev); err == nil {
			t.Fatal("resolveLegacyImplGit must propagate the missing-blob register error, got nil")
		}
	})
}

// gitRevParseBlob returns the object SHA for a <rev>:<path> spec in dir's repo.
func gitRevParseBlob(t *testing.T, dir, spec string) string {
	t.Helper()
	c := exec.Command("git", "rev-parse", spec)
	c.Dir = dir
	out, err := c.Output()
	if err != nil {
		t.Fatalf("git rev-parse %s: %v", spec, err)
	}
	return strings.TrimSpace(string(out))
}

// removeLooseObject deletes the loose object file for sha from dir's object
// store. A freshly-initialized repository stores every object loose, so the
// blob remains reachable from its tree entry while its object is unreadable.
func removeLooseObject(t *testing.T, dir, sha string) {
	t.Helper()
	obj := filepath.Join(dir, ".git", "objects", sha[:2], sha[2:])
	if err := os.Remove(obj); err != nil {
		t.Fatalf("remove loose object %s: %v (object may be packed; test needs a loose-object fresh repo)", sha, err)
	}
}
