package dcdn_test

import (
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudDcdnKvAccountDataSource(t *testing.T) {
	resourceId := "data.alicloud_dcdn_kv_account.current"
	testAccCheck := acctest.ResourceAttrInit(resourceId, map[string]string{}).ResourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudDcdnKvAccountDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"id":     acctest.CHECKSET,
						"status": "online",
					}),
				),
			},
		},
	})
}

const testAccCheckAlicloudDcdnKvAccountDataSource = `
data "alicloud_dcdn_kv_account" "current" {
		status = "online"
}
`
