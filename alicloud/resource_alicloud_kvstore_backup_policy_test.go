package alicloud

import (
	"fmt"
	"testing"

	r_kvstore "github.com/aliyun/alibaba-cloud-sdk-go/services/r-kvstore"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var redisInstanceClassForTest = "redis.master.small.default"
var memcacheInstanceClassForTest = "memcache.master.small.default"

func TestAccAlicloudKVStoreRedisBackupPolicy_classic(t *testing.T) {
	var policy *r_kvstore.DescribeBackupPolicyResponse

	resourceId := "alicloud_kvstore_backup_policy.default"
	ra := resourceAttrInit(resourceId, kvStoreMap)
	serviceFunc := func() interface{} {
		return &KvstoreService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &policy, serviceFunc, "DescribeKVstoreBackupPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.KVstoreClassicNetworkInstanceSupportRegions)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKVStoreBackupPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKVStoreBackupPolicy_classic(string(KVStoreRedis), redisInstanceClassForTest, string(KVStore4Dot0)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccKVStoreBackupPolicy_classicUpdatePeriod(string(KVStoreRedis), redisInstanceClassForTest, string(KVStore4Dot0)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_period.#": "3",
					}),
				),
			},
			{
				Config: testAccKVStoreBackupPolicy_classicUpdateTime(string(KVStoreRedis), redisInstanceClassForTest, string(KVStore4Dot0)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_time": "12:00Z-13:00Z",
					}),
				),
			},
			{
				Config: testAccKVStoreBackupPolicy_classicUpdateAll(string(KVStoreRedis), redisInstanceClassForTest, string(KVStore4Dot0)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_time":     "13:00Z-14:00Z",
						"backup_period.#": "1",
					}),
				),
			},
		},
	})

}

func TestAccAlicloudKVStoreMemcacheBackupPolicy_classic(t *testing.T) {
	var policy *r_kvstore.DescribeBackupPolicyResponse

	resourceId := "alicloud_kvstore_backup_policy.default"
	ra := resourceAttrInit(resourceId, kvStoreMap)
	serviceFunc := func() interface{} {
		return &KvstoreService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &policy, serviceFunc, "DescribeKVstoreBackupPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.KVstoreClassicNetworkInstanceSupportRegions)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKVStoreBackupPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKVStoreBackupPolicy_classic(string(KVStoreMemcache), memcacheInstanceClassForTest, string(KVStore4Dot0)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccKVStoreBackupPolicy_classicUpdatePeriod(string(KVStoreMemcache), memcacheInstanceClassForTest, string(KVStore4Dot0)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_period.#": "3",
					}),
				),
			},
			{
				Config: testAccKVStoreBackupPolicy_classicUpdateTime(string(KVStoreMemcache), memcacheInstanceClassForTest, string(KVStore4Dot0)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_time": "12:00Z-13:00Z",
					}),
				),
			},
			{
				Config: testAccKVStoreBackupPolicy_classicUpdateAll(string(KVStoreMemcache), memcacheInstanceClassForTest, string(KVStore4Dot0)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_time":     "13:00Z-14:00Z",
						"backup_period.#": "1",
					}),
				),
			},
		},
	})

}

func TestAccAlicloudKVStoreRedisBackupPolicy_vpc(t *testing.T) {
	var policy *r_kvstore.DescribeBackupPolicyResponse

	resourceId := "alicloud_kvstore_backup_policy.default"
	ra := resourceAttrInit(resourceId, kvStoreMap)
	serviceFunc := func() interface{} {
		return &KvstoreService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &policy, serviceFunc, "DescribeKVstoreBackupPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.KVstoreClassicNetworkInstanceSupportRegions)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKVStoreBackupPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKVStoreBackupPolicy_vpc(KVStoreCommonTestCase, string(KVStoreRedis), redisInstanceClassForTest, string(KVStore2Dot8)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccKVStoreBackupPolicy_vpcUpdatePeriod(KVStoreCommonTestCase, string(KVStoreRedis), redisInstanceClassForTest, string(KVStore2Dot8)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_period.#": "3",
					}),
				),
			},
			{
				Config: testAccKVStoreBackupPolicy_vpcUpdateTime(KVStoreCommonTestCase, string(KVStoreRedis), redisInstanceClassForTest, string(KVStore2Dot8)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_time": "11:00Z-12:00Z",
					}),
				),
			},
			{
				Config: testAccKVStoreBackupPolicy_vpcUpdateAll(KVStoreCommonTestCase, string(KVStoreRedis), redisInstanceClassForTest, string(KVStore2Dot8)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_time":     "12:00Z-13:00Z",
						"backup_period.#": "1",
					}),
				),
			},
		},
	})

}

func TestAccAlicloudKVStoreMemcacheBackupPolicy_vpc(t *testing.T) {
	var policy *r_kvstore.DescribeBackupPolicyResponse

	resourceId := "alicloud_kvstore_backup_policy.default"
	ra := resourceAttrInit(resourceId, kvStoreMap)
	serviceFunc := func() interface{} {
		return &KvstoreService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &policy, serviceFunc, "DescribeKVstoreBackupPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.KVstoreClassicNetworkInstanceSupportRegions)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKVStoreBackupPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKVStoreBackupPolicy_vpc(KVStoreCommonTestCase, string(KVStoreMemcache), memcacheInstanceClassForTest, string(KVStore4Dot0)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccKVStoreBackupPolicy_vpcUpdatePeriod(KVStoreCommonTestCase, string(KVStoreMemcache), memcacheInstanceClassForTest, string(KVStore4Dot0)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_period.#": "3",
					}),
				),
			},
			{
				Config: testAccKVStoreBackupPolicy_vpcUpdateTime(KVStoreCommonTestCase, string(KVStoreMemcache), memcacheInstanceClassForTest, string(KVStore4Dot0)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_time": "11:00Z-12:00Z",
					}),
				),
			},
			{
				Config: testAccKVStoreBackupPolicy_vpcUpdateAll(KVStoreCommonTestCase, string(KVStoreMemcache), memcacheInstanceClassForTest, string(KVStore4Dot0)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_time":     "12:00Z-13:00Z",
						"backup_period.#": "1",
					}),
				),
			},
		},
	})

}

func testAccCheckKVStoreBackupPolicyDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	kvstoreService := KvstoreService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_kvstore_instance" {
			continue
		}

		if _, err := kvstoreService.DescribeKVstoreBackupPolicy(rs.Primary.ID); err != nil {
			if NotFoundError(err) {
				continue
			}
			return fmt.Errorf("Error Describe DB backup policy: %#v", err)
		}
		return fmt.Errorf("KVStore Instance %s Policy sitll exists.", rs.Primary.ID)
	}

	return nil
}

var kvStoreMap = map[string]string{
	"instance_id":     CHECKSET,
	"backup_time":     "10:00Z-11:00Z",
	"backup_period.#": "2",
}

func testAccKVStoreBackupPolicy_classic(instanceType, instanceClass, engineVersion string) string {
	return fmt.Sprintf(`
	data "alicloud_kvstore_zones" "default"{
		instance_charge_type = "PostPaid"
	}
	variable "name" {
		default = "tf-testAccKVStoreBackupPolicy_classic"
	}

	resource "alicloud_kvstore_instance" "default" {
		availability_zone = data.alicloud_kvstore_zones.default.zones[length(data.alicloud_kvstore_zones.default.ids) - 1].id
		instance_name  = "${var.name}"
		security_ips = ["10.0.0.1"]
		instance_type = "%s"
		instance_class = "%s"
		engine_version = "%s"
	}
	resource "alicloud_kvstore_backup_policy" "default" {
		instance_id = "${alicloud_kvstore_instance.default.id}"
		backup_period = ["Tuesday", "Wednesday"]
		backup_time = "10:00Z-11:00Z"
	}
	`, instanceType, instanceClass, engineVersion)
}

func testAccKVStoreBackupPolicy_classicUpdatePeriod(instanceType, instanceClass, engineVersion string) string {
	return fmt.Sprintf(`
	data "alicloud_kvstore_zones" "default"{
		instance_charge_type = "PostPaid"
	}
	variable "name" {
		default = "tf-testAccKVStoreBackupPolicy_classic"
	}

	resource "alicloud_kvstore_instance" "default" {
		availability_zone = data.alicloud_kvstore_zones.default.zones[length(data.alicloud_kvstore_zones.default.ids) - 1].id
		instance_name  = "${var.name}"
		security_ips = ["10.0.0.1"]
		instance_type = "%s"
		instance_class = "%s"
		engine_version = "%s"
	}
	resource "alicloud_kvstore_backup_policy" "default" {
		instance_id = "${alicloud_kvstore_instance.default.id}"
		backup_period = ["Tuesday", "Wednesday", "Sunday"]
		backup_time = "10:00Z-11:00Z"
	}
	`, instanceType, instanceClass, engineVersion)
}

func testAccKVStoreBackupPolicy_classicUpdateTime(instanceType, instanceClass, engineVersion string) string {
	return fmt.Sprintf(`
	data "alicloud_kvstore_zones" "default"{
		instance_charge_type = "PostPaid"
	}
	variable "name" {
		default = "tf-testAccKVStoreBackupPolicy_classic"
	}

	resource "alicloud_kvstore_instance" "default" {
		availability_zone = data.alicloud_kvstore_zones.default.zones[length(data.alicloud_kvstore_zones.default.ids) - 1].id
		instance_name  = "${var.name}"
		security_ips = ["10.0.0.1"]
		instance_type = "%s"
		instance_class = "%s"
		engine_version = "%s"
	}
	resource "alicloud_kvstore_backup_policy" "default" {
		instance_id = "${alicloud_kvstore_instance.default.id}"
		backup_period = ["Tuesday", "Wednesday", "Sunday"]
		backup_time = "12:00Z-13:00Z"
	}
	`, instanceType, instanceClass, engineVersion)
}

func testAccKVStoreBackupPolicy_classicUpdateAll(instanceType, instanceClass, engineVersion string) string {
	return fmt.Sprintf(`
	data "alicloud_kvstore_zones" "default"{
		instance_charge_type = "PostPaid"
	}
	variable "name" {
		default = "tf-testAccKVStoreBackupPolicy_classic"
	}

	resource "alicloud_kvstore_instance" "default" {
		availability_zone = data.alicloud_kvstore_zones.default.zones[length(data.alicloud_kvstore_zones.default.ids) - 1].id
		instance_name  = "${var.name}"
		security_ips = ["10.0.0.1"]
		instance_type = "%s"
		instance_class = "%s"
		engine_version = "%s"
	}
	resource "alicloud_kvstore_backup_policy" "default" {
		instance_id = "${alicloud_kvstore_instance.default.id}"
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
	resource "alicloud_kvstore_instance" "default" {
		instance_class = "%s"
		instance_name  = "${var.name}"
		vswitch_id     = data.alicloud_vswitches.default.ids.0
		security_ips = ["10.0.0.1"]
		instance_type = "%s"
		engine_version = "%s"
	}
	resource "alicloud_kvstore_backup_policy" "default" {
		instance_id = "${alicloud_kvstore_instance.default.id}"
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
	resource "alicloud_kvstore_instance" "default" {
		instance_class = "%s"
		instance_name  = "${var.name}"
		vswitch_id     = data.alicloud_vswitches.default.ids.0
		security_ips = ["10.0.0.1"]
		instance_type = "%s"
		engine_version = "%s"
	}
	resource "alicloud_kvstore_backup_policy" "default" {
		instance_id = "${alicloud_kvstore_instance.default.id}"
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
	resource "alicloud_kvstore_instance" "default" {
		instance_class = "%s"
		instance_name  = "${var.name}"
		vswitch_id     = data.alicloud_vswitches.default.ids.0
		security_ips = ["10.0.0.1"]
		instance_type = "%s"
		engine_version = "%s"
	}
	resource "alicloud_kvstore_backup_policy" "default" {
		instance_id = "${alicloud_kvstore_instance.default.id}"
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
	resource "alicloud_kvstore_instance" "default" {
		instance_class = "%s"
		instance_name  = "${var.name}"
		vswitch_id     = data.alicloud_vswitches.default.ids.0
		security_ips = ["10.0.0.1"]
		instance_type = "%s"
		engine_version = "%s"
	}
	resource "alicloud_kvstore_backup_policy" "default" {
		instance_id = "${alicloud_kvstore_instance.default.id}"
		backup_period = ["Tuesday"]
		backup_time = "12:00Z-13:00Z"
	}
	`, common, instanceClass, instanceType, engineVersion)
}
