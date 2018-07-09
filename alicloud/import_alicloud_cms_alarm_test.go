package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudCmsAlarm_import(t *testing.T) {
	resourceName := "alicloud_cms_alarm.basic"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCmsAlarmDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCmsAlarm_basic,
			},

			resource.TestStep{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dimensions"},
			},
		},
	})
}
