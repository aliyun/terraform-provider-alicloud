package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudNetworkInterfaceAttachment_importBasic(t *testing.T) {
	resourceName := "alicloud_network_interface_attachment.att"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkInterfaceAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkInterfaceAttachmentConfig,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
