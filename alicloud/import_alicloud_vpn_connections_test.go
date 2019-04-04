package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudVpnConnection_importBasic(t *testing.T) {
	resourceName := "alicloud_vpn_connection.foo"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithAccountSiteType(t, IntlSite)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVpnConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpnConnConfig,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
