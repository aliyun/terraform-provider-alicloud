package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

// At present, only white list users can operate Common Bandwidth Package Resource.
func SkipTestAccAlicloudCommonBandwidthPackageAttachment_importBasic(t *testing.T) {
	resourceName := "alicloud_common_bandwidth_package_attachment.foo"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCommonBandwidthPackageAttachmentDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCommonBandwidthPackageAttachmentConfig,
			},

			resource.TestStep{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
