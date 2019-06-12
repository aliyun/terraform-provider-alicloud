package alicloud

import (
	"fmt"
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
				Config: testAccElasticsearchDatasourceProfile(BasicElasticsearchDatasourceDescription),
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
					resource.TestCheckResourceAttr("data.alicloud_elasticsearch_instances.default", "ids.#", "1"),
				),
			},
			{
				Config: testAccElasticsearchDatasourceProfile(BasicElasticsearchDatasourceDescriptionEmpty),
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
					resource.TestCheckResourceAttr("data.alicloud_elasticsearch_instances.default", "ids.#", "0"),
				),
			},
			{
				Config: testAccElasticsearchDatasourceProfile(BasicElasticsearchDatasourceDescriptionExistsRegex),
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
					resource.TestCheckResourceAttr("data.alicloud_elasticsearch_instances.default", "ids.#", "1"),
				),
			},
			{
				Config: testAccElasticsearchDatasourceProfileWithVersion(BasicElasticsearchDatasourceDescription, string(ESVersion553WithXPack)),
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
					resource.TestCheckResourceAttr("data.alicloud_elasticsearch_instances.default", "ids.#", "1"),
				),
			},
			{
				Config: testAccElasticsearchDatasourceProfileWithVersion(BasicElasticsearchDatasourceDescription, string(ESVersion632WithXPack)),
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
					resource.TestCheckResourceAttr("data.alicloud_elasticsearch_instances.default", "ids.#", "0"),
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
					resource.TestCheckResourceAttr("data.alicloud_elasticsearch_instances.default", "ids.#", "0"),
				),
			},
		},
	})
}

func testAccElasticsearchDatasourceProfile(descriptionRegex string) string {
	return fmt.Sprintf(`
data "alicloud_elasticsearch_instances" "default" {
  description_regex = "%s"
}

data "alicloud_zones" "default" {
  available_resource_creation = "Elasticsearch"
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
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  name              = "${var.name}"
}

resource "alicloud_elasticsearch_instance" "default" {
  description          = "${var.name}"
  password             = "Yourpassword1234"
  vswitch_id           = "${alicloud_vswitch.default.id}"
  data_node_amount     = "2"
  data_node_spec       = "elasticsearch.sn2ne.large"
  data_node_disk_size  = "20"
  data_node_disk_type  = "cloud_ssd"
  instance_charge_type = "PostPaid"
  version              = "5.5.3_with_X-Pack"
}
`, descriptionRegex)
}

func testAccElasticsearchDatasourceProfileWithVersion(description_regex string, version string) string {
	return fmt.Sprintf(`
data "alicloud_elasticsearch_instances" "default" {
  description_regex = "%s"
  version           = "%s"
}

data "alicloud_zones" "default" {
  available_resource_creation = "Elasticsearch"
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
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  name              = "${var.name}"
}

resource "alicloud_elasticsearch_instance" "default" {
  description          = "${var.name}"
  password             = "Yourpassword1234"
  vswitch_id           = "${alicloud_vswitch.default.id}"
  data_node_amount     = "2"
  data_node_spec       = "elasticsearch.sn2ne.large"
  data_node_disk_size  = "20"
  data_node_disk_type  = "cloud_ssd"
  instance_charge_type = "PostPaid"
  version              = "5.5.3_with_X-Pack"
}
`, description_regex, version)
}

const BasicElasticsearchDatasourceDescription = "${alicloud_elasticsearch_instance.default.description}"
const BasicElasticsearchDatasourceDescriptionEmpty = "^tf-testAccESDataSourcea.*"
const BasicElasticsearchDatasourceDescriptionExistsRegex = "^tf-testAccESDataSour.*"

const testAccCheckAlicloudElasticsearchDataSourceEmpty = `
data "alicloud_elasticsearch_instances" "default" {
  description_regex = "^tf-testAcc-fake-name"
}
`
