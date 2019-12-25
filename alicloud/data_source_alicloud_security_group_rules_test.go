package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudSecurityGroupRulesDataSourceWithDirection(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudSecurityGroupRulesDataSourceConfigDirection,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_security_group_rules.ingress"),
					resource.TestCheckResourceAttr("data.alicloud_security_group_rules.ingress", "group_name", "tf-testAccCheckAlicloudSecurityGroupRulesDataSourceConfig_1"),
					resource.TestCheckResourceAttr("data.alicloud_security_group_rules.ingress", "group_desc", "alicloud security group"),
					resource.TestCheckResourceAttr("data.alicloud_security_group_rules.ingress", "rules.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_security_group_rules.ingress", "rules.0.direction", "ingress"),
					resource.TestCheckResourceAttr("data.alicloud_security_group_rules.ingress", "rules.0.ip_protocol", "tcp"),
					resource.TestCheckResourceAttr("data.alicloud_security_group_rules.ingress", "rules.0.port_range", "5000/5001"),
					resource.TestCheckNoResourceAttr("data.alicloud_security_group_rules.ingress", "rules.0.cidr_ip"),
					resource.TestCheckResourceAttrSet("data.alicloud_security_group_rules.ingress", "rules.0.source_cidr_ip"),
					resource.TestCheckResourceAttr("data.alicloud_security_group_rules.ingress", "rules.0.source_group_id", ""),
					resource.TestCheckResourceAttr("data.alicloud_security_group_rules.ingress", "rules.0.source_group_owner_account", ""),
					resource.TestCheckResourceAttr("data.alicloud_security_group_rules.ingress", "rules.0.dest_cidr_ip", ""),
					resource.TestCheckResourceAttr("data.alicloud_security_group_rules.ingress", "rules.0.dest_group_id", ""),
					resource.TestCheckResourceAttr("data.alicloud_security_group_rules.ingress", "rules.0.dest_group_owner_account", ""),
					resource.TestCheckResourceAttrSet("data.alicloud_security_group_rules.ingress", "rules.0.priority"),
					resource.TestCheckResourceAttr("data.alicloud_security_group_rules.ingress", "rules.0.description", ""),
					resource.TestCheckResourceAttrSet("data.alicloud_security_group_rules.ingress", "rules.0.nic_type"),
				),
			},
		},
	})
}

func TestAccAlicloudSecurityGroupRulesDataSourceWithGroupId(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudSecurityGroupRulesDataSourceConfigGroup_id,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_security_group_rules.egress"),
					resource.TestCheckResourceAttr("data.alicloud_security_group_rules.egress", "group_name", "tf-testAccCheckAlicloudSecurityGroupRulesDataSourceConfig0"),
					resource.TestCheckResourceAttr("data.alicloud_security_group_rules.egress", "group_desc", "alicloud security group"),
					resource.TestCheckResourceAttr("data.alicloud_security_group_rules.egress", "rules.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_security_group_rules.egress", "rules.0.direction", "egress"),
					resource.TestCheckResourceAttr("data.alicloud_security_group_rules.egress", "rules.0.ip_protocol", "udp"),
					resource.TestCheckResourceAttr("data.alicloud_security_group_rules.egress", "rules.0.nic_type", "intranet"),
					resource.TestCheckNoResourceAttr("data.alicloud_security_group_rules.egress", "rules.0.cidr_ip"),
					resource.TestCheckResourceAttr("data.alicloud_security_group_rules.egress", "rules.0.port_range", "6000/6001"),
					resource.TestCheckResourceAttr("data.alicloud_security_group_rules.egress", "rules.0.source_cidr_ip", ""),
					resource.TestCheckResourceAttr("data.alicloud_security_group_rules.egress", "rules.0.source_group_id", ""),
					resource.TestCheckResourceAttr("data.alicloud_security_group_rules.egress", "rules.0.source_group_owner_account", ""),
					resource.TestCheckResourceAttr("data.alicloud_security_group_rules.egress", "rules.0.dest_cidr_ip", "0.0.0.0/0"),
					resource.TestCheckResourceAttr("data.alicloud_security_group_rules.egress", "rules.0.dest_group_id", ""),
					resource.TestCheckResourceAttr("data.alicloud_security_group_rules.egress", "rules.0.dest_group_owner_account", ""),
					resource.TestCheckResourceAttrSet("data.alicloud_security_group_rules.egress", "rules.0.priority"),
					resource.TestCheckResourceAttr("data.alicloud_security_group_rules.egress", "rules.0.description", ""),
				),
			},
		},
	})
}

func TestAccAlicloudSecurityGroupRulesDataSourceWithNic_Type(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, connectivity.EcsClassicSupportedRegions)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudSecurityGroupRulesDataSourceConfigNic_Type,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_security_group_rules.egress"),
					resource.TestCheckResourceAttr("data.alicloud_security_group_rules.egress", "rules.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_security_group_rules.egress", "group_name", "tf-testAccCheckAlicloudSecurityGroupRulesDataSourceConfig1"),
					resource.TestCheckResourceAttr("data.alicloud_security_group_rules.egress", "group_desc", "alicloud security group"),
					resource.TestCheckResourceAttr("data.alicloud_security_group_rules.egress", "rules.0.direction", "egress"),
					resource.TestCheckResourceAttr("data.alicloud_security_group_rules.egress", "rules.0.ip_protocol", "udp"),
					resource.TestCheckResourceAttr("data.alicloud_security_group_rules.egress", "rules.0.nic_type", "intranet"),
					resource.TestCheckResourceAttr("data.alicloud_security_group_rules.egress", "rules.0.port_range", "6000/6001"),
					resource.TestCheckNoResourceAttr("data.alicloud_security_group_rules.egress", "rules.0.cidr_ip"),
					resource.TestCheckResourceAttr("data.alicloud_security_group_rules.egress", "rules.0.source_cidr_ip", ""),
					resource.TestCheckResourceAttr("data.alicloud_security_group_rules.egress", "rules.0.source_group_id", ""),
					resource.TestCheckResourceAttr("data.alicloud_security_group_rules.egress", "rules.0.source_group_owner_account", ""),
					resource.TestCheckResourceAttr("data.alicloud_security_group_rules.egress", "rules.0.dest_cidr_ip", "0.0.0.0/0"),
					resource.TestCheckResourceAttr("data.alicloud_security_group_rules.egress", "rules.0.dest_group_id", ""),
					resource.TestCheckResourceAttr("data.alicloud_security_group_rules.egress", "rules.0.dest_group_owner_account", ""),
					resource.TestCheckResourceAttrSet("data.alicloud_security_group_rules.egress", "rules.0.priority"),
					resource.TestCheckResourceAttr("data.alicloud_security_group_rules.egress", "rules.0.description", ""),
				),
			},
		},
	})
}

func TestAccAlicloudSecurityGroupRulesDataSourceWithPolicy(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudSecurityGroupRulesDataSourceConfigPolicy,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_security_group_rules.egress"),
					resource.TestCheckResourceAttr("data.alicloud_security_group_rules.egress", "group_name", "tf-testAccCheckAlicloudSecurityGroupRulesDataSourceConfig3"),
					resource.TestCheckResourceAttr("data.alicloud_security_group_rules.egress", "group_desc", "alicloud security group"),
					resource.TestCheckResourceAttr("data.alicloud_security_group_rules.egress", "rules.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_security_group_rules.egress", "rules.0.direction", "egress"),
					resource.TestCheckResourceAttr("data.alicloud_security_group_rules.egress", "rules.0.ip_protocol", "udp"),
					resource.TestCheckResourceAttr("data.alicloud_security_group_rules.egress", "rules.0.nic_type", "intranet"),
					resource.TestCheckResourceAttr("data.alicloud_security_group_rules.egress", "rules.0.port_range", "6000/6001"),
					resource.TestCheckNoResourceAttr("data.alicloud_security_group_rules.egress", "rules.0.cidr_ip"),
					resource.TestCheckResourceAttr("data.alicloud_security_group_rules.egress", "rules.0.policy", "drop"),
					resource.TestCheckResourceAttr("data.alicloud_security_group_rules.egress", "rules.0.source_cidr_ip", ""),
					resource.TestCheckResourceAttr("data.alicloud_security_group_rules.egress", "rules.0.source_group_id", ""),
					resource.TestCheckResourceAttr("data.alicloud_security_group_rules.egress", "rules.0.source_group_owner_account", ""),
					resource.TestCheckResourceAttr("data.alicloud_security_group_rules.egress", "rules.0.dest_cidr_ip", "0.0.0.0/0"),
					resource.TestCheckResourceAttr("data.alicloud_security_group_rules.egress", "rules.0.dest_group_id", ""),
					resource.TestCheckResourceAttr("data.alicloud_security_group_rules.egress", "rules.0.dest_group_owner_account", ""),
					resource.TestCheckResourceAttrSet("data.alicloud_security_group_rules.egress", "rules.0.priority"),
					resource.TestCheckResourceAttr("data.alicloud_security_group_rules.egress", "rules.0.description", ""),
				),
			},
		},
	})
}

func TestAccAlicloudSecurityGroupRulesDataSourceWithIp_Protocol(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudSecurityGroupRulesDataSourceConfigIp_Protocol,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_security_group_rules.egress"),
					resource.TestCheckResourceAttr("data.alicloud_security_group_rules.egress", "group_name", "tf-testAccCheckAlicloudSecurityGroupRulesDataSourceConfig2"),
					resource.TestCheckResourceAttr("data.alicloud_security_group_rules.egress", "group_desc", "alicloud security group"),
					resource.TestCheckResourceAttr("data.alicloud_security_group_rules.egress", "rules.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_security_group_rules.egress", "rules.0.direction", "egress"),
					resource.TestCheckResourceAttr("data.alicloud_security_group_rules.egress", "rules.0.ip_protocol", "udp"),
					resource.TestCheckResourceAttr("data.alicloud_security_group_rules.egress", "rules.0.nic_type", "intranet"),
					resource.TestCheckResourceAttr("data.alicloud_security_group_rules.egress", "rules.0.port_range", "6000/6001"),
					resource.TestCheckNoResourceAttr("data.alicloud_security_group_rules.egress", "rules.0.cidr_ip"),
					resource.TestCheckResourceAttr("data.alicloud_security_group_rules.egress", "rules.0.source_cidr_ip", ""),
					resource.TestCheckResourceAttr("data.alicloud_security_group_rules.egress", "rules.0.source_group_id", ""),
					resource.TestCheckResourceAttr("data.alicloud_security_group_rules.egress", "rules.0.source_group_owner_account", ""),
					resource.TestCheckResourceAttr("data.alicloud_security_group_rules.egress", "rules.0.dest_cidr_ip", "0.0.0.0/0"),
					resource.TestCheckResourceAttr("data.alicloud_security_group_rules.egress", "rules.0.dest_group_id", ""),
					resource.TestCheckResourceAttr("data.alicloud_security_group_rules.egress", "rules.0.dest_group_owner_account", ""),
					resource.TestCheckResourceAttrSet("data.alicloud_security_group_rules.egress", "rules.0.priority"),
					resource.TestCheckResourceAttr("data.alicloud_security_group_rules.egress", "rules.0.description", ""),
				),
			},
		},
	})
}

func TestAccAlicloudSecurityGroupRulesDataSourceEmpty(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudSecurityGroupRulesDataSourceConfigEmpty,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_security_group_rules.empty"),
					resource.TestCheckResourceAttr("data.alicloud_security_group_rules.empty", "group_name", "tf-testAccCheckAlicloudSecurityGroupRulesDataSourceConfigEgress"),
					resource.TestCheckResourceAttr("data.alicloud_security_group_rules.empty", "group_desc", "alicloud security group"),
					resource.TestCheckResourceAttr("data.alicloud_security_group_rules.empty", "rules.#", "0"),
					resource.TestCheckNoResourceAttr("data.alicloud_security_group_rules.empty", "rules.0.direction"),
					resource.TestCheckNoResourceAttr("data.alicloud_security_group_rules.empty", "rules.0.ip_protocol"),
					resource.TestCheckNoResourceAttr("data.alicloud_security_group_rules.empty", "rules.0.nic_type"),
					resource.TestCheckNoResourceAttr("data.alicloud_security_group_rules.empty", "rules.0.policy"),
					resource.TestCheckNoResourceAttr("data.alicloud_security_group_rules.empty", "rules.0.port_range"),
					resource.TestCheckNoResourceAttr("data.alicloud_security_group_rules.empty", "rules.0.priority"),
					resource.TestCheckNoResourceAttr("data.alicloud_security_group_rules.empty", "rules.0.source_cidr_ip"),
				),
			},
		},
	})
}

const testAccCheckAlicloudSecurityGroupRulesDataSourceConfigDirection = `
variable "name" {
	default = "tf-testAccCheckAlicloudSecurityGroupRulesDataSourceConfig_1"
}
resource "alicloud_vpc" "foo" {
  cidr_block = "172.16.0.0/12"
  name = var.name
}

resource "alicloud_security_group" "group" {
  name = var.name
  description = "alicloud security group"
  vpc_id      = alicloud_vpc.foo.id
}

resource "alicloud_security_group_rule" "rule_ingress" {
  type              = "ingress"
  ip_protocol       = "tcp"
  port_range        = "5000/5001"
  security_group_id = alicloud_security_group.group.id
  cidr_ip           = "0.0.0.0/0"
}

resource "alicloud_security_group_rule" "rule_egress" {
  type              = "egress"
  ip_protocol       = "udp"
  port_range        = "6000/6001"
  security_group_id = alicloud_security_group.group.id
  cidr_ip           = "0.0.0.0/0"
}

data "alicloud_security_group_rules" "ingress" {
  direction   = "ingress"
  group_id    = alicloud_security_group_rule.rule_ingress.security_group_id
}
`

const testAccCheckAlicloudSecurityGroupRulesDataSourceConfigGroup_id = `
variable "name" {
	default = "tf-testAccCheckAlicloudSecurityGroupRulesDataSourceConfig0"
}
resource "alicloud_vpc" "foo" {
  cidr_block = "172.16.0.0/12"
  name = var.name
}

resource "alicloud_security_group" "group" {
  name = var.name
  description = "alicloud security group"
  vpc_id      = alicloud_vpc.foo.id
}

resource "alicloud_security_group" "bar" {
  name = "tf-testAccCheckAlicloudSecurityGroupRules"
  description = "alicloud security group"
  vpc_id      = alicloud_vpc.foo.id
}

resource "alicloud_security_group_rule" "rule_ingress" {
  type              = "ingress"
  ip_protocol       = "tcp"
  port_range        = "5000/5001"
  security_group_id = alicloud_security_group.bar.id
  cidr_ip           = "0.0.0.0/0"
}

resource "alicloud_security_group_rule" "rule_egress" {
  type              = "egress"
  ip_protocol       = "udp"
  port_range        = "6000/6001"
  security_group_id = alicloud_security_group.group.id
  cidr_ip           = "0.0.0.0/0"
}

data "alicloud_security_group_rules" "egress" {
  group_id    = alicloud_security_group_rule.rule_egress.security_group_id
}
`

const testAccCheckAlicloudSecurityGroupRulesDataSourceConfigNic_Type = `
variable "name" {
	default = "tf-testAccCheckAlicloudSecurityGroupRulesDataSourceConfig1"
}

resource "alicloud_security_group" "group" {
  name = var.name
  description = "alicloud security group"
}

resource "alicloud_security_group_rule" "rule_ingress" {
  type              = "ingress"
  ip_protocol       = "tcp"
  port_range        = "5000/5001"
  security_group_id = alicloud_security_group.group.id
  cidr_ip           = "0.0.0.0/0"
  nic_type          = "internet"
}

resource "alicloud_security_group_rule" "rule_egress" {
  type              = "egress"
  ip_protocol       = "udp"
  port_range        = "6000/6001"
  security_group_id = alicloud_security_group.group.id
  cidr_ip           = "0.0.0.0/0"
  nic_type          = "intranet"
}

data "alicloud_security_group_rules" "egress" {
  nic_type   = "intranet"
  group_id    = alicloud_security_group_rule.rule_egress.security_group_id
}
`

const testAccCheckAlicloudSecurityGroupRulesDataSourceConfigIp_Protocol = `
variable "name" {
	default = "tf-testAccCheckAlicloudSecurityGroupRulesDataSourceConfig2"
}
resource "alicloud_vpc" "foo" {
  cidr_block = "172.16.0.0/12"
  name = var.name
}

resource "alicloud_security_group" "group" {
  name = var.name
  description = "alicloud security group"
  vpc_id      = alicloud_vpc.foo.id
}

resource "alicloud_security_group_rule" "rule_ingress" {
  type              = "ingress"
  ip_protocol       = "tcp"
  port_range        = "5000/5001"
  security_group_id = alicloud_security_group.group.id
  cidr_ip           = "0.0.0.0/0"
  nic_type          = "intranet"
}

resource "alicloud_security_group_rule" "rule_egress" {
  type              = "egress"
  ip_protocol       = "udp"
  port_range        = "6000/6001"
  security_group_id = alicloud_security_group.group.id
  cidr_ip           = "0.0.0.0/0"
  nic_type          = "intranet"
}

data "alicloud_security_group_rules" "egress" {
  ip_protocol   = "udp"
  group_id    = alicloud_security_group_rule.rule_egress.security_group_id
}
`

const testAccCheckAlicloudSecurityGroupRulesDataSourceConfigPolicy = `
variable "name" {
	default = "tf-testAccCheckAlicloudSecurityGroupRulesDataSourceConfig3"
}
resource "alicloud_vpc" "foo" {
  cidr_block = "172.16.0.0/12"
  name = var.name
}

resource "alicloud_security_group" "group" {
  name = var.name
  description = "alicloud security group"
  vpc_id      = alicloud_vpc.foo.id
}

resource "alicloud_security_group_rule" "rule_ingress" {
  type              = "ingress"
  ip_protocol       = "tcp"
  port_range        = "5000/5001"
  security_group_id = alicloud_security_group.group.id
  cidr_ip           = "0.0.0.0/0"
  nic_type          = "intranet"
}

resource "alicloud_security_group_rule" "rule_egress" {
  type              = "egress"
  ip_protocol       = "udp"
  port_range        = "6000/6001"
  policy            = "drop"
  security_group_id = alicloud_security_group.group.id
  cidr_ip           = "0.0.0.0/0"
  nic_type          = "intranet"
}

data "alicloud_security_group_rules" "egress" {
  policy   = "drop"
  group_id   =alicloud_security_group_rule.rule_egress.security_group_id
}
`

const testAccCheckAlicloudSecurityGroupRulesDataSourceConfigEmpty = `
variable "name" {
	default = "tf-testAccCheckAlicloudSecurityGroupRulesDataSourceConfigEgress"
}
resource "alicloud_vpc" "foo" {
  cidr_block = "172.16.0.0/12"
  name = var.name
}

resource "alicloud_security_group" "group" {
  name = var.name
  description = "alicloud security group"
  vpc_id      = alicloud_vpc.foo.id
}

data "alicloud_security_group_rules" "empty" {
  group_id    = alicloud_security_group.group.id
}
`
