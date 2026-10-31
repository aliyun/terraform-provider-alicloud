package cen_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAlicloudCenInterRegionTrafficQosPoliciesDataSource(t *testing.T) {
	rand := tfacctest.RandInt()
	acctest.CheckoutSupportedRegions(t, true, connectivity.TestSalveRegions)
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudCenInterRegionTrafficQosPoliciesDataSourceName(rand, map[string]string{
			"ids":                            `["${alicloud_cen_inter_region_traffic_qos_policy.default.id}"]`,
			"name_regex":                     `"${alicloud_cen_inter_region_traffic_qos_policy.default.inter_region_traffic_qos_policy_name}"`,
			"traffic_qos_policy_id":          `"${alicloud_cen_inter_region_traffic_qos_policy.default.id}"`,
			"traffic_qos_policy_name":        `"${alicloud_cen_inter_region_traffic_qos_policy.default.inter_region_traffic_qos_policy_name}"`,
			"traffic_qos_policy_description": `"${alicloud_cen_inter_region_traffic_qos_policy.default.inter_region_traffic_qos_policy_description}"`,
			"status":                         `"Active"`,
		}),
		FakeConfig: testAccCheckAlicloudCenInterRegionTrafficQosPoliciesDataSourceName(rand, map[string]string{
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
			"policies.0.id":                acctest.CHECKSET,
			"policies.0.transit_router_id": acctest.CHECKSET,
			"policies.0.transit_router_attachment_id":                acctest.CHECKSET,
			"policies.0.inter_region_traffic_qos_policy_id":          acctest.CHECKSET,
			"policies.0.inter_region_traffic_qos_policy_name":        acctest.CHECKSET,
			"policies.0.inter_region_traffic_qos_policy_description": acctest.CHECKSET,
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
	var alicloudCenInterRegionTrafficQosPoliciesCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_cen_inter_region_traffic_qos_policies.default",
		ExistMapFunc: existAlicloudCenInterRegionTrafficQosPoliciesDataSourceNameMapFunc,
		FakeMapFunc:  fakeAlicloudCenInterRegionTrafficQosPoliciesDataSourceNameMapFunc,
	}

	existUpdate := acctest.ResourceAttrInit(alicloudCenInterRegionTrafficQosPoliciesCheckInfo.ResourceId, alicloudCenInterRegionTrafficQosPoliciesCheckInfo.ExistMapFunc(rand)).ResourceAttrMapUpdateSet()
	emptyUpdate := acctest.ResourceAttrInit(alicloudCenInterRegionTrafficQosPoliciesCheckInfo.ResourceId, alicloudCenInterRegionTrafficQosPoliciesCheckInfo.FakeMapFunc(rand)).ResourceAttrMapUpdateSet()
	var steps []resource.TestStep
	if allConf.ExistConfig != "" {
		steps = append(steps, resource.TestStep{Config: allConf.ExistConfig, Check: resource.ComposeTestCheckFunc(existUpdate(allConf.ExistChangMap))})
	}
	if allConf.FakeConfig != "" {
		steps = append(steps, resource.TestStep{Config: allConf.FakeConfig, Check: resource.ComposeTestCheckFunc(emptyUpdate(allConf.FakeChangMap))})
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactoriesAlternate(),
		CheckDestroy:      nil,
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

	%s

	resource "alicloud_cen_instance" "default" {
  		cen_instance_name = var.name
	}

	resource "alicloud_cen_bandwidth_package" "default" {
  		bandwidth              = 5
  		geographic_region_a_id = "China"
  		geographic_region_b_id = "China"
	}

	resource "alicloud_cen_bandwidth_package_attachment" "default" {
  		instance_id          = alicloud_cen_instance.default.id
  		bandwidth_package_id = alicloud_cen_bandwidth_package.default.id
	}

	resource "alicloud_cen_transit_router" "hz" {
  		cen_id   = alicloud_cen_bandwidth_package_attachment.default.instance_id
	}

	resource "alicloud_cen_transit_router" "bj" {
  		provider = alicloudalt
  		cen_id   = alicloud_cen_transit_router.hz.cen_id
	}

	resource "alicloud_cen_transit_router_peer_attachment" "default" {
  		cen_id                        = alicloud_cen_instance.default.id
  		transit_router_id             = alicloud_cen_transit_router.hz.transit_router_id
  		peer_transit_router_region_id = "cn-beijing"
  		peer_transit_router_id        = alicloud_cen_transit_router.bj.transit_router_id
  		cen_bandwidth_package_id      = alicloud_cen_bandwidth_package_attachment.default.bandwidth_package_id
  		bandwidth                     = 5
	}

	resource "alicloud_cen_inter_region_traffic_qos_policy" "default" {
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
`, rand, acctest.ConfigAlternateRegionProvider("cn-beijing"), strings.Join(pairs, " \n "))
	return config
}
