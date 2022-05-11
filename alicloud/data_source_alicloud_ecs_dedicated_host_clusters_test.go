package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudECSDedicatedHostClustersDataSource(t *testing.T) {
	resourceId := "data.alicloud_ecs_dedicated_host_clusters.default"
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc-ecsdedicatedhostcluster-%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceEcsDedicatedHostClustersDependence)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_ecs_dedicated_host_cluster.default.dedicated_host_cluster_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_ecs_dedicated_host_cluster.default.dedicated_host_cluster_name}-fake",
		}),
	}
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_ecs_dedicated_host_cluster.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_ecs_dedicated_host_cluster.default.id}-fake"},
		}),
	}
	clusterIdsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"dedicated_host_cluster_ids": []string{"${alicloud_ecs_dedicated_host_cluster.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"dedicated_host_cluster_ids": []string{"${alicloud_ecs_dedicated_host_cluster.default.id}-fake"},
		}),
	}
	clusterNameConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"dedicated_host_cluster_name": "${alicloud_ecs_dedicated_host_cluster.default.dedicated_host_cluster_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"dedicated_host_cluster_name": "${alicloud_ecs_dedicated_host_cluster.default.dedicated_host_cluster_name}-fake",
		}),
	}
	zoneIdConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":     []string{"${alicloud_ecs_dedicated_host_cluster.default.id}"},
			"zone_id": "${alicloud_ecs_dedicated_host_cluster.default.zone_id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":     []string{"${alicloud_ecs_dedicated_host_cluster.default.id}"},
			"zone_id": "${alicloud_ecs_dedicated_host_cluster.default.zone_id}-fake",
		}),
	}
	tagsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_ecs_dedicated_host_cluster.default.id}"},
			"tags": map[string]string{
				"Create": "TF",
				"For":    "DDH_Cluster_Test",
			},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_ecs_dedicated_host_cluster.default.id}"},
			"tags": map[string]string{
				"Create": "DDH_Cluster_Test",
				"For":    "TF",
			},
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex":                  "${alicloud_ecs_dedicated_host_cluster.default.dedicated_host_cluster_name}",
			"dedicated_host_cluster_name": "${alicloud_ecs_dedicated_host_cluster.default.dedicated_host_cluster_name}",
			"dedicated_host_cluster_ids":  []string{"${alicloud_ecs_dedicated_host_cluster.default.id}"},
			"ids":                         []string{"${alicloud_ecs_dedicated_host_cluster.default.id}"},
			"zone_id":                     "${alicloud_ecs_dedicated_host_cluster.default.zone_id}",
			"tags": map[string]string{
				"Create": "TF",
				"For":    "DDH_Cluster_Test",
			},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex":                  "${alicloud_ecs_dedicated_host_cluster.default.dedicated_host_cluster_name}-fake",
			"dedicated_host_cluster_name": "${alicloud_ecs_dedicated_host_cluster.default.dedicated_host_cluster_name}-fake",
			"dedicated_host_cluster_ids":  []string{"${alicloud_ecs_dedicated_host_cluster.default.id}-fake"},
			"ids":                         []string{"${alicloud_ecs_dedicated_host_cluster.default.id}"},
			"zone_id":                     "${alicloud_ecs_dedicated_host_cluster.default.zone_id}-fake",
			"tags": map[string]string{
				"Create": "DDH_Cluster_Test",
				"For":    "TF",
			},
		}),
	}
	var existEcsDedicatedHostClusterMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                        "1",
			"ids.0":                                        CHECKSET,
			"names.#":                                      "1",
			"clusters.#":                                   "1",
			"clusters.0.dedicated_host_cluster_name":       fmt.Sprintf("tf-testacc-ecsdedicatedhostcluster-%d", rand),
			"clusters.0.description":                       fmt.Sprintf("tf-testacc-ecsdedicatedhostcluster-%d", rand),
			"clusters.0.dedicated_host_cluster_id":         CHECKSET,
			"clusters.0.id":                                CHECKSET,
			"clusters.0.resource_group_id":                 "",
			"clusters.0.zone_id":                           CHECKSET,
			"clusters.0.tags.Create":                       "TF",
			"clusters.0.tags.For":                          "DDH_Cluster_Test",
			"clusters.0.dedicated_host_ids.#":              "0",
			"clusters.0.dedicated_host_cluster_capacity.#": "1",
			"clusters.0.dedicated_host_cluster_capacity.0.total_memory":               "0",
			"clusters.0.dedicated_host_cluster_capacity.0.available_memory":           "0",
			"clusters.0.dedicated_host_cluster_capacity.0.available_vcpus":            "0",
			"clusters.0.dedicated_host_cluster_capacity.0.total_vcpus":                "0",
			"clusters.0.dedicated_host_cluster_capacity.0.local_storage_capacities.#": "0",
		}
	}

	var fakeEcsDedicatedHostClusterMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":      "0",
			"names.#":    "0",
			"clusters.#": "0",
		}
	}

	var EcsDedicatedHostClusterCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existEcsDedicatedHostClusterMapFunc,
		fakeMapFunc:  fakeEcsDedicatedHostClusterMapFunc,
	}

	EcsDedicatedHostClusterCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, idsConf, clusterIdsConf, clusterNameConf, zoneIdConf, tagsConf, allConf)
}

func dataSourceEcsDedicatedHostClustersDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_zones" default {}

resource "alicloud_ecs_dedicated_host_cluster" "default" {
  dedicated_host_cluster_name = var.name
  description                 = var.name
  zone_id                     = data.alicloud_zones.default.zones.0.id
  tags                        = {
    Create = "TF"
    For    = "DDH_Cluster_Test",
  }
}`, name)
}
