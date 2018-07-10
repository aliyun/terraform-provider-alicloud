package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAlicloudDBInstance_basic(t *testing.T) {
	var instance rds.DBInstanceAttribute

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_db_instance.foo",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDBInstanceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDBInstanceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDBInstanceExists(
						"alicloud_db_instance.foo", &instance),
					resource.TestCheckResourceAttr(
						"alicloud_db_instance.foo",
						"instance_storage",
						"10"),
					resource.TestCheckResourceAttr(
						"alicloud_db_instance.foo",
						"engine_version",
						"5.6"),
					resource.TestCheckResourceAttr(
						"alicloud_db_instance.foo",
						"engine",
						"MySQL"),
					resource.TestCheckResourceAttr(
						"alicloud_db_instance.foo",
						"instance_name",
						"testAccDBInstanceConfig"),
				),
			},
		},
	})

}

func TestAccAlicloudDBInstance_vpc(t *testing.T) {
	var instance rds.DBInstanceAttribute

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_db_instance.foo",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDBInstanceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDBInstance_vpc,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDBInstanceExists(
						"alicloud_db_instance.foo", &instance),
					resource.TestCheckResourceAttr(
						"alicloud_db_instance.foo",
						"instance_storage",
						"10"),
					resource.TestCheckResourceAttr(
						"alicloud_db_instance.foo",
						"engine_version",
						"5.6"),
					resource.TestCheckResourceAttr(
						"alicloud_db_instance.foo",
						"engine",
						"MySQL"),
				),
			},
		},
	})

}

func TestAccAlicloudDBInstance_multiAZ(t *testing.T) {
	var instance rds.DBInstanceAttribute

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_db_instance.foo",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDBInstanceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDBInstance_multiAZ,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDBInstanceExists(
						"alicloud_db_instance.foo", &instance),
					testAccCheckDBInstanceMultiIZ(&instance),
				),
			},
		},
	})

}

func TestAccAlicloudDBInstance_securityIps(t *testing.T) {
	var ips []map[string]interface{}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_db_instance.foo",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDBInstanceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDBInstance_securityIps,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityIpExists(
						"alicloud_db_instance.foo", ips),
					testAccCheckKeyValueInMaps(ips, "security ip", "security_ips", "127.0.0.1"),
				),
			},

			resource.TestStep{
				Config: testAccDBInstance_securityIpsUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityIpExists(
						"alicloud_db_instance.foo", ips),
					testAccCheckKeyValueInMaps(ips, "security ip", "security_ips", "10.168.1.12,100.69.7.112"),
				),
			},
		},
	})

}

func TestAccAlicloudDBInstance_upgradeClass(t *testing.T) {
	var instance rds.DBInstanceAttribute

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_db_instance.foo",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDBInstanceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDBInstance_class,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDBInstanceExists(
						"alicloud_db_instance.foo", &instance),
					resource.TestCheckResourceAttr("alicloud_db_instance.foo", "instance_type", "rds.mysql.t1.small"),
				),
			},

			resource.TestStep{
				Config: testAccDBInstance_classUpgrade,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDBInstanceExists(
						"alicloud_db_instance.foo", &instance),
					resource.TestCheckResourceAttr("alicloud_db_instance.foo", "instance_type", "rds.mysql.s1.small"),
				),
			},
		},
	})

}

func testAccCheckSecurityIpExists(n string, ips []map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No DB Instance ID is set")
		}

		resp, err := testAccProvider.Meta().(*AliyunClient).DescribeDBSecurityIps(rs.Primary.ID)
		log.Printf("[DEBUG] check instance %s security ip %#v", rs.Primary.ID, resp)

		if err != nil {
			return err
		}

		if len(resp) < 1 {
			return fmt.Errorf("DB security ip not found")
		}

		ips = flattenDBSecurityIPs(resp)
		return nil
	}
}

func testAccCheckDBInstanceMultiIZ(i *rds.DBInstanceAttribute) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if !strings.Contains(i.ZoneId, MULTI_IZ_SYMBOL) {
			return fmt.Errorf("Current region does not support multiIZ.")
		}
		return nil
	}
}

func testAccCheckDBInstanceExists(n string, d *rds.DBInstanceAttribute) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No DB Instance ID is set")
		}

		client := testAccProvider.Meta().(*AliyunClient)
		attr, err := client.DescribeDBInstanceById(rs.Primary.ID)
		log.Printf("[DEBUG] check instance %s attribute %#v", rs.Primary.ID, attr)

		if err != nil {
			return err
		}

		if attr == nil {
			return fmt.Errorf("DB Instance not found")
		}

		*d = *attr
		return nil
	}
}

func testAccCheckKeyValueInMaps(ps []map[string]interface{}, propName, key, value string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, policy := range ps {
			if policy[key].(string) != value {
				return fmt.Errorf("DB %s attribute '%s' expected %#v, got %#v", propName, key, value, policy[key])
			}
		}
		return nil
	}
}

func testAccCheckDBInstanceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*AliyunClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_db_instance" {
			continue
		}

		ins, err := client.DescribeDBInstanceById(rs.Primary.ID)

		if ins != nil {
			return fmt.Errorf("Error DB Instance still exist")
		}

		// Verify the error is what we want
		if err != nil {
			if NotFoundError(err) || IsExceptedError(err, InvalidDBInstanceIdNotFound) || IsExceptedError(err, InvalidDBInstanceNameNotFound) {
				continue
			}
			return err
		}
	}

	return nil
}

const testAccDBInstanceConfig = `
resource "alicloud_db_instance" "foo" {
	engine = "MySQL"
	engine_version = "5.6"
	instance_type = "rds.mysql.t1.small"
	instance_storage = "10"
	instance_charge_type = "Postpaid"
	instance_name = "testAccDBInstanceConfig"
}
`

const testAccDBInstance_vpc = `
data "alicloud_zones" "default" {
	available_resource_creation = "Rds"
}
variable "name" {
	default = "testAccDBInstance_vpc"
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

resource "alicloud_db_instance" "foo" {
	engine = "MySQL"
	engine_version = "5.6"
	instance_type = "rds.mysql.t1.small"
	instance_storage = "10"
	instance_charge_type = "Postpaid"
	instance_name = "${var.name}"
	vswitch_id = "${alicloud_vswitch.foo.id}"
	security_ips = ["10.168.1.12", "100.69.7.112"]
}
`
const testAccDBInstance_multiAZ = `
provider "alicloud" {
  region = "cn-shanghai"
}

data "alicloud_zones" "default" {
  available_resource_creation= "Rds"
  multi = true
  output_file = "zone.json"
}
variable "name" {
	default = "testAccDBInstance_multiAZ"
}
resource "alicloud_db_instance" "foo" {
	engine = "MySQL"
	engine_version = "5.6"
	instance_type = "rds.mysql.t1.small"
	instance_storage = "10"
	zone_id = "${data.alicloud_zones.default.zones.0.id}"
	instance_name = "${var.name}"
}
`

const testAccDBInstance_securityIps = `
variable "name" {
	default = "testAccDBInstance_securityIps"
}
resource "alicloud_db_instance" "foo" {
	engine = "MySQL"
	engine_version = "5.6"
	instance_type = "rds.mysql.t1.small"
	instance_storage = "10"
	instance_charge_type = "Postpaid"
	instance_name = "${var.name}"
}
`
const testAccDBInstance_securityIpsUpdate = `
variable "name" {
	default = "testAccDBInstance_securityIpsUpdate"
}
resource "alicloud_db_instance" "foo" {
	engine = "MySQL"
	engine_version = "5.6"
	instance_type = "rds.mysql.t1.small"
	instance_storage = "10"
	instance_charge_type = "Postpaid"
	instance_name = "${var.name}"
	security_ips = ["10.168.1.12", "100.69.7.112"]
}
`

const testAccDBInstance_class = `
variable "name" {
	default = "testAccDBInstance_class"
}
resource "alicloud_db_instance" "foo" {
	engine = "MySQL"
	engine_version = "5.6"
	instance_type = "rds.mysql.t1.small"
	instance_storage = "10"
	instance_name = "${var.name}"
}
`
const testAccDBInstance_classUpgrade = `
variable "name" {
	default = "testAccDBInstance_class"
}
resource "alicloud_db_instance" "foo" {
	engine = "MySQL"
	engine_version = "5.6"
	instance_type = "rds.mysql.s1.small"
	instance_storage = "10"
	instance_name = "${var.name}"
}
`
