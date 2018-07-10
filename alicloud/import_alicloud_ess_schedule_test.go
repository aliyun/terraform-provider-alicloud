package alicloud

import (
	"testing"

	"time"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudEssSchedule_importBasic(t *testing.T) {
	resourceName := "alicloud_ess_schedule.foo"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEssScheduleDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccEssScheduleConfig(time.Now().Format("2006-01-02T15:04Z")),
			},

			resource.TestStep{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
