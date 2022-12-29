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

func TestAccAlicloudCenInterRegionTrafficQosQueuesDataSource(t *testing.T) {
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
		existConfig: testAccCheckAlicloudCenInterRegionTrafficQosQueuesDataSourceName(rand, map[string]string{
			"ids":                   `["${alicloud_cen_inter_region_traffic_qos_queue.default.id}"]`,
			"name_regex":            `"${alicloud_cen_inter_region_traffic_qos_queue.default.inter_region_traffic_qos_queue_name}"`,
			"traffic_qos_policy_id": `"${alicloud_cen_inter_region_traffic_qos_policy.default.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudCenInterRegionTrafficQosQueuesDataSourceName(rand, map[string]string{
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
			"queues.0.id": CHECKSET,
			"queues.0.inter_region_traffic_qos_queue_description": CHECKSET,
			"queues.0.inter_region_traffic_qos_queue_name":        CHECKSET,
			"queues.0.remain_bandwidth_percent":                   CHECKSET,
			"queues.0.traffic_qos_policy_id":                      CHECKSET,
		}
	}
	var fakeAlicloudCenInterRegionTrafficQosQueuesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":    "0",
			"names.#":  "0",
			"queues.#": "0",
		}
	}
	var alicloudCenInterRegionTrafficQosQueuesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_cen_inter_region_traffic_qos_queues.default",
		existMapFunc: existAlicloudCenInterRegionTrafficQosQueuesDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudCenInterRegionTrafficQosQueuesDataSourceNameMapFunc,
	}

	steps := allConf.buildDataSourceSteps(t, &alicloudCenInterRegionTrafficQosQueuesCheckInfo, rand)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckCenInterRegionTrafficQosQueueDestroyWithProviders(&providers),
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
`, rand, strings.Join(pairs, " \n "))
	return config
}
