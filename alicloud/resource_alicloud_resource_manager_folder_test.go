package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/resourcemanager"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_resource_manager_folder", &resource.Sweeper{
		Name: "alicloud_resource_manager_folder",
		F:    testSweepResourceManagerFolder,
	})
}

func testSweepResourceManagerFolder(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "Error getting Alicloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)
	resourceManagerService := ResourcemanagerService{client}

	prefixes := []string{
		"tf-testAcc",
	}

	request := resourcemanager.CreateListFoldersForParentRequest()
	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)
	var folderIds []string
	for {
		raw, err := resourceManagerService.client.WithResourcemanagerClient(func(resourceManagerClient *resourcemanager.Client) (interface{}, error) {
			return resourceManagerClient.ListFoldersForParent(request)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to retrieve resoure manager folder in service list: %s", err)
		}

		response, _ := raw.(*resourcemanager.ListFoldersForParentResponse)

		for _, v := range response.Folders.Folder {
			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(v.FolderName), strings.ToLower(prefix)) {
					skip = false
				}
			}
			if skip {
				log.Printf("[INFO] Skipping resource manager folder: %s ", v.FolderName)
			} else {
				folderIds = append(folderIds, v.FolderId)
			}
		}
		if len(response.Folders.Folder) < PageSizeLarge {
			break
		}
		page, err := getNextpageNumber(request.PageNumber)
		if err != nil {
			return WrapError(err)
		}
		request.PageNumber = page
	}

	for _, folderId := range folderIds {
		log.Printf("[INFO] Delete resource manager folder: %s ", folderId)

		request := resourcemanager.CreateDeleteFolderRequest()

		_, err = resourceManagerService.client.WithResourcemanagerClient(func(resourceManagerClient *resourcemanager.Client) (interface{}, error) {
			return resourceManagerClient.DeleteFolder(request)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete resource manager folder(%s): %s", folderId, err)
		}
	}
	return nil
}

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
	name := fmt.Sprintf("tf-testAcc-%d", rand)
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
					"folder_name": "tf-testAccFolder-change",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"folder_name": "tf-testAccFolder-change",
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
