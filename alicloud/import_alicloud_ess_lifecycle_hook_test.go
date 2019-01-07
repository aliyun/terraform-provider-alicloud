package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudEssLifecycleHook_import(t *testing.T) {
	resourceName := "alicloud_ess_lifecycle_hook.foo"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEssLifecycleHookDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccEssLifecycleHook(EcsInstanceCommonTestCase, acctest.RandIntRange(1000, 999999)),
			},

			resource.TestStep{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
