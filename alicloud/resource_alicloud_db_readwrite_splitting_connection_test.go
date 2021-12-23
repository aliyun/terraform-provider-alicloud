package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

var DBReadWriteMap = map[string]string{
	"port":              "3306",
	"distribution_type": "Standard",
	"weight":            NOSET,
	"max_delay_time":    "30",
	"instance_id":       CHECKSET,
	"connection_string": CHECKSET,
}

func TestAccAlicloudRdsDBReadWriteSplittingConnection_update(t *testing.T) {
	var connection map[string]interface{}
	var primary map[string]interface{}
	var readonly map[string]interface{}

	resourceId := "alicloud_db_read_write_splitting_connection.default"
	ra := resourceAttrInit(resourceId, DBReadWriteMap)

	rc_connection := resourceCheckInitWithDescribeMethod(resourceId, &connection, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDBReadWriteSplittingConnection")
	rc_primary := resourceCheckInitWithDescribeMethod("alicloud_db_instance.default", &primary, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDBInstance")
	rc_readonly := resourceCheckInitWithDescribeMethod("alicloud_db_readonly_instance.default", &readonly, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDBReadonlyInstance")
	rand := acctest.RandIntRange(10000, 999999)

	rac := resourceAttrCheckInit(rc_connection, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccDBReadWrite%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDBReadWriteSplittingConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.DBReadwriteSplittingConnectionSupportedRegions)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":       "${alicloud_db_readonly_instance.default.master_db_instance_id}",
					"connection_prefix": "${var.name}",
					"distribution_type": "Standard",
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
					"max_delay_time":    "300",
					"distribution_type": "Custom",
					"weight": `${map(
						"${alicloud_db_instance.default.id}", "0",
						"${alicloud_db_readonly_instance.default.id}", "500"
					)}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					rc_primary.checkResourceExists(),
					rc_readonly.checkResourceExists(),
					testAccCheck(map[string]string{
						"max_delay_time":    "300",
						"weight.%":          "2",
						"distribution_type": "Custom",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":       "${alicloud_db_readonly_instance.default.master_db_instance_id}",
					"connection_prefix": "${var.name}",
					"distribution_type": "Standard",
					"max_delay_time":    "30",
					"weight":            REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"port":              "3306",
						"distribution_type": "Standard",
						"weight.%":          REMOVEKEY,
						"max_delay_time":    "30",
						"instance_id":       CHECKSET,
						"connection_string": CHECKSET,
					}),
				),
			},
		},
	})
}

func resourceDBReadWriteSplittingConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
data "alicloud_db_zones" "default"{
	engine = "MySQL"
	engine_version = "8.0"
	instance_charge_type = "PostPaid"
	category = "HighAvailability"
 	db_instance_storage_type = "cloud_essd"
}

data "alicloud_db_instance_classes" "default" {
    zone_id = data.alicloud_db_zones.default.zones.0.id
	engine = "MySQL"
	engine_version = "8.0"
    category = "HighAvailability"
 	db_instance_storage_type = "cloud_essd"
	instance_charge_type = "PostPaid"
}

data "alicloud_vpcs" "default" {
 name_regex = "^default-NODELETING"
}
data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_db_zones.default.zones.0.id
}

resource "alicloud_vswitch" "this" {
 count = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
 vswitch_name = var.name
 vpc_id = data.alicloud_vpcs.default.ids.0
 zone_id = data.alicloud_db_zones.default.ids.0
 cidr_block = cidrsubnet(data.alicloud_vpcs.default.vpcs.0.cidr_block, 8, 4)
}
locals {
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids.0 : concat(alicloud_vswitch.this.*.id, [""])[0]
  zone_id = data.alicloud_db_zones.default.ids[length(data.alicloud_db_zones.default.ids)-1]
}

data "alicloud_resource_manager_resource_groups" "default" {
	status = "OK"
}

resource "alicloud_security_group" "default" {
	name   = var.name
	vpc_id = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_db_instance" "default" {
    engine = "MySQL"
	engine_version = "8.0"
 	db_instance_storage_type = "cloud_essd"
	instance_type = data.alicloud_db_instance_classes.default.instance_classes.0.instance_class
	instance_storage = data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.min
	vswitch_id = local.vswitch_id
	instance_name = var.name
	security_ips = ["10.168.1.12", "100.69.7.112"]
}
resource "alicloud_db_readonly_instance" "default" {
	master_db_instance_id = alicloud_db_instance.default.id
	zone_id = alicloud_db_instance.default.zone_id
	engine_version = alicloud_db_instance.default.engine_version
	instance_type = alicloud_db_instance.default.instance_type
	instance_storage = alicloud_db_instance.default.instance_storage
	instance_name = "${var.name}_ro"
	vswitch_id = "${local.vswitch_id}"
}

`, name)
}
