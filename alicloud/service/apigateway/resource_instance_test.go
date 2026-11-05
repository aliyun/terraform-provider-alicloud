package apigateway_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/apigateway"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test ApiGateway Instance. >>> Resource test cases, automatically generated.
// Case 对接Terraform_NORMAL 5800
func TestAccAliCloudApiGatewayInstance_basic5800(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_api_gateway_instance.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudApiGatewayInstanceMap5800)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return apigateway.NewApiGatewayServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeApiGatewayInstance")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sapigate%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudApiGatewayInstanceBasicDependence5800)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
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
					"zone_id":                 "cn-hangzhou-MAZ6(i,j,k)",
					"payment_type":            "PayAsYouGo",
					"instance_type":           "normal",
					"egress_ipv6_enable":      "true",
					"vpc_slb_intranet_enable": "true",
					"ipv6_enabled":            "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name":           name + "_update",
						"instance_spec":           "api.s1.small",
						"https_policy":            "HTTPS2_TLS1_0",
						"zone_id":                 "cn-hangzhou-MAZ6(i,j,k)",
						"payment_type":            "PayAsYouGo",
						"instance_type":           "normal",
						"egress_ipv6_enable":      "true",
						"vpc_slb_intranet_enable": "true",
						"ipv6_enabled":            "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_spec":    "api.s1.medium",
					"skip_wait_switch": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_spec": "api.s1.medium",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true, ImportStateVerifyIgnore: []string{"skip_wait_switch"},
			},
		},
	})
}

var AlicloudApiGatewayInstanceMap5800 = map[string]string{
	"status":      acctest.CHECKSET,
	"create_time": acctest.CHECKSET,
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
	ra := acctest.ResourceAttrInit(resourceId, AlicloudApiGatewayInstanceMap5806)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return apigateway.NewApiGatewayServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeApiGatewayInstance")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sapigate%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudApiGatewayInstanceBasicDependence5806)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name": name,
					"instance_spec": "api.s1.small",
					"https_policy":  "HTTPS2_TLS1_0",
					"payment_type":  "Subscription",
					"pricing_cycle": "month",
					"duration":      "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name": name,
						"instance_spec": "api.s1.small",
						"https_policy":  "HTTPS2_TLS1_0",
						"payment_type":  "Subscription",
						"pricing_cycle": "month",
						"duration":      "1",
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
					"zone_id":                 "cn-hangzhou-MAZ6(i,j,k)",
					"payment_type":            "Subscription",
					"egress_ipv6_enable":      "true",
					"vpc_slb_intranet_enable": "true",
					"ipv6_enabled":            "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name":           name + "_update",
						"instance_spec":           "api.s1.small",
						"https_policy":            "HTTPS2_TLS1_0",
						"zone_id":                 "cn-hangzhou-MAZ6(i,j,k)",
						"payment_type":            "Subscription",
						"egress_ipv6_enable":      "true",
						"vpc_slb_intranet_enable": "true",
						"ipv6_enabled":            "true",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"pricing_cycle", "duration"},
			},
		},
	})
}

var AlicloudApiGatewayInstanceMap5806 = map[string]string{
	"status":      acctest.CHECKSET,
	"create_time": acctest.CHECKSET,
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
	ra := acctest.ResourceAttrInit(resourceId, AlicloudApiGatewayInstanceMap5800)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return apigateway.NewApiGatewayServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeApiGatewayInstance")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sapigatewayinstance%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudApiGatewayInstanceBasicDependence5800)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name":           name,
					"instance_spec":           "api.s1.small",
					"https_policy":            "HTTPS2_TLS1_2",
					"zone_id":                 "cn-hangzhou-MAZ6(i,j,k)",
					"payment_type":            "PayAsYouGo",
					"instance_type":           "normal",
					"egress_ipv6_enable":      "false",
					"vpc_slb_intranet_enable": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name":           name,
						"instance_spec":           "api.s1.small",
						"https_policy":            "HTTPS2_TLS1_2",
						"zone_id":                 "cn-hangzhou-MAZ6(i,j,k)",
						"payment_type":            "PayAsYouGo",
						"instance_type":           "normal",
						"egress_ipv6_enable":      "false",
						"vpc_slb_intranet_enable": "false",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Case 对接Terraform_PREPAY 5806  twin
func TestAccAliCloudApiGatewayInstance_basic5806_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_api_gateway_instance.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudApiGatewayInstanceMap5806)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return apigateway.NewApiGatewayServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeApiGatewayInstance")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sapigatewayinstance%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudApiGatewayInstanceBasicDependence5806)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name":           name,
					"instance_spec":           "api.s1.small",
					"https_policy":            "HTTPS2_TLS1_2",
					"zone_id":                 "cn-hangzhou-MAZ6(i,j,k)",
					"payment_type":            "Subscription",
					"pricing_cycle":           "month",
					"duration":                "1",
					"egress_ipv6_enable":      "false",
					"vpc_slb_intranet_enable": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name":           name,
						"instance_spec":           "api.s1.small",
						"https_policy":            "HTTPS2_TLS1_2",
						"zone_id":                 "cn-hangzhou-MAZ6(i,j,k)",
						"payment_type":            "Subscription",
						"pricing_cycle":           "month",
						"duration":                "1",
						"egress_ipv6_enable":      "false",
						"vpc_slb_intranet_enable": "false",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"pricing_cycle", "duration"},
			},
		},
	})
}

// Test ApiGateway Instance. <<< Resource test cases, automatically generated.

var AlicloudApiGatewayVpcConnectInstanceCheckMap = map[string]string{
	"status":      acctest.CHECKSET,
	"create_time": acctest.CHECKSET,
}

func AlicloudApiGatewayVpcConnectInstanceBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alicloud_vpc" "vpc" {
  cidr_block = "172.16.0.0/12"
  vpc_name   = var.name
}

resource "alicloud_vswitch" "vswitch_1" {
  vpc_id       = alicloud_vpc.vpc.id
  cidr_block   = "172.16.0.0/16"
  zone_id      = "cn-hangzhou-j"
  vswitch_name = "${var.name}_1"
}

resource "alicloud_vswitch" "vswitch_2" {
  vpc_id       = alicloud_vpc.vpc.id
  cidr_block   = "172.17.0.0/16"
  zone_id      = "cn-hangzhou-k"
  vswitch_name = "${var.name}_2"
}

resource "alicloud_security_group" "security_group_1" {
  vpc_id              = alicloud_vpc.vpc.id
  security_group_name = "${var.name}_1"
}

resource "alicloud_security_group" "security_group_2" {
  vpc_id              = alicloud_vpc.vpc.id
  security_group_name = "${var.name}_2"
}

data "alicloud_account" "current" {
}
`, name)
}

func TestAccAliCloudApiGatewayVpcConnectInstance(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_api_gateway_instance.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudApiGatewayVpcConnectInstanceCheckMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return apigateway.NewApiGatewayServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeApiGatewayInstance")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sapigate%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudApiGatewayVpcConnectInstanceBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name": name,
					"instance_spec": "api.s1.small",
					"https_policy":  "HTTPS2_TLS1_0",
					"payment_type":  "PayAsYouGo",
					"instance_type": "vpc_connect",
					"instance_cidr": "192.168.0.0/16",
					"user_vpc_id":   "${alicloud_vpc.vpc.id}",
					"zone_vswitch_security_group": []map[string]interface{}{
						{
							"zone_id":        "${alicloud_vswitch.vswitch_1.zone_id}",
							"vswitch_id":     "${alicloud_vswitch.vswitch_1.id}",
							"cidr_block":     "${alicloud_vswitch.vswitch_1.cidr_block}",
							"security_group": "${alicloud_security_group.security_group_1.id}",
						},
						{
							"zone_id":        "${alicloud_vswitch.vswitch_2.zone_id}",
							"vswitch_id":     "${alicloud_vswitch.vswitch_2.id}",
							"cidr_block":     "${alicloud_vswitch.vswitch_2.cidr_block}",
							"security_group": "${alicloud_security_group.security_group_2.id}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name":       name,
						"instance_spec":       "api.s1.small",
						"https_policy":        "HTTPS2_TLS1_0",
						"payment_type":        "PayAsYouGo",
						"instance_cidr":       "192.168.0.0/16",
						"user_vpc_id":         acctest.CHECKSET,
						"connect_cidr_blocks": `["172.16.0.0/16","172.17.0.0/16"]`,
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
					"to_connect_vpc_ip_block": []map[string]interface{}{
						{
							"cidr_block": "10.1.0.0/24",
							"customized": "true",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"connect_cidr_blocks": `["172.16.0.0/16","172.17.0.0/16","10.1.0.0/24"]`,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"to_connect_vpc_ip_block": acctest.REMOVEKEY,
					"delete_vpc_ip_block":     "10.1.0.0/24",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"connect_cidr_blocks": `["172.16.0.0/16","172.17.0.0/16"]`,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name":           name + "_update",
					"https_policy":            "HTTPS2_TLS1_0",
					"egress_ipv6_enable":      "true",
					"vpc_slb_intranet_enable": "true",
					"ipv6_enabled":            "true",
					"delete_vpc_ip_block":     "10.1.0.0/24",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name":           name + "_update",
						"https_policy":            "HTTPS2_TLS1_0",
						"egress_ipv6_enable":      "true",
						"vpc_slb_intranet_enable": "true",
						"ipv6_enabled":            "true",
						"connect_cidr_blocks":     `["172.16.0.0/16","172.17.0.0/16"]`,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ingress_vpc_id":       "${alicloud_vpc.vpc.id}",
					"ingress_vpc_owner_id": "${data.alicloud_account.current.id}",
					"ingress_vswitch_id":   "${alicloud_vswitch.vswitch_1.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ingress_vpc_id":       acctest.CHECKSET,
						"ingress_vpc_owner_id": acctest.CHECKSET,
						"ingress_vswitch_id":   acctest.CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"to_connect_vpc_ip_block", "delete_vpc_ip_block"},
			},
		},
	})
}

func TestAccAliCloudApiGatewayVpcConnectInstance_twin1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_api_gateway_instance.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudApiGatewayVpcConnectInstanceCheckMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return apigateway.NewApiGatewayServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeApiGatewayInstance")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sapigate%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudApiGatewayVpcConnectInstanceBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name":           name,
					"instance_spec":           "api.s1.small",
					"https_policy":            "HTTPS2_TLS1_0",
					"payment_type":            "PayAsYouGo",
					"instance_type":           "vpc_connect",
					"instance_cidr":           "192.168.0.0/16",
					"egress_ipv6_enable":      "true",
					"user_vpc_id":             "${alicloud_vpc.vpc.id}",
					"vpc_slb_intranet_enable": "true",
					"ipv6_enabled":            "true",
					"zone_vswitch_security_group": []map[string]interface{}{
						{
							"zone_id":        "${alicloud_vswitch.vswitch_1.zone_id}",
							"vswitch_id":     "${alicloud_vswitch.vswitch_1.id}",
							"cidr_block":     "${alicloud_vswitch.vswitch_1.cidr_block}",
							"security_group": "${alicloud_security_group.security_group_1.id}",
						},
						{
							"zone_id":        "${alicloud_vswitch.vswitch_2.zone_id}",
							"vswitch_id":     "${alicloud_vswitch.vswitch_2.id}",
							"cidr_block":     "${alicloud_vswitch.vswitch_2.cidr_block}",
							"security_group": "${alicloud_security_group.security_group_2.id}",
						},
					},
					"delete_vpc_ip_block": "10.1.0.0/24",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name":           name,
						"instance_spec":           "api.s1.small",
						"https_policy":            "HTTPS2_TLS1_0",
						"payment_type":            "PayAsYouGo",
						"instance_cidr":           "192.168.0.0/16",
						"user_vpc_id":             acctest.CHECKSET,
						"egress_ipv6_enable":      "true",
						"vpc_slb_intranet_enable": "true",
						"ipv6_enabled":            "true",
						"connect_cidr_blocks":     `["172.16.0.0/16","172.17.0.0/16"]`,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"to_connect_vpc_ip_block", "delete_vpc_ip_block"},
			},
		},
	})
}

func TestAccAliCloudApiGatewayVpcConnectInstance_twin2(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_api_gateway_instance.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudApiGatewayVpcConnectInstanceCheckMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return apigateway.NewApiGatewayServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeApiGatewayInstance")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sapigate%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudApiGatewayVpcConnectInstanceBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name":           name,
					"instance_spec":           "api.s1.small",
					"https_policy":            "HTTPS2_TLS1_0",
					"payment_type":            "PayAsYouGo",
					"instance_type":           "vpc_connect",
					"instance_cidr":           "192.168.0.0/16",
					"egress_ipv6_enable":      "false",
					"user_vpc_id":             "${alicloud_vpc.vpc.id}",
					"vpc_slb_intranet_enable": "true",
					"ipv6_enabled":            "true",
					"zone_vswitch_security_group": []map[string]interface{}{
						{
							"zone_id":        "${alicloud_vswitch.vswitch_1.zone_id}",
							"vswitch_id":     "${alicloud_vswitch.vswitch_1.id}",
							"cidr_block":     "${alicloud_vswitch.vswitch_1.cidr_block}",
							"security_group": "${alicloud_security_group.security_group_1.id}",
						},
						{
							"zone_id":        "${alicloud_vswitch.vswitch_2.zone_id}",
							"vswitch_id":     "${alicloud_vswitch.vswitch_2.id}",
							"cidr_block":     "${alicloud_vswitch.vswitch_2.cidr_block}",
							"security_group": "${alicloud_security_group.security_group_2.id}",
						},
					},
					"to_connect_vpc_ip_block": []map[string]interface{}{
						{
							"cidr_block": "10.1.0.0/24",
							"customized": "true",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name":           name,
						"instance_spec":           "api.s1.small",
						"https_policy":            "HTTPS2_TLS1_0",
						"payment_type":            "PayAsYouGo",
						"instance_cidr":           "192.168.0.0/16",
						"user_vpc_id":             acctest.CHECKSET,
						"egress_ipv6_enable":      "false",
						"vpc_slb_intranet_enable": "true",
						"ipv6_enabled":            "true",
						"connect_cidr_blocks":     `["172.16.0.0/16","172.17.0.0/16","10.1.0.0/24"]`,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"to_connect_vpc_ip_block", "delete_vpc_ip_block"},
			},
		},
	})
}

// TestAccAliCloudApiGatewayInstance_zoneVswitchSecurityGroupOrder verifies that a
// reordered zone_vswitch_security_group list (same members, different order)
// converges after apply. zone_vswitch_security_group is ForceNew, so reordering
// triggers replacement; the post-apply plan is empty because the recreated
// instance reflects the reordered config. Covers TypeListOrderCoverage for this
// TypeList and exercises real modify coverage of its nested attributes.
func TestAccAliCloudApiGatewayInstance_zoneVswitchSecurityGroupOrder(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_api_gateway_instance.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudApiGatewayVpcConnectInstanceCheckMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return apigateway.NewApiGatewayServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeApiGatewayInstance")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sapigate%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudApiGatewayVpcConnectInstanceBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name": name,
					"instance_spec": "api.s1.small",
					"https_policy":  "HTTPS2_TLS1_0",
					"payment_type":  "PayAsYouGo",
					"instance_type": "vpc_connect",
					"instance_cidr": "192.168.0.0/16",
					"user_vpc_id":   "${alicloud_vpc.vpc.id}",
					"zone_vswitch_security_group": []map[string]interface{}{
						{
							"zone_id":        "${alicloud_vswitch.vswitch_1.zone_id}",
							"vswitch_id":     "${alicloud_vswitch.vswitch_1.id}",
							"cidr_block":     "${alicloud_vswitch.vswitch_1.cidr_block}",
							"security_group": "${alicloud_security_group.security_group_1.id}",
						},
						{
							"zone_id":        "${alicloud_vswitch.vswitch_2.zone_id}",
							"vswitch_id":     "${alicloud_vswitch.vswitch_2.id}",
							"cidr_block":     "${alicloud_vswitch.vswitch_2.cidr_block}",
							"security_group": "${alicloud_security_group.security_group_2.id}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name":       name,
						"instance_spec":       "api.s1.small",
						"https_policy":        "HTTPS2_TLS1_0",
						"payment_type":        "PayAsYouGo",
						"instance_cidr":       "192.168.0.0/16",
						"user_vpc_id":         acctest.CHECKSET,
						"connect_cidr_blocks": `["172.16.0.0/16","172.17.0.0/16"]`,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"zone_vswitch_security_group": []map[string]interface{}{
						{
							"zone_id":        "${alicloud_vswitch.vswitch_2.zone_id}",
							"vswitch_id":     "${alicloud_vswitch.vswitch_2.id}",
							"cidr_block":     "${alicloud_vswitch.vswitch_2.cidr_block}",
							"security_group": "${alicloud_security_group.security_group_2.id}",
						},
						{
							"zone_id":        "${alicloud_vswitch.vswitch_1.zone_id}",
							"vswitch_id":     "${alicloud_vswitch.vswitch_1.id}",
							"cidr_block":     "${alicloud_vswitch.vswitch_1.cidr_block}",
							"security_group": "${alicloud_security_group.security_group_1.id}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"connect_cidr_blocks": acctest.CHECKSET,
					}),
				),
			},
		},
	})
}

// TestAccAliCloudApiGatewayInstance_toConnectVpcIpBlockOrder verifies that a
// reordered to_connect_vpc_ip_block list (same members, different order)
// converges after apply. to_connect_vpc_ip_block is write-only (not read back,
// listed in ImportStateVerifyIgnore), so state follows config and the post-apply
// plan is empty. Update applies only the [0] element, so reordering changes
// which block is applied (a real modify). Covers TypeListOrderCoverage for this
// TypeList and provides real modify coverage of its nested vswitch_id/zone_id
// attributes, which were previously uncovered.
func TestAccAliCloudApiGatewayInstance_toConnectVpcIpBlockOrder(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_api_gateway_instance.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudApiGatewayVpcConnectInstanceCheckMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return apigateway.NewApiGatewayServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeApiGatewayInstance")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sapigate%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudApiGatewayVpcConnectInstanceBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name": name,
					"instance_spec": "api.s1.small",
					"https_policy":  "HTTPS2_TLS1_0",
					"payment_type":  "PayAsYouGo",
					"instance_type": "vpc_connect",
					"instance_cidr": "192.168.0.0/16",
					"user_vpc_id":   "${alicloud_vpc.vpc.id}",
					"zone_vswitch_security_group": []map[string]interface{}{
						{
							"zone_id":        "${alicloud_vswitch.vswitch_1.zone_id}",
							"vswitch_id":     "${alicloud_vswitch.vswitch_1.id}",
							"cidr_block":     "${alicloud_vswitch.vswitch_1.cidr_block}",
							"security_group": "${alicloud_security_group.security_group_1.id}",
						},
						{
							"zone_id":        "${alicloud_vswitch.vswitch_2.zone_id}",
							"vswitch_id":     "${alicloud_vswitch.vswitch_2.id}",
							"cidr_block":     "${alicloud_vswitch.vswitch_2.cidr_block}",
							"security_group": "${alicloud_security_group.security_group_2.id}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name":       name,
						"instance_spec":       "api.s1.small",
						"https_policy":        "HTTPS2_TLS1_0",
						"payment_type":        "PayAsYouGo",
						"instance_cidr":       "192.168.0.0/16",
						"user_vpc_id":         acctest.CHECKSET,
						"connect_cidr_blocks": `["172.16.0.0/16","172.17.0.0/16"]`,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"to_connect_vpc_ip_block": []map[string]interface{}{
						{
							"vswitch_id": "${alicloud_vswitch.vswitch_1.id}",
							"zone_id":    "${alicloud_vswitch.vswitch_1.zone_id}",
							"cidr_block": "10.1.0.0/24",
							"customized": "true",
						},
						{
							"vswitch_id": "${alicloud_vswitch.vswitch_2.id}",
							"zone_id":    "${alicloud_vswitch.vswitch_2.zone_id}",
							"cidr_block": "10.2.0.0/24",
							"customized": "false",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"connect_cidr_blocks": `["172.16.0.0/16","172.17.0.0/16","10.1.0.0/24"]`,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"to_connect_vpc_ip_block": []map[string]interface{}{
						{
							"vswitch_id": "${alicloud_vswitch.vswitch_2.id}",
							"zone_id":    "${alicloud_vswitch.vswitch_2.zone_id}",
							"cidr_block": "10.2.0.0/24",
							"customized": "false",
						},
						{
							"vswitch_id": "${alicloud_vswitch.vswitch_1.id}",
							"zone_id":    "${alicloud_vswitch.vswitch_1.zone_id}",
							"cidr_block": "10.1.0.0/24",
							"customized": "true",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"connect_cidr_blocks": acctest.CHECKSET,
					}),
				),
			},
		},
	})
}

func TestUnitApiGatewayInstanceStateUpgradeV0(t *testing.T) {
	cases := []struct {
		name     string
		input    map[string]interface{}
		expected []interface{}
	}{
		{
			name: "map with data",
			input: map[string]interface{}{
				"id":            "api-inst-123",
				"instance_name": "test-instance",
				"to_connect_vpc_ip_block": map[string]interface{}{
					"cidr_block": "172.16.0.0/24",
					"customized": "true",
				},
			},
			expected: []interface{}{
				map[string]interface{}{
					"cidr_block": "172.16.0.0/24",
					"customized": "true",
				},
			},
		},
		{
			name: "empty map",
			input: map[string]interface{}{
				"id":                      "api-inst-456",
				"to_connect_vpc_ip_block": map[string]interface{}{},
			},
			expected: []interface{}{},
		},
		{
			name: "nil value",
			input: map[string]interface{}{
				"id":                      "api-inst-789",
				"to_connect_vpc_ip_block": nil,
			},
			expected: nil,
		},
		{
			name: "field not present",
			input: map[string]interface{}{
				"id": "api-inst-000",
			},
			expected: nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := apigateway.ResourceInstanceStateUpgradeV0(context.Background(), tc.input, nil)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			got := result["to_connect_vpc_ip_block"]
			if tc.expected == nil {
				if got != nil {
					t.Errorf("expected nil, got %v", got)
				}
				return
			}
			gotList, ok := got.([]interface{})
			if !ok {
				t.Fatalf("expected []interface{}, got %T", got)
			}
			if len(gotList) != len(tc.expected) {
				t.Fatalf("expected length %d, got %d", len(tc.expected), len(gotList))
			}
			for i, item := range gotList {
				gotMap := item.(map[string]interface{})
				expMap := tc.expected[i].(map[string]interface{})
				for k, v := range expMap {
					if gotMap[k] != v {
						t.Errorf("item[%d].%s: expected %v, got %v", i, k, v, gotMap[k])
					}
				}
			}
		})
	}
}

func TestUnitApiGatewayInstanceSchemaVersionV0ToV1(t *testing.T) {
	r := apigateway.ResourceInstance()
	if r.SchemaVersion != 1 {
		t.Errorf("expected SchemaVersion 1, got %d", r.SchemaVersion)
	}
	if len(r.StateUpgraders) != 1 {
		t.Fatalf("expected 1 StateUpgrader, got %d", len(r.StateUpgraders))
	}
	if r.StateUpgraders[0].Version != 0 {
		t.Errorf("expected StateUpgrader version 0, got %d", r.StateUpgraders[0].Version)
	}
}

func TestAccAliCloudApiGatewayInstance_StateMigrationV0ToV1(t *testing.T) {
	if os.Getenv("ALICLOUD_STATE_MIGRATION_V0_V1") == "" {
		t.Skip("ALICLOUD_STATE_MIGRATION_V0_V1 not set. Test for V0 -> V1 state migration")
	}
	var v map[string]interface{}
	resourceId := "alicloud_api_gateway_instance.default"
	ra := acctest.ResourceAttrInit(resourceId, map[string]string{})
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return apigateway.NewApiGatewayServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeApiGatewayInstance")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc-statemig-%d", rand)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		CheckDestroy: rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"alicloud": {
						Source:            "aliyun/alicloud",
						VersionConstraint: "1.282.0",
					},
				},
				Config: testAccApiGatewayInstanceStateMigrationConfigV0(name),
			},
			{
				ProviderFactories: acctest.SDKProviderFactories,
				Config:            testAccApiGatewayInstanceStateMigrationConfigV1(name),
				//PlanOnly:          true,
			},
		},
	})
}

func testAccApiGatewayInstanceStateMigrationConfigV0(name string) string {
	return testAccApiGatewayInstanceStateMigrationDeps(name) + fmt.Sprintf(`
resource "alicloud_api_gateway_instance" "default" {
  instance_name           = "%s"
  instance_spec           = "api.s1.small"
  https_policy            = "HTTPS2_TLS1_0"
  zone_id                 = "cn-hangzhou-MAZ6(i,j,k)"
  payment_type            = "PayAsYouGo"
  instance_type           = "vpc_connect"
  user_vpc_id             = alicloud_vpc.vpc.id
  vpc_slb_intranet_enable = "true"
  ipv6_enabled            = "true"
  to_connect_vpc_ip_block = {
    cidr_block = "10.1.0.0/24"
    customized = "true"
  }
  zone_vswitch_security_group {
    zone_id        = alicloud_vswitch.vswitch_1.zone_id
    vswitch_id     = alicloud_vswitch.vswitch_1.id
    cidr_block     = alicloud_vswitch.vswitch_1.cidr_block
    security_group = alicloud_security_group.security_group_1.id
  }
  zone_vswitch_security_group {
    zone_id        = alicloud_vswitch.vswitch_2.zone_id
    vswitch_id     = alicloud_vswitch.vswitch_2.id
    cidr_block     = alicloud_vswitch.vswitch_2.cidr_block
    security_group = alicloud_security_group.security_group_2.id
  }
}

output "to_connect_vpc_ip_block" {
  value = alicloud_api_gateway_instance.default.to_connect_vpc_ip_block["cidr_block"]
}
`, name)
}

func testAccApiGatewayInstanceStateMigrationConfigV1(name string) string {
	return testAccApiGatewayInstanceStateMigrationDeps(name) + fmt.Sprintf(`
resource "alicloud_api_gateway_instance" "default" {
  instance_name           = "%s"
  instance_spec           = "api.s1.small"
  https_policy            = "HTTPS2_TLS1_0"
  zone_id                 = "cn-hangzhou-MAZ6(i,j,k)"
  payment_type            = "PayAsYouGo"
  instance_type           = "vpc_connect"
  user_vpc_id             = alicloud_vpc.vpc.id
  vpc_slb_intranet_enable = "true"
  ipv6_enabled            = "true"
  to_connect_vpc_ip_block {
    cidr_block = "10.1.0.0/24"
    customized = true
  }
  zone_vswitch_security_group {
    zone_id        = alicloud_vswitch.vswitch_1.zone_id
    vswitch_id     = alicloud_vswitch.vswitch_1.id
    cidr_block     = alicloud_vswitch.vswitch_1.cidr_block
    security_group = alicloud_security_group.security_group_1.id
  }
  zone_vswitch_security_group {
    zone_id        = alicloud_vswitch.vswitch_2.zone_id
    vswitch_id     = alicloud_vswitch.vswitch_2.id
    cidr_block     = alicloud_vswitch.vswitch_2.cidr_block
    security_group = alicloud_security_group.security_group_2.id
  }
}

output "to_connect_vpc_ip_block" {
  value = alicloud_api_gateway_instance.default.to_connect_vpc_ip_block.0.cidr_block
}
`, name)
}

func testAccApiGatewayInstanceStateMigrationDeps(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}
resource "alicloud_vpc" "vpc" {
  vpc_name   = var.name
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "vswitch_1" {
  vpc_id       = alicloud_vpc.vpc.id
  cidr_block   = "172.16.0.0/16"
  zone_id      = "cn-hangzhou-j"
  vswitch_name = "${var.name}_1"
}

resource "alicloud_vswitch" "vswitch_2" {
  vpc_id       = alicloud_vpc.vpc.id
  cidr_block   = "172.17.0.0/16"
  zone_id      = "cn-hangzhou-k"
  vswitch_name = "${var.name}_2"
}

resource "alicloud_security_group" "security_group_1" {
  vpc_id              = alicloud_vpc.vpc.id
  security_group_name = "${var.name}_1"
}

resource "alicloud_security_group" "security_group_2" {
  vpc_id              = alicloud_vpc.vpc.id
  security_group_name = "${var.name}_2"
}
`, name)
}
