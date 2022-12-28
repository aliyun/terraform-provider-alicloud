package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudRdsDBReadWriteSplittingConnectionMssql_create(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_db_read_write_splitting_connection.create"

	var dbRwSplitConnMap = map[string]string{
		"instance_id":       CHECKSET,
		"distribution_type": "Standard",
	}

	ra := resourceAttrInit(resourceId, dbRwSplitConnMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDBInstanceNetInfo")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := "tf-testAccDBRwSplitConnMssql_create"
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDBReadWriteSplittingConfigDependenceMssql)
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
					"instance_id":       "${alicloud_db_readonly_instance.default.master_db_instance_id}",
					"distribution_type": "Standard",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":       "${alicloud_db_readonly_instance.default.master_db_instance_id}",
					"distribution_type": "Custom",
					"weight": `${map(
						"${alicloud_db_readonly_instance.default.id}", "200",
						"master", "200",
						"slave", "400"
					)}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					rc.checkResourceExists(),
					testAccCheck(map[string]string{
						"weight.%":          "3",
						"distribution_type": "Custom",
					}),
				),
			},
		},
	})

}

func resourceDBReadWriteSplittingConfigDependenceMssql(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
data "alicloud_db_zones" "default"{
	engine                   = "SQLServer"
	engine_version           = "2017_ent"
	instance_charge_type     = "PostPaid"
 	db_instance_storage_type = "cloud_essd"
	category                 = "AlwaysOn"
}

data "alicloud_db_instance_classes" "master" {
    zone_id                  = data.alicloud_db_zones.default.zones.0.id
	engine                   = "SQLServer"
	engine_version           = "2017_ent"
 	db_instance_storage_type = "cloud_essd"
	instance_charge_type     = "PostPaid"
	category                 = "AlwaysOn"
}

data "alicloud_vswitches" "default" {
  zone_id = data.alicloud_db_zones.default.zones.0.id
}

data "alicloud_resource_manager_resource_groups" "default" {
	status = "OK"
}

resource "alicloud_security_group" "default" {
	name   = var.name
	vpc_id = data.alicloud_vswitches.default.vswitches.0.vpc_id
}

resource "alicloud_db_instance" "default" {
    engine                   = "SQLServer"
	engine_version           = "2017_ent"
 	db_instance_storage_type = "cloud_essd"
	instance_type            = data.alicloud_db_instance_classes.master.instance_classes.0.instance_class
	instance_storage         = data.alicloud_db_instance_classes.master.instance_classes.0.storage_range.min
	vswitch_id               = data.alicloud_vswitches.default.vswitches.0.vswitch_id
	instance_name            = var.name
}

resource "alicloud_db_readonly_instance" "default" {
	zone_id               = alicloud_db_instance.default.zone_id
	master_db_instance_id = alicloud_db_instance.default.id
	engine_version        = alicloud_db_instance.default.engine_version
	instance_storage      = alicloud_db_instance.default.instance_storage
	instance_type         = "mssql.x8.large.ro"
	instance_name         = "${var.name}_ro"
	vswitch_id            = data.alicloud_vswitches.default.vswitches.0.vswitch_id
}

`, name)
}
