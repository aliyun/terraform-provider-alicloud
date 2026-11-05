package cloudmonitorservice_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudCloudMonitorServiceSiteMonitorDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(1000000, 9999999)

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudCloudMonitorServiceSiteMonitorSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_cms_site_monitor.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudCloudMonitorServiceSiteMonitorSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_cms_site_monitor.default.id}_fake"]`,
		}),
	}

	TaskTypeConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudCloudMonitorServiceSiteMonitorSourceConfig(rand, map[string]string{
			"ids":       `["${alicloud_cms_site_monitor.default.id}"]`,
			"task_type": `"HTTP"`,
		}),
		FakeConfig: testAccCheckAlicloudCloudMonitorServiceSiteMonitorSourceConfig(rand, map[string]string{
			"ids":       `["${alicloud_cms_site_monitor.default.id}_fake"]`,
			"task_type": `"HTTP"`,
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudCloudMonitorServiceSiteMonitorSourceConfig(rand, map[string]string{
			"ids":       `["${alicloud_cms_site_monitor.default.id}"]`,
			"task_type": `"HTTP"`,
		}),
		FakeConfig: testAccCheckAlicloudCloudMonitorServiceSiteMonitorSourceConfig(rand, map[string]string{
			"ids":       `["${alicloud_cms_site_monitor.default.id}_fake"]`,
			"task_type": `"HTTP"`,
		}),
	}

	CloudMonitorServiceSiteMonitorCheckInfo.DataSourceTestCheck(t, rand, idsConf, TaskTypeConf, allConf)
}

var existCloudMonitorServiceSiteMonitorMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"monitors.#":         "1",
		"monitors.0.task_id": acctest.CHECKSET,
	}
}

var fakeCloudMonitorServiceSiteMonitorMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"monitors.#": "0",
	}
}

var CloudMonitorServiceSiteMonitorCheckInfo = acctest.DataSourceAttr{
	ResourceId:   "data.alicloud_cms_site_monitors.default",
	ExistMapFunc: existCloudMonitorServiceSiteMonitorMapFunc,
	FakeMapFunc:  fakeCloudMonitorServiceSiteMonitorMapFunc,
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
