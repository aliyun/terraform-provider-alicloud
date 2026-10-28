package mscopensubscription_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudMscSubWebhooksDataSource(t *testing.T) {
	rand := tfacctest.RandInt()
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudMscSubWebhooksDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_msc_sub_webhook.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudMscSubWebhooksDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_msc_sub_webhook.default.id}_fake"]`,
		}),
	}
	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudMscSubWebhooksDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_msc_sub_webhook.default.webhook_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudMscSubWebhooksDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_msc_sub_webhook.default.webhook_name}_fake"`,
		}),
	}
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudMscSubWebhooksDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_msc_sub_webhook.default.id}"]`,
			"name_regex": `"${alicloud_msc_sub_webhook.default.webhook_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudMscSubWebhooksDataSourceName(rand, map[string]string{
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

	var alicloudMscSubWebhookCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_msc_sub_webhooks.default",
		ExistMapFunc: existMscSubWebhooksMapFunc,
		FakeMapFunc:  fakeMscSubWebhooksMapFunc,
	}

	preCheck := func() {
		acctest.PreCheck(t)
		acctest.PreCheckWithEnvVariable(t, "ALICLOUD_MSC_SUB_WEBHOOK_TOKEN")
	}

	alicloudMscSubWebhookCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, nameRegexConf, idsConf, allConf)
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
