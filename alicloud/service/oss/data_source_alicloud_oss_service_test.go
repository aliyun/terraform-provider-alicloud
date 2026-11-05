package oss_test

import (
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudOssServiceDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudOssServiceDataSource,
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckAlicloudDataSourceID("data.alicloud_oss_service.current"),
					resource.TestCheckResourceAttrSet("data.alicloud_oss_service.current", "id"),
					resource.TestCheckResourceAttr("data.alicloud_oss_service.current", "status", "Opened"),
				),
			},
		},
	})
}

const testAccCheckAlicloudOssServiceDataSource = `
data "alicloud_oss_service" "current" {
	enable = "On"
}
`
