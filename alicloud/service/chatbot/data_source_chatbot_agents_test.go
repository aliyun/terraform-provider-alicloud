package chatbot_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudChatbotAgentsDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(100, 999)
	acctest.CheckoutSupportedRegions(t, true, connectivity.ChatbotSupportRegions)
	AgentNameConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudChatbotAgentsDataSourceName(rand, map[string]string{
			"agent_name": `"default-NODELETING"`,
		}),
		FakeConfig: "",
	}
	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudChatbotAgentsDataSourceName(rand, map[string]string{
			"agent_name": `"default-NODELETING"`,
		}),
		FakeConfig: "",
	}
	var existAlicloudChatbotAgentsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":               "1",
			"names.#":             "1",
			"agents.#":            "1",
			"agents.0.agent_name": acctest.CHECKSET,
			"agents.0.id":         acctest.CHECKSET,
			"agents.0.agent_id":   acctest.CHECKSET,
			"agents.0.agent_key":  acctest.CHECKSET,
		}
	}
	var fakeAlicloudChatbotAgentsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":    "0",
			"names.#":  "0",
			"agents.#": "0",
		}
	}
	var alicloudChatbotAgentsBusesCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_chatbot_agents.default",
		ExistMapFunc: existAlicloudChatbotAgentsDataSourceNameMapFunc,
		FakeMapFunc:  fakeAlicloudChatbotAgentsDataSourceNameMapFunc,
	}
	preCheck := func() {
		acctest.PreCheck(t)
	}
	alicloudChatbotAgentsBusesCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, AgentNameConf, nameRegexConf)
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
