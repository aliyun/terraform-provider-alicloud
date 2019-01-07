package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudSlbServerCertificate_import(t *testing.T) {
	resourceName := "alicloud_slb_server_certificate.foo"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSlbServerCertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSlbServerCertificateBasicConfig,
			},

			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"server_certificate", "private_key"},
			},
		},
	})
}
