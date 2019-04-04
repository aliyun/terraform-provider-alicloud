package alicloud

import (
	"testing"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudVpnCustomerGateway_importBasic(t *testing.T) {
	resourceName := "alicloud_vpn_customer_gateway.foo"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { 
			testAccPreCheck(t) 
			testAccPreCheckWithAccountSiteType(t, IntlSite)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVpnCustomerGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpnCustomerGatewayConfig,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
