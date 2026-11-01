package cloudmonitorservice_test

import (
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudCmsServiceDataSource(t *testing.T) {
	resourceId := "data.alicloud_cms_service.current"
	testAccCheck := acctest.ResourceAttrInit(resourceId, map[string]string{}).ResourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudCmsServiceDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"id":     acctest.CHECKSET,
						"status": "Opened",
					}),
				),
			},
		},
	})
}

const testAccCheckAlicloudCmsServiceDataSource = `
data "alicloud_cms_service" "current" {
	enable = "On"
}
`
