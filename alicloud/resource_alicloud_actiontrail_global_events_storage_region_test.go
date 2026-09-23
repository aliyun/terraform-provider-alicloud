package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudActiontrailGlobalEventsStorageRegion_basic0(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.ActiontrailGlobalEventsStorageRegionSupportRegions)
	resourceId := "alicloud_actiontrail_global_events_storage_region.default"
	ra := resourceAttrInit(resourceId, AlicloudActiontrailGlobalEventsStorageRegionMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ActiontrailService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeActiontrailGlobalEventsStorageRegion")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testaccactiontrail%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudActiontrailGlobalEventsStorageRegionBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"storage_region": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"storage_region": defaultRegionToTest,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"storage_region": defaultRegionToTest,
					}),
				),
			},
		},
	})
}

var AlicloudActiontrailGlobalEventsStorageRegionMap0 = map[string]string{}

func AlicloudActiontrailGlobalEventsStorageRegionBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}

`, name)
}
