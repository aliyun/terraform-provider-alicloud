package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudApiGatewayBackendModel_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_api_gateway_backend_model.default"
	ra := resourceAttrInit(resourceId, AlicloudApiGatewayBackendModelMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudApiService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeApiGatewayBackendModel")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testacc%sapigatewaybackendmodel%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudApiGatewayBackendModelBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"backend_id":         "${alicloud_api_gateway_backend.default.id}",
					"backend_type":       "HTTP",
					"stage_name":         "RELEASE",
					"description":        "tf-testAcc-desc",
					"backend_model_data": `{\"ServiceAddress\":\"http://apigateway.alicloudapi.com:8080\"}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"stage_name":         "RELEASE",
						"description":        "tf-testAcc-desc",
						"backend_model_data": CHECKSET,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "tf-testAcc-desc-update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "tf-testAcc-desc-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"backend_model_data": `{\"ServiceAddress\":\"http://apigateway.alicloudapi.com:9090\"}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backend_model_data": CHECKSET,
					}),
				),
			},
		},
	})
}

var AlicloudApiGatewayBackendModelMap0 = map[string]string{
	"backend_model_id": CHECKSET,
}

func AlicloudApiGatewayBackendModelBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}
resource "alicloud_api_gateway_backend" "default" {
  backend_name = var.name
  backend_type = "HTTP"
  description  = var.name
}
`, name)
}

func TestAccAliCloudApiGatewayBackendModel_vpcFullChain(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_api_gateway_backend_model.default"
	ra := resourceAttrInit(resourceId, AlicloudApiGatewayBackendModelMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudApiService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeApiGatewayBackendModel")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testacc%sapigatewaybackendmodel%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudApiGatewayBackendModelVpcFullChainDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"backend_id":         "${alicloud_api_gateway_backend.default.id}",
					"backend_type":       "VPC",
					"stage_name":         "RELEASE",
					"description":        "tf-testAcc-vpc-desc",
					"backend_model_data": `{\"VpcConfig\":{\"VpcAccessId\":\"${alicloud_api_gateway_vpc_access.default.vpc_access_id}\",\"VpcTargetHostName\":\"www.host.com\"}}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"stage_name":         "RELEASE",
						"description":        "tf-testAcc-vpc-desc",
						"backend_model_data": CHECKSET,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"backend_model_data": `{\"VpcConfig\":{\"VpcAccessId\":\"${alicloud_api_gateway_vpc_access.default.vpc_access_id}\",\"VpcScheme\":\"https\",\"VpcTargetHostName\":\"www.update.com\"}}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backend_model_data": CHECKSET,
					}),
				),
			},
		},
	})
}

func AlicloudApiGatewayBackendModelVpcFullChainDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vswitch_name = var.name
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "192.168.192.0/24"
  zone_id      = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_slb_load_balancer" "default" {
  load_balancer_name = var.name
  address_type       = "intranet"
  vswitch_id         = alicloud_vswitch.default.id
  load_balancer_spec = "slb.s1.small"
}

resource "alicloud_api_gateway_group" "default" {
  name        = var.name
  description = "tf_testAcc_api group description"
}

resource "alicloud_api_gateway_vpc_access" "default" {
  name        = var.name
  vpc_id      = alicloud_vpc.default.id
  instance_id = alicloud_slb_load_balancer.default.id
  port        = "80"
}

resource "alicloud_api_gateway_backend" "default" {
  backend_name = "${var.name}_backend"
  backend_type = "VPC"
  description  = var.name
}

resource "alicloud_api_gateway_api" "default" {
  name        = var.name
  group_id    = alicloud_api_gateway_group.default.id
  description = "tf_testAcc_api with vpc backend"
  auth_type   = "ANONYMOUS"

  request_config {
    protocol = "HTTP"
    method   = "GET"
    path     = "/test/path"
    mode     = "MAPPING"
  }

  service_type = "HTTP-VPC"

  http_vpc_service_config {
    name    = alicloud_api_gateway_vpc_access.default.name
    path    = "/api/v1/test"
    method  = "GET"
    timeout = 30
  }

  backend_id = alicloud_api_gateway_backend.default.id
}
`, name)
}
