package dataworks_test

import (
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAlicloudDataWorksServiceDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudDataWorksServiceDataSource,
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckAlicloudDataSourceID("data.alicloud_dataworks_service.current"),
					resource.TestCheckResourceAttrSet("data.alicloud_dataworks_service.current", "id"),
					resource.TestCheckResourceAttr("data.alicloud_dataworks_service.current", "status", "Opened"),
				),
			},
		},
	})
}

const testAccCheckAlicloudDataWorksServiceDataSource = `
data "alicloud_dataworks_service" "current" {
	enable = "On"
}
`
