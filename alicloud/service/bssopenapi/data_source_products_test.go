package bssopenapi_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudBssOpenApiProductDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(1000000, 9999999)

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudBssOpenApiProductSourceConfig(rand, map[string]string{
			"ids": `["cdn:CDN:PayAsYouGo"]`,
		}),
		FakeConfig: testAccCheckAlicloudBssOpenApiProductSourceConfig(rand, map[string]string{
			"ids": `["cdn:CDN:PayAsYouGo_fake"]`,
		}),
	}

	regexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudBssOpenApiProductSourceConfig(rand, map[string]string{
			"name_regex": `"内容分发网络CDN"`,
		}),
		FakeConfig: testAccCheckAlicloudBssOpenApiProductSourceConfig(rand, map[string]string{
			"name_regex": `"内容分发网络CDN_fake"`,
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudBssOpenApiProductSourceConfig(rand, map[string]string{
			"ids":        `["cdn:CDN:PayAsYouGo"]`,
			"name_regex": `"内容分发网络CDN"`,
		}),
		FakeConfig: testAccCheckAlicloudBssOpenApiProductSourceConfig(rand, map[string]string{
			"ids":        `["cdn:CDN:PayAsYouGo"]`,
			"name_regex": `"内容分发网络CDN_fake"`,
		}),
	}

	BssOpenApiProductCheckInfo.DataSourceTestCheck(t, rand, idsConf, regexConf, allConf)
}

var existBssOpenApiProductMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"products.#":                   "1",
		"products.0.id":                acctest.CHECKSET,
		"products.0.product_code":      acctest.CHECKSET,
		"products.0.product_name":      acctest.CHECKSET,
		"products.0.product_type":      acctest.CHECKSET,
		"products.0.subscription_type": acctest.CHECKSET,
	}
}

var fakeBssOpenApiProductMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"products.#": "0",
	}
}

var BssOpenApiProductCheckInfo = acctest.DataSourceAttr{
	ResourceId:   "data.alicloud_bss_open_api_products.default",
	ExistMapFunc: existBssOpenApiProductMapFunc,
	FakeMapFunc:  fakeBssOpenApiProductMapFunc,
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
