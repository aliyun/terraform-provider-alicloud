package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudDBConnection_import(t *testing.T) {
	resourceName := "alicloud_db_connection.foo"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDBConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDBConnection_basic(RdsCommonTestCase),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
