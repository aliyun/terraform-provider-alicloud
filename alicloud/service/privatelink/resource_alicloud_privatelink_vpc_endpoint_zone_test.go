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
func TestAccAliCloudPrivatelinkVpcEndpointZone_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_privatelink_vpc_endpoint_zone.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudPrivatelinkVpcEndpointZoneMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return privatelink.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribePrivatelinkVpcEndpointZone")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccPrivatelinkVpcEndpointZoneTest%d", rand)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudPrivatelinkVpcEndpointZoneBasicDependence)
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
					"endpoint_id": "${alicloud_privatelink_vpc_endpoint.default.id}",
					"vswitch_id":  "${alicloud_vswitch.default.id}",
					"zone_id":     "${alicloud_vswitch.default.zone_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
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

var AlicloudPrivatelinkVpcEndpointZoneMap = map[string]string{
	"status": acctest.CHECKSET,
}

func AlicloudPrivatelinkVpcEndpointZoneBasicDependence(name string) string {
	return fmt.Sprintf(`
	data "alicloud_slb_zones" "default" {}

	resource "alicloud_vpc" "default" {
	  description = "test-terraform-service"
	  cidr_block  = "10.0.0.0/8"
	  vpc_name    = "%[1]s"
	}

	resource "alicloud_vswitch" "default" {
	  vpc_id     = alicloud_vpc.default.id
	  zone_id    = data.alicloud_slb_zones.default.zones.0.id
	  cidr_block = "10.1.0.0/16"
	}

	resource "alicloud_slb_load_balancer" "default" {
	  load_balancer_name = "%[1]s"
	  load_balancer_spec  = "slb.s2.small"
      address_type = "intranet"
      instance_charge_type = "PayBySpec"
      delete_protection = "off"
      vswitch_id = alicloud_vswitch.default.id
      master_zone_id = data.alicloud_slb_zones.default.zones.0.id
      slave_zone_id = data.alicloud_slb_zones.default.zones.1.id
      lifecycle {
        ignore_changes = [delete_protection]
      }
	}

	data "alicloud_vswitches" "default" {
	 	 is_default = true
	}
	resource "alicloud_security_group" "default" {
	 name = "%[1]s"
	 description = "privatelink test security group"
	 vpc_id = alicloud_vpc.default.id
	}
	resource "alicloud_privatelink_vpc_endpoint_service" "default" {
	  service_description   = "test for privatelink connection"
	  service_resource_type = "slb"
      auto_accept_connection = false
	}
	resource "alicloud_privatelink_vpc_endpoint_service_resource" "default" {
	 service_id    =  "${alicloud_privatelink_vpc_endpoint_service.default.id}"
	 resource_id   =  "${alicloud_slb_load_balancer.default.id}"
	 resource_type = "slb"
	}
	resource "alicloud_privatelink_vpc_endpoint" "default" {
	 service_id = alicloud_privatelink_vpc_endpoint_service_resource.default.service_id
	 vpc_id = alicloud_vpc.default.id
	 security_group_ids = [alicloud_security_group.default.id]
	 vpc_endpoint_name = "%[1]s"
	 depends_on = [alicloud_privatelink_vpc_endpoint_service.default]
	}
`, name)
}

func TestAccAliCloudPrivateLinkVpcEndpointZone_basic_without_zoneid(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_privatelink_vpc_endpoint_zone.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudPrivateLinkVpcEndpointZoneMap4898)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return privatelink.NewPrivateLinkServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribePrivateLinkVpcEndpointZone")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sprivatelinkvpcendpointzone%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudPrivateLinkVpcEndpointZoneBasicDependence4898)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"endpoint_id": "${alicloud_privatelink_vpc_endpoint.defaulti9F95i.id}",
					"vswitch_id":  "${alicloud_vswitch.defaultmEwUAc.id}",
					"eni_ip":      "10.1.0.245",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"zone_id":     acctest.CHECKSET,
						"endpoint_id": acctest.CHECKSET,
						"vswitch_id":  acctest.CHECKSET,
						"eni_ip":      "10.1.0.245",
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

func TestAccAliCloudPrivateLinkVpcEndpointZone_basic4898(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_privatelink_vpc_endpoint_zone.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudPrivateLinkVpcEndpointZoneMap4898)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return privatelink.NewPrivateLinkServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribePrivateLinkVpcEndpointZone")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sprivatelinkvpcendpointzone%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudPrivateLinkVpcEndpointZoneBasicDependence4898)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"zone_id":     "${data.alicloud_nlb_zones.default.zones.0.id}",
					"endpoint_id": "${alicloud_privatelink_vpc_endpoint.defaulti9F95i.id}",
					"vswitch_id":  "${alicloud_vswitch.defaultmEwUAc.id}",
					"eni_ip":      "10.1.0.245",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"zone_id":     acctest.CHECKSET,
						"endpoint_id": acctest.CHECKSET,
						"vswitch_id":  acctest.CHECKSET,
						"eni_ip":      "10.1.0.245",
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

var AlicloudPrivateLinkVpcEndpointZoneMap4898 = map[string]string{
	"status": acctest.CHECKSET,
}

func AlicloudPrivateLinkVpcEndpointZoneBasicDependence4898(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_nlb_zones" "default" {}

resource "alicloud_vpc" "defaultbFzA4a" {
  description = "test-terraform"
  cidr_block  = "10.0.0.0/8"
  vpc_name    = var.name

}

resource "alicloud_vswitch" "defaultmEwUAc" {
  vpc_id     = alicloud_vpc.defaultbFzA4a.id
  zone_id    = data.alicloud_nlb_zones.default.zones.0.id
  cidr_block = "10.1.0.0/16"
}

resource "alicloud_security_group" "default1FTFrP" {
  name = var.name
  vpc_id = alicloud_vpc.defaultbFzA4a.id
}

resource "alicloud_privatelink_vpc_endpoint_service" "defaultr0WBYX" {
  service_description   = "test-zejun-service"
  connect_bandwidth     = "3072"
  service_resource_type = "nlb"
}

resource "alicloud_privatelink_vpc_endpoint" "defaulti9F95i" {
  vpc_id     = alicloud_vpc.defaultbFzA4a.id
  service_id = alicloud_privatelink_vpc_endpoint_service.defaultr0WBYX.id
  security_group_ids = [alicloud_security_group.default1FTFrP.id]
  # The NLB-backed endpoint service only reports a zone as supported after the
  # service_resource (the NLB registration) is created. Without this dependency
  # the endpoint and service_resource race in parallel, and AddZoneToVpcEndpoint
  # can run before the zone is propagated, failing with
  # EndpointServiceNotSupportedZone. The service_resource Create waits on its own
  # state refresh, so depending on it here lets endpoint_zone inherit the ordering.
  depends_on = [alicloud_privatelink_vpc_endpoint_service_resource.defaultdTPOne]
}

resource "alicloud_vpc" "defaultVpcService" {
  description = "test-terraform-service"
  cidr_block  = "10.0.0.0/8"
  vpc_name    = var.name

}

resource "alicloud_vswitch" "defaultuYH1VC" {
  vpc_id     = alicloud_vpc.defaultVpcService.id
  zone_id    = data.alicloud_nlb_zones.default.zones.0.id
  cidr_block = "10.1.0.0/16"
}

resource "alicloud_vswitch" "defaultTJZ8ud" {
  vpc_id     = alicloud_vpc.defaultVpcService.id
  zone_id    = data.alicloud_nlb_zones.default.zones.1.id
  cidr_block = "10.10.0.0/16"
}

resource "alicloud_nlb_load_balancer" "defaultyuY5jZ" {
  zone_mappings {
    vswitch_id = alicloud_vswitch.defaultuYH1VC.id
    zone_id    = data.alicloud_nlb_zones.default.zones.0.id
  }
  zone_mappings {
    vswitch_id = alicloud_vswitch.defaultTJZ8ud.id
    zone_id    = data.alicloud_nlb_zones.default.zones.1.id
  }
  load_balancer_name = var.name

  vpc_id       = alicloud_vpc.defaultVpcService.id
  address_type = "Intranet"
}

resource "alicloud_privatelink_vpc_endpoint_service_resource" "defaultdTPOne" {
  zone_id       = data.alicloud_nlb_zones.default.zones.0.id
  resource_id   = alicloud_nlb_load_balancer.defaultyuY5jZ.id
  resource_type = "nlb"
  service_id    = alicloud_privatelink_vpc_endpoint_service.defaultr0WBYX.id
}


`, name)
}
