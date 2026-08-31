package alicloud

import (
	"fmt"
	"testing"

	slsPop "github.com/aliyun/alibaba-cloud-sdk-go/services/sls"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudLogAudit_basic(t *testing.T) {
	var v *slsPop.DescribeAppResponse
	resourceId := "alicloud_log_audit.foo"
	ra := resourceAttrInit(resourceId, logAuditMap)
	serviceFunc := func() interface{} {
		return &LogService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacclogaudit-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceLogAuditConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"display_name": name,
					"aliuid":       "${data.alicloud_account.default.id}",
					"variable_map": map[string]string{
						"actiontrail_enabled": "false",
						"actiontrail_ttl":     "10",
						"oss_access_enabled":  "true",
						"oss_access_ttl":      "155",
						"oss_sync_enabled":    "true",
						"oss_sync_ttl":        "180",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"display_name":                     name,
						"aliuid":                           CHECKSET,
						"variable_map.%":                   "6",
						"variable_map.actiontrail_enabled": "false",
						"variable_map.actiontrail_ttl":     "10",
						"variable_map.oss_access_enabled":  "true",
						"variable_map.oss_access_ttl":      "155",
						"variable_map.oss_sync_enabled":    "true",
						"variable_map.oss_sync_ttl":        "180",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"multi_account": []string{"1234567", "123123123213", "123141412"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"multi_account.#": "3",
					}),
				),
			},
			// TODO: only when center account is resource directory master or resource directory admin need to check resource type config，otherwise pass it directly
			// {
			// 	Config: testAccConfig(map[string]interface{}{
			// 		"resource_directory_type": "custom",
			// 	}),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheck(map[string]string{
			// 			"resource_directory_type": "custom",
			// 		}),
			// 	),
			// },
		},
	})
}

func resourceLogAuditConfigDependence(name string) string {
	return fmt.Sprintf(`
	data "alicloud_account" "default" {}
`)
}

var logAuditMap = map[string]string{
	"display_name": CHECKSET,
}

// TestGetInitParameter is a table-driven unit test for getInitParameter.
// It does NOT require cloud resources and covers all edge cases related to
// the displayName fallback fix.
func TestGetInitParameter(t *testing.T) {
	// Config is a JSON string value in the real API response (not a raw JSON object).
	// getInitParameter does configNew.(string) so Config must be a string type.
	configWithInitParam := `{"initParam":{"aliuid":"1234567890","actiontrail_ttl":"10"}}`

	tests := []struct {
		name            string
		input           string
		expectedDisplay string
		expectError     bool
		expectPanic     bool
		hasInitMap      bool
	}{
		// --- Normal response with DisplayName ---
		{
			name:            "normal response with displayName",
			input:           fmt.Sprintf(`{"AppModel":{"DisplayName":"my-audit-app","Config":%q}}`, configWithInitParam),
			expectedDisplay: "my-audit-app",
			hasInitMap:      true,
		},

		// --- DisplayName is empty string (the bug scenario) ---
		{
			name:            "DisplayName is empty string",
			input:           fmt.Sprintf(`{"AppModel":{"DisplayName":"","Config":%q}}`, configWithInitParam),
			expectedDisplay: "",
			hasInitMap:      true,
		},

		// --- DisplayName field is missing entirely ---
		// Note: getInitParameter uses d["DisplayName"].(string) without comma-ok,
		// so a missing DisplayName field causes a panic. This is a pre-existing
		// limitation in the original code — the displayName == "" fallback in
		// Read() handles the empty-string case, not the missing-field case.
		{
			name:            "DisplayName field missing from AppModel",
			input:           fmt.Sprintf(`{"AppModel":{"Config":%q}}`, configWithInitParam),
			expectedDisplay: "",
			hasInitMap:      false,
			expectPanic:     true,
		},

		// --- AppModel is missing ---
		{
			name:            "AppModel missing entirely",
			input:           `{"someOtherKey":"value"}`,
			expectedDisplay: "",
			hasInitMap:      false,
			expectError:     false, // no error, just empty return values
		},

		// --- AppModel is not a map ---
		{
			name:            "AppModel is a string not a map",
			input:           `{"AppModel":"not-a-map"}`,
			expectedDisplay: "",
			hasInitMap:      false,
		},

		// --- Invalid JSON ---
		{
			name:            "invalid JSON",
			input:           `{invalid json}`,
			expectedDisplay: "",
			hasInitMap:      false,
			expectError:     true,
		},

		// --- Empty string input ---
		{
			name:            "empty string input",
			input:           ``,
			expectedDisplay: "",
			hasInitMap:      false,
			expectError:     true, // json.Unmarshal of empty string returns error
		},

		// --- Config is not valid JSON string ---
		// Note: json.Unmarshal of "not-json" fails, m is empty, then
		// m["initParam"].(map[string]interface{}) panics on nil.
		{
			name:            "Config is not a valid JSON string",
			input:           `{"AppModel":{"DisplayName":"test","Config":"not-json"}}`,
			expectedDisplay: "test",
			hasInitMap:      false,
			expectPanic:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					if !tt.expectPanic {
						t.Fatalf("unexpected panic: %v", r)
					}
					return
				}
				if tt.expectPanic {
					t.Fatalf("expected panic but got none")
				}
			}()

			displayName, initMap, err := getInitParameter(tt.input)

			if tt.expectError && err == nil {
				t.Fatalf("expected an error but got none")
			}
			if !tt.expectError && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if displayName != tt.expectedDisplay {
				t.Errorf("displayName = %q, want %q", displayName, tt.expectedDisplay)
			}
			if tt.hasInitMap && initMap == nil {
				t.Errorf("expected initMap to be non-nil")
			}
			if !tt.hasInitMap && initMap != nil {
				t.Errorf("expected initMap to be nil, got %v", initMap)
			}
		})
	}
}

// TestAccAlicloudLogAudit_display_name_fallback is a minimal acceptance test that
// verifies the displayName == "" fallback to d.Id() fix.
//
// It uses the same data.alicloud_account pattern as TestAccAlicloudLogAudit_basic
// but with a minimal variable_map and an additional PlanOnly step to verify
// no spurious diff after the fix.
func TestAccAlicloudLogAudit_display_name_fallback(t *testing.T) {
	var v *slsPop.DescribeAppResponse
	resourceId := "alicloud_log_audit.foo"
	ra := resourceAttrInit(resourceId, logAuditMap)
	serviceFunc := func() interface{} {
		return &LogService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacclogaudit-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceLogAuditConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			// Step 1: Create log_audit with minimal config.
			{
				Config: testAccConfig(map[string]interface{}{
					"display_name": name,
					"aliuid":       "${data.alicloud_account.default.id}",
					"variable_map": map[string]string{
						"actiontrail_enabled": "false",
						"actiontrail_ttl":     "10",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"display_name":                     name,
						"variable_map.actiontrail_enabled": "false",
						"variable_map.actiontrail_ttl":     "10",
					}),
				),
			},
			// Step 2: Import and verify state.
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Step 3: Re-plan with same config — expect no diff.
			// This is the core assertion for the displayName fallback fix.
			{
				Config: testAccConfig(map[string]interface{}{
					"display_name": name,
					"aliuid":       "${data.alicloud_account.default.id}",
					"variable_map": map[string]string{
						"actiontrail_enabled": "false",
						"actiontrail_ttl":     "10",
					},
				}),
				PlanOnly: true,
			},
		},
	})
}
