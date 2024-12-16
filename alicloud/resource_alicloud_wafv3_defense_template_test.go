package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Wafv3 DefenseTemplate. >>> Resource test cases, automatically generated.
// Case 接入terraform 5993
func TestAccAliCloudWafv3DefenseTemplate_basic5993(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.WAFV3SupportRegions)
	resourceId := "alicloud_wafv3_defense_template.default"
	ra := resourceAttrInit(resourceId, AlicloudWafv3DefenseTemplateMap5993)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &Wafv3ServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeWafv3DefenseTemplate")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%swafv3defensetemplate%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudWafv3DefenseTemplateBasicDependence5993)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckForCleanUpInstances(t, string(connectivity.Hangzhou), "waf", "waf", "waf", "waf")
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"status":                             "0",
					"instance_id":                        "${alicloud_wafv3_instance.default.id}",
					"defense_template_name":              name,
					"template_type":                      "user_custom",
					"template_origin":                    "custom",
					"defense_scene":                      "antiscan",
					"resource_manager_resource_group_id": "test",
					"description":                        "update_template",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status":                             "0",
						"instance_id":                        CHECKSET,
						"defense_template_name":              name,
						"template_type":                      "user_custom",
						"template_origin":                    "custom",
						"defense_scene":                      "antiscan",
						"resource_manager_resource_group_id": "test",
						"description":                        "update_template",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "createTestDescription",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "createTestDescription",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"defense_template_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"defense_template_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "update_template",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "update_template",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status":                             "1",
					"instance_id":                        "${alicloud_wafv3_instance.default.id}",
					"defense_template_name":              name + "_update",
					"template_type":                      "user_custom",
					"template_origin":                    "custom",
					"defense_scene":                      "antiscan",
					"resource_manager_resource_group_id": "test",
					"description":                        "createTestDescription",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status":                             "1",
						"instance_id":                        CHECKSET,
						"defense_template_name":              name + "_update",
						"template_type":                      "user_custom",
						"template_origin":                    "custom",
						"defense_scene":                      "antiscan",
						"resource_manager_resource_group_id": "test",
						"description":                        "createTestDescription",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"resource_manager_resource_group_id"},
			},
		},
	})
}

var AlicloudWafv3DefenseTemplateMap5993 = map[string]string{
	"defense_template_id": CHECKSET,
}

func AlicloudWafv3DefenseTemplateBasicDependence5993(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_wafv3_instance" "default" {
}


`, name)
}

// Case 接入terraform 5993  twin
func TestAccAliCloudWafv3DefenseTemplate_basic5993_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_wafv3_defense_template.default"
	ra := resourceAttrInit(resourceId, AlicloudWafv3DefenseTemplateMap5993)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &Wafv3ServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeWafv3DefenseTemplate")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%swafv3defensetemplate%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudWafv3DefenseTemplateBasicDependence5993)
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
					"status":                             "0",
					"instance_id":                        "${alicloud_wafv3_instance.default.id}",
					"defense_template_name":              name,
					"template_type":                      "user_custom",
					"template_origin":                    "custom",
					"defense_scene":                      "antiscan",
					"resource_manager_resource_group_id": "test",
					"description":                        "update_template",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status":                             "0",
						"instance_id":                        CHECKSET,
						"defense_template_name":              name,
						"template_type":                      "user_custom",
						"template_origin":                    "custom",
						"defense_scene":                      "antiscan",
						"resource_manager_resource_group_id": "test",
						"description":                        "update_template",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"resource_manager_resource_group_id"},
			},
		},
	})
}

// Test Wafv3 DefenseTemplate. <<< Resource test cases, automatically generated.
