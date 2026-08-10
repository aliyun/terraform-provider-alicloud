package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudEnsNetworkPeerConnection_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ens_network_peer_connection.default"
	ra := resourceAttrInit(resourceId, AlicloudEnsNetworkPeerConnectionMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EnsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEnsNetworkPeerConnection")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc-ensnpc-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEnsNetworkPeerConnectionBasicDependence)

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
					"network_id":                   "${alicloud_ens_network.default.id}",
					"accepting_network_id":         "${alicloud_ens_network.accepting.id}",
					"network_peer_connection_name": name,
					"description":                  "description_test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"network_peer_connection_name": name,
						"description":                  "description_test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"network_peer_connection_name": name + "_update",
					"description":                  "description_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"network_peer_connection_name": name + "_update",
						"description":                  "description_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "description_updated_again",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "description_updated_again",
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

var AlicloudEnsNetworkPeerConnectionMap = map[string]string{
	"status":               CHECKSET,
	"create_time":          CHECKSET,
	"instance_id":          CHECKSET,
	"ens_region_id":        CHECKSET,
	"network_id":           CHECKSET,
	"accepting_network_id": CHECKSET,
}

func AlicloudEnsNetworkPeerConnectionBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_ens_network" "default" {
  network_name  = "${var.name}-net1"
  description   = var.name
  cidr_block    = "192.168.2.0/24"
  ens_region_id = "cn-chenzhou-telecom_unicom_cmcc"
}

resource "alicloud_ens_network" "accepting" {
  network_name  = "${var.name}-net2"
  description   = var.name
  cidr_block    = "192.168.3.0/24"
  ens_region_id = "cn-chenzhou-telecom_unicom_cmcc"
}

`, name)
}
