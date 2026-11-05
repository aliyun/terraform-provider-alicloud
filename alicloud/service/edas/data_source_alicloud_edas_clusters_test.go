package edas_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudEdasClustersDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(1000, 9999)
	resourceId := "data.alicloud_edas_clusters.default"
	name := fmt.Sprintf("tf-testacc-edas-clusters%v", rand)

	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId, name, dataSourceEdasClustersConfigDependence)

	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"name_regex":        "${alicloud_edas_cluster.default.cluster_name}",
			"logical_region_id": os.Getenv("ALICLOUD_REGION"),
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"name_regex":        "fake_tf-testacc*",
			"logical_region_id": os.Getenv("ALICLOUD_REGION"),
		}),
	}

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids":               []string{"${alicloud_edas_cluster.default.id}"},
			"logical_region_id": os.Getenv("ALICLOUD_REGION"),
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ids":               []string{"${alicloud_edas_cluster.default.id}_fake"},
			"logical_region_id": os.Getenv("ALICLOUD_REGION"),
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids":               []string{"${alicloud_edas_cluster.default.id}"},
			"logical_region_id": os.Getenv("ALICLOUD_REGION"),
			"name_regex":        "${alicloud_edas_cluster.default.cluster_name}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ids":               []string{"${alicloud_edas_cluster.default.id}_fake"},
			"logical_region_id": os.Getenv("ALICLOUD_REGION"),
			"name_regex":        "${alicloud_edas_cluster.default.cluster_name}",
		}),
	}

	var existEdasClustersMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"clusters.#":              "1",
			"clusters.0.cluster_id":   acctest.CHECKSET,
			"clusters.0.cluster_name": name,
			"clusters.0.cluster_type": "2",
			"clusters.0.network_mode": "2",
			"clusters.0.vpc_id":       acctest.CHECKSET,
			"clusters.0.region_id":    acctest.CHECKSET,
		}
	}

	var fakeEdasClustersMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":      "0",
			"clusters.#": "0",
		}
	}

	var edasApplicationCheckInfo = acctest.DataSourceAttr{
		ResourceId:   resourceId,
		ExistMapFunc: existEdasClustersMapFunc,
		FakeMapFunc:  fakeEdasClustersMapFunc,
	}

	preCheck := func() {
		acctest.PreCheckWithRegions(t, true, connectivity.EdasSupportedRegions)
	}

	edasApplicationCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, nameRegexConf, idsConf, allConf)
}

func dataSourceEdasClustersConfigDependence(name string) string {
	return fmt.Sprintf(`
		variable "name" {
		 default = "%v"
		}

		resource "alicloud_vpc" "default" {
		  cidr_block = "172.16.0.0/12"
		  name       = "${var.name}"
		}

		resource "alicloud_edas_cluster" "default" {
		  cluster_name = "${var.name}"
		  cluster_type = "2"
		  network_mode = "2"
		  vpc_id       = "${alicloud_vpc.default.id}"
		}
		`, name)
}
