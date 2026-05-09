package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func TestAccAliCloudArmsEnvFeaturesDataSource_basic0(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)
	resourceId := "data.alicloud_arms_env_features.default"
	name := fmt.Sprintf("tf-testacc%sarmsenvfeature%d", defaultRegionToTest, rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceArmsEnvFeaturesConfig0)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"environment_id": "${alicloud_arms_env_feature.default.environment_id}",
			"ids":            []string{"${alicloud_arms_env_feature.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"environment_id": "${alicloud_arms_env_feature.default.environment_id}",
			"ids":            []string{"${alicloud_arms_env_feature.default.id}_fake"},
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"environment_id": "${alicloud_arms_env_feature.default.environment_id}",
			"name_regex":     "${alicloud_arms_env_feature.default.env_feature_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"environment_id": "${alicloud_arms_env_feature.default.environment_id}",
			"name_regex":     "${alicloud_arms_env_feature.default.env_feature_name}_fake",
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"environment_id": "${alicloud_arms_env_feature.default.environment_id}",
			"ids":            []string{"${alicloud_arms_env_feature.default.id}"},
			"name_regex":     "${alicloud_arms_env_feature.default.env_feature_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"environment_id": "${alicloud_arms_env_feature.default.environment_id}",
			"ids":            []string{"${alicloud_arms_env_feature.default.id}_fake"},
			"name_regex":     "${alicloud_arms_env_feature.default.env_feature_name}_fake",
		}),
	}
	var existAliCloudArmsEnvFeaturesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                       "1",
			"names.#":                     "1",
			"features.#":                  "1",
			"features.0.id":               CHECKSET,
			"features.0.aliyun_lang":      CHECKSET,
			"features.0.env_feature_name": CHECKSET,
			"features.0.environment_id":   CHECKSET,
			"features.0.feature_version":  CHECKSET,
			"features.0.status":           CHECKSET,
		}
	}
	var fakeAliCloudArmsEnvFeaturesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":      "0",
			"names.#":    "0",
			"features.#": "0",
		}
	}
	var alicloudArmsEnvFeaturesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_arms_env_features.default",
		existMapFunc: existAliCloudArmsEnvFeaturesDataSourceNameMapFunc,
		fakeMapFunc:  fakeAliCloudArmsEnvFeaturesDataSourceNameMapFunc,
	}
	preCheck := func() {
		testAccPreCheck(t)
	}
	alicloudArmsEnvFeaturesCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, allConf)
}

func dataSourceArmsEnvFeaturesConfig0(name string) string {
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
  		pod_cidr             = "10.131.0.0/16"
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

	resource "alicloud_arms_env_feature" "default" {
  		env_feature_name = "metric-agent"
  		environment_id   = alicloud_arms_environment.default.id
  		feature_version  = "1.1.17"
	}
`, name)
}
