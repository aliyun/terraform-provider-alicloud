package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudRamLoginProfile_importBasic(t *testing.T) {
	resourceName := "alicloud_ram_login_profile.default"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRamLoginProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRamLoginProfileCreateConfig(acctest.RandIntRange(1000000, 9999999)),
			},

			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}
