package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test RealtimeCompute SqlFile. >>> Resource test cases.
func TestAccAliCloudRealtimeComputeSqlFile_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_realtime_compute_sql_file.default"
	ra := resourceAttrInit(resourceId, AlicloudRealtimeComputeSqlFileMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RealtimeComputeServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRealtimeComputeSqlFile")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccrealtimecompute%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudRealtimeComputeSqlFileBasicDependence0)
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
					"workspace":            "${alicloud_realtime_compute_vvp_instance.default.resource_id}",
					"namespace":            "${alicloud_realtime_compute_vvp_instance.default.vvp_instance_name}-default",
					"name":                 name,
					"sql_script":           "CREATE TABLE datagen (id VARCHAR, name VARCHAR) WITH ('connector' = 'datagen'); INSERT INTO blackhole SELECT * FROM datagen;",
					"description":          "This is a test sql file.",
					"batch_mode":           "false",
					"session_cluster_name": "test-session-cluster",
					"parent_id":            "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                 name,
						"description":          "This is a test sql file.",
						"batch_mode":           "false",
						"session_cluster_name": "test-session-cluster",
						"workspace":            CHECKSET,
						"namespace":            CHECKSET,
						"sql_file_id":          CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"workspace":            "${alicloud_realtime_compute_vvp_instance.default.resource_id}",
					"namespace":            "${alicloud_realtime_compute_vvp_instance.default.vvp_instance_name}-default",
					"name":                 fmt.Sprintf("%s-updated", name),
					"sql_script":           "CREATE TABLE datagen2 (id VARCHAR, name VARCHAR) WITH ('connector' = 'datagen'); INSERT INTO blackhole SELECT * FROM datagen2;",
					"description":          "This is an updated sql file.",
					"batch_mode":           "true",
					"session_cluster_name": "test-session-cluster-updated",
					"parent_id":            "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                 fmt.Sprintf("%s-updated", name),
						"description":          "This is an updated sql file.",
						"batch_mode":           "true",
						"session_cluster_name": "test-session-cluster-updated",
						"parent_id":            "1",
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

var AlicloudRealtimeComputeSqlFileMap0 = map[string]string{
	"sql_file_id": CHECKSET,
	"region_id":   CHECKSET,
}

func AlicloudRealtimeComputeSqlFileBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alicloud_vpc" "default" {
  is_default = false
  cidr_block = "172.16.0.0/16"
  vpc_name   = "test-tf-vpc-sql-file"
}

resource "alicloud_vswitch" "default" {
  is_default   = false
  vpc_id       = alicloud_vpc.default.id
  zone_id      = "cn-beijing-g"
  cidr_block   = "172.16.0.0/24"
  vswitch_name = "test-tf-vswitch-sql-file"
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

// Test RealtimeCompute SqlFile. <<< Resource test cases.
