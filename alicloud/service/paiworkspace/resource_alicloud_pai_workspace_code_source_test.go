package paiworkspace_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/paiworkspace"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test PaiWorkspace CodeSource. >>> Resource test cases, automatically generated.
// Case TestCodeSource 9010
func TestAccAliCloudPaiWorkspaceCodeSource_basic9010(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_pai_workspace_code_source.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudPaiWorkspaceCodeSourceMap9010)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return paiworkspace.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribePaiWorkspaceCodeSource")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tf_testacc%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudPaiWorkspaceCodeSourceBasicDependence9010)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-shenzhen"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"mount_path":             "/mnt/code/dir_01/",
					"code_repo":              "https://github.com/mattn/go-sqlite3.git",
					"description":            "desc-01",
					"code_repo_access_token": "token-01",
					"accessibility":          "PRIVATE",
					"display_name":           "codesource-test-01",
					"workspace_id":           "${alicloud_pai_workspace_workspace.defaultgklBnM.id}",
					"code_branch":            "master",
					"code_repo_user_name":    "user-01",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mount_path":             "/mnt/code/dir_01/",
						"code_repo":              "https://github.com/mattn/go-sqlite3.git",
						"description":            "desc-01",
						"code_repo_access_token": "token-01",
						"accessibility":          "PRIVATE",
						"display_name":           "codesource-test-01",
						"workspace_id":           acctest.CHECKSET,
						"code_branch":            "master",
						"code_repo_user_name":    "user-01",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mount_path":             "/mnt/code/dir_02/",
					"code_repo":              "https://github.com/mattn-02/go-sqlite3.git",
					"description":            "desc-02",
					"code_repo_access_token": "token-02",
					"display_name":           "codesource-test-02",
					"code_branch":            "master-02",
					"code_repo_user_name":    "user-02",
					"code_commit":            "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mount_path":             "/mnt/code/dir_02/",
						"code_repo":              "https://github.com/mattn-02/go-sqlite3.git",
						"description":            "desc-02",
						"code_repo_access_token": "token-02",
						"display_name":           "codesource-test-02",
						"code_branch":            "master-02",
						"code_repo_user_name":    "user-02",
						"code_commit":            "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"accessibility": "PUBLIC",
					"code_commit":   "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"accessibility": "PUBLIC",
						"code_commit":   "1",
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

var AlicloudPaiWorkspaceCodeSourceMap9010 = map[string]string{
	"create_time": acctest.CHECKSET,
}

func AlicloudPaiWorkspaceCodeSourceBasicDependence9010(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_pai_workspace_workspace" "defaultgklBnM" {
  description    = "for-pop-test"
  display_name   = "CodeSourceTest_1732796226"
  workspace_name = var.name
  env_types      = ["prod"]
}


`, name)
}

// Test PaiWorkspace CodeSource. <<< Resource test cases, automatically generated.
