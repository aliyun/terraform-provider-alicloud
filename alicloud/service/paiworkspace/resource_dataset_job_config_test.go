package paiworkspace_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"

	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// This test lives in package paiworkspace_test (NOT package alicloud) so it may
// import terraform-plugin-sdk/v2/helper/resource. Per the acctest.go
// invariant, the -sweep flag conflict only arises when a single test binary
// links BOTH terraform-plugin-sdk/v2/helper/resource and
// terraform-plugin-testing/helper/resource. Using the SDK v2 TestCase here
// means only the SDK v2 helper/resource is linked (polardb_polar_fs production
// code already imports it), so no "flag redefined: sweep" panic, and the
// muxed provider is still served via ProtoV5ProviderFactories.

func TestAccAliCloudPaiWorkspaceDatasetJobConfig_basic(t *testing.T) {
	rand := sdkacctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tf_testacc_djc_%d", rand)
	resourceId := "alicloud_pai_workspace_dataset_job_config.default"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		IDRefreshName:            resourceId,
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckPaiWorkspaceDatasetJobConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPaiWorkspaceDatasetJobConfigBasic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "config_type", "MultimodalIntelligentTag"),
					resource.TestCheckResourceAttrSet(resourceId, "dataset_job_config_id"),
					resource.TestCheckResourceAttrSet(resourceId, "create_time"),
					resource.TestCheckResourceAttrSet(resourceId, "workspace_id"),
				),
			},
			{
				Config: testAccPaiWorkspaceDatasetJobConfigUpdate(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "config_type", "MultimodalIntelligentTag"),
					resource.TestCheckResourceAttr(resourceId, "config", `{"apiKey":"sk-updatedkey000000"}`),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dataset_version"},
			},
		},
	})
}

// testAccPaiWorkspaceDatasetJobConfigDependence builds the prerequisite chain
// (workspace + dataset) the job config hangs off. The dataset mirrors the
// known-working PAI ACC fixture (NAS-backed, cn-hangzhou).
func testAccPaiWorkspaceDatasetJobConfigDependence(name string) string {
	return fmt.Sprintf(`
resource "alicloud_pai_workspace_workspace" "defaultWorkspace" {
  workspace_name = "%[1]s_ws"
  display_name   = "DatasetJobConfigTest"
  description    = "tf testacc dataset job config"
  env_types      = ["prod"]
}

resource "alicloud_pai_workspace_dataset" "default" {
  workspace_id    = alicloud_pai_workspace_workspace.defaultWorkspace.id
  dataset_name    = "%[1]s_ds"
  data_source_type = "NAS"
  data_type       = "PIC"
  property        = "DIRECTORY"
  uri             = "nas://086b649545.cn-hangzhou/"
  accessibility   = "PRIVATE"
  source_type     = "ITAG"
  source_id       = "d-xxxxx_v1"
  user_id         = "1511928242963727"
}
`, name)
}

func testAccPaiWorkspaceDatasetJobConfigBasic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "alicloud_pai_workspace_dataset_job_config" "default" {
  workspace_id = alicloud_pai_workspace_workspace.defaultWorkspace.id
  dataset_id   = alicloud_pai_workspace_dataset.default.id
  config_type  = "MultimodalIntelligentTag"
  config       = jsonencode({ apiKey = "sk-testkey000000" })
  dataset_version = "v1"
}
`, testAccPaiWorkspaceDatasetJobConfigDependence(name))
}

func testAccPaiWorkspaceDatasetJobConfigUpdate(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "alicloud_pai_workspace_dataset_job_config" "default" {
  workspace_id = alicloud_pai_workspace_workspace.defaultWorkspace.id
  dataset_id   = alicloud_pai_workspace_dataset.default.id
  config_type  = "MultimodalIntelligentTag"
  config       = jsonencode({ apiKey = "sk-updatedkey000000" })
  dataset_version = "v1"
}
`, testAccPaiWorkspaceDatasetJobConfigDependence(name))
}

// testAccCheckPaiWorkspaceDatasetJobConfigDestroy verifies the config is gone
// after terraform destroy: it calls GetDatasetJobConfig directly and requires a
// not-found response. It mirrors the production Read/Delete not-found detection
// (tea.SDKError 404 via errs.NotFoundError plus the PAI-specific code/message
// fallback), because this external test package cannot reach the unexported
// helpers.
func testAccCheckPaiWorkspaceDatasetJobConfigDestroy(s *terraform.State) error {
	client := acctest.Provider.Meta().(*connectivity.AliyunClient)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_pai_workspace_dataset_job_config" {
			continue
		}
		parts := strings.Split(rs.Primary.ID, ":")
		if len(parts) != 3 {
			return fmt.Errorf("unexpected ID format %q, expected <workspace_id>:<dataset_id>:<dataset_job_config_id>", rs.Primary.ID)
		}
		workspaceId, datasetId, datasetJobConfigId := parts[0], parts[1], parts[2]
		action := fmt.Sprintf("/api/v1/datasets/%s/datasetjobconfigs/%s", datasetId, datasetJobConfigId)
		ws := workspaceId
		query := map[string]*string{"WorkspaceId": &ws}
		if _, err := client.RoaGet("AIWorkSpace", "2021-02-04", action, query, nil, nil); err != nil {
			if errs.NotFoundError(err) || isDatasetJobConfigGoneMessage(err) {
				continue
			}
			return fmt.Errorf("checking destroy of dataset job config %s: %w", rs.Primary.ID, err)
		}
		return fmt.Errorf("dataset job config %s still exists", rs.Primary.ID)
	}
	return nil
}

// isDatasetJobConfigGoneMessage covers the PAI-specific not-found signals that
// arrive as plain message strings (error codes 201300002/201300003, the
// "not exist"/"notfound" substrings, and the NoPermissionError gone signal PAI
// returns for a deleted config) rather than a transport-level 404. It mirrors
// the production isDatasetJobConfigNotExist detection.
func isDatasetJobConfigGoneMessage(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "201300002") || strings.Contains(msg, "201300003") ||
		strings.Contains(msg, "notfound") || strings.Contains(msg, "not exist") ||
		strings.Contains(msg, "nopermissionerror") || strings.Contains(msg, "no permission for this action")
}
