package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAliCloudAgentLoopAgentSpacesDataSource(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccagentloop%d", rand)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudAgentLoopAgentSpacesDataSourceConfig(name, map[string]string{
			"ids": `[alicloud_agentloop_agent_space.default.id]`,
		}),
		fakeConfig: testAccCheckAliCloudAgentLoopAgentSpacesDataSourceConfig(name, map[string]string{
			"ids": `["agentspace-fake00000"]`,
		}),
	}

	agentSpaceConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudAgentLoopAgentSpacesDataSourceConfig(name, map[string]string{
			"agent_space": `"${alicloud_agentloop_agent_space.default.agent_space}"`,
		}),
		fakeConfig: testAccCheckAliCloudAgentLoopAgentSpacesDataSourceConfig(name, map[string]string{
			"agent_space": `"${alicloud_agentloop_agent_space.default.agent_space}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudAgentLoopAgentSpacesDataSourceConfig(name, map[string]string{
			"ids":         `[alicloud_agentloop_agent_space.default.id]`,
			"agent_space": `"${alicloud_agentloop_agent_space.default.agent_space}"`,
		}),
		fakeConfig: testAccCheckAliCloudAgentLoopAgentSpacesDataSourceConfig(name, map[string]string{
			"ids":         `["agentspace-fake00000"]`,
			"agent_space": `"${alicloud_agentloop_agent_space.default.agent_space}_fake"`,
		}),
	}

	var existAliCloudAgentLoopAgentSpacesDataSourceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                "1",
			"spaces.#":             "1",
			"spaces.0.id":          CHECKSET,
			"spaces.0.agent_space": CHECKSET,
			"spaces.0.description": CHECKSET,
			"spaces.0.create_time": CHECKSET,
			"spaces.0.region_id":   CHECKSET,
			"spaces.0.update_time": CHECKSET,
		}
	}

	var fakeAliCloudAgentLoopAgentSpacesDataSourceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":    "0",
			"spaces.#": "0",
		}
	}

	aliCloudAgentLoopAgentSpacesCheckInfo := dataSourceAttr{
		resourceId:   "data.alicloud_agentloop_agent_spaces.default",
		existMapFunc: existAliCloudAgentLoopAgentSpacesDataSourceMapFunc,
		fakeMapFunc:  fakeAliCloudAgentLoopAgentSpacesDataSourceMapFunc,
	}

	preCheck := func() {
		testAccPreCheckAgentLoopRegion(t)
	}
	aliCloudAgentLoopAgentSpacesCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, agentSpaceConf, allConf)
}

func testAccCheckAliCloudAgentLoopAgentSpacesDataSourceConfig(name string, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alicloud_agentloop_agent_space" "default" {
  agent_space = var.name
  description = var.name
}

data "alicloud_agentloop_agent_spaces" "default" {
  %s
}
`, name, strings.Join(pairs, "\n  "))
}
