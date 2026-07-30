// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
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
					"project_name":  "${var.name}",
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
  default = "logstore-test-%s"
}

variable "project_name" {
  default = "tf-testacc-sls-project-%s"
}

resource "alicloud_log_project" "defaultdCM1bA" {
  description = "terrafrom test"
  name        = var.name
}

resource "alicloud_log_store" "default7MW26R" {
  hot_ttl          = "7"
  retention_period = "30"
  shard_count      = "2"
  project          = alicloud_log_project.defaultdCM1bA.project_name
  name             = var.logstore_name
}


`, name, name, name)
}

// Test Sls Index. <<< Resource test cases, automatically generated.

// TestAccAliCloudSlsIndex_keys_default_value verifies that SLS API defaults in
// an index key do not cause a perpetual diff.
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

resource "alicloud_sls_index" "default" {
    project_name  = alicloud_log_project.default.project_name
    logstore_name = alicloud_log_store.default.logstore_name

    line {
        chn            = true
        case_sensitive = false
        token          = ["\n", "\t", ",", " ", ";"]
    }

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

resource "alicloud_sls_index" "default" {
    project_name  = alicloud_log_project.default.project_name
    logstore_name = alicloud_log_store.default.logstore_name

    line {
        chn            = true
        case_sensitive = false
        token          = ["\n", "\t", ",", " ", ";"]
    }

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

func TestAreKeysJsonEquivalent(t *testing.T) {
	tests := []struct {
		name        string
		old         string
		new         string
		expected    bool
		expectError bool
	}{
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
		{
			name:     "identical compound field values",
			old:      `{"message":{"type":"text","token":[",",";"]}}`,
			new:      `{"message":{"type":"text","token":[",",";"]}}`,
			expected: true,
		},
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
			name:     "user deletes alias with empty value is suppressed",
			old:      `{"cost_ms":{"type":"long","doc_value":true,"alias":""}}`,
			new:      `{"cost_ms":{"type":"long","doc_value":true}}`,
			expected: true,
		},
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
			name:        "invalid JSON in old falls back to original comparison",
			old:         `{invalid}`,
			new:         `{"cost_ms":{"type":"long","doc_value":true}}`,
			expected:    false,
			expectError: true,
		},
		{
			name:        "invalid JSON in new falls back to original comparison",
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
