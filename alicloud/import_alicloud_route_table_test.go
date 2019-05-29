package alicloud

import (
	"github.com/hashicorp/terraform/helper/acctest"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudRouteTable_importBasic(t *testing.T) {
	resourceName := "alicloud_route_table.default"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckWithRegions(t, false, connectivity.RouteTableNoSupportedRegions) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRouteTableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRouteTableConfigBasic(acctest.RandInt()),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
