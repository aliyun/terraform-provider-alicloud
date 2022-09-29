package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudGPDBDBInstance_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_gpdb_instance.default"
	checkoutSupportedRegions(t, true, connectivity.GPDBSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudGPDBDBInstanceMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GpdbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGpdbDbInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sgpdbdbinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudGPDBDBInstanceBasicDependence0)
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
					"db_instance_category":  "HighAvailability",
					"db_instance_class":     "gpdb.group.segsdx1",
					"db_instance_mode":      "StorageElastic",
					"description":           name,
					"engine":                "gpdb",
					"engine_version":        "6.0",
					"zone_id":               "${data.alicloud_gpdb_zones.default.ids.0}",
					"instance_network_type": "VPC",
					"instance_spec":         "2C16G",
					"master_node_num":       "1",
					"payment_type":          "PayAsYouGo",
					"private_ip_address":    "1.1.1.1",
					"seg_storage_type":      "cloud_essd",
					"seg_node_num":          "4",
					"storage_size":          "50",
					"vpc_id":                "${data.alicloud_vpcs.default.ids.0}",
					"vswitch_id":            "${local.vswitch_id}",
					"create_sample_data":    `false`,
					"ip_whitelist": []map[string]interface{}{
						{
							"security_ip_list": "127.0.0.1",
						},
					},
					"tags": map[string]string{
						"Created": "TF",
						"For":     "acceptance test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_category":  "HighAvailability",
						"db_instance_mode":      "StorageElastic",
						"description":           name,
						"engine":                "gpdb",
						"engine_version":        "6.0",
						"zone_id":               CHECKSET,
						"instance_network_type": "VPC",
						"instance_spec":         "2C16G",
						"master_node_num":       "1",
						"payment_type":          "PayAsYouGo",
						"private_ip_address":    "1.1.1.1",
						"seg_node_num":          "4",
						"storage_size":          "50",
						"vpc_id":                CHECKSET,
						"vswitch_id":            CHECKSET,
						"tags.%":                "2",
						"tags.Created":          "TF",
						"tags.For":              "acceptance test",
						"ip_whitelist.#":        "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name + "Update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name + "Update",
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
					"maintain_start_time": "08:00Z",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"maintain_start_time": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"maintain_end_time": "12:00Z",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"maintain_end_time": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ip_whitelist": []map[string]interface{}{
						{
							"security_ip_list": "1.1.1.1",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ip_whitelist.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"seg_node_num": "8",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"seg_node_num": "8",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"master_node_num": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"master_node_num": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_spec": "4C32G",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_spec": "4C32G",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"storage_size": "100",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"storage_size": "100",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description":         name + "Update",
					"instance_spec":       "2C16G",
					"master_node_num":     "1",
					"seg_node_num":        "12",
					"storage_size":        "200",
					"maintain_start_time": "09:00Z",
					"maintain_end_time":   "13:00Z",
					"ip_whitelist": []map[string]interface{}{
						{
							"security_ip_list": "127.0.0.1",
						},
					},
					"tags": map[string]string{
						"Created": "TF2",
						"For":     "acceptance test2",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":         name + "Update",
						"instance_spec":       "2C16G",
						"master_node_num":     "1",
						"seg_node_num":        "12",
						"storage_size":        "200",
						"maintain_start_time": CHECKSET,
						"maintain_end_time":   CHECKSET,
						"ip_whitelist.#":      "1",
						"tags.%":              "2",
						"tags.Created":        "TF2",
						"tags.For":            "acceptance test2",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"period", "used_time", "seg_storage_type", "private_ip_address", "instance_spec", "db_instance_class", "resource_group_id", "create_sample_data"},
			},
		},
	})
}

func TestAccAlicloudGPDBDBInstanceServerless(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_gpdb_instance.default"
	checkoutSupportedRegions(t, true, connectivity.GPDBSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudGPDBDBInstanceMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GpdbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGpdbDbInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sgpdbdbinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudGPDBDBInstanceBasicDependence1)
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
					"db_instance_mode":      "Serverless",
					"description":           name,
					"engine":                "gpdb",
					"engine_version":        "6.0",
					"zone_id":               "${data.alicloud_gpdb_zones.default.ids.2}",
					"instance_network_type": "VPC",
					"instance_spec":         "4C16G",
					"master_node_num":       "1",
					"payment_type":          "PayAsYouGo",
					"private_ip_address":    "1.1.1.1",
					"seg_node_num":          "2",
					"vpc_id":                "${data.alicloud_vpcs.default.ids.0}",
					"vswitch_id":            "${local.vswitch_id}",
					"create_sample_data":    `false`,
					"ip_whitelist": []map[string]interface{}{
						{
							"security_ip_list": "127.0.0.1",
						},
					},
					"tags": map[string]string{
						"Created": "TF",
						"For":     "acceptance test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_mode":      "Serverless",
						"description":           name,
						"engine":                "gpdb",
						"engine_version":        "6.0",
						"zone_id":               CHECKSET,
						"instance_network_type": "VPC",
						"instance_spec":         "4C16G",
						"master_node_num":       "1",
						"payment_type":          "PayAsYouGo",
						"private_ip_address":    "1.1.1.1",
						"seg_node_num":          "2",
						"vpc_id":                CHECKSET,
						"vswitch_id":            CHECKSET,
						"ip_whitelist.#":        "1",
						"tags.%":                "2",
						"tags.Created":          "TF",
						"tags.For":              "acceptance test",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"period", "used_time", "seg_storage_type", "private_ip_address", "instance_spec", "db_instance_class", "resource_group_id", "create_sample_data"},
			},
		},
	})
}

func TestAccAlicloudGPDBDBInstancePrepaid(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_gpdb_instance.default"
	checkoutSupportedRegions(t, true, connectivity.GPDBSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudGPDBDBInstanceMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GpdbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGpdbDbInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sgpdbdbinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudGPDBDBInstanceBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_category":  "HighAvailability",
					"db_instance_class":     "gpdb.group.segsdx1",
					"db_instance_mode":      "StorageElastic",
					"description":           name,
					"engine":                "gpdb",
					"engine_version":        "6.0",
					"zone_id":               "${data.alicloud_gpdb_zones.default.ids.0}",
					"instance_network_type": "VPC",
					"instance_spec":         "2C16G",
					"master_node_num":       "1",
					"payment_type":          "Subscription",
					"private_ip_address":    "1.1.1.1",
					"seg_storage_type":      "cloud_essd",
					"seg_node_num":          "4",
					"storage_size":          "50",
					"vpc_id":                "${data.alicloud_vpcs.default.ids.0}",
					"vswitch_id":            "${local.vswitch_id}",
					"period":                "Month",
					"used_time":             "1",
					"create_sample_data":    `false`,
					"ip_whitelist": []map[string]interface{}{
						{
							"security_ip_list": "127.0.0.1",
						},
					},
					"tags": map[string]string{
						"Created": "TF",
						"For":     "acceptance test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_category":  "HighAvailability",
						"db_instance_mode":      "StorageElastic",
						"description":           name,
						"engine":                "gpdb",
						"engine_version":        "6.0",
						"zone_id":               CHECKSET,
						"instance_network_type": "VPC",
						"instance_spec":         "2C16G",
						"master_node_num":       "1",
						"payment_type":          "Subscription",
						"private_ip_address":    "1.1.1.1",
						"seg_node_num":          "4",
						"storage_size":          "50",
						"vpc_id":                CHECKSET,
						"vswitch_id":            CHECKSET,
						"tags.%":                "2",
						"tags.Created":          "TF",
						"tags.For":              "acceptance test",
						"ip_whitelist.#":        "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"period", "used_time", "seg_storage_type", "private_ip_address", "instance_spec", "db_instance_class", "resource_group_id", "create_sample_data"},
			},
		},
	})
}

var AlicloudGPDBDBInstanceMap0 = map[string]string{
	"status": CHECKSET,
}

func AlicloudGPDBDBInstanceBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
data "alicloud_resource_manager_resource_groups" "default" {}
data "alicloud_gpdb_zones" "default" {}
data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}
data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_gpdb_zones.default.ids.0
}
resource "alicloud_vswitch" "vswitch" {
  count        = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
  vpc_id       = data.alicloud_vpcs.default.ids.0
  cidr_block   = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 8)
  zone_id      = data.alicloud_gpdb_zones.default.ids.0
  vswitch_name = var.name
}
locals {
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.vswitch.*.id, [""])[0]
}
`, name)
}

func AlicloudGPDBDBInstanceBasicDependence1(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
data "alicloud_resource_manager_resource_groups" "default" {}
data "alicloud_gpdb_zones" "default" {}
data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}
data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_gpdb_zones.default.ids.2
}
resource "alicloud_vswitch" "vswitch" {
  count        = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
  vpc_id       = data.alicloud_vpcs.default.ids.0
  cidr_block   = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 8)
  zone_id      = data.alicloud_gpdb_zones.default.ids.2
  vswitch_name = var.name
}
locals {
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.vswitch.*.id, [""])[0]
}
`, name)
}
