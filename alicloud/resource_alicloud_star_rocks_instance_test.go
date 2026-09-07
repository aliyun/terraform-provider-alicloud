// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test StarRocks Instance. >>> Resource test cases, automatically generated.
// Case Instance 11073
func TestAccAliCloudStarRocksInstance_basic11073(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_star_rocks_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudStarRocksInstanceMap11073)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &StarRocksServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeStarRocksInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccstarrocks%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudStarRocksInstanceBasicDependence11073)
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
					"instance_name": name,
					"auto_renew":    "false",
					"frontend_node_groups": []map[string]interface{}{
						{
							"cu":                          "8",
							"storage_size":                "100",
							"resident_node_number":        "3",
							"storage_performance_level":   "pl1",
							"spec_type":                   "standard",
							"disk_number":                 "1",
							"zone_id":                     "cn-hangzhou-i",
							"local_storage_instance_type": "null",
						},
					},
					"vswitches": []map[string]interface{}{
						{
							"vswitch_id": "${alicloud_vswitch.default106DkE.id}",
							"zone_id":    "cn-hangzhou-i",
						},
					},
					"backend_node_groups": []map[string]interface{}{
						{
							"cu":                          "8",
							"storage_size":                "100",
							"resident_node_number":        "3",
							"disk_number":                 "1",
							"storage_performance_level":   "pl1",
							"spec_type":                   "standard",
							"zone_id":                     "cn-hangzhou-i",
							"local_storage_instance_type": "null",
						},
					},
					"cluster_zone_id":         "cn-hangzhou-i",
					"duration":                "1",
					"pay_type":                "postPaid",
					"vpc_id":                  "${alicloud_vpc.defaultB21JUD.id}",
					"version":                 "3.3",
					"run_mode":                "shared_data",
					"package_type":            "official",
					"admin_password":          "1qaz@QAZ",
					"oss_accessing_role_name": "AliyunEMRStarRocksAccessingOSSRole",
					"pricing_cycle":           "Month",
					"kms_key_id":              "123",
					"promotion_option_no":     "123",
					"encrypted":               "false",
					"resource_group_id":       "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"observer_node_groups": []map[string]interface{}{
						{
							"cu":                          "8",
							"storage_size":                "100",
							"storage_performance_level":   "pl1",
							"disk_number":                 "1",
							"resident_node_number":        "1",
							"spec_type":                   "standard",
							"local_storage_instance_type": "null",
							"zone_id":                     "cn-hangzhou-h",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name":           name,
						"auto_renew":              "false",
						"frontend_node_groups.#":  "1",
						"vswitches.#":             "1",
						"backend_node_groups.#":   "1",
						"cluster_zone_id":         "cn-hangzhou-i",
						"duration":                "1",
						"pay_type":                "postPaid",
						"vpc_id":                  CHECKSET,
						"version":                 CHECKSET,
						"run_mode":                "shared_data",
						"package_type":            "official",
						"admin_password":          "1qaz@QAZ",
						"oss_accessing_role_name": "AliyunEMRStarRocksAccessingOSSRole",
						"pricing_cycle":           "Month",
						"kms_key_id":              CHECKSET,
						"promotion_option_no":     CHECKSET,
						"encrypted":               "false",
						"resource_group_id":       CHECKSET,
						"observer_node_groups.#":  "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
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
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
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
				ImportStateVerifyIgnore: []string{"admin_password", "auto_renew", "backend_node_groups", "cluster_zone_id", "duration", "frontend_node_groups", "observer_node_groups", "oss_accessing_role_name", "pricing_cycle", "promotion_option_no"},
			},
		},
	})
}

var AlicloudStarRocksInstanceMap11073 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
	"region_id":   CHECKSET,
}

func AlicloudStarRocksInstanceBasicDependence11073(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_vpc" "defaultB21JUD" {
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "default106DkE" {
  vpc_id       = alicloud_vpc.defaultB21JUD.id
  cidr_block   = "172.16.1.0/24"
  vswitch_name = "sr-test"
  zone_id      = "cn-hangzhou-i"
}


`, name)
}

// Case Instance_online 11101
func TestAccAliCloudStarRocksInstance_basic11101(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_star_rocks_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudStarRocksInstanceMap11101)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &StarRocksServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeStarRocksInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccstarrocks%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudStarRocksInstanceBasicDependence11101)
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
					"instance_name": name,
					"auto_renew":    "false",
					"frontend_node_groups": []map[string]interface{}{
						{
							"cu":                          "8",
							"storage_size":                "100",
							"resident_node_number":        "3",
							"storage_performance_level":   "pl1",
							"spec_type":                   "standard",
							"disk_number":                 "1",
							"zone_id":                     "cn-hangzhou-i",
							"local_storage_instance_type": "null",
						},
					},
					"vswitches": []map[string]interface{}{
						{
							"vswitch_id": "${alicloud_vswitch.default106DkE.id}",
							"zone_id":    "cn-hangzhou-i",
						},
					},
					"backend_node_groups": []map[string]interface{}{
						{
							"cu":                          "8",
							"storage_size":                "100",
							"resident_node_number":        "3",
							"disk_number":                 "1",
							"storage_performance_level":   "pl1",
							"spec_type":                   "standard",
							"zone_id":                     "cn-hangzhou-i",
							"local_storage_instance_type": "null",
						},
					},
					"cluster_zone_id":         "cn-hangzhou-i",
					"duration":                "1",
					"pay_type":                "postPaid",
					"vpc_id":                  "${alicloud_vpc.defaultB21JUD.id}",
					"version":                 "3.3",
					"run_mode":                "shared_data",
					"package_type":            "official",
					"admin_password":          "1qaz@QAZ",
					"oss_accessing_role_name": "AliyunEMRStarRocksAccessingOSSRole",
					"pricing_cycle":           "Month",
					"kms_key_id":              "123",
					"promotion_option_no":     "123",
					"encrypted":               "false",
					"resource_group_id":       "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"observer_node_groups": []map[string]interface{}{
						{
							"cu":                          "8",
							"storage_size":                "100",
							"storage_performance_level":   "pl1",
							"disk_number":                 "1",
							"resident_node_number":        "1",
							"spec_type":                   "standard",
							"local_storage_instance_type": "null",
							"zone_id":                     "cn-hangzhou-h",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name":           name,
						"auto_renew":              "false",
						"frontend_node_groups.#":  "1",
						"vswitches.#":             "1",
						"backend_node_groups.#":   "1",
						"cluster_zone_id":         "cn-hangzhou-i",
						"duration":                "1",
						"pay_type":                "postPaid",
						"vpc_id":                  CHECKSET,
						"version":                 CHECKSET,
						"run_mode":                "shared_data",
						"package_type":            "official",
						"admin_password":          "1qaz@QAZ",
						"oss_accessing_role_name": "AliyunEMRStarRocksAccessingOSSRole",
						"pricing_cycle":           "Month",
						"kms_key_id":              CHECKSET,
						"promotion_option_no":     CHECKSET,
						"encrypted":               "false",
						"resource_group_id":       CHECKSET,
						"observer_node_groups.#":  "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
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
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
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
				ImportStateVerifyIgnore: []string{"admin_password", "auto_renew", "backend_node_groups", "cluster_zone_id", "duration", "frontend_node_groups", "observer_node_groups", "oss_accessing_role_name", "pricing_cycle", "promotion_option_no"},
			},
		},
	})
}

var AlicloudStarRocksInstanceMap11101 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
	"region_id":   CHECKSET,
}

func AlicloudStarRocksInstanceBasicDependence11101(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_vpc" "defaultB21JUD" {
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "default106DkE" {
  vpc_id       = alicloud_vpc.defaultB21JUD.id
  cidr_block   = "172.16.1.0/24"
  vswitch_name = "sr-test"
  zone_id      = "cn-hangzhou-i"
}


`, name)
}

// Case instance_online 11093
func TestAccAliCloudStarRocksInstance_basic11093(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_star_rocks_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudStarRocksInstanceMap11093)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &StarRocksServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeStarRocksInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccstarrocks%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudStarRocksInstanceBasicDependence11093)
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
					"instance_name": name,
					"auto_renew":    "false",
					"frontend_node_groups": []map[string]interface{}{
						{
							"cu":                          "8",
							"storage_size":                "100",
							"resident_node_number":        "3",
							"storage_performance_level":   "pl1",
							"spec_type":                   "standard",
							"disk_number":                 "1",
							"zone_id":                     "cn-hangzhou-i",
							"local_storage_instance_type": "null",
						},
					},
					"vswitches": []map[string]interface{}{
						{
							"vswitch_id": "${alicloud_vswitch.default106DkE.id}",
							"zone_id":    "cn-hangzhou-i",
						},
					},
					"backend_node_groups": []map[string]interface{}{
						{
							"cu":                          "8",
							"storage_size":                "100",
							"resident_node_number":        "3",
							"disk_number":                 "1",
							"storage_performance_level":   "pl1",
							"spec_type":                   "standard",
							"zone_id":                     "cn-hangzhou-i",
							"local_storage_instance_type": "null",
						},
					},
					"cluster_zone_id":         "cn-hangzhou-i",
					"duration":                "1",
					"pay_type":                "postPaid",
					"vpc_id":                  "${alicloud_vpc.defaultB21JUD.id}",
					"version":                 "3.3",
					"run_mode":                "shared_data",
					"package_type":            "official",
					"admin_password":          "1qaz@QAZ",
					"oss_accessing_role_name": "AliyunEMRStarRocksAccessingOSSRole",
					"pricing_cycle":           "Month",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name":           name,
						"auto_renew":              "false",
						"frontend_node_groups.#":  "1",
						"vswitches.#":             "1",
						"backend_node_groups.#":   "1",
						"cluster_zone_id":         "cn-hangzhou-i",
						"duration":                "1",
						"pay_type":                "postPaid",
						"vpc_id":                  CHECKSET,
						"version":                 CHECKSET,
						"run_mode":                "shared_data",
						"package_type":            "official",
						"admin_password":          "1qaz@QAZ",
						"oss_accessing_role_name": "AliyunEMRStarRocksAccessingOSSRole",
						"pricing_cycle":           "Month",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name": name + "_update",
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
				ImportStateVerifyIgnore: []string{"admin_password", "auto_renew", "backend_node_groups", "cluster_zone_id", "duration", "frontend_node_groups", "observer_node_groups", "oss_accessing_role_name", "pricing_cycle", "promotion_option_no"},
			},
		},
	})
}

var AlicloudStarRocksInstanceMap11093 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
	"region_id":   CHECKSET,
}

func AlicloudStarRocksInstanceBasicDependence11093(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_vpc" "defaultB21JUD" {
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "default106DkE" {
  vpc_id       = alicloud_vpc.defaultB21JUD.id
  cidr_block   = "172.16.1.0/24"
  vswitch_name = "sr-test"
  zone_id      = "cn-hangzhou-i"
}


`, name)
}

// Test StarRocks Instance. <<< Resource test cases, automatically generated.
