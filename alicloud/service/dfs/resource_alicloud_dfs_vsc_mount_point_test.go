package dfs_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/dfs"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test Dfs VscMountPoint. >>> Resource test cases, automatically generated.
// Case VscMountPoint资源测试用例 5268
func TestAccAliCloudDfsVscMountPoint_basic5268(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_dfs_vsc_mount_point.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudDfsVscMountPointMap5268)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return dfs.NewServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeDfsVscMountPoint")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testaccdfsvscmountpoint%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudDfsVscMountPointBasicDependence5268)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{connectivity.Hangzhou})
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"file_system_id": "${alicloud_dfs_file_system.DefaultFsForRMCVscMp.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"file_system_id": acctest.CHECKSET,
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
	"mount_point_id": acctest.CHECKSET,
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
	ra := acctest.ResourceAttrInit(resourceId, AliCloudDfsVscMountPointMap5268)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return dfs.NewServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeDfsVscMountPoint")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testaccdfsvscmountpoint%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudDfsVscMountPointBasicDependence5268)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{connectivity.Hangzhou})
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"file_system_id": "${alicloud_dfs_file_system.DefaultFsForRMCVscMp.id}",
					"alias_prefix":   name,
					"description":    name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"file_system_id": acctest.CHECKSET,
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
