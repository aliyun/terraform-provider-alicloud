package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
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
			{
				Config: testAccDBBackupPolicy_basic(DatabaseCommonTestCase),
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
		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		rdsService := RdsService{client}
		resp, err := rdsService.DescribeBackupPolicy(rs.Primary.ID)
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
	client := testAccProvider.Meta().(*connectivity.AliyunClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_db_account" {
			continue
		}

		_, err := client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
			return rdsClient.DescribeBackupPolicy(&rds.DescribeBackupPolicyRequest{
				DBInstanceId: rs.Primary.ID,
			})
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

func testAccDBBackupPolicy_basic(common string) string {
	return fmt.Sprintf(`
	%s
	variable "creation" {
		default = "Rds"
	}
	variable "multi_az" {
		default = "false"
	}
	variable "name" {
		default = "tf-testAccDBbackuppolicy_basic"
	}

	resource "alicloud_db_instance" "instance" {
		engine = "MySQL"
		engine_version = "5.6"
		instance_type = "rds.mysql.s1.small"
		instance_storage = "10"
		vswitch_id = "${alicloud_vswitch.default.id}"
		instance_name = "${var.name}"
	}

	resource "alicloud_db_backup_policy" "policy" {
		  instance_id = "${alicloud_db_instance.instance.id}"
		  backup_period = ["Tuesday", "Wednesday"]
		  backup_time = "10:00Z-11:00Z"
		  retention_period = "10"
	}
	`, common)
}
