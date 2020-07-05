package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/config"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudConfigRule_basic(t *testing.T) {
	var v config.ConfigRule
	resourceId := "alicloud_config_rule.default"
	ra := resourceAttrInit(resourceId, ConfigRuleMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ConfigService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeConfigRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccConfigRule%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, ConfigRuleBasicdependence)
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
					"config_rule_name":                "terraform-test",
					"risk_level":                      "1",
					"scope_compliance_resource_types": []string{"testscope_compliance_resource_types-1", "testscope_compliance_resource_types-1"},
					"source_detail_message_type":      "ConfigurationItemChangeNotification",
					"source_identifier":               "terraform-test",
					"source_owner":                    "ALIYUN",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"config_rule_name":                  "terraform-test",
						"risk_level":                        "1",
						"scope_compliance_resource_types.#": "2",
						"source_detail_message_type":        "ConfigurationItemChangeNotification",
						"source_identifier":                 "terraform-test",
						"source_owner":                      "ALIYUN",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"InputParameters", "ScopeComplianceResourceId", "ScopeComplianceResourceTypes", "SourceDetailMessageType", "SourceMaximumExecutionFrequency"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"config_rule_name": "terraform-test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"config_rule_name": "terraform-test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"input_parameters": "terraform-test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"input_parameters": "terraform-test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"risk_level": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"risk_level": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scope_compliance_resource_id": "terraform-test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scope_compliance_resource_id": "terraform-test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scope_compliance_resource_types": []string{"testscope_compliance_resource_types-1", "testscope_compliance_resource_types-1"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scope_compliance_resource_types.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"source_detail_message_type": "ConfigurationItemChangeNotification",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source_detail_message_type": "ConfigurationItemChangeNotification",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"source_identifier": "terraform-test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source_identifier": "terraform-test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"source_maximum_execution_frequency": "One_Hour",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source_maximum_execution_frequency": "One_Hour",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"source_owner": "ALIYUN",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source_owner": "ALIYUN",
					}),
				),
			},
		},
	})
}

var ConfigRuleMap = map[string]string{}

func ConfigRuleBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
`, name)
}
