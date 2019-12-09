package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliclouStsCallerIdentityDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: TestAccAliclouStsCallerIdentityDataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_sts_caller_identity.cur"),
					resource.TestCheckResourceAttrSet("data.alicloud_sts_caller_identity.cur", "arn"),
					resource.TestCheckResourceAttrSet("data.alicloud_sts_caller_identity.cur", "account_id"),
					resource.TestCheckResourceAttrSet("data.alicloud_sts_caller_identity.cur", "identity_type"),
					resource.TestCheckResourceAttrSet("data.alicloud_sts_caller_identity.cur", "principal_id"),
				),
			},
		},
	})
}

const TestAccAliclouStsCallerIdentityDataSourceBasic = `
data "alicloud_sts_caller_identity" "cur" {
}
`
