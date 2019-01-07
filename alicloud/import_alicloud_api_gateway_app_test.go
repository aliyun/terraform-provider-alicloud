package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

// At present, One account only support create 50 apps totally.
func SkipTestAccAlicloudApigatewayApp_importBasic(t *testing.T) {
	resourceName := "alicloud_api_gateway_app.appTest"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckWithRegions(t, false, connectivity.ApiGatewayNoSupportedRegions) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAlicloudApigatewayAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAlicloudApigatwayAppBasic,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
