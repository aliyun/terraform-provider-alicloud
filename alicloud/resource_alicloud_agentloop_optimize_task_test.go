package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Agentloop OptimizeTask. >>> Resource test cases, automatically generated.
// Case OptimizeTask全生命周期 11085
func TestAccAliCloudAgentloopOptimizeTask_basic11085(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_agentloop_optimize_task.default"
	ra := resourceAttrInit(resourceId, AlicloudAgentloopOptimizeTaskMap11085)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AgentloopServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAgentloopOptimizeTask")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tftestaccopttask%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudAgentloopOptimizeTaskBasicDependence11085)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"agent_space":        "${alicloud_agentloop_agent_space.default.agent_space}",
					"optimize_task_name": name,
					"type":               "EfficiencyAnalysis",
					"status":             "Enabled",
					"config": map[string]interface{}{
						"config1": "value1",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"agent_space":        name + "-as",
						"optimize_task_name": name,
						"type":               "EfficiencyAnalysis",
						"status":             "Enabled",
						"config.%":           "1",
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
	"create_time": CHECKSET,
	"update_time": CHECKSET,
	"region_id":   CHECKSET,
}

func AlicloudAgentloopOptimizeTaskBasicDependence11085(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_agentloop_agent_space" "default" {
    agent_space = "${var.name}-as"
}

`, name)
}

// Test Agentloop OptimizeTask. <<< Resource test cases, automatically generated.
