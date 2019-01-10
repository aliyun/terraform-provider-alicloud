package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudDBAccountPrivilege_import(t *testing.T) {
	resourceName := "alicloud_db_account_privilege.privilege"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDBAccountPrivilegeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDBAccountPrivilege_basic(DatabaseCommonTestCase),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
