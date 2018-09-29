package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudCenBandwidthLimit_importBasic(t *testing.T) {
	resourceName := "alicloud_cen_bandwidth_limit.foo"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCenBandwidthLimitDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCenBandwidthLimitConfig,
			},

			resource.TestStep{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
