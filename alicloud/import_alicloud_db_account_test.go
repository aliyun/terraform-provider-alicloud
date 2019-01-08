package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudDBAccount_import(t *testing.T) {
	resourceName := "alicloud_db_account.account"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDBAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDBAccount_mysql(DatabaseCommonTestCase),
			},

			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}
