package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudRamLoginProfile_importBasic(t *testing.T) {
	resourceName := "alicloud_ram_login_profile.profile"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRamLoginProfileDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRamLoginProfileConfig,
			},

			resource.TestStep{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}
