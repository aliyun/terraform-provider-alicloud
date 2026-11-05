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
//
// Beyond the basic create/update/import flow it exercises:
//   - R002: clearing the optional description (config omits it, the API resets it and the
//     empty-string-to-null plan modifier keeps the plan converged).
//   - R003: updating owner to a RAM user that is added to the workspace as a project member,
//     verifying the updatable (no ForceNew) owner path.
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
					resource.TestCheckResourceAttrPair("alicloud_data_works_parameter.default", "owner", "data.alicloud_account.default", "id"),
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
					// R003: owner changed from the master account to a RAM user workspace member.
					resource.TestCheckResourceAttrPair("alicloud_data_works_parameter.default", "owner", "alicloud_ram_user.member", "id"),
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
				// R002: description is omitted so the plan modifier maps it to null; UpdateParameter
				// is sent with an empty Description and the server-side value is reset. The owner and
				// properties from the previous step are unchanged.
				Config: testAccAliCloudDataWorksParameterConfigClearDescription(r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr("alicloud_data_works_parameter.default", "description"),
					resource.TestCheckResourceAttrPair("alicloud_data_works_parameter.default", "owner", "alicloud_ram_user.member", "id"),
					resource.TestCheckResourceAttr("alicloud_data_works_parameter.default", "properties.#", "2"),
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

// TestAccAliCloudDataWorksParameterSecretConstant covers the SecretConstant type (R004), which
// stores the value as ciphertext and requires a Prod properties entry. It verifies create/read
// and import for the secret-constant path that the PlainConstant test does not touch.
func TestAccAliCloudDataWorksParameterSecretConstant(t *testing.T) {
	r := tfacctest.RandIntRange(10000, 99999)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-chengdu"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudDataWorksParameterConfigSecretConstant(r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("alicloud_data_works_parameter.default", "id"),
					resource.TestCheckResourceAttr("alicloud_data_works_parameter.default", "type", "SecretConstant"),
					resource.TestCheckResourceAttr("alicloud_data_works_parameter.default", "scope", "Project"),
					resource.TestCheckResourceAttr("alicloud_data_works_parameter.default", "description", "tf_secret"),
					resource.TestCheckResourceAttrPair("alicloud_data_works_parameter.default", "owner", "data.alicloud_account.default", "id"),
					resource.TestCheckResourceAttr("alicloud_data_works_parameter.default", "properties.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs("alicloud_data_works_parameter.default", "properties.*", map[string]string{
						"env_type": "Prod",
						"value":    "secret_value",
					}),
					resource.TestCheckResourceAttrSet("alicloud_data_works_parameter.default", "version"),
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

// TestAccAliCloudDataWorksParameterVariable covers the Variable type (R004), exercising the
// variable-parameter create/read and import path with Prod and Dev value configurations.
func TestAccAliCloudDataWorksParameterVariable(t *testing.T) {
	r := tfacctest.RandIntRange(10000, 99999)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-chengdu"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudDataWorksParameterConfigVariable(r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("alicloud_data_works_parameter.default", "id"),
					resource.TestCheckResourceAttr("alicloud_data_works_parameter.default", "type", "Variable"),
					resource.TestCheckResourceAttr("alicloud_data_works_parameter.default", "scope", "Project"),
					resource.TestCheckResourceAttr("alicloud_data_works_parameter.default", "description", "tf_variable"),
					resource.TestCheckResourceAttrPair("alicloud_data_works_parameter.default", "owner", "data.alicloud_account.default", "id"),
					resource.TestCheckResourceAttr("alicloud_data_works_parameter.default", "properties.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs("alicloud_data_works_parameter.default", "properties.*", map[string]string{
						"env_type": "Prod",
						"value":    "var_prod",
					}),
					resource.TestCheckTypeSetElemNestedAttrs("alicloud_data_works_parameter.default", "properties.*", map[string]string{
						"env_type": "Dev",
						"value":    "var_dev",
					}),
					resource.TestCheckResourceAttrSet("alicloud_data_works_parameter.default", "version"),
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

// testAccAliCloudDataWorksParameterMember creates a RAM user and registers it as a workspace
// member so that it is a valid owner (责任人) for the parameter. It is only included in steps
// that change owner to the RAM user, and the parameter carries an explicit depends_on so that
// the member is created before the owner change is applied.
func testAccAliCloudDataWorksParameterMember(r int) string {
	return fmt.Sprintf(`
resource "alicloud_ram_user" "member" {
  name         = "${var.name}-member-${random_integer.default.result}"
  display_name = "${var.name}-member"
}

resource "alicloud_data_works_project_member" "member" {
  project_id = alicloud_data_works_project.default.id
  user_id    = alicloud_ram_user.member.id

  roles {
    code = "role_project_deploy"
  }
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
	return testAccAliCloudDataWorksParameterBase(r) + testAccAliCloudDataWorksParameterMember(r) + fmt.Sprintf(`
resource "alicloud_data_works_parameter" "default" {
  project_id  = alicloud_data_works_project.default.id
  name        = "workspace.tfacc%[1]d"
  type        = "PlainConstant"
  owner       = alicloud_ram_user.member.id
  description = "tf_updated"
  depends_on  = [alicloud_data_works_project_member.member]

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

// testAccAliCloudDataWorksParameterConfigClearDescription drops the description attribute so the
// empty-string-to-null plan modifier maps the configured null onto the API's empty value.
func testAccAliCloudDataWorksParameterConfigClearDescription(r int) string {
	return testAccAliCloudDataWorksParameterBase(r) + testAccAliCloudDataWorksParameterMember(r) + fmt.Sprintf(`
resource "alicloud_data_works_parameter" "default" {
  project_id = alicloud_data_works_project.default.id
  name       = "workspace.tfacc%[1]d"
  type       = "PlainConstant"
  owner      = alicloud_ram_user.member.id
  depends_on = [alicloud_data_works_project_member.member]

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

func testAccAliCloudDataWorksParameterConfigSecretConstant(r int) string {
	return testAccAliCloudDataWorksParameterBase(r) + fmt.Sprintf(`
resource "alicloud_data_works_parameter" "default" {
  project_id  = alicloud_data_works_project.default.id
  name        = "workspace.tfaccsec%[1]d"
  type        = "SecretConstant"
  owner       = data.alicloud_account.default.id
  description = "tf_secret"

  properties {
    env_type = "Prod"
    value    = "secret_value"
  }
}
`, r)
}

func testAccAliCloudDataWorksParameterConfigVariable(r int) string {
	return testAccAliCloudDataWorksParameterBase(r) + fmt.Sprintf(`
resource "alicloud_data_works_parameter" "default" {
  project_id  = alicloud_data_works_project.default.id
  name        = "workspace.tfaccvar%[1]d"
  type        = "Variable"
  owner       = data.alicloud_account.default.id
  description = "tf_variable"

  properties {
    env_type = "Prod"
    value    = "var_prod"
  }

  properties {
    env_type = "Dev"
    value    = "var_dev"
  }
}
`, r)
}
