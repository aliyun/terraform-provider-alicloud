package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudNasAccessGroup_importBasic(t *testing.T) {
	resourceName := "alicloud_nas_access_group.foo"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAccessGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNasAccessGroupConfig,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
