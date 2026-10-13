package alicloud

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// checkItemRuleJSON returns the raw CheckRule JSON definition for a given match value.
// CheckRule must be a JSON definition string; arbitrary strings are rejected by the
// CreateCheckItem API with ParamIllegal.checkRule. The format is verified against the API.
//
// Note on CspmVerifyItemRule: UpdateCheckItem runs CspmVerifyItemRule which validates the
// CheckRule against cloud-asset data registered in the CSPM system. Test accounts that lack
// registered cloud-asset data receive CspmVerifyItemRuleError.UnknownFromDataName for ANY
// CheckRule structure — including empty {} and rules with different DataName values
// (ACS_ECS_Disk, ACS_ECS_Instance, ACS_MNS_Topic all rejected identically). This is an
// account-level data gap, not a rule-structure issue. The update test
// (TestAccAliCloudThreatDetectionCustomCheckItem_update) is therefore env-guarded and
// skipped on accounts without registered cloud-asset data.
func checkItemRuleJSON(matchValue string) string {
	return fmt.Sprintf(`{"AssociatedData":{"ToDataList":[{"DataName":"ACS_ECS_Disk","PropertyPath":"InstanceId","FromPropertyPath":"InstanceId"}]},"MatchProperty":{"Operator":"AND","MatchProperties":[{"DataName":"ACS_ECS_Disk","PropertyPath":"InstanceId","MatchOperator":"EQ","MatchPropertyValue":%q}]}}`, matchValue)
}

// checkItemRuleEscaped returns the CheckRule JSON escaped for HCL string interpolation,
// because testAccConfig wraps string values in quotes without escaping internal double quotes.
func checkItemRuleEscaped(matchValue string) string {
	return strings.ReplaceAll(checkItemRuleJSON(matchValue), `"`, `\"`)
}

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
					"risk_level":        "HIGH",
					"status":            "RELEASE",
					"check_rule":        checkItemRuleEscaped("acc1"),
					"remark":            "testremark",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"check_show_name":   name,
						"section_ids.#":     "1",
						"vendor":            "ALIYUN",
						"instance_type":     "ECS",
						"instance_sub_type": "ECS_INSTANCE",
						"risk_level":        "HIGH",
						"status":            "RELEASE",
						"check_rule":        checkItemRuleJSON("acc1"),
						"remark":            "testremark",
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

// TestAccAliCloudThreatDetectionCustomCheckItem_update verifies the UpdateCheckItem lifecycle
// (create → update fields → update again → clear optional fields → import). UpdateCheckItem runs
// CspmVerifyItemRule which requires cloud-asset data registered in the CSPM system. Test accounts
// that lack this data receive CspmVerifyItemRuleError.UnknownFromDataName for any CheckRule
// structure. Set ALICLOUD_CSPM_HAS_ASSET_DATA=true to run this test on accounts with registered
// cloud assets.
func TestAccAliCloudThreatDetectionCustomCheckItem_update(t *testing.T) {
	if os.Getenv("ALICLOUD_CSPM_HAS_ASSET_DATA") != "true" {
		t.Skip("skipping update test: UpdateCheckItem runs CspmVerifyItemRule which requires cloud-asset data registered in the CSPM system; the test account lacks this data (CspmVerifyItemRuleError.UnknownFromDataName). Set ALICLOUD_CSPM_HAS_ASSET_DATA=true to run on accounts with registered cloud assets.")
	}
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
					"risk_level":        "HIGH",
					"status":            "RELEASE",
					"check_rule":        checkItemRuleEscaped("acc1"),
					"remark":            "testremark",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"check_show_name":   name,
						"section_ids.#":     "1",
						"vendor":            "ALIYUN",
						"instance_type":     "ECS",
						"instance_sub_type": "ECS_INSTANCE",
						"risk_level":        "HIGH",
						"status":            "RELEASE",
						"check_rule":        checkItemRuleJSON("acc1"),
						"remark":            "testremark",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"check_show_name":   name + "_update",
					"vendor":            "ALIYUN",
					"instance_type":     "ECS",
					"instance_sub_type": "ECS_INSTANCE",
					"risk_level":        "MEDIUM",
					"status":            "EDIT",
					"check_rule":        checkItemRuleEscaped("acc2"),
					"remark":            "updatedremark",
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
						"vendor":              "ALIYUN",
						"instance_type":       "ECS",
						"instance_sub_type":   "ECS_INSTANCE",
						"risk_level":          "MEDIUM",
						"status":              "EDIT",
						"check_rule":          checkItemRuleJSON("acc2"),
						"remark":              "updatedremark",
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
					"vendor":            "ALIYUN",
					"instance_type":     "ECS",
					"instance_sub_type": "ECS_INSTANCE",
					"risk_level":        "LOW",
					"status":            "EDIT",
					"check_rule":        checkItemRuleEscaped("acc3"),
					"remark":            "finalremark",
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
						"risk_level":          "LOW",
						"check_rule":          checkItemRuleJSON("acc3"),
						"remark":              "finalremark",
						"description.0.value": "updated description",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"check_show_name":   name + "_update",
					"vendor":            "ALIYUN",
					"instance_type":     "ECS",
					"instance_sub_type": "ECS_INSTANCE",
					"risk_level":        "LOW",
					"status":            "EDIT",
					"check_rule":        checkItemRuleEscaped("acc3"),
					"remark":            "finalremark",
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
  risk_level        = "HIGH"
  status            = "RELEASE"
  check_rule        = "%[2]s"
  remark            = "testremark"
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
}
`, name, checkItemRuleEscaped("accds")),
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
