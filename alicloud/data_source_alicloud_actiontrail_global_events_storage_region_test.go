package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudActiontrailGlobalEventsStorageRegionDataSource(t *testing.T) {
	checkoutSupportedRegions(t, true, connectivity.ActiontrailGlobalEventsStorageRegionSupportRegions)
	resourceId := "data.alicloud_actiontrail_global_events_storage_region.current"
	testAccCheck := resourceAttrInit(resourceId, map[string]string{}).resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudActiontrailGlobalEventsStorageRegionDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"storage_region": CHECKSET,
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
