package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAliCloudVPNGatewayVcoRoute_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vpn_gateway_vco_route.default"
	checkoutSupportedRegions(t, true, connectivity.VpnGatewayVpnAttachmentSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudVPNGatewayVcoRouteMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVpnGatewayVcoRoute")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%svpngatewayvcoroute%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudVPNGatewayVcoRouteBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"route_dest":        "192.168.10.0/24",
					"next_hop":          "${alicloud_cen_transit_router_vpn_attachment.default.vpn_id}",
					"vpn_connection_id": "${alicloud_cen_transit_router_vpn_attachment.default.vpn_id}",
					"weight":            "100",
					"overlay_mode":      "Ipsec",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"route_dest":        "192.168.10.0/24",
						"next_hop":          CHECKSET,
						"vpn_connection_id": CHECKSET,
						"weight":            "100",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"overlay_mode"},
			},
		},
	})
}

var AlicloudVPNGatewayVcoRouteMap0 = map[string]string{
	"status": CHECKSET,
}

func AlicloudVPNGatewayVcoRouteBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
resource "alicloud_cen_instance" "default" {
	cen_instance_name = var.name
}
resource "alicloud_cen_transit_router" "default" {
	cen_id = alicloud_cen_instance.default.id
	transit_router_description = "desd"
	transit_router_name = var.name
}
resource "alicloud_cen_transit_router_cidr" "default" {
  transit_router_id        = alicloud_cen_transit_router.default.transit_router_id
  cidr                     = "192.168.0.0/16"
  transit_router_cidr_name = var.name
  description              = var.name
  publish_cidr_route       = true
}

data "alicloud_cen_transit_router_available_resources" "default" {}
resource "alicloud_vpn_customer_gateway" "default" {
  name        = "${var.name}"
  ip_address  = "42.104.22.210"
  asn         = "45014"
  description = "testAccVpnConnectionDesc"
}
resource "alicloud_vpn_gateway_vpn_attachment" "default" {
  customer_gateway_id = alicloud_vpn_customer_gateway.default.id
  network_type        = "public"
  local_subnet        = "0.0.0.0/0"
  remote_subnet       = "0.0.0.0/0"
  effect_immediately  = false
  ike_config {
    ike_auth_alg = "md5"
    ike_enc_alg  = "des"
    ike_version  = "ikev2"
    ike_mode     = "main"
    ike_lifetime = 86400
    psk          = "tf-testvpn2"
    ike_pfs      = "group1"
    remote_id    = "testbob2"
    local_id     = "testalice2"
  }
  ipsec_config {
    ipsec_pfs      = "group5"
    ipsec_enc_alg  = "des"
    ipsec_auth_alg = "md5"
    ipsec_lifetime = 86400
  }
  bgp_config {
    enable       = true
    local_asn    = 45014
    tunnel_cidr  = "169.254.11.0/30"
    local_bgp_ip = "169.254.11.1"
  }
  health_check_config {
    enable   = true
    sip      = "192.168.1.1"
    dip      = "10.0.0.1"
    interval = 10
    retry    = 10
    policy   = "revoke_route"
  }
  enable_dpd           = true
  enable_nat_traversal = true
  vpn_attachment_name  = var.name
}
resource "alicloud_cen_transit_router_vpn_attachment" "default" {
	auto_publish_route_enabled = false
	transit_router_attachment_description = var.name
	transit_router_attachment_name = var.name
	cen_id = alicloud_cen_transit_router.default.cen_id
	transit_router_id = alicloud_cen_transit_router_cidr.default.transit_router_id
	vpn_id = alicloud_vpn_gateway_vpn_attachment.default.id
	zone {
		zone_id = data.alicloud_cen_transit_router_available_resources.default.resources.0.master_zones.0
	}
}
`, name)
}

// lintignore: R001
