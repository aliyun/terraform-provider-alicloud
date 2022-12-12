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

func TestAccAlicloudGaBasicAccelerateIpEndpointRelationsDataSource(t *testing.T) {
	rand := acctest.RandInt()
	checkoutSupportedRegions(t, true, connectivity.GaSupportRegions)
	var providers []*schema.Provider
	providerFactories := map[string]terraform.ResourceProviderFactory{
		"alicloud": func() (terraform.ResourceProvider, error) {
			p := Provider()
			providers = append(providers, p.(*schema.Provider))
			return p, nil
		},
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudGaBasicAccelerateIpEndpointRelationsDataSourceName(rand, map[string]string{
			"ids":              `["${alicloud_ga_basic_accelerate_ip_endpoint_relation.default.id}"]`,
			"accelerate_ip_id": `"${alicloud_ga_basic_accelerate_ip_endpoint_relation.default.accelerate_ip_id}"`,
			"endpoint_id":      `"${alicloud_ga_basic_accelerate_ip_endpoint_relation.default.endpoint_id}"`,
			"status":           `"active"`,
		}),
		fakeConfig: testAccCheckAlicloudGaBasicAccelerateIpEndpointRelationsDataSourceName(rand, map[string]string{
			"ids":              `["${alicloud_ga_basic_accelerate_ip_endpoint_relation.default.id}_fake"]`,
			"accelerate_ip_id": `"${alicloud_ga_basic_accelerate_ip_endpoint_relation.default.accelerate_ip_id}_fake"`,
			"endpoint_id":      `"${alicloud_ga_basic_accelerate_ip_endpoint_relation.default.endpoint_id}_fake"`,
		}),
	}
	var existAlicloudGaBasicAccelerateIpEndpointRelationsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                 "1",
			"relations.#":                           "1",
			"relations.0.id":                        CHECKSET,
			"relations.0.accelerator_id":            CHECKSET,
			"relations.0.accelerate_ip_id":          CHECKSET,
			"relations.0.endpoint_id":               CHECKSET,
			"relations.0.endpoint_type":             "ENI",
			"relations.0.endpoint_address":          CHECKSET,
			"relations.0.endpoint_sub_address_type": "primary",
			"relations.0.endpoint_sub_address":      "192.168.0.1",
			"relations.0.ip_address":                CHECKSET,
			"relations.0.basic_endpoint_name":       CHECKSET,
			"relations.0.status":                    "active",
		}
	}
	var fakeAlicloudGaBasicAccelerateIpEndpointRelationsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":       "0",
			"relations.#": "0",
		}
	}
	var alicloudGaBasicAccelerateIpEndpointRelationsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_ga_basic_accelerate_ip_endpoint_relations.default",
		existMapFunc: existAlicloudGaBasicAccelerateIpEndpointRelationsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudGaBasicAccelerateIpEndpointRelationsDataSourceNameMapFunc,
	}

	steps := allConf.buildDataSourceSteps(t, &alicloudGaBasicAccelerateIpEndpointRelationsCheckInfo, rand)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithTime(t, []int{1})
		},
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckGaBasicAccelerateIpEndpointRelationDestroyWithProviders(&providers),
		Steps:             steps,
	})
}

func testAccCheckAlicloudGaBasicAccelerateIpEndpointRelationsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
	variable "name" {
  		default = "tf-testAccGaBasicAccelerateIpEndpointRelation-%d"
	}

	provider "alicloud" {
  		alias  = "sz"
  		region = "cn-shenzhen"
	}

	provider "alicloud" {
  		alias  = "hz"
  		region = "cn-hangzhou"
	}

	data "alicloud_vpcs" "default" {
  		provider   = "alicloud.sz"
  		name_regex = "default-NODELETING"
	}

	data "alicloud_vswitches" "default" {
  		provider = "alicloud.sz"
  		vpc_id   = data.alicloud_vpcs.default.ids.0
	}

	resource "alicloud_security_group" "default" {
  		provider = "alicloud.sz"
  		vpc_id   = data.alicloud_vpcs.default.ids.0
	}

	resource "alicloud_ecs_network_interface" "default" {
  		provider           = "alicloud.sz"
  		vswitch_id         = data.alicloud_vswitches.default.ids.0
  		security_group_ids = [alicloud_security_group.default.id]
	}

	resource "alicloud_ga_basic_accelerator" "default" {
  		duration               = 1
  		pricing_cycle          = "Month"
  		basic_accelerator_name = var.name
  		description            = var.name
  		bandwidth_billing_type = "CDT"
  		auto_pay               = true
  		auto_use_coupon        = "true"
  		auto_renew             = false
  		auto_renew_duration    = 1
	}

	resource "alicloud_ga_basic_ip_set" "default" {
  		accelerator_id       = alicloud_ga_basic_accelerator.default.id
  		accelerate_region_id = "cn-hangzhou"
  		isp_type             = "BGP"
  		bandwidth            = "5"
	}

	resource "alicloud_ga_basic_accelerate_ip" "default" {
  		accelerator_id = alicloud_ga_basic_ip_set.default.accelerator_id
  		ip_set_id      = alicloud_ga_basic_ip_set.default.id
	}

	resource "alicloud_ga_basic_endpoint_group" "default" {
  		accelerator_id        = alicloud_ga_basic_accelerator.default.id
  		endpoint_group_region = "cn-shenzhen"
	}

	resource "alicloud_ga_basic_endpoint" "default" {
  		accelerator_id            = alicloud_ga_basic_accelerator.default.id
  		endpoint_group_id         = alicloud_ga_basic_endpoint_group.default.id
  		endpoint_type             = "ENI"
  		endpoint_address          = alicloud_ecs_network_interface.default.id
  		endpoint_sub_address_type = "primary"
  		endpoint_sub_address      = "192.168.0.1"
  		basic_endpoint_name       = var.name
	}

	resource "alicloud_ga_basic_accelerate_ip_endpoint_relation" "default" {
  		accelerator_id   = alicloud_ga_basic_accelerate_ip.default.accelerator_id
  		accelerate_ip_id = alicloud_ga_basic_accelerate_ip.default.id
  		endpoint_id      = alicloud_ga_basic_endpoint.default.endpoint_id
	}

	data "alicloud_ga_basic_accelerate_ip_endpoint_relations" "default" {
  		accelerator_id = alicloud_ga_basic_accelerate_ip_endpoint_relation.default.accelerator_id
		%s
	}
`, rand, strings.Join(pairs, " \n "))
	return config
}
