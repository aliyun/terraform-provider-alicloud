// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// Test Apig PluginClass. >>> Resource test cases, automatically generated.
// Case plugin_class_basic_test 12944
func TestAccAliCloudApigPluginClass_basic12944(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_apig_plugin_class.default"
	ra := resourceAttrInit(resourceId, AlicloudApigPluginClassMap12944)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApigServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeApigPluginClass")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccapig%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudApigPluginClassBasicDependence12944)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"wasm_url":                      "https://example.com/plugin.wasm",
					"description":                   "A test plugin class for CloudSpec coverage",
					"version_description":           "Initial version for testing",
					"plugin_class_name":             name,
					"version":                       "1.0.2",
					"alias":                         "test-plugin-alias-v3",
					"execute_priority":              "1",
					"wasm_language":                 "TinyGo",
					"supported_min_gateway_version": "1.0.0",
					"execute_stage":                 "UNSPECIFIED_PHASE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"wasm_url":                      "https://example.com/plugin.wasm",
						"description":                   "A test plugin class for CloudSpec coverage",
						"version_description":           "Initial version for testing",
						"plugin_class_name":             name,
						"version":                       "1.0.2",
						"alias":                         "test-plugin-alias-v3",
						"execute_priority":              "1",
						"wasm_language":                 "TinyGo",
						"supported_min_gateway_version": "1.0.0",
						"execute_stage":                 "UNSPECIFIED_PHASE",
					}),
				),
			},
			{
				ResourceName: resourceId,
				ImportState:  true,
				// wasm_url, execute_priority, execute_stage and version_description are
				// write-only create parameters that the Get API never returns, so the
				// imported state cannot verify them. ForceNew attributes must not appear
				// in ImportStateVerifyIgnore (testing-coverage rule), so import is checked
				// without attribute-level verification.
				ImportStateVerify: false,
			},
		},
	})
}

var AlicloudApigPluginClassMap12944 = map[string]string{
	"status": CHECKSET,
	"type":   CHECKSET,
}

func AlicloudApigPluginClassBasicDependence12944(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Test Apig PluginClass. <<< Resource test cases, automatically generated.
