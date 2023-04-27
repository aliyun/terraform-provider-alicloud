package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudCloudMonitorServiceHybridDoubleWriteDataSource(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCloudMonitorServiceHybridDoubleWriteSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_cms_hybrid_double_write.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudCloudMonitorServiceHybridDoubleWriteSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_cms_hybrid_double_write.default.id}_fake"]`,
		}),
	}

	SourceNamespaceConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCloudMonitorServiceHybridDoubleWriteSourceConfig(rand, map[string]string{
			"ids":              `["${alicloud_cms_hybrid_double_write.default.id}"]`,
			"source_namespace": `"${alicloud_cms_hybrid_double_write.default.source_namespace}"`,
		}),
		fakeConfig: testAccCheckAlicloudCloudMonitorServiceHybridDoubleWriteSourceConfig(rand, map[string]string{
			"ids":              `["${alicloud_cms_hybrid_double_write.default.id}_fake"]`,
			"source_namespace": `"${alicloud_cms_hybrid_double_write.default.source_namespace}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCloudMonitorServiceHybridDoubleWriteSourceConfig(rand, map[string]string{
			"ids":              `["${alicloud_cms_hybrid_double_write.default.id}"]`,
			"source_namespace": `"${alicloud_cms_hybrid_double_write.default.source_namespace}"`,
		}),
		fakeConfig: testAccCheckAlicloudCloudMonitorServiceHybridDoubleWriteSourceConfig(rand, map[string]string{
			"ids":              `["${alicloud_cms_hybrid_double_write.default.id}_fake"]`,
			"source_namespace": `"${alicloud_cms_hybrid_double_write.default.source_namespace}_fake"`,
		}),
	}

	CloudMonitorServiceHybridDoubleWriteCheckInfo.dataSourceTestCheck(t, rand, idsConf, SourceNamespaceConf, allConf)
}

var existCloudMonitorServiceHybridDoubleWriteMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":                                   "1",
		"hybrid_double_writes.#":                  "1",
		"hybrid_double_writes.0.id":               CHECKSET,
		"hybrid_double_writes.0.namespace":        CHECKSET,
		"hybrid_double_writes.0.source_namespace": CHECKSET,
		"hybrid_double_writes.0.source_user_id":   CHECKSET,
		"hybrid_double_writes.0.user_id":          CHECKSET,
	}
}

var fakeCloudMonitorServiceHybridDoubleWriteMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":                  "0",
		"hybrid_double_writes.#": "0",
	}
}

var CloudMonitorServiceHybridDoubleWriteCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_cms_hybrid_double_writes.default",
	existMapFunc: existCloudMonitorServiceHybridDoubleWriteMapFunc,
	fakeMapFunc:  fakeCloudMonitorServiceHybridDoubleWriteMapFunc,
}

func testAccCheckAlicloudCloudMonitorServiceHybridDoubleWriteSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testacccmshybriddoublewrite%d"
}

data "alicloud_account" "default" {}

resource "alicloud_cms_namespace" "default" {
	count = 2
	description = var.name
	namespace = "${var.name}-${count.index}"
	specification = "cms.s1.large"
}

resource "alicloud_cms_hybrid_double_write" "default" {
  user_id          = data.alicloud_account.default.id
  source_namespace = alicloud_cms_namespace.default.0.namespace
  namespace        = alicloud_cms_namespace.default.1.namespace
}

data "alicloud_cms_hybrid_double_writes" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
