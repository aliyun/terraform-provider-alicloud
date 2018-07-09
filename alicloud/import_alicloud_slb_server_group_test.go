package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudSlbServerGroup_import(t *testing.T) {
	resourceName := "alicloud_slb_server_group.group"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSlbServerGroupDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccSlbServerGroupClassic,
			},

			resource.TestStep{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"load_balancer_id"},
			},
		},
	})
}
