package alicloud

import (
	"fmt"
	"github.com/denverdino/aliyungo/common"
	"github.com/denverdino/aliyungo/ecs"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"log"
	"regexp"
	"strings"
	"testing"
)

func TestAccAlicloudSecurityGroupRule_Ingress(t *testing.T) {
	var pt ecs.PermissionType

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_security_group_rule.ingress",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSecurityGroupRuleDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccSecurityGroupRuleIngress,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupRuleExists(
						"alicloud_security_group_rule.ingress", &pt),
					resource.TestCheckResourceAttr(
						"alicloud_security_group_rule.ingress",
						"priority",
						"1"),
					resource.TestCheckResourceAttr(
						"alicloud_security_group_rule.ingress",
						"nic_type",
						"internet"),
					resource.TestCheckResourceAttr(
						"alicloud_security_group_rule.ingress",
						"ip_protocol",
						"tcp"),
				),
			},
		},
	})

}

func TestAccAlicloudSecurityGroupRule_Egress(t *testing.T) {
	var pt ecs.PermissionType

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_security_group_rule.egress",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSecurityGroupRuleDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccSecurityGroupRuleEgress,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupRuleExists(
						"alicloud_security_group_rule.egress", &pt),
					resource.TestCheckResourceAttr(
						"alicloud_security_group_rule.egress",
						"port_range",
						"80/80"),
					resource.TestCheckResourceAttr(
						"alicloud_security_group_rule.egress",
						"ip_protocol",
						"udp"),
				),
			},
		},
	})

}

func TestAccAlicloudSecurityGroupRule_EgressDefaultNicType(t *testing.T) {
	var pt ecs.PermissionType

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_security_group_rule.egress",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSecurityGroupRuleDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccSecurityGroupRuleEgress_emptyNicType,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupRuleExists(
						"alicloud_security_group_rule.egress", &pt),
					resource.TestCheckResourceAttr(
						"alicloud_security_group_rule.egress",
						"port_range",
						"80/80"),
					resource.TestCheckResourceAttr(
						"alicloud_security_group_rule.egress",
						"nic_type",
						"internet"),
				),
			},
		},
	})

}

func TestAccAlicloudSecurityGroupRule_Vpc_Ingress(t *testing.T) {
	var pt ecs.PermissionType

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_security_group_rule.ingress",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSecurityGroupRuleDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccSecurityGroupRuleVpcIngress,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupRuleExists(
						"alicloud_security_group_rule.ingress", &pt),
					resource.TestCheckResourceAttr(
						"alicloud_security_group_rule.ingress",
						"port_range",
						"1/200"),
					resource.TestCheckResourceAttr(
						"alicloud_security_group_rule.ingress",
						"ip_protocol",
						"udp"),
				),
			},
		},
	})

}

func TestAccAlicloudSecurityGroupRule_MissParameterSourceCidrIp(t *testing.T) {
	var pt ecs.PermissionType

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_security_group_rule.egress",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSecurityGroupRuleDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccSecurityGroupRule_missingSourceCidrIp,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupRuleExists(
						"alicloud_security_group_rule.egress", &pt),
					resource.TestCheckResourceAttr(
						"alicloud_security_group_rule.egress",
						"port_range",
						"80/80"),
					resource.TestCheckResourceAttr(
						"alicloud_security_group_rule.egress",
						"nic_type",
						"internet"),
					resource.TestCheckResourceAttr(
						"alicloud_security_group_rule.egress",
						"ip_protocol",
						"udp"),
				),
			},
		},
	})

}

func TestAccAlicloudSecurityGroupRule_SourceSecurityGroup(t *testing.T) {
	var pt ecs.PermissionType

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_security_group_rule.ingress",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSecurityGroupRuleDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccSecurityGroupRuleSourceSecurityGroup,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupRuleExists(
						"alicloud_security_group_rule.ingress", &pt),
					resource.TestCheckResourceAttr(
						"alicloud_security_group_rule.ingress",
						"port_range",
						"3306/3306"),
					resource.TestMatchResourceAttr(
						"alicloud_security_group_rule.ingress",
						"source_security_group_id",
						regexp.MustCompile("^sg-[a-zA-Z0-9_]+")),
					resource.TestCheckResourceAttr(
						"alicloud_security_group_rule.ingress",
						"cidr_ip",
						""),
				),
			},
		},
	})

}

func TestAccAlicloudSecurityGroupRule_Multi(t *testing.T) {
	var pt ecs.PermissionType

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_security_group_rule.ingresses.0",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSecurityGroupRuleDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccSecurityGroupRuleMulti,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupRuleExists(
						"alicloud_security_group_rule.ingresses.0", &pt),
					resource.TestCheckResourceAttr(
						"alicloud_security_group_rule.ingresses.0",
						"port_range",
						"1/200"),
					resource.TestCheckResourceAttr(
						"alicloud_security_group_rule.ingresses.0",
						"ip_protocol",
						"udp"),
				),
			},
			resource.TestStep{
				Config: testAccSecurityGroupRuleMulti,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupRuleExists(
						"alicloud_security_group_rule.egresses.0", &pt),
					resource.TestCheckResourceAttr(
						"alicloud_security_group_rule.egresses.0",
						"port_range",
						"3306/3306"),
					resource.TestCheckResourceAttr(
						"alicloud_security_group_rule.egresses.0",
						"ip_protocol",
						"tcp"),
				),
			},
		},
	})

}

func testAccCheckSecurityGroupRuleExists(n string, m *ecs.PermissionType) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No SecurityGroup Rule ID is set")
		}

		client := testAccProvider.Meta().(*AliyunClient)
		log.Printf("[WARN]get sg rule %s", rs.Primary.ID)
		parts := strings.Split(rs.Primary.ID, ":")
		// securityGroupId, direction, nicType, ipProtocol, portRange
		rule, err := client.DescribeSecurityGroupRule(parts[0], parts[1], parts[4], parts[2], parts[3])

		if err != nil {
			return err
		}

		if rule == nil {
			return fmt.Errorf("SecurityGroup not found")
		}

		*m = *rule
		return nil
	}
}

func testAccCheckSecurityGroupRuleDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*AliyunClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_security_group_rule" {
			continue
		}

		parts := strings.Split(rs.Primary.ID, ":")
		rule, err := client.DescribeSecurityGroupRule(parts[0], parts[1], parts[4], parts[2], parts[3])

		if rule != nil {
			return fmt.Errorf("Error SecurityGroup Rule still exist")
		}

		// Verify the error is what we want
		if err != nil {
			// Verify the error is what we want
			e, _ := err.(*common.Error)
			if e.ErrorResponse.Code == InvalidSecurityGroupIdNotFound {
				continue
			}
			return err
		}
	}

	return nil
}

const testAccSecurityGroupRuleIngress = `
resource "alicloud_security_group" "foo" {
  name = "sg_foo"
}

resource "alicloud_security_group_rule" "ingress" {
  type = "ingress"
  ip_protocol = "tcp"
  nic_type = "internet"
  policy = "accept"
  port_range = "1/200"
  priority = 1
  security_group_id = "${alicloud_security_group.foo.id}"
  cidr_ip = "10.159.6.18/12"
}


`

const testAccSecurityGroupRuleEgress = `
resource "alicloud_security_group" "foo" {
  name = "sg_foo"
}


resource "alicloud_security_group_rule" "egress" {
  type = "egress"
  ip_protocol = "udp"
  nic_type = "internet"
  policy = "accept"
  port_range = "80/80"
  priority = 1
  security_group_id = "${alicloud_security_group.foo.id}"
  cidr_ip = "10.159.6.18/12"
}

`

const testAccSecurityGroupRuleEgress_emptyNicType = `
resource "alicloud_security_group" "foo" {
  name = "sg_foo"
}

resource "alicloud_security_group_rule" "egress" {
  type = "egress"
  ip_protocol = "udp"
  policy = "accept"
  port_range = "80/80"
  priority = 1
  security_group_id = "${alicloud_security_group.foo.id}"
  cidr_ip = "10.159.6.18/12"
}

`

const testAccSecurityGroupRuleVpcIngress = `
resource "alicloud_security_group" "foo" {
  vpc_id = "${alicloud_vpc.vpc.id}"
  name = "sg_foo"
}

resource "alicloud_vpc" "vpc" {
  cidr_block = "10.1.0.0/21"
}

resource "alicloud_security_group_rule" "ingress" {
  type = "ingress"
  ip_protocol = "udp"
  nic_type = "intranet"
  policy = "accept"
  port_range = "1/200"
  priority = 1
  security_group_id = "${alicloud_security_group.foo.id}"
  cidr_ip = "10.159.6.18/12"
}

`
const testAccSecurityGroupRule_missingSourceCidrIp = `
resource "alicloud_security_group" "foo" {
  name = "sg_foo"
}

resource "alicloud_security_group_rule" "egress" {
  security_group_id = "${alicloud_security_group.foo.id}"
  type = "egress"
  cidr_ip= "0.0.0.0/0"
  policy = "accept"
  ip_protocol= "udp"
  port_range= "80/80"
  priority= 1
}

`

const testAccSecurityGroupRuleMulti = `
variable "cidr_ip_list" {
  type = "list"
  default = ["50.255.255.255/32", "75.250.250.250/32", "45.20.250.240/32"]
}
variable "cidr_ip_list_2" {
  type = "list"
  default = ["10.159.6.18/12", "127.0.1.18/16"]
}
resource "alicloud_vpc" "main" {
  cidr_block = "10.1.0.0/21"
}

resource "alicloud_vswitch" "main" {
  vpc_id = "${alicloud_vpc.main.id}"
  cidr_block = "10.1.1.0/24"
  availability_zone = "cn-beijing-a"
  depends_on = [
    "alicloud_vpc.main"]
}

resource "alicloud_security_group" "foo" {
  name = "test_rules"
  description = "Security group for rules"
  vpc_id = "${alicloud_vpc.main.id}"
}

resource "alicloud_security_group_rule" "ingresses" {
  count = "${length(compact(var.cidr_ip_list))}"
  security_group_id = "${alicloud_security_group.foo.id}"

  type = "ingress"
  policy = "accept"
  port_range = "1/200"
  ip_protocol = "udp"
  nic_type = "intranet"
  priority = 1

  cidr_ip = "${element(var.cidr_ip_list, count.index)}"
}

resource "alicloud_security_group_rule" "egresses" {
  count = "${length(compact(var.cidr_ip_list_2))}"
  type = "egress"
  ip_protocol = "tcp"
  nic_type = "intranet"
  policy = "accept"
  port_range = "3306/3306"
  priority = 1
  security_group_id = "${alicloud_security_group.foo.id}"
  cidr_ip = "${element(var.cidr_ip_list_2, count.index)}"
}
`

const testAccSecurityGroupRuleSourceSecurityGroup = `
resource "alicloud_security_group" "foo" {
  name = "sg_foo"
}

resource "alicloud_security_group" "bar" {
  name = "sg_bar"
}

resource "alicloud_security_group_rule" "ingress" {
  type = "ingress"
  ip_protocol = "tcp"
  nic_type = "intranet"
  policy = "accept"
  port_range = "3306/3306"
  priority = 50
  security_group_id = "${alicloud_security_group.bar.id}"
  source_security_group_id = "${alicloud_security_group.foo.id}"
}
`
