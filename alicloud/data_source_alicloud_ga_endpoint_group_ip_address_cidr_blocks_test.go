package alicloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"strings"

	"testing"
)

// TODO: There is an api bug that it always return some cidr blocks whatever accelerator_id whether correct
// aone id: 62065946.
func SkipTestAccAliCloudGaEndpointGroupIpAddressCidrBlocksDataSource(t *testing.T) {
	rand := acctest.RandInt()

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudGaEndpointGroupIpAddressCidrBlocksDataSourceName(rand, map[string]string{
			"accelerator_id": `"${data.alicloud_ga_accelerators.exist.accelerators.0.id}"`,
		}),
		fakeConfig: testAccCheckAliCloudGaEndpointGroupIpAddressCidrBlocksDataSourceName(rand, map[string]string{
			"accelerator_id": `"${data.alicloud_ga_accelerators.fake.accelerators.0.id}"`,
		}),
	}

	var existAliCloudGaEndpointGroupIpAddressCidrBlocksDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"endpoint_group_ip_address_cidr_blocks.#":                          "1",
			"endpoint_group_ip_address_cidr_blocks.0.endpoint_group_region":    CHECKSET,
			"endpoint_group_ip_address_cidr_blocks.0.ip_address_cidr_blocks.#": CHECKSET,
			"endpoint_group_ip_address_cidr_blocks.0.status":                   CHECKSET,
		}
	}

	var fakeAliCloudGaEndpointGroupIpAddressCidrBlocksDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"endpoint_group_ip_address_cidr_blocks.#":                          "1",
			"endpoint_group_ip_address_cidr_blocks.0.endpoint_group_region":    CHECKSET,
			"endpoint_group_ip_address_cidr_blocks.0.ip_address_cidr_blocks.#": "0",
			"endpoint_group_ip_address_cidr_blocks.0.status":                   CHECKSET,
		}
	}

	var aliCloudGaEndpointGroupIpAddressCidrBlocksCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_ga_endpoint_group_ip_address_cidr_blocks.default",
		existMapFunc: existAliCloudGaEndpointGroupIpAddressCidrBlocksDataSourceNameMapFunc,
		fakeMapFunc:  fakeAliCloudGaEndpointGroupIpAddressCidrBlocksDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}

	aliCloudGaEndpointGroupIpAddressCidrBlocksCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, allConf)
}

func testAccCheckAliCloudGaEndpointGroupIpAddressCidrBlocksDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
	data "alicloud_ga_accelerators" "exist" {
  		status                 = "active"
  		bandwidth_billing_type = "BandwidthPackage"
	}

	data "alicloud_ga_accelerators" "fake" {
  		status                 = "active"
  		bandwidth_billing_type = "CDT"
	}

	data "alicloud_ga_endpoint_group_ip_address_cidr_blocks" "default" {
  		endpoint_group_region = "cn-hangzhou"
		%s
	}
`, strings.Join(pairs, " \n "))
	return config
}
