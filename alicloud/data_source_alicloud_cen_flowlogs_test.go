package alicloud

import (
	"strings"
	"testing"

	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform/helper/acctest"

	"fmt"
)

func TestAccAlicloudCenFlowlogsDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 99999999)
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenFlowlogsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_cen_flowlog.default.flow_log_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudCenFlowlogsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_cen_flowlog.default.flow_log_name}-fake"`,
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
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenFlowlogsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_cen_flowlog.default.flow_log_name}"`,
			"ids":        `["${alicloud_cen_flowlog.default.id}"]`,
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
	cenFlowlogsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, nameRegexConf, idsConf, allConf)
}

func testAccCheckAlicloudCenFlowlogsDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	  default = "tf-testAcc%sCenFlowlogsDataSource-%d"
	}
resource "alicloud_cen_instance" "default" {
	name = "${var.name}"
	description = "tf-testAccCenConfigDescription"
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
resource "alicloud_cen_flowlog" "default" {
	cen_id = "${alicloud_cen_instance.default.id}"
	project_name = "${alicloud_log_project.default.name}"
	log_store_name = "${alicloud_log_store.default.name}"
	flow_log_name = "${var.name}"
	description = "${var.name}"
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
