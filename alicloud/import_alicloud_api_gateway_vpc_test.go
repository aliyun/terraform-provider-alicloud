package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudApigatewayVpc_importBasic(t *testing.T) {
	resourceName := "alicloud_api_gateway_vpc_access.foo"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckWithRegions(t, false, connectivity.ApiGatewayNoSupportedRegions) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAlicloudApigatewayVpcDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAlicloudApigatwaVpcBasic,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
