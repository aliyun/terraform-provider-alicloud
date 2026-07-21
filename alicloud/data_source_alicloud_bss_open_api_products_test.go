package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudBssOpenApiProductDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudBssOpenApiProductSourceConfig(rand, map[string]string{
			"ids": `["cdn:CDN:PayAsYouGo"]`,
		}),
		fakeConfig: testAccCheckAlicloudBssOpenApiProductSourceConfig(rand, map[string]string{
			"ids": `["cdn:CDN:PayAsYouGo_fake"]`,
		}),
	}

	regexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudBssOpenApiProductSourceConfig(rand, map[string]string{
			"name_regex": `"内容分发网络CDN"`,
		}),
		fakeConfig: testAccCheckAlicloudBssOpenApiProductSourceConfig(rand, map[string]string{
			"name_regex": `"内容分发网络CDN_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudBssOpenApiProductSourceConfig(rand, map[string]string{
			"ids":        `["cdn:CDN:PayAsYouGo"]`,
			"name_regex": `"内容分发网络CDN"`,
		}),
		fakeConfig: testAccCheckAlicloudBssOpenApiProductSourceConfig(rand, map[string]string{
			"ids":        `["cdn:CDN:PayAsYouGo"]`,
			"name_regex": `"内容分发网络CDN_fake"`,
		}),
	}

	BssOpenApiProductCheckInfo.dataSourceTestCheck(t, rand, idsConf, regexConf, allConf)
}

var existBssOpenApiProductMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"products.#":                   "1",
		"products.0.id":                CHECKSET,
		"products.0.product_code":      CHECKSET,
		"products.0.product_name":      CHECKSET,
		"products.0.product_type":      CHECKSET,
		"products.0.subscription_type": CHECKSET,
	}
}

var fakeBssOpenApiProductMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"products.#": "0",
	}
}

var BssOpenApiProductCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_bss_open_api_products.default",
	existMapFunc: existBssOpenApiProductMapFunc,
	fakeMapFunc:  fakeBssOpenApiProductMapFunc,
}

func testAccCheckAlicloudBssOpenApiProductSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccBssOpenApiProduct%d"
}

data "alicloud_bss_open_api_products" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
