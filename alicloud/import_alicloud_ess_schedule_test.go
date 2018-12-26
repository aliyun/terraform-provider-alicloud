package alicloud

import (
	"testing"

	"time"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudEssSchedule_importBasic(t *testing.T) {
	resourceName := "alicloud_ess_schedule.foo"
	// Setting schedule time to more than one day
	oneDay, _ := time.ParseDuration("24h")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEssScheduleDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccEssScheduleConfig(time.Now().Add(oneDay).Format("2006-01-02T15:04Z")),
			},

			resource.TestStep{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
