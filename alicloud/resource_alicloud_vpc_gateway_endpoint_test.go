package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Vpc GatewayEndpoint. >>> Resource test cases, automatically generated.
// Case 20250210GatewayEndpoint 10173
func TestAccAliCloudVpcGatewayEndpoint_basic10173(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vpc_gateway_endpoint.default"
	ra := resourceAttrInit(resourceId, AlicloudVpcGatewayEndpointMap10173)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVpcGatewayEndpoint")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccvpc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudVpcGatewayEndpointBasicDependence10173)
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
					"gateway_endpoint_descrption": "test-gateway-endpoint",
					"gateway_endpoint_name":       name,
					"vpc_id":                      "${alicloud_vpc.defaultVpc.id}",
					"resource_group_id":           "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"service_name":                "${var.domain}",
					"policy_document":             "{ \\\"Version\\\" : \\\"1\\\", \\\"Statement\\\" : [ { \\\"Effect\\\" : \\\"Allow\\\", \\\"Resource\\\" : [ \\\"*\\\" ], \\\"Action\\\" : [ \\\"*\\\" ], \\\"Principal\\\" : [ \\\"*\\\" ] } ] }",
					"route_tables": []string{
						"${alicloud_route_table.routeTable2.id}", "${alicloud_route_table.routeTable3.id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"gateway_endpoint_descrption": "test-gateway-endpoint",
						"gateway_endpoint_name":       name,
						"vpc_id":                      CHECKSET,
						"resource_group_id":           CHECKSET,
						"service_name":                CHECKSET,
						"policy_document":             CHECKSET,
						"route_tables.#":              "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"gateway_endpoint_descrption": "testupdate",
					"gateway_endpoint_name":       name + "_update",
					"resource_group_id":           "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"policy_document":             "{ \\\"Version\\\" : \\\"1\\\", \\\"Statement\\\" : [ { \\\"Effect\\\" : \\\"Deny\\\", \\\"Resource\\\" : [ \\\"*\\\" ], \\\"Action\\\" : [ \\\"*\\\" ], \\\"Principal\\\" : [ \\\"*\\\" ] } ] }",
					"route_tables": []string{
						"${alicloud_route_table.routeTable2.id}", "${alicloud_route_table.routeTable1.id}", "${alicloud_route_table.routeTable3.id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"gateway_endpoint_descrption": "testupdate",
						"gateway_endpoint_name":       name + "_update",
						"resource_group_id":           CHECKSET,
						"policy_document":             CHECKSET,
						"route_tables.#":              "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"route_tables": []string{},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"route_tables.#": "0",
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

var AlicloudVpcGatewayEndpointMap10173 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
}

func AlicloudVpcGatewayEndpointBasicDependence10173(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "domain" {
  default = "com.aliyun.cn-hangzhou.oss"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_vpc" "defaultVpc" {
  description = "tf-example"
}

resource "alicloud_route_table" "routeTable1" {
  vpc_id           = alicloud_vpc.defaultVpc.id
  route_table_name = "testGatewayEndpointAssociationRoutetable1"
}

resource "alicloud_route_table" "routeTable2" {
  vpc_id           = alicloud_vpc.defaultVpc.id
  route_table_name = "testGatewayEndpointAssociationRoutetable2"
}

resource "alicloud_route_table" "routeTable3" {
  vpc_id           = alicloud_vpc.defaultVpc.id
  route_table_name = "testGatewayEndpointAssociationRoutetable3"
}


`, name)
}

// Case 订正资源 3630
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
	name := fmt.Sprintf("tfaccvpc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudVpcGatewayEndpointBasicDependence3630)
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
					"gateway_endpoint_descrption": "test-gateway-endpoint",
					"gateway_endpoint_name":       name,
					"vpc_id":                      "${alicloud_vpc.defaultVpc.id}",
					"resource_group_id":           "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"service_name":                "${var.domain}",
					"policy_document":             "{ \\\"Version\\\" : \\\"1\\\", \\\"Statement\\\" : [ { \\\"Effect\\\" : \\\"Allow\\\", \\\"Resource\\\" : [ \\\"*\\\" ], \\\"Action\\\" : [ \\\"*\\\" ], \\\"Principal\\\" : [ \\\"*\\\" ] } ] }",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"gateway_endpoint_descrption": "test-gateway-endpoint",
						"gateway_endpoint_name":       name,
						"vpc_id":                      CHECKSET,
						"resource_group_id":           CHECKSET,
						"service_name":                CHECKSET,
						"policy_document":             CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"gateway_endpoint_descrption": "testupdate",
					"gateway_endpoint_name":       name + "_update",
					"resource_group_id":           "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"policy_document":             "{ \\\"Version\\\" : \\\"1\\\", \\\"Statement\\\" : [ { \\\"Effect\\\" : \\\"Deny\\\", \\\"Resource\\\" : [ \\\"*\\\" ], \\\"Action\\\" : [ \\\"*\\\" ], \\\"Principal\\\" : [ \\\"*\\\" ] } ] }",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"gateway_endpoint_descrption": "testupdate",
						"gateway_endpoint_name":       name + "_update",
						"resource_group_id":           CHECKSET,
						"policy_document":             CHECKSET,
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

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_vpc" "defaultVpc" {
  description = "tf-example"
}


`, name)
}

// Case 对接Terraform_换账号_覆盖率_文档修订3 3621
func TestAccAliCloudVpcGatewayEndpoint_basic3621(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vpc_gateway_endpoint.default"
	ra := resourceAttrInit(resourceId, AlicloudVpcGatewayEndpointMap3621)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVpcGatewayEndpoint")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccvpc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudVpcGatewayEndpointBasicDependence3621)
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
					"gateway_endpoint_descrption": "test-gateway-endpoint",
					"gateway_endpoint_name":       name,
					"service_name":                "${var.demoin}",
					"vpc_id":                      "${alicloud_vpc.defaultvpc.id}",
					"policy_document":             "{ \\\"Version\\\" : \\\"1\\\", \\\"Statement\\\" : [ { \\\"Effect\\\" : \\\"Allow\\\", \\\"Resource\\\" : [ \\\"*\\\" ], \\\"Action\\\" : [ \\\"*\\\" ], \\\"Principal\\\" : [ \\\"*\\\" ] } ] }",
					"resource_group_id":           "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"route_tables": []string{
						"${alicloud_route_table.defaultRt.id}", "${alicloud_route_table.defaultrt1.id}", "${alicloud_route_table.defaultrt2.id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"gateway_endpoint_descrption": "test-gateway-endpoint",
						"gateway_endpoint_name":       name,
						"service_name":                CHECKSET,
						"vpc_id":                      CHECKSET,
						"policy_document":             CHECKSET,
						"resource_group_id":           CHECKSET,
						"route_tables.#":              "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"gateway_endpoint_descrption": "testupdate",
					"gateway_endpoint_name":       name + "_update",
					"policy_document":             "{ \\\"Version\\\" : \\\"1\\\", \\\"Statement\\\" : [ { \\\"Effect\\\" : \\\"Deny\\\", \\\"Resource\\\" : [ \\\"*\\\" ], \\\"Action\\\" : [ \\\"*\\\" ], \\\"Principal\\\" : [ \\\"*\\\" ] } ] }",
					"resource_group_id":           "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"route_tables": []string{
						"${alicloud_route_table.defaultrt2.id}", "${alicloud_route_table.defaultrt1.id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"gateway_endpoint_descrption": "testupdate",
						"gateway_endpoint_name":       name + "_update",
						"policy_document":             CHECKSET,
						"resource_group_id":           CHECKSET,
						"route_tables.#":              "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"gateway_endpoint_descrption": "test-mod-description",
					"gateway_endpoint_name":       name + "_update",
					"resource_group_id":           "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"route_tables":                []string{},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"gateway_endpoint_descrption": "test-mod-description",
						"gateway_endpoint_name":       name + "_update",
						"resource_group_id":           CHECKSET,
						"route_tables.#":              "0",
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
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudVpcGatewayEndpointMap3621 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
}

func AlicloudVpcGatewayEndpointBasicDependence3621(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "demoin" {
  default = "com.aliyun.cn-hangzhou.oss"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_vpc" "defaultvpc" {
  dry_run     = false
  enable_ipv6 = false
}

resource "alicloud_route_table" "defaultRt" {
  vpc_id      = alicloud_vpc.defaultvpc.id
  description = "tf-testacc"
}

resource "alicloud_route_table" "defaultrt1" {
  vpc_id      = alicloud_vpc.defaultvpc.id
  description = "tf-testacc1"
}

resource "alicloud_route_table" "defaultrt2" {
  vpc_id      = alicloud_vpc.defaultvpc.id
  description = "tf-testacc2"
}


`, name)
}

// Case 对接Terraform_换账号_覆盖率_文档修订2 3620
func TestAccAliCloudVpcGatewayEndpoint_basic3620(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vpc_gateway_endpoint.default"
	ra := resourceAttrInit(resourceId, AlicloudVpcGatewayEndpointMap3620)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVpcGatewayEndpoint")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccvpc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudVpcGatewayEndpointBasicDependence3620)
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
					"gateway_endpoint_descrption": "test-gateway-endpoint",
					"gateway_endpoint_name":       name,
					"service_name":                "${var.demoin}",
					"vpc_id":                      "${alicloud_vpc.defaultvpc.id}",
					"policy_document":             "{ \\\"Version\\\" : \\\"1\\\", \\\"Statement\\\" : [ { \\\"Effect\\\" : \\\"Allow\\\", \\\"Resource\\\" : [ \\\"*\\\" ], \\\"Action\\\" : [ \\\"*\\\" ], \\\"Principal\\\" : [ \\\"*\\\" ] } ] }",
					"resource_group_id":           "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"route_tables": []string{
						"${alicloud_route_table.defaultRt.id}", "${alicloud_route_table.defaultrt1.id}", "${alicloud_route_table.defaultrt2.id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"gateway_endpoint_descrption": "test-gateway-endpoint",
						"gateway_endpoint_name":       name,
						"service_name":                CHECKSET,
						"vpc_id":                      CHECKSET,
						"policy_document":             CHECKSET,
						"resource_group_id":           CHECKSET,
						"route_tables.#":              "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"gateway_endpoint_descrption": "testupdate",
					"gateway_endpoint_name":       name + "_update",
					"policy_document":             "{ \\\"Version\\\" : \\\"1\\\", \\\"Statement\\\" : [ { \\\"Effect\\\" : \\\"Deny\\\", \\\"Resource\\\" : [ \\\"*\\\" ], \\\"Action\\\" : [ \\\"*\\\" ], \\\"Principal\\\" : [ \\\"*\\\" ] } ] }",
					"resource_group_id":           "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"route_tables": []string{
						"${alicloud_route_table.defaultrt2.id}", "${alicloud_route_table.defaultrt1.id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"gateway_endpoint_descrption": "testupdate",
						"gateway_endpoint_name":       name + "_update",
						"policy_document":             CHECKSET,
						"resource_group_id":           CHECKSET,
						"route_tables.#":              "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"gateway_endpoint_descrption": "test-mod-description",
					"gateway_endpoint_name":       name + "_update",
					"resource_group_id":           "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"route_tables":                []string{},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"gateway_endpoint_descrption": "test-mod-description",
						"gateway_endpoint_name":       name + "_update",
						"resource_group_id":           CHECKSET,
						"route_tables.#":              "0",
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
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudVpcGatewayEndpointMap3620 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
}

func AlicloudVpcGatewayEndpointBasicDependence3620(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "demoin" {
  default = "com.aliyun.cn-hangzhou.oss"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_vpc" "defaultvpc" {
  dry_run     = false
  enable_ipv6 = false
}

resource "alicloud_route_table" "defaultRt" {
  vpc_id      = alicloud_vpc.defaultvpc.id
  description = "tf-testacc"
}

resource "alicloud_route_table" "defaultrt1" {
  vpc_id      = alicloud_vpc.defaultvpc.id
  description = "tf-testacc1"
}

resource "alicloud_route_table" "defaultrt2" {
  vpc_id      = alicloud_vpc.defaultvpc.id
  description = "tf-testacc2"
}


`, name)
}

// Case 对接Terraform_换账号_覆盖率_文档修订 3619
func TestAccAliCloudVpcGatewayEndpoint_basic3619(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vpc_gateway_endpoint.default"
	ra := resourceAttrInit(resourceId, AlicloudVpcGatewayEndpointMap3619)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVpcGatewayEndpoint")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccvpc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudVpcGatewayEndpointBasicDependence3619)
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
					"gateway_endpoint_descrption": "test-gateway-endpoint",
					"gateway_endpoint_name":       name,
					"service_name":                "${var.demoin}",
					"vpc_id":                      "${alicloud_vpc.defaultvpc.id}",
					"policy_document":             "{ \\\"Version\\\" : \\\"1\\\", \\\"Statement\\\" : [ { \\\"Effect\\\" : \\\"Allow\\\", \\\"Resource\\\" : [ \\\"*\\\" ], \\\"Action\\\" : [ \\\"*\\\" ], \\\"Principal\\\" : [ \\\"*\\\" ] } ] }",
					"resource_group_id":           "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"route_tables": []string{
						"${alicloud_route_table.defaultRt.id}", "${alicloud_route_table.defaultrt1.id}", "${alicloud_route_table.defaultrt2.id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"gateway_endpoint_descrption": "test-gateway-endpoint",
						"gateway_endpoint_name":       name,
						"service_name":                CHECKSET,
						"vpc_id":                      CHECKSET,
						"policy_document":             CHECKSET,
						"resource_group_id":           CHECKSET,
						"route_tables.#":              "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"gateway_endpoint_descrption": "testupdate",
					"gateway_endpoint_name":       name + "_update",
					"policy_document":             "{ \\\"Version\\\" : \\\"1\\\", \\\"Statement\\\" : [ { \\\"Effect\\\" : \\\"Deny\\\", \\\"Resource\\\" : [ \\\"*\\\" ], \\\"Action\\\" : [ \\\"*\\\" ], \\\"Principal\\\" : [ \\\"*\\\" ] } ] }",
					"resource_group_id":           "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"route_tables": []string{
						"${alicloud_route_table.defaultrt2.id}", "${alicloud_route_table.defaultrt1.id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"gateway_endpoint_descrption": "testupdate",
						"gateway_endpoint_name":       name + "_update",
						"policy_document":             CHECKSET,
						"resource_group_id":           CHECKSET,
						"route_tables.#":              "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"gateway_endpoint_descrption": "test-mod-description",
					"gateway_endpoint_name":       name + "_update",
					"resource_group_id":           "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"route_tables":                []string{},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"gateway_endpoint_descrption": "test-mod-description",
						"gateway_endpoint_name":       name + "_update",
						"resource_group_id":           CHECKSET,
						"route_tables.#":              "0",
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
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudVpcGatewayEndpointMap3619 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
}

func AlicloudVpcGatewayEndpointBasicDependence3619(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "demoin" {
  default = "com.aliyun.cn-hangzhou.oss"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_vpc" "defaultvpc" {
  dry_run     = false
  enable_ipv6 = false
}

resource "alicloud_route_table" "defaultRt" {
  vpc_id      = alicloud_vpc.defaultvpc.id
  description = "tf-testacc"
}

resource "alicloud_route_table" "defaultrt1" {
  vpc_id      = alicloud_vpc.defaultvpc.id
  description = "tf-testacc1"
}

resource "alicloud_route_table" "defaultrt2" {
  vpc_id      = alicloud_vpc.defaultvpc.id
  description = "tf-testacc2"
}


`, name)
}

// Case 对接Terraform_换账号_覆盖率 3617
func TestAccAliCloudVpcGatewayEndpoint_basic3617(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vpc_gateway_endpoint.default"
	ra := resourceAttrInit(resourceId, AlicloudVpcGatewayEndpointMap3617)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVpcGatewayEndpoint")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccvpc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudVpcGatewayEndpointBasicDependence3617)
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
					"gateway_endpoint_descrption": "test-gateway-endpoint",
					"gateway_endpoint_name":       name,
					"service_name":                "${var.demoin}",
					"vpc_id":                      "${alicloud_vpc.defaultvpc.id}",
					"policy_document":             "{ \\\"Version\\\" : \\\"1\\\", \\\"Statement\\\" : [ { \\\"Effect\\\" : \\\"Allow\\\", \\\"Resource\\\" : [ \\\"*\\\" ], \\\"Action\\\" : [ \\\"*\\\" ], \\\"Principal\\\" : [ \\\"*\\\" ] } ] }",
					"resource_group_id":           "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"route_tables": []string{
						"${alicloud_route_table.defaultRt.id}", "${alicloud_route_table.defaultrt1.id}", "${alicloud_route_table.defaultrt2.id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"gateway_endpoint_descrption": "test-gateway-endpoint",
						"gateway_endpoint_name":       name,
						"service_name":                CHECKSET,
						"vpc_id":                      CHECKSET,
						"policy_document":             CHECKSET,
						"resource_group_id":           CHECKSET,
						"route_tables.#":              "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"gateway_endpoint_descrption": "testupdate",
					"gateway_endpoint_name":       name + "_update",
					"policy_document":             "{ \\\"Version\\\" : \\\"1\\\", \\\"Statement\\\" : [ { \\\"Effect\\\" : \\\"Deny\\\", \\\"Resource\\\" : [ \\\"*\\\" ], \\\"Action\\\" : [ \\\"*\\\" ], \\\"Principal\\\" : [ \\\"*\\\" ] } ] }",
					"resource_group_id":           "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"route_tables": []string{
						"${alicloud_route_table.defaultrt2.id}", "${alicloud_route_table.defaultrt1.id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"gateway_endpoint_descrption": "testupdate",
						"gateway_endpoint_name":       name + "_update",
						"policy_document":             CHECKSET,
						"resource_group_id":           CHECKSET,
						"route_tables.#":              "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"gateway_endpoint_descrption": "test-mod-description",
					"gateway_endpoint_name":       name + "_update",
					"resource_group_id":           "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"route_tables":                []string{},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"gateway_endpoint_descrption": "test-mod-description",
						"gateway_endpoint_name":       name + "_update",
						"resource_group_id":           CHECKSET,
						"route_tables.#":              "0",
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
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudVpcGatewayEndpointMap3617 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
}

func AlicloudVpcGatewayEndpointBasicDependence3617(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "demoin" {
  default = "com.aliyun.cn-hangzhou.oss"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_vpc" "defaultvpc" {
  dry_run     = false
  enable_ipv6 = false
}

resource "alicloud_route_table" "defaultRt" {
  vpc_id      = alicloud_vpc.defaultvpc.id
  description = "tf-testacc"
}

resource "alicloud_route_table" "defaultrt1" {
  vpc_id      = alicloud_vpc.defaultvpc.id
  description = "tf-testacc1"
}

resource "alicloud_route_table" "defaultrt2" {
  vpc_id      = alicloud_vpc.defaultvpc.id
  description = "tf-testacc2"
}


`, name)
}

// Case 对接Terraform_换账号 3615
func TestAccAliCloudVpcGatewayEndpoint_basic3615(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vpc_gateway_endpoint.default"
	ra := resourceAttrInit(resourceId, AlicloudVpcGatewayEndpointMap3615)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVpcGatewayEndpoint")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccvpc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudVpcGatewayEndpointBasicDependence3615)
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
					"gateway_endpoint_descrption": "test-gateway-endpoint",
					"gateway_endpoint_name":       name,
					"service_name":                "${var.demoin}",
					"vpc_id":                      "${alicloud_vpc.defaultvpc.id}",
					"policy_document":             "{ \\\"Version\\\" : \\\"1\\\", \\\"Statement\\\" : [ { \\\"Effect\\\" : \\\"Allow\\\", \\\"Resource\\\" : [ \\\"*\\\" ], \\\"Action\\\" : [ \\\"*\\\" ], \\\"Principal\\\" : [ \\\"*\\\" ] } ] }",
					"resource_group_id":           "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"gateway_endpoint_descrption": "test-gateway-endpoint",
						"gateway_endpoint_name":       name,
						"service_name":                CHECKSET,
						"vpc_id":                      CHECKSET,
						"policy_document":             CHECKSET,
						"resource_group_id":           CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"gateway_endpoint_descrption": "testupdate",
					"gateway_endpoint_name":       name + "_update",
					"resource_group_id":           "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"gateway_endpoint_descrption": "testupdate",
						"gateway_endpoint_name":       name + "_update",
						"resource_group_id":           CHECKSET,
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

var AlicloudVpcGatewayEndpointMap3615 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
}

func AlicloudVpcGatewayEndpointBasicDependence3615(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "demoin" {
  default = "com.aliyun.cn-hangzhou.oss"
}

resource "alicloud_vpc" "defaultvpc" {
  dry_run     = false
  enable_ipv6 = false
}

data "alicloud_resource_manager_resource_groups" "default"{
	status = "OK"
}
`, name)
}

// Test Vpc GatewayEndpoint. <<< Resource test cases, automatically generated.
