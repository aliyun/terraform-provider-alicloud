package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAlicloudDBBackupPolicy_basic(t *testing.T) {
	var policy rds.DescribeBackupPolicyResponse

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_db_backup_policy.policy",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDBBackupPolicyDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDBBackupPolicy_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDBBackupPolicyExists(
						"alicloud_db_backup_policy.policy", &policy),
					resource.TestCheckResourceAttr("alicloud_db_backup_policy.policy", "backup_time", "10:00Z-11:00Z"),
					resource.TestCheckResourceAttr("alicloud_db_backup_policy.policy", "retention_period", "10"),
				),
			},
		},
	})

}

func testAccCheckDBBackupPolicyExists(n string, d *rds.DescribeBackupPolicyResponse) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No DB account ID is set")
		}

		resp, err := testAccProvider.Meta().(*AliyunClient).DescribeBackupPolicy(rs.Primary.ID)
		if err != nil {

			return fmt.Errorf("Error Describe DB backup policy: %#v", err)
		}

		if resp == nil {
			return fmt.Errorf("Backup policy is not found in the instance %s.", rs.Primary.ID)
		}

		*d = *resp
		return nil
	}
}

func testAccCheckDBBackupPolicyDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*AliyunClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_db_account" {
			continue
		}

		_, err := client.rdsconn.DescribeBackupPolicy(&rds.DescribeBackupPolicyRequest{
			DBInstanceId: rs.Primary.ID,
		})
		if err != nil {
			if IsExceptedError(err, InvalidDBInstanceIdNotFound) || IsExceptedError(err, InvalidDBInstanceNameNotFound) {
				continue
			}
			return fmt.Errorf("Error Describe DB backup policy: %#v", err)
		}
	}

	return nil
}

const testAccDBBackupPolicy_basic = `
variable "name" {
	default = "testaccdbbackuppolicy_basic"
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

resource "alicloud_db_backup_policy" "policy" {
  	instance_id = "${alicloud_db_instance.instance.id}"
  	backup_period = ["Tuesday", "Wednesday"]
  	backup_time = "10:00Z-11:00Z"
  	retention_period = "10"
}
`
