package alicloud

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudMscSubWebhooksDataSource(t *testing.T) {
	rand := acctest.RandInt()
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudMscSubWebhooksDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_msc_sub_webhook.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudMscSubWebhooksDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_msc_sub_webhook.default.id}_fake"]`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudMscSubWebhooksDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_msc_sub_webhook.default.webhook_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudMscSubWebhooksDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_msc_sub_webhook.default.webhook_name}_fake"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudMscSubWebhooksDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_msc_sub_webhook.default.id}"]`,
			"name_regex": `"${alicloud_msc_sub_webhook.default.webhook_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudMscSubWebhooksDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_msc_sub_webhook.default.id}_fake"]`,
			"name_regex": `"${alicloud_msc_sub_webhook.default.webhook_name}fake"`,
		}),
	}

	var existMscSubWebhooksMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                   "1",
			"webhooks.#":              "1",
			"webhooks.0.webhook_name": "testtfac",
			"webhooks.0.server_url":   "https://oapi.dingtalk.com/robot/send?access_token=" + os.Getenv("ALICLOUD_MSC_SUB_WEBHOOK_TOKEN"),
		}
	}

	var fakeMscSubWebhooksMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"webhooks.#": "0",
		}
	}

	var alicloudMscSubWebhookCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_msc_sub_webhooks.default",
		existMapFunc: existMscSubWebhooksMapFunc,
		fakeMapFunc:  fakeMscSubWebhooksMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
		testAccPreCheckWithEnvVariable(t, "ALICLOUD_MSC_SUB_WEBHOOK_TOKEN")
	}

	alicloudMscSubWebhookCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, nameRegexConf, idsConf, allConf)
}

func testAccCheckAlicloudMscSubWebhooksDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
variable "name" {	
  default = "testtfac"
}
resource "alicloud_msc_sub_webhook" "default" {
  webhook_name = var.name
  server_url = "https://oapi.dingtalk.com/robot/send?access_token=%s"
}
data "alicloud_msc_sub_webhooks" "default" {
  %s
}
`, os.Getenv("ALICLOUD_MSC_SUB_WEBHOOK_TOKEN"), strings.Join(pairs, "\n  "))
	return config
}
