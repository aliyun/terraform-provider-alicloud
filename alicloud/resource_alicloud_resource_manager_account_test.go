package alicloud

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/resourcemanager"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudResourceManagerAccount_basic(t *testing.T) {
	var v resourcemanager.Account
	resourceId := "alicloud_resource_manager_account.default"
	ra := resourceAttrInit(resourceId, ResourceManagerAccountMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ResourcemanagerService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeResourceManagerAccount")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, ResourceManagerAccountBasicdependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			// The created account resource cannot be deleted, and the created dependent resource folder will also be deleted. Therefore, the existing folder is specified in the environment variable for testing. If it is not specified, skip the test.
			testAccPreCheckWithResourceManagerFloderIdSetting(t)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"display_name": name,
					"folder_id":    os.Getenv("ALICLOUD_RESOURCE_MANAGER_FOLDER_ID1"),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"display_name": name,
						"folder_id":    os.Getenv("ALICLOUD_RESOURCE_MANAGER_FOLDER_ID1"),
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
					"display_name": "tf-1233",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"display_name": "tf-1233",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"folder_id": os.Getenv("ALICLOUD_RESOURCE_MANAGER_FOLDER_ID2"),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"folder_id": os.Getenv("ALICLOUD_RESOURCE_MANAGER_FOLDER_ID2"),
					}),
				),
			},
		},
	})
}

var ResourceManagerAccountMap = map[string]string{
	"join_method":           CHECKSET,
	"join_time":             CHECKSET,
	"modify_time":           CHECKSET,
	"resource_directory_id": CHECKSET,
	"status":                CHECKSET,
	"type":                  CHECKSET,
}

func ResourceManagerAccountBasicdependence(name string) string {
	return ""
}
