package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudBssOpenApiPricingModuleDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudBssOpenApiPricingModuleSourceConfig(rand, map[string]string{
			"name_regex":        `"国内月均日峰值带宽"`,
			"product_code":      `"cdn"`,
			"subscription_type": `"PayAsYouGo"`,
			"product_type":      `"CDN"`,
		}),
		fakeConfig: testAccCheckAlicloudBssOpenApiPricingModuleSourceConfig(rand, map[string]string{
			"name_regex":        `"国内月均日峰值带宽_fake"`,
			"product_code":      `"cdn"`,
			"subscription_type": `"PayAsYouGo"`,
			"product_type":      `"CDN"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudBssOpenApiPricingModuleSourceConfig(rand, map[string]string{
			"ids":               `["cdn:CDN:PayAsYouGo:Cdn95China"]`,
			"product_code":      `"cdn"`,
			"subscription_type": `"PayAsYouGo"`,
			"product_type":      `"CDN"`,
		}),
		fakeConfig: testAccCheckAlicloudBssOpenApiPricingModuleSourceConfig(rand, map[string]string{
			"ids":               `["cdn:CDN:PayAsYouGo:Cdn95China_fake"]`,
			"product_code":      `"cdn"`,
			"subscription_type": `"PayAsYouGo"`,
			"product_type":      `"CDN"`,
		}),
	}

	BssOpenApiPricingModuleCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, allConf)
}

var existBssOpenApiPricingModuleMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"modules.#":                     "1",
		"modules.0.id":                  CHECKSET,
		"modules.0.code":                CHECKSET,
		"modules.0.pricing_module_name": CHECKSET,
		"modules.0.product_code":        CHECKSET,
		"modules.0.product_type":        CHECKSET,
		"modules.0.subscription_type":   CHECKSET,
	}
}

var fakeBssOpenApiPricingModuleMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"modules.#": "0",
	}
}

var BssOpenApiPricingModuleCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_bss_open_api_pricing_modules.default",
	existMapFunc: existBssOpenApiPricingModuleMapFunc,
	fakeMapFunc:  fakeBssOpenApiPricingModuleMapFunc,
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
