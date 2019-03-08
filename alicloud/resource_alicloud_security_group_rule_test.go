package alicloud

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudSecurityGroupRule_Ingress(t *testing.T) {
	var pt ecs.Permission

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, connectivity.EcsClassicSupportedRegions)
		},

		// module name
		IDRefreshName: "alicloud_security_group_rule.ingress",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSecurityGroupRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityGroupRuleIngress,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupRuleExists("alicloud_security_group_rule.ingress", &pt),
					resource.TestCheckResourceAttr("alicloud_security_group_rule.ingress", "type", "ingress"),
					resource.TestCheckResourceAttr("alicloud_security_group_rule.ingress", "ip_protocol", "tcp"),
					resource.TestCheckResourceAttrSet("alicloud_security_group_rule.ingress", "security_group_id"),
					resource.TestCheckResourceAttr("alicloud_security_group_rule.ingress", "port_range", "22/22"),
					resource.TestCheckResourceAttr("alicloud_security_group_rule.ingress", "cidr_ip", "10.159.6.18/12"),
					resource.TestCheckResourceAttr("alicloud_security_group_rule.ingress", "priority", "1"),
					resource.TestCheckResourceAttr("alicloud_security_group_rule.ingress", "policy", "accept"),
					resource.TestCheckResourceAttr("alicloud_security_group_rule.ingress", "nic_type", "intranet"),
				),
			},
		},
	})

}

func TestAccAlicloudSecurityGroupRule_Egress(t *testing.T) {
	var pt ecs.Permission

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, connectivity.EcsClassicSupportedRegions)
		},

		// module name
		IDRefreshName: "alicloud_security_group_rule.egress",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSecurityGroupRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityGroupRuleEgress,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupRuleExists("alicloud_security_group_rule.egress", &pt),
					resource.TestCheckResourceAttr("alicloud_security_group_rule.egress", "type", "egress"),
					resource.TestCheckResourceAttr("alicloud_security_group_rule.egress", "ip_protocol", "udp"),
					resource.TestCheckResourceAttrSet("alicloud_security_group_rule.egress", "security_group_id"),
					resource.TestCheckResourceAttr("alicloud_security_group_rule.egress", "port_range", "22/22"),
					resource.TestCheckResourceAttr("alicloud_security_group_rule.egress", "cidr_ip", "10.159.6.18/12"),
					resource.TestCheckResourceAttr("alicloud_security_group_rule.egress", "nic_type", "intranet"),
					resource.TestCheckResourceAttr("alicloud_security_group_rule.egress", "priority", "1"),
					resource.TestCheckResourceAttr("alicloud_security_group_rule.egress", "policy", "accept"),
				),
			},
		},
	})
}

func TestAccAlicloudSecurityGroupRule_EgressWithPolicy(t *testing.T) {
	var pt ecs.Permission

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_security_group_rule.egress",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSecurityGroupRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityGroupRuleEgress_WithPolicy,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupRuleExists("alicloud_security_group_rule.egress", &pt),
					resource.TestCheckResourceAttr("alicloud_security_group_rule.egress", "type", "egress"),
					resource.TestCheckResourceAttr("alicloud_security_group_rule.egress", "ip_protocol", "udp"),
					resource.TestCheckResourceAttr("alicloud_security_group_rule.egress", "policy", "accept"),
					resource.TestCheckResourceAttrSet("alicloud_security_group_rule.egress", "security_group_id"),
					resource.TestCheckResourceAttr("alicloud_security_group_rule.egress", "port_range", "22/22"),
					resource.TestCheckResourceAttr("alicloud_security_group_rule.egress", "cidr_ip", "10.159.6.18/12"),
					resource.TestCheckResourceAttr("alicloud_security_group_rule.egress", "nic_type", "intranet"),
					resource.TestCheckResourceAttr("alicloud_security_group_rule.egress", "priority", "1"),
				),
			},
		},
	})

}

func TestAccAlicloudSecurityGroupRule_IngressWithNic_Type(t *testing.T) {
	var pt ecs.Permission

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_security_group_rule.ingress",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSecurityGroupRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityGroupRuleIngress_Withnic_type,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupRuleExists("alicloud_security_group_rule.ingress", &pt),
					resource.TestCheckResourceAttr("alicloud_security_group_rule.ingress", "type", "ingress"),
					resource.TestCheckResourceAttr("alicloud_security_group_rule.ingress", "ip_protocol", "udp"),
					resource.TestCheckResourceAttr("alicloud_security_group_rule.ingress", "nic_type", "intranet"),
					resource.TestCheckResourceAttrSet("alicloud_security_group_rule.ingress", "security_group_id"),
					resource.TestCheckResourceAttr("alicloud_security_group_rule.ingress", "port_range", "22/22"),
					resource.TestCheckResourceAttr("alicloud_security_group_rule.ingress", "cidr_ip", "10.159.6.18/12"),
					resource.TestCheckResourceAttr("alicloud_security_group_rule.ingress", "policy", "accept"),
					resource.TestCheckResourceAttr("alicloud_security_group_rule.ingress", "priority", "1"),
				),
			},
		},
	})

}

func TestAccAlicloudSecurityGroupRule_EgressWithPriority(t *testing.T) {
	var pt ecs.Permission

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_security_group_rule.egress",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSecurityGroupRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityGroupRuleEgress_Withpriority,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupRuleExists("alicloud_security_group_rule.egress", &pt),
					resource.TestCheckResourceAttr("alicloud_security_group_rule.egress", "type", "egress"),
					resource.TestCheckResourceAttr("alicloud_security_group_rule.egress", "ip_protocol", "udp"),
					resource.TestCheckResourceAttr("alicloud_security_group_rule.egress", "priority", "1"),
					resource.TestCheckResourceAttrSet("alicloud_security_group_rule.egress", "security_group_id"),
					resource.TestCheckResourceAttr("alicloud_security_group_rule.egress", "port_range", "22/22"),
					resource.TestCheckResourceAttr("alicloud_security_group_rule.egress", "cidr_ip", "10.159.6.18/12"),
					resource.TestCheckResourceAttr("alicloud_security_group_rule.egress", "nic_type", "intranet"),
					resource.TestCheckResourceAttr("alicloud_security_group_rule.egress", "policy", "accept"),
				),
			},
		},
	})

}

func TestAccAlicloudSecurityGroupRule_EgressWithAll(t *testing.T) {
	var pt ecs.Permission

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_security_group_rule.egress",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSecurityGroupRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityGroupRuleEgress_WithAll,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupRuleExists("alicloud_security_group_rule.egress", &pt),
					resource.TestCheckResourceAttr("alicloud_security_group_rule.egress", "type", "egress"),
					resource.TestCheckResourceAttr("alicloud_security_group_rule.egress", "ip_protocol", "udp"),
					resource.TestCheckResourceAttr("alicloud_security_group_rule.egress", "priority", "1"),
					resource.TestCheckResourceAttrSet("alicloud_security_group_rule.egress", "security_group_id"),
					resource.TestCheckResourceAttr("alicloud_security_group_rule.egress", "nic_type", "intranet"),
					resource.TestCheckResourceAttr("alicloud_security_group_rule.egress", "port_range", "80/80"),
					resource.TestCheckResourceAttr("alicloud_security_group_rule.egress", "policy", "accept"),
					resource.TestCheckResourceAttr("alicloud_security_group_rule.egress", "cidr_ip", "10.159.6.18/12")),
			},
		},
	})

}

func TestAccAlicloudSecurityGroupRule_IngressWithAll(t *testing.T) {
	var pt ecs.Permission

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_security_group_rule.ingress",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSecurityGroupRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityGroupRuleIngress_WithAll,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupRuleExists("alicloud_security_group_rule.ingress", &pt),
					resource.TestCheckResourceAttr("alicloud_security_group_rule.ingress", "type", "ingress"),
					resource.TestCheckResourceAttr("alicloud_security_group_rule.ingress", "ip_protocol", "udp"),
					resource.TestCheckResourceAttr("alicloud_security_group_rule.ingress", "priority", "1"),
					resource.TestCheckResourceAttrSet("alicloud_security_group_rule.ingress", "security_group_id"),
					resource.TestCheckResourceAttr("alicloud_security_group_rule.ingress", "nic_type", "intranet"),
					resource.TestCheckResourceAttr("alicloud_security_group_rule.ingress", "port_range", "80/80"),
					resource.TestCheckResourceAttr("alicloud_security_group_rule.ingress", "policy", "accept"),
					resource.TestCheckResourceAttr("alicloud_security_group_rule.ingress", "cidr_ip", "10.159.6.18/12")),
			},
		},
	})

}

func TestAccAlicloudSecurityGroupRule_IngressWithSourceSecurityGroup(t *testing.T) {
	var pt ecs.Permission

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_security_group_rule.ingress",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSecurityGroupRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityGroupRuleSourceSecurityGroup,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupRuleExists("alicloud_security_group_rule.ingress", &pt),
					resource.TestMatchResourceAttr("alicloud_security_group_rule.ingress", "source_security_group_id", regexp.MustCompile("^sg-[a-zA-Z0-9_]+")),
					resource.TestCheckResourceAttr("alicloud_security_group_rule.ingress", "type", "ingress"),
					resource.TestCheckResourceAttr("alicloud_security_group_rule.ingress", "ip_protocol", "tcp"),
					resource.TestCheckResourceAttrSet("alicloud_security_group_rule.ingress", "security_group_id"),
					resource.TestCheckResourceAttr("alicloud_security_group_rule.ingress", "port_range", "22/22"),
					resource.TestCheckResourceAttr("alicloud_security_group_rule.ingress", "cidr_ip", ""),
					resource.TestCheckResourceAttr("alicloud_security_group_rule.ingress", "priority", "1"),
					resource.TestCheckResourceAttr("alicloud_security_group_rule.ingress", "policy", "accept"),
				),
			},
		},
	})

}

func TestAccAlicloudSecurityGroupRule_WithMulti(t *testing.T) {
	var pt ecs.Permission

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_security_group_rule.ingresses.0",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSecurityGroupRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityGroupRuleMulti,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupRuleExists("alicloud_security_group_rule.ingresses.0", &pt),
					resource.TestCheckResourceAttr("alicloud_security_group_rule.ingresses.0", "port_range", "1/200"),
					resource.TestCheckResourceAttr("alicloud_security_group_rule.ingresses.0", "ip_protocol", "udp"),
					resource.TestCheckResourceAttr("alicloud_security_group_rule.ingresses.0", "type", "ingress"),
					resource.TestCheckResourceAttr("alicloud_security_group_rule.ingresses.0", "priority", "1"),
					resource.TestCheckResourceAttrSet("alicloud_security_group_rule.ingresses.0", "security_group_id"),
					resource.TestCheckResourceAttr("alicloud_security_group_rule.ingresses.0", "nic_type", "intranet"),
					resource.TestCheckResourceAttr("alicloud_security_group_rule.ingresses.0", "policy", "accept"),
					resource.TestCheckResourceAttr("alicloud_security_group_rule.ingresses.0", "cidr_ip", "50.255.255.255/32")),
			},
			{
				Config: testAccSecurityGroupRuleMulti,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupRuleExists("alicloud_security_group_rule.egresses.0", &pt),
					resource.TestCheckResourceAttr("alicloud_security_group_rule.egresses.0", "port_range", "3306/3306"),
					resource.TestCheckResourceAttr("alicloud_security_group_rule.egresses.0", "ip_protocol", "tcp"),
					resource.TestCheckResourceAttr("alicloud_security_group_rule.egresses.0", "type", "egress"),
					resource.TestCheckResourceAttr("alicloud_security_group_rule.egresses.0", "priority", "1"),
					resource.TestCheckResourceAttrSet("alicloud_security_group_rule.egresses.0", "security_group_id"),
					resource.TestCheckResourceAttr("alicloud_security_group_rule.egresses.0", "nic_type", "intranet"),
					resource.TestCheckResourceAttr("alicloud_security_group_rule.egresses.0", "policy", "accept"),
					resource.TestCheckResourceAttr("alicloud_security_group_rule.egresses.0", "cidr_ip", "10.159.6.18/12")),
			},
		},
	})

}

func TestAccAlicloudSecurityGroupRule_MultiAttri(t *testing.T) {
	var pt ecs.Permission

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_security_group_rule.ingress_allow_tcp_22_sg",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSecurityGroupRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityGroupRuleMultiAttri,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupRuleExists("alicloud_security_group_rule.all", &pt),
					testAccCheckSecurityGroupRuleExists("alicloud_security_group_rule.gre", &pt),
					testAccCheckSecurityGroupRuleExists("alicloud_security_group_rule.ingress_allow_tcp_22_sg", &pt),
					testAccCheckSecurityGroupRuleExists("alicloud_security_group_rule.ingress_allow_tcp_22", &pt),
					testAccCheckSecurityGroupRuleExists("alicloud_security_group_rule.ingress_deny_tcp_22", &pt),
					testAccCheckSecurityGroupRuleExists("alicloud_security_group_rule.ingress_allow_tcp_22_prior", &pt),
					testAccCheckSecurityGroupRuleExists("alicloud_security_group_rule.ingress_deny_tcp_22_prior", &pt),
					resource.TestCheckResourceAttr("alicloud_security_group_rule.ingress_allow_tcp_22_sg", "port_range", "22/22"),
					resource.TestCheckResourceAttr("alicloud_security_group_rule.ingress_deny_tcp_22_prior", "port_range", "22/22"),
					resource.TestCheckResourceAttr("alicloud_security_group_rule.all", "port_range", "-1/-1"),
					resource.TestCheckResourceAttr("alicloud_security_group_rule.gre", "port_range", "-1/-1"),
				),
			},
		},
	})

}

func testAccCheckSecurityGroupRuleExists(n string, m *ecs.Permission) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No SecurityGroup Rule ID is set")
		}

		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		ecsService := EcsService{client}
		log.Printf("[WARN]get sg rule %s", rs.Primary.ID)
		parts := strings.Split(rs.Primary.ID, ":")
		prior, err := strconv.Atoi(parts[7])
		if err != nil {
			return fmt.Errorf("testSecrityGroupRuleExists parse rule id gets an error: %#v", err)
		}
		rule, err := ecsService.DescribeSecurityGroupRule(parts[0], parts[1], parts[2], parts[3], parts[4], parts[5], parts[6], prior)

		if err != nil {
			return err
		}

		*m = rule
		return nil
	}
}

func testAccCheckSecurityGroupRuleDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	ecsService := EcsService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_security_group_rule" {
			continue
		}

		parts := strings.Split(rs.Primary.ID, ":")
		prior, err := strconv.Atoi(parts[7])
		if err != nil {
			return fmt.Errorf("testSecrityGroupRuleDestroy parse rule id gets an error: %#v", err)
		}
		_, err = ecsService.DescribeSecurityGroupRule(parts[0], parts[1], parts[2], parts[3], parts[4], parts[5], parts[6], prior)

		// Verify the error is what we want
		if err != nil && !IsExceptedErrors(err, []string{InvalidSecurityGroupIdNotFound}) {
			return err
		}
	}

	return nil
}

const testAccSecurityGroupRuleIngress = `
variable "name" {
  default = "tf-testAccSecurityGroupRuleEgress_emptyNicType"
}

resource "alicloud_vpc" "vpc" {
  name = "${var.name}"
  cidr_block = "172.16.0.0/24"
}

resource "alicloud_security_group" "foo" {
  name = "${var.name}"
  vpc_id = "${alicloud_vpc.vpc.id}"
}

resource "alicloud_security_group_rule" "ingress" {
  type = "ingress"
  ip_protocol = "tcp"
  port_range = "22/22"
  cidr_ip = "10.159.6.18/12"
  security_group_id = "${alicloud_security_group.foo.id}"
}
`

const testAccSecurityGroupRuleEgress = `
variable "name" {
  default = "tf-testAccSecurityGroupRuleEgress_emptyNicType"
}

resource "alicloud_vpc" "vpc" {
  name = "${var.name}"
  cidr_block = "172.16.0.0/24"
}

resource "alicloud_security_group" "foo" {
  name = "${var.name}"
  vpc_id = "${alicloud_vpc.vpc.id}"
}

resource "alicloud_security_group_rule" "egress" {
  type = "egress"
  ip_protocol = "udp"
  port_range = "22/22"
  cidr_ip = "10.159.6.18/12"
  security_group_id = "${alicloud_security_group.foo.id}"
}
`

const testAccSecurityGroupRuleEgress_WithPolicy = `
variable "name" {
  default = "tf-testAccSecurityGroupRuleEgress_emptyNicType"
}

resource "alicloud_vpc" "vpc" {
  name = "${var.name}"
  cidr_block = "172.16.0.0/24"
}
resource "alicloud_security_group" "foo" {
  name = "${var.name}"
  vpc_id = "${alicloud_vpc.vpc.id}"
}

resource "alicloud_security_group_rule" "egress" {
  type = "egress"
  ip_protocol = "udp"
  policy = "accept"
  port_range = "22/22"
  cidr_ip = "10.159.6.18/12"
  security_group_id = "${alicloud_security_group.foo.id}"
  
}
`

const testAccSecurityGroupRuleIngress_Withnic_type = `
variable "name" {
  default = "tf-testAccSecurityGroupRuleVpcIngress"
}
resource "alicloud_security_group" "foo" {
  vpc_id = "${alicloud_vpc.vpc.id}"
  name = "${var.name}"
}

resource "alicloud_vpc" "vpc" {
  name = "${var.name}"
  cidr_block = "10.1.0.0/21"
}

resource "alicloud_security_group_rule" "ingress" {
  type = "ingress"
  ip_protocol = "udp"
  nic_type = "intranet"
  port_range = "22/22"
  cidr_ip = "10.159.6.18/12"
  security_group_id = "${alicloud_security_group.foo.id}"
}
`

const testAccSecurityGroupRuleEgress_Withpriority = `
variable "name" {
  default = "tf-testAccSecurityGroupRuleVpcIngress"
}
resource "alicloud_security_group" "foo" {
  vpc_id = "${alicloud_vpc.vpc.id}"
  name = "${var.name}"
}

resource "alicloud_vpc" "vpc" {
  name = "${var.name}"
  cidr_block = "10.1.0.0/21"
}

resource "alicloud_security_group_rule" "egress" {
  type = "egress"
  ip_protocol = "udp"
  priority = 1
  port_range = "22/22"
  cidr_ip = "10.159.6.18/12"
  security_group_id = "${alicloud_security_group.foo.id}"
}
`

const testAccSecurityGroupRuleIngress_WithAll = `
variable "name" {
  default = "tf-testAccSecurityGroupRule_missingSourceCidrIp"
}
resource "alicloud_security_group" "foo" {
  name = "${var.name}"
  vpc_id = "${alicloud_vpc.vpc.id}"
}

resource "alicloud_vpc" "vpc" {
  name = "${var.name}"
  cidr_block = "172.16.0.0/24"
}

resource "alicloud_security_group_rule" "ingress" {
  security_group_id = "${alicloud_security_group.foo.id}"
  type = "ingress"
  cidr_ip= "10.159.6.18/12"
  policy = "accept"
  ip_protocol= "udp"
  port_range= "80/80"
  priority= 1
  nic_type = "intranet"
}
`

const testAccSecurityGroupRuleEgress_WithAll = `
variable "name" {
  default = "tf-testAccSecurityGroupRule_missingSourceCidrIp"
}
resource "alicloud_security_group" "foo" {
  name = "${var.name}"
  vpc_id = "${alicloud_vpc.vpc.id}"
}

resource "alicloud_vpc" "vpc" {
  name = "${var.name}"
  cidr_block = "172.16.0.0/24"
}

resource "alicloud_security_group_rule" "egress" {
  security_group_id = "${alicloud_security_group.foo.id}"
  type = "egress"
  cidr_ip= "10.159.6.18/12"
  policy = "accept"
  ip_protocol= "udp"
  port_range= "80/80"
  priority= 1
  nic_type = "intranet"
}
`

const testAccSecurityGroupRuleMulti = `
variable "name" {
  default = "tf-testAccSecurityGroupRuleMulti"
}
variable "cidr_ip_list" {
  type = "list"
  default = ["50.255.255.255/32", "75.250.250.250/32", "45.20.250.240/32"]
}
variable "cidr_ip_list_2" {
  type = "list"
  default = ["10.159.6.18/12", "127.0.1.18/16"]
}
resource "alicloud_vpc" "main" {
  name = "${var.name}"
  cidr_block = "10.1.0.0/21"
}

resource "alicloud_security_group" "foo" {
  name = "${var.name}"
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
variable "name" {
  default = "tf-testAccSecurityGroupRuleSourceSecurityGroup"
}
resource "alicloud_vpc" "vpc" {
  name = "${var.name}"
  cidr_block = "172.16.0.0/24"
}
resource "alicloud_security_group" "foo" {
  name = "${var.name}-foo"
  vpc_id = "${alicloud_vpc.vpc.id}"
}

resource "alicloud_security_group" "bar" {
  name = "${var.name}_bar"
  vpc_id = "${alicloud_vpc.vpc.id}"
}

resource "alicloud_security_group_rule" "ingress" {
  type = "ingress"
  ip_protocol = "tcp"
  port_range  = "22/22"
  security_group_id = "${alicloud_security_group.bar.id}"
  source_security_group_id = "${alicloud_security_group.foo.id}"
}
`

const testAccSecurityGroupRuleMultiAttri = `
variable "name" {
  default = "tf-testAccSecurityGroupRuleMultiAttri"
}

variable "source_cidr_blocks" {
  type = "list"
  default = ["0.0.0.0/0"]
}


resource "alicloud_vpc" "vpc" {
  name = "${var.name}"
  cidr_block = "172.16.0.0/24"
}

resource "alicloud_security_group" "main" {
  name = "${var.name}"
  vpc_id = "${alicloud_vpc.vpc.id}"
}
resource "alicloud_security_group" "source" {
  name = "${var.name}-2"
  vpc_id = "${alicloud_vpc.vpc.id}"
}

resource "alicloud_security_group_rule" "ingress_allow_tcp_22_sg" {
  type              = "ingress"
  ip_protocol       = "tcp"
  nic_type          = "intranet"
  policy            = "accept"
  port_range        = "22/22"
  priority          = 1
  security_group_id = "${alicloud_security_group.main.id}"
  source_security_group_id = "${alicloud_security_group.source.id}"
}

resource "alicloud_security_group_rule" "ingress_allow_tcp_22" {
  count             = "${length(var.source_cidr_blocks)}"
  type              = "ingress"
  ip_protocol       = "tcp"
  nic_type          = "intranet"
  policy            = "accept"
  port_range        = "22/22"
  priority          = 1
  security_group_id = "${alicloud_security_group.main.id}"
  cidr_ip           = "${element(var.source_cidr_blocks, count.index)}"
}

resource "alicloud_security_group_rule" "ingress_deny_tcp_22" {
  type = "ingress"
  ip_protocol = "tcp"
  nic_type = "intranet"
  policy = "drop"
  port_range = "22/22"
  priority = 1
  security_group_id = "${alicloud_security_group.main.id}"
  cidr_ip = "0.0.0.0/0"
}

resource "alicloud_security_group_rule" "ingress_allow_tcp_22_prior" {
  type = "ingress"
  ip_protocol = "tcp"
  nic_type = "intranet"
  policy = "accept"
  port_range = "22/22"
  priority = 100
  security_group_id = "${alicloud_security_group.main.id}"
  cidr_ip = "0.0.0.0/0"
}
resource "alicloud_security_group_rule" "ingress_deny_tcp_22_prior" {
  type = "ingress"
  ip_protocol = "tcp"
  nic_type = "intranet"
  policy = "drop"
  port_range = "22/22"
  priority = 100
  security_group_id = "${alicloud_security_group.main.id}"
  cidr_ip = "0.0.0.0/0"
}
resource "alicloud_security_group_rule" "all" {
  type = "ingress"
  ip_protocol = "all"
  nic_type = "intranet"
  policy = "accept"
  priority = 100
  security_group_id = "${alicloud_security_group.main.id}"
  cidr_ip = "0.0.0.0/0"
}
resource "alicloud_security_group_rule" "gre" {
  type = "ingress"
  ip_protocol = "gre"
  nic_type = "intranet"
  policy = "accept"
  priority = 100
  security_group_id = "${alicloud_security_group.main.id}"
  cidr_ip = "0.0.0.0/0"
}
`
