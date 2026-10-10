package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// TestAccAlicloudRemoteRunValidation_regions exercises the remote runner using
// a read-only data source. This draft PR is a validation fixture, not for merge.
func TestAccAlicloudRemoteRunValidation_regions(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudRegionsDataSourceRegionsConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_regions.region"),
					resource.TestCheckResourceAttr("data.alicloud_regions.region", "regions.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_regions.region", "regions.0.id", "cn-beijing"),
				),
			},
		},
	})
}
