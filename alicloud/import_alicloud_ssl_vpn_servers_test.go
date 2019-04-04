package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudSslVpnServer_importBasic(t *testing.T) {
	resourceName := "alicloud_ssl_vpn_server.foo"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithAccountSiteType(t, IntlSite)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSslVpnServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSslVpnServerConfig,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
