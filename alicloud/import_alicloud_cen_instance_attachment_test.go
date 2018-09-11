package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudCenInstanceAttachment_importBasic(t *testing.T) {
	resourceName := "alicloud_cen_instance_attachment.foo"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCenInstanceAttachmentDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCenInstanceAttachmentBasic,
			},

			resource.TestStep{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
