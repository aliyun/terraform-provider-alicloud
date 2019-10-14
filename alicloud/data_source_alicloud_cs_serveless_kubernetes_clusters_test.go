package alicloud

import (
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform/helper/acctest"
)

func TestAccAlicloudCSServelessKubernetesClustersDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)
	resourceId := "data.alicloud_cs_serveless_kubernetes_clusters.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testaccservelessk8s-%d", rand),
		dataSourceCSServelessKubernetesClustersConfigDependence)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"enable_details": "true",
			"ids":            []string{"${alicloud_cs_serveless_kubernetes.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"enable_details": "true",
			"ids":            []string{"${alicloud_cs_serveless_kubernetes.default.id}-fake"},
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"enable_details": "true",
			"name_regex":     "${alicloud_cs_serveless_kubernetes.default.name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"enable_details": "true",
			"name_regex":     "${alicloud_cs_serveless_kubernetes.default.name}-fake",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"enable_details": "true",
			"ids":            []string{"${alicloud_cs_serveless_kubernetes.default.id}"},
			"name_regex":     "${alicloud_cs_serveless_kubernetes.default.name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"enable_details": "true",
			"ids":            []string{"${alicloud_cs_serveless_kubernetes.default.id}"},
			"name_regex":     "${alicloud_cs_serveless_kubernetes.default.name}-fake",
		}),
	}

	var existCSServelessKubernetesClustersMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                          "1",
			"ids.0":                          CHECKSET,
			"names.#":                        "1",
			"names.0":                        REGEXMATCH + fmt.Sprintf("tf-testaccservelessk8s-%d", rand),
			"clusters.#":                     "1",
			"clusters.0.id":                  CHECKSET,
			"clusters.0.name":                REGEXMATCH + fmt.Sprintf("tf-testaccservelessk8s-%d", rand),
			"clusters.0.security_group_id":   CHECKSET,
			"clusters.0.nat_gateway_id":      CHECKSET,
			"clusters.0.vpc_id":              CHECKSET,
			"clusters.0.vswitch_id":          CHECKSET,
			"clusters.0.deletion_protection": CHECKSET,
			"clusters.0.tags.%":              CHECKSET,
			"clusters.0.enndpoint_public_access_enabled": CHECKSET,
			"clusters.0.connections.%":                   CHECKSET,
		}
	}

	var fakeCSServelessKubernetesClustersMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":      "0",
			"names.#":    "0",
			"clusters.#": "0",
		}
	}

	var csServelessKubernetesClustersCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existCSServelessKubernetesClustersMapFunc,
		fakeMapFunc:  fakeCSServelessKubernetesClustersMapFunc,
	}
	preCheck := func() {
		testAccPreCheckWithRegions(t, true, connectivity.ServelessKubernetesSupportedRegions)
	}
	csServelessKubernetesClustersCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, allConf)
}

func dataSourceCSServelessKubernetesClustersConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

data "alicloud_zones" default {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "default" {
  name = "${var.name}"
  cidr_block = "10.1.0.0/21"
}

resource "alicloud_vswitch" "default" {
  name = "${var.name}"
  vpc_id = "${alicloud_vpc.default.id}"
  cidr_block = "10.1.1.0/24"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
}

resource "alicloud_cs_serveless_kubernetes" "default" {
  name_prefix = "${var.name}"
  vpc_id = "${alicloud_vpc.default.id}"
  vswitch_id = "${alicloud_vswitch.default.id}"
  new_nat_gateway = true
  enndpoint_public_access_enabled = true
  private_zone = false
  deletion_protection = false
  tags = {
		"k-aa":"v-aa"
		"k-bb":"v-aa",
  }
}
`, name)
}
