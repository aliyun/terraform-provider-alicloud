package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudMSEClustersDataSource(t *testing.T) {
	resourceId := "data.alicloud_mse_clusters.default"
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccMseCluster-%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceMseClustersDependence)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"enable_details": "true",
			"name_regex":     "${alicloud_mse_cluster.default.cluster_alias_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"enable_details": "true",
			"name_regex":     "${alicloud_mse_cluster.default.cluster_alias_name}-fake",
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"enable_details": "true",
			"ids":            []string{"${alicloud_mse_cluster.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"enable_details": "true",
			"ids":            []string{"${alicloud_mse_cluster.default.id}-fake"},
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alicloud_mse_cluster.default.id}"},
			"status":         "INIT_SUCCESS",
			"enable_details": "true",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alicloud_mse_cluster.default.id}"},
			"status":         "DESTROY_FAILED",
			"enable_details": "true",
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"enable_details": "true",
			"ids":            []string{"${alicloud_mse_cluster.default.id}"},
			"status":         "INIT_SUCCESS",
			"name_regex":     "${alicloud_mse_cluster.default.cluster_alias_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"enable_details": "true",
			"ids":            []string{"${alicloud_mse_cluster.default.id}-fake"},
			"status":         "INIT_SUCCESS",
			"name_regex":     "${alicloud_mse_cluster.default.cluster_alias_name}",
		}),
	}
	var existMseClusterMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                       "1",
			"ids.0":                       CHECKSET,
			"names.#":                     "1",
			"names.0":                     name,
			"clusters.#":                  "1",
			"clusters.0.app_version":      CHECKSET,
			"clusters.0.cluster_name":     name,
			"clusters.0.cluster_id":       CHECKSET,
			"clusters.0.cluster_type":     "Nacos-Ans",
			"clusters.0.id":               CHECKSET,
			"clusters.0.instance_id":      CHECKSET,
			"clusters.0.internet_address": CHECKSET,
			"clusters.0.intranet_address": CHECKSET,
			"clusters.0.internet_domain":  CHECKSET,
			"clusters.0.intranet_domain":  CHECKSET,
			"clusters.0.acl_id":           CHECKSET,
			"clusters.0.health_status":    CHECKSET,
			"clusters.0.init_cost_time":   CHECKSET,
			"clusters.0.instance_count":   "1",
			"clusters.0.internet_port":    CHECKSET,
			"clusters.0.intranet_port":    CHECKSET,
			"clusters.0.memory_capacity":  CHECKSET,
			"clusters.0.pay_info":         "按量付费",
			"clusters.0.pub_network_flow": "1",
			"clusters.0.status":           "INIT_SUCCESS",
			"clusters.0.cpu":              CHECKSET,
		}
	}

	var fakeMseClusterMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":      "0",
			"names.#":    "0",
			"clusters.#": "0",
		}
	}

	var mseClustersInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existMseClusterMapFunc,
		fakeMapFunc:  fakeMseClusterMapFunc,
	}
	mseClustersInfo.dataSourceTestCheck(t, 0, nameRegexConf, idsConf, statusConf, allConf)
}

func dataSourceMseClustersDependence(name string) string {
	return fmt.Sprintf(`
	data "alicloud_vpcs" "default" {
	  name_regex = "default-NODELETING"
	}
	data "alicloud_vswitches" "default" {
	  vpc_id = data.alicloud_vpcs.default.ids.0
	}
	
	resource "alicloud_mse_cluster" "default" {
	  cluster_specification = "MSE_SC_1_2_200_c"
	  cluster_type = "Nacos-Ans"
	  cluster_version = "NACOS_ANS_1_2_1"
	  instance_count = 1
	  net_type = "privatenet"
	  vswitch_id = data.alicloud_vswitches.default.ids.0
	  pub_network_flow = "1"
	  acl_entry_list= ["127.0.0.1/32"]
	  cluster_alias_name= "%s"
	}
	`, name)
}
