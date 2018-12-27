package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudEssAttachment_import(t *testing.T) {
	resourceName := "alicloud_ess_attachment.attach"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEssAttachmentDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccEssAttachmentConfig(EcsInstanceCommonTestCase, acctest.RandIntRange(1000, 99999)),
			},

			resource.TestStep{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force"},
			},
		},
	})
}
