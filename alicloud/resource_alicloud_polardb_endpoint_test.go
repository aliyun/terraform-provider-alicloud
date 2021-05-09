package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/polardb"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

var privateConnectionStringRegexp = "^[a-z-A-Z-0-9]+.rwlb.([a-z-A-Z-0-9]+.){0,1}rds.aliyuncs.com"

func TestAccAlicloudPolarDBEndpointConfigUpdate(t *testing.T) {
	var v *polardb.DBEndpoint
	rand := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	name := fmt.Sprintf("tf-testAccPolarDBendpoint-%s", rand)
	var basicMap = map[string]string{
		"db_cluster_id": CHECKSET,
	}
	resourceId := "alicloud_polardb_endpoint.default"
	ra := resourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return &PolarDBService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribePolarDBClusterEndpoint")
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
					"db_cluster_id": "${alicloud_polardb_cluster.cluster.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"endpoint_type": "Custom",
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
					"auto_add_new_nodes": "Enable",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_add_new_nodes": "Enable",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"read_write_mode": "ReadWrite",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"read_write_mode": "ReadWrite",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"endpoint_config": map[string]string{
						"ConsistLevel": "1",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"endpoint_config.ConsistLevel": "1",
					}),
				),
			},
			//todo: After resource polardb_node is supported, it is necessary to add a modification check on the “nodes” parameter
			{
				Config: testAccConfig(map[string]interface{}{
					"ssl_enabled": "Enable",
					"net_type":    "Private",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ssl_enabled":           "Enable",
						"net_type":              "Private",
						"ssl_connection_string": REGEXMATCH + privateConnectionStringRegexp,
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

        data "alicloud_vpcs" "vpcs_ds"{
                is_default = "true"
        }

        resource "alicloud_polardb_cluster" "cluster" {
                db_type = "MySQL"
                db_version = "8.0"
                pay_type = "PostPaid"
                db_node_class = "polar.mysql.x4.medium"
                vswitch_id = "${data.alicloud_vpcs.vpcs_ds.vpcs.0.vswitch_ids.0}"
                description = "${var.name}"
        }
        `, name)
}
