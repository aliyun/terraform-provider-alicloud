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

func TestAccAliCloudPolarDBZonalDatabase_life(t *testing.T) {
	var v *polardb.Database
	name := "tf-testAccPolarDBClusterZonalDatabase"
	var basicMap = map[string]string{
		"db_cluster_id": CHECKSET,
	}
	resourceId := "alicloud_polardb_on_ens_database.default"
	ra := resourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return &PolarDbServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribePolarDBDatabase")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourcePolarDBClusterZonalDatabaseDependence)
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
					"db_cluster_id":      "${alicloud_polardb_on_ens_cluster.cluster.id}",
					"db_description":     "test terraform",
					"db_name":            "from_terraform",
					"character_set_name": "utf8",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_description":     "test terraform",
						"db_name":            "from_terraform",
						"character_set_name": "utf8",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_cluster_id":      "${alicloud_polardb_on_ens_cluster.cluster.id}",
					"db_description":     "test terraform modify",
					"db_name":            "from_terraform",
					"character_set_name": "utf8",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_description": "test terraform modify",
					}),
				),
			},
		},
	})
}

func resourcePolarDBClusterZonalDatabaseDependence(name string) string {
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
			db_type = "MySQL"
			db_version = "8.0"
			pay_type = "PrePaid"
			db_node_class = "polar.mysql.x4.medium.c"
			ens_region_id = "sg-singapore-9"
			vpc_id = alicloud_ens_network.network.id
			vswitch_id = alicloud_ens_vswitch.switch.id
			description = var.name
			storage_space = "20"
			db_minor_version = "8.0.2"
			renewal_status = "AutoRenewal"
			target_minor_version = "innovate_x86#20250311"
			storage_type = "ESSDPL0"
			period = "1"
			auto_renew_period = "1"
			db_cluster_nodes_configs = {"nodeWriter":"{\"db_node_class\":\"polar.mysql.x4.medium.c\",\"db_node_role\":\"Writer\"}"}
        }
        `, name)
}
