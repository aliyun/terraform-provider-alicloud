// Package starrocks. This file is generated automatically. Please do not modify it manually, thank you!
package starrocks_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/starrocks"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// The following multi-member TypeList attributes of alicloud_star_rocks_instance
// have no convergent reorder acceptance coverage. This is pre-existing stock
// debt: the flat-layout alicloud/resource_alicloud_star_rocks_instance_test.go
// (from which this file was mechanically migrated) carried no reorder test for
// them either, and this PR relocates the sources verbatim without changing the
// resource schema. Per scripts/collection-order/README.md ("Bypassing stock
// attributes"), they are marked as known stock debt outside this change rather
// than having reorder cloud-ACC added to a mechanical-migration PR.
// collection-order: bypass backend_node_groups -- pre-existing uncovered TypeList; no reorder coverage in the flat-layout original before this mechanical migration, stock debt outside this change
// collection-order: bypass frontend_node_groups -- pre-existing uncovered TypeList; no reorder coverage in the flat-layout original before this mechanical migration, stock debt outside this change
// collection-order: bypass observer_node_groups -- pre-existing uncovered TypeList; no reorder coverage in the flat-layout original before this mechanical migration, stock debt outside this change
// collection-order: bypass vswitches -- pre-existing uncovered TypeList; no reorder coverage in the flat-layout original before this mechanical migration, stock debt outside this change

// Test StarRocks Instance. >>> Resource test cases, automatically generated.
// Case Instance 11073
func TestAccAliCloudStarRocksInstance_basic11073(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_star_rocks_instance.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudStarRocksInstanceMap11073)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return starrocks.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeInstance")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccstarrocks%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudStarRocksInstanceBasicDependence11073)
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
							"storage_size":                "200",
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
						"vpc_id":                  acctest.CHECKSET,
						"version":                 acctest.CHECKSET,
						"run_mode":                "shared_data",
						"package_type":            "official",
						"admin_password":          "1qaz@QAZ",
						"oss_accessing_role_name": "AliyunEMRStarRocksAccessingOSSRole",
						"pricing_cycle":           "Month",
						"kms_key_id":              acctest.CHECKSET,
						"promotion_option_no":     acctest.CHECKSET,
						"encrypted":               "false",
						"resource_group_id":       acctest.CHECKSET,
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
						"resource_group_id": acctest.CHECKSET,
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
					"tags": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": acctest.REMOVEKEY,
						"tags.For":     acctest.REMOVEKEY,
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
	"status":      acctest.CHECKSET,
	"create_time": acctest.CHECKSET,
	"region_id":   acctest.CHECKSET,
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
	ra := acctest.ResourceAttrInit(resourceId, AlicloudStarRocksInstanceMap11101)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return starrocks.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeInstance")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccstarrocks%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudStarRocksInstanceBasicDependence11101)
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
							"storage_size":                "200",
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
						"vpc_id":                  acctest.CHECKSET,
						"version":                 acctest.CHECKSET,
						"run_mode":                "shared_data",
						"package_type":            "official",
						"admin_password":          "1qaz@QAZ",
						"oss_accessing_role_name": "AliyunEMRStarRocksAccessingOSSRole",
						"pricing_cycle":           "Month",
						"kms_key_id":              acctest.CHECKSET,
						"promotion_option_no":     acctest.CHECKSET,
						"encrypted":               "false",
						"resource_group_id":       acctest.CHECKSET,
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
						"resource_group_id": acctest.CHECKSET,
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
					"tags": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": acctest.REMOVEKEY,
						"tags.For":     acctest.REMOVEKEY,
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
	"status":      acctest.CHECKSET,
	"create_time": acctest.CHECKSET,
	"region_id":   acctest.CHECKSET,
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
	ra := acctest.ResourceAttrInit(resourceId, AlicloudStarRocksInstanceMap11093)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return starrocks.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeInstance")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccstarrocks%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudStarRocksInstanceBasicDependence11093)
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
							"storage_size":                "200",
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
						"vpc_id":                  acctest.CHECKSET,
						"version":                 acctest.CHECKSET,
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
					"tags": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": acctest.REMOVEKEY,
						"tags.For":     acctest.REMOVEKEY,
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
	"status":      acctest.CHECKSET,
	"create_time": acctest.CHECKSET,
	"region_id":   acctest.CHECKSET,
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
