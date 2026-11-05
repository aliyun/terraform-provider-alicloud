package videosurveillancesystem_test

import (
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudVsServiceDataSource(t *testing.T) {
	defer acctest.CheckoutAccount(t, false)
	acctest.CheckoutAccount(t, true)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudVsServiceDataSource,
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckAlicloudDataSourceID("data.alicloud_vs_service.current"),
					resource.TestCheckResourceAttrSet("data.alicloud_vs_service.current", "id"),
					resource.TestCheckResourceAttr("data.alicloud_vs_service.current", "status", "Opened"),
				),
			},
		},
	})
}

const testAccCheckAlicloudVsServiceDataSource = `
data "alicloud_vs_service" "current" {
	enable = "On"
}
`
