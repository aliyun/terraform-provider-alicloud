package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudElasticsearchZonesDataSource_basic(t *testing.T) {
	rand := acctest.RandInt()
	resourceId := "data.alicloud_elasticsearch_zones.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, "", dataSourceElasticsearchZonesConfigDependence)

	multiConfig := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"multi": "true",
		}),
	}

	var existElasticsearchZonesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   CHECKSET,
			"zones.#": CHECKSET,
		}
	}

	var fakeElasticsearchZonesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"zones.#": "0",
		}
	}

	var elasticsearchZonesCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existElasticsearchZonesMapFunc,
		fakeMapFunc:  fakeElasticsearchZonesMapFunc,
	}

	elasticsearchZonesCheckInfo.dataSourceTestCheck(t, rand, multiConfig)
}

func dataSourceElasticsearchZonesConfigDependence(name string) string {
	return ""
}
