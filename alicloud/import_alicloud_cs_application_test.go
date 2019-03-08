package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudCSApplication_import(t *testing.T) {
	resourceName := "alicloud_cs_application.env"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckWithRegions(t, true, connectivity.SwarmSupportedRegions) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerApplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCSApplication_basic(testWebTemplate, testMultiTemplate),
			},

			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"latest_image", "blue_green", "blue_green_confirm"},
			},
		},
	})
}
