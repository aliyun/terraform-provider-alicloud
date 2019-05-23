package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudEssScalingruleImport(t *testing.T) {
	resourceName := "alicloud_ess_scaling_rule.default"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEssScalingRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEssScalingRuleConfig(EcsInstanceCommonTestCase, acctest.RandIntRange(1000, 999999)),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
