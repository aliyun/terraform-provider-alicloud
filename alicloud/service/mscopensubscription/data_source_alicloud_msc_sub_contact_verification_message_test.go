package mscopensubscription_test

import (
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudMscSubContactVerificationMessageeDataSource(t *testing.T) {
	resourceId := "data.alicloud_msc_sub_contact_verification_message.default"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithEnvVariable(t, "NEED_REAL_MOBILE")
		},
		ProviderFactories: acctest.SDKProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudMscSubContactVerificationMessageDataSource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceId, "id"),
					resource.TestCheckResourceAttr(resourceId, "status", "Success"),
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
