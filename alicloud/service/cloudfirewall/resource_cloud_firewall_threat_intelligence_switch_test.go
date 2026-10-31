// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package cloudfirewall_test

import (
	"fmt"
	"testing"

	acctest "github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	connectivity "github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	cloudfirewall "github.com/aliyun/terraform-provider-alicloud/alicloud/service/cloudfirewall"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	resource "github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test CloudFirewall ThreatIntelligenceSwitch. >>> Resource test cases, automatically generated.
// Case ThreatIntelligenceSwitch 11212
// lintignore: AT001
func TestAccAliCloudCloudFirewallThreatIntelligenceSwitch_basic11212(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_firewall_threat_intelligence_switch.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudCloudFirewallThreatIntelligenceSwitchMap11212)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cloudfirewall.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCloudFirewallThreatIntelligenceSwitch")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccloudfirewall%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudCloudFirewallThreatIntelligenceSwitchBasicDependence11212)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
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
