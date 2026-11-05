package cen_test

import (
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// This testcase can not work in the multi region.
// The current resource does not need to support same region.
func SkipTestAccAlicloudCenBandwidthLimitImportBasic(t *testing.T) {
	resourceName := "alicloud_cen_bandwidth_limit.default"
	rand := tfacctest.RandIntRange(1000000, 9999999)

	// multi provideris
	var providers []*schema.Provider
	providerFactories := map[string]func() (*schema.Provider, error){
		"alicloud": func() (*schema.Provider, error) {
			p := acctest.Provider
			providers = append(providers, p)
			return p, nil
		},
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithAccountSiteType(t, acctest.DomesticSite)
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
