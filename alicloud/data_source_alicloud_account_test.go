package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudAccountDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudAccountDataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_account.current"),
					resource.TestCheckResourceAttrSet("data.alicloud_account.current", "id"),
				),
			},
		},
	})
}

const testAccCheckAlicloudAccountDataSourceBasic = `
data "alicloud_account" "current" {
}
`
