package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudDBInstance_import(t *testing.T) {
	resourceName := "alicloud_db_instance.default"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDBInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDBInstance_vpc(RdsCommonTestCase),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
