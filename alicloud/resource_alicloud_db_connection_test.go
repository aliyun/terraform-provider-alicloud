package alicloud

import (
	"fmt"
	"testing"

	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudDBConnection_basic(t *testing.T) {
	var connection rds.DBInstanceNetInfo
	rand := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	connectionStringRegexp := regexp.MustCompile(fmt.Sprintf("^tf-testacc%s.mysql.([a-z-A-Z-0-9]+.){0,1}rds.aliyuncs.com", rand))

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
				Config: testAccDBConnection_basic(RdsCommonTestCase, rand),
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
				Config: testAccDBConnection_update(RdsCommonTestCase, rand),
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
		object, err := rdsService.DescribeDBConnection(rs.Primary.ID)

		if err != nil {
			return WrapError(err)
		}

		*d = *object
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

		_, err := rdsService.DescribeDBConnection(rs.Primary.ID)

		if err != nil {
			if rdsService.NotFoundDBInstance(err) {
				continue
			}
			return WrapError(err)
		}

		return WrapError(fmt.Errorf("Error db connection %s still existed.", rs.Primary.ID))
	}

	return nil
}

func testAccDBConnection_basic(common, rand string) string {
	return fmt.Sprintf(`
	%s
	variable "creation" {
		default = "Rds"
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
	  connection_prefix = "tf-testacc%s"
	}
	`, common, rand)
}
func testAccDBConnection_update(common, rand string) string {
	return fmt.Sprintf(`
	%s
	variable "creation" {
		default = "Rds"
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
	  connection_prefix = "tf-testacc%s"
	  port = 3333
	}
	`, common, rand)
}
