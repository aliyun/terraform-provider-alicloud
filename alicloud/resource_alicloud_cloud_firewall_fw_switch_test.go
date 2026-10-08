// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test CloudFirewall FwSwitch. >>> Resource test cases, automatically generated.
// lintignore: AT001
func TestAccAliCloudCloudFirewallFwSwitch_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_firewall_fw_switch.default"
	ra := resourceAttrInit(resourceId, AlicloudCloudFirewallFwSwitchMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudFirewallServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudFirewallFwSwitch")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccloudfirewall%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCloudFirewallFwSwitchBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"internet_address": "${alicloud_eip_address.default.ip_address}",
					"ip_version":       "4",
					"lang":             "en",
					"member_uid":       "${data.alicloud_account.current.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"internet_address": CHECKSET,
						"ip_version":       "4",
						"member_uid":       CHECKSET,
						"resource_type":    CHECKSET,
						"protect_status":   "open",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"lang"},
			},
		},
	})
}

var AlicloudCloudFirewallFwSwitchMap = map[string]string{}

func AlicloudCloudFirewallFwSwitchBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_account" "current" {
}

resource "alicloud_eip_address" "default" {
}

`, name)
}
