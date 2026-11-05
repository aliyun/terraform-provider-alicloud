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

// Test CloudFirewall PolicyAdvancedConfig. >>> Resource test cases, automatically generated.
// Case 云防火墙Terraform启用严格模式Strict Mode 10921
// lintignore: AT001
func TestAccAliCloudCloudFirewallPolicyAdvancedConfig_basic10921(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_firewall_policy_advanced_config.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudCloudFirewallPolicyAdvancedConfigMap10921)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cloudfirewall.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCloudFirewallPolicyAdvancedConfig")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccloudfirewall%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudCloudFirewallPolicyAdvancedConfigBasicDependence10921)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
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
