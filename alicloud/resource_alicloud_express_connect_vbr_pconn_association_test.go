package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Case 1
func TestAccAliCloudExpressConnectVbrPconnAssociation_basic2042(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_express_connect_vbr_pconn_association.default"
	checkoutSupportedRegions(t, true, connectivity.VbrSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudExpressConnectVbrPconnAssociationMap2042)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ExpressConnectService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeExpressConnectVbrPconnAssociation")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testaccExpressConnectVbrPconnAssociation%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudExpressConnectVbrPconnAssociationBasicDependence2042)
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
					"peer_gateway_ip":          "10.0.0.6",
					"local_gateway_ip":         "10.0.0.5",
					"physical_connection_id":   "${data.alicloud_express_connect_physical_connections.nameRegex.connections.1.id}",
					"vbr_id":                   "${alicloud_express_connect_virtual_border_router.default.id}",
					"peering_subnet_mask":      "255.255.255.252",
					"vlan_id":                  "1122",
					"enable_ipv6":              "true",
					"local_ipv6_gateway_ip":    "2408:4004:cc::3",
					"peer_ipv6_gateway_ip":     "2408:4004:cc::4",
					"peering_ipv6_subnet_mask": "2408:4004:cc::/56",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"peer_gateway_ip":          "10.0.0.6",
						"local_gateway_ip":         "10.0.0.5",
						"physical_connection_id":   CHECKSET,
						"vbr_id":                   CHECKSET,
						"peering_subnet_mask":      "255.255.255.252",
						"vlan_id":                  "1122",
						"enable_ipv6":              "true",
						"local_ipv6_gateway_ip":    "2408:4004:cc::3",
						"peer_ipv6_gateway_ip":     "2408:4004:cc::4",
						"peering_ipv6_subnet_mask": "2408:4004:cc::/56",
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

var AlicloudExpressConnectVbrPconnAssociationMap2042 = map[string]string{}

func AlicloudExpressConnectVbrPconnAssociationBasicDependence2042(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_express_connect_physical_connections" "nameRegex" {
  name_regex = "^preserved-NODELETING"
}

resource "alicloud_express_connect_virtual_border_router" "default" {
  local_gateway_ip           = "10.0.0.1"
  peer_gateway_ip            = "10.0.0.2"
  peering_subnet_mask        = "255.255.255.252"
  physical_connection_id     = data.alicloud_express_connect_physical_connections.nameRegex.connections.0.id
  virtual_border_router_name = var.name
  vlan_id                    = 110
  min_rx_interval            = 1000
  min_tx_interval            = 1000
  detect_multiplier          = 10
  enable_ipv6                = true
  local_ipv6_gateway_ip      = "2408:4004:cc:400::1"
  peer_ipv6_gateway_ip       = "2408:4004:cc:400::2"
  peering_ipv6_subnet_mask   = "2408:4004:cc:400::/56"
}

`, name)
}
