package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudCassandraBackupPlan_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cassandra_backup_plan.default"
	ra := resourceAttrInit(resourceId, AlicloudCassandraBackupPlanMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CassandraService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCassandraBackupPlan")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scassandrabackupplan%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCassandraBackupPlanBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"data_center_id": "${alicloud_cassandra_cluster.default.zone_id}",
					"cluster_id":     "${alicloud_cassandra_cluster.default.id}",
					"backup_time":    "00:10Z",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"data_center_id": CHECKSET,
						"cluster_id":     CHECKSET,
						"backup_time":    "00:10Z",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"active": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"active": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"backup_period": "Tuesday",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_period": "Tuesday",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"backup_time": "00:30Z",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_time": "00:30Z",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"retention_period": "7",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"retention_period": "7",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"active":           "false",
					"backup_period":    "Monday",
					"backup_time":      "00:10Z",
					"retention_period": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"active":           "false",
						"backup_period":    "Monday",
						"backup_time":      "00:10Z",
						"retention_period": "1",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

var AlicloudCassandraBackupPlanMap0 = map[string]string{

	"retention_period": CHECKSET,
	"backup_period":    CHECKSET,
	"data_center_id":   CHECKSET,
	"cluster_id":       CHECKSET,
	"backup_time":      "00:10Z",
	"active":           CHECKSET,
}

func AlicloudCassandraBackupPlanBasicDependence0(name string) string {
	return fmt.Sprintf(` 
		variable "name" {
		  default = "%s"
		}

		data "alicloud_cassandra_zones" "default" {
		}
		
		data "alicloud_vpcs" "default" {
		  is_default = true
		}
		
		data "alicloud_vswitches" "default_1" {
		  vpc_id = data.alicloud_vpcs.default.ids[0]
		  zone_id = data.alicloud_cassandra_zones.default.zones[length(data.alicloud_cassandra_zones.default.ids)+(-1)].id
		}
		
		resource "alicloud_vswitch" "this_1" {
		  count = length(data.alicloud_vswitches.default_1.ids) > 0 ? 0 : 1
		  vswitch_name = var.name
		  vpc_id = data.alicloud_vpcs.default.ids.0
		  zone_id = data.alicloud_cassandra_zones.default.zones[length(data.alicloud_cassandra_zones.default.ids)+(-1)].id
		  cidr_block = cidrsubnet(data.alicloud_vpcs.default.vpcs.0.cidr_block, 8, 4)
		}
		resource "alicloud_cassandra_cluster" "default" {
		  cluster_name = var.name
		  data_center_name = var.name
		  auto_renew = "false"
		  instance_type = "cassandra.c.large"
		  major_version = "3.11"
		  node_count = "2"
		  pay_type = "PayAsYouGo"
		  vswitch_id = length(data.alicloud_vswitches.default_1.ids) > 0 ? data.alicloud_vswitches.default_1.ids[0] : alicloud_vswitch.this_1[0].id
		  disk_size = "160"
		  disk_type = "cloud_ssd"
		  maintain_start_time = "18:00Z"
		  maintain_end_time = "20:00Z"
		  ip_white = "127.0.0.1"
		}
	`, name)
}
