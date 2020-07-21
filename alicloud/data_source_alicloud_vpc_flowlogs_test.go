package alicloud

import (
	"strings"
	"testing"

	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform/helper/acctest"

	"fmt"
)

func TestAccAlicloudVpcFlowlogsDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 99999999)
	descriptionConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpcFlowlogsDataSourceConfig(rand, map[string]string{
			"description": `"${alicloud_vpc_flowlog.default.description}"`,
		}),
		fakeConfig: testAccCheckAlicloudVpcFlowlogsDataSourceConfig(rand, map[string]string{
			"description": `"${alicloud_vpc_flowlog.default.description}-fake"`,
		}),
	}
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpcFlowlogsDataSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_vpc_flowlog.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudVpcFlowlogsDataSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_vpc_flowlog.default.id}-fake"]`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpcFlowlogsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_vpc_flowlog.default.flow_log_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudVpcFlowlogsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_vpc_flowlog.default.flow_log_name}-fake"`,
		}),
	}
	logStoreNameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpcFlowlogsDataSourceConfig(rand, map[string]string{
			"log_store_name": `"${alicloud_vpc_flowlog.default.log_store_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudVpcFlowlogsDataSourceConfig(rand, map[string]string{
			"log_store_name": `"${alicloud_vpc_flowlog.default.log_store_name}-fake"`,
		}),
	}
	projectNameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpcFlowlogsDataSourceConfig(rand, map[string]string{
			"project_name": `"${alicloud_vpc_flowlog.default.project_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudVpcFlowlogsDataSourceConfig(rand, map[string]string{
			"project_name": `"${alicloud_vpc_flowlog.default.project_name}-fake"`,
		}),
	}
	resourceIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpcFlowlogsDataSourceConfig(rand, map[string]string{
			"resource_id": `"${alicloud_vpc_flowlog.default.resource_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudVpcFlowlogsDataSourceConfig(rand, map[string]string{
			"resource_id": `"${alicloud_vpc_flowlog.default.resource_id}-fake"`,
		}),
	}
	resourceTypeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpcFlowlogsDataSourceConfig(rand, map[string]string{
			"resource_type": `"VPC"`,
		}),
		fakeConfig: testAccCheckAlicloudVpcFlowlogsDataSourceConfig(rand, map[string]string{
			"resource_type": `"VSwitch"`,
		}),
	}
	trafficTypeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpcFlowlogsDataSourceConfig(rand, map[string]string{
			"traffic_type": `"All"`,
		}),
		fakeConfig: testAccCheckAlicloudVpcFlowlogsDataSourceConfig(rand, map[string]string{
			"traffic_type": `"Allow"`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpcFlowlogsDataSourceConfig(rand, map[string]string{
			"status":     `"Active"`,
			"name_regex": `"${alicloud_vpc_flowlog.default.flow_log_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudVpcFlowlogsDataSourceConfig(rand, map[string]string{
			"status":     `"Inactive"`,
			"name_regex": `"${alicloud_vpc_flowlog.default.flow_log_name}"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpcFlowlogsDataSourceConfig(rand, map[string]string{
			"description":    `"${alicloud_vpc_flowlog.default.description}"`,
			"name_regex":     `"${alicloud_vpc_flowlog.default.flow_log_name}"`,
			"ids":            `["${alicloud_vpc_flowlog.default.id}"]`,
			"log_store_name": `"${alicloud_vpc_flowlog.default.log_store_name}"`,
			"project_name":   `"${alicloud_vpc_flowlog.default.project_name}"`,
			"resource_id": 	  `"${alicloud_vpc_flowlog.default.resource_id}"`,
			"resource_type":  `"VPC"`,
			"traffic_type":   `"All"`,
			"status":         `"Active"`,
		}),
		fakeConfig: testAccCheckAlicloudVpcFlowlogsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_vpc_flowlog.default.flow_log_name}-fake"`,
			"ids":        `["${alicloud_vpc_flowlog.default.id}"]`,
		}),
	}
	preCheck := func() {
		testAccPreCheckWithRegions(t, true, connectivity.VpcFlowLogNoSkipRegions)
	}
	vpcFlowlogsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, descriptionConf, idsConf, nameRegexConf, logStoreNameConf, projectNameConf, resourceIdConf, resourceTypeConf, trafficTypeConf, statusConf, allConf)
}

func testAccCheckAlicloudVpcFlowlogsDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
  default = "tf-testAcc%sVpcFlowlogsDataSource-%d"
}
resource "alicloud_vpc" "default" {
  cidr_block = "192.168.0.0/24"
  name = "${var.name}"
}
resource "alicloud_log_project" "default"{
  name = "${lower(var.name)}"
  description = "create by terraform"
}	
resource "alicloud_log_store" "default"{
  project = "${alicloud_log_project.default.name}"
  name = "${lower(var.name)}"
  retention_period = 3650
  shard_count = 3
  auto_split = true
  max_split_shard_count = 60
  append_meta = true
}		
resource "alicloud_vpc_flowlog" "default" {
  resource_id = "${alicloud_vpc.default.id}"
  resource_type = "VPC"
  traffic_type = "All"
  log_store_name = "${alicloud_log_store.default.name}"
  project_name = "${alicloud_log_project.default.name}"
  flow_log_name = "${var.name}"
  description = "${var.name}"
}

data "alicloud_vpc_flowlogs" "default" {
	%s
}
`, defaultRegionToTest, rand, strings.Join(pairs, "\n  "))
	return config
}

var existVpcFlowlogsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":                     "1",
		"flow_logs.#":                "1",
		"flow_logs.0.description":    CHECKSET,
		"flow_logs.0.id":             CHECKSET,
		"flow_logs.0.flow_log_name":  CHECKSET,
		"flow_logs.0.log_store_name": CHECKSET,
		"flow_logs.0.project_name":   CHECKSET,
		"flow_logs.0.region_id":   CHECKSET,
		"flow_logs.0.creation_time":   CHECKSET,
	}
}

var fakeVpcFlowlogsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":      "0",
		"flow_logs.#": "0",
	}
}

var vpcFlowlogsCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_vpc_flowlogs.default",
	existMapFunc: existVpcFlowlogsMapFunc,
	fakeMapFunc:  fakeVpcFlowlogsMapFunc,
}
