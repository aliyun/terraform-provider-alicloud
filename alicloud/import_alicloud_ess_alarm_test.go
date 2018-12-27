package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudEssAlarm_import(t *testing.T) {
	resourceName := "alicloud_ess_alarm.foo"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEssAlarmDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEssAlarm_basic(EcsInstanceCommonTestCase, acctest.RandIntRange(1000, 999999)),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
