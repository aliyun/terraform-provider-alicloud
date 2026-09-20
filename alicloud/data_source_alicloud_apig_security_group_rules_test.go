package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudApigSecurityGroupRules_basic(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sapigsecuritygrouprules%d", defaultRegionToTest, rand)
	resourceId := "data.alicloud_apig_security_group_rules.default"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
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

resource "alicloud_apig_security_group_rule" "default" {
  gateway_id        = alicloud_apig_gateway.defaultgateway.id
  security_group_id = alicloud_security_group.default.id
  port_ranges       = ["8080/8080"]
  description       = var.name
}

data "alicloud_apig_security_group_rules" "default" {
  gateway_id = alicloud_apig_gateway.defaultgateway.id
  depends_on = [alicloud_apig_security_group_rule.default]
}
`, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceId, "rules.#"),
					resource.TestCheckResourceAttrSet(resourceId, "ids.#"),
				),
			},
		},
	})
}
