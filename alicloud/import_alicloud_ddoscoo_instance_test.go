package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudDdoscoo_import(t *testing.T) {
	resourceName := "alicloud_ddoscoo_instance.foo"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDdoscooDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDdoscooInstanceConfig(),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
