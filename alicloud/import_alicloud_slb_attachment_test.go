package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudSlbAttachment_import(t *testing.T) {
	resourceName := "alicloud_slb_attachment.default"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSlbDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSlbAttachment,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
