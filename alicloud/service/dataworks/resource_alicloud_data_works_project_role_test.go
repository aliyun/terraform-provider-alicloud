package dataworks_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/dataworks"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// Test DataWorks ProjectRole. >>> Resource test cases.
const dataWorksProjectRoleResourceId = "alicloud_data_works_project_role.default"

func TestAccAliCloudDataWorksProjectRole_basic(t *testing.T) {
	rand := tfacctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tf_testacc%d", rand)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-chengdu"})
			acctest.PreCheck(t)
		},
		IDRefreshName:            dataWorksProjectRoleResourceId,
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckDataWorksProjectRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataWorksProjectRoleConfig(name, "Read", "Write"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataWorksProjectRoleResourceId, "id"),
					resource.TestCheckResourceAttr(dataWorksProjectRoleResourceId, "role_name", name),
					resource.TestCheckResourceAttrSet(dataWorksProjectRoleResourceId, "code"),
					resource.TestCheckResourceAttrSet(dataWorksProjectRoleResourceId, "type"),
					resource.TestCheckResourceAttr(dataWorksProjectRoleResourceId, "module_permissions.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(dataWorksProjectRoleResourceId, "module_permissions.*", map[string]string{
						"module_id":       "9",
						"permission_type": "Read",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(dataWorksProjectRoleResourceId, "module_permissions.*", map[string]string{
						"module_id":       "11",
						"permission_type": "Write",
					}),
				),
			},
			{
				// Update only ModulePermissions: drop module 11 and switch module 9 to Write.
				Config: testAccDataWorksProjectRoleConfig(name, "Write", ""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataWorksProjectRoleResourceId, "module_permissions.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(dataWorksProjectRoleResourceId, "module_permissions.*", map[string]string{
						"module_id":       "9",
						"permission_type": "Write",
					}),
				),
			},
			{
				// Second plan with no change must produce no diff.
				Config:             testAccDataWorksProjectRoleConfig(name, "Write", ""),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
			{
				ResourceName:      dataWorksProjectRoleResourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// testAccCheckDataWorksProjectRoleDestroy confirms every ProjectRole left in
// state is gone. A role already removed (1101030502 or a typed 404) surfaces as
// a provider NotFound; any other failure must be reported instead of masked.
func testAccCheckDataWorksProjectRoleDestroy(s *terraform.State) error {
	client := acctest.Provider.Meta().(*connectivity.AliyunClient)
	service := dataworks.NewService(client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_data_works_project_role" {
			continue
		}
		_, err := service.DescribeDataWorksProjectRole(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("project role %s still exists", rs.Primary.ID)
		}
		if errs.IsProviderNotFound(err) {
			continue
		}
		return fmt.Errorf("describing project role %s failed: %s", rs.Primary.ID, err)
	}
	return nil
}

// testAccDataWorksProjectRoleConfig builds a workspace then a project role with
// two module permissions (module 9 = 数据集成, module 11 = 数据开发). When the
// second permission type is empty, only the module 9 permission is configured,
// which exercises the UpdateProjectRole whole-set replacement path.
func testAccDataWorksProjectRoleConfig(name, firstPermission, secondPermission string) string {
	secondBlock := ""
	if secondPermission != "" {
		secondBlock = fmt.Sprintf(`
  module_permissions {
    module_id       = "11"
    permission_type = "%s"
  }
`, secondPermission)
	}
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_data_works_project" "default" {
  status             = "Available"
  description        = "tf_desc"
  project_name       = var.name
  pai_task_enabled   = "false"
  display_name       = "tf_project_role_display"
  dev_role_disabled  = "true"
  resource_group_id  = data.alicloud_resource_manager_resource_groups.default.ids.0
}

resource "alicloud_data_works_project_role" "default" {
  project_id = alicloud_data_works_project.default.id
  role_name  = var.name

  module_permissions {
    module_id       = "9"
    permission_type = "%s"
  }
%s}
`, name, firstPermission, secondBlock)
}

// Test DataWorks ProjectRole. <<< Resource test cases.
