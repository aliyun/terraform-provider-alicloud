package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudCloudfirewallVpcFirewall_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_firewall_vpc_firewall.default"
	ra := resourceAttrInit(resourceId, AlicloudCloudfirewallVpcFirewallMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudfwService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudFirewallVpcFirewall")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	checkoutSupportedRegions(t, true, connectivity.CloudFirewallVpcFirewallSupportRegions)
	name := fmt.Sprintf("tf-testacc%sCloudfirewallVpcFirewall%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCloudfirewallVpcFirewallBasicDependence)
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
					"vpc_firewall_name": name,
					"member_uid":        "${data.alicloud_account.current.id}",
					"peer_vpc": []map[string]interface{}{
						{
							"vpc_id":    "${data.alicloud_vpcs.vpcs_ds_peer.vpcs.0.id}",
							"region_no": defaultRegionToTest,
							"peer_vpc_cidr_table_list": []map[string]interface{}{
								{
									"peer_route_table_id": "${data.alicloud_route_tables.local_peer.tables.0.id}",
									"peer_route_entry_list": []map[string]interface{}{
										{
											"peer_destination_cidr":     "${data.alicloud_vpcs.vpcs_ds.vpcs.0.cidr_block}",
											"peer_next_hop_instance_id": "${data.alicloud_vpc_peer_connections.cfw_vpc_peer.connections.0.id}",
										},
									},
								},
							},
						},
					},
					"local_vpc": []map[string]interface{}{
						{
							"vpc_id":    "${data.alicloud_vpcs.vpcs_ds.vpcs.0.id}",
							"region_no": defaultRegionToTest,
							"local_vpc_cidr_table_list": []map[string]interface{}{
								{
									"local_route_table_id": "${data.alicloud_route_tables.local_vpc.tables.0.id}",
									"local_route_entry_list": []map[string]interface{}{
										{
											"local_destination_cidr":     "${data.alicloud_vpcs.vpcs_ds_peer.vpcs.0.cidr_block}",
											"local_next_hop_instance_id": "${data.alicloud_vpc_peer_connections.cfw_vpc_peer.connections.0.id}",
										},
									},
								},
							},
						},
					},
					"status": "open",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpc_firewall_name": CHECKSET,
						"member_uid":        CHECKSET,
						"status":            "open",
					}),
				),
			}, {
				Config: testAccConfig(map[string]interface{}{
					"vpc_firewall_name": "tf-test-update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpc_firewall_name": "tf-test-update",
					}),
				),
			}, {
				Config: testAccConfig(map[string]interface{}{
					"status": "close",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "close",
					}),
				),
			}, {
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"lang"},
			},
		},
	})
}

var AlicloudCloudfirewallVpcFirewallMap = map[string]string{}

func AlicloudCloudfirewallVpcFirewallBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}
data "alicloud_account" "current" {
}
data "alicloud_vpcs" "vpcs_ds" {
  name_regex = "^cfw-1-default-NODELETING"
}
data "alicloud_route_tables" "local_vpc" {
  vpc_id = "${data.alicloud_vpcs.vpcs_ds.vpcs.0.id}"
}
data "alicloud_vpcs" "vpcs_ds_peer" {
  name_regex = "^cfw-2-default-NODELETING"
}
data "alicloud_route_tables" "local_peer" {
  vpc_id = "${data.alicloud_vpcs.vpcs_ds_peer.vpcs.0.id}"
}
data "alicloud_vpc_peer_connections" "cfw_vpc_peer" {
  name_regex = "^cfw-default-NODELETING"
}

`, name)
}
