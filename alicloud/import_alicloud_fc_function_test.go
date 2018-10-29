package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudFCFunction_import(t *testing.T) {
	if !isRegionSupports(FunctionCompute) {
		logTestSkippedBecauseOfUnsupportedRegionalFeatures(t.Name(), FunctionCompute)
		return
	}

	resourceName := "alicloud_fc_function.foo"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAlicloudFCFunctionDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAlicloudFCFunctionBasic(testFCRoleTemplate, acctest.RandInt()),
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
