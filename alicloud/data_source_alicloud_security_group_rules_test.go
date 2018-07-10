package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudSecurityGroupRulesDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudSecurityGroupRulesDataSourceConfigIngress,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_security_group_rules.ingress"),
					resource.TestCheckResourceAttr("data.alicloud_security_group_rules.ingress", "group_name", "testAccCheckAlicloudSecurityGroupRulesDataSourceConfigIngress"),
					resource.TestCheckResourceAttr("data.alicloud_security_group_rules.ingress", "group_desc", "alicloud security group"),
					resource.TestCheckResourceAttr("data.alicloud_security_group_rules.ingress", "rules.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_security_group_rules.ingress", "rules.0.direction", "ingress"),
					resource.TestCheckResourceAttr("data.alicloud_security_group_rules.ingress", "rules.0.ip_protocol", "tcp"),
					resource.TestCheckResourceAttr("data.alicloud_security_group_rules.ingress", "rules.0.nic_type", "intranet"),
					resource.TestCheckResourceAttr("data.alicloud_security_group_rules.ingress", "rules.0.policy", "accept"),
					resource.TestCheckResourceAttr("data.alicloud_security_group_rules.ingress", "rules.0.port_range", "5000/5001"),
					resource.TestCheckResourceAttr("data.alicloud_security_group_rules.ingress", "rules.0.priority", "1"),
					resource.TestCheckResourceAttr("data.alicloud_security_group_rules.ingress", "rules.0.source_cidr_ip", "0.0.0.0/0"),
				),
			},
			{
				Config: testAccCheckAlicloudSecurityGroupRulesDataSourceConfigEgress,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_security_group_rules.egress"),
					resource.TestCheckResourceAttr("data.alicloud_security_group_rules.egress", "rules.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_security_group_rules.egress", "rules.0.direction", "egress"),
					resource.TestCheckResourceAttr("data.alicloud_security_group_rules.egress", "rules.0.ip_protocol", "udp"),
					resource.TestCheckResourceAttr("data.alicloud_security_group_rules.egress", "rules.0.nic_type", "intranet"),
				),
			},
		},
	})
}

const testAccCheckAlicloudSecurityGroupRulesDataSourceConfigIngress = `
variable "name" {
	default = "testAccCheckAlicloudSecurityGroupRulesDataSourceConfigIngress"
}
resource "alicloud_vpc" "foo" {
  cidr_block = "172.16.0.0/12"
  name = "${var.name}"
}

resource "alicloud_security_group" "group" {
  name = "${var.name}"
  description = "alicloud security group"
  vpc_id      = "${alicloud_vpc.foo.id}"
}

resource "alicloud_security_group_rule" "rule_ingress" {
  type              = "ingress"
  ip_protocol       = "tcp"
  nic_type          = "intranet"
  policy            = "accept"
  port_range        = "5000/5001"
  priority          = 1
  security_group_id = "${alicloud_security_group.group.id}"
  cidr_ip           = "0.0.0.0/0"
}

resource "alicloud_security_group_rule" "rule_egress" {
  type              = "egress"
  ip_protocol       = "udp"
  port_range        = "6000/6001"
  security_group_id = "${alicloud_security_group.group.id}"
  cidr_ip           = "0.0.0.0/0"
}

data "alicloud_security_group_rules" "ingress" {
  group_id    = "${alicloud_security_group_rule.rule_ingress.security_group_id}"
  nic_type    = "intranet"
  direction   = "ingress"
  ip_protocol = "tcp"
  policy      = "accept"
}
`
const testAccCheckAlicloudSecurityGroupRulesDataSourceConfigEgress = `
variable "name" {
	default = "testAccCheckAlicloudSecurityGroupRulesDataSourceConfigEgress"
}
resource "alicloud_vpc" "foo" {
  cidr_block = "172.16.0.0/12"
  name = "${var.name}"
}

resource "alicloud_security_group" "group" {
  name = "${var.name}"
  description = "alicloud security group"
  vpc_id      = "${alicloud_vpc.foo.id}"
}

resource "alicloud_security_group_rule" "rule_ingress" {
  type              = "ingress"
  ip_protocol       = "tcp"
  port_range        = "5000/5001"
  security_group_id = "${alicloud_security_group.group.id}"
  cidr_ip           = "0.0.0.0/0"
}

resource "alicloud_security_group_rule" "rule_egress" {
  type              = "egress"
  ip_protocol       = "udp"
  port_range        = "6000/6001"
  security_group_id = "${alicloud_security_group.group.id}"
  cidr_ip           = "0.0.0.0/0"
}

data "alicloud_security_group_rules" "egress" {
  group_id    = "${alicloud_security_group_rule.rule_ingress.security_group_id}"
  nic_type    = "intranet"
  direction   = "egress"
  ip_protocol = "udp"
}
`
