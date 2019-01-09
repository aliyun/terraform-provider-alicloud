package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudPvtzZone_importBasic(t *testing.T) {
	resourceName := "alicloud_pvtz_zone.foo"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccAlicloudPvtzZoneDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPvtzZoneConfig(acctest.RandIntRange(10000, 999999)),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
