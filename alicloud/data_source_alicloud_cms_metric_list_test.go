package alicloud

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

// TestAccAlicloudCmsMetricListDataSource verifies the data source that wraps
// DescribeMetricList (Cms/2019-01-01). It queries the read-only
// acs_ecs_dashboard CPUUtilization metric so the test does not depend on any
// provisioned resource. The exist configs assert the millisecond timestamp is
// a decimal string (not scientific notation) and that string fields are never
// "<nil>". The refresh step re-applies the identical exist config right after
// the first read ended on the last page (NextToken null): the stored
// next_token must be "" and must not be resent as "<nil>" on re-read. The
// "fake" config queries a non-existent metric name and expects an empty
// datapoints list.
func TestAccAlicloudCmsMetricListDataSource(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)

	// A short window keeps the datapoint count far below the default Length
	// (1000), so the query is guaranteed to end on the last page with a null
	// NextToken. basicConf and refreshConf are exist-only so the framework
	// runs their identical configs on consecutive steps and the second step
	// re-reads the datasource in place with the stored next_token.
	existConfig := testAccCheckAlicloudCmsMetricListDataSourceNameShort(rand, map[string]string{
		"metric_name":      `"CPUUtilization"`,
		"metric_namespace": `"acs_ecs_dashboard"`,
		"period":           `"60"`,
	})
	basicConf := dataSourceTestAccConfig{
		existConfig: existConfig,
		existChangMap: map[string]string{
			"next_token": "",
		},
	}
	refreshConf := dataSourceTestAccConfig{
		existConfig: existConfig,
		existChangMap: map[string]string{
			"next_token": "",
		},
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
			"datapoints.#":             CHECKSET,
			"datapoints.0.timestamp":   REGEXMATCH + "^\\d+$",
			"datapoints.0.metric_name": "CPUUtilization",
			"datapoints.0.uuid":        "",
			"datapoints.0.average":     CHECKSET,
			"datapoints.0.maximum":     CHECKSET,
			"datapoints.0.minimum":     CHECKSET,
			"datapoints.0.user_id":     CHECKSET,
			"datapoints.0.instance_id": CHECKSET,
			"actual_period":            CHECKSET,
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
	alicloudCmsMetricListCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, basicConf, refreshConf, periodConf)
}

func testAccCheckAlicloudCmsMetricListDataSourceName(rand int, attrMap map[string]string) string {
	now := time.Now()
	start := now.Add(-2 * time.Hour).Format(time.RFC3339)
	end := now.Add(-1 * time.Hour).Format(time.RFC3339)
	return testAccCheckAlicloudCmsMetricListDataSourceConfig(start, end, attrMap)
}

func testAccCheckAlicloudCmsMetricListDataSourceNameShort(rand int, attrMap map[string]string) string {
	now := time.Now()
	start := now.Add(-10 * time.Minute).Format(time.RFC3339)
	end := now.Add(-5 * time.Minute).Format(time.RFC3339)
	return testAccCheckAlicloudCmsMetricListDataSourceConfig(start, end, attrMap)
}

func testAccCheckAlicloudCmsMetricListDataSourceConfig(start, end string, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
data "alicloud_cms_metric_list" "default" {
	start_time = "%s"
	end_time = "%s"
	%s
}
`, start, end, strings.Join(pairs, " \n "))
	return config
}
