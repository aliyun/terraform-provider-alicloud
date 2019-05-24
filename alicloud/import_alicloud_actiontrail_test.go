package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudActionTrail_importBasic(t *testing.T) {
	resourceName := "alicloud_actiontrail.default"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckWithRegions(t, false, connectivity.ApiGatewayNoSupportedRegions) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAlicloudApigatewayApiDestroy,
		Steps: []resource.TestStep{
			{
				Config: testActionTrailBasicConfig(acctest.RandIntRange(10000, 999999)),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
