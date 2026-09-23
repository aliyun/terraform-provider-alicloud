package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudCRServiceDataSource(t *testing.T) {
	resourceId := "data.alicloud_cr_service.current"
	testAccCheck := resourceAttrInit(resourceId, map[string]string{}).resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudCrServiceDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"id":       CHECKSET,
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
