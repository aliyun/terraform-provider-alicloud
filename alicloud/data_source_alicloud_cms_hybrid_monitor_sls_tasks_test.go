package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudCmsHybridMonitorSlsTasksDataSource(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)
	checkoutSupportedRegions(t, true, connectivity.CloudMonitorServiceSupportRegions)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCmsHybridMonitorSlsTasksDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_cms_hybrid_monitor_sls_task.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudCmsHybridMonitorSlsTasksDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_cms_hybrid_monitor_sls_task.default.id}_fake"]`,
		}),
	}
	keywordConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCmsHybridMonitorSlsTasksDataSourceName(rand, map[string]string{
			"ids":     `["${alicloud_cms_hybrid_monitor_sls_task.default.id}"]`,
			"keyword": `"${alicloud_cms_hybrid_monitor_sls_task.default.task_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudCmsHybridMonitorSlsTasksDataSourceName(rand, map[string]string{
			"ids":     `["${alicloud_cms_hybrid_monitor_sls_task.default.id}"]`,
			"keyword": `"${alicloud_cms_hybrid_monitor_sls_task.default.task_name}_fake"`,
		}),
	}
	namespaceConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCmsHybridMonitorSlsTasksDataSourceName(rand, map[string]string{
			"ids":       `["${alicloud_cms_hybrid_monitor_sls_task.default.id}"]`,
			"namespace": `"${alicloud_cms_hybrid_monitor_sls_task.default.namespace}"`,
		}),
		fakeConfig: testAccCheckAlicloudCmsHybridMonitorSlsTasksDataSourceName(rand, map[string]string{
			"ids":       `["${alicloud_cms_hybrid_monitor_sls_task.default.id}"]`,
			"namespace": `"${alicloud_cms_hybrid_monitor_sls_task.default.namespace}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCmsHybridMonitorSlsTasksDataSourceName(rand, map[string]string{
			"ids":       `["${alicloud_cms_hybrid_monitor_sls_task.default.id}"]`,
			"keyword":   `"${alicloud_cms_hybrid_monitor_sls_task.default.task_name}"`,
			"namespace": `"${alicloud_cms_hybrid_monitor_sls_task.default.namespace}"`,
		}),
		fakeConfig: testAccCheckAlicloudCmsHybridMonitorSlsTasksDataSourceName(rand, map[string]string{
			"ids":       `["${alicloud_cms_hybrid_monitor_sls_task.default.id}_fake"]`,
			"keyword":   `"${alicloud_cms_hybrid_monitor_sls_task.default.task_name}_fake"`,
			"namespace": `"${alicloud_cms_hybrid_monitor_sls_task.default.namespace}_fake"`,
		}),
	}
	var existAlicloudCmsHybridMonitorSlsTasksDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                                        "1",
			"tasks.#":                                                      "1",
			"tasks.0.attach_labels.#":                                      "1",
			"tasks.0.attach_labels.0.name":                                 "app_service",
			"tasks.0.attach_labels.0.value":                                "testValue",
			"tasks.0.collect_interval":                                     "60",
			"tasks.0.collect_target_type":                                  CHECKSET,
			"tasks.0.description":                                          CHECKSET,
			"tasks.0.group_id":                                             CHECKSET,
			"tasks.0.task_name":                                            CHECKSET,
			"tasks.0.task_type":                                            "aliyun_sls",
			"tasks.0.sls_process_config.#":                                 "1",
			"tasks.0.sls_process_config.0.filter.#":                        "1",
			"tasks.0.sls_process_config.0.filter.0.relation":               "and",
			"tasks.0.sls_process_config.0.filter.0.filters.#":              "1",
			"tasks.0.sls_process_config.0.filter.0.filters.0.operator":     "=",
			"tasks.0.sls_process_config.0.filter.0.filters.0.value":        "200",
			"tasks.0.sls_process_config.0.filter.0.filters.0.sls_key_name": "code",
			"tasks.0.sls_process_config.0.statistics.#":                    "1",
			"tasks.0.sls_process_config.0.statistics.0.function":           "count",
			"tasks.0.sls_process_config.0.statistics.0.alias":              "level_count",
			"tasks.0.sls_process_config.0.statistics.0.sls_key_name":       "name",
			"tasks.0.sls_process_config.0.statistics.0.parameter_one":      "200",
			"tasks.0.sls_process_config.0.statistics.0.parameter_two":      "299",
			"tasks.0.sls_process_config.0.express.#":                       "1",
			"tasks.0.sls_process_config.0.express.0.express":               "success_count",
			"tasks.0.sls_process_config.0.express.0.alias":                 "SuccRate",
			"tasks.0.sls_process_config.0.group_by.#":                      "1",
			"tasks.0.sls_process_config.0.group_by.0.alias":                "code",
			"tasks.0.sls_process_config.0.group_by.0.sls_key_name":         "ApiResult",
			"tasks.0.match_express.#":                                      "0",
			"tasks.0.id":                                                   CHECKSET,
			"tasks.0.collect_target_endpoint":                              "",
			"tasks.0.collect_target_path":                                  "",
			"tasks.0.collect_timout":                                       CHECKSET,
			"tasks.0.create_time":                                          CHECKSET,
			"tasks.0.extra_info":                                           "",
			"tasks.0.hybrid_monitor_sls_task_id":                           CHECKSET,
			"tasks.0.instances.#":                                          "0",
			"tasks.0.log_file_path":                                        "",
			"tasks.0.log_process":                                          "",
			"tasks.0.log_sample":                                           "",
			"tasks.0.log_split":                                            "",
			"tasks.0.match_express_relation":                               "",
			"tasks.0.namespace":                                            CHECKSET,
			"tasks.0.network_type":                                         CHECKSET,
			"tasks.0.sls_process":                                          CHECKSET,
			"tasks.0.upload_region":                                        CHECKSET,
			"tasks.0.yarm_config":                                          "",
		}
	}
	var fakeAlicloudCmsHybridMonitorSlsTasksDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#": "0",
		}
	}
	var alicloudCmsHybridMonitorSlsTasksCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_cms_hybrid_monitor_sls_tasks.default",
		existMapFunc: existAlicloudCmsHybridMonitorSlsTasksDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudCmsHybridMonitorSlsTasksDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}
	alicloudCmsHybridMonitorSlsTasksCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, keywordConf, namespaceConf, allConf)
}
func testAccCheckAlicloudCmsHybridMonitorSlsTasksDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {
	default = "tf_testacc_cms_slstask%d"
}

data "alicloud_account" "this" {}

resource "alicloud_cms_sls_group" "default" {
	sls_group_config {
		sls_user_id = data.alicloud_account.this.id
		sls_logstore = "Logstore-ECS"
		sls_project = "aliyun-project"
		sls_region = "cn-hangzhou"
	}
	sls_group_description = var.name
	sls_group_name = var.name
}

resource "alicloud_cms_namespace" "default" {
	description = var.name
	namespace = "tf-testacc-cloudmonitorservicenamespace"
	specification = "cms.s1.large"
}

resource "alicloud_cms_hybrid_monitor_sls_task" "default" {
  sls_process_config {
    filter {
      relation = "and"
      filters {
        operator     = "="
        value        = "200"
        sls_key_name = "code"
      }
    }
    statistics {
      function      = "count"
      alias         = "level_count"
      sls_key_name  = "name"
      parameter_one = "200"
      parameter_two = "299"
    }
    group_by {
      alias        = "code"
      sls_key_name = "ApiResult"
    }
    express {
      express = "success_count"
      alias   = "SuccRate"
    }
  }
  task_name           = var.name
  namespace           = alicloud_cms_namespace.default.id
  description         = var.name
  collect_interval    = 60
  collect_target_type = alicloud_cms_sls_group.default.id
  attach_labels {
    name  = "app_service"
    value = "testValue"
  }
}

data "alicloud_cms_hybrid_monitor_sls_tasks" "default" {	
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
