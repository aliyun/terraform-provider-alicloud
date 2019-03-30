package alicloud

import (
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

func TestAccAlicloudSnapshot_importBasic(t *testing.T) {
	resourceName := "alicloud_snapshot.snapshot"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSnapshotDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccSnapshotConfig,
			},

			resource.TestStep{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
