package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudKMSServiceDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudKmsServiceDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_kms_service.current"),
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
