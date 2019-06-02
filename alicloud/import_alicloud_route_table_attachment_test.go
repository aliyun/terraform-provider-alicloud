package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudRouteTableAttachment_importBasic(t *testing.T) {
	resourceName := "alicloud_route_table_attachment.default"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckWithRegions(t, false, connectivity.RouteTableNoSupportedRegions) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRouteTableAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRouteTableAttachmentConfigBasic(acctest.RandInt()),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
