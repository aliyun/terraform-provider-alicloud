package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAliCloudRealtimeComputeSessionClusterDataSource(t *testing.T) {
	testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
	rand := acctest.RandIntRange(10000, 99999)

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudRealtimeComputeSessionClusterSourceConfig(rand),
		fakeConfig:  "",
	}

	RealtimeComputeSessionClusterCheckInfo.dataSourceTestCheck(t, rand, allConf)
}

var existRealtimeComputeSessionClusterMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":                                              "1",
		"names.#":                                            "1",
		"clusters.#":                                         "1",
		"clusters.0.id":                                      CHECKSET,
		"clusters.0.workspace":                               CHECKSET,
		"clusters.0.namespace":                               CHECKSET,
		"clusters.0.session_cluster_name":                    CHECKSET,
		"clusters.0.session_cluster_id":                      CHECKSET,
		"clusters.0.engine_version":                          CHECKSET,
		"clusters.0.deployment_target_name":                  CHECKSET,
		"clusters.0.basic_resource_setting.#":                "1",
		"clusters.0.basic_resource_setting.0.parallelism":    CHECKSET,
		"clusters.0.status.#":                                "1",
		"clusters.0.status.0.current_session_cluster_status": CHECKSET,
	}
}

var fakeRealtimeComputeSessionClusterMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":      "0",
		"names.#":    "0",
		"clusters.#": "0",
	}
}

var RealtimeComputeSessionClusterCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_realtime_compute_session_clusters.default",
	existMapFunc: existRealtimeComputeSessionClusterMapFunc,
	fakeMapFunc:  fakeRealtimeComputeSessionClusterMapFunc,
}

func testAccCheckAlicloudRealtimeComputeSessionClusterSourceConfig(rand int) string {
	return fmt.Sprintf(`
variable "name" {
	default = "tfaccrealtimecompute%d"
}

resource "alicloud_vpc" "default" {
  is_default = false
  cidr_block = "172.16.0.0/16"
  vpc_name   = "test-tf-vpc-session-cluster-ds"
}

resource "alicloud_vswitch" "default" {
  is_default   = false
  vpc_id       = alicloud_vpc.default.id
  zone_id      = "cn-beijing-g"
  cidr_block   = "172.16.0.0/24"
  vswitch_name = "test-tf-vswitch-session-cluster-ds"
}

resource "alicloud_oss_bucket" "default" {
}

resource "alicloud_realtime_compute_vvp_instance" "default" {
  vvp_instance_name = var.name
  storage {
    oss {
      bucket = alicloud_oss_bucket.default.id
    }
  }
  vpc_id      = alicloud_vpc.default.id
  vswitch_ids = ["${alicloud_vswitch.default.id}"]
  resource_spec {
    cpu       = "4"
    memory_gb = "16"
  }
  payment_type = "PayAsYouGo"
  zone_id     = alicloud_vswitch.default.zone_id
}

resource "alicloud_realtime_compute_session_cluster" "default" {
  workspace              = alicloud_realtime_compute_vvp_instance.default.resource_id
  namespace              = "${alicloud_realtime_compute_vvp_instance.default.vvp_instance_name}-default"
  session_cluster_name   = var.name
  deployment_target_name = "default-queue"
  engine_version         = "vvr-8.0.11-flink-1.17"
  basic_resource_setting {
    parallelism = 1
    jobmanager_resource_setting_spec {
      cpu    = 1
      memory = "1Gi"
    }
    taskmanager_resource_setting_spec {
      cpu    = 1
      memory = "1Gi"
    }
  }
}

data "alicloud_realtime_compute_session_clusters" "default" {
  workspace      = alicloud_realtime_compute_vvp_instance.default.resource_id
  namespace      = "${alicloud_realtime_compute_vvp_instance.default.vvp_instance_name}-default"
  ids            = [alicloud_realtime_compute_session_cluster.default.id]
  enable_details = true
}
`, rand)
}
