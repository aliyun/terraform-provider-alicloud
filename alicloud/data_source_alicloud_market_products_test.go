package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
)

func TestAccAlicloudMarketProductsDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)
	resourceId := "data.alicloud_market_products.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		"",
		dataSourceMarketProductsConfigDependence)

	basicConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"product_type": "MIRROR",
			"category_id":  "53366009",
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex":   "BatchCompute Ubuntu14.04",
			"product_type": "MIRROR",
			"category_id":  "53366009",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex":   "BatchCompute_fake",
			"product_type": "MIRROR",
			"category_id":  "53366009",
		}),
	}
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":          []string{"cmjj022644"},
			"product_type": "MIRROR",
			"category_id":  "53366009",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":          []string{"cmjj022644_fake"},
			"product_type": "MIRROR",
			"category_id":  "53366009",
		}),
	}
	sortConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"sort":         "user_count-desc",
			"product_type": "MIRROR",
			"category_id":  "53366009",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"sort":         "created_on-desc",
			"category_id":  "56024006",
			"product_type": "MIRROR",
			"name_regex":   "SQLServer2016_Ent_FullFeature_winupdate",
			"ids":          []string{"cmjj031537"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"sort":         "created_on-desc",
			"category_id":  "56024006",
			"product_type": "MIRROR",
			"name_regex":   "SQLServer2016_Ent_FullFeature_winupdate_fake",
			"ids":          []string{"cmjj031537"},
		}),
	}

	var existMarketProductsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                 CHECKSET,
			"ids.0":                 CHECKSET,
			"products.#":            CHECKSET,
			"products.0.code":       CHECKSET,
			"products.0.name":       CHECKSET,
			"products.0.target_url": CHECKSET,
		}
	}

	var fakeMarketProductsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":      "0",
			"products.#": "0",
		}
	}

	var pvtzZoneRecordsCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existMarketProductsMapFunc,
		fakeMapFunc:  fakeMarketProductsMapFunc,
	}

	pvtzZoneRecordsCheckInfo.dataSourceTestCheck(t, rand, basicConf, nameRegexConf, idsConf, sortConf, allConf)
}

func dataSourceMarketProductsConfigDependence(name string) string {
	return ""
}
