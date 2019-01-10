package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudRouterInterface_import(t *testing.T) {
	resourceName := "alicloud_router_interface.interface"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVpcDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRouterInterfaceConfig,
			},

			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"period"},
			},
		},
	})
}
