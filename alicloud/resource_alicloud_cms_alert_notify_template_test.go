package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Cms AlertNotifyTemplate. >>> Resource test cases.
// Case AlertNotifyTemplate basic test
func TestAccAliCloudCmsAlertNotifyTemplate_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_alert_notify_template.default"
	ra := resourceAttrInit(resourceId, AlicloudCmsAlertNotifyTemplateMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCmsAlertNotifyTemplate")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccms%d", rand)
	templateId := fmt.Sprintf("tf-test-ant-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCmsAlertNotifyTemplateBasicDependence)
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
					"alert_notify_template_id":   templateId,
					"alert_notify_template_name": name,
					"type":                       "DING",
					"program_lang":               "zh_CN",
					"templates": []map[string]interface{}{
						{
							"channel": "ding",
							"title":   "Alert Title",
							"content": "Alert Content",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"alert_notify_template_id":   templateId,
						"alert_notify_template_name": name,
						"type":                       "DING",
						"program_lang":               "zh_CN",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"alert_notify_template_id":   templateId,
					"alert_notify_template_name": name + "-updated",
					"type":                       "DING",
					"program_lang":               "zh_CN",
					"templates": []map[string]interface{}{
						{
							"channel": "mail",
							"title":   "Updated Alert Title",
							"content": "Updated Alert Content",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"alert_notify_template_name": name + "-updated",
						"templates.0.channel":        "mail",
						"templates.0.title":          "Updated Alert Title",
						"templates.0.content":        "Updated Alert Content",
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

func AlicloudCmsAlertNotifyTemplateBasicDependence(name string) string {
	return ""
}

var AlicloudCmsAlertNotifyTemplateMap = map[string]string{
	"alert_notify_template_id":   CHECKSET,
	"alert_notify_template_name": CHECKSET,
	"type":                       CHECKSET,
	"program_lang":               CHECKSET,
	"region_id":                  CHECKSET,
}
