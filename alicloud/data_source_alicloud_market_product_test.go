package alicloud

import (
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudMarketProductDataSource(t *testing.T) {
	checkoutSupportedRegions(t, true, connectivity.TestSalveRegions)
	rand := acctest.RandIntRange(1000000, 9999999)
	resourceId := "data.alicloud_market_product.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		"",
		dataSourceMarketProductConfigDependence)
	basicConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"product_code": "${data.alicloud_market_products.default.ids.0}",
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

	marketProductCheckInfo.dataSourceTestCheck(t, rand, basicConf)
}

func dataSourceMarketProductConfigDependence(name string) string {
	return `
		data "alicloud_market_products" "default" {
			name_regex = "BatchCompute"
			product_type = "MIRROR"
		}
`
}
