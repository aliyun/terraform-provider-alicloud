package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
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
			testAccPreCheckWithRegions(t, true, []connectivity.Region{connectivity.Hangzhou})
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      testAccCheckEbsSolutionInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "description",
					"solution_id": "mysql_compare",
					"parameters": []map[string]interface{}{
						{
							"parameter_key":   "dataDiskType",
							"parameter_value": "cloud_essd",
						},
						{
							"parameter_key":   "mysqlPort",
							"parameter_value": "3306",
						},
						{
							"parameter_key":   "dataDiskPerformance",
							"parameter_value": "PL1",
						},
						{
							"parameter_key":   "ecsPassword",
							"parameter_value": "DHJfuebf123",
						},
						{
							"parameter_key":   "zoneId",
							"parameter_value": "cn-hangzhou-b",
						},
						{
							"parameter_key":   "ecsType",
							"parameter_value": "ecs.c7.xlarge",
						},
						{
							"parameter_key":   "sysDiskType",
							"parameter_value": "cloud_essd",
						},
						{
							"parameter_key":   "sysDiskSize",
							"parameter_value": "40",
						},
						{
							"parameter_key":   "dataDiskSize",
							"parameter_value": "500",
						},
						{
							"parameter_key":   "mysqlPassword",
							"parameter_value": "DHJfuebf123",
						},
					},
					"solution_instance_name": name,
					"resource_group_id":      "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":            "description",
						"solution_id":            "mysql_compare",
						"parameters.#":           "10",
						"solution_instance_name": name,
						"resource_group_id":      CHECKSET,
						"tags.%":                 "2",
						"tags.Created":           "TF",
						"tags.For":               "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
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
					"solution_id": "mysql_compare",
					"parameters": []map[string]interface{}{
						{
							"parameter_key":   "dataDiskType",
							"parameter_value": "cloud_essd",
						},
						{
							"parameter_key":   "mysqlPort",
							"parameter_value": "3306",
						},
						{
							"parameter_key":   "dataDiskPerformance",
							"parameter_value": "PL1",
						},
						{
							"parameter_key":   "ecsPassword",
							"parameter_value": "DHJfuebf123",
						},
						{
							"parameter_key":   "zoneId",
							"parameter_value": "cn-hangzhou-b",
						},
						{
							"parameter_key":   "ecsType",
							"parameter_value": "ecs.c7.xlarge",
						},
						{
							"parameter_key":   "sysDiskType",
							"parameter_value": "cloud_essd",
						},
						{
							"parameter_key":   "sysDiskSize",
							"parameter_value": "40",
						},
						{
							"parameter_key":   "dataDiskSize",
							"parameter_value": "500",
						},
						{
							"parameter_key":   "mysqlPassword",
							"parameter_value": "DHJfuebf123",
						},
					},
					"solution_instance_name": name + "_update",
					"resource_group_id":      "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":            "description",
						"solution_id":            "mysql_compare",
						"parameters.#":           "10",
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
			testAccPreCheckWithRegions(t, true, []connectivity.Region{connectivity.Hangzhou})
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      testAccCheckEbsSolutionInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "description",
					"solution_id": "mysql_compare",
					"parameters": []map[string]interface{}{
						{
							"parameter_key":   "dataDiskType",
							"parameter_value": "cloud_essd",
						},
						{
							"parameter_key":   "mysqlPort",
							"parameter_value": "3306",
						},
						{
							"parameter_key":   "dataDiskPerformance",
							"parameter_value": "PL1",
						},
						{
							"parameter_key":   "ecsPassword",
							"parameter_value": "DHJfuebf123",
						},
						{
							"parameter_key":   "zoneId",
							"parameter_value": "cn-hangzhou-b",
						},
						{
							"parameter_key":   "ecsType",
							"parameter_value": "ecs.c7.xlarge",
						},
						{
							"parameter_key":   "sysDiskType",
							"parameter_value": "cloud_essd",
						},
						{
							"parameter_key":   "sysDiskSize",
							"parameter_value": "40",
						},
						{
							"parameter_key":   "dataDiskSize",
							"parameter_value": "500",
						},
						{
							"parameter_key":   "mysqlPassword",
							"parameter_value": "DHJfuebf123",
						},
					},
					"solution_instance_name": name,
					"resource_group_id":      "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":            "description",
						"solution_id":            "mysql_compare",
						"parameters.#":           "10",
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

func testAccCheckEbsSolutionInstanceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	ebsService := EbsServiceV2{client}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_ebs_solution_instance" {
			continue
		}
		if resp, err := ebsService.DescribeEbsSolutionInstance(rs.Primary.ID); err != nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		} else {
			// check ID and status
			if rs.Primary.ID == fmt.Sprint(resp["SolutionInstanceId"]) && resp["Status"] == "DELETE_COMPLETE" {
				continue
			}
		}
		return fmt.Errorf("EBS Solution Instance %s still exists", rs.Primary.ID)
	}
	return nil
}

// Test Ebs SolutionInstance. <<< Resource test cases, automatically generated.
