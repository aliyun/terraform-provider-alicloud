package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudEssScalingGroup_importBasic(t *testing.T) {
	resourceName := "alicloud_ess_scaling_group.foo"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEssScalingGroupDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccEssScalingGroupConfig,
			},

			resource.TestStep{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
