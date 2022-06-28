package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudGaForwardingRulesDataSource(t *testing.T) {
	resourceId := "data.alicloud_ga_forwarding_rules.default"
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccAlicloudGaListenerDataSource%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceGaForwardingRulesConfigDependence)
	checkoutSupportedRegions(t, true, connectivity.GaSupportRegions)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"accelerator_id": "${alicloud_ga_forwarding_rule.default.accelerator_id}",
			"listener_id":    "${alicloud_ga_forwarding_rule.default.listener_id}",
			"ids":            []string{"${alicloud_ga_forwarding_rule.default.forwarding_rule_id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"accelerator_id": "${alicloud_ga_forwarding_rule.default.accelerator_id}",
			"listener_id":    "${alicloud_ga_forwarding_rule.default.listener_id}",
			"ids":            []string{"${alicloud_ga_forwarding_rule.default.forwarding_rule_id}_fake"},
		}),
	}

	statusConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"accelerator_id": "${alicloud_ga_forwarding_rule.default.accelerator_id}",
			"listener_id":    "${alicloud_ga_forwarding_rule.default.listener_id}",
			"ids":            []string{"${alicloud_ga_forwarding_rule.default.forwarding_rule_id}"},
			"status":         "active",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"accelerator_id": "${alicloud_ga_forwarding_rule.default.accelerator_id}",
			"listener_id":    "${alicloud_ga_forwarding_rule.default.listener_id}",
			"ids":            []string{"${alicloud_ga_forwarding_rule.default.forwarding_rule_id}_fake"},
			"status":         "configuring",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"accelerator_id": "${alicloud_ga_forwarding_rule.default.accelerator_id}",
			"listener_id":    "${alicloud_ga_forwarding_rule.default.listener_id}",
			"ids":            []string{"${alicloud_ga_forwarding_rule.default.forwarding_rule_id}"},
			"status":         "active",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"accelerator_id": "${alicloud_ga_forwarding_rule.default.accelerator_id}",
			"listener_id":    "${alicloud_ga_forwarding_rule.default.listener_id}",
			"ids":            []string{"${alicloud_ga_forwarding_rule.default.forwarding_rule_id}_fake"},
			"status":         "configuring",
		}),
	}

	var existMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                     "1",
			"forwarding_rules.#":                        "1",
			"forwarding_rules.0.priority":               "1",
			"forwarding_rules.0.forwarding_rule_id":     CHECKSET,
			"forwarding_rules.0.forwarding_rule_name":   "",
			"forwarding_rules.0.forwarding_rule_status": "active",
			"forwarding_rules.0.listener_id":            CHECKSET,
			"forwarding_rules.0.rule_conditions.#":      "1",
			"forwarding_rules.0.rule_actions.#":         "1",
			"forwarding_rules.0.id":                     CHECKSET,
		}
	}

	var fakeMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"forwarding_rules.#": "0",
			"ids.#":              "0",
		}
	}

	var CheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existMapFunc,
		fakeMapFunc:  fakeMapFunc,
	}
	preCheck := func() {}

	CheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, statusConf, allConf)
}

func dataSourceGaForwardingRulesConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default  = "%s"
}

data "alicloud_ga_accelerators" "default" {
  status = "active"
}

resource "alicloud_ga_bandwidth_package" "default" {
   	bandwidth              =  100
  	type                   = "Basic"
  	bandwidth_type         = "Basic"
	payment_type           = "PayAsYouGo"
  	billing_type           = "PayBy95"
	ratio       = 30
	bandwidth_package_name = var.name
    auto_pay               = true
    auto_use_coupon        = true
}

resource "alicloud_ga_bandwidth_package_attachment" "default" {
	// Please run resource ga_accelerator test case to ensure this account has at least one accelerator before run this case.
	accelerator_id = data.alicloud_ga_accelerators.default.ids.0
	bandwidth_package_id = alicloud_ga_bandwidth_package.default.id
}

resource "alicloud_ga_listener" "default" {
  port_ranges{
    from_port="70"
    to_port="70"
  }
  accelerator_id=alicloud_ga_bandwidth_package_attachment.default.accelerator_id
  client_affinity="SOURCE_IP"
  protocol="HTTP"
  name=var.name
}

resource "alicloud_eip_address" "default" {
  bandwidth            = "10"
  internet_charge_type = "PayByBandwidth"
  address_name = var.name
}

resource "alicloud_ga_endpoint_group" "default" {
  accelerator_id=alicloud_ga_bandwidth_package_attachment.default.accelerator_id
  endpoint_configurations{
    endpoint=alicloud_eip_address.default.ip_address
    type="PublicIp"
    weight="20"
  }
  description=var.name
  name=var.name
  threshold_count=4
  endpoint_group_region="%s"
  health_check_interval_seconds="3"
  health_check_path="/healthcheck"
  health_check_port="9999"
  health_check_protocol="http"
  port_overrides{
    endpoint_port="10"
    listener_port="70"
  }
  traffic_percentage=20
  listener_id=alicloud_ga_listener.default.id
  endpoint_group_type = "virtual"
}

resource "alicloud_ga_forwarding_rule" "default" {
  accelerator_id = alicloud_ga_endpoint_group.default.accelerator_id
  listener_id    = alicloud_ga_endpoint_group.default.listener_id
  rule_conditions {
    rule_condition_type = "Path"
    path_config {
      values = ["/test"]
    }
  }
  rule_actions {
    order            = "30"
    rule_action_type = "ForwardGroup"
    forward_group_config {
      server_group_tuples {
        endpoint_group_id = alicloud_ga_endpoint_group.default.id
      }
    }
  }
}
`, name, defaultRegionToTest)
}
