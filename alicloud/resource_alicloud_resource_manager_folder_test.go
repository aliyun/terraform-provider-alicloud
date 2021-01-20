package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"

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

	prefixes := []string{
		"tf-testAcc",
	}

	action := "ListFoldersForParent"
	request := make(map[string]interface{})

	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1

	var response map[string]interface{}
	conn, err := client.NewResourcemanagerClient()
	if err != nil {
		return WrapError(err)
	}

	var folderIds []string
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-03-31"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			log.Printf("[ERROR] Failed to retrieve resoure manager folder in service list: %s", err)
		}

		resp, err := jsonpath.Get("$.Folders.Folder", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Folders.Folder", response)
		}
		for _, v := range resp.([]interface{}) {
			item := v.(map[string]interface{})
			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(item["FolderName"].(string)), strings.ToLower(prefix)) {
					skip = false
				}
			}
			if skip {
				log.Printf("[INFO] Skipping resource manager folder: %s ", item["FolderName"].(string))
			} else {
				folderIds = append(folderIds, item["FolderId"].(string))
			}
		}
		if len(resp.([]interface{})) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}

	for _, folderId := range folderIds {
		log.Printf("[INFO] Delete resource manager folder: %s ", folderId)

		request := map[string]interface{}{
			"FolderId": folderId,
		}

		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-03-31"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete resource manager folder(%s): %s", folderId, err)
		}
	}
	return nil
}

func TestAccAlicloudResourceManagerFolder_basic(t *testing.T) {
	var v map[string]interface{}
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
