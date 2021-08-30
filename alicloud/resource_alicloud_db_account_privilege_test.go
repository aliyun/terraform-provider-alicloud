package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudDBAccountPrivilege_mysql(t *testing.T) {

	var v map[string]interface{}
	name := "tf-testAccDBAccountPrivilege_mysql"
	resourceId := "alicloud_db_account_privilege.default"
	var basicMap = map[string]string{
		"instance_id":  CHECKSET,
		"account_name": "tftestprivilege",
		"privilege":    "ReadOnly",
		"db_names.#":   "2",
	}
	ra := resourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeDBAccountPrivilege")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDBAccountPrivilegeConfigDependenceForMySql)

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
					"instance_id":  "${alicloud_db_instance.default.id}",
					"account_name": "${alicloud_db_account.default.name}",
					"privilege":    "ReadOnly",
					"db_names":     "${alicloud_db_database.default.*.name}",
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
					"instance_id":  "${alicloud_db_instance.default.id}",
					"account_name": "${alicloud_db_account.default.name}",
					"privilege":    "ReadOnly",
					"db_names":     []string{"${alicloud_db_database.default.0.name}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_names.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":  "${alicloud_db_instance.default.id}",
					"account_name": "${alicloud_db_account.default.name}",
					"privilege":    "ReadOnly",
					"db_names":     "${alicloud_db_database.default.*.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_names.#": "2",
					}),
				),
			},
		},
	})

}

func TestAccAlicloudDBAccountPrivilege_PostgreSql(t *testing.T) {

	var v map[string]interface{}
	name := "tf-testAccDBAccountPrivilege_PostgreSql"
	resourceId := "alicloud_db_account_privilege.default"
	var basicMap = map[string]string{
		"instance_id":  CHECKSET,
		"account_name": "tftestprivilege",
		"privilege":    "DBOwner",
		"db_names.#":   "1",
	}
	ra := resourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeDBAccountPrivilege")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDBAccountPrivilegeConfigDependenceForPostgreSql)

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
					"instance_id":  "${alicloud_db_instance.default.id}",
					"account_name": "${alicloud_db_account.default.name}",
					"privilege":    "DBOwner",
					"db_names":     []string{"${alicloud_db_database.default.0.name}"},
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
					"instance_id":  "${alicloud_db_instance.default.id}",
					"account_name": "${alicloud_db_account.default.name}",
					"privilege":    "DBOwner",
					"db_names":     "${alicloud_db_database.default.*.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_names.#": "2",
					}),
				),
			},
		},
	})

}

func resourceDBAccountPrivilegeConfigDependenceForMySql(name string) string {
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

	resource "alicloud_db_instance" "default" {
		engine = "${data.alicloud_db_instance_engines.default.instance_engines.0.engine}"
		engine_version = "${data.alicloud_db_instance_engines.default.instance_engines.0.engine_version}"
		instance_type = "${data.alicloud_db_instance_classes.default.instance_classes.0.instance_class}"
		instance_storage = "${data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.min}"
		vswitch_id = "${alicloud_vswitch.default.id}"
		instance_name = "${var.name}"
	}

	resource "alicloud_db_database" "default" {
	  count = 2
	  instance_id = "${alicloud_db_instance.default.id}"
	  name = "tfaccountpri_${count.index}"
	  description = "from terraform"
	}

	resource "alicloud_db_account" "default" {
	  instance_id = "${alicloud_db_instance.default.id}"
	  name = "tftestprivilege"
	  password = "Test12345"
	  description = "from terraform"
	}
`, RdsCommonTestCase, name)
}

func resourceDBAccountPrivilegeConfigDependenceForPostgreSql(name string) string {
	return fmt.Sprintf(`
%s
	variable "creation" {
		default = "Rds"
	}

	variable "name" {
		default = "%s"
	}
	
	data "alicloud_db_instance_classes" "default" {
		instance_charge_type = "PostPaid"
		engine               = "PostgreSQL"
		engine_version       = "10.0"
		storage_type         = "cloud_ssd"
	}

	resource "alicloud_db_instance" "default" {
		engine = "PostgreSQL"
		engine_version = "10.0"
		instance_type = "${data.alicloud_db_instance_classes.default.instance_classes.0.instance_class}"
		instance_storage = "${data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.min}"
		vswitch_id = "${alicloud_vswitch.default.id}"
		instance_name = "${var.name}"
	}

	resource "alicloud_db_database" "default" {
	  count = 2
	  instance_id = "${alicloud_db_instance.default.id}"
	  name = "tfaccountpri_${count.index}"
	  description = "from terraform"
      character_set = "UTF8"
	}

	resource "alicloud_db_account" "default" {
	  instance_id = "${alicloud_db_instance.default.id}"
	  name = "tftestprivilege"
	  password = "Test12345"
	  description = "from terraform"
	}
`, RdsCommonTestCase, name)
}
