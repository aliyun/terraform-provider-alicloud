package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudRdsDBDatabaseMySQL(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_db_database.default"

	var dbDatabaseBasicMap = map[string]string{
		"instance_id":   CHECKSET,
		"name":          "tftestdatabase",
		"character_set": "utf8",
		"description":   "",
	}

	ra := resourceAttrInit(resourceId, dbDatabaseBasicMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDBDatabase")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := "tf-testAccDBdatabase_basic"
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDBDatabaseConfigDependence)
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
					"instance_id": "${alicloud_db_instance.default.id}",
					"name":        "tftestdatabase",
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
					"description": "from terraform",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{"description": "from terraform"}),
				),
			},
		},
	})

}

func TestAccAlicloudRdsDBDatabaseMySQL1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_db_database.default"

	var dbDatabaseBasicMap = map[string]string{
		"instance_id":   CHECKSET,
		"name":          "tftestdatabase",
		"character_set": "utf8",
		"description":   "",
	}

	ra := resourceAttrInit(resourceId, dbDatabaseBasicMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDBDatabase")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := "tf-testAccDBdatabase_basic"
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDBDatabaseConfigDependence)
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
					"instance_id": "${alicloud_db_instance.default.id}",
					"name":        "tftestdatabase",
					"description": "from terraform",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{"description": "from terraform"}),
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

func resourceDBDatabaseConfigDependence(name string) string {
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
}
`, name)
}
