package pai_workspace_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccAliCloudPaiWorkspaceConnection exercises the connection lifecycle:
// create with the five CloudSpec-mapped governance fields, update the settable
// fields (verifying the create-only capability fields survive unchanged), and
// import. The workspace dependency is created in-config because a connection
// belongs to a PAI workspace.
func TestAccAliCloudPaiWorkspaceConnection(t *testing.T) {
	r := tfacctest.RandIntRange(10000, 99999)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccPaiWorkspaceConnectionConfig(r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("alicloud_pai_workspace_connection.default", "id"),
					resource.TestCheckResourceAttrSet("alicloud_pai_workspace_connection.default", "connection_id"),
					resource.TestCheckResourceAttr("alicloud_pai_workspace_connection.default", "connection_name", fmt.Sprintf("tf_testaccconn%d", r)),
					resource.TestCheckResourceAttr("alicloud_pai_workspace_connection.default", "connection_type", "MilvusConnection"),
					resource.TestCheckResourceAttrSet("alicloud_pai_workspace_connection.default", "workspace_id"),
					resource.TestCheckResourceAttr("alicloud_pai_workspace_connection.default", "accessibility", "PRIVATE"),
					resource.TestCheckResourceAttr("alicloud_pai_workspace_connection.default", "description", fmt.Sprintf("tf-testacc-desc-%d", r)),
					resource.TestCheckResourceAttr("alicloud_pai_workspace_connection.default", "configs.database", "default"),
					resource.TestCheckResourceAttr("alicloud_pai_workspace_connection.default", "configs.uri", fmt.Sprintf("tf-testacc-milvus-%d.example.com", r)),
					resource.TestCheckResourceAttr("alicloud_pai_workspace_connection.default", "secrets.token", fmt.Sprintf("tf-testacc-token-%d", r)),
					resource.TestCheckResourceAttr("alicloud_pai_workspace_connection.default", "models.0.model", fmt.Sprintf("tf-testacc-model-%d", r)),
					resource.TestCheckResourceAttr("alicloud_pai_workspace_connection.default", "models.0.display_name", fmt.Sprintf("tf-testacc-model-display-%d", r)),
					resource.TestCheckResourceAttr("alicloud_pai_workspace_connection.default", "models.0.model_type", "OPEN_SOURCE"),
					resource.TestCheckResourceAttr("alicloud_pai_workspace_connection.default", "models.0.tool_call", "false"),
					// Five governance fields: written at create, read back.
					resource.TestCheckResourceAttr("alicloud_pai_workspace_connection.default", "models.0.support_reasoning", "true"),
					resource.TestCheckResourceAttr("alicloud_pai_workspace_connection.default", "models.0.support_response_schema", "true"),
					resource.TestCheckResourceAttr("alicloud_pai_workspace_connection.default", "models.0.support_vision", "false"),
					resource.TestCheckResourceAttr("alicloud_pai_workspace_connection.default", "models.0.max_model_length", "4096"),
					resource.TestCheckResourceAttr("alicloud_pai_workspace_connection.default", "resource_meta.extra", fmt.Sprintf("tf-testacc-extra-%d", r)),
				),
			},
			{
				Config: testAccPaiWorkspaceConnectionConfigUpdate(r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_pai_workspace_connection.default", "description", fmt.Sprintf("tf-testacc-desc-upd-%d", r)),
					resource.TestCheckResourceAttr("alicloud_pai_workspace_connection.default", "configs.uri", fmt.Sprintf("tf-testacc-milvus-upd-%d.example.com", r)),
					resource.TestCheckResourceAttr("alicloud_pai_workspace_connection.default", "secrets.token", fmt.Sprintf("tf-testacc-token-upd-%d", r)),
					resource.TestCheckResourceAttr("alicloud_pai_workspace_connection.default", "models.0.model", fmt.Sprintf("tf-testacc-model-upd-%d", r)),
					resource.TestCheckResourceAttr("alicloud_pai_workspace_connection.default", "models.0.display_name", fmt.Sprintf("tf-testacc-model-display-upd-%d", r)),
					resource.TestCheckResourceAttr("alicloud_pai_workspace_connection.default", "models.0.tool_call", "true"),
					// Create-only governance fields survive the update unchanged.
					resource.TestCheckResourceAttr("alicloud_pai_workspace_connection.default", "models.0.support_reasoning", "true"),
					resource.TestCheckResourceAttr("alicloud_pai_workspace_connection.default", "models.0.max_model_length", "4096"),
					resource.TestCheckResourceAttr("alicloud_pai_workspace_connection.default", "resource_meta.extra", fmt.Sprintf("tf-testacc-extra-%d", r)),
				),
			},
			{
				// secrets are masked by the API on read-back; max_model_length and the
				// support_reasoning/support_response_schema create-only capability fields
				// are not echoed by GetConnection, so none can be reconstructed from a
				// fresh import (no prior state to preserve). gmt_modified_time is bumped
				// by the API on every update and may advance between the last apply and
				// the import read, so it is also excluded from import verification.
				ResourceName:            "alicloud_pai_workspace_connection.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"secrets", "models.0.max_model_length", "gmt_modified_time", "models.0.support_reasoning", "models.0.support_response_schema"},
			},
		},
	})
}

func testAccPaiWorkspaceConnectionConfig(r int) string {
	return fmt.Sprintf(`
resource "alicloud_pai_workspace_workspace" "default" {
  workspace_name = "tf_testaccconnws%[1]d"
  description    = "tf-testacc-conn-ws-%[1]d"
  env_types      = ["prod"]
}

resource "alicloud_pai_workspace_connection" "default" {
  connection_name = "tf_testaccconn%[1]d"
  connection_type = "MilvusConnection"
  workspace_id    = alicloud_pai_workspace_workspace.default.id
  accessibility   = "PRIVATE"
  description     = "tf-testacc-desc-%[1]d"

  configs = {
    database = "default"
    uri      = "tf-testacc-milvus-%[1]d.example.com"
  }

  secrets = {
    token = "tf-testacc-token-%[1]d"
  }

  models = [
    {
      model                   = "tf-testacc-model-%[1]d"
      display_name            = "tf-testacc-model-display-%[1]d"
      model_type              = "OPEN_SOURCE"
      tool_call               = false
      support_reasoning       = true
      support_response_schema = true
      support_vision          = false
      max_model_length        = 4096
    }
  ]

  resource_meta = {
    instance_name = "tf-testacc-inst-%[1]d"
    instance_id   = "tf-testacc-instid-%[1]d"
    extra         = "tf-testacc-extra-%[1]d"
  }
}
`, r)
}

func testAccPaiWorkspaceConnectionConfigUpdate(r int) string {
	return fmt.Sprintf(`
resource "alicloud_pai_workspace_workspace" "default" {
  workspace_name = "tf_testaccconnws%[1]d"
  description    = "tf-testacc-conn-ws-%[1]d"
  env_types      = ["prod"]
}

resource "alicloud_pai_workspace_connection" "default" {
  connection_name = "tf_testaccconn%[1]d"
  connection_type = "MilvusConnection"
  workspace_id    = alicloud_pai_workspace_workspace.default.id
  accessibility   = "PRIVATE"
  description     = "tf-testacc-desc-upd-%[1]d"

  configs = {
    database = "default"
    uri      = "tf-testacc-milvus-upd-%[1]d.example.com"
  }

  secrets = {
    token = "tf-testacc-token-upd-%[1]d"
  }

  models = [
    {
      model                   = "tf-testacc-model-upd-%[1]d"
      display_name            = "tf-testacc-model-display-upd-%[1]d"
      model_type              = "OPEN_SOURCE"
      tool_call               = true
      support_reasoning       = true
      support_response_schema = true
      support_vision          = false
      max_model_length        = 4096
    }
  ]

  resource_meta = {
    instance_name = "tf-testacc-inst-%[1]d"
    instance_id   = "tf-testacc-instid-%[1]d"
    extra         = "tf-testacc-extra-%[1]d"
  }
}
`, r)
}
