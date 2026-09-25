package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudDataWorksRemind_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_data_works_remind.default"
	ra := resourceAttrInit(resourceId, AlicloudDataWorksRemindMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DataworksPublicService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDataWorksRemind")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sdataworksremind%d", defaultRegionToTest, rand)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccAlicloudDataWorksRemindBasicConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"remind_name":     name,
						"remind_unit":     "NODE",
						"remind_type":     "FINISHED",
						"alert_unit":      "OWNER",
						"useflag":         "true",
						"max_alert_times": "3",
						"alert_interval":  "300",
					}),
				),
			},
			{
				Config: testAccAlicloudDataWorksRemindBasicConfigUpdate(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"remind_name":     fmt.Sprintf("%s-updated", name),
						"remind_unit":     "BASELINE",
						"remind_type":     "ERROR",
						"alert_unit":      "OTHER",
						"useflag":         "false",
						"max_alert_times": "5",
						"alert_interval":  "600",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"useflag"},
			},
		},
	})
}

var AlicloudDataWorksRemindMap0 = map[string]string{
	"remind_id":       CHECKSET,
	"remind_name":     "",
	"remind_unit":     "",
	"remind_type":     "",
	"alert_unit":      "",
	"dnd_end":         "",
	"alert_methods":   CHECKSET,
	"nodes":           CHECKSET,
	"baselines":       CHECKSET,
	"alert_targets":   CHECKSET,
	"useflag":         "",
	"biz_processes":   CHECKSET,
	"max_alert_times": "",
	"alert_interval":  "",
	"detail":          "",
	"robots":          CHECKSET,
	"webhooks":        CHECKSET,
	"projects":        CHECKSET,
	"founder":         NOSET,
	"dnd_start":       NOSET,
	"region_id":       NOSET,
}

func testAccAlicloudDataWorksRemindBasicConfig(name string) string {
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
  detail          = "test remind rule created by terraform"
  webhooks        = ["https://oapi.dingtalk.com/robot/send?access_token=tf-test-token"]
  nodes           = []
  baselines       = []
  biz_processes   = []
  robots          = []
  projects        = []
}
`, name)
}

func testAccAlicloudDataWorksRemindBasicConfigUpdate(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alicloud_data_works_remind" "default" {
  remind_name     = "${var.name}-updated"
  remind_unit     = "BASELINE"
  remind_type     = "ERROR"
  alert_unit      = "OTHER"
  dnd_end         = "23:00"
  alert_methods   = ["MAIL"]
  alert_targets   = ["testuser2@example.com"]
  useflag         = false
  max_alert_times = 5
  alert_interval  = 600
  detail          = "updated remind rule by terraform"
  webhooks        = ["https://oapi.dingtalk.com/robot/send?access_token=tf-test-token-2"]
  nodes           = []
  baselines       = []
  biz_processes   = []
  robots          = []
  projects        = []
}
`, name)
}
