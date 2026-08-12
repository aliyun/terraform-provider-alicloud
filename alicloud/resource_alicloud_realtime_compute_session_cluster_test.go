package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test RealtimeCompute SessionCluster. >>> Resource test cases, automatically generated.
// Case SessionCluster lifecycle test
func TestAccAliCloudRealtimeComputeSessionCluster_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_realtime_compute_session_cluster.default"
	ra := resourceAttrInit(resourceId, AlicloudRealtimeComputeSessionClusterMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RealtimeComputeServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRealtimeComputeSessionCluster")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccrealtimecompute%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudRealtimeComputeSessionClusterBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"workspace":              "${alicloud_realtime_compute_vvp_instance.default.resource_id}",
					"namespace":              "${alicloud_realtime_compute_vvp_instance.default.vvp_instance_name}-default",
					"session_cluster_name":   name,
					"deployment_target_name": "default-queue",
					"engine_version":         "vvr-8.0.11-flink-1.17",
					"basic_resource_setting": []map[string]interface{}{
						{
							"parallelism": "1",
							"jobmanager_resource_setting_spec": []map[string]interface{}{
								{
									"cpu":    "1",
									"memory": "1Gi",
								},
							},
							"taskmanager_resource_setting_spec": []map[string]interface{}{
								{
									"cpu":    "1",
									"memory": "1Gi",
								},
							},
						},
					},
					"labels": map[string]interface{}{
						"env": "test",
					},
					"flink_conf": map[string]interface{}{
						"\"execution.checkpointing.interval\"": "180s",
						"\"restart-strategy\"":                 "fixed-delay",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"session_cluster_name":   name,
						"deployment_target_name": "default-queue",
						"engine_version":         "vvr-8.0.11-flink-1.17",
						"workspace":              CHECKSET,
						"namespace":              CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"basic_resource_setting": []map[string]interface{}{
						{
							"parallelism": "2",
							"jobmanager_resource_setting_spec": []map[string]interface{}{
								{
									"cpu":    "2",
									"memory": "2Gi",
								},
							},
							"taskmanager_resource_setting_spec": []map[string]interface{}{
								{
									"cpu":    "2",
									"memory": "2Gi",
								},
							},
						},
					},
					"labels": map[string]interface{}{
						"env": "prod",
					},
					"flink_conf": map[string]interface{}{
						"\"execution.checkpointing.interval\"": "300s",
					},
					"logging": []map[string]interface{}{
						{
							"logging_profile":               "default",
							"log4j2_configuration_template": "test-template",
							"log4j_loggers": []map[string]interface{}{
								{
									"logger_name":  "StdOut",
									"logger_level": "DEBUG",
								},
							},
							"log_reserve_policy": []map[string]interface{}{
								{
									"open_history":    "true",
									"expiration_days": "7",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"logging.#": "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudRealtimeComputeSessionClusterMap0 = map[string]string{
	"session_cluster_id": CHECKSET,
	"created_at":         CHECKSET,
	"status.#":           "1",
}

func AlicloudRealtimeComputeSessionClusterBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_vpc" "default" {
  is_default = false
  cidr_block = "172.16.0.0/16"
  vpc_name   = "test-tf-vpc-session-cluster"
}

resource "alicloud_vswitch" "default" {
  is_default   = false
  vpc_id       = alicloud_vpc.default.id
  zone_id      = "cn-beijing-g"
  cidr_block   = "172.16.0.0/24"
  vswitch_name = "test-tf-vswitch-session-cluster"
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
`, name)
}

// Test RealtimeCompute SessionCluster. <<< Resource test cases, automatically generated.
