package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

// At present, One account only support create 50 apps totally.
func SkipTestAccAlicloudApigatewayApp_importBasic(t *testing.T) {
	resourceName := "alicloud_api_gateway_app.appTest"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAlicloudApigatewayAppDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccAlicloudApigatwayAppBasic,
			},

			resource.TestStep{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
