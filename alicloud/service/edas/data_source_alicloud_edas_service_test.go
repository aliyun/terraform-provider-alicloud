package edas_test

import (
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudEdasServiceDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudEdasServiceDataSource,
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckAlicloudDataSourceID("data.alicloud_edas_service.current"),
					resource.TestCheckResourceAttrSet("data.alicloud_edas_service.current", "id"),
					resource.TestCheckResourceAttr("data.alicloud_edas_service.current", "status", "Opened"),
				),
			},
		},
	})
}

const testAccCheckAlicloudEdasServiceDataSource = `
data "alicloud_edas_service" "current" {
	enable = "On"
}
`
