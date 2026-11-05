package swasopen_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudSimpleApplicationServerFirewallRulesDataSource(t *testing.T) {
	resourceId := "data.alicloud_simple_application_server_firewall_rules.default"
	rand := tfacctest.RandIntRange(1000000, 9999999)
	acctest.CheckoutSupportedRegions(t, true, connectivity.SWASSupportRegions)
	name := fmt.Sprintf("tf-testacc-swasfirewallrule-%d", rand)
	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId, name, dataSourceSimpleApplicationServerFirewallRulesDependence)

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alicloud_simple_application_server_firewall_rule.default.instance_id}",
			"ids":         []string{"${alicloud_simple_application_server_firewall_rule.default.id}"},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alicloud_simple_application_server_firewall_rule.default.instance_id}",
			"ids":         []string{"${alicloud_simple_application_server_firewall_rule.default.id}-fake"},
		}),
	}
	var existSimpleApplicationServerFirewallRuleMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                    "1",
			"ids.0":                    acctest.CHECKSET,
			"rules.#":                  "1",
			"rules.0.id":               acctest.CHECKSET,
			"rules.0.firewall_rule_id": acctest.CHECKSET,
			"rules.0.instance_id":      acctest.CHECKSET,
			"rules.0.port":             "9999",
			"rules.0.rule_protocol":    "Tcp",
			"rules.0.remark":           name,
		}
	}

	var fakeSimpleApplicationServerFirewallRuleMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"rules.#": "0",
		}
	}

	var SimpleApplicationServerFirewallRuleCheckInfo = acctest.DataSourceAttr{
		ResourceId:   resourceId,
		ExistMapFunc: existSimpleApplicationServerFirewallRuleMapFunc,
		FakeMapFunc:  fakeSimpleApplicationServerFirewallRuleMapFunc,
	}

	preCheck := func() {
		acctest.PreCheck(t)
		acctest.PreCheckPrePaidResources(t)
		acctest.PreCheckWithRegions(t, false, connectivity.SimpleApplicationServerNotSupportRegions)
	}
	SimpleApplicationServerFirewallRuleCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf)
}

func dataSourceSimpleApplicationServerFirewallRulesDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_simple_application_server_images" "default" {
	platform = "Linux"
}
data "alicloud_simple_application_server_plans" "default" {
	platform = "Linux"
}

resource "alicloud_simple_application_server_instance" "default" {
  payment_type   = "Subscription"
  plan_id        = data.alicloud_simple_application_server_plans.default.plans.0.id
  instance_name  = var.name
  image_id       = data.alicloud_simple_application_server_images.default.images.0.id
  period         = 1
  data_disk_size = 100
}

resource "alicloud_simple_application_server_firewall_rule" "default" {
  instance_id   = alicloud_simple_application_server_instance.default.id
  rule_protocol = "Tcp"
  port          = "9999"
  remark        = var.name
}`, name)
}
