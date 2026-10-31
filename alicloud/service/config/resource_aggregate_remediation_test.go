// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package config_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/config"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test Config AggregateRemediation. >>> Resource test cases, automatically generated.
// Case AggregateRemediation-resource-test-251202 11937
func TestAccAliCloudConfigAggregateRemediation_basic11937(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_config_aggregate_remediation.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudConfigAggregateRemediationMap11937)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return config.NewConfigServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeConfigAggregateRemediation")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccconfig%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudConfigAggregateRemediationBasicDependence11937)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
						"config_rule_id":            acctest.CHECKSET,
						"remediation_template_id":   "ACS-TAG-TagResources",
						"remediation_source_type":   "ALIYUN",
						"invoke_type":               "MANUAL_EXECUTION",
						"remediation_type":          "OOS",
						"aggregator_id":             acctest.CHECKSET,
						"remediation_origin_params": acctest.CHECKSET,
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
						"remediation_origin_params": acctest.CHECKSET,
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
	"remediation_id": acctest.CHECKSET,
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

// Representative ACC that does NOT create a new aggregator: it reads the
// account's existing aggregator via the alicloud_config_aggregators data
// source, creates an aggregate config rule against it, then a remediation on
// that rule. Avoids the hardcoded "rd" aggregator-name collision of the
// Dependence11937 fixture (RDAggregatorAlreadyExist) without changing
// Schema/CRUD or deleting the existing basic11937 assertions.
func TestAccAliCloudConfigAggregateRemediation_basic_readonlyAggregator(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_config_aggregate_remediation.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudConfigAggregateRemediationMap11937)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return config.NewConfigServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeConfigAggregateRemediation")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccconfigro%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudConfigAggregateRemediationBasicDependenceReadOnly)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"config_rule_id":            "${alicloud_config_aggregate_config_rule.create-rule.config_rule_id}",
					"remediation_template_id":   "ACS-TAG-TagResources",
					"remediation_source_type":   "ALIYUN",
					"invoke_type":               "MANUAL_EXECUTION",
					"remediation_type":          "OOS",
					"aggregator_id":             "${data.alicloud_config_aggregators.existing.ids.0}",
					"remediation_origin_params": "{\\\"properties\\\":[{\\\"name\\\":\\\"regionId\\\",\\\"type\\\":\\\"String\\\",\\\"value\\\":\\\"{regionId}\\\",\\\"allowedValues\\\":[],\\\"description\\\":\\\"地域ID\\\"},{\\\"name\\\":\\\"tags\\\",\\\"type\\\":\\\"Json\\\",\\\"value\\\":\\\"{\\\\\\\"aaa\\\\\\\":\\\\\\\"bbb\\\\\\\"}\\\",\\\"allowedValues\\\":[],\\\"description\\\":\\\"资源标签（例：{\\\\\\\"k1\\\\\\\":\\\\\\\"v1\\\\\\\",\\\\\\\"k2\\\\\\\":\\\\\\\"v2\\\\\\\"}）。\\\"},{\\\"name\\\":\\\"resourceType\\\",\\\"type\\\":\\\"String\\\",\\\"value\\\":\\\"{resourceType}\\\",\\\"allowedValues\\\":[],\\\"description\\\":\\\"资源类型\\\"},{\\\"name\\\":\\\"resourceIds\\\",\\\"type\\\":\\\"ARRAY\\\",\\\"value\\\":\\\"[\\\\\\\"{resourceId}\\\\\\\"]\\\",\\\"allowedValues\\\":[],\\\"description\\\":\\\"资源ID列表\\\"}]}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"config_rule_id":            acctest.CHECKSET,
						"remediation_template_id":   "ACS-TAG-TagResources",
						"remediation_source_type":   "ALIYUN",
						"invoke_type":               "MANUAL_EXECUTION",
						"remediation_type":          "OOS",
						"aggregator_id":             acctest.CHECKSET,
						"remediation_origin_params": acctest.CHECKSET,
					}),
				),
			},
		},
	})
}

// AlicloudConfigAggregateRemediationBasicDependenceReadOnly reads the
// account's existing aggregator (no new aggregator, no hardcoded "rd"
// collision) and creates an aggregate config rule against it; the remediation
// references that rule.
func AlicloudConfigAggregateRemediationBasicDependenceReadOnly(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_config_aggregators" "existing" {}

resource "alicloud_config_aggregate_config_rule" "create-rule" {
  source_owner               = "ALIYUN"
  source_identifier          = "required-tags"
  aggregate_config_rule_name = var.name
  config_rule_trigger_types  = "ConfigurationItemChangeNotification"
  risk_level                 = "1"
  resource_types_scope       = ["ACS::OSS::Bucket"]
  aggregator_id              = data.alicloud_config_aggregators.existing.ids.0
  input_parameters = {
    tag1Key   = "aaa"
    tag1Value = "bbb"
  }
}

`, name)
}

// Test Config AggregateRemediation. <<< Resource test cases, automatically generated.
