package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudDBConnection_basic(t *testing.T) {
	var connection rds.DBInstanceNetInfo

	connectionStringRegexp := regexp.MustCompile("^test-connection.mysql.([a-z-A-Z-0-9]+.){0,1}rds.aliyuncs.com")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_db_connection.foo",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDBConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDBConnection_basic(DatabaseCommonTestCase),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDBConnectionExists(
						"alicloud_db_connection.foo", &connection),
					resource.TestMatchResourceAttr(
						"alicloud_db_connection.foo",
						"connection_string",
						connectionStringRegexp),
					resource.TestCheckResourceAttr(
						"alicloud_db_connection.foo",
						"port", "3306"),
				),
			},
			{
				Config: testAccDBConnection_update(DatabaseCommonTestCase),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDBConnectionExists(
						"alicloud_db_connection.foo", &connection),
					resource.TestMatchResourceAttr(
						"alicloud_db_connection.foo",
						"connection_string",
						connectionStringRegexp),
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

		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		rdsService := RdsService{client}
		parts := strings.Split(rs.Primary.ID, COLON_SEPARATED)
		conn, err := rdsService.DescribeDBInstanceNetInfoByIpType(parts[0], Public)

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
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	rdsService := RdsService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_db_connection" {
			continue
		}

		parts := strings.Split(rs.Primary.ID, COLON_SEPARATED)

		conn, err := rdsService.DescribeDBInstanceNetInfoByIpType(parts[0], Public)

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

func testAccDBConnection_basic(common string) string {
	return fmt.Sprintf(`
	%s
	variable "creation" {
		default = "Rds"
	}
	variable "multi_az" {
		default = "false"
	}
	variable "name" {
		default = "tf-testAccDBconnection_basic"
	}

	resource "alicloud_db_instance" "instance" {
		engine = "MySQL"
		engine_version = "5.6"
		instance_type = "rds.mysql.t1.small"
		instance_storage = "10"
		  vswitch_id = "${alicloud_vswitch.default.id}"
		instance_name = "${var.name}"
	}

	resource "alicloud_db_connection" "foo" {
	  instance_id = "${alicloud_db_instance.instance.id}"
	  connection_prefix = "test-connection"
	}
	`, common)
}
func testAccDBConnection_update(common string) string {
	return fmt.Sprintf(`
	%s
	variable "creation" {
		default = "Rds"
	}
	variable "multi_az" {
		default = "false"
	}
	variable "name" {
		default = "tf-testAccDBconnection_basic"
	}

	resource "alicloud_db_instance" "instance" {
		engine = "MySQL"
		engine_version = "5.6"
		instance_type = "rds.mysql.t1.small"
		instance_storage = "10"
		vswitch_id = "${alicloud_vswitch.default.id}"
		instance_name = "${var.name}"
	}

	resource "alicloud_db_connection" "foo" {
	  instance_id = "${alicloud_db_instance.instance.id}"
	  connection_prefix = "test-connection"
	  port = 3333
	}
	`, common)
}
