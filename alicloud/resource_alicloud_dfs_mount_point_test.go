package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudDFSMountPoint_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_dfs_mount_point.default"
	checkoutSupportedRegions(t, true, connectivity.DfsSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudDFSMountPointMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DfsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDfsMountPoint")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sdfsmountpoint%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudDFSMountPointBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"network_type":    "VPC",
					"vpc_id":          "${local.vpc_id}",
					"vswitch_id":      "${local.vswitch_id}",
					"file_system_id":  "${alicloud_dfs_file_system.default.id}",
					"description":     name,
					"access_group_id": "${alicloud_dfs_access_group.default[0].id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"network_type":    "VPC",
						"vpc_id":          CHECKSET,
						"vswitch_id":      CHECKSET,
						"file_system_id":  CHECKSET,
						"description":     name,
						"access_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name + "update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"access_group_id": "${alicloud_dfs_access_group.default[1].id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"access_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Inactive",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Inactive",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status":          "Active",
					"description":     name,
					"access_group_id": "${alicloud_dfs_access_group.default[0].id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status":          "Active",
						"description":     name,
						"access_group_id": CHECKSET,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

var AlicloudDFSMountPointMap0 = map[string]string{}

func AlicloudDFSMountPointBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}

data "alicloud_dfs_zones" "default" {}

data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_dfs_zones.default.zones.0.zone_id
}

resource "alicloud_vswitch" "default" {
  count        = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
  vpc_id       = data.alicloud_vpcs.default.ids.0
  cidr_block   = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 2)
  zone_id      = data.alicloud_dfs_zones.default.zones.0.zone_id
  vswitch_name = var.name
}

locals {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.default.*.id, [""])[0]
}

resource "alicloud_dfs_file_system" "default" {
  storage_type     = data.alicloud_dfs_zones.default.zones.0.options.0.storage_type
  zone_id          = data.alicloud_dfs_zones.default.zones.0.zone_id
  protocol_type    = "HDFS"
  description      = var.name
  file_system_name = var.name
  throughput_mode  = "Standard"
  space_capacity   = "1024"
}

resource "alicloud_dfs_access_group" "default" {
  count             = 2
  network_type      = "VPC"
  access_group_name = join("", [var.name, count.index])
  description       = var.name
}


`, name)
}
