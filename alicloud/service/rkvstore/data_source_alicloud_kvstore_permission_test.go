package rkvstore_test

import (
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudKvstorePermissionDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudKVStorePermissionDataSource,
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckAlicloudDataSourceID("data.alicloud_kvstore_permission.current"),
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
