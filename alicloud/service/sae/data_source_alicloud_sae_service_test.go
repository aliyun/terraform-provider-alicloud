package sae_test

import (
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudSAEServiceDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudSaeServiceDataSource,
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckAlicloudDataSourceID("data.alicloud_sae_service.current"),
					resource.TestCheckResourceAttrSet("data.alicloud_sae_service.current", "id"),
					resource.TestCheckResourceAttr("data.alicloud_sae_service.current", "status", "Opened"),
				),
			},
		},
	})
}

const testAccCheckAlicloudSaeServiceDataSource = `
data "alicloud_sae_service" "current" {
	enable = "On"
}
`
