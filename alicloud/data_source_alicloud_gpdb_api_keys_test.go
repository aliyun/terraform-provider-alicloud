package alicloud

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// TestAccAlicloudGPDBApiKeysDataSource covers the alicloud_gpdb_api_keys data
// source. A GPDB ApiKey requires a pre-existing GPDB workspace which has no
// Terraform resource, so the workspace is provisioned out of band (reusing the
// helpers shared with the resource test) and the returned id is interpolated
// into the config.
//
// Each scenario is split into two TestSteps:
//
//	Step 1 — create the resource only (populates state).
//	Step 2 — add the data source config (resource already in state, so the
//	         data source can be read during plan without depends_on).
//
// This avoids the Terraform SDK v1 quirk where depends_on on a data source
// defers its Read to the apply phase, leaving all attributes as <computed>
// after refresh and causing plan-non-empty failures.
func TestAccAlicloudGPDBApiKeysDataSource(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)

	var workspaceId string
	var wsClient *connectivity.AliyunClient
	if os.Getenv("TF_ACC") != "" {
		wsClient, workspaceId = gpdbApiKeyTestCreateWorkspace(t, rand)
		defer gpdbApiKeyTestDeleteWorkspace(wsClient, workspaceId)
	}

	preCheck := func() {
		testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
		testAccPreCheck(t)
	}

	name := fmt.Sprintf("tfaccgpdbapikey%d", rand)

	// Resource-only config — used as Step 1 to populate state.
	resourceConf := fmt.Sprintf(`
resource "alicloud_gpdb_api_key" "default" {
  workspace_id = "%s"
  key_name     = "%s"
  description  = "terraform test"
}
`, workspaceId, name)

	// Data-source-only configs — used as Step 2. The resource already
	// exists in state from Step 1, so the data source can reference it
	// via interpolation (implicit dependency) without depends_on.
	dsBasicConf := fmt.Sprintf(`
resource "alicloud_gpdb_api_key" "default" {
  workspace_id = "%s"
  key_name     = "%s"
  description  = "terraform test"
}

data "alicloud_gpdb_api_keys" "default" {
  workspace_id = "%s"
}
`, workspaceId, name, workspaceId)

	dsIdsExistConf := fmt.Sprintf(`
resource "alicloud_gpdb_api_key" "default" {
  workspace_id = "%s"
  key_name     = "%s"
  description  = "terraform test"
}

data "alicloud_gpdb_api_keys" "default" {
  workspace_id = "%s"
  ids          = ["${alicloud_gpdb_api_key.default.key_id}"]
}
`, workspaceId, name, workspaceId)

	dsIdsFakeConf := fmt.Sprintf(`
resource "alicloud_gpdb_api_key" "default" {
  workspace_id = "%s"
  key_name     = "%s"
  description  = "terraform test"
}

data "alicloud_gpdb_api_keys" "default" {
  workspace_id = "%s"
  ids          = ["${alicloud_gpdb_api_key.default.key_id}_fake"]
}
`, workspaceId, name, workspaceId)

	dsOutputFileConf := fmt.Sprintf(`
resource "alicloud_gpdb_api_key" "default" {
  workspace_id = "%s"
  key_name     = "%s"
  description  = "terraform test"
}

data "alicloud_gpdb_api_keys" "default" {
  workspace_id = "%s"
  output_file  = "./test_output_file"
}
`, workspaceId, name, workspaceId)

	// Expected attributes when data source finds one key.
	existCheck := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("data.alicloud_gpdb_api_keys.default", "ids.#", "1"),
		resource.TestCheckResourceAttr("data.alicloud_gpdb_api_keys.default", "keys.#", "1"),
		resource.TestCheckResourceAttrSet("data.alicloud_gpdb_api_keys.default", "keys.0.id"),
		resource.TestCheckResourceAttrSet("data.alicloud_gpdb_api_keys.default", "keys.0.key_id"),
		resource.TestCheckResourceAttr("data.alicloud_gpdb_api_keys.default", "keys.0.key_name", name),
	)
	// Expected attributes when data source finds no keys.
	emptyCheck := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("data.alicloud_gpdb_api_keys.default", "ids.#", "0"),
		resource.TestCheckResourceAttr("data.alicloud_gpdb_api_keys.default", "keys.#", "0"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck:  preCheck,
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// Case 1: basic list — Step 1 creates resource, Step 2 reads data source.
			{Config: resourceConf},
			{Config: dsBasicConf, Check: existCheck},
			// Case 2: ids filter (exist) — resource persists in state from
			// Case 1, so we only need the data-source step.
			{Config: dsIdsExistConf, Check: existCheck},
			// Case 3: ids filter (fake) — non-matching id yields empty set.
			{Config: dsIdsFakeConf, Check: emptyCheck},
			// Case 4: output_file — results are written to a file.
			{Config: dsOutputFileConf, Check: existCheck},
		},
	})
}
