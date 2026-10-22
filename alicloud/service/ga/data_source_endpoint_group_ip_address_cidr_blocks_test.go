package ga_test

import (
	"os"
	"regexp"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// GA ListEndpointGroupIpAddressCidrBlocks is Hangzhou-only for the endpoint
// group CIDR lookup; the flat suite pinned the region through
// testAccPreCheckWithRegions, so pin it here.
const testAccEndpointGroupIpAddressCidrBlocksRegion = "cn-hangzhou"

// The non-existent accelerator id reported in aone bug 62065946. The 2024-12-23
// defect (a fake AcceleratorId returned real data instead of erroring) no
// longer reproduces: the API now returns 400 NotExist.Accelerator.
const testAccFakeAcceleratorId = "ga-bp1fs9fu0f1a7dnl6hr3g-fake"

func TestAccAliCloudGaEndpointGroupIpAddressCidrBlocksDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			os.Setenv("ALICLOUD_REGION", testAccEndpointGroupIpAddressCidrBlocksRegion)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			// AC1: a literally non-existent AcceleratorId must make the data
			// source read fail with NotExist.Accelerator and produce no CIDR
			// data, instead of returning real endpoint group CIDR blocks.
			{
				Config:      testAccEndpointGroupIpAddressCidrBlocksConfig(testAccFakeAcceleratorId),
				ExpectError: regexp.MustCompile(`NotExist\.Accelerator`),
			},
			// AC2: omitting the optional accelerator_id returns the cn-hangzhou
			// endpoint group CIDR blocks: State=active, non-empty CIDR list.
			{
				Config: testAccEndpointGroupIpAddressCidrBlocksConfig(""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.alicloud_ga_endpoint_group_ip_address_cidr_blocks.default", "endpoint_group_ip_address_cidr_blocks.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_ga_endpoint_group_ip_address_cidr_blocks.default", "endpoint_group_ip_address_cidr_blocks.0.endpoint_group_region", testAccEndpointGroupIpAddressCidrBlocksRegion),
					resource.TestCheckResourceAttr("data.alicloud_ga_endpoint_group_ip_address_cidr_blocks.default", "endpoint_group_ip_address_cidr_blocks.0.status", "active"),
					resource.TestMatchResourceAttr("data.alicloud_ga_endpoint_group_ip_address_cidr_blocks.default", "endpoint_group_ip_address_cidr_blocks.0.ip_address_cidr_blocks.#", regexp.MustCompile(`[1-9]\d*`)),
				),
			},
		},
	})
}

func testAccEndpointGroupIpAddressCidrBlocksConfig(acceleratorId string) string {
	if acceleratorId == "" {
		return `
data "alicloud_ga_endpoint_group_ip_address_cidr_blocks" "default" {
  endpoint_group_region = "cn-hangzhou"
}
`
	}
	return `
data "alicloud_ga_endpoint_group_ip_address_cidr_blocks" "default" {
  endpoint_group_region = "cn-hangzhou"
  accelerator_id         = "` + acceleratorId + `"
}
`
}
