package cloudmonitorservice_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudCmsHybridMonitorFcTasksDataSource(t *testing.T) {
	rand := tfacctest.RandInt()
	acctest.CheckoutSupportedRegions(t, true, connectivity.CloudMonitorServiceSupportRegions)
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudCmsHybridMonitorFcTasksDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_cms_hybrid_monitor_fc_task.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudCmsHybridMonitorFcTasksDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_cms_hybrid_monitor_fc_task.default.id}_fake"]`,
		}),
	}
	namespaceConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudCmsHybridMonitorFcTasksDataSourceName(rand, map[string]string{
			"ids":       `["${alicloud_cms_hybrid_monitor_fc_task.default.id}"]`,
			"namespace": `"${alicloud_cms_hybrid_monitor_fc_task.default.namespace}"`,
		}),
		FakeConfig: testAccCheckAlicloudCmsHybridMonitorFcTasksDataSourceName(rand, map[string]string{
			"ids":       `["${alicloud_cms_hybrid_monitor_fc_task.default.id}"]`,
			"namespace": `"${alicloud_cms_hybrid_monitor_fc_task.default.namespace}_fake"`,
		}),
	}
	var existAlicloudCmsHybridMonitorFcTasksDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                             "1",
			"tasks.#":                           "1",
			"tasks.0.namespace":                 acctest.CHECKSET,
			"tasks.0.target_user_id":            acctest.CHECKSET,
			"tasks.0.yarm_config":               acctest.CHECKSET,
			"tasks.0.create_time":               "",
			"tasks.0.id":                        acctest.CHECKSET,
			"tasks.0.hybrid_monitor_fc_task_id": acctest.CHECKSET,
		}
	}
	var fakeAlicloudCmsHybridMonitorFcTasksDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#": "0",
		}
	}
	var alicloudCmsHybridMonitorFcTasksCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_cms_hybrid_monitor_fc_tasks.default",
		ExistMapFunc: existAlicloudCmsHybridMonitorFcTasksDataSourceNameMapFunc,
		FakeMapFunc:  fakeAlicloudCmsHybridMonitorFcTasksDataSourceNameMapFunc,
	}

	preCheck := func() {
		acctest.PreCheck(t)
	}
	alicloudCmsHybridMonitorFcTasksCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, namespaceConf)
}
func testAccCheckAlicloudCmsHybridMonitorFcTasksDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testacchybridmonitorfctask-%d"
}
data "alicloud_account" "this" {}

resource "alicloud_cms_namespace" "default" {
	description = var.name
	namespace = var.name
	specification = "cms.s1.large"
}

resource "alicloud_cms_hybrid_monitor_fc_task" "default" {
	namespace = alicloud_cms_namespace.default.id
	target_user_id = data.alicloud_account.this.id
	yarm_config = "---\nproducts:\n- namespace: \"acs_ecs_dashboard\"\n  metric_info:\n  - metric_list:\n    - \"cpu_total\"\n    - \"diskusage_utilization\"\n    - \"memory_usedutilization\"\n    - \"net_tcpconnection\"\n- namespace: \"acs_ecs_dashboard1\"\n  metric_info:\n  - metric_list:\n    - \"CPUUtilization\"\n    - \"DiskReadBPS\"\n    - \"InternetOut\"\n    - \"IntranetOut\"\n    - \"cpu_idle\"\n    - \"cpu_system\"\n    - \"cpu_total\"\n    - \"diskusage_utilization\"\n- namespace: \"acs_rds_dashboard\"\n  metric_info:\n  - metric_list:\n    - \"MySQL_QPS\"\n    - \"MySQL_TPS\"\n"
}

data "alicloud_cms_hybrid_monitor_fc_tasks" "default" {
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
