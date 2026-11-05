package bssopenapi_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudBssOpenApiPricingModuleDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(1000000, 9999999)

	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudBssOpenApiPricingModuleSourceConfig(rand, map[string]string{
			"name_regex":        `"国内月均日峰值带宽"`,
			"product_code":      `"cdn"`,
			"subscription_type": `"PayAsYouGo"`,
			"product_type":      `"CDN"`,
		}),
		FakeConfig: testAccCheckAlicloudBssOpenApiPricingModuleSourceConfig(rand, map[string]string{
			"name_regex":        `"国内月均日峰值带宽_fake"`,
			"product_code":      `"cdn"`,
			"subscription_type": `"PayAsYouGo"`,
			"product_type":      `"CDN"`,
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudBssOpenApiPricingModuleSourceConfig(rand, map[string]string{
			"ids":               `["cdn:CDN:PayAsYouGo:Cdn95China"]`,
			"product_code":      `"cdn"`,
			"subscription_type": `"PayAsYouGo"`,
			"product_type":      `"CDN"`,
		}),
		FakeConfig: testAccCheckAlicloudBssOpenApiPricingModuleSourceConfig(rand, map[string]string{
			"ids":               `["cdn:CDN:PayAsYouGo:Cdn95China_fake"]`,
			"product_code":      `"cdn"`,
			"subscription_type": `"PayAsYouGo"`,
			"product_type":      `"CDN"`,
		}),
	}

	BssOpenApiPricingModuleCheckInfo.DataSourceTestCheck(t, rand, nameRegexConf, allConf)
}

var existBssOpenApiPricingModuleMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"modules.#":                     "1",
		"modules.0.id":                  acctest.CHECKSET,
		"modules.0.code":                acctest.CHECKSET,
		"modules.0.pricing_module_name": acctest.CHECKSET,
		"modules.0.product_code":        acctest.CHECKSET,
		"modules.0.product_type":        acctest.CHECKSET,
		"modules.0.subscription_type":   acctest.CHECKSET,
	}
}

var fakeBssOpenApiPricingModuleMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"modules.#": "0",
	}
}

var BssOpenApiPricingModuleCheckInfo = acctest.DataSourceAttr{
	ResourceId:   "data.alicloud_bss_open_api_pricing_modules.default",
	ExistMapFunc: existBssOpenApiPricingModuleMapFunc,
	FakeMapFunc:  fakeBssOpenApiPricingModuleMapFunc,
}

func testAccCheckAlicloudBssOpenApiPricingModuleSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccBssOpenApiPricingModule%d"
}

data "alicloud_bss_open_api_pricing_modules" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
