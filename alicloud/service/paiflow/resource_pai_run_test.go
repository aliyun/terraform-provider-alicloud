package paiflow_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccAliCloudPaiRun exercises the PAI Flow pipeline run resource end to end:
// a workspace, a pipeline defined by manifest, and a run submitted from that
// pipeline. The run is fire-and-forget: Create submits and starts it, status is
// read back from GetPipelineRun. Downstream execution may fail on accounts
// without MaxCompute resources, but the resource (the run record) is created.
func TestAccAliCloudPaiRun(t *testing.T) {
	r := tfacctest.RandIntRange(10000, 99999)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudPaiRunConfig(r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("alicloud_pai_run.default", "id"),
					resource.TestCheckResourceAttrSet("alicloud_pai_run.default", "status"),
					resource.TestCheckResourceAttrSet("alicloud_pai_run.default", "workspace_id"),
					resource.TestCheckResourceAttrSet("alicloud_pai_run.default", "gmt_create_time"),
				),
			},
			{
				Config:            testAccAliCloudPaiRunConfigWithSource(r),
				ResourceName:      "alicloud_pai_run.default",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					// pipeline_manifest is a pure input the API does not echo back;
					// it is preserved from config on refresh, but import has no
					// config to recover it from.
					"pipeline_manifest",
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_pai_run.default", "source_type", "DESIGNER"),
					resource.TestCheckResourceAttr("alicloud_pai_run.default", "accessibility", "PRIVATE"),
				),
			},
		},
	})
}

func testAccAliCloudPaiRunConfig(r int) string {
	return fmt.Sprintf(`
variable "name" {
  default = "tf-testacc-pai-run-%[1]d"
}

data "alicloud_account" "default" {
}

resource "alicloud_pai_workspace_workspace" "default" {
  description    = var.name
  display_name   = var.name
  workspace_name = var.name
  env_types      = ["dev"]
}

resource "alicloud_pai_flow_pipeline" "default" {
  workspace_id = alicloud_pai_workspace_workspace.default.id
  manifest     = <<EOF
apiVersion: "core/v1"
metadata:
  provider: "${data.alicloud_account.default.id}"
  version: "v1"
  identifier: "my_pipeline"
  name: "source-transform"
  uuid: "terraformuuid%[1]d"
  annotations: {}
  labels: {}
spec:
  inputs:
    artifacts: []
    parameters: []
  outputs:
    artifacts: []
    parameters: []
  arguments:
    artifacts: []
    parameters: []
  dependencies: []
  initContainers: []
  sideCarContainers: []
  pipelines: []
  volumes: []
EOF
}

resource "alicloud_pai_run" "default" {
  workspace_id = alicloud_pai_workspace_workspace.default.id
  pipeline_id  = alicloud_pai_flow_pipeline.default.id
}
`, r)
}

func testAccAliCloudPaiRunConfigWithSource(r int) string {
	return fmt.Sprintf(`
variable "name" {
  default = "tf-testacc-pai-run-%[1]d"
}

data "alicloud_account" "default" {
}

resource "alicloud_pai_workspace_workspace" "default" {
  description    = var.name
  display_name   = var.name
  workspace_name = var.name
  env_types      = ["dev"]
}

resource "alicloud_pai_flow_pipeline" "default" {
  workspace_id = alicloud_pai_workspace_workspace.default.id
  manifest     = <<EOF
apiVersion: "core/v1"
metadata:
  provider: "${data.alicloud_account.default.id}"
  version: "v1"
  identifier: "my_pipeline"
  name: "source-transform"
  uuid: "terraformuuid%[1]d"
  annotations: {}
  labels: {}
spec:
  inputs:
    artifacts: []
    parameters: []
  outputs:
    artifacts: []
    parameters: []
  arguments:
    artifacts: []
    parameters: []
  dependencies: []
  initContainers: []
  sideCarContainers: []
  pipelines: []
  volumes: []
EOF
}

resource "alicloud_pai_run" "default" {
  workspace_id  = alicloud_pai_workspace_workspace.default.id
  pipeline_id   = alicloud_pai_flow_pipeline.default.id
  source_type   = "DESIGNER"
  accessibility = "PRIVATE"
}
`, r)
}
