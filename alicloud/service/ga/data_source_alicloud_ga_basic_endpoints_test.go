package ga_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudGaBasicEndpointsDataSource(t *testing.T) {
	rand := tfacctest.RandInt()
	acctest.CheckoutSupportedRegions(t, true, connectivity.GaSupportRegions)
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudGaBasicEndpointsDataSourceName(rand, map[string]string{
			"ids":           `["${alicloud_ga_basic_endpoint.default.id}"]`,
			"name_regex":    `"${alicloud_ga_basic_endpoint.default.basic_endpoint_name}"`,
			"endpoint_id":   `"${alicloud_ga_basic_endpoint.default.endpoint_id}"`,
			"endpoint_type": `"ENI"`,
			"name":          `"${alicloud_ga_basic_endpoint.default.basic_endpoint_name}"`,
			"status":        `"active"`,
		}),
		FakeConfig: testAccCheckAlicloudGaBasicEndpointsDataSourceName(rand, map[string]string{
			"ids":           `["${alicloud_ga_basic_endpoint.default.id}_fake"]`,
			"name_regex":    `"${alicloud_ga_basic_endpoint.default.basic_endpoint_name}_fake"`,
			"endpoint_id":   `"${alicloud_ga_basic_endpoint.default.endpoint_id}_fake"`,
			"endpoint_type": `"SLB"`,
			"name":          `"${alicloud_ga_basic_endpoint.default.basic_endpoint_name}_fake"`,
			"status":        `"deleting"`,
		}),
	}
	var existAlicloudGaBasicEndpointsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                 "1",
			"names.#":                               "1",
			"endpoints.#":                           "1",
			"endpoints.0.id":                        acctest.CHECKSET,
			"endpoints.0.endpoint_group_id":         acctest.CHECKSET,
			"endpoints.0.endpoint_id":               acctest.CHECKSET,
			"endpoints.0.accelerator_id":            acctest.CHECKSET,
			"endpoints.0.endpoint_type":             "ENI",
			"endpoints.0.endpoint_address":          acctest.CHECKSET,
			"endpoints.0.endpoint_sub_address_type": "secondary",
			"endpoints.0.endpoint_sub_address":      "192.168.0.1",
			"endpoints.0.basic_endpoint_name":       acctest.CHECKSET,
			"endpoints.0.status":                    "active",
		}
	}
	var fakeAlicloudGaBasicEndpointsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":       "0",
			"names.#":     "0",
			"endpoints.#": "0",
		}
	}
	var alicloudGaBasicEndpointsCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_ga_basic_endpoints.default",
		ExistMapFunc: existAlicloudGaBasicEndpointsDataSourceNameMapFunc,
		FakeMapFunc:  fakeAlicloudGaBasicEndpointsDataSourceNameMapFunc,
	}

	steps := gaBuildDataSourceSteps(allConf, &alicloudGaBasicEndpointsCheckInfo, rand)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithTime(t, []int{1})
		},
		ProviderFactories: acctest.SDKProviderFactoriesAlternate(),
		CheckDestroy:      testAccCheckGaBasicEndpointDestroy,
		Steps:             steps,
	})
}

func testAccCheckAlicloudGaBasicEndpointsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
	variable "name" {
  		default = "tf-testAccGaBasicEndpoint-%d"
	}

	`, rand) + acctest.ConfigAlternateRegionProvider("cn-shenzhen") + `

	data "alicloud_vpcs" "default" {
  		provider   = alicloudalt
  		name_regex = "default-NODELETING"
	}

	data "alicloud_vswitches" "default" {
  		provider = alicloudalt
  		vpc_id   = data.alicloud_vpcs.default.ids.0
	}

	resource "alicloud_security_group" "default" {
  		provider = alicloudalt
  		vpc_id   = data.alicloud_vpcs.default.ids.0
	}

	resource "alicloud_ecs_network_interface" "default" {
  		provider           = alicloudalt
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

	resource "alicloud_ga_basic_endpoint_group" "default" {
  		accelerator_id        = alicloud_ga_basic_accelerator.default.id
  		endpoint_group_region = "cn-shenzhen"
	}

	resource "alicloud_ga_basic_endpoint" "default" {
  		accelerator_id            = alicloud_ga_basic_accelerator.default.id
  		endpoint_group_id         = alicloud_ga_basic_endpoint_group.default.id
  		endpoint_type             = "ENI"
  		endpoint_address          = alicloud_ecs_network_interface.default.id
  		endpoint_sub_address_type = "secondary"
  		endpoint_sub_address      = "192.168.0.1"
  		basic_endpoint_name       = var.name
	}

	data "alicloud_ga_basic_endpoints" "default" {
  		endpoint_group_id = alicloud_ga_basic_endpoint.default.endpoint_group_id
		` + strings.Join(pairs, " \n ") + `
	}
`
	return config
}
