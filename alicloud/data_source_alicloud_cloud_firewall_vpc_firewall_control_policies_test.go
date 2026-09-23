package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAliCloudCloudFirewallVpcFirewallControlPoliciesDataSource(t *testing.T) {
	rand := acctest.RandInt()
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCloudFirewallVpcFirewallControlPoliciesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_cloud_firewall_vpc_firewall_control_policy.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudCloudFirewallVpcFirewallControlPoliciesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_cloud_firewall_vpc_firewall_control_policy.default.id}_fake"]`,
		}),
	}
	aclActionConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCloudFirewallVpcFirewallControlPoliciesDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_cloud_firewall_vpc_firewall_control_policy.default.id}"]`,
			"acl_action": `"accept"`,
		}),
		fakeConfig: testAccCheckAlicloudCloudFirewallVpcFirewallControlPoliciesDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_cloud_firewall_vpc_firewall_control_policy.default.id}_fake"]`,
			"acl_action": `"drop"`,
		}),
	}
	aclUuidConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCloudFirewallVpcFirewallControlPoliciesDataSourceName(rand, map[string]string{
			"ids":      `["${alicloud_cloud_firewall_vpc_firewall_control_policy.default.id}"]`,
			"acl_uuid": `"${alicloud_cloud_firewall_vpc_firewall_control_policy.default.acl_uuid}"`,
		}),
		fakeConfig: testAccCheckAlicloudCloudFirewallVpcFirewallControlPoliciesDataSourceName(rand, map[string]string{
			"ids":      `["${alicloud_cloud_firewall_vpc_firewall_control_policy.default.id}"]`,
			"acl_uuid": `"${alicloud_cloud_firewall_vpc_firewall_control_policy.default.acl_uuid}_fake"`,
		}),
	}
	descriptionConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCloudFirewallVpcFirewallControlPoliciesDataSourceName(rand, map[string]string{
			"ids":         `["${alicloud_cloud_firewall_vpc_firewall_control_policy.default.id}"]`,
			"description": `"${alicloud_cloud_firewall_vpc_firewall_control_policy.default.description}"`,
		}),
		fakeConfig: testAccCheckAlicloudCloudFirewallVpcFirewallControlPoliciesDataSourceName(rand, map[string]string{
			"ids":         `["${alicloud_cloud_firewall_vpc_firewall_control_policy.default.id}"]`,
			"description": `"${alicloud_cloud_firewall_vpc_firewall_control_policy.default.description}_fake"`,
		}),
	}
	destinationConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCloudFirewallVpcFirewallControlPoliciesDataSourceName(rand, map[string]string{
			"ids":         `["${alicloud_cloud_firewall_vpc_firewall_control_policy.default.id}"]`,
			"destination": `"${alicloud_cloud_firewall_vpc_firewall_control_policy.default.destination}"`,
		}),
		fakeConfig: testAccCheckAlicloudCloudFirewallVpcFirewallControlPoliciesDataSourceName(rand, map[string]string{
			"ids":         `["${alicloud_cloud_firewall_vpc_firewall_control_policy.default.id}"]`,
			"destination": `"${alicloud_cloud_firewall_vpc_firewall_control_policy.default.destination}_fake"`,
		}),
	}
	memberUidConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCloudFirewallVpcFirewallControlPoliciesDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_cloud_firewall_vpc_firewall_control_policy.default.id}"]`,
			"member_uid": `"${alicloud_cloud_firewall_vpc_firewall_control_policy.default.member_uid}"`,
		}),
		fakeConfig: "",
	}
	protoConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCloudFirewallVpcFirewallControlPoliciesDataSourceName(rand, map[string]string{
			"ids":   `["${alicloud_cloud_firewall_vpc_firewall_control_policy.default.id}"]`,
			"proto": `"${alicloud_cloud_firewall_vpc_firewall_control_policy.default.proto}"`,
		}),
		fakeConfig: testAccCheckAlicloudCloudFirewallVpcFirewallControlPoliciesDataSourceName(rand, map[string]string{
			"ids":   `["${alicloud_cloud_firewall_vpc_firewall_control_policy.default.id}"]`,
			"proto": `"UDP"`,
		}),
	}
	releaseConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCloudFirewallVpcFirewallControlPoliciesDataSourceName(rand, map[string]string{
			"ids":     `["${alicloud_cloud_firewall_vpc_firewall_control_policy.default.id}"]`,
			"release": `"${alicloud_cloud_firewall_vpc_firewall_control_policy.default.release}"`,
		}),
		fakeConfig: testAccCheckAlicloudCloudFirewallVpcFirewallControlPoliciesDataSourceName(rand, map[string]string{
			"ids":     `["${alicloud_cloud_firewall_vpc_firewall_control_policy.default.id}"]`,
			"release": `"false"`,
		}),
	}
	sourceConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCloudFirewallVpcFirewallControlPoliciesDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_cloud_firewall_vpc_firewall_control_policy.default.id}"]`,
			"source": `"${alicloud_cloud_firewall_vpc_firewall_control_policy.default.source}"`,
		}),
		fakeConfig: testAccCheckAlicloudCloudFirewallVpcFirewallControlPoliciesDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_cloud_firewall_vpc_firewall_control_policy.default.id}"]`,
			"source": `"127.0.0.2/32"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCloudFirewallVpcFirewallControlPoliciesDataSourceName(rand, map[string]string{
			"ids":         `["${alicloud_cloud_firewall_vpc_firewall_control_policy.default.id}"]`,
			"acl_action":  `"accept"`,
			"acl_uuid":    `"${alicloud_cloud_firewall_vpc_firewall_control_policy.default.acl_uuid}"`,
			"description": `"${alicloud_cloud_firewall_vpc_firewall_control_policy.default.description}"`,
			"member_uid":  `"${alicloud_cloud_firewall_vpc_firewall_control_policy.default.member_uid}"`,
			"proto":       `"${alicloud_cloud_firewall_vpc_firewall_control_policy.default.proto}"`,
			"release":     `"${alicloud_cloud_firewall_vpc_firewall_control_policy.default.release}"`,
			"source":      `"${alicloud_cloud_firewall_vpc_firewall_control_policy.default.source}"`,
		}),
		fakeConfig: testAccCheckAlicloudCloudFirewallVpcFirewallControlPoliciesDataSourceName(rand, map[string]string{
			"ids":         `["${alicloud_cloud_firewall_vpc_firewall_control_policy.default.id}_fake"]`,
			"acl_action":  `"drop"`,
			"acl_uuid":    `"${alicloud_cloud_firewall_vpc_firewall_control_policy.default.acl_uuid}_fake"`,
			"description": `"${alicloud_cloud_firewall_vpc_firewall_control_policy.default.description}_fake"`,
			"proto":       `"UDP"`,
			"release":     `"false"`,
			"source":      `"127.0.0.2/32"`,
		}),
	}
	var existAlicloudCloudFirewallVpcFirewallControlPoliciesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                "1",
			"policies.#":                           "1",
			"policies.0.id":                        CHECKSET,
			"policies.0.acl_action":                "accept",
			"policies.0.acl_uuid":                  CHECKSET,
			"policies.0.application_id":            CHECKSET,
			"policies.0.application_name":          "ANY",
			"policies.0.description":               fmt.Sprintf("tf-testAccVpcFirewallControlPolicy-%d", rand),
			"policies.0.dest_port":                 "80/88",
			"policies.0.dest_port_group":           "",
			"policies.0.dest_port_group_ports.#":   "0",
			"policies.0.dest_port_type":            "port",
			"policies.0.destination":               "127.0.0.2/32",
			"policies.0.destination_group_cidrs.#": "0",
			"policies.0.destination_group_type":    "",
			"policies.0.destination_type":          "net",
			"policies.0.hit_times":                 CHECKSET,
			"policies.0.member_uid":                CHECKSET,
			"policies.0.order":                     "1",
			"policies.0.proto":                     "TCP",
			"policies.0.release":                   "true",
			"policies.0.source":                    "127.0.0.1/32",
			"policies.0.source_group_cidrs.#":      "0",
			"policies.0.source_group_type":         "",
			"policies.0.source_type":               "net",
			"policies.0.vpc_firewall_id":           CHECKSET,
		}
	}
	var fakeAlicloudCloudFirewallVpcFirewallControlPoliciesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":      "0",
			"policies.#": "0",
		}
	}
	var AlicloudCloudFirewallVpcFirewallControlPoliciesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_cloud_firewall_vpc_firewall_control_policies.default",
		existMapFunc: existAlicloudCloudFirewallVpcFirewallControlPoliciesDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudCloudFirewallVpcFirewallControlPoliciesDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}
	AlicloudCloudFirewallVpcFirewallControlPoliciesCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, aclActionConf, aclUuidConf, descriptionConf, destinationConf, memberUidConf, protoConf, releaseConf, sourceConf, allConf)
}
func testAccCheckAlicloudCloudFirewallVpcFirewallControlPoliciesDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {
	default = "tf-testAccVpcFirewallControlPolicy-%d"
}

data "alicloud_account" "default" {}

resource "alicloud_cen_instance" "default" {
	cen_instance_name = var.name
	description = "tf-testAccCenConfigDescription"
	tags 		= {
		Created = "TF"
		For 	= "acceptance test"
	}
}

resource "alicloud_cloud_firewall_vpc_firewall_control_policy" "default" {
  order            = "1"
  destination      = "127.0.0.2/32"
  application_name = "ANY"
  description      = var.name
  source_type      = "net"
  dest_port        = "80/88"
  acl_action       = "accept"
  lang             = "zh"
  destination_type = "net"
  source           = "127.0.0.1/32"
  dest_port_type   = "port"
  proto            = "TCP"
  release          = true
  member_uid       = data.alicloud_account.default.id
  vpc_firewall_id  = alicloud_cen_instance.default.id
}

data "alicloud_cloud_firewall_vpc_firewall_control_policies" "default" {
	vpc_firewall_id = alicloud_cen_instance.default.id
	lang =             "zh"
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
