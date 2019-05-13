package alicloud

import (
	"github.com/hashicorp/terraform/helper/acctest"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudCommonBandwidthPackage_importBasic(t *testing.T) {
	resourceName := "alicloud_common_bandwidth_package.default"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCommonBandwidthPackageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCommonBandwidthPackageBasic(acctest.RandInt(), "PayByTraffic"),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
