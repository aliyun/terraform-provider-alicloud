package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudApiGatewayPluginsDataSource(t *testing.T) {
	rand := acctest.RandInt()
	checkoutSupportedRegions(t, true, connectivity.ApiGatewaySupportRegions)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudApiGatewayPluginsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_api_gateway_plugin.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudApiGatewayPluginsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_api_gateway_plugin.default.id}_fake"]`,
		}),
	}
	pluginNameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudApiGatewayPluginsDataSourceName(rand, map[string]string{
			"ids":         `["${alicloud_api_gateway_plugin.default.id}"]`,
			"plugin_name": `"${alicloud_api_gateway_plugin.default.plugin_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudApiGatewayPluginsDataSourceName(rand, map[string]string{
			"ids":         `["${alicloud_api_gateway_plugin.default.id}"]`,
			"plugin_name": `"${alicloud_api_gateway_plugin.default.plugin_name}_fake"`,
		}),
	}
	pluginTypeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudApiGatewayPluginsDataSourceName(rand, map[string]string{
			"ids":         `["${alicloud_api_gateway_plugin.default.id}"]`,
			"plugin_type": `"${alicloud_api_gateway_plugin.default.plugin_type}"`,
		}),
		fakeConfig: testAccCheckAlicloudApiGatewayPluginsDataSourceName(rand, map[string]string{
			"ids":         `["${alicloud_api_gateway_plugin.default.id}"]`,
			"plugin_type": `"backendSignature"`,
		}),
	}
	tagsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudApiGatewayPluginsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_api_gateway_plugin.default.id}"]`,
			"tags": `{
				"Created" = "TF"
				"For" = "Acceptance-test"
		}`,
		}),
		fakeConfig: testAccCheckAlicloudApiGatewayPluginsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_api_gateway_plugin.default.id}"]`,
			"tags": `{
				"Created" = "TF-fake"
				"For" = "Acceptance-test"
			}`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudApiGatewayPluginsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_api_gateway_plugin.default.plugin_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudApiGatewayPluginsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_api_gateway_plugin.default.plugin_name}_fake"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudApiGatewayPluginsDataSourceName(rand, map[string]string{
			"ids":         `["${alicloud_api_gateway_plugin.default.id}"]`,
			"name_regex":  `"${alicloud_api_gateway_plugin.default.plugin_name}"`,
			"plugin_name": `"${alicloud_api_gateway_plugin.default.plugin_name}"`,
			"plugin_type": `"${alicloud_api_gateway_plugin.default.plugin_type}"`,
			"tags": `{
				"Created" = "TF"
				"For" = "Acceptance-test"
		}`,
		}),
		fakeConfig: testAccCheckAlicloudApiGatewayPluginsDataSourceName(rand, map[string]string{
			"ids":         `["${alicloud_api_gateway_plugin.default.id}_fake"]`,
			"name_regex":  `"${alicloud_api_gateway_plugin.default.plugin_name}_fake"`,
			"plugin_name": `"${alicloud_api_gateway_plugin.default.plugin_name}_fake"`,
			"plugin_type": `"backendSignature"`,
			"tags": `{
				"Created" = "TF-fake"
				"For" = "Acceptance-test"
			}`,
		}),
	}
	var existAlicloudApiGatewayPluginsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                   "1",
			"names.#":                 "1",
			"plugins.#":               "1",
			"plugins.0.description":   fmt.Sprintf("tf_testAccPlugin_%d", rand),
			"plugins.0.plugin_data":   "{\"allowOrigins\": \"api.foo.com\",\"allowMethods\": \"GET,POST,PUT,DELETE,HEAD,OPTIONS,PATCH\",\"allowHeaders\": \"Authorization,Accept,Accept-Ranges,Cache-Control,Range,Date,Content-Type,Content-Length,Content-MD5,User-Agent,X-Ca-Signature,X-Ca-Signature-Headers,X-Ca-Signature-Method,X-Ca-Key,X-Ca-Timestamp,X-Ca-Nonce,X-Ca-Stage,X-Ca-Request-Mode,x-ca-deviceid\",\"exposeHeaders\": \"Content-MD5,Server,Date,Latency,X-Ca-Request-Id,X-Ca-Error-Code,X-Ca-Error-Message\",\"maxAge\": 172800,\"allowCredentials\": true}",
			"plugins.0.plugin_name":   fmt.Sprintf("tf_testAccPlugin_%d", rand),
			"plugins.0.plugin_type":   "cors",
			"plugins.0.tags.%":        "2",
			"plugins.0.tags.Created":  "TF",
			"plugins.0.tags.For":      "Acceptance-test",
			"plugins.0.id":            CHECKSET,
			"plugins.0.plugin_id":     CHECKSET,
			"plugins.0.create_time":   CHECKSET,
			"plugins.0.modified_time": CHECKSET,
		}
	}
	var fakeAlicloudApiGatewayPluginsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudApiGatewayPluginsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_api_gateway_plugins.default",
		existMapFunc: existAlicloudApiGatewayPluginsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudApiGatewayPluginsDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}
	alicloudApiGatewayPluginsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, pluginNameConf, pluginTypeConf, tagsConf, nameRegexConf, allConf)
}
func testAccCheckAlicloudApiGatewayPluginsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf_testAccPlugin_%d"
}

resource "alicloud_api_gateway_plugin" "default" {
	description = var.name
	plugin_name = var.name
	plugin_data = "{\"allowOrigins\": \"api.foo.com\",\"allowMethods\": \"GET,POST,PUT,DELETE,HEAD,OPTIONS,PATCH\",\"allowHeaders\": \"Authorization,Accept,Accept-Ranges,Cache-Control,Range,Date,Content-Type,Content-Length,Content-MD5,User-Agent,X-Ca-Signature,X-Ca-Signature-Headers,X-Ca-Signature-Method,X-Ca-Key,X-Ca-Timestamp,X-Ca-Nonce,X-Ca-Stage,X-Ca-Request-Mode,x-ca-deviceid\",\"exposeHeaders\": \"Content-MD5,Server,Date,Latency,X-Ca-Request-Id,X-Ca-Error-Code,X-Ca-Error-Message\",\"maxAge\": 172800,\"allowCredentials\": true}"
	plugin_type = "cors"
	tags = {
		Created = "TF"
		For = "Acceptance-test"
	}
}

data "alicloud_api_gateway_plugins" "default" {	
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
