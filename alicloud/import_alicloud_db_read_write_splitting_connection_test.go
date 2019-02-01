package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudDBReadWriteSplittingConnection_import(t *testing.T) {
	resourceName := "alicloud_db_read_write_splitting_connection.foo"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDBReadWriteSplittingConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDBReadWriteSplittingConnection_update(RdsCommonTestCase),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
