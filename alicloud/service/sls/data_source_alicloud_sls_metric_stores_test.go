package sls_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudSlsMetricStoresDataSource(t *testing.T) {
	acctest.PreCheckWithRegions(t, true, connectivity.SlsTestRegions)
	rand := tfacctest.RandIntRange(10000, 99999)

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudSlsMetricStoresDataSourceConfig(rand, map[string]string{
			"ids":          `["${alicloud_sls_metric_store.default.id}"]`,
			"project_name": `"${alicloud_log_project.default.project_name}"`,
		}),
		FakeConfig: testAccCheckAliCloudSlsMetricStoresDataSourceConfig(rand, map[string]string{
			"ids":          `["${alicloud_sls_metric_store.default.id}_fake"]`,
			"project_name": `"${alicloud_log_project.default.project_name}"`,
		}),
	}

	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudSlsMetricStoresDataSourceConfig(rand, map[string]string{
			"name_regex":   `"${alicloud_sls_metric_store.default.metric_store_name}"`,
			"project_name": `"${alicloud_log_project.default.project_name}"`,
		}),
		FakeConfig: testAccCheckAliCloudSlsMetricStoresDataSourceConfig(rand, map[string]string{
			"name_regex":   `"${alicloud_sls_metric_store.default.metric_store_name}_fake"`,
			"project_name": `"${alicloud_log_project.default.project_name}"`,
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudSlsMetricStoresDataSourceConfig(rand, map[string]string{
			"ids":          `["${alicloud_sls_metric_store.default.id}"]`,
			"name_regex":   `"${alicloud_sls_metric_store.default.metric_store_name}"`,
			"project_name": `"${alicloud_log_project.default.project_name}"`,
		}),
		FakeConfig: testAccCheckAliCloudSlsMetricStoresDataSourceConfig(rand, map[string]string{
			"ids":          `["${alicloud_sls_metric_store.default.id}_fake"]`,
			"name_regex":   `"${alicloud_sls_metric_store.default.metric_store_name}_fake"`,
			"project_name": `"${alicloud_log_project.default.project_name}"`,
		}),
	}

	SlsMetricStoresCheckInfo.DataSourceTestCheck(t, rand, idsConf, nameRegexConf, allConf)
}

var existSlsMetricStoresMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":                                 "1",
		"names.#":                               "1",
		"metric_stores.#":                       "1",
		"metric_stores.0.id":                    acctest.CHECKSET,
		"metric_stores.0.metric_store_name":     acctest.CHECKSET,
		"metric_stores.0.ttl":                   acctest.CHECKSET,
		"metric_stores.0.shard_count":           acctest.CHECKSET,
		"metric_stores.0.auto_split":            acctest.CHECKSET,
		"metric_stores.0.max_split_shard_count": acctest.CHECKSET,
		"metric_stores.0.mode":                  acctest.CHECKSET,
		"metric_stores.0.create_time":           acctest.CHECKSET,
		"metric_stores.0.last_modify_time":      acctest.CHECKSET,
	}
}

var fakeSlsMetricStoresMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":           "0",
		"names.#":         "0",
		"metric_stores.#": "0",
	}
}

var SlsMetricStoresCheckInfo = acctest.DataSourceAttr{
	ResourceId:   "data.alicloud_sls_metric_stores.default",
	ExistMapFunc: existSlsMetricStoresMapFunc,
	FakeMapFunc:  fakeSlsMetricStoresMapFunc,
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
