package cms_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

// TestAccAliCloudCmsAlertWebhooks_basic0 exercises the alert-webhooks data
// source against the alert webhook created in the same run, filtered by ids.
//
// The ids filter maps to the ListAlertWebhooks webhookIds server-side query
// parameter (a direct lookup, the same path the resource Read uses after
// Create). The bare GET /webhooks index is observed to omit a webhook created
// without a workspace, while webhookIds returns it immediately, so the test
// filters by the created webhook's id rather than name_regex (which would run
// against the bare index and find nothing for this account).
func TestAccAliCloudCmsAlertWebhooks_basic0(t *testing.T) {
	r := tfacctest.RandIntRange(100000, 999999)
	name := fmt.Sprintf("tf-testacc-%d", r)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckCmsAlertWebhookDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCmsAlertWebhooksConfigExist(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("alicloud_cms_alert_webhook.default", "id"),
					resource.TestCheckResourceAttrSet("data.alicloud_cms_alert_webhooks.default", "ids.#"),
					resource.TestCheckResourceAttrSet("data.alicloud_cms_alert_webhooks.default", "names.#"),
					resource.TestCheckResourceAttrSet("data.alicloud_cms_alert_webhooks.default", "webhooks.#"),
					resource.TestCheckResourceAttr("data.alicloud_cms_alert_webhooks.default", "webhooks.0.alert_webhook_name", name),
					resource.TestCheckResourceAttrSet("data.alicloud_cms_alert_webhooks.default", "webhooks.0.alert_webhook_id"),
					resource.TestCheckResourceAttr("data.alicloud_cms_alert_webhooks.default", "webhooks.0.url", "https://example.com/alert-webhook-ds"),
				),
			},
		},
	})
}

func testAccCmsAlertWebhooksConfigExist(name string) string {
	return fmt.Sprintf(`
resource "alicloud_cms_alert_webhook" "default" {
  alert_webhook_name = "%[1]s"
  url                = "https://example.com/alert-webhook-ds"
  content_type       = "JSON"
  method             = "POST"
  lang               = "en_US"
}

data "alicloud_cms_alert_webhooks" "default" {
  ids = [alicloud_cms_alert_webhook.default.id]
}
`, name)
}
