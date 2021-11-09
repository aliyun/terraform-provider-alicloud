package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudSimpleApplicationServerFirewallRule_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_simple_application_server_firewall_rule.default"
	ra := resourceAttrInit(resourceId, AlicloudSimpleApplicationServerFirewallRuleMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SwasOpenService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSimpleApplicationServerFirewallRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc-simpleapplicationserverfirewallrule%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudSimpleApplicationServerFirewallRuleBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":   "${local.instance_id}",
					"rule_protocol": "TcpAndUdp",
					"port":          "1024/1055",
					"remark":        "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":   CHECKSET,
						"rule_protocol": "TcpAndUdp",
						"port":          "1024/1055",
						"remark":        name,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

var AlicloudSimpleApplicationServerFirewallRuleMap0 = map[string]string{
	"instance_id":      CHECKSET,
	"firewall_rule_id": CHECKSET,
	"port":             CHECKSET,
	"rule_protocol":    CHECKSET,
}

func AlicloudSimpleApplicationServerFirewallRuleBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}

data "alicloud_simple_application_server_instances" "default" {}

data "alicloud_simple_application_server_images" "default" {}

data "alicloud_simple_application_server_plans" "default" {}

resource "alicloud_simple_application_server_instance" "default" {
  count          = length(data.alicloud_simple_application_server_instances.default.ids) > 0 ? 0 : 1
  payment_type   = "Subscription"
  plan_id        = data.alicloud_simple_application_server_plans.default.plans.0.id
  instance_name  = "tf-testaccswas-firewallrule"
  image_id       = data.alicloud_simple_application_server_images.default.images.0.id
  period         = 1
  data_disk_size = 100
}

locals {
  instance_id = length(data.alicloud_simple_application_server_instances.default.ids) > 0 ? data.alicloud_simple_application_server_instances.default.ids.0 : alicloud_simple_application_server_instance.default.0.id
}
`, name)
}
