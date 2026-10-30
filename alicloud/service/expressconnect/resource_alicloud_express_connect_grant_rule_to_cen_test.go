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

// Test ExpressConnect GrantRuleToCen. >>> Resource test cases, automatically generated.
// Case 高速通道跨账号授权_副本1758786253950 11567
func TestAccAliCloudExpressConnectGrantRuleToCen_basic11567(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_express_connect_grant_rule_to_cen.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudExpressConnectGrantRuleToCenMap11567)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return expressconnect.NewExpressConnectServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeExpressConnectGrantRuleToCen")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccexpressconnect%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudExpressConnectGrantRuleToCenBasicDependence11567)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":  "${alicloud_express_connect_virtual_border_router.default.id}",
					"cen_owner_id": "${data.alicloud_account.default.id}",
					"cen_id":       "${alicloud_cen_instance.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":  acctest.CHECKSET,
						"cen_owner_id": acctest.CHECKSET,
						"cen_id":       acctest.CHECKSET,
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

var AliCloudExpressConnectGrantRuleToCenMap11567 = map[string]string{
	"create_time": acctest.CHECKSET,
}

func AliCloudExpressConnectGrantRuleToCenBasicDependence11567(name string) string {
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
`, name, tfacctest.RandIntRange(1, 2999))
}

// Test ExpressConnect GrantRuleToCen. <<< Resource test cases, automatically generated.
