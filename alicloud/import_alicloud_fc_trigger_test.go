package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

// Import function does not support read account_id from provider.
func SkipTestAccAlicloudFCTrigger_import(t *testing.T) {
	if !isRegionSupports(FunctionCompute) {
		logTestSkippedBecauseOfUnsupportedRegionalFeatures(t.Name(), FunctionCompute)
		return
	}

	resourceName := "alicloud_fc_trigger.foo"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAlicloudFCTriggerDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAlicloudFCTriggerLog(testTriggerLogTemplate, testFCLogRoleTemplate, testFCLogPolicyTemplate, acctest.RandInt()),
			},

			resource.TestStep{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name_prefix", "filename", "oss_bucket", "oss_key"},
			},
		},
	})
}
