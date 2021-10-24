package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudMhubProductsDataSource(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudMhubProductsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_mhub_product.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudMhubProductsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_mhub_product.default.id}_fake"]`,
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudMhubProductsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_mhub_product.default.product_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudMhubProductsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_mhub_product.default.product_name}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudMhubProductsDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_mhub_product.default.id}"]`,
			"name_regex": `"${alicloud_mhub_product.default.product_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudMhubProductsDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_mhub_product.default.id}_fake"]`,
			"name_regex": `"${alicloud_mhub_product.default.product_name}_fake"`,
		}),
	}

	var existAlicloudMhubProductsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":      "1",
			"names.#":    "1",
			"products.#": "1",
		}
	}
	var fakeAlicloudMhubProductsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":      "0",
			"names.#":    "0",
			"products.#": "0",
		}
	}

	var AlicloudMhubProductsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_mhub_products.default",
		existMapFunc: existAlicloudMhubProductsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudMhubProductsDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
		testAccPreCheckWithRegions(t, true, connectivity.MHUBSupportRegions)
	}
	AlicloudMhubProductsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, allConf)
}
func testAccCheckAlicloudMhubProductsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf_testaccmhubproduct%d"
}

resource "alicloud_mhub_product" "default"{
  product_name = var.name
}

data "alicloud_mhub_products" "default" {	
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
