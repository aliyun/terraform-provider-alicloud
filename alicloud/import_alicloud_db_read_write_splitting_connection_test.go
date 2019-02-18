package alicloud

import (
	"github.com/hashicorp/terraform/helper/acctest"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudDBReadWriteSplittingConnection_import(t *testing.T) {
	resourceName := "alicloud_db_read_write_splitting_connection.foo"
	randomPrefix := acctest.RandIntRange(10000, 999999)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDBReadWriteSplittingConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDBReadWriteSplittingConnection_basic(testAccDBReadonlyInstance_vpc(testAccDBRInstance_vpc(RdsCommonTestCase)), randomPrefix),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
