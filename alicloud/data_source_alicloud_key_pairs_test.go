package alicloud

import (
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

func TestAccAlicloudKeyPairsDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudKeyPairsDataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_key_pairs.name_regex"),
				),
			},
		},
	})
}

const testAccCheckAlicloudKeyPairsDataSourceBasic = `
data "alicloud_key_pairs" "name_regex" {
	name_regex = "test"
}
`
