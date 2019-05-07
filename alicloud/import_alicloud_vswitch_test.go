package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudVSwitch_importBasic(t *testing.T) {
	resourceName := "alicloud_vswitch.default"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVSwitchDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVSwitchConfigBasic(acctest.RandInt()),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
