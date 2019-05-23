package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudEssLifecycleHook_import(t *testing.T) {
	resourceName := "alicloud_ess_lifecycle_hook.default"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEssLifecycleHookDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEssLifecycleHook(EcsInstanceCommonTestCase, acctest.RandIntRange(1000, 999999)),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
