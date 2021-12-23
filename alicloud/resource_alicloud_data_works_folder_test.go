package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudDataWorksFolder_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_data_works_folder.default"
	ra := resourceAttrInit(resourceId, AlicloudDataWorksFolderMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DataworksPublicService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDataWorksFolder")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sdataworksfolder%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudDataWorksFolderBasicDependence0)
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
					"project_id":  "34051",
					"folder_path": "Business Flow/tfTestAcc/folderDi/tftest1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"project_id":  "34051",
						"folder_path": "Business Flow/tfTestAcc/folderDi/tftest1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"folder_path": "Business Flow/tfTestAcc/folderDi/tftest2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"folder_path": "Business Flow/tfTestAcc/folderDi/tftest2",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true, ImportStateVerifyIgnore: []string{"project_identifier"},
			},
		},
	})
}

var AlicloudDataWorksFolderMap0 = map[string]string{
	"folder_id":          CHECKSET,
	"folder_path":        "",
	"project_identifier": NOSET,
	"project_id":         "34051",
}

func AlicloudDataWorksFolderBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
`, name)
}
