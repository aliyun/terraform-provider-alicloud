package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudElasticsearchDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudElasticsearchDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_elasticsearch.instances"),
					resource.TestCheckResourceAttr("data.alicloud_elasticsearch.instances", "instances.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_elasticsearch.instances", "instances.0.id"),
					resource.TestCheckResourceAttr("data.alicloud_elasticsearch.instances", "instances.0.instance_charge_type", string(PostPaid)),
					resource.TestCheckResourceAttr("data.alicloud_elasticsearch.instances", "instances.0.description", "tf-testESDataSource"),
					resource.TestCheckResourceAttr("data.alicloud_elasticsearch.instances", "instances.0.data_node_amount", "2"),
					resource.TestCheckResourceAttr("data.alicloud_elasticsearch.instances", "instances.0.data_node_spec", "elasticsearch.sn2ne.large"),
					resource.TestCheckResourceAttr("data.alicloud_elasticsearch.instances", "instances.0.status", string(ElasticsearchStatusActive)),
					resource.TestCheckResourceAttrSet("data.alicloud_elasticsearch.instances", "instances.0.es_version"),
					resource.TestCheckResourceAttrSet("data.alicloud_elasticsearch.instances", "instances.0.created_at"),
					resource.TestCheckResourceAttrSet("data.alicloud_elasticsearch.instances", "instances.0.updated_at"),
					resource.TestCheckResourceAttrSet("data.alicloud_elasticsearch.instances", "instances.0.vswitch_id"),
				),
			},
		},
	})
}

const testAccCheckAlicloudElasticsearchDataSourceConfig = `
data "alicloud_elasticsearch" "instances" {
  description = "${alicloud_elasticsearch.instance.description}"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

variable "creation" {
  default = "Elasticsearch"
}

variable "name" {
  default = "tf-testESDataSource"
}

resource "alicloud_vpc" "default" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vpc_id            = "${alicloud_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  availability_zone = "${data.alicloud_zones.default.zones.1.id}"
  name              = "${var.name}"
}

resource "alicloud_elasticsearch" "instance" {
  description          = "${var.name}"
  es_admin_password    = "Test@12345"
  vswitch_id           = "${alicloud_vswitch.default.id}"
  data_node_amount     = "2"
  data_node_spec       = "elasticsearch.sn2ne.large"
  data_node_disk       = "20"
  data_node_disk_type  = "cloud_ssd"
  instance_charge_type = "PostPaid"
  es_version           = "5.5.3_with_X-Pack"
}
`
