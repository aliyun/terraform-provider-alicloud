package alicloud

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/polardb"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func init() {
}

func TestAccAliCloudPolarDBZonalEndpoint_life(t *testing.T) {
	var v *polardb.DBEndpoint
	name := "tf-testAccPolarDBClusterZonalEndpoint"
	var basicMap = map[string]string{
		"db_cluster_id": CHECKSET,
	}
	resourceId := "alicloud_polardb_on_ens_endpoint.default"
	ra := resourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return &PolarDbServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeDBClusterEndpointsZonal")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourcePolarDBClusterZonalEndpointConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_cluster_id":        "${alicloud_polardb_on_ens_cluster.cluster.id}",
					"db_cluster_nodes_ids": "${alicloud_polardb_on_ens_cluster.cluster.db_cluster_nodes_ids}",
					"endpoint_config":      map[string]string{"MasterAcceptReads": "on"},
					"read_write_mode":      "ReadWrite",
					"nodes_key":            []string{"node1", "nodeWriter"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"net_type":      "Private",
						"endpoint_type": "Custom",
						"nodes_key.#":   "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_cluster_id":           "${alicloud_polardb_on_ens_cluster.cluster.id}",
					"db_cluster_nodes_ids":    "${alicloud_polardb_on_ens_cluster.cluster.db_cluster_nodes_ids}",
					"endpoint_config":         map[string]string{"MasterAcceptReads": "on"},
					"db_endpoint_description": "123",
					"read_write_mode":         "ReadWrite",
					"nodes_key":               []string{"node1", "nodeWriter"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_endpoint_description": "123",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_cluster_id":           "${alicloud_polardb_on_ens_cluster.cluster.id}",
					"db_cluster_nodes_ids":    "${alicloud_polardb_on_ens_cluster.cluster.db_cluster_nodes_ids}",
					"endpoint_config":         map[string]string{"MasterAcceptReads": "on"},
					"db_endpoint_description": "123",
					"read_write_mode":         "ReadWrite",
					"nodes_key":               []string{"nodeWriter"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"nodes_key.#": "1",
					}),
				),
			},
		},
	})
}

func resourcePolarDBClusterZonalEndpointConfigDependence(name string) string {
	return fmt.Sprintf(`
		variable "name" {
			default = "%s"
		}
	
		resource "alicloud_ens_network" "network" {
		  network_name = var.name
		
		  description   = "LoadBalancerNetworkDescription_autotest"
		  cidr_block    = "192.168.2.0/24"
		  ens_region_id = "sg-singapore-9"
		}
		
		resource "alicloud_ens_vswitch" "switch" {
		  description  = "LoadBalancerVSwitchDescription_autotest"
		  cidr_block   = "192.168.2.0/24"
		  vswitch_name = var.name
		
		  ens_region_id = "sg-singapore-9"
		  network_id    = alicloud_ens_network.network.id
		}
		
		locals {
			vpc_id = alicloud_ens_network.network.id
			vswitch_id = alicloud_ens_vswitch.switch.id
		}

        resource "alicloud_polardb_on_ens_cluster" "cluster" {
			db_node_class = "polar.mysql.x4.medium.c"
			ens_region_id = "sg-singapore-9"
			vpc_id = alicloud_ens_network.network.id
			vswitch_id = alicloud_ens_vswitch.switch.id
			db_cluster_nodes_configs = {"nodeWriter":"{\"db_node_class\":\"polar.mysql.x4.medium.c\",\"db_node_role\":\"Writer\"}","node1": "{\"db_node_class\":\"polar.mysql.x4.medium.c\",\"db_node_role\":\"Reader\"}"}
        }

        `, name)
}
