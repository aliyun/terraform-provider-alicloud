package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Apig SecurityGroupRule. >>> Resource test cases, hand-written.
func TestAccAliCloudApigSecurityGroupRule_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_apig_security_group_rule.default"
	ra := resourceAttrInit(resourceId, AlicloudApigSecurityGroupRuleMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApigServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeApigSecurityGroupRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sapigsecuritygrouprule%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudApigSecurityGroupRuleBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"gateway_id":        "${alicloud_apig_gateway.defaultgateway.id}",
					"security_group_id": "${alicloud_security_group.default.id}",
					"port_ranges":       []string{"8080/8080"},
					"description":       name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name,
					}),
					resource.TestCheckResourceAttrSet(resourceId, "security_group_rule_id"),
					resource.TestCheckResourceAttrSet(resourceId, "port_range"),
					resource.TestCheckResourceAttrSet(resourceId, "ip_protocol"),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"port_ranges": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"port_ranges"},
			},
		},
	})
}

var AlicloudApigSecurityGroupRuleMap = map[string]string{}

func AlicloudApigSecurityGroupRuleBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}
data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_security_group" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_apig_gateway" "defaultgateway" {
  network_access_config {
    type = "Intranet"
  }
  vswitch {
    vswitch_id = data.alicloud_vswitches.default.ids.0
  }
  zone_config {
    select_option = "Auto"
  }
  vpc {
    vpc_id = data.alicloud_vpcs.default.ids.0
  }
  payment_type = "PayAsYouGo"
  gateway_name = format("%%s2", var.name)
  spec         = "apigw.small.x1"
  log_config {
    sls {
    }
  }
}
`, name)
}
