package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudRosTemplateEstimateCostDataSource_rosTemplate(t *testing.T) {
	checkoutSupportedRegions(t, true, connectivity.ROSSupportRegions)
	rand := acctest.RandIntRange(100000, 999999)
	resourceId := "data.alicloud_ros_template_estimate_cost.default"
	name := fmt.Sprintf("tf-testacc-roscost-%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceRosTemplateEstimateCostDependence)

	templateBodyConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"template_body": "${local.estimate_template}",
			"client_token":  name,
			"parameters": []map[string]interface{}{
				{
					"parameter_key":   "Bandwidth",
					"parameter_value": "5",
				},
			},
			"output_file": "./tf-testacc-ros-template-estimate-cost.txt",
		}),
	}
	templateIdConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"template_id":      "${alicloud_ros_template.default.id}",
			"template_version": "v1",
		}),
	}
	stackConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"template_body": "${local.estimate_template}",
			"stack_id":      "${alicloud_ros_stack.default.id}",
			"parameters": []map[string]interface{}{
				{
					"parameter_key":   "Bandwidth",
					"parameter_value": "10",
				},
			},
		}),
	}

	var existRosTemplateEstimateCostMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"resources.#":                               "1",
			"resources.0.resource_name":                 "NewEip",
			"resources.0.resource_type":                 "ALIYUN::VPC::EIP",
			"resources.0.success":                       "true",
			"resources.0.result.#":                      "1",
			"resources.0.result.0.inquiry_type":         "Buy",
			"resources.0.result.0.order.#":              "1",
			"resources.0.result.0.order.0.currency":     CHECKSET,
			"resources.0.result.0.order.0.trade_amount": CHECKSET,
			"resources.0.result.0.order_details.#":      CHECKSET,
			"resources.0.result.0.order_supplement.#":   "1",
		}
	}

	var fakeRosTemplateEstimateCostMapFunc = func(rand int) map[string]string {
		return map[string]string{}
	}

	var RosTemplateEstimateCostCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existRosTemplateEstimateCostMapFunc,
		fakeMapFunc:  fakeRosTemplateEstimateCostMapFunc,
	}

	RosTemplateEstimateCostCheckInfo.dataSourceTestCheck(t, rand, templateBodyConf, templateIdConf, stackConf)
}

func TestAccAlicloudRosTemplateEstimateCostDataSource_terraformTemplate(t *testing.T) {
	checkoutSupportedRegions(t, true, connectivity.ROSSupportRegions)
	rand := acctest.RandIntRange(100000, 999999)
	resourceId := "data.alicloud_ros_template_estimate_cost.default"
	name := fmt.Sprintf("tf-testacc-roscosttf-%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceRosTemplateEstimateCostTerraformDependence)

	terraformTemplateConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"template_body": "${local.tf_estimate_template}",
		}),
	}

	var existRosTemplateEstimateCostTerraformMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"resources.#":                               "1",
			"resources.0.resource_name":                 "example",
			"resources.0.resource_type":                 "alicloud_eip_address",
			"resources.0.alias_type":                    "ALIYUN::VPC::EIP",
			"resources.0.success":                       "true",
			"resources.0.result.#":                      "1",
			"resources.0.result.0.inquiry_type":         "Buy",
			"resources.0.result.0.order.#":              "1",
			"resources.0.result.0.order.0.currency":     CHECKSET,
			"resources.0.result.0.order.0.trade_amount": CHECKSET,
			"resources.0.result.0.order_details.#":      CHECKSET,
			"resources.0.result.0.order_supplement.#":   "1",
		}
	}

	var fakeRosTemplateEstimateCostTerraformMapFunc = func(rand int) map[string]string {
		return map[string]string{}
	}

	var RosTemplateEstimateCostTerraformCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existRosTemplateEstimateCostTerraformMapFunc,
		fakeMapFunc:  fakeRosTemplateEstimateCostTerraformMapFunc,
	}

	RosTemplateEstimateCostTerraformCheckInfo.dataSourceTestCheck(t, rand, terraformTemplateConf)
}

func TestAccAlicloudRosTemplateEstimateCostDataSource_templateScratch(t *testing.T) {
	checkoutSupportedRegions(t, true, connectivity.ROSSupportRegions)
	rand := acctest.RandIntRange(100000, 999999)
	resourceId := "data.alicloud_ros_template_estimate_cost.default"
	name := fmt.Sprintf("tf-testacc-roscostts-%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceRosTemplateEstimateCostScratchDependence)

	templateScratchConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"template_scratch_id":        "${alicloud_ros_template_scratch.default.id}",
			"template_scratch_region_id": "${data.alicloud_regions.default.regions.0.id}",
		}),
	}

	var existRosTemplateEstimateCostScratchMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"resources.#": "0",
		}
	}

	var fakeRosTemplateEstimateCostScratchMapFunc = func(rand int) map[string]string {
		return map[string]string{}
	}

	var RosTemplateEstimateCostScratchCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existRosTemplateEstimateCostScratchMapFunc,
		fakeMapFunc:  fakeRosTemplateEstimateCostScratchMapFunc,
	}

	RosTemplateEstimateCostScratchCheckInfo.dataSourceTestCheck(t, rand, templateScratchConf)
}

func dataSourceRosTemplateEstimateCostDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

locals {
  estimate_template = jsonencode({
    ROSTemplateFormatVersion = "2015-09-01"
    Parameters = {
      Bandwidth = {
        Type    = "Number"
        Default = 5
      }
    }
    Resources = {
      NewEip = {
        Type = "ALIYUN::VPC::EIP"
        Properties = {
          InstanceChargeType = "Postpaid"
          InternetChargeType = "PayByTraffic"
          Bandwidth = {
            Ref = "Bandwidth"
          }
        }
      }
    }
  })
}

resource "alicloud_ros_template" "default" {
  template_name = var.name
  template_body = local.estimate_template
}

resource "alicloud_ros_stack" "default" {
  stack_name    = var.name
  template_body = local.estimate_template
  parameters {
    parameter_key   = "Bandwidth"
    parameter_value = "5"
  }
}
`, name)
}

func dataSourceRosTemplateEstimateCostTerraformDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

locals {
  tf_estimate_template = jsonencode({
    ROSTemplateFormatVersion = "2015-09-01"
    Transform                = "Aliyun::Terraform-v1.2"
    Workspace = {
      "main.tf" = <<-EOT
        provider "alicloud" {
          region = "cn-hangzhou"
        }
        resource "alicloud_eip_address" "example" {
          bandwidth            = 5
          internet_charge_type = "PayByTraffic"
          payment_type         = "Subscription"
        }
      EOT
    }
  })
}
`, name)
}

func dataSourceRosTemplateEstimateCostScratchDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_regions" "default" {
  current = true
}

data "alicloud_resource_manager_resource_groups" "default" {
}

resource "alicloud_ros_template_scratch" "default" {
  description           = var.name
  template_scratch_type = "ArchitectureReplication"
  source_resource_group {
    resource_group_id    = data.alicloud_resource_manager_resource_groups.default.ids.0
    resource_type_filter = ["ALIYUN::ECS::VPC"]
  }
}
`, name)
}
