package alicloud

import (
	"github.com/hashicorp/terraform/helper/acctest"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudCommonBandwidthPackageAttachment_importBasic(t *testing.T) {
	resourceName := "alicloud_common_bandwidth_package_attachment.default"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCommonBandwidthPackageAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCommonBandwidthPackageAttachmentConfigBasic(acctest.RandInt()),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
