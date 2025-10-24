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

func TestAccAliCloudPolarDBZonalAccount_life(t *testing.T) {
	var v *polardb.DBAccount
	name := "tf-testAccPolarDBClusterZonalAccount"
	var basicMap = map[string]string{
		"db_cluster_id": CHECKSET,
	}
	resourceId := "alicloud_polardb_zonal_account.default"
	ra := resourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return &PolarDbServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribePolarDBAccount")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourcePolarDBClusterZonalAccountDependence)
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
					"db_cluster_id":       "${alicloud_polardb_zonal_db_cluster.cluster.id}",
					"account_name":        "from_terraform",
					"account_type":        "Normal",
					"account_password":    "Ali123789",
					"account_description": "from_terraform",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_name":        "from_terraform",
						"account_type":        "Normal",
						"account_description": "from_terraform",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"account_password": "Ali123798",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_password": "Ali123798",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"account_description": "from_terraform_modify",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_description": "from_terraform_modify",
					}),
				),
			},
		},
	})
}

func resourcePolarDBClusterZonalAccountDependence(name string) string {
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
			db_cluster_nodes_configs = {"db_node_1":"{\"db_node_class\":\"polar.mysql.x4.medium.c\",\"db_node_role\":\"Writer\"}"}
        }

        `, name)
}
