package alicloud

import (
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

func TestAccAlicloudNetworkInterfaceAttachment_importBasic(t *testing.T) {
	resourceName := "alicloud_network_interface_attachment.att"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkInterfaceAttachmentDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccNetworkInterfaceAttachmentConfig,
			},

			resource.TestStep{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
