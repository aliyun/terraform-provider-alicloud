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

// Test CloudFirewall AiTrafficAnalysisStatus. >>> Resource test cases, automatically generated.
// Case AiTrafficAnalysisStatus 11208
// lintignore: AT001
func TestAccAliCloudCloudFirewallAiTrafficAnalysisStatus_basic11208(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_firewall_ai_traffic_analysis_status.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudCloudFirewallAiTrafficAnalysisStatusMap11208)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cloudfirewall.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCloudFirewallAiTrafficAnalysisStatus")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccloudfirewall%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudCloudFirewallAiTrafficAnalysisStatusBasicDependence11208)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
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
