package ga_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

// TODO: There is an api bug that it always return some cidr blocks whatever accelerator_id whether correct
// aone id: 62065946.
func SkipTestAccAliCloudGaEndpointGroupIpAddressCidrBlocksDataSource(t *testing.T) {
	rand := tfacctest.RandInt()

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudGaEndpointGroupIpAddressCidrBlocksDataSourceName(rand, map[string]string{
			"accelerator_id": `"${data.alicloud_ga_accelerators.exist.accelerators.0.id}"`,
		}),
		FakeConfig: testAccCheckAliCloudGaEndpointGroupIpAddressCidrBlocksDataSourceName(rand, map[string]string{
			"accelerator_id": `"${data.alicloud_ga_accelerators.fake.accelerators.0.id}"`,
		}),
	}

	var existAliCloudGaEndpointGroupIpAddressCidrBlocksDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"endpoint_group_ip_address_cidr_blocks.#":                          "1",
			"endpoint_group_ip_address_cidr_blocks.0.endpoint_group_region":    acctest.CHECKSET,
			"endpoint_group_ip_address_cidr_blocks.0.ip_address_cidr_blocks.#": acctest.CHECKSET,
			"endpoint_group_ip_address_cidr_blocks.0.status":                   acctest.CHECKSET,
		}
	}

	var fakeAliCloudGaEndpointGroupIpAddressCidrBlocksDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"endpoint_group_ip_address_cidr_blocks.#":                          "1",
			"endpoint_group_ip_address_cidr_blocks.0.endpoint_group_region":    acctest.CHECKSET,
			"endpoint_group_ip_address_cidr_blocks.0.ip_address_cidr_blocks.#": "0",
			"endpoint_group_ip_address_cidr_blocks.0.status":                   acctest.CHECKSET,
		}
	}

	var aliCloudGaEndpointGroupIpAddressCidrBlocksCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_ga_endpoint_group_ip_address_cidr_blocks.default",
		ExistMapFunc: existAliCloudGaEndpointGroupIpAddressCidrBlocksDataSourceNameMapFunc,
		FakeMapFunc:  fakeAliCloudGaEndpointGroupIpAddressCidrBlocksDataSourceNameMapFunc,
	}

	preCheck := func() {
		acctest.PreCheck(t)
	}

	aliCloudGaEndpointGroupIpAddressCidrBlocksCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, allConf)
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
