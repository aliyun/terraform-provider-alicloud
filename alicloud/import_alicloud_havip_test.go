package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

// At present, only white list users can operate HaVip Resource.
func SkipTestAccAlicloudHaVip_importBasic(t *testing.T) {
	resourceName := "alicloud_havip.foo"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckHaVipDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccHaVipConfig,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
