package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudCmsAlertHistoriesDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudCmsAlertHistoriesDataSourceConfig(acctest.RandInt()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.alicloud_cms_alert_histories.default", "id"),
					resource.TestCheckResourceAttr("data.alicloud_cms_alert_histories.default", "output_file", "alert_histories.json"),
				),
			},
		},
	})
}

func testAccCheckAlicloudCmsAlertHistoriesDataSourceConfig(rand int) string {
	return fmt.Sprintf(`
variable "name" {
  default = "tftestaccms-alerthist-%d"
}

data "alicloud_cms_alert_histories" "default" {
  alert_history_id     = "ah-fake-${var.name}"
  alert_rule_id        = "arule-fake-${var.name}"
  biz_source           = "fake-source"
  workspace            = "fake-ws-${var.name}"
  max_level            = "P1"
  latest_level         = "P2"
  instance_key         = "fake-key-${var.name}"
  display_name_keyword = "fake-keyword"
  ids                  = ["ah-fake-${var.name}"]
  label_filter {
    labels = {
      env = "fake"
    }
    opt = "AND"
  }
  output_file = "alert_histories.json"
}
`, rand)
}
