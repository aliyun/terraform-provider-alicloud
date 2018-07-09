package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAlicloudDBAccountPrivilege_basic(t *testing.T) {

	var account rds.DBInstanceAccount

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_db_account_privilege.privilege",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDBAccountPrivilegeDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDBAccountPrivilege_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDBAccountPrivilegeExists(
						"alicloud_db_account_privilege.privilege", &account),
					resource.TestCheckResourceAttr("alicloud_db_account_privilege.privilege", "account_name", "tftestprivilege"),
					resource.TestCheckResourceAttr("alicloud_db_account_privilege.privilege", "db_names.#", "2"),
				),
			},
		},
	})

}

func testAccCheckDBAccountPrivilegeExists(n string, d *rds.DBInstanceAccount) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No DB account ID is set")
		}

		client := testAccProvider.Meta().(*AliyunClient)
		parts := strings.Split(rs.Primary.ID, COLON_SEPARATED)
		account, err := client.DescribeDatabaseAccount(parts[0], parts[1])

		if err != nil {
			return err
		}

		if account == nil {
			return fmt.Errorf("account is not found in the instance %s.", parts[0])
		}

		*d = *account
		return nil
	}
}

func testAccCheckDBAccountPrivilegeDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*AliyunClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_db_account_privilege" {
			continue
		}

		parts := strings.Split(rs.Primary.ID, COLON_SEPARATED)

		account, err := client.DescribeDatabaseAccount(parts[0], parts[1])

		// Verify the error is what we want
		if err != nil {
			if NotFoundError(err) || IsExceptedError(err, InvalidDBInstanceIdNotFound) || IsExceptedError(err, InvalidAccountNameNotFound) {
				continue
			}
			return err
		}

		if account != nil {
			return fmt.Errorf("Error db account %s is still existing.", parts[1])
		}
	}

	return nil
}

const testAccDBAccountPrivilege_basic = `
variable "name" {
	default = "testaccdbaccountprivilege_basic"
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

resource "alicloud_db_database" "db" {
  count = 2
  instance_id = "${alicloud_db_instance.instance.id}"
  name = "${var.name}_${count.index}"
  description = "from terraform"
}

resource "alicloud_db_account" "account" {
  instance_id = "${alicloud_db_instance.instance.id}"
  name = "tftestprivilege"
  password = "Test12345"
  description = "from terraform"
}

resource "alicloud_db_account_privilege" "privilege" {
  instance_id = "${alicloud_db_instance.instance.id}"
  account_name = "${alicloud_db_account.account.name}"
  privilege = "ReadOnly"
  db_names = ["${alicloud_db_database.db.*.name}"]
}
`
