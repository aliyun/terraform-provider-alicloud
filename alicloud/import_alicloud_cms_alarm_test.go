package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

// At present, the provider does not support creating contact group resource, so you should add a contact group called "tf-acc-test-group"
// by web console manually before running the following test case.
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
