package alicloud

import (
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudBrainIndustrialServiceDataSource(t *testing.T) {
	checkoutSupportedRegions(t, true, connectivity.BrainIndustrialSupportRegions)
	resourceId := "data.alicloud_brain_industrial_service.current"
	testAccCheck := resourceAttrInit(resourceId, map[string]string{}).resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudBrainIndustrialServiceDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"id":     CHECKSET,
						"status": "Opened",
					}),
				),
			},
		},
	})
}

const testAccCheckAlicloudBrainIndustrialServiceDataSource = `
data "alicloud_brain_industrial_service" "current" {
	enable = "On"
}
`
