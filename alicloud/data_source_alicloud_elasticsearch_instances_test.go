package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudElasticsearchDataSource(t *testing.T) {
	rand := acctest.RandInt()
	resourceId := "data.alicloud_elasticsearch_instances.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testAccES%s%d", defaultRegionToTest, rand)[:30],
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

	tagsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"tags": map[string]interface{}{
				"Created": "TF",
				"For":     "acceptance test",
			},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"tags": map[string]interface{}{
				"Created": "TF-fake",
				"For":     "acceptance test",
			},
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"description_regex": "${alicloud_elasticsearch_instance.default.description}",
			"ids":               []string{"${alicloud_elasticsearch_instance.default.id}"},
			"version":           "5.5.3_with_X-Pack",
			"tags": map[string]interface{}{
				"Created": "TF",
				"For":     "acceptance test",
			},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"description_regex": "${alicloud_elasticsearch_instance.default.description}-F",
			"ids":               []string{"${alicloud_elasticsearch_instance.default.id}"},
			"version":           "6.7.0_with_X-Pack",
			"tags": map[string]interface{}{
				"Created": "TF-fake",
				"For":     "acceptance test",
			},
		}),
	}

	var elasticsearchCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existElasticsearchMapFunc,
		fakeMapFunc:  fakeElasticsearchMapFunc,
	}
	preCheck := func() {
		testAccPreCheckWithNoDefaultVpc(t)
	}
	elasticsearchCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, descriptionRegexConf, idsConf, tagsConf, allConf)
}

var existElasticsearchMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":                            "1",
		"ids.0":                            CHECKSET,
		"descriptions.#":                   "1",
		"descriptions.0":                   fmt.Sprintf("tf-testAccES%s%d", defaultRegionToTest, rand)[:30],
		"instances.#":                      "1",
		"instances.0.id":                   CHECKSET,
		"instances.0.instance_charge_type": string(PostPaid),
		"instances.0.description":          fmt.Sprintf("tf-testAccES%s%d", defaultRegionToTest, rand)[:30],
		"instances.0.data_node_amount":     "2",
		"instances.0.data_node_spec":       "elasticsearch.sn2ne.large",
		"instances.0.status":               "active",
		"instances.0.version":              CHECKSET,
		"instances.0.tags.%":               CHECKSET,
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

data "alicloud_vswitches" "default" {
  zone_id = data.alicloud_zones.default.ids[0]
  name_regex = "default-tf--testAcc-00"
}

resource "alicloud_elasticsearch_instance" "default" {
  description          = "${var.name}"
  password             = "Yourpassword1234"
  vswitch_id           = "${data.alicloud_vswitches.default.ids.0}"
  data_node_amount     = "2"
  data_node_spec       = "elasticsearch.sn2ne.large"
  data_node_disk_size  = "20"
  data_node_disk_type  = "cloud_ssd"
  instance_charge_type = "PostPaid"
  version              = "5.5.3_with_X-Pack"
  tags                 = {
	  "Created": "TF",
	  "For":     "acceptance test",
  }
}
`, name)
}
