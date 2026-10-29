package market_test

import (
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudMarketProductDataSource(t *testing.T) {
	acctest.CheckoutSupportedRegions(t, true, connectivity.TestSalveRegions)
	rand := tfacctest.RandIntRange(1000000, 9999999)
	resourceId := "data.alicloud_market_product.default"

	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId,
		"",
		dataSourceMarketProductConfigDependence)
	basicConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"product_code": "${data.alicloud_market_products.default.ids.0}",
		}),
	}
	var existMarketProductMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"product.#":             "1",
			"product.0.code":        acctest.CHECKSET,
			"product.0.name":        acctest.CHECKSET,
			"product.0.description": acctest.CHECKSET,
			"product.0.skus.#":      acctest.CHECKSET,
		}
	}

	var fakeMarketProductMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"product.#": "0",
		}
	}

	var marketProductCheckInfo = acctest.DataSourceAttr{
		ResourceId:   resourceId,
		ExistMapFunc: existMarketProductMapFunc,
		FakeMapFunc:  fakeMarketProductMapFunc,
	}

	marketProductCheckInfo.DataSourceTestCheck(t, rand, basicConf)
}

func dataSourceMarketProductConfigDependence(name string) string {
	return `
		data "alicloud_market_products" "default" {
			name_regex = "BatchCompute"
			product_type = "MIRROR"
		}
`
}
