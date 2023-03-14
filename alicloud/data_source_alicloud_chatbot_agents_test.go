package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudChatbotAgentsDataSource(t *testing.T) {
	rand := acctest.RandIntRange(100, 999)
	checkoutSupportedRegions(t, true, connectivity.ChatbotSupportRegions)
	AgentNameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudChatbotAgentsDataSourceName(rand, map[string]string{
			"agent_name": `"default-NODELETING"`,
		}),
		fakeConfig: "",
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudChatbotAgentsDataSourceName(rand, map[string]string{
			"agent_name": `"default-NODELETING"`,
		}),
		fakeConfig: "",
	}
	var existAlicloudChatbotAgentsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":               "1",
			"names.#":             "1",
			"agents.#":            "1",
			"agents.0.agent_name": CHECKSET,
			"agents.0.id":         CHECKSET,
			"agents.0.agent_id":   CHECKSET,
			"agents.0.agent_key":  CHECKSET,
		}
	}
	var fakeAlicloudChatbotAgentsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":    "0",
			"names.#":  "0",
			"agents.#": "0",
		}
	}
	var alicloudChatbotAgentsBusesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_chatbot_agents.default",
		existMapFunc: existAlicloudChatbotAgentsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudChatbotAgentsDataSourceNameMapFunc,
	}
	preCheck := func() {
		testAccPreCheck(t)
	}
	alicloudChatbotAgentsBusesCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, AgentNameConf, nameRegexConf)
}
func testAccCheckAlicloudChatbotAgentsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
variable "name" {
  default = "tf_testacc%d"
}
data "alicloud_chatbot_agents" "default" {	
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
