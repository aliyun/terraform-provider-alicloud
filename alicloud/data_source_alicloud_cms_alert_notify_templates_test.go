package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudCmsAlertNotifyTemplates_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "data.alicloud_cms_alert_notify_templates.default"
	ra := resourceAttrInit(resourceId, AlicloudCmsAlertNotifyTemplatesMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCmsAlertNotifyTemplate")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-test-ant-%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, AlicloudCmsAlertNotifyTemplatesBasicDependence)
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
					"ids": []string{name},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ids.#":       "1",
						"templates.#": "1",
					}),
				),
			},
		},
	})
}

func AlicloudCmsAlertNotifyTemplatesBasicDependence(name string) string {
	return fmt.Sprintf(`
resource "alicloud_cms_alert_notify_template" "default" {
  alert_notify_template_id   = "%[1]s"
  alert_notify_template_name = "%[1]s"
  type                       = "DING"
  program_lang               = "zh_CN"
  templates {
    channel = "ding"
    title   = "Alert Title"
    content = "Alert Content"
  }
}
`, name)
}

var AlicloudCmsAlertNotifyTemplatesMap = map[string]string{
	"ids.#": CHECKSET,
}
