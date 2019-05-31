package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudNasFileSystem_importBasic(t *testing.T) {
	resourceName := "alicloud_nas_file_system.default"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNasDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNasConfig(),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
