// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Apig Route. >>> Resource test cases, automatically generated.
// Case resource_Route_arr 12909
func TestAccAliCloudApigRoute_basic12909(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_apig_route.default"
	ra := resourceAttrInit(resourceId, AlicloudApigRouteMap12909)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApigServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeApigRoute")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccapig%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudApigRouteBasicDependence12909)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"backend": []map[string]interface{}{
						{
							"services": []map[string]interface{}{
								{
									"version":    "v1.0",
									"port":       "8080",
									"protocol":   "HTTP",
									"weight":     "70",
									"service_id": "${alicloud_apig_service.route_svc_arr_1.id}",
								},
								{
									"port":       "8081",
									"protocol":   "HTTP",
									"weight":     "30",
									"service_id": "${alicloud_apig_service.route_svc_arr_2.id}",
								},
							},
							"scene": "MultiServiceByRatio",
						},
					},
					"description": "array route description",
					"route_name":  name,
					"http_api_id": "${alicloud_apig_http_api.route_httpapi_arr.id}",
					"environment_info": []map[string]interface{}{
						{
							"environment_id": "${alicloud_apig_gateway.route_gateway_arr.environments.0.environment_id}",
						},
					},
					"match": []map[string]interface{}{
						{
							"path": []map[string]interface{}{
								{
									"type":  "Prefix",
									"value": "/arr-path",
								},
							},
							"headers": []map[string]interface{}{
								{
									"type":  "Exact",
									"value": "v1",
									"name":  "h1",
								},
								{
									"type":  "Prefix",
									"value": "v2",
									"name":  "h2",
								},
							},
							"query_params": []map[string]interface{}{
								{
									"type":  "Exact",
									"value": "v1",
									"name":  "q1",
								},
								{
									"type":  "Prefix",
									"value": "v2",
									"name":  "q2",
								},
							},
							"methods": []string{
								"GET", "POST", "PUT"},
							"ignore_uri_case": "false",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "array route description",
						"route_name":  name,
						"http_api_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "arr updated",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "arr updated",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"match": []map[string]interface{}{
						{
							"path": []map[string]interface{}{
								{
									"type":  "Exact",
									"value": "/arr-updated",
								},
							},
							"headers": []map[string]interface{}{
								{
									"type":  "Exact",
									"value": "v3",
									"name":  "h3",
								},
							},
							"query_params": []map[string]interface{}{
								{
									"type":  "Exact",
									"value": "v3",
									"name":  "q3",
								},
							},
							"methods": []string{
								"GET"},
							"ignore_uri_case": "true",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"backend": []map[string]interface{}{
						{
							"services": []map[string]interface{}{
								{
									"weight":     "100",
									"service_id": "${alicloud_apig_service.route_svc_arr_3.id}",
								},
							},
							"scene": "SingleService",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"domain_ids"},
			},
		},
	})
}

var AlicloudApigRouteMap12909 = map[string]string{
	"status":           CHECKSET,
	"create_time":      CHECKSET,
	"route_id":         CHECKSET,
	"gateway_status.%": CHECKSET,
	"update_time":      CHECKSET,
}

func AlicloudApigRouteBasicDependence12909(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_apig_gateway" "route_gateway_arr" {
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

resource "alicloud_apig_service" "route_svc_arr_1" {
  service_name = "${var.name}svc1"
  source_type  = "DNS"
  gateway_id   = alicloud_apig_gateway.route_gateway_arr.id
  addresses    = ["httpbin.org:80"]
}

resource "alicloud_apig_service" "route_svc_arr_2" {
  service_name = "${var.name}svc2"
  source_type  = "DNS"
  gateway_id   = alicloud_apig_gateway.route_gateway_arr.id
  addresses    = ["httpbin.org:80"]
}

resource "alicloud_apig_service" "route_svc_arr_3" {
  service_name = "${var.name}svc3"
  source_type  = "DNS"
  gateway_id   = alicloud_apig_gateway.route_gateway_arr.id
  addresses    = ["httpbin.org:80"]
}

resource "alicloud_apig_http_api" "route_httpapi_arr" {
  http_api_name = "${var.name}api"
  protocols     = ["HTTP"]
  type          = "Rest"
  description   = "array test httpapi"
  base_path     = "/${var.name}"
}


`, name)
}

// Case resource_Route_test 12910
func TestAccAliCloudApigRoute_basic12910(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_apig_route.default"
	ra := resourceAttrInit(resourceId, AlicloudApigRouteMap12910)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApigServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeApigRoute")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccapig%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudApigRouteBasicDependence12910)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"backend": []map[string]interface{}{
						{
							"scene": "Mock",
						},
					},
					"description": "test route description",
					"environment_info": []map[string]interface{}{
						{
							"environment_id": "${alicloud_apig_gateway.route_gateway_pre.environments.0.environment_id}",
						},
					},
					"route_name":  name,
					"http_api_id": "${alicloud_apig_http_api.route_httpapi_pre.id}",
					"match": []map[string]interface{}{
						{
							"path": []map[string]interface{}{
								{
									"type":  "Prefix",
									"value": "/test-path",
								},
							},
							"methods": []string{
								"GET"},
							"ignore_uri_case": "false",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "test route description",
						"route_name":  name,
						"http_api_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "updated route description",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "updated route description",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"match": []map[string]interface{}{
						{
							"path": []map[string]interface{}{
								{
									"type":  "Exact",
									"value": "/updated-path",
								},
							},
							"methods": []string{
								"GET", "POST"},
							"ignore_uri_case": "true",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"domain_ids"},
			},
		},
	})
}

var AlicloudApigRouteMap12910 = map[string]string{
	"status":           CHECKSET,
	"create_time":      CHECKSET,
	"route_id":         CHECKSET,
	"gateway_status.%": CHECKSET,
	"update_time":      CHECKSET,
}

func AlicloudApigRouteBasicDependence12910(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_apig_gateway" "route_gateway_pre" {
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
  gateway_name = format("%%sgw", var.name)
  spec         = "apigw.small.x1"
  log_config {
    sls {
      enable = true
    }
  }
}

resource "alicloud_apig_domain" "route_domain_pre" {
  domain_name = format("%%s.example.com", var.name)
  force_https = false
  protocol    = "HTTP"
}

resource "alicloud_apig_http_api" "route_httpapi_pre" {
  http_api_name = format("%%sapi", var.name)
  protocols     = ["HTTP"]
  type          = "Rest"
  description   = "test description for cspec route"
  base_path     = format("/%%s", var.name)
}


`, name)
}

// Case resource_Route_svc 12911
func TestAccAliCloudApigRoute_basic12911(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_apig_route.default"
	ra := resourceAttrInit(resourceId, AlicloudApigRouteMap12911)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApigServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeApigRoute")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccapig%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudApigRouteBasicDependence12911)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"backend": []map[string]interface{}{
						{
							"services": []map[string]interface{}{
								{
									"port":       "80",
									"protocol":   "HTTP",
									"weight":     "100",
									"service_id": "${alicloud_apig_service.route_svc_svc_1.id}",
								},
							},
							"scene": "SingleService",
						},
					},
					"description": "svc route description",
					"route_name":  name,
					"http_api_id": "${alicloud_apig_http_api.route_httpapi_svc.id}",
					"environment_info": []map[string]interface{}{
						{
							"environment_id": "${alicloud_apig_gateway.route_gateway_svc.environments.0.environment_id}",
						},
					},
					"domain_ids": []string{
						"${alicloud_apig_domain.route_domain_svc.id}"},
					"match": []map[string]interface{}{
						{
							"path": []map[string]interface{}{
								{
									"type":  "Prefix",
									"value": "/svc-path",
								},
							},
							"headers": []map[string]interface{}{
								{
									"type":  "Exact",
									"value": "custom-val",
									"name":  "x-custom",
								},
							},
							"query_params": []map[string]interface{}{
								{
									"type":  "Exact",
									"value": "active",
									"name":  "filter",
								},
							},
							"methods": []string{
								"GET", "POST"},
							"ignore_uri_case": "false",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":  "svc route description",
						"route_name":   name,
						"http_api_id":  CHECKSET,
						"domain_ids.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "svc updated description",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "svc updated description",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"backend": []map[string]interface{}{
						{
							"services": []map[string]interface{}{
								{
									"service_id": "${alicloud_apig_service.route_svc_svc_2.id}",
								},
							},
							"scene": "SingleService",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"match": []map[string]interface{}{
						{
							"path": []map[string]interface{}{
								{
									"type":  "Exact",
									"value": "/svc-updated",
								},
							},
							"headers": []map[string]interface{}{
								{
									"type":  "Prefix",
									"value": "updated",
									"name":  "x-updated",
								},
							},
							"query_params": []map[string]interface{}{
								{
									"type":  "Prefix",
									"value": "1",
									"name":  "page",
								},
							},
							"methods": []string{
								"GET"},
							"ignore_uri_case": "true",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"domain_ids"},
			},
		},
	})
}

var AlicloudApigRouteMap12911 = map[string]string{
	"status":           CHECKSET,
	"create_time":      CHECKSET,
	"route_id":         CHECKSET,
	"gateway_status.%": CHECKSET,
	"update_time":      CHECKSET,
}

func AlicloudApigRouteBasicDependence12911(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_apig_gateway" "route_gateway_svc" {
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

resource "alicloud_apig_domain" "route_domain_svc" {
  domain_name = "${var.name}.example.com"
  force_https = false
  protocol    = "HTTP"
}

resource "alicloud_apig_service" "route_svc_svc_1" {
  service_name = "${var.name}svc1"
  source_type  = "DNS"
  gateway_id   = alicloud_apig_gateway.route_gateway_svc.id
  addresses    = ["httpbin.org:80"]
}

resource "alicloud_apig_service" "route_svc_svc_2" {
  service_name = "${var.name}svc2"
  source_type  = "DNS"
  gateway_id   = alicloud_apig_gateway.route_gateway_svc.id
  addresses    = ["httpbin.org:80"]
}

resource "alicloud_apig_http_api" "route_httpapi_svc" {
  http_api_name = "${var.name}api"
  protocols     = ["HTTP"]
  type          = "Rest"
  description   = "svc test httpapi"
  base_path     = "/${var.name}"
}


`, name)
}

// Case resource_Route_alt2 12912
func TestAccAliCloudApigRoute_basic12912(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_apig_route.default"
	ra := resourceAttrInit(resourceId, AlicloudApigRouteMap12912)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApigServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeApigRoute")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccapig%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudApigRouteBasicDependence12912)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"backend": []map[string]interface{}{
						{
							"scene": "Mock",
						},
					},
					"description": "alt2 route description",
					"route_name":  name,
					"http_api_id": "${alicloud_apig_http_api.route_httpapi_alt2.id}",
					"environment_info": []map[string]interface{}{
						{
							"environment_id": "${alicloud_apig_gateway.route_gateway_alt2.environments.0.environment_id}",
						},
					},
					"match": []map[string]interface{}{
						{
							"path": []map[string]interface{}{
								{
									"type":  "Prefix",
									"value": "/alt2-path",
								},
							},
							"headers": []map[string]interface{}{
								{
									"type":  "Exact",
									"value": "test",
									"name":  "x-test",
								},
							},
							"query_params": []map[string]interface{}{
								{
									"type":  "Exact",
									"value": "test",
									"name":  "q",
								},
							},
							"methods": []string{
								"GET", "POST"},
							"ignore_uri_case": "false",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "alt2 route description",
						"route_name":  name,
						"http_api_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "alt2 updated",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "alt2 updated",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"match": []map[string]interface{}{
						{
							"path": []map[string]interface{}{
								{
									"type":  "Exact",
									"value": "/alt2-updated",
								},
							},
							"methods": []string{
								"GET"},
							"ignore_uri_case": "true",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"domain_ids"},
			},
		},
	})
}

var AlicloudApigRouteMap12912 = map[string]string{
	"status":           CHECKSET,
	"create_time":      CHECKSET,
	"route_id":         CHECKSET,
	"gateway_status.%": CHECKSET,
	"update_time":      CHECKSET,
}

func AlicloudApigRouteBasicDependence12912(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_apig_gateway" "route_gateway_alt2" {
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
  gateway_name = format("%%sgw", var.name)
  spec         = "apigw.small.x1"
  log_config {
    sls {
      enable = true
    }
  }
}

resource "alicloud_apig_http_api" "route_httpapi_alt2" {
  http_api_name = format("%%sapi", var.name)
  protocols     = ["HTTP"]
  type          = "Rest"
  description   = "alt2 test httpapi"
  base_path     = format("/%%s", var.name)
}


`, name)
}

// Case resource_Route_arr2 12913
func TestAccAliCloudApigRoute_basic12913(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_apig_route.default"
	ra := resourceAttrInit(resourceId, AlicloudApigRouteMap12913)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApigServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeApigRoute")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccapig%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudApigRouteBasicDependence12913)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"backend": []map[string]interface{}{
						{
							"services": []map[string]interface{}{
								{
									"port":       "80",
									"protocol":   "HTTP",
									"weight":     "100",
									"service_id": "${alicloud_apig_service.route_svc_arr2_1.id}",
								},
							},
							"scene": "SingleService",
						},
					},
					"description": "array2 route description",
					"route_name":  name,
					"http_api_id": "${alicloud_apig_http_api.route_httpapi_arr2.id}",
					"environment_info": []map[string]interface{}{
						{
							"environment_id": "${alicloud_apig_gateway.route_gateway_arr2.environments.0.environment_id}",
						},
					},
					"match": []map[string]interface{}{
						{
							"path": []map[string]interface{}{
								{
									"type":  "Prefix",
									"value": "/arr2-path",
								},
							},
							"headers": []map[string]interface{}{
								{
									"type":  "Exact",
									"value": "v1",
									"name":  "h1",
								},
								{
									"type":  "Prefix",
									"value": "v2",
									"name":  "h2",
								},
								{
									"type":  "Regex",
									"value": "v3",
									"name":  "h3",
								},
							},
							"query_params": []map[string]interface{}{
								{
									"type":  "Exact",
									"value": "v1",
									"name":  "q1",
								},
								{
									"type":  "Prefix",
									"value": "v2",
									"name":  "q2",
								},
								{
									"type":  "Regex",
									"value": "v3",
									"name":  "q3",
								},
							},
							"methods": []string{
								"GET", "POST", "PUT", "DELETE"},
							"ignore_uri_case": "false",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "array2 route description",
						"route_name":  name,
						"http_api_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "arr2 updated",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "arr2 updated",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"match": []map[string]interface{}{
						{
							"path": []map[string]interface{}{
								{
									"type":  "Exact",
									"value": "/arr2-updated",
								},
							},
							"headers": []map[string]interface{}{
								{
									"type":  "Exact",
									"value": "v-new",
									"name":  "h-new",
								},
							},
							"query_params": []map[string]interface{}{
								{
									"type":  "Exact",
									"value": "v-new",
									"name":  "q-new",
								},
							},
							"methods": []string{
								"GET"},
							"ignore_uri_case": "true",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"backend": []map[string]interface{}{
						{
							"services": []map[string]interface{}{
								{
									"weight":     "50",
									"service_id": "${alicloud_apig_service.route_svc_arr2_2.id}",
								},
								{
									"weight":     "50",
									"service_id": "${alicloud_apig_service.route_svc_arr2_3.id}",
								},
							},
							"scene": "MultiServiceByRatio",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"domain_ids"},
			},
		},
	})
}

var AlicloudApigRouteMap12913 = map[string]string{
	"status":           CHECKSET,
	"create_time":      CHECKSET,
	"route_id":         CHECKSET,
	"gateway_status.%": CHECKSET,
	"update_time":      CHECKSET,
}

func AlicloudApigRouteBasicDependence12913(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_apig_gateway" "route_gateway_arr2" {
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

resource "alicloud_apig_service" "route_svc_arr2_1" {
  service_name = "${var.name}svc1"
  source_type  = "DNS"
  gateway_id   = alicloud_apig_gateway.route_gateway_arr2.id
  addresses    = ["httpbin.org:80"]
}

resource "alicloud_apig_service" "route_svc_arr2_2" {
  service_name = "${var.name}svc2"
  source_type  = "DNS"
  gateway_id   = alicloud_apig_gateway.route_gateway_arr2.id
  addresses    = ["httpbin.org:80"]
}

resource "alicloud_apig_service" "route_svc_arr2_3" {
  service_name = "${var.name}svc3"
  source_type  = "DNS"
  gateway_id   = alicloud_apig_gateway.route_gateway_arr2.id
  addresses    = ["httpbin.org:80"]
}

resource "alicloud_apig_http_api" "route_httpapi_arr2" {
  http_api_name = "${var.name}api"
  protocols     = ["HTTP"]
  type          = "Rest"
  description   = "array2 test httpapi"
  base_path     = "/${var.name}"
}


`, name)
}

// Case resource_Route_alt3 12914
func TestAccAliCloudApigRoute_basic12914(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_apig_route.default"
	ra := resourceAttrInit(resourceId, AlicloudApigRouteMap12914)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApigServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeApigRoute")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccapig%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudApigRouteBasicDependence12914)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"backend": []map[string]interface{}{
						{
							"scene": "Mock",
						},
					},
					"description": "alt3 route description",
					"route_name":  name,
					"http_api_id": "${alicloud_apig_http_api.route_httpapi_alt3.id}",
					"environment_info": []map[string]interface{}{
						{
							"environment_id": "${alicloud_apig_gateway.route_gateway_alt3.environments.0.environment_id}",
						},
					},
					"match": []map[string]interface{}{
						{
							"path": []map[string]interface{}{
								{
									"type":  "Prefix",
									"value": "/alt3-path",
								},
							},
							"methods": []string{
								"GET"},
							"ignore_uri_case": "false",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "alt3 route description",
						"route_name":  name,
						"http_api_id": CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"domain_ids"},
			},
		},
	})
}

var AlicloudApigRouteMap12914 = map[string]string{
	"status":           CHECKSET,
	"create_time":      CHECKSET,
	"route_id":         CHECKSET,
	"gateway_status.%": CHECKSET,
	"update_time":      CHECKSET,
}

func AlicloudApigRouteBasicDependence12914(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_apig_gateway" "route_gateway_alt3" {
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
  gateway_name = format("%%sgw", var.name)
  spec         = "apigw.small.x1"
  log_config {
    sls {
      enable = true
    }
  }
}

resource "alicloud_apig_http_api" "route_httpapi_alt3" {
  http_api_name = format("%%sapi", var.name)
  protocols     = ["HTTP"]
  type          = "Rest"
  description   = "alt3 test httpapi"
  base_path     = format("/%%s", var.name)
}


`, name)
}

// Case test-route-singleservice 8931
func TestAccAliCloudApigRoute_basic8931(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_apig_route.default"
	ra := resourceAttrInit(resourceId, AlicloudApigRouteMap8931)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApigServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeApigRoute")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccapig%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudApigRouteBasicDependence8931)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"backend": []map[string]interface{}{
						{
							"scene": "SingleService",
							"services": []map[string]interface{}{
								{
									"service_id": "${alicloud_apig_service.defaultservice1.id}",
								},
							},
						},
					},
					"environment_info": []map[string]interface{}{
						{
							"environment_id": "${alicloud_apig_environment.defaultenvironment.id}",
						},
					},
					"route_name": name,
					"match": []map[string]interface{}{
						{
							"path": []map[string]interface{}{
								{
									"type":  "Exact",
									"value": "/",
								},
							},
							"methods": []string{
								"GET"},
							"headers": []map[string]interface{}{
								{
									"type":  "Exact",
									"name":  "k1",
									"value": "v1",
								},
							},
							"query_params": []map[string]interface{}{
								{
									"type":  "Exact",
									"name":  "p1",
									"value": "v1",
								},
							},
						},
					},
					"http_api_id": "${alicloud_apig_http_api.defaultapi.id}",
					"domain_ids": []string{
						"${alicloud_apig_domain.defaultdomain.id}"},
					"description": "test-route",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"route_name":   name,
						"http_api_id":  CHECKSET,
						"domain_ids.#": "1",
						"description":  "test-route",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"match": []map[string]interface{}{
						{
							"ignore_uri_case": "true",
							"path": []map[string]interface{}{
								{
									"type":  "Exact",
									"value": "/v1",
								},
							},
							"methods": []string{
								"POST", "GET", "PUT"},
							"headers": []map[string]interface{}{
								{
									"type":  "Prefix",
									"name":  "k1",
									"value": "v1",
								},
								{
									"type":  "Exact",
									"name":  "k2",
									"value": "v2",
								},
								{
									"type":  "Exact",
									"name":  "k3",
									"value": "v3",
								},
							},
							"query_params": []map[string]interface{}{
								{
									"type":  "Prefix",
									"name":  "p1",
									"value": "v1",
								},
								{
									"type":  "Exact",
									"name":  "p2",
									"value": "v2",
								},
								{
									"type":  "Exact",
									"name":  "p3",
									"value": "v3",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"backend": []map[string]interface{}{
						{
							"scene": "MultiServiceByRatio",
							"services": []map[string]interface{}{
								{
									"service_id": "${alicloud_apig_service.defaultservice1.id}",
									"weight":     "20",
								},
								{
									"service_id": "${alicloud_apig_service.defaultservice2.id}",
									"weight":     "80",
								},
							},
						},
					},
					"match": []map[string]interface{}{
						{
							"ignore_uri_case": "true",
							"path": []map[string]interface{}{
								{
									"type":  "Prefix",
									"value": "/v1",
								},
							},
							"headers": []map[string]interface{}{
								{
									"type":  "Exact",
									"name":  "k1",
									"value": "v1",
								},
								{
									"type":  "Exact",
									"name":  "k2",
									"value": "v2",
								},
							},
							"methods": []string{
								"GET", "POST"},
							"query_params": []map[string]interface{}{
								{
									"type":  "Exact",
									"name":  "p1",
									"value": "v1",
								},
								{
									"type":  "Exact",
									"name":  "p2",
									"value": "v2",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"backend": []map[string]interface{}{
						{
							"scene": "SingleService",
							"services": []map[string]interface{}{
								{
									"service_id": "${alicloud_apig_service.defaultservice1.id}",
								},
							},
						},
					},
					"match": []map[string]interface{}{
						{
							"path": []map[string]interface{}{
								{
									"type":  "Exact",
									"value": "/v2",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"domain_ids"},
			},
		},
	})
}

var AlicloudApigRouteMap8931 = map[string]string{
	"status":           CHECKSET,
	"create_time":      CHECKSET,
	"route_id":         CHECKSET,
	"gateway_status.%": CHECKSET,
	"update_time":      CHECKSET,
}

func AlicloudApigRouteBasicDependence8931(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "protocol" {
  default = "HTTP"
}

variable "dns_address" {
  default = "httpbin.org:80"
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_apig_gateway" "defaultgateway" {
  network_access_config {
    type = "Intranet"
  }
  vswitch {
    vswitch_id = data.alicloud_vswitches.default.ids.0
  }
  zone_config {
    select_option = "Auto"
  }
  vpc {
    vpc_id = data.alicloud_vpcs.default.ids.0
  }
  payment_type = "PayAsYouGo"
  gateway_name = "${var.name}gw"
  spec         = "apigw.small.x1"
}

resource "alicloud_apig_domain" "defaultdomain" {
  domain_name = "${var.name}.example.com"
  protocol    = "HTTP"
}

resource "alicloud_apig_environment" "defaultenvironment" {
  description      = "镇元测试环境"
  environment_name = "${var.name}env"
  gateway_id       = alicloud_apig_gateway.defaultgateway.id
}

resource "alicloud_apig_service" "defaultservice1" {
  service_name = "${var.name}svc1"
  source_type  = "DNS"
  gateway_id   = alicloud_apig_gateway.defaultgateway.id
  addresses    = ["${var.dns_address}"]
}

resource "alicloud_apig_service" "defaultservice2" {
  service_name = "${var.name}svc2"
  source_type  = "DNS"
  gateway_id   = alicloud_apig_gateway.defaultgateway.id
  addresses    = ["${var.dns_address}"]
}

resource "alicloud_apig_http_api" "defaultapi" {
  http_api_name = "${var.name}api"
  type          = "Http"
  description   = "test api"
  protocols     = ["${var.protocol}"]
}


`, name)
}

// Case test-route-multiservice 8980
func TestAccAliCloudApigRoute_basic8980(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_apig_route.default"
	ra := resourceAttrInit(resourceId, AlicloudApigRouteMap8980)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApigServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeApigRoute")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccapig%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudApigRouteBasicDependence8980)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"backend": []map[string]interface{}{
						{
							"scene": "MultiServiceByRatio",
							"services": []map[string]interface{}{
								{
									"service_id": "${alicloud_apig_service.defaultservice1.id}",
									"weight":     "20",
									"port":       "80",
									"protocol":   "HTTP",
								},
								{
									"weight":     "80",
									"service_id": "${alicloud_apig_service.defaultservice2.id}",
								},
							},
						},
					},
					"environment_info": []map[string]interface{}{
						{
							"environment_id": "${alicloud_apig_environment.defaultenvironment.id}",
						},
					},
					"route_name": name,
					"match": []map[string]interface{}{
						{
							"path": []map[string]interface{}{
								{
									"type":  "Exact",
									"value": "/",
								},
							},
							"methods": []string{
								"GET"},
							"headers": []map[string]interface{}{
								{
									"type":  "Exact",
									"name":  "k1",
									"value": "v1",
								},
							},
							"query_params": []map[string]interface{}{
								{
									"type":  "Exact",
									"name":  "p1",
									"value": "v1",
								},
							},
							"ignore_uri_case": "false",
						},
					},
					"http_api_id": "${alicloud_apig_http_api.defaultapi.id}",
					"domain_ids": []string{
						"${alicloud_apig_domain.defaultdomain.id}"},
					"description": "test-route",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"route_name":   name,
						"http_api_id":  CHECKSET,
						"domain_ids.#": "1",
						"description":  "test-route",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"backend": []map[string]interface{}{
						{
							"scene": "MultiServiceByRatio",
							"services": []map[string]interface{}{
								{
									"service_id": "${alicloud_apig_service.defaultservice1.id}",
									"weight":     "50",
								},
								{
									"service_id": "${alicloud_apig_service.defaultservice2.id}",
									"weight":     "20",
								},
								{
									"service_id": "${alicloud_apig_service.defaultservice3.id}",
									"weight":     "30",
								},
							},
						},
					},
					"match": []map[string]interface{}{
						{
							"ignore_uri_case": "true",
							"path": []map[string]interface{}{
								{
									"type":  "Exact",
									"value": "/v1",
								},
							},
							"methods": []string{
								"POST", "GET", "PUT"},
							"headers": []map[string]interface{}{
								{
									"type":  "Prefix",
									"name":  "k1",
									"value": "v1",
								},
							},
							"query_params": []map[string]interface{}{
								{
									"type":  "Prefix",
									"name":  "p1",
									"value": "v1",
								},
							},
						},
					},
					"description": "test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"backend": []map[string]interface{}{
						{
							"scene": "SingleService",
							"services": []map[string]interface{}{
								{
									"service_id": "${alicloud_apig_service.defaultservice1.id}",
								},
							},
						},
					},
					"match": []map[string]interface{}{
						{
							"path": []map[string]interface{}{
								{
									"type":  "Exact",
									"value": "/v2",
								},
							},
							"ignore_uri_case": "false",
							"methods":         []string{},
						},
					},
					"description": "test-route",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "test-route",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"domain_ids"},
			},
		},
	})
}

var AlicloudApigRouteMap8980 = map[string]string{
	"status":           CHECKSET,
	"create_time":      CHECKSET,
	"route_id":         CHECKSET,
	"gateway_status.%": CHECKSET,
	"update_time":      CHECKSET,
}

func AlicloudApigRouteBasicDependence8980(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "protocol" {
  default = "HTTP"
}

variable "dns_address" {
  default = "httpbin.org:80"
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_apig_gateway" "defaultgateway" {
  network_access_config {
    type = "Intranet"
  }
  vswitch {
    vswitch_id = data.alicloud_vswitches.default.ids.0
  }
  zone_config {
    select_option = "Auto"
  }
  vpc {
    vpc_id = data.alicloud_vpcs.default.ids.0
  }
  payment_type = "PayAsYouGo"
  gateway_name = "${var.name}gw"
  spec         = "apigw.small.x1"
}

resource "alicloud_apig_domain" "defaultdomain" {
  domain_name = "${var.name}.example.com"
  protocol    = "HTTP"
}

resource "alicloud_apig_environment" "defaultenvironment" {
  description      = "镇元测试环境"
  environment_name = "${var.name}env"
  gateway_id       = alicloud_apig_gateway.defaultgateway.id
}

resource "alicloud_apig_service" "defaultservice1" {
  service_name = "${var.name}svc1"
  source_type  = "DNS"
  gateway_id   = alicloud_apig_gateway.defaultgateway.id
  addresses    = ["${var.dns_address}"]
}

resource "alicloud_apig_service" "defaultservice2" {
  service_name = "${var.name}svc2"
  source_type  = "DNS"
  gateway_id   = alicloud_apig_gateway.defaultgateway.id
  addresses    = ["${var.dns_address}"]
}

resource "alicloud_apig_service" "defaultservice3" {
  service_name = "${var.name}svc3"
  source_type  = "DNS"
  gateway_id   = alicloud_apig_gateway.defaultgateway.id
  addresses    = ["${var.dns_address}"]
}

resource "alicloud_apig_http_api" "defaultapi" {
  http_api_name = "${var.name}api"
  type          = "Http"
  description   = "test api"
  protocols     = ["${var.protocol}"]
}


`, name)
}

// Case test-route-updatedomain 8981
func TestAccAliCloudApigRoute_basic8981(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_apig_route.default"
	ra := resourceAttrInit(resourceId, AlicloudApigRouteMap8981)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApigServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeApigRoute")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccapig%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudApigRouteBasicDependence8981)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"backend": []map[string]interface{}{
						{
							"scene": "MultiServiceByRatio",
							"services": []map[string]interface{}{
								{
									"service_id": "${alicloud_apig_service.defaultservice1.id}",
									"weight":     "20",
									"port":       "80",
									"protocol":   "HTTP",
								},
								{
									"weight":     "80",
									"service_id": "${alicloud_apig_service.defaultservice2.id}",
								},
							},
						},
					},
					"environment_info": []map[string]interface{}{
						{
							"environment_id": "${alicloud_apig_environment.defaultenvironment.id}",
						},
					},
					"route_name": name,
					"match": []map[string]interface{}{
						{
							"path": []map[string]interface{}{
								{
									"type":  "Exact",
									"value": "/",
								},
							},
							"methods": []string{
								"GET"},
							"headers": []map[string]interface{}{
								{
									"type":  "Exact",
									"name":  "k1",
									"value": "v1",
								},
							},
							"query_params": []map[string]interface{}{
								{
									"type":  "Exact",
									"name":  "p1",
									"value": "v1",
								},
							},
							"ignore_uri_case": "false",
						},
					},
					"http_api_id": "${alicloud_apig_http_api.defaultapi.id}",
					"domain_ids": []string{
						"${alicloud_apig_domain.defaultdomain.id}"},
					"description": "test-route",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"route_name":   name,
						"http_api_id":  CHECKSET,
						"domain_ids.#": "1",
						"description":  "test-route",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"backend": []map[string]interface{}{
						{
							"scene": "MultiServiceByRatio",
							"services": []map[string]interface{}{
								{
									"service_id": "${alicloud_apig_service.defaultservice1.id}",
									"weight":     "50",
								},
								{
									"service_id": "${alicloud_apig_service.defaultservice2.id}",
									"weight":     "20",
								},
								{
									"service_id": "${alicloud_apig_service.defaultservice3.id}",
									"weight":     "30",
								},
							},
						},
					},
					"match": []map[string]interface{}{
						{
							"ignore_uri_case": "true",
							"path": []map[string]interface{}{
								{
									"type":  "Exact",
									"value": "/v1",
								},
							},
							"methods": []string{
								"POST", "GET", "PUT"},
							"headers": []map[string]interface{}{
								{
									"type":  "Prefix",
									"name":  "k1",
									"value": "v1",
								},
							},
							"query_params": []map[string]interface{}{
								{
									"type":  "Prefix",
									"name":  "p1",
									"value": "v1",
								},
							},
						},
					},
					"description": "test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"backend": []map[string]interface{}{
						{
							"scene": "SingleService",
							"services": []map[string]interface{}{
								{
									"service_id": "${alicloud_apig_service.defaultservice1.id}",
								},
							},
						},
					},
					"match": []map[string]interface{}{
						{
							"ignore_uri_case": "false",
							"path": []map[string]interface{}{
								{
									"type":  "Exact",
									"value": "/",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"domain_ids"},
			},
		},
	})
}

var AlicloudApigRouteMap8981 = map[string]string{
	"status":           CHECKSET,
	"create_time":      CHECKSET,
	"route_id":         CHECKSET,
	"gateway_status.%": CHECKSET,
	"update_time":      CHECKSET,
}

func AlicloudApigRouteBasicDependence8981(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "protocol" {
  default = "HTTP"
}

variable "dns_address" {
  default = "httpbin.org:80"
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_apig_gateway" "defaultgateway" {
  network_access_config {
    type = "Intranet"
  }
  vswitch {
    vswitch_id = data.alicloud_vswitches.default.ids.0
  }
  zone_config {
    select_option = "Auto"
  }
  vpc {
    vpc_id = data.alicloud_vpcs.default.ids.0
  }
  payment_type = "PayAsYouGo"
  gateway_name = "${var.name}gw"
  spec         = "apigw.small.x1"
}

resource "alicloud_apig_domain" "defaultdomain" {
  domain_name = "${var.name}a.example.com"
  protocol    = "HTTP"
}

resource "alicloud_apig_domain" "defaultdomain2" {
  domain_name = "${var.name}b.example.com"
  protocol    = "HTTP"
}

resource "alicloud_apig_domain" "defaultdomain3" {
  domain_name = "${var.name}c.example.com"
  protocol    = "HTTP"
}

resource "alicloud_apig_environment" "defaultenvironment" {
  description      = "镇元测试环境"
  environment_name = "${var.name}env"
  gateway_id       = alicloud_apig_gateway.defaultgateway.id
}

resource "alicloud_apig_service" "defaultservice1" {
  service_name = "${var.name}svc1"
  source_type  = "DNS"
  gateway_id   = alicloud_apig_gateway.defaultgateway.id
  addresses    = ["${var.dns_address}"]
}

resource "alicloud_apig_service" "defaultservice2" {
  service_name = "${var.name}svc2"
  source_type  = "DNS"
  gateway_id   = alicloud_apig_gateway.defaultgateway.id
  addresses    = ["${var.dns_address}"]
}

resource "alicloud_apig_service" "defaultservice3" {
  service_name = "${var.name}svc3"
  source_type  = "DNS"
  gateway_id   = alicloud_apig_gateway.defaultgateway.id
  addresses    = ["${var.dns_address}"]
}

resource "alicloud_apig_http_api" "defaultapi" {
  http_api_name = "${var.name}api"
  type          = "Http"
  description   = "test api"
  protocols     = ["${var.protocol}"]
}


`, name)
}

// Case test-route-singleservice_副本1732859242408 9263
func TestAccAliCloudApigRoute_basic9263(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_apig_route.default"
	ra := resourceAttrInit(resourceId, AlicloudApigRouteMap9263)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApigServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeApigRoute")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccapig%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudApigRouteBasicDependence9263)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"backend": []map[string]interface{}{
						{
							"scene": "SingleService",
							"services": []map[string]interface{}{
								{
									"service_id": "${alicloud_apig_service.defaultservice1.id}",
								},
							},
						},
					},
					"environment_info": []map[string]interface{}{
						{
							"environment_id": "${alicloud_apig_environment.defaultenvironment.id}",
						},
					},
					"route_name": name,
					"match": []map[string]interface{}{
						{
							"path": []map[string]interface{}{
								{
									"type":  "Exact",
									"value": "/",
								},
							},
							"methods": []string{
								"GET"},
							"headers": []map[string]interface{}{
								{
									"type":  "Exact",
									"name":  "k1",
									"value": "v1",
								},
							},
							"query_params": []map[string]interface{}{
								{
									"type":  "Exact",
									"name":  "p1",
									"value": "v1",
								},
							},
						},
					},
					"http_api_id": "${alicloud_apig_http_api.defaultapi.id}",
					"domain_ids": []string{
						"${alicloud_apig_domain.defaultdomain.id}"},
					"description": "test-route",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"route_name":   name,
						"http_api_id":  CHECKSET,
						"domain_ids.#": "1",
						"description":  "test-route",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"match": []map[string]interface{}{
						{
							"ignore_uri_case": "true",
							"path": []map[string]interface{}{
								{
									"type":  "Exact",
									"value": "/v1",
								},
							},
							"methods": []string{
								"POST", "GET", "PUT"},
							"headers": []map[string]interface{}{
								{
									"type":  "Prefix",
									"name":  "k1",
									"value": "v1",
								},
								{
									"type":  "Exact",
									"name":  "k2",
									"value": "v2",
								},
								{
									"type":  "Exact",
									"name":  "k3",
									"value": "v3",
								},
							},
							"query_params": []map[string]interface{}{
								{
									"type":  "Prefix",
									"name":  "p1",
									"value": "v1",
								},
								{
									"type":  "Exact",
									"name":  "p2",
									"value": "v2",
								},
								{
									"type":  "Exact",
									"name":  "p3",
									"value": "v3",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"backend": []map[string]interface{}{
						{
							"scene": "MultiServiceByRatio",
							"services": []map[string]interface{}{
								{
									"service_id": "${alicloud_apig_service.defaultservice1.id}",
									"weight":     "20",
								},
								{
									"service_id": "${alicloud_apig_service.defaultservice2.id}",
									"weight":     "80",
								},
							},
						},
					},
					"match": []map[string]interface{}{
						{
							"ignore_uri_case": "true",
							"path": []map[string]interface{}{
								{
									"type":  "Prefix",
									"value": "/v1",
								},
							},
							"headers": []map[string]interface{}{
								{
									"type":  "Exact",
									"name":  "k1",
									"value": "v1",
								},
								{
									"type":  "Exact",
									"name":  "k2",
									"value": "v2",
								},
							},
							"methods": []string{
								"GET", "POST"},
							"query_params": []map[string]interface{}{
								{
									"type":  "Exact",
									"name":  "p1",
									"value": "v1",
								},
								{
									"type":  "Exact",
									"name":  "p2",
									"value": "v2",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"backend": []map[string]interface{}{
						{
							"scene": "SingleService",
							"services": []map[string]interface{}{
								{
									"service_id": "${alicloud_apig_service.defaultservice1.id}",
								},
							},
						},
					},
					"match": []map[string]interface{}{
						{
							"path": []map[string]interface{}{
								{
									"type":  "Exact",
									"value": "/v2",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"domain_ids"},
			},
		},
	})
}

var AlicloudApigRouteMap9263 = map[string]string{
	"status":           CHECKSET,
	"create_time":      CHECKSET,
	"route_id":         CHECKSET,
	"gateway_status.%": CHECKSET,
	"update_time":      CHECKSET,
}

func AlicloudApigRouteBasicDependence9263(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "protocol" {
  default = "HTTP"
}

variable "dns_address" {
  default = "httpbin.org:80"
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_apig_gateway" "defaultgateway" {
  network_access_config {
    type = "Intranet"
  }
  vswitch {
    vswitch_id = data.alicloud_vswitches.default.ids.0
  }
  zone_config {
    select_option = "Auto"
  }
  vpc {
    vpc_id = data.alicloud_vpcs.default.ids.0
  }
  payment_type = "PayAsYouGo"
  gateway_name = "${var.name}gw"
  spec         = "apigw.small.x1"
}

resource "alicloud_apig_domain" "defaultdomain" {
  domain_name = "${var.name}.example.com"
  protocol    = "HTTP"
}

resource "alicloud_apig_environment" "defaultenvironment" {
  description      = "镇元测试环境"
  environment_name = "${var.name}env"
  gateway_id       = alicloud_apig_gateway.defaultgateway.id
}

resource "alicloud_apig_service" "defaultservice1" {
  service_name = "${var.name}svc1"
  source_type  = "DNS"
  gateway_id   = alicloud_apig_gateway.defaultgateway.id
  addresses    = ["${var.dns_address}"]
}

resource "alicloud_apig_service" "defaultservice2" {
  service_name = "${var.name}svc2"
  source_type  = "DNS"
  gateway_id   = alicloud_apig_gateway.defaultgateway.id
  addresses    = ["${var.dns_address}"]
}

resource "alicloud_apig_http_api" "defaultapi" {
  http_api_name = "${var.name}api"
  type          = "Http"
  description   = "test api"
  protocols     = ["${var.protocol}"]
}


`, name)
}

// Test Apig Route domain_ids attribute (scene=Mock, does not depend on apig_service).
func TestAccAliCloudApigRoute_domainIds(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_apig_route.default"
	ra := resourceAttrInit(resourceId, AlicloudApigRouteMapDomainIds)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApigServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeApigRoute")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccapig%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudApigRouteBasicDependenceDomainIds)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"backend": []map[string]interface{}{
						{"scene": "Mock"},
					},
					"description": "domain_ids coverage",
					"route_name":  name,
					"http_api_id": "${alicloud_apig_http_api.route_httpapi_dids.id}",
					"environment_info": []map[string]interface{}{
						{"environment_id": "${alicloud_apig_gateway.route_gateway_dids.environments.0.environment_id}"},
					},
					"match": []map[string]interface{}{
						{
							"path": []map[string]interface{}{
								{"type": "Prefix", "value": "/dids"},
							},
							"methods":         []string{"GET"},
							"ignore_uri_case": "false",
						},
					},
					"domain_ids": []string{
						"${alicloud_apig_domain.route_domain_dids_a.id}",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":  "domain_ids coverage",
						"route_name":   name,
						"http_api_id":  CHECKSET,
						"domain_ids.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"domain_ids": []string{
						"${alicloud_apig_domain.route_domain_dids_a.id}",
						"${alicloud_apig_domain.route_domain_dids_b.id}",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"domain_ids.#": "2",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"domain_ids"},
			},
		},
	})
}

var AlicloudApigRouteMapDomainIds = map[string]string{
	"status":           CHECKSET,
	"create_time":      CHECKSET,
	"route_id":         CHECKSET,
	"gateway_status.%": CHECKSET,
	"update_time":      CHECKSET,
}

func AlicloudApigRouteBasicDependenceDomainIds(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_apig_gateway" "route_gateway_dids" {
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
  gateway_name = format("%%sgw", var.name)
  spec         = "apigw.small.x1"
  log_config {
    sls {
      enable = true
    }
  }
}

resource "alicloud_apig_domain" "route_domain_dids_a" {
  domain_name = format("%%sa.example.com", var.name)
  force_https = false
  protocol    = "HTTP"
}

resource "alicloud_apig_domain" "route_domain_dids_b" {
  domain_name = format("%%sb.example.com", var.name)
  force_https = false
  protocol    = "HTTP"
}

resource "alicloud_apig_http_api" "route_httpapi_dids" {
  http_api_name = format("%%sapi", var.name)
  protocols     = ["HTTP"]
  type          = "Rest"
  description   = "dids test httpapi"
  base_path     = format("/%%s", var.name)
}


`, name)
}

// Test Apig Route. <<< Resource test cases, automatically generated.
