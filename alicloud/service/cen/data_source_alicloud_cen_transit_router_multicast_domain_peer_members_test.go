package cen_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAlicloudCenTransitRouterMulticastDomainPeerMemberDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(1000000, 9999999)
	acctest.CheckoutSupportedRegions(t, true, connectivity.CENTransitRouterMulticastDomainPeerMemberSupportRegions)

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudCenTransitRouterMulticastDomainPeerMemberSourceConfig(rand, map[string]string{
			"ids":                                `["${alicloud_cen_transit_router_multicast_domain_peer_member.default.id}"]`,
			"transit_router_multicast_domain_id": `"${alicloud_cen_transit_router_multicast_domain_peer_member.default.transit_router_multicast_domain_id}"`,
		}),
		FakeConfig: testAccCheckAlicloudCenTransitRouterMulticastDomainPeerMemberSourceConfig(rand, map[string]string{
			"ids":                                `["${alicloud_cen_transit_router_multicast_domain_peer_member.default.id}_fake"]`,
			"transit_router_multicast_domain_id": `"${alicloud_cen_transit_router_multicast_domain_peer_member.default.transit_router_multicast_domain_id}"`,
		}),
	}
	existUpdate := acctest.ResourceAttrInit(CenTransitRouterMulticastDomainPeerMemberCheckInfo.ResourceId, CenTransitRouterMulticastDomainPeerMemberCheckInfo.ExistMapFunc(rand)).ResourceAttrMapUpdateSet()
	emptyUpdate := acctest.ResourceAttrInit(CenTransitRouterMulticastDomainPeerMemberCheckInfo.ResourceId, CenTransitRouterMulticastDomainPeerMemberCheckInfo.FakeMapFunc(rand)).ResourceAttrMapUpdateSet()
	var steps []resource.TestStep
	if allConf.ExistConfig != "" {
		steps = append(steps, resource.TestStep{Config: allConf.ExistConfig, Check: resource.ComposeTestCheckFunc(existUpdate(allConf.ExistChangMap))})
	}
	if allConf.FakeConfig != "" {
		steps = append(steps, resource.TestStep{Config: allConf.FakeConfig, Check: resource.ComposeTestCheckFunc(emptyUpdate(allConf.FakeChangMap))})
	}

	// multi provideris
	var providers []*schema.Provider
	providerFactories := map[string]func() (*schema.Provider, error){
		"alicloud": func() (*schema.Provider, error) {
			p := acctest.Provider
			providers = append(providers, p)
			return p, nil
		},
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithAccountSiteType(t, acctest.DomesticSite)
		},

		// module name
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckCenInterRegionTransitRouterMulticastDomainPeerMemberDestroyWithProviders(&providers),
		Steps:             steps,
	})
}

var existCenTransitRouterMulticastDomainPeerMemberMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"members.#":                  "1",
		"members.0.id":               acctest.CHECKSET,
		"members.0.group_ip_address": acctest.CHECKSET,
		"members.0.peer_transit_router_multicast_domain_id": acctest.CHECKSET,
		"members.0.transit_router_multicast_domain_id":      acctest.CHECKSET,
		"members.0.status": acctest.CHECKSET,
	}
}

var fakeCenTransitRouterMulticastDomainPeerMemberMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"members.#": "0",
	}
}

var CenTransitRouterMulticastDomainPeerMemberCheckInfo = acctest.DataSourceAttr{
	ResourceId:   "data.alicloud_cen_transit_router_multicast_domain_peer_members.default",
	ExistMapFunc: existCenTransitRouterMulticastDomainPeerMemberMapFunc,
	FakeMapFunc:  fakeCenTransitRouterMulticastDomainPeerMemberMapFunc,
}

func testAccCheckAlicloudCenTransitRouterMulticastDomainPeerMemberSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
  default = "tf-testAccCenTransitRouterMulticastDomainPeerMember%d"
}
provider "alicloud" {
  alias  = "hz"
  region = "cn-hangzhou"
}
provider "alicloud" {
  alias  = "qd"
  region = "cn-qingdao"
}
data "alicloud_cen_instances" "default" {
  provider   = alicloud.hz
  name_regex = "no-deleting-cen"
}
data "alicloud_cen_transit_routers" "default" {
  provider   = alicloud.hz
  cen_id     = data.alicloud_cen_instances.default.instances.0.id
  name_regex = "no-deleting-cen"
}
data "alicloud_cen_transit_routers" "peer" {
  provider   = alicloud.qd
  cen_id     = data.alicloud_cen_instances.default.instances.0.id
  name_regex = "qingdao-no-delete-cen"
}
data "alicloud_cen_transit_router_multicast_domains" "default" {
  provider          = alicloud.hz
  transit_router_id = data.alicloud_cen_transit_routers.default.transit_routers.0.transit_router_id
  name_regex        = "no-deleting-cen"
}
data "alicloud_cen_transit_router_multicast_domains" "peer" {
  provider          = alicloud.qd
  transit_router_id = data.alicloud_cen_transit_routers.peer.transit_routers.0.transit_router_id
  name_regex        = "no-deleting-cen"
}
resource "alicloud_cen_transit_router_multicast_domain_peer_member" "default" {
  transit_router_multicast_domain_id      = data.alicloud_cen_transit_router_multicast_domains.default.ids.0
  peer_transit_router_multicast_domain_id = data.alicloud_cen_transit_router_multicast_domains.peer.ids.0
  group_ip_address                        = "239.1.1.1"
  provider                                = alicloud.hz
}

data "alicloud_cen_transit_router_multicast_domain_peer_members" "default" {
  provider                                = alicloud.hz
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
