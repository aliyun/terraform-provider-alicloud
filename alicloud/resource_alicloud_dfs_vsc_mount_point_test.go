package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Dfs VscMountPoint. >>> Resource test cases, automatically generated.
// Case VscMountPoint资源测试用例 5268
func TestAccAliCloudDfsVscMountPoint_basic5268(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_dfs_vsc_mount_point.default"
	ra := resourceAttrInit(resourceId, AlicloudDfsVscMountPointMap5268)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DfsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDfsVscMountPoint")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sdfsvscmountpoint%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudDfsVscMountPointBasicDependence5268)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{connectivity.Hangzhou})
			testAccPreCheckWithTime(t, []int{1})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"file_system_id": "${alicloud_dfs_file_system.DefaultFsForRMCVscMp.id}",
					"alias_prefix":   "VscMpRMCTestAlias656",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"file_system_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "VscMpRMCTestAlias656",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "VscMpRMCTestAlias656",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "VscMpRMCTestAliasUpdate",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "VscMpRMCTestAliasUpdate",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"file_system_id": "${alicloud_dfs_file_system.DefaultFsForRMCVscMp.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"file_system_id": CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"alias_prefix"},
			},
		},
	})
}

var AlicloudDfsVscMountPointMap5268 = map[string]string{
	"mount_point_id": CHECKSET,
}

func AlicloudDfsVscMountPointBasicDependence5268(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_dfs_zones" "default" {}

resource "alicloud_dfs_file_system" "DefaultFsForRMCVscMp" {
  space_capacity       = "1024"
  description          = "for vsc mountpoint RMC test"
  storage_type         = "STANDARD"
  zone_id              = data.alicloud_dfs_zones.default.zones.0.zone_id
  protocol_type        = "HDFS"
  data_redundancy_type = "LRS"
  file_system_name     = var.name

}


`, name)
}

// Case VscMountPoint资源测试用例 5268  twin
func TestAccAliCloudDfsVscMountPoint_basic5268_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_dfs_vsc_mount_point.default"
	ra := resourceAttrInit(resourceId, AlicloudDfsVscMountPointMap5268)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DfsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDfsVscMountPoint")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sdfsvscmountpoint%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudDfsVscMountPointBasicDependence5268)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{connectivity.Hangzhou})
			testAccPreCheckWithTime(t, []int{1})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"file_system_id": "${alicloud_dfs_file_system.DefaultFsForRMCVscMp.id}",
					"alias_prefix":   "VscMpRMCTestAlias656",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"file_system_id": CHECKSET,
						"alias_prefix":   CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"alias_prefix"},
			},
		},
	})
}

// Test Dfs VscMountPoint. <<< Resource test cases, automatically generated.
