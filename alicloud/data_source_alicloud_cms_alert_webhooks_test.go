// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Cms AlertWebhooks data source. >>> Data source test cases, automatically generated.
// Case AlertWebhooks data source basic test
func TestAccAliCloudCmsAlertWebhooks_basic0(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccms%d", rand)
	datasourceId := "data.alicloud_cms_alert_webhooks.default"

	idsExistConfig := testAccCheckAliCloudCmsAlertWebhooksSourceConfig(name, `ids = ["${alicloud_cms_alert_webhook.default.id}"]`)
	idsFakeConfig := testAccCheckAliCloudCmsAlertWebhooksSourceConfig(name, `ids = ["${alicloud_cms_alert_webhook.default.id}_fake"]`)
	nameRegexExistConfig := testAccCheckAliCloudCmsAlertWebhooksSourceConfig(name, fmt.Sprintf(`
	ids        = ["${alicloud_cms_alert_webhook.default.id}"]
	name_regex = "^%s$"`, name))
	nameRegexFakeConfig := testAccCheckAliCloudCmsAlertWebhooksSourceConfig(name, `
	ids        = ["${alicloud_cms_alert_webhook.default.id}"]
	name_regex = "^tf-nonexistent-webhook-regex$"`)
	nameFilterConfig := testAccCheckAliCloudCmsAlertWebhooksSourceConfig(name, `
	ids                = ["${alicloud_cms_alert_webhook.default.id}"]
	alert_webhook_name = "`+name+`"`)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: idsExistConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceId, "ids.#", "1"),
					resource.TestCheckResourceAttr(datasourceId, "names.#", "1"),
					resource.TestCheckResourceAttr(datasourceId, "webhooks.#", "1"),
					resource.TestCheckResourceAttr(datasourceId, "webhooks.0.alert_webhook_name", name),
					resource.TestCheckResourceAttrSet(datasourceId, "webhooks.0.alert_webhook_id"),
					resource.TestCheckResourceAttrSet(datasourceId, "webhooks.0.url"),
				),
			},
			{
				Config: idsFakeConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceId, "ids.#", "0"),
					resource.TestCheckResourceAttr(datasourceId, "webhooks.#", "0"),
				),
			},
			{
				Config: nameRegexExistConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceId, "ids.#", "1"),
					resource.TestCheckResourceAttr(datasourceId, "webhooks.#", "1"),
				),
			},
			{
				Config: nameRegexFakeConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceId, "ids.#", "0"),
					resource.TestCheckResourceAttr(datasourceId, "webhooks.#", "0"),
				),
			},
			{
				Config: nameFilterConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceId, "ids.#", "1"),
					resource.TestCheckResourceAttr(datasourceId, "webhooks.#", "1"),
					resource.TestCheckResourceAttr(datasourceId, "webhooks.0.alert_webhook_name", name),
				),
			},
		},
	})
}

func testAccCheckAliCloudCmsAlertWebhooksSourceConfig(name string, attr string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_cms_alert_webhook" "default" {
  alert_webhook_name = var.name
  url                = "https://example.com/alert-webhook-ds"
}

data "alicloud_cms_alert_webhooks" "default" {
%s
}
`, name, attr)
}

// Test Cms AlertWebhooks data source. <<< Data source test cases, automatically generated.
