package cloudmonitorservice_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudCloudMonitorServiceHybridDoubleWritesDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(100, 999)

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudCloudMonitorServiceHybridDoubleWritesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_cloud_monitor_service_hybrid_double_write.default.id}"]`,
		}),
		FakeConfig: testAccCheckAliCloudCloudMonitorServiceHybridDoubleWritesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_cloud_monitor_service_hybrid_double_write.default.id}_fake"]`,
		}),
	}
	sourceNamespaceConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudCloudMonitorServiceHybridDoubleWritesDataSourceName(rand, map[string]string{
			"source_namespace": `"${alicloud_cloud_monitor_service_hybrid_double_write.default.source_namespace}"`,
		}),
		FakeConfig: testAccCheckAliCloudCloudMonitorServiceHybridDoubleWritesDataSourceName(rand, map[string]string{
			"source_namespace": `"${alicloud_cloud_monitor_service_hybrid_double_write.default.source_namespace}_fake"`,
		}),
	}
	sourceUserIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudCloudMonitorServiceHybridDoubleWritesDataSourceName(rand, map[string]string{
			"source_user_id": `"${alicloud_cloud_monitor_service_hybrid_double_write.default.source_user_id}"`,
		}),
	}
	namespaceConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudCloudMonitorServiceHybridDoubleWritesDataSourceName(rand, map[string]string{
			"namespace": `"${alicloud_cloud_monitor_service_hybrid_double_write.default.namespace}"`,
		}),
		FakeConfig: testAccCheckAliCloudCloudMonitorServiceHybridDoubleWritesDataSourceName(rand, map[string]string{
			"namespace": `"${alicloud_cloud_monitor_service_hybrid_double_write.default.namespace}_fake"`,
		}),
	}
	userIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudCloudMonitorServiceHybridDoubleWritesDataSourceName(rand, map[string]string{
			"user_id": `"${alicloud_cloud_monitor_service_hybrid_double_write.default.user_id}"`,
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudCloudMonitorServiceHybridDoubleWritesDataSourceName(rand, map[string]string{
			"ids":              `["${alicloud_cloud_monitor_service_hybrid_double_write.default.id}"]`,
			"source_namespace": `"${alicloud_cloud_monitor_service_hybrid_double_write.default.source_namespace}"`,
			"source_user_id":   `"${alicloud_cloud_monitor_service_hybrid_double_write.default.source_user_id}"`,
			"namespace":        `"${alicloud_cloud_monitor_service_hybrid_double_write.default.namespace}"`,
			"user_id":          `"${alicloud_cloud_monitor_service_hybrid_double_write.default.user_id}"`,
		}),
		FakeConfig: testAccCheckAliCloudCloudMonitorServiceHybridDoubleWritesDataSourceName(rand, map[string]string{
			"ids":              `["${alicloud_cloud_monitor_service_hybrid_double_write.default.id}_fake"]`,
			"source_namespace": `"${alicloud_cloud_monitor_service_hybrid_double_write.default.source_namespace}_fake"`,
			"namespace":        `"${alicloud_cloud_monitor_service_hybrid_double_write.default.namespace}_fake"`,
		}),
	}
	var existAliCloudCloudMonitorServiceHybridDoubleWritesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                   "1",
			"hybrid_double_writes.#":                  "1",
			"hybrid_double_writes.0.id":               acctest.CHECKSET,
			"hybrid_double_writes.0.source_namespace": acctest.CHECKSET,
			"hybrid_double_writes.0.source_user_id":   acctest.CHECKSET,
			"hybrid_double_writes.0.namespace":        acctest.CHECKSET,
			"hybrid_double_writes.0.user_id":          acctest.CHECKSET,
		}
	}
	var fakeAliCloudCloudMonitorServiceHybridDoubleWritesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                  "0",
			"hybrid_double_writes.#": "0",
		}
	}
	var alicloudCloudMonitorServiceHybridDoubleWritesCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_cloud_monitor_service_hybrid_double_writes.default",
		ExistMapFunc: existAliCloudCloudMonitorServiceHybridDoubleWritesDataSourceNameMapFunc,
		FakeMapFunc:  fakeAliCloudCloudMonitorServiceHybridDoubleWritesDataSourceNameMapFunc,
	}
	preCheck := func() {
		acctest.PreCheck(t)
	}
	alicloudCloudMonitorServiceHybridDoubleWritesCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, sourceNamespaceConf, sourceUserIdConf, namespaceConf, userIdConf, allConf)
}

func testAccCheckAliCloudCloudMonitorServiceHybridDoubleWritesDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
	variable "name" {
  		default = "tf-testacc-chw-%d"
	}

	data "alicloud_account" "default" {
	}

	resource "alicloud_cms_namespace" "source" {
  		namespace = "${var.name}-source"
	}

	resource "alicloud_cms_namespace" "default" {
  		namespace = var.name
	}

	resource "alicloud_cloud_monitor_service_hybrid_double_write" "default" {
  		source_namespace = alicloud_cms_namespace.source.id
  		source_user_id   = data.alicloud_account.default.id
  		namespace        = alicloud_cms_namespace.default.id
  		user_id          = data.alicloud_account.default.id
	}

	data "alicloud_cloud_monitor_service_hybrid_double_writes" "default" {
		%s
	}
`, rand, strings.Join(pairs, " \n "))
	return config
}
