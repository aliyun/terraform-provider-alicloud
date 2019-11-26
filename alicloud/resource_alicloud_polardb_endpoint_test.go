package alicloud

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/polardb"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudPolarDBEndpointConfigUpdate(t *testing.T) {
	var v *polardb.Address
	rand := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	name := fmt.Sprintf("tf-testAccPolarDBconnection%s", rand)
	var basicMap = map[string]string{
		"db_cluster_id":  CHECKSET,
		"db_endpoint_id": CHECKSET,
		"net_type":       "Public",
	}
	resourceId := "alicloud_polardb_endpoint.default"
	ra := resourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return &PolarDBService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribePolarDBConnection")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourcePolarDBEndpointConfigDependence)
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
					"db_cluster_id":  "${alicloud_polardb_cluster.default.id}",
					"db_endpoint_id": "${data.alicloud_polardb_endpoints.default1.endpoints[0].db_endpoint_id}",
					"net_type":       "Public",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"connection_prefix"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"connection_prefix": "pe-t4n8mug7j22nsf2eb111",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"connection_prefix": "pe-t4n8mug7j22nsf2eb111",
					}),
				),
			},
		},
	})
}

func resourcePolarDBEndpointConfigDependence(name string) string {
	return fmt.Sprintf(`
	

	variable "creation" {
		default = "PolarDB"
	}

	variable "name" {
		default = "%s"
	}

	variable "instancechargetype" {
		default = "PostPaid"
	}

	variable "engine" {
		default = "MySQL"
	}

	variable "engineversion" {
		default = "8.0"
	}

	variable "instanceclass" {
		default = "polar.mysql.x4.large"
	}

	resource "alicloud_polardb_cluster" "default" {
		db_type = "${var.engine}"
		db_version = "${var.engineversion}"
		pay_type = "${var.instancechargetype}"
		db_node_class = "${var.instanceclass}"
		vswitch_id = "vsw-t4nlb8goj0as5zaau0cm0"
		description = "${var.name}"
	}

	data "alicloud_polardb_endpoints" "default1" {
	  db_cluster_id    = "${alicloud_polardb_cluster.default.id}"
	}
	`, name)
}
