// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Apig Policy. >>> Resource test cases, automatically generated.
// Case policy_crud_test_0 13010
func TestAccAliCloudApigPolicy_basic13010(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_apig_policy.default"
	ra := resourceAttrInit(resourceId, AlicloudApigPolicyMap13010)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApigServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeApigPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccapig%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudApigPolicyBasicDependence13010)
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
					"attach_resource_ids": []string{
						"${alicloud_apig_route.policy_route.route_id}"},
					"policy_config":        "{\\\"unit\\\":\\\"Second\\\",\\\"limit\\\":100,\\\"responseStatusCode\\\":429}",
					"environment_id":       "${alicloud_apig_gateway.policy_gateway.environments.0.environment_id}",
					"policy_class_name":    "RateLimit",
					"policy_name":          name,
					"gateway_id":           "${alicloud_apig_gateway.policy_gateway.id}",
					"attach_resource_type": "GatewayRoute",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"attach_resource_ids.#": "1",
						"policy_config":         CHECKSET,
						"environment_id":        CHECKSET,
						"policy_class_name":     "RateLimit",
						"policy_name":           name,
						"gateway_id":            CHECKSET,
						"attach_resource_type":  "GatewayRoute",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_config": "{\\\"unit\\\":\\\"Minute\\\",\\\"limit\\\":1000,\\\"responseStatusCode\\\":429}",
					"policy_name":   name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy_config": CHECKSET,
						"policy_name":   name + "_update",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"attach_resource_ids", "attach_resource_type", "environment_id", "gateway_id"},
			},
		},
	})
}

var AlicloudApigPolicyMap13010 = map[string]string{
	"policy_class_id": CHECKSET,
}

func AlicloudApigPolicyBasicDependence13010(name string) string {
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

resource "alicloud_apig_gateway" "policy_gateway" {
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

resource "alicloud_apig_service" "policy_svc" {
  service_name = "${var.name}svc"
  source_type  = "DNS"
  gateway_id   = alicloud_apig_gateway.policy_gateway.id
  addresses    = ["httpbin.org:80"]
}

resource "alicloud_apig_http_api" "policy_api" {
  http_api_name = "${var.name}api"
  protocols     = ["HTTP"]
  type          = "Rest"
  description   = "policy test httpapi"
  base_path     = "/${var.name}"
}

resource "alicloud_apig_route" "policy_route" {
  route_name  = "${var.name}route"
  http_api_id = alicloud_apig_http_api.policy_api.id
  environment_info {
    environment_id = alicloud_apig_gateway.policy_gateway.environments.0.environment_id
  }
  match {
    path {
      type  = "Prefix"
      value = "/policy-path"
    }
  }
  backend {
    scene = "SingleService"
    services {
      port       = 8080
      protocol   = "HTTP"
      service_id = alicloud_apig_service.policy_svc.id
    }
  }
}


`, name)
}

// Case test-policy-service 9256
func TestAccAliCloudApigPolicy_basic9256(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_apig_policy.default"
	ra := resourceAttrInit(resourceId, AlicloudApigPolicyMap9256)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApigServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeApigPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccapig%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudApigPolicyBasicDependence9256)
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
					"attach_resource_ids": []string{
						"${alicloud_apig_service.defaultservice.id}"},
					"policy_config":        "{\\\"mode\\\":\\\"SIMPLE\\\",\\\"sni\\\":\\\"aaaa\\\",\\\"enable\\\":true}",
					"policy_name":          name,
					"gateway_id":           "${alicloud_apig_gateway.defaultgateway.id}",
					"attach_resource_type": "GatewayService",
					"policy_class_name":    "ServiceTls",
					"environment_id":       "${alicloud_apig_gateway.defaultgateway.environments.0.environment_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"attach_resource_ids.#": "1",
						"policy_config":         CHECKSET,
						"policy_name":           name,
						"gateway_id":            CHECKSET,
						"attach_resource_type":  "GatewayService",
						"policy_class_name":     "ServiceTls",
						"environment_id":        CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_config": "{\\\"mode\\\":\\\"SIMPLE\\\",\\\"sni\\\":\\\"bbbb\\\",\\\"enable\\\":false}",
					"policy_name":   name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy_config": CHECKSET,
						"policy_name":   name + "_update",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"attach_resource_ids", "attach_resource_type", "environment_id", "gateway_id"},
			},
		},
	})
}

var AlicloudApigPolicyMap9256 = map[string]string{
	"policy_class_id": CHECKSET,
}

func AlicloudApigPolicyBasicDependence9256(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_vpc" "defaultvpc" {
        cidr_block = "192.168.0.0/16"
        vpc_name = "zhenyuan-test"
}

resource "alicloud_vswitch" "defaultvswitch" {
        vpc_id = "${alicloud_vpc.defaultvpc.id}"
        zone_id = "cn-hangzhou-b"
        cidr_block = "192.168.15.0/24"
        vswitch_name = "zhenyuan-test"
}

resource "alicloud_apig_gateway" "defaultgateway" {
            network_access_config  {
                type = "Intranet"
        }
            vswitch  {
                vswitch_id = "${alicloud_vswitch.defaultvswitch.id}"
        }
            zone_config  {
                select_option = "Auto"
        }
            vpc  {
                vpc_id = "${alicloud_vpc.defaultvpc.id}"
        }
        payment_type = "PayAsYouGo"
        gateway_name = "test"
        spec = "apigw.small.x1"
}


resource "alicloud_apig_service" "defaultservice" {
                        addresses = [            "httpbin.org:8080"         ]
            service_name = "zhenyuan-test-1787970491"
        source_type = "DNS"
        gateway_id = "${alicloud_apig_gateway.defaultgateway.id}"
}


`, name)
}

// Case test-policy-service 7105
func TestAccAliCloudApigPolicy_basic7105(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_apig_policy.default"
	ra := resourceAttrInit(resourceId, AlicloudApigPolicyMap7105)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApigServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeApigPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccapig%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudApigPolicyBasicDependence7105)
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
					"attach_resource_ids": []string{
						"${alicloud_apig_service.defaultservice.id}"},
					"policy_config":        "{\\\"mode\\\":\\\"SIMPLE\\\",\\\"sni\\\":\\\"aaaa\\\",\\\"enable\\\":true}",
					"policy_name":          name,
					"gateway_id":           "${alicloud_apig_gateway.defaultgateway.id}",
					"attach_resource_type": "GatewayService",
					"policy_class_name":    "ServiceTls",
					"environment_id":       "${alicloud_apig_gateway.defaultgateway.environments.0.environment_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"attach_resource_ids.#": "1",
						"policy_config":         CHECKSET,
						"policy_name":           name,
						"gateway_id":            CHECKSET,
						"attach_resource_type":  "GatewayService",
						"policy_class_name":     "ServiceTls",
						"environment_id":        CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_config": "{\\\"mode\\\":\\\"SIMPLE\\\",\\\"sni\\\":\\\"bbbb\\\",\\\"enable\\\":false}",
					"policy_name":   name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy_config": CHECKSET,
						"policy_name":   name + "_update",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"attach_resource_ids", "attach_resource_type", "environment_id", "gateway_id"},
			},
		},
	})
}

var AlicloudApigPolicyMap7105 = map[string]string{
	"policy_class_id": CHECKSET,
}

func AlicloudApigPolicyBasicDependence7105(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_vpc" "defaultvpc" {
        cidr_block = "192.168.0.0/16"
        vpc_name = "zhenyuan-test"
}

resource "alicloud_vswitch" "defaultvswitch" {
        vpc_id = "${alicloud_vpc.defaultvpc.id}"
        zone_id = "cn-hangzhou-b"
        cidr_block = "192.168.15.0/24"
        vswitch_name = "zhenyuan-test"
}

resource "alicloud_apig_gateway" "defaultgateway" {
            network_access_config  {
                type = "Intranet"
        }
            vswitch  {
                vswitch_id = "${alicloud_vswitch.defaultvswitch.id}"
        }
            zone_config  {
                select_option = "Auto"
        }
            vpc  {
                vpc_id = "${alicloud_vpc.defaultvpc.id}"
        }
        payment_type = "PayAsYouGo"
        gateway_name = "test"
        spec = "apigw.small.x1"
}


resource "alicloud_apig_service" "defaultservice" {
                        addresses = [            "httpbin.org:8080"         ]
            service_name = "zhenyuan-test-1787970497"
        source_type = "DNS"
        gateway_id = "${alicloud_apig_gateway.defaultgateway.id}"
}


`, name)
}

// Test Apig Policy. <<< Resource test cases, automatically generated.
