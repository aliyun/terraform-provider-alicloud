package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Apig Gateway Security Group Rule. >>> Resource test cases, hand-written.
// Case apig-sgr-basic
func TestAccAliCloudApigGatewaySecurityGroupRule_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_apig_gateway_security_group_rule.default"
	ra := resourceAttrInit(resourceId, AlicloudApigGatewaySecurityGroupRuleMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApigServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeApigGatewaySecurityGroupRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sapigsgr%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudApigGatewaySecurityGroupRuleBasicDependence)
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
					"port_range":        "8080/8080",
					"description":       "tf-testacc-apig-sgr-basic",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"gateway_id":             CHECKSET,
						"security_group_id":      CHECKSET,
						"port_range":             "8080/8080",
						"description":            "tf-testacc-apig-sgr-basic",
						"security_group_rule_id": CHECKSET,
						"ip_protocol":            CHECKSET,
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

var AlicloudApigGatewaySecurityGroupRuleMap = map[string]string{}

func AlicloudApigGatewaySecurityGroupRuleBasicDependence(name string) string {
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
  name   = format("%%s-sg", var.name)
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

// Test Apig Gateway Security Group Rule. <<< Resource test cases, hand-written.
