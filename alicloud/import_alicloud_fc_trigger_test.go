package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudFCTrigger_import(t *testing.T) {
	resourceName := "alicloud_fc_trigger.foo"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckWithRegions(t, false, connectivity.FcNoSupportedRegions) },
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
