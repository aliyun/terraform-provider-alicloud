package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
)

func TestAccAlicloudMarketProductDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)
	resourceId := "data.alicloud_market_product.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		"",
		dataSourceMarketProductConfigDependence)
	basicConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"product_code": "cmapi022206",
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"product_code": "cmapi033136",
		}),
	}

	var existMarketProductMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"product.#":             "1",
			"product.0.code":        CHECKSET,
			"product.0.name":        CHECKSET,
			"product.0.description": CHECKSET,
			"product.0.skus.#":      CHECKSET,
		}
	}

	var fakeMarketProductMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"product.#": "0",
		}
	}

	var marketProductCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existMarketProductMapFunc,
		fakeMapFunc:  fakeMarketProductMapFunc,
	}

	marketProductCheckInfo.dataSourceTestCheck(t, rand, basicConf, allConf)
}

func dataSourceMarketProductConfigDependence(name string) string {
	return ""
}
