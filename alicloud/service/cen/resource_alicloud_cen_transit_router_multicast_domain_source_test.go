package cen_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/cen"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Case 1
func TestAccAliCloudCenTransitRouterMulticastDomainSource_basic1903(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cen_transit_router_multicast_domain_source.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudCenTransitRouterMulticastDomainSourceMap1903)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cen.NewCbnService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCenTransitRouterMulticastDomainSource")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sCenTransitRouterMulticastDomainSource%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudCenTransitRouterMulticastDomainSourceBasicDependence1903)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.CenSupportRegions)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"vpc_id":                             "${alicloud_vpc.default.id}",
					"transit_router_multicast_domain_id": "${alicloud_cen_transit_router_multicast_domain_association.default.transit_router_multicast_domain_id}",
					"network_interface_id":               "${alicloud_ecs_network_interface.default.id}",
					"group_ip_address":                   "232.1.1.1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpc_id":                             acctest.CHECKSET,
						"transit_router_multicast_domain_id": acctest.CHECKSET,
						"network_interface_id":               acctest.CHECKSET,
						"group_ip_address":                   "232.1.1.1",
					}),
				),
			}, {
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudCenTransitRouterMulticastDomainSourceMap1903 = map[string]string{}

func AlicloudCenTransitRouterMulticastDomainSourceBasicDependence1903(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

variable "instance_name" {
  default = "tf-testacc-cen_instance"
}

data "alicloud_cen_transit_router_available_resources" "default" {
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.instance_name
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_vswitch" "default_master" {
  vswitch_name = var.instance_name
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "192.168.1.0/24"
  zone_id      = "cn-hangzhou-i"
}

resource "alicloud_vswitch" "default_slave" {
  vswitch_name = var.instance_name
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "192.168.2.0/24"
  zone_id      = "cn-hangzhou-j"
}

resource "alicloud_cen_instance" "default" {
  cen_instance_name = var.instance_name
  protection_level  = "REDUCED"
}

resource "alicloud_cen_transit_router" "default" {
  cen_id = alicloud_cen_instance.default.id
  support_multicast = true
}

resource "alicloud_cen_transit_router_vpc_attachment" "default" {
  cen_id            = alicloud_cen_instance.default.id
  transit_router_id = alicloud_cen_transit_router.default.transit_router_id
  vpc_id            = alicloud_vpc.default.id
  zone_mappings {
    zone_id    = alicloud_vswitch.default_master.zone_id
    vswitch_id = alicloud_vswitch.default_master.id
  }
  zone_mappings {
    zone_id    = alicloud_vswitch.default_slave.zone_id
    vswitch_id = alicloud_vswitch.default_slave.id
  }
  transit_router_attachment_name        = var.instance_name
  transit_router_attachment_description = var.instance_name
}


resource "alicloud_security_group" "default" {
    name = var.name
    vpc_id = alicloud_vpc.default.id
}

data "alicloud_resource_manager_resource_groups" "default"{
  status = "OK"
}

resource "alicloud_cen_transit_router_multicast_domain" "default" {
  depends_on =           ["alicloud_cen_transit_router_vpc_attachment.default"]
  transit_router_id                           = alicloud_cen_transit_router.default.transit_router_id
  transit_router_multicast_domain_name        = var.name
  transit_router_multicast_domain_description = var.name
}

resource "alicloud_ecs_network_interface" "default" {
    depends_on =           ["alicloud_cen_transit_router_multicast_domain.default"]
    network_interface_name = var.name
    vswitch_id = alicloud_vswitch.default_master.id
    security_group_ids = [alicloud_security_group.default.id]
  description = "Basic test"
  primary_ip_address = cidrhost(alicloud_vswitch.default_master.cidr_block, 100)
  tags = {
    Created = "TF",
    For =    "Test",
  }
  resource_group_id = data.alicloud_resource_manager_resource_groups.default.ids.0
}

resource "alicloud_cen_transit_router_multicast_domain_association" "default" {
	transit_router_multicast_domain_id = alicloud_cen_transit_router_multicast_domain.default.id
	transit_router_attachment_id       = alicloud_cen_transit_router_vpc_attachment.default.transit_router_attachment_id
	vswitch_id                         = alicloud_vswitch.default_master.id
}

`, name)
}

// Case 2
func TestAccAliCloudCenTransitRouterMulticastDomainSource_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cen_transit_router_multicast_domain_source.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudCenTransitRouterMulticastDomainSourceMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cen.NewCbnService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCenTransitRouterMulticastDomainSource")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sCenTransitRouterMulticastDomainSource%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudCenTransitRouterMulticastDomainSourceBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.CenSupportRegions)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"transit_router_multicast_domain_id": "${alicloud_cen_transit_router_multicast_domain_association.default.transit_router_multicast_domain_id}",
					"network_interface_id":               "${alicloud_ecs_network_interface.default.id}",
					"group_ip_address":                   "232.1.1.1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"transit_router_multicast_domain_id": acctest.CHECKSET,
						"network_interface_id":               acctest.CHECKSET,
						"group_ip_address":                   "232.1.1.1",
					}),
				),
			}, {
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudCenTransitRouterMulticastDomainSourceMap0 = map[string]string{}

func AlicloudCenTransitRouterMulticastDomainSourceBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "instance_name" {
  default = "tf-testacc-cen_instance"
}

data "alicloud_cen_transit_router_available_resources" "default" {
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.instance_name
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_vswitch" "default_master" {
  vswitch_name = var.instance_name
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "192.168.1.0/24"
  zone_id      = "cn-hangzhou-i"
}

resource "alicloud_vswitch" "default_slave" {
  vswitch_name = var.instance_name
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "192.168.2.0/24"
  zone_id      = "cn-hangzhou-j"
}

resource "alicloud_cen_instance" "default" {
  cen_instance_name = var.instance_name
  protection_level  = "REDUCED"
}

resource "alicloud_cen_transit_router" "default" {
  cen_id = alicloud_cen_instance.default.id
  support_multicast = true
}

resource "alicloud_cen_transit_router_vpc_attachment" "default" {
  cen_id            = alicloud_cen_instance.default.id
  transit_router_id = alicloud_cen_transit_router.default.transit_router_id
  vpc_id            = alicloud_vpc.default.id
  zone_mappings {
    zone_id    = alicloud_vswitch.default_master.zone_id
    vswitch_id = alicloud_vswitch.default_master.id
  }
  zone_mappings {
    zone_id    = alicloud_vswitch.default_slave.zone_id
    vswitch_id = alicloud_vswitch.default_slave.id
  }
  transit_router_attachment_name        = var.instance_name
  transit_router_attachment_description = var.instance_name
}


resource "alicloud_security_group" "default" {
    name = var.name
    vpc_id = alicloud_vpc.default.id
}

data "alicloud_resource_manager_resource_groups" "default"{
  status = "OK"
}

resource "alicloud_cen_transit_router_multicast_domain" "default" {
  depends_on =           ["alicloud_cen_transit_router_vpc_attachment.default"]
  transit_router_id                           = alicloud_cen_transit_router.default.transit_router_id
  transit_router_multicast_domain_name        = var.name
  transit_router_multicast_domain_description = var.name
}

resource "alicloud_ecs_network_interface" "default" {
    depends_on =           ["alicloud_cen_transit_router_multicast_domain.default"]
    network_interface_name = var.name
    vswitch_id = alicloud_vswitch.default_master.id
    security_group_ids = [alicloud_security_group.default.id]
  description = "Basic test"
  primary_ip_address = cidrhost(alicloud_vswitch.default_master.cidr_block, 100)
  tags = {
    Created = "TF",
    For =    "Test",
  }
  resource_group_id = data.alicloud_resource_manager_resource_groups.default.ids.0
}

resource "alicloud_cen_transit_router_multicast_domain_association" "default" {
	transit_router_multicast_domain_id = alicloud_cen_transit_router_multicast_domain.default.id
	transit_router_attachment_id       = alicloud_cen_transit_router_vpc_attachment.default.transit_router_attachment_id
	vswitch_id                         = alicloud_vswitch.default_master.id
}

`, name)
}
