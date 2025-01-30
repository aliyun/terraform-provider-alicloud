package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudExpressConnectGrantRuleToCen_basic0(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.VbrSupportRegions)
	resourceId := "alicloud_express_connect_grant_rule_to_cen.default"
	ra := resourceAttrInit(resourceId, resourceAlicloudExpressConnectGrantRuleToCenMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeExpressConnectGrantRuleToCen")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccExpressConnectGrantRuleToCen-name%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceAlicloudExpressConnectGrantRuleToCenBasicDependence)
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
					"cen_id":       "${alicloud_cen_instance.default.id}",
					"cen_owner_id": "${data.alicloud_account.default.id}",
					"instance_id":  "${alicloud_express_connect_virtual_border_router.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cen_id":       CHECKSET,
						"cen_owner_id": CHECKSET,
						"instance_id":  CHECKSET,
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

var resourceAlicloudExpressConnectGrantRuleToCenMap = map[string]string{}

func resourceAlicloudExpressConnectGrantRuleToCenBasicDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_account" "default" {
	}

	data "alicloud_express_connect_physical_connections" "default" {
  		name_regex = "^preserved-NODELETING"
	}

	resource "alicloud_cen_instance" "default" {
  		cen_instance_name = var.name
	}

	resource "alicloud_express_connect_virtual_border_router" "default" {
  		local_gateway_ip       = "10.0.0.1"
  		peer_gateway_ip        = "10.0.0.2"
  		peering_subnet_mask    = "255.255.255.252"
  		physical_connection_id = data.alicloud_express_connect_physical_connections.default.connections.0.id
  		vlan_id                = %d
  		min_rx_interval        = 1000
  		min_tx_interval        = 1000
  		detect_multiplier      = 10
	}
`, name, acctest.RandIntRange(1, 2999))
}
