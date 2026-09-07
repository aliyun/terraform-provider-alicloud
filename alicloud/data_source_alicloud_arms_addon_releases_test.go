package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAliCloudArmsAddonReleasesDataSource_basic0(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)
	resourceId := "data.alicloud_arms_addon_releases.default"
	name := fmt.Sprintf("tf-testacc%sarmsaddonrelease%d", defaultRegionToTest, rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceArmsAddonReleasesConfig0)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"environment_id": "${alicloud_arms_addon_release.default.environment_id}",
			"ids":            []string{"${alicloud_arms_addon_release.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"environment_id": "${alicloud_arms_addon_release.default.environment_id}",
			"ids":            []string{"${alicloud_arms_addon_release.default.id}_fake"},
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"environment_id": "${alicloud_arms_addon_release.default.environment_id}",
			"name_regex":     "${alicloud_arms_addon_release.default.addon_release_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"environment_id": "${alicloud_arms_addon_release.default.environment_id}",
			"name_regex":     "${alicloud_arms_addon_release.default.addon_release_name}_fake",
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"environment_id": "${alicloud_arms_addon_release.default.environment_id}",
			"ids":            []string{"${alicloud_arms_addon_release.default.id}"},
			"name_regex":     "${alicloud_arms_addon_release.default.addon_release_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"environment_id": "${alicloud_arms_addon_release.default.environment_id}",
			"ids":            []string{"${alicloud_arms_addon_release.default.id}_fake"},
			"name_regex":     "${alicloud_arms_addon_release.default.addon_release_name}_fake",
		}),
	}
	var existAliCloudArmsAddonReleasesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                         "1",
			"names.#":                       "1",
			"releases.#":                    "1",
			"releases.0.id":                 CHECKSET,
			"releases.0.addon_name":         CHECKSET,
			"releases.0.addon_release_name": CHECKSET,
			"releases.0.addon_version":      CHECKSET,
			"releases.0.alert_rule_count":   CHECKSET,
			"releases.0.aliyun_lang":        CHECKSET,
			"releases.0.create_time":        CHECKSET,
			"releases.0.dashboard_count":    CHECKSET,
			"releases.0.environment_id":     CHECKSET,
			"releases.0.exporter_count":     CHECKSET,
			"releases.0.region_id":          CHECKSET,
		}
	}
	var fakeAliCloudArmsAddonReleasesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":      "0",
			"names.#":    "0",
			"releases.#": "0",
		}
	}
	var alicloudArmsAddonReleasesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_arms_addon_releases.default",
		existMapFunc: existAliCloudArmsAddonReleasesDataSourceNameMapFunc,
		fakeMapFunc:  fakeAliCloudArmsAddonReleasesDataSourceNameMapFunc,
	}
	preCheck := func() {
		testAccPreCheck(t)
	}
	alicloudArmsAddonReleasesCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, allConf)
}

func dataSourceArmsAddonReleasesConfig0(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_enhanced_nat_available_zones" "enhanced" {
	}

	resource "alicloud_vpc" "vpc" {
  		description = "api-resource-test1-hz"
  		cidr_block  = "192.168.0.0/16"
  		vpc_name    = var.name
	}

	resource "alicloud_vswitch" "vswitch" {
  		description  = "api-resource-test1-hz"
  		vpc_id       = alicloud_vpc.vpc.id
  		vswitch_name = var.name
  		zone_id    = data.alicloud_enhanced_nat_available_zones.enhanced.zones.0.zone_id
  		cidr_block = cidrsubnet(alicloud_vpc.vpc.cidr_block, 8, 8)
	}


	resource "alicloud_snapshot_policy" "default" {
  		name            = var.name
  		repeat_weekdays = ["1", "2", "3"]
  		retention_days  = -1
  		time_points     = ["1", "22", "23"]
	}

	data "alicloud_instance_types" "default" {
  		availability_zone    = alicloud_vswitch.vswitch.zone_id
  		cpu_core_count       = 2
  		memory_size          = 4
  		kubernetes_node_role = "Worker"
  		instance_type_family = "ecs.sn1ne"
	}

	resource "alicloud_cs_managed_kubernetes" "default" {
  		name               = var.name
  		cluster_spec       = "ack.pro.small"
  		version            = "1.24.6-aliyun.1"
  		new_nat_gateway    = true
  		node_cidr_mask     = 26
  		proxy_mode         = "ipvs"
  		service_cidr       = "172.23.0.0/16"
  		pod_cidr           = "10.95.0.0/16"
  		worker_vswitch_ids = [alicloud_vswitch.vswitch.id]
	}

	resource "alicloud_key_pair" "default" {
  		key_pair_name = var.name
	}

	resource "alicloud_cs_kubernetes_node_pool" "default" {
  		name                 = "desired_size"
  		cluster_id           = alicloud_cs_managed_kubernetes.default.id
  		vswitch_ids          = [alicloud_vswitch.vswitch.id]
  		instance_types       = [data.alicloud_instance_types.default.instance_types.0.id]
  		system_disk_category = "cloud_efficiency"
  		system_disk_size     = 40
  		key_name             = alicloud_key_pair.default.key_pair_name
  		desired_size         = 2
	}

	resource "alicloud_arms_environment" "default" {
  		environment_type     = "CS"
  		environment_name     = var.name
  		bind_resource_id     = alicloud_cs_kubernetes_node_pool.default.cluster_id
  		environment_sub_type = "ManagedKubernetes"
	}

	resource "alicloud_arms_addon_release" "default" {
  		aliyun_lang    = "zh"
  		addon_name     = "mysql"
  		environment_id = alicloud_arms_environment.default.id
  		addon_version  = "0.0.2"
  		values = jsonencode(
    		{
      			host     = "mysql-service.default"
      			password = "roots"
      			port     = 3306
      			username = "root"
    		}
  		)
	}
`, name)
}
