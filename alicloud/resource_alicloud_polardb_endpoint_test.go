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

func TestAccAliCloudPolarDBEndpointConfigUpdate(t *testing.T) {
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
					"endpoint_type": "Custom",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"endpoint_type": "Custom",
					}),
				),
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
					"nodes": []string{"${data.alicloud_polardb_clusters.default.clusters.0.db_nodes.0.db_node_id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"nodes.#": "1",
					}),
				),
			},
			//todo: After resource polardb_node is supported, it is necessary to add a modification check on the “nodes” parameter
			{
				Config: testAccConfig(map[string]interface{}{
					"ssl_enabled":     "Enable",
					"net_type":        "Private",
					"ssl_auto_rotate": "Enable",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ssl_enabled":           "Enable",
						"net_type":              "Private",
						"ssl_connection_string": REGEXMATCH + privateConnectionStringRegexp,
						"ssl_auto_rotate":       "Enable",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_endpoint_description": "polar_db_endpoint_test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_endpoint_description": "polar_db_endpoint_test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"connection_prefix": "tf-privatetestacc-3",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"connection_prefix": "tf-privatetestacc-3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"port": "3307",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"port": "3307",
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
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"ssl_enabled", "net_type", "endpoint_config"},
			},
		},
	})
}

func TestAccAliCloudPolarDBEndpointConfigUpdate_SslConnectionStringAndConnectionPrefix(t *testing.T) {
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
					"endpoint_type": "Custom",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"endpoint_type": "Custom",
					}),
				),
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
					"nodes": []string{"${data.alicloud_polardb_clusters.default.clusters.0.db_nodes.0.db_node_id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"nodes.#": "1",
					}),
				),
			},
			//todo: After resource polardb_node is supported, it is necessary to add a modification check on the “nodes” parameter
			{
				Config: testAccConfig(map[string]interface{}{
					"connection_prefix": "tf-privatetestacc-4",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"connection_prefix": "tf-privatetestacc-4",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ssl_enabled":     "Enable",
					"net_type":        "Private",
					"ssl_auto_rotate": "Enable",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ssl_enabled":           "Enable",
						"net_type":              "Private",
						"ssl_connection_string": REGEXMATCH + privateConnectionStringRegexp,
						"ssl_auto_rotate":       "Enable",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"ssl_enabled", "net_type", "endpoint_config"},
			},
		},
	})
}

func resourcePolarDBEndpointConfigDependence(name string) string {
	return fmt.Sprintf(`
        variable "name" {
                default = "%s"
        }

		data "alicloud_polardb_zones" "default"{}
		data "alicloud_vpcs" "default" {
			name_regex = "^default-NODELETING$"
		}
		data "alicloud_vswitches" "default" {
			zone_id = data.alicloud_polardb_zones.default.ids[length(data.alicloud_polardb_zones.default.ids) - 1]
			vpc_id = data.alicloud_vpcs.default.ids.0
		}

        resource "alicloud_polardb_cluster" "cluster" {
                db_type = "MySQL"
                db_version = "8.0"
                pay_type = "PostPaid"
                db_node_class = "polar.mysql.x4.medium"
                vswitch_id = data.alicloud_vswitches.default.ids.0
                description = "${var.name}"
        }

		data "alicloud_polardb_clusters" "default" {
		  	ids = [alicloud_polardb_cluster.cluster.id]
		}
        `, name)
}
