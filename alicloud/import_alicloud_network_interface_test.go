package alicloud

import (
	"github.com/hashicorp/terraform/helper/acctest"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudNetworkInterface_importBasic(t *testing.T) {
	resourceName := "alicloud_network_interface.default"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkInterfaceConfigBasic(acctest.RandIntRange(1000, 9999)),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
