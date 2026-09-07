package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAliCloudCloudFirewallControlPoliciesDataSource(t *testing.T) {
	rand := acctest.RandInt()
	aclActionConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudCloudFirewallControlPoliciesDataSourceName(rand, map[string]string{
			"acl_action": `"${alicloud_cloud_firewall_control_policy.default.acl_action}"`,
		}),
		fakeConfig: testAccCheckAliCloudCloudFirewallControlPoliciesDataSourceName(rand, map[string]string{
			"acl_action": `"drop"`,
		}),
	}
	descriptionConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudCloudFirewallControlPoliciesDataSourceName(rand, map[string]string{
			"description": `"${alicloud_cloud_firewall_control_policy.default.description}"`,
		}),
		fakeConfig: testAccCheckAliCloudCloudFirewallControlPoliciesDataSourceName(rand, map[string]string{
			"description": `"${alicloud_cloud_firewall_control_policy.default.description}_fake"`,
		}),
	}
	destinationConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudCloudFirewallControlPoliciesDataSourceName(rand, map[string]string{
			"destination": `"${alicloud_cloud_firewall_control_policy.default.destination}"`,
		}),
		fakeConfig: testAccCheckAliCloudCloudFirewallControlPoliciesDataSourceName(rand, map[string]string{
			"destination": `"${alicloud_cloud_firewall_control_policy.default.destination}_fake"`,
		}),
	}
	ipVersionConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudCloudFirewallControlPoliciesDataSourceName(rand, map[string]string{
			"ip_version": `"${alicloud_cloud_firewall_control_policy.default.ip_version}"`,
		}),
		fakeConfig: testAccCheckAliCloudCloudFirewallControlPoliciesDataSourceName(rand, map[string]string{
			"destination": `"${alicloud_cloud_firewall_control_policy.default.destination}_fake"`,
			"ip_version":  `"6"`,
		}),
	}
	protoConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudCloudFirewallControlPoliciesDataSourceName(rand, map[string]string{
			"proto": `"${alicloud_cloud_firewall_control_policy.default.proto}"`,
		}),
		fakeConfig: testAccCheckAliCloudCloudFirewallControlPoliciesDataSourceName(rand, map[string]string{
			"proto": `"TCP"`,
		}),
	}
	sourceConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudCloudFirewallControlPoliciesDataSourceName(rand, map[string]string{
			"source": `"${alicloud_cloud_firewall_control_policy.default.source}"`,
		}),
		fakeConfig: testAccCheckAliCloudCloudFirewallControlPoliciesDataSourceName(rand, map[string]string{
			"source": `"${alicloud_cloud_firewall_control_policy.default.source}_fake"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudCloudFirewallControlPoliciesDataSourceName(rand, map[string]string{
			"acl_action":  `"${alicloud_cloud_firewall_control_policy.default.acl_action}"`,
			"description": `"${alicloud_cloud_firewall_control_policy.default.description}"`,
			"destination": `"${alicloud_cloud_firewall_control_policy.default.destination}"`,
			"ip_version":  `"${alicloud_cloud_firewall_control_policy.default.ip_version}"`,
			"proto":       `"${alicloud_cloud_firewall_control_policy.default.proto}"`,
			"source":      `"${alicloud_cloud_firewall_control_policy.default.source}"`,
		}),
		fakeConfig: testAccCheckAliCloudCloudFirewallControlPoliciesDataSourceName(rand, map[string]string{
			"acl_action":  `"drop"`,
			"description": `"${alicloud_cloud_firewall_control_policy.default.description}_fake"`,
			"destination": `"${alicloud_cloud_firewall_control_policy.default.destination}_fake"`,
			"ip_version":  `"6"`,
			"proto":       `"TCP"`,
			"source":      `"${alicloud_cloud_firewall_control_policy.default.source}_fake"`,
		}),
	}
	var existAliCloudCloudFirewallControlPoliciesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                              "1",
			"policies.#":                         "1",
			"policies.0.id":                      CHECKSET,
			"policies.0.acl_uuid":                CHECKSET,
			"policies.0.direction":               "out",
			"policies.0.acl_action":              "accept",
			"policies.0.application_id":          CHECKSET,
			"policies.0.application_name":        "ANY",
			"policies.0.description":             CHECKSET,
			"policies.0.dest_port":               CHECKSET,
			"policies.0.dest_port_group":         "",
			"policies.0.dest_port_group_ports":   NOSET,
			"policies.0.dest_port_type":          CHECKSET,
			"policies.0.destination":             "100.1.1.0/24",
			"policies.0.destination_group_cidrs": NOSET,
			"policies.0.destination_group_type":  "",
			"policies.0.destination_type":        "net",
			"policies.0.dns_result":              "",
			"policies.0.dns_result_time":         CHECKSET,
			"policies.0.hit_times":               CHECKSET,
			"policies.0.order":                   CHECKSET,
			"policies.0.proto":                   "ANY",
			"policies.0.release":                 CHECKSET,
			"policies.0.source":                  "1.2.3.0/24",
			"policies.0.source_group_cidrs":      NOSET,
			"policies.0.source_group_type":       "",
			"policies.0.source_type":             "net",
		}
	}
	var fakeAliCloudCloudFirewallControlPoliciesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":      "0",
			"policies.#": "0",
		}
	}
	var alicloudCloudFirewallControlPoliciesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_cloud_firewall_control_policies.default",
		existMapFunc: existAliCloudCloudFirewallControlPoliciesDataSourceNameMapFunc,
		fakeMapFunc:  fakeAliCloudCloudFirewallControlPoliciesDataSourceNameMapFunc,
	}
	preCheck := func() {
		testAccPreCheck(t)
	}
	alicloudCloudFirewallControlPoliciesCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, aclActionConf, descriptionConf, destinationConf, ipVersionConf, protoConf, sourceConf, allConf)
}

func testAccCheckAliCloudCloudFirewallControlPoliciesDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
	variable "description" {
  		default = "tf-testAccCloudFirewallControlPolicies-%d"
	}

	resource "alicloud_cloud_firewall_control_policy" "default" {
  		application_name = "ANY"
  		acl_action       = "accept"
  		description      = var.description
		destination_type = "net"
  		destination      = "100.1.1.0/24"
  		direction        = "out"
  		proto            = "ANY"
  		source           = "1.2.3.0/24"
  		source_type      = "net"
  		ip_version       = "4"
  		lang             = "zh"
	}

	data "alicloud_cloud_firewall_control_policies" "default" {
		direction = alicloud_cloud_firewall_control_policy.default.direction
  		acl_uuid = alicloud_cloud_firewall_control_policy.default.acl_uuid
		%s
	}
`, rand, strings.Join(pairs, " \n "))
	return config
}
