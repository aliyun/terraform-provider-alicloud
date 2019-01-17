package alicloud

import (
	"testing"

	"fmt"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudKVStoreInstancesDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudRKVInstancesDataSourceConfig(KVStoreCommonTestCase, string(KVStoreRedis), redisInstanceClassForTest, string(KVStore2Dot8)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_kvstore_instances.rkvs"),
					resource.TestCheckResourceAttr("data.alicloud_kvstore_instances.rkvs", "instances.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_kvstore_instances.rkvs", "instances.0.id"),
					resource.TestCheckResourceAttr("data.alicloud_kvstore_instances.rkvs", "instances.0.instance_class", redisInstanceClassForTest),
					resource.TestCheckResourceAttr("data.alicloud_kvstore_instances.rkvs", "instances.0.name", "tf-testAccCheckAlicloudRKVInstancesDataSourceConfig"),
					resource.TestCheckResourceAttr("data.alicloud_kvstore_instances.rkvs", "instances.0.instance_type", string(KVStoreRedis)),
					resource.TestCheckResourceAttr("data.alicloud_kvstore_instances.rkvs", "instances.0.charge_type", string(PostPaid)),
					resource.TestCheckResourceAttrSet("data.alicloud_kvstore_instances.rkvs", "instances.0.region_id"),
					resource.TestCheckResourceAttrSet("data.alicloud_kvstore_instances.rkvs", "instances.0.create_time"),
					resource.TestCheckResourceAttr("data.alicloud_kvstore_instances.rkvs", "instances.0.expire_time", ""),
					resource.TestCheckResourceAttr("data.alicloud_kvstore_instances.rkvs", "instances.0.status", string(Normal)),
					resource.TestCheckResourceAttrSet("data.alicloud_kvstore_instances.rkvs", "instances.0.availability_zone"),
					resource.TestCheckResourceAttrSet("data.alicloud_kvstore_instances.rkvs", "instances.0.vpc_id"),
					resource.TestCheckResourceAttrSet("data.alicloud_kvstore_instances.rkvs", "instances.0.vswitch_id"),
					resource.TestCheckResourceAttrSet("data.alicloud_kvstore_instances.rkvs", "instances.0.private_ip"),
					resource.TestCheckResourceAttrSet("data.alicloud_kvstore_instances.rkvs", "instances.0.port"),
					resource.TestCheckResourceAttrSet("data.alicloud_kvstore_instances.rkvs", "instances.0.user_name"),
					resource.TestCheckResourceAttrSet("data.alicloud_kvstore_instances.rkvs", "instances.0.capacity"),
					resource.TestCheckResourceAttrSet("data.alicloud_kvstore_instances.rkvs", "instances.0.bandwidth"),
					resource.TestCheckResourceAttrSet("data.alicloud_kvstore_instances.rkvs", "instances.0.connections"),
					resource.TestCheckResourceAttrSet("data.alicloud_kvstore_instances.rkvs", "instances.0.connection_domain"),
				),
			},
		},
	})
}

func TestAccAlicloudKVStoreInstancesDataSourceEmpty(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudRKVInstancesDataSourceEmpty,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_kvstore_instances.rkvs"),
					resource.TestCheckResourceAttr("data.alicloud_kvstore_instances.rkvs", "instances.#", "0"),
					resource.TestCheckNoResourceAttr("data.alicloud_kvstore_instances.rkvs", "instances.0.id"),
					resource.TestCheckNoResourceAttr("data.alicloud_kvstore_instances.rkvs", "instances.0.instance_class"),
					resource.TestCheckNoResourceAttr("data.alicloud_kvstore_instances.rkvs", "instances.0.name"),
					resource.TestCheckNoResourceAttr("data.alicloud_kvstore_instances.rkvs", "instances.0.instance_type"),
					resource.TestCheckNoResourceAttr("data.alicloud_kvstore_instances.rkvs", "instances.0.charge_type"),
					resource.TestCheckNoResourceAttr("data.alicloud_kvstore_instances.rkvs", "instances.0.region_id"),
					resource.TestCheckNoResourceAttr("data.alicloud_kvstore_instances.rkvs", "instances.0.create_time"),
					resource.TestCheckNoResourceAttr("data.alicloud_kvstore_instances.rkvs", "instances.0.expire_time"),
					resource.TestCheckNoResourceAttr("data.alicloud_kvstore_instances.rkvs", "instances.0.status"),
					resource.TestCheckNoResourceAttr("data.alicloud_kvstore_instances.rkvs", "instances.0.availability_zone"),
					resource.TestCheckNoResourceAttr("data.alicloud_kvstore_instances.rkvs", "instances.0.vpc_id"),
					resource.TestCheckNoResourceAttr("data.alicloud_kvstore_instances.rkvs", "instances.0.vswitch_id"),
					resource.TestCheckNoResourceAttr("data.alicloud_kvstore_instances.rkvs", "instances.0.private_ip"),
					resource.TestCheckNoResourceAttr("data.alicloud_kvstore_instances.rkvs", "instances.0.port"),
					resource.TestCheckNoResourceAttr("data.alicloud_kvstore_instances.rkvs", "instances.0.user_name"),
					resource.TestCheckNoResourceAttr("data.alicloud_kvstore_instances.rkvs", "instances.0.capacity"),
					resource.TestCheckNoResourceAttr("data.alicloud_kvstore_instances.rkvs", "instances.0.bandwidth"),
					resource.TestCheckNoResourceAttr("data.alicloud_kvstore_instances.rkvs", "instances.0.connections"),
					resource.TestCheckNoResourceAttr("data.alicloud_kvstore_instances.rkvs", "instances.0.connection_domain"),
				),
			},
		},
	})
}

func testAccCheckAlicloudRKVInstancesDataSourceConfig(common, instanceType, instanceClass, engineVersion string) string {
	return fmt.Sprintf(`
	%s
	data "alicloud_kvstore_instances" "rkvs" {
	  name_regex = "${alicloud_kvstore_instance.rkv.instance_name}"
	}
	variable "creation" {
		default = "KVStore"
	}
	variable "multi_az" {
		default = "false"
	}
	variable "name" {
		default = "tf-testAccCheckAlicloudRKVInstancesDataSourceConfig"
	}

	resource "alicloud_kvstore_instance" "rkv" {
		instance_class = "%s"
		instance_name  = "${var.name}"
		vswitch_id     = "${alicloud_vswitch.default.id}"
		private_ip     = "172.16.0.10"
		security_ips = ["10.0.0.1"]
		instance_type = "%s"
		engine_version = "%s"
	}
	`, common, instanceClass, instanceType, engineVersion)
}

const testAccCheckAlicloudRKVInstancesDataSourceEmpty = `
data "alicloud_kvstore_instances" "rkvs" {
  name_regex = "^tf-testacc-fake-name"
}
`
