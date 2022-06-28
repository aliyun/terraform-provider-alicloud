package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudKVStorePermissionDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudKVStorePermissionDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_kvstore_permission.current"),
					resource.TestCheckResourceAttrSet("data.alicloud_kvstore_permission.current", "id"),
					resource.TestCheckResourceAttr("data.alicloud_kvstore_permission.current", "status", "Initialized"),
				),
			},
		},
	})
}

const testAccCheckAlicloudKVStorePermissionDataSource = `
data "alicloud_kvstore_permission" "current" {
	enable = "On"
}
`
