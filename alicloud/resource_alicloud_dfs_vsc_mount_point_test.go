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
	ra := resourceAttrInit(resourceId, AliCloudDfsVscMountPointMap5268)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DfsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDfsVscMountPoint")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testaccdfsvscmountpoint%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudDfsVscMountPointBasicDependence5268)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{connectivity.Hangzhou})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
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
				Config: testAccConfig(map[string]interface{}{
					"alias_prefix": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"alias_prefix": name,
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
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

var AliCloudDfsVscMountPointMap5268 = map[string]string{
	"mount_point_id": CHECKSET,
}

func AliCloudDfsVscMountPointBasicDependence5268(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	resource "alicloud_dfs_file_system" "DefaultFsForRMCVscMp" {
  		space_capacity       = "1024"
  		description          = "for vsc mountpoint RMC test"
  		storage_type         = "PERFORMANCE"
  		zone_id              = "cn-hangzhou-b"
  		protocol_type        = "PANGU"
  		data_redundancy_type = "LRS"
  		file_system_name     = var.name
	}
`, name)
}

// Case VscMountPoint资源测试用例 5268  twin
func TestAccAliCloudDfsVscMountPoint_basic5268_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_dfs_vsc_mount_point.default"
	ra := resourceAttrInit(resourceId, AliCloudDfsVscMountPointMap5268)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DfsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDfsVscMountPoint")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testaccdfsvscmountpoint%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudDfsVscMountPointBasicDependence5268)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{connectivity.Hangzhou})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"file_system_id": "${alicloud_dfs_file_system.DefaultFsForRMCVscMp.id}",
					"alias_prefix":   name,
					"description":    name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"file_system_id": CHECKSET,
						"alias_prefix":   name,
						"description":    name,
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

// Test Dfs VscMountPoint. <<< Resource test cases, automatically generated.
