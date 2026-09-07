// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Cms AlertWebhook. >>> Resource test cases, automatically generated.
// Case AlertWebhook basic test
func TestAccAliCloudCmsAlertWebhook_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_alert_webhook.default"
	ra := resourceAttrInit(resourceId, AliCloudCmsAlertWebhookMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCmsAlertWebhook")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccms%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCmsAlertWebhookBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"alert_webhook_name": name,
					"url":                "https://example.com/alert-webhook",
					"content_type":       "JSON",
					"method":             "POST",
					"lang":               "zh_CN",
					"headers": map[string]interface{}{
						"X-Custom-Token": "token-initial",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"alert_webhook_name":     name,
						"url":                    "https://example.com/alert-webhook",
						"content_type":           "JSON",
						"method":                 "POST",
						"lang":                   "zh_CN",
						"headers.%":              "1",
						"headers.X-Custom-Token": "token-initial",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"alert_webhook_name": name + "_update",
					"url":                "https://example.com/alert-webhook-updated",
					"content_type":       "FORM",
					"method":             "GET",
					"lang":               "en_US",
					"headers": map[string]interface{}{
						"X-Custom-Token":  "token-updated",
						"X-Custom-Source": "terraform",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"alert_webhook_name":      name + "_update",
						"url":                     "https://example.com/alert-webhook-updated",
						"content_type":            "FORM",
						"method":                  "GET",
						"lang":                    "en_US",
						"headers.%":               "2",
						"headers.X-Custom-Token":  "token-updated",
						"headers.X-Custom-Source": "terraform",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

var AliCloudCmsAlertWebhookMap0 = map[string]string{
	"workspace": CHECKSET,
}

func AliCloudCmsAlertWebhookBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}
`, name)
}

// Test Cms AlertWebhook. >>> Resource test cases, automatically generated. <<<
