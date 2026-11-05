// Package starrocks. This file is generated automatically. Please do not modify it manually, thank you!
package starrocks_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/starrocks"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdktf "github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
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
			groups, err := service.DescribeNodeGroupsByComponentType(instance.Primary.ID, componentType, time.Now().Add(time.Minute))
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
				groups, err := service.DescribeBackendNodeGroups(instance.Primary.ID, time.Now().Add(time.Minute))
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
				groups, err := service.DescribeNodeGroupsByComponentType(instance.Primary.ID, componentType, time.Now().Add(time.Minute))
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

// --- White-box style checks driven through the public resource entry point. ---

type starRocksNodeGroupsHandler struct {
	t      *testing.T
	groups []map[string]interface{}
	total  int
	// pageErrors injects a failure for one page number (1-based).
	pageErrors map[int]int
	// alwaysFullPage keeps returning a full page without Total to exercise the guard.
	alwaysFullPage bool
	// convergeAfter keeps reporting the old size for this many reads after a resize.
	convergeAfter int
	convergeReads int
	pendingSize   map[string]int
	// delayedAction delays one action, e.g. to exhaust the update budget.
	delayedAction   string
	delayedStatus   int
	delay           time.Duration
	modifyResponses []int
	modifyCalls     int
	actions         []string
	modifyQueries   []map[string]string
	listRequests    []map[string]string
	byIdRequests    []map[string]string
	instanceStatus  string
}

func (h *starRocksNodeGroupsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	action := r.URL.Path
	h.actions = append(h.actions, action)
	if action == h.delayedAction && h.delay > 0 {
		time.Sleep(h.delay)
		if h.delayedStatus != 0 {
			w.WriteHeader(h.delayedStatus)
			writeTestJSON(w, map[string]interface{}{"Code": "InternalError", "Message": "code: 500, injected slow failure", "RequestId": "rid"})
			return
		}
	}
	switch action {
	case "/webapi/nodegroup/describeNodeGroups":
		body := map[string]interface{}{}
		if raw, err := io.ReadAll(r.Body); err == nil && len(raw) > 0 {
			_ = json.Unmarshal(raw, &body)
		}
		if ids, byId := body["nodeGroupIds"]; byId {
			list := ids.([]interface{})
			nodeGroupId := fmt.Sprint(list[0])
			h.byIdRequests = append(h.byIdRequests, map[string]string{
				"ClusterId": query.Get("ClusterId"), "nodeGroupIds": nodeGroupId,
			})
			item := h.groupById(nodeGroupId)
			item["Status"] = "RUNNING"
			if h.instanceStatus != "" {
				item["Status"] = h.instanceStatus
			}
			if target, pending := h.pendingSize[nodeGroupId]; pending {
				if h.convergeReads <= 1 {
					item["StorageSize"] = float64(target)
					delete(h.pendingSize, nodeGroupId)
				} else {
					h.convergeReads--
				}
			}
			writeTestJSON(w, map[string]interface{}{"Data": []interface{}{item}})
			return
		}
		pageNumber := 1
		pageSize := 100
		fmt.Sscan(query.Get("PageNumber"), &pageNumber)
		fmt.Sscan(query.Get("PageSize"), &pageSize)
		componentType := fmt.Sprint(body["componentType"])
		h.listRequests = append(h.listRequests, map[string]string{
			"ClusterId": query.Get("ClusterId"), "PageNumber": query.Get("PageNumber"),
			"PageSize": query.Get("PageSize"), "componentType": componentType,
		})
		if status, ok := h.pageErrors[pageNumber]; ok {
			w.WriteHeader(status)
			writeTestJSON(w, map[string]interface{}{"Code": "InternalError", "Message": "code: 500, injected page failure", "RequestId": "rid"})
			return
		}
		if h.alwaysFullPage {
			page := make([]interface{}, 0, pageSize)
			for i := 0; i < pageSize; i++ {
				page = append(page, testBackendGroup(fmt.Sprintf("ng-loop-%d-%03d", pageNumber, i), "cn-hangzhou-h"))
			}
			writeTestJSON(w, map[string]interface{}{"Data": page})
			return
		}
		filtered := make([]map[string]interface{}, 0, len(h.groups))
		for _, group := range h.groups {
			if componentType != "" && fmt.Sprint(group["ComponentType"]) != componentType {
				continue
			}
			filtered = append(filtered, group)
		}
		start := (pageNumber - 1) * pageSize
		end := start + pageSize
		if start > len(filtered) {
			start = len(filtered)
		}
		if end > len(filtered) {
			end = len(filtered)
		}
		page := make([]interface{}, 0, end-start)
		for _, group := range filtered[start:end] {
			page = append(page, group)
		}
		total := h.total
		if total > len(filtered) {
			total = len(filtered)
		}
		writeTestJSON(w, map[string]interface{}{"Data": page, "Total": total})
	case "/webapi/resourceChange/modifyDiskSize":
		h.modifyCalls++
		h.modifyQueries = append(h.modifyQueries, map[string]string{
			"InstanceId": query.Get("InstanceId"), "NodeGroupId": query.Get("NodeGroupId"), "Target": query.Get("Target"),
		})
		target := 0
		fmt.Sscan(query.Get("Target"), &target)
		if h.convergeAfter > 0 {
			if h.pendingSize == nil {
				h.pendingSize = make(map[string]int)
			}
			h.pendingSize[query.Get("NodeGroupId")] = target
			h.convergeReads = h.convergeAfter
		} else {
			for _, group := range h.groups {
				if fmt.Sprint(group["NodeGroupId"]) == query.Get("NodeGroupId") {
					group["StorageSize"] = float64(target)
				}
			}
		}
		status := http.StatusOK
		if len(h.modifyResponses) > 0 {
			status = h.modifyResponses[0]
			h.modifyResponses = h.modifyResponses[1:]
		}
		if status != http.StatusOK {
			w.WriteHeader(status)
			writeTestJSON(w, map[string]interface{}{"Code": "InternalError", "Message": "code: 500, injected failure", "RequestId": "rid"})
			return
		}
		writeTestJSON(w, map[string]interface{}{"Success": true})
	case "/webapi/starrocks/describeInstances":
		writeTestJSON(w, map[string]interface{}{"Data": []interface{}{map[string]interface{}{
			"InstanceId": "c-fixture", "InstanceStatus": "running", "InstanceName": "tf-testacc-starrocks",
			"Tags": []interface{}{}, "VSwitches": []interface{}{},
		}}})
	default:
		writeTestJSON(w, map[string]interface{}{"Success": true})
	}
}

func (h *starRocksNodeGroupsHandler) groupById(id string) map[string]interface{} {
	for _, group := range h.groups {
		if fmt.Sprint(group["NodeGroupId"]) == id {
			return group
		}
	}
	return testBackendGroup(id, "")
}

func (h *starRocksNodeGroupsHandler) countAction(action string) int {
	count := 0
	for _, item := range h.actions {
		if item == action {
			count++
		}
	}
	return count
}

func writeTestJSON(w http.ResponseWriter, body map[string]interface{}) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(body)
}

func starRocksTestClient(t *testing.T, handler http.HandlerFunc) *connectivity.AliyunClient {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)
	credential, err := credentials.NewCredential(new(credentials.Config).
		SetType("access_key").SetAccessKeyId("test-key").SetAccessKeySecret("test-secret"))
	if err != nil {
		t.Fatalf("build credential: %s", err)
	}
	endpoints := new(sync.Map)
	endpoint := strings.TrimPrefix(server.URL, "http://")
	t.Setenv("NO_PROXY", endpoint)
	config := &connectivity.Config{
		AccessKey: "test-key", SecretKey: "test-secret", Credential: credential,
		RegionId: "cn-hangzhou", AccountType: "test", Protocol: "http",
		Endpoints: endpoints, SignVersion: new(sync.Map), SkipRegionValidation: true,
	}
	client, err := config.Client()
	if err != nil {
		t.Fatalf("build client: %s", err)
	}
	endpoints.Store("starrocks", endpoint)
	return client
}

func testBackendGroup(id, zone string) map[string]interface{} {
	return map[string]interface{}{
		"NodeGroupId": id, "ComponentType": "CN", "ZoneId": zone,
		"Cu": float64(8), "DiskNumber": float64(1), "ResidentNodeNumber": float64(3),
		"SpecType": "standard", "StoragePerformanceLevel": "pl1",
		"LocalStorageInstanceType": "", "StorageSize": float64(200), "Status": "RUNNING",
	}
}

func backendGroupStateAttributes(size string) map[string]string {
	return map[string]string{
		"backend_node_groups.#":                             "1",
		"backend_node_groups.0.cu":                          "8",
		"backend_node_groups.0.zone_id":                     "cn-hangzhou-i",
		"backend_node_groups.0.spec_type":                   "standard",
		"backend_node_groups.0.disk_number":                 "1",
		"backend_node_groups.0.storage_performance_level":   "pl1",
		"backend_node_groups.0.resident_node_number":        "3",
		"backend_node_groups.0.storage_size":                size,
		"backend_node_groups.0.local_storage_instance_type": "null",
		"instance_name":                                     "tf-testacc-starrocks",
	}
}

func twoBackendGroupStateAttributes(sizeA, sizeB string) map[string]string {
	return map[string]string{
		"backend_node_groups.#":                             "2",
		"backend_node_groups.0.cu":                          "8",
		"backend_node_groups.0.zone_id":                     "cn-hangzhou-i",
		"backend_node_groups.0.spec_type":                   "standard",
		"backend_node_groups.0.disk_number":                 "1",
		"backend_node_groups.0.storage_performance_level":   "pl1",
		"backend_node_groups.0.resident_node_number":        "3",
		"backend_node_groups.0.storage_size":                sizeA,
		"backend_node_groups.0.local_storage_instance_type": "null",
		"backend_node_groups.1.cu":                          "16",
		"backend_node_groups.1.zone_id":                     "cn-hangzhou-h",
		"backend_node_groups.1.spec_type":                   "standard",
		"backend_node_groups.1.disk_number":                 "1",
		"backend_node_groups.1.storage_performance_level":   "pl1",
		"backend_node_groups.1.resident_node_number":        "3",
		"backend_node_groups.1.storage_size":                sizeB,
		"backend_node_groups.1.local_storage_instance_type": "null",
		"instance_name":                                     "tf-testacc-starrocks",
	}
}

func backendGroupConfig(storageSize interface{}, overrides map[string]interface{}) map[string]interface{} {
	group := map[string]interface{}{
		"cu": 8, "zone_id": "cn-hangzhou-i", "spec_type": "standard", "disk_number": 1,
		"storage_performance_level": "pl1", "resident_node_number": 3, "local_storage_instance_type": "null",
	}
	if storageSize != nil {
		group["storage_size"] = storageSize
	}
	for key, value := range overrides {
		group[key] = value
	}
	return map[string]interface{}{
		"admin_password": "1qaz@QAZ", "cluster_zone_id": "cn-hangzhou-i",
		"instance_name": "tf-testacc-starrocks", "backend_node_groups": []interface{}{group},
	}
}

// applyStarRocksInstance drives the update through the public resource entry point.
func applyStarRocksInstance(t *testing.T, client *connectivity.AliyunClient, attributes map[string]string, config map[string]interface{}, updateTimeout time.Duration) error {
	t.Helper()
	resource := starrocks.ResourceInstance()
	state := &sdktf.InstanceState{ID: "c-fixture", Attributes: attributes}
	diff, err := resource.Diff(context.Background(), state, sdktf.NewResourceConfigRaw(config), nil)
	if err != nil {
		t.Fatalf("diff: %s", err)
	}
	if updateTimeout > 0 {
		// Apply reads the resource timeouts from the diff or state meta, which is how a
		// real plan carries the practitioner's timeouts block.
		diff.Meta = map[string]interface{}{schema.TimeoutKey: map[string]interface{}{"update": updateTimeout.String()}}
	}
	_, diags := resource.Apply(context.Background(), state, diff, client)
	if diags.HasError() {
		return fmt.Errorf("%v", diags)
	}
	return nil
}

func TestUnitDescribeNodeGroupsPagination(t *testing.T) {
	groups := make([]map[string]interface{}, 0, 150)
	for i := 0; i < 149; i++ {
		filler := testBackendGroup(fmt.Sprintf("ng-filler-%03d", i), "cn-hangzhou-h")
		filler["Cu"] = float64(4)
		groups = append(groups, filler)
	}
	groups = append(groups, testBackendGroup("ng-last-page", "cn-hangzhou-i"))
	handler := &starRocksNodeGroupsHandler{t: t, groups: groups, total: len(groups)}
	client := starRocksTestClient(t, handler.ServeHTTP)

	listed, err := starrocks.NewService(client).DescribeNodeGroupsByComponentType("c-fixture", "CN", time.Now().Add(10*time.Second))
	if err != nil {
		t.Fatalf("list node groups: %s", err)
	}
	if len(listed) != len(groups) {
		t.Fatalf("expected %d groups across all pages, got %d", len(groups), len(listed))
	}
	if len(handler.listRequests) < 2 {
		t.Fatal("expected the listing to read a second page")
	}
	for _, req := range handler.listRequests {
		if req["ClusterId"] != "c-fixture" || req["componentType"] != "CN" || req["PageSize"] != "100" || req["PageNumber"] == "" {
			t.Fatalf("unexpected list request: %#v", req)
		}
	}
	found := false
	for _, group := range listed {
		if fmt.Sprint(group["NodeGroupId"]) == "ng-last-page" && fmt.Sprint(group["ComponentType"]) == "CN" {
			found = true
		}
	}
	if !found {
		t.Fatal("the group that only exists on the second page was not returned")
	}
}

func TestUnitBackendNodeGroupResizeRequest(t *testing.T) {
	handler := &starRocksNodeGroupsHandler{
		t: t, groups: []map[string]interface{}{testBackendGroup("ng-target", "cn-hangzhou-i")},
		total: 1, instanceStatus: "RUNNING",
	}
	client := starRocksTestClient(t, handler.ServeHTTP)
	config := backendGroupConfig(300, nil)
	if err := applyStarRocksInstance(t, client, backendGroupStateAttributes("200"), config, 0); err != nil {
		t.Fatalf("apply: %s", err)
	}
	if handler.modifyCalls != 1 {
		t.Fatalf("expected exactly one ModifyDiskSize call, got %d", handler.modifyCalls)
	}
	sent := handler.modifyQueries[0]
	if sent["InstanceId"] != "c-fixture" || sent["NodeGroupId"] != "ng-target" || sent["Target"] != "300" {
		t.Fatalf("unexpected ModifyDiskSize request: %#v", sent)
	}
	if size := fmt.Sprint(handler.groups[0]["StorageSize"]); size != "300" {
		t.Fatalf("expected the remote group to reach 300 GiB, got %s", size)
	}
	if len(handler.byIdRequests) == 0 || handler.byIdRequests[0]["nodeGroupIds"] != "ng-target" || handler.byIdRequests[0]["ClusterId"] != "c-fixture" {
		t.Fatalf("unexpected by-id describe request: %#v", handler.byIdRequests)
	}
	if handler.countAction("/webapi/cluster/update_name") != 0 {
		t.Fatal("the name must not be updated when only storage_size changes")
	}
}

func TestUnitBackendNodeGroupAmbiguityWritesNothing(t *testing.T) {
	handler := &starRocksNodeGroupsHandler{
		t: t,
		groups: []map[string]interface{}{
			testBackendGroup("ng-1", "cn-hangzhou-i"), testBackendGroup("ng-2", "cn-hangzhou-i"),
		},
		total: 2, instanceStatus: "RUNNING",
	}
	client := starRocksTestClient(t, handler.ServeHTTP)
	config := backendGroupConfig(300, nil)
	config["instance_name"] = "tf-testacc-starrocks-renamed"
	config["resource_group_id"] = "rg-new"
	config["tags"] = map[string]interface{}{"env": "test"}

	err := applyStarRocksInstance(t, client, backendGroupStateAttributes("200"), config, 0)
	if err == nil || !strings.Contains(err.Error(), "refusing to guess") {
		t.Fatalf("expected an ambiguity error, got %v", err)
	}
	if handler.modifyCalls != 0 || handler.countAction("/webapi/cluster/update_name") != 0 || handler.countAction("/webapi/resourceGroup/change") != 0 || handler.countAction("/webapi/tags") != 0 {
		t.Fatalf("no modifying call may be sent on ambiguity: %#v", handler.actions)
	}
}

func TestUnitBackendNodeGroupZoneNotUsedForMatching(t *testing.T) {
	withZone := testBackendGroup("ng-with-zone", "cn-hangzhou-i")
	withoutZone := testBackendGroup("ng-without-zone", "")
	handler := &starRocksNodeGroupsHandler{
		t: t, groups: []map[string]interface{}{withZone, withoutZone}, total: 2, instanceStatus: "RUNNING",
	}
	client := starRocksTestClient(t, handler.ServeHTTP)
	err := applyStarRocksInstance(t, client, backendGroupStateAttributes("200"), backendGroupConfig(300, nil), 0)
	if err == nil || !strings.Contains(err.Error(), "refusing to guess") {
		t.Fatalf("a same-attribute group without a zone must make the match ambiguous, got %v", err)
	}
	if handler.modifyCalls != 0 {
		t.Fatalf("no resize may be sent when the zone cannot separate the candidates, got %d", handler.modifyCalls)
	}
}

func TestUnitBackendNodeGroupLaterAmbiguityWritesNothing(t *testing.T) {
	groupA := testBackendGroup("ng-a", "cn-hangzhou-i")
	groupB1 := testBackendGroup("ng-b1", "cn-hangzhou-h")
	groupB1["Cu"] = float64(16)
	groupB2 := testBackendGroup("ng-b2", "cn-hangzhou-h")
	groupB2["Cu"] = float64(16)
	handler := &starRocksNodeGroupsHandler{
		t: t, groups: []map[string]interface{}{groupA, groupB1, groupB2}, total: 3, instanceStatus: "RUNNING",
	}
	client := starRocksTestClient(t, handler.ServeHTTP)
	config := backendGroupConfig(300, nil)
	second := config["backend_node_groups"].([]interface{})[0].(map[string]interface{})
	secondCopy := map[string]interface{}{}
	for key, value := range second {
		secondCopy[key] = value
	}
	secondCopy["cu"] = 16
	secondCopy["zone_id"] = "cn-hangzhou-h"
	config["backend_node_groups"] = []interface{}{second, secondCopy}

	err := applyStarRocksInstance(t, client, twoBackendGroupStateAttributes("200", "200"), config, 0)
	if err == nil || !strings.Contains(err.Error(), "refusing to guess") {
		t.Fatalf("expected an ambiguity error, got %v", err)
	}
	if handler.modifyCalls != 0 {
		t.Fatalf("the resolvable first group must not be resized before the ambiguity fails, got %d", handler.modifyCalls)
	}
}

func TestUnitBackendNodeGroupLostResponseDoesNotResend(t *testing.T) {
	handler := &starRocksNodeGroupsHandler{
		t: t, groups: []map[string]interface{}{testBackendGroup("ng-target", "cn-hangzhou-i")},
		total: 1, instanceStatus: "RUNNING", modifyResponses: []int{http.StatusInternalServerError},
	}
	client := starRocksTestClient(t, handler.ServeHTTP)
	if err := applyStarRocksInstance(t, client, backendGroupStateAttributes("200"), backendGroupConfig(300, nil), 0); err != nil {
		t.Fatalf("apply: %s", err)
	}
	if handler.modifyCalls != 1 {
		t.Fatalf("the lost response must not be resent, got %d writes", handler.modifyCalls)
	}
}

func TestUnitBackendNodeGroupUnknownResultDoesNotResend(t *testing.T) {
	previous := starrocks.BackendNodeGroupResizeConvergenceWindowForTest(150 * time.Millisecond)
	defer starrocks.BackendNodeGroupResizeConvergenceWindowForTest(previous)
	handler := &starRocksNodeGroupsHandler{
		t: t, groups: []map[string]interface{}{testBackendGroup("ng-target", "cn-hangzhou-i")},
		total: 1, instanceStatus: "RUNNING", modifyResponses: []int{http.StatusInternalServerError}, convergeAfter: 1000,
	}
	client := starRocksTestClient(t, handler.ServeHTTP)
	err := applyStarRocksInstance(t, client, backendGroupStateAttributes("200"), backendGroupConfig(300, nil), 0)
	if err == nil || !strings.Contains(err.Error(), "is unknown") {
		t.Fatalf("expected an unknown result error, got %v", err)
	}
	if handler.modifyCalls != 1 {
		t.Fatalf("the unknown result must not trigger another billable write, got %d", handler.modifyCalls)
	}
}

func TestUnitBackendNodeGroupSharedUpdateBudget(t *testing.T) {
	handler := &starRocksNodeGroupsHandler{
		t: t, groups: []map[string]interface{}{testBackendGroup("ng-target", "cn-hangzhou-i")},
		total: 1, instanceStatus: "RUNNING",
		delayedAction: "/webapi/cluster/update_name", delayedStatus: http.StatusInternalServerError, delay: 3 * time.Second,
	}
	client := starRocksTestClient(t, handler.ServeHTTP)
	config := backendGroupConfig(300, nil)
	config["instance_name"] = "tf-testacc-starrocks-slow"
	err := applyStarRocksInstance(t, client, backendGroupStateAttributes("200"), config, time.Second)
	if err == nil {
		t.Fatal("the slow name update must exhaust the configured update budget")
	}
	if handler.modifyCalls != 0 {
		t.Fatalf("no resize may run after the update budget is exhausted, got %d", handler.modifyCalls)
	}
}

func TestUnitResourceInstanceBackendNodeGroupPlan(t *testing.T) {
	resource := starrocks.ResourceInstance()
	state := func(cfg map[string]interface{}) *sdktf.InstanceState {
		data := schema.TestResourceDataRaw(t, resource.Schema, cfg)
		data.SetId("c-fixture")
		return data.State()
	}
	planDiff := func(t *testing.T, oldCfg, newCfg map[string]interface{}) *sdktf.InstanceDiff {
		t.Helper()
		diff, err := resource.Diff(context.Background(), state(oldCfg), sdktf.NewResourceConfigRaw(newCfg), nil)
		if err != nil {
			t.Fatalf("diff: %s", err)
		}
		return diff
	}
	diffRequiresNew := func(t *testing.T, oldCfg, newCfg map[string]interface{}) bool {
		diff := planDiff(t, oldCfg, newCfg)
		return diff != nil && diff.RequiresNew()
	}

	t.Run("growing is in place", func(t *testing.T) {
		diff := planDiff(t, backendGroupConfig(200, nil), backendGroupConfig(300, nil))
		if diff == nil || len(diff.Attributes) == 0 || diff.RequiresNew() {
			t.Fatal("growing storage_size must be an in-place change")
		}
	})
	t.Run("shrinking replaces", func(t *testing.T) {
		if !diffRequiresNew(t, backendGroupConfig(300, nil), backendGroupConfig(200, nil)) {
			t.Fatal("shrinking storage_size must replace the instance")
		}
	})
	t.Run("removing replaces", func(t *testing.T) {
		if !diffRequiresNew(t, backendGroupConfig(300, nil), backendGroupConfig(nil, nil)) {
			t.Fatal("removing storage_size must replace the instance")
		}
	})
	t.Run("setting a previously unset value replaces", func(t *testing.T) {
		if !diffRequiresNew(t, backendGroupConfig(nil, nil), backendGroupConfig(300, nil)) {
			t.Fatal("setting storage_size without a previous value must replace the instance")
		}
	})
	for _, field := range []struct {
		name  string
		value interface{}
	}{
		{"cu", 16}, {"spec_type", "localSSD"}, {"disk_number", 2}, {"resident_node_number", 6},
		{"storage_performance_level", "pl2"}, {"zone_id", "cn-hangzhou-h"}, {"local_storage_instance_type", "local_ssd"},
	} {
		field := field
		t.Run("changing "+field.name+" replaces", func(t *testing.T) {
			if !diffRequiresNew(t, backendGroupConfig(200, nil), backendGroupConfig(200, map[string]interface{}{field.name: field.value})) {
				t.Fatalf("changing %s must replace the instance", field.name)
			}
		})
	}
	t.Run("adding a backend group replaces", func(t *testing.T) {
		second := backendGroupConfig(200, map[string]interface{}{"zone_id": "cn-hangzhou-h"})["backend_node_groups"].([]interface{})[0]
		cfg := backendGroupConfig(200, nil)
		cfg["backend_node_groups"] = []interface{}{cfg["backend_node_groups"].([]interface{})[0], second}
		if !diffRequiresNew(t, backendGroupConfig(200, nil), cfg) {
			t.Fatal("adding a backend group must replace the instance")
		}
	})
	t.Run("removing a backend group replaces", func(t *testing.T) {
		second := backendGroupConfig(200, map[string]interface{}{"zone_id": "cn-hangzhou-h"})["backend_node_groups"].([]interface{})[0]
		oldCfg := backendGroupConfig(200, nil)
		oldCfg["backend_node_groups"] = []interface{}{oldCfg["backend_node_groups"].([]interface{})[0], second}
		if !diffRequiresNew(t, oldCfg, backendGroupConfig(200, nil)) {
			t.Fatal("removing a backend group must replace the instance")
		}
	})
	groupWith := func(cu int) []interface{} {
		return []interface{}{map[string]interface{}{
			"cu": cu, "zone_id": "cn-hangzhou-i", "spec_type": "standard", "disk_number": 1,
			"storage_performance_level": "pl1", "resident_node_number": 3, "storage_size": 100,
		}}
	}
	t.Run("changing frontend groups replaces", func(t *testing.T) {
		oldCfg := backendGroupConfig(200, nil)
		oldCfg["frontend_node_groups"] = groupWith(8)
		newCfg := backendGroupConfig(200, nil)
		newCfg["frontend_node_groups"] = groupWith(16)
		if !diffRequiresNew(t, oldCfg, newCfg) {
			t.Fatal("changing frontend_node_groups must replace the instance")
		}
	})
	t.Run("changing observer groups replaces", func(t *testing.T) {
		oldCfg := backendGroupConfig(200, nil)
		oldCfg["observer_node_groups"] = groupWith(8)
		newCfg := backendGroupConfig(200, nil)
		newCfg["observer_node_groups"] = groupWith(16)
		if !diffRequiresNew(t, oldCfg, newCfg) {
			t.Fatal("changing observer_node_groups must replace the instance")
		}
	})
}

func TestUnitBackendNodeGroupSlowDiscoveryExhaustsBudget(t *testing.T) {
	handler := &starRocksNodeGroupsHandler{
		t: t, groups: []map[string]interface{}{testBackendGroup("ng-target", "cn-hangzhou-i")},
		total: 1, instanceStatus: "RUNNING",
		delayedAction: "/webapi/nodegroup/describeNodeGroups", delayedStatus: http.StatusInternalServerError, delay: 3 * time.Second,
	}
	client := starRocksTestClient(t, handler.ServeHTTP)
	config := backendGroupConfig(300, nil)
	config["instance_name"] = "tf-testacc-starrocks-renamed"
	err := applyStarRocksInstance(t, client, backendGroupStateAttributes("200"), config, 2*time.Second)
	if err == nil {
		t.Fatal("the slow group discovery must exhaust the shared update budget")
	}
	if handler.countAction("/webapi/cluster/update_name") != 0 || handler.modifyCalls != 0 {
		t.Fatalf("no write may run after the shared budget is exhausted: %#v", handler.actions)
	}
}
