package paiworkspace_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// PAI Workspace permissions are read-only catalog entries. The acceptance test
// brings up a workspace (env_types=["prod"], which surfaces the PaiDLC
// permission catalog), then applies the permission resource against the
// documented `PaiDLC:ListJobs` code (the GetPermission response example) and
// verifies the computed permission_rules serialize with real content. The
// import step resolves the permission by enumerating ListPermissions for the
// workspace and matching the code, so importing one entry among the catalog
// resolves the correct identity — the collection-Read distinguishing power
// (R4.23/R7.16) — and the post-import verify checks it matches the apply path.
//
// The test uses the terraform-plugin-sdk/v2 helper/resource harness (rather
// than terraform-plugin-testing) because the alicloud root package, pulled in
// transitively via acctest.Provider, contains non-test files that import
// terraform-plugin-sdk/v2/helper/resource (e.g. resource_alicloud_polardb_polar_fs).
// Both helper/resource packages register a -sweep flag in init, so a test
// binary that links both dies at startup with "flag redefined: sweep". Using
// the SDK v2 harness keeps a single -sweep registration while still driving the
// muxed ProtoV5 provider (Framework half + SDK v2 half) via
// ProtoV5ProviderFactories.
//
// The permission_rules catalog is flattened into two parallel computed list
// attributes (permission_rule_accessibilities, permission_rule_entity_access_types)
// rather than a nested attribute: see resource_permission.go for the v5 mux
// constraint that rules out both nested attributes and blocks.
func TestAccAliCloudPaiWorkspacePermission_basic(t *testing.T) {
	preCheckRegion(t, "cn-hangzhou")

	// PAI Workspace workspace_name rules (per the alicloud_pai_workspace_workspace
	// docs): 3-23 chars, letters/underscores/numbers, must start with a letter.
	// Hyphens and >23-char names are rejected by the API (ErrorCode 1101080002 /
	// 100400010 "Workspace name error"). Keep the prefix short so prefix+5-digit
	// suffix stays <=23 (e.g. "tf_testpai_perm"+5 = 20), and use a random suffix
	// so repeated runs do not collide on a reserved/recycled name.
	rand := sdkacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testpai_perm%d", rand)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudPaiWorkspacePermissionConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("alicloud_pai_workspace_permission.default", "id"),
					resource.TestCheckResourceAttrSet("alicloud_pai_workspace_permission.default", "workspace_id"),
					resource.TestCheckResourceAttr("alicloud_pai_workspace_permission.default", "permission_code", "PaiDLC:ListJobs"),
					// The flattened rule lists carry real content (R4.23/R7.16
					// distinguishing power), not just a count.
					resource.TestCheckResourceAttrSet("alicloud_pai_workspace_permission.default", "permission_rule_accessibilities.0"),
					resource.TestCheckResourceAttrSet("alicloud_pai_workspace_permission.default", "permission_rule_entity_access_types.0"),
				),
			},
			{
				ResourceName:            "alicloud_pai_workspace_permission.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

// preCheckRegion pins the test region (PAI Workspace acceptance runs in
// cn-hangzhou, matching the existing alicloud_pai_workspace_member suite) and
// skips when the harness targets a different region, so the suite does not flap
// on regions where PAI Workspace is unavailable.
func preCheckRegion(t *testing.T, region string) {
	t.Helper()
	if v := os.Getenv("ALICLOUD_REGION"); v != "" && v != region {
		t.Skipf("PAI Workspace acceptance tests require region %s; got %s", region, v)
	}
	if os.Getenv("ALICLOUD_REGION") == "" {
		_ = os.Setenv("ALICLOUD_REGION", region)
	}
}

func testAccAliCloudPaiWorkspacePermissionConfig(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alicloud_pai_workspace_workspace" "Workspace" {
  workspace_name = var.name
  description    = "tf-acc-pai-permission"
  env_types      = ["prod"]
}

resource "alicloud_pai_workspace_permission" "default" {
  workspace_id    = alicloud_pai_workspace_workspace.Workspace.id
  permission_code = "PaiDLC:ListJobs"
}
`, name)
}
