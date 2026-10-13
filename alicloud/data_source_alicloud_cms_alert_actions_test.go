package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudCmsAlertActionsDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudCmsAlertActionsDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.alicloud_cms_alert_actions.current", "alert_actions.#"),
				),
			},
		},
	})
}

func testAccAliCloudCmsAlertActionsDataSourceConfig() string {
	return `
data "alicloud_cms_alert_actions" "current" {
  type = "WEBHOOK"
  output_file = "alert_actions_output.json"
}
`
}
