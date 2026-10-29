package cloudfirewall_test

import (
	"fmt"
	"strings"
	"testing"

	acctest "github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudCloudFirewallControlPoliciesDataSource(t *testing.T) {
	rand := tfacctest.RandInt()
	aclActionConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudCloudFirewallControlPoliciesDataSourceName(rand, map[string]string{
			"acl_action": `"${alicloud_cloud_firewall_control_policy.default.acl_action}"`,
		}),
		FakeConfig: testAccCheckAliCloudCloudFirewallControlPoliciesDataSourceName(rand, map[string]string{
			"acl_action": `"drop"`,
		}),
	}
	descriptionConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudCloudFirewallControlPoliciesDataSourceName(rand, map[string]string{
			"description": `"${alicloud_cloud_firewall_control_policy.default.description}"`,
		}),
		FakeConfig: testAccCheckAliCloudCloudFirewallControlPoliciesDataSourceName(rand, map[string]string{
			"description": `"${alicloud_cloud_firewall_control_policy.default.description}_fake"`,
		}),
	}
	destinationConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudCloudFirewallControlPoliciesDataSourceName(rand, map[string]string{
			"destination": `"${alicloud_cloud_firewall_control_policy.default.destination}"`,
		}),
		FakeConfig: testAccCheckAliCloudCloudFirewallControlPoliciesDataSourceName(rand, map[string]string{
			"destination": `"${alicloud_cloud_firewall_control_policy.default.destination}_fake"`,
		}),
	}
	ipVersionConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudCloudFirewallControlPoliciesDataSourceName(rand, map[string]string{
			"ip_version": `"${alicloud_cloud_firewall_control_policy.default.ip_version}"`,
		}),
		FakeConfig: testAccCheckAliCloudCloudFirewallControlPoliciesDataSourceName(rand, map[string]string{
			"destination": `"${alicloud_cloud_firewall_control_policy.default.destination}_fake"`,
			"ip_version":  `"6"`,
		}),
	}
	protoConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudCloudFirewallControlPoliciesDataSourceName(rand, map[string]string{
			"proto": `"${alicloud_cloud_firewall_control_policy.default.proto}"`,
		}),
		FakeConfig: testAccCheckAliCloudCloudFirewallControlPoliciesDataSourceName(rand, map[string]string{
			"proto": `"TCP"`,
		}),
	}
	sourceConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudCloudFirewallControlPoliciesDataSourceName(rand, map[string]string{
			"source": `"${alicloud_cloud_firewall_control_policy.default.source}"`,
		}),
		FakeConfig: testAccCheckAliCloudCloudFirewallControlPoliciesDataSourceName(rand, map[string]string{
			"source": `"${alicloud_cloud_firewall_control_policy.default.source}_fake"`,
		}),
	}
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudCloudFirewallControlPoliciesDataSourceName(rand, map[string]string{
			"acl_action":  `"${alicloud_cloud_firewall_control_policy.default.acl_action}"`,
			"description": `"${alicloud_cloud_firewall_control_policy.default.description}"`,
			"destination": `"${alicloud_cloud_firewall_control_policy.default.destination}"`,
			"ip_version":  `"${alicloud_cloud_firewall_control_policy.default.ip_version}"`,
			"proto":       `"${alicloud_cloud_firewall_control_policy.default.proto}"`,
			"source":      `"${alicloud_cloud_firewall_control_policy.default.source}"`,
		}),
		FakeConfig: testAccCheckAliCloudCloudFirewallControlPoliciesDataSourceName(rand, map[string]string{
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
			"ids.#":                             "1",
			"policies.#":                        "1",
			"policies.0.id":                     acctest.CHECKSET,
			"policies.0.acl_uuid":               acctest.CHECKSET,
			"policies.0.direction":              "out",
			"policies.0.acl_action":             "accept",
			"policies.0.application_id":         acctest.CHECKSET,
			"policies.0.application_name":       "ANY",
			"policies.0.description":            acctest.CHECKSET,
			"policies.0.dest_port":              acctest.CHECKSET,
			"policies.0.dest_port_group":        "",
			"policies.0.dest_port_type":         acctest.CHECKSET,
			"policies.0.destination":            "100.1.1.0/24",
			"policies.0.destination_group_type": "",
			"policies.0.destination_type":       "net",
			"policies.0.dns_result":             "",
			"policies.0.dns_result_time":        acctest.CHECKSET,
			"policies.0.hit_times":              acctest.CHECKSET,
			"policies.0.order":                  acctest.CHECKSET,
			"policies.0.proto":                  "ANY",
			"policies.0.release":                acctest.CHECKSET,
			"policies.0.source":                 "1.2.3.0/24",
			"policies.0.source_group_type":      "",
			"policies.0.source_type":            "net",
		}
	}
	var fakeAliCloudCloudFirewallControlPoliciesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":      "0",
			"policies.#": "0",
		}
	}
	var alicloudCloudFirewallControlPoliciesCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_cloud_firewall_control_policies.default",
		ExistMapFunc: existAliCloudCloudFirewallControlPoliciesDataSourceNameMapFunc,
		FakeMapFunc:  fakeAliCloudCloudFirewallControlPoliciesDataSourceNameMapFunc,
	}
	preCheck := func() {
		acctest.PreCheck(t)
	}
	alicloudCloudFirewallControlPoliciesCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, aclActionConf, descriptionConf, destinationConf, ipVersionConf, protoConf, sourceConf, allConf)
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
