// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test CloudFirewall AiTrafficAnalysisStatus. >>> Resource test cases, automatically generated.
// Case AiTrafficAnalysisStatus 11208
func TestAccAliCloudCloudFirewallAiTrafficAnalysisStatus_basic11208(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_firewall_ai_traffic_analysis_status.default"
	ra := resourceAttrInit(resourceId, AlicloudCloudFirewallAiTrafficAnalysisStatusMap11208)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudFirewallServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudFirewallAiTrafficAnalysisStatus")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccloudfirewall%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCloudFirewallAiTrafficAnalysisStatusBasicDependence11208)
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
					"status": "Open",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Open",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Close",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Close",
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

var AlicloudCloudFirewallAiTrafficAnalysisStatusMap11208 = map[string]string{}

func AlicloudCloudFirewallAiTrafficAnalysisStatusBasicDependence11208(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Test CloudFirewall AiTrafficAnalysisStatus. <<< Resource test cases, automatically generated.
