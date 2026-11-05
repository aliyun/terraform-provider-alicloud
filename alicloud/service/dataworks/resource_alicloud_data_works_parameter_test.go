package dataworks_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccAliCloudDataWorksParameter covers the Create/Read/Update/Delete and Import lifecycle of
// the alicloud_data_works_parameter Framework resource against dataworks-public/2024-05-18.
func TestAccAliCloudDataWorksParameter(t *testing.T) {
	r := tfacctest.RandIntRange(10000, 99999)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-chengdu"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudDataWorksParameterConfig(r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("alicloud_data_works_parameter.default", "id"),
					resource.TestCheckResourceAttr("alicloud_data_works_parameter.default", "name", fmt.Sprintf("workspace.tfacc%d", r)),
					resource.TestCheckResourceAttr("alicloud_data_works_parameter.default", "type", "PlainConstant"),
					resource.TestCheckResourceAttr("alicloud_data_works_parameter.default", "scope", "Project"),
					resource.TestCheckResourceAttr("alicloud_data_works_parameter.default", "description", "tf_basic"),
					resource.TestCheckResourceAttr("alicloud_data_works_parameter.default", "properties.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs("alicloud_data_works_parameter.default", "properties.*", map[string]string{
						"env_type": "Prod",
						"value":    "prod_value",
					}),
					resource.TestCheckTypeSetElemNestedAttrs("alicloud_data_works_parameter.default", "properties.*", map[string]string{
						"env_type": "Dev",
						"value":    "dev_value",
					}),
					resource.TestCheckResourceAttrSet("alicloud_data_works_parameter.default", "create_time"),
					resource.TestCheckResourceAttrSet("alicloud_data_works_parameter.default", "version"),
				),
			},
			{
				Config: testAccAliCloudDataWorksParameterConfigUpdate(r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_data_works_parameter.default", "description", "tf_updated"),
					resource.TestCheckResourceAttr("alicloud_data_works_parameter.default", "properties.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs("alicloud_data_works_parameter.default", "properties.*", map[string]string{
						"env_type": "Prod",
						"value":    "prod_value_new",
					}),
					resource.TestCheckTypeSetElemNestedAttrs("alicloud_data_works_parameter.default", "properties.*", map[string]string{
						"env_type": "Dev",
						"value":    "dev_value_new",
					}),
					resource.TestCheckResourceAttrSet("alicloud_data_works_parameter.default", "modify_time"),
				),
			},
			{
				ResourceName:      "alicloud_data_works_parameter.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccAliCloudDataWorksParameterBase(r int) string {
	return fmt.Sprintf(`
provider "alicloud" {
  region = "cn-chengdu"
}

variable "name" {
  default = "tf-acc-%[1]d"
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

data "alicloud_account" "default" {}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_data_works_project" "default" {
  project_name            = "${var.name}-${random_integer.default.result}"
  display_name            = var.name
  description             = "tf_for_parameter"
  pai_task_enabled        = false
  dev_environment_enabled = true
  dev_role_disabled       = true
  resource_group_id       = data.alicloud_resource_manager_resource_groups.default.ids.0
  status                  = "Available"
}
`, r)
}

func testAccAliCloudDataWorksParameterConfig(r int) string {
	return testAccAliCloudDataWorksParameterBase(r) + fmt.Sprintf(`
resource "alicloud_data_works_parameter" "default" {
  project_id  = alicloud_data_works_project.default.id
  name        = "workspace.tfacc%[1]d"
  type        = "PlainConstant"
  owner       = data.alicloud_account.default.id
  description = "tf_basic"

  properties {
    env_type = "Prod"
    value    = "prod_value"
  }

  properties {
    env_type = "Dev"
    value    = "dev_value"
  }
}
`, r)
}

func testAccAliCloudDataWorksParameterConfigUpdate(r int) string {
	return testAccAliCloudDataWorksParameterBase(r) + fmt.Sprintf(`
resource "alicloud_data_works_parameter" "default" {
  project_id  = alicloud_data_works_project.default.id
  name        = "workspace.tfacc%[1]d"
  type        = "PlainConstant"
  owner       = data.alicloud_account.default.id
  description = "tf_updated"

  properties {
    env_type = "Prod"
    value    = "prod_value_new"
  }

  properties {
    env_type = "Dev"
    value    = "dev_value_new"
  }
}
`, r)
}
