package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudFCService_import(t *testing.T) {
	if !isRegionSupports(FunctionCompute) {
		logTestSkippedBecauseOfUnsupportedRegionalFeatures(t.Name(), FunctionCompute)
		return
	}

	resourceName := "alicloud_fc_service.foo"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAlicloudFCServiceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAlicloudFCServiceBasic("tf-testaccalicloudfcserviceimport", testFCRoleTemplate),
			},

			resource.TestStep{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name_prefix"},
			},
		},
	})
}
