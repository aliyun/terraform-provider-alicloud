package alicloud

import (
	"fmt"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"testing"
)

func TestAccAlicloudCSClusterLogsDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)
	resourceId := "data.alicloud_cs_cluster_logs.default"

	testAccConfig := dataSourceTestAccConfigFunc(
		resourceId,
		fmt.Sprintf("tf-testaccinternetk8s-%d", rand),
		dataSourceCSClusterLogsConfigDependence,
	)

	logConfig := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"cluster_id": "${alicloud_cs_managed_kubernetes.default.id}",
			"type":       "log",
			"entries":    "2",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"cluster_id": "${alicloud_cs_managed_kubernetes.default.id}",
			"type":       "log",
			"entries":    "0",
		}),
	}

	eventConfig := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"cluster_id": "${alicloud_cs_managed_kubernetes.default.id}",
			"type":       "event",
			"entries":    "2",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"cluster_id": "${alicloud_cs_managed_kubernetes.default.id}-fake",
			"type":       "event",
			"entries":    "2",
		}),
	}

	taskConfig := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"cluster_id": "${alicloud_cs_managed_kubernetes.default.id}",
			"type":       "task",
			"entries":    "2",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"cluster_id": "${alicloud_cs_managed_kubernetes.default.id}",
			"type":       "task",
			"entries":    "0",
		}),
	}

	var existCSClusterLogsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"entries": "2",
		}
	}

	var fakeCSClusterLogsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"entries": "0",
		}
	}

	var csClusterLogsCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existCSClusterLogsMapFunc,
		fakeMapFunc:  fakeCSClusterLogsMapFunc,
	}

	preCheck := func() {
		testAccPreCheckWithRegions(t, true, connectivity.KubernetesSupportedRegions)
	}
	csClusterLogsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, logConfig, eventConfig, taskConfig)
}

func dataSourceCSClusterLogsConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}
data "alicloud_zones" default {
  available_resource_creation = "VSwitch"
}
data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}
data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_zones.default.zones.0.id
}
resource "alicloud_vswitch" "vswitch" {
  count             = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
  vpc_id            = data.alicloud_vpcs.default.ids.0
  cidr_block        = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 8)
  zone_id           = data.alicloud_zones.default.zones.0.id
  vswitch_name      = var.name
}
locals {
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.vswitch.*.id, [""])[0]
}
resource "alicloud_cs_managed_kubernetes" "default" {
  name_prefix                 = "${var.name}"
  cluster_spec                = "ack.pro.small"
  worker_vswitch_ids          = [local.vswitch_id]
  new_nat_gateway             = true
  node_port_range             = "30000-32767"
  password                    = "Hello1234"
  pod_cidr                    = cidrsubnet("10.0.0.0/8", 8, 35)
  service_cidr                = cidrsubnet("172.16.0.0/16", 4, 6)
  install_cloud_monitor       = true
  slb_internet_enabled        = true
  worker_disk_category        = "cloud_efficiency"
  worker_data_disk_category   = "cloud_ssd"
  worker_data_disk_size       = 200
  worker_disk_size            = 40
  worker_instance_charge_type = "PostPaid"
}
`, name)
}
