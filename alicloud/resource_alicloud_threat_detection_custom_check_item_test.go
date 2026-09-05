package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test ThreatDetection CustomCheckItem. >>> Resource test cases.
// Case new resource alicloud_threat_detection_custom_check_item
func TestAccAliCloudThreatDetectionCustomCheckItem_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_threat_detection_custom_check_item.default"
	ra := resourceAttrInit(resourceId, AliCloudThreatDetectionCustomCheckItemMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ThreatDetectionServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeThreatDetectionCustomCheckItem")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccthreatdetectioncustomcheckitem%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudThreatDetectionCustomCheckItemBasicDependence)
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
					"check_show_name":   name,
					"section_ids":       []string{"515"},
					"vendor":            "ALIYUN",
					"instance_type":     "ECS",
					"instance_sub_type": "ECS_INSTANCE",
					"risk_level":        "high",
					"status":            "RELEASE",
					"check_rule":        "test_check_rule",
					"remark":            "test remark",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"check_show_name":   name,
						"section_ids.#":     "1",
						"vendor":            "ALIYUN",
						"instance_type":     "ECS",
						"instance_sub_type": "ECS_INSTANCE",
						"risk_level":        "high",
						"status":            "RELEASE",
						"check_rule":        "test_check_rule",
						"remark":            "test remark",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"check_show_name":   name + "_update",
					"vendor":            "TENCENT",
					"instance_type":     "RDS",
					"instance_sub_type": "RDS_INSTANCE",
					"risk_level":        "medium",
					"status":            "EDIT",
					"check_rule":        "test_check_rule_updated",
					"remark":            "updated remark",
					"description": []map[string]interface{}{
						{"type": "text", "value": "test description"},
					},
					"assist_info": []map[string]interface{}{
						{"type": "text", "value": "test assist info"},
					},
					"solution": []map[string]interface{}{
						{"type": "text", "value": "test solution"},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"check_show_name":     name + "_update",
						"vendor":              "TENCENT",
						"instance_type":       "RDS",
						"instance_sub_type":   "RDS_INSTANCE",
						"risk_level":          "medium",
						"status":              "EDIT",
						"check_rule":          "test_check_rule_updated",
						"remark":              "updated remark",
						"description.#":       "1",
						"description.0.type":  "text",
						"description.0.value": "test description",
						"assist_info.#":       "1",
						"solution.#":          "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"check_show_name":   name + "_update",
					"vendor":            "TENCENT",
					"instance_type":     "RDS",
					"instance_sub_type": "RDS_INSTANCE",
					"risk_level":        "low",
					"status":            "EDIT",
					"check_rule":        "test_check_rule_final",
					"remark":            "final remark",
					"description": []map[string]interface{}{
						{"type": "text", "value": "updated description"},
					},
					"assist_info": []map[string]interface{}{
						{"type": "text", "value": "updated assist info"},
					},
					"solution": []map[string]interface{}{
						{"type": "text", "value": "updated solution"},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"check_show_name":     name + "_update",
						"risk_level":          "low",
						"check_rule":          "test_check_rule_final",
						"remark":              "final remark",
						"description.0.value": "updated description",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"check_show_name":   name + "_update",
					"vendor":            "TENCENT",
					"instance_type":     "RDS",
					"instance_sub_type": "RDS_INSTANCE",
					"risk_level":        "low",
					"status":            "EDIT",
					"check_rule":        "test_check_rule_final",
					"remark":            "final remark",
					"description":       REMOVEKEY,
					"assist_info":       REMOVEKEY,
					"solution":          REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description.#": "0",
						"assist_info.#": "0",
						"solution.#":    "0",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"section_ids"},
			},
		},
	})
}

// lintignore: AT001
func TestAccAliCloudThreatDetectionCustomCheckItem_datasource(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccthreatdetectioncustomcheckitem%d", rand)
	dependence := AliCloudThreatDetectionCustomCheckItemBasicDependence(name)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: dependence + fmt.Sprintf(`
resource "alicloud_threat_detection_custom_check_item" "default" {
  check_show_name   = "%[1]s"
  section_ids       = [515]
  vendor            = "ALIYUN"
  instance_type     = "ECS"
  instance_sub_type = "ECS_INSTANCE"
  risk_level        = "high"
  status            = "RELEASE"
  check_rule        = "test_check_rule"
  remark            = "test remark"
  description {
    type  = "text"
    value = "test description"
  }
  assist_info {
    type  = "text"
    value = "test assist info"
  }
  solution {
    type  = "text"
    value = "test solution"
  }
}

data "alicloud_threat_detection_custom_check_items" "default" {
  check_id     = alicloud_threat_detection_custom_check_item.default.check_id
  current_page = 1
  depends_on   = [alicloud_threat_detection_custom_check_item.default]
}
`, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.alicloud_threat_detection_custom_check_items.default", "items.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_threat_detection_custom_check_items.default", "items.0.check_id"),
				),
			},
		},
	})
}

var AliCloudThreatDetectionCustomCheckItemMap = map[string]string{}

func AliCloudThreatDetectionCustomCheckItemBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Test ThreatDetection CustomCheckItem. <<< Resource test cases.
