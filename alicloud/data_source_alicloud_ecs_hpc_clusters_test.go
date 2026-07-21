package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudECSHpcClustersDataSource(t *testing.T) {
	resourceId := "data.alicloud_ecs_hpc_clusters.default"
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccEcsHpcClustersTest%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceEcsHpcClustersDependence)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_ecs_hpc_cluster.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_ecs_hpc_cluster.default.id}-fake"},
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": name,
			"ids":        []string{"${alicloud_ecs_hpc_cluster.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": name,
			"ids":        []string{"${alicloud_ecs_hpc_cluster.default.id}-fake"},
		}),
	}
	var existEcsClustersMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                     "1",
			"ids.0":                     CHECKSET,
			"names.#":                   "1",
			"names.0":                   name,
			"clusters.#":                "1",
			"clusters.0.id":             CHECKSET,
			"clusters.0.hpc_cluster_id": CHECKSET,
			"clusters.0.description":    "For Terraform Test",
			"clusters.0.name":           name,
		}
	}

	var fakeEcsClustersMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":      "0",
			"names.#":    "0",
			"clusters.#": "0",
		}
	}

	var EcsClustersInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existEcsClustersMapFunc,
		fakeMapFunc:  fakeEcsClustersMapFunc,
	}

	EcsClustersInfo.dataSourceTestCheck(t, 0, idsConf, nameRegexConf)
}

func dataSourceEcsHpcClustersDependence(name string) string {
	return fmt.Sprintf(`
	resource "alicloud_ecs_hpc_cluster" "default" {
		name              = "%s"
		description       = "For Terraform Test"
	}`, name)
}
