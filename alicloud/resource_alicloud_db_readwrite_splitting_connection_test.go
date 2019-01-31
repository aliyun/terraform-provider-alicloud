package alicloud

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudDBReadWriteSplittingConnection_basic(t *testing.T) {
	var connection rds.DBInstanceNetInfo

	connectionStringRegexp := regexp.MustCompile("^test-connection.mysql.([a-z-A-Z-0-9]+.){0,1}rds.aliyuncs.com")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_db_read_write_splitting_connection.foo",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDBReadWriteSplittingConnectionDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDBReadWriteSplittingConnection_basic(RdsCommonTestCase),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDBReadWriteSplittingConnectionExists(
						"alicloud_db_read_write_splitting_connection.foo", &connection),
					resource.TestMatchResourceAttr(
						"alicloud_db_read_write_splitting_connection.foo",
						"connection_string",
						connectionStringRegexp),
					resource.TestCheckResourceAttr(
						"alicloud_db_read_write_splitting_connection.foo",
						"port", "3306"),
				),
			},
			resource.TestStep{
				Config: testAccDBReadWriteSplittingConnection_update(RdsCommonTestCase),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDBReadWriteSplittingConnectionExists(
						"alicloud_db_read_write_splitting_connection.foo", &connection),
					resource.TestMatchResourceAttr(
						"alicloud_db_read_write_splitting_connection.foo",
						"connection_string",
						connectionStringRegexp),
					resource.TestCheckResourceAttr(
						"alicloud_db_read_write_splitting_connection.foo",
						"max_delay_time", "300"),
					resource.TestMatchResourceAttr(
						"alicloud_db_read_write_splitting_connection.foo",
						"weight", regexp.MustCompile(".+500.+")),
				),
			},
		},
	})

}

func testAccCheckDBReadWriteSplittingConnectionExists(n string, d *rds.DBInstanceNetInfo) resource.TestCheckFunc {
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

		resp, err := rdsService.DescribeDBInstanceNetInfos(parts[0])

		if err != nil {
			return err
		}

		if resp == nil {
			return GetNotFoundErrorFromString(fmt.Sprintf("DB instance %s does not have any connection.", parts[0]))
		}

		found := false
		for _, conn := range resp {
			if conn.ConnectionStringType == "ReadWriteSplitting" {
				*d = conn
				found = true
				break
			}
		}

		if !found {
			return fmt.Errorf("Connection string is not found in the instance %s.", parts[0])
		}

		return nil
	}
}

func testAccCheckDBReadWriteSplittingConnectionDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	rdsService := RdsService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_db_read_write_splitting_connection" {
			continue
		}

		parts := strings.Split(rs.Primary.ID, COLON_SEPARATED)

		resp, err := rdsService.DescribeDBInstanceNetInfos(parts[0])

		if err != nil {
			return err
		}

		if resp == nil {
			return nil
		}

		found := false
		for _, conn := range resp {
			if conn.ConnectionStringType == "ReadWriteSplitting" {
				found = true
				break
			}
		}

		if found {
			return fmt.Errorf("Error db connection string prefix %s is still existing.", parts[1])
		}
	}

	return nil
}

func testAccDBReadWriteSplittingConnection_basic(common string) string {
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

	resource "alicloud_db_instance" "foo" {
		engine = "MySQL"
		engine_version = "5.6"
		instance_type = "rds.mysql.t1.small"
		instance_storage = "10"
		vswitch_id = "${alicloud_vswitch.default.id}"
		instance_name = "${var.name}"
	}

	resource "alicloud_db_readonly_instance" "foo" {
		master_db_instance_id = "${alicloud_db_instance.foo.id}"
		zone_id = "${alicloud_db_instance.foo.zone_id}"
		engine_version = "${alicloud_db_instance.foo.engine_version}"
		instance_type = "${alicloud_db_instance.foo.instance_type}"
		instance_storage = "10"
		instance_name = "${var.name}ro"
		vswitch_id = "${alicloud_vswitch.default.id}"
	}

	resource "alicloud_db_read_write_splitting_connection" "foo" {
	    instance_id = "${alicloud_db_instance.foo.id}"
	    connection_prefix = "test-connection"
		distribution_type = "Standard"
		
		depends_on = ["alicloud_db_readonly_instance.foo"]
	}
	`, common)
}
func testAccDBReadWriteSplittingConnection_update(common string) string {
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

	resource "alicloud_db_instance" "foo" {
		engine = "MySQL"
		engine_version = "5.6"
		instance_type = "rds.mysql.t1.small"
		instance_storage = "10"
		vswitch_id = "${alicloud_vswitch.default.id}"
		instance_name = "${var.name}"
	}

	resource "alicloud_db_readonly_instance" "foo" {
		master_db_instance_id = "${alicloud_db_instance.foo.id}"
		zone_id = "${alicloud_db_instance.foo.zone_id}"
		engine_version = "${alicloud_db_instance.foo.engine_version}"
		instance_type = "${alicloud_db_instance.foo.instance_type}"
		instance_storage = "10"
		instance_name = "${var.name}ro"
		vswitch_id = "${alicloud_vswitch.default.id}"
	}

	resource "alicloud_db_read_write_splitting_connection" "foo" {
	    instance_id = "${alicloud_db_instance.foo.id}"
	    connection_prefix = "test-connection"
		distribution_type = "Custom"
		max_delay_time = 300
		weight = "{\"${alicloud_db_readonly_instance.foo.id}\":\"500\"}"
		
		depends_on = ["alicloud_db_readonly_instance.foo"]
	}
	`, common)
}
