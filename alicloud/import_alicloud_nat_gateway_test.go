package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudNatGateway_importSpec(t *testing.T) {
	resourceName := "alicloud_nat_gateway.foo"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNatGatewayDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccNatGatewayConfigSpec,
			},

			resource.TestStep{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
