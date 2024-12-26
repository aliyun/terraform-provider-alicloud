package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
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
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
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
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
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

// Test Apig Gateway. <<< Resource test cases, automatically generated.
