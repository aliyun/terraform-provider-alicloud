package alicloud

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test Apig Gateway. >>> Resource test cases, automatically generated.
// Case 资源组接入_prepay 9249
func TestAccAliCloudApigGateway_basic9249(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_apig_gateway.default"
	ra := resourceAttrInit(resourceId, AlicloudApigGatewayMap9249)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApigServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeApigGateway")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sapiggateway%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudApigGatewayBasicDependence9249)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"gateway_name": name,
					"spec":         "apigw.small.x1",
					"vpc": []map[string]interface{}{
						{
							"vpc_id": "${data.alicloud_vpcs.default.ids.0}",
						},
					},
					"network_access_config": []map[string]interface{}{
						{
							"type": "Intranet",
						},
					},
					"zone_config": []map[string]interface{}{
						{
							"select_option": "Auto",
						},
					},
					"vswitch": []map[string]interface{}{
						{
							"vswitch_id": "${data.alicloud_vswitches.default.ids.0}",
						},
					},
					"log_config": []map[string]interface{}{
						{
							"sls": []map[string]interface{}{
								{
									"enable": "false",
								},
							},
						},
					},
					"payment_type":      "Subscription",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"gateway_name":      name,
						"spec":              "apigw.small.x1",
						"payment_type":      "Subscription",
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"log_config", "network_access_config", "zone_config"},
			},
		},
	})
}

var AlicloudApigGatewayMap9249 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
}

func AlicloudApigGatewayBasicDependence9249(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}
data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
}


`, name)
}

// Case 资源组接入_postpay 9246
func TestAccAliCloudApigGateway_basic9246(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_apig_gateway.default"
	ra := resourceAttrInit(resourceId, AlicloudApigGatewayMap9246)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApigServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeApigGateway")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sapiggateway%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudApigGatewayBasicDependence9246)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"gateway_name": name,
					"spec":         "apigw.small.x1",
					"vpc": []map[string]interface{}{
						{
							"vpc_id": "${data.alicloud_vpcs.default.ids.0}",
						},
					},
					"network_access_config": []map[string]interface{}{
						{
							"type": "Intranet",
						},
					},
					"zone_config": []map[string]interface{}{
						{
							"select_option": "Auto",
						},
					},
					"vswitch": []map[string]interface{}{
						{
							"vswitch_id": "${data.alicloud_vswitches.default.ids.0}",
						},
					},
					"log_config": []map[string]interface{}{
						{
							"sls": []map[string]interface{}{
								{
									"enable": "false",
								},
							},
						},
					},
					"payment_type":      "PayAsYouGo",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"gateway_name":      name,
						"spec":              "apigw.small.x1",
						"payment_type":      "PayAsYouGo",
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"gateway_name":      name + "_update",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"gateway_name":      name + "_update",
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"log_config", "network_access_config", "zone_config"},
			},
		},
	})
}

var AlicloudApigGatewayMap9246 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
}

func AlicloudApigGatewayBasicDependence9246(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}
data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
}


`, name)
}

// Case aigateway_postpay 10903
func TestAccAliCloudApigGateway_basic10903(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_apig_gateway.default"
	ra := resourceAttrInit(resourceId, AlicloudApigGatewayMap10903)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApigServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeApigGateway")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccapig%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudApigGatewayBasicDependence10903)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"gateway_name": name,
					"spec":         "aigw.small.x1",
					"vpc": []map[string]interface{}{
						{
							"vpc_id": "${alicloud_vswitch.j.vpc_id}",
						},
					},
					"network_access_config": []map[string]interface{}{
						{
							"type": "Intranet",
						},
					},
					"zone_config": []map[string]interface{}{
						{
							"select_option": "Manual",
						},
					},
					"vswitch": []map[string]interface{}{
						{
							"vswitch_id": "${alicloud_vswitch.j.id}",
						},
					},
					"log_config": []map[string]interface{}{
						{
							"sls": []map[string]interface{}{
								{
									"enable": "false",
								},
							},
						},
					},
					"zones": []map[string]interface{}{
						{
							"vswitch_id": "${alicloud_vswitch.j.id}",
							"zone_id":    "${alicloud_vswitch.j.zone_id}",
						},
						{
							"vswitch_id": "${alicloud_vswitch.k.id}",
							"zone_id":    "${alicloud_vswitch.k.zone_id}",
						},
					},
					"payment_type":      "PayAsYouGo",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"gateway_type":      "AI",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"gateway_name":      name,
						"spec":              "aigw.small.x1",
						"payment_type":      "PayAsYouGo",
						"resource_group_id": CHECKSET,
						"gateway_type":      "AI",
						"zones.#":           "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"gateway_name":      name + "_update",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"gateway_name":      name + "_update",
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"log_config", "network_access_config", "zone_config"},
			},
		},
	})
}

var AlicloudApigGatewayMap10903 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
}

func AlicloudApigGatewayBasicDependence10903(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_vpc" "default" {
  cidr_block = "172.16.0.0/12"
  vpc_name   = var.name
}

resource "alicloud_vswitch" "k" {
  vpc_id     = alicloud_vpc.default.id
  zone_id    = "cn-hangzhou-k"
  cidr_block = "172.16.1.0/24"
}

resource "alicloud_vswitch" "j" {
  vpc_id     = alicloud_vpc.default.id
  zone_id    = "cn-hangzhou-j"
  cidr_block = "172.16.2.0/24"
}


`, name)
}

// collection-order: bypass zones -- stock schema is not modified by this change; cn-hangzhou exposes only two APIG zones (ListZones, 2026-10-20) and the gateway-level vSwitch always follows the first zone, so a reorder apply cannot keep every non-target attribute unchanged
// Test Apig Gateway. zones TypeList reorder convergence coverage.
func TestAccAliCloudApigGateway_zoneOrder(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_apig_gateway.default"
	ra := resourceAttrInit(resourceId, AlicloudApigGatewayMap10903)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApigServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeApigGateway")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccapig%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudApigGatewayBasicDependence10903)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		// tags is Optional-only and only populated by Read, while Create ends with an Update call,
		// so the post-apply state legitimately lacks tags.% and the ID-only refresh would add it.
		IDRefreshIgnore: []string{"tags"},
		Providers:       testAccProviders,
		CheckDestroy:    rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"gateway_name": name,
					"spec":         "aigw.small.x1",
					"vpc": []map[string]interface{}{
						{
							"vpc_id": "${alicloud_vswitch.j.vpc_id}",
						},
					},
					"network_access_config": []map[string]interface{}{
						{
							"type": "Intranet",
						},
					},
					"zone_config": []map[string]interface{}{
						{
							"select_option": "Manual",
						},
					},
					"vswitch": []map[string]interface{}{
						{
							"vswitch_id": "${alicloud_vswitch.j.id}",
						},
					},
					"log_config": []map[string]interface{}{
						{
							"sls": []map[string]interface{}{
								{
									"enable": "false",
								},
							},
						},
					},
					"zones": []map[string]interface{}{
						{
							"vswitch_id": "${alicloud_vswitch.j.id}",
							"zone_id":    "${alicloud_vswitch.j.zone_id}",
						},
						{
							"vswitch_id": "${alicloud_vswitch.k.id}",
							"zone_id":    "${alicloud_vswitch.k.zone_id}",
						},
					},
					"payment_type":      "PayAsYouGo",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"gateway_type":      "AI",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"zones.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"zones": []map[string]interface{}{
						{
							"vswitch_id": "${alicloud_vswitch.k.id}",
							"zone_id":    "${alicloud_vswitch.k.zone_id}",
						},
						{
							"vswitch_id": "${alicloud_vswitch.j.id}",
							"zone_id":    "${alicloud_vswitch.j.zone_id}",
						},
					},
				}),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"vswitch": []map[string]interface{}{
						{
							"vswitch_id": "${alicloud_vswitch.k.id}",
						},
					},
					"zones": []map[string]interface{}{
						{
							"vswitch_id": "${alicloud_vswitch.k.id}",
							"zone_id":    "${alicloud_vswitch.k.zone_id}",
						},
						{
							"vswitch_id": "${alicloud_vswitch.j.id}",
							"zone_id":    "${alicloud_vswitch.j.zone_id}",
						},
					},
				}),
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

// Test Apig Gateway. <<< Resource test cases, automatically generated.

// Test Apig Gateway. gateway_edition coverage.
func TestAccAliCloudApigGateway_basic_edition(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_apig_gateway.default"
	ra := resourceAttrInit(resourceId, AliCloudApigGatewayEditionMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApigServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeApigGateway")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sapiggwedition%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudApigGatewayEditionDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"gateway_name":    name,
					"gateway_edition": "Professional",
					"spec":            "apigw.small.x1",
					"vpc": []map[string]interface{}{
						{
							"vpc_id": "${data.alicloud_vpcs.default.ids.0}",
						},
					},
					"network_access_config": []map[string]interface{}{
						{
							"type": "Intranet",
						},
					},
					"zone_config": []map[string]interface{}{
						{
							"select_option": "Auto",
						},
					},
					"vswitch": []map[string]interface{}{
						{
							"vswitch_id": "${data.alicloud_vswitches.default.ids.0}",
						},
					},
					"log_config": []map[string]interface{}{
						{
							"sls": []map[string]interface{}{
								{
									"enable": "false",
								},
							},
						},
					},
					"payment_type":      "PayAsYouGo",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"gateway_name":    name,
						"gateway_edition": "Professional",
						"spec":            "apigw.small.x1",
						"payment_type":    "PayAsYouGo",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"log_config", "network_access_config", "zone_config"},
			},
		},
	})
}

var AliCloudApigGatewayEditionMap = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
}

func AliCloudApigGatewayEditionDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
}
`, name)
}

// Case 自动续费_renewal: covers Subscription auto-renew fields via BssOpenApi SetRenewal/QueryAvailableInstances.
func TestAccAliCloudApigGateway_basicRenewal(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_apig_gateway.default"
	ra := resourceAttrInit(resourceId, AlicloudApigGatewayMapRenewal)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApigServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeApigGateway")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sapiggateway%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudApigGatewayBasicDependenceRenewal)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		// tags is Optional-only and only populated by Read, while Create ends with an Update call,
		// so the post-apply state legitimately lacks tags.% and the ID-only refresh would add it.
		IDRefreshIgnore: []string{"tags"},
		Providers:       testAccProviders,
		CheckDestroy:    nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"gateway_name": name,
					"spec":         "apigw.small.x1",
					"vpc": []map[string]interface{}{
						{
							"vpc_id": "${data.alicloud_vpcs.default.ids.0}",
						},
					},
					"network_access_config": []map[string]interface{}{
						{
							"type": "Intranet",
						},
					},
					"zone_config": []map[string]interface{}{
						{
							"select_option": "Auto",
						},
					},
					"vswitch": []map[string]interface{}{
						{
							"vswitch_id": "${data.alicloud_vswitches.default.ids.0}",
						},
					},
					"log_config": []map[string]interface{}{
						{
							"sls": []map[string]interface{}{
								{
									"enable": "false",
								},
							},
						},
					},
					"payment_type":        "Subscription",
					"resource_group_id":   "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"renewal_status":      "AutoRenewal",
					"renew_period":        "1",
					"renewal_period_unit": "M",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"gateway_name":        name,
						"spec":                "apigw.small.x1",
						"payment_type":        "Subscription",
						"resource_group_id":   CHECKSET,
						"renewal_status":      "AutoRenewal",
						"renew_period":        "1",
						"renewal_period_unit": "M",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"renew_period": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"renewal_status":      "AutoRenewal",
						"renew_period":        "2",
						"renewal_period_unit": "M",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"renew_period":        "1",
					"renewal_period_unit": "Y",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"renewal_status":      "AutoRenewal",
						"renew_period":        "1",
						"renewal_period_unit": "Y",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"renewal_status": "ManualRenewal",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"renewal_status": "ManualRenewal",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"log_config", "network_access_config", "zone_config"},
			},
		},
	})
}

func TestUnitAliCloudApigGatewayRenewalInstanceResponse(t *testing.T) {
	requestErr := errors.New("SignatureDoesNotMatch")
	if _, err := apigGatewayRenewalInstance("gw-renewal", nil, requestErr); !errors.Is(err, requestErr) {
		t.Fatalf("expected the original request error to be returned, got %v", err)
	} else if NotFoundError(err) {
		t.Fatalf("expected the request error not to be reported as NotFound, got %v", err)
	}

	if _, err := apigGatewayRenewalInstance("gw-renewal", nil, nil); err == nil {
		t.Fatal("expected an error for an empty BssOpenApi response")
	} else if NotFoundError(err) {
		t.Fatalf("expected an empty response not to be reported as NotFound, got %v", err)
	}

	emptyResponse := map[string]interface{}{"Data": map[string]interface{}{"InstanceList": []interface{}{}}}
	if _, err := apigGatewayRenewalInstance("gw-renewal", emptyResponse, nil); err == nil {
		t.Fatal("expected an error for an empty BssOpenApi InstanceList")
	} else if NotFoundError(err) {
		t.Fatalf("expected an empty InstanceList not to be reported as NotFound, got %v", err)
	}

	entry := map[string]interface{}{"InstanceID": "gw-renewal", "RenewStatus": "AutoRenewal", "RenewalDuration": "2", "RenewalDurationUnit": "M"}
	response := map[string]interface{}{"Data": map[string]interface{}{"InstanceList": []interface{}{entry}}}
	instance, err := apigGatewayRenewalInstance("gw-renewal", response, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(instance, entry) {
		t.Fatalf("expected %v, got %v", entry, instance)
	}
}

func TestUnitAliCloudApigGatewayRenewalInstanceSelection(t *testing.T) {
	other := map[string]interface{}{"InstanceID": "gw-other", "RenewStatus": "ManualRenewal", "RenewalDuration": "1", "RenewalDurationUnit": "Y"}
	target := map[string]interface{}{"InstanceID": "gw-renewal", "RenewStatus": "AutoRenewal", "RenewalDuration": "2", "RenewalDurationUnit": "M"}

	unordered := map[string]interface{}{"Data": map[string]interface{}{"InstanceList": []interface{}{other, target}}}
	instance, err := apigGatewayRenewalInstance("gw-renewal", unordered, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(instance, target) {
		t.Fatalf("expected the entry of gw-renewal to be selected regardless of order, got %v", instance)
	}

	missing := map[string]interface{}{"Data": map[string]interface{}{"InstanceList": []interface{}{other}}}
	if _, err := apigGatewayRenewalInstance("gw-renewal", missing, nil); err == nil {
		t.Fatal("expected an error when the response holds no entry for the gateway")
	} else if !strings.Contains(err.Error(), "gw-renewal") {
		t.Fatalf("expected the error to name the gateway, got %v", err)
	} else if NotFoundError(err) {
		t.Fatalf("expected a missing entry not to be reported as NotFound, got %v", err)
	}

	unexpected := map[string]interface{}{"Data": map[string]interface{}{"InstanceList": []interface{}{"gw-renewal"}}}
	if _, err := apigGatewayRenewalInstance("gw-renewal", unexpected, nil); err == nil {
		t.Fatal("expected an error for an unexpected instance entry")
	}
}

func TestUnitAliCloudApigGatewayReadRenewalStateProtection(t *testing.T) {
	// Mirrors the renewal branch of resourceAliCloudApigGatewayRead: the service
	// query is replaced by its pure response helper, the state writer is the real
	// function the Read path calls.
	readRenewalState := func(d *schema.ResourceData, response map[string]interface{}, responseErr error) error {
		instance, queryErr := apigGatewayRenewalInstance("gw-renewal", response, responseErr)
		return setApigGatewayRenewalState(d, instance, queryErr)
	}
	newResourceData := func(t *testing.T) *schema.ResourceData {
		return schema.TestResourceDataRaw(t, resourceAliCloudApigGateway().Schema, map[string]interface{}{
			"payment_type":        "Subscription",
			"renewal_status":      "AutoRenewal",
			"renew_period":        2,
			"renewal_period_unit": "M",
		})
	}
	assertStateKept := func(t *testing.T, d *schema.ResourceData) {
		t.Helper()
		if got := d.Get("renewal_status").(string); got != "AutoRenewal" {
			t.Fatalf("renewal_status must survive a failed BssOpenApi query, got %q", got)
		}
		if got := d.Get("renew_period").(int); got != 2 {
			t.Fatalf("renew_period must survive a failed BssOpenApi query, got %d", got)
		}
		if got := d.Get("renewal_period_unit").(string); got != "M" {
			t.Fatalf("renewal_period_unit must survive a failed BssOpenApi query, got %q", got)
		}
	}

	requestErr := errors.New("SignatureDoesNotMatch")
	emptyResponse := map[string]interface{}{"Data": map[string]interface{}{"InstanceList": []interface{}{}}}
	cases := []struct {
		name     string
		response map[string]interface{}
		err      error
	}{
		{name: "request-error", response: nil, err: requestErr},
		{name: "empty-response", response: nil, err: nil},
		{name: "empty-instance-list", response: emptyResponse, err: nil},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			d := newResourceData(t)
			err := readRenewalState(d, c.response, c.err)
			if err == nil {
				t.Fatal("expected the failed BssOpenApi query to be reported")
			}
			if NotFoundError(err) {
				t.Fatalf("expected the failure not to be reported as NotFound, got %v", err)
			}
			if c.err != nil && !strings.Contains(err.Error(), c.err.Error()) {
				t.Fatalf("expected the original error to be preserved, got %v", err)
			}
			assertStateKept(t, d)
		})
	}

	t.Run("successful-query-refreshes-state", func(t *testing.T) {
		d := newResourceData(t)
		response := map[string]interface{}{"Data": map[string]interface{}{"InstanceList": []interface{}{
			map[string]interface{}{"InstanceID": "gw-renewal", "RenewStatus": "ManualRenewal", "RenewalDuration": 1, "RenewalDurationUnit": "Y"},
		}}}
		if err := readRenewalState(d, response, nil); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got := d.Get("renewal_status").(string); got != "ManualRenewal" {
			t.Fatalf("expected renewal_status to be refreshed, got %q", got)
		}
		if got := d.Get("renew_period").(int); got != 1 {
			t.Fatalf("expected renew_period to be refreshed, got %d", got)
		}
		if got := d.Get("renewal_period_unit").(string); got != "Y" {
			t.Fatalf("expected renewal_period_unit to be refreshed, got %q", got)
		}
	})
}

func TestUnitAliCloudApigGatewayRenewalConverged(t *testing.T) {
	autoRenewal := map[string]interface{}{"RenewalStatus": "AutoRenewal", "RenewalPeriod": "2", "RenewalPeriodUnit": "M"}

	renewalRaw := map[string]interface{}{"RenewStatus": "AutoRenewal", "RenewalDuration": "2", "RenewalDurationUnit": "M"}
	if !apigGatewayRenewalConverged(renewalRaw, autoRenewal) {
		t.Fatal("expected matching AutoRenewal settings to be converged")
	}

	staleRaw := map[string]interface{}{"RenewStatus": "AutoRenewal", "RenewalDuration": "1", "RenewalDurationUnit": "M"}
	if apigGatewayRenewalConverged(staleRaw, autoRenewal) {
		t.Fatal("expected a stale renewal duration not to be converged")
	}

	staleUnitRaw := map[string]interface{}{"RenewStatus": "AutoRenewal", "RenewalDuration": "2", "RenewalDurationUnit": "Y"}
	if apigGatewayRenewalConverged(staleUnitRaw, autoRenewal) {
		t.Fatal("expected a stale renewal duration unit not to be converged")
	}

	if apigGatewayRenewalConverged(staleRaw, map[string]interface{}{"RenewalStatus": "ManualRenewal"}) {
		t.Fatal("expected a stale renewal status not to be converged")
	}

	manualRenewal := map[string]interface{}{"RenewalStatus": "ManualRenewal", "RenewalPeriod": "2", "RenewalPeriodUnit": "M"}
	manualRaw := map[string]interface{}{"RenewStatus": "ManualRenewal", "RenewalDuration": "1", "RenewalDurationUnit": "Y"}
	if !apigGatewayRenewalConverged(manualRaw, manualRenewal) {
		t.Fatal("expected renewal duration and unit not to be compared for ManualRenewal")
	}

	if !apigGatewayRenewalConverged(map[string]interface{}{"RenewStatus": "Any"}, map[string]interface{}{}) {
		t.Fatal("expected convergence to be skipped when no renewal status is requested")
	}
}

func TestUnitAliCloudApigGatewayWaitRenewalConvergence(t *testing.T) {
	autoRenewal := map[string]interface{}{"RenewalStatus": "AutoRenewal", "RenewalPeriod": "2", "RenewalPeriodUnit": "M"}
	values := []map[string]interface{}{
		{"RenewStatus": "AutoRenewal", "RenewalDuration": "1", "RenewalDurationUnit": "M"},
		{"RenewStatus": "AutoRenewal", "RenewalDuration": "2", "RenewalDurationUnit": "M"},
	}
	calls := 0
	waits := 0
	err := waitForApigGatewayRenewalConvergence("gw-renewal", autoRenewal, time.Minute, func() {
		waits++
	}, func() (map[string]interface{}, error) {
		renewalRaw := values[calls]
		calls++
		return renewalRaw, nil
	})
	if err != nil {
		t.Fatalf("expected the stale value to converge on the second query, got %v", err)
	}
	if calls != 2 || waits != 1 {
		t.Fatalf("expected 2 queries and 1 wait, got %d queries and %d waits", calls, waits)
	}

	calls = 0
	waits = 0
	err = waitForApigGatewayRenewalConvergence("gw-renewal", autoRenewal, 0, func() {
		waits++
	}, func() (map[string]interface{}, error) {
		renewalRaw := values[0]
		calls++
		return renewalRaw, nil
	})
	if err == nil {
		t.Fatal("expected a timeout error when the renewal settings never converge")
	}
	if !strings.Contains(err.Error(), "converge") || !strings.Contains(err.Error(), "gw-renewal") {
		t.Fatalf("expected a clear convergence timeout error, got %v", err)
	}
	if calls != 1 || waits != 0 {
		t.Fatalf("expected 1 query and no wait before the timeout, got %d queries and %d waits", calls, waits)
	}

	queryErr := errors.New("BssOpenApi is unavailable")
	err = waitForApigGatewayRenewalConvergence("gw-renewal", autoRenewal, 0, func() {}, func() (map[string]interface{}, error) {
		return nil, queryErr
	})
	if !errors.Is(err, queryErr) {
		t.Fatalf("expected the last query error to be reported on timeout, got %v", err)
	}
	if !strings.Contains(err.Error(), "converge") || !strings.Contains(err.Error(), "gw-renewal") || !strings.Contains(err.Error(), "expected ") || !strings.Contains(err.Error(), "RenewalStatus:AutoRenewal") {
		t.Fatalf("expected the timeout error to carry the gateway and expected renewal settings, got %v", err)
	}
}

var AlicloudApigGatewayMapRenewal = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
}

func AlicloudApigGatewayBasicDependenceRenewal(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}
data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
}


`, name)
}

// TestUnitAliCloudApigGatewayPublicReadRenewalStateProtection drives the public
// resourceAliCloudApigGatewayRead with an isolated APIG/BSS client fixture. It
// guards the original regression: an already known renewal configuration must
// survive a BSS failure instead of being cleared by the read path.
func TestUnitAliCloudApigGatewayPublicReadRenewalStateProtection(t *testing.T) {
	type scenario struct {
		name       string
		id         string
		bssStatus  int
		bssBody    string
		expectErr  bool
		wantPeriod int
		wantUnit   string
		wantStatus string
	}
	scenarios := []scenario{
		{
			name:       "bss-request-error",
			id:         "gw-unit-request-error",
			bssStatus:  http.StatusForbidden,
			bssBody:    `{"Code":"SignatureDoesNotMatch","Message":"The request signature we calculated does not match the signature you provided."}`,
			expectErr:  true,
			wantPeriod: 2,
			wantUnit:   "M",
			wantStatus: "AutoRenewal",
		},
		{
			name:       "bss-empty-response",
			id:         "gw-unit-empty-response",
			bssStatus:  http.StatusOK,
			bssBody:    `{}`,
			expectErr:  true,
			wantPeriod: 2,
			wantUnit:   "M",
			wantStatus: "AutoRenewal",
		},
		{
			name:       "bss-empty-instance-list",
			id:         "gw-unit-empty-list",
			bssStatus:  http.StatusOK,
			bssBody:    `{"Data":{"InstanceList":[]}}`,
			expectErr:  true,
			wantPeriod: 2,
			wantUnit:   "M",
			wantStatus: "AutoRenewal",
		},
		{
			name:       "bss-success-control",
			id:         "gw-unit-success",
			bssStatus:  http.StatusOK,
			bssBody:    "",
			wantPeriod: 3,
			wantUnit:   "Y",
			wantStatus: "AutoRenewal",
		},
	}

	byID := make(map[string]scenario, len(scenarios))
	for _, sc := range scenarios {
		byID[sc.id] = sc
	}
	// The isolated fixture records every request so a failure shows whether the
	// public Read reached APIG and BSS instead of only reporting a missing call.
	var callsMutex sync.Mutex
	apigCalls := 0
	var bssCalls []string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if strings.Contains(r.URL.Path, "/gateways/") {
			callsMutex.Lock()
			apigCalls++
			callsMutex.Unlock()
			// Keep the gateway Subscription so the public Read refreshes the renewal settings.
			_, _ = w.Write([]byte(`{"data":{"chargeType":"PREPAY"}}`))
			return
		}
		rawBody, _ := io.ReadAll(r.Body)
		_ = r.ParseForm()
		instanceIDs := r.Form.Get("InstanceIDs")
		if instanceIDs == "" {
			instanceIDs = r.URL.Query().Get("InstanceIDs")
		}
		if instanceIDs == "" {
			// The form value may arrive in the body depending on the signing path.
			params, _ := url.ParseQuery(string(rawBody))
			instanceIDs = params.Get("InstanceIDs")
		}
		if instanceIDs == "" {
			for _, pair := range strings.Split(r.URL.RawQuery+"&"+string(rawBody), "&") {
				if strings.Contains(pair, "InstanceIDs=") {
					instanceIDs = pair[strings.Index(pair, "InstanceIDs=")+len("InstanceIDs="):]
					break
				}
			}
		}
		callsMutex.Lock()
		bssCalls = append(bssCalls, fmt.Sprintf("%s %s rawQuery=%q body=%q ids=%q", r.Method, r.URL.Path, r.URL.RawQuery, string(rawBody), instanceIDs))
		callsMutex.Unlock()
		sc, ok := byID[instanceIDs]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"Code":"InvalidInstanceIDs"}`))
			return
		}
		body := sc.bssBody
		if body == "" {
			body = fmt.Sprintf(`{"Data":{"InstanceList":[{"InstanceID":"%s","RenewStatus":"AutoRenewal","RenewalDuration":3,"RenewalDurationUnit":"Y"}]}}`, instanceIDs)
		}
		w.WriteHeader(sc.bssStatus)
		_, _ = w.Write([]byte(body))
	}))
	t.Cleanup(server.Close)

	credential, err := credentials.NewCredential(new(credentials.Config).
		SetType("access_key").SetAccessKeyId("test-key").SetAccessKeySecret("test-secret"))
	if err != nil {
		t.Fatal(err)
	}
	endpoints := new(sync.Map)
	endpoint := strings.TrimPrefix(server.URL, "http://")
	t.Setenv("NO_PROXY", endpoint)
	config := &connectivity.Config{
		AccessKey: "test-key", SecretKey: "test-secret", Credential: credential,
		RegionId: "cn-hangzhou", AccountType: "test", Protocol: "http",
		Endpoints: endpoints, SignVersion: new(sync.Map), SkipRegionValidation: true,
	}
	client, err := config.Client()
	if err != nil {
		t.Fatal(err)
	}
	endpoints.Store("apig", endpoint)
	endpoints.Store("bssopenapi", endpoint)

	// The fixture must actually intercept the client: a silently unused endpoint
	// would send the test to the real API instead of the isolated server.
	probeResponse, probeErr := client.RoaGet("APIG", "2024-03-27", "/v1/gateways/fixture-probe", nil, nil, nil)
	callsMutex.Lock()
	probeCalls := apigCalls
	callsMutex.Unlock()
	if probeErr != nil || probeCalls != 1 {
		t.Fatalf("the isolated APIG fixture must intercept the client, got response=%v err=%v apigCalls=%d", probeResponse, probeErr, probeCalls)
	}

	for _, sc := range scenarios {
		t.Run(sc.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, resourceAliCloudApigGateway().Schema, map[string]interface{}{
				"gateway_name":        "tf-testacc-apig-gateway",
				"payment_type":        "Subscription",
				"renewal_status":      "AutoRenewal",
				"renew_period":        2,
				"renewal_period_unit": "M",
			})
			d.SetId(sc.id)

			readErr := resourceAliCloudApigGatewayRead(d, client)

			callsMutex.Lock()
			apigCount := apigCalls
			bssSummary := strings.Join(bssCalls, "; ")
			callsMutex.Unlock()

			if sc.expectErr {
				if readErr == nil {
					t.Fatalf("expected the public Read to fail when BssOpenApi fails for %s (id after read=%q, apig calls=%d, bss calls=%s)", sc.id, d.Id(), apigCount, bssSummary)
				}
				if NotFoundError(readErr) {
					t.Fatalf("expected the BssOpenApi failure not to be reported as NotFound, got %v", readErr)
				}
				if !strings.Contains(readErr.Error(), "QueryAvailableInstances") {
					t.Fatalf("expected the original BssOpenApi error to be preserved, got %v (apig calls=%d, bss calls=%s)", readErr, apigCount, bssSummary)
				}
				if got := d.Get("renew_period").(int); got != sc.wantPeriod {
					t.Fatalf("expected the known renewal period %d to be kept, got %d", sc.wantPeriod, got)
				}
				if got := d.Get("renewal_period_unit").(string); got != sc.wantUnit {
					t.Fatalf("expected the known renewal period unit %q to be kept, got %q", sc.wantUnit, got)
				}
				if got := d.Get("renewal_status").(string); got != sc.wantStatus {
					t.Fatalf("expected the known renewal status %q to be kept, got %q", sc.wantStatus, got)
				}
				return
			}

			if readErr != nil {
				t.Fatalf("expected the public Read to succeed for %s, got %v (id after read=%q, apig calls=%d, bss calls=%s)", sc.id, readErr, d.Id(), apigCount, bssSummary)
			}
			if got := d.Get("renew_period").(int); got != sc.wantPeriod {
				t.Fatalf("expected the renewal period to be refreshed to %d, got %d", sc.wantPeriod, got)
			}
			if got := d.Get("renewal_period_unit").(string); got != sc.wantUnit {
				t.Fatalf("expected the renewal period unit to be refreshed to %q, got %q", sc.wantUnit, got)
			}
			if got := d.Get("renewal_status").(string); got != sc.wantStatus {
				t.Fatalf("expected the renewal status to be refreshed to %q, got %q", sc.wantStatus, got)
			}
		})
	}

	callsMutex.Lock()
	defer callsMutex.Unlock()
	if len(bssCalls) == 0 {
		t.Fatalf("expected the public Read to call BssOpenApi at least once, observed none")
	}
}

// TestUnitAliCloudApigGatewayRenewalInstancePagination covers the unfiltered
// defense path of the BssOpenApi collection query: the entry is matched by
// InstanceID, so every reported page must be walked before concluding that the
// gateway is not visible yet.
func TestUnitAliCloudApigGatewayRenewalInstancePagination(t *testing.T) {
	type pageSpec struct {
		totalCount int
		instances  []interface{}
	}
	var mutex sync.Mutex
	pagesByID := make(map[string]map[int]pageSpec)
	// attemptsByID keeps one page sequence per walk attempt, so assertions stay
	// independent of how many times the retry loop re-walks the list.
	attemptsByID := make(map[string][][]int)
	malformedIDs := make(map[string]bool)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := r.ParseForm(); err != nil {
			t.Errorf("unexpected BssOpenApi request: %v", err)
			return
		}
		id := r.Form.Get("InstanceIDs")
		page, _ := strconv.Atoi(r.Form.Get("PageNum"))

		mutex.Lock()
		if page == 1 || len(attemptsByID[id]) == 0 {
			attemptsByID[id] = append(attemptsByID[id], []int{})
		}
		last := len(attemptsByID[id]) - 1
		attemptsByID[id][last] = append(attemptsByID[id][last], page)
		malformed := malformedIDs[id]
		spec := pagesByID[id][page]
		mutex.Unlock()

		if malformed {
			_, _ = w.Write([]byte(`{"Data":{"InstanceList":"not-a-list"}}`))
			return
		}
		instances := spec.instances
		if instances == nil {
			instances = []interface{}{}
		}
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"Success": true,
			"Data": map[string]interface{}{
				"TotalCount":   spec.totalCount,
				"InstanceList": instances,
			},
		})
	}))
	t.Cleanup(server.Close)

	credential, err := credentials.NewCredential(new(credentials.Config).
		SetType("access_key").SetAccessKeyId("test-key").SetAccessKeySecret("test-secret"))
	if err != nil {
		t.Fatal(err)
	}
	endpoints := new(sync.Map)
	endpoint := strings.TrimPrefix(server.URL, "http://")
	t.Setenv("NO_PROXY", endpoint)
	config := &connectivity.Config{
		AccessKey: "test-key", SecretKey: "test-secret", Credential: credential,
		RegionId: "cn-hangzhou", AccountType: "test", Protocol: "http",
		Endpoints: endpoints, SignVersion: new(sync.Map), SkipRegionValidation: true,
	}
	client, err := config.Client()
	if err != nil {
		t.Fatal(err)
	}
	endpoints.Store("apig", endpoint)
	endpoints.Store("bssopenapi", endpoint)
	service := &ApigServiceV2{client}

	firstAttempt := func(t *testing.T, id string) []int {
		t.Helper()
		mutex.Lock()
		defer mutex.Unlock()
		if len(attemptsByID[id]) == 0 {
			return nil
		}
		return append([]int(nil), attemptsByID[id][0]...)
	}

	t.Run("target on a later page is found", func(t *testing.T) {
		id := "gw-page-target"
		mutex.Lock()
		pagesByID[id] = map[int]pageSpec{
			1: {totalCount: apigGatewayRenewalPageSize + 1, instances: []interface{}{map[string]interface{}{"InstanceID": "gw-other-page1"}}},
			2: {totalCount: apigGatewayRenewalPageSize + 1, instances: []interface{}{map[string]interface{}{"InstanceID": id, "RenewStatus": "AutoRenewal"}}},
		}
		mutex.Unlock()

		instance, err := service.DescribeGatewayQueryAvailableInstancesByID(id, 30*time.Second)
		if err != nil {
			t.Fatalf("expected the paged lookup to find the target, got %v", err)
		}
		if fmt.Sprint(instance["InstanceID"]) != id {
			t.Fatalf("expected instance %q, got %v", id, instance["InstanceID"])
		}
		if requested := firstAttempt(t, id); len(requested) != 2 || requested[0] != 1 || requested[1] != 2 {
			t.Fatalf("expected PageNum 1 then 2 exactly once, got %v", requested)
		}
	})

	t.Run("target missing from every page is retried", func(t *testing.T) {
		id := "gw-page-missing"
		mutex.Lock()
		pagesByID[id] = map[int]pageSpec{
			1: {totalCount: apigGatewayRenewalPageSize + 1, instances: []interface{}{map[string]interface{}{"InstanceID": "gw-other-p1"}}},
			2: {totalCount: apigGatewayRenewalPageSize + 1, instances: []interface{}{map[string]interface{}{"InstanceID": "gw-other-p2"}}},
		}
		mutex.Unlock()

		if _, err := service.DescribeGatewayQueryAvailableInstancesByID(id, 3*time.Second); err == nil {
			t.Fatal("expected an error when the target is missing from every page")
		} else if NotFoundError(err) {
			t.Fatalf("expected the paged miss not to be reported as NotFound, got %v", err)
		}

		requested := firstAttempt(t, id)
		seen := map[int]bool{}
		for _, page := range requested {
			seen[page] = true
		}
		if !seen[1] || !seen[2] {
			t.Fatalf("expected the walk to request every reported page, got %v", requested)
		}
	})

	t.Run("malformed page response fails without paging", func(t *testing.T) {
		id := "gw-page-malformed"
		mutex.Lock()
		malformedIDs[id] = true
		mutex.Unlock()

		if _, err := service.DescribeGatewayQueryAvailableInstancesByID(id, 3*time.Second); err == nil {
			t.Fatal("expected an error for a malformed InstanceList")
		} else if NotFoundError(err) {
			t.Fatalf("expected the malformed response not to be reported as NotFound, got %v", err)
		}
		requested := firstAttempt(t, id)
		for _, page := range requested {
			if page != 1 {
				t.Fatalf("expected a malformed first page to stop the walk, got %v", requested)
			}
		}
	})

	t.Run("target past the fifty page mark is still found", func(t *testing.T) {
		id := "gw-page-deep-target"
		const lastPage = 51
		mutex.Lock()
		pages := make(map[int]pageSpec, lastPage)
		for page := 1; page <= lastPage; page++ {
			entry := map[string]interface{}{"InstanceID": fmt.Sprintf("gw-other-%d", page)}
			if page == lastPage {
				entry = map[string]interface{}{"InstanceID": id, "RenewStatus": "ManualRenewal"}
			}
			pages[page] = pageSpec{totalCount: lastPage, instances: []interface{}{entry}}
		}
		pagesByID[id] = pages
		mutex.Unlock()

		instance, err := service.DescribeGatewayQueryAvailableInstancesByID(id, 30*time.Second)
		if err != nil {
			t.Fatalf("expected the deep lookup to find the target, got %v", err)
		}
		if fmt.Sprint(instance["InstanceID"]) != id {
			t.Fatalf("expected instance %q, got %v", id, instance["InstanceID"])
		}
		requested := firstAttempt(t, id)
		if len(requested) != lastPage {
			t.Fatalf("expected the walk to reach page %d, got %v", lastPage, requested)
		}
		for index, page := range requested {
			if page != index+1 {
				t.Fatalf("expected a contiguous walk up to page %d, got %v", lastPage, requested)
			}
		}
	})

	t.Run("walk stops when the reported total is reached", func(t *testing.T) {
		id := "gw-page-total-boundary"
		fullPage := make([]interface{}, 0, apigGatewayRenewalPageSize)
		for index := 0; index < apigGatewayRenewalPageSize; index++ {
			fullPage = append(fullPage, map[string]interface{}{"InstanceID": fmt.Sprintf("gw-other-full-%d", index)})
		}
		mutex.Lock()
		pagesByID[id] = map[int]pageSpec{
			1: {totalCount: apigGatewayRenewalPageSize + 1, instances: fullPage},
			2: {totalCount: apigGatewayRenewalPageSize + 1, instances: []interface{}{map[string]interface{}{"InstanceID": "gw-other-b2"}}},
		}
		mutex.Unlock()

		if _, err := service.DescribeGatewayQueryAvailableInstancesByID(id, 2*time.Second); err == nil {
			t.Fatal("expected an error when the target is missing from the reported result set")
		}
		if requested := firstAttempt(t, id); len(requested) != 2 || requested[0] != 1 || requested[1] != 2 {
			t.Fatalf("expected the walk to stop at the reported total after pages 1 and 2, got %v", requested)
		}
	})

	t.Run("stale small TotalCount does not cut the walk short", func(t *testing.T) {
		id := "gw-page-stale-total"
		mutex.Lock()
		pagesByID[id] = map[int]pageSpec{
			1: {totalCount: 1, instances: []interface{}{map[string]interface{}{"InstanceID": "gw-other-s1"}}},
			2: {totalCount: 1, instances: []interface{}{map[string]interface{}{"InstanceID": id, "RenewStatus": "ManualRenewal"}}},
		}
		mutex.Unlock()

		instance, err := service.DescribeGatewayQueryAvailableInstancesByID(id, 30*time.Second)
		if err != nil {
			t.Fatalf("expected a stale TotalCount not to hide the target, got %v", err)
		}
		if fmt.Sprint(instance["InstanceID"]) != id {
			t.Fatalf("expected instance %q, got %v", id, instance["InstanceID"])
		}
		if requested := firstAttempt(t, id); len(requested) != 2 || requested[0] != 1 || requested[1] != 2 {
			t.Fatalf("expected the walk to continue past the stale total, got %v", requested)
		}
	})

	t.Run("walk stops on an empty page when TotalCount is absent", func(t *testing.T) {
		id := "gw-page-empty-boundary"
		mutex.Lock()
		pagesByID[id] = map[int]pageSpec{
			1: {totalCount: 0, instances: []interface{}{map[string]interface{}{"InstanceID": "gw-other-e1"}}},
		}
		mutex.Unlock()

		if _, err := service.DescribeGatewayQueryAvailableInstancesByID(id, 2*time.Second); err == nil {
			t.Fatal("expected an error when the target is missing from the reported result set")
		}
		if requested := firstAttempt(t, id); len(requested) != 2 || requested[0] != 1 || requested[1] != 2 {
			t.Fatalf("expected the walk to end on the empty second page, got %v", requested)
		}
	})
}

// TestUnitAliCloudApigGatewayRenewalTotalCountParsing pins the decoding of the
// BssOpenApi TotalCount value: the SDK reads response numbers with UseNumber, so
// in production the value is a json.Number instead of a float64.
func TestUnitAliCloudApigGatewayRenewalTotalCountParsing(t *testing.T) {
	cases := []struct {
		name     string
		set      bool
		raw      interface{}
		expected int
	}{
		{name: "json.Number from the SDK decoder", set: true, raw: json.Number("301"), expected: 301},
		{name: "float64", set: true, raw: float64(12), expected: 12},
		{name: "int", set: true, raw: 7, expected: 7},
		{name: "numeric string", set: true, raw: " 42 ", expected: 42},
		{name: "malformed string", set: true, raw: "many", expected: 0},
		{name: "missing value", set: false, expected: 0},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			data := map[string]interface{}{}
			if c.set {
				data["TotalCount"] = c.raw
			}
			response := map[string]interface{}{"Data": data}
			if got := apigGatewayRenewalTotalCount(response); got != c.expected {
				t.Fatalf("expected TotalCount %d, got %d", c.expected, got)
			}
		})
	}
}

// TestAccAliCloudApigGateway_renewalRequiresSubscription pins the user-visible
// behaviour of the renewal attributes on a non-prepaid gateway: the configuration
// must be rejected at plan time, before the gateway is created, instead of being
// applied while nothing in the cloud changes.
func TestAccAliCloudApigGateway_renewalRequiresSubscription(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: `
resource "alicloud_apig_gateway" "default" {
  gateway_name = "tf-testacc-apig-gateway-renewal-guard"
  spec         = "apigw.small.x1"
  payment_type = "PayAsYouGo"

  vpc {
    vpc_id = "vpc-bp1apiggatewayrenewalguard"
  }

  vswitch {
    vswitch_id = "vsw-bp1apiggatewayrenewalguard"
  }

  zone_config {
    select_option = "Auto"
  }

  renewal_status      = "AutoRenewal"
  renew_period        = 1
  renewal_period_unit = "M"
}
`,
				ExpectError: regexp.MustCompile(`can only be set when payment_type is Subscription`),
			},
		},
	})
}

func TestUnitAliCloudApigGatewayRenewalPaymentTypeGuard(t *testing.T) {
	renewalConfig := func(renewPeriod, renewalPeriodUnit, renewalStatus cty.Value) cty.Value {
		return cty.ObjectVal(map[string]cty.Value{
			"renew_period":        renewPeriod,
			"renewal_period_unit": renewalPeriodUnit,
			"renewal_status":      renewalStatus,
		})
	}
	noRenewalConfig := renewalConfig(cty.NullVal(cty.Number), cty.NullVal(cty.String), cty.NullVal(cty.String))

	// State carries the renewal values of every Subscription gateway because Read
	// populates them, so only what the configuration sets may be rejected. The raw
	// config is the only source that distinguishes the two.
	configCases := []struct {
		name      string
		rawConfig cty.Value
		want      []string
	}{
		{name: "absent raw config", rawConfig: cty.NilVal},
		{name: "null raw config", rawConfig: cty.NullVal(cty.String)},
		{name: "renewal values only in state", rawConfig: noRenewalConfig},
		{name: "renew_period configured", rawConfig: renewalConfig(cty.NumberIntVal(1), cty.NullVal(cty.String), cty.NullVal(cty.String)), want: []string{"renew_period"}},
		{name: "renew_period zero configured", rawConfig: renewalConfig(cty.NumberIntVal(0), cty.NullVal(cty.String), cty.NullVal(cty.String)), want: []string{"renew_period"}},
		{name: "renewal_period_unit configured", rawConfig: renewalConfig(cty.NullVal(cty.Number), cty.StringVal("M"), cty.NullVal(cty.String)), want: []string{"renewal_period_unit"}},
		{name: "renewal_status configured", rawConfig: renewalConfig(cty.NullVal(cty.Number), cty.NullVal(cty.String), cty.StringVal("AutoRenewal")), want: []string{"renewal_status"}},
		{name: "every renewal attribute configured", rawConfig: renewalConfig(cty.NumberIntVal(2), cty.StringVal("M"), cty.StringVal("AutoRenewal")), want: []string{"renew_period", "renewal_period_unit", "renewal_status"}},
		{name: "unknown value is not counted as configured", rawConfig: renewalConfig(cty.UnknownVal(cty.Number), cty.NullVal(cty.String), cty.NullVal(cty.String))},
	}
	for _, c := range configCases {
		t.Run("config presence: "+c.name, func(t *testing.T) {
			got := apigGatewayRenewalAttributesInConfig(c.rawConfig)
			if len(got) != len(c.want) {
				t.Fatalf("expected the configured renewal attributes %v, got %v", c.want, got)
			}
			for index := range got {
				if got[index] != c.want[index] {
					t.Fatalf("expected the configured renewal attributes %v, got %v", c.want, got)
				}
			}
		})
	}

	guardCases := []struct {
		name        string
		paymentType string
		configured  []string
		wantErr     string
	}{
		{name: "subscription accepts every renewal attribute", paymentType: "Subscription", configured: []string{"renew_period", "renewal_period_unit", "renewal_status"}},
		{name: "pay-as-you-go without renewal attributes", paymentType: "PayAsYouGo"},
		{name: "pay-as-you-go with renewal_status", paymentType: "PayAsYouGo", configured: []string{"renewal_status"}, wantErr: `renewal_status can only be set when payment_type is Subscription, but payment_type is "PayAsYouGo"`},
		{name: "pay-as-you-go with renew_period", paymentType: "PayAsYouGo", configured: []string{"renew_period"}, wantErr: `renew_period can only be set when payment_type is Subscription, but payment_type is "PayAsYouGo"`},
		{name: "pay-as-you-go with renewal_period_unit", paymentType: "PayAsYouGo", configured: []string{"renewal_period_unit"}, wantErr: `renewal_period_unit can only be set when payment_type is Subscription, but payment_type is "PayAsYouGo"`},
		{name: "pay-as-you-go with every renewal attribute", paymentType: "PayAsYouGo", configured: []string{"renew_period", "renewal_period_unit", "renewal_status"}, wantErr: `renew_period, renewal_period_unit, renewal_status can only be set when payment_type is Subscription, but payment_type is "PayAsYouGo"`},
	}
	for _, c := range guardCases {
		t.Run(c.name, func(t *testing.T) {
			err := checkApigGatewayRenewalPaymentType(c.paymentType, c.configured)
			if c.wantErr == "" {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				return
			}
			if err == nil {
				t.Fatalf("expected the rejection %q, got no error", c.wantErr)
			}
			if err.Error() != c.wantErr {
				t.Fatalf("expected the rejection %q, got %q", c.wantErr, err.Error())
			}
		})
	}
}

// TestUnitAliCloudApigGatewayRenewalRejectedBeforeWrite drives the public Update
// with a nil client: reaching any API would panic, so a returned error proves the
// non-Subscription rejection happens before the first call instead of after a
// partial apply.
func TestUnitAliCloudApigGatewayRenewalRejectedBeforeWrite(t *testing.T) {
	cases := []struct {
		name   string
		config map[string]interface{}
	}{
		{name: "renewal_status", config: map[string]interface{}{"renewal_status": "AutoRenewal"}},
		{name: "renew_period", config: map[string]interface{}{"renew_period": 1}},
		{name: "renewal_period_unit", config: map[string]interface{}{"renewal_period_unit": "M"}},
		{name: "every renewal attribute", config: map[string]interface{}{"renewal_status": "AutoRenewal", "renew_period": 1, "renewal_period_unit": "M"}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			raw := map[string]interface{}{
				"gateway_name": "tf-testacc-apig-gateway",
				"payment_type": "PayAsYouGo",
			}
			for key, value := range c.config {
				raw[key] = value
			}
			d := schema.TestResourceDataRaw(t, resourceAliCloudApigGateway().Schema, raw)
			d.SetId("gw-unit-pay-as-you-go")

			var updateErr error
			func() {
				defer func() {
					if recovered := recover(); recovered != nil {
						t.Fatalf("expected the renewal rejection before any API call, but the update reached the client: %v", recovered)
					}
				}()
				updateErr = resourceAliCloudApigGatewayUpdate(d, nil)
			}()
			if updateErr == nil {
				t.Fatal("expected an explicit error for a non-Subscription renewal configuration instead of a silent no-op")
			}
			if !strings.Contains(updateErr.Error(), "can only be set when payment_type is Subscription") {
				t.Fatalf("expected the error to name the Subscription requirement, got %v", updateErr)
			}
		})
	}
}

// TestUnitAliCloudApigGatewaySetRenewalRequest captures the serialized BssOpenApi
// SetRenewal request, so the renewal attributes are proven to reach the wire under
// the documented parameter names and not only to exist in the request map.
func TestUnitAliCloudApigGatewaySetRenewalRequest(t *testing.T) {
	const gatewayID = "gw-unit-set-renewal"
	var mutex sync.Mutex
	setRenewalRequests := make([]map[string]string, 0, 1)
	apigWriteCalls := make([]string, 0, 1)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if strings.Contains(r.URL.Path, "/gateways/") {
			if r.Method != http.MethodGet {
				mutex.Lock()
				apigWriteCalls = append(apigWriteCalls, r.Method+" "+r.URL.Path)
				mutex.Unlock()
			}
			_, _ = w.Write([]byte(`{"data":{"chargeType":"PREPAY"}}`))
			return
		}
		if err := r.ParseForm(); err != nil {
			t.Errorf("unexpected BssOpenApi request: %v", err)
			return
		}
		if r.Form.Get("Action") == "SetRenewal" || r.Form.Get("RenewalStatus") != "" {
			captured := make(map[string]string)
			for _, key := range []string{"InstanceIDs", "RenewalStatus", "RenewalPeriod", "RenewalPeriodUnit", "ProductCode", "SubscriptionType"} {
				captured[key] = r.Form.Get(key)
			}
			mutex.Lock()
			setRenewalRequests = append(setRenewalRequests, captured)
			mutex.Unlock()
			_, _ = w.Write([]byte(`{"Success":true,"Code":"Success"}`))
			return
		}
		// QueryAvailableInstances reports the requested settings so the convergence
		// wait finishes on its first poll.
		_, _ = w.Write([]byte(`{"Success":true,"Data":{"TotalCount":1,"InstanceList":[{"InstanceID":"` + gatewayID + `","RenewStatus":"AutoRenewal","RenewalDuration":2,"RenewalDurationUnit":"M"}]}}`))
	}))
	t.Cleanup(server.Close)

	client := newApigGatewayUnitTestClient(t, server)

	d := schema.TestResourceDataRaw(t, resourceAliCloudApigGateway().Schema, map[string]interface{}{
		"gateway_name":        "tf-testacc-apig-gateway",
		"payment_type":        "Subscription",
		"renewal_status":      "AutoRenewal",
		"renew_period":        2,
		"renewal_period_unit": "M",
	})
	d.SetId(gatewayID)

	if err := setApigGatewayRenewal(d, client); err != nil {
		t.Fatalf("expected the Subscription renewal update to succeed, got %v", err)
	}

	mutex.Lock()
	defer mutex.Unlock()
	if len(setRenewalRequests) != 1 {
		t.Fatalf("expected exactly one SetRenewal request, got %d (%v)", len(setRenewalRequests), setRenewalRequests)
	}
	expected := map[string]string{
		"InstanceIDs":       gatewayID,
		"RenewalStatus":     "AutoRenewal",
		"RenewalPeriod":     "2",
		"RenewalPeriodUnit": "M",
		"ProductCode":       "apigateway",
		"SubscriptionType":  "Subscription",
	}
	for key, value := range expected {
		if got := setRenewalRequests[0][key]; got != value {
			t.Fatalf("expected the serialized SetRenewal %s to be %q, got %q (request=%v)", key, value, got, setRenewalRequests[0])
		}
	}
	if len(apigWriteCalls) != 0 {
		t.Fatalf("expected the renewal path not to issue APIG write calls, got %v", apigWriteCalls)
	}
}

// TestUnitAliCloudApigGatewayRenewalPayAsYouGoGuard keeps the explicit rejection
// and the untouched gateway apart: the same PayAsYouGo gateway is accepted when no
// renewal attribute changes and fails without a write when one does.
func TestUnitAliCloudApigGatewayRenewalPayAsYouGoGuard(t *testing.T) {
	var mutex sync.Mutex
	apigWriteCalls := make([]string, 0, 1)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method != http.MethodGet {
			mutex.Lock()
			apigWriteCalls = append(apigWriteCalls, r.Method+" "+r.URL.Path)
			mutex.Unlock()
		}
		_, _ = w.Write([]byte(`{"data":{"chargeType":"POSTPAY"}}`))
	}))
	t.Cleanup(server.Close)

	client := newApigGatewayUnitTestClient(t, server)

	t.Run("a renewal change is rejected without a write", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, resourceAliCloudApigGateway().Schema, map[string]interface{}{
			"gateway_name":        "tf-testacc-apig-gateway",
			"payment_type":        "PayAsYouGo",
			"renewal_status":      "AutoRenewal",
			"renew_period":        1,
			"renewal_period_unit": "M",
		})
		d.SetId("gw-unit-pay-as-you-go")

		err := setApigGatewayRenewal(d, client)
		if err == nil {
			t.Fatal("expected an explicit error for a renewal change on a PayAsYouGo gateway instead of a silent no-op")
		}
		if !strings.Contains(err.Error(), "can only be managed for Subscription gateways") {
			t.Fatalf("expected the error to name the Subscription requirement, got %v", err)
		}
		mutex.Lock()
		defer mutex.Unlock()
		if len(apigWriteCalls) != 0 {
			t.Fatalf("expected the rejection to happen before any write API, got %v", apigWriteCalls)
		}
	})

	t.Run("a gateway without a renewal change is left alone", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, resourceAliCloudApigGateway().Schema, map[string]interface{}{
			"gateway_name": "tf-testacc-apig-gateway",
			"payment_type": "PayAsYouGo",
		})
		d.SetId("gw-unit-pay-as-you-go-plain")

		if err := setApigGatewayRenewal(d, client); err != nil {
			t.Fatalf("expected a PayAsYouGo gateway without renewal attributes to be accepted, got %v", err)
		}
	})
}

func newApigGatewayUnitTestClient(t *testing.T, server *httptest.Server) *connectivity.AliyunClient {
	t.Helper()
	credential, err := credentials.NewCredential(new(credentials.Config).
		SetType("access_key").SetAccessKeyId("test-key").SetAccessKeySecret("test-secret"))
	if err != nil {
		t.Fatal(err)
	}
	endpoints := new(sync.Map)
	endpoint := strings.TrimPrefix(server.URL, "http://")
	t.Setenv("NO_PROXY", endpoint)
	config := &connectivity.Config{
		AccessKey: "test-key", SecretKey: "test-secret", Credential: credential,
		RegionId: "cn-hangzhou", AccountType: "test", Protocol: "http",
		Endpoints: endpoints, SignVersion: new(sync.Map), SkipRegionValidation: true,
	}
	client, err := config.Client()
	if err != nil {
		t.Fatal(err)
	}
	endpoints.Store("apig", endpoint)
	endpoints.Store("bssopenapi", endpoint)
	return client
}
