//nolint:all
package main

import (
	"fmt"
	"strings"
	"testing"
)

// TestNewRequiredFieldDetection tests if we can detect new required fields
func TestUnitNewRequiredFieldDetection(t *testing.T) {
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
func TestUnitAttributeDeletion(t *testing.T) {
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
func TestUnitOptionalToRequired(t *testing.T) {
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
func TestUnitTypeChange(t *testing.T) {
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
func TestUnitForceNewAdded(t *testing.T) {
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
func TestUnitSafeChange(t *testing.T) {
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
func TestUnitRetryCodeRemoval(t *testing.T) {
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
func TestUnitRetryCodeAddition(t *testing.T) {
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
func TestUnitRetryCodeCompleteRemoval(t *testing.T) {
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
func TestUnitRetryCodeParsingFromContent(t *testing.T) {
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

func TestUnitParseRetryErrorCodesFromHunk(t *testing.T) {
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

func TestUnitIsRetryCodeBreaking(t *testing.T) {
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

func TestUnitRegexPatterns(t *testing.T) {
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

func TestUnitExactRdsErrorClassifierCompatibility(t *testing.T) {
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

func TestUnitParseFrameworkSchemaFromAST(t *testing.T) {
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

func TestUnitParseSchemasMergesSDKv2AndFramework(t *testing.T) {
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

func TestUnitFrameworkSchemaBreakingChanges(t *testing.T) {
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
func TestUnitMigratedEntityNoPhantomTypeBreak(t *testing.T) {
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

// TestNestedRequiredChildUnderExistingOptionalParent (path 2): an existing
// Optional parent that gains a new Required child is a breaking change because
// existing users of the parent must now supply the new child.
func TestUnitNestedRequiredChildUnderExistingOptionalParent(t *testing.T) {
	oldAttrs := map[string]map[string]interface{}{
		"resource_tag": {
			"Name":     "resource_tag",
			"Type":     "TypeSet",
			"Optional": true,
			"Nested": map[string]map[string]interface{}{
				"key": {"Name": "key", "Type": "TypeString", "Required": true},
			},
		},
	}
	newAttrs := map[string]map[string]interface{}{
		"resource_tag": {
			"Name":     "resource_tag",
			"Type":     "TypeSet",
			"Optional": true,
			"Nested": map[string]map[string]interface{}{
				"key":   {"Name": "key", "Type": "TypeString", "Required": true},
				"value": {"Name": "value", "Type": "TypeString", "Required": true},
			},
		},
	}
	if !IsBreakingChange(oldAttrs, newAttrs) {
		t.Error("Expected breaking change when a Required child is added to an existing Optional parent")
	} else {
		fmt.Println("✓ Test passed: new Required child under existing Optional parent is breaking")
	}
}

// TestNestedRequiredChildUnderNewOptionalParent (path 3): a brand-new Optional
// parent whose nested children are Required is NOT a global breaking change —
// the whole block is new and Optional. This is the ecs_invocation resource_tag
// case that the flattening bug previously misreported.
func TestUnitNestedRequiredChildUnderNewOptionalParent(t *testing.T) {
	oldAttrs := map[string]map[string]interface{}{
		"instance_id": {"Name": "instance_id", "Type": "TypeList", "Required": true},
	}
	newAttrs := map[string]map[string]interface{}{
		"instance_id": {"Name": "instance_id", "Type": "TypeList", "Optional": true},
		"resource_tag": {
			"Name":     "resource_tag",
			"Type":     "TypeSet",
			"Optional": true,
			"Nested": map[string]map[string]interface{}{
				"key":   {"Name": "key", "Type": "TypeString", "Required": true},
				"value": {"Name": "value", "Type": "TypeString", "Required": true},
			},
		},
	}
	if IsBreakingChange(oldAttrs, newAttrs) {
		t.Error("A brand-new Optional parent with nested Required children must not be flagged as breaking")
	} else {
		fmt.Println("✓ Test passed: new Optional parent with nested Required children is safe")
	}
}

// TestParseSchemaNestedNotFlattened verifies ParseSchemaFromAST records nested
// Elem children under their parent field's "Nested" key instead of flattening
// them into the top-level attribute map, and that flat Elem contributes no
// nested attributes.
func TestUnitParseSchemaNestedNotFlattened(t *testing.T) {
	src := `package main
import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
func res() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"resource_tag": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key":   {Type: schema.TypeString, Required: true, ForceNew: true},
						"value": {Type: schema.TypeString, Required: true, ForceNew: true},
					},
				},
			},
			"instance_id": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}`
	attrs := ParseSchemaFromAST(src)
	if _, top := attrs["key"]; top {
		t.Error("nested child 'key' must not be flattened to top level")
	}
	if _, top := attrs["value"]; top {
		t.Error("nested child 'value' must not be flattened to top level")
	}
	rt, ok := attrs["resource_tag"]
	if !ok {
		t.Fatal("resource_tag not parsed at top level")
	}
	nested, ok := rt["Nested"].(map[string]map[string]interface{})
	if !ok || len(nested) != 2 {
		t.Fatalf("expected 2 nested children under resource_tag, got %v", rt["Nested"])
	}
	if c, ok := nested["key"]; !ok || c["Required"] != true {
		t.Errorf("nested key Required not captured: %v", nested["key"])
	}
	if iid, ok := attrs["instance_id"]; ok {
		if _, has := iid["Nested"]; has {
			t.Error("flat Elem (schema.Schema) must not produce a Nested map")
		}
	} else {
		t.Error("instance_id not parsed")
	}
}

// TestEcsInvocationResourceTagNotBreaking is the end-to-end regression for the
// real ecs_invocation change: master had instance_id (Required, flat Elem) and
// no resource_tag; the candidate adds a brand-new Optional resource_tag with
// nested Required key/value plus other new Optional fields, and relaxes
// instance_id to Optional. This must not be reported as a breaking change.
func TestUnitEcsInvocationResourceTagNotBreaking(t *testing.T) {
	oldSrc := `package main
import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
func res() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"command_id": {Type: schema.TypeString, Required: true, ForceNew: true},
			"instance_id": {Type: schema.TypeList, Required: true, Elem: &schema.Schema{Type: schema.TypeString}, ForceNew: true},
			"repeat_mode": {Type: schema.TypeString, Optional: true, ForceNew: true, Computed: true},
			"status": {Type: schema.TypeString, Computed: true},
		},
	}
}`
	newSrc := `package main
import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
func res() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"command_id": {Type: schema.TypeString, Required: true, ForceNew: true},
			"instance_id": {Type: schema.TypeList, Optional: true, Computed: true, Elem: &schema.Schema{Type: schema.TypeString}, ForceNew: true},
			"repeat_mode": {Type: schema.TypeString, Optional: true, ForceNew: true, Computed: true},
			"resource_tag": {
				Type: schema.TypeSet, Optional: true, ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key":   {Type: schema.TypeString, Required: true, ForceNew: true},
						"value": {Type: schema.TypeString, Required: true, ForceNew: true},
					},
				},
			},
			"resource_group_id": {Type: schema.TypeString, Optional: true, ForceNew: true},
			"oss_output_delivery": {Type: schema.TypeString, Optional: true, ForceNew: true},
			"working_dir": {Type: schema.TypeString, Optional: true, Computed: true, ForceNew: true},
			"status": {Type: schema.TypeString, Computed: true},
		},
	}
}`
	oldAttrs := ParseSchemaFromAST(oldSrc)
	newAttrs := ParseSchemaFromAST(newSrc)
	if IsBreakingChange(oldAttrs, newAttrs) {
		t.Error("ecs_invocation: adding a new Optional resource_tag with nested Required children must not be breaking")
	} else {
		fmt.Println("✓ Test passed: ecs_invocation resource_tag addition is non-breaking")
	}
}

// schemaThreeLevel renders a three-level nested schema: top field "rules"
// (TypeList, Optional, Elem &schema.Resource) whose child "match" (TypeList,
// Optional, Elem &schema.Resource) has the given depth-2 grandchildren. It is
// consumed by ParseSchemaFromAST, exercising the real AST pipeline rather than
// hand-crafted attribute maps. Changing a grandchild under the existing
// rules->match parent chain is the case the previous single-level comparison
// silently missed.
func schemaThreeLevel(innerFields string) string {
	return `package main
import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
func res() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"rules": {
				Type: schema.TypeList, Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"match": {
							Type: schema.TypeList, Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
` + innerFields + `
								},
							},
						},
					},
				},
			},
		},
	}
}`
}

// nestedGrandchild walks attrs -> "rules"[Nested] -> "match"[Nested] and returns
// the depth-2 grandchild map plus ok. It documents that the grandchild level is
// actually parsed (otherwise the deep regression below would be vacuous).
func nestedGrandchild(attrs map[string]map[string]interface{}) (map[string]map[string]interface{}, bool) {
	rules, ok := attrs["rules"]
	if !ok {
		return nil, false
	}
	matchLevel, ok := rules["Nested"].(map[string]map[string]interface{})
	if !ok {
		return nil, false
	}
	match, ok := matchLevel["match"]
	if !ok {
		return nil, false
	}
	grand, ok := match["Nested"].(map[string]map[string]interface{})
	return grand, ok
}

// TestDeepNestedRequiredAddIsBreaking: a new Required grandchild added under an
// existing parent chain (rules->match both present in old) must be detected as
// breaking. The previous single-level comparison only inspected rules' direct
// children (match) and never descended into match's Nested, so this was missed.
func TestUnitDeepNestedRequiredAddIsBreaking(t *testing.T) {
	old := schemaThreeLevel(`"field":   {Type: schema.TypeString, Optional: true},`)
	new := schemaThreeLevel(`"field":   {Type: schema.TypeString, Optional: true},
									"pattern": {Type: schema.TypeString, Required: true},`)
	oldAttrs := ParseSchemaFromAST(old)
	newAttrs := ParseSchemaFromAST(new)

	oldGrand, ok := nestedGrandchild(oldAttrs)
	if !ok || len(oldGrand) != 1 {
		t.Fatalf("expected depth-2 grandchild parsed under rules->match, got %v", oldAttrs)
	}
	newGrand, ok := nestedGrandchild(newAttrs)
	if !ok || len(newGrand) != 2 {
		t.Fatalf("expected new depth-2 grandchild 'pattern' parsed, got %v", newAttrs)
	}
	if _, has := newGrand["pattern"]; !has {
		t.Fatal("new Required grandchild 'pattern' not parsed at depth 2")
	}

	if !IsBreakingChange(oldAttrs, newAttrs) {
		t.Error("adding a Required grandchild under an existing parent chain must be breaking")
	} else {
		fmt.Println("✓ Test passed: deep Required grandchild addition detected")
	}
}

// TestDeepNestedDeletionIsBreaking: deleting a grandchild that existed under an
// existing parent chain (rules->match) must be detected as breaking.
func TestUnitDeepNestedDeletionIsBreaking(t *testing.T) {
	old := schemaThreeLevel(`"field":   {Type: schema.TypeString, Optional: true},
								"pattern": {Type: schema.TypeString, Optional: true},`)
	new := schemaThreeLevel(`"field":   {Type: schema.TypeString, Optional: true},`)
	oldAttrs := ParseSchemaFromAST(old)
	newAttrs := ParseSchemaFromAST(new)

	oldGrand, ok := nestedGrandchild(oldAttrs)
	if !ok || len(oldGrand) != 2 {
		t.Fatalf("expected two depth-2 grandchildren in old, got %v", oldAttrs)
	}
	newGrand, ok := nestedGrandchild(newAttrs)
	if !ok || len(newGrand) != 1 {
		t.Fatalf("expected one depth-2 grandchild in new, got %v", newAttrs)
	}
	if _, has := newGrand["pattern"]; has {
		t.Fatal("grandchild 'pattern' should be removed in new schema")
	}

	if !IsBreakingChange(oldAttrs, newAttrs) {
		t.Error("deleting a grandchild under an existing parent chain must be breaking")
	} else {
		fmt.Println("✓ Test passed: deep grandchild deletion detected")
	}
}

// TestDeepNestedTypeChangeIsBreaking: changing a grandchild's type under an
// existing parent chain (rules->match) must be detected as breaking.
func TestUnitDeepNestedTypeChangeIsBreaking(t *testing.T) {
	old := schemaThreeLevel(`"field": {Type: schema.TypeString, Optional: true},`)
	new := schemaThreeLevel(`"field": {Type: schema.TypeInt, Optional: true},`)
	oldAttrs := ParseSchemaFromAST(old)
	newAttrs := ParseSchemaFromAST(new)

	oldGrand, ok := nestedGrandchild(oldAttrs)
	if !ok {
		t.Fatalf("expected depth-2 grandchild parsed, got %v", oldAttrs)
	}
	if oldGrand["field"]["Type"] != "TypeString" {
		t.Fatalf("old grandchild type not captured, got %v", oldGrand["field"])
	}
	newGrand, ok := nestedGrandchild(newAttrs)
	if !ok || newGrand["field"]["Type"] != "TypeInt" {
		t.Fatalf("new grandchild type not captured, got %v", newAttrs)
	}

	if !IsBreakingChange(oldAttrs, newAttrs) {
		t.Error("changing a grandchild type under an existing parent chain must be breaking")
	} else {
		fmt.Println("✓ Test passed: deep grandchild type change detected")
	}
}

// TestNewOptionalMultiLevelSubtreeNotBreaking: a brand-new Optional top-level
// parent whose nested subtree is several levels deep and contains Required
// leaves must NOT be flagged as breaking. The whole block is new and Optional,
// so its internal Required children are legal (path 3) at every depth. This
// guards against a regression that would over-report a new Optional block.
func TestUnitNewOptionalMultiLevelSubtreeNotBreaking(t *testing.T) {
	oldSrc := `package main
import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
func res() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"instance_id": {Type: schema.TypeList, Required: true, Elem: &schema.Schema{Type: schema.TypeString}},
		},
	}
}`
	newSrc := `package main
import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
func res() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"instance_id": {Type: schema.TypeList, Optional: true, Elem: &schema.Schema{Type: schema.TypeString}},
			"policy": {
				Type: schema.TypeList, Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"rule": {
							Type: schema.TypeList, Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"action":   {Type: schema.TypeString, Required: true},
									"resource": {Type: schema.TypeString, Required: true},
								},
							},
						},
					},
				},
			},
		},
	}
}`
	oldAttrs := ParseSchemaFromAST(oldSrc)
	newAttrs := ParseSchemaFromAST(newSrc)

	// Sanity: the new Optional multi-level subtree is actually parsed.
	policy, ok := newAttrs["policy"]
	if !ok {
		t.Fatal("new Optional parent 'policy' not parsed")
	}
	ruleLevel, ok := policy["Nested"].(map[string]map[string]interface{})
	if !ok || len(ruleLevel) != 1 {
		t.Fatalf("expected depth-1 child 'rule' under policy, got %v", policy["Nested"])
	}
	leafLevel, ok := ruleLevel["rule"]["Nested"].(map[string]map[string]interface{})
	if !ok || len(leafLevel) != 2 {
		t.Fatalf("expected depth-2 Required leaves under rule, got %v", ruleLevel["rule"])
	}
	if leafLevel["action"]["Required"] != true || leafLevel["resource"]["Required"] != true {
		t.Fatal("deep Required leaves under the new Optional parent must be parsed")
	}

	if IsBreakingChange(oldAttrs, newAttrs) {
		t.Error("a brand-new Optional multi-level subtree with nested Required leaves must not be breaking")
	} else {
		fmt.Println("✓ Test passed: new Optional multi-level subtree is non-breaking")
	}
}

// TestNestedFullSubtreeDeletionIsBreaking: an existing parent whose nested
// subtree is fully removed (Elem &schema.Resource{} retained but its Schema
// emptied) is a breaking deletion of the removed child attributes. The
// previous len(newNested)==0 guard skipped this case and silently missed it.
func TestUnitNestedFullSubtreeDeletionIsBreaking(t *testing.T) {
	oldSrc := `package main
import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
func res() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"rules": {
				Type: schema.TypeList, Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"field":   {Type: schema.TypeString, Optional: true},
						"pattern": {Type: schema.TypeString, Optional: true},
					},
				},
			},
		},
	}
}`
	newSrc := `package main
import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
func res() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"rules": {
				Type: schema.TypeList, Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{},
				},
			},
		},
	}
}`
	oldAttrs := ParseSchemaFromAST(oldSrc)
	newAttrs := ParseSchemaFromAST(newSrc)
	oldNested, _ := oldAttrs["rules"]["Nested"].(map[string]map[string]interface{})
	if len(oldNested) != 2 {
		t.Fatalf("expected 2 children under rules in old, got %v", oldAttrs["rules"]["Nested"])
	}
	if _, has := newAttrs["rules"]["Nested"]; has {
		t.Fatal("rules must have no Nested children in new schema (subtree fully removed)")
	}
	if !IsBreakingChange(oldAttrs, newAttrs) {
		t.Error("fully removing an existing parent's nested subtree must be breaking")
	} else {
		fmt.Println("✓ Test passed: full nested subtree deletion detected")
	}
}

// TestDeepNestedFullSubtreeDeletionIsBreaking: the deep variant — the
// grandchild level under an existing rules->match parent chain is fully
// removed (match retained, its children emptied), so newNested at the
// grandchild level is empty while oldNested is non-empty. This was missed by
// the previous len(newNested)==0 skip. The partial-deletion case (delete one
// leaf, keep another) is already covered by TestDeepNestedDeletionIsBreaking;
// this covers the full-deletion complement.
func TestUnitDeepNestedFullSubtreeDeletionIsBreaking(t *testing.T) {
	old := schemaThreeLevel(`"field":   {Type: schema.TypeString, Optional: true},
								"pattern": {Type: schema.TypeString, Optional: true},`)
	new := schemaThreeLevel(``)
	oldAttrs := ParseSchemaFromAST(old)
	newAttrs := ParseSchemaFromAST(new)
	oldGrand, ok := nestedGrandchild(oldAttrs)
	if !ok || len(oldGrand) != 2 {
		t.Fatalf("expected two depth-2 grandchildren in old, got %v", oldAttrs)
	}
	if newGrand, ok := nestedGrandchild(newAttrs); ok {
		t.Fatalf("expected no depth-2 grandchildren in new (subtree fully removed), got %v", newGrand)
	}
	if !IsBreakingChange(oldAttrs, newAttrs) {
		t.Error("fully removing the grandchild level under an existing parent chain must be breaking")
	} else {
		fmt.Println("✓ Test passed: deep full nested subtree deletion detected")
	}
}

// TestNestedResourceToFlatIsBreaking: an existing parent whose Elem changes
// from &schema.Resource{Schema:{...}} (nested block) to a flat
// &schema.Schema{Type: ...} (scalar list) silently drops the nested child
// attributes. newNested is empty while oldNested is non-empty, so this is a
// breaking full-subtree deletion the previous skip missed.
func TestUnitNestedResourceToFlatIsBreaking(t *testing.T) {
	oldSrc := `package main
import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
func res() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"rules": {
				Type: schema.TypeList, Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"field": {Type: schema.TypeString, Optional: true},
					},
				},
			},
		},
	}
}`
	newSrc := `package main
import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
func res() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"rules": {
				Type: schema.TypeList, Optional: true,
				Elem: &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}`
	oldAttrs := ParseSchemaFromAST(oldSrc)
	newAttrs := ParseSchemaFromAST(newSrc)
	oldNested, _ := oldAttrs["rules"]["Nested"].(map[string]map[string]interface{})
	if len(oldNested) != 1 {
		t.Fatalf("expected 1 child under rules in old, got %v", oldAttrs["rules"]["Nested"])
	}
	if _, has := newAttrs["rules"]["Nested"]; has {
		t.Fatal("flat Elem (schema.Schema) must not contribute a Nested map in new schema")
	}
	if !IsBreakingChange(oldAttrs, newAttrs) {
		t.Error("flattening an existing nested block (Resource->flat Elem) must be breaking")
	} else {
		fmt.Println("✓ Test passed: Resource->flat Elem deletion detected")
	}
}

// TestUnchangedNestedSubtreeNotBreaking guards against a regression from
// removing the len(newNested)==0 skip: a schema compared against itself (an
// existing parent with an unchanged nested subtree, plus an existing flat
// parent whose Nested is empty on both sides) must remain non-breaking. The
// now-unconditional recursion must be a no-op when both sides are empty.
func TestUnitUnchangedNestedSubtreeNotBreaking(t *testing.T) {
	src := `package main
import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
func res() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"rules": {
				Type: schema.TypeList, Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"field": {Type: schema.TypeString, Optional: true},
					},
				},
			},
			"tags": {Type: schema.TypeList, Optional: true, Elem: &schema.Schema{Type: schema.TypeString}},
		},
	}
}`
	attrs := ParseSchemaFromAST(src)
	if IsBreakingChange(attrs, attrs) {
		t.Error("a schema compared against itself must not be flagged as breaking")
	} else {
		fmt.Println("✓ Test passed: unchanged nested subtree is non-breaking")
	}
}
