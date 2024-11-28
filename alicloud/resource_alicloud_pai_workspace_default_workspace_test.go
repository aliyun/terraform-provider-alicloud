package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test PaiWorkspace DefaultWorkspace. >>> Resource test cases, automatically generated.
// Case DefaultWorkspace 用例测试01 6819
func TestAccAliCloudPaiWorkspaceDefaultWorkspace_basic6819(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_pai_workspace_default_workspace.default"
	ra := resourceAttrInit(resourceId, AlicloudPaiWorkspaceDefaultWorkspaceMap6819)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &PaiWorkspaceServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribePaiWorkspaceDefaultWorkspace")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tf_testacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudPaiWorkspaceDefaultWorkspaceBasicDependence6819)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "defaultWorkspace",
					"env_types": []string{
						"prod"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "defaultWorkspace",
						"env_types.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"workspace_id": "${alicloud_pai_workspace_workspace.defaultWorkspace.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"workspace_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"workspace_id": "${alicloud_pai_workspace_workspace.changeWorkspace.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"workspace_id": CHECKSET,
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

var AlicloudPaiWorkspaceDefaultWorkspaceMap6819 = map[string]string{
	"status": CHECKSET,
}

func AlicloudPaiWorkspaceDefaultWorkspaceBasicDependence6819(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_pai_workspace_workspace" "defaultWorkspace" {
  description    = "defaultWorkspace"
  display_name   = "DatasetResouceTest_16"
  workspace_name = "DatasetResouceTest_16"
  env_types      = ["prod"]
}

resource "alicloud_pai_workspace_workspace" "changeWorkspace" {
  description    = "defaultWorkspaced'f'c"
  display_name   = "DatasetResouceTest_146"
  workspace_name = "DatasetResouceTest_146"
  env_types      = ["prod"]
}

`, name)
}

// Test PaiWorkspace DefaultWorkspace. <<< Resource test cases, automatically generated.
