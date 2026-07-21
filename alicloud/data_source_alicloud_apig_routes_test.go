// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudApigRouteDataSource(t *testing.T) {
	testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	rand := acctest.RandIntRange(1000000, 9999999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudApigRouteSourceConfig(rand, map[string]string{
			"ids":         `["${alicloud_apig_route.default.id}"]`,
			"http_api_id": `"${alicloud_apig_http_api.route_httpapi_arr.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudApigRouteSourceConfig(rand, map[string]string{
			"ids":         `["${alicloud_apig_route.default.id}_fake"]`,
			"http_api_id": `"${alicloud_apig_http_api.route_httpapi_arr.id}"`,
		}),
	}

	RouteNameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudApigRouteSourceConfig(rand, map[string]string{
			"ids":         `["${alicloud_apig_route.default.id}"]`,
			"route_name":  `"${var.name}"`,
			"http_api_id": `"${alicloud_apig_http_api.route_httpapi_arr.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudApigRouteSourceConfig(rand, map[string]string{
			"ids":         `["${alicloud_apig_route.default.id}_fake"]`,
			"route_name":  `"${var.name}_fake"`,
			"http_api_id": `"${alicloud_apig_http_api.route_httpapi_arr.id}"`,
		}),
	}
	HttpApiIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudApigRouteSourceConfig(rand, map[string]string{
			"ids":         `["${alicloud_apig_route.default.id}"]`,
			"http_api_id": `"${alicloud_apig_http_api.route_httpapi_arr.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudApigRouteSourceConfig(rand, map[string]string{
			"ids":         `["${alicloud_apig_route.default.id}_fake"]`,
			"http_api_id": `"${alicloud_apig_http_api.route_httpapi_arr.id}"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudApigRouteSourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_apig_route.default.id}"]`,
			"route_name": `"${var.name}"`,

			"http_api_id": `"${alicloud_apig_http_api.route_httpapi_arr.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudApigRouteSourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_apig_route.default.id}_fake"]`,
			"route_name": `"${var.name}_fake"`,

			"http_api_id": `"${alicloud_apig_http_api.route_httpapi_arr.id}"`,
		}),
	}

	ApigRouteCheckInfo.dataSourceTestCheck(t, rand, idsConf, RouteNameConf, HttpApiIdConf, allConf)
}

var existApigRouteMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"routes.#":                    "1",
		"routes.0.status":             CHECKSET,
		"routes.0.description":        CHECKSET,
		"routes.0.create_time":        CHECKSET,
		"routes.0.domain_infos.#":     CHECKSET,
		"routes.0.route_id":           CHECKSET,
		"routes.0.match.#":            CHECKSET,
		"routes.0.gateway_status.%":   CHECKSET,
		"routes.0.backend.#":          CHECKSET,
		"routes.0.environment_info.#": CHECKSET,
		"routes.0.update_time":        CHECKSET,
	}
}

var fakeApigRouteMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"routes.#": "0",
	}
}

var ApigRouteCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_apig_routes.default",
	existMapFunc: existApigRouteMapFunc,
	fakeMapFunc:  fakeApigRouteMapFunc,
}

func testAccCheckAlicloudApigRouteSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testaccapigroute%d"
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_apig_gateway" "route_gateway_ds" {
  network_access_config {
    type = "Internet"
  }
  zone_config {
    select_option = "Auto"
  }
  vpc {
    vpc_id = data.alicloud_vpcs.default.ids.0
  }
  vswitch {
    vswitch_id = data.alicloud_vswitches.default.ids.0
  }
  gateway_type = "API"
  payment_type = "PayAsYouGo"
  gateway_name = "${var.name}gw"
  spec         = "apigw.small.x1"
  log_config {
    sls {
      enable = true
    }
  }
}

resource "alicloud_apig_service" "route_svc_ds_1" {
  service_name = "${var.name}svc1"
  source_type  = "DNS"
  gateway_id   = alicloud_apig_gateway.route_gateway_ds.id
  addresses    = ["httpbin.org:80"]
}

resource "alicloud_apig_service" "route_svc_ds_2" {
  service_name = "${var.name}svc2"
  source_type  = "DNS"
  gateway_id   = alicloud_apig_gateway.route_gateway_ds.id
  addresses    = ["httpbin.org:80"]
}

resource "alicloud_apig_http_api" "route_httpapi_arr" {
  http_api_name = "${var.name}api"
  protocols     = ["HTTP"]
  type          = "Http"
  description   = "array test httpapi"
  base_path     = "/${var.name}"
}



resource "alicloud_apig_route" "default" {
  backend {
    services {
      port       = "8080"
      protocol   = "HTTP"
      weight     = "70"
      service_id = alicloud_apig_service.route_svc_ds_1.id
    }
    services {
      port       = "8081"
      protocol   = "HTTP"
      weight     = "30"
      service_id = alicloud_apig_service.route_svc_ds_2.id
    }
    scene = "MultiServiceByRatio"
  }
  description = "array route description"
  route_name  = "${var.name}"
  http_api_id = alicloud_apig_http_api.route_httpapi_arr.id
  environment_info {
    environment_id = alicloud_apig_gateway.route_gateway_ds.environments.0.environment_id
  }
  match {
    path {
      type  = "Prefix"
      value = "/arr-path"
    }
    headers {
      type  = "Exact"
      value = "v1"
      name  = "h1"
    }
    headers {
      type  = "Prefix"
      value = "v2"
      name  = "h2"
    }
    query_params {
      type  = "Exact"
      value = "v1"
      name  = "q1"
    }
    query_params {
      type  = "Prefix"
      value = "v2"
      name  = "q2"
    }
    methods         = ["GET", "POST", "PUT"]
    ignore_uri_case = false
  }
}

data "alicloud_apig_routes" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
