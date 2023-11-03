package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccCheckAlicloudCloudFirewallControlPoliciesDataSource(t *testing.T) {
	rand := acctest.RandInt()

	var existAlicloudCloudFirewallControlPoliciesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                              "1",
			"policies.#":                         "1",
			"policies.0.description":             fmt.Sprintf("tf-testAccCloudFirewallControlPolicies-%d", rand),
			"policies.0.application_name":        "ANY",
			"policies.0.acl_action":              "accept",
			"policies.0.destination_type":        "net",
			"policies.0.destination":             "100.1.1.0/24",
			"policies.0.direction":               "out",
			"policies.0.proto":                   "ANY",
			"policies.0.source":                  "1.2.3.0/24",
			"policies.0.source_type":             "net",
			"policies.0.release":                 CHECKSET,
			"policies.0.acl_uuid":                CHECKSET,
			"policies.0.application_id":          CHECKSET,
			"policies.0.dest_port":               CHECKSET,
			"policies.0.dest_port_group":         CHECKSET,
			"policies.0.dest_port_group_ports":   CHECKSET,
			"policies.0.dest_port_type":          CHECKSET,
			"policies.0.destination_group_cidrs": CHECKSET,
			"policies.0.destination_group_type":  CHECKSET,
			"dns_result":                         CHECKSET,
			"policies.0.dns_result_time":         CHECKSET,
			"policies.0.hit_times":               CHECKSET,
			"policies.0.order":                   CHECKSET,
			"policies.0.source_group_cidrs":      CHECKSET,
			"policies.0.source_group_type":       CHECKSET,
		}
	}
	var fakeAlicloudCloudFirewallControlPoliciesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudEventBridgeEventBusesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_cloud_firewall_control_policies.default",
		existMapFunc: existAlicloudCloudFirewallControlPoliciesDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudCloudFirewallControlPoliciesDataSourceNameMapFunc,
	}
	alicloudEventBridgeEventBusesCheckInfo.dataSourceTestCheck(t, rand)
}
func testAccCheckAlicloudCloudFirewallControlPoliciesDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "description" {	
	default = "tf-testAccCloudFirewallControlPolicies-%d"
}

resource "alicloud_cloud_firewall_control_policy" "default" {
	application_name =  "ANY"
	acl_action       =  "accept"
	description      =  var.description
	destination_type =  "net"
	destination      =  "100.1.1.0/24"
	direction        =  "out"
	proto            =  "ANY"
	source           =  "1.2.3.0/24"
	source_type      =  "net"
}

data "alicloud_cloud_firewall_control_policies" "default" {	
	direction = alicloud_cloud_firewall_control_policy.default.direction
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
