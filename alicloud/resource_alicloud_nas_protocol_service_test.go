// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Nas ProtocolService. >>> Resource test cases, automatically generated.
// Case create_protocol_service_General 12170
func TestAccAliCloudNasProtocolService_basic12170(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_nas_protocol_service.default"
	ra := resourceAttrInit(resourceId, AlicloudNasProtocolServiceMap12170)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &NasServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeNasProtocolService")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccnas%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudNasProtocolServiceBasicDependence12170)
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
					"vpc_id":         "${alicloud_vpc.createEVpc_Cpfs1.id}",
					"protocol_type":  "NFS",
					"protocol_spec":  "General",
					"vswitch_id":     "${alicloud_vswitch.CreateVswitch1.id}",
					"dry_run":        "false",
					"file_system_id": "${alicloud_nas_file_system.create_cpfs_file_system_General.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpc_id":         CHECKSET,
						"protocol_type":  "NFS",
						"protocol_spec":  "General",
						"vswitch_id":     CHECKSET,
						"dry_run":        "false",
						"file_system_id": CHECKSET,
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

var AlicloudNasProtocolServiceMap12170 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
}

func AlicloudNasProtocolServiceBasicDependence12170(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_vpc" "createEVpc_Cpfs1" {
  is_default  = false
  cidr_block  = "192.168.0.0/16"
  vpc_name    = "nas-teste1031-vpc"
  enable_ipv6 = true
}

resource "alicloud_vswitch" "CreateVswitch1" {
  is_default   = false
  vpc_id       = alicloud_vpc.createEVpc_Cpfs1.id
  zone_id      = "cn-beijing-i"
  cidr_block   = "192.168.2.0/24"
  vswitch_name = "nas-teste1031-vsw1sdw-F"
}

resource "alicloud_nas_file_system" "create_cpfs_file_system_General" {
  description      = "cpfs-文件系统本地冗余-protocol_service测试"
  storage_type     = "advance_100"
  zone_id          = "cn-beijing-i"
  encrypt_type     = "0"
  vpc_id           = alicloud_vpc.createEVpc_Cpfs1.id
  capacity         = "3600"
  protocol_type    = "cpfs"
  vswitch_id       = alicloud_vswitch.CreateVswitch1.id
  file_system_type = "cpfs"
}


`, name)
}

// Case create_protocol_service_CL2 12172
func TestAccAliCloudNasProtocolService_basic12172(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_nas_protocol_service.default"
	ra := resourceAttrInit(resourceId, AlicloudNasProtocolServiceMap12172)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &NasServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeNasProtocolService")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccnas%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudNasProtocolServiceBasicDependence12172)
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
					"protocol_throughput": "16000",
					"vpc_id":              "${alicloud_vpc.createEVpc_Cpfs1.id}",
					"protocol_type":       "NFS",
					"protocol_spec":       "CL2",
					"vswitch_id":          "${alicloud_vswitch.CreateVswitch1.id}",
					"dry_run":             "false",
					"file_system_id":      "${alicloud_nas_file_system.create_cpfs_file_system_CL2.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"protocol_throughput": "16000",
						"vpc_id":              CHECKSET,
						"protocol_type":       "NFS",
						"protocol_spec":       "CL2",
						"vswitch_id":          CHECKSET,
						"dry_run":             "false",
						"file_system_id":      CHECKSET,
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

var AlicloudNasProtocolServiceMap12172 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
}

func AlicloudNasProtocolServiceBasicDependence12172(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_vpc" "createEVpc_Cpfs1" {
  is_default  = false
  cidr_block  = "192.168.0.0/16"
  vpc_name    = "nas-teste1031-vpc"
  enable_ipv6 = true
}

resource "alicloud_vswitch" "CreateVswitch1" {
  is_default   = false
  vpc_id       = alicloud_vpc.createEVpc_Cpfs1.id
  zone_id      = "cn-beijing-i"
  cidr_block   = "192.168.2.0/24"
  vswitch_name = "nas-teste1031-vsw1sdw-F"
}

resource "alicloud_nas_file_system" "create_cpfs_file_system_CL2" {
  description      = "cpfs-文件系统本地冗余-protocol_service测试"
  storage_type     = "advance_100"
  zone_id          = "cn-beijing-i"
  encrypt_type     = "0"
  vpc_id           = alicloud_vpc.createEVpc_Cpfs1.id
  capacity         = "3600"
  protocol_type    = "cpfs"
  vswitch_id       = alicloud_vswitch.CreateVswitch1.id
  file_system_type = "cpfs"
}


`, name)
}

// Case create_protocol_service_CL1 12174
func TestAccAliCloudNasProtocolService_basic12174(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_nas_protocol_service.default"
	ra := resourceAttrInit(resourceId, AlicloudNasProtocolServiceMap12174)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &NasServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeNasProtocolService")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccnas%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudNasProtocolServiceBasicDependence12174)
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
					"protocol_throughput": "8000",
					"description":         "wyf协议服务描述test",
					"vpc_id":              "${alicloud_vpc.createEVpc_Cpfs1.id}",
					"protocol_type":       "NFS",
					"protocol_spec":       "CL1",
					"vswitch_id":          "${alicloud_vswitch.CreateVswitch1.id}",
					"dry_run":             "false",
					"file_system_id":      "${alicloud_nas_file_system.create_cpfs_file_system_CL1.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"protocol_throughput": "8000",
						"description":         "wyf协议服务描述test",
						"vpc_id":              CHECKSET,
						"protocol_type":       "NFS",
						"protocol_spec":       "CL1",
						"vswitch_id":          CHECKSET,
						"dry_run":             "false",
						"file_system_id":      CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "wyf更新协议服务描述信息test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "wyf更新协议服务描述信息test",
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

var AlicloudNasProtocolServiceMap12174 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
}

func AlicloudNasProtocolServiceBasicDependence12174(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_vpc" "createEVpc_Cpfs1" {
  is_default  = false
  cidr_block  = "192.168.0.0/16"
  vpc_name    = "nas-teste1031-vpc"
  enable_ipv6 = true
}

resource "alicloud_vswitch" "CreateVswitch1" {
  is_default   = false
  vpc_id       = alicloud_vpc.createEVpc_Cpfs1.id
  zone_id      = "cn-beijing-i"
  cidr_block   = "192.168.2.0/24"
  vswitch_name = "nas-teste1031-vsw1sdw-F"
}

resource "alicloud_nas_file_system" "create_cpfs_file_system_CL1" {
  description      = "cpfs-文件系统本地冗余-protocol_service测试"
  storage_type     = "advance_100"
  zone_id          = "cn-beijing-i"
  encrypt_type     = "0"
  vpc_id           = alicloud_vpc.createEVpc_Cpfs1.id
  capacity         = "3600"
  protocol_type    = "cpfs"
  vswitch_id       = alicloud_vswitch.CreateVswitch1.id
  file_system_type = "cpfs"
}


`, name)
}

// Case create_cpfsse_protocol_service_General 12186
func TestAccAliCloudNasProtocolService_basic12186(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_nas_protocol_service.default"
	ra := resourceAttrInit(resourceId, AlicloudNasProtocolServiceMap12186)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &NasServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeNasProtocolService")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccnas%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudNasProtocolServiceBasicDependence12186)
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
					"protocol_type":  "NFS",
					"protocol_spec":  "General",
					"dry_run":        "false",
					"file_system_id": "${alicloud_nas_file_system.create_cpfsse_file_system_General.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"protocol_type":  "NFS",
						"protocol_spec":  "General",
						"dry_run":        "false",
						"file_system_id": CHECKSET,
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

var AlicloudNasProtocolServiceMap12186 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
}

func AlicloudNasProtocolServiceBasicDependence12186(name string) string {
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

resource "alicloud_nas_file_system" "create_cpfsse_file_system_General" {
  description      = "峰-cpfsse-ZRS-protocol_service测试"
  storage_type     = "advance_100"
  encrypt_type     = "0"
  vpc_id           = alicloud_vpc.createEVpc_Cpfs.id
  redundancy_vswitch_ids = [alicloud_vswitch.CreateVswitchC.id, alicloud_vswitch.CreateVswitchD.id, alicloud_vswitch.CreateVswitchF.id]
  capacity         = "500"
  protocol_type    = "cpfs"
  file_system_type = "cpfsse"
  redundancy_type  = "ZRS"
}


`, name)
}

// Test Nas ProtocolService. <<< Resource test cases, automatically generated.
