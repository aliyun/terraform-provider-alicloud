package expressconnect_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/expressconnect"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudExpressConnectVbrHa_basic0(t *testing.T) {
	acctest.CheckoutSupportedRegions(t, true, connectivity.VbrSupportRegions)
	var v map[string]interface{}
	resourceId := "alicloud_vpc_vbr_ha.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudVPCVbrHaMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return expressconnect.NewExpressConnectServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeVpcVbrHa")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%svpcvbrha%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudVPCVbrHaBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"depends_on":  []string{"alicloud_cen_instance_attachment.default"},
					"description": "${var.name}",
					"vbr_id":      "${alicloud_express_connect_virtual_border_router.default[0].id}",
					"peer_vbr_id": "${alicloud_express_connect_virtual_border_router.default[1].id}",
					"vbr_ha_name": "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name,
						"vbr_id":      acctest.CHECKSET,
						"peer_vbr_id": acctest.CHECKSET,
						"vbr_ha_name": name,
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

var AlicloudVPCVbrHaMap0 = map[string]string{
	"dry_run": acctest.NOSET,
	"status":  acctest.CHECKSET,
}

func AlicloudVPCVbrHaBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}

data "alicloud_express_connect_physical_connections" "default" {
  name_regex = "^preserved-NODELETING-"
}

resource "alicloud_express_connect_virtual_border_router" "default" {
  count                      = 2
  local_gateway_ip           = "10.0.0.1"
  peer_gateway_ip            = "10.0.0.2"
  peering_subnet_mask        = "255.255.255.252"
  physical_connection_id     = data.alicloud_express_connect_physical_connections.default.connections[count.index].id
  virtual_border_router_name = var.name
  vlan_id                    = %d
  min_rx_interval            = 1000
  min_tx_interval            = 1000
  detect_multiplier          = 10
}

resource "alicloud_cen_instance" "default" {
  cen_instance_name = var.name
  description       = var.name
}

resource "alicloud_cen_instance_attachment" "default" {
  count                    = 2
  instance_id              = alicloud_cen_instance.default.id
  child_instance_id        = alicloud_express_connect_virtual_border_router.default[count.index].id
  child_instance_type      = "VBR"
  child_instance_region_id = "%s"
}

`, name, tfacctest.RandIntRange(1, 2999), acctest.DefaultRegionToTest)
}
