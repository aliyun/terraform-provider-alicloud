package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudDataWorksRemindDataSource_basic0(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sdataworksremindds%d", defaultRegionToTest, rand)
	dataSourceId := "data.alicloud_data_works_remind.default"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccAlicloudDataWorksRemindDataSourceConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceId, "remind_name", name),
					resource.TestCheckResourceAttr(dataSourceId, "remind_unit", "NODE"),
					resource.TestCheckResourceAttr(dataSourceId, "remind_type", "FINISHED"),
					resource.TestCheckResourceAttr(dataSourceId, "alert_unit", "OWNER"),
					resource.TestCheckResourceAttr(dataSourceId, "max_alert_times", "3"),
					resource.TestCheckResourceAttr(dataSourceId, "alert_interval", "300"),
					resource.TestCheckResourceAttr(dataSourceId, "useflag", "true"),
				),
			},
		},
	})
}

func testAccAlicloudDataWorksRemindDataSourceConfig(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alicloud_data_works_remind" "default" {
  remind_name     = var.name
  remind_unit     = "NODE"
  remind_type     = "FINISHED"
  alert_unit      = "OWNER"
  dnd_end         = "22:00"
  alert_methods   = ["DING_TALK"]
  alert_targets   = ["testuser@example.com"]
  useflag         = true
  max_alert_times = 3
  alert_interval  = 300
  detail          = "test remind rule for data source"
  webhooks        = ["https://oapi.dingtalk.com/robot/send?access_token=tf-test-token"]
  nodes           = []
  baselines       = []
  biz_processes   = []
  robots          = []
  projects        = []
}

data "alicloud_data_works_remind" "default" {
  remind_id = alicloud_data_works_remind.default.remind_id
}
`, name)
}
