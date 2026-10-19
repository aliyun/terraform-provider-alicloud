package cms_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/cms"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

// TestAccAliCloudCmsAlertWebhook_basic0 creates an alert webhook with every
// configurable attribute, updates each attribute (name/url/content_type/
// method/lang/headers) and then imports the state to verify — 100% coverage of
// the Optional/Required surface, inherited from PR #10446's SDKv2 test.
//
// Note: PR #10446's SDKv2 test also asserted workspace via CHECKSET, but the
// CMS ListAlertWebhooks response does not populate workspace for a webhook
// created without one (verified end-to-end: id/name/url/content_type/method/
// lang/headers all read back, workspace is null). The Computed workspace field
// is retained in the schema (null when the API omits it, as the resource Read
// maps an empty/absent value to types.StringNull); only the set-assertion is
// dropped because it does not match actual API behavior and was never
// verified — PR #10446 was never merged.
//
// The harness is terraform-plugin-sdk/v2/helper/resource (not
// terraform-plugin-testing/helper/resource): the test binary links
// acctest.ProtoV5ProviderFactories → package alicloud, which (via polar_fs
// since #10623) already registers the -sweep flag from the SDK v2
// helper/resource init. Linking the testing fork's helper/resource on top
// redeclares -sweep and panics the binary before any test runs ("flag
// redefined: sweep", confirmed end-to-end). resource.TestCheckFunc therefore
// takes a *terraform.State from terraform-plugin-sdk/v2/terraform, and
// testAccCheckCmsAlertWebhookDestroy below uses that same State type.
func TestAccAliCloudCmsAlertWebhook_basic0(t *testing.T) {
	r := tfacctest.RandIntRange(100000, 999999)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckCmsAlertWebhookDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCmsAlertWebhookConfig(r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("alicloud_cms_alert_webhook.default", "id"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_webhook.default", "alert_webhook_name", fmt.Sprintf("tf-testacc-%d", r)),
					resource.TestCheckResourceAttr("alicloud_cms_alert_webhook.default", "url", fmt.Sprintf("https://example.com/alert-webhook-%d", r)),
					resource.TestCheckResourceAttr("alicloud_cms_alert_webhook.default", "content_type", "JSON"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_webhook.default", "method", "POST"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_webhook.default", "lang", "en_US"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_webhook.default", "headers.X-Token", fmt.Sprintf("tf-%d", r)),
				),
			},
			{
				Config: testAccCmsAlertWebhookConfigUpdate(r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_cms_alert_webhook.default", "alert_webhook_name", fmt.Sprintf("tf-testacc-%d-upd", r)),
					resource.TestCheckResourceAttr("alicloud_cms_alert_webhook.default", "url", fmt.Sprintf("https://example.com/alert-webhook-%d-upd", r)),
					resource.TestCheckResourceAttr("alicloud_cms_alert_webhook.default", "content_type", "FORM"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_webhook.default", "method", "GET"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_webhook.default", "lang", "zh_CN"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_webhook.default", "headers.X-Token", fmt.Sprintf("tf-%d-2", r)),
				),
			},
			{
				ResourceName:      "alicloud_cms_alert_webhook.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCmsAlertWebhookConfig(r int) string {
	return fmt.Sprintf(`
resource "alicloud_cms_alert_webhook" "default" {
  alert_webhook_name = "tf-testacc-%[1]d"
  url                = "https://example.com/alert-webhook-%[1]d"
  content_type       = "JSON"
  method             = "POST"
  lang               = "en_US"
  headers = {
    X-Token = "tf-%[1]d"
  }
}
`, r)
}

func testAccCmsAlertWebhookConfigUpdate(r int) string {
	return fmt.Sprintf(`
resource "alicloud_cms_alert_webhook" "default" {
  alert_webhook_name = "tf-testacc-%[1]d-upd"
  url                = "https://example.com/alert-webhook-%[1]d-upd"
  content_type       = "FORM"
  method             = "GET"
  lang               = "zh_CN"
  headers = {
    X-Token = "tf-%[1]d-2"
  }
}
`, r)
}

// testAccCheckCmsAlertWebhookDestroy verifies that every alert webhook created
// during the test has been removed after `terraform destroy`. The provider
// client is reached via acctest.Provider.Meta() — the SDK v2 half of the muxed
// provider shares the *connectivity.AliyunClient that the Framework half uses,
// and ProtoV5ProviderFactories drives ConfigureProvider against it, so Meta is
// populated by the time CheckDestroy runs. DescribeCmsAlertWebhook returns
// (nil, nil) for a gone webhook (the inherited not-found contract), so a
// non-nil object means the resource survived destroy.
func testAccCheckCmsAlertWebhookDestroy(s *terraform.State) error {
	client := acctest.Provider.Meta().(*connectivity.AliyunClient)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_cms_alert_webhook" {
			continue
		}
		obj, err := cms.NewService(client).DescribeCmsAlertWebhook(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("checking cms alert webhook %s destruction: %w", rs.Primary.ID, err)
		}
		if obj != nil {
			return fmt.Errorf("cms alert webhook %s still exists", rs.Primary.ID)
		}
	}
	return nil
}
