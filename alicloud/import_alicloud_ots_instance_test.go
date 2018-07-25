package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudOtsInstance_import(t *testing.T) {
	resourceName := "alicloud_ots_instance.basic"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOtsInstanceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOtsInstance,
			},
			resource.TestStep{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
