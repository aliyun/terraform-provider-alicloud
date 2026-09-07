// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test StarRocks NodeGroup. >>> Resource test cases, automatically generated.
// Case NodeGroup_副本1760322662795 11601
func TestAccAliCloudStarRocksNodeGroup_basic11601(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_star_rocks_node_group.default"
	ra := resourceAttrInit(resourceId, AlicloudStarRocksNodeGroupMap11601)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &StarRocksServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeStarRocksNodeGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tfacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudStarRocksNodeGroupBasicDependence11601)
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
					"description":                 "test_desc",
					"node_group_name":             name,
					"instance_id":                 "${alicloud_star_rocks_instance.defaultvjnpM0.id}",
					"spec_type":                   "standard",
					"storage_performance_level":   "pl1",
					"pricing_cycle":               "1",
					"auto_renew":                  "false",
					"storage_size":                "200",
					"duration":                    "1",
					"pay_type":                    "postPaid",
					"cu":                          "8",
					"disk_number":                 "1",
					"resident_node_number":        "1",
					"local_storage_instance_type": "non_local_storage",
					"promotion_option_no":         "blank",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":                 "test_desc",
						"node_group_name":             name,
						"instance_id":                 CHECKSET,
						"spec_type":                   "standard",
						"storage_performance_level":   "pl1",
						"pricing_cycle":               CHECKSET,
						"auto_renew":                  "false",
						"storage_size":                "200",
						"duration":                    "1",
						"pay_type":                    "postPaid",
						"cu":                          "8",
						"disk_number":                 "1",
						"resident_node_number":        "1",
						"local_storage_instance_type": "non_local_storage",
						"promotion_option_no":         "blank",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cu":                  "16",
					"promotion_option_no": "2345",
					"fast_mode":           "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cu":                  "16",
						"promotion_option_no": CHECKSET,
						"fast_mode":           "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"storage_size": "500",
					"spec_type":    "ramEnhanced",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"storage_size": "500",
						"spec_type":    "ramEnhanced",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"storage_performance_level": "pl2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"storage_performance_level": "pl2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"disk_number": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"disk_number": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resident_node_number": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resident_node_number": "2",
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
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_renew", "duration", "fast_mode", "pricing_cycle", "promotion_option_no"},
			},
		},
	})
}

var AlicloudStarRocksNodeGroupMap11601 = map[string]string{
	"status":        CHECKSET,
	"create_time":   CHECKSET,
	"node_group_id": CHECKSET,
	"region_id":     CHECKSET,
}

func AlicloudStarRocksNodeGroupBasicDependence11601(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_vpc" "defaultq6pcFe" {
  cidr_block = "172.16.0.0/12"
  vpc_name   = "test-vpc-511"
}

resource "alicloud_vswitch" "defaultujlpyG" {
  vpc_id       = alicloud_vpc.defaultq6pcFe.id
  zone_id      = "cn-hangzhou-i"
  cidr_block   = "172.16.0.0/24"
  vswitch_name = "sr-test-ng"
}

resource "alicloud_star_rocks_instance" "defaultvjnpM0" {
  cluster_zone_id = "cn-hangzhou-i"
  encrypted       = false
  auto_renew      = false
  pay_type        = "postPaid"
  frontend_node_groups {
    cu                        = "8"
    storage_size              = "100"
    storage_performance_level = "pl1"
    disk_number               = "1"
    zone_id                   = "cn-hangzhou-i"
    spec_type                 = "standard"
    resident_node_number      = "1"
  }
  instance_name = "t1"
  vswitches {
    zone_id    = "cn-hangzhou-i"
    vswitch_id = alicloud_vswitch.defaultujlpyG.id
  }
  vpc_id                  = alicloud_vpc.defaultq6pcFe.id
  version                 = "3.3"
  run_mode                = "shared_data"
  package_type            = "official"
  oss_accessing_role_name = "AliyunEMRStarRocksAccessingOSSRolecn"
  admin_password          = "1qaz@QAZ"
  backend_node_groups {
    cu                        = "8"
    storage_size              = "200"
    zone_id                   = "cn-hangzhou-i"
    spec_type                 = "standard"
    resident_node_number      = "3"
    disk_number               = "1"
    storage_performance_level = "pl1"
  }
}


`, name)
}

// Case NodeGroup 11074
func TestAccAliCloudStarRocksNodeGroup_basic11074(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_star_rocks_node_group.default"
	ra := resourceAttrInit(resourceId, AlicloudStarRocksNodeGroupMap11074)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &StarRocksServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeStarRocksNodeGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tfacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudStarRocksNodeGroupBasicDependence11074)
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
					"description":                 "test_desc",
					"node_group_name":             name,
					"instance_id":                 "${alicloud_star_rocks_instance.defaultvjnpM0.id}",
					"spec_type":                   "standard",
					"storage_performance_level":   "pl1",
					"pricing_cycle":               "1",
					"auto_renew":                  "false",
					"storage_size":                "200",
					"duration":                    "1",
					"pay_type":                    "postPaid",
					"cu":                          "8",
					"disk_number":                 "1",
					"resident_node_number":        "1",
					"local_storage_instance_type": "non_local_storage",
					"promotion_option_no":         "blank",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":                 "test_desc",
						"node_group_name":             name,
						"instance_id":                 CHECKSET,
						"spec_type":                   "standard",
						"storage_performance_level":   "pl1",
						"pricing_cycle":               CHECKSET,
						"auto_renew":                  "false",
						"storage_size":                "200",
						"duration":                    "1",
						"pay_type":                    "postPaid",
						"cu":                          "8",
						"disk_number":                 "1",
						"resident_node_number":        "1",
						"local_storage_instance_type": "non_local_storage",
						"promotion_option_no":         "blank",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cu":                  "16",
					"promotion_option_no": "2345",
					"fast_mode":           "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cu":                  "16",
						"promotion_option_no": CHECKSET,
						"fast_mode":           "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"storage_size": "500",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"storage_size": "500",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"storage_performance_level": "pl2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"storage_performance_level": "pl2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"disk_number": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"disk_number": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resident_node_number": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resident_node_number": "2",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_renew", "duration", "fast_mode", "pricing_cycle", "promotion_option_no"},
			},
		},
	})
}

var AlicloudStarRocksNodeGroupMap11074 = map[string]string{
	"status":        CHECKSET,
	"create_time":   CHECKSET,
	"node_group_id": CHECKSET,
	"region_id":     CHECKSET,
}

func AlicloudStarRocksNodeGroupBasicDependence11074(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_vpc" "defaultq6pcFe" {
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "defaultujlpyG" {
  vpc_id       = alicloud_vpc.defaultq6pcFe.id
  zone_id      = "cn-hangzhou-i"
  cidr_block   = "172.16.0.0/24"
  vswitch_name = "sr-test-ng"
}

resource "alicloud_star_rocks_instance" "defaultvjnpM0" {
  cluster_zone_id = "cn-hangzhou-i"
  encrypted       = false
  auto_renew      = false
  pay_type        = "postPaid"
  frontend_node_groups {
    cu                        = "8"
    storage_size              = "100"
    storage_performance_level = "pl1"
    disk_number               = "1"
    zone_id                   = "cn-hangzhou-i"
    spec_type                 = "standard"
    resident_node_number      = "1"
  }
  instance_name = "t1"
  vswitches {
    zone_id    = "cn-hangzhou-i"
    vswitch_id = alicloud_vswitch.defaultujlpyG.id
  }
  vpc_id                  = alicloud_vpc.defaultq6pcFe.id
  version                 = "3.3"
  run_mode                = "shared_data"
  package_type            = "official"
  oss_accessing_role_name = "AliyunEMRStarRocksAccessingOSSRole"
  admin_password          = "1qaz@QAZ"
  backend_node_groups {
    cu                        = "8"
    storage_size              = "200"
    zone_id                   = "cn-hangzhou-i"
    spec_type                 = "standard"
    resident_node_number      = "3"
    disk_number               = "1"
    storage_performance_level = "pl1"
  }
}


`, name)
}

// Test StarRocks NodeGroup. <<< Resource test cases, automatically generated.
