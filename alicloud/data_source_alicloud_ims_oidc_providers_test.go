// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudImsOidcProviderDataSource(t *testing.T) {
	testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	rand := acctest.RandIntRange(1000000, 9999999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudImsOidcProviderSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_ims_oidc_provider.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudImsOidcProviderSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_ims_oidc_provider.default.id}_fake"]`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudImsOidcProviderSourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_ims_oidc_provider.default.id}"]`,
			"name_regex": `"${alicloud_ims_oidc_provider.default.oidc_provider_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudImsOidcProviderSourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_ims_oidc_provider.default.id}_fake"]`,
			"name_regex": `"${alicloud_ims_oidc_provider.default.oidc_provider_name}"`,
		}),
	}

	ImsOidcProviderCheckInfo.dataSourceTestCheck(t, rand, idsConf, allConf)
}

var existImsOidcProviderMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"providers.#":                     "1",
		"providers.0.description":         CHECKSET,
		"providers.0.issuer_url":          CHECKSET,
		"providers.0.fingerprints.#":      "1",
		"providers.0.create_time":         CHECKSET,
		"providers.0.update_time":         CHECKSET,
		"providers.0.issuance_limit_time": "12",
		"providers.0.oidc_provider_name":  CHECKSET,
		"providers.0.arn":                 CHECKSET,
		"providers.0.client_ids.#":        "2",
	}
}

var fakeImsOidcProviderMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"providers.#": "0",
	}
}

var ImsOidcProviderCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_ims_oidc_providers.default",
	existMapFunc: existImsOidcProviderMapFunc,
	fakeMapFunc:  fakeImsOidcProviderMapFunc,
}

func testAccCheckAlicloudImsOidcProviderSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccImsOidcProvider%d"
}
variable "oidc_provider_name" {
  default = "amp-resource-example-oidc-provider"
}

resource "alicloud_ims_oidc_provider" "default" {
  description = var.oidc_provider_name
  issuer_url  = "https://oauth.aliyun.com"
  fingerprints = [
    "0BBFAB97059595E8D1EC48E89EB8657C0E5AAE71"
  ]
  issuance_limit_time = "12"
  oidc_provider_name  = "tfaccims81016"
  client_ids = [
    "123",
    "456"
  ]
}

data "alicloud_ims_oidc_providers" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
