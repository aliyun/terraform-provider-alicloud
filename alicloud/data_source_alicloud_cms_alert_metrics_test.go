// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAliCloudCmsAlertMetricsDataSource(t *testing.T) {
	testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	rand := acctest.RandIntRange(1000000, 9999999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudCmsAlertMetricsSourceConfig(rand, map[string]string{
			"ids": `["${data.alicloud_cms_alert_metrics.all.metrics.0.id}"]`,
		}),
		fakeConfig: testAccCheckAliCloudCmsAlertMetricsSourceConfig(rand, map[string]string{
			"ids": `["${data.alicloud_cms_alert_metrics.all.metrics.0.id}_fake"]`,
		}),
	}

	groupConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudCmsAlertMetricsSourceConfig(rand, map[string]string{
			"ids":   `["${data.alicloud_cms_alert_metrics.all.metrics.0.id}"]`,
			"group": `"${data.alicloud_cms_alert_metrics.all.metrics.0.group}"`,
		}),
		fakeConfig: testAccCheckAliCloudCmsAlertMetricsSourceConfig(rand, map[string]string{
			"ids":   `["${data.alicloud_cms_alert_metrics.all.metrics.0.id}_fake"]`,
			"group": `"${data.alicloud_cms_alert_metrics.all.metrics.0.group}"`,
		}),
	}

	includeDetailsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudCmsAlertMetricsSourceConfig(rand, map[string]string{
			"ids":             `["${data.alicloud_cms_alert_metrics.all.metrics.0.id}"]`,
			"include_details": `"true"`,
		}),
		fakeConfig: testAccCheckAliCloudCmsAlertMetricsSourceConfig(rand, map[string]string{
			"ids":             `["${data.alicloud_cms_alert_metrics.all.metrics.0.id}_fake"]`,
			"include_details": `"true"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudCmsAlertMetricsSourceConfig(rand, map[string]string{
			"ids":             `["${data.alicloud_cms_alert_metrics.all.metrics.0.id}"]`,
			"group":           `"${data.alicloud_cms_alert_metrics.all.metrics.0.group}"`,
			"include_details": `"true"`,
		}),
		fakeConfig: testAccCheckAliCloudCmsAlertMetricsSourceConfig(rand, map[string]string{
			"ids":             `["${data.alicloud_cms_alert_metrics.all.metrics.0.id}_fake"]`,
			"group":           `"${data.alicloud_cms_alert_metrics.all.metrics.0.group}_fake"`,
			"include_details": `"true"`,
		}),
	}

	CmsAlertMetricsCheckInfo.dataSourceTestCheck(t, rand, idsConf, groupConf, includeDetailsConf, allConf)
}

var existCmsAlertMetricsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"metrics.#":                 "1",
		"metrics.0.id":              CHECKSET,
		"metrics.0.alert_metric_id": CHECKSET,
		"ids.#":                     "1",
	}
}

var fakeCmsAlertMetricsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"metrics.#": "0",
	}
}

var CmsAlertMetricsCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_cms_alert_metrics.default",
	existMapFunc: existCmsAlertMetricsMapFunc,
	fakeMapFunc:  fakeCmsAlertMetricsMapFunc,
}

func testAccCheckAliCloudCmsAlertMetricsSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
data "alicloud_cms_alert_metrics" "all" {
}

data "alicloud_cms_alert_metrics" "default" {
%s
}
`, strings.Join(pairs, "\n   "))
	return config
}
