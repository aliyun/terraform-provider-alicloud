package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudMscSubContactVerificationMessageeDataSource(t *testing.T) {
	resourceId := "data.alicloud_msc_sub_contact_verification_message.default"
	testAccCheck := resourceAttrInit(resourceId, map[string]string{}).resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithEnvVariable(t, "NEED_REAL_MOBILE")
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudMscSubContactVerificationMessageDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"id":     CHECKSET,
						"status": "Success",
					}),
				),
			},
		},
	})
}

const testAccCheckAlicloudMscSubContactVerificationMessageDataSource = `
resource "alicloud_msc_sub_contact" "default" {
  contact_name = "tftest"
  position     = "Other"
  email        = "123@163.com"
  mobile       = "153xxxxx906"
}

data "alicloud_msc_sub_contact_verification_message" "default" {
  contact_id = alicloud_msc_sub_contact.default.id
  type       = 1
}
`
