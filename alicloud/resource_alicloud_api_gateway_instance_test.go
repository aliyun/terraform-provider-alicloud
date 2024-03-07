package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test ApiGateway Instance. >>> Resource test cases, automatically generated.
// Case 对接Terraform_NORMAL 5800
func TestAccAliCloudApiGatewayInstance_basic5800(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_api_gateway_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudApiGatewayInstanceMap5800)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApiGatewayServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeApiGatewayInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sapigate%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudApiGatewayInstanceBasicDependence5800)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name": name,
					"instance_spec": "api.s1.small",
					"https_policy":  "HTTPS2_TLS1_0",
					"payment_type":  "PayAsYouGo",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name": name,
						"instance_spec": "api.s1.small",
						"https_policy":  "HTTPS2_TLS1_0",
						"payment_type":  "PayAsYouGo",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"https_policy": "HTTPS2_TLS1_2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"https_policy": "HTTPS2_TLS1_2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name":           name + "_update",
					"instance_spec":           "api.s1.small",
					"https_policy":            "HTTPS2_TLS1_0",
					"zone_id":                 "cn-hangzhou-MAZ6",
					"payment_type":            "PayAsYouGo",
					"instance_type":           "normal",
					"user_vpc_id":             "1706841299",
					"egress_ipv6_enable":      "true",
					"support_ipv6":            "true",
					"vpc_slb_intranet_enable": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name":           name + "_update",
						"instance_spec":           "api.s1.small",
						"https_policy":            "HTTPS2_TLS1_0",
						"zone_id":                 "cn-hangzhou-MAZ6",
						"payment_type":            "PayAsYouGo",
						"instance_type":           "normal",
						"user_vpc_id":             CHECKSET,
						"egress_ipv6_enable":      "true",
						"support_ipv6":            "true",
						"vpc_slb_intranet_enable": "true",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"egress_ipv6_enable", "support_ipv6", "user_vpc_id", "vpc_slb_intranet_enable"},
			},
		},
	})
}

var AlicloudApiGatewayInstanceMap5800 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
}

func AlicloudApiGatewayInstanceBasicDependence5800(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case 对接Terraform_PREPAY 5806
func TestAccAliCloudApiGatewayInstance_basic5806(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_api_gateway_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudApiGatewayInstanceMap5806)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApiGatewayServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeApiGatewayInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sapigate%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudApiGatewayInstanceBasicDependence5806)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name": name,
					"instance_spec": "api.s1.small",
					"https_policy":  "HTTPS2_TLS1_0",
					"payment_type":  "Subscription",
					"pricing_cycle": "month",
					"duration":      "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name": name,
						"instance_spec": "api.s1.small",
						"https_policy":  "HTTPS2_TLS1_0",
						"payment_type":  "Subscription",
						"pricing_cycle": "month",
						"duration":      "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"https_policy": "HTTPS2_TLS1_2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"https_policy": "HTTPS2_TLS1_2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name":           name + "_update",
					"instance_spec":           "api.s1.small",
					"https_policy":            "HTTPS2_TLS1_0",
					"zone_id":                 "cn-hangzhou-MAZ6",
					"payment_type":            "Subscription",
					"user_vpc_id":             "1706841299",
					"egress_ipv6_enable":      "true",
					"support_ipv6":            "true",
					"vpc_slb_intranet_enable": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name":           name + "_update",
						"instance_spec":           "api.s1.small",
						"https_policy":            "HTTPS2_TLS1_0",
						"zone_id":                 "cn-hangzhou-MAZ6",
						"payment_type":            "Subscription",
						"user_vpc_id":             CHECKSET,
						"egress_ipv6_enable":      "true",
						"support_ipv6":            "true",
						"vpc_slb_intranet_enable": "true",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"egress_ipv6_enable", "support_ipv6", "user_vpc_id", "vpc_slb_intranet_enable", "pricing_cycle", "duration"},
			},
		},
	})
}

var AlicloudApiGatewayInstanceMap5806 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
}

func AlicloudApiGatewayInstanceBasicDependence5806(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case 对接Terraform_NORMAL 5800  twin
func TestAccAliCloudApiGatewayInstance_basic5800_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_api_gateway_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudApiGatewayInstanceMap5800)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApiGatewayServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeApiGatewayInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sapigatewayinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudApiGatewayInstanceBasicDependence5800)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name":           name,
					"instance_spec":           "api.s1.small",
					"https_policy":            "HTTPS2_TLS1_2",
					"zone_id":                 "cn-hangzhou-MAZ6",
					"payment_type":            "PayAsYouGo",
					"instance_type":           "normal",
					"user_vpc_id":             "1706841299",
					"egress_ipv6_enable":      "false",
					"support_ipv6":            "false",
					"vpc_slb_intranet_enable": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name":           name,
						"instance_spec":           "api.s1.small",
						"https_policy":            "HTTPS2_TLS1_2",
						"zone_id":                 "cn-hangzhou-MAZ6",
						"payment_type":            "PayAsYouGo",
						"instance_type":           "normal",
						"user_vpc_id":             CHECKSET,
						"egress_ipv6_enable":      "false",
						"support_ipv6":            "false",
						"vpc_slb_intranet_enable": "false",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"egress_ipv6_enable", "support_ipv6", "user_vpc_id", "vpc_slb_intranet_enable"},
			},
		},
	})
}

// Case 对接Terraform_PREPAY 5806  twin
func TestAccAliCloudApiGatewayInstance_basic5806_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_api_gateway_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudApiGatewayInstanceMap5806)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApiGatewayServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeApiGatewayInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sapigatewayinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudApiGatewayInstanceBasicDependence5806)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name":           name,
					"instance_spec":           "api.s1.small",
					"https_policy":            "HTTPS2_TLS1_2",
					"zone_id":                 "cn-hangzhou-MAZ6",
					"payment_type":            "Subscription",
					"user_vpc_id":             "1706841300",
					"egress_ipv6_enable":      "false",
					"support_ipv6":            "false",
					"vpc_slb_intranet_enable": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name":           name,
						"instance_spec":           "api.s1.small",
						"https_policy":            "HTTPS2_TLS1_2",
						"zone_id":                 "cn-hangzhou-MAZ6",
						"payment_type":            "Subscription",
						"user_vpc_id":             CHECKSET,
						"egress_ipv6_enable":      "false",
						"support_ipv6":            "false",
						"vpc_slb_intranet_enable": "false",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"egress_ipv6_enable", "support_ipv6", "user_vpc_id", "vpc_slb_intranet_enable"},
			},
		},
	})
}

// Test ApiGateway Instance. <<< Resource test cases, automatically generated.
