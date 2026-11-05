package polardb_test

import (
	"fmt"
	"testing"

	sdkpolardb "github.com/aliyun/alibaba-cloud-sdk-go/services/polardb"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	polardb "github.com/aliyun/terraform-provider-alicloud/alicloud/service/polardb"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func init() {
}

func TestAccAliCloudPolarDBZonalEndpoint_life(t *testing.T) {
	var v *sdkpolardb.DBEndpoint
	name := "tf-testAccPolarDBClusterZonalEndpoint"
	var basicMap = map[string]string{
		"db_cluster_id": CHECKSET,
	}
	resourceId := "alicloud_polardb_zonal_endpoint.default"
	ra := acctest.ResourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return polardb.NewPolarDbServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeDBClusterEndpointsZonal")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourcePolarDBClusterZonalEndpointConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},

		// module name

		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_cluster_id":           "${alicloud_polardb_zonal_db_cluster.cluster.id}",
					"db_cluster_nodes_ids":    "${alicloud_polardb_zonal_db_cluster.cluster.db_cluster_nodes_ids}",
					"endpoint_config":         map[string]string{},
					"read_write_mode":         "ReadWrite",
					"endpoint_type":           "Custom",
					"auto_add_new_nodes":      "Enable",
					"net_type":                "Private",
					"db_endpoint_description": "",
					"nodes_key":               []string{"db_node_1", "db_node_2"},
					"vpc_id":                  "${alicloud_ens_network.network.id}",
					"vswitch_id":              "${alicloud_ens_vswitch.switch.id}",
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
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"nodes_key": []string{"db_node_1"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"nodes_key.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_endpoint_description": "123",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_endpoint_description": "123",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auto_add_new_nodes": "Disable",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_add_new_nodes": "Disable",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"endpoint_config": map[string]string{"MasterAcceptReads": "off"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"endpoint_config.MasterAcceptReads": "off",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"read_write_mode": "ReadOnly",
					"endpoint_config": map[string]string{"CausalConsistRead": "0", "MasterAcceptReads": "off"},
					"nodes_key":       []string{"db_node_2"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"read_write_mode": "ReadOnly",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_cluster_nodes_ids": map[string]string{},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_cluster_nodes_ids": NOSET,
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
		  ens_region_id = "tr-Istanbul-1"
		}
		
		resource "alicloud_ens_vswitch" "switch" {
		  description  = "LoadBalancerVSwitchDescription_autotest"
		  cidr_block   = "192.168.2.0/24"
		  vswitch_name = var.name
		
		  ens_region_id = "tr-Istanbul-1"
		  network_id    = alicloud_ens_network.network.id
		}
		
		locals {
			vpc_id = alicloud_ens_network.network.id
			vswitch_id = alicloud_ens_vswitch.switch.id
		}

        resource "alicloud_polardb_zonal_db_cluster" "cluster" {
			db_node_class = "polar.mysql.x4.medium.c"
			ens_region_id = "tr-Istanbul-1"
			vpc_id = alicloud_ens_network.network.id
			vswitch_id = alicloud_ens_vswitch.switch.id
			db_cluster_nodes_configs = {"db_node_1":"{\"db_node_class\":\"polar.mysql.x4.medium.c\",\"db_node_role\":\"Writer\"}","db_node_2": "{\"db_node_class\":\"polar.mysql.x4.medium.c\",\"db_node_role\":\"Reader\"}"}
        }

        `, name)
}
