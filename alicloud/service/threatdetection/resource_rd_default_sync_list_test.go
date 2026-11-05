package threatdetection_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/threatdetection"
	tFacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test ThreatDetection RdDefaultSyncList. >>> Resource test cases.
// lintignore: AT001
func TestAccAliCloudThreatDetectionRdDefaultSyncList_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_threat_detection_rd_default_sync_list.default"
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return threatdetection.NewServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeThreatDetectionRdDefaultSyncList")
	rac := acctest.ResourceAttrCheckInit(rc, acctest.ResourceAttrInit(resourceId, map[string]string{}))
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tFacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sthreatdetectionrddefaultsynclist%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudThreatDetectionRdDefaultSyncListBasicDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				// Create: omit folder_ids so the singleton is adopted without
				// sending an empty FolderIds that would clear an existing
				// synchronization list. CreateRdDefaultSyncList only sends
				// FolderIds when the user declares it in the configuration,
				// and an account that never configured the list reads back
				// as an empty folder list (the ListRdDefaultSyncList
				// response carries no Data field in that case), not as a
				// missing resource.
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"folder_ids.#": "0",
					}),
				),
			},
			{
				// Create: set the default sync list to a single resource directory folder.
				Config: testAccConfig(map[string]interface{}{
					"folder_ids": []string{alicloudResourceManagerFolderRef("default")},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"folder_ids.#": "1",
					}),
				),
			},
			{
				// Update: expand the list to two resource directory folders
				// (FolderIds accepts multiple comma-separated folder ids).
				Config: testAccConfig(map[string]interface{}{
					"folder_ids": []string{
						alicloudResourceManagerFolderRef("default"),
						alicloudResourceManagerFolderRef("second"),
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"folder_ids.#": "2",
					}),
					resource.TestCheckResourceAttrPair(resourceId, "folder_ids.0", "alicloud_resource_manager_folder.default", "id"),
					resource.TestCheckResourceAttrPair(resourceId, "folder_ids.1", "alicloud_resource_manager_folder.second", "id"),
				),
			},
			{
				// Update: shrink the list back to a single folder (set/replace semantics).
				Config: testAccConfig(map[string]interface{}{
					"folder_ids": []string{alicloudResourceManagerFolderRef("default")},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"folder_ids.#": "1",
					}),
					resource.TestCheckResourceAttrPair(resourceId, "folder_ids.0", "alicloud_resource_manager_folder.default", "id"),
				),
			},
			{
				// Update: clear the list (empty FolderIds disables the auto-sync).
				Config: testAccConfig(map[string]interface{}{
					"folder_ids": []string{},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"folder_ids.#": "0",
					}),
				),
			},
			{
				// Update: re-set the list (set/replace semantics).
				Config: testAccConfig(map[string]interface{}{
					"folder_ids": []string{alicloudResourceManagerFolderRef("default")},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"folder_ids.#": "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

func alicloudResourceManagerFolderRef(name string) string {
	return fmt.Sprintf("${alicloud_resource_manager_folder.%s.id}", name)
}

func AlicloudThreatDetectionRdDefaultSyncListBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alicloud_resource_manager_folder" "default" {
  folder_name = "tf-testacc-rd-sync-folder-${var.name}"
}

resource "alicloud_resource_manager_folder" "second" {
  folder_name = "tf-testacc-rd-sync-folder2-${var.name}"
}
`, name)
}

// Test ThreatDetection RdDefaultSyncList. <<< Resource test cases.
