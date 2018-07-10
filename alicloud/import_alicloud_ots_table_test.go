package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

// At present, OTS sdk does not support creating OTS instance and import test case does not support provider config,
// so this test can not be run successfully.

func SkipTestAccAlicloudOtsTable_importBasic(t *testing.T) {
	resourceName := "alicloud_ots_table.basic"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOtsTableDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOtsTable,
			},
			resource.TestStep{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
