package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAliCloudApiGatewayApisDataSource(t *testing.T) {
	rand := acctest.RandIntRange(100000, 999999)
	checkoutSupportedRegions(t, true, connectivity.TestSalveRegions)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudApiGatewayApisDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_api_gateway_api.default.id}"]`,
		}),
		fakeConfig: testAccCheckAliCloudApiGatewayApisDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_api_gateway_api.default.id}_fake"]`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudApiGatewayApisDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_api_gateway_api.default.name}"`,
		}),
		fakeConfig: testAccCheckAliCloudApiGatewayApisDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_api_gateway_api.default.name}_fake"`,
		}),
	}
	groupIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudApiGatewayApisDataSourceName(rand, map[string]string{
			"group_id": `"${alicloud_api_gateway_api.default.group_id}"`,
		}),
		fakeConfig: testAccCheckAliCloudApiGatewayApisDataSourceName(rand, map[string]string{
			"group_id":   `"${alicloud_api_gateway_api.default.group_id}"`,
			"name_regex": `"${alicloud_api_gateway_api.default.name}_fake"`,
		}),
	}
	apiIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudApiGatewayApisDataSourceName(rand, map[string]string{
			"api_id": `"${alicloud_api_gateway_api.default.api_id}"`,
		}),
		fakeConfig: testAccCheckAliCloudApiGatewayApisDataSourceName(rand, map[string]string{
			"api_id":     `"${alicloud_api_gateway_api.default.api_id}"`,
			"name_regex": `"${alicloud_api_gateway_api.default.name}_fake"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudApiGatewayApisDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_api_gateway_api.default.id}"]`,
			"name_regex": `"${alicloud_api_gateway_api.default.name}"`,
			"group_id":   `"${alicloud_api_gateway_api.default.group_id}"`,
			"api_id":     `"${alicloud_api_gateway_api.default.api_id}"`,
		}),
		fakeConfig: testAccCheckAliCloudApiGatewayApisDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_api_gateway_api.default.id}_fake"]`,
			"name_regex": `"${alicloud_api_gateway_api.default.name}_fake"`,
			"group_id":   `"${alicloud_api_gateway_api.default.group_id}"`,
			"api_id":     `"${alicloud_api_gateway_api.default.api_id}"`,
		}),
	}

	var existAliCloudApiGatewayApisDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":              "1",
			"names.#":            "1",
			"apis.#":             "1",
			"apis.0.id":          CHECKSET,
			"apis.0.group_id":    CHECKSET,
			"apis.0.api_id":      CHECKSET,
			"apis.0.name":        CHECKSET,
			"apis.0.description": CHECKSET,
			"apis.0.group_name":  CHECKSET,
			"apis.0.region_id":   CHECKSET,
		}
	}

	var fakeAliCloudApiGatewayApisDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
			"apis.#":  "0",
		}
	}

	var alicloudCenTransitRouterVpcAttachmentsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_api_gateway_apis.default",
		existMapFunc: existAliCloudApiGatewayApisDataSourceNameMapFunc,
		fakeMapFunc:  fakeAliCloudApiGatewayApisDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}

	alicloudCenTransitRouterVpcAttachmentsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, groupIdConf, apiIdConf, allConf)
}

func testAccCheckAliCloudApiGatewayApisDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
	variable "name" {
  		default = "tf-testAcc-ApiGatewayApi-%d"
	}

	resource "alicloud_api_gateway_group" "default" {
  		name        = var.name
  		description = var.name
	}

	resource "alicloud_api_gateway_api" "default" {
  		group_id     = alicloud_api_gateway_group.default.id
  		name         = var.name
  		description  = var.name
  		auth_type    = "APP"
  		service_type = "HTTP"
  		request_config {
    		protocol = "HTTP"
    		method   = "GET"
    		path     = "/test/path"
    		mode     = "MAPPING"
  		}
  		http_service_config {
    		address   = "http://apigateway-backend.alicloudapi.com:8080"
    		method    = "GET"
    		path      = "/web/cloudapi"
    		timeout   = 20
    		aone_name = "cloudapi-openapi"
  		}
  		request_parameters {
    		name         = var.name
    		type         = "STRING"
    		required     = "OPTIONAL"
    		in           = "QUERY"
    		in_service   = "QUERY"
    		name_service = var.name
  		}
	}

	data "alicloud_api_gateway_apis" "default" {
		%s
	}
`, rand, strings.Join(pairs, " \n "))
	return config
}
