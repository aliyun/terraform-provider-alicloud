// Package actiontrail_test.
package actiontrail_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudActiontrailGlobalEventsStorageRegionDataSource(t *testing.T) {
	acctest.CheckoutSupportedRegions(t, true, connectivity.ActiontrailGlobalEventsStorageRegionSupportRegions)
	resourceId := "data.alicloud_actiontrail_global_events_storage_region.current"
	testAccCheck := acctest.ResourceAttrInit(resourceId, map[string]string{}).ResourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudActiontrailGlobalEventsStorageRegionDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"storage_region": acctest.CHECKSET,
					}),
				),
			},
		},
	})
}

var testAccCheckAlicloudActiontrailGlobalEventsStorageRegionDataSource = fmt.Sprintf(`
data "alicloud_actiontrail_global_events_storage_region" "current" {
}
`)
