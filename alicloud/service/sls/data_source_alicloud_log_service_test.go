package sls_test

import (
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudLogServiceDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudLogServiceDataSource,
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckAlicloudDataSourceID("data.alicloud_log_service.current"),
					resource.TestCheckResourceAttrSet("data.alicloud_log_service.current", "id"),
					resource.TestCheckResourceAttr("data.alicloud_log_service.current", "status", "Opened"),
				),
			},
		},
	})
}

const testAccCheckAlicloudLogServiceDataSource = `
data "alicloud_log_service" "current" {
	enable = "On"
}
`
