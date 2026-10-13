package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAliCloudDcdnKvAccountDataSource(t *testing.T) {
	resourceId := "data.alicloud_dcdn_kv_account.current"
	testAccCheck := resourceAttrInit(resourceId, map[string]string{}).resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudDcdnKvAccountDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"id":     CHECKSET,
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
