package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudServiceMeshExtensionProvidersDataSource(t *testing.T) {
	rand := acctest.RandInt()
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudServiceMeshExtensionProvidersDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_service_mesh_extension_provider.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudServiceMeshExtensionProvidersDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_service_mesh_extension_provider.default.id}_fake"]`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudServiceMeshExtensionProvidersDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_service_mesh_extension_provider.default.extension_provider_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudServiceMeshExtensionProvidersDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_service_mesh_extension_provider.default.extension_provider_name}_fake"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudServiceMeshExtensionProvidersDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_service_mesh_extension_provider.default.id}"]`,
			"name_regex": `"${alicloud_service_mesh_extension_provider.default.extension_provider_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudServiceMeshExtensionProvidersDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_service_mesh_extension_provider.default.id}_fake"]`,
			"name_regex": `"${alicloud_service_mesh_extension_provider.default.extension_provider_name}_fake"`,
		}),
	}
	var existAlicloudServiceMeshExtensionProvidersDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                               "1",
			"names.#":                             "1",
			"providers.#":                         "1",
			"providers.0.id":                      CHECKSET,
			"providers.0.service_mesh_id":         CHECKSET,
			"providers.0.extension_provider_name": "httpextauth-tf-example",
			"providers.0.type":                    "httpextauth",
			"providers.0.config":                  "{\"headersToDownstreamOnDeny\":[\"content-type\",\"set-cookie\"],\"headersToUpstreamOnAllow\":[\"authorization\",\"cookie\",\"path\",\"x-auth-request-access-token\",\"x-forwarded-access-token\"],\"includeRequestHeadersInCheck\":[\"cookie\",\"x-forward-access-token\"],\"oidc\":{\"clientID\":\"qweqweqwewqeqwe\",\"clientSecret\":\"asdasdasdasdsadas\",\"cookieExpire\":\"1000\",\"cookieRefresh\":\"500\",\"cookieSecret\":\"scxzcxzcxzcxzcxz\",\"issuerURI\":\"qweqwewqeqweqweqwe\",\"redirectDomain\":\"www.baidu.com\",\"redirectProtocol\":\"http\",\"scopes\":[\"profile\"]},\"port\":4180,\"service\":\"asm-oauth2proxy-httpextauth-tf-example.istio-system.svc.cluster.local\",\"timeout\":\"10s\"}",
		}
	}
	var fakeAlicloudServiceMeshExtensionProvidersDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":       "0",
			"names.#":     "0",
			"providers.#": "0",
		}
	}
	var alicloudServiceMeshExtensionProvidersCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_service_mesh_extension_providers.default",
		existMapFunc: existAlicloudServiceMeshExtensionProvidersDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudServiceMeshExtensionProvidersDataSourceNameMapFunc,
	}
	alicloudServiceMeshExtensionProvidersCheckInfo.dataSourceTestCheck(t, rand, idsConf, nameRegexConf, allConf)
}

func testAccCheckAlicloudServiceMeshExtensionProvidersDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
	variable "name" {
  		default = "tf-testAcc-%d"
	}

	data "alicloud_zones" "default" {
  		available_resource_creation = "VSwitch"
	}

	data "alicloud_vpcs" "default" {
  		name_regex = "default-NODELETING"
	}

	resource "alicloud_vpc" "default" {
  		count    = length(data.alicloud_vpcs.default.ids) > 0 ? 0 : 1
	}

	data "alicloud_vswitches" "default" {
  		vpc_id = length(data.alicloud_vpcs.default.ids) > 0 ? data.alicloud_vpcs.default.ids[0] : alicloud_vpc.default[0].id
	}

	resource "alicloud_vswitch" "default" {
  		count        = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
  		vpc_id       = length(data.alicloud_vpcs.default.ids) > 0 ? data.alicloud_vpcs.default.ids[0] : alicloud_vpc.default[0].id
  		cidr_block   = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 2)
  		zone_id      = data.alicloud_zones.default.zones.0.id
	}

	resource "alicloud_service_mesh_service_mesh" "default" {
  		network {
    		vpc_id        = length(data.alicloud_vpcs.default.ids) > 0 ? data.alicloud_vpcs.default.ids[0] : alicloud_vpc.default[0].id
    		vswitche_list = [length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : alicloud_vswitch.default[0].id]
  		}
	}

	resource "alicloud_service_mesh_extension_provider" "default" {
  		service_mesh_id         = alicloud_service_mesh_service_mesh.default.id
  		extension_provider_name = "httpextauth-tf-example"
  		type                    = "httpextauth"
  		config                  = "{\"headersToDownstreamOnDeny\":[\"content-type\",\"set-cookie\"],\"headersToUpstreamOnAllow\":[\"authorization\",\"cookie\",\"path\",\"x-auth-request-access-token\",\"x-forwarded-access-token\"],\"includeRequestHeadersInCheck\":[\"cookie\",\"x-forward-access-token\"],\"oidc\":{\"clientID\":\"qweqweqwewqeqwe\",\"clientSecret\":\"asdasdasdasdsadas\",\"cookieExpire\":\"1000\",\"cookieRefresh\":\"500\",\"cookieSecret\":\"scxzcxzcxzcxzcxz\",\"issuerURI\":\"qweqwewqeqweqweqwe\",\"redirectDomain\":\"www.baidu.com\",\"redirectProtocol\":\"http\",\"scopes\":[\"profile\"]},\"port\":4180,\"service\":\"asm-oauth2proxy-httpextauth-tf-example.istio-system.svc.cluster.local\",\"timeout\":\"10s\"}"
	}

	data "alicloud_service_mesh_extension_providers" "default" {
  		service_mesh_id = alicloud_service_mesh_extension_provider.default.service_mesh_id
  		type            = alicloud_service_mesh_extension_provider.default.type
		%s
	}
`, rand, strings.Join(pairs, " \n "))
	return config
}
