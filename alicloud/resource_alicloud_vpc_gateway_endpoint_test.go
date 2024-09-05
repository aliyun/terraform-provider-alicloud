package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Vpc GatewayEndpoint. >>> Resource test cases, automatically generated.
// Case 3630
func TestAccAliCloudVpcGatewayEndpoint_basic3630(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vpc_gateway_endpoint.default"
	ra := resourceAttrInit(resourceId, AlicloudVpcGatewayEndpointMap3630)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVpcGatewayEndpoint")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%svpcgatewayendpoint%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudVpcGatewayEndpointBasicDependence3630)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.VPCGatewayEndpointSupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"vpc_id":                "${alicloud_vpc.defaultVpc.id}",
					"service_name":          "${var.domain}",
					"gateway_endpoint_name": name,
					"policy_document":       "{ \\\"Version\\\" : \\\"1\\\", \\\"Statement\\\" : [ { \\\"Effect\\\" : \\\"Allow\\\", \\\"Resource\\\" : [ \\\"*\\\" ], \\\"Action\\\" : [ \\\"*\\\" ], \\\"Principal\\\" : [ \\\"*\\\" ] } ] }",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpc_id":                CHECKSET,
						"service_name":          CHECKSET,
						"gateway_endpoint_name": name,
						"policy_document":       "{ \"Version\" : \"1\", \"Statement\" : [ { \"Effect\" : \"Allow\", \"Resource\" : [ \"*\" ], \"Action\" : [ \"*\" ], \"Principal\" : [ \"*\" ] } ] }",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"gateway_endpoint_descrption": "test-gateway-endpoint",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"gateway_endpoint_descrption": "test-gateway-endpoint",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_document": "{ \\\"Version\\\" : \\\"1\\\", \\\"Statement\\\" : [ { \\\"Effect\\\" : \\\"Allow\\\", \\\"Resource\\\" : [ \\\"*\\\" ], \\\"Action\\\" : [ \\\"*\\\" ], \\\"Principal\\\" : [ \\\"*\\\" ] } ] }",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy_document": "{ \"Version\" : \"1\", \"Statement\" : [ { \"Effect\" : \"Allow\", \"Resource\" : [ \"*\" ], \"Action\" : [ \"*\" ], \"Principal\" : [ \"*\" ] } ] }",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"gateway_endpoint_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"gateway_endpoint_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"gateway_endpoint_descrption": "testupdate",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"gateway_endpoint_descrption": "testupdate",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"gateway_endpoint_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"gateway_endpoint_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_document": "{ \\\"Version\\\" : \\\"1\\\", \\\"Statement\\\" : [ { \\\"Effect\\\" : \\\"Deny\\\", \\\"Resource\\\" : [ \\\"*\\\" ], \\\"Action\\\" : [ \\\"*\\\" ], \\\"Principal\\\" : [ \\\"*\\\" ] } ] }",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy_document": "{ \"Version\" : \"1\", \"Statement\" : [ { \"Effect\" : \"Deny\", \"Resource\" : [ \"*\" ], \"Action\" : [ \"*\" ], \"Principal\" : [ \"*\" ] } ] }",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"gateway_endpoint_descrption": "test-gateway-endpoint",
					"gateway_endpoint_name":       name + "_update",
					"vpc_id":                      "${alicloud_vpc.defaultVpc.id}",
					"resource_group_id":           "${data.alicloud_resource_manager_resource_groups.default.groups.1.id}",
					"service_name":                "${var.domain}",
					"policy_document":             "{ \\\"Version\\\" : \\\"1\\\", \\\"Statement\\\" : [ { \\\"Effect\\\" : \\\"Allow\\\", \\\"Resource\\\" : [ \\\"*\\\" ], \\\"Action\\\" : [ \\\"*\\\" ], \\\"Principal\\\" : [ \\\"*\\\" ] } ] }",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"gateway_endpoint_descrption": "test-gateway-endpoint",
						"gateway_endpoint_name":       name + "_update",
						"vpc_id":                      CHECKSET,
						"resource_group_id":           CHECKSET,
						"service_name":                CHECKSET,
						"policy_document":             "{ \"Version\" : \"1\", \"Statement\" : [ { \"Effect\" : \"Allow\", \"Resource\" : [ \"*\" ], \"Action\" : [ \"*\" ], \"Principal\" : [ \"*\" ] } ] }",
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
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudVpcGatewayEndpointMap3630 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
}

func AlicloudVpcGatewayEndpointBasicDependence3630(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "domain" {
  default = "com.aliyun.cn-hangzhou.oss"
}

resource "alicloud_vpc" "defaultVpc" {
  description = "tf-example"
}

data "alicloud_resource_manager_resource_groups" "default" {
  status = "OK"
}

`, name)
}

// Case 3630  twin
func TestAccAliCloudVpcGatewayEndpoint_basic3630_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vpc_gateway_endpoint.default"
	ra := resourceAttrInit(resourceId, AlicloudVpcGatewayEndpointMap3630)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVpcGatewayEndpoint")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%svpcgatewayendpoint%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudVpcGatewayEndpointBasicDependence3630)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.VPCGatewayEndpointSupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"gateway_endpoint_descrption": "test-mod-description",
					"gateway_endpoint_name":       name,
					"vpc_id":                      "${alicloud_vpc.defaultVpc.id}",
					"resource_group_id":           "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
					"service_name":                "${var.domain}",
					"policy_document":             "{ \\\"Version\\\" : \\\"1\\\", \\\"Statement\\\" : [ { \\\"Effect\\\" : \\\"Deny\\\", \\\"Resource\\\" : [ \\\"*\\\" ], \\\"Action\\\" : [ \\\"*\\\" ], \\\"Principal\\\" : [ \\\"*\\\" ] } ] }",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"gateway_endpoint_descrption": "test-mod-description",
						"gateway_endpoint_name":       name,
						"vpc_id":                      CHECKSET,
						"resource_group_id":           CHECKSET,
						"service_name":                CHECKSET,
						"policy_document":             "{ \"Version\" : \"1\", \"Statement\" : [ { \"Effect\" : \"Deny\", \"Resource\" : [ \"*\" ], \"Action\" : [ \"*\" ], \"Principal\" : [ \"*\" ] } ] }",
						"tags.%":                      "2",
						"tags.Created":                "TF",
						"tags.For":                    "Test",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

// Test Vpc GatewayEndpoint. <<< Resource test cases, automatically generated.
