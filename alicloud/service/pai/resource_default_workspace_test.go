package pai_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// The default workspace is a singleton per region per account: it cannot be
// deleted and is one-per-account, so the acceptance tests use a STABLE
// description/env_types pair (not a random suffix) across every create step.
// This keeps the resource re-runnable: a second run adopts the already-enabled
// default workspace without drift, because the configuration matches the cloud.
//
// The test uses the terraform-plugin-sdk/v2 harness (not terraform-plugin-testing)
// because the repo's root package links terraform-plugin-sdk/v2/helper/resource
// (which registers the SDK "sweep" flag); importing terraform-plugin-testing's
// helper/resource in the same test binary would re-register "sweep" and panic.
// The SDK TestCase still accepts ProtoV5ProviderFactories, so the muxed
// Framework provider is exercised.

func TestAccAliCloudPaiWorkspaceDefaultWorkspace_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             func(*terraform.State) error { return nil },
		Steps: []resource.TestStep{
			{
				// description matches the immutable singleton's existing value
				// ("defaultWorkspace"): the default workspace is one per account,
				// cannot be deleted, and Description is ForceNew, so on an account
				// that already hosts a default workspace the only drift-free
				// description is the one the cloud already holds.
				Config: testAccPaiWorkspaceDefaultWorkspaceConfig("defaultWorkspace", []string{"prod"}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("alicloud_pai_workspace_default_workspace.default", "id"),
					resource.TestCheckResourceAttrSet("alicloud_pai_workspace_default_workspace.default", "workspace_id"),
					resource.TestCheckResourceAttr("alicloud_pai_workspace_default_workspace.default", "description", "defaultWorkspace"),
					resource.TestCheckResourceAttr("alicloud_pai_workspace_default_workspace.default", "env_types.#", "1"),
					resource.TestCheckResourceAttr("alicloud_pai_workspace_default_workspace.default", "env_types.0", "prod"),
					resource.TestCheckResourceAttr("alicloud_pai_workspace_default_workspace.default", "status", "ENABLED"),
					resource.TestCheckResourceAttrSet("alicloud_pai_workspace_default_workspace.default", "workspace_name"),
					resource.TestCheckResourceAttrSet("alicloud_pai_workspace_default_workspace.default", "region_id"),
				),
			},
			{
				ResourceName:            "alicloud_pai_workspace_default_workspace.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"gmt_modified_time"},
			},
		},
	})
}

// TestAccAliCloudPaiWorkspaceDefaultWorkspace_forceNew asserts that changing
// the RequiresReplace description produces a replacement plan. The step is
// PlanOnly: the singleton default workspace cannot be deleted and re-created,
// so the replacement is planned but not applied, which is the honest signal a
// user would receive when trying to change an immutable field.
func TestAccAliCloudPaiWorkspaceDefaultWorkspace_forceNew(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             func(*terraform.State) error { return nil },
		Steps: []resource.TestStep{
			{
				Config: testAccPaiWorkspaceDefaultWorkspaceConfig("defaultWorkspace", []string{"prod"}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_pai_workspace_default_workspace.default", "status", "ENABLED"),
				),
			},
			{
				Config:             testAccPaiWorkspaceDefaultWorkspaceConfig("terraform-pai-default-workspace-changed", []string{"prod"}),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

// TestAccAliCloudPaiWorkspaceDefaultWorkspace_invalidDescription exercises the
// failure path deterministically: a description longer than 80 characters is
// rejected by provider-side schema validation, so the failure does not depend on
// the account already hosting a default workspace.
func TestAccAliCloudPaiWorkspaceDefaultWorkspace_invalidDescription(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             func(*terraform.State) error { return nil },
		Steps: []resource.TestStep{
			{
				Config:      testAccPaiWorkspaceDefaultWorkspaceConfig(longDescription(81), []string{"prod"}),
				ExpectError: regexp.MustCompile(`description`),
			},
		},
	})
}

func testAccPaiWorkspaceDefaultWorkspaceConfig(description string, envTypes []string) string {
	envs := ""
	for i, e := range envTypes {
		if i > 0 {
			envs += ", "
		}
		envs += fmt.Sprintf("%q", e)
	}
	return fmt.Sprintf(`
resource "alicloud_pai_workspace_default_workspace" "default" {
  description = %[1]q
  env_types   = [%[2]s]
}
`, description, envs)
}

// longDescription returns a string of n lowercase 'a' characters, used to exceed
// the 80-character description cap.
func longDescription(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = 'a'
	}
	return string(b)
}
