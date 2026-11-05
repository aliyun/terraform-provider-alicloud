package mhub_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudMhubProductsDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(10000, 99999)
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudMhubProductsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_mhub_product.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudMhubProductsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_mhub_product.default.id}_fake"]`,
		}),
	}

	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudMhubProductsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_mhub_product.default.product_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudMhubProductsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_mhub_product.default.product_name}_fake"`,
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudMhubProductsDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_mhub_product.default.id}"]`,
			"name_regex": `"${alicloud_mhub_product.default.product_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudMhubProductsDataSourceName(rand, map[string]string{
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

	var AlicloudMhubProductsCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_mhub_products.default",
		ExistMapFunc: existAlicloudMhubProductsDataSourceNameMapFunc,
		FakeMapFunc:  fakeAlicloudMhubProductsDataSourceNameMapFunc,
	}

	preCheck := func() {
		acctest.PreCheck(t)
		acctest.PreCheckWithRegions(t, true, connectivity.MHUBSupportRegions)
	}
	AlicloudMhubProductsCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, allConf)
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
