package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
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

	data "alicloud_vpcs" "default" {
  		name_regex = "^default-NODELETING$"
	}

	data "alicloud_vswitches" "default" {
  		vpc_id = data.alicloud_vpcs.default.ids.0
	}

	data "alicloud_cs_managed_kubernetes_clusters" "default" {
  		name_regex = "^Default"
	}

	resource "alicloud_cs_managed_kubernetes" "default" {
  		count                = length(data.alicloud_cs_managed_kubernetes_clusters.default.ids) > 0 ? 0 : 1
  		name                 = var.name
  		cluster_spec         = "ack.pro.small"
  		worker_vswitch_ids   = [data.alicloud_vswitches.default.ids.0]
  		new_nat_gateway      = false
  		pod_cidr             = "10.130.0.0/16"
  		service_cidr         = "192.168.0.0/16"
  		slb_internet_enabled = true
  		is_enterprise_security_group = true
	}

	locals {
  		cluster_id = length(data.alicloud_cs_managed_kubernetes_clusters.default.ids) > 0 ? data.alicloud_cs_managed_kubernetes_clusters.default.ids.0 : alicloud_cs_managed_kubernetes.default.0.id
	}

	resource "alicloud_arms_environment" "default" {
  		environment_type     = "CS"
  		environment_name     = var.name
  		bind_resource_id     = local.cluster_id
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
