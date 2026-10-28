package kms_test

import (
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudKmsServiceDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudKmsServiceDataSource,
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckAlicloudDataSourceID("data.alicloud_kms_service.current"),
					resource.TestCheckResourceAttrSet("data.alicloud_kms_service.current", "id"),
					resource.TestCheckResourceAttr("data.alicloud_kms_service.current", "status", "Opened"),
				),
			},
		},
	})
}

const testAccCheckAlicloudKmsServiceDataSource = `
data "alicloud_kms_service" "current" {
	enable = "On"
}
`
