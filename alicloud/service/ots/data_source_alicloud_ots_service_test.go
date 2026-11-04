package ots_test

import (
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAlicloudOtsServiceDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudOtsServiceDataSource,
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckAlicloudDataSourceID("data.alicloud_ots_service.current"),
					resource.TestCheckResourceAttrSet("data.alicloud_ots_service.current", "id"),
					resource.TestCheckResourceAttr("data.alicloud_ots_service.current", "status", "Opened"),
				),
			},
		},
	})
}

const testAccCheckAlicloudOtsServiceDataSource = `
data "alicloud_ots_service" "current" {
	enable = "On"
}
`
