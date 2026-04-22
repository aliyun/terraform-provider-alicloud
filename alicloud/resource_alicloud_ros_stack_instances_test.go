package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAliCloudRosStackInstances_basic(t *testing.T) {
	var (
		resourceId = "alicloud_ros_stack_instances.default"
		rand       = acctest.RandIntRange(10000, 99999)
		name       = fmt.Sprintf("tf-testacc%srosstackinstances%d", defaultRegionToTest, rand)
	)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudROSStackInstancesBasicDependence0)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRosStackInstancesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"stack_group_name":      "${alicloud_ros_stack_group.default.stack_group_name}",
					"region_ids":            []string{"${data.alicloud_ros_regions.default.regions.0.region_id}"},
					"account_ids":           []string{"${data.alicloud_account.this.id}"},
					"timeout_in_minutes":    "60",
					"operation_description": "Initial deployment",
					"disable_rollback":      false,
					"operation_preferences": []map[string]interface{}{
						{
							"max_concurrent_count":    5,
							"failure_tolerance_count": 1,
							"region_concurrency_type": "SEQUENTIAL",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceId, "id"),
					resource.TestCheckResourceAttr(resourceId, "timeout_in_minutes", "60"),
					resource.TestCheckResourceAttr(resourceId, "operation_preferences.#", "1"),
					resource.TestCheckResourceAttr(resourceId, "operation_preferences.0.max_concurrent_count", "5"),
					resource.TestCheckResourceAttr(resourceId, "operation_preferences.0.failure_tolerance_count", "1"),
					resource.TestCheckResourceAttr(resourceId, "operation_preferences.0.region_concurrency_type", "SEQUENTIAL"),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"stack_group_name":      "${alicloud_ros_stack_group.default.stack_group_name}",
					"region_ids":            []string{"${data.alicloud_ros_regions.default.regions.0.region_id}"},
					"account_ids":           []string{"${data.alicloud_account.this.id}"},
					"timeout_in_minutes":    "90",
					"operation_description": "Updated deployment",
					"operation_preferences": []map[string]interface{}{
						{
							"max_concurrent_percentage":    80,
							"failure_tolerance_percentage": 20,
							"region_concurrency_type":      "PARALLEL",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "timeout_in_minutes", "90"),
					resource.TestCheckResourceAttr(resourceId, "operation_preferences.0.max_concurrent_percentage", "80"),
					resource.TestCheckResourceAttr(resourceId, "operation_preferences.0.failure_tolerance_percentage", "20"),
					resource.TestCheckResourceAttr(resourceId, "operation_preferences.0.region_concurrency_type", "PARALLEL"),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"stack_group_name":      "${alicloud_ros_stack_group.default.stack_group_name}",
					"region_ids":            []string{"${data.alicloud_ros_regions.default.regions.0.region_id}"},
					"account_ids":           []string{"${data.alicloud_account.this.id}"},
					"timeout_in_minutes":    "120",
					"operation_description": "Updated deployment v2",
					"operation_preferences": []map[string]interface{}{
						{
							"max_concurrent_count":    3,
							"failure_tolerance_count": 2,
							"region_concurrency_type": "SEQUENTIAL",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "timeout_in_minutes", "120"),
					resource.TestCheckResourceAttr(resourceId, "operation_preferences.0.max_concurrent_count", "3"),
					resource.TestCheckResourceAttr(resourceId, "operation_preferences.0.failure_tolerance_count", "2"),
					resource.TestCheckResourceAttr(resourceId, "operation_preferences.0.region_concurrency_type", "SEQUENTIAL"),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"stack_group_name":   "${alicloud_ros_stack_group.default.stack_group_name}",
					"region_ids":         []string{"${data.alicloud_ros_regions.default.regions.0.region_id}"},
					"account_ids":        []string{"${data.alicloud_account.this.id}"},
					"timeout_in_minutes": "120",
					"operation_preferences": []map[string]interface{}{
						{
							"max_concurrent_count":    3,
							"failure_tolerance_count": 2,
							"region_concurrency_type": "SEQUENTIAL",
						},
					},
					"parameter_overrides": []map[string]interface{}{
						{
							"parameter_key":   "VpcName",
							"parameter_value": "TestVpcNameUpdated",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "parameter_overrides.#", "1"),
				),
			},
		},
	})
}

func TestAccAliCloudRosStackInstances_serviceManaged(t *testing.T) {
	var (
		resourceId = "alicloud_ros_stack_instances.default"
		rand       = acctest.RandIntRange(10000, 99999)
		name       = fmt.Sprintf("tf-testacc%srosstackinstances%d", defaultRegionToTest, rand)
	)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudROSStackInstancesServiceManagedDependence0)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRosStackInstancesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"stack_group_name":      "${alicloud_ros_stack_group.default.stack_group_name}",
					"region_ids":            []string{"${data.alicloud_ros_regions.default.regions.0.region_id}"},
					"operation_description": "Updated deployment",
					"deployment_targets": []map[string]interface{}{
						{
							"rd_folder_ids": []string{"${alicloud_resource_manager_folder.test_folder.id}"},
						},
					},
					"deployment_options": []string{"IgnoreExisting"},
					"disable_rollback":   true,
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceId, "id"),
					resource.TestCheckResourceAttr(resourceId, "deployment_targets.#", "1"),
					resource.TestCheckResourceAttr(resourceId, "deployment_targets.0.rd_folder_ids.#", "1"),
					resource.TestCheckResourceAttr(resourceId, "disable_rollback", "true"),
					resource.TestCheckResourceAttr(resourceId, "deployment_options.#", "1"),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"stack_group_name":      "${alicloud_ros_stack_group.default.stack_group_name}",
					"region_ids":            []string{"${data.alicloud_ros_regions.default.regions.0.region_id}"},
					"operation_description": "Updated deployment v2",
					"deployment_targets": []map[string]interface{}{
						{
							"account_ids":   []interface{}{},
							"rd_folder_ids": []string{"${alicloud_resource_manager_folder.test_folder2.id}"},
						},
					},
					"deployment_options": []string{"IgnoreExisting"},
					"disable_rollback":   true,
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "deployment_targets.0.account_ids.#", "0"),
					resource.TestCheckResourceAttr(resourceId, "deployment_targets.0.rd_folder_ids.#", "1"),
				),
			},
		},
	})
}

func testAccCheckRosStackInstancesDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_ros_stack_instances" {
			continue
		}
		stackGroupName := rs.Primary.Attributes["stack_group_name"]
		if stackGroupName == "" {
			stackGroupName = rs.Primary.ID
		}
		if stackGroupName == "" {
			continue
		}

		pageNum := 1
		for {
			req := map[string]interface{}{
				"RegionId":       client.RegionId,
				"StackGroupName": stackGroupName,
				"PageNumber":     pageNum,
				"PageSize":       50,
			}
			resp, err := client.RpcPost("ROS", "2019-09-10", "ListStackInstances", nil, req, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"StackGroupNotFound"}) {
					break
				}
				return err
			}
			if instances, ok := resp["StackInstances"].([]interface{}); ok && len(instances) > 0 {
				return fmt.Errorf("Stack Instances still exist for group %s", stackGroupName)
			}
			if total, ok := resp["TotalCount"].(float64); ok && int(total) <= pageNum*50 {
				break
			}
			pageNum++
		}
	}
	return nil
}

func AlicloudROSStackInstancesBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
data "alicloud_account" "this" {}
data "alicloud_ros_regions" "default" {}

resource "alicloud_ros_stack_group" "default" {
  stack_group_name = var.name
  template_body    = jsonencode({
    "ROSTemplateFormatVersion": "2015-09-01",
    "Parameters": {
      "VpcName": { "Type": "String", "Default": "DefaultVpc" },
      "InstanceType": { "Type": "String", "Default": "ecs.g6.large" }
    },
    "Resources": {}
  })
  description = "test for stack groups"
  parameters {
    parameter_key   = "InstanceType"
    parameter_value = "ecs.g6.large"
  }
  parameters {
    parameter_key   = "VpcName"
    parameter_value = "GroupLevelVpc"
  }
}
`, name)
}

func AlicloudROSStackInstancesServiceManagedDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_account" "this" {}
data "alicloud_ros_regions" "default" {}

resource "alicloud_resource_manager_folder" "test_folder" {
  folder_name = "tf-test-ros-${substr(md5(var.name), 0, 6)}"
}

resource "alicloud_resource_manager_folder" "test_folder2" {
  folder_name = "tf-test-ros2-${substr(md5(var.name), 0, 6)}"
}

resource "alicloud_ros_stack_group" "default" {
  stack_group_name   = var.name
  auto_deployment {
	enabled = true
}
  permission_model   = "SERVICE_MANAGED"

  template_body = jsonencode({
    "ROSTemplateFormatVersion": "2015-09-01",
    "Resources": {}
  })
  description = "test for stack groups service managed"
}
`, name)
}
