package cr_test

import (
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudCRServiceDataSource(t *testing.T) {
	resourceId := "data.alicloud_cr_service.current"
	testAccCheck := acctest.ResourceAttrInit(resourceId, map[string]string{}).ResourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudCrServiceDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"id":       acctest.CHECKSET,
						"status":   "Opened",
						"password": "1111aaaa",
					}),
				),
			},
		},
	})
}

const testAccCheckAlicloudCrServiceDataSource = `
data "alicloud_cr_service" "current" {
	enable 		= "On"
	password   	= "1111aaaa"
}
`
