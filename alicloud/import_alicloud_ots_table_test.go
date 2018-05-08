package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudOtsTable_importBasic(t *testing.T) {
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
