package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAliCloudSlsMetricStoresDataSource(t *testing.T) {
	testAccPreCheckWithRegions(t, true, connectivity.SlsTestRegions)
	rand := acctest.RandIntRange(10000, 99999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudSlsMetricStoresDataSourceConfig(rand, map[string]string{
			"ids":          `["${alicloud_sls_metric_store.default.id}"]`,
			"project_name": `"${alicloud_log_project.default.project_name}"`,
		}),
		fakeConfig: testAccCheckAliCloudSlsMetricStoresDataSourceConfig(rand, map[string]string{
			"ids":          `["${alicloud_sls_metric_store.default.id}_fake"]`,
			"project_name": `"${alicloud_log_project.default.project_name}"`,
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudSlsMetricStoresDataSourceConfig(rand, map[string]string{
			"name_regex":   `"${alicloud_sls_metric_store.default.metric_store_name}"`,
			"project_name": `"${alicloud_log_project.default.project_name}"`,
		}),
		fakeConfig: testAccCheckAliCloudSlsMetricStoresDataSourceConfig(rand, map[string]string{
			"name_regex":   `"${alicloud_sls_metric_store.default.metric_store_name}_fake"`,
			"project_name": `"${alicloud_log_project.default.project_name}"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudSlsMetricStoresDataSourceConfig(rand, map[string]string{
			"ids":          `["${alicloud_sls_metric_store.default.id}"]`,
			"name_regex":   `"${alicloud_sls_metric_store.default.metric_store_name}"`,
			"project_name": `"${alicloud_log_project.default.project_name}"`,
		}),
		fakeConfig: testAccCheckAliCloudSlsMetricStoresDataSourceConfig(rand, map[string]string{
			"ids":          `["${alicloud_sls_metric_store.default.id}_fake"]`,
			"name_regex":   `"${alicloud_sls_metric_store.default.metric_store_name}_fake"`,
			"project_name": `"${alicloud_log_project.default.project_name}"`,
		}),
	}

	SlsMetricStoresCheckInfo.dataSourceTestCheck(t, rand, idsConf, nameRegexConf, allConf)
}

var existSlsMetricStoresMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":                                 "1",
		"names.#":                               "1",
		"metric_stores.#":                       "1",
		"metric_stores.0.id":                    CHECKSET,
		"metric_stores.0.metric_store_name":     CHECKSET,
		"metric_stores.0.ttl":                   CHECKSET,
		"metric_stores.0.shard_count":           CHECKSET,
		"metric_stores.0.auto_split":            CHECKSET,
		"metric_stores.0.max_split_shard_count": CHECKSET,
		"metric_stores.0.mode":                  CHECKSET,
		"metric_stores.0.create_time":           CHECKSET,
		"metric_stores.0.last_modify_time":      CHECKSET,
	}
}

var fakeSlsMetricStoresMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":           "0",
		"names.#":         "0",
		"metric_stores.#": "0",
	}
}

var SlsMetricStoresCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_sls_metric_stores.default",
	existMapFunc: existSlsMetricStoresMapFunc,
	fakeMapFunc:  fakeSlsMetricStoresMapFunc,
}

func testAccCheckAliCloudSlsMetricStoresDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testacc-sls-msds-%d"
}

resource "alicloud_log_project" "default" {
  project_name = var.name
  description  = "terraform test project for metric store datasource"
}

resource "alicloud_sls_metric_store" "default" {
  project_name      = alicloud_log_project.default.project_name
  metric_store_name = var.name
  ttl               = 30
  shard_count       = 2
}

data "alicloud_sls_metric_stores" "default" {
%s
}
`, rand, strings.Join(pairs, "\n  "))
	return config
}
