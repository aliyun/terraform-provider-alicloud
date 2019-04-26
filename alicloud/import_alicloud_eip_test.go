package alicloud

import (
	"github.com/hashicorp/terraform/helper/acctest"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudEip_importBasic(t *testing.T) {
	resourceName := "alicloud_eip.default"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckEipConfigBasic(acctest.RandInt(), "PayByTraffic"),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
