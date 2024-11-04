package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/polardb"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

var primaryConnectionStringRegexp_mysql = "^[a-z-A-Z-0-9]+.mysql.([a-z-A-Z-0-9]+.){0,2}rds.aliyuncs.com"

var primaryConnectionStringRegexp_postgresql = "^[a-z-A-Z-0-9]+.pg.([a-z-A-Z-0-9]+.){0,1}rds.aliyuncs.com"

func TestAccAliCloudPolarDBPrimaryEndpointConfigUpdate(t *testing.T) {
	var v *polardb.DBEndpoint
	rand := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	name := fmt.Sprintf("tf-testAccPolarDBPrimaryEndpoint-%s", rand)
	var basicMap = map[string]string{
		"db_cluster_id": CHECKSET,
	}
	resourceId := "alicloud_polardb_primary_endpoint.default"
	ra := resourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return &PolarDBService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribePolarDBClusterEndpoint")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourcePolarDBPrimaryEndpointConfigDependence)
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
						"endpoint_type": "Primary",
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
						"ssl_enabled":     "Enable",
						"net_type":        "Private",
						"ssl_auto_rotate": "Enable",
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
					"connection_prefix": "tf-privatetestacc",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"connection_prefix": "tf-privatetestacc",
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
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"net_type"},
			},
		},
	})
}

func resourcePolarDBPrimaryEndpointConfigDependence(name string) string {
	return fmt.Sprintf(`
        variable "name" {
			default = "%s"
        }

		data "alicloud_vpcs" "default" {
			name_regex = "^default-NODELETING$"
		}
		resource "alicloud_vpc" "default" {
			vpc_name = var.name
		}
	
		resource "alicloud_vswitch" "default" {
			zone_id = data.alicloud_polardb_node_classes.default.classes.0.zone_id
			vpc_id = alicloud_vpc.default.id
			cidr_block = cidrsubnet(alicloud_vpc.default.cidr_block, 8, 4)
		}

		data "alicloud_polardb_node_classes" "default" {
			pay_type   = "PostPaid"
			db_type    = "MySQL"
			db_version = "8.0"
			category   = "Normal"
		}

        resource "alicloud_polardb_cluster" "cluster" {
                db_type = "MySQL"
                db_version = "8.0"
                pay_type = "PostPaid"
                db_node_class = "polar.mysql.x4.medium"
                vswitch_id = alicloud_vswitch.default.id
                description = "${var.name}"
        }
        `, name)
}
