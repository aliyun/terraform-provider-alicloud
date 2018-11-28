package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudRouteTableAttachment_importBasic(t *testing.T) {
	resourceName := "alicloud_route_table_attachment.foo"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckWithRegions(t, false, connectivity.RouteTableNoSupportedRegions) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRouteTableAttachmentDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRouteTableAttachmentConfig,
			},

			resource.TestStep{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
