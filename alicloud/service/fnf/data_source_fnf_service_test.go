package fnf_test

import (
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAlicloudFnfServiceDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudFnfServiceDataSource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.alicloud_fnf_service.current", "id"),
					resource.TestCheckResourceAttr("data.alicloud_fnf_service.current", "status", "Opened"),
				),
			},
		},
	})
}

const testAccCheckAlicloudFnfServiceDataSource = `
data "alicloud_fnf_service" "current" {
	enable = "On"
}
`
