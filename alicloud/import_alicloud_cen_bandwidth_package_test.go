package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudCenBandwidthPackage_importBasic(t *testing.T) {
	resourceName := "alicloud_cen_bandwidth_package.foo"
	ignoreFields := []string{"period"}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithAccountSiteType(t, DomesticSite)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCenBandwidthPackageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCenBandwidthPackageConfig,
			},

			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: ignoreFields,
			},
		},
	})
}
