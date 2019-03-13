package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudRamGroup_importBasic(t *testing.T) {
	resourceName := "alicloud_ram_group.group"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRamGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRamGroupConfig(acctest.RandInt()),
			},

			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force"},
			},
		},
	})
}
