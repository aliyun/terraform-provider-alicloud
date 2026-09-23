package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Config Remediation. >>> Resource test cases, automatically generated.
// Case 2979
func TestAccAliCloudConfigRemediation_basic2979(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_config_remediation.default"
	ra := resourceAttrInit(resourceId, AlicloudConfigRemediationMap2979)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ConfigServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeConfigRemediation")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sConfigRemediation%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudConfigRemediationBasicDependence2979)
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
					"config_rule_id":          "${alicloud_config_rule.prerequirement-rule.config_rule_id}",
					"remediation_template_id": "ACS-TAG-TagResources",
					"invoke_type":             "MANUAL_EXECUTION",
					"params":                  "{\\\"regionId\\\":\\\"{regionId}\\\",\\\"tags\\\":\\\"{\\\\\\\"terraform\\\\\\\":\\\\\\\"terraform\\\\\\\"}\\\",\\\"resourceType\\\":\\\"{resourceType}\\\",\\\"resourceIds\\\":\\\"{resourceId}\\\"}",
					"remediation_type":        "OOS",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"config_rule_id":          CHECKSET,
						"remediation_template_id": "ACS-TAG-TagResources",
						"invoke_type":             "MANUAL_EXECUTION",
						"params":                  "{\"regionId\":\"{regionId}\",\"tags\":\"{\\\"terraform\\\":\\\"terraform\\\"}\",\"resourceType\":\"{resourceType}\",\"resourceIds\":\"{resourceId}\"}",
						"remediation_type":        "OOS",
					}),
				),
			}, {
				Config: testAccConfig(map[string]interface{}{
					"remediation_template_id": "ACS-TAG-TagResourcesIgnoreCaseSensitive",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"remediation_template_id": "ACS-TAG-TagResourcesIgnoreCaseSensitive",
					}),
				),
			}, {
				Config: testAccConfig(map[string]interface{}{
					"invoke_type": "AUTO_EXECUTION",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"invoke_type": "AUTO_EXECUTION",
					}),
				),
			}, {
				Config: testAccConfig(map[string]interface{}{
					"params": "{\\\"regionId\\\":\\\"{regionId}\\\",\\\"tags\\\":\\\"{\\\\\\\"terraform\\\\\\\":\\\\\\\"terraform_update\\\\\\\"}\\\",\\\"resourceType\\\":\\\"{resourceType}\\\",\\\"resourceIds\\\":\\\"{resourceId}\\\"}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"params": "{\"regionId\":\"{regionId}\",\"tags\":\"{\\\"terraform\\\":\\\"terraform_update\\\"}\",\"resourceType\":\"{resourceType}\",\"resourceIds\":\"{resourceId}\"}",
					}),
				),
			}, {
				Config: testAccConfig(map[string]interface{}{
					"config_rule_id":          "${alicloud_config_rule.prerequirement-rule.config_rule_id}",
					"remediation_template_id": "ACS-TAG-TagResources",
					"remediation_source_type": "ALIYUN",
					"invoke_type":             "MANUAL_EXECUTION",
					"params":                  "{\\\"regionId\\\":\\\"{regionId}\\\",\\\"tags\\\":\\\"{\\\\\\\"terraform\\\\\\\":\\\\\\\"terraform\\\\\\\"}\\\",\\\"resourceType\\\":\\\"{resourceType}\\\",\\\"resourceIds\\\":\\\"{resourceId}\\\"}",
					"remediation_type":        "OOS",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"config_rule_id":          CHECKSET,
						"remediation_template_id": "ACS-TAG-TagResources",
						"remediation_source_type": "ALIYUN",
						"invoke_type":             "MANUAL_EXECUTION",
						"params":                  "{\"regionId\":\"{regionId}\",\"tags\":\"{\\\"terraform\\\":\\\"terraform\\\"}\",\"resourceType\":\"{resourceType}\",\"resourceIds\":\"{resourceId}\"}",
						"remediation_type":        "OOS",
					}),
				),
			}, {
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudConfigRemediationMap2979 = map[string]string{}

func AlicloudConfigRemediationBasicDependence2979(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_config_rule" "prerequirement-rule" {
  description                = "关联的资源类型下实体资源均已有指定标签，存在没有指定标签的资源则视为“不合规”。"
  source_owner               = "ALIYUN"
  source_identifier          = "required-tags"
  risk_level                 = 1
  tag_value_scope            = "test"
  tag_key_scope              = "test"
  exclude_resource_ids_scope = "test"
  region_ids_scope           = "cn-hangzhou"
  config_rule_trigger_types  = "ConfigurationItemChangeNotification"
  resource_group_ids_scope   = "rg-acfmvoh45rhcfly"
  resource_types_scope = [
  "ACS::RDS::DBInstance"]
  rule_name = "tf-cicd-rule-by-required-tags"
  input_parameters = {
    tag1Key = "terraform"
    tag1Value = "terraform"
  }
}


`, name)
}

// Test Config Remediation. <<< Resource test cases, automatically generated.
