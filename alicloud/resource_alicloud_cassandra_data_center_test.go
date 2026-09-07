package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cassandra"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func SkipTestAccAlicloudCassandraDataCenter_basic(t *testing.T) {
	// Cassandra has been offline
	t.Skip("Cassandra has been offline")
	var v cassandra.DescribeDataCenterResponse
	resourceId := "alicloud_cassandra_data_center.default"
	ra := resourceAttrInit(resourceId, CassandraDataCenterMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CassandraService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCassandraDataCenter")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccCassandraDataCenter%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, CassandraDataCenterBasicdependence)
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
					"cluster_id":       "${alicloud_cassandra_cluster.default.id}",
					"instance_type":    "cassandra.c.large",
					"node_count":       "2",
					"pay_type":         "PayAsYouGo",
					"vswitch_id":       "${length(data.alicloud_vswitches.default_2.ids) > 0 ? data.alicloud_vswitches.default_2.ids[0] : alicloud_vswitch.this_2[0].id}",
					"disk_size":        "160",
					"disk_type":        "cloud_ssd",
					"data_center_name": "dc2-2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_type":    "cassandra.c.large",
						"node_count":       "2",
						"pay_type":         "PayAsYouGo",
						"vswitch_id":       CHECKSET,
						"disk_size":        "160",
						"disk_type":        "cloud_ssd",
						"data_center_name": "dc2-2",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"data_center_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"data_center_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_public": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_public": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"disk_size": "240",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"disk_size": "240",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_type": "cassandra.c.xlarge",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_type": "cassandra.c.xlarge",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"node_count": "3",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"node_count": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_public":    "false",
					"data_center_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_public":    "false",
						"data_center_name": name + "_update",
					}),
				),
			},
		},
	})
}

var CassandraDataCenterMap = map[string]string{
	"auto_renew":     "false",
	"data_center_id": CHECKSET,
	"status":         CHECKSET,
}

func CassandraDataCenterBasicdependence(name string) string {
	return fmt.Sprintf(`
		variable "name" {
		  default = "%s"
		}
		data "alicloud_cassandra_zones" "default" {
		}
		
		data "alicloud_vpcs" "default" {
		  name_regex = "default-NODELETING"
		}
		
		data "alicloud_vswitches" "default_1" {
		  vpc_id = data.alicloud_vpcs.default.ids[0]
		  zone_id = data.alicloud_cassandra_zones.default.zones[length(data.alicloud_cassandra_zones.default.ids)+(-1)].id
		}
		
		resource "alicloud_vswitch" "this_1" {
		  count = "${length(data.alicloud_vswitches.default_1.ids) > 0 ? 0 : 1}"
		  name = "${var.name}"
		  vpc_id = "${data.alicloud_vpcs.default.ids.0}"
		  availability_zone = data.alicloud_cassandra_zones.default.zones[length(data.alicloud_cassandra_zones.default.ids)+(-1)].id
		  cidr_block = "${cidrsubnet(data.alicloud_vpcs.default.vpcs.0.cidr_block, 8, 4)}"
		}
		data "alicloud_vswitches" "default_2" {
		  vpc_id = data.alicloud_vpcs.default.ids[0]
		  zone_id = data.alicloud_cassandra_zones.default.zones[length(data.alicloud_cassandra_zones.default.ids)+(-2)].id
		}
		
		resource "alicloud_vswitch" "this_2" {
		  count = "${length(data.alicloud_vswitches.default_2.ids) > 0 ? 0 : 1}"
		  name = "${var.name}_2"
		  vpc_id = "${data.alicloud_vpcs.default.ids.0}"
		  availability_zone = data.alicloud_cassandra_zones.default.zones[length(data.alicloud_cassandra_zones.default.ids)+(-2)].id
		  cidr_block = "${cidrsubnet(data.alicloud_vpcs.default.vpcs.0.cidr_block, 8, 10)}"
		}
		resource "alicloud_cassandra_cluster" "default" {
		  cluster_name = "${var.name}"
		  data_center_name = "${var.name}"
		  auto_renew = "false"
		  instance_type = "cassandra.c.large"
		  major_version = "3.11"
		  node_count = "2"
		  pay_type = "PayAsYouGo"
		  vswitch_id = "${length(data.alicloud_vswitches.default_1.ids) > 0 ? data.alicloud_vswitches.default_1.ids[0] : alicloud_vswitch.this_1[0].id}"
		  disk_size = "160"
		  disk_type = "cloud_ssd"
		  maintain_start_time = "18:00Z"
		  maintain_end_time = "20:00Z"
		  ip_white = "127.0.0.1"
		}
		`, name)
}
