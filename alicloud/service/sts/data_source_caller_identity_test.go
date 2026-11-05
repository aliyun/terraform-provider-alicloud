package sts_test

import (
	"testing"

	tfacctest "github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliclouCallerIdentityDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tfacctest.PreCheck(t)
		},
		ProviderFactories: tfacctest.SDKProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: TestAccAliclouCallerIdentityDataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					tfacctest.CheckAlicloudDataSourceID("data.alicloud_caller_identity.cur"),
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
