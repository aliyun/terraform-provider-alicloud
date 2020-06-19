package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cassandra"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudCassandraCluster_basic(t *testing.T) {
	var v cassandra.Cluster
	resourceId := "alicloud_cassandra_cluster.default"
	ra := resourceAttrInit(resourceId, CassandraClusterMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CassandraService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCassandraCluster")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccCassandraCluster%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, CassandraClusterBasicdependence)
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
					"instance_type":       "cassandra.c.large",
					"major_version":       "3.11",
					"node_count":          "2",
					"pay_type":            "PayAsYouGo",
					"vswitch_id":          "${length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : alicloud_vswitch.this[0].id}",
					"disk_size":           "160",
					"disk_type":           "cloud_ssd",
					"maintain_start_time": "01:00Z",
					"maintain_end_time":   "03:00Z",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_type":       "cassandra.c.large",
						"major_version":       "3.11",
						"pay_type":            "PayAsYouGo",
						"ip_white":            "127.0.0.1",
						"node_count":          "2",
						"vswitch_id":          CHECKSET,
						"disk_size":           "160",
						"disk_type":           "cloud_ssd",
						"maintain_start_time": "01:00Z",
						"maintain_end_time":   "03:00Z",
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
					"cluster_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ip_white": "127.0.0.4,127.0.0.3",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ip_white": "127.0.0.4,127.0.0.3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_groups": []string{"${alicloud_security_group.default.id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_groups.#": "1",
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
					"data_center_name": "dc-1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"data_center_name": "dc-1",
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
					"maintain_start_time": "08:00Z",
					"maintain_end_time":   "10:00Z",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"maintain_start_time": "08:00Z",
						"maintain_end_time":   "10:00Z",
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
					"maintain_start_time": "18:00Z",
					"maintain_end_time":   "20:00Z",
					"ip_white":            "127.0.0.1",
					"cluster_name":        name + "_update",
					"data_center_name":    "dc-1—update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"maintain_start_time": "18:00Z",
						"maintain_end_time":   "20:00Z",
						"ip_white":            "127.0.0.1",
						"cluster_name":        name + "_update",
						"data_center_name":    "dc-1—update",
					}),
				),
			},
		},
	})
}

var CassandraClusterMap = map[string]string{
	"auto_renew": "false",
	"status":     CHECKSET,
}

func CassandraClusterBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
data "alicloud_cassandra_zones" "default" {
}

data "alicloud_vpcs" "default" {
}

data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids[0]
  zone_id = data.alicloud_cassandra_zones.default.zones[length(data.alicloud_cassandra_zones.default.ids)-1].id
}

resource "alicloud_security_group" "default" {
  name = "terraform-test-group"
  description = "New security group"
  vpc_id = data.alicloud_vpcs.default.ids[0]
}

resource "alicloud_vswitch" "this" {
  count = "${length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1}"
  name = "tf_testAccCassandra_vpc"
  vpc_id = "${data.alicloud_vpcs.default.ids.0}"
  availability_zone = data.alicloud_cassandra_zones.default.zones[length(data.alicloud_cassandra_zones.default.ids)-1].id
  cidr_block = "${cidrsubnet(data.alicloud_vpcs.default.vpcs.0.cidr_block, 8, 4)}"
}
`, name)
}
