package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudDBReadonlyInstance_import(t *testing.T) {
	resourceName := "alicloud_db_readonly_instance.default"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDBReadonlyInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDBReadonlyInstance_vpc(testAccDBRInstance_vpc(RdsCommonTestCase)),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
