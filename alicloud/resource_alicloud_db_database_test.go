package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudDBDatabase_update(t *testing.T) {
	var database *rds.Database
	resourceId := "alicloud_db_database.default"

	var dbDatabaseBasicMap = map[string]string{
		"instance_id":   CHECKSET,
		"name":          "tftestdatabase",
		"character_set": "utf8",
		"description":   "",
	}

	ra := resourceAttrInit(resourceId, dbDatabaseBasicMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &database, func() interface{} {
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
					"instance_id": "${alicloud_db_instance.instance.id}",
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

func resourceDBDatabaseConfigDependence(name string) string {
	return fmt.Sprintf(`
	%s
	variable "creation" {
		default = "Rds"
	}

	variable "name" {
		default = "%s"
	}

	data "alicloud_db_instance_engines" "default" {
  		instance_charge_type = "PostPaid"
  		engine               = "MySQL"
  		engine_version       = "5.6"
	}

	data "alicloud_db_instance_classes" "default" {
 	 	engine = "${data.alicloud_db_instance_engines.default.instance_engines.0.engine}"
		engine_version = "${data.alicloud_db_instance_engines.default.instance_engines.0.engine_version}"
	}

	resource "alicloud_db_instance" "instance" {
		engine = "${data.alicloud_db_instance_engines.default.instance_engines.0.engine}"
		engine_version = "${data.alicloud_db_instance_engines.default.instance_engines.0.engine_version}"
		instance_type = "${data.alicloud_db_instance_classes.default.instance_classes.0.instance_class}"
		instance_storage = "${data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.min}"
		vswitch_id = "${alicloud_vswitch.default.id}"
		instance_name = "${var.name}"
	}`, RdsCommonTestCase, name)
}
