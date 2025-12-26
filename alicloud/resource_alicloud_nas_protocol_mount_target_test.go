package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Case create_protocol_mount_target_with_fileset 12169
func TestAccAliCloudNasProtocolMountTarget_basic12169(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_nas_protocol_mount_target.default"
	ra := resourceAttrInit(resourceId, AlicloudNasProtocolMountTargetMap12169)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &NasServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeNasProtocolMountTarget")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccnas%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudNasProtocolMountTargetBasicDependence12169)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"fset_id":             "${alicloud_nas_fileset.cpfs_LRS_fileset01.fileset_id}",
					"description":         "wyf通过Fileset创建的协议服务挂载点test",
					"vpc_id":              "${alicloud_vpc.createEVpc_Cpfs.id}",
					"vswitch_id":          "${alicloud_vswitch.CreateVswitch1.id}",
					"access_group_name":   "DEFAULT_VPC_GROUP_NAME",
					"dry_run":             "false",
					"file_system_id":      "${alicloud_nas_file_system.create_cpfs_file_system_General01.id}",
					"protocol_service_id": "${alicloud_nas_protocol_service.create_protocol_service_General01.protocol_service_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"fset_id":             CHECKSET,
						"description":         "wyf通过Fileset创建的协议服务挂载点test",
						"vpc_id":              CHECKSET,
						"vswitch_id":          CHECKSET,
						"access_group_name":   CHECKSET,
						"dry_run":             "false",
						"file_system_id":      CHECKSET,
						"protocol_service_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "wyf更新挂载点描述信息test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "wyf更新挂载点描述信息test",
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
				ImportStateVerifyIgnore: []string{"dry_run"},
			},
		},
	})
}

var AlicloudNasProtocolMountTargetMap12169 = map[string]string{
	"status":      CHECKSET,
	"export_id":   CHECKSET,
	"create_time": CHECKSET,
}

func AlicloudNasProtocolMountTargetBasicDependence12169(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_vpc" "createEVpc_Cpfs" {
  is_default  = false
  cidr_block  = "192.168.0.0/16"
  vpc_name    = "nas-teste1223-vpc"
  enable_ipv6 = true
}

resource "alicloud_vswitch" "CreateVswitch1" {
  is_default   = false
  vpc_id       = alicloud_vpc.createEVpc_Cpfs.id
  zone_id      = "cn-beijing-i"
  cidr_block   = "192.168.3.0/24"
  vswitch_name = "nas-teste1223-vsw2sdw-C"
}

resource "alicloud_nas_file_system" "create_cpfs_file_system_General01" {
  description      = "cpfs-文件系统本地冗余-ProtocolMountTarget-fileset测试"
  storage_type     = "advance_100"
  zone_id          = "cn-beijing-i"
  vpc_id           = alicloud_vpc.createEVpc_Cpfs.id
  capacity         = "3600"
  protocol_type    = "cpfs"
  vswitch_id       = alicloud_vswitch.CreateVswitch1.id
  file_system_type = "cpfs"
}

resource "alicloud_nas_protocol_service" "create_protocol_service_General01" {
  vpc_id         = alicloud_vpc.createEVpc_Cpfs.id
  protocol_type  = "NFS"
  protocol_spec  = "General"
  vswitch_id     = alicloud_vswitch.CreateVswitch1.id
  dry_run        = false
  file_system_id = alicloud_nas_file_system.create_cpfs_file_system_General01.id
}

resource "alicloud_nas_fileset" "cpfs_LRS_fileset01" {
  file_system_path = "/testfileset/"
  description      = "cpfs-LRS-fileset测试-wyf"
  file_system_id   = alicloud_nas_file_system.create_cpfs_file_system_General01.id
}


`, name)
}

// Case create_protocol_mount_target_with_path 12171
func TestAccAliCloudNasProtocolMountTarget_basic12171(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_nas_protocol_mount_target.default"
	ra := resourceAttrInit(resourceId, AlicloudNasProtocolMountTargetMap12171)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &NasServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeNasProtocolMountTarget")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccnas%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudNasProtocolMountTargetBasicDependence12171)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"path":                "/",
					"description":         "wyf服务协议目录导出test",
					"vpc_id":              "${alicloud_vpc.createEVpc_Cpfs.id}",
					"vswitch_id":          "${alicloud_vswitch.CreateVswitch1.id}",
					"access_group_name":   "DEFAULT_VPC_GROUP_NAME",
					"dry_run":             "false",
					"file_system_id":      "${alicloud_nas_file_system.create_cpfs_file_system_General02.id}",
					"protocol_service_id": "${alicloud_nas_protocol_service.create_protocol_service_General02.protocol_service_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"path":                "/",
						"description":         "wyf服务协议目录导出test",
						"vpc_id":              CHECKSET,
						"vswitch_id":          CHECKSET,
						"access_group_name":   CHECKSET,
						"dry_run":             "false",
						"file_system_id":      CHECKSET,
						"protocol_service_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "wyf更新服务协议目录导出描述信息test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "wyf更新服务协议目录导出描述信息test",
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
				ImportStateVerifyIgnore: []string{"dry_run"},
			},
		},
	})
}

var AlicloudNasProtocolMountTargetMap12171 = map[string]string{
	"status":      CHECKSET,
	"export_id":   CHECKSET,
	"create_time": CHECKSET,
}

func AlicloudNasProtocolMountTargetBasicDependence12171(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_vpc" "createEVpc_Cpfs" {
  is_default  = false
  cidr_block  = "192.168.0.0/16"
  vpc_name    = "nas-teste1223-vpc"
  enable_ipv6 = true
}

resource "alicloud_vswitch" "CreateVswitch1" {
  is_default   = false
  vpc_id       = alicloud_vpc.createEVpc_Cpfs.id
  zone_id      = "cn-beijing-i"
  cidr_block   = "192.168.3.0/24"
  vswitch_name = "nas-teste1223-vsw2sdw-C"
}

resource "alicloud_nas_file_system" "create_cpfs_file_system_General02" {
  description      = "cpfs-文件系统本地冗余-ProtocolMountTarget-path测试"
  storage_type     = "advance_100"
  zone_id          = "cn-beijing-i"
  vpc_id           = alicloud_vpc.createEVpc_Cpfs.id
  capacity         = "3600"
  protocol_type    = "cpfs"
  vswitch_id       = alicloud_vswitch.CreateVswitch1.id
  file_system_type = "cpfs"
}

resource "alicloud_nas_protocol_service" "create_protocol_service_General02" {
  vpc_id         = alicloud_vpc.createEVpc_Cpfs.id
  protocol_type  = "NFS"
  protocol_spec  = "General"
  vswitch_id     = alicloud_vswitch.CreateVswitch1.id
  dry_run        = false
  file_system_id = alicloud_nas_file_system.create_cpfs_file_system_General02.id
}


`, name)
}

// Case create_protocol_mount_target_with_fileset04 12183
func TestAccAliCloudNasProtocolMountTarget_basic12183(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_nas_protocol_mount_target.default"
	ra := resourceAttrInit(resourceId, AlicloudNasProtocolMountTargetMap12183)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &NasServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeNasProtocolMountTarget")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccnas%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudNasProtocolMountTargetBasicDependence12183)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"fset_id":     "${alicloud_nas_fileset.cpfsse_ZRS_fileset01.fileset_id}",
					"description": "wyf通过Fileset创建的协议服务挂载点test",
					"vpc_id":      "${alicloud_vpc.createEVpc_Cpfs.id}",
					"vswitch_ids": []string{
						"${alicloud_vswitch.CreateVswitchC.id}", "${alicloud_vswitch.CreateVswitchD.id}", "${alicloud_vswitch.CreateVswitchF.id}"},
					"access_group_name":   "DEFAULT_VPC_GROUP_NAME",
					"dry_run":             "false",
					"file_system_id":      "${alicloud_nas_file_system.create_cpfsse_file_system_ZRS.id}",
					"protocol_service_id": "${alicloud_nas_protocol_service.create_protocol_service_General04.protocol_service_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"fset_id":             CHECKSET,
						"description":         "wyf通过Fileset创建的协议服务挂载点test",
						"vpc_id":              CHECKSET,
						"vswitch_ids.#":       "3",
						"access_group_name":   CHECKSET,
						"dry_run":             "false",
						"file_system_id":      CHECKSET,
						"protocol_service_id": CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dry_run"},
			},
		},
	})
}

var AlicloudNasProtocolMountTargetMap12183 = map[string]string{
	"status":      CHECKSET,
	"export_id":   CHECKSET,
	"create_time": CHECKSET,
}

func AlicloudNasProtocolMountTargetBasicDependence12183(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_vpc" "createEVpc_Cpfs" {
  is_default  = false
  cidr_block  = "192.168.0.0/16"
  vpc_name    = "nas-teste1223-vpc"
  enable_ipv6 = true
}

resource "alicloud_vswitch" "CreateVswitchF" {
  is_default   = false
  vpc_id       = alicloud_vpc.createEVpc_Cpfs.id
  zone_id      = "cn-beijing-h"
  cidr_block   = "192.168.2.0/24"
  vswitch_name = "nas-teste1223-vsw1sdw-F"
}

resource "alicloud_vswitch" "CreateVswitchC" {
  is_default   = false
  vpc_id       = alicloud_vpc.createEVpc_Cpfs.id
  zone_id      = "cn-beijing-i"
  cidr_block   = "192.168.3.0/24"
  vswitch_name = "nas-teste1223-vsw2sdw-C"
}

resource "alicloud_vswitch" "CreateVswitchD" {
  is_default   = false
  vpc_id       = alicloud_vpc.createEVpc_Cpfs.id
  zone_id      = "cn-beijing-l"
  cidr_block   = "192.168.4.0/24"
  vswitch_name = "nas-teste1223-vsw3sdw-D"
}

resource "alicloud_nas_file_system" "create_cpfsse_file_system_ZRS" {
  description      = "cpfsse-文件系统ProtocolMountTarget-wyf"
  storage_type     = "advance_100"
  vpc_id           = alicloud_vpc.createEVpc_Cpfs.id
  capacity         = "500"
  protocol_type    = "cpfs"
  file_system_type = "cpfsse"
  encrypt_type     = "0"
  redundancy_vswitch_ids = [alicloud_vswitch.CreateVswitchC.id, alicloud_vswitch.CreateVswitchD.id, alicloud_vswitch.CreateVswitchF.id]
  redundancy_type  = "ZRS"
}

resource "alicloud_nas_fileset" "cpfsse_ZRS_fileset01" {
  file_system_path = "/testfileset/"
  description      = "cpfs-ZRS-fileset测试-wyf"
  file_system_id   = alicloud_nas_file_system.create_cpfsse_file_system_ZRS.id
}

resource "alicloud_nas_protocol_service" "create_protocol_service_General04" {
  protocol_type  = "NFS"
  protocol_spec  = "General"
  dry_run        = false
  file_system_id = alicloud_nas_file_system.create_cpfsse_file_system_ZRS.id
}


`, name)
}
