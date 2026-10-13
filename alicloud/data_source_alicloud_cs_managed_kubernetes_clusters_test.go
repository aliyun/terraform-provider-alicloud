package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func TestAccAliCloudCSManagedKubernetesClustersDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)
	resourceId := "data.alicloud_cs_managed_kubernetes_clusters.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testaccmanagedk8s-%d", rand),
		dataSourceCSManagedKubernetesClustersConfigDependence)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"enable_details": "true",
			"ids":            []string{"${alicloud_cs_managed_kubernetes.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"enable_details": "true",
			"ids":            []string{"${alicloud_cs_managed_kubernetes.default.id}-fake"},
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"enable_details": "true",
			"name_regex":     "${alicloud_cs_managed_kubernetes.default.name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"enable_details": "true",
			"name_regex":     "${alicloud_cs_managed_kubernetes.default.name}-fake",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"enable_details": "true",
			"ids":            []string{"${alicloud_cs_managed_kubernetes.default.id}"},
			"name_regex":     "${alicloud_cs_managed_kubernetes.default.name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"enable_details": "true",
			"ids":            []string{"${alicloud_cs_managed_kubernetes.default.id}"},
			"name_regex":     "${alicloud_cs_managed_kubernetes.default.name}-fake",
		}),
	}

	var existCSManagedKubernetesClustersMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                            "1",
			"ids.0":                            CHECKSET,
			"names.#":                          "1",
			"names.0":                          REGEXMATCH + fmt.Sprintf("tf-testaccmanagedk8s-%d", rand),
			"clusters.#":                       "1",
			"clusters.0.id":                    CHECKSET,
			"clusters.0.name":                  REGEXMATCH + fmt.Sprintf("tf-testaccmanagedk8s-%d", rand),
			"clusters.0.availability_zone":     CHECKSET,
			"clusters.0.security_group_id":     CHECKSET,
			"clusters.0.vpc_id":                CHECKSET,
			"clusters.0.connections.#":         CHECKSET,
			"clusters.0.state":                 CHECKSET,
			"clusters.0.rrsa_config.0.enabled": "true",
			"clusters.0.rrsa_config.0.rrsa_oidc_issuer_url":   CHECKSET,
			"clusters.0.rrsa_config.0.ram_oidc_provider_name": CHECKSET,
			"clusters.0.rrsa_config.0.ram_oidc_provider_arn":  CHECKSET,
		}
	}

	var fakeCSManagedKubernetesClustersMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":      "0",
			"names.#":    "0",
			"clusters.#": "0",
		}
	}

	var csManagedKubernetesClustersCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existCSManagedKubernetesClustersMapFunc,
		fakeMapFunc:  fakeCSManagedKubernetesClustersMapFunc,
	}
	preCheck := func() {
		testAccPreCheckWithRegions(t, true, connectivity.KubernetesSupportedRegions)
	}
	csManagedKubernetesClustersCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, allConf)
}

func dataSourceCSManagedKubernetesClustersConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_zones" default {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "10.4.0.0/16"
}

resource "alicloud_vswitch" "vswitch" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = cidrsubnet(alicloud_vpc.default.cidr_block, 8, 8)
  zone_id      = data.alicloud_zones.default.zones.0.id
  vswitch_name = var.name
}

locals {
  vswitch_id = alicloud_vswitch.vswitch.id
}

resource "alicloud_cs_managed_kubernetes" "default" {
  name                 = var.name
  cluster_spec         = "ack.pro.small"
  worker_vswitch_ids   = [local.vswitch_id]
  new_nat_gateway      = false
  pod_cidr             = cidrsubnet("10.0.0.0/8", 8, 37)
  service_cidr         = cidrsubnet("172.17.0.0/16", 4, 7)
  slb_internet_enabled = true
  enable_rrsa          = true
}
`, name)
}
