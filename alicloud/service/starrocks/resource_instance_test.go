// Package starrocks. This file is generated automatically. Please do not modify it manually, thank you!
package starrocks_test

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/starrocks"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
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

// TestAccAliCloudStarRocksInstance_backendNodeGroupsStorageSize verifies that
// increasing backend_node_groups.*.storage_size expands the backend node group disks
// in place without touching the other compute groups, and that shrinking still
// replaces the instance.
func TestAccAliCloudStarRocksInstance_backendNodeGroupsStorageSize(t *testing.T) {
	var v map[string]interface{}
	var instanceId string
	var backendNodeGroupId string
	var otherComputeGroups map[string]string
	resourceId := "alicloud_star_rocks_instance.default"
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccStarRocks%d", rand)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return starrocks.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeInstance")
	rac := acctest.ResourceAttrCheckInit(rc, acctest.ResourceAttrInit(resourceId, map[string]string{}))
	checkInstanceNotReplaced := func(s *terraform.State) error {
		instance, ok := s.RootModule().Resources[resourceId]
		if !ok {
			return fmt.Errorf("resource %s not found in state", resourceId)
		}
		if instanceId == "" {
			instanceId = instance.Primary.ID
			return nil
		}
		if instance.Primary.ID != instanceId {
			return fmt.Errorf("instance was replaced during in-place disk expansion: %s -> %s", instanceId, instance.Primary.ID)
		}
		return nil
	}
	checkInstanceReplaced := func(s *terraform.State) error {
		instance, ok := s.RootModule().Resources[resourceId]
		if !ok {
			return fmt.Errorf("resource %s not found in state", resourceId)
		}
		if instanceId == "" {
			return fmt.Errorf("instance id was not captured before the shrink step")
		}
		if instance.Primary.ID == instanceId {
			return fmt.Errorf("shrinking storage_size did not replace the instance %s", instanceId)
		}
		// The replacement instance has its own node groups.
		instanceId = instance.Primary.ID
		backendNodeGroupId = ""
		return nil
	}
	logCleanupReferences := func(s *terraform.State) {
		vpcId := ""
		if resource, ok := s.RootModule().Resources["alicloud_vpc.defaultB21JUD"]; ok {
			vpcId = resource.Primary.ID
		}
		vSwitchId := ""
		if resource, ok := s.RootModule().Resources["alicloud_vswitch.default106DkE"]; ok {
			vSwitchId = resource.Primary.ID
		}
		t.Logf("cleanup reference: instance=%s vpc=%s vswitch=%s", instanceId, vpcId, vSwitchId)
	}
	// captureGroups records every compute group so the next step can prove that only
	// the targeted backend group was resized.
	captureGroups := func(s *terraform.State) error {
		instance, ok := s.RootModule().Resources[resourceId]
		if !ok {
			return fmt.Errorf("resource %s not found in state", resourceId)
		}
		service := starrocks.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
		snapshot := make(map[string]string)
		for _, componentType := range []string{"FE", "BE", "CN"} {
			groups, err := service.DescribeNodeGroupsByComponentType(instance.Primary.ID, componentType, time.Minute)
			if err != nil {
				return err
			}
			for _, group := range groups {
				snapshot[componentType+"/"+fmt.Sprint(group["NodeGroupId"])] = fmt.Sprint(group["StorageSize"])
			}
		}
		otherComputeGroups = snapshot
		logCleanupReferences(s)
		return nil
	}
	checkTargetStorageSize := func(expected int) resource.TestCheckFunc {
		return func(s *terraform.State) error {
			instance, ok := s.RootModule().Resources[resourceId]
			if !ok {
				return fmt.Errorf("resource %s not found in state", resourceId)
			}
			service := starrocks.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
			deadline := time.Now().Add(5 * time.Minute)
			for {
				groups, err := service.DescribeBackendNodeGroups(instance.Primary.ID, time.Minute)
				if err != nil {
					return err
				}
				matches := configuredBackendGroups(groups)
				if len(matches) == 1 {
					size, _ := strconv.Atoi(fmt.Sprint(matches[0]["StorageSize"]))
					nodeGroupId := fmt.Sprint(matches[0]["NodeGroupId"])
					if size == expected {
						if backendNodeGroupId != "" && nodeGroupId != backendNodeGroupId {
							return fmt.Errorf("disk expansion hit node group %s instead of %s", nodeGroupId, backendNodeGroupId)
						}
						backendNodeGroupId = nodeGroupId
						return nil
					}
				}
				if time.Now().After(deadline) {
					return fmt.Errorf("configured backend node group did not reach %d GiB, groups: %v", expected, groups)
				}
				time.Sleep(10 * time.Second)
			}
		}
	}
	checkOtherComputeGroupsUnchanged := func() resource.TestCheckFunc {
		return func(s *terraform.State) error {
			instance, ok := s.RootModule().Resources[resourceId]
			if !ok {
				return fmt.Errorf("resource %s not found in state", resourceId)
			}
			if len(otherComputeGroups) == 0 {
				return fmt.Errorf("no compute group snapshot was captured before the resizing step")
			}
			service := starrocks.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
			current := make(map[string]string)
			for _, componentType := range []string{"FE", "BE", "CN"} {
				groups, err := service.DescribeNodeGroupsByComponentType(instance.Primary.ID, componentType, time.Minute)
				if err != nil {
					return err
				}
				for _, group := range groups {
					current[componentType+"/"+fmt.Sprint(group["NodeGroupId"])] = fmt.Sprint(group["StorageSize"])
				}
			}
			for key, size := range otherComputeGroups {
				if strings.HasSuffix(key, "/"+backendNodeGroupId) {
					continue
				}
				actual, ok := current[key]
				if !ok {
					return fmt.Errorf("compute group %s disappeared during the disk expansion", key)
				}
				if actual != size {
					return fmt.Errorf("compute group %s changed from %s to %s GiB during the disk expansion", key, size, actual)
				}
			}
			for key := range current {
				if _, ok := otherComputeGroups[key]; !ok && !strings.HasSuffix(key, "/"+backendNodeGroupId) {
					return fmt.Errorf("unexpected compute group %s appeared during the disk expansion", key)
				}
			}
			return nil
		}
	}
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: starRocksInstanceBackendNodeGroupsStorageSizeConfig(name, 200),
				Check: resource.ComposeTestCheckFunc(
					checkInstanceNotReplaced,
					captureGroups,
					resource.TestCheckResourceAttr(resourceId, "backend_node_groups.#", "1"),
					resource.TestCheckResourceAttr(resourceId, "backend_node_groups.0.storage_size", "200"),
					checkTargetStorageSize(200),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"admin_password", "auto_renew", "backend_node_groups", "cluster_zone_id", "duration", "frontend_node_groups", "observer_node_groups", "oss_accessing_role_name", "pricing_cycle", "promotion_option_no"},
			},
			{
				Config: starRocksInstanceBackendNodeGroupsStorageSizeConfig(name, 300),
				Check: resource.ComposeTestCheckFunc(
					checkInstanceNotReplaced,
					resource.TestCheckResourceAttr(resourceId, "backend_node_groups.#", "1"),
					resource.TestCheckResourceAttr(resourceId, "backend_node_groups.0.storage_size", "300"),
					checkTargetStorageSize(300),
					checkOtherComputeGroupsUnchanged(),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"admin_password", "auto_renew", "backend_node_groups", "cluster_zone_id", "duration", "frontend_node_groups", "observer_node_groups", "oss_accessing_role_name", "pricing_cycle", "promotion_option_no"},
			},
			{
				// Shrinking is not supported in place, so the plan must replace the
				// instance and the new instance gets a new id. 200 GiB is the smallest
				// backend disk size the API accepts, so the shrink target cannot be lower.
				Config: starRocksInstanceBackendNodeGroupsStorageSizeConfig(name, 200),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceId, plancheck.ResourceActionReplace),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					checkInstanceReplaced,
					resource.TestCheckResourceAttr(resourceId, "backend_node_groups.0.storage_size", "200"),
					checkTargetStorageSize(200),
					func(s *terraform.State) error {
						logCleanupReferences(s)
						return nil
					},
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

// configuredBackendGroups filters the backend groups down to the group the test
// configuration describes. The zone is not part of the filter because the API does
// not report a reliable zone for ECI based compute groups.
func configuredBackendGroups(groups []map[string]interface{}) []map[string]interface{} {
	matches := make([]map[string]interface{}, 0)
	for _, group := range groups {
		if fmt.Sprint(group["SpecType"]) != "standard" ||
			fmt.Sprint(group["Cu"]) != "8" ||
			fmt.Sprint(group["DiskNumber"]) != "1" ||
			fmt.Sprint(group["ResidentNodeNumber"]) != "3" ||
			fmt.Sprint(group["StoragePerformanceLevel"]) != "pl1" {
			continue
		}
		matches = append(matches, group)
	}
	return matches
}

func starRocksInstanceBackendNodeGroupsStorageSizeConfig(name string, storageSize int) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_vpc" "defaultB21JUD" {
  cidr_block = "172.16.0.0/12"
  vpc_name   = "tf-testAccStarRocks-vpc"
}

resource "alicloud_vswitch" "default106DkE" {
  vpc_id       = alicloud_vpc.defaultB21JUD.id
  cidr_block   = "172.16.1.0/24"
  vswitch_name = "tf-testAccStarRocks-vswitch"
  zone_id      = "cn-hangzhou-i"
}

resource "alicloud_star_rocks_instance" "default" {
  admin_password          = "1qaz@QAZ"
  cluster_zone_id         = "cn-hangzhou-i"
  instance_name           = var.name
  duration                = 1
  pay_type                = "postPaid"
  vpc_id                  = alicloud_vpc.defaultB21JUD.id
  version                 = "3.3"
  run_mode                = "shared_data"
  package_type            = "official"
  resource_group_id       = data.alicloud_resource_manager_resource_groups.default.ids.0
  oss_accessing_role_name = "AliyunEMRStarRocksAccessingOSSRole"
  pricing_cycle           = "Month"

  vswitches {
    vswitch_id = alicloud_vswitch.default106DkE.id
    zone_id    = "cn-hangzhou-i"
  }

  frontend_node_groups {
    cu                        = 8
    storage_size              = 100
    storage_performance_level = "pl1"
    disk_number               = 1
    resident_node_number      = 3
    spec_type                 = "standard"
    zone_id                   = "cn-hangzhou-i"
  }

  backend_node_groups {
    cu                        = 8
    storage_size              = %d
    storage_performance_level = "pl1"
    disk_number               = 1
    resident_node_number      = 3
    spec_type                 = "standard"
    zone_id                   = "cn-hangzhou-i"
  }
}
`, name, storageSize)
}

// Test StarRocks Instance. <<< Resource test cases, automatically generated.
