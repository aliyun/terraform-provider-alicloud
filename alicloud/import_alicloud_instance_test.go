package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudInstance_import(t *testing.T) {
	resourceName := "alicloud_instance.foo"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceConfigVPC(EcsInstanceCommonTestCase),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
