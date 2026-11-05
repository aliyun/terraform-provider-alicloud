package apig_test

import (
	"fmt"
	"testing"

	acctest "github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	apig "github.com/aliyun/terraform-provider-alicloud/alicloud/service/apig"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test Apig Gateway. >>> Resource test cases, automatically generated.
// Case 资源组接入_prepay 9249
func TestAccAliCloudApigGateway_basic9249(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_apig_gateway.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudApigGatewayMap9249)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return apig.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeApigGateway")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sapiggateway%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudApigGatewayBasicDependence9249)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
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
						"resource_group_id": acctest.CHECKSET,
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
					"tags": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": acctest.REMOVEKEY,
						"tags.For":     acctest.REMOVEKEY,
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
	"status":      acctest.CHECKSET,
	"create_time": acctest.CHECKSET,
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
	ra := acctest.ResourceAttrInit(resourceId, AlicloudApigGatewayMap9246)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return apig.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeApigGateway")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sapiggateway%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudApigGatewayBasicDependence9246)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
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
						"resource_group_id": acctest.CHECKSET,
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
						"resource_group_id": acctest.CHECKSET,
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
					"tags": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": acctest.REMOVEKEY,
						"tags.For":     acctest.REMOVEKEY,
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
	"status":      acctest.CHECKSET,
	"create_time": acctest.CHECKSET,
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
	ra := acctest.ResourceAttrInit(resourceId, AlicloudApigGatewayMap10903)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return apig.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeApigGateway")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccapig%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudApigGatewayBasicDependence10903)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
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
						"resource_group_id": acctest.CHECKSET,
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
						"resource_group_id": acctest.CHECKSET,
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
					"tags": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": acctest.REMOVEKEY,
						"tags.For":     acctest.REMOVEKEY,
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
	"status":      acctest.CHECKSET,
	"create_time": acctest.CHECKSET,
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

// Test Apig Gateway. <<< Resource test cases, automatically generated.

// Test Apig Gateway. gateway_edition coverage.
func TestAccAliCloudApigGateway_basic_edition(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_apig_gateway.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudApigGatewayEditionMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return apig.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeApigGateway")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sapiggwedition%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudApigGatewayEditionDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
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
	"status":      acctest.CHECKSET,
	"create_time": acctest.CHECKSET,
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

// collection-order: bypass zones -- ForceNew TypeList (zones.ForceNew=true at resource_alicloud_apig_gateway.go:291); any member reorder forces replacement, so memberOrder convergence is impossible by schema contract. ACC TestAccAliCloudApigGateway_memberOrder step2 confirmed -/+ destroy and create replacement (zones reorder reassigns the server-coupled primary vswitch whose vswitch_id is also ForceNew, forcing full replacement). Per-property adjudication, not a uniform historical bypass.
