// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test PaiWorkspace Member. >>> Resource test cases, automatically generated.
// Case 推荐转换-资源-54: Member_副本1732083979031 9039
func TestAccAliCloudPaiWorkspaceMember_basic9039(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_pai_workspace_member.default"
	ra := resourceAttrInit(resourceId, AliCloudPaiWorkspaceMemberMap9039)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &PaiWorkspaceServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribePaiWorkspaceMember")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tfacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudPaiWorkspaceMemberBasicDependence9039)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"user_id":      "${alicloud_ram_user.default.id}",
					"workspace_id": "${alicloud_pai_workspace_workspace.Workspace.id}",
					"roles":        []string{"PAI.AlgoDeveloper", "PAI.AlgoOperator", "PAI.LabelManager"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"user_id":      CHECKSET,
						"workspace_id": CHECKSET,
						"roles.#":      "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"roles": []string{"PAI.AlgoOperator"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"roles.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"roles": []string{"PAI.AlgoDeveloper", "PAI.AlgoOperator", "PAI.LabelManager"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"roles.#": "3",
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

func TestAccAliCloudPaiWorkspaceMember_basic9039_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_pai_workspace_member.default"
	ra := resourceAttrInit(resourceId, AliCloudPaiWorkspaceMemberMap9039)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &PaiWorkspaceServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribePaiWorkspaceMember")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tfacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudPaiWorkspaceMemberBasicDependence9039)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"user_id":      "${alicloud_ram_user.default.id}",
					"workspace_id": "${alicloud_pai_workspace_workspace.Workspace.id}",
					"roles":        []string{"PAI.AlgoDeveloper", "PAI.AlgoOperator", "PAI.LabelManager"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"user_id":      CHECKSET,
						"workspace_id": CHECKSET,
						"roles.#":      "3",
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

var AliCloudPaiWorkspaceMemberMap9039 = map[string]string{
	"member_id":   CHECKSET,
	"create_time": CHECKSET,
}

func AliCloudPaiWorkspaceMemberBasicDependence9039(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_pai_workspace_workspace" "Workspace" {
  description    = "800"
  workspace_name = var.name
  env_types      = ["prod"]
  display_name   = var.name
}

resource "alicloud_ram_user" "default" {
  name = var.name
}
`, name)
}

// Test PaiWorkspace Member. <<< Resource test cases, automatically generated.
