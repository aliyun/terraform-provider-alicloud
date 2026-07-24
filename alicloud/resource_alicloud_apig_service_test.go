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
					"service_name":            name,
					"source_type":             "VIP",
					"gateway_id":              "${alicloud_apig_gateway.defaultFsRKYn.id}",
					"namespace":               "default",
					"express_type":            "Standard",
					"protocol":                "HTTP",
					"healthy_panic_threshold": "50",
					"health_check_config": []map[string]interface{}{
						{
							"enable":              "true",
							"protocol":            "HTTP",
							"http_path":           "/health",
							"http_host":           "www.example.com",
							"healthy_threshold":   "2",
							"unhealthy_threshold": "2",
							"timeout":             "5",
							"interval":            "10",
							"expected_statuses":   []string{"200"},
						},
					},
					"outlier_detection_config": []map[string]interface{}{
						{
							"enable":                           "true",
							"base_ejection_time":               "30",
							"interval":                         "30",
							"failure_percentage_threshold":     "80",
							"failure_percentage_minimum_hosts": "1",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"addresses.#":                "1",
						"service_name":               name,
						"source_type":                "VIP",
						"gateway_id":                 CHECKSET,
						"namespace":                  "default",
						"protocol":                   "HTTP",
						"healthy_panic_threshold":    "50",
						"ports.#":                    "1",
						"health_check_config.#":      "1",
						"outlier_detection_config.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"addresses": []string{
						"${var.address_1}"},
					"healthy_panic_threshold": "60",
					"health_check_config": []map[string]interface{}{
						{
							"enable":              "true",
							"protocol":            "HTTP",
							"http_path":           "/healthz",
							"http_host":           "api.example.com",
							"healthy_threshold":   "3",
							"unhealthy_threshold": "3",
							"timeout":             "8",
							"interval":            "15",
							"expected_statuses":   []string{"200", "204"},
						},
					},
					"outlier_detection_config": []map[string]interface{}{
						{
							"enable":                           "true",
							"base_ejection_time":               "60",
							"interval":                         "60",
							"failure_percentage_threshold":     "90",
							"failure_percentage_minimum_hosts": "2",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"addresses.#":             "1",
						"healthy_panic_threshold": "60",
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
				Config: testAccConfig(map[string]interface{}{
					"health_check_config": []map[string]interface{}{
						{
							"enable":              "false",
							"protocol":            "TCP",
							"http_path":           "/healthz",
							"http_host":           "api.example.com",
							"healthy_threshold":   "3",
							"unhealthy_threshold": "3",
							"timeout":             "8",
							"interval":            "15",
							"expected_statuses":   []string{"200", "204"},
						},
					},
					"outlier_detection_config": []map[string]interface{}{
						{
							"enable":                           "false",
							"base_ejection_time":               "60",
							"interval":                         "60",
							"failure_percentage_threshold":     "90",
							"failure_percentage_minimum_hosts": "2",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_config.0.enable":      "false",
						"health_check_config.0.protocol":    "TCP",
						"outlier_detection_config.0.enable": "false",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"protocol", "update_timestamp",
					"unhealthy_endpoints", "outlier_endpoints",
				},
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
					"dns_servers": []string{
						"100.100.2.136:53"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"addresses.#":       "1",
						"service_name":      name,
						"source_type":       "DNS",
						"gateway_id":        CHECKSET,
						"resource_group_id": CHECKSET,
						"dns_servers.#":     "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"addresses": []string{
						"${var.address1}", "${var.address}", "${var.address2}"},
					"dns_servers": []string{
						"100.100.2.136:53", "100.100.2.138:53"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"addresses.#":   "3",
						"dns_servers.#": "2",
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
				ImportStateVerifyIgnore: []string{"protocol", "update_timestamp", "unhealthy_endpoints", "outlier_endpoints"},
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
					"service_name": "${alicloud_fcv3_function.defaultfc3.function_name}",
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
				ImportStateVerifyIgnore: []string{"protocol", "update_timestamp", "unhealthy_endpoints", "outlier_endpoints"},
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
  cpu         = 0.5
  handler     = "index.Handler"
  code {
    zip_file = "UEsDBBQACAAIAAAAAAAAAAAAAAAAAAAAAAAIAAAAaW5kZXgucHmEkEFKxEAQRfd9ig9ZTCJOooIwDMwNXLqXnnQlaalUhU5lRj2KZ/FOXkESGR114bJ/P/7jV4b1xRq1hijtFpM1682cuNgPmgysbRulPT0fRxXnMtwrSPyeCdYRokSLnuMLJTTkbUqEvDMbxm1VdcRD6Tk+T1LW2ldB66knsYdA5iNX17ebm6tN2VnPhcswMPmREPuBacb+CiapLarAj9gT6/H97dVlCNScY3mtYvRkxdZlwDKDEnanPWVLdrdkeXEGlFEazVdfPVHaVeHc3N15CUwppwOJXeK7HshAB8NuOU7J6sP4SRXuH/EvbUfMiqMmDqv5M5FNSfAj/wgAAP//UEsHCPl//NYAAQAArwEAAFBLAQIUABQACAAIAAAAAAD5f/zWAAEAAK8BAAAIAAAAAAAAAAAAAAAAAAAAAABpbmRleC5weVBLBQYAAAAAAQABADYAAAA2AQAAAAA="
  }
  function_name = var.name
  runtime       = "python3.9"
  disk_size     = "512"
  log_config {
    log_begin_rule = "None"
  }
}


`, name)
}

// Test Apig Service. <<< Resource test cases, automatically generated.
