package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/hashcode"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudDBReadonlyInstance_vpc(t *testing.T) {
	var instance rds.DBInstanceAttribute

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_db_readonly_instance.foo",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDBReadonlyInstanceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDBReadonlyInstance_vpc(testAccDBRInstance_vpc(RdsCommonTestCase)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDBInstanceExists(
						"alicloud_db_readonly_instance.foo", &instance),
					resource.TestCheckResourceAttr(
						"alicloud_db_readonly_instance.foo", "instance_storage", "30"),
					resource.TestCheckResourceAttr(
						"alicloud_db_readonly_instance.foo", "engine_version", "5.6"),
					resource.TestCheckResourceAttr(
						"alicloud_db_readonly_instance.foo", "engine", "MySQL"),
					resource.TestCheckResourceAttr(
						"alicloud_db_readonly_instance.foo", "port", "3306"),
					resource.TestCheckResourceAttr(
						"alicloud_db_readonly_instance.foo", "instance_name", "tf-testAccDBInstance_vpc_ro"),
					resource.TestCheckResourceAttr(
						"alicloud_db_readonly_instance.foo", "instance_type", "rds.mysql.t1.small"),
					resource.TestCheckNoResourceAttr(
						"alicloud_db_readonly_instance.foo", "parameters"),
					resource.TestCheckResourceAttrSet(
						"alicloud_db_readonly_instance.foo", "master_db_instance_id"),
					resource.TestCheckResourceAttrSet(
						"alicloud_db_readonly_instance.foo", "zone_id"),
					resource.TestCheckResourceAttrSet(
						"alicloud_db_readonly_instance.foo", "vswitch_id"),
					resource.TestCheckResourceAttrSet(
						"alicloud_db_readonly_instance.foo", "connection_string"),
				),
			},
		},
	})

}

func TestAccAlicloudDBReadonlyInstance_multiAZ_classic(t *testing.T) {
	var instance rds.DBInstanceAttribute

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_db_readonly_instance.foo",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDBReadonlyInstanceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDBReadonlyInstance_multiAZ(testAccDBInstance_multiAZ),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDBInstanceExists(
						"alicloud_db_readonly_instance.foo", &instance),
					resource.TestCheckResourceAttr(
						"alicloud_db_readonly_instance.foo", "instance_storage", "30"),
					resource.TestCheckResourceAttr(
						"alicloud_db_readonly_instance.foo", "engine_version", "5.6"),
					resource.TestCheckResourceAttr(
						"alicloud_db_readonly_instance.foo", "engine", "MySQL"),
					resource.TestCheckResourceAttr(
						"alicloud_db_readonly_instance.foo", "port", "3306"),
					resource.TestCheckResourceAttr(
						"alicloud_db_readonly_instance.foo", "instance_name", "tf-testAccDBInstance_multiAZ_ro"),
					resource.TestCheckResourceAttr(
						"alicloud_db_readonly_instance.foo", "instance_type", "rds.mysql.t1.small"),
					resource.TestCheckNoResourceAttr(
						"alicloud_db_readonly_instance.foo", "parameters"),
					resource.TestCheckResourceAttrSet(
						"alicloud_db_readonly_instance.foo", "master_db_instance_id"),
					resource.TestCheckResourceAttrSet(
						"alicloud_db_readonly_instance.foo", "zone_id"),
					resource.TestCheckResourceAttr(
						"alicloud_db_readonly_instance.foo", "vswitch_id", ""),
					resource.TestCheckResourceAttrSet(
						"alicloud_db_readonly_instance.foo", "connection_string"),
				),
			},
		},
	})

}

func TestAccAlicloudDBReadonlyInstance_multiAZ_vpc(t *testing.T) {
	var instance rds.DBInstanceAttribute

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, false, connectivity.RdsMultiAzNoSupportedRegions)
		},

		// module name
		IDRefreshName: "alicloud_db_readonly_instance.foo",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDBReadonlyInstanceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDBReadonlyInstance_multiAZ_vpc(testAccDBInstance_vpc_multiAZ(DBMultiAZCommonTestCase)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDBInstanceExists(
						"alicloud_db_readonly_instance.foo", &instance),
					resource.TestCheckResourceAttr(
						"alicloud_db_readonly_instance.foo", "instance_storage", "30"),
					resource.TestCheckResourceAttr(
						"alicloud_db_readonly_instance.foo", "engine_version", "5.6"),
					resource.TestCheckResourceAttr(
						"alicloud_db_readonly_instance.foo", "engine", "MySQL"),
					resource.TestCheckResourceAttr(
						"alicloud_db_readonly_instance.foo", "port", "3306"),
					resource.TestCheckResourceAttr(
						"alicloud_db_readonly_instance.foo", "instance_name", "tf-testAccDBInstance_multiAZ_ro"),
					resource.TestCheckResourceAttr(
						"alicloud_db_readonly_instance.foo", "instance_type", "rds.mysql.t1.small"),
					resource.TestCheckNoResourceAttr(
						"alicloud_db_readonly_instance.foo", "parameters"),
					resource.TestCheckResourceAttrSet(
						"alicloud_db_readonly_instance.foo", "master_db_instance_id"),
					resource.TestCheckResourceAttrSet(
						"alicloud_db_readonly_instance.foo", "zone_id"),
					resource.TestCheckResourceAttrSet(
						"alicloud_db_readonly_instance.foo", "vswitch_id"),
					resource.TestCheckResourceAttrSet(
						"alicloud_db_readonly_instance.foo", "connection_string"),
				),
			},
		},
	})

}

func TestAccAlicloudDBReadonlyInstance_upgrade(t *testing.T) {
	var instance rds.DBInstanceAttribute

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_db_readonly_instance.foo",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDBReadonlyInstanceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDBReadonlyInstance_vpc(testAccDBRInstance_vpc(RdsCommonTestCase)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDBInstanceExists(
						"alicloud_db_readonly_instance.foo", &instance),
					resource.TestCheckResourceAttr(
						"alicloud_db_readonly_instance.foo", "instance_storage", "30"),
					resource.TestCheckResourceAttr(
						"alicloud_db_readonly_instance.foo", "engine_version", "5.6"),
					resource.TestCheckResourceAttr(
						"alicloud_db_readonly_instance.foo", "engine", "MySQL"),
					resource.TestCheckResourceAttr(
						"alicloud_db_readonly_instance.foo", "port", "3306"),
					resource.TestCheckResourceAttr(
						"alicloud_db_readonly_instance.foo", "instance_name", "tf-testAccDBInstance_vpc_ro"),
					resource.TestCheckResourceAttr(
						"alicloud_db_readonly_instance.foo", "instance_type", "rds.mysql.t1.small"),
					resource.TestCheckNoResourceAttr(
						"alicloud_db_readonly_instance.foo", "parameters"),
					resource.TestCheckResourceAttrSet(
						"alicloud_db_readonly_instance.foo", "master_db_instance_id"),
					resource.TestCheckResourceAttrSet(
						"alicloud_db_readonly_instance.foo", "zone_id"),
					resource.TestCheckResourceAttrSet(
						"alicloud_db_readonly_instance.foo", "vswitch_id"),
					resource.TestCheckResourceAttrSet(
						"alicloud_db_readonly_instance.foo", "connection_string"),
				),
			},
			// upgrade storage
			resource.TestStep{
				Config: testAccDBReadonlyInstance_upgrade(testAccDBRInstance_vpc(RdsCommonTestCase),
					"${alicloud_db_instance.foo.instance_type}", "40"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDBInstanceExists(
						"alicloud_db_readonly_instance.foo", &instance),
					resource.TestCheckResourceAttr(
						"alicloud_db_readonly_instance.foo", "instance_storage", "40"),
				),
			},
			// upgrade instanceType
			resource.TestStep{
				Config: testAccDBReadonlyInstance_upgrade(testAccDBRInstance_vpc(RdsCommonTestCase),
					"rds.mysql.s1.small", "40"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDBInstanceExists(
						"alicloud_db_readonly_instance.foo", &instance),
					resource.TestCheckResourceAttr(
						"alicloud_db_readonly_instance.foo", "instance_storage", "40"),
					resource.TestCheckResourceAttr(
						"alicloud_db_readonly_instance.foo", "instance_type", "rds.mysql.s1.small"),
				),
			},
			// upgrade storage and instanceType
			resource.TestStep{
				Config: testAccDBReadonlyInstance_upgrade(testAccDBRInstance_vpc(RdsCommonTestCase),
					"rds.mysql.s1.large", "50"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDBInstanceExists(
						"alicloud_db_readonly_instance.foo", &instance),
					resource.TestCheckResourceAttr(
						"alicloud_db_readonly_instance.foo", "instance_storage", "50"),
					resource.TestCheckResourceAttr(
						"alicloud_db_readonly_instance.foo", "instance_type", "rds.mysql.s1.large"),
				),
			},
		},
	})

}

func TestAccAlicloudDBReadonlyInstance_parameter(t *testing.T) {
	var instance rds.DBInstanceAttribute

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_db_readonly_instance.foo",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDBReadonlyInstanceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDBReadonlyInstance_vpc(testAccDBRInstance_vpc(RdsCommonTestCase)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDBInstanceExists(
						"alicloud_db_readonly_instance.foo", &instance),
					resource.TestCheckResourceAttr(
						"alicloud_db_readonly_instance.foo", "instance_storage", "30"),
					resource.TestCheckResourceAttr(
						"alicloud_db_readonly_instance.foo", "engine_version", "5.6"),
					resource.TestCheckResourceAttr(
						"alicloud_db_readonly_instance.foo", "engine", "MySQL"),
					resource.TestCheckResourceAttr(
						"alicloud_db_readonly_instance.foo", "port", "3306"),
					resource.TestCheckResourceAttr(
						"alicloud_db_readonly_instance.foo", "instance_name", "tf-testAccDBInstance_vpc_ro"),
					resource.TestCheckResourceAttr(
						"alicloud_db_readonly_instance.foo", "instance_type", "rds.mysql.t1.small"),
					resource.TestCheckNoResourceAttr(
						"alicloud_db_readonly_instance.foo", "parameters"),
					resource.TestCheckResourceAttrSet(
						"alicloud_db_readonly_instance.foo", "master_db_instance_id"),
					resource.TestCheckResourceAttrSet(
						"alicloud_db_readonly_instance.foo", "zone_id"),
					resource.TestCheckResourceAttrSet(
						"alicloud_db_readonly_instance.foo", "vswitch_id"),
					resource.TestCheckResourceAttrSet(
						"alicloud_db_readonly_instance.foo", "connection_string"),
				),
			},
			// update parameter
			resource.TestStep{
				Config: testAccDBReadonlyInstance_parameter(testAccDBRInstance_vpc(RdsCommonTestCase)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_db_readonly_instance.foo",
						fmt.Sprintf("parameters.%d.value",
							hashcode.String("innodb_large_prefix")), "ON"),
				),
			},
			// update multi parameter
			resource.TestStep{
				Config: testAccDBReadonlyInstance_parameterMulti(testAccDBRInstance_vpc(RdsCommonTestCase)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_db_readonly_instance.foo",
						fmt.Sprintf("parameters.%d.value",
							hashcode.String("innodb_large_prefix")), "ON"),
					resource.TestCheckResourceAttr("alicloud_db_readonly_instance.foo",
						fmt.Sprintf("parameters.%d.value",
							hashcode.String("connect_timeout")), "50"),
				),
			},
			// remove parameter definition, parameter value not change
			resource.TestStep{
				Config: testAccDBReadonlyInstance_parameter(testAccDBRInstance_vpc(RdsCommonTestCase)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_db_readonly_instance.foo",
						fmt.Sprintf("parameters.%d.value",
							hashcode.String("innodb_large_prefix")), "ON"),
					testAccCheckDBParameterExpects(
						"alicloud_db_readonly_instance.foo", "connect_timeout", "50"),
				),
			},
		},
	})

}

func TestAccAlicloudDBReadonlyInstance_updateInstanceName(t *testing.T) {
	var instance rds.DBInstanceAttribute

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_db_readonly_instance.foo",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDBReadonlyInstanceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDBReadonlyInstance_vpc(testAccDBRInstance_vpc(RdsCommonTestCase)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDBInstanceExists(
						"alicloud_db_readonly_instance.foo", &instance),
					resource.TestCheckResourceAttr(
						"alicloud_db_readonly_instance.foo", "instance_storage", "30"),
					resource.TestCheckResourceAttr(
						"alicloud_db_readonly_instance.foo", "engine_version", "5.6"),
					resource.TestCheckResourceAttr(
						"alicloud_db_readonly_instance.foo", "engine", "MySQL"),
					resource.TestCheckResourceAttr(
						"alicloud_db_readonly_instance.foo", "port", "3306"),
					resource.TestCheckResourceAttr(
						"alicloud_db_readonly_instance.foo", "instance_name", "tf-testAccDBInstance_vpc_ro"),
					resource.TestCheckResourceAttr(
						"alicloud_db_readonly_instance.foo", "instance_type", "rds.mysql.t1.small"),
					resource.TestCheckNoResourceAttr(
						"alicloud_db_readonly_instance.foo", "parameters"),
					resource.TestCheckResourceAttrSet(
						"alicloud_db_readonly_instance.foo", "master_db_instance_id"),
					resource.TestCheckResourceAttrSet(
						"alicloud_db_readonly_instance.foo", "zone_id"),
					resource.TestCheckResourceAttrSet(
						"alicloud_db_readonly_instance.foo", "vswitch_id"),
					resource.TestCheckResourceAttrSet(
						"alicloud_db_readonly_instance.foo", "connection_string"),
				),
			},
			// update instanceName
			resource.TestStep{
				Config: testAccDBReadonlyInstance_vpc_instanceName(
					testAccDBRInstance_vpc(RdsCommonTestCase), "some_other_name"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDBInstanceExists(
						"alicloud_db_readonly_instance.foo", &instance),
					resource.TestCheckResourceAttr(
						"alicloud_db_readonly_instance.foo", "instance_storage", "30"),
					resource.TestCheckResourceAttr(
						"alicloud_db_readonly_instance.foo", "engine_version", "5.6"),
					resource.TestCheckResourceAttr(
						"alicloud_db_readonly_instance.foo", "engine", "MySQL"),
					resource.TestCheckResourceAttr(
						"alicloud_db_readonly_instance.foo", "port", "3306"),
					resource.TestCheckResourceAttr(
						"alicloud_db_readonly_instance.foo", "instance_name", "some_other_name"),
					resource.TestCheckResourceAttr(
						"alicloud_db_readonly_instance.foo", "instance_type", "rds.mysql.t1.small"),
					resource.TestCheckNoResourceAttr(
						"alicloud_db_readonly_instance.foo", "parameters"),
					resource.TestCheckResourceAttrSet(
						"alicloud_db_readonly_instance.foo", "master_db_instance_id"),
					resource.TestCheckResourceAttrSet(
						"alicloud_db_readonly_instance.foo", "zone_id"),
					resource.TestCheckResourceAttrSet(
						"alicloud_db_readonly_instance.foo", "vswitch_id"),
					resource.TestCheckResourceAttrSet(
						"alicloud_db_readonly_instance.foo", "connection_string"),
				),
			},
		},
	})

}

func testAccDBRInstance_vpc(common string) string {
	return fmt.Sprintf(`
	%s
	variable "creation" {
		default = "Rds"
	}
	variable "multi_az" {
		default = "false"
	}
	variable "name" {
		default = "tf-testAccDBInstance_vpc"
	}

	resource "alicloud_db_instance" "foo" {
		engine = "MySQL"
		engine_version = "5.6"
		instance_type = "rds.mysql.t1.small"
		instance_storage = "20"
		instance_charge_type = "Postpaid"
		instance_name = "${var.name}"
		vswitch_id = "${alicloud_vswitch.default.id}"
		security_ips = ["10.168.1.12", "100.69.7.112"]
	}
	`, common)
}

func testAccDBReadonlyInstance_vpc(common string) string {
	return fmt.Sprintf(`
	%s
	resource "alicloud_db_readonly_instance" "foo" {
		master_db_instance_id = "${alicloud_db_instance.foo.id}"
		zone_id = "${alicloud_db_instance.foo.zone_id}"
		engine_version = "${alicloud_db_instance.foo.engine_version}"
		instance_type = "${alicloud_db_instance.foo.instance_type}"
		instance_storage = "30"
		instance_name = "${var.name}_ro"
		vswitch_id = "${alicloud_vswitch.default.id}"
	}
	`, common)
}

func testAccDBReadonlyInstance_vpc_instanceName(common, instanceName string) string {
	return fmt.Sprintf(`
	%s
	resource "alicloud_db_readonly_instance" "foo" {
		master_db_instance_id = "${alicloud_db_instance.foo.id}"
		zone_id = "${alicloud_db_instance.foo.zone_id}"
		engine_version = "${alicloud_db_instance.foo.engine_version}"
		instance_type = "${alicloud_db_instance.foo.instance_type}"
		instance_storage = "30"
		instance_name = "%s"
		vswitch_id = "${alicloud_vswitch.default.id}"
	}
	`, common, instanceName)
}

func testAccDBReadonlyInstance_upgrade(common, instanceType, storage string) string {
	return fmt.Sprintf(`
	%s
	resource "alicloud_db_readonly_instance" "foo" {
		master_db_instance_id = "${alicloud_db_instance.foo.id}"
		zone_id = "${alicloud_db_instance.foo.zone_id}"
		engine_version = "${alicloud_db_instance.foo.engine_version}"
		instance_type = "%s"
		instance_storage = "%s"
		instance_name = "${var.name}_ro"
		vswitch_id = "${alicloud_vswitch.default.id}"
	}
	`, common, instanceType, storage)
}

func testAccDBReadonlyInstance_multiAZ(common string) string {
	return fmt.Sprintf(`
	%s
	resource "alicloud_db_readonly_instance" "foo" {
		master_db_instance_id = "${alicloud_db_instance.foo.id}"
		zone_id = "${alicloud_db_instance.foo.zone_id}"
		engine_version = "${alicloud_db_instance.foo.engine_version}"
		instance_type = "${alicloud_db_instance.foo.instance_type}"
		instance_storage = "30"
		instance_name = "${var.name}_ro"
	}
	`, common)
}

func testAccDBReadonlyInstance_multiAZ_vpc(common string) string {
	return fmt.Sprintf(`
	%s
	resource "alicloud_db_readonly_instance" "foo" {
		master_db_instance_id = "${alicloud_db_instance.foo.id}"
		zone_id = "${alicloud_db_instance.foo.zone_id}"
		engine_version = "${alicloud_db_instance.foo.engine_version}"
		instance_type = "${alicloud_db_instance.foo.instance_type}"
		instance_storage = "30"
		instance_name = "${var.name}_ro"
		vswitch_id = "${alicloud_vswitch.default.id}"
	}
	`, common)
}

func testAccDBReadonlyInstance_parameter(common string) string {
	return fmt.Sprintf(`
	%s
	resource "alicloud_db_readonly_instance" "foo" {
		master_db_instance_id = "${alicloud_db_instance.foo.id}"
		zone_id = "${alicloud_db_instance.foo.zone_id}"
		engine_version = "${alicloud_db_instance.foo.engine_version}"
		instance_type = "${alicloud_db_instance.foo.instance_type}"
		instance_storage = "30"
		instance_name = "${var.name}_ro"
		vswitch_id = "${alicloud_vswitch.default.id}"
		parameters = [{
			name = "innodb_large_prefix"
			value = "ON"
		}]
	}
	`, common)
}

func testAccDBReadonlyInstance_parameterMulti(common string) string {
	return fmt.Sprintf(`
	%s
	resource "alicloud_db_readonly_instance" "foo" {
		master_db_instance_id = "${alicloud_db_instance.foo.id}"
		zone_id = "${alicloud_db_instance.foo.zone_id}"
		engine_version = "${alicloud_db_instance.foo.engine_version}"
		instance_type = "${alicloud_db_instance.foo.instance_type}"
		instance_storage = "30"
		instance_name = "${var.name}_ro"
		vswitch_id = "${alicloud_vswitch.default.id}"
		parameters = [{
			name = "innodb_large_prefix"
			value = "ON"
		},{
			name = "connect_timeout"
			value = "50"
		}]
	}
	`, common)
}

func testAccCheckDBReadonlyInstanceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	rdsService := RdsService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_db_readonly_instance" {
			continue
		}

		ins, err := rdsService.DescribeDBInstanceById(rs.Primary.ID)

		if ins != nil {
			return fmt.Errorf("Error DB Instance still exist")
		}

		// Verify the error is what we want
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return err
		}
	}

	return nil
}
