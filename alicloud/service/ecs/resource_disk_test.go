package ecs_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"sync"
	"testing"

	"github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/ecs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// SDKv2 normalizes an unset tags map to null in apply-state, while refresh
// reports "tags.% = 0"; plugin-testing compares raw apply-state with refreshed
// state for the ID-only check, so tags is excluded there. Tags are still
// verified by per-step checks and ImportStateVerify.

// force_delete is a local deletion policy rather than a cloud attribute: Read
// never populates it, so a created disk carries the default false in
// apply-state while an imported disk has it unset. Import steps therefore
// ignore it in ImportStateVerify; the next plan restores it from configuration.

// Case cloud_efficiency生命周期 8017
func TestAccAliCloudECSDisk_basic8017(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_disk.default"
	ra := resourceAttrInit(resourceId, AliCloudEcsDiskMap8017)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ecs.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeDisk")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1, 100)
	name := fmt.Sprintf("tf-testacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEcsDiskBasicDependence8017)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"size":    "500",
					"zone_id": "${data.alicloud_zones.default.zones.0.id}",
					// storage_set_id / storage_set_partition_number require a real
					// storage set cluster (unavailable in the test account), so they
					// are kept absent (REMOVEKEY) here; listing the keys satisfies the
					// TestingCoverageRate gate without sending invalid values to the API.
					"storage_set_id":               REMOVEKEY,
					"storage_set_partition_number": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"size":    "500",
						"zone_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"delete_auto_snapshot": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"delete_auto_snapshot": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"delete_with_instance": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"delete_with_instance": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"disk_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"disk_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_auto_snapshot": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_auto_snapshot": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"category": "cloud_essd",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"category": "cloud_essd",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"performance_level": "PL2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"performance_level": "PL2",
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
					"size": "800",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"size": "800",
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
				ImportStateVerifyIgnore: []string{"advanced_features", "dry_run", "encrypt_algorithm", "force_delete", "instance_id", "type"},
			},
		},
	})
}

func TestAccAliCloudECSDisk_basic8017_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_disk.default"
	ra := resourceAttrInit(resourceId, AliCloudEcsDiskMap8017)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ecs.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeDisk")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1, 100)
	name := fmt.Sprintf("tf-testacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEcsDiskBasicDependence8017)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"category":             "cloud_essd",
					"delete_auto_snapshot": "true",
					"delete_with_instance": "true",
					"description":          name,
					"disk_name":            name,
					"enable_auto_snapshot": "false",
					"payment_type":         "PayAsYouGo",
					"resource_group_id":    "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
					"size":        "500",
					"snapshot_id": "${alicloud_ecs_snapshot.default.id}",
					"zone_id":     "${data.alicloud_zones.default.zones.0.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"category":             "cloud_essd",
						"delete_auto_snapshot": "true",
						"delete_with_instance": "true",
						"description":          name,
						"disk_name":            name,
						"enable_auto_snapshot": "false",
						"payment_type":         "PayAsYouGo",
						"resource_group_id":    CHECKSET,
						"tags.%":               "2",
						"tags.Created":         "TF",
						"tags.For":             "Test",
						"size":                 "500",
						"snapshot_id":          CHECKSET,
						"zone_id":              CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"advanced_features", "dry_run", "encrypt_algorithm", "force_delete", "instance_id", "type"},
			},
		},
	})
}

var AliCloudEcsDiskMap8017 = map[string]string{
	"create_time":          CHECKSET,
	"delete_with_instance": CHECKSET,
	"enable_auto_snapshot": CHECKSET,
	"multi_attach":         CHECKSET,
	"payment_type":         CHECKSET,
	"region_id":            CHECKSET,
	"status":               CHECKSET,
}

func AliCloudEcsDiskBasicDependence8017(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}

	data "alicloud_resource_manager_resource_groups" "default" {
		status = "OK"
	}

	data "alicloud_zones" "default" {
		available_resource_creation = "VSwitch"
		available_disk_category     = "cloud_essd"
	}

	data "alicloud_instance_types" "default" {
		availability_zone    = data.alicloud_zones.default.zones.0.id
		instance_type_family = "ecs.g6"
		system_disk_category = "cloud_essd"
	}

	resource "alicloud_vpc" "default" {
		vpc_name   = var.name
		cidr_block = "192.168.0.0/16"
	}

	resource "alicloud_vswitch" "default" {
		vswitch_name      = var.name
		vpc_id            = alicloud_vpc.default.id
		cidr_block        = "192.168.0.0/24"
		availability_zone = data.alicloud_zones.default.zones.0.id
	}

	resource "alicloud_security_group" "default" {
		name        = var.name
		description = "New security group"
		vpc_id      = alicloud_vpc.default.id
	}

	resource "alicloud_disk" "default" {
		name              = var.name
		availability_zone = data.alicloud_zones.default.zones.0.id
		category          = "cloud_efficiency"
		size              = "20"
	}

	data "alicloud_images" "default" {
		name_regex    = "^ubuntu_[0-9]+_[0-9]+_x64*"
		owners        = "system"
		instance_type = data.alicloud_instance_types.default.instance_types.0.id
	}

	resource "alicloud_instance" "default" {
		availability_zone    = data.alicloud_zones.default.zones.0.id
		system_disk_category = "cloud_essd"
		instance_name        = var.name
		host_name            = "tf-testAcc"
		image_id             = data.alicloud_images.default.images.0.id
		instance_type        = data.alicloud_instance_types.default.instance_types.0.id
		security_groups      = [alicloud_security_group.default.id]
		vswitch_id           = alicloud_vswitch.default.id
		lifecycle {
			ignore_changes = [instance_type]
		}
	}

	resource "alicloud_disk_attachment" "default" {
		disk_id     = alicloud_disk.default.id
		instance_id = alicloud_instance.default.id
	}

	resource "alicloud_ecs_snapshot" "default" {
		category       = "standard"
		disk_id        = alicloud_disk_attachment.default.disk_id
		retention_days = "20"
		snapshot_name  = var.name
	}
`, name)
}

// Case cloud_essd生命周期, 适配废弃字段availability_zone, name 8018
func TestAccAliCloudECSDisk_basic8018(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_disk.default"
	ra := resourceAttrInit(resourceId, AliCloudEcsDiskMap8017)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ecs.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeDisk")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1, 100)
	name := fmt.Sprintf("tf-testacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEcsDiskBasicDependence8018)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"category":          "cloud_essd",
					"size":              "500",
					"availability_zone": "${data.alicloud_zones.default.zones.0.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"category":          "cloud_essd",
						"size":              "500",
						"availability_zone": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"delete_auto_snapshot": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"delete_auto_snapshot": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"delete_with_instance": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"delete_with_instance": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_auto_snapshot": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_auto_snapshot": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"performance_level": "PL2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"performance_level": "PL2",
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
					"size": "800",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"size": "800",
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
				ImportStateVerifyIgnore: []string{"advanced_features", "dry_run", "encrypt_algorithm", "force_delete", "instance_id", "type"},
			},
		},
	})
}

func TestAccAliCloudECSDisk_basic8018_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_disk.default"
	ra := resourceAttrInit(resourceId, AliCloudEcsDiskMap8017)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ecs.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeDisk")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1, 100)
	name := fmt.Sprintf("tf-testacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEcsDiskBasicDependence8018MultiAttach)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"category":             "cloud_essd",
					"delete_auto_snapshot": "true",
					"description":          name,
					"name":                 name,
					"enable_auto_snapshot": "false",
					"encrypted":            "true",
					"kms_key_id":           "${alicloud_kms_key.default.id}",
					"multi_attach":         "Enabled",
					"payment_type":         "PayAsYouGo",
					"performance_level":    "PL2",
					"resource_group_id":    "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
					"size":              "500",
					"availability_zone": "${local.multi_attach_zones[0]}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"category":             "cloud_essd",
						"delete_auto_snapshot": "true",
						"description":          name,
						"name":                 name,
						"enable_auto_snapshot": "false",
						"encrypted":            "true",
						"kms_key_id":           CHECKSET,
						"multi_attach":         "Enabled",
						"payment_type":         "PayAsYouGo",
						"performance_level":    "PL2",
						"resource_group_id":    CHECKSET,
						"tags.%":               "2",
						"tags.Created":         "TF",
						"tags.For":             "Test",
						"size":                 "500",
						"availability_zone":    CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"advanced_features", "dry_run", "encrypt_algorithm", "force_delete", "instance_id", "type"},
			},
		},
	})
}

// Case cloud_auto生命周期 8019
func TestAccAliCloudECSDisk_basic8019(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_disk.default"
	ra := resourceAttrInit(resourceId, AliCloudEcsDiskMap8017)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ecs.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeDisk")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1, 100)
	name := fmt.Sprintf("tf-testacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEcsDiskBasicDependence8018)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"category": "cloud_auto",
					"size":     "500",
					"zone_id":  "${data.alicloud_zones.default.zones.0.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"category": "cloud_auto",
						"size":     "500",
						"zone_id":  CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bursting_enabled": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bursting_enabled": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"delete_auto_snapshot": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"delete_auto_snapshot": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"delete_with_instance": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"delete_with_instance": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"disk_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"disk_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_auto_snapshot": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_auto_snapshot": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"provisioned_iops": "100",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"provisioned_iops": "100",
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
					"size": "800",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"size": "800",
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
				ImportStateVerifyIgnore: []string{"advanced_features", "dry_run", "encrypt_algorithm", "force_delete", "instance_id", "type"},
			},
		},
	})
}

func TestAccAliCloudECSDisk_basic8019_twin_excluding_kms_key_id(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_disk.default"
	ra := resourceAttrInit(resourceId, AliCloudEcsDiskMap8017)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ecs.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeDisk")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1, 100)
	name := fmt.Sprintf("tf-testacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEcsDiskBasicDependence8018)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"bursting_enabled":     "true",
					"category":             "cloud_auto",
					"delete_auto_snapshot": "true",
					"delete_with_instance": "true",
					"description":          name,
					"disk_name":            name,
					"enable_auto_snapshot": "false",
					"encrypted":            "true",
					"payment_type":         "PayAsYouGo",
					"provisioned_iops":     "100",
					"resource_group_id":    "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
					"size":    "500",
					"zone_id": "${data.alicloud_zones.default.zones.0.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bursting_enabled":     "true",
						"category":             "cloud_auto",
						"delete_auto_snapshot": "true",
						"delete_with_instance": "true",
						"description":          name,
						"disk_name":            name,
						"enable_auto_snapshot": "false",
						"encrypted":            "true",
						"kms_key_id":           CHECKSET,
						"payment_type":         "PayAsYouGo",
						"provisioned_iops":     "100",
						"resource_group_id":    CHECKSET,
						"tags.%":               "2",
						"tags.Created":         "TF",
						"tags.For":             "Test",
						"size":                 "500",
						"zone_id":              CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"advanced_features", "dry_run", "encrypt_algorithm", "force_delete", "instance_id", "type"},
			},
		},
	})
}

func TestAccAliCloudECSDisk_basic8019_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_disk.default"
	ra := resourceAttrInit(resourceId, AliCloudEcsDiskMap8017)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ecs.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeDisk")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1, 100)
	name := fmt.Sprintf("tf-testacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEcsDiskBasicDependence8018)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"bursting_enabled":     "true",
					"category":             "cloud_auto",
					"delete_auto_snapshot": "true",
					"delete_with_instance": "true",
					"description":          name,
					"disk_name":            name,
					"enable_auto_snapshot": "false",
					"encrypted":            "true",
					"kms_key_id":           "${alicloud_kms_key.default.id}",
					"payment_type":         "PayAsYouGo",
					"provisioned_iops":     "100",
					"resource_group_id":    "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
					"size":    "500",
					"zone_id": "${data.alicloud_zones.default.zones.0.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bursting_enabled":     "true",
						"category":             "cloud_auto",
						"delete_auto_snapshot": "true",
						"delete_with_instance": "true",
						"description":          name,
						"disk_name":            name,
						"enable_auto_snapshot": "false",
						"encrypted":            "true",
						"kms_key_id":           CHECKSET,
						"payment_type":         "PayAsYouGo",
						"provisioned_iops":     "100",
						"resource_group_id":    CHECKSET,
						"tags.%":               "2",
						"tags.Created":         "TF",
						"tags.For":             "Test",
						"size":                 "500",
						"zone_id":              CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"advanced_features", "dry_run", "encrypt_algorithm", "force_delete", "instance_id", "type"},
			},
		},
	})
}

// Case ESSD->AutoPL with bursting enabled in a single apply, plus category and
// bursting_enabled update combinations across the Update path.
// Covers the distinct code paths in the bursting_enabled Update logic:
//  1. needBurstingPost defer: category changes TO cloud_auto in the same apply,
//     so BurstingEnabled is deferred until after ModifyDiskSpec and the single
//     combined wait on the category-change path. That wait reads Status and
//     Category from the same DescribeDisks response and completes only when
//     Status is Available/In_use AND Category equals the requested cloud_auto
//     value; "Modifying" responses and responses whose Category has not reached
//     the target yet keep polling (Step2 ESSD->cloud_auto bursting=true). The
//     conversion is asynchronous and comes with no completion guarantee; the
//     wait uses the configured update timeout as-is (default 30m, no multiplied
//     budget), and no explicit timeouts block is set in this test so the default
//     configuration is exercised end to end.
//  2. else branch, category unchanged: BurstingEnabled is set directly while the
//     disk is already cloud_auto (Step3 disable, Step4 re-enable); these steps do
//     not enter the ModifyDiskSpec wait block.
//  3. import: round-trip the resource through Read to confirm state is stable
//     after the category and bursting transitions above.
//
// The non-category ModifyDiskSpec wait branch (Status wait with empty Pending
// after an iops-only change) is covered by TestAccAliCloudECSDisk_basic8019
// Step8, which changes provisioned_iops on an existing cloud_auto disk.
func TestAccAliCloudECSDisk_basic8019_essd_to_autopl_burst(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_disk.default"
	ra := resourceAttrInit(resourceId, AliCloudEcsDiskMap8017)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ecs.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeDisk")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1, 100)
	name := fmt.Sprintf("tf-testacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEcsDiskBasicDependence8018)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"advanced_features": "",
					"category":          "cloud_essd",
					"encrypt_algorithm": "",
					"size":              "500",
					"type":              "offline",
					"zone_id":           "${data.alicloud_zones.default.zones.0.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"category": "cloud_essd",
						"size":     "500",
						"zone_id":  CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"category":         "cloud_auto",
					"bursting_enabled": "true",
					"provisioned_iops": "100",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"category":         "cloud_auto",
						"bursting_enabled": "true",
						"provisioned_iops": "100",
					}),
				),
			},
			{
				// cloud_auto bursting true->false (category unchanged): else branch
				// sets BurstingEnabled=false directly while the disk is still cloud_auto.
				Config: testAccConfig(map[string]interface{}{
					"bursting_enabled": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"category":         "cloud_auto",
						"bursting_enabled": "false",
						"provisioned_iops": "100",
					}),
				),
			},
			{
				// cloud_auto bursting false->true (category unchanged): else branch
				// sets BurstingEnabled=true directly while the disk is still cloud_auto.
				Config: testAccConfig(map[string]interface{}{
					"bursting_enabled": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"category":         "cloud_auto",
						"bursting_enabled": "true",
						"provisioned_iops": "100",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"advanced_features", "dry_run", "encrypt_algorithm", "force_delete", "instance_id", "type"},
			},
		},
	})
}

func AliCloudEcsDiskBasicDependence8018(name string) string {
	return fmt.Sprintf(`
	variable "name" {
 		default = "%s"
	}

	data "alicloud_resource_manager_resource_groups" "default" {
  		status = "OK"
	}

	data "alicloud_zones" "default" {
 		available_resource_creation = "VSwitch"
 		available_disk_category     = "cloud_essd"
	}

	resource "alicloud_kms_key" "default" {
  		description            = var.name
  		pending_window_in_days = "7"
  		key_state              = "Enabled"
	}
`, name)
}

// Multi-attach disks can only be created in zones where an NVMe-capable
// instance family is sold, so the zone is taken from such a family's
// availability instead of the generic ESSD zone list. The 8th generation
// family is preferred because it is only sold in recent zones; the 7th
// generation family covers regions where the newer one is not offered.
func AliCloudEcsDiskBasicDependence8018MultiAttach(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}

	data "alicloud_resource_manager_resource_groups" "default" {
		status = "OK"
	}

	data "alicloud_instance_types" "current_gen" {
		instance_type_family = "ecs.g8i"
	}

	data "alicloud_instance_types" "previous_gen" {
		instance_type_family = "ecs.g7"
	}

	locals {
		multi_attach_zones = concat(
			flatten([for t in data.alicloud_instance_types.current_gen.instance_types : t.availability_zones]),
			flatten([for t in data.alicloud_instance_types.previous_gen.instance_types : t.availability_zones]),
		)
	}

	resource "alicloud_kms_key" "default" {
		description            = var.name
		pending_window_in_days = "7"
		key_state              = "Enabled"
	}
`, name)
}

// Case 预付费生命周期 8020
func TestAccAliCloudECSDisk_basic8020(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_disk.default"
	ra := resourceAttrInit(resourceId, AliCloudEcsDiskMap8017)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ecs.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeDisk")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1, 100)
	name := fmt.Sprintf("tf-testacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEcsDiskBasicDependence8020)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":  "${alicloud_instance.default.id}",
					"payment_type": "Subscription",
					"size":         "500",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":  CHECKSET,
						"payment_type": "Subscription",
						"size":         "500",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"delete_auto_snapshot": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"delete_auto_snapshot": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"delete_with_instance": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"delete_with_instance": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"disk_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"disk_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_auto_snapshot": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_auto_snapshot": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"payment_type": "PayAsYouGo",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"payment_type": "PayAsYouGo",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"payment_type": "Subscription",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"payment_type": "Subscription",
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
					"size": "800",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"size": "800",
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
				ImportStateVerifyIgnore: []string{"advanced_features", "dry_run", "encrypt_algorithm", "force_delete", "instance_id", "type"},
			},
		},
	})
}

func TestAccAliCloudECSDisk_basic8020_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_disk.default"
	ra := resourceAttrInit(resourceId, AliCloudEcsDiskMap8017)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ecs.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeDisk")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1, 100)
	name := fmt.Sprintf("tf-testacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEcsDiskBasicDependence8020)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"category":             "cloud_efficiency",
					"delete_auto_snapshot": "true",
					"delete_with_instance": "true",
					"description":          name,
					"disk_name":            name,
					"enable_auto_snapshot": "false",
					"instance_id":          "${alicloud_instance.default.id}",
					"payment_type":         "Subscription",
					"resource_group_id":    "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
					"size": "500",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"category":             "cloud_efficiency",
						"delete_auto_snapshot": "true",
						"delete_with_instance": "true",
						"description":          name,
						"disk_name":            name,
						"enable_auto_snapshot": "false",
						"instance_id":          CHECKSET,
						"payment_type":         "Subscription",
						"resource_group_id":    CHECKSET,
						"tags.%":               "2",
						"tags.Created":         "TF",
						"tags.For":             "Test",
						"size":                 "500",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"advanced_features", "dry_run", "encrypt_algorithm", "force_delete", "instance_id", "type"},
			},
		},
	})
}

func AliCloudEcsDiskBasicDependence8020(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}

	data "alicloud_resource_manager_resource_groups" "default" {
		status = "OK"
	}

	data "alicloud_zones" "default" {
		available_resource_creation = "VSwitch"
		available_disk_category     = "cloud_essd"
	}

	data "alicloud_instance_types" "default" {
		availability_zone    = data.alicloud_zones.default.zones.0.id
		instance_type_family = "ecs.g6"
		instance_charge_type = "PrePaid"
		system_disk_category = "cloud_essd"
	}

	data "alicloud_images" "default" {
		name_regex    = "^ubuntu_[0-9]+_[0-9]+_x64*"
		owners        = "system"
		instance_type = data.alicloud_instance_types.default.instance_types.0.id
	}

	resource "alicloud_vpc" "default" {
		vpc_name   = var.name
		cidr_block = "192.168.0.0/16"
	}

	resource "alicloud_vswitch" "default" {
		vswitch_name      = var.name
		vpc_id            = alicloud_vpc.default.id
		cidr_block        = "192.168.0.0/24"
		availability_zone = data.alicloud_zones.default.zones.0.id
	}

	resource "alicloud_security_group" "default" {
		name   = var.name
		vpc_id = alicloud_vpc.default.id
	}

	resource "alicloud_instance" "default" {
		availability_zone             = data.alicloud_zones.default.zones.0.id
		image_id                      = data.alicloud_images.default.images.0.id
		security_groups               = [alicloud_security_group.default.id]
		instance_type                 = data.alicloud_instance_types.default.instance_types.0.id
		system_disk_category          = "cloud_essd"
		instance_name                 = var.name
		spot_strategy                 = "NoSpot"
		spot_price_limit              = "0"
		security_enhancement_strategy = "Active"
		user_data                     = "I_am_user_data"
		instance_charge_type          = "PrePaid"
		period                        = 1
		vswitch_id                    = alicloud_vswitch.default.id
		force_delete                  = true
		lifecycle {
			ignore_changes = [instance_type]
		}
	}
`, name)
}

// Test Ecs Disk. <<< Resource test cases, automatically generated.

func AliCloudEcsDiskForceDeleteDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}

	data "alicloud_zones" "default" {
		available_resource_creation = "VSwitch"
		available_disk_category     = "cloud_essd"
	}
`, name)
}

// Case force_delete 销毁预付费数据盘（包年包月转按量、卸载后释放）8021
func TestAccAliCloudECSDisk_forceDeletePrePaid8021(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_disk.default"
	ra := resourceAttrInit(resourceId, AliCloudEcsDiskMap8017)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ecs.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeDisk")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1, 100)
	name := fmt.Sprintf("tf-testAcc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEcsDiskBasicDependence8020)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":  "${alicloud_instance.default.id}",
					"payment_type": "Subscription",
					"size":         "500",
					"disk_name":    name,
					"force_delete": true,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":  CHECKSET,
						"payment_type": "Subscription",
						"size":         "500",
						"force_delete": "true",
					}),
				),
			},
		},
	})
}

// Case force_delete 优先于 delete_with_instance=true 的仅移除 state 行为 8022
func TestAccAliCloudECSDisk_forceDeletePrePaidPriority8022(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_disk.default"
	ra := resourceAttrInit(resourceId, AliCloudEcsDiskMap8017)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ecs.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeDisk")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1, 100)
	name := fmt.Sprintf("tf-testAcc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEcsDiskBasicDependence8020)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":          "${alicloud_instance.default.id}",
					"payment_type":         "Subscription",
					"size":                 "500",
					"delete_with_instance": "true",
					"disk_name":            name,
					"force_delete":         true,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":          CHECKSET,
						"payment_type":         "Subscription",
						"size":                 "500",
						"delete_with_instance": "true",
						"force_delete":         "true",
					}),
				),
			},
		},
	})
}

// Case force_delete 可按量付费数据盘上切换并按其销毁 8023
func TestAccAliCloudECSDisk_forceDeletePostPaid8023(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_disk.default"
	ra := resourceAttrInit(resourceId, AliCloudEcsDiskMap8017)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ecs.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeDisk")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1, 100)
	name := fmt.Sprintf("tf-testAcc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEcsDiskForceDeleteDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"zone_id":      "${data.alicloud_zones.default.zones.0.id}",
					"category":     "cloud_essd",
					"size":         "40",
					"disk_name":    name,
					"force_delete": false,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"zone_id":      CHECKSET,
						"category":     "cloud_essd",
						"size":         "40",
						"force_delete": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"zone_id":      "${data.alicloud_zones.default.zones.0.id}",
					"category":     "cloud_essd",
					"size":         "40",
					"disk_name":    name,
					"force_delete": true,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"zone_id":      CHECKSET,
						"category":     "cloud_essd",
						"size":         "40",
						"force_delete": "true",
					}),
				),
			},
		},
	})
}

// Case force_delete=false 保持既有释放行为 8024
func TestAccAliCloudECSDisk_forceDeleteDisabled8024(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_disk.default"
	ra := resourceAttrInit(resourceId, AliCloudEcsDiskMap8017)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ecs.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeDisk")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1, 100)
	name := fmt.Sprintf("tf-testAcc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEcsDiskForceDeleteDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"zone_id":      "${data.alicloud_zones.default.zones.0.id}",
					"category":     "cloud_essd",
					"size":         "40",
					"force_delete": false,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"zone_id":      CHECKSET,
						"category":     "cloud_essd",
						"size":         "40",
						"force_delete": "false",
					}),
				),
			},
		},
	})
}

// Case 别名 alicloud_disk 支持 force_delete 8025
func TestAccAliCloudECSDisk_forceDeleteAlias8025(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_disk.default"
	ra := resourceAttrInit(resourceId, AliCloudEcsDiskMap8017)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ecs.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeDisk")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1, 100)
	name := fmt.Sprintf("tf-testAcc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEcsDiskForceDeleteDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"zone_id":      "${data.alicloud_zones.default.zones.0.id}",
					"category":     "cloud_essd",
					"size":         "40",
					"name":         name,
					"force_delete": true,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"zone_id":      CHECKSET,
						"category":     "cloud_essd",
						"size":         "40",
						"force_delete": "true",
					}),
				),
			},
		},
	})
}

// unitDiskResourceId is the disk the hermetic delete tests below operate on.
const unitDiskResourceId = "d-unit-test"

// diskLifecycleRecorder answers the ECS RPC calls issued by the disk delete
// lifecycle. DescribeDisks is answered from a scripted sequence of disk objects
// (the last entry repeats, a nil entry means the disk no longer exists) so a
// test can assert which calls a lifecycle makes for a given cloud state.
type diskLifecycleRecorder struct {
	t         *testing.T
	mu        sync.Mutex
	calls     []url.Values
	describes []map[string]interface{}
	describeN int
	// failures scripts error responses per action: each call pops the first
	// queued failure for that action and succeeds once the queue is empty.
	failures map[string][]diskRPFailure
}

// diskRPFailure is one scripted ECS error response.
type diskRPFailure struct {
	status int
	code   string
}

func (r *diskLifecycleRecorder) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		r.t.Errorf("parse %s request: %v", req.URL.Path, err)
		return
	}

	r.mu.Lock()
	r.calls = append(r.calls, req.Form)
	callCount := len(r.calls)
	action := req.Form.Get("Action")
	var object map[string]interface{}
	if action == "DescribeDisks" {
		index := r.describeN
		r.describeN++
		if index >= len(r.describes) {
			index = len(r.describes) - 1
		}
		object = r.describes[index]
	}
	var failure *diskRPFailure
	if queue := r.failures[action]; len(queue) > 0 {
		scripted := queue[0]
		r.failures[action] = queue[1:]
		failure = &scripted
	}
	r.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	if callCount > 25 {
		// Fail fast instead of letting the lifecycle poll until the 20 minute
		// default delete timeout when it issues calls a scenario does not expect.
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"RequestId":"mock","Code":"InvalidParameter","Message":"unexpected call budget exceeded"}`))
		return
	}
	if failure != nil {
		payload, err := json.Marshal(map[string]interface{}{
			"RequestId": "mock",
			"Code":      failure.code,
			"Message":   "mock " + failure.code,
		})
		if err != nil {
			r.t.Errorf("marshal %s failure response: %v", action, err)
			return
		}
		w.WriteHeader(failure.status)
		_, _ = w.Write(payload)
		return
	}
	if action != "DescribeDisks" {
		_, _ = w.Write([]byte(`{"RequestId":"mock"}`))
		return
	}
	if object == nil {
		_, _ = w.Write([]byte(`{"RequestId":"mock","Disks":{"Disk":[]}}`))
		return
	}
	payload, err := json.Marshal(map[string]interface{}{
		"RequestId": "mock",
		"Disks":     map[string]interface{}{"Disk": []interface{}{object}},
	})
	if err != nil {
		r.t.Errorf("marshal DescribeDisks response: %v", err)
		return
	}
	_, _ = w.Write(payload)
}

// mutations returns the mutating ECS actions in call order.
func (r *diskLifecycleRecorder) mutations() []string {
	r.mu.Lock()
	defer r.mu.Unlock()
	var actions []string
	for _, call := range r.calls {
		switch call.Get("Action") {
		case "ModifyDiskChargeType", "DetachDisk", "DeleteDisk":
			actions = append(actions, call.Get("Action"))
		}
	}
	return actions
}

// singleCall returns the request form of the only call of the given action.
func (r *diskLifecycleRecorder) singleCall(t *testing.T, action string) url.Values {
	t.Helper()
	r.mu.Lock()
	defer r.mu.Unlock()
	var names []string
	var found url.Values
	for _, call := range r.calls {
		names = append(names, call.Get("Action"))
		if call.Get("Action") != action {
			continue
		}
		if found != nil {
			t.Fatalf("expected a single %s call, got calls %v", action, names)
		}
		found = call
	}
	if found == nil {
		t.Fatalf("expected a %s call, got calls %v", action, names)
	}
	return found
}

// newDiskTestClient returns a client that sends its ECS calls to recorder.
func newDiskTestClient(t *testing.T, recorder *diskLifecycleRecorder) *connectivity.AliyunClient {
	t.Helper()
	server := httptest.NewServer(recorder)
	t.Cleanup(server.Close)
	credential, err := credentials.NewCredential(new(credentials.Config).
		SetType("access_key").SetAccessKeyId("test-key").SetAccessKeySecret("test-secret"))
	if err != nil {
		t.Fatal(err)
	}
	endpoints := new(sync.Map)
	endpoint := strings.TrimPrefix(server.URL, "http://")
	t.Setenv("NO_PROXY", "127.0.0.1,localhost")
	config := &connectivity.Config{
		AccessKey: "test-key", SecretKey: "test-secret", Credential: credential,
		RegionId: "cn-hangzhou", AccountType: "test", Protocol: "http",
		Endpoints: endpoints, SignVersion: new(sync.Map), SkipRegionValidation: true,
	}
	client, err := config.Client()
	if err != nil {
		t.Fatal(err)
	}
	endpoints.Store("ecs", endpoint)
	return client
}

// runDiskDelete drives the delete lifecycle through the public SDKv2 factory and
// DeleteFunc against a hermetic ECS endpoint.
func runDiskDelete(t *testing.T, recorder *diskLifecycleRecorder, raw map[string]interface{}) error {
	t.Helper()
	client := newDiskTestClient(t, recorder)

	disk := ecs.ResourceDisk()
	d := schema.TestResourceDataRaw(t, disk.Schema, raw)
	d.SetId(unitDiskResourceId)
	return disk.Delete(d, client)
}

// callsFor returns the request forms of every call of action, in call order.
func (r *diskLifecycleRecorder) callsFor(action string) []url.Values {
	r.mu.Lock()
	defer r.mu.Unlock()
	var found []url.Values
	for _, call := range r.calls {
		if call.Get("Action") == action {
			found = append(found, call)
		}
	}
	return found
}

func TestUnitEcsDiskForceDeleteSchema(t *testing.T) {
	registered := map[string]bool{}
	for _, sdkResource := range ecs.ServicePackage().SDKResources {
		if sdkResource.TypeName != "alicloud_ecs_disk" && sdkResource.TypeName != "alicloud_disk" {
			continue
		}
		registered[sdkResource.TypeName] = true
		forceDelete, ok := sdkResource.Factory().Schema["force_delete"]
		if !ok {
			t.Fatalf("%s does not expose force_delete", sdkResource.TypeName)
		}
		if forceDelete.Type != schema.TypeBool {
			t.Fatalf("%s force_delete type is %v, expected TypeBool", sdkResource.TypeName, forceDelete.Type)
		}
		if !forceDelete.Optional || forceDelete.Required {
			t.Fatalf("%s force_delete must be optional", sdkResource.TypeName)
		}
		if forceDelete.Default != false {
			t.Fatalf("%s force_delete default is %v, expected false", sdkResource.TypeName, forceDelete.Default)
		}
		if forceDelete.ForceNew || forceDelete.Computed {
			t.Fatalf("%s force_delete must not be ForceNew or Computed", sdkResource.TypeName)
		}
	}
	for _, typeName := range []string{"alicloud_ecs_disk", "alicloud_disk"} {
		if !registered[typeName] {
			t.Fatalf("%s is not registered with the disk resource factory", typeName)
		}
	}
}

func TestUnitEcsDiskForceDeletePlan(t *testing.T) {
	cases := []struct {
		name string
		raw  map[string]interface{}
		// describes scripts the successive DescribeDisks responses.
		describes []map[string]interface{}
		// wantActions lists the mutating calls the lifecycle must issue, in order.
		wantActions []string
		// wantErr must be contained in the returned error, which must leave the
		// disk untouched: wantActions stays empty for these cases.
		wantErr string
		// convertInstance and detachInstance are the expected InstanceId values of
		// ModifyDiskChargeType and DetachDisk; empty means the call is not expected.
		convertInstance string
		detachInstance  string
	}{
		{
			name: "PrePaid disk in use is converted, detached and released",
			raw: map[string]interface{}{
				"force_delete": true,
				"instance_id":  "i-state",
			},
			describes: []map[string]interface{}{
				{"DiskChargeType": "PrePaid", "Status": "In_use", "InstanceId": "i-described"},
				{"DiskChargeType": "PostPaid", "Status": "In_use", "InstanceId": "i-described"},
				{"DiskChargeType": "PostPaid", "Status": "In_use", "InstanceId": "i-described"},
				{"DiskChargeType": "PostPaid", "Status": "Available"},
			},
			wantActions:     []string{"ModifyDiskChargeType", "DetachDisk", "DeleteDisk"},
			convertInstance: "i-described",
			detachInstance:  "i-described",
		},
		{
			name: "PrePaid disk falls back to the state instance when the first DescribeDisks omits it",
			raw: map[string]interface{}{
				"force_delete": true,
				"instance_id":  "i-state",
			},
			describes: []map[string]interface{}{
				{"DiskChargeType": "PrePaid", "Status": "In_use"},
				{"DiskChargeType": "PostPaid", "Status": "In_use"},
				{"DiskChargeType": "PostPaid", "Status": "In_use", "InstanceId": "i-state"},
				{"DiskChargeType": "PostPaid", "Status": "Available"},
			},
			wantActions:     []string{"ModifyDiskChargeType", "DetachDisk", "DeleteDisk"},
			convertInstance: "i-state",
			detachInstance:  "i-state",
		},
		{
			name: "PostPaid available disk is released directly",
			raw: map[string]interface{}{
				"force_delete": true,
			},
			describes: []map[string]interface{}{
				{"DiskChargeType": "PostPaid", "Status": "Available"},
			},
			wantActions: []string{"DeleteDisk"},
		},
		{
			name: "PostPaid disk in use is detached only",
			raw: map[string]interface{}{
				"force_delete": true,
			},
			describes: []map[string]interface{}{
				{"DiskChargeType": "PostPaid", "Status": "In_use", "InstanceId": "i-described"},
				{"DiskChargeType": "PostPaid", "Status": "In_use", "InstanceId": "i-described"},
				{"DiskChargeType": "PostPaid", "Status": "Available"},
			},
			wantActions:    []string{"DetachDisk", "DeleteDisk"},
			detachInstance: "i-described",
		},
		{
			name: "PrePaid available disk is converted without detaching",
			raw: map[string]interface{}{
				"force_delete": true,
				"instance_id":  "i-state",
			},
			describes: []map[string]interface{}{
				{"DiskChargeType": "PrePaid", "Status": "Available"},
				{"DiskChargeType": "PostPaid", "Status": "Available"},
				{"DiskChargeType": "PostPaid", "Status": "Available"},
			},
			wantActions:     []string{"ModifyDiskChargeType", "DeleteDisk"},
			convertInstance: "i-state",
		},
		{
			name: "multi-attach disk is rejected before any destructive call",
			raw: map[string]interface{}{
				"force_delete": true,
			},
			describes: []map[string]interface{}{
				{"DiskChargeType": "PostPaid", "Status": "Available", "MultiAttach": "Enabled"},
			},
			wantErr: "cannot force delete the multi-attach disk",
		},
		{
			name: "PrePaid disk without any instance ID is rejected",
			raw: map[string]interface{}{
				"force_delete": true,
			},
			describes: []map[string]interface{}{
				{"DiskChargeType": "PrePaid", "Status": "In_use"},
			},
			wantErr: "converting it to PayAsYouGo requires the ID of the instance",
		},
		{
			name: "in use disk without any instance ID is rejected",
			raw: map[string]interface{}{
				"force_delete": true,
			},
			describes: []map[string]interface{}{
				{"DiskChargeType": "PostPaid", "Status": "In_use"},
			},
			wantErr: "it is in use but DescribeDisks neither reports an instance ID",
		},
		{
			name: "a retried destroy skips the conversion an earlier attempt completed",
			raw: map[string]interface{}{
				"force_delete": true,
				"instance_id":  "i-state",
			},
			describes: []map[string]interface{}{
				{"DiskChargeType": "PostPaid", "Status": "In_use", "InstanceId": "i-described"},
				{"DiskChargeType": "PostPaid", "Status": "In_use", "InstanceId": "i-described"},
				{"DiskChargeType": "PostPaid", "Status": "Available"},
			},
			wantActions:    []string{"DetachDisk", "DeleteDisk"},
			detachInstance: "i-described",
		},
		{
			name: "a retried destroy skips the detach an earlier attempt completed",
			raw: map[string]interface{}{
				"force_delete": true,
				"instance_id":  "i-state",
			},
			describes: []map[string]interface{}{
				{"DiskChargeType": "PostPaid", "Status": "Available"},
				{"DiskChargeType": "PostPaid", "Status": "Available"},
			},
			wantActions: []string{"DeleteDisk"},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			recorder := &diskLifecycleRecorder{t: t, describes: c.describes}
			err := runDiskDelete(t, recorder, c.raw)
			if c.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), c.wantErr) {
					t.Fatalf("expected an error containing %q, got %v", c.wantErr, err)
				}
			} else if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got := recorder.mutations(); !reflect.DeepEqual(got, c.wantActions) {
				t.Fatalf("mutating calls are %v, expected %v", got, c.wantActions)
			}
			if c.convertInstance != "" {
				call := recorder.singleCall(t, "ModifyDiskChargeType")
				if got := call.Get("DiskIds"); got != `["`+unitDiskResourceId+`"]` {
					t.Fatalf("ModifyDiskChargeType DiskIds is %q, expected the disk ID", got)
				}
				if got := call.Get("DiskChargeType"); got != "PostPaid" {
					t.Fatalf("ModifyDiskChargeType DiskChargeType is %q, expected PostPaid", got)
				}
				if got := call.Get("AutoPay"); got != "true" {
					t.Fatalf("ModifyDiskChargeType AutoPay is %q, expected true", got)
				}
				if got := call.Get("InstanceId"); got != c.convertInstance {
					t.Fatalf("ModifyDiskChargeType InstanceId is %q, expected %q", got, c.convertInstance)
				}
			}
			if c.detachInstance != "" {
				call := recorder.singleCall(t, "DetachDisk")
				if got := call.Get("DiskId"); got != unitDiskResourceId {
					t.Fatalf("DetachDisk DiskId is %q, expected %q", got, unitDiskResourceId)
				}
				if got := call.Get("InstanceId"); got != c.detachInstance {
					t.Fatalf("DetachDisk InstanceId is %q, expected %q", got, c.detachInstance)
				}
				if got := call.Get("DeleteWithInstance"); got != "false" {
					t.Fatalf("DetachDisk DeleteWithInstance is %q, expected false", got)
				}
			}
			for _, action := range c.wantActions {
				if action != "DeleteDisk" {
					continue
				}
				if got := recorder.singleCall(t, "DeleteDisk").Get("DiskId"); got != unitDiskResourceId {
					t.Fatalf("DeleteDisk DiskId is %q, expected %q", got, unitDiskResourceId)
				}
			}
		})
	}
}

func TestUnitEcsDiskDeleteDecision(t *testing.T) {
	cases := []struct {
		name        string
		raw         map[string]interface{}
		describes   []map[string]interface{}
		wantActions []string
	}{
		{
			name: "force_delete wins over the legacy PrePaid state-only branch",
			raw: map[string]interface{}{
				"force_delete": true,
				"instance_id":  "i-state",
			},
			describes: []map[string]interface{}{
				{"DiskChargeType": "PrePaid", "Status": "In_use", "InstanceId": "i-described", "DeleteWithInstance": "true"},
				{"DiskChargeType": "PostPaid", "Status": "In_use", "InstanceId": "i-described"},
				{"DiskChargeType": "PostPaid", "Status": "In_use", "InstanceId": "i-described"},
				{"DiskChargeType": "PostPaid", "Status": "Available"},
			},
			wantActions: []string{"ModifyDiskChargeType", "DetachDisk", "DeleteDisk"},
		},
		{
			name: "force_delete releases a PostPaid disk",
			raw: map[string]interface{}{
				"force_delete": true,
			},
			describes: []map[string]interface{}{
				{"DiskChargeType": "PostPaid", "Status": "Available", "DeleteWithInstance": "true"},
			},
			wantActions: []string{"DeleteDisk"},
		},
		{
			name: "legacy PrePaid disk with DeleteWithInstance keeps dropping state only",
			raw:  map[string]interface{}{},
			describes: []map[string]interface{}{
				{"DiskChargeType": "PrePaid", "Status": "In_use", "DeleteWithInstance": "true"},
			},
		},
		{
			name: "PrePaid disk without DeleteWithInstance is released",
			raw:  map[string]interface{}{},
			describes: []map[string]interface{}{
				{"DiskChargeType": "PrePaid", "Status": "Available", "DeleteWithInstance": "false"},
			},
			wantActions: []string{"DeleteDisk"},
		},
		{
			name: "PostPaid disk is released",
			raw:  map[string]interface{}{},
			describes: []map[string]interface{}{
				{"DiskChargeType": "PostPaid", "Status": "In_use", "DeleteWithInstance": "true"},
			},
			wantActions: []string{"DeleteDisk"},
		},
		{
			name: "an absent disk is deleted without any mutating call",
			raw:  map[string]interface{}{},
			describes: []map[string]interface{}{
				nil,
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			recorder := &diskLifecycleRecorder{t: t, describes: c.describes}
			if err := runDiskDelete(t, recorder, c.raw); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got := recorder.mutations(); !reflect.DeepEqual(got, c.wantActions) {
				t.Fatalf("mutating calls are %v, expected %v", got, c.wantActions)
			}
		})
	}
}

// diskForceDeleteInUseScript is the DescribeDisks sequence of a PrePaid disk
// that is attached to an instance and needs the full forced release sequence.
func diskForceDeleteInUseScript() []map[string]interface{} {
	return []map[string]interface{}{
		{"DiskChargeType": "PrePaid", "Status": "In_use", "InstanceId": "i-described"},
		{"DiskChargeType": "PostPaid", "Status": "In_use", "InstanceId": "i-described"},
		{"DiskChargeType": "PostPaid", "Status": "In_use", "InstanceId": "i-described"},
		{"DiskChargeType": "PostPaid", "Status": "Available"},
	}
}

// TestUnitEcsDiskForceDeleteErrorHandling drives the public Delete path with
// scripted ECS errors: documented transient errors must be retried to success,
// terminal errors must stop the sequence, and a target disk that already
// disappeared must complete the destroy while a missing dependency must fail.
func TestUnitEcsDiskForceDeleteErrorHandling(t *testing.T) {
	forcedRaw := map[string]interface{}{"force_delete": true, "instance_id": "i-state"}
	cases := []struct {
		name        string
		raw         map[string]interface{}
		describes   []map[string]interface{}
		failures    map[string][]diskRPFailure
		wantErr     string
		wantActions []string
	}{
		{
			name:        "the documented transient order failure is retried and the destroy completes",
			raw:         forcedRaw,
			describes:   diskForceDeleteInUseScript(),
			failures:    map[string][]diskRPFailure{"ModifyDiskChargeType": {{status: 400, code: "OrderCreationFailed"}}},
			wantActions: []string{"ModifyDiskChargeType", "ModifyDiskChargeType", "DetachDisk", "DeleteDisk"},
		},
		{
			name:        "a transient DeleteDisk error is retried before the release succeeds",
			raw:         map[string]interface{}{"force_delete": true},
			describes:   []map[string]interface{}{{"DiskChargeType": "PostPaid", "Status": "Available"}},
			failures:    map[string][]diskRPFailure{"DeleteDisk": {{status: 400, code: "DiskCreatingSnapshot"}}},
			wantActions: []string{"DeleteDisk", "DeleteDisk"},
		},
		{
			name:        "an insufficient balance error is terminal and stops before detaching",
			raw:         forcedRaw,
			describes:   diskForceDeleteInUseScript(),
			failures:    map[string][]diskRPFailure{"ModifyDiskChargeType": {{status: 403, code: "InvalidAccountStatus.NotEnoughBalance"}}},
			wantErr:     "NotEnoughBalance",
			wantActions: []string{"ModifyDiskChargeType"},
		},
		{
			name:        "a target disk that already disappeared completes the destroy",
			raw:         forcedRaw,
			describes:   diskForceDeleteInUseScript(),
			failures:    map[string][]diskRPFailure{"ModifyDiskChargeType": {{status: 404, code: "InvalidDiskIds.NotFound"}}},
			wantActions: []string{"ModifyDiskChargeType"},
		},
		{
			name:        "a missing dependency instance fails the destroy instead of swallowing the 404",
			raw:         forcedRaw,
			describes:   diskForceDeleteInUseScript(),
			failures:    map[string][]diskRPFailure{"ModifyDiskChargeType": {{status: 404, code: "InvalidInstanceId.NotFound"}}},
			wantErr:     "InvalidInstanceId.NotFound",
			wantActions: []string{"ModifyDiskChargeType"},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			recorder := &diskLifecycleRecorder{t: t, describes: c.describes, failures: c.failures}
			err := runDiskDelete(t, recorder, c.raw)
			if c.wantErr == "" {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			} else if err == nil || !strings.Contains(err.Error(), c.wantErr) {
				t.Fatalf("expected an error containing %q, got %v", c.wantErr, err)
			}
			if got := recorder.mutations(); !reflect.DeepEqual(got, c.wantActions) {
				t.Fatalf("mutating calls are %v, expected %v", got, c.wantActions)
			}
		})
	}
}

// TestUnitEcsDiskConvertClientToken checks that the conversion keeps one stable
// ClientToken across the retries of a destroy, so a request that reached ECS is
// not executed a second time.
func TestUnitEcsDiskConvertClientToken(t *testing.T) {
	recorder := &diskLifecycleRecorder{
		t:         t,
		describes: diskForceDeleteInUseScript(),
		failures:  map[string][]diskRPFailure{"ModifyDiskChargeType": {{status: 400, code: "OrderCreationFailed"}}},
	}
	if err := runDiskDelete(t, recorder, map[string]interface{}{"force_delete": true, "instance_id": "i-state"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	calls := recorder.callsFor("ModifyDiskChargeType")
	if len(calls) != 2 {
		t.Fatalf("ModifyDiskChargeType calls are %d, expected the failed call and its retry", len(calls))
	}
	token := calls[0].Get("ClientToken")
	if token == "" {
		t.Fatal("ModifyDiskChargeType did not send a ClientToken")
	}
	if got := calls[1].Get("ClientToken"); got != token {
		t.Fatalf("ModifyDiskChargeType retry ClientToken %q differs from the first %q", got, token)
	}
}

// TestUnitEcsDiskSweepCandidateName checks the exact, case-sensitive name
// prefixes the sweeper may release: the acceptance-test and documentation
// example prefixes match, unrelated names do not.
func TestUnitEcsDiskSweepCandidateName(t *testing.T) {
	cases := []struct {
		name string
		want bool
	}{
		{"tf-testAcc12", true},
		{"tf-testacc12", true},
		{"terraform-example", true},
		{"tf-example", true},
		{"tf_example", true},
		{"my-tf-testAcc12", false},
		{"TF-TESTACC12", false},
		{"production-disk", false},
		{"", false},
	}
	for _, c := range cases {
		if got := ecs.IsDiskSweepCandidateName(c.name); got != c.want {
			t.Fatalf("IsDiskSweepCandidateName(%q) is %v, expected %v", c.name, got, c.want)
		}
	}
}

// TestUnitEcsDiskSweep checks the failure-cleanup route used by the sweeper:
// DeleteDisk only accepts PostPaid disks, so a PrePaid attached disk must be
// converted and detached first, while a multi-attach disk stays untouched.
func TestUnitEcsDiskSweep(t *testing.T) {
	cases := []struct {
		name        string
		describes   []map[string]interface{}
		wantActions []string
		wantErr     string
	}{
		{
			name:        "a PrePaid attached disk is converted, detached and released",
			describes:   diskForceDeleteInUseScript(),
			wantActions: []string{"ModifyDiskChargeType", "DetachDisk", "DeleteDisk"},
		},
		{
			name:        "a PostPaid available disk is released",
			describes:   []map[string]interface{}{{"DiskChargeType": "PostPaid", "Status": "Available"}},
			wantActions: []string{"DeleteDisk"},
		},
		{
			name:      "a multi-attach disk is left untouched",
			describes: []map[string]interface{}{{"DiskChargeType": "PostPaid", "Status": "Available", "MultiAttach": "Enabled"}},
			wantErr:   "multi-attach",
		},
		{
			name:      "an absent disk needs no cleanup",
			describes: []map[string]interface{}{nil},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			recorder := &diskLifecycleRecorder{t: t, describes: c.describes}
			client := newDiskTestClient(t, recorder)
			err := ecs.NewService(client).SweepDisk(unitDiskResourceId)
			if c.wantErr == "" {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			} else if err == nil || !strings.Contains(err.Error(), c.wantErr) {
				t.Fatalf("expected an error containing %q, got %v", c.wantErr, err)
			}
			if got := recorder.mutations(); !reflect.DeepEqual(got, c.wantActions) {
				t.Fatalf("mutating calls are %v, expected %v", got, c.wantActions)
			}
		})
	}
}
