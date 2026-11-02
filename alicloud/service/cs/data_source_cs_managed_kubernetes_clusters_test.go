package cs_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudCSManagedKubernetesClustersDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(1000000, 9999999)
	resourceId := "data.alicloud_cs_managed_kubernetes_clusters.default"

	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testaccmanagedk8s-%d", rand),
		dataSourceCSManagedKubernetesClustersConfigDependence)

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"enable_details": "true",
			"ids":            []string{"${alicloud_cs_managed_kubernetes.default.id}"},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"enable_details": "true",
			"ids":            []string{"${alicloud_cs_managed_kubernetes.default.id}-fake"},
		}),
	}

	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"enable_details": "true",
			"name_regex":     "${alicloud_cs_managed_kubernetes.default.name}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"enable_details": "true",
			"name_regex":     "${alicloud_cs_managed_kubernetes.default.name}-fake",
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"enable_details": "true",
			"ids":            []string{"${alicloud_cs_managed_kubernetes.default.id}"},
			"name_regex":     "${alicloud_cs_managed_kubernetes.default.name}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"enable_details": "true",
			"ids":            []string{"${alicloud_cs_managed_kubernetes.default.id}"},
			"name_regex":     "${alicloud_cs_managed_kubernetes.default.name}-fake",
		}),
	}

	var existCSManagedKubernetesClustersMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                            "1",
			"ids.0":                            acctest.CHECKSET,
			"names.#":                          "1",
			"names.0":                          acctest.REGEXMATCH + fmt.Sprintf("tf-testaccmanagedk8s-%d", rand),
			"clusters.#":                       "1",
			"clusters.0.id":                    acctest.CHECKSET,
			"clusters.0.name":                  acctest.REGEXMATCH + fmt.Sprintf("tf-testaccmanagedk8s-%d", rand),
			"clusters.0.availability_zone":     acctest.CHECKSET,
			"clusters.0.security_group_id":     acctest.CHECKSET,
			"clusters.0.vpc_id":                acctest.CHECKSET,
			"clusters.0.connections.#":         acctest.CHECKSET,
			"clusters.0.state":                 acctest.CHECKSET,
			"clusters.0.rrsa_config.0.enabled": "true",
			"clusters.0.rrsa_config.0.rrsa_oidc_issuer_url":   acctest.CHECKSET,
			"clusters.0.rrsa_config.0.ram_oidc_provider_name": acctest.CHECKSET,
			"clusters.0.rrsa_config.0.ram_oidc_provider_arn":  acctest.CHECKSET,
		}
	}

	var fakeCSManagedKubernetesClustersMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":      "0",
			"names.#":    "0",
			"clusters.#": "0",
		}
	}

	var csManagedKubernetesClustersCheckInfo = acctest.DataSourceAttr{
		ResourceId:   resourceId,
		ExistMapFunc: existCSManagedKubernetesClustersMapFunc,
		FakeMapFunc:  fakeCSManagedKubernetesClustersMapFunc,
	}
	preCheck := func() {
		acctest.PreCheckWithRegions(t, true, connectivity.KubernetesSupportedRegions)
	}
	csManagedKubernetesClustersCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, allConf)
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
