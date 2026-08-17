package alicloud

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

// TestAccAlicloudCmsMetricListDataSource verifies the data source that wraps
// DescribeMetricList (Cms/2019-01-01). It queries the read-only
// acs_ecs_dashboard CPUUtilization metric within the last 1-2 hours so the
// test does not depend on any provisioned resource. The "exist" config uses
// a real metric name (expecting a non-empty response when the account has
// ECS instances, otherwise still a valid empty list), while the "fake" config
// queries a non-existent metric name and expects an empty datapoints list.
func TestAccAlicloudCmsMetricListDataSource(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)

	basicConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCmsMetricListDataSourceName(rand, map[string]string{
			"metric_name":      `"CPUUtilization"`,
			"metric_namespace": `"acs_ecs_dashboard"`,
			"period":           `"60"`,
		}),
		fakeConfig: testAccCheckAlicloudCmsMetricListDataSourceName(rand, map[string]string{
			"metric_name":      `"FakeMetricNotExist"`,
			"metric_namespace": `"acs_ecs_dashboard"`,
			"period":           `"60"`,
		}),
	}

	periodConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCmsMetricListDataSourceName(rand, map[string]string{
			"metric_name":      `"CPUUtilization"`,
			"metric_namespace": `"acs_ecs_dashboard"`,
			"period":           `"300"`,
			"length":           "100",
		}),
		fakeConfig: testAccCheckAlicloudCmsMetricListDataSourceName(rand, map[string]string{
			"metric_name":      `"FakeMetricNotExist"`,
			"metric_namespace": `"acs_ecs_dashboard"`,
			"period":           `"300"`,
			"length":           "100",
		}),
	}

	var existAlicloudCmsMetricListDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"datapoints.#":  CHECKSET,
			"actual_period": CHECKSET,
		}
	}
	var fakeAlicloudCmsMetricListDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"datapoints.#": "0",
		}
	}
	var alicloudCmsMetricListCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_cms_metric_list.default",
		existMapFunc: existAlicloudCmsMetricListDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudCmsMetricListDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}
	alicloudCmsMetricListCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, basicConf, periodConf)
}

func testAccCheckAlicloudCmsMetricListDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	// Use a dynamic time window (1h-2h ago to 1h ago) so the test stays valid
	// over time. CMS DescribeMetricList accepts ISO 8601 timestamps.
	now := time.Now()
	start := now.Add(-2 * time.Hour).Format(time.RFC3339)
	end := now.Add(-1 * time.Hour).Format(time.RFC3339)

	config := fmt.Sprintf(`
data "alicloud_cms_metric_list" "default" {
	start_time = "%s"
	end_time = "%s"
	%s
}
`, start, end, strings.Join(pairs, " \n "))
	return config
}
