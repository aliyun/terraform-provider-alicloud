package cs_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudAckNodepoolDataSource(t *testing.T) {
	// eu-central-1 no longer has 4c8g worker instance stock; use cn-hangzhou,
	// the same proven region as the other ACK nodepool tests
	acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	rand := tfacctest.RandIntRange(1000000, 9999999)

	ClusterIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudAckNodepoolSourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_cs_kubernetes_node_pool.default.node_pool_id}"]`,
			"cluster_id": `"${alicloud_cs_managed_kubernetes.default.id}"`,
		}),
		FakeConfig: testAccCheckAlicloudAckNodepoolSourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_cs_kubernetes_node_pool.default.node_pool_id}_fake"]`,
			"cluster_id": `"${alicloud_cs_managed_kubernetes.default.id}"`,
		}),
	}

	NodePoolNameConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudAckNodepoolSourceConfig(rand, map[string]string{
			"node_pool_name": `"${var.name}"`,
			"cluster_id":     `"${alicloud_cs_managed_kubernetes.default.id}"`,
		}),
		FakeConfig: testAccCheckAlicloudAckNodepoolSourceConfig(rand, map[string]string{
			"node_pool_name": `"${var.name}_fake"`,
			"cluster_id":     `"${alicloud_cs_managed_kubernetes.default.id}"`,
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudAckNodepoolSourceConfig(rand, map[string]string{
			"ids":            `["${alicloud_cs_kubernetes_node_pool.default.node_pool_id}"]`,
			"node_pool_name": `"${var.name}"`,
			"cluster_id":     `"${alicloud_cs_managed_kubernetes.default.id}"`,
		}),
		FakeConfig: testAccCheckAlicloudAckNodepoolSourceConfig(rand, map[string]string{
			"ids":            `["${alicloud_cs_kubernetes_node_pool.default.node_pool_id}_fake"]`,
			"node_pool_name": `"${var.name}_fake"`,
			"cluster_id":     `"${alicloud_cs_managed_kubernetes.default.id}"`,
		}),
	}

	AckNodepoolCheckInfo.DataSourceTestCheck(t, rand, ClusterIdConf, NodePoolNameConf, allConf)
}

var existAckNodepoolMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"nodepools.#":                              "1",
		"nodepools.0.spot_instance_remedy":         acctest.CHECKSET,
		"nodepools.0.soc_enabled":                  acctest.CHECKSET,
		"nodepools.0.data_disks.#":                 acctest.CHECKSET,
		"nodepools.0.platform":                     acctest.CHECKSET,
		"nodepools.0.instance_charge_type":         acctest.CHECKSET,
		"nodepools.0.runtime_name":                 acctest.CHECKSET,
		"nodepools.0.system_disk_provisioned_iops": acctest.CHECKSET,
		"nodepools.0.image_type":                   acctest.CHECKSET,
		"nodepools.0.tee_config.#":                 acctest.CHECKSET,
		"nodepools.0.ram_role_name":                acctest.CHECKSET,
		"nodepools.0.private_pool_options.#":       acctest.CHECKSET,
		"nodepools.0.node_name_mode":               acctest.CHECKSET,
		"nodepools.0.image_id":                     acctest.CHECKSET,
		"nodepools.0.install_cloud_monitor":        acctest.CHECKSET,
		"nodepools.0.tags.%":                       acctest.CHECKSET,
		"nodepools.0.multi_az_policy":              acctest.CHECKSET,
		"nodepools.0.cpu_policy":                   acctest.CHECKSET,
		"nodepools.0.node_pool_name":               acctest.CHECKSET,
		"nodepools.0.scaling_group_id":             acctest.CHECKSET,
		"nodepools.0.period":                       acctest.CHECKSET,
		"nodepools.0.runtime_version":              acctest.CHECKSET,
		"nodepools.0.spot_instance_pools":          acctest.CHECKSET,
		"nodepools.0.labels.#":                     acctest.CHECKSET,
		"nodepools.0.security_group_ids.#":         acctest.CHECKSET,
		"nodepools.0.taints.#":                     acctest.CHECKSET,
		"nodepools.0.internet_max_bandwidth_out":   acctest.CHECKSET,
		"nodepools.0.login_as_non_root":            acctest.CHECKSET,
		"nodepools.0.compensate_with_on_demand":    acctest.CHECKSET,
		"nodepools.0.system_disk_size":             acctest.CHECKSET,
		"nodepools.0.auto_renew":                   acctest.CHECKSET,
		"nodepools.0.system_disk_encrypted":        acctest.CHECKSET,
		"nodepools.0.security_hardening_os":        acctest.CHECKSET,
		"nodepools.0.system_disk_categories.#":     acctest.CHECKSET,
		"nodepools.0.vswitch_ids.#":                acctest.CHECKSET,
		"nodepools.0.spot_price_limit.#":           acctest.CHECKSET,
		"nodepools.0.instance_types.#":             acctest.CHECKSET,
		"nodepools.0.unschedulable":                acctest.CHECKSET,
		"nodepools.0.spot_strategy":                acctest.CHECKSET,
		"nodepools.0.auto_renew_period":            acctest.CHECKSET,
		"nodepools.0.scaling_policy":               acctest.CHECKSET,
		"nodepools.0.scaling_config.#":             acctest.CHECKSET,
		"nodepools.0.security_group_id":            acctest.CHECKSET,
		"nodepools.0.management.#":                 acctest.CHECKSET,
		"nodepools.0.cis_enabled":                  acctest.CHECKSET,
		"nodepools.0.system_disk_category":         acctest.CHECKSET,
		"nodepools.0.key_name":                     acctest.CHECKSET,
		"nodepools.0.system_disk_bursting_enabled": acctest.CHECKSET,
		"nodepools.0.rds_instances.#":              acctest.CHECKSET,
		"nodepools.0.kubelet_configuration.#":      acctest.CHECKSET,
		"nodepools.0.node_pool_id":                 acctest.CHECKSET,
	}
}

var fakeAckNodepoolMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"nodepools.#": "0",
	}
}

var AckNodepoolCheckInfo = acctest.DataSourceAttr{
	ResourceId:   "data.alicloud_cs_kubernetes_node_pools.default",
	ExistMapFunc: existAckNodepoolMapFunc,
	FakeMapFunc:  fakeAckNodepoolMapFunc,
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

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

data "alicloud_instance_types" "cloud_efficiency" {
  availability_zone    = data.alicloud_zones.default.zones.0.id
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
  zone_id      = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_cs_managed_kubernetes" "default" {
  name_prefix          = var.name
  cluster_spec         = "ack.pro.small"
  worker_vswitch_ids   = [alicloud_vswitch.default.id]
  new_nat_gateway      = false
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
