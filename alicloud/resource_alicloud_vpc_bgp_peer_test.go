package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudVPCBgpPeer_basic0(t *testing.T) {
	checkoutSupportedRegions(t, true, connectivity.VbrSupportRegions)
	var v map[string]interface{}
	resourceId := "alicloud_vpc_bgp_peer.default"
	checkoutSupportedRegions(t, true, connectivity.VPCBgpGroupSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudVPCBgpPeerMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVpcBgpPeer")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%svpcbgppeer%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudVPCBgpPeerBasicDependence0)
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
					"peer_ip_address": "1.1.1.1",
					"bgp_group_id":    "${alicloud_vpc_bgp_group.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"peer_ip_address": "1.1.1.1",
						"bgp_group_id":    CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"peer_ip_address": "1.1.1.2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"peer_ip_address": "1.1.1.2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_bfd":    "true",
					"bfd_multi_hop": "10",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_bfd":    "true",
						"bfd_multi_hop": "10",
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
func TestAccAlicloudVPCBgpPeer_basic1(t *testing.T) {
	checkoutSupportedRegions(t, true, connectivity.VbrSupportRegions)
	var v map[string]interface{}
	resourceId := "alicloud_vpc_bgp_peer.default"
	checkoutSupportedRegions(t, true, connectivity.VPCBgpGroupSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudVPCBgpPeerMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVpcBgpPeer")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%svpcbgppeer%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudVPCBgpPeerBasicDependence0)
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
					"peer_ip_address": "1.1.1.1",
					"enable_bfd":      "true",
					"bfd_multi_hop":   "10",
					"ip_version":      "IPV4",
					"bgp_group_id":    "${alicloud_vpc_bgp_group.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"peer_ip_address": "1.1.1.1",
						"enable_bfd":      "true",
						"bfd_multi_hop":   "10",
						"ip_version":      "IPV4",
						"bgp_group_id":    CHECKSET,
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

var AlicloudVPCBgpPeerMap0 = map[string]string{
	"status": CHECKSET,
}

func AlicloudVPCBgpPeerBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}

data "alicloud_express_connect_physical_connections" "default" {
	name_regex = "^preserved-NODELETING"
}

resource "alicloud_express_connect_virtual_border_router" "default" {
  local_gateway_ip           = "10.0.0.1"
  peer_gateway_ip            = "10.0.0.2"
  peering_subnet_mask        = "255.255.255.252"
  physical_connection_id     = data.alicloud_express_connect_physical_connections.default.connections.0.id
  virtual_border_router_name = var.name
  vlan_id                    = %d
  min_rx_interval            = 1000
  min_tx_interval            = 1000
  detect_multiplier          = 10
}

resource "alicloud_vpc_bgp_group" "default" {
  auth_key       = "YourPassword+12345678"
  bgp_group_name = var.name
  description    = var.name
  local_asn      = 64512
  peer_asn       = 1111
  router_id      = alicloud_express_connect_virtual_border_router.default.id
}
`, name, acctest.RandIntRange(1, 2999))
}
