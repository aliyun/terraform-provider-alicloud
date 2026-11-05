package agentloop_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/agentloop"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test Agentloop OptimizeTask. >>> Resource test cases, automatically generated.
// Case OptimizeTask全生命周期 11085
func TestAccAliCloudAgentloopOptimizeTask_basic11085(t *testing.T) {
	resourceId := "alicloud_agentloop_optimize_task.default"
	var v map[string]interface{}
	ra := acctest.ResourceAttrInit(resourceId, AlicloudAgentloopOptimizeTaskMap11085)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return agentloop.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeOptimizeTask")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tftestaccopttask%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudAgentloopOptimizeTaskBasicDependence11085)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{connectivity.Hangzhou})
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"agent_space":        "${alicloud_agentloop_agent_space.default.agent_space}",
					"optimize_task_name": name,
					"type":               "EfficiencyAnalysis",
					"status":             "Enabled",
					"config":             `{\"ioVersion\":2,\"inputs\":{\"trajectory\":{\"type\":\"sls\",\"project\":\"${alicloud_agentloop_agent_space.default.sls_project}\",\"logstore\":\"agent-trajectory\",\"format\":\"atif\"},\"featureEvents\":{\"type\":\"sls\",\"project\":\"${alicloud_agentloop_agent_space.default.sls_project}\",\"logstore\":\"trajectory-feature-event\",\"format\":\"trajectory.feature-event\"}},\"outputs\":{\"analysis\":{\"type\":\"sls\",\"project\":\"${alicloud_agentloop_agent_space.default.sls_project}\",\"logstore\":\"mining-token-analysis\",\"format\":\"token.analysis/v1\"},\"suggestions\":{\"type\":\"dataset\",\"dataset\":\"${alicloud_agentloop_dataset.suggestions.dataset_name}\"}}}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"agent_space":        name + "-as",
						"optimize_task_name": name,
						"type":               "EfficiencyAnalysis",
						"status":             "Enabled",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Disabled",
					"config": `{\"ioVersion\":2,\"inputs\":{\"trajectory\":{\"type\":\"sls\",\"project\":\"${alicloud_agentloop_agent_space.default.sls_project}\",\"logstore\":\"agent-trajectory\",\"format\":\"atif\"},\"featureEvents\":{\"type\":\"sls\",\"project\":\"${alicloud_agentloop_agent_space.default.sls_project}\",\"logstore\":\"trajectory-feature-event\",\"format\":\"trajectory.feature-event\"}},\"outputs\":{\"analysis\":{\"type\":\"sls\",\"project\":\"${alicloud_agentloop_agent_space.default.sls_project}\",\"logstore\":\"mining-token-analysis\",\"format\":\"token.analysis/v1\"},\"suggestions\":{\"type\":\"dataset\",\"dataset\":\"${alicloud_agentloop_dataset.suggestions_updated.dataset_name}\"}}}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Disabled",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Disabled",
					"config": `{\"ioVersion\":2,\"inputs\":{\"trajectory\":{\"type\":\"sls\",\"project\":\"${alicloud_agentloop_agent_space.default.sls_project}\",\"logstore\":\"agent-trajectory\",\"format\":\"atif\"},\"featureEvents\":{\"type\":\"sls\",\"project\":\"${alicloud_agentloop_agent_space.default.sls_project}\",\"logstore\":\"trajectory-feature-event\",\"format\":\"trajectory.feature-event\"}},\"outputs\":{\"analysis\":{\"type\":\"sls\",\"project\":\"${alicloud_agentloop_agent_space.default.sls_project}\",\"logstore\":\"mining-token-analysis\",\"format\":\"token.analysis/v1\"},\"suggestions\":{\"type\":\"dataset\",\"dataset\":\"${alicloud_agentloop_dataset.suggestions_updated.dataset_name}\"}}}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Disabled",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
				// update_time is refreshed by the backend asynchronously, so the
				// value captured at apply time may differ from the import read.
				ImportStateVerifyIgnore: []string{"update_time"},
			},
		},
	})
}

var AlicloudAgentloopOptimizeTaskMap11085 = map[string]string{
	"create_time": acctest.CHECKSET,
	"update_time": acctest.CHECKSET,
	"region_id":   acctest.CHECKSET,
}

func AlicloudAgentloopOptimizeTaskBasicDependence11085(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alicloud_agentloop_agent_space" "default" {
  agent_space = "${var.name}-as"
}

locals {
  optimize_task_suggestions_schema = {
    candidate_id    = jsonencode({ type = "text" })
    schema_version  = jsonencode({ type = "text" })
    task_type       = jsonencode({ type = "text" })
    run_id          = jsonencode({ type = "text" })
    task_id         = jsonencode({ type = "text" })
    output_name     = jsonencode({ type = "text" })
    artifact_schema = jsonencode({ type = "text" })
    payload         = jsonencode({ type = "text", chn = true })
    created_at      = jsonencode({ type = "long" })
  }
}

resource "alicloud_agentloop_dataset" "suggestions" {
  agent_space  = alicloud_agentloop_agent_space.default.agent_space
  dataset_name = "${var.name}_suggestions"
  description  = "OptimizeTask suggestions output"
  schema       = local.optimize_task_suggestions_schema
}

resource "alicloud_agentloop_dataset" "suggestions_updated" {
  agent_space  = alicloud_agentloop_agent_space.default.agent_space
  dataset_name = "${var.name}_suggestions_updated"
  description  = "Updated OptimizeTask suggestions output"
  schema       = local.optimize_task_suggestions_schema
}`, name)
}

// Test Agentloop OptimizeTask. <<< Resource test cases, automatically generated.
