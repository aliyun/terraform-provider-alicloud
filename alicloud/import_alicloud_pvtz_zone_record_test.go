package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudPvtzZoneRecord_importBasic(t *testing.T) {
	resourceName := "alicloud_pvtz_zone_record.foo"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccAlicloudPvtzZoneRecordDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccPvtzZoneRecordConfig,
			},

			resource.TestStep{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
