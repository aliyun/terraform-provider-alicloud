//nolint:all
package main

import (
	"fmt"
	"strings"
	"testing"
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

// TestUnitNestedOptionalBlockRequiredNotBreaking covers acceptance (a): a Required
// subfield declared inside a newly added Optional parent block must not be
// reported as a new top-level required attribute, because configurations
// written before the block existed never reference it.
func TestUnitNestedOptionalBlockRequiredNotBreaking(t *testing.T) {
	oldSrc := `package ecs

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func ResourceEip() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"bandwidth": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}
`
	newSrc := `package ecs

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func ResourceEip() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"bandwidth": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"additional_bandwidth": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bandwidth": {
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
			},
		},
	}
}
`
	newAttrs := ParseSchemaFromAST(newSrc)
	if _, ok := newAttrs["additional_bandwidth.bandwidth"]; !ok {
		t.Fatalf("nested required subfield was flattened away; want key additional_bandwidth.bandwidth, got %v", newAttrs)
	}
	if _, ok := newAttrs["bandwidth"]; !ok {
		t.Fatalf("top-level bandwidth must remain under its bare name; got %v", newAttrs)
	}
	if IsBreakingChange(ParseSchemaFromAST(oldSrc), newAttrs) {
		t.Error("a Required subfield under a newly added Optional block must not be a breaking change")
	} else {
		fmt.Println("✓ Test passed: new optional block with a required subfield is safe")
	}
}

// TestUnitNestedRequiredUnderExistingBlockBreaking covers acceptance (b): a new
// Required subfield added under a block that already existed does invalidate
// old configurations that use the block, so it must stay a breaking change.
func TestUnitNestedRequiredUnderExistingBlockBreaking(t *testing.T) {
	oldSrc := `package ecs

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func ResourceEip() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"additional_bandwidth": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"other": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}
`
	newSrc := `package ecs

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func ResourceEip() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"additional_bandwidth": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"other": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"bandwidth": {
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
			},
		},
	}
}
`
	if !IsBreakingChange(ParseSchemaFromAST(oldSrc), ParseSchemaFromAST(newSrc)) {
		t.Error("a new Required subfield under an existing block must be a breaking change")
	} else {
		fmt.Println("✓ Test passed: required subfield under existing block is breaking")
	}
}

// TestUnitNestedSameNameFieldsComparedSeparately covers acceptance (c): two
// subfields that share a name but live under different parent blocks are
// captured and judged under their full parent path. The same name "shared" is
// breaking under block_a (an existing field flipped Optional->Required) while
// it is waived under block_b (a new Optional block); a flat map would have
// collapsed them and could not tell the two apart.
func TestUnitNestedSameNameFieldsComparedSeparately(t *testing.T) {
	oldSrc := `package ecs

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func ResourceEip() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"block_a": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"shared": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}
`
	newSrc := `package ecs

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func ResourceEip() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"block_a": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"shared": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"block_b": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"shared": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}
`
	newAttrs := ParseSchemaFromAST(newSrc)
	if _, ok := newAttrs["block_a.shared"]; !ok {
		t.Error("block_a.shared must be captured with its parent path")
	}
	if _, ok := newAttrs["block_b.shared"]; !ok {
		t.Error("block_b.shared must be captured with its parent path")
	}
	if _, ok := newAttrs["shared"]; ok {
		t.Error("same-named nested subfields must not collapse onto the bare name 'shared'")
	}
	if !IsBreakingChange(ParseSchemaFromAST(oldSrc), newAttrs) {
		t.Error("block_a.shared flipped Optional->Required must be breaking even when block_b.shared is waived")
	} else {
		fmt.Println("✓ Test passed: same-named nested fields compared per parent path")
	}
}

// TestUnitBreakingChangeSmoke is a TestUnit* so ci-check-quick compiles the
// compatibility package: the breaking-change checker is launched by PR CI via
// go run rather than by local-ci-check.sh, so a unit test is what brings it
// under the quick gate alongside the resource-name helper. It pins the
// non-breaking baseline for empty schemas.
func TestUnitBreakingChangeSmoke(t *testing.T) {
	if IsBreakingChange(map[string]map[string]interface{}{}, map[string]map[string]interface{}{}) {
		t.Fatal("empty schemas must not be reported as a breaking change")
	}
}
