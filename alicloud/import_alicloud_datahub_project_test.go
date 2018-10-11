package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudDatahubProject_import(t *testing.T) {
	resourceName := "alicloud_datahub_project.basic"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDatahubProjectDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDatahubProject,
			},
			resource.TestStep{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
