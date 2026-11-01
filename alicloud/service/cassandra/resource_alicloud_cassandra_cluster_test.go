// collection-order: bypass security_groups -- pre-existing stock attribute outside this change; cassandra delisted, original test never carried reorder coverage
package cassandra_test

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cassandra"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	cassvc "github.com/aliyun/terraform-provider-alicloud/alicloud/service/cassandra"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/sweep"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_cassandra_cluster", &resource.Sweeper{
		Name: "alicloud_cassandra_cluster",
		F:    testSweepCassandraCluster,
	})
}

// TestMain bootstraps the external cassandra_test package so registered
// sweepers can be selected via -sweep, mirroring the root test entry.
func TestMain(m *testing.M) {
	resource.TestMain(m)
}

func testSweepCassandraCluster(region string) error {
	rawClient, err := sweep.SharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testacc",
	}

	request := cassandra.CreateDescribeClustersRequest()
	request.PageSize = requests.NewInteger(helper.PageSizeXLarge)
	raw, err := client.WithCassandraClient(func(cassandraClient *cassandra.Client) (interface{}, error) {
		return cassandraClient.DescribeClusters(request)
	})
	if err != nil {
		log.Printf("[ERROR] Error retrieving Cassandra Clusters: %s", errs.WrapError(err))
	}
	response, _ := raw.(*cassandra.DescribeClustersResponse)

	sweeped := false
	for _, v := range response.Clusters.Cluster {
		id := v.ClusterId
		name := v.ClusterName
		skip := true
		if !sweep.SweepAll() {
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
					skip = false
					break
				}
			}
			if skip {
				log.Printf("[INFO] Skipping Cassandra Clusters: %s (%s)", name, id)
				continue
			}
		}
		sweeped = true
		log.Printf("[INFO] Deleting Cassandra Clusters: %s (%s)", name, id)
		req := cassandra.CreateDeleteClusterRequest()
		req.ClusterId = id
		_, err := client.WithCassandraClient(func(cassandraClient *cassandra.Client) (interface{}, error) {
			return cassandraClient.DeleteCluster(req)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete Cassandra Clusters (%s (%s)): %s", name, id, err)
		}
		log.Printf("[INFO] Purging Cassandra Clusters: %s (%s)", name, id)
		reqPurge := cassandra.CreatePurgeClusterRequest()
		reqPurge.ClusterId = id
		_, errPurge := client.WithCassandraClient(func(cassandraClient *cassandra.Client) (interface{}, error) {
			return cassandraClient.PurgeCluster(reqPurge)
		})
		if errPurge != nil {
			log.Printf("[ERROR] Failed to purge Cassandra Clusters (%s (%s)): %s", name, id, err)
		}
	}
	if sweeped {
		// Waiting 30 seconds to ensure these Cassandra Clusters have been deleted.
		time.Sleep(30 * time.Second)
	}
	return nil
}

func SkipTestAccAlicloudCassandraCluster_basic(t *testing.T) {
	t.Skip("The cloud database Cassandra has been closed for sale. For the use of Cassandra in the future, Lindorm is recommended")
	var v cassandra.Cluster
	resourceId := "alicloud_cassandra_cluster.default"
	ra := acctest.ResourceAttrInit(resourceId, CassandraClusterMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cassvc.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCassandraCluster")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccCassandraCluster%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, CassandraClusterBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},

		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
						"vswitch_id":          acctest.CHECKSET,
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
	"status":     acctest.CHECKSET,
}

func CassandraClusterBasicdependence(name string) string {
	return fmt.Sprintf(`
		variable "name" {
			default = "%s"
		}
		data "alicloud_cassandra_zones" "default" {
		}
		
		data "alicloud_vpcs" "default" {
			name_regex = "^default-NODELETING$"
		}
		
		data "alicloud_vswitches" "default" {
		  vpc_id = data.alicloud_vpcs.default.ids[0]
		  zone_id = data.alicloud_cassandra_zones.default.zones[length(data.alicloud_cassandra_zones.default.ids)-1].id
		}
		
		resource "alicloud_security_group" "default" {
		  name = "${var.name}"
		  description = "New security group"
		  vpc_id = data.alicloud_vpcs.default.ids[0]
		}
		
		resource "alicloud_vswitch" "this" {
		  count = "${length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1}"
		  name = "${var.name}"
		  vpc_id = "${data.alicloud_vpcs.default.ids.0}"
		  availability_zone = data.alicloud_cassandra_zones.default.zones[length(data.alicloud_cassandra_zones.default.ids)-1].id
		  cidr_block = "${cidrsubnet(data.alicloud_vpcs.default.vpcs.0.cidr_block, 8, 4)}"
		}
		`, name)
}
