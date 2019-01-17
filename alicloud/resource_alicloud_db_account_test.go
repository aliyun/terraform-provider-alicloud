package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudDBAccount_mysql(t *testing.T) {
	var account rds.DBInstanceAccount

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_db_account.account",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDBAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDBAccount_mysql(RdsCommonTestCase),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDBAccountExists("alicloud_db_account.account", &account),
					resource.TestCheckResourceAttrSet("alicloud_db_account.account", "instance_id"),
					resource.TestCheckResourceAttr("alicloud_db_account.account", "name", "tftestnormal"),
					resource.TestCheckResourceAttr("alicloud_db_account.account", "description", "from terraform"),
					resource.TestCheckResourceAttr("alicloud_db_account.account", "type", string(DBAccountNormal)),
					resource.TestCheckResourceAttrSet("alicloud_db_account.super", "instance_id"),
					resource.TestCheckResourceAttr("alicloud_db_account.super", "name", "tftestsuper"),
					resource.TestCheckResourceAttr("alicloud_db_account.super", "description", "from terraform"),
					resource.TestCheckResourceAttr("alicloud_db_account.super", "type", string(DBAccountSuper)),
				),
			},
		},
	})

}

func testAccCheckDBAccountExists(n string, d *rds.DBInstanceAccount) resource.TestCheckFunc {
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
		parts := strings.Split(rs.Primary.ID, COLON_SEPARATED)
		account, err := rdsService.DescribeDatabaseAccount(parts[0], parts[1])

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

func testAccCheckDBAccountDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	rdsService := RdsService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_db_account" {
			continue
		}

		parts := strings.Split(rs.Primary.ID, COLON_SEPARATED)

		account, err := rdsService.DescribeDatabaseAccount(parts[0], parts[1])

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

func testAccDBAccount_mysql(common string) string {
	return fmt.Sprintf(`
	%s
	variable "creation" {
		default = "Rds"
	}
	variable "name" {
		default = "tf-testAccDBaccount_mysql"
	}

	resource "alicloud_db_instance" "instance" {
		engine = "MySQL"
		engine_version = "5.6"
		instance_type = "rds.mysql.s1.small"
		instance_storage = "10"
		vswitch_id = "${alicloud_vswitch.default.id}"
	        instance_name = "${var.name}"
	}

	resource "alicloud_db_account" "account" {
	  	instance_id = "${alicloud_db_instance.instance.id}"
	  	name = "tftestnormal"
	  	password = "Test12345"
	  	description = "from terraform"
	}
	resource "alicloud_db_account" "super" {
	  	instance_id = "${alicloud_db_instance.instance.id}"
	  	name = "tftestsuper"
	  	password = "Test12345"
	  	description = "from terraform"
	  	type = "Super"
	}
	`, common)
}
