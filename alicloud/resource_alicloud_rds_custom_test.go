package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Rds Custom. >>> Resource test cases, automatically generated.
// Case rdscustom_test_tf_tmk1113_0.1 8923
func TestAccAliCloudRdsCustom_basic8923(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_rds_custom.default"
	ra := resourceAttrInit(resourceId, AlicloudRdsCustomMap8923)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RdsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRdsCustom")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%srdscustom%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudRdsCustomBasicDependence8923)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-chengdu"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"amount":     "1",
					"vswitch_id": "${alicloud_vswitch.vSwitchId.id}",
					//"vswitch_id":    "${data.alicloud_vswitches.default.ids.0}",
					"auto_renew":    "false",
					"period":        "1",
					"auto_pay":      "true",
					"instance_type": "mysql.x2.xlarge.6cm",
					"data_disk": []map[string]interface{}{
						{
							"category":          "cloud_essd",
							"size":              "50",
							"performance_level": "PL1",
						},
					},
					"deployment_set_id": "${alicloud_ecs_deployment_set.deploymentSet.id}",
					"status":            "Running",
					"security_group_ids": []string{
						"${alicloud_security_group.securityGroupId.id}"},
					"io_optimized":                  "optimized",
					"description":                   "jingyi_test",
					"key_pair_name":                 "${alicloud_ecs_key_pair.KeyPairName.key_pair_name}",
					"zone_id":                       "${var.test_zone_id}",
					"instance_charge_type":          "Prepaid",
					"internet_charge_type":          "PayByTraffic",
					"internet_max_bandwidth_out":    "0",
					"image_id":                      "aliyun_2_1903_x64_20G_alibase_20240628.vhd",
					"security_enhancement_strategy": "Active",
					"period_unit":                   "Month",
					"password":                      "jingyiTEST@123",
					"system_disk": []map[string]interface{}{
						{
							"size":     "40",
							"category": "cloud_essd",
						},
					},
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"host_name":         "1731641300",
					"create_mode":       "0",
					"force":             "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"amount":                        "1",
						"vswitch_id":                    CHECKSET,
						"auto_renew":                    "false",
						"period":                        "1",
						"auto_pay":                      "true",
						"instance_type":                 "mysql.x2.xlarge.6cm",
						"data_disk.#":                   "1",
						"deployment_set_id":             CHECKSET,
						"status":                        "Running",
						"security_group_ids.#":          "1",
						"io_optimized":                  "optimized",
						"description":                   "jingyi_test",
						"key_pair_name":                 CHECKSET,
						"zone_id":                       CHECKSET,
						"instance_charge_type":          "Prepaid",
						"internet_charge_type":          "PayByTraffic",
						"internet_max_bandwidth_out":    "0",
						"image_id":                      "aliyun_2_1903_x64_20G_alibase_20240628.vhd",
						"security_enhancement_strategy": "Active",
						"period_unit":                   "Month",
						"password":                      "jingyiTEST@123",
						"resource_group_id":             CHECKSET,
						"host_name":                     CHECKSET,
						"create_mode":                   "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_type":     "mysql.x8.2xlarge.7cm",
					"status":            "Stopped",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"force_stop":        "false",
					"direction":         "Up",
					"dry_run":           "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_type":     "mysql.x8.2xlarge.7cm",
						"status":            "Stopped",
						"resource_group_id": CHECKSET,
						"force_stop":        "false",
						"direction":         "Up",
						"dry_run":           "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status":            "Running",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status":            "Running",
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auto_pay":          "false",
					"instance_type":     "mysql.x2.xlarge.7cm",
					"status":            "Stopped",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"force_stop":        "true",
					"direction":         "Down",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_pay":          "false",
						"instance_type":     "mysql.x2.xlarge.7cm",
						"status":            "Stopped",
						"resource_group_id": CHECKSET,
						"force_stop":        "true",
						"direction":         "Down",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
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
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"amount", "auto_pay", "auto_renew", "create_mode", "direction", "dry_run", "force_stop", "host_name", "image_id", "instance_charge_type", "internet_charge_type", "internet_max_bandwidth_out", "io_optimized", "key_pair_name", "password", "period", "period_unit", "security_enhancement_strategy", "system_disk", "system_disk.category", "system_disk.size", "force"},
			},
		},
	})
}

var AlicloudRdsCustomMap8923 = map[string]string{
	"region_id": CHECKSET,
}

func AlicloudRdsCustomBasicDependence8923(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "test_region_id" {
  default = "cn-chengdu"
}

variable "test_zone_id" {
  default = "cn-chengdu-b"
}

data "alicloud_resource_manager_resource_groups" "default" {}

data "alicloud_vpcs" "default" {
}
data "alicloud_vswitches" "default" {
	vpc_id = data.alicloud_vpcs.default.ids.0
	zone_id      = var.test_zone_id
}

resource "alicloud_vpc" "vpcId" {
 cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "vSwitchId" {
 vpc_id       = alicloud_vpc.vpcId.id
 cidr_block   = "172.16.5.0/24"
 zone_id      = var.test_zone_id
 vswitch_name = format("%%s1", var.name)
}

resource "alicloud_security_group" "securityGroupId" {
  vpc_id = alicloud_vpc.vpcId.id
  //vpc_id = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_ecs_deployment_set" "deploymentSet" {
  domain      = "Default"
  granularity = "Host"
  strategy    = "Availability"
}

resource "alicloud_ecs_key_pair" "KeyPairName" {
  key_pair_name = format("%%s4", var.name)
}


`, name)
}

// Test Rds Custom. <<< Resource test cases, automatically generated.
