// Package ims. This file is generated automatically. Please do not modify it manually, thank you!
package ims_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudImsOidcProviderDataSource(t *testing.T) {
	acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	rand := tfacctest.RandIntRange(1000000, 9999999)
	// IMS validates Fingerprints against the issuer's real HTTPS CA
	// certificate, so the dependent resource must use the live issuer root CA
	// SHA-1 (see imsOidcProviderFingerprint), not the synthetic doc-example
	// value.
	fp := imsOidcProviderFingerprint(t)

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudImsOidcProviderSourceConfig(rand, fp, map[string]string{
			"ids": `["${alicloud_ims_oidc_provider.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudImsOidcProviderSourceConfig(rand, fp, map[string]string{
			"ids": `["${alicloud_ims_oidc_provider.default.id}_fake"]`,
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudImsOidcProviderSourceConfig(rand, fp, map[string]string{
			"ids":        `["${alicloud_ims_oidc_provider.default.id}"]`,
			"name_regex": `"${alicloud_ims_oidc_provider.default.oidc_provider_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudImsOidcProviderSourceConfig(rand, fp, map[string]string{
			"ids":        `["${alicloud_ims_oidc_provider.default.id}_fake"]`,
			"name_regex": `"${alicloud_ims_oidc_provider.default.oidc_provider_name}"`,
		}),
	}

	ImsOidcProviderCheckInfo.DataSourceTestCheck(t, rand, idsConf, allConf)
}

var existImsOidcProviderMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"providers.#":                     "1",
		"providers.0.description":         acctest.CHECKSET,
		"providers.0.issuer_url":          acctest.CHECKSET,
		"providers.0.fingerprints.#":      "1",
		"providers.0.create_time":         acctest.CHECKSET,
		"providers.0.update_time":         acctest.CHECKSET,
		"providers.0.issuance_limit_time": "12",
		"providers.0.oidc_provider_name":  acctest.CHECKSET,
		"providers.0.arn":                 acctest.CHECKSET,
		"providers.0.client_ids.#":        "2",
	}
}

var fakeImsOidcProviderMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"providers.#": "0",
	}
}

var ImsOidcProviderCheckInfo = acctest.DataSourceAttr{
	ResourceId:   "data.alicloud_ims_oidc_providers.default",
	ExistMapFunc: existImsOidcProviderMapFunc,
	FakeMapFunc:  fakeImsOidcProviderMapFunc,
}

func testAccCheckAlicloudImsOidcProviderSourceConfig(rand int, fingerprint string, attrMap map[string]string) string {
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
    "%s"
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
`, rand, fingerprint, strings.Join(pairs, "\n   "))
	return config
}
