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

func TestAccAlicloudCenInterRegionTrafficQosPoliciesDataSource(t *testing.T) {
	rand := acctest.RandInt()
	checkoutSupportedRegions(t, true, connectivity.TestSalveRegions)
	var providers []*schema.Provider
	providerFactories := map[string]terraform.ResourceProviderFactory{
		"alicloud": func() (terraform.ResourceProvider, error) {
			p := Provider()
			providers = append(providers, p.(*schema.Provider))
			return p, nil
		},
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenInterRegionTrafficQosPoliciesDataSourceName(rand, map[string]string{
			"ids":                            `["${alicloud_cen_inter_region_traffic_qos_policy.default.id}"]`,
			"name_regex":                     `"${alicloud_cen_inter_region_traffic_qos_policy.default.inter_region_traffic_qos_policy_name}"`,
			"traffic_qos_policy_id":          `"${alicloud_cen_inter_region_traffic_qos_policy.default.id}"`,
			"traffic_qos_policy_name":        `"${alicloud_cen_inter_region_traffic_qos_policy.default.inter_region_traffic_qos_policy_name}"`,
			"traffic_qos_policy_description": `"${alicloud_cen_inter_region_traffic_qos_policy.default.inter_region_traffic_qos_policy_description}"`,
			"status":                         `"Active"`,
		}),
		fakeConfig: testAccCheckAlicloudCenInterRegionTrafficQosPoliciesDataSourceName(rand, map[string]string{
			"ids":                            `["${alicloud_cen_inter_region_traffic_qos_policy.default.id}_fake"]`,
			"name_regex":                     `"${alicloud_cen_inter_region_traffic_qos_policy.default.inter_region_traffic_qos_policy_name}_fake"`,
			"traffic_qos_policy_id":          `"${alicloud_cen_inter_region_traffic_qos_policy.default.id}_fake"`,
			"traffic_qos_policy_name":        `"${alicloud_cen_inter_region_traffic_qos_policy.default.inter_region_traffic_qos_policy_name}_fake"`,
			"traffic_qos_policy_description": `"${alicloud_cen_inter_region_traffic_qos_policy.default.inter_region_traffic_qos_policy_description}_fake"`,
			"status":                         `"Deleted"`,
		}),
	}
	var existAlicloudCenInterRegionTrafficQosPoliciesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                        "1",
			"names.#":                      "1",
			"policies.#":                   "1",
			"policies.0.id":                CHECKSET,
			"policies.0.transit_router_id": CHECKSET,
			"policies.0.transit_router_attachment_id":                CHECKSET,
			"policies.0.inter_region_traffic_qos_policy_id":          CHECKSET,
			"policies.0.inter_region_traffic_qos_policy_name":        CHECKSET,
			"policies.0.inter_region_traffic_qos_policy_description": CHECKSET,
			"policies.0.status": "Active",
		}
	}
	var fakeAlicloudCenInterRegionTrafficQosPoliciesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":      "0",
			"names.#":    "0",
			"policies.#": "0",
		}
	}
	var alicloudCenInterRegionTrafficQosPoliciesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_cen_inter_region_traffic_qos_policies.default",
		existMapFunc: existAlicloudCenInterRegionTrafficQosPoliciesDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudCenInterRegionTrafficQosPoliciesDataSourceNameMapFunc,
	}

	steps := allConf.buildDataSourceSteps(t, &alicloudCenInterRegionTrafficQosPoliciesCheckInfo, rand)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckCenInterRegionTrafficQosPolicyDestroyWithProviders(&providers),
		Steps:             steps,
	})
}

func testAccCheckAlicloudCenInterRegionTrafficQosPoliciesDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
	variable "name" {
  		default = "tf-testAccCenInterRegionTrafficQosPolicy-%d"
	}

	provider "alicloud" {
  		alias  = "bj"
  		region = "cn-beijing"
	}

	provider "alicloud" {
  		alias  = "hz"
  		region = "cn-hangzhou"
	}

	resource "alicloud_cen_instance" "default" {
  		provider          = alicloud.hz
  		cen_instance_name = var.name
	}

	resource "alicloud_cen_bandwidth_package" "default" {
  		provider               = alicloud.hz
  		bandwidth              = 5
  		geographic_region_a_id = "China"
  		geographic_region_b_id = "China"
	}

	resource "alicloud_cen_bandwidth_package_attachment" "default" {
  		provider             = alicloud.hz
  		instance_id          = alicloud_cen_instance.default.id
  		bandwidth_package_id = alicloud_cen_bandwidth_package.default.id
	}

	resource "alicloud_cen_transit_router" "hz" {
  		provider = alicloud.hz
  		cen_id   = alicloud_cen_bandwidth_package_attachment.default.instance_id
	}

	resource "alicloud_cen_transit_router" "bj" {
  		provider = alicloud.bj
  		cen_id   = alicloud_cen_transit_router.hz.cen_id
	}

	resource "alicloud_cen_transit_router_peer_attachment" "default" {
  		provider                      = alicloud.hz
  		cen_id                        = alicloud_cen_instance.default.id
  		transit_router_id             = alicloud_cen_transit_router.hz.transit_router_id
  		peer_transit_router_region_id = "cn-beijing"
  		peer_transit_router_id        = alicloud_cen_transit_router.bj.transit_router_id
  		cen_bandwidth_package_id      = alicloud_cen_bandwidth_package_attachment.default.bandwidth_package_id
  		bandwidth                     = 5
	}

	resource "alicloud_cen_inter_region_traffic_qos_policy" "default" {
  		provider                                    = alicloud.hz
  		transit_router_id                           = alicloud_cen_transit_router.hz.transit_router_id
  		transit_router_attachment_id                = alicloud_cen_transit_router_peer_attachment.default.transit_router_attachment_id
  		inter_region_traffic_qos_policy_name        = var.name
  		inter_region_traffic_qos_policy_description = var.name
	}

	data "alicloud_cen_inter_region_traffic_qos_policies" "default" {
  		transit_router_id            = alicloud_cen_inter_region_traffic_qos_policy.default.transit_router_id
  		transit_router_attachment_id = alicloud_cen_inter_region_traffic_qos_policy.default.transit_router_attachment_id
		%s
	}
`, rand, strings.Join(pairs, " \n "))
	return config
}
