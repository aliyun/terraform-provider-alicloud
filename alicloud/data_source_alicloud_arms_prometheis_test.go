package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAliCloudArmsPrometheisDataSource_basic0(t *testing.T) {
	rand := acctest.RandIntRange(10, 99)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudArmsPrometheisDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_arms_prometheus.default.id}"]`,
		}),
		fakeConfig: testAccCheckAliCloudArmsPrometheisDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_arms_prometheus.default.id}_fake"]`,
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudArmsPrometheisDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_arms_prometheus.default.cluster_name}"`,
		}),
		fakeConfig: testAccCheckAliCloudArmsPrometheisDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_arms_prometheus.default.cluster_name}_fake"`,
		}),
	}

	resourceGroupIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudArmsPrometheisDataSourceName(rand, map[string]string{
			"name_regex":        `"${alicloud_arms_prometheus.default.cluster_name}"`,
			"resource_group_id": `"${alicloud_arms_prometheus.default.resource_group_id}"`,
		}),
		fakeConfig: testAccCheckAliCloudArmsPrometheisDataSourceName(rand, map[string]string{
			"resource_group_id": `"${alicloud_arms_prometheus.default.resource_group_id}_fake"`,
		}),
	}

	tagsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudArmsPrometheisDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_arms_prometheus.default.cluster_name}"`,
			"tags": `{
							Created = "TF"
							For 	= "Prometheus"
					  }`,
		}),
		fakeConfig: testAccCheckAliCloudArmsPrometheisDataSourceName(rand, map[string]string{
			"tags": `{
							Created = "TF_Update"
							For 	= "Prometheus_Update"
					  }`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudArmsPrometheisDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_arms_prometheus.default.id}"]`,
			"name_regex":        `"${alicloud_arms_prometheus.default.cluster_name}"`,
			"resource_group_id": `"${alicloud_arms_prometheus.default.resource_group_id}"`,
			"tags": `{
							Created = "TF"
							For 	= "Prometheus"
					  }`,
		}),
		fakeConfig: testAccCheckAliCloudArmsPrometheisDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_arms_prometheus.default.id}_fake"]`,
			"name_regex":        `"${alicloud_arms_prometheus.default.cluster_name}_fake"`,
			"resource_group_id": `"${alicloud_arms_prometheus.default.resource_group_id}_fake"`,
			"tags": `{
							Created = "TF_Update"
							For 	= "Prometheus_Update"
					  }`,
		}),
	}

	var existAliCloudArmsPrometheisDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                "1",
			"names.#":                              "1",
			"prometheis.#":                         "1",
			"prometheis.0.id":                      CHECKSET,
			"prometheis.0.cluster_id":              CHECKSET,
			"prometheis.0.cluster_type":            "ecs",
			"prometheis.0.cluster_name":            CHECKSET,
			"prometheis.0.vpc_id":                  CHECKSET,
			"prometheis.0.vswitch_id":              CHECKSET,
			"prometheis.0.security_group_id":       CHECKSET,
			"prometheis.0.sub_clusters_json":       "",
			"prometheis.0.grafana_instance_id":     "free",
			"prometheis.0.resource_group_id":       CHECKSET,
			"prometheis.0.remote_read_intra_url":   "",
			"prometheis.0.remote_read_inter_url":   "",
			"prometheis.0.remote_write_intra_url":  "",
			"prometheis.0.remote_write_inter_url":  "",
			"prometheis.0.push_gate_way_intra_url": "",
			"prometheis.0.push_gate_way_inter_url": "",
			"prometheis.0.http_api_intra_url":      "",
			"prometheis.0.http_api_inter_url":      "",
			"prometheis.0.auth_token":              "",
			"prometheis.0.tags.%":                  "2",
			"prometheis.0.tags.Created":            "TF",
			"prometheis.0.tags.For":                "Prometheus",
		}
	}

	var fakeAliCloudArmsPrometheisDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":        "0",
			"names.#":      "0",
			"prometheis.#": "0",
		}
	}

	var aliCloudArmsPrometheisCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_arms_prometheus.default",
		existMapFunc: existAliCloudArmsPrometheisDataSourceNameMapFunc,
		fakeMapFunc:  fakeAliCloudArmsPrometheisDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}

	aliCloudArmsPrometheisCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, resourceGroupIdConf, tagsConf, allConf)
}

func TestAccAliCloudArmsPrometheisDataSource_basic1(t *testing.T) {
	rand := acctest.RandIntRange(10, 99)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudArmsPrometheisDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_arms_prometheus.default.id}"]`,
			"enable_details": `true`,
		}),
		fakeConfig: testAccCheckAliCloudArmsPrometheisDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_arms_prometheus.default.id}"]`,
			"enable_details": `false`,
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudArmsPrometheisDataSourceName(rand, map[string]string{
			"name_regex":     `"${alicloud_arms_prometheus.default.cluster_name}"`,
			"enable_details": `true`,
		}),
		fakeConfig: testAccCheckAliCloudArmsPrometheisDataSourceName(rand, map[string]string{
			"name_regex":     `"${alicloud_arms_prometheus.default.cluster_name}"`,
			"enable_details": `false`,
		}),
	}

	resourceGroupIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudArmsPrometheisDataSourceName(rand, map[string]string{
			"name_regex":        `"${alicloud_arms_prometheus.default.cluster_name}"`,
			"resource_group_id": `"${alicloud_arms_prometheus.default.resource_group_id}"`,
			"enable_details":    `true`,
		}),
		fakeConfig: testAccCheckAliCloudArmsPrometheisDataSourceName(rand, map[string]string{
			"name_regex":        `"${alicloud_arms_prometheus.default.cluster_name}"`,
			"resource_group_id": `"${alicloud_arms_prometheus.default.resource_group_id}"`,
			"enable_details":    `false`,
		}),
	}

	tagsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudArmsPrometheisDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_arms_prometheus.default.cluster_name}"`,
			"tags": `{
							Created = "TF"
							For 	= "Prometheus"
					  }`,
			"enable_details": `true`,
		}),
		fakeConfig: testAccCheckAliCloudArmsPrometheisDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_arms_prometheus.default.cluster_name}"`,
			"tags": `{
							Created = "TF"
							For 	= "Prometheus"
					  }`,
			"enable_details": `false`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudArmsPrometheisDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_arms_prometheus.default.id}"]`,
			"name_regex":        `"${alicloud_arms_prometheus.default.cluster_name}"`,
			"resource_group_id": `"${alicloud_arms_prometheus.default.resource_group_id}"`,
			"tags": `{
							Created = "TF"
							For 	= "Prometheus"
					  }`,
			"enable_details": `true`,
		}),
		fakeConfig: testAccCheckAliCloudArmsPrometheisDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_arms_prometheus.default.id}"]`,
			"name_regex":        `"${alicloud_arms_prometheus.default.cluster_name}"`,
			"resource_group_id": `"${alicloud_arms_prometheus.default.resource_group_id}"`,
			"tags": `{
							Created = "TF"
							For 	= "Prometheus"
					  }`,
			"enable_details": `false`,
		}),
	}

	var existAliCloudArmsPrometheisDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                "1",
			"names.#":                              "1",
			"prometheis.#":                         "1",
			"prometheis.0.id":                      CHECKSET,
			"prometheis.0.cluster_id":              CHECKSET,
			"prometheis.0.cluster_type":            "ecs",
			"prometheis.0.cluster_name":            CHECKSET,
			"prometheis.0.vpc_id":                  CHECKSET,
			"prometheis.0.vswitch_id":              CHECKSET,
			"prometheis.0.security_group_id":       CHECKSET,
			"prometheis.0.sub_clusters_json":       "",
			"prometheis.0.grafana_instance_id":     "free",
			"prometheis.0.resource_group_id":       CHECKSET,
			"prometheis.0.remote_read_intra_url":   CHECKSET,
			"prometheis.0.remote_read_inter_url":   CHECKSET,
			"prometheis.0.remote_write_intra_url":  CHECKSET,
			"prometheis.0.remote_write_inter_url":  CHECKSET,
			"prometheis.0.push_gate_way_intra_url": CHECKSET,
			"prometheis.0.push_gate_way_inter_url": CHECKSET,
			"prometheis.0.http_api_intra_url":      CHECKSET,
			"prometheis.0.http_api_inter_url":      CHECKSET,
			"prometheis.0.auth_token":              "",
			"prometheis.0.tags.%":                  "2",
			"prometheis.0.tags.Created":            "TF",
			"prometheis.0.tags.For":                "Prometheus",
		}
	}

	var fakeAliCloudArmsPrometheisDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                "1",
			"names.#":                              "1",
			"prometheis.#":                         "1",
			"prometheis.0.id":                      CHECKSET,
			"prometheis.0.cluster_id":              CHECKSET,
			"prometheis.0.cluster_type":            "ecs",
			"prometheis.0.cluster_name":            CHECKSET,
			"prometheis.0.vpc_id":                  CHECKSET,
			"prometheis.0.vswitch_id":              CHECKSET,
			"prometheis.0.security_group_id":       CHECKSET,
			"prometheis.0.sub_clusters_json":       "",
			"prometheis.0.grafana_instance_id":     "free",
			"prometheis.0.resource_group_id":       CHECKSET,
			"prometheis.0.remote_read_intra_url":   "",
			"prometheis.0.remote_read_inter_url":   "",
			"prometheis.0.remote_write_intra_url":  "",
			"prometheis.0.remote_write_inter_url":  "",
			"prometheis.0.push_gate_way_intra_url": "",
			"prometheis.0.push_gate_way_inter_url": "",
			"prometheis.0.http_api_intra_url":      "",
			"prometheis.0.http_api_inter_url":      "",
			"prometheis.0.auth_token":              "",
			"prometheis.0.tags.%":                  "2",
			"prometheis.0.tags.Created":            "TF",
			"prometheis.0.tags.For":                "Prometheus",
		}
	}

	var aliCloudArmsPrometheisCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_arms_prometheus.default",
		existMapFunc: existAliCloudArmsPrometheisDataSourceNameMapFunc,
		fakeMapFunc:  fakeAliCloudArmsPrometheisDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}

	aliCloudArmsPrometheisCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, resourceGroupIdConf, tagsConf, allConf)
}

func testAccCheckAliCloudArmsPrometheisDataSourceName(rand int, attrMap map[string]string) string {
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
  		tags = {
			Created = "TF"
    		For     = "Prometheus"
  		}
	}

	data "alicloud_arms_prometheus" "default" {
		%s
	}
`, rand, strings.Join(pairs, " \n "))
	return config
}
