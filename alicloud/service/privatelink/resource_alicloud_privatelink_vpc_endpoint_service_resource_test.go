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

// Currently, Private network slb can only be created through the console.
func TestAccAliCloudPrivateLinkVpcEndpointServiceResource_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_privatelink_vpc_endpoint_service_resource.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudPrivatelinkVpcEndpointServiceResourceMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return privatelink.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribePrivatelinkVpcEndpointServiceResource")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, "", AlicloudPrivatelinkVpcEndpointServiceResourceBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.SlbPrivateNetSupportRegions)
		},

		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"service_id":    "${alicloud_privatelink_vpc_endpoint_service.default.id}",
					"resource_id":   "${alicloud_slb_load_balancer.default.id}",
					"resource_type": "slb",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"service_id":    acctest.CHECKSET,
						"resource_id":   acctest.CHECKSET,
						"resource_type": "slb",
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

var AlicloudPrivatelinkVpcEndpointServiceResourceMap = map[string]string{}

func AlicloudPrivatelinkVpcEndpointServiceResourceBasicDependence(name string) string {
	return fmt.Sprintf(`
		variable "name" {
		  default = "%s"
		}

		data "alicloud_slb_zones" "default" {}

		resource "alicloud_vpc" "default" {
		  description = "test-terraform-service"
		  cidr_block  = "10.0.0.0/8"
		  vpc_name    = var.name
		}

		resource "alicloud_vswitch" "default" {
		  vpc_id     = alicloud_vpc.default.id
		  zone_id    = data.alicloud_slb_zones.default.zones.0.id
		  cidr_block = "10.1.0.0/16"
		}

		resource "alicloud_privatelink_vpc_endpoint_service" "default" {
		  service_description   = "test-zejun"
		  service_resource_type = "slb"
	      auto_accept_connection = false
		}

		resource "alicloud_slb_load_balancer" "default" {
		  load_balancer_name = "${var.name}"
		  load_balancer_spec  = "slb.lcu.elastic"
	      address_type = "intranet"
	      instance_charge_type = "PayByCLCU"
	      # eu-west-1 defaults intranet SLB delete_protection to "on". The SLB
	      # resource Create only sends DeleteProtection when d.GetOk("delete_protection")
	      # is true; "off" is the schema Default, so setting "off" does not send it
	      # (GetOk returns false), the server defaults to "on", and the refresh plan
	      # drifts delete_protection on->off (observed in task 21825). Setting "on"
	      # (non-default) makes GetOk return true, so Create sends DeleteProtection="on";
	      # the SLB is "on", matching config -> empty refresh plan. (basic4 sets no
	      # delete_protection in a default region where the server defaults to "off";
	      # eu-west-1 defaults to "on", so "on" is the matching value here.)
	      delete_protection = "on"
	      vswitch_id = alicloud_vswitch.default.id
	      master_zone_id = data.alicloud_slb_zones.default.zones.0.id
	      slave_zone_id = data.alicloud_slb_zones.default.zones.1.id
		}
`, name)
}

// Test PrivateLink VpcEndpointServiceResource. >>> Resource test cases, automatically generated.
// Case 生命周期测试-VpcEndpointServiceResource_gwlb 9625
func TestAccAliCloudPrivateLinkVpcEndpointServiceResource_basic9625(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_privatelink_vpc_endpoint_service_resource.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudPrivateLinkVpcEndpointServiceResourceMap9625)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return privatelink.NewPrivateLinkServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribePrivateLinkVpcEndpointServiceResource")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sprivatelinkvpcendpointserviceresource%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudPrivateLinkVpcEndpointServiceResourceBasicDependence9625)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-wulanchabu"})
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_id":   "${alicloud_gwlb_load_balancer.defaultGllPJd.id}",
					"resource_type": "gwlb",
					"service_id":    "${alicloud_privatelink_vpc_endpoint_service.defaultQtVkqH.id}",
					"zone_id":       "${var.zone_id_1}",
					"dry_run":       "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_id":   acctest.CHECKSET,
						"resource_type": "gwlb",
						"service_id":    acctest.CHECKSET,
						"zone_id":       acctest.CHECKSET,
						"dry_run":       "false",
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

var AlicloudPrivateLinkVpcEndpointServiceResourceMap9625 = map[string]string{
	"region_id": acctest.CHECKSET,
}

func AlicloudPrivateLinkVpcEndpointServiceResourceBasicDependence9625(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "zone_id_1" {
  default = "cn-wulanchabu-b"
}

variable "region_id" {
  default = "cn-wulanchabu"
}

resource "alicloud_vpc" "defaultvVpc" {
  description = "test"
  cidr_block  = "10.0.0.0/8"
  vpc_name    = var.name
}

resource "alicloud_vswitch" "defaultVSwitch1" {
  vpc_id     = alicloud_vpc.defaultvVpc.id
  zone_id    = var.zone_id_1
  cidr_block = "10.2.0.0/16"
}

resource "alicloud_gwlb_load_balancer" "defaultGllPJd" {
  load_balancer_name = format("%%s2", var.name)
  address_ip_version = "Ipv4"
  vpc_id             = alicloud_vpc.defaultvVpc.id
  zone_mappings {
    vswitch_id = alicloud_vswitch.defaultVSwitch1.id
    zone_id    = var.zone_id_1
  }
}

resource "alicloud_privatelink_vpc_endpoint_service" "defaultQtVkqH" {
  service_description   = "test-lengqing"
  service_resource_type = "gwlb"
}


`, name)
}

// Test PrivateLink VpcEndpointServiceResource. <<< Resource test cases, automatically generated.
