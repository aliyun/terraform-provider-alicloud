package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Ebs SolutionInstance. >>> Resource test cases, automatically generated.
// Case 创建并更新用例 5796
func TestAccAliCloudEbsSolutionInstance_basic5796(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ebs_solution_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudEbsSolutionInstanceMap5796)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EbsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEbsSolutionInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sebssolutioninstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEbsSolutionInstanceBasicDependence5796)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{connectivity.Shanghai})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "description",
					"solution_id": "mysql",
					"parameters": []map[string]interface{}{
						{
							"parameter_key":   "zoneId",
							"parameter_value": "${var.zone_id}",
						},
						{
							"parameter_key":   "ecsType",
							"parameter_value": "ecs.c6.large",
						},
						{
							"parameter_key":   "ecsImageId",
							"parameter_value": "CentOS_7",
						},
						{
							"parameter_key":   "internetMaxBandwidthOut",
							"parameter_value": "100",
						},
						{
							"parameter_key":   "internetChargeType",
							"parameter_value": "PayByTraffic",
						},
						{
							"parameter_key":   "ecsPassword",
							"parameter_value": "Ebs12345",
						},
						{
							"parameter_key":   "sysDiskType",
							"parameter_value": "cloud_essd",
						},
						{
							"parameter_key":   "sysDiskPerformance",
							"parameter_value": "PL0",
						},
						{
							"parameter_key":   "sysDiskSize",
							"parameter_value": "40",
						},
						{
							"parameter_key":   "dataDiskType",
							"parameter_value": "cloud_essd",
						},
						{
							"parameter_key":   "dataDiskPerformance",
							"parameter_value": "PL0",
						},
						{
							"parameter_key":   "dataDiskSize",
							"parameter_value": "40",
						},
						{
							"parameter_key":   "mysqlVersion",
							"parameter_value": "MySQL80",
						},
						{
							"parameter_key":   "mysqlUser",
							"parameter_value": "root",
						},
						{
							"parameter_key":   "mysqlPassword",
							"parameter_value": "Ebs12345",
						},
					},
					"solution_instance_name": name,
					"resource_group_id":      "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":            "description",
						"solution_id":            "mysql",
						"parameters.#":           "15",
						"solution_instance_name": name,
						"resource_group_id":      CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "description",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "description",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "newDescription",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "newDescription",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"solution_instance_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"solution_instance_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "description",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "description",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"solution_instance_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"solution_instance_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "description",
					"solution_id": "mysql",
					"parameters": []map[string]interface{}{
						{
							"parameter_key":   "zoneId",
							"parameter_value": "${var.zone_id}",
						},
						{
							"parameter_key":   "ecsType",
							"parameter_value": "ecs.c6.large",
						},
						{
							"parameter_key":   "ecsImageId",
							"parameter_value": "CentOS_7",
						},
						{
							"parameter_key":   "internetMaxBandwidthOut",
							"parameter_value": "100",
						},
						{
							"parameter_key":   "internetChargeType",
							"parameter_value": "PayByTraffic",
						},
						{
							"parameter_key":   "ecsPassword",
							"parameter_value": "Ebs12345",
						},
						{
							"parameter_key":   "sysDiskType",
							"parameter_value": "cloud_essd",
						},
						{
							"parameter_key":   "sysDiskPerformance",
							"parameter_value": "PL0",
						},
						{
							"parameter_key":   "sysDiskSize",
							"parameter_value": "40",
						},
						{
							"parameter_key":   "dataDiskType",
							"parameter_value": "cloud_essd",
						},
						{
							"parameter_key":   "dataDiskPerformance",
							"parameter_value": "PL0",
						},
						{
							"parameter_key":   "dataDiskSize",
							"parameter_value": "40",
						},
						{
							"parameter_key":   "mysqlVersion",
							"parameter_value": "MySQL80",
						},
						{
							"parameter_key":   "mysqlUser",
							"parameter_value": "root",
						},
						{
							"parameter_key":   "mysqlPassword",
							"parameter_value": "Ebs12345",
						},
					},
					"solution_instance_name": name + "_update",
					"resource_group_id":      "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":            "description",
						"solution_id":            "mysql",
						"parameters.#":           "15",
						"solution_instance_name": name + "_update",
						"resource_group_id":      CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parameters"},
			},
		},
	})
}

var AlicloudEbsSolutionInstanceMap5796 = map[string]string{
	"status":                 CHECKSET,
	"create_time":            CHECKSET,
	"solution_instance_name": CHECKSET,
}

func AlicloudEbsSolutionInstanceBasicDependence5796(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "zone_id" {
  default = "cn-shanghai-l"
}

variable "region_id" {
  default = "cn-shanghai"
}

data "alicloud_resource_manager_resource_groups" "default" {}


`, name)
}

// Case 创建并更新用例 5796  twin
func TestAccAliCloudEbsSolutionInstance_basic5796_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ebs_solution_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudEbsSolutionInstanceMap5796)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EbsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEbsSolutionInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sebssolutioninstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEbsSolutionInstanceBasicDependence5796)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{connectivity.Shanghai})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "description",
					"solution_id": "mysql",
					"parameters": []map[string]interface{}{
						{
							"parameter_key":   "zoneId",
							"parameter_value": "${var.zone_id}",
						},
						{
							"parameter_key":   "ecsType",
							"parameter_value": "ecs.c6.large",
						},
						{
							"parameter_key":   "ecsImageId",
							"parameter_value": "CentOS_7",
						},
						{
							"parameter_key":   "internetMaxBandwidthOut",
							"parameter_value": "100",
						},
						{
							"parameter_key":   "internetChargeType",
							"parameter_value": "PayByTraffic",
						},
						{
							"parameter_key":   "ecsPassword",
							"parameter_value": "Ebs12345",
						},
						{
							"parameter_key":   "sysDiskType",
							"parameter_value": "cloud_essd",
						},
						{
							"parameter_key":   "sysDiskPerformance",
							"parameter_value": "PL0",
						},
						{
							"parameter_key":   "sysDiskSize",
							"parameter_value": "40",
						},
						{
							"parameter_key":   "dataDiskType",
							"parameter_value": "cloud_essd",
						},
						{
							"parameter_key":   "dataDiskPerformance",
							"parameter_value": "PL0",
						},
						{
							"parameter_key":   "dataDiskSize",
							"parameter_value": "40",
						},
						{
							"parameter_key":   "mysqlVersion",
							"parameter_value": "MySQL80",
						},
						{
							"parameter_key":   "mysqlUser",
							"parameter_value": "root",
						},
						{
							"parameter_key":   "mysqlPassword",
							"parameter_value": "Ebs12345",
						},
					},
					"solution_instance_name": name,
					"resource_group_id":      "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":            "description",
						"solution_id":            "mysql",
						"parameters.#":           "15",
						"solution_instance_name": name,
						"resource_group_id":      CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parameters"},
			},
		},
	})
}

// Test Ebs SolutionInstance. <<< Resource test cases, automatically generated.
