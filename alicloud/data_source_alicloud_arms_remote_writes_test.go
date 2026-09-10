package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

// data source alicloud_arms_remote_writes has been deprecated from version 1.228.0
func SkipTestAccAliCloudArmsRemoteWritesDataSource(t *testing.T) {
	rand := acctest.RandIntRange(0, 10)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudArmsRemoteWritesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_arms_remote_write.default.id}"]`,
		}),
		fakeConfig: testAccCheckAliCloudArmsRemoteWritesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_arms_remote_write.default.id}_fake"]`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudArmsRemoteWritesDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_arms_remote_write.default.remote_write_name}"`,
		}),
		fakeConfig: testAccCheckAliCloudArmsRemoteWritesDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_arms_remote_write.default.remote_write_name}_fake"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudArmsRemoteWritesDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_arms_remote_write.default.id}"]`,
			"name_regex": `"${alicloud_arms_remote_write.default.remote_write_name}"`,
		}),
		fakeConfig: testAccCheckAliCloudArmsRemoteWritesDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_arms_remote_write.default.id}_fake"]`,
			"name_regex": `"${alicloud_arms_remote_write.default.remote_write_name}_fake"`,
		}),
	}
	var existAliCloudArmsRemoteWritesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                             "1",
			"names.#":                           "1",
			"remote_writes.#":                   "1",
			"remote_writes.0.id":                CHECKSET,
			"remote_writes.0.cluster_id":        CHECKSET,
			"remote_writes.0.remote_write_name": "ArmsRemoteWrite",
			"remote_writes.0.remote_write_yaml": "remote_write:\n- name: ArmsRemoteWrite\n  url: http://47.96.227.137:8080/prometheus/xxx/yyy/cn-hangzhou/api/v3/write\n  basic_auth: {username: 666, password: '******'}\n  write_relabel_configs:\n  - source_labels: [instance_id]\n    separator: ;\n    regex: si-6e2ca86444db4e55a7c1\n    replacement: $1\n    action: keep\n",
		}
	}
	var fakeAliCloudArmsRemoteWritesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":           "0",
			"names.#":         "0",
			"remote_writes.#": "0",
		}
	}
	var alicloudArmsRemoteWritesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_arms_remote_writes.default",
		existMapFunc: existAliCloudArmsRemoteWritesDataSourceNameMapFunc,
		fakeMapFunc:  fakeAliCloudArmsRemoteWritesDataSourceNameMapFunc,
	}
	preCheck := func() {
		testAccPreCheck(t)
	}
	alicloudArmsRemoteWritesCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, allConf)
}

func testAccCheckAliCloudArmsRemoteWritesDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
	variable "name" {
  		default = "tf-testacc-ArmsRW-%d"
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
  		resource_group_id   = data.alicloud_resource_manager_resource_groups.default.groups.0.id
	}

	resource "alicloud_arms_remote_write" "default" {
  		cluster_id        = alicloud_arms_prometheus.default.id
  		remote_write_yaml = "remote_write:\n- name: ArmsRemoteWrite\n  url: http://47.96.227.137:8080/prometheus/xxx/yyy/cn-hangzhou/api/v3/write\n  basic_auth: {username: 666, password: '******'}\n  write_relabel_configs:\n  - source_labels: [instance_id]\n    separator: ;\n    regex: si-6e2ca86444db4e55a7c1\n    replacement: $1\n    action: keep\n"
	}

	data "alicloud_arms_remote_writes" "default" {
  		cluster_id = alicloud_arms_remote_write.default.cluster_id
		%s
	}
`, rand, strings.Join(pairs, " \n "))
	return config
}
