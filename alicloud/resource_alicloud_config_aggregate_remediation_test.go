// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Config AggregateRemediation. >>> Resource test cases, automatically generated.
// Case AggregateRemediation-resource-test-251202 11937
func TestAccAliCloudConfigAggregateRemediation_basic11937(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_config_aggregate_remediation.default"
	ra := resourceAttrInit(resourceId, AlicloudConfigAggregateRemediationMap11937)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ConfigServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeConfigAggregateRemediation")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccconfig%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudConfigAggregateRemediationBasicDependence11937)
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
					"config_rule_id":            "${alicloud_config_aggregate_config_rule.create-rule.config_rule_id}",
					"remediation_template_id":   "ACS-TAG-TagResources",
					"remediation_source_type":   "ALIYUN",
					"invoke_type":               "MANUAL_EXECUTION",
					"remediation_type":          "OOS",
					"aggregator_id":             "${alicloud_config_aggregator.create-agg.id}",
					"remediation_origin_params": "{\\\"properties\\\":[{\\\"name\\\":\\\"regionId\\\",\\\"type\\\":\\\"String\\\",\\\"value\\\":\\\"{regionId}\\\",\\\"allowedValues\\\":[],\\\"description\\\":\\\"地域ID\\\"},{\\\"name\\\":\\\"tags\\\",\\\"type\\\":\\\"Json\\\",\\\"value\\\":\\\"{\\\\\\\"aaa\\\\\\\":\\\\\\\"bbb\\\\\\\"}\\\",\\\"allowedValues\\\":[],\\\"description\\\":\\\"资源标签（例：{\\\\\\\"k1\\\\\\\":\\\\\\\"v1\\\\\\\",\\\\\\\"k2\\\\\\\":\\\\\\\"v2\\\\\\\"}）。\\\"},{\\\"name\\\":\\\"resourceType\\\",\\\"type\\\":\\\"String\\\",\\\"value\\\":\\\"{resourceType}\\\",\\\"allowedValues\\\":[],\\\"description\\\":\\\"资源类型\\\"},{\\\"name\\\":\\\"resourceIds\\\",\\\"type\\\":\\\"ARRAY\\\",\\\"value\\\":\\\"[\\\\\\\"{resourceId}\\\\\\\"]\\\",\\\"allowedValues\\\":[],\\\"description\\\":\\\"资源ID列表\\\"}]}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"config_rule_id":            CHECKSET,
						"remediation_template_id":   "ACS-TAG-TagResources",
						"remediation_source_type":   "ALIYUN",
						"invoke_type":               "MANUAL_EXECUTION",
						"remediation_type":          "OOS",
						"aggregator_id":             CHECKSET,
						"remediation_origin_params": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"remediation_template_id":   "ACS-TAG-TagResourcesIgnoreCaseSensitive",
					"invoke_type":               "AUTO_EXECUTION",
					"remediation_origin_params": "{\\\"properties\\\":[{\\\"name\\\":\\\"regionId\\\",\\\"type\\\":\\\"String\\\",\\\"value\\\":\\\"{regionId}\\\",\\\"allowedValues\\\":[],\\\"description\\\":\\\"地域ID\\\"},{\\\"name\\\":\\\"tags\\\",\\\"type\\\":\\\"String\\\",\\\"value\\\":\\\"{\\\\\\\"aaa\\\\\\\":\\\\\\\"bbb\\\\\\\"}\\\",\\\"allowedValues\\\":[],\\\"description\\\":\\\"\\\\n模版tag参数占位符，不用单独设置。\\\"},{\\\"name\\\":\\\"resourceType\\\",\\\"type\\\":\\\"String\\\",\\\"value\\\":\\\"{resourceType}\\\",\\\"allowedValues\\\":[],\\\"description\\\":\\\"资源类型\\\"},{\\\"name\\\":\\\"resourceIds\\\",\\\"type\\\":\\\"ARRAY\\\",\\\"value\\\":\\\"[\\\\\\\"{resourceId}\\\\\\\"]\\\",\\\"allowedValues\\\":[],\\\"description\\\":\\\"资源ID的列表\\\"}]}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"remediation_template_id":   "ACS-TAG-TagResourcesIgnoreCaseSensitive",
						"invoke_type":               "AUTO_EXECUTION",
						"remediation_origin_params": CHECKSET,
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

var AlicloudConfigAggregateRemediationMap11937 = map[string]string{
	"remediation_id": CHECKSET,
}

func AlicloudConfigAggregateRemediationBasicDependence11937(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_config_aggregator" "create-agg" {
  aggregator_name = "rd"
  description     = "rd"
  aggregator_type = "RD"
}

resource "alicloud_config_aggregate_config_rule" "create-rule" {
  source_owner               = "ALIYUN"
  source_identifier          = "required-tags"
  aggregate_config_rule_name = "agg-rule-name"
  config_rule_trigger_types  = "ConfigurationItemChangeNotification"
  risk_level                 = "1"
  resource_types_scope       = ["ACS::OSS::Bucket"]
  aggregator_id              = alicloud_config_aggregator.create-agg.id
  input_parameters = {
    tag1Key   = "aaa"
    tag1Value = "bbb"
  }
}


`, name)
}

// Test Config AggregateRemediation. <<< Resource test cases, automatically generated.
