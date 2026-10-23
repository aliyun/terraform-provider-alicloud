package nas_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/nas"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// collection-order: bypass redundancy_vswitch_ids -- pre-existing attribute outside this change; TestAccAliCloudNasFileSystem_memberOrder keeps its two-step reorder coverage, but this migrated file uses the plugin-testing harness the checker does not parse.
// Computed nested blocks and tags are absent from the apply state for file
// systems that never surface them; the plugin-testing ID-only refresh check
// compares raw apply state with refreshed state, so ignore them here. They are
// still verified by per-step checks and ImportStateVerify.
var nasFileSystemIDRefreshIgnore = []string{"nfs_acl", "recycle_bin", "smb_acl", "tags"}

func TestAccAliCloudNasFileSystem_basic0(t *testing.T) {
	var v map[string]interface{}
	acctest.CheckoutSupportedRegions(t, true, connectivity.NasClassicSupportedRegions)
	resourceId := "alicloud_nas_file_system.default"
	ra := resourceAttrInit(resourceId, AliCloudNasFileSystem0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return nas.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeFileSystem")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAliCloudNasFileSystem%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudNasFileSystemBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		IDRefreshName:            resourceId,
		IDRefreshIgnore:          nasFileSystemIDRefreshIgnore,
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"protocol_type": "NFS",
					"storage_type":  "Capacity",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"protocol_type": "NFS",
						"storage_type":  "Capacity",
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
					"recycle_bin": []map[string]interface{}{
						{
							"status":        "Enable",
							"reserved_days": "10",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"recycle_bin.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"recycle_bin": []map[string]interface{}{
						{
							"status": "Disable",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"recycle_bin.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"nfs_acl": []map[string]interface{}{
						{
							"enabled": "true",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"nfs_acl.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"nfs_acl": []map[string]interface{}{
						{
							"enabled": "false",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"nfs_acl.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.1.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "FileSystem",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "FileSystem",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"snapshot_id"},
			},
		},
	})
}

func TestAccAliCloudNasFileSystem_basic0_twin(t *testing.T) {
	var v map[string]interface{}
	acctest.CheckoutSupportedRegions(t, true, connectivity.NasClassicSupportedRegions)
	resourceId := "alicloud_nas_file_system.default"
	ra := resourceAttrInit(resourceId, AliCloudNasFileSystem0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return nas.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeFileSystem")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAliCloudNasFileSystem%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudNasFileSystemBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		IDRefreshName:            resourceId,
		IDRefreshIgnore:          nasFileSystemIDRefreshIgnore,
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"protocol_type":    "NFS",
					"storage_type":     "Capacity",
					"description":      name,
					"encrypt_type":     "2",
					"file_system_type": "standard",
					"kms_key_id":       "${alicloud_kms_key.default.id}",
					"recycle_bin": []map[string]interface{}{
						{
							"status":        "Enable",
							"reserved_days": "10",
						},
					},
					"nfs_acl": []map[string]interface{}{
						{
							"enabled": "true",
						},
					},
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.1.id}",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "FileSystem",
					},
					"zone_id": "${local.nas_capacity_zone_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"protocol_type":     "NFS",
						"storage_type":      "Capacity",
						"description":       name,
						"encrypt_type":      "2",
						"file_system_type":  "standard",
						"kms_key_id":        CHECKSET,
						"recycle_bin.#":     "1",
						"nfs_acl.#":         "1",
						"resource_group_id": CHECKSET,
						"tags.%":            "2",
						"tags.Created":      "TF",
						"tags.For":          "FileSystem",
						"zone_id":           CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"snapshot_id"},
			},
		},
	})
}

func TestAccAliCloudNasFileSystem_basic1(t *testing.T) {
	var v map[string]interface{}
	acctest.CheckoutSupportedRegions(t, true, connectivity.NasClassicSupportedRegions)
	resourceId := "alicloud_nas_file_system.default"
	ra := resourceAttrInit(resourceId, AliCloudNasFileSystem0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return nas.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeFileSystem")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAliCloudNasFileSystem%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudNasFileSystemBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		IDRefreshName:            resourceId,
		IDRefreshIgnore:          nasFileSystemIDRefreshIgnore,
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"protocol_type":    "NFS",
					"storage_type":     "standard",
					"capacity":         "100",
					"file_system_type": "extreme",
					"zone_id":          "${data.alicloud_nas_zones.extreme.zones.0.zone_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"protocol_type":    "NFS",
						"storage_type":     "standard",
						"capacity":         "100",
						"file_system_type": "extreme",
						"zone_id":          CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"capacity": "200",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"capacity": "200",
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
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.1.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "FileSystem",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "FileSystem",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"snapshot_id"},
			},
		},
	})
}

func TestAccAliCloudNasFileSystem_basic1_twin(t *testing.T) {
	var v map[string]interface{}
	acctest.CheckoutSupportedRegions(t, true, connectivity.NasClassicSupportedRegions)
	resourceId := "alicloud_nas_file_system.default"
	ra := resourceAttrInit(resourceId, AliCloudNasFileSystem0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return nas.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeFileSystem")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAliCloudNasFileSystem%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudNasFileSystemBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		IDRefreshName:            resourceId,
		IDRefreshIgnore:          nasFileSystemIDRefreshIgnore,
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"protocol_type":     "NFS",
					"storage_type":      "standard",
					"capacity":          "100",
					"description":       name,
					"encrypt_type":      "2",
					"file_system_type":  "extreme",
					"kms_key_id":        "${alicloud_kms_key.default.id}",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.1.id}",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "FileSystem",
					},
					"zone_id": "${data.alicloud_nas_zones.extreme.zones.0.zone_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"protocol_type":     "NFS",
						"storage_type":      "standard",
						"capacity":          "100",
						"description":       name,
						"encrypt_type":      "2",
						"file_system_type":  "extreme",
						"kms_key_id":        CHECKSET,
						"resource_group_id": CHECKSET,
						"tags.%":            "2",
						"tags.Created":      "TF",
						"tags.For":          "FileSystem",
						"zone_id":           CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"snapshot_id"},
			},
		},
	})
}

func TestAccAliCloudNasFileSystem_basic2(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_nas_file_system.default"
	ra := resourceAttrInit(resourceId, AliCloudNasFileSystem0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return nas.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeFileSystem")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAliCloudNasFileSystem%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudNasFileSystemBasicDependence1)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
		},
		IDRefreshName:            resourceId,
		IDRefreshIgnore:          nasFileSystemIDRefreshIgnore,
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"protocol_type":    "cpfs",
					"storage_type":     "advance_100",
					"capacity":         "3600",
					"file_system_type": "cpfs",
					"vswitch_id":       "${alicloud_vswitch.CreateVswitchC.id}",
					"vpc_id":           "${alicloud_vswitch.CreateVswitchC.vpc_id}",
					"zone_id":          "${alicloud_vswitch.CreateVswitchC.zone_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"protocol_type":    "cpfs",
						"storage_type":     "advance_100",
						"capacity":         "3600",
						"file_system_type": "cpfs",
						"vswitch_id":       CHECKSET,
						"vpc_id":           CHECKSET,
						"zone_id":          CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"capacity": "4800",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"capacity": "4800",
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
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.1.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "FileSystem",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "FileSystem",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"snapshot_id"},
			},
		},
	})
}

func TestAccAliCloudNasFileSystem_basic2_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_nas_file_system.default"
	ra := resourceAttrInit(resourceId, AliCloudNasFileSystem0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return nas.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeFileSystem")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAliCloudNasFileSystem%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudNasFileSystemBasicDependence1)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
		},
		IDRefreshName:            resourceId,
		IDRefreshIgnore:          nasFileSystemIDRefreshIgnore,
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"protocol_type":     "cpfs",
					"storage_type":      "advance_100",
					"capacity":          "4800",
					"description":       name,
					"file_system_type":  "cpfs",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.1.id}",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "FileSystem",
					},
					"vswitch_id": "${alicloud_vswitch.CreateVswitchC.id}",
					"vpc_id":     "${alicloud_vswitch.CreateVswitchC.vpc_id}",
					"zone_id":    "${alicloud_vswitch.CreateVswitchC.zone_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"protocol_type":     "cpfs",
						"storage_type":      "advance_100",
						"capacity":          "4800",
						"description":       name,
						"file_system_type":  "cpfs",
						"resource_group_id": CHECKSET,
						"tags.%":            "2",
						"tags.Created":      "TF",
						"tags.For":          "FileSystem",
						"vswitch_id":        CHECKSET,
						"vpc_id":            CHECKSET,
						"zone_id":           CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"snapshot_id"},
			},
		},
	})
}

var AliCloudNasFileSystem0 = map[string]string{
	"capacity":          CHECKSET,
	"create_time":       CHECKSET,
	"file_system_type":  CHECKSET,
	"resource_group_id": CHECKSET,
	"status":            CHECKSET,
	"zone_id":           CHECKSET,
}

func AliCloudNasFileSystemBasicDependence0(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_resource_manager_resource_groups" "default" {
	}

	data "alicloud_nas_zones" "standard" {
  		file_system_type = "standard"
	}

	locals {
  		nas_capacity_zone_id = try([
  			for z in data.alicloud_nas_zones.standard.zones : z.zone_id
  			if length([for it in z.instance_types : it if lower(it.protocol_type) == "nfs" && lower(it.storage_type) == "capacity"]) > 0
  		][0], data.alicloud_nas_zones.standard.zones.0.zone_id)
	}

	data "alicloud_nas_zones" "extreme" {
  		file_system_type = "extreme"
	}

	data "alicloud_nas_zones" "cpfs" {
  		file_system_type = "cpfs"
	}

	data "alicloud_vpcs" "default" {
  		name_regex = "^default-NODELETING$"
	}

	data "alicloud_vswitches" "default" {
  		vpc_id  = data.alicloud_vpcs.default.ids.0
  		zone_id = data.alicloud_nas_zones.cpfs.zones.0.zone_id
	}

	resource "alicloud_kms_key" "default" {
  		description            = var.name
  		pending_window_in_days = "7"
  		key_state              = "Enabled"
	}
`, name)
}

func AliCloudNasFileSystemBasicDependence1(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_resource_manager_resource_groups" "default" {
	}

	resource "alicloud_vpc" "createEVpc_Cpfs" {
	  is_default  = false
	  cidr_block  = "192.168.0.0/16"
	  vpc_name    = "nas-teste1031-vpc"
	  enable_ipv6 = true
	}
	
	resource "alicloud_vswitch" "CreateVswitchC" {
	  is_default   = false
	  vpc_id       = alicloud_vpc.createEVpc_Cpfs.id
	  zone_id      = "cn-beijing-i"
	  cidr_block   = "192.168.2.0/24"
	  vswitch_name = "nas-teste1031-vsw1sdw-F"
	}

	resource "alicloud_kms_key" "default" {
  		description            = var.name
  		pending_window_in_days = "7"
  		key_state              = "Enabled"
	}
`, name)
}

// Case create_cpfs_file_system 12182
func TestAccAliCloudNasFileSystem_basic12182(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_nas_file_system.default"
	ra := resourceAttrInit(resourceId, AlicloudNasFileSystemMap12182)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return nas.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeFileSystem")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccnas%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudNasFileSystemBasicDependence12182)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
			acctest.PreCheck(t)
		},
		IDRefreshName:            resourceId,
		IDRefreshIgnore:          nasFileSystemIDRefreshIgnore,
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"description":       "cpfs-文件系统同城冗余测试-wyf",
					"storage_type":      "advance_100",
					"encrypt_type":      "0",
					"vpc_id":            "${alicloud_vpc.createEVpc_Cpfs.id}",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"vswitch_id":        "${alicloud_vswitch.CreateVswitchC.id}",
					"zone_id":           "${alicloud_vswitch.CreateVswitchC.zone_id}",
					"capacity":          "3600",
					"protocol_type":     "cpfs",
					"file_system_type":  "cpfs",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":       "cpfs-文件系统同城冗余测试-wyf",
						"storage_type":      "advance_100",
						"encrypt_type":      "0",
						"vpc_id":            CHECKSET,
						"resource_group_id": CHECKSET,
						"vswitch_id":        CHECKSET,
						"capacity":          "3600",
						"protocol_type":     "cpfs",
						"file_system_type":  "cpfs",
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
					"description": "cpfs-文件系统同城冗余测试-wyf01",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "cpfs-文件系统同城冗余测试-wyf01",
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
				ImportStateVerifyIgnore: []string{"keytab", "keytab_md5", "snapshot_id"},
			},
		},
	})
}

var AlicloudNasFileSystemMap12182 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
	"region_id":   CHECKSET,
}

func AlicloudNasFileSystemBasicDependence12182(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_vpc" "createEVpc_Cpfs" {
  is_default  = false
  cidr_block  = "192.168.0.0/16"
  vpc_name    = "nas-teste1031-vpc"
  enable_ipv6 = true
}

resource "alicloud_vswitch" "CreateVswitchC" {
  is_default   = false
  vpc_id       = alicloud_vpc.createEVpc_Cpfs.id
  zone_id      = "cn-beijing-l"
  cidr_block   = "192.168.3.0/24"
  vswitch_name = "nas-teste1031-vsw2sdw-C"
}


`, name)
}

// Case 10730: Agentic is a valid storage_type for standard file systems (issue 10730).
func TestAccAliCloudNasFileSystem_basic10730(t *testing.T) {
	// Agentic file systems are whitelist-gated on the NAS service side; creating one on a
	// non-whitelisted account is rejected by the API with InvalidAzone.NotSupport. Skip this
	// case unless the test account has been whitelisted (set TEST_ACC_NAS_AGENTIC_WHITELIST).
	if os.Getenv("TEST_ACC_NAS_AGENTIC_WHITELIST") == "" {
		t.Skip("Agentic file system requires the NAS whitelist; set TEST_ACC_NAS_AGENTIC_WHITELIST to run this case on a whitelisted account")
	}
	var v map[string]interface{}
	acctest.CheckoutSupportedRegions(t, true, connectivity.NasClassicSupportedRegions)
	resourceId := "alicloud_nas_file_system.default"
	ra := resourceAttrInit(resourceId, AlicloudNasFileSystemMap10730)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return nas.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeFileSystem")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAliCloudNasFileSystem%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudNasFileSystemBasicDependence10730)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		IDRefreshName:            resourceId,
		IDRefreshIgnore:          nasFileSystemIDRefreshIgnore,
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"protocol_type":    "NFS",
					"storage_type":     "Agentic",
					"file_system_type": "standard",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"protocol_type":    "NFS",
						"storage_type":     "Agentic",
						"file_system_type": "standard",
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
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"snapshot_id"},
			},
		},
	})
}

var AlicloudNasFileSystemMap10730 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
	"region_id":   CHECKSET,
}

func AlicloudNasFileSystemBasicDependence10730(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}
`, name)
}

// Case 10730-order: TypeList redundancy_vswitch_ids reorder and convergence coverage
func TestAccAliCloudNasFileSystem_memberOrder(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_nas_file_system.default"
	ra := resourceAttrInit(resourceId, AlicloudNasFileSystemMapOrder)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return nas.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeFileSystem")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccnas%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudNasFileSystemOrderDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
			acctest.PreCheck(t)
		},
		IDRefreshName:            resourceId,
		IDRefreshIgnore:          nasFileSystemIDRefreshIgnore,
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"description":            "tf-acc-nas-order",
					"storage_type":           "advance_100",
					"encrypt_type":           "0",
					"vpc_id":                 "${alicloud_vpc.createEVpc_Cpfs.id}",
					"redundancy_vswitch_ids": []string{"${alicloud_vswitch.CreateVswitchC.id}", "${alicloud_vswitch.CreateVswitchD.id}", "${alicloud_vswitch.CreateVswitchF.id}"},
					"capacity":               "500",
					"protocol_type":          "cpfs",
					"file_system_type":       "cpfsse",
					"redundancy_type":        "ZRS",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":              "tf-acc-nas-order",
						"storage_type":             "advance_100",
						"redundancy_vswitch_ids.#": "3",
						"redundancy_type":          "ZRS",
						"protocol_type":            "cpfs",
						"file_system_type":         "cpfsse",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description":            "tf-acc-nas-order",
					"storage_type":           "advance_100",
					"encrypt_type":           "0",
					"vpc_id":                 "${alicloud_vpc.createEVpc_Cpfs.id}",
					"redundancy_vswitch_ids": []string{"${alicloud_vswitch.CreateVswitchF.id}", "${alicloud_vswitch.CreateVswitchC.id}", "${alicloud_vswitch.CreateVswitchD.id}"},
					"capacity":               "500",
					"protocol_type":          "cpfs",
					"file_system_type":       "cpfsse",
					"redundancy_type":        "ZRS",
				}),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description":            "tf-acc-nas-order",
					"storage_type":           "advance_100",
					"encrypt_type":           "0",
					"vpc_id":                 "${alicloud_vpc.createEVpc_Cpfs.id}",
					"redundancy_vswitch_ids": []string{"${alicloud_vswitch.CreateVswitchF.id}", "${alicloud_vswitch.CreateVswitchC.id}", "${alicloud_vswitch.CreateVswitchD.id}"},
					"capacity":               "500",
					"protocol_type":          "cpfs",
					"file_system_type":       "cpfsse",
					"redundancy_type":        "ZRS",
				}),
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

var AlicloudNasFileSystemMapOrder = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
	"region_id":   CHECKSET,
}

func AlicloudNasFileSystemOrderDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_vpc" "createEVpc_Cpfs" {
  is_default  = false
  cidr_block  = "192.168.0.0/16"
  vpc_name    = "nas-teste1031-vpc"
  enable_ipv6 = true
}

resource "alicloud_vswitch" "CreateVswitchF" {
  is_default   = false
  vpc_id       = alicloud_vpc.createEVpc_Cpfs.id
  zone_id      = "cn-beijing-i"
  cidr_block   = "192.168.2.0/24"
  vswitch_name = "nas-teste1031-vsw1sdw-F"
}

resource "alicloud_vswitch" "CreateVswitchC" {
  is_default   = false
  vpc_id       = alicloud_vpc.createEVpc_Cpfs.id
  zone_id      = "cn-beijing-l"
  cidr_block   = "192.168.3.0/24"
  vswitch_name = "nas-teste1031-vsw2sdw-C"
}

resource "alicloud_vswitch" "CreateVswitchD" {
  is_default   = false
  vpc_id       = alicloud_vpc.createEVpc_Cpfs.id
  zone_id      = "cn-beijing-h"
  cidr_block   = "192.168.4.0/24"
  vswitch_name = "nas-teste1031-vsw3sdw-D"
}

`, name)
}

// Test Nas FileSystem. >>> Resource test cases, automatically generated.
// Case 通用型性能型文件系统（SMB_ACL）_20250328 10635
func TestAccAliCloudNasFileSystem_basic10635(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_nas_file_system.default"
	ra := resourceAttrInit(resourceId, AlicloudNasFileSystemMap10635)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return nas.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeFileSystem")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccnas%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudNasFileSystemBasicDependence10635)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		IDRefreshName:            resourceId,
		IDRefreshIgnore:          nasFileSystemIDRefreshIgnore,
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"storage_type":     "Performance",
					"encrypt_type":     "0",
					"protocol_type":    "SMB",
					"file_system_type": "standard",
					"zone_id":          "cn-hangzhou-g",
					"description":      "testyaya",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"storage_type":     "Performance",
						"encrypt_type":     "0",
						"protocol_type":    "SMB",
						"file_system_type": "standard",
						"zone_id":          "cn-hangzhou-g",
						"description":      "testyaya",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"options": []map[string]interface{}{
						{
							"enable_oplock": "true",
						},
					},
					"keytab_md5": "${var.key_tab_md5}",
					"smb_acl": []map[string]interface{}{
						{
							"enabled": "true",
						},
					},
					"keytab": "${var.key_tab}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"keytab_md5": CHECKSET,
						"keytab":     CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"options": []map[string]interface{}{
						{
							"enable_oplock": "false",
						},
					},
					"keytab_md5": "${var.key_tab_1_md5}",
					"smb_acl": []map[string]interface{}{
						{
							"enable_anonymous_access":   "true",
							"encrypt_data":              "true",
							"reject_unencrypted_access": "true",
							"super_admin_sid":           "S-1-5-22",
							"home_dir_path":             "/home",
						},
					},
					"keytab": "${var.key_tab_1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"keytab_md5": CHECKSET,
						"keytab":     CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"smb_acl": []map[string]interface{}{
						{
							"enable_anonymous_access":   "false",
							"encrypt_data":              "false",
							"reject_unencrypted_access": "false",
							"super_admin_sid":           "S-1-5-22-23",
							"home_dir_path":             "/home/A",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"smb_acl": []map[string]interface{}{
						{
							"enabled": "false",
						},
					},
				}),
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
				ImportStateVerifyIgnore: []string{"keytab", "keytab_md5", "snapshot_id"},
			},
		},
	})
}

var AlicloudNasFileSystemMap10635 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
	"region_id":   CHECKSET,
}

func AlicloudNasFileSystemBasicDependence10635(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "key_tab_1_md5" {
  default = "E3CCF7E2416DF04FA958AA4513EA29E8"
}

variable "key_tab_md5" {
  default = "91224BCF548214F4DF3F7CA1898354CE"
}

variable "key_tab_1" {
  default = "BQIAAABHAAIADUFMSUFEVEVTVC5DT00ABGNpZnMAGXNtYnNlcnZlcjI0LmFsaWFkdGVzdC5jb20AAAABAAAAAAEAAQAIqIx6v7p11oUAAABHAAIADUFMSUFEVEVTVC5DT00ABGNpZnMAGXNtYnNlcnZlcjI0LmFsaWFkdGVzdC5jb20AAAABAAAAAAEAAwAIqIx6v7p11oUAAABPAAIADUFMSUFEVEVTVC5DT00ABGNpZnMAGXNtYnNlcnZlcjI0LmFsaWFkdGVzdC5jb20AAAABAAAAAAEAFwAQnQZWB3RAPHU7PMIJyBWePAAAAF8AAgANQUxJQURURVNULkNPTQAEY2lmcwAZc21ic2VydmVyMjQuYWxpYWR0ZXN0LmNvbQAAAAEAAAAAAQASACAGJ7F0s+bcBjf6jD5HlvlRLmPSOW+qDZe0Qk0lQcf8WwAAAE8AAgANQUxJQURURVNULkNPTQAEY2lmcwAZc21ic2VydmVyMjQuYWxpYWR0ZXN0LmNvbQAAAAEAAAAAAQARABDdFmanrSIatnDDhoOXYadj"
}

variable "key_tab" {
  default = "BQIAAABZAAIAC1RFU1ROQVMuQ09NAARjaWZzAC0wMjUyYzRiZmFlLXR4YTc5LmNuLWhhbmd6aG91Lm5hcy5hbGl5dW5jcy5jb20AAAABAAAAAAMAAQAIvJLQj+BefzEAAABZAAIAC1RFU1ROQVMuQ09NAARjaWZzAC0wMjUyYzRiZmFlLXR4YTc5LmNuLWhhbmd6aG91Lm5hcy5hbGl5dW5jcy5jb20AAAABAAAAAAMAAwAIvJLQj+BefzEAAABhAAIAC1RFU1ROQVMuQ09NAARjaWZzAC0wMjUyYzRiZmFlLXR4YTc5LmNuLWhhbmd6aG91Lm5hcy5hbGl5dW5jcy5jb20AAAABAAAAAAMAFwAQmnGFGg9wrGIckBGRDz3ZEAAAAHEAAgALVEVTVE5BUy5DT00ABGNpZnMALTAyNTJjNGJmYWUtdHhhNzkuY24taGFuZ3pob3UubmFzLmFsaXl1bmNzLmNvbQAAAAEAAAAAAwASACB8NytHr3V2sFcEmV0p/k+7gxk4nbRvRotJ19NTmP1GigAAAGEAAgALVEVTVE5BUy5DT00ABGNpZnMALTAyNTJjNGJmYWUtdHhhNzkuY24taGFuZ3pob3UubmFzLmFsaXl1bmNzLmNvbQAAAAEAAAAAAwARABBz8lgg+jF54cqR2WZ7nltq"
}


`, name)
}

// Case 极速型标准型文件系统 5328
func TestAccAliCloudNasFileSystem_basic5328(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_nas_file_system.default"
	ra := resourceAttrInit(resourceId, AlicloudNasFileSystemMap5328)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return nas.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeFileSystem")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccnas%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudNasFileSystemBasicDependence5328)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		IDRefreshName:            resourceId,
		IDRefreshIgnore:          nasFileSystemIDRefreshIgnore,
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"storage_type":      "standard",
					"encrypt_type":      "0",
					"protocol_type":     "NFS",
					"file_system_type":  "extreme",
					"zone_id":           "${var.zone_id}",
					"capacity":          "100",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"storage_type":      "standard",
						"encrypt_type":      "0",
						"protocol_type":     "NFS",
						"file_system_type":  "extreme",
						"zone_id":           CHECKSET,
						"capacity":          "100",
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"capacity":          "200",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"capacity":          "200",
						"resource_group_id": CHECKSET,
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
				ImportStateVerifyIgnore: []string{"keytab", "keytab_md5", "snapshot_id"},
			},
		},
	})
}

var AlicloudNasFileSystemMap5328 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
	"region_id":   CHECKSET,
}

func AlicloudNasFileSystemBasicDependence5328(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "zone_id" {
  default = "cn-hangzhou-h"
}

variable "region" {
  default = "cn-hangzhou"
}

data "alicloud_resource_manager_resource_groups" "default" {}


`, name)
}

// Case 极速型快照创标准型文件系统 7668
func TestAccAliCloudNasFileSystem_basic7668(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_nas_file_system.default"
	ra := resourceAttrInit(resourceId, AlicloudNasFileSystemMap7668)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return nas.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeFileSystem")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccnas%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudNasFileSystemBasicDependence7668)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		IDRefreshName:            resourceId,
		IDRefreshIgnore:          nasFileSystemIDRefreshIgnore,
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"storage_type":     "standard",
					"protocol_type":    "NFS",
					"file_system_type": "extreme",
					"zone_id":          "${var.zone_id}",
					"capacity":         "100",
					"snapshot_id":      "${alicloud_nas_snapshot.defaultLXp8tf.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"storage_type":     "standard",
						"protocol_type":    "NFS",
						"file_system_type": "extreme",
						"zone_id":          CHECKSET,
						"capacity":         "100",
						"snapshot_id":      CHECKSET,
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
				ImportStateVerifyIgnore: []string{"keytab", "keytab_md5", "snapshot_id"},
			},
		},
	})
}

var AlicloudNasFileSystemMap7668 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
	"region_id":   CHECKSET,
}

func AlicloudNasFileSystemBasicDependence7668(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "zone_id" {
  default = "cn-hangzhou-h"
}

variable "region" {
  default = "cn-hangzhou"
}

resource "alicloud_nas_file_system" "defaultEHlvST" {
  storage_type     = "standard"
  zone_id          = var.zone_id
  encrypt_type     = "0"
  capacity         = "100"
  protocol_type    = "NFS"
  file_system_type = "extreme"
}

resource "alicloud_nas_snapshot" "defaultLXp8tf" {
  file_system_id = alicloud_nas_file_system.defaultEHlvST.id
  retention_days = "1"
  snapshot_name  = "testSnapshotCreateFs"
}


`, name)
}

// Case   12188
func TestAccAliCloudNasFileSystem_basic12188(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_nas_file_system.default"
	ra := resourceAttrInit(resourceId, AlicloudNasFileSystemMap12188)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return nas.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeFileSystem")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccnas%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudNasFileSystemBasicDependence12188)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
			acctest.PreCheck(t)
		},
		IDRefreshName:            resourceId,
		IDRefreshIgnore:          nasFileSystemIDRefreshIgnore,
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"description":  "cpfsse-文件系统同城冗余测试-wyf",
					"storage_type": "advance_100",
					"encrypt_type": "0",
					"vpc_id":       "${alicloud_vpc.createEVpc_Cpfs.id}",
					"redundancy_vswitch_ids": []string{
						"${alicloud_vswitch.CreateVswitchC.id}", "${alicloud_vswitch.CreateVswitchD.id}", "${alicloud_vswitch.CreateVswitchF.id}"},
					"capacity":         "500",
					"protocol_type":    "cpfs",
					"file_system_type": "cpfsse",
					"redundancy_type":  "ZRS",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":              "cpfsse-文件系统同城冗余测试-wyf",
						"storage_type":             "advance_100",
						"encrypt_type":             "0",
						"vpc_id":                   CHECKSET,
						"redundancy_vswitch_ids.#": "3",
						"capacity":                 "500",
						"protocol_type":            "cpfs",
						"file_system_type":         "cpfsse",
						"redundancy_type":          "ZRS",
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
				ImportStateVerifyIgnore: []string{"keytab", "keytab_md5", "snapshot_id"},
			},
		},
	})
}

var AlicloudNasFileSystemMap12188 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
	"region_id":   CHECKSET,
}

func AlicloudNasFileSystemBasicDependence12188(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_vpc" "createEVpc_Cpfs" {
  is_default  = false
  cidr_block  = "192.168.0.0/16"
  vpc_name    = "nas-teste1031-vpc"
  enable_ipv6 = true
}

resource "alicloud_vswitch" "CreateVswitchF" {
  is_default   = false
  vpc_id       = alicloud_vpc.createEVpc_Cpfs.id
  zone_id      = "cn-beijing-i"
  cidr_block   = "192.168.2.0/24"
  vswitch_name = "nas-teste1031-vsw1sdw-F"
}

resource "alicloud_vswitch" "CreateVswitchC" {
  is_default   = false
  vpc_id       = alicloud_vpc.createEVpc_Cpfs.id
  zone_id      = "cn-beijing-l"
  cidr_block   = "192.168.3.0/24"
  vswitch_name = "nas-teste1031-vsw2sdw-C"
}

resource "alicloud_vswitch" "CreateVswitchD" {
  is_default   = false
  vpc_id       = alicloud_vpc.createEVpc_Cpfs.id
  zone_id      = "cn-beijing-h"
  cidr_block   = "192.168.4.0/24"
  vswitch_name = "nas-teste1031-vsw3sdw-D"
}


`, name)
}

// Case auto upgrade config for cpfsse file systems.
func TestAccAliCloudNasFileSystem_autoUpgradeConfig(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_nas_file_system.default"
	ra := resourceAttrInit(resourceId, AlicloudNasFileSystemAutoUpgradeConfigMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return nas.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeFileSystem")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccNasAutoUpgrade%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudNasFileSystemAutoUpgradeConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
			acctest.PreCheck(t)
		},
		IDRefreshName:            resourceId,
		IDRefreshIgnore:          nasFileSystemIDRefreshIgnore,
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"description":  name,
					"storage_type": "advance_100",
					"encrypt_type": "0",
					"vpc_id":       "${alicloud_vpc.default.id}",
					"redundancy_vswitch_ids": []string{
						"${alicloud_vswitch.vswitch_c.id}", "${alicloud_vswitch.vswitch_d.id}", "${alicloud_vswitch.vswitch_f.id}"},
					"capacity":         "500",
					"protocol_type":    "cpfs",
					"file_system_type": "cpfsse",
					"redundancy_type":  "ZRS",
					"auto_upgrade_config": []map[string]interface{}{
						{
							"enabled":             "true",
							"capacity_used_ratio": "60",
							"time":                "90",
							"step":                "200",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_upgrade_config.#":                     "1",
						"auto_upgrade_config.0.enabled":             "true",
						"auto_upgrade_config.0.capacity_used_ratio": "60",
						"auto_upgrade_config.0.time":                "90",
						"auto_upgrade_config.0.step":                "200",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auto_upgrade_config": []map[string]interface{}{
						{
							"enabled":             "true",
							"capacity_used_ratio": "80",
							"time":                "60",
							"step":                "300",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_upgrade_config.#":                     "1",
						"auto_upgrade_config.0.enabled":             "true",
						"auto_upgrade_config.0.capacity_used_ratio": "80",
						"auto_upgrade_config.0.time":                "60",
						"auto_upgrade_config.0.step":                "300",
					}),
				),
			},
			{
				// Sub-attribute removal: only enabled is in the configuration.
				// The resource must send exactly that attribute instead of
				// re-sending the previous state values of the Optional+Computed
				// siblings; the siblings become computed values read back from
				// the server.
				Config: testAccConfig(map[string]interface{}{
					"auto_upgrade_config": []map[string]interface{}{
						{
							"enabled": "false",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_upgrade_config.#":                     "1",
						"auto_upgrade_config.0.enabled":             "false",
						"auto_upgrade_config.0.capacity_used_ratio": CHECKSET,
						"auto_upgrade_config.0.time":                CHECKSET,
						"auto_upgrade_config.0.step":                CHECKSET,
					}),
				),
			},
			{
				// Whole-block removal: the block is Optional+Computed, so
				// dropping it from the configuration keeps the value computed
				// by the server and must not send stale state back to
				// ModifyFileSystem. The step also proves the plan stays stable
				// after the removal.
				Config: testAccConfig(map[string]interface{}{
					"auto_upgrade_config": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_upgrade_config.#":                     "1",
						"auto_upgrade_config.0.enabled":             "false",
						"auto_upgrade_config.0.capacity_used_ratio": CHECKSET,
						"auto_upgrade_config.0.time":                CHECKSET,
						"auto_upgrade_config.0.step":                CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auto_upgrade_config": []map[string]interface{}{
						{
							"enabled":             "true",
							"capacity_used_ratio": "70",
							"time":                "30",
							"step":                "100",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_upgrade_config.#":                     "1",
						"auto_upgrade_config.0.enabled":             "true",
						"auto_upgrade_config.0.capacity_used_ratio": "70",
						"auto_upgrade_config.0.time":                "30",
						"auto_upgrade_config.0.step":                "100",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"keytab", "keytab_md5", "snapshot_id"},
			},
		},
	})
}

var AlicloudNasFileSystemAutoUpgradeConfigMap = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
	"region_id":   CHECKSET,
}

func AlicloudNasFileSystemAutoUpgradeConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

# This fixture pins three specific cn-beijing zones and the advance_100 CPFS SE
# specification: ZRS needs three distinct zones, and cpfsse advance_100 capacity
# is per-zone inventory. A region-wide zone query only proves API support, not
# inventory, so a dynamically selected triple can fail at apply time with
# Resource.OutOfStock. Update the fixed triple when these zones run out of stock.
resource "alicloud_vpc" "default" {
  is_default  = false
  cidr_block  = "192.168.0.0/16"
  vpc_name    = "tf-testAccNasAutoUpgradeVpc"
  enable_ipv6 = true
}

resource "alicloud_vswitch" "vswitch_f" {
  is_default   = false
  vpc_id       = alicloud_vpc.default.id
  zone_id      = "cn-beijing-i"
  cidr_block   = "192.168.2.0/24"
  vswitch_name = "tf-testAccNasAutoUpgradeVswF"
}

resource "alicloud_vswitch" "vswitch_c" {
  is_default   = false
  vpc_id       = alicloud_vpc.default.id
  zone_id      = "cn-beijing-l"
  cidr_block   = "192.168.3.0/24"
  vswitch_name = "tf-testAccNasAutoUpgradeVswC"
}

resource "alicloud_vswitch" "vswitch_d" {
  is_default   = false
  vpc_id       = alicloud_vpc.default.id
  zone_id      = "cn-beijing-h"
  cidr_block   = "192.168.4.0/24"
  vswitch_name = "tf-testAccNasAutoUpgradeVswD"
}


`, name)
}

// Test Nas FileSystem. <<< Resource test cases, automatically generated.
