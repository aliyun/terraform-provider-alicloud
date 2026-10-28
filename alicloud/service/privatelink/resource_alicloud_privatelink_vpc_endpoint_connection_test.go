package privatelink_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	privatelink "github.com/aliyun/terraform-provider-alicloud/alicloud/service/privatelink"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudPrivatelinkVpcEndpointConnection_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_privatelink_vpc_endpoint_connection.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudPrivatelinkVpcEndpointConnectionMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return privatelink.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribePrivatelinkVpcEndpointConnection")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testaccprivatelinkvpcendpointconnection%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudPrivatelinkVpcEndpointConnectionBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},

		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"endpoint_id":          "${alicloud_privatelink_vpc_endpoint.default.id}",
					"service_id":           "${alicloud_privatelink_vpc_endpoint_service.default.id}",
					"bandwidth":            "1024",
					"traffic_control_mode": "BandwidthLimit",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"endpoint_id":          acctest.CHECKSET,
						"service_id":           acctest.CHECKSET,
						"bandwidth":            "1024",
						"traffic_control_mode": "BandwidthLimit",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bandwidth": "1000",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth": "1000",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dry_run"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"endpoint_id":          "${alicloud_privatelink_vpc_endpoint.default.id}",
					"service_id":           "${alicloud_privatelink_vpc_endpoint_service.default.id}",
					"traffic_control_mode": "Scalability",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"traffic_control_mode": "Scalability",
					}),
				),
			},
		},
	})
}

var AlicloudPrivatelinkVpcEndpointConnectionMap = map[string]string{
	"endpoint_owner_id": acctest.CHECKSET,
	"status":            acctest.CHECKSET,
}

func AlicloudPrivatelinkVpcEndpointConnectionBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" { default = "%s" }
data "alicloud_nlb_zones" "default" {}
resource "alicloud_vpc" "default" {
  description = "test-terraform"
  cidr_block  = "10.0.0.0/8"
  vpc_name    = var.name
}
resource "alicloud_vswitch" "default_a" {
  vpc_id     = alicloud_vpc.default.id
  zone_id    = data.alicloud_nlb_zones.default.zones.0.id
  cidr_block = "10.1.0.0/16"
}
resource "alicloud_security_group" "default" {
  name        = var.name
  description = "privatelink test security group"
  vpc_id      = alicloud_vpc.default.id
}
resource "alicloud_privatelink_vpc_endpoint_service" "default" {
  service_description    = "test for privatelink connection"
  connect_bandwidth      = 3072
  auto_accept_connection = false
  service_resource_type  = "nlb"
}
resource "alicloud_vpc" "defaultService" {
  description = "test-terraform-service"
  cidr_block  = "172.16.0.0/12"
  vpc_name    = var.name
}
resource "alicloud_vswitch" "defaultService_a" {
  vpc_id     = alicloud_vpc.defaultService.id
  zone_id    = data.alicloud_nlb_zones.default.zones.0.id
  cidr_block = "172.16.1.0/24"
}
resource "alicloud_vswitch" "defaultService_b" {
  vpc_id     = alicloud_vpc.defaultService.id
  zone_id    = data.alicloud_nlb_zones.default.zones.1.id
  cidr_block = "172.16.2.0/24"
}
resource "alicloud_nlb_load_balancer" "default" {
  load_balancer_name = var.name
  vpc_id       = alicloud_vpc.defaultService.id
  address_type = "Intranet"
  zone_mappings {
    vswitch_id = alicloud_vswitch.defaultService_a.id
    zone_id    = data.alicloud_nlb_zones.default.zones.0.id
  }
  zone_mappings {
    vswitch_id = alicloud_vswitch.defaultService_b.id
    zone_id    = data.alicloud_nlb_zones.default.zones.1.id
  }
}
resource "alicloud_privatelink_vpc_endpoint_service_resource" "default" {
  zone_id       = data.alicloud_nlb_zones.default.zones.0.id
  resource_id   = alicloud_nlb_load_balancer.default.id
  resource_type = "nlb"
  service_id    = alicloud_privatelink_vpc_endpoint_service.default.id
}
resource "alicloud_privatelink_vpc_endpoint" "default" {
  service_id         = alicloud_privatelink_vpc_endpoint_service.default.id
  vpc_id             = alicloud_vpc.default.id
  security_group_ids = [alicloud_security_group.default.id]
  vpc_endpoint_name  = "testformaintf"
  depends_on         = [alicloud_privatelink_vpc_endpoint_service_resource.default]
}
`, name)
}

// Test PrivateLink VpcEndpointConnection. >>> Resource test cases, automatically generated.
// Case 4902
func TestAccAliCloudPrivateLinkVpcEndpointConnection_basic4902(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_privatelink_vpc_endpoint_connection.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudPrivateLinkVpcEndpointConnectionMap4902)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return privatelink.NewPrivateLinkServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribePrivateLinkVpcEndpointConnection")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sprivatelinkvpcendpointconnection%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudPrivateLinkVpcEndpointConnectionBasicDependence4902)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"endpoint_id": "${alicloud_privatelink_vpc_endpoint.defaultEndpoint.id}",
					"service_id":  "${alicloud_privatelink_vpc_endpoint_service.defaultService.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"endpoint_id": acctest.CHECKSET,
						"service_id":  acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bandwidth": "3072",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth": "3072",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bandwidth": "4096",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth": "4096",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bandwidth": "3072",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth": "3072",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"endpoint_id":          "${alicloud_privatelink_vpc_endpoint.defaultEndpoint.id}",
					"service_id":           "${alicloud_privatelink_vpc_endpoint_service.defaultService.id}",
					"bandwidth":            "3072",
					"traffic_control_mode": "BandwidthLimit",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"endpoint_id":          acctest.CHECKSET,
						"service_id":           acctest.CHECKSET,
						"bandwidth":            "3072",
						"traffic_control_mode": "BandwidthLimit",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dry_run"},
			},
		},
	})
}

var AlicloudPrivateLinkVpcEndpointConnectionMap4902 = map[string]string{
	"endpoint_owner_id": acctest.CHECKSET,
	"status":            acctest.CHECKSET,
}

func AlicloudPrivateLinkVpcEndpointConnectionBasicDependence4902(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "region_id" {
  default = "eu-west-1"
}

resource "alicloud_vpc" "defaultVpc" {
  description = "test-zejun"
  cidr_block  = "10.0.0.0/8"
  vpc_name    = var.name

}

resource "alicloud_security_group" "defaultSg" {
  name = var.name

  vpc_id = alicloud_vpc.defaultVpc.id
}

resource "alicloud_privatelink_vpc_endpoint_service" "defaultService" {
  service_description   = "test-zejun"
  connect_bandwidth     = "3072"
  service_resource_type = "slb"
}

resource "alicloud_privatelink_vpc_endpoint" "defaultEndpoint" {
  vpc_id     = alicloud_vpc.defaultVpc.id
  service_id = alicloud_privatelink_vpc_endpoint_service.defaultService.id
  security_group_ids = [alicloud_security_group.defaultSg.id]
}


`, name)
}

// Case 4902  twin
func TestAccAliCloudPrivateLinkVpcEndpointConnection_basic4902_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_privatelink_vpc_endpoint_connection.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudPrivateLinkVpcEndpointConnectionMap4902)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return privatelink.NewPrivateLinkServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribePrivateLinkVpcEndpointConnection")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sprivatelinkvpcendpointconnection%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudPrivateLinkVpcEndpointConnectionBasicDependence4902)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"endpoint_id":          "${alicloud_privatelink_vpc_endpoint.defaultEndpoint.id}",
					"service_id":           "${alicloud_privatelink_vpc_endpoint_service.defaultService.id}",
					"bandwidth":            "3072",
					"traffic_control_mode": "BandwidthLimit",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"endpoint_id":          acctest.CHECKSET,
						"service_id":           acctest.CHECKSET,
						"bandwidth":            "3072",
						"traffic_control_mode": "BandwidthLimit",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dry_run"},
			},
		},
	})
}

// Test PrivateLink VpcEndpointConnection. <<< Resource test cases, automatically generated.
