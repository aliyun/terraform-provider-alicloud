// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test CloudFirewall ThreatIntelligenceSwitch. >>> Resource test cases, automatically generated.
// Case ThreatIntelligenceSwitch 11212
func TestAccAliCloudCloudFirewallThreatIntelligenceSwitch_basic11212(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_firewall_threat_intelligence_switch.default"
	ra := resourceAttrInit(resourceId, AlicloudCloudFirewallThreatIntelligenceSwitchMap11212)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudFirewallServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudFirewallThreatIntelligenceSwitch")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccloudfirewall%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCloudFirewallThreatIntelligenceSwitchBasicDependence11212)
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
					"action":        "alert",
					"enable_status": "0",
					"category_id":   "IpOutThreatTorExit",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"action":        "alert",
						"enable_status": "0",
						"category_id":   "IpOutThreatTorExit",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"action":        "drop",
					"enable_status": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"action":        "drop",
						"enable_status": "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudCloudFirewallThreatIntelligenceSwitchMap11212 = map[string]string{}

func AlicloudCloudFirewallThreatIntelligenceSwitchBasicDependence11212(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Test CloudFirewall ThreatIntelligenceSwitch. <<< Resource test cases, automatically generated.
