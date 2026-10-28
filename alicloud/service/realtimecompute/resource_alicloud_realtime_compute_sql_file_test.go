package realtimecompute_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/realtimecompute"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test RealtimeCompute SqlFile. >>> Resource test cases, automatically generated.
// Case SqlFile lifecycle test
func TestAccAliCloudRealtimeComputeSqlFile_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_realtime_compute_sql_file.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudRealtimeComputeSqlFileMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return realtimecompute.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeRealtimeComputeSqlFile")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccrealtimecompute%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudRealtimeComputeSqlFileBasicDependence0)
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
					"workspace":            "${alicloud_realtime_compute_vvp_instance.default.resource_id}",
					"namespace":            "${alicloud_realtime_compute_vvp_instance.default.vvp_instance_name}-default",
					"name":                 name,
					"sql_script":           "SELECT 1;",
					"description":          "tf acc test sql file",
					"batch_mode":           "true",
					"session_cluster_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                 name,
						"sql_script":           "SELECT 1;",
						"description":          "tf acc test sql file",
						"batch_mode":           "true",
						"session_cluster_name": name,
						"workspace":            acctest.CHECKSET,
						"namespace":            acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"workspace":            "${alicloud_realtime_compute_vvp_instance.default.resource_id}",
					"namespace":            "${alicloud_realtime_compute_vvp_instance.default.vvp_instance_name}-default",
					"name":                 name + "_update",
					"sql_script":           "SELECT 2;",
					"description":          "tf acc test sql file update",
					"batch_mode":           "false",
					"session_cluster_name": name + "_update",
					"parent_id":            "${alicloud_realtime_compute_sql_file.probe.parent_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                 name + "_update",
						"sql_script":           "SELECT 2;",
						"description":          "tf acc test sql file update",
						"batch_mode":           "false",
						"session_cluster_name": name + "_update",
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
	"sql_file_id": acctest.CHECKSET,
}

func AlicloudRealtimeComputeSqlFileBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_vpc" "default" {
  is_default = false
  cidr_block = "172.16.0.0/16"
  vpc_name   = var.name
}

resource "alicloud_vswitch" "default" {
  is_default   = false
  vpc_id       = alicloud_vpc.default.id
  zone_id      = "cn-beijing-g"
  cidr_block   = "172.16.0.0/24"
  vswitch_name = var.name
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

resource "alicloud_realtime_compute_sql_file" "probe" {
  workspace            = alicloud_realtime_compute_vvp_instance.default.resource_id
  namespace            = "${alicloud_realtime_compute_vvp_instance.default.vvp_instance_name}-default"
  name                 = "${var.name}-probe"
  sql_script           = "SELECT 1;"
  session_cluster_name = "${var.name}-probe"
}
`, name)
}

// Test RealtimeCompute SqlFile. <<< Resource test cases, automatically generated.
