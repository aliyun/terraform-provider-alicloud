package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudDisk_importBasic(t *testing.T) {
	resourceName := "alicloud_disk.foo"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDiskDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDiskConfig_all,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAlicloudDisk_importWithTags(t *testing.T) {
	resourceName := "alicloud_disk.foo"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDiskDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDiskConfig_tags,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
