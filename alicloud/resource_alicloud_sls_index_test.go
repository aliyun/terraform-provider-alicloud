// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Sls Index. >>> Resource test cases, automatically generated.
// Case index_terraform 10982
func TestAccAliCloudSlsIndex_basic10982(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_sls_index.default"
	ra := resourceAttrInit(resourceId, AlicloudSlsIndexMap10982)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SlsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSlsIndex")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccsls%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudSlsIndexBasicDependence10982)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-nanjing"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"line": []map[string]interface{}{
						{
							"chn":            "true",
							"case_sensitive": "true",
							"token": []string{
								"a"},
							"exclude_keys": []string{
								"t"},
						},
					},
					"logstore_name": "${alicloud_log_store.default7MW26R.logstore_name}",
					"project_name":  "${var.project_name}",
					"keys":          "{\\\"test\\\":{\\\"caseSensitive\\\":false,\\\"token\\\":[\\\"\\\\n\\\",\\\"\\\\t\\\",\\\",\\\",\\\" \\\",\\\";\\\",\\\"\\\\\\\"\\\",\\\"'\\\",\\\"(\\\",\\\")\\\",\\\"{\\\",\\\"}\\\",\\\"[\\\",\\\"]\\\",\\\"<\\\",\\\">\\\",\\\"?\\\",\\\"/\\\",\\\"#\\\",\\\":\\\"],\\\"type\\\":\\\"text\\\",\\\"doc_value\\\":false,\\\"alias\\\":\\\"\\\",\\\"chn\\\":false}}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"logstore_name": CHECKSET,
						"project_name":  CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"line": []map[string]interface{}{
						{
							"chn":            "false",
							"case_sensitive": "false",
							"token": []string{
								"tt"},
							"include_keys": []string{
								"tt"},
						},
					},
					"max_text_len": "500",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"max_text_len": "500",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudSlsIndexMap10982 = map[string]string{}

func AlicloudSlsIndexBasicDependence10982(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "logstore_name" {
  default = "logstore-test"
}

variable "project_name" {
  default = "project-for-index-terraform-test"
}

resource "alicloud_log_project" "defaultdCM1bA" {
  description = "terrafrom test"
  name        = var.project_name
}

resource "alicloud_log_store" "default7MW26R" {
  hot_ttl          = "7"
  retention_period = "30"
  shard_count      = "2"
  project          = alicloud_log_project.defaultdCM1bA.project_name
  name             = var.logstore_name
}


`, name)
}

// Test Sls Index. <<< Resource test cases, automatically generated.

// TestAccAliCloudSlsIndex_keys_default_value is an acceptance test that verifies
// the SLS Index keys default value diff suppression logic.
//
// Problem scenario:
// 1. User configures keys with only doc_value and type fields
// 2. SLS API automatically adds default values like alias:"", caseSensitive:false
// 3. Read function writes the full response to tfstate
// 4. Next terraform plan detects a diff, causing unnecessary changes
//
// After fix: DiffSuppressFunc uses areKeysJsonEquivalent which ignores
// known API defaults (alias="", caseSensitive=false), so no diff is produced.
func TestAccAliCloudSlsIndex_keys_default_value(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_sls_index.default"
	ra := resourceAttrInit(resourceId, map[string]string{})
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SlsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSlsIndex")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccsls%d", rand)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-shanghai"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			// Step 1: Create index, keys only configures doc_value and type (user real scenario)
			{
				Config: fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_log_project" "default" {
    name        = "tf-testacc-log-project-%s"
    description = "Test project for SLS index keys default value problem"
}

resource "alicloud_log_store" "default" {
    project          = alicloud_log_project.default.project_name
    name             = "tf-testacc-log-store-%s"
    retention_period = "30"
    shard_count      = "2"
}

# SLS Index - keys only configures doc_value and type, this is the trigger point
resource "alicloud_sls_index" "default" {
    project_name  = alicloud_log_project.default.project_name
    logstore_name = alicloud_log_store.default.logstore_name

    line {
        chn            = true
        case_sensitive = false
        token          = ["\n", "\t", ",", " ", ";"]
    }

    # User config: only doc_value and type, no alias, caseSensitive etc default values
    keys = jsonencode({
        "cost_ms" : {
            "doc_value" : true,
            "type" : "long"
        }
    })
}
`, name, name, name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"logstore_name": CHECKSET,
						"project_name":  CHECKSET,
						"keys":          CHECKSET,
					}),
				),
			},
			// Step 2: No config change, expect plan has no diff
			{
				Config: fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_log_project" "default" {
    name        = "tf-testacc-log-project-%s"
    description = "Test project for SLS index keys default value problem"
}

resource "alicloud_log_store" "default" {
    project          = alicloud_log_project.default.project_name
    name             = "tf-testacc-log-store-%s"
    retention_period = "30"
    shard_count      = "2"
}

# SLS Index - config unchanged, expect no diff
resource "alicloud_sls_index" "default" {
    project_name  = alicloud_log_project.default.project_name
    logstore_name = alicloud_log_store.default.logstore_name

    line {
        chn            = true
        case_sensitive = false
        token          = ["\n", "\t", ",", " ", ";"]
    }

    # Config unchanged, after fix DiffSuppressFunc ignores API defaults
    keys = jsonencode({
        "cost_ms" : {
            "doc_value" : true,
            "type" : "long"
        }
    })
}
`, name, name, name),
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

// TestAreKeysJsonEquivalent is a table-driven unit test for areKeysJsonEquivalent.
// It does NOT require cloud resources and covers all edge cases.
func TestAreKeysJsonEquivalent(t *testing.T) {
	tests := []struct {
		name        string
		old         string
		new         string
		expected    bool
		expectError bool
	}{
		// --- API adds known default values ---
		{
			name:     "API adds alias and caseSensitive defaults",
			old:      `{"cost_ms":{"type":"long","doc_value":true,"alias":"","caseSensitive":false}}`,
			new:      `{"cost_ms":{"type":"long","doc_value":true}}`,
			expected: true,
		},
		{
			name:     "API adds only alias default",
			old:      `{"cost_ms":{"type":"long","doc_value":true,"alias":""}}`,
			new:      `{"cost_ms":{"type":"long","doc_value":true}}`,
			expected: true,
		},
		{
			name:     "API adds only caseSensitive default",
			old:      `{"cost_ms":{"type":"text","doc_value":true,"caseSensitive":false}}`,
			new:      `{"cost_ms":{"type":"text","doc_value":true}}`,
			expected: true,
		},
		{
			name:     "API adds chn default for text type field",
			old:      `{"message":{"type":"text","doc_value":true,"chn":false,"caseSensitive":false,"alias":""}}`,
			new:      `{"message":{"type":"text","doc_value":true}}`,
			expected: true,
		},
		{
			name:     "API adds only chn default",
			old:      `{"message":{"type":"text","doc_value":true,"chn":false}}`,
			new:      `{"message":{"type":"text","doc_value":true}}`,
			expected: true,
		},
		{
			name:     "user changes chn from false to true should surface diff",
			old:      `{"message":{"type":"text","doc_value":true,"chn":true}}`,
			new:      `{"message":{"type":"text","doc_value":true}}`,
			expected: false,
		},
		{
			name:     "both configs identical with no API defaults",
			old:      `{"cost_ms":{"type":"long","doc_value":true}}`,
			new:      `{"cost_ms":{"type":"long","doc_value":true}}`,
			expected: true,
		},

		// --- User modifies a field ---
		{
			name:     "user changes type from long to text",
			old:      `{"cost_ms":{"type":"long","doc_value":true}}`,
			new:      `{"cost_ms":{"type":"text","doc_value":true}}`,
			expected: false,
		},
		{
			name:     "user changes doc_value from true to false",
			old:      `{"cost_ms":{"type":"long","doc_value":true}}`,
			new:      `{"cost_ms":{"type":"long","doc_value":false}}`,
			expected: false,
		},
		{
			name:     "user changes alias to non-empty value",
			old:      `{"cost_ms":{"type":"long","doc_value":true,"alias":"old"}}`,
			new:      `{"cost_ms":{"type":"long","doc_value":true,"alias":"new"}}`,
			expected: false,
		},

		// --- User deletes a field ---
		{
			name:     "user deletes alias with non-empty value should surface diff",
			old:      `{"cost_ms":{"type":"long","doc_value":true,"alias":"legacy"}}`,
			new:      `{"cost_ms":{"type":"long","doc_value":true}}`,
			expected: false,
		},
		{
			name:     "user deletes caseSensitive with non-default value should surface diff",
			old:      `{"cost_ms":{"type":"long","doc_value":true,"caseSensitive":true}}`,
			new:      `{"cost_ms":{"type":"long","doc_value":true}}`,
			expected: false,
		},
		{
			name:     "user deletes doc_value should surface diff",
			old:      `{"cost_ms":{"type":"long","doc_value":true}}`,
			new:      `{"cost_ms":{"type":"long"}}`,
			expected: false,
		},
		{
			name:     "user deletes alias with empty value (API default) should be suppressed",
			old:      `{"cost_ms":{"type":"long","doc_value":true,"alias":""}}`,
			new:      `{"cost_ms":{"type":"long","doc_value":true}}`,
			expected: true,
		},

		// --- Add / delete keys ---
		{
			name:     "user adds a new key",
			old:      `{"cost_ms":{"type":"long","doc_value":true}}`,
			new:      `{"cost_ms":{"type":"long","doc_value":true},"duration_ms":{"type":"long","doc_value":true}}`,
			expected: false,
		},
		{
			name:     "user deletes an existing key",
			old:      `{"cost_ms":{"type":"long","doc_value":true},"duration_ms":{"type":"long","doc_value":true}}`,
			new:      `{"cost_ms":{"type":"long","doc_value":true}}`,
			expected: false,
		},

		// --- Multiple keys with mixed API defaults ---
		{
			name:     "multiple keys with API defaults on one key",
			old:      `{"cost_ms":{"type":"long","doc_value":true,"alias":""},"duration_ms":{"type":"long","doc_value":true,"caseSensitive":false,"alias":""}}`,
			new:      `{"cost_ms":{"type":"long","doc_value":true},"duration_ms":{"type":"long","doc_value":true}}`,
			expected: true,
		},
		{
			name:     "multiple keys one has non-default old-only field",
			old:      `{"cost_ms":{"type":"long","doc_value":true,"alias":""},"duration_ms":{"type":"long","doc_value":true,"alias":"custom"}}`,
			new:      `{"cost_ms":{"type":"long","doc_value":true},"duration_ms":{"type":"long","doc_value":true}}`,
			expected: false,
		},

		// --- Edge cases ---
		{
			name:     "both empty strings",
			old:      ``,
			new:      ``,
			expected: true,
		},
		{
			name:     "one empty one not",
			old:      `{"cost_ms":{"type":"long","doc_value":true}}`,
			new:      ``,
			expected: false,
		},
		{
			name:        "invalid JSON in old falls back to compareJsonTemplateAreEquivalent",
			old:         `{invalid}`,
			new:         `{"cost_ms":{"type":"long","doc_value":true}}`,
			expected:    false,
			expectError: true,
		},
		{
			name:        "invalid JSON in new falls back to compareJsonTemplateAreEquivalent",
			old:         `{"cost_ms":{"type":"long","doc_value":true}}`,
			new:         `{invalid}`,
			expected:    false,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := areKeysJsonEquivalent(tt.old, tt.new)
			if tt.expectError {
				if err == nil {
					t.Fatalf("expected an error but got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("areKeysJsonEquivalent(%q, %q) = %v, want %v", tt.old, tt.new, result, tt.expected)
			}
		})
	}
}
