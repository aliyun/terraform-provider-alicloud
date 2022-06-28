package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudSmartagFlowLogsDataSource(t *testing.T) {
	rand := acctest.RandInt()
	checkoutSupportedRegions(t, true, connectivity.SmartagSupportedRegions)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSmartagFlowLogsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_smartag_flow_log.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudSmartagFlowLogsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_smartag_flow_log.default.id}_fake"]`,
		}),
	}
	descriptionConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSmartagFlowLogsDataSourceName(rand, map[string]string{
			"ids":         `["${alicloud_smartag_flow_log.default.id}"]`,
			"description": `"${alicloud_smartag_flow_log.default.description}"`,
		}),
		fakeConfig: testAccCheckAlicloudSmartagFlowLogsDataSourceName(rand, map[string]string{
			"ids":         `["${alicloud_smartag_flow_log.default.id}"]`,
			"description": `"${alicloud_smartag_flow_log.default.description}_fake"`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSmartagFlowLogsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_smartag_flow_log.default.flow_log_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudSmartagFlowLogsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_smartag_flow_log.default.flow_log_name}_fake"`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSmartagFlowLogsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_smartag_flow_log.default.id}"]`,
			"status": `"Active"`,
		}),
		fakeConfig: testAccCheckAlicloudSmartagFlowLogsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_smartag_flow_log.default.id}"]`,
			"status": `"Inactive"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSmartagFlowLogsDataSourceName(rand, map[string]string{
			"description": `"${alicloud_smartag_flow_log.default.description}"`,
			"ids":         `["${alicloud_smartag_flow_log.default.id}"]`,
			"name_regex":  `"${alicloud_smartag_flow_log.default.flow_log_name}"`,
			"status":      `"Active"`,
		}),
		fakeConfig: testAccCheckAlicloudSmartagFlowLogsDataSourceName(rand, map[string]string{
			"description": `"${alicloud_smartag_flow_log.default.description}_fake"`,
			"ids":         `["${alicloud_smartag_flow_log.default.id}_fake"]`,
			"name_regex":  `"${alicloud_smartag_flow_log.default.flow_log_name}_fake"`,
			"status":      `"Inactive"`,
		}),
	}
	var existAlicloudSmartagFlowLogsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                      "1",
			"names.#":                    "1",
			"logs.#":                     "1",
			"logs.0.id":                  CHECKSET,
			"logs.0.flow_log_id":         CHECKSET,
			"logs.0.total_sag_num":       CHECKSET,
			"logs.0.status":              "Active",
			"logs.0.resource_group_id":   CHECKSET,
			"logs.0.active_aging":        "300",
			"logs.0.description":         fmt.Sprintf("tf-testAccFlowLog-%d", rand),
			"logs.0.flow_log_name":       fmt.Sprintf("tf-testAccFlowLog-%d", rand),
			"logs.0.inactive_aging":      "15",
			"logs.0.logstore_name":       fmt.Sprintf("tf-testAccFlowLog-%d", rand),
			"logs.0.netflow_server_ip":   "192.168.0.2",
			"logs.0.netflow_server_port": "9995",
			"logs.0.netflow_version":     "V9",
			"logs.0.output_type":         "all",
			"logs.0.project_name":        fmt.Sprintf("tf-testAccFlowLog-%d", rand),
			"logs.0.sls_region_id":       defaultRegionToTest,
		}
	}
	var fakeAlicloudSmartagFlowLogsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudSmartagFlowLogsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_smartag_flow_logs.default",
		existMapFunc: existAlicloudSmartagFlowLogsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudSmartagFlowLogsDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}
	alicloudSmartagFlowLogsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, descriptionConf, nameRegexConf, statusConf, allConf)
}
func testAccCheckAlicloudSmartagFlowLogsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAccFlowLog-%d"
}

resource "alicloud_smartag_flow_log" "default" {
	netflow_server_port = 9995
	logstore_name =       var.name
	description =         var.name
	active_aging =        300
	project_name =        var.name
	netflow_server_ip =   "192.168.0.2"
	netflow_version =     "V9"
	inactive_aging =      15
	flow_log_name =       var.name
	sls_region_id =       "%s"
	output_type =         "all"
}

data "alicloud_smartag_flow_logs" "default" {	
	%s
}
`, rand, strings.Join(pairs, " \n "), defaultRegionToTest)
	return config
}
