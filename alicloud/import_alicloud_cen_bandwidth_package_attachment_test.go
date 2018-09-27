package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudCenBandwidthPackageAttachment_importBasic(t *testing.T) {
	resourceName := "alicloud_cen_bandwidth_package_attachment.foo"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCenBandwidthPackageAttachmentDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCenBandwidthPackageAttachmentConfig,
			},

			resource.TestStep{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
