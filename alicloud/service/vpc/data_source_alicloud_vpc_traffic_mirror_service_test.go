package vpc_test

import (
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudVPCTrafficMirrorServiceDataSource(t *testing.T) {
	resourceId := "data.alicloud_vpc_traffic_mirror_service.default"
	checkoutSupportedRegions(t, true, connectivity.VpcTrafficMirrorSupportRegions)
	testAccCheck := acctest.ResourceAttrInit(resourceId, map[string]string{}).ResourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudVpcTrafficMirrorServiceDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"id":     acctest.CHECKSET,
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
