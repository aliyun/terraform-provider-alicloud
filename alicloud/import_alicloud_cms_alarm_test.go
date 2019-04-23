package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

// At present, the provider does not support creating contact group resource, so you should create manually a contact group
// by web console and set it by environment variable ALICLOUD_CMS_CONTACT_GROUP before running the following test case.
func TestAccAlicloudCmsAlarm_import(t *testing.T) {
	resourceName := "alicloud_cms_alarm.basic"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithCmsContactGroupSetting(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCmsAlarmDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCmsAlarm_basic(cmsContactGroup),
			},

			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dimensions"},
			},
		},
	})
}
