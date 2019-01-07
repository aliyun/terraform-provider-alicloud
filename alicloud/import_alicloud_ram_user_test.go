package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudRamUser_importBasic(t *testing.T) {
	resourceName := "alicloud_ram_user.user"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRamUserDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRamUserConfig(acctest.RandInt()),
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
