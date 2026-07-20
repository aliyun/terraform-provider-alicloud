// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Apig Service. >>> Resource test cases, automatically generated.
// Case test-service-vip 8886
func TestAccAliCloudApigService_basic8886(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_apig_service.default"
	ra := resourceAttrInit(resourceId, AlicloudApigServiceMap8886)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApigServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeApigService")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccapig%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudApigServiceBasicDependence8886)
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
					"addresses": []string{
						"${var.address}"},
					"service_name": name,
					"source_type":  "VIP",
					"gateway_id":   "${alicloud_apig_gateway.defaultFsRKYn.id}",
					"namespace":    "default",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"addresses.#":  "1",
						"service_name": name,
						"source_type":  "VIP",
						"gateway_id":   CHECKSET,
						"namespace":    "default",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"addresses": []string{
						"${var.address_1}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"addresses.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"addresses": []string{
						"${var.address}", "${var.address_1}", "${var.address_2}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"addresses.#": "3",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"ports", "protocol"},
			},
		},
	})
}

var AlicloudApigServiceMap8886 = map[string]string{
	"update_timestamp":      CHECKSET,
	"create_timestamp":      CHECKSET,
	"unhealthy_endpoints.#": CHECKSET,
	"health_status":         CHECKSET,
	"outlier_endpoints.#":   CHECKSET,
	"runtime_detail_status": CHECKSET,
}

func AlicloudApigServiceBasicDependence8886(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "address" {
  default = "127.0.0.1:8080"
}

variable "address_1" {
  default = "127.0.0.1:7891"
}

variable "address_2" {
  default = "127.0.0.1:7890"
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_apig_gateway" "defaultFsRKYn" {
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
  gateway_name = "zhenyuantest"
  spec         = "apigw.small.x1"
  log_config {
    sls {
      enable = false
    }
  }
}


`, name)
}

// Case 资源组接入_DNS 9199
func TestAccAliCloudApigService_basic9199(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_apig_service.default"
	ra := resourceAttrInit(resourceId, AlicloudApigServiceMap9199)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApigServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeApigService")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccapig%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudApigServiceBasicDependence9199)
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
					"addresses": []string{
						"${var.address}"},
					"service_name":      name,
					"source_type":       "DNS",
					"gateway_id":        "${alicloud_apig_gateway.defaultgateway.id}",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"addresses.#":       "1",
						"service_name":      name,
						"source_type":       "DNS",
						"gateway_id":        CHECKSET,
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"addresses": []string{
						"${var.address1}", "${var.address}", "${var.address2}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"addresses.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"addresses": []string{
						"${var.address2}", "${var.address1}"},
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"addresses.#":       "2",
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"ports", "protocol"},
			},
		},
	})
}

var AlicloudApigServiceMap9199 = map[string]string{
	"update_timestamp":      CHECKSET,
	"create_timestamp":      CHECKSET,
	"unhealthy_endpoints.#": CHECKSET,
	"health_status":         CHECKSET,
	"outlier_endpoints.#":   CHECKSET,
	"runtime_detail_status": CHECKSET,
}

func AlicloudApigServiceBasicDependence9199(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "address" {
  default = "httpbin.org:8080"
}

variable "address2" {
  default = "taobao.com:80"
}

variable "address1" {
  default = "baidu.com:80"
}

data "alicloud_resource_manager_resource_groups" "default" {}

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
  gateway_name = "zhenyuantest"
  spec         = "apigw.small.x1"
  log_config {
    sls {
      enable = false
    }
  }
}


`, name)
}

// Case test-service-fc 10174
func TestAccAliCloudApigService_basic10174(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_apig_service.default"
	ra := resourceAttrInit(resourceId, AlicloudApigServiceMap10174)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApigServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeApigService")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccapig%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudApigServiceBasicDependence10174)
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
					"service_name": name,
					"source_type":  "FC3",
					"gateway_id":   "${alicloud_apig_gateway.defaultgateway.id}",
					"qualifier":    "LATEST",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"service_name": name,
						"source_type":  "FC3",
						"gateway_id":   CHECKSET,
						"qualifier":    "LATEST",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"ports", "protocol"},
			},
		},
	})
}

var AlicloudApigServiceMap10174 = map[string]string{
	"update_timestamp":      CHECKSET,
	"create_timestamp":      CHECKSET,
	"unhealthy_endpoints.#": CHECKSET,
	"health_status":         CHECKSET,
	"outlier_endpoints.#":   CHECKSET,
	"runtime_detail_status": CHECKSET,
}

func AlicloudApigServiceBasicDependence10174(name string) string {
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
  gateway_name = "zhenyuantest"
  spec         = "apigw.small.x1"
  log_config {
    sls {
      enable = false
    }
  }
}

resource "alicloud_fcv3_function" "defaultfc3" {
  memory_size = "512"
  description = "CreateFcFunction"
  timeout     = "3"
  instance_lifecycle_config {
    initializer {
      timeout = "3"
    }
    pre_stop {
      timeout = "3"
    }
  }
  cpu     = 0.5
  handler = "index.handler"
  code {
    zip_file = "UEsDBBQACAAIAAAAAAAAAAAAAAAAAAAAAAAIAAAAaW5kZXgucHmEkEFKxEAQRfd9ig9ZTCJOooIwDMwNXLqXnnQlaalUhU5lRj2KZ/FOXkESGR114bJ/P/7jV4b1xRq1hijtFpM1682cuNgPmgysbRulPT0fRxXnMtwrSPyeCdYRokSLnuMLJTTkbUqEvDMbxm1VdcRD6Tk+T1LW2ldB66knsYdA5iNX17ebm6tN2VnPhcswMPmREPuBacb+CiapLarAj9gT6/H97dVlCNScY3mtYvRkxdZlwDKDEnanPWVLdrdkeXEGlFEazVdfPVHaVeHc3N15CUwppwOJXeK7HshAB8NuOU7J6sP4SRXuH/EvbUfMiqMmDqv5M5FNSfAj/wgAAP//UEsHCPl//NYAAQAArwEAAFBLAQIUABQACAAIAAAAAAD5f/zWAAEAAK8BAAAIAAAAAAAAAAAAAAAAAAAAAABpbmRleC5weVBLBQYAAAAAAQABADYAAAA2AQAAAAA="
  }
  function_name = "ZhenyuanTestFunction-688"
  runtime       = "python3.9"
  disk_size     = "512"
}


`, name)
}

// Case test-service-ai 10175
func TestAccAliCloudApigService_basic10175(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_apig_service.default"
	ra := resourceAttrInit(resourceId, AlicloudApigServiceMap10175)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApigServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeApigService")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccapig%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudApigServiceBasicDependence10175)
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
					"service_name": name,
					"source_type":  "AI",
					"gateway_id":   "${alicloud_apig_gateway.defaultgateway.id}",
					"ai_service_config": []map[string]interface{}{
						{
							"provider": "qwen",
							"protocols": []string{
								"${var.ai_protocol}"},
							"address": "https://dashscope.aliyuncs.com/compatible-mode/v1",
							"api_keys": []string{
								"${var.apikey1}", "${var.apikey3}"},
							"enable_health_check": "true",
						},
					},
					"group_name": "test-group",
					"namespace":  "test-namespace",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"service_name": name,
						"source_type":  "AI",
						"gateway_id":   CHECKSET,
						"group_name":   "test-group",
						"namespace":    "test-namespace",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ai_service_config": []map[string]interface{}{
						{
							"provider": "openai",
							"protocols": []string{
								"${var.ai_protocol}", "HTTP"},
							"address": "https://api.openai.com/v1",
							"api_keys": []string{
								"${var.apikey1}", "${var.apikey2}", "${var.apikey3}"},
							"enable_health_check": "false",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ai_service_config": []map[string]interface{}{
						{
							"protocols": []string{
								"${var.ai_protocol}"},
							"address": "https://api.openai.com/v1",
							"api_keys": []string{
								"${var.apikey2}"},
							"enable_health_check": "false",
							"provider":            "openai",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
		},
	})
}

var AlicloudApigServiceMap10175 = map[string]string{
	"update_timestamp":      CHECKSET,
	"create_timestamp":      CHECKSET,
	"unhealthy_endpoints.#": CHECKSET,
	"health_status":         CHECKSET,
	"outlier_endpoints.#":   CHECKSET,
	"runtime_detail_status": CHECKSET,
}

func AlicloudApigServiceBasicDependence10175(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "apikey2" {
  default = "bbbbbbbbbb"
}

variable "apikey1" {
  default = "aaaaaaaaaa"
}

variable "ai_protocol" {
  default = "OpenAI/v1"
}

variable "apikey3" {
  default = "cccccccccc"
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
  gateway_name = "zhenyuantest"
  spec         = "apigw.small.x1"
  log_config {
    sls {
      enable = false
    }
  }
}


`, name)
}

// Case test-service-full-attr 12000
func TestAccAliCloudApigService_basic12000(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_apig_service.default"
	ra := resourceAttrInit(resourceId, AlicloudApigServiceMap12000)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApigServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeApigService")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccapig%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudApigServiceBasicDependence12000)
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
					"addresses":               []string{"${var.address}"},
					"service_name":            name,
					"source_type":             "VIP",
					"gateway_id":              "${alicloud_apig_gateway.defaultFsRKYn.id}",
					"namespace":               "default",
					"dns_servers":             []string{"8.8.8.8"},
					"express_type":            "express",
					"healthy_panic_threshold": 0.5,
					"model_provider_id":       "mp-test",
					"protocol":                "HTTP",
					"agent_service_config": []map[string]interface{}{
						{
							"address":                  "127.0.0.1:8080",
							"protocols":                []string{"HTTP"},
							"provider":                 "agent",
							"enable_health_check":      "true",
							"enable_outlier_detection": "true",
						},
					},
					"health_check_config": []map[string]interface{}{
						{
							"enable":              "true",
							"expected_statuses":   []string{"200"},
							"healthy_threshold":   3,
							"http_host":           "example.com",
							"http_path":           "/",
							"interval":            10,
							"protocol":            "HTTP",
							"timeout":             5,
							"unhealthy_threshold": 3,
						},
					},
					"outlier_detection_config": []map[string]interface{}{
						{
							"base_ejection_time":               30,
							"enable":                           "true",
							"failure_percentage_minimum_hosts": 5,
							"failure_percentage_threshold":     80,
							"interval":                         10,
						},
					},
					"ports": []map[string]interface{}{
						{
							"name":     "port1",
							"port":     80,
							"protocol": "TCP",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"service_name": name,
						"source_type":  "VIP",
						"gateway_id":   CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"dns_servers":             []string{"8.8.4.4"},
					"healthy_panic_threshold": 0.6,
					"model_provider_id":       "mp-test-2",
					"agent_service_config": []map[string]interface{}{
						{
							"address":                  "127.0.0.1:8081",
							"protocols":                []string{"HTTP", "HTTPS"},
							"provider":                 "agent2",
							"enable_health_check":      "false",
							"enable_outlier_detection": "false",
						},
					},
					"health_check_config": []map[string]interface{}{
						{
							"enable":              "false",
							"expected_statuses":   []string{"200", "201"},
							"healthy_threshold":   5,
							"http_host":           "test.com",
							"http_path":           "/health",
							"interval":            20,
							"protocol":            "TCP",
							"timeout":             10,
							"unhealthy_threshold": 5,
						},
					},
					"outlier_detection_config": []map[string]interface{}{
						{
							"base_ejection_time":               60,
							"enable":                           "false",
							"failure_percentage_minimum_hosts": 10,
							"failure_percentage_threshold":     90,
							"interval":                         20,
						},
					},
					"ports": []map[string]interface{}{
						{
							"name":     "port2",
							"port":     81,
							"protocol": "UDP",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
		},
	})
}

var AlicloudApigServiceMap12000 = map[string]string{
	"update_timestamp":      CHECKSET,
	"create_timestamp":      CHECKSET,
	"unhealthy_endpoints.#": CHECKSET,
	"health_status":         CHECKSET,
	"outlier_endpoints.#":   CHECKSET,
	"runtime_detail_status": CHECKSET,
}

func AlicloudApigServiceBasicDependence12000(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "address" {
  default = "127.0.0.1:8080"
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_apig_gateway" "defaultFsRKYn" {
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
  gateway_name = "zhenyuantest"
  spec         = "apigw.small.x1"
  log_config {
    sls {
      enable = false
    }
  }
}


`, name)
}

// Test Apig Service. <<< Resource test cases, automatically generated.
