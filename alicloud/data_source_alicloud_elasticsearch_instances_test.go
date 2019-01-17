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
					testAccCheckAlicloudDataSourceID("data.alicloud_elasticsearch_instances.default"),
					resource.TestCheckResourceAttr("data.alicloud_elasticsearch_instances.default", "instances.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_elasticsearch_instances.default", "instances.0.id"),
					resource.TestCheckResourceAttr("data.alicloud_elasticsearch_instances.default", "instances.0.instance_charge_type", string(PostPaid)),
					resource.TestCheckResourceAttr("data.alicloud_elasticsearch_instances.default", "instances.0.description", "tf-testAccESDataSource"),
					resource.TestCheckResourceAttr("data.alicloud_elasticsearch_instances.default", "instances.0.data_node_amount", "2"),
					resource.TestCheckResourceAttr("data.alicloud_elasticsearch_instances.default", "instances.0.data_node_spec", "elasticsearch.sn2ne.large"),
					resource.TestCheckResourceAttr("data.alicloud_elasticsearch_instances.default", "instances.0.status", string(ElasticsearchStatusActive)),
					resource.TestCheckResourceAttrSet("data.alicloud_elasticsearch_instances.default", "instances.0.version"),
					resource.TestCheckResourceAttrSet("data.alicloud_elasticsearch_instances.default", "instances.0.created_at"),
					resource.TestCheckResourceAttrSet("data.alicloud_elasticsearch_instances.default", "instances.0.updated_at"),
					resource.TestCheckResourceAttrSet("data.alicloud_elasticsearch_instances.default", "instances.0.vswitch_id"),
				),
			},
		},
	})
}

func TestAccAlicloudElasticsearchDataSource_empty(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudElasticsearchDataSourceEmpty,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_elasticsearch_instances.default"),
					resource.TestCheckResourceAttr("data.alicloud_elasticsearch_instances.default", "instances.#", "0"),
					resource.TestCheckNoResourceAttr("data.alicloud_elasticsearch_instances.default", "instances.0.id"),
					resource.TestCheckNoResourceAttr("data.alicloud_elasticsearch_instances.default", "instances.0.instance_charge_type"),
					resource.TestCheckNoResourceAttr("data.alicloud_elasticsearch_instances.default", "instances.0.description"),
					resource.TestCheckNoResourceAttr("data.alicloud_elasticsearch_instances.default", "instances.0.data_node_amount"),
					resource.TestCheckNoResourceAttr("data.alicloud_elasticsearch_instances.default", "instances.0.data_node_spec"),
					resource.TestCheckNoResourceAttr("data.alicloud_elasticsearch_instances.default", "instances.0.status"),
					resource.TestCheckNoResourceAttr("data.alicloud_elasticsearch_instances.default", "instances.0.version"),
					resource.TestCheckNoResourceAttr("data.alicloud_elasticsearch_instances.default", "instances.0.created_at"),
					resource.TestCheckNoResourceAttr("data.alicloud_elasticsearch_instances.default", "instances.0.updated_at"),
					resource.TestCheckNoResourceAttr("data.alicloud_elasticsearch_instances.default", "instances.0.vswitch_id"),
				),
			},
		},
	})
}

const testAccCheckAlicloudElasticsearchDataSourceConfig = `
data "alicloud_elasticsearch_instances" "default" {
  description_regex = "${alicloud_elasticsearch_instance.default.description}"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

variable "name" {
  default = "tf-testAccESDataSource"
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

resource "alicloud_elasticsearch_instance" "default" {
  description          = "${var.name}"
  password             = "Test@12345"
  vswitch_id           = "${alicloud_vswitch.default.id}"
  data_node_amount     = "2"
  data_node_spec       = "elasticsearch.sn2ne.large"
  data_node_disk_size  = "20"
  data_node_disk_type  = "cloud_ssd"
  instance_charge_type = "PostPaid"
  version              = "5.5.3_with_X-Pack"
}
`

const testAccCheckAlicloudElasticsearchDataSourceEmpty = `
data "alicloud_elasticsearch_instances" "default" {
  description_regex = "^tf-testAcc-fake-name"
}
`