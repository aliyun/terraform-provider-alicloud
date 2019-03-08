package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudApigatewayGroup_importBasic(t *testing.T) {
	resourceName := "alicloud_api_gateway_group.apiGroupTest"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckWithRegions(t, false, connectivity.ApiGatewayNoSupportedRegions) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAlicloudApigatewayGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAlicloudApigatwayGroupBasic(acctest.RandIntRange(10000, 999999)),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
