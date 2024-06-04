package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudCloudMonitorServiceSiteMonitorDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCloudMonitorServiceSiteMonitorSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_cms_site_monitor.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudCloudMonitorServiceSiteMonitorSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_cms_site_monitor.default.id}_fake"]`,
		}),
	}

	TaskTypeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCloudMonitorServiceSiteMonitorSourceConfig(rand, map[string]string{
			"ids":       `["${alicloud_cms_site_monitor.default.id}"]`,
			"task_type": `"HTTP"`,
		}),
		fakeConfig: testAccCheckAlicloudCloudMonitorServiceSiteMonitorSourceConfig(rand, map[string]string{
			"ids":       `["${alicloud_cms_site_monitor.default.id}_fake"]`,
			"task_type": `"HTTP"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCloudMonitorServiceSiteMonitorSourceConfig(rand, map[string]string{
			"ids":       `["${alicloud_cms_site_monitor.default.id}"]`,
			"task_type": `"HTTP"`,
		}),
		fakeConfig: testAccCheckAlicloudCloudMonitorServiceSiteMonitorSourceConfig(rand, map[string]string{
			"ids":       `["${alicloud_cms_site_monitor.default.id}_fake"]`,
			"task_type": `"HTTP"`,
		}),
	}

	CloudMonitorServiceSiteMonitorCheckInfo.dataSourceTestCheck(t, rand, idsConf, TaskTypeConf, allConf)
}

var existCloudMonitorServiceSiteMonitorMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"monitors.#":         "1",
		"monitors.0.task_id": CHECKSET,
	}
}

var fakeCloudMonitorServiceSiteMonitorMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"monitors.#": "0",
	}
}

var CloudMonitorServiceSiteMonitorCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_cms_site_monitors.default",
	existMapFunc: existCloudMonitorServiceSiteMonitorMapFunc,
	fakeMapFunc:  fakeCloudMonitorServiceSiteMonitorMapFunc,
}

func testAccCheckAlicloudCloudMonitorServiceSiteMonitorSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccCloudMonitorServiceSiteMonitor%d"
}

resource "alicloud_cms_site_monitor" "default" {
  address   = "http://www.alibabacloud.com"
  task_name = var.name
  task_type = "HTTP"
  interval  = 5
  isp_cities {
    city = "546"
    isp  = "465"
  }
  options_json = <<EOT
{
    "http_method": "get",
    "waitTime_after_completion": null,
    "ipv6_task": false,
    "diagnosis_ping": false,
    "diagnosis_mtr": false,
    "assertions": [
        {
            "operator": "lessThan",
            "type": "response_time",
            "target": 1000
        }
    ],
    "time_out": 30000
}
EOT
}

data "alicloud_cms_site_monitors" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
