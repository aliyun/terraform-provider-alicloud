// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Cr ScanRule. >>> Resource test cases, automatically generated.
// Case ScanRule-1_pl 11745
func TestAccAliCloudCrScanRule_basic11745(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cr_scan_rule.default"
	ra := resourceAttrInit(resourceId, AlicloudCrScanRuleMap11745)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CrServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCrScanRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccr%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCrScanRuleBasicDependence11745)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"repo_tag_filter_pattern": ".*",
					"scan_scope":              "REPO",
					"trigger_type":            "MANUAL",
					"scan_type":               "VUL",
					"rule_name":               "302",
					"namespaces": []string{
						"aa"},
					"repo_names": []string{
						"bb", "cc", "dd", "ee"},
					"instance_id": "cri-j3ak232wjg7u4r1p",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"repo_tag_filter_pattern": ".*",
						"scan_scope":              "REPO",
						"trigger_type":            "MANUAL",
						"scan_type":               "VUL",
						"rule_name":               CHECKSET,
						"namespaces.#":            "1",
						"repo_names.#":            "4",
						"instance_id":             CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"repo_tag_filter_pattern": "a",
					"scan_scope":              "NAMESPACE",
					"trigger_type":            "AUTO",
					"rule_name":               "aa",
					"namespaces": []string{
						"aa", "bb", "cc", "dd"},
					"repo_names": []string{},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"repo_tag_filter_pattern": "a",
						"scan_scope":              "NAMESPACE",
						"trigger_type":            "AUTO",
						"rule_name":               "aa",
						"namespaces.#":            "4",
						"repo_names.#":            "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"repo_tag_filter_pattern": "bb",
					"scan_scope":              "REPO",
					"trigger_type":            "MANUAL",
					"rule_name":               "bb",
					"namespaces": []string{
						"aa"},
					"repo_names": []string{
						"aa"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"repo_tag_filter_pattern": "bb",
						"scan_scope":              "REPO",
						"trigger_type":            "MANUAL",
						"rule_name":               "bb",
						"namespaces.#":            "1",
						"repo_names.#":            "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"repo_tag_filter_pattern": "cc",
					"scan_scope":              "INSTANCE",
					"trigger_type":            "AUTO",
					"rule_name":               "dd",
					"namespaces":              []string{},
					"repo_names":              []string{},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"repo_tag_filter_pattern": "cc",
						"scan_scope":              "INSTANCE",
						"trigger_type":            "AUTO",
						"rule_name":               "dd",
						"namespaces.#":            "0",
						"repo_names.#":            "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"repo_tag_filter_pattern": "dd",
					"scan_scope":              "NAMESPACE",
					"trigger_type":            "MANUAL",
					"rule_name":               "ff",
					"namespaces": []string{
						"aa", "bb"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"repo_tag_filter_pattern": "dd",
						"scan_scope":              "NAMESPACE",
						"trigger_type":            "MANUAL",
						"rule_name":               "ff",
						"namespaces.#":            "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"repo_tag_filter_pattern": "aa",
					"scan_scope":              "INSTANCE",
					"trigger_type":            "AUTO",
					"rule_name":               "gg",
					"namespaces":              []string{},
					"repo_names":              []string{},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"repo_tag_filter_pattern": "aa",
						"scan_scope":              "INSTANCE",
						"trigger_type":            "AUTO",
						"rule_name":               "gg",
						"namespaces.#":            "0",
						"repo_names.#":            "0",
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

var AlicloudCrScanRuleMap11745 = map[string]string{
	"scan_rule_id": CHECKSET,
	"create_time":  CHECKSET,
}

func AlicloudCrScanRuleBasicDependence11745(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {}


`, name)
}

// Test Cr ScanRule. <<< Resource test cases, automatically generated.
