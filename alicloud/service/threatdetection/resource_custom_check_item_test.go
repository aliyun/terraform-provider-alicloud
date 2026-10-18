package threatdetection_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// checkItemRuleJSON returns the raw CheckRule JSON definition for a given match
// value. CheckRule must be a JSON definition string; arbitrary strings are
// rejected by the CreateCheckItem API with ParamIllegal.checkRule.
//
// UpdateCheckItem runs CspmVerifyItemRule which validates the CheckRule against
// cloud-asset data registered in the CSPM system. Test accounts that lack
// registered cloud-asset data receive CspmVerifyItemRuleError.UnknownFromDataName
// for any CheckRule structure. The update test is therefore env-guarded.
func checkItemRuleJSON(matchValue string) string {
	return fmt.Sprintf(`{"AssociatedData":{"ToDataList":[{"DataName":"ACS_ECS_Disk","PropertyPath":"InstanceId","FromPropertyPath":"InstanceId"}]},"MatchProperty":{"Operator":"AND","MatchProperties":[{"DataName":"ACS_ECS_Disk","PropertyPath":"InstanceId","MatchOperator":"EQ","MatchPropertyValue":%q}]}}`, matchValue)
}

// checkItemRuleEscaped returns the CheckRule JSON escaped for HCL string
// interpolation, because the test config wraps the value in double quotes.
func checkItemRuleEscaped(matchValue string) string {
	return strings.ReplaceAll(checkItemRuleJSON(matchValue), `"`, `\"`)
}

func TestAccAliCloudThreatDetectionCustomCheckItem_basic(t *testing.T) {
	r := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccthreatdetectioncustomcheckitem%d", r)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccThreatDetectionCustomCheckItemConfig(name, checkItemRuleEscaped("acc1"), "HIGH", "RELEASE", "testremark", nil),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("alicloud_threat_detection_custom_check_item.default", "id"),
					resource.TestCheckResourceAttrSet("alicloud_threat_detection_custom_check_item.default", "check_id"),
					resource.TestCheckResourceAttr("alicloud_threat_detection_custom_check_item.default", "check_show_name", name),
					resource.TestCheckResourceAttr("alicloud_threat_detection_custom_check_item.default", "section_ids.#", "1"),
					resource.TestCheckResourceAttr("alicloud_threat_detection_custom_check_item.default", "section_ids.0", "515"),
					resource.TestCheckResourceAttr("alicloud_threat_detection_custom_check_item.default", "vendor", "ALIYUN"),
					resource.TestCheckResourceAttr("alicloud_threat_detection_custom_check_item.default", "instance_type", "ECS"),
					resource.TestCheckResourceAttr("alicloud_threat_detection_custom_check_item.default", "instance_sub_type", "ECS_INSTANCE"),
					resource.TestCheckResourceAttr("alicloud_threat_detection_custom_check_item.default", "risk_level", "HIGH"),
					resource.TestCheckResourceAttr("alicloud_threat_detection_custom_check_item.default", "status", "RELEASE"),
					resource.TestCheckResourceAttr("alicloud_threat_detection_custom_check_item.default", "check_rule", checkItemRuleJSON("acc1")),
					resource.TestCheckResourceAttr("alicloud_threat_detection_custom_check_item.default", "remark", "testremark"),
				),
			},
			{
				ResourceName:            "alicloud_threat_detection_custom_check_item.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"section_ids"},
			},
		},
	})
}

// TestAccAliCloudThreatDetectionCustomCheckItem_update verifies the UpdateCheckItem
// lifecycle. UpdateCheckItem runs CspmVerifyItemRule which requires cloud-asset data
// registered in the CSPM system. Set ALICLOUD_CSPM_HAS_ASSET_DATA=true to run this
// test on accounts with registered cloud assets.
//
// lintignore: AT001
func TestAccAliCloudThreatDetectionCustomCheckItem_update(t *testing.T) {
	if os.Getenv("ALICLOUD_CSPM_HAS_ASSET_DATA") != "true" {
		t.Skip("skipping update test: UpdateCheckItem runs CspmVerifyItemRule which requires cloud-asset data registered in the CSPM system; the test account lacks this data (CspmVerifyItemRuleError.UnknownFromDataName). Set ALICLOUD_CSPM_HAS_ASSET_DATA=true to run on accounts with registered cloud assets.")
	}
	r := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccthreatdetectioncustomcheckitem%d", r)
	blocks := map[string]string{"type": "text", "value": "test description"}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccThreatDetectionCustomCheckItemConfig(name, checkItemRuleEscaped("acc1"), "HIGH", "RELEASE", "testremark", nil),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_threat_detection_custom_check_item.default", "check_show_name", name),
					resource.TestCheckResourceAttr("alicloud_threat_detection_custom_check_item.default", "section_ids.#", "1"),
					resource.TestCheckResourceAttr("alicloud_threat_detection_custom_check_item.default", "risk_level", "HIGH"),
					resource.TestCheckResourceAttr("alicloud_threat_detection_custom_check_item.default", "check_rule", checkItemRuleJSON("acc1")),
					resource.TestCheckResourceAttr("alicloud_threat_detection_custom_check_item.default", "remark", "testremark"),
				),
			},
			{
				Config: testAccThreatDetectionCustomCheckItemConfig(name+"_update", checkItemRuleEscaped("acc2"), "MEDIUM", "EDIT", "updatedremark", blocks),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_threat_detection_custom_check_item.default", "check_show_name", name+"_update"),
					resource.TestCheckResourceAttr("alicloud_threat_detection_custom_check_item.default", "risk_level", "MEDIUM"),
					resource.TestCheckResourceAttr("alicloud_threat_detection_custom_check_item.default", "status", "EDIT"),
					resource.TestCheckResourceAttr("alicloud_threat_detection_custom_check_item.default", "check_rule", checkItemRuleJSON("acc2")),
					resource.TestCheckResourceAttr("alicloud_threat_detection_custom_check_item.default", "remark", "updatedremark"),
					resource.TestCheckResourceAttr("alicloud_threat_detection_custom_check_item.default", "description.#", "1"),
					resource.TestCheckResourceAttr("alicloud_threat_detection_custom_check_item.default", "description.0.type", "text"),
					resource.TestCheckResourceAttr("alicloud_threat_detection_custom_check_item.default", "description.0.value", "test description"),
					resource.TestCheckResourceAttr("alicloud_threat_detection_custom_check_item.default", "assist_info.#", "1"),
					resource.TestCheckResourceAttr("alicloud_threat_detection_custom_check_item.default", "solution.#", "1"),
				),
			},
			{
				Config: testAccThreatDetectionCustomCheckItemConfig(name+"_update", checkItemRuleEscaped("acc3"), "LOW", "EDIT", "finalremark", map[string]string{"type": "text", "value": "updated description"}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_threat_detection_custom_check_item.default", "risk_level", "LOW"),
					resource.TestCheckResourceAttr("alicloud_threat_detection_custom_check_item.default", "check_rule", checkItemRuleJSON("acc3")),
					resource.TestCheckResourceAttr("alicloud_threat_detection_custom_check_item.default", "remark", "finalremark"),
					resource.TestCheckResourceAttr("alicloud_threat_detection_custom_check_item.default", "description.0.value", "updated description"),
				),
			},
			{
				Config: testAccThreatDetectionCustomCheckItemConfig(name+"_update", checkItemRuleEscaped("acc3"), "LOW", "EDIT", "finalremark", nil),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_threat_detection_custom_check_item.default", "description.#", "0"),
					resource.TestCheckResourceAttr("alicloud_threat_detection_custom_check_item.default", "assist_info.#", "0"),
					resource.TestCheckResourceAttr("alicloud_threat_detection_custom_check_item.default", "solution.#", "0"),
				),
			},
			{
				ResourceName:            "alicloud_threat_detection_custom_check_item.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"section_ids"},
			},
		},
	})
}

// TestAccAliCloudThreatDetectionCustomCheckItem_datasource verifies the data
// source reads back the created check item by check_id.
//
// lintignore: AT001
func TestAccAliCloudThreatDetectionCustomCheckItem_datasource(t *testing.T) {
	r := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccthreatdetectioncustomcheckitem%d", r)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccThreatDetectionCustomCheckItemDatasourceConfig(name, checkItemRuleEscaped("accds")),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.alicloud_threat_detection_custom_check_items.default", "items.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_threat_detection_custom_check_items.default", "items.0.check_id"),
				),
			},
		},
	})
}

// testAccThreatDetectionCustomCheckItemConfig renders the resource config. When
// block is nil the description / assist_info / solution blocks are omitted.
func testAccThreatDetectionCustomCheckItemConfig(name, rule, riskLevel, status, remark string, block map[string]string) string {
	blocks := ""
	if block != nil {
		blocks = fmt.Sprintf(`
  description {
    type  = "%[1]s"
    value = "%[2]s"
  }
  assist_info {
    type  = "%[1]s"
    value = "test assist info"
  }
  solution {
    type  = "%[1]s"
    value = "test solution"
  }
`, block["type"], block["value"])
	}
	return fmt.Sprintf(`
resource "alicloud_threat_detection_custom_check_item" "default" {
  check_show_name   = "%[1]s"
  section_ids       = [515]
  vendor            = "ALIYUN"
  instance_type     = "ECS"
  instance_sub_type = "ECS_INSTANCE"
  risk_level        = "%[2]s"
  status            = "%[3]s"
  check_rule        = "%[4]s"
  remark            = "%[5]s"%[6]s
}
`, name, riskLevel, status, rule, remark, blocks)
}

func testAccThreatDetectionCustomCheckItemDatasourceConfig(name, rule string) string {
	return fmt.Sprintf(`
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
`, name, rule)
}
