package alicloud

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudROSStackInstance_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ros_stack_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudROSStackInstanceMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RosService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRosStackInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%srosstackinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudROSStackInstanceBasicDependence0)
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
					"stack_instance_account_id": "${var.account_id}",
					"stack_group_name":          "${alicloud_ros_stack_group.default.stack_group_name}",
					"stack_instance_region_id":  "${data.alicloud_ros_regions.default.regions.0.region_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"stack_instance_account_id": CHECKSET,
						"stack_group_name":          CHECKSET,
						"stack_instance_region_id":  CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"parameter_overrides": []map[string]interface{}{
						{
							"parameter_value": "VpcName",
							"parameter_key":   "VpcName",
						},
						{
							"parameter_value": "InstanceType",
							"parameter_key":   "InstanceType",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"parameter_overrides.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"timeout_in_minutes": "60",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"operation_preferences": `{\"FailureToleranceCount\": 1,\"MaxConcurrentCount\": 2}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"operation_description": "tf-acctest",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeout_in_minutes", "retain_stacks", "operation_description", "operation_preferences"},
			},
		},
	})
}

func TestAccAlicloudROSStackInstance_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ros_stack_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudROSStackInstanceMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RosService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRosStackInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%srosstackinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudROSStackInstanceBasicDependence0)
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
					"stack_instance_account_id": "${var.account_id}",
					"stack_group_name":          "${alicloud_ros_stack_group.default.stack_group_name}",
					"stack_instance_region_id":  "${data.alicloud_ros_regions.default.regions.0.region_id}",
					"parameter_overrides": []map[string]interface{}{
						{
							"parameter_value": "VpcName",
							"parameter_key":   "VpcName",
						},
						{
							"parameter_value": "InstanceType",
							"parameter_key":   "InstanceType",
						},
					},
					"timeout_in_minutes":    "60",
					"operation_preferences": `{\"FailureToleranceCount\": 1,\"MaxConcurrentCount\": 2}`,
					"operation_description": "tf-acctest",
					"retain_stacks":         "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"stack_instance_account_id": CHECKSET,
						"stack_group_name":          CHECKSET,
						"stack_instance_region_id":  CHECKSET,
						"parameter_overrides.#":     "2",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeout_in_minutes", "retain_stacks", "operation_description", "operation_preferences"},
			},
		},
	})
}

var AlicloudROSStackInstanceMap0 = map[string]string{
	"parameter_overrides.#":     CHECKSET,
	"stack_instance_account_id": CHECKSET,
	"stack_group_name":          CHECKSET,
	"status":                    CHECKSET,
	"stack_instance_region_id":  CHECKSET,
}

func AlicloudROSStackInstanceBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}

variable "account_id" {
  default = "%s"
}

data "alicloud_ros_regions" "default" {}

resource "alicloud_ros_stack_group" "default" {
  stack_group_name=   var.name
  template_body=   "{\"ROSTemplateFormatVersion\":\"2015-09-01\", \"Parameters\": {\"VpcName\": {\"Type\": \"String\"},\"InstanceType\": {\"Type\": \"String\"}}}"
  description = "test for stack groups"
  parameters {
	parameter_key =   "VpcName"
	parameter_value = "VpcName"
  }
  parameters {
	parameter_key =   "InstanceType"
	parameter_value = "InstanceType"
  }
}
`, name, os.Getenv("ALICLOUD_ACCOUNT_ID"))
}
