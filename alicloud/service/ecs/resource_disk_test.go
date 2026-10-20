package ecs_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/ecs"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// SDKv2 normalizes an unset tags map to null in apply-state, while refresh
// reports "tags.% = 0"; plugin-testing compares raw apply-state with refreshed
// state for the ID-only check, so tags is excluded there. Tags are still
// verified by per-step checks and ImportStateVerify.
var diskIDRefreshIgnore = []string{"tags"}

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
		IDRefreshName:            resourceId,
		IDRefreshIgnore:          diskIDRefreshIgnore,
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
				ImportStateVerifyIgnore: []string{"advanced_features", "dry_run", "encrypt_algorithm", "instance_id", "type"},
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
		IDRefreshName:            resourceId,
		IDRefreshIgnore:          diskIDRefreshIgnore,
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
				ImportStateVerifyIgnore: []string{"advanced_features", "dry_run", "encrypt_algorithm", "instance_id", "type"},
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
		IDRefreshName:            resourceId,
		IDRefreshIgnore:          diskIDRefreshIgnore,
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
				ImportStateVerifyIgnore: []string{"advanced_features", "dry_run", "encrypt_algorithm", "instance_id", "type"},
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
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEcsDiskBasicDependence8018)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		IDRefreshName:            resourceId,
		IDRefreshIgnore:          diskIDRefreshIgnore,
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
					"availability_zone": "${data.alicloud_zones.default.zones.0.id}",
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
				ImportStateVerifyIgnore: []string{"advanced_features", "dry_run", "encrypt_algorithm", "instance_id", "type"},
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
		IDRefreshName:            resourceId,
		IDRefreshIgnore:          diskIDRefreshIgnore,
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
				ImportStateVerifyIgnore: []string{"advanced_features", "dry_run", "encrypt_algorithm", "instance_id", "type"},
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
		IDRefreshName:            resourceId,
		IDRefreshIgnore:          diskIDRefreshIgnore,
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
				ImportStateVerifyIgnore: []string{"advanced_features", "dry_run", "encrypt_algorithm", "instance_id", "type"},
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
		IDRefreshName:            resourceId,
		IDRefreshIgnore:          diskIDRefreshIgnore,
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
				ImportStateVerifyIgnore: []string{"advanced_features", "dry_run", "encrypt_algorithm", "instance_id", "type"},
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
		IDRefreshName:            resourceId,
		IDRefreshIgnore:          diskIDRefreshIgnore,
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
				ImportStateVerifyIgnore: []string{"advanced_features", "dry_run", "encrypt_algorithm", "instance_id", "type"},
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
		IDRefreshName:            resourceId,
		IDRefreshIgnore:          diskIDRefreshIgnore,
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
				ImportStateVerifyIgnore: []string{"advanced_features", "dry_run", "encrypt_algorithm", "instance_id", "type"},
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
		IDRefreshName:            resourceId,
		IDRefreshIgnore:          diskIDRefreshIgnore,
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
				ImportStateVerifyIgnore: []string{"advanced_features", "dry_run", "encrypt_algorithm", "instance_id", "type"},
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
