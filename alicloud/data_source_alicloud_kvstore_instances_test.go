package alicloud

import (
	"testing"

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
				Config: testAccCheckAlicloudRKVInstancesDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_kvstore_instances.rkvs"),
					resource.TestCheckResourceAttr("data.alicloud_kvstore_instances.rkvs", "instances.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_kvstore_instances.rkvs", "instances.0.id"),
					resource.TestCheckResourceAttr("data.alicloud_kvstore_instances.rkvs", "instances.0.instance_class", "redis.master.small.default"),
					resource.TestCheckResourceAttr("data.alicloud_kvstore_instances.rkvs", "instances.0.name", "tf-testAccCheckAlicloudRKVInstancesDataSourceConfig"),
					resource.TestCheckResourceAttr("data.alicloud_kvstore_instances.rkvs", "instances.0.instance_type", "Redis"),
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

const testAccCheckAlicloudRKVInstancesDataSourceConfig = `
data "alicloud_kvstore_instances" "rkvs" {
  name_regex = "${alicloud_kvstore_instance.rkv.instance_name}"
}

data "alicloud_zones" "default" {
	available_resource_creation = "KVStore"
}
variable "name" {
	default = "tf-testAccCheckAlicloudRKVInstancesDataSourceConfig"
}
resource "alicloud_vpc" "foo" {
	name = "${var.name}"
	cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "foo" {
 	vpc_id = "${alicloud_vpc.foo.id}"
 	cidr_block = "172.16.0.0/24"
 	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
 	name = "${var.name}"
}

resource "alicloud_kvstore_instance" "rkv" {
	instance_class = "redis.master.small.default"
	instance_name  = "${var.name}"
	password       = "Test12345"
	vswitch_id     = "${alicloud_vswitch.foo.id}"
	private_ip     = "172.16.0.10"
}
`

const testAccCheckAlicloudRKVInstancesDataSourceEmpty = `
data "alicloud_kvstore_instances" "rkvs" {
  name_regex = "^tf-testacc-fake-name"
}
`
