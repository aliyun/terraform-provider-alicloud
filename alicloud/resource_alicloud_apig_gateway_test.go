package alicloud

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
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
	name := fmt.Sprintf("tfaccapig%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudApigGatewayBasicDependence10903)
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
	name := fmt.Sprintf("tf-testacc%sapiggateway%d", defaultRegionToTest, rand)
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
