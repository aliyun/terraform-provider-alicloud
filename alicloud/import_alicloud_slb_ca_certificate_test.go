package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudSlbCACertificate_import(t *testing.T) {
	resourceName := "alicloud_slb_ca_certificate.default"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSlbCACertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSlbCACertificateBasicConfig,
			},

			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"ca_certificate"},
			},
		},
	})
}
