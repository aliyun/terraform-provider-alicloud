package alicloud

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/polardb"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudPolarDBEndpointAddressConfigUpdate(t *testing.T) {
	var v *polardb.Address
	rand := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	name := fmt.Sprintf("tf-testAccPolarDBconnection%s", rand)
	var basicMap = map[string]string{
		"db_cluster_id":  CHECKSET,
		"db_endpoint_id": CHECKSET,
		"net_type":       "Public",
	}
	resourceId := "alicloud_polardb_endpoint_address.default"
	ra := resourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return &PolarDBService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribePolarDBConnection")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourcePolarDBEndpointAddressConfigDependence)
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
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"connection_prefix": "tf-testacc",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"connection_prefix": "tf-testacc",
					}),
				),
			},
		},
	})
}

func resourcePolarDBEndpointAddressConfigDependence(name string) string {
	return fmt.Sprintf(`
	%s
	variable "creation" {
		default = "PolarDB"
	}

	variable "name" {
		default = "%s"
	}

	resource "alicloud_polardb_cluster" "default" {
		db_type = "MySQL"
		db_version = "8.0"
		pay_type = "PostPaid"
		db_node_class = "polar.mysql.x4.large"
		vswitch_id = "${alicloud_vswitch.default.id}"
		description = "${var.name}"
	}

	data "alicloud_polardb_endpoints" "default1" {
	  db_cluster_id    = "${alicloud_polardb_cluster.default.id}"
	}
	`, PolarDBCommonTestCase, name)
}
