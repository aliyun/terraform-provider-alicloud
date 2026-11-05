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

func TestAccAlicloudCenInterRegionTrafficQosQueuesDataSource(t *testing.T) {
	rand := tfacctest.RandInt()
	acctest.CheckoutSupportedRegions(t, true, connectivity.TestSalveRegions)
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudCenInterRegionTrafficQosQueuesDataSourceName(rand, map[string]string{
			"ids":                   `["${alicloud_cen_inter_region_traffic_qos_queue.default.id}"]`,
			"name_regex":            `"${alicloud_cen_inter_region_traffic_qos_queue.default.inter_region_traffic_qos_queue_name}"`,
			"traffic_qos_policy_id": `"${alicloud_cen_inter_region_traffic_qos_policy.default.id}"`,
		}),
		FakeConfig: testAccCheckAlicloudCenInterRegionTrafficQosQueuesDataSourceName(rand, map[string]string{
			"ids":                   `["${alicloud_cen_inter_region_traffic_qos_queue.default.id}_fake"]`,
			"name_regex":            `"${alicloud_cen_inter_region_traffic_qos_queue.default.inter_region_traffic_qos_queue_name}_fake"`,
			"traffic_qos_policy_id": `"${alicloud_cen_inter_region_traffic_qos_policy.default.id}_fake"`,
		}),
	}
	var existAlicloudCenInterRegionTrafficQosQueuesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":       "1",
			"names.#":     "1",
			"queues.#":    "1",
			"queues.0.id": acctest.CHECKSET,
			"queues.0.inter_region_traffic_qos_queue_description": acctest.CHECKSET,
			"queues.0.inter_region_traffic_qos_queue_name":        acctest.CHECKSET,
			"queues.0.remain_bandwidth_percent":                   acctest.CHECKSET,
			"queues.0.traffic_qos_policy_id":                      acctest.CHECKSET,
		}
	}
	var fakeAlicloudCenInterRegionTrafficQosQueuesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":    "0",
			"names.#":  "0",
			"queues.#": "0",
		}
	}
	var alicloudCenInterRegionTrafficQosQueuesCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_cen_inter_region_traffic_qos_queues.default",
		ExistMapFunc: existAlicloudCenInterRegionTrafficQosQueuesDataSourceNameMapFunc,
		FakeMapFunc:  fakeAlicloudCenInterRegionTrafficQosQueuesDataSourceNameMapFunc,
	}

	existUpdate := acctest.ResourceAttrInit(alicloudCenInterRegionTrafficQosQueuesCheckInfo.ResourceId, alicloudCenInterRegionTrafficQosQueuesCheckInfo.ExistMapFunc(rand)).ResourceAttrMapUpdateSet()
	emptyUpdate := acctest.ResourceAttrInit(alicloudCenInterRegionTrafficQosQueuesCheckInfo.ResourceId, alicloudCenInterRegionTrafficQosQueuesCheckInfo.FakeMapFunc(rand)).ResourceAttrMapUpdateSet()
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

func testAccCheckAlicloudCenInterRegionTrafficQosQueuesDataSourceName(rand int, attrMap map[string]string) string {
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

	resource "alicloud_cen_inter_region_traffic_qos_queue" "default" {
	  remain_bandwidth_percent = 20
	  traffic_qos_policy_id    = alicloud_cen_inter_region_traffic_qos_policy.default.id
	  dscps = [1, 2]
      inter_region_traffic_qos_queue_name = "test"
	  inter_region_traffic_qos_queue_description = "test"
	}

	data "alicloud_cen_inter_region_traffic_qos_queues" "default" {
		%s
	}
`, rand, acctest.ConfigAlternateRegionProvider("cn-beijing"), strings.Join(pairs, " \n "))
	return config
}
