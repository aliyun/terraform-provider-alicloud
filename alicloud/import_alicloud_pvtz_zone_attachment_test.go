package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudPvtzZoneAttachment_importBasic(t *testing.T) {
	resourceName := "alicloud_pvtz_zone_attachment.zone-attachment"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccAlicloudPvtzZoneAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPvtzZoneAttachmentConfig(acctest.RandIntRange(10000, 999999)),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
