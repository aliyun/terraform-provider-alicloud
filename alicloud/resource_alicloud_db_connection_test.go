package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAlicloudDBConnection_basic(t *testing.T) {
	var connection rds.DBInstanceNetInfo

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_db_connection.foo",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDBConnectionDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDBConnection_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDBConnectionExists(
						"alicloud_db_connection.foo", &connection),
					resource.TestCheckResourceAttr(
						"alicloud_db_connection.foo",
						"connection_string",
						"test-connection.mysql.rds.aliyuncs.com"),
					resource.TestCheckResourceAttr(
						"alicloud_db_connection.foo",
						"port", "3306"),
				),
			},
			resource.TestStep{
				Config: testAccDBConnection_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDBConnectionExists(
						"alicloud_db_connection.foo", &connection),
					resource.TestCheckResourceAttr(
						"alicloud_db_connection.foo",
						"connection_string",
						"test-connection.mysql.rds.aliyuncs.com"),
					resource.TestCheckResourceAttr(
						"alicloud_db_connection.foo",
						"port", "3333"),
				),
			},
		},
	})

}

func testAccCheckDBConnectionExists(n string, d *rds.DBInstanceNetInfo) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No DB connection ID is set")
		}

		client := testAccProvider.Meta().(*AliyunClient)
		parts := strings.Split(rs.Primary.ID, COLON_SEPARATED)
		conn, err := client.DescribeDBInstanceNetInfoByIpType(parts[0], Public)

		if err != nil {
			return err
		}

		if conn == nil {
			return fmt.Errorf("Connection string is not found in the instance %s.", parts[0])
		}

		*d = *conn
		return nil
	}
}

func testAccCheckDBConnectionDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*AliyunClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_db_connection" {
			continue
		}

		parts := strings.Split(rs.Primary.ID, COLON_SEPARATED)

		conn, err := client.DescribeDBInstanceNetInfoByIpType(parts[0], Public)

		if err != nil {
			if NotFoundError(err) || IsExceptedError(err, InvalidDBInstanceIdNotFound) || IsExceptedError(err, InvalidCurrentConnectionStringNotFound) {
				continue
			}
			return err
		}

		if conn != nil {
			return fmt.Errorf("Error db connection string prefix %s is still existing.", parts[1])
		}
	}

	return nil
}

const testAccDBConnection_basic = `
variable "name" {
	default = "testaccdbconnection_basic"
}
data "alicloud_zones" "default" {
	"available_resource_creation"= "Rds"
}

resource "alicloud_vpc" "foo" {
	name = "${var.name}"
	cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "foo" {
 	vpc_id = "${alicloud_vpc.foo.id}"
 	cidr_block = "172.16.0.0/21"
 	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
}

resource "alicloud_db_instance" "instance" {
	engine = "MySQL"
	engine_version = "5.6"
	instance_type = "rds.mysql.t1.small"
	instance_storage = "10"
  	vswitch_id = "${alicloud_vswitch.foo.id}"
	instance_name = "${var.name}"
}

resource "alicloud_db_connection" "foo" {
  instance_id = "${alicloud_db_instance.instance.id}"
  connection_prefix = "test-connection"
}
`
const testAccDBConnection_update = `
variable "name" {
	default = "testaccdbconnection_basic"
}
data "alicloud_zones" "default" {
	"available_resource_creation"= "Rds"
}

resource "alicloud_vpc" "foo" {
	name = "${var.name}"
	cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "foo" {
 	vpc_id = "${alicloud_vpc.foo.id}"
 	cidr_block = "172.16.0.0/21"
 	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
}

resource "alicloud_db_instance" "instance" {
	engine = "MySQL"
	engine_version = "5.6"
	instance_type = "rds.mysql.t1.small"
	instance_storage = "10"
  	vswitch_id = "${alicloud_vswitch.foo.id}"
  	instance_name = "${var.name}"
}

resource "alicloud_db_connection" "foo" {
  instance_id = "${alicloud_db_instance.instance.id}"
  connection_prefix = "test-connection"
  port = 3333
}
`
