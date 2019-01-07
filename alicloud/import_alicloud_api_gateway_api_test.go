package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudApigatewayApi_importBasic(t *testing.T) {
	resourceName := "alicloud_api_gateway_api.apiTest"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckWithRegions(t, false, connectivity.ApiGatewayNoSupportedRegions) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAlicloudApigatewayApiDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccAlicloudApigatwayApiBasic(acctest.RandIntRange(10000, 999999)),
			},

			resource.TestStep{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
