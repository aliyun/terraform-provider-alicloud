package alicloud

import (
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"fmt"
)

func TestAccAlicloudCenFlowlogsDataSource(t *testing.T) {
	// flow log has been offline
	// t.Skip("From January 30, 2022, the cloud enterprise network will take the old console flow log function offline. If you need to continue to use the flow log function, you can enter the new version console to use the flow log function of the enterprise version forwarding router. The Enterprise Edition Forwarding Router Flow Log feature provides the same capabilities as the Legacy Console Flow Log feature")
	rand := acctest.RandIntRange(1000000, 99999999)
	cenIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenFlowlogsDataSourceConfig(rand, map[string]string{
			"cen_id": `"${alicloud_cen_flowlog.default.cen_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudCenFlowlogsDataSourceConfig(rand, map[string]string{
			"cen_id": `"${alicloud_cen_flowlog.default.cen_id}-fake"`,
		}),
	}
	descriptionConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenFlowlogsDataSourceConfig(rand, map[string]string{
			"description": `"${alicloud_cen_flowlog.default.description}"`,
		}),
		fakeConfig: testAccCheckAlicloudCenFlowlogsDataSourceConfig(rand, map[string]string{
			"description": `"${alicloud_cen_flowlog.default.description}-fake"`,
		}),
	}
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenFlowlogsDataSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_cen_flowlog.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudCenFlowlogsDataSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_cen_flowlog.default.id}-fake"]`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenFlowlogsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_cen_flowlog.default.flow_log_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudCenFlowlogsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_cen_flowlog.default.flow_log_name}-fake"`,
		}),
	}
	logStoreNameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenFlowlogsDataSourceConfig(rand, map[string]string{
			"log_store_name": `"${alicloud_cen_flowlog.default.log_store_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudCenFlowlogsDataSourceConfig(rand, map[string]string{
			"log_store_name": `"${alicloud_cen_flowlog.default.log_store_name}-fake"`,
		}),
	}
	projectNameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenFlowlogsDataSourceConfig(rand, map[string]string{
			"project_name": `"${alicloud_cen_flowlog.default.project_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudCenFlowlogsDataSourceConfig(rand, map[string]string{
			"project_name": `"${alicloud_cen_flowlog.default.project_name}-fake"`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenFlowlogsDataSourceConfig(rand, map[string]string{
			"status":     `"Active"`,
			"name_regex": `"${alicloud_cen_flowlog.default.flow_log_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudCenFlowlogsDataSourceConfig(rand, map[string]string{
			"status":     `"Inactive"`,
			"name_regex": `"${alicloud_cen_flowlog.default.flow_log_name}"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenFlowlogsDataSourceConfig(rand, map[string]string{
			"cen_id":         `"${alicloud_cen_flowlog.default.cen_id}"`,
			"description":    `"${alicloud_cen_flowlog.default.description}"`,
			"name_regex":     `"${alicloud_cen_flowlog.default.flow_log_name}"`,
			"ids":            `["${alicloud_cen_flowlog.default.id}"]`,
			"log_store_name": `"${alicloud_cen_flowlog.default.log_store_name}"`,
			"project_name":   `"${alicloud_cen_flowlog.default.project_name}"`,
			"status":         `"Active"`,
		}),
		fakeConfig: testAccCheckAlicloudCenFlowlogsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_cen_flowlog.default.flow_log_name}-fake"`,
			"ids":        `["${alicloud_cen_flowlog.default.id}"]`,
		}),
	}
	preCheck := func() {
		testAccPreCheckWithAccountSiteType(t, DomesticSite)
		testAccPreCheckWithRegions(t, true, connectivity.CenNoSkipRegions)
	}
	cenFlowlogsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, cenIdConf, descriptionConf, idsConf, nameRegexConf, logStoreNameConf, projectNameConf, statusConf, allConf)
}

func testAccCheckAlicloudCenFlowlogsDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	  default = "tf-testacc%s-%d"
	}

resource "alicloud_cen_instance" "defaultc5kxyC" {
  cen_instance_name = var.name
}

resource "alicloud_cen_transit_router" "defaultVw2U9u" {
  cen_id = alicloud_cen_instance.defaultc5kxyC.id
}

resource "alicloud_log_project" "default" {
  project_name = var.name
  description  = "terraform-example"
}

resource "alicloud_log_store" "default" {
  project_name          = alicloud_log_project.default.project_name
  logstore_name         = var.name
  shard_count           = 3
  auto_split            = true
  max_split_shard_count = 60
  append_meta           = true
}

resource "alicloud_cen_flowlog" "default" {
  project_name = "${alicloud_log_store.default.project_name}"
  flow_log_name = "tf-testacceu-central-1cenflowlog41302_update"
  log_format_string = "$${srcaddr}$${dstaddr}$${bytes}"
  cen_id = "${alicloud_cen_instance.defaultc5kxyC.id}"
  log_store_name = "${alicloud_log_store.default.logstore_name}"
  interval = "600"
  status = "Active"
  transit_router_id = "${alicloud_cen_transit_router.defaultVw2U9u.transit_router_id}"
  description = "flowlog-resource-test-1"
}

data "alicloud_cen_flowlogs" "default" {
	%s
}
`, defaultRegionToTest, rand, strings.Join(pairs, "\n  "))
	return config
}

var existCenFlowlogsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":                     "1",
		"flowlogs.#":                "1",
		"flowlogs.0.cen_id":         CHECKSET,
		"flowlogs.0.description":    CHECKSET,
		"flowlogs.0.id":             CHECKSET,
		"flowlogs.0.flow_log_id":    CHECKSET,
		"flowlogs.0.flow_log_name":  CHECKSET,
		"flowlogs.0.log_store_name": CHECKSET,
		"flowlogs.0.project_name":   CHECKSET,
		"flowlogs.0.status":         "Active",
	}
}

var fakeCenFlowlogsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":      "0",
		"flowlogs.#": "0",
	}
}

var cenFlowlogsCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_cen_flowlogs.default",
	existMapFunc: existCenFlowlogsMapFunc,
	fakeMapFunc:  fakeCenFlowlogsMapFunc,
}
