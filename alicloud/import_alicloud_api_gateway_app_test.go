package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudApigatewayApp_importBasic(t *testing.T) {
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
