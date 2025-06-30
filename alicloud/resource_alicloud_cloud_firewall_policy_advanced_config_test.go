// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test CloudFirewall PolicyAdvancedConfig. >>> Resource test cases, automatically generated.
// Case 云防火墙Terraform启用严格模式Strict Mode 10921
func TestAccAliCloudCloudFirewallPolicyAdvancedConfig_basic10921(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_firewall_policy_advanced_config.default"
	ra := resourceAttrInit(resourceId, AlicloudCloudFirewallPolicyAdvancedConfigMap10921)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudFirewallServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudFirewallPolicyAdvancedConfig")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccloudfirewall%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCloudFirewallPolicyAdvancedConfigBasicDependence10921)
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
					"internet_switch": "off",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"internet_switch": "off",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"internet_switch": "on",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"internet_switch": "on",
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

var AlicloudCloudFirewallPolicyAdvancedConfigMap10921 = map[string]string{}

func AlicloudCloudFirewallPolicyAdvancedConfigBasicDependence10921(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Test CloudFirewall PolicyAdvancedConfig. <<< Resource test cases, automatically generated.
