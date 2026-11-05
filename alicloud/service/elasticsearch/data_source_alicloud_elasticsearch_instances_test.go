package elasticsearch_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudElasticsearchDataSource(t *testing.T) {
	defer acctest.CheckoutAccount(t, false)
	acctest.CheckoutAccount(t, true)
	rand := tfacctest.RandIntRange(1000000, 9999999)
	resourceId := "data.alicloud_elasticsearch_instances.default"

	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testAccES%d", rand),
		dataSourceElasticsearchConfigDependence)

	descriptionRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"description_regex": "${alicloud_elasticsearch_instance.default.description}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"description_regex": "${alicloud_elasticsearch_instance.default.description}-F",
		}),
	}

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_elasticsearch_instance.default.id}"},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_elasticsearch_instance.default.id}-F"},
		}),
	}

	tagsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"tags": map[string]interface{}{
				"Created": "TF",
				"For":     "acceptance test",
			},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"tags": map[string]interface{}{
				"Created": "TF-fake",
				"For":     "acceptance test",
			},
		}),
	}

	// The version filter dimension is omitted here: the data source `version`
	// argument's schema ValidateFunc only permits EOL/unpurchasable versions
	// (5.5.3/6.3.2/6.7.0_with_X-Pack), none of which can be ordered today
	// (instance_*_test.go uses only 7.x; 5.5.3 ordering failed in cn-beijing with
	// both cloud_ssd and cloud_essd). Modifying the production schema is out of
	// scope, so the combined filter keeps description_regex/ids/tags, which are
	// the dimensions exercisable against a real purchasable instance.
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"description_regex": "${alicloud_elasticsearch_instance.default.description}",
			"ids":               []string{"${alicloud_elasticsearch_instance.default.id}"},
			"tags": map[string]interface{}{
				"Created": "TF",
				"For":     "acceptance test",
			},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"description_regex": "${alicloud_elasticsearch_instance.default.description}-F",
			"ids":               []string{"${alicloud_elasticsearch_instance.default.id}"},
			"tags": map[string]interface{}{
				"Created": "TF-fake",
				"For":     "acceptance test",
			},
		}),
	}

	var elasticsearchCheckInfo = acctest.DataSourceAttr{
		ResourceId:   resourceId,
		ExistMapFunc: existElasticsearchMapFunc,
		FakeMapFunc:  fakeElasticsearchMapFunc,
	}
	// ES 5.5.3 is not orderable in the default cn-beijing region; force cn-chengdu,
	// where Elasticsearch is provisionable (see TestAccAliCloudElasticsearchInstance_basic).
	elasticsearchCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, func() {
		acctest.PreCheckWithRegions(t, true, []connectivity.Region{connectivity.ChengDu})
	}, descriptionRegexConf, idsConf, tagsConf, allConf)
}

var existElasticsearchMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":                            "1",
		"ids.0":                            acctest.CHECKSET,
		"descriptions.#":                   "1",
		"descriptions.0":                   fmt.Sprintf("tf-testAccES%d", rand),
		"instances.#":                      "1",
		"instances.0.id":                   acctest.CHECKSET,
		"instances.0.instance_charge_type": "PostPaid",
		"instances.0.description":          fmt.Sprintf("tf-testAccES%d", rand),
		"instances.0.data_node_amount":     "2",
		"instances.0.data_node_spec":       "elasticsearch.sn1ne.large",
		"instances.0.status":               "active",
		"instances.0.version":              acctest.CHECKSET,
		"instances.0.tags.%":               acctest.CHECKSET,
		"instances.0.created_at":           acctest.CHECKSET,
		"instances.0.updated_at":           acctest.CHECKSET,
		"instances.0.vswitch_id":           acctest.CHECKSET,
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
variable "name" {
  default = "%s"
}

data "alicloud_elasticsearch_zones" "default" {}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "default" {
  vswitch_name = var.name
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "172.16.0.0/21"
  zone_id      = data.alicloud_elasticsearch_zones.default.ids.0
}

resource "alicloud_elasticsearch_instance" "default" {
  description                      = var.name
  password                         = "Yourpassword1234"
  vswitch_id                       = alicloud_vswitch.default.id
  data_node_amount                 = "2"
  data_node_spec                   = "elasticsearch.sn1ne.large"
  data_node_disk_size              = "20"
  data_node_disk_type              = "cloud_essd"
  data_node_disk_performance_level = "PL1"
  instance_charge_type             = "PostPaid"
  zone_count                       = "1"
  version                          = "7.16.2_with_X-Pack"
  kibana_node_spec                 = "elasticsearch.sn1ne.large"
  tags = {
    "Created" = "TF"
    "For"     = "acceptance test"
  }
}
`, name)
}
