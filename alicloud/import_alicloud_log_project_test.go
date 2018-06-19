package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudLogProject_import(t *testing.T) {
	resourceName := "alicloud_log_project.foo"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAlicloudLogProjectDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAlicloudLogProjectBasic,
			},

			resource.TestStep{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
