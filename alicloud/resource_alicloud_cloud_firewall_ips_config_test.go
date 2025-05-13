package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test CloudFirewall IPSConfig. >>> Resource test cases, automatically generated.
// Case 修改IPS拦截模式 10240
func TestAccAliCloudCloudFirewallIPSConfig_basic10240(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_firewall_ips_config.default"
	ra := resourceAttrInit(resourceId, AlicloudCloudFirewallIPSConfigMap10240)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudFirewallServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudFirewallIPSConfig")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccloudfirewall%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCloudFirewallIPSConfigBasicDependence10240)
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
					"max_sdl":     "1000",
					"basic_rules": "1",
					"run_mode":    "1",
					"cti_rules":   "0",
					"patch_rules": "0",
					"rule_class":  "1",
					"lang":        "zh",
					"depends_on": []string{
						"alicloud_cloud_firewall_instance.default",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"max_sdl":     "1000",
						"basic_rules": "1",
						"run_mode":    "1",
						"cti_rules":   "0",
						"patch_rules": "0",
						"rule_class":  "1",
						"lang":        "zh",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"max_sdl":     "1000",
					"basic_rules": "1",
					"run_mode":    "1",
					"cti_rules":   "1",
					"patch_rules": "1",
					"rule_class":  "1",
					"lang":        "zh",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"max_sdl":     "1000",
						"basic_rules": "1",
						"run_mode":    "1",
						"cti_rules":   "1",
						"patch_rules": "1",
						"rule_class":  "1",
						"lang":        "zh",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"max_sdl": "2000",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"max_sdl": "2000",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"run_mode":    "0",
					"rule_class":  "0",
					"basic_rules": "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"run_mode":    "0",
						"rule_class":  "0",
						"basic_rules": "0",
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

var AlicloudCloudFirewallIPSConfigMap10240 = map[string]string{}

func AlicloudCloudFirewallIPSConfigBasicDependence10240(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_cloud_firewall_instance" "default" {
  payment_type = "PayAsYouGo"
}
`, name)
}

// Test CloudFirewall IPSConfig. <<< Resource test cases, automatically generated.
