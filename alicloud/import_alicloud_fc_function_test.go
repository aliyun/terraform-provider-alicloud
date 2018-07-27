package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudFCFunction_import(t *testing.T) {
	resourceName := "alicloud_fc_function.foo"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAlicloudFCFunctionDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAlicloudFCFunctionBasic(testFCRoleTemplate),
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
