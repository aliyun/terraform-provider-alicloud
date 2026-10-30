package smartag_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudSmartagFlowLogsDataSource(t *testing.T) {
	rand := tfacctest.RandInt()
	acctest.CheckoutSupportedRegions(t, true, connectivity.SmartagSupportedRegions)
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudSmartagFlowLogsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_smartag_flow_log.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudSmartagFlowLogsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_smartag_flow_log.default.id}_fake"]`,
		}),
	}
	descriptionConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudSmartagFlowLogsDataSourceName(rand, map[string]string{
			"ids":         `["${alicloud_smartag_flow_log.default.id}"]`,
			"description": `"${alicloud_smartag_flow_log.default.description}"`,
		}),
		FakeConfig: testAccCheckAlicloudSmartagFlowLogsDataSourceName(rand, map[string]string{
			"ids":         `["${alicloud_smartag_flow_log.default.id}"]`,
			"description": `"${alicloud_smartag_flow_log.default.description}_fake"`,
		}),
	}
	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudSmartagFlowLogsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_smartag_flow_log.default.flow_log_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudSmartagFlowLogsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_smartag_flow_log.default.flow_log_name}_fake"`,
		}),
	}
	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudSmartagFlowLogsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_smartag_flow_log.default.id}"]`,
			"status": `"Active"`,
		}),
		FakeConfig: testAccCheckAlicloudSmartagFlowLogsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_smartag_flow_log.default.id}"]`,
			"status": `"Inactive"`,
		}),
	}
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudSmartagFlowLogsDataSourceName(rand, map[string]string{
			"description": `"${alicloud_smartag_flow_log.default.description}"`,
			"ids":         `["${alicloud_smartag_flow_log.default.id}"]`,
			"name_regex":  `"${alicloud_smartag_flow_log.default.flow_log_name}"`,
			"status":      `"Active"`,
		}),
		FakeConfig: testAccCheckAlicloudSmartagFlowLogsDataSourceName(rand, map[string]string{
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
			"logs.0.id":                  acctest.CHECKSET,
			"logs.0.flow_log_id":         acctest.CHECKSET,
			"logs.0.total_sag_num":       acctest.CHECKSET,
			"logs.0.status":              "Active",
			"logs.0.resource_group_id":   acctest.CHECKSET,
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
			"logs.0.sls_region_id":       acctest.DefaultRegionToTest,
		}
	}
	var fakeAlicloudSmartagFlowLogsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudSmartagFlowLogsCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_smartag_flow_logs.default",
		ExistMapFunc: existAlicloudSmartagFlowLogsDataSourceNameMapFunc,
		FakeMapFunc:  fakeAlicloudSmartagFlowLogsDataSourceNameMapFunc,
	}

	preCheck := func() {
		acctest.PreCheck(t)
	}
	alicloudSmartagFlowLogsCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, descriptionConf, nameRegexConf, statusConf, allConf)
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
`, rand, acctest.DefaultRegionToTest, strings.Join(pairs, " \n "))
	return config
}
