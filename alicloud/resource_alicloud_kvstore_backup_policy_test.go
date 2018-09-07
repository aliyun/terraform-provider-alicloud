package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/r-kvstore"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAlicloudKVStoreBackupPolicy_basic(t *testing.T) {
	var policy r_kvstore.DescribeBackupPolicyResponse

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_kvstore_backup_policy.policy",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKVStoreBackupPolicyDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccKVStoreBackupPolicy_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKVStoreBackupPolicyExists(
						"alicloud_kvstore_backup_policy.policy", &policy),
					resource.TestCheckResourceAttr("alicloud_kvstore_backup_policy.policy", "backup_time", "10:00Z-11:00Z"),
				),
			},
		},
	})

}

func testAccCheckKVStoreBackupPolicyExists(n string, d *r_kvstore.DescribeBackupPolicyResponse) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No KVStore Instance backup policy ID is set")
		}

		client := testAccProvider.Meta().(*AliyunClient)
		policy, err := client.DescribeRKVInstancebackupPolicy(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error Describe KVStore Instance backup policy: %#v", err)
		}

		*d = *policy
		return nil
	}
}

func testAccCheckKVStoreBackupPolicyDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*AliyunClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_kvstore_instance" {
			continue
		}

		if _, err := client.DescribeRKVInstancebackupPolicy(rs.Primary.ID); err != nil {
			if NotFoundError(err) {
				continue
			}
			return fmt.Errorf("Error Describe DB backup policy: %#v", err)
		}
		return fmt.Errorf("KVStore Instance %s Policy sitll exists.", rs.Primary.ID)
	}

	return nil
}

const testAccKVStoreBackupPolicy_basic = `
data "alicloud_zones" "default" {
	available_resource_creation = "KVStore"
}
variable "name" {
	default = "tf-testacckvstorebackuppolicy_basic"
}
resource "alicloud_vpc" "foo" {
	name = "${var.name}"
	cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "foo" {
 	vpc_id = "${alicloud_vpc.foo.id}"
 	cidr_block = "172.16.0.0/21"
 	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
 	name = "${var.name}"
}

resource "alicloud_kvstore_instance" "foo" {
	instance_class = "redis.master.small.default"
	instance_name  = "${var.name}"
	password       = "Test12345"
	vswitch_id     = "${alicloud_vswitch.foo.id}"
}

resource "alicloud_kvstore_backup_policy" "policy" {
	instance_id = "${alicloud_kvstore_instance.foo.id}"
	backup_period = ["Tuesday", "Wednesday"]
	backup_time = "10:00Z-11:00Z"
}
`
