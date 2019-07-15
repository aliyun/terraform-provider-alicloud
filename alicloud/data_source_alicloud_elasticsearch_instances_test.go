package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
)

func TestAccAlicloudElasticsearchDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000, 9999)
	resourceId := "data.alicloud_elasticsearch_instances.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testAcc%s%d", defaultRegionToTest, rand),
		dataSourceElasticsearchConfigDependence)

	descriptionRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"description_regex": "${alicloud_elasticsearch_instance.default.description}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"description_regex": "${alicloud_elasticsearch_instance.default.description}-F",
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_elasticsearch_instance.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_elasticsearch_instance.default.id}-F"},
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"description_regex": "${alicloud_elasticsearch_instance.default.description}",
			"ids":               []string{"${alicloud_elasticsearch_instance.default.id}"},
			"version":           "5.5.3_with_X-Pack",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"description_regex": "${alicloud_elasticsearch_instance.default.description}-F",
			"ids":               []string{"${alicloud_elasticsearch_instance.default.id}"},
			"version":           "6.7.0_with_X-Pack",
		}),
	}

	var elasticsearchCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existElasticsearchMapFunc,
		fakeMapFunc:  fakeElasticsearchMapFunc,
	}

	elasticsearchCheckInfo.dataSourceTestCheck(t, rand, descriptionRegexConf, idsConf, allConf)
}

var existElasticsearchMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":                            "1",
		"ids.0":                            CHECKSET,
		"descriptions.#":                   "1",
		"descriptions.0":                   fmt.Sprintf("tf-testAcc%s%d", defaultRegionToTest, rand),
		"instances.#":                      "1",
		"instances.0.id":                   CHECKSET,
		"instances.0.instance_charge_type": string(PostPaid),
		"instances.0.description":          fmt.Sprintf("tf-testAcc%s%d", defaultRegionToTest, rand),
		"instances.0.data_node_amount":     "2",
		"instances.0.data_node_spec":       "elasticsearch.sn2ne.large",
		"instances.0.status":               string(ElasticsearchStatusActive),
		"instances.0.version":              CHECKSET,
		"instances.0.created_at":           CHECKSET,
		"instances.0.updated_at":           CHECKSET,
		"instances.0.vswitch_id":           CHECKSET,
	}
}

var fakeElasticsearchMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"instances.#":    "0",
		"ids.#":          "0",
		"descriptions.#": "0",
	}
}

func dataSourceElasticsearchConfigDependence(name string) string {
	return fmt.Sprintf(`

data "alicloud_zones" "default" {
  available_resource_creation = "Elasticsearch"
}

variable "name" {
  default = "%s"
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
`, name)
}
