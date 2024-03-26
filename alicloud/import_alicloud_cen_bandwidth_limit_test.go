package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// This testcase can not work in the multi region.
// The current resource does not need to support same region.
func SkipTestAccAlicloudCenBandwidthLimitImportBasic(t *testing.T) {
	resourceName := "alicloud_cen_bandwidth_limit.default"
	rand := acctest.RandIntRange(1000000, 9999999)

	// multi provideris
	var providers []*schema.Provider
	providerFactories := map[string]func() (*schema.Provider, error){
		"alicloud": func() (*schema.Provider, error) {
			p := Provider()
			providers = append(providers, p.(*schema.Provider))
			return p.(*schema.Provider), nil
		},
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithAccountSiteType(t, DomesticSite)
		},

		// module name
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckCenBandwidthLimitDestroyWithProviders(&providers),
		Steps: []resource.TestStep{
			{
				Config: testAccCenBandwidthLimitCreateConfig(rand),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
