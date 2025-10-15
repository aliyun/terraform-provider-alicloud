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

func TestAccAliCloudPolarDBZonalAccountPrivilege_life(t *testing.T) {
	var v *polardb.DBAccount
	name := "tf-testAccPolarDBClusterZonalAccountPrivilege"
	var basicMap = map[string]string{
		"db_cluster_id": CHECKSET,
	}
	resourceId := "alicloud_polardb_on_ens_account_privilege.default"
	ra := resourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return &PolarDbServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribePolarDBAccountPrivilege")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourcePolarDBClusterZonalAccountPrivilegeDependence)
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
					"db_cluster_id":     "${alicloud_polardb_on_ens_cluster.cluster.id}",
					"account_privilege": "ReadWrite",
					"account_name":      "${alicloud_polardb_on_ens_account.account_1.account_name}",
					"db_names":          []string{"${alicloud_polardb_on_ens_database.database_1.db_name}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_names.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_cluster_id":     "${alicloud_polardb_on_ens_cluster.cluster.id}",
					"account_privilege": "ReadWrite",
					"account_name":      "${alicloud_polardb_on_ens_account.account_1.account_name}",
					"db_names":          []string{"${alicloud_polardb_on_ens_database.database_1.db_name}", "${alicloud_polardb_on_ens_database.database_2.db_name}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_names.#": "2",
					}),
				),
			},
		},
	})
}

func resourcePolarDBClusterZonalAccountPrivilegeDependence(name string) string {
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
			db_cluster_nodes_configs = {"nodeWriter":"{\"db_node_class\":\"polar.mysql.x4.medium.c\",\"db_node_role\":\"Writer\"}"}
        }

		resource "alicloud_polardb_on_ens_account" "account_1" {
		  db_cluster_id          = alicloud_polardb_on_ens_cluster.cluster.id
		  account_name           = "for_terraform1"
		  account_password       = "Ali123789"
		  account_description    = "from_terraform_modify"
		}
		
		resource "alicloud_polardb_on_ens_database" "database_1" {
		  db_cluster_id         = alicloud_polardb_on_ens_cluster.cluster.id
		  db_description        = "test terraforms"
		  db_name               = "for_terraform"
		}
		
		resource "alicloud_polardb_on_ens_database" "database_2" {
		  db_cluster_id         = alicloud_polardb_on_ens_cluster.cluster.id
		  db_description        = "test terraforms"
		  db_name               = "for_terraform_2"
		}
        `, name)
}
