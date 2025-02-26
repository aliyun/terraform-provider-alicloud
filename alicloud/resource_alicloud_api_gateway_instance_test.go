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
	"status":      CHECKSET,
	"create_time": CHECKSET,
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
    vpc_id = alicloud_vpc.vpc.id
    name   = "${var.name}_1"
}

resource "alicloud_security_group" "security_group_2" {
    vpc_id = alicloud_vpc.vpc.id
    name   = "${var.name}_2"
}
`, name)
}

func TestAccAliCloudApiGatewayVpcConnectInstance(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_api_gateway_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudApiGatewayVpcConnectInstanceCheckMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApiGatewayServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeApiGatewayInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sapigate%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudApiGatewayVpcConnectInstanceBasicDependence)
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
						"user_vpc_id":         CHECKSET,
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
					"to_connect_vpc_ip_block": map[string]interface{}{
						"cidr_block": "10.1.0.0/24",
						"customized": "true",
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
					"to_connect_vpc_ip_block": REMOVEKEY,
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
	ra := resourceAttrInit(resourceId, AlicloudApiGatewayVpcConnectInstanceCheckMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApiGatewayServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeApiGatewayInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sapigate%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudApiGatewayVpcConnectInstanceBasicDependence)
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
						"user_vpc_id":             CHECKSET,
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
	ra := resourceAttrInit(resourceId, AlicloudApiGatewayVpcConnectInstanceCheckMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApiGatewayServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeApiGatewayInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sapigate%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudApiGatewayVpcConnectInstanceBasicDependence)
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
					"to_connect_vpc_ip_block": map[string]interface{}{
						"cidr_block": "10.1.0.0/24",
						"customized": "true",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name":           name,
						"instance_spec":           "api.s1.small",
						"https_policy":            "HTTPS2_TLS1_0",
						"payment_type":            "PayAsYouGo",
						"instance_cidr":           "192.168.0.0/16",
						"user_vpc_id":             CHECKSET,
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
