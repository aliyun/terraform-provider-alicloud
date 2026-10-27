// Package actiontrail_test.
package actiontrail_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/actiontrail"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudActiontrailGlobalEventsStorageRegion_basic0(t *testing.T) {
	var v map[string]interface{}
	acctest.CheckoutSupportedRegions(t, true, connectivity.ActiontrailGlobalEventsStorageRegionSupportRegions)
	resourceId := "alicloud_actiontrail_global_events_storage_region.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudActiontrailGlobalEventsStorageRegionMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return actiontrail.NewActiontrailService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeActiontrailGlobalEventsStorageRegion")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testaccactiontrail%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudActiontrailGlobalEventsStorageRegionBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"storage_region": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"storage_region": acctest.DefaultRegionToTest,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"storage_region": acctest.DefaultRegionToTest,
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
