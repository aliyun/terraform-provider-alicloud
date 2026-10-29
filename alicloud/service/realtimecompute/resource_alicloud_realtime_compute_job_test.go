package realtimecompute_test

// local_variables is pre-existing stock carried verbatim from the flat
// alicloud/ layout by this product migration; the move does not add or change
// it, so convergent reorder coverage is known stock debt.
// collection-order: bypass local_variables -- pre-existing attribute outside this change

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/realtimecompute"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudRealtimeComputeJob_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_realtime_compute_job.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudRealtimeComputeJobMap11874)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return realtimecompute.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeRealtimeComputeJob")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccrealtimecompute%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudRealtimeComputeJobBasicDependence11874)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_queue_name": "default-queue",
					"resource_id":         "${alicloud_realtime_compute_vvp_instance.create_VvpInstance.resource_id}",
					"restore_strategy": []map[string]interface{}{
						{
							"kind": "NONE",
						},
					},
					"namespace":     "${alicloud_realtime_compute_vvp_instance.create_VvpInstance.vvp_instance_name}-default",
					"deployment_id": "${alicloud_realtime_compute_deployment.create_Deployment5.deployment_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_queue_name": "default-queue",
						"resource_id":         acctest.CHECKSET,
						"namespace":           acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"stop_strategy": "NONE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"stop_strategy": "NONE",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": []map[string]interface{}{
						{
							"current_job_status": "CANCELLED",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"resource_queue_name", "stop_strategy"},
			},
		},
	})
}

func TestAccAliCloudRealtimeComputeJob_basic11874(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_realtime_compute_job.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudRealtimeComputeJobMap11874)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return realtimecompute.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeRealtimeComputeJob")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccrealtimecompute%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudRealtimeComputeJobBasicDependence11874)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_queue_name": "default-queue",
					"resource_id":         "${alicloud_realtime_compute_vvp_instance.create_VvpInstance.resource_id}",
					"local_variables": []map[string]interface{}{
						{
							"value": "qq",
							"name":  "tt",
						},
					},
					"restore_strategy": []map[string]interface{}{
						{
							"kind":                 "NONE",
							"job_start_time_in_ms": "1763694521254",
						},
					},
					"namespace":     "${alicloud_realtime_compute_vvp_instance.create_VvpInstance.vvp_instance_name}-default",
					"deployment_id": "${alicloud_realtime_compute_deployment.create_Deployment5.deployment_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_queue_name": "default-queue",
						"resource_id":         acctest.CHECKSET,
						"local_variables.#":   "1",
						"namespace":           acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"stop_strategy": "NONE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"stop_strategy": "NONE",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": []map[string]interface{}{
						{
							"current_job_status": "CANCELLED",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"resource_queue_name", "stop_strategy"},
			},
		},
	})
}

var AlicloudRealtimeComputeJobMap11874 = map[string]string{
	"status.#": acctest.CHECKSET,
}

func AlicloudRealtimeComputeJobBasicDependence11874(name string) string {
	sqlScript := `create temporary table ` + "`datagen`" + ` ( id varchar, name varchar ) with ( 'connector' = 'datagen' );
create temporary table ` + "`blackhole`" + ` ( id varchar, name varchar ) with ( 'connector' = 'blackhole' );
insert into blackhole select * from datagen;`

	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alicloud_vpc" "create_Vpc" {
  is_default = false
  cidr_block = "172.16.0.0/16"
  vpc_name   = "test-tf-vpc-deployment"
}

resource "alicloud_vswitch" "create_Vswitch" {
  is_default   = false
  vpc_id       = alicloud_vpc.create_Vpc.id
  zone_id      = "cn-beijing-g"
  cidr_block   = "172.16.0.0/24"
  vswitch_name = "test-tf-vSwitch-deployment"
}

resource "alicloud_oss_bucket" "create_bucket" {
}

resource "alicloud_realtime_compute_vvp_instance" "create_VvpInstance" {
  vvp_instance_name = "${var.name}"
  storage {
    oss {
      bucket = alicloud_oss_bucket.create_bucket.id
    }
  }
  vpc_id      = alicloud_vpc.create_Vpc.id
  vswitch_ids = ["${alicloud_vswitch.create_Vswitch.id}"]
  resource_spec {
    cpu       = "8"
    memory_gb = "32"
  }
  payment_type = "PayAsYouGo"
  zone_id     = alicloud_vswitch.create_Vswitch.zone_id
}

resource "alicloud_realtime_compute_deployment" "create_Deployment5" {
  deployment_name = "tf-test-deployment-sql-24"
  engine_version  = "vvr-8.0.11-flink-1.17"
  resource_id     = alicloud_realtime_compute_vvp_instance.create_VvpInstance.resource_id
  execution_mode  = "STREAMING"
  
  deployment_target {
    mode = "PER_JOB"
    name = "default-queue"
  }
  
  namespace = "${alicloud_realtime_compute_vvp_instance.create_VvpInstance.vvp_instance_name}-default"
  
  artifact {
    kind = "SQLSCRIPT"
    sql_artifact {
      sql_script = %q
    }
  }
}`, name, sqlScript)
}
