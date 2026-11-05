package account_test

import (
	"testing"

	tfacctest "github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAlicloudAccountDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tfacctest.PreCheck(t)
		},
		ProviderFactories: tfacctest.SDKProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudAccountDataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					tfacctest.CheckAlicloudDataSourceID("data.alicloud_account.current"),
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
