package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func init() {
	resource.AddTestSweepers("alicloud_db_readonly_instance", &resource.Sweeper{
		Name: "alicloud_db_readonly_instance",
		F:    testSweepDBInstances,
	})
}

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
				Config: testAccDBReadonlyInstance_vpc(RdsCommonTestCase),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDBInstanceExists(
						"alicloud_db_readonly_instance.foo", &instance),
					resource.TestCheckResourceAttr(
						"alicloud_db_readonly_instance.foo",
						"instance_storage",
						"30"),
					resource.TestCheckResourceAttr(
						"alicloud_db_readonly_instance.foo",
						"engine_version",
						"5.6"),
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
				Config: testAccDBReadonlyInstance_vpc(RdsCommonTestCase),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDBInstanceExists(
						"alicloud_db_readonly_instance.foo", &instance),
					resource.TestCheckResourceAttr(
						"alicloud_db_readonly_instance.foo",
						"instance_storage",
						"30"),
					resource.TestCheckResourceAttr(
						"alicloud_db_readonly_instance.foo",
						"engine_version",
						"5.6"),
				),
			},
			resource.TestStep{
				Config: testAccDBReadonlyInstance_upgrade(RdsCommonTestCase),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDBInstanceExists(
						"alicloud_db_readonly_instance.foo", &instance),
					resource.TestCheckResourceAttr(
						"alicloud_db_readonly_instance.foo",
						"instance_storage",
						"40"),
					resource.TestCheckResourceAttr(
						"alicloud_db_readonly_instance.foo",
						"engine_version",
						"5.6"),
				),
			},
		},
	})

}

func testAccDBReadonlyInstance_vpc(common string) string {
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
		instance_type = "rds.mysql.s2.large"
		instance_storage = "20"
		instance_charge_type = "Postpaid"
		instance_name = "${var.name}"
		vswitch_id = "${alicloud_vswitch.default.id}"
		security_ips = ["10.168.1.12", "100.69.7.112"]
	}

	resource "alicloud_db_readonly_instance" "foo" {
		master_db_instance_id = "${alicloud_db_instance.foo.id}"
		zone_id = "${alicloud_db_instance.foo.zone_id}"
		engine_version = "${alicloud_db_instance.foo.engine_version}"
		instance_type = "${alicloud_db_instance.foo.instance_type}"
		instance_storage = "30"
		instance_name = "${var.name}ro"
		vswitch_id = "${alicloud_vswitch.default.id}"
	}
	`, common)
}

func testAccDBReadonlyInstance_upgrade(common string) string {
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
		instance_type = "rds.mysql.s2.large"
		instance_storage = "20"
		instance_charge_type = "Postpaid"
		instance_name = "${var.name}"
		vswitch_id = "${alicloud_vswitch.default.id}"
		security_ips = ["10.168.1.12", "100.69.7.112"]
	}

	resource "alicloud_db_readonly_instance" "foo" {
		master_db_instance_id = "${alicloud_db_instance.foo.id}"
		zone_id = "${alicloud_db_instance.foo.zone_id}"
		engine_version = "${alicloud_db_instance.foo.engine_version}"
		instance_type = "${alicloud_db_instance.foo.instance_type}"
		instance_storage = "40"
		instance_name = "${var.name}ro"
		vswitch_id = "${alicloud_vswitch.default.id}"
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
