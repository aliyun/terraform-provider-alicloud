package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

// At present, only white list users can operate HaVip Resource.
func SkipTestAccAlicloudHaVipAttachment_importBasic(t *testing.T) {
	resourceName := "alicloud_havip_attachment.foo"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckHaVipAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccHaVipAttachmentConfig,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
