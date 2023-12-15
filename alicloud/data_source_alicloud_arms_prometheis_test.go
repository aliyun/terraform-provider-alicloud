package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudArmsPrometheisDataSource(t *testing.T) {
	rand := acctest.RandIntRange(10, 99)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudArmsPrometheisDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_arms_prometheus.default.id}"]`,
			"enable_details": `true`,
		}),
		fakeConfig: testAccCheckAlicloudArmsPrometheisDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_arms_prometheus.default.id}_fake"]`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudArmsPrometheisDataSourceName(rand, map[string]string{
			"name_regex":     `"${alicloud_arms_prometheus.default.cluster_name}"`,
			"enable_details": `true`,
		}),
		fakeConfig: testAccCheckAlicloudArmsPrometheisDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_arms_prometheus.default.cluster_name}_fake"`,
		}),
	}
	resourceGroupIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudArmsPrometheisDataSourceName(rand, map[string]string{
			"name_regex":        `"${alicloud_arms_prometheus.default.cluster_name}"`,
			"resource_group_id": `"${alicloud_arms_prometheus.default.resource_group_id}"`,
			"enable_details":    `true`,
		}),
		fakeConfig: testAccCheckAlicloudArmsPrometheisDataSourceName(rand, map[string]string{
			"resource_group_id": `"${alicloud_arms_prometheus.default.resource_group_id}_fake"`,
		}),
	}
	tagsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudArmsPrometheisDataSourceName(rand, map[string]string{
			"name_regex":     `"${alicloud_arms_prometheus.default.cluster_name}"`,
			"enable_details": `true`,
			"tags": `{
							Created = "TF"
							For 	= "Prometheus"
					  }`,
		}),
		fakeConfig: testAccCheckAlicloudArmsPrometheisDataSourceName(rand, map[string]string{
			"tags": `{
							Created = "TF_Update"
							For 	= "Prometheus_Update"
					  }`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudArmsPrometheisDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_arms_prometheus.default.id}"]`,
			"name_regex":        `"${alicloud_arms_prometheus.default.cluster_name}"`,
			"enable_details":    `true`,
			"resource_group_id": `"${alicloud_arms_prometheus.default.resource_group_id}"`,
			"tags": `{
							Created = "TF"
							For 	= "Prometheus"
					  }`,
		}),
		fakeConfig: testAccCheckAlicloudArmsPrometheisDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_arms_prometheus.default.id}_fake"]`,
			"name_regex":        `"${alicloud_arms_prometheus.default.cluster_name}_fake"`,
			"resource_group_id": `"${alicloud_arms_prometheus.default.resource_group_id}_fake"`,
			"tags": `{
							Created = "TF_Update"
							For 	= "Prometheus_Update"
					  }`,
		}),
	}
	var existAlicloudArmsPrometheisDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                               "1",
			"names.#":                             "1",
			"prometheis.#":                        "1",
			"prometheis.0.id":                     CHECKSET,
			"prometheis.0.cluster_id":             CHECKSET,
			"prometheis.0.cluster_type":           "ecs",
			"prometheis.0.cluster_name":           CHECKSET,
			"prometheis.0.vpc_id":                 CHECKSET,
			"prometheis.0.vswitch_id":             CHECKSET,
			"prometheis.0.security_group_id":      CHECKSET,
			"prometheis.0.sub_clusters_json":      "",
			"prometheis.0.grafana_instance_id":    "free",
			"prometheis.0.resource_group_id":      CHECKSET,
			"prometheis.0.tags.%":                 "2",
			"prometheis.0.tags.Created":           "TF",
			"prometheis.0.tags.For":               "Prometheus",
			"prometheis.0.remote_read_intra_url":  CHECKSET,
			"prometheis.0.remote_read_inter_url":  CHECKSET,
			"prometheis.0.remote_write_intra_url": CHECKSET,
			"prometheis.0.remote_write_inter_url": CHECKSET,
			//"prometheis.0.push_gate_way_intra_url": CHECKSET,
			//"prometheis.0.push_gate_way_inter_url": CHECKSET,
			"prometheis.0.http_api_intra_url": CHECKSET,
			"prometheis.0.http_api_inter_url": CHECKSET,
			//"prometheis.0.auth_token":              CHECKSET,
		}
	}
	var fakeAlicloudArmsPrometheisDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":        "0",
			"names.#":      "0",
			"prometheis.#": "0",
		}
	}
	var alicloudArmsPrometheisCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_arms_prometheis.default",
		existMapFunc: existAlicloudArmsPrometheisDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudArmsPrometheisDataSourceNameMapFunc,
	}
	preCheck := func() {
		testAccPreCheck(t)
	}
	alicloudArmsPrometheisCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, resourceGroupIdConf, tagsConf, allConf)
}

func testAccCheckAlicloudArmsPrometheisDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
	variable "name" {
  		default = "tf-ArmsPrometheus-%d"
	}

	data "alicloud_vpcs" "default" {
  		name_regex = "default-NODELETING"
	}

	data "alicloud_vswitches" "default" {
 		 vpc_id = data.alicloud_vpcs.default.ids.0
	}

	data "alicloud_resource_manager_resource_groups" "default" {
	}

	resource "alicloud_security_group" "default" {
  		vpc_id = data.alicloud_vpcs.default.ids.0
	}

	resource "alicloud_arms_prometheus" "default" {
  		cluster_type        = "ecs"
  		grafana_instance_id = "free"
  		vpc_id              = data.alicloud_vpcs.default.ids.0
  		vswitch_id          = data.alicloud_vswitches.default.ids.0
  		security_group_id   = alicloud_security_group.default.id
  		cluster_name        = "${var.name}-${data.alicloud_vpcs.default.ids.0}"
  		resource_group_id   = data.alicloud_resource_manager_resource_groups.default.groups.1.id
  		tags = {
			Created = "TF"
    		For     = "Prometheus"
  		}
	}

	data "alicloud_arms_prometheis" "default" {
		%s
	}
`, rand, strings.Join(pairs, " \n "))
	return config
}
