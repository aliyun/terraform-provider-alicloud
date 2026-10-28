// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package cr_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	crsvc "github.com/aliyun/terraform-provider-alicloud/alicloud/service/cr"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test Cr ScanRule. >>> Resource test cases, automatically generated.
// Case ScanRule-1_pl 11745
func TestAccAliCloudCrScanRule_basic11745(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cr_scan_rule.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudCrScanRuleMap11745)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return crsvc.NewCrServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCrScanRule")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccr%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudCrScanRuleBasicDependence11745)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
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
					"instance_id": "${alicloud_cr_ee_instance.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"repo_tag_filter_pattern": ".*",
						"scan_scope":              "REPO",
						"trigger_type":            "MANUAL",
						"scan_type":               "VUL",
						"rule_name":               acctest.CHECKSET,
						"namespaces.#":            "1",
						"repo_names.#":            "4",
						"instance_id":             acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rule_name": "aab",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule_name": "aab",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"namespaces": []string{
						"cc"},
					"repo_names": []string{
						"bb", "cc", "dd"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"namespaces.#": "1",
						"repo_names.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"namespaces": []string{
						"aa"},
					"repo_names": []string{
						"bb", "cc"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"namespaces.#": "1",
						"repo_names.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"namespaces": []string{
						"bb"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"namespaces.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"repo_names": []string{
						"bb", "aa"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"repo_names.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"trigger_type": "MANUAL",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"trigger_type": "MANUAL",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"repo_tag_filter_pattern": "cc",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"repo_tag_filter_pattern": "cc",
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
	"scan_rule_id": acctest.CHECKSET,
	"create_time":  acctest.CHECKSET,
}

func AlicloudCrScanRuleBasicDependence11745(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_cr_ee_instance" "default" {
  payment_type   = "Subscription"
  period         = 1
  renew_period   = 1
  renewal_status = "AutoRenewal"
  instance_type  = "Advanced"
  instance_name  = var.name
}
`, name)
}

// Test Cr ScanRule. <<< Resource test cases, automatically generated.
