// Package modelstudio. This file is generated automatically. Please do not modify it manually, thank you!
package modelstudio_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

// Test ModelStudio Workspace. >>> Resource test cases, automatically generated.
// Case Resource_Workspace_Init 12769
func TestAccAliCloudModelStudioWorkspace_basic12769(t *testing.T) {
	r := tfacctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tfacc%d", r)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudModelStudioWorkspace_basic12769_step0(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_model_studio_workspace.default", "service_site", "global"),

					resource.TestCheckResourceAttr("alicloud_model_studio_workspace.default", "workspace_name", name),
					resource.TestCheckResourceAttrSet("alicloud_model_studio_workspace.default", "create_time"),
				),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("alicloud_model_studio_workspace.default", tfjsonpath.New("api_host"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_model_studio_workspace.default", tfjsonpath.New("region_id"), knownvalue.NotNull()),
				},
			},
			{
				Config: testAccAliCloudModelStudioWorkspace_basic12769_step1(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_model_studio_workspace.default", "service_site", "global"),

					resource.TestCheckResourceAttr("alicloud_model_studio_workspace.default", "workspace_name", name+"upd"),
				),
			},
			{
				ResourceName:            "alicloud_model_studio_workspace.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

func testAccAliCloudModelStudioWorkspace_basic12769_step0(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%[1]s"
}

resource "alicloud_model_studio_workspace" "default" {
        service_site = "global"

        workspace_name = "%[1]s"
}
`, name)
}

func testAccAliCloudModelStudioWorkspace_basic12769_step1(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%[1]s"
}

resource "alicloud_model_studio_workspace" "default" {
        service_site = "global"

        workspace_name = "%[1]supd"
}
`, name)
}

// Test ModelStudio Workspace. <<< Resource test cases, automatically generated.
