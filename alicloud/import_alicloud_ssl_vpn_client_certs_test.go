package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudSslVpnClientCert_importBasic(t *testing.T) {
	resourceName := "alicloud_ssl_vpn_client_cert.foo"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithAccountSiteType(t, IntlSite)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSslVpnClientCertDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSslVpnClientCertConfig,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
