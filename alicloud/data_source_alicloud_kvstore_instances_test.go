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
					resource.TestCheckResourceAttr("data.alicloud_kvstore_instances.rkvs", "instances.0.instance_class", "redis.master.small.default"),
					resource.TestCheckResourceAttr("data.alicloud_kvstore_instances.rkvs", "instances.0.name", "tf-testAccCheckAlicloudRKVInstancesDataSourceConfig"),
					resource.TestCheckResourceAttr("data.alicloud_kvstore_instances.rkvs", "instances.0.instance_type", "Redis"),
					resource.TestCheckResourceAttr("data.alicloud_kvstore_instances.rkvs", "instances.0.charge_type", string(PostPaid)),
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
