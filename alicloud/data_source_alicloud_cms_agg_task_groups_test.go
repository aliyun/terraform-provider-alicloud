package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAliCloudCmsAggTaskGroupsDataSource_basic(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudCmsAggTaskGroupsDataSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_cms_agg_task_group.default.agg_task_group_id}"]`,
		}),
		fakeConfig: testAccCheckAliCloudCmsAggTaskGroupsDataSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_cms_agg_task_group.default.agg_task_group_id}_fake"]`,
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudCmsAggTaskGroupsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_cms_agg_task_group.default.agg_task_group_name}"`,
		}),
		fakeConfig: testAccCheckAliCloudCmsAggTaskGroupsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_cms_agg_task_group.default.agg_task_group_name}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudCmsAggTaskGroupsDataSourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_cms_agg_task_group.default.agg_task_group_id}"]`,
			"name_regex": `"${alicloud_cms_agg_task_group.default.agg_task_group_name}"`,
		}),
		fakeConfig: testAccCheckAliCloudCmsAggTaskGroupsDataSourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_cms_agg_task_group.default.agg_task_group_id}_fake"]`,
			"name_regex": `"${alicloud_cms_agg_task_group.default.agg_task_group_name}_fake"`,
		}),
	}

	var existCmsAggTaskGroupsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                         "1",
			"names.#":                       "1",
			"groups.#":                      "1",
			"groups.0.id":                   CHECKSET,
			"groups.0.agg_task_group_id":    CHECKSET,
			"groups.0.agg_task_group_name":  CHECKSET,
			"groups.0.source_prometheus_id": CHECKSET,
			"groups.0.target_prometheus_id": CHECKSET,
			"groups.0.status":               CHECKSET,
		}
	}

	var fakeCmsAggTaskGroupsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":    "0",
			"names.#":  "0",
			"groups.#": "0",
		}
	}

	var cmsAggTaskGroupsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_cms_agg_task_groups.default",
		existMapFunc: existCmsAggTaskGroupsMapFunc,
		fakeMapFunc:  fakeCmsAggTaskGroupsMapFunc,
	}

	cmsAggTaskGroupsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, func() {
		testAccPreCheck(t)
		testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	}, idsConf, nameRegexConf, allConf)
}

func testAccCheckAliCloudCmsAggTaskGroupsDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
variable "name" {
  default = "tfacccms%d"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_log_project" "default" {
  project_name = var.name
}

resource "alicloud_cms_workspace" "default" {
  workspace_name = var.name
  sls_project    = alicloud_log_project.default.project_name
}

resource "alicloud_cms_prometheus_instance" "default" {
  count                    = 2
  prometheus_instance_name = "${var.name}_${count.index}"
  workspace                = alicloud_cms_workspace.default.id
}

resource "alicloud_cms_agg_task_group" "default" {
  source_prometheus_id  = alicloud_cms_prometheus_instance.default.0.id
  target_prometheus_id  = alicloud_cms_prometheus_instance.default.1.id
  agg_task_group_name   = var.name
  agg_task_group_config = <<EOF
groups:
- name: "node.rules"
  interval: "60s"
  rules:
  - record: "node_namespace_pod:kube_pod_info:"
    expr: "max(label_replace(kube_pod_info{job=\"kubernetes-pods-kube-state-metrics\" }, \"pod\", \"$1\", \"pod\", \"(.*)\")) by (node, namespace, pod, cluster)"
EOF
}

data "alicloud_cms_agg_task_groups" "default" {
  source_prometheus_id = alicloud_cms_agg_task_group.default.source_prometheus_id
  %s
}
`, rand, strings.Join(pairs, "\n  "))
	return config
}
