package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/config"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
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
			testAccPreCheckWithRegions(t, true, connectivity.CloudConfigSupportedRegions)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"rule_name":                       name,
					"risk_level":                      "1",
					"scope_compliance_resource_types": []string{"ACS::ECS::Instance"},
					"source_detail_message_type":      "ConfigurationItemChangeNotification",
					"source_identifier":               "ecs-instances-in-vpc",
					"source_owner":                    "ALIYUN",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule_name":                         name,
						"risk_level":                        "1",
						"scope_compliance_resource_types.#": "1",
						"source_detail_message_type":        "ConfigurationItemChangeNotification",
						"source_identifier":                 "ecs-instances-in-vpc",
						"source_owner":                      "ALIYUN",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       false,
				ImportStateVerifyIgnore: []string{"scope_compliance_resource_id", "scope_compliance_resource_types", "source_detail_message_type", "source_maximum_execution_frequency"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "test_rule",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "test_rule",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"input_parameters": map[string]string{
						"vpcIds": "${data.alicloud_instances.default.instances[0].vpc_id}",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"input_parameters.%": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"risk_level": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"risk_level": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scope_compliance_resource_id": "${data.alicloud_instances.default.instances[0].id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scope_compliance_resource_id": CHECKSET,
					}),
				),
			},
			// Can not Update when source_owner is ALIYUN.
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"scope_compliance_resource_types": []string{"ACS::OSS::Bucket"},
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"scope_compliance_resource_type.#": "1",
			//		}),
			//	),
			//},
			{
				Config: testAccConfig(map[string]interface{}{
					"source_detail_message_type": "ScheduledNotification",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source_detail_message_type": "ScheduledNotification",
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
					"rule_name":                       name,
					"risk_level":                      "1",
					"scope_compliance_resource_types": []string{"ACS::ECS::Instance"},
					"source_detail_message_type":      "ConfigurationItemChangeNotification",
					"source_identifier":               "ecs-instances-in-vpc",
					"source_owner":                    "ALIYUN",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule_name":                         name,
						"risk_level":                        "1",
						"scope_compliance_resource_types.#": "1",
						"source_detail_message_type":        "ConfigurationItemChangeNotification",
						"source_identifier":                 "ecs-instances-in-vpc",
						"source_owner":                      "ALIYUN",
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

data "alicloud_instances" "default"{}

`, name)
}
