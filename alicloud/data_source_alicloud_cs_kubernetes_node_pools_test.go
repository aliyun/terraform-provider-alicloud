package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudAckNodepoolDataSource(t *testing.T) {
	testAccPreCheckWithRegions(t, true, []connectivity.Region{"eu-central-1"})
	rand := acctest.RandIntRange(1000000, 9999999)

	ClusterIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAckNodepoolSourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_cs_kubernetes_node_pool.default.node_pool_id}"]`,
			"cluster_id": `"${alicloud_cs_managed_kubernetes.default.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudAckNodepoolSourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_cs_kubernetes_node_pool.default.node_pool_id}_fake"]`,
			"cluster_id": `"${alicloud_cs_managed_kubernetes.default.id}"`,
		}),
	}

	NodePoolNameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAckNodepoolSourceConfig(rand, map[string]string{
			"node_pool_name": `"${var.name}"`,
			"cluster_id":     `"${alicloud_cs_managed_kubernetes.default.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudAckNodepoolSourceConfig(rand, map[string]string{
			"node_pool_name": `"${var.name}_fake"`,
			"cluster_id":     `"${alicloud_cs_managed_kubernetes.default.id}"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAckNodepoolSourceConfig(rand, map[string]string{
			"ids":            `["${alicloud_cs_kubernetes_node_pool.default.node_pool_id}"]`,
			"node_pool_name": `"${var.name}"`,
			"cluster_id":     `"${alicloud_cs_managed_kubernetes.default.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudAckNodepoolSourceConfig(rand, map[string]string{
			"ids":            `["${alicloud_cs_kubernetes_node_pool.default.node_pool_id}_fake"]`,
			"node_pool_name": `"${var.name}_fake"`,
			"cluster_id":     `"${alicloud_cs_managed_kubernetes.default.id}"`,
		}),
	}

	AckNodepoolCheckInfo.dataSourceTestCheck(t, rand, ClusterIdConf, NodePoolNameConf, allConf)
}

var existAckNodepoolMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"nodepools.#":                              "1",
		"nodepools.0.spot_instance_remedy":         CHECKSET,
		"nodepools.0.soc_enabled":                  CHECKSET,
		"nodepools.0.data_disks.#":                 CHECKSET,
		"nodepools.0.platform":                     CHECKSET,
		"nodepools.0.instance_charge_type":         CHECKSET,
		"nodepools.0.runtime_name":                 CHECKSET,
		"nodepools.0.system_disk_provisioned_iops": CHECKSET,
		"nodepools.0.image_type":                   CHECKSET,
		"nodepools.0.tee_config.#":                 CHECKSET,
		"nodepools.0.ram_role_name":                CHECKSET,
		"nodepools.0.private_pool_options.#":       CHECKSET,
		"nodepools.0.node_name_mode":               CHECKSET,
		"nodepools.0.image_id":                     CHECKSET,
		"nodepools.0.install_cloud_monitor":        CHECKSET,
		"nodepools.0.tags.%":                       CHECKSET,
		"nodepools.0.multi_az_policy":              CHECKSET,
		"nodepools.0.cpu_policy":                   CHECKSET,
		"nodepools.0.node_pool_name":               CHECKSET,
		"nodepools.0.scaling_group_id":             CHECKSET,
		"nodepools.0.period":                       CHECKSET,
		"nodepools.0.runtime_version":              CHECKSET,
		"nodepools.0.spot_instance_pools":          CHECKSET,
		"nodepools.0.labels.#":                     CHECKSET,
		"nodepools.0.security_group_ids.#":         CHECKSET,
		"nodepools.0.taints.#":                     CHECKSET,
		"nodepools.0.internet_max_bandwidth_out":   CHECKSET,
		"nodepools.0.login_as_non_root":            CHECKSET,
		"nodepools.0.compensate_with_on_demand":    CHECKSET,
		"nodepools.0.system_disk_size":             CHECKSET,
		"nodepools.0.auto_renew":                   CHECKSET,
		"nodepools.0.system_disk_encrypted":        CHECKSET,
		"nodepools.0.security_hardening_os":        CHECKSET,
		"nodepools.0.system_disk_categories.#":     CHECKSET,
		"nodepools.0.vswitch_ids.#":                CHECKSET,
		"nodepools.0.spot_price_limit.#":           CHECKSET,
		"nodepools.0.instance_types.#":             CHECKSET,
		"nodepools.0.unschedulable":                CHECKSET,
		"nodepools.0.spot_strategy":                CHECKSET,
		"nodepools.0.auto_renew_period":            CHECKSET,
		"nodepools.0.scaling_policy":               CHECKSET,
		"nodepools.0.scaling_config.#":             CHECKSET,
		"nodepools.0.security_group_id":            CHECKSET,
		"nodepools.0.management.#":                 CHECKSET,
		"nodepools.0.cis_enabled":                  CHECKSET,
		"nodepools.0.system_disk_category":         CHECKSET,
		"nodepools.0.key_name":                     CHECKSET,
		"nodepools.0.system_disk_bursting_enabled": CHECKSET,
		"nodepools.0.rds_instances.#":              CHECKSET,
		"nodepools.0.kubelet_configuration.#":      CHECKSET,
		"nodepools.0.node_pool_id":                 CHECKSET,
	}
}

var fakeAckNodepoolMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"nodepools.#": "0",
	}
}

var AckNodepoolCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_cs_kubernetes_node_pools.default",
	existMapFunc: existAckNodepoolMapFunc,
	fakeMapFunc:  fakeAckNodepoolMapFunc,
}

func testAccCheckAlicloudAckNodepoolSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccAckNodepool%d"
}

data "alicloud_enhanced_nat_available_zones" "enhanced" {
}

data "alicloud_instance_types" "cloud_efficiency" {
  availability_zone    = data.alicloud_enhanced_nat_available_zones.enhanced.zones.0.zone_id
  cpu_core_count       = 4
  memory_size          = 8
  kubernetes_node_role = "Worker"
  system_disk_category = "cloud_efficiency"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "10.4.0.0/16"
}
resource "alicloud_vswitch" "default" {
  vswitch_name = var.name
  cidr_block   = "10.4.0.0/24"
  vpc_id       = alicloud_vpc.default.id
  zone_id      = data.alicloud_enhanced_nat_available_zones.enhanced.zones.0.zone_id
}

resource "alicloud_cs_managed_kubernetes" "default" {
  name_prefix          = var.name
  cluster_spec         = "ack.pro.small"
  worker_vswitch_ids   = [alicloud_vswitch.default.id]
  new_nat_gateway      = true
  pod_cidr             = cidrsubnet("10.0.0.0/8", 8, 36)
  service_cidr         = cidrsubnet("172.16.0.0/16", 4, 7)
  slb_internet_enabled = true
  enable_rrsa          = true
}

resource "alicloud_key_pair" "default" {
  key_pair_name = var.name
}

resource "alicloud_cs_kubernetes_node_pool" "default" {
  node_pool_name       = var.name
  cluster_id           = alicloud_cs_managed_kubernetes.default.id
  vswitch_ids          = [alicloud_vswitch.default.id]
  instance_types       = [data.alicloud_instance_types.cloud_efficiency.instance_types.0.id]
  system_disk_category = "cloud_efficiency"
  system_disk_size     = 40
  key_name             = alicloud_key_pair.default.key_pair_name

  # automatic scaling node pool configuration.
  scaling_config {
    min_size = 1
    max_size = 10
    type     = "spot"
  }
  # spot price config
  spot_strategy = "SpotWithPriceLimit"
  spot_price_limit {
    instance_type = data.alicloud_instance_types.cloud_efficiency.instance_types.0.id
    price_limit   = "0.70"
  }
}

data "alicloud_cs_kubernetes_node_pools" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
