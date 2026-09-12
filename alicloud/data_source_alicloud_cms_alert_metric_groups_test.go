package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

// Test Cms AlertMetricGroups data source. >>> Data source test cases.

// TestAccAliCloudCmsAlertMetricGroupsDataSource exercises every Optional filter
// of the data source: ids, datasource_type, include_details and output_file.
//
// Alert metric groups are predefined, read-only system entries shared by all
// accounts, so the test does not create any resource: an unfiltered companion
// data source ("all") provides a real alert metric group id that each config
// then filters on, keeping the exist assertions deterministic.
func TestAccAliCloudCmsAlertMetricGroupsDataSource(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCmsAlertMetricGroupsSourceConfig(rand, map[string]string{
			"ids": `[data.alicloud_cms_alert_metric_groups.all.groups.0.id]`,
		}),
		fakeConfig: testAccCheckAlicloudCmsAlertMetricGroupsSourceConfig(rand, map[string]string{
			"ids": `["${data.alicloud_cms_alert_metric_groups.all.groups.0.id}_fake"]`,
		}),
	}

	datasourceTypeConf := dataSourceTestAccConfig{
		// datasource_type is a server-side filter; the value is derived from the
		// first group's own datasource_types (a JSON array string such as
		// ["arms_metrics"]) so the pinned id always matches the filter.
		existConfig: testAccCheckAlicloudCmsAlertMetricGroupsSourceConfig(rand, map[string]string{
			"ids":             `[data.alicloud_cms_alert_metric_groups.all.groups.0.id]`,
			"datasource_type": `jsondecode(data.alicloud_cms_alert_metric_groups.all.groups.0.datasource_types)[0]`,
		}),
		fakeConfig: testAccCheckAlicloudCmsAlertMetricGroupsSourceConfig(rand, map[string]string{
			"ids":             `["${data.alicloud_cms_alert_metric_groups.all.groups.0.id}_fake"]`,
			"datasource_type": `"nonexistent_datasource_type"`,
		}),
	}

	includeDetailsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCmsAlertMetricGroupsSourceConfig(rand, map[string]string{
			"ids":             `[data.alicloud_cms_alert_metric_groups.all.groups.0.id]`,
			"include_details": `true`,
		}),
		fakeConfig: testAccCheckAlicloudCmsAlertMetricGroupsSourceConfig(rand, map[string]string{
			"ids":             `["${data.alicloud_cms_alert_metric_groups.all.groups.0.id}_fake"]`,
			"include_details": `true`,
		}),
	}

	outputFileConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCmsAlertMetricGroupsSourceConfig(rand, map[string]string{
			"ids":         `[data.alicloud_cms_alert_metric_groups.all.groups.0.id]`,
			"output_file": `"/tmp/cms_alert_metric_groups_ds_output.txt"`,
		}),
		fakeConfig: testAccCheckAlicloudCmsAlertMetricGroupsSourceConfig(rand, map[string]string{
			"ids":         `["${data.alicloud_cms_alert_metric_groups.all.groups.0.id}_fake"]`,
			"output_file": `"/tmp/cms_alert_metric_groups_ds_output_fake.txt"`,
		}),
	}

	CmsAlertMetricGroupsCheckInfo.dataSourceTestCheck(t, rand, idsConf, datasourceTypeConf, includeDetailsConf, outputFileConf)
}

var existCmsAlertMetricGroupsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"groups.#":                       "1",
		"ids.#":                          "1",
		"groups.0.id":                    CHECKSET,
		"groups.0.alert_metric_group_id": CHECKSET,
		"groups.0.datasource_types":      CHECKSET,
		"groups.0.display_name_cn":       CHECKSET,
		"groups.0.display_name_en":       CHECKSET,
		// filters/params/order_index/description_* are sparse predefined-system
		// fields (only populated for some groups, or only when include_details is
		// true), so they are intentionally not asserted.
	}
}

var fakeCmsAlertMetricGroupsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"groups.#": "0",
		"ids.#":    "0",
	}
}

var CmsAlertMetricGroupsCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_cms_alert_metric_groups.default",
	existMapFunc: existCmsAlertMetricGroupsMapFunc,
	fakeMapFunc:  fakeCmsAlertMetricGroupsMapFunc,
}

func testAccCheckAlicloudCmsAlertMetricGroupsSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	return fmt.Sprintf(`
data "alicloud_cms_alert_metric_groups" "all" {
}

data "alicloud_cms_alert_metric_groups" "default" {
  %s
}
`, strings.Join(pairs, "\n  "))
}

// Test Cms AlertMetricGroups data source. <<< Data source test cases.
