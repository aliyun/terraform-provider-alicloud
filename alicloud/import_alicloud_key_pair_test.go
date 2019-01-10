package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudKeyPair_importBasic(t *testing.T) {
	resourceName := "alicloud_key_pair.basic"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKeyPairDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKeyPairConfig,
			},

			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"key_file"},
			},
		},
	})
}
