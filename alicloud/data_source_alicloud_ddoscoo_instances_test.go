package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudDdoscooInstanceDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudDdoscooInstanceDataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_ddoscoo_instances.foo"),
				),
			},
		},
	})
}

const testAccCheckAlicloudDdoscooInstanceDataSourceBasic = `
data "alicloud_ddoscoo_instances" "foo" {
    id = "tf-AccCheckAlicloudDdoscooInstanceDataSourceBasic"
}
`
