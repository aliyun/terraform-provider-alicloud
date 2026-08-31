package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudRdsDBInstanceEndpointMySql(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_rds_db_instance_endpoint.default"
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDBInstanceEndpoints")
	ra := resourceAttrInit(resourceId, nil)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := "tf-testAccRdsDBInstanceEndpoint"
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceRdsDBInstanceEndpointMysqlDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_id":                   "${alicloud_rds_db_node.default.db_instance_id}",
					"vpc_id":                           "${alicloud_db_instance.default.vpc_id}",
					"vswitch_id":                       "${alicloud_db_instance.default.vswitch_id}",
					"connection_string_prefix":         "zcctest",
					"port":                             "3306",
					"db_instance_endpoint_description": "test",
					"node_items": []map[string]interface{}{
						{
							"node_id": "${alicloud_rds_db_node.default.node_id}",
							"weight":  "10",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_id":                   CHECKSET,
						"vpc_id":                           CHECKSET,
						"vswitch_id":                       CHECKSET,
						"connection_string_prefix":         CHECKSET,
						"port":                             CHECKSET,
						"db_instance_endpoint_description": CHECKSET,
						"node_items.#":                     "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"vpc_id":     "${alicloud_vpc.default.id}",
					"vswitch_id": "${alicloud_vswitch.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpc_id":     CHECKSET,
						"vswitch_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"node_items": []map[string]interface{}{
						{
							"node_id": "${alicloud_rds_db_node.default.node_id}",
							"weight":  "15",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"node_items.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_endpoint_description": "test001",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_endpoint_description": "test001",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"port": "3308",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"port": "3308",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"connection_string_prefix": "zcctest001",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"connection_string_prefix": "zcctest001",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func resourceRdsDBInstanceEndpointMysqlDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_db_zones" "default" {
  engine                   = "MySQL"
  engine_version           = "8.0"
  instance_charge_type     = "PostPaid"
  category                 = "cluster"
  db_instance_storage_type = "cloud_essd"
}

data "alicloud_db_instance_classes" "default" {
  zone_id                  = data.alicloud_db_zones.default.zones.0.id
  engine                   = "MySQL"
  engine_version           = "8.0"
  category                 = "cluster"
  db_instance_storage_type = "cloud_essd"
  instance_charge_type     = "PostPaid"
}

data "alicloud_vpcs" "default" {
	name_regex = "^default-NODELETING$"
}
data "alicloud_vswitches" "default" {
	vpc_id 		 = data.alicloud_vpcs.default.ids.0
	zone_id      = data.alicloud_db_zones.default.ids.0
}

resource "alicloud_vpc" "default" {
  vpc_name       = var.name
  cidr_block = "172.16.0.0/16"
}
resource "alicloud_vswitch" "default" {
  vpc_id            = alicloud_vpc.default.id
  cidr_block        = "172.16.0.0/24"
  zone_id = data.alicloud_db_zones.default.ids.0
  vswitch_name              = var.name
}

resource "alicloud_db_instance" "default" {
  engine                   = "MySQL"
  engine_version           = "8.0"
  db_instance_storage_type = "cloud_essd"
  instance_type            = data.alicloud_db_instance_classes.default.instance_classes.0.instance_class
  instance_storage         = data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.min
  vswitch_id               = data.alicloud_vswitches.default.ids.0
  instance_name            = var.name
  zone_id 				   = data.alicloud_db_zones.default.ids.0
  zone_id_slave_a          = data.alicloud_db_zones.default.ids.1
}

resource "alicloud_rds_db_node" "default" {
  db_instance_id = alicloud_db_instance.default.id
  class_code     = alicloud_db_instance.default.instance_type
  zone_id        = alicloud_db_instance.default.zone_id
}
`, name)
}
