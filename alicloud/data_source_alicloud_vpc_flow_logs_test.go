package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudVPCFlowLogsDataSource(t *testing.T) {
	rand := acctest.RandInt()
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpcFlowLogsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_vpc_flow_log.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudVpcFlowLogsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_vpc_flow_log.default.id}_fake"]`,
		}),
	}
	descriptionConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpcFlowLogsDataSourceName(rand, map[string]string{
			"ids":         `["${alicloud_vpc_flow_log.default.id}"]`,
			"description": `"${alicloud_vpc_flow_log.default.description}"`,
		}),
		fakeConfig: testAccCheckAlicloudVpcFlowLogsDataSourceName(rand, map[string]string{
			"ids":         `["${alicloud_vpc_flow_log.default.id}_fake"]`,
			"description": `"${alicloud_vpc_flow_log.default.description}_fake"`,
		}),
	}
	flowLogNameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpcFlowLogsDataSourceName(rand, map[string]string{
			"ids":           `["${alicloud_vpc_flow_log.default.id}"]`,
			"flow_log_name": `"${alicloud_vpc_flow_log.default.flow_log_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudVpcFlowLogsDataSourceName(rand, map[string]string{
			"ids":           `["${alicloud_vpc_flow_log.default.id}_fake"]`,
			"flow_log_name": `"${alicloud_vpc_flow_log.default.flow_log_name}_fake"`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpcFlowLogsDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_vpc_flow_log.default.id}"]`,
			"name_regex": `"${alicloud_vpc_flow_log.default.flow_log_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudVpcFlowLogsDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_vpc_flow_log.default.id}_fake"]`,
			"name_regex": `"${alicloud_vpc_flow_log.default.flow_log_name}_fake"`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpcFlowLogsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_vpc_flow_log.default.id}"]`,
			"status": `"Active"`,
		}),
		fakeConfig: testAccCheckAlicloudVpcFlowLogsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_vpc_flow_log.default.id}_fake"]`,
			"status": `"Inactive"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpcFlowLogsDataSourceName(rand, map[string]string{
			"description":   `"${alicloud_vpc_flow_log.default.description}"`,
			"flow_log_name": `"${alicloud_vpc_flow_log.default.flow_log_name}"`,
			"ids":           `["${alicloud_vpc_flow_log.default.id}"]`,
			"name_regex":    `"${alicloud_vpc_flow_log.default.flow_log_name}"`,
			"status":        `"Active"`,
		}),
		fakeConfig: testAccCheckAlicloudVpcFlowLogsDataSourceName(rand, map[string]string{
			"description":   `"${alicloud_vpc_flow_log.default.description}_fake"`,
			"flow_log_name": `"${alicloud_vpc_flow_log.default.flow_log_name}_fake"`,
			"ids":           `["${alicloud_vpc_flow_log.default.id}_fake"]`,
			"name_regex":    `"${alicloud_vpc_flow_log.default.flow_log_name}_fake"`,
			"status":        `"Inactive"`,
		}),
	}
	var existAlicloudVpcFlowLogsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                 "1",
			"names.#":               "1",
			"logs.#":                "1",
			"logs.0.description":    fmt.Sprintf("tf-testAccFlowLog-%d", rand),
			"logs.0.flow_log_name":  fmt.Sprintf("tf-testAccFlowLog-%d", rand),
			"logs.0.log_store_name": `vpc-flow-log-for-vpc`,
			"logs.0.project_name":   `vpc-flow-log-for-vpc`,
			"logs.0.resource_id":    CHECKSET,
			"logs.0.resource_type":  `VPC`,
			"logs.0.traffic_type":   `All`,
			"logs.0.status":         `Active`,
		}
	}
	var fakeAlicloudVpcFlowLogsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudVpcFlowLogsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_vpc_flow_logs.default",
		existMapFunc: existAlicloudVpcFlowLogsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudVpcFlowLogsDataSourceNameMapFunc,
	}
	alicloudVpcFlowLogsCheckInfo.dataSourceTestCheck(t, rand, idsConf, descriptionConf, flowLogNameConf, nameRegexConf, statusConf, allConf)
}
func testAccCheckAlicloudVpcFlowLogsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAccFlowLog-%d"
}

resource "alicloud_vpc" "default" {
  cidr_block = "192.168.0.0/24"
  name       = var.name
}

resource "alicloud_vpc_flow_log" "default" {
	description = var.name
	flow_log_name = var.name
	log_store_name = "vpc-flow-log-for-vpc" 
	project_name = "vpc-flow-log-for-vpc"
	resource_id = "${alicloud_vpc.default.id}"
	resource_type = "VPC"
	traffic_type = "All"
	status = "Active"
}

data "alicloud_vpc_flow_logs" "default" {	
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
