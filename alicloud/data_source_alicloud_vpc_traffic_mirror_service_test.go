package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudVPCTrafficMirrorServiceDataSource(t *testing.T) {
	resourceId := "data.alicloud_vpc_traffic_mirror_service.default"
	checkoutSupportedRegions(t, true, connectivity.VpcTrafficMirrorSupportRegions)
	testAccCheck := resourceAttrInit(resourceId, map[string]string{}).resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudVpcTrafficMirrorServiceDataSource,
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

const testAccCheckAlicloudVpcTrafficMirrorServiceDataSource = `
data "alicloud_vpc_traffic_mirror_service" "default" {
	enable = "On"
}
`
