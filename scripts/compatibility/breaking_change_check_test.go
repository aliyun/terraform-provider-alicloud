//nolint:all
package main

import (
	"fmt"
	"strings"
	"testing"

	"github.com/waigani/diffparser"
)

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
			// Parse the diff
			diff, err := diffparser.Parse(tt.diffContent)
			if err != nil {
				t.Fatalf("Failed to parse diff: %v", err)
			}

			if len(diff.Files) == 0 {
				t.Fatal("No files in diff")
			}

			// Process the hunks
			oldRetryCodes := make(map[string]map[string]struct{})
			newRetryCodes := make(map[string]map[string]struct{})

			for _, hunk := range diff.Files[0].Hunks {
				if hunk != nil {
					ParseRetryErrorCodesFromHunk(hunk, oldRetryCodes, newRetryCodes)
				}
			}

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
			name:    "standard_expected_errors_call",
			content: `if IsExpectedErrors(err, []string{"Throttling", "SystemBusy"}) {`,
			expected: map[string][]string{
				"line_1": {"Throttling", "SystemBusy"},
			},
		},
		{
			name:    "expected_errors_call_with_spaces",
			content: `if IsExpectedErrors(err,  []string{ "Throttling" , "SystemBusy" }) {`,
			expected: map[string][]string{
				"line_1": {"Throttling", "SystemBusy"},
			},
		},
		{
			name:    "single_error_code",
			content: `if IsExpectedErrors(err, []string{"InvalidInstanceId.NotFound"}) {`,
			expected: map[string][]string{
				"line_1": {"InvalidInstanceId.NotFound"},
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
			name: "function_context_without_action_definition",
			content: `func resourceAliCloudEcsInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	
	var request map[string]interface{}
	var response map[string]interface{}
	var err error
	
	wait := incrementalWait(1*time.Second, 1*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Ecs", "2014-05-26", "RunInstances", nil, request, false)
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
				"resourceAliCloudEcsInstanceCreate": {"IncorrectVSwitchStatus"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock hunk with the content
			lines := strings.Split(tt.content, "\n")
			var diffLines []*diffparser.DiffLine

			for i, line := range lines {
				diffLines = append(diffLines, &diffparser.DiffLine{
					Mode:    diffparser.UNCHANGED,
					Number:  i + 1,
					Content: line,
				})
			}

			hunk := &diffparser.DiffHunk{
				WholeRange: diffparser.DiffRange{
					Lines: diffLines,
				},
			}

			oldRetryCodes := make(map[string]map[string]struct{})
			newRetryCodes := make(map[string]map[string]struct{})

			ParseRetryErrorCodesFromHunk(hunk, oldRetryCodes, newRetryCodes)

			// Since it's UNCHANGED, both old and new should have the same codes
			if !verifyExpectedCodes(t, "parsed", newRetryCodes, tt.expected) {
				t.Errorf("Failed to parse expected codes from content: %s", tt.content)
			}
		})
	}
}
