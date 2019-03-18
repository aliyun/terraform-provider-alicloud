package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudCommonBandwidthPackage_importBasic(t *testing.T) {
	resourceName := "alicloud_common_bandwidth_package.foo"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCommonBandwidthPackageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCommonBandwidthPackagePayByTraffic,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
