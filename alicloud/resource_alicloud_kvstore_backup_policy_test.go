package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/r-kvstore"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudKVStoreRedisBackupPolicy_classic(t *testing.T) {
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
				Config: testAccKVStoreBackupPolicy_classic(string(KVStoreRedis), redisInstanceClassForTest, string(KVStore4Dot0)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKVStoreBackupPolicyExists("alicloud_kvstore_backup_policy.policy", &policy),
					resource.TestCheckResourceAttrSet("alicloud_kvstore_backup_policy.policy", "instance_id"),
					resource.TestCheckResourceAttr("alicloud_kvstore_backup_policy.policy", "backup_time", "10:00Z-11:00Z"),
					resource.TestCheckResourceAttr("alicloud_kvstore_backup_policy.policy", "backup_period.#", "2"),
				),
			},
			resource.TestStep{
				Config: testAccKVStoreBackupPolicy_classicUpdatePeriod(string(KVStoreRedis), redisInstanceClassForTest, string(KVStore4Dot0)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKVStoreBackupPolicyExists("alicloud_kvstore_backup_policy.policy", &policy),
					resource.TestCheckResourceAttrSet("alicloud_kvstore_backup_policy.policy", "instance_id"),
					resource.TestCheckResourceAttr("alicloud_kvstore_backup_policy.policy", "backup_time", "10:00Z-11:00Z"),
					resource.TestCheckResourceAttr("alicloud_kvstore_backup_policy.policy", "backup_period.#", "3"),
				),
			},
			resource.TestStep{
				Config: testAccKVStoreBackupPolicy_classicUpdateTime(string(KVStoreRedis), redisInstanceClassForTest, string(KVStore4Dot0)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKVStoreBackupPolicyExists("alicloud_kvstore_backup_policy.policy", &policy),
					resource.TestCheckResourceAttrSet("alicloud_kvstore_backup_policy.policy", "instance_id"),
					resource.TestCheckResourceAttr("alicloud_kvstore_backup_policy.policy", "backup_time", "12:00Z-13:00Z"),
					resource.TestCheckResourceAttr("alicloud_kvstore_backup_policy.policy", "backup_period.#", "3"),
				),
			},
			resource.TestStep{
				Config: testAccKVStoreBackupPolicy_classicUpdateAll(string(KVStoreRedis), redisInstanceClassForTest, string(KVStore4Dot0)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKVStoreBackupPolicyExists("alicloud_kvstore_backup_policy.policy", &policy),
					resource.TestCheckResourceAttrSet("alicloud_kvstore_backup_policy.policy", "instance_id"),
					resource.TestCheckResourceAttr("alicloud_kvstore_backup_policy.policy", "backup_time", "13:00Z-14:00Z"),
					resource.TestCheckResourceAttr("alicloud_kvstore_backup_policy.policy", "backup_period.#", "1"),
				),
			},
		},
	})

}

func TestAccAlicloudKVStoreMemcacheBackupPolicy_classic(t *testing.T) {
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
				Config: testAccKVStoreBackupPolicy_classic(string(KVStoreMemcache), memcacheInstanceClassForTest, string(KVStore2Dot8)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKVStoreBackupPolicyExists("alicloud_kvstore_backup_policy.policy", &policy),
					resource.TestCheckResourceAttrSet("alicloud_kvstore_backup_policy.policy", "instance_id"),
					resource.TestCheckResourceAttr("alicloud_kvstore_backup_policy.policy", "backup_time", "10:00Z-11:00Z"),
					resource.TestCheckResourceAttr("alicloud_kvstore_backup_policy.policy", "backup_period.#", "2"),
				),
			},
			resource.TestStep{
				Config: testAccKVStoreBackupPolicy_classicUpdatePeriod(string(KVStoreMemcache), memcacheInstanceClassForTest, string(KVStore2Dot8)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKVStoreBackupPolicyExists("alicloud_kvstore_backup_policy.policy", &policy),
					resource.TestCheckResourceAttrSet("alicloud_kvstore_backup_policy.policy", "instance_id"),
					resource.TestCheckResourceAttr("alicloud_kvstore_backup_policy.policy", "backup_time", "10:00Z-11:00Z"),
					resource.TestCheckResourceAttr("alicloud_kvstore_backup_policy.policy", "backup_period.#", "3"),
				),
			},
			resource.TestStep{
				Config: testAccKVStoreBackupPolicy_classicUpdateTime(string(KVStoreMemcache), memcacheInstanceClassForTest, string(KVStore2Dot8)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKVStoreBackupPolicyExists("alicloud_kvstore_backup_policy.policy", &policy),
					resource.TestCheckResourceAttrSet("alicloud_kvstore_backup_policy.policy", "instance_id"),
					resource.TestCheckResourceAttr("alicloud_kvstore_backup_policy.policy", "backup_time", "12:00Z-13:00Z"),
					resource.TestCheckResourceAttr("alicloud_kvstore_backup_policy.policy", "backup_period.#", "3"),
				),
			},
			resource.TestStep{
				Config: testAccKVStoreBackupPolicy_classicUpdateAll(string(KVStoreMemcache), memcacheInstanceClassForTest, string(KVStore2Dot8)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKVStoreBackupPolicyExists("alicloud_kvstore_backup_policy.policy", &policy),
					resource.TestCheckResourceAttrSet("alicloud_kvstore_backup_policy.policy", "instance_id"),
					resource.TestCheckResourceAttr("alicloud_kvstore_backup_policy.policy", "backup_time", "13:00Z-14:00Z"),
					resource.TestCheckResourceAttr("alicloud_kvstore_backup_policy.policy", "backup_period.#", "1"),
				),
			},
		},
	})

}

func TestAccAlicloudKVStoreRedisBackupPolicy_vpc(t *testing.T) {
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
				Config: testAccKVStoreBackupPolicy_vpc(DatabaseCommonTestCase, string(KVStoreRedis), redisInstanceClassForTest, string(KVStore2Dot8)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKVStoreBackupPolicyExists("alicloud_kvstore_backup_policy.policy", &policy),
					resource.TestCheckResourceAttrSet("alicloud_kvstore_backup_policy.policy", "instance_id"),
					resource.TestCheckResourceAttr("alicloud_kvstore_backup_policy.policy", "backup_time", "10:00Z-11:00Z"),
					resource.TestCheckResourceAttr("alicloud_kvstore_backup_policy.policy", "backup_period.#", "2"),
				),
			},
			resource.TestStep{
				Config: testAccKVStoreBackupPolicy_vpcUpdatePeriod(DatabaseCommonTestCase, string(KVStoreRedis), redisInstanceClassForTest, string(KVStore2Dot8)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKVStoreBackupPolicyExists("alicloud_kvstore_backup_policy.policy", &policy),
					resource.TestCheckResourceAttrSet("alicloud_kvstore_backup_policy.policy", "instance_id"),
					resource.TestCheckResourceAttr("alicloud_kvstore_backup_policy.policy", "backup_time", "10:00Z-11:00Z"),
					resource.TestCheckResourceAttr("alicloud_kvstore_backup_policy.policy", "backup_period.#", "3"),
				),
			},
			resource.TestStep{
				Config: testAccKVStoreBackupPolicy_vpcUpdateTime(DatabaseCommonTestCase, string(KVStoreRedis), redisInstanceClassForTest, string(KVStore2Dot8)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKVStoreBackupPolicyExists("alicloud_kvstore_backup_policy.policy", &policy),
					resource.TestCheckResourceAttrSet("alicloud_kvstore_backup_policy.policy", "instance_id"),
					resource.TestCheckResourceAttr("alicloud_kvstore_backup_policy.policy", "backup_time", "11:00Z-12:00Z"),
					resource.TestCheckResourceAttr("alicloud_kvstore_backup_policy.policy", "backup_period.#", "3"),
				),
			},
			resource.TestStep{
				Config: testAccKVStoreBackupPolicy_vpcUpdateAll(DatabaseCommonTestCase, string(KVStoreRedis), redisInstanceClassForTest, string(KVStore2Dot8)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKVStoreBackupPolicyExists("alicloud_kvstore_backup_policy.policy", &policy),
					resource.TestCheckResourceAttrSet("alicloud_kvstore_backup_policy.policy", "instance_id"),
					resource.TestCheckResourceAttr("alicloud_kvstore_backup_policy.policy", "backup_time", "12:00Z-13:00Z"),
					resource.TestCheckResourceAttr("alicloud_kvstore_backup_policy.policy", "backup_period.#", "1"),
				),
			},
		},
	})

}

func TestAccAlicloudKVStoreMemcacheBackupPolicy_vpc(t *testing.T) {
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
				Config: testAccKVStoreBackupPolicy_vpc(DatabaseCommonTestCase, string(KVStoreMemcache), memcacheInstanceClassForTest, string(KVStore2Dot8)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKVStoreBackupPolicyExists("alicloud_kvstore_backup_policy.policy", &policy),
					resource.TestCheckResourceAttrSet("alicloud_kvstore_backup_policy.policy", "instance_id"),
					resource.TestCheckResourceAttr("alicloud_kvstore_backup_policy.policy", "backup_time", "10:00Z-11:00Z"),
					resource.TestCheckResourceAttr("alicloud_kvstore_backup_policy.policy", "backup_period.#", "2"),
				),
			},
			resource.TestStep{
				Config: testAccKVStoreBackupPolicy_vpcUpdatePeriod(DatabaseCommonTestCase, string(KVStoreMemcache), memcacheInstanceClassForTest, string(KVStore2Dot8)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKVStoreBackupPolicyExists("alicloud_kvstore_backup_policy.policy", &policy),
					resource.TestCheckResourceAttrSet("alicloud_kvstore_backup_policy.policy", "instance_id"),
					resource.TestCheckResourceAttr("alicloud_kvstore_backup_policy.policy", "backup_time", "10:00Z-11:00Z"),
					resource.TestCheckResourceAttr("alicloud_kvstore_backup_policy.policy", "backup_period.#", "3"),
				),
			},
			resource.TestStep{
				Config: testAccKVStoreBackupPolicy_vpcUpdateTime(DatabaseCommonTestCase, string(KVStoreMemcache), memcacheInstanceClassForTest, string(KVStore2Dot8)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKVStoreBackupPolicyExists("alicloud_kvstore_backup_policy.policy", &policy),
					resource.TestCheckResourceAttrSet("alicloud_kvstore_backup_policy.policy", "instance_id"),
					resource.TestCheckResourceAttr("alicloud_kvstore_backup_policy.policy", "backup_time", "11:00Z-12:00Z"),
					resource.TestCheckResourceAttr("alicloud_kvstore_backup_policy.policy", "backup_period.#", "3"),
				),
			},
			resource.TestStep{
				Config: testAccKVStoreBackupPolicy_vpcUpdateAll(DatabaseCommonTestCase, string(KVStoreMemcache), memcacheInstanceClassForTest, string(KVStore2Dot8)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKVStoreBackupPolicyExists("alicloud_kvstore_backup_policy.policy", &policy),
					resource.TestCheckResourceAttrSet("alicloud_kvstore_backup_policy.policy", "instance_id"),
					resource.TestCheckResourceAttr("alicloud_kvstore_backup_policy.policy", "backup_time", "12:00Z-13:00Z"),
					resource.TestCheckResourceAttr("alicloud_kvstore_backup_policy.policy", "backup_period.#", "1"),
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

		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		kvstoreService := KvstoreService{client}
		policy, err := kvstoreService.DescribeRKVInstancebackupPolicy(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error Describe KVStore Instance backup policy: %#v", err)
		}

		*d = *policy
		return nil
	}
}

func testAccCheckKVStoreBackupPolicyDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	kvstoreService := KvstoreService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_kvstore_instance" {
			continue
		}

		if _, err := kvstoreService.DescribeRKVInstancebackupPolicy(rs.Primary.ID); err != nil {
			if NotFoundError(err) {
				continue
			}
			return fmt.Errorf("Error Describe DB backup policy: %#v", err)
		}
		return fmt.Errorf("KVStore Instance %s Policy sitll exists.", rs.Primary.ID)
	}

	return nil
}

func testAccKVStoreBackupPolicy_classic(instanceType, instanceClass, engineVersion string) string {
	return fmt.Sprintf(`
	data "alicloud_zones" "default" {
		available_resource_creation = "KVStore"
	}
	variable "name" {
		default = "tf-testAccKVStoreBackupPolicy_classic"
	}

	resource "alicloud_kvstore_instance" "foo" {
		availability_zone = "${data.alicloud_zones.default.zones.0.id}"
		instance_name  = "${var.name}"
		security_ips = ["10.0.0.1"]
		instance_type = "%s"
		instance_class = "%s"
		engine_version = "%s"
	}
	resource "alicloud_kvstore_backup_policy" "policy" {
		instance_id = "${alicloud_kvstore_instance.foo.id}"
		backup_period = ["Tuesday", "Wednesday"]
		backup_time = "10:00Z-11:00Z"
	}
	`, instanceType, instanceClass, engineVersion)
}

func testAccKVStoreBackupPolicy_classicUpdatePeriod(instanceType, instanceClass, engineVersion string) string {
	return fmt.Sprintf(`
	data "alicloud_zones" "default" {
		available_resource_creation = "KVStore"
	}
	variable "name" {
		default = "tf-testAccKVStoreBackupPolicy_classic"
	}

	resource "alicloud_kvstore_instance" "foo" {
		availability_zone = "${data.alicloud_zones.default.zones.0.id}"
		instance_name  = "${var.name}"
		security_ips = ["10.0.0.1"]
		instance_type = "%s"
		instance_class = "%s"
		engine_version = "%s"
	}
	resource "alicloud_kvstore_backup_policy" "policy" {
		instance_id = "${alicloud_kvstore_instance.foo.id}"
		backup_period = ["Tuesday", "Wednesday", "Sunday"]
		backup_time = "10:00Z-11:00Z"
	}
	`, instanceType, instanceClass, engineVersion)
}

func testAccKVStoreBackupPolicy_classicUpdateTime(instanceType, instanceClass, engineVersion string) string {
	return fmt.Sprintf(`
	data "alicloud_zones" "default" {
		available_resource_creation = "KVStore"
	}
	variable "name" {
		default = "tf-testAccKVStoreBackupPolicy_classic"
	}

	resource "alicloud_kvstore_instance" "foo" {
		availability_zone = "${data.alicloud_zones.default.zones.0.id}"
		instance_name  = "${var.name}"
		security_ips = ["10.0.0.1"]
		instance_type = "%s"
		instance_class = "%s"
		engine_version = "%s"
	}
	resource "alicloud_kvstore_backup_policy" "policy" {
		instance_id = "${alicloud_kvstore_instance.foo.id}"
		backup_period = ["Tuesday", "Wednesday", "Sunday"]
		backup_time = "12:00Z-13:00Z"
	}
	`, instanceType, instanceClass, engineVersion)
}

func testAccKVStoreBackupPolicy_classicUpdateAll(instanceType, instanceClass, engineVersion string) string {
	return fmt.Sprintf(`
	data "alicloud_zones" "default" {
		available_resource_creation = "KVStore"
	}
	variable "name" {
		default = "tf-testAccKVStoreBackupPolicy_classic"
	}

	resource "alicloud_kvstore_instance" "foo" {
		availability_zone = "${data.alicloud_zones.default.zones.0.id}"
		instance_name  = "${var.name}"
		security_ips = ["10.0.0.1"]
		instance_type = "%s"
		instance_class = "%s"
		engine_version = "%s"
	}
	resource "alicloud_kvstore_backup_policy" "policy" {
		instance_id = "${alicloud_kvstore_instance.foo.id}"
		backup_period = ["Sunday"]
		backup_time = "13:00Z-14:00Z"
	}
	`, instanceType, instanceClass, engineVersion)
}

func testAccKVStoreBackupPolicy_vpc(common, instanceType, instanceClass, engineVersion string) string {
	return fmt.Sprintf(`
	%s
	variable "creation" {
		default = "KVStore"
	}
	variable "multi_az" {
		default = "false"
	}
	variable "name" {
		default = "tf-testAccKVStoreBackupPolicy_vpc"
	}
	resource "alicloud_kvstore_instance" "foo" {
		instance_class = "%s"
		instance_name  = "${var.name}"
		vswitch_id     = "${alicloud_vswitch.default.id}"
		private_ip     = "172.16.0.10"
		security_ips = ["10.0.0.1"]
		instance_type = "%s"
		engine_version = "%s"
	}
	resource "alicloud_kvstore_backup_policy" "policy" {
		instance_id = "${alicloud_kvstore_instance.foo.id}"
		backup_period = ["Tuesday", "Wednesday"]
		backup_time = "10:00Z-11:00Z"
	}
	`, common, instanceClass, instanceType, engineVersion)
}

func testAccKVStoreBackupPolicy_vpcUpdatePeriod(common, instanceType, instanceClass, engineVersion string) string {
	return fmt.Sprintf(`
	%s
	variable "creation" {
		default = "KVStore"
	}
	variable "multi_az" {
		default = "false"
	}
	variable "name" {
		default = "tf-testAccKVStoreBackupPolicy_vpc"
	}
	resource "alicloud_kvstore_instance" "foo" {
		instance_class = "%s"
		instance_name  = "${var.name}"
		vswitch_id     = "${alicloud_vswitch.default.id}"
		private_ip     = "172.16.0.10"
		security_ips = ["10.0.0.1"]
		instance_type = "%s"
		engine_version = "%s"
	}
	resource "alicloud_kvstore_backup_policy" "policy" {
		instance_id = "${alicloud_kvstore_instance.foo.id}"
		backup_period = ["Tuesday", "Wednesday", "Sunday"]
		backup_time = "10:00Z-11:00Z"
	}
	`, common, instanceClass, instanceType, engineVersion)
}
func testAccKVStoreBackupPolicy_vpcUpdateTime(common, instanceType, instanceClass, engineVersion string) string {
	return fmt.Sprintf(`
	%s
	variable "creation" {
		default = "KVStore"
	}
	variable "multi_az" {
		default = "false"
	}
	variable "name" {
		default = "tf-testAccKVStoreBackupPolicy_vpc"
	}
	resource "alicloud_kvstore_instance" "foo" {
		instance_class = "%s"
		instance_name  = "${var.name}"
		vswitch_id     = "${alicloud_vswitch.default.id}"
		private_ip     = "172.16.0.10"
		security_ips = ["10.0.0.1"]
		instance_type = "%s"
		engine_version = "%s"
	}
	resource "alicloud_kvstore_backup_policy" "policy" {
		instance_id = "${alicloud_kvstore_instance.foo.id}"
		backup_period = ["Tuesday", "Wednesday", "Sunday"]
		backup_time = "11:00Z-12:00Z"
	}
	`, common, instanceClass, instanceType, engineVersion)
}
func testAccKVStoreBackupPolicy_vpcUpdateAll(common, instanceType, instanceClass, engineVersion string) string {
	return fmt.Sprintf(`
	%s
	variable "creation" {
		default = "KVStore"
	}
	variable "multi_az" {
		default = "false"
	}
	variable "name" {
		default = "tf-testAccKVStoreBackupPolicy_vpc"
	}
	resource "alicloud_kvstore_instance" "foo" {
		instance_class = "%s"
		instance_name  = "${var.name}"
		vswitch_id     = "${alicloud_vswitch.default.id}"
		private_ip     = "172.16.0.10"
		security_ips = ["10.0.0.1"]
		instance_type = "%s"
		engine_version = "%s"
	}
	resource "alicloud_kvstore_backup_policy" "policy" {
		instance_id = "${alicloud_kvstore_instance.foo.id}"
		backup_period = ["Tuesday"]
		backup_time = "12:00Z-13:00Z"
	}
	`, common, instanceClass, instanceType, engineVersion)
}
