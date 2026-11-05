package paiworkspace_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/paiworkspace"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// Test PaiWorkspace Workspace. >>> Resource test cases, automatically generated.
// Case CC发布PAIWorkspace_副本1726115966862 7875
func TestAccAliCloudPaiWorkspaceWorkspace_basic7875(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_pai_workspace_workspace.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudPaiWorkspaceWorkspaceMap7875)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return paiworkspace.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribePaiWorkspaceWorkspace")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testaccworkspace%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudPaiWorkspaceWorkspaceBasicDependence7875)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"description":    "452",
					"workspace_name": name,
					"env_types": []string{
						"prod"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":    acctest.CHECKSET,
						"workspace_name": name,
						"env_types.#":    "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"display_name": "new_test_pop_559",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"display_name": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "986",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"display_name": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": acctest.CHECKSET,
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

func TestAccAliCloudPaiWorkspaceWorkspace_basic7875_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_pai_workspace_workspace.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudPaiWorkspaceWorkspaceMap7875)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return paiworkspace.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribePaiWorkspaceWorkspace")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testaccworkspace%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudPaiWorkspaceWorkspaceBasicDependence7875)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"description":       "452",
					"workspace_name":    name,
					"display_name":      "new_test_pop_559",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"env_types": []string{
						"prod"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":       acctest.CHECKSET,
						"workspace_name":    name,
						"display_name":      acctest.CHECKSET,
						"resource_group_id": acctest.CHECKSET,
						"env_types.#":       "1",
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

var AliCloudPaiWorkspaceWorkspaceMap7875 = map[string]string{
	"status":            acctest.CHECKSET,
	"create_time":       acctest.CHECKSET,
	"resource_group_id": acctest.CHECKSET,
	"display_name":      acctest.CHECKSET,
}

func AliCloudPaiWorkspaceWorkspaceBasicDependence7875(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}
data "alicloud_resource_manager_resource_groups" "default" {
}

`, name)
}

func TestAccAliCloudPaiWorkspaceWorkspace_duplicateName(t *testing.T) {
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testaccworkspace%d", rand)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckPaiWorkspaceWorkspaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPaiWorkspaceWorkspaceNameConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_pai_workspace_workspace.default", "workspace_name", name),
				),
			},
			{
				Config:      testAccPaiWorkspaceWorkspaceDuplicateNameConfig(name),
				ExpectError: regexp.MustCompile("already taken by the workspace.*terraform import alicloud_pai_workspace_workspace"),
			},
		},
	})
}

func TestAccAliCloudPaiWorkspaceWorkspace_nameRetainedAfterDestroy(t *testing.T) {
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testaccworkspace%d", rand)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckPaiWorkspaceWorkspaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPaiWorkspaceWorkspaceNameConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_pai_workspace_workspace.default", "workspace_name", name),
				),
			},
			{
				Config:  testAccPaiWorkspaceWorkspaceNameConfig(name),
				Destroy: true,
			},
			{
				Config:      testAccPaiWorkspaceWorkspaceNameConfig(name),
				ExpectError: regexp.MustCompile("name is still reserved.*recycle bin"),
			},
		},
	})
}

func testAccPaiWorkspaceWorkspaceNameConfig(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alicloud_pai_workspace_workspace" "default" {
  workspace_name = var.name
  description    = "terraform test duplicate workspace name"
  env_types      = ["prod"]
}
`, name)
}

func testAccPaiWorkspaceWorkspaceDuplicateNameConfig(name string) string {
	return testAccPaiWorkspaceWorkspaceNameConfig(name) + `
resource "alicloud_pai_workspace_workspace" "duplicate" {
  workspace_name = var.name
  description    = "terraform test duplicate workspace name"
  env_types      = ["prod"]
}
`
}

func testAccCheckPaiWorkspaceWorkspaceDestroy(s *terraform.State) error {
	client := acctest.Provider.Meta().(*connectivity.AliyunClient)
	paiWorkspaceServiceV2 := paiworkspace.NewService(client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_pai_workspace_workspace" {
			continue
		}
		if _, err := paiWorkspaceServiceV2.DescribePaiWorkspaceWorkspace(rs.Primary.ID); err != nil {
			if errs.NotFoundError(err) {
				continue
			}
			return err
		}
		return fmt.Errorf("PAI workspace %s still exists", rs.Primary.ID)
	}
	return nil
}

// Test PaiWorkspace Workspace. <<< Resource test cases, automatically generated.
