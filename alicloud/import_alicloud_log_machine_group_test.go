package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudLogMachineGroup_import(t *testing.T) {
	resourceName := "alicloud_log_machine_group.foo"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAlicloudLogMachineGroupDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAlicloudLogMachineGroupIp(acctest.RandInt()),
			},

			resource.TestStep{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
