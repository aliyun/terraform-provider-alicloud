// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test RealtimeCompute Job. >>> Resource test cases, automatically generated.
// Case fix_bug 11874
func TestAccAliCloudRealtimeComputeJob_basic11874(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_realtime_compute_job.default"
	ra := resourceAttrInit(resourceId, AlicloudRealtimeComputeJobMap11874)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RealtimeComputeServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRealtimeComputeJob")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccrealtimecompute%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudRealtimeComputeJobBasicDependence11874)
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
					"resource_queue_name": "default-queue",
					"resource_id":         "e7d4d4f4510947",
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
					"namespace":     "code-test-tf-deployment-default",
					"deployment_id": "${alicloud_realtime_compute_deployment.create_Deployment5.deployment_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_queue_name": "default-queue",
						"resource_id":         CHECKSET,
						"local_variables.#":   "1",
						"namespace":           CHECKSET,
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
	"status.#": CHECKSET,
}

func AlicloudRealtimeComputeJobBasicDependence11874(name string) string {
	sqlScript := `create temporary table ` + "`datagen`" + ` ( id varchar, name varchar ) with ( 'connector' = 'datagen' );
create temporary table ` + "`blackhole`" + ` ( id varchar, name varchar ) with ( 'connector' = 'blackhole' );
insert into blackhole select * from datagen;`

	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alicloud_realtime_compute_deployment" "create_Deployment5" {
  deployment_name = "tf-test-deployment-sql-24"
  engine_version  = "vvr-8.0.10-flink-1.17"
  resource_id     = "e7d4d4f4510947"
  execution_mode  = "STREAMING"
  
  deployment_target {
    mode = "PER_JOB"
    name = "default-queue"
  }
  
  namespace = "code-test-tf-deployment-default"
  
  artifact {
    kind = "SQLSCRIPT"
    sql_artifact {
      sql_script = %q
    }
  }
}`, name, sqlScript)
}

// Test RealtimeCompute Job. <<< Resource test cases, automatically generated.
