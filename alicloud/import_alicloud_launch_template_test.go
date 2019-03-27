package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudLaunchTemplate_importBasic(t *testing.T) {
	resourceName := "alicloud_launch_template.template"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLaunchTemplateDestroy,
		Steps: []resource.TestStep{
			{
				Config: lauchTemplateConfig,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
