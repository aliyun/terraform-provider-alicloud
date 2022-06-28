package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudCenTransitRouterPeerAttachmentsDataSource(t *testing.T) {
	rand := acctest.RandInt()
	var providers []*schema.Provider
	providerFactories := map[string]terraform.ResourceProviderFactory{
		"alicloud": func() (terraform.ResourceProvider, error) {
			p := Provider()
			providers = append(providers, p.(*schema.Provider))
			return p, nil
		},
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenTransitRouterPeerAttachmentsDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_cen_transit_router_peer_attachment.default.transit_router_attachment_id}"]`,
			"name_regex": `"${alicloud_cen_transit_router_peer_attachment.default.transit_router_attachment_name}"`,
			"status":     `"Attached"`,
		}),
		fakeConfig: testAccCheckAlicloudCenTransitRouterPeerAttachmentsDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_cen_transit_router_peer_attachment.default.transit_router_attachment_id}_fake"]`,
			"name_regex": `"fake"`,
			"status":     `"Attaching"`,
		}),
	}
	var existAlicloudCenTransitRouterPeerAttachmentsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":         "1",
			"names.#":       "1",
			"attachments.#": "1",
			"attachments.0.auto_publish_route_enabled":            CHECKSET,
			"attachments.0.bandwidth":                             `5`,
			"attachments.0.cen_bandwidth_package_id":              CHECKSET,
			"attachments.0.peer_transit_router_id":                CHECKSET,
			"attachments.0.peer_transit_router_region_id":         `cn-beijing`,
			"attachments.0.transit_router_attachment_description": fmt.Sprintf("tf-testAccTransitRouterPeerAttachment-%d", rand),
			"attachments.0.transit_router_attachment_name":        fmt.Sprintf("tf-testAccTransitRouterPeerAttachment-%d", rand),
			"attachments.0.transit_router_id":                     CHECKSET,
			"attachments.0.status":                                "Attached",
			"attachments.0.resource_type":                         "TR",
		}
	}
	var fakeAlicloudCenTransitRouterPeerAttachmentsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudCenTransitRouterPeerAttachmentsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_cen_transit_router_peer_attachments.default",
		existMapFunc: existAlicloudCenTransitRouterPeerAttachmentsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudCenTransitRouterPeerAttachmentsDataSourceNameMapFunc,
	}

	steps := allConf.buildDataSourceSteps(t, &alicloudCenTransitRouterPeerAttachmentsCheckInfo, rand)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.CenTRSupportRegions)
		},
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckCenTransitRouterPeerAttachmentDestroyWithProviders(&providers),
		Steps:             steps,
	})
}
func testAccCheckAlicloudCenTransitRouterPeerAttachmentsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
variable "name" {	
	default = "tf-testAccTransitRouterPeerAttachment-%d"
}

provider "alicloud" {
  alias = "bj"
  region = "cn-beijing"
}

provider "alicloud" {
  alias = "cn"
  region = "cn-hangzhou"
}

resource "alicloud_cen_instance" "default" {
  provider = alicloud.cn
  name = var.name
  protection_level = "REDUCED"
}

resource "alicloud_cen_bandwidth_package" "default" {
  provider = alicloud.cn
  bandwidth                  = 5
  cen_bandwidth_package_name = var.name
  geographic_region_a_id     = "China"
  geographic_region_b_id     = "China"
}

resource "alicloud_cen_bandwidth_package_attachment" "default" {
  provider = alicloud.cn
  instance_id          = alicloud_cen_instance.default.id
  bandwidth_package_id = alicloud_cen_bandwidth_package.default.id
}

resource "alicloud_cen_transit_router" "default_0" {
  provider = alicloud.cn
  cen_id = alicloud_cen_bandwidth_package_attachment.default.instance_id
  transit_router_name = "${var.name}-00"
}

resource "alicloud_cen_transit_router" "default_1" {
  provider = alicloud.bj
  cen_id = alicloud_cen_transit_router.default_0.cen_id
  transit_router_name = "${var.name}-01"
}

resource "alicloud_cen_transit_router_peer_attachment" "default" {
  provider = alicloud.cn
  cen_id                                = alicloud_cen_instance.default.id
  transit_router_id                     = alicloud_cen_transit_router.default_0.transit_router_id
  peer_transit_router_region_id         = "cn-beijing"
  peer_transit_router_id                = alicloud_cen_transit_router.default_1.transit_router_id
  cen_bandwidth_package_id              = alicloud_cen_bandwidth_package_attachment.default.bandwidth_package_id
  bandwidth                             = 5
  transit_router_attachment_description = var.name
  transit_router_attachment_name        = var.name
}

data "alicloud_cen_transit_router_peer_attachments" "default" {
  provider = alicloud.cn
  cen_id = alicloud_cen_transit_router_peer_attachment.default.cen_id
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
