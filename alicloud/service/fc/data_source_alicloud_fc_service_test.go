package fc_test

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"testing"
)

func TestAccAliCloudFCServiceDataSource(t *testing.T) {
	resourceId := "data.alicloud_fc_service.current"
	testAccCheck := acctest.ResourceAttrInit(resourceId, map[string]string{}).ResourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudFcServiceDataSource,
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

const testAccCheckAlicloudFcServiceDataSource = `
data "alicloud_fc_service" "current" {
	enable = "On"
}
`
