package market_test

import (
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudMarketProductsDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(1000000, 9999999)
	resourceId := "data.alicloud_market_products.default"

	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId,
		"",
		dataSourceMarketProductsConfigDependence)

	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"name_regex":   "BatchCompute",
			"product_type": "MIRROR",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"name_regex":   "BatchCompute_fake",
			"product_type": "MIRROR",
		}),
	}
	//idsConf := acctest.DataSourceTestAccConfig{
	//	ExistConfig: testAccConfig(map[string]interface{}{
	//		"ids":          []string{"cmjj022644"},
	//		"product_type": "MIRROR",
	//	}),
	//	FakeConfig: testAccConfig(map[string]interface{}{
	//		"ids":          []string{"cmjj022644_fake"},
	//		"product_type": "MIRROR",
	//	}),
	//}
	sortConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"sort":         "user_count-desc",
			"product_type": "MIRROR",
		}),
	}
	categoryIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"category_id": "53366009",
		}),
	}
	productTypeConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"product_type": "MIRROR",
		}),
	}
	searchTermConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"search_term": "Image",
		}),
	}

	//allConf := acctest.DataSourceTestAccConfig{
	//	ExistConfig: testAccConfig(map[string]interface{}{
	//		"sort":         "created_on-desc",
	//		"category_id":  "56024006",
	//		"product_type": "MIRROR",
	//		"name_regex":   "BatchCompute",
	//		"search_term":  "Image",
	//	}),
	//	FakeConfig: testAccConfig(map[string]interface{}{
	//		"sort":         "created_on-desc",
	//		"category_id":  "56024006",
	//		"product_type": "MIRROR",
	//		"name_regex":   "BatchCompute_fake",
	//		"search_term":  "Image",
	//	}),
	//}

	var existMarketProductsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                        acctest.CHECKSET,
			"ids.0":                        acctest.CHECKSET,
			"products.#":                   acctest.CHECKSET,
			"products.0.category_id":       acctest.CHECKSET,
			"products.0.code":              acctest.CHECKSET,
			"products.0.delivery_way":      acctest.CHECKSET,
			"products.0.image_url":         acctest.CHECKSET,
			"products.0.name":              acctest.CHECKSET,
			"products.0.score":             acctest.CHECKSET,
			"products.0.short_description": acctest.CHECKSET,
			"products.0.supplier_id":       acctest.CHECKSET,
			"products.0.supplier_name":     acctest.CHECKSET,
			"products.0.target_url":        acctest.CHECKSET,
		}
	}

	var fakeMarketProductsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":      "0",
			"products.#": "0",
		}
	}

	var pvtzZoneRecordsCheckInfo = acctest.DataSourceAttr{
		ResourceId:   resourceId,
		ExistMapFunc: existMarketProductsMapFunc,
		FakeMapFunc:  fakeMarketProductsMapFunc,
	}

	pvtzZoneRecordsCheckInfo.DataSourceTestCheck(t, rand, nameRegexConf, sortConf, categoryIdConf, productTypeConf, searchTermConf)
}

func dataSourceMarketProductsConfigDependence(name string) string {
	return ""
}
