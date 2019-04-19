package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudDBDatabase_import(t *testing.T) {
	resourceName := "alicloud_db_database.default"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDBDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDBDatabase_basic(RdsCommonTestCase),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
