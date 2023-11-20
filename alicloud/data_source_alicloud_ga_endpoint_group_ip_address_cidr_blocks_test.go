package alicloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"testing"
)

func TestAccAliCloudGaEndpointGroupIpAddressCidrBlocksDataSource(t *testing.T) {
	resourceId := "data.alicloud_ga_endpoint_group_ip_address_cidr_blocks.default"
	testAccCheck := resourceAttrInit(resourceId, map[string]string{}).resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAliCloudGaEndpointGroupIpAddressCidrBlocksDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"endpoint_group_ip_address_cidr_blocks.#":                          "1",
						"endpoint_group_ip_address_cidr_blocks.0.endpoint_group_region":    CHECKSET,
						"endpoint_group_ip_address_cidr_blocks.0.ip_address_cidr_blocks.#": CHECKSET,
						"endpoint_group_ip_address_cidr_blocks.0.status":                   CHECKSET,
					}),
				),
			},
		},
	})
}

const testAccCheckAliCloudGaEndpointGroupIpAddressCidrBlocksDataSource = `
	data "alicloud_ga_endpoint_group_ip_address_cidr_blocks" "default" {
  		endpoint_group_region = "cn-hangzhou"
	}
`
