package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliclouCallerIdentityDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: TestAccAliclouCallerIdentityDataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_caller_identity.cur"),
					resource.TestCheckResourceAttrSet("data.alicloud_caller_identity.cur", "arn"),
					resource.TestCheckResourceAttrSet("data.alicloud_caller_identity.cur", "account_id"),
					resource.TestCheckResourceAttrSet("data.alicloud_caller_identity.cur", "identity_type"),
				),
			},
		},
	})
}

const TestAccAliclouCallerIdentityDataSourceBasic = `
data "alicloud_caller_identity" "cur" {
}
`
