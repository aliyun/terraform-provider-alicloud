package paiflow_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// alicloud_pai_flow_pipeline_run is a framework-native resource, so it is only
// exposed by the muxed (SDK v2 + framework) provider. This suite lives in the
// service package and serves the muxed provider through
// acctest.ProtoV5ProviderFactories, the same entry the root alicloud_test
// suites use.
//
// Harness choice: this suite uses the SDK v2 helper/resource (not
// terraform-plugin-testing/helper/resource). acctest imports package alicloud
// for Provider(); the baseline still has non-test root-package files importing
// SDK v2 helper/resource, which registers the -sweep flag in its init. Linking
// terraform-plugin-testing/helper/resource too would re-register -sweep and
// panic with "flag redefined: sweep" before any test runs. Staying on the SDK v2
// harness keeps a single -sweep registration. SDK v2 helper/resource.TestCase
// carries ProtoV5ProviderFactories, so the framework resource is exercised via
// the muxed server the same way the root suites do.

// PAIFlow pipeline manifest mirroring the proven alicloud_pai_flow_pipeline
// fixture (basic10139). The spec.pipelines[] entries reference real registered
// PAIFlow components by their fixed uuids (data_source, type_transform); a
// fabricated sub-pipeline uuid resolves to null and PAIFlow rejects the
// pipeline with "NO SUCH PIPELINE". The top-level pipeline uuid is made
// unique per run. The run may not reach Succeeded when the backing MaxCompute
// table is absent; the resource is an execution record, so Create settles on
// any terminal status and the test asserts the lifecycle.
const testAccPaiflowPipelineRunManifest = `
apiVersion: "core/v1"
metadata:
  provider: "%s"
  version: "v1"
  identifier: "my_pipeline"
  name: "source-transform"
  uuid: "tfpipelinerun%d"
  annotations: {}
  labels: {}
spec:
  inputs:
    artifacts: []
    parameters:
    - name: "execution_maxcompute"
      type: "Map"
      value:
        spec:
          endpoint: "http://service.cn.maxcompute.aliyun-inc.com/api"
          odpsProject: "test_i****"
  outputs:
    artifacts: []
    parameters: []
  arguments:
    artifacts: []
    parameters: []
  dependencies: []
  initContainers: []
  sideCarContainers: []
  pipelines:
  - apiVersion: "core/v1"
    metadata:
      provider: "pai"
      version: "v1"
      identifier: "data_source"
      name: "data-source"
      uuid: "2ftahdnzcod2rt6u9q"
      displayName: "read-table-1"
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
        parameters:
        - name: "inputTableName"
          value: "pai_online_project.wumai_data"
        - name: "execution"
          from: "{{inputs.parameters.execution_maxcompute}}"
      dependencies: []
      initContainers: []
      sideCarContainers: []
      pipelines: []
      volumes: []
  - apiVersion: "core/v1"
    metadata:
      provider: "pai"
      version: "v1"
      identifier: "type_transform"
      name: "type-transform"
      uuid: "gacnnnl4ksvbabfh6l"
      displayName: "type-transform-1"
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
        artifacts:
        - name: "inputTable"
          from: "{{pipelines.data_source.outputs.artifacts.outputTable}}"
        parameters:
        - name: "cols_to_double"
          value: "time,hour,pm2,pm10,so2,co,no2"
        - name: "execution"
          from: "{{inputs.parameters.execution_maxcompute}}"
      dependencies:
      - "data_source"
      initContainers: []
      sideCarContainers: []
      pipelines: []
      volumes: []
  volumes: []
`

func TestAccAliCloudPaiFlowPipelineRun(t *testing.T) {
	r := tfacctest.RandIntRange(10000, 99999)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			// The manifest's MaxCompute endpoint and the proven basic10139
			// fixture are scoped to cn-hangzhou.
			os.Setenv("ALICLOUD_REGION", "cn-hangzhou")
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		// SDK v2 helper/resource.TestCase.CheckDestroy expects a func typed
		// against the SDK v2 terraform.State, so the cross-harness
		// acctest.CheckDestroyNoop (typed for terraform-plugin-testing/terraform)
		// does not assign here. A local noop keeps the same no-verification
		// semantics used by other suites while satisfying tflint AT001.
		CheckDestroy: testAccCheckPaiFlowPipelineRunDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPaiflowPipelineRunConfig(r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("alicloud_pai_flow_pipeline_run.default", "id"),
					resource.TestCheckResourceAttrSet("alicloud_pai_flow_pipeline_run.default", "pipeline_run_id"),
					resource.TestCheckResourceAttrSet("alicloud_pai_flow_pipeline_run.default", "workspace_id"),
					resource.TestCheckResourceAttrSet("alicloud_pai_flow_pipeline_run.default", "pipeline_id"),
					resource.TestCheckResourceAttrSet("alicloud_pai_flow_pipeline_run.default", "status"),
					resource.TestCheckResourceAttr("alicloud_pai_flow_pipeline_run.default", "name", fmt.Sprintf("tfaccrun%d", r)),
					resource.TestCheckResourceAttr("alicloud_pai_flow_pipeline_run.default", "no_confirm_required", "true"),
				),
			},
			{
				Config: testAccPaiflowPipelineRunConfigUpdate(r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_pai_flow_pipeline_run.default", "name", fmt.Sprintf("tfaccrun%dupd", r)),
					resource.TestCheckResourceAttrSet("alicloud_pai_flow_pipeline_run.default", "pipeline_run_id"),
				),
			},
			{
				ResourceName:      "alicloud_pai_flow_pipeline_run.default",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					// GetPipelineRun does not return these round-trip:
					// manifest/arguments/no_confirm_required are create-only;
					// name is config-tracked because PAIFlow appends a suffix.
					"manifest", "arguments", "no_confirm_required", "name",
				},
			},
		},
	})
}

// testAccCheckPaiFlowPipelineRunDestroy is the post-destroy verification. The
// resource's Delete is idempotent and tolerant of not-found; the harness
// already drives the destroy step, so this noop mirrors the
// acctest.CheckDestroyNoop semantics used by other suites (tflint AT001
// requires a non-nil CheckDestroy). It is typed for the SDK v2 terraform.State
// to match this suite's SDK v2 harness (see the harness-choice note above).
func testAccCheckPaiFlowPipelineRunDestroy(s *terraform.State) error {
	return nil
}

func testAccPaiflowPipelineRunConfig(r int) string {
	return fmt.Sprintf(`
variable "name" {
  default = "tf_testaccworkspace%[1]d"
}

data "alicloud_account" "default" {}

resource "alicloud_pai_workspace_workspace" "default" {
  description    = "paiflow pipeline run test"
  display_name   = var.name
  workspace_name = var.name
  env_types      = ["prod"]
}

resource "alicloud_pai_flow_pipeline" "default" {
  workspace_id = alicloud_pai_workspace_workspace.default.id
  manifest     = <<-EOT
%[2]s
  EOT
}

resource "alicloud_pai_flow_pipeline_run" "default" {
  workspace_id        = alicloud_pai_workspace_workspace.default.id
  pipeline_id         = alicloud_pai_flow_pipeline.default.id
  name                = "tfaccrun%[1]d"
  no_confirm_required = true
}
`, r, fmt.Sprintf(testAccPaiflowPipelineRunManifest, "${data.alicloud_account.default.id}", r))
}

func testAccPaiflowPipelineRunConfigUpdate(r int) string {
	return fmt.Sprintf(`
variable "name" {
  default = "tf_testaccworkspace%[1]d"
}

data "alicloud_account" "default" {}

resource "alicloud_pai_workspace_workspace" "default" {
  description    = "paiflow pipeline run test"
  display_name   = var.name
  workspace_name = var.name
  env_types      = ["prod"]
}

resource "alicloud_pai_flow_pipeline" "default" {
  workspace_id = alicloud_pai_workspace_workspace.default.id
  manifest     = <<-EOT
%[2]s
  EOT
}

resource "alicloud_pai_flow_pipeline_run" "default" {
  workspace_id        = alicloud_pai_workspace_workspace.default.id
  pipeline_id         = alicloud_pai_flow_pipeline.default.id
  name                = "tfaccrun%[1]dupd"
  no_confirm_required = true
}
`, r, fmt.Sprintf(testAccPaiflowPipelineRunManifest, "${data.alicloud_account.default.id}", r))
}
