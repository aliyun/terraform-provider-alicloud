package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/resourcemanager"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudResourceManagerFolder_basic(t *testing.T) {
	var v resourcemanager.Folder
	resourceId := "alicloud_resource_manager_folder.default"
	ra := resourceAttrInit(resourceId, ResourceManagerFolderMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ResourcemanagerService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeResourceManagerFolder")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, ResourceManagerFolderBasicdependence)
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
					"folder_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"folder_name": name,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"folder_name": "tf-test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"folder_name": "tf-test",
					}),
				),
			},
		},
	})
}

var ResourceManagerFolderMap = map[string]string{
	"parent_folder_id": CHECKSET,
}

func ResourceManagerFolderBasicdependence(name string) string {
	return ""
}
