package alicloud

import (
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAlicloudECSSecurityGroupRuleBasic(t *testing.T) {
	var v ecs.Permission
	resourceId := "alicloud_security_group_rule.default"
	ra := resourceAttrInit(resourceId, testAccCheckSecurityGroupRuleBasicMap)
	serviceFunc := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSecurityGroupRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityGroupRuleBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "abc",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccSecurityGroupRule_cidrIp,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source_security_group_id": REMOVEKEY,
						"cidr_ip":                  "0.0.0.0/0",
						"description":              "abcd",
					}),
				),
			},
			{
				Config: testAccSecurityGroupRule_description,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "description",
					}),
				),
			},
			{
				Config: testAccSecurityGroupRule_all,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "abcd",
					}),
				),
			},
		},
	})

}

func TestAccAlicloudECSSecurityGroupEgressRule(t *testing.T) {
	var v ecs.Permission
	resourceId := "alicloud_security_group_rule.default"
	ra := resourceAttrInit(resourceId, map[string]string{
		"type":        "egress",
		"policy":      "accept",
		"description": "SHDRP-7513",
		"port_range":  "443/443",
		"priority":    "1",
		"cidr_ip":     "182.254.11.243/32",
	})
	serviceFunc := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSecurityGroupRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityGroupEgressRule,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "SHDRP-7513",
					}),
				),
			},
			{
				Config: testAccSecurityGroupEgressRule_description,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "SHDRP-7512",
					}),
				),
			},
		},
	})

}

func TestAccAlicloudECSSecurityGroupRuleMulti(t *testing.T) {
	var v ecs.Permission
	resourceId := "alicloud_security_group_rule.default.2"
	ra := resourceAttrInit(resourceId, testAccCheckSecurityGroupRuleBasicMap)
	serviceFunc := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSecurityGroupRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityGroupRuleMulti,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cidr_ip":                  "45.20.250.240/32",
						"source_security_group_id": REMOVEKEY,
					}),
				),
			},
		},
	})

}

func TestAccAlicloudECSSecurityGroupRulePrefixList(t *testing.T) {
	var v ecs.Permission
	resourceId := "alicloud_security_group_rule.default"
	ra := resourceAttrInit(resourceId, testAccCheckSecurityGroupRulePrefixList)
	serviceFunc := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSecurityGroupRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityGroupRulePrefix,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "abc",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})

}

func TestAccAlicloudECSSecurityGroupEgressRuleIpv6(t *testing.T) {
	var v ecs.Permission
	resourceId := "alicloud_security_group_rule.default"
	ra := resourceAttrInit(resourceId, map[string]string{
		"type":         "egress",
		"policy":       "accept",
		"description":  "SHDRP-7513",
		"port_range":   "443/443",
		"priority":     "1",
		"ipv6_cidr_ip": "2408:4004:cc:400::/56",
	})
	serviceFunc := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSecurityGroupRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityGroupEgressRuleIpv6,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "SHDRP-7513",
					}),
				),
			},
		},
	})

}

func TestAccAlicloudECSSecurityGroupIngressRuleIpv6(t *testing.T) {
	var v ecs.Permission
	resourceId := "alicloud_security_group_rule.default"
	ra := resourceAttrInit(resourceId, testAccCheckSecurityGroupIngressRuleIpv6Map)
	serviceFunc := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSecurityGroupRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityGroupIngressRuleIpv6,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ipv6_cidr_ip": "2408:4004:cc:400::/56",
					}),
				),
			},
		},
	})

}

const testAccSecurityGroupRuleBasic = `
variable "name" {
  default = "tf-testAccSecurityGroupRuleBasic"
}

data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}

resource "alicloud_security_group" "default" {
  count  = 2
  vpc_id = data.alicloud_vpcs.default.ids.0
  name = "${var.name}"
}

resource "alicloud_security_group_rule" "default" {
  type = "ingress"
  ip_protocol = "tcp"
  nic_type = "intranet"
  policy = "drop"
  port_range = "22/22"
  priority = 100
  security_group_id = "${alicloud_security_group.default.0.id}"
  source_security_group_id = "${alicloud_security_group.default.1.id}"
  description = "abc"
}
`

const testAccSecurityGroupRulePrefix = `
variable "name" {
  default = "tf-testAccSecurityGroupRuleBasic"
}

data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}


resource "alicloud_security_group" "default" {
  vpc_id = "${data.alicloud_vpcs.default.vpcs.0.id}"
  name = "${var.name}"
}
resource "alicloud_ecs_prefix_list" "default"{
	address_family = "IPv4"
	max_entries = 2
	prefix_list_name = "tftest"
	description = "description"
	entry {
		cidr = "192.168.0.0/24"
		description = "description"
	}
}

resource "alicloud_security_group_rule" "default" {
  type = "ingress"
  ip_protocol = "tcp"
  prefix_list_id = "${alicloud_ecs_prefix_list.default.id}"
  nic_type = "intranet"
  policy = "accept"
  port_range = "22/22"
  priority = 100
  security_group_id = "${alicloud_security_group.default.id}"
  description = "abc"
}
`

const testAccSecurityGroupRule_cidrIp = `

variable "name" {
  default = "tf-testAccSecurityGroupRuleBasic"
}

data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}

resource "alicloud_security_group" "default" {
  count = 2
  vpc_id = data.alicloud_vpcs.default.ids.0
  name = "${var.name}"
}

resource "alicloud_security_group_rule" "default" {
  type = "ingress"
  ip_protocol = "tcp"
  nic_type = "intranet"
  policy = "drop"
  port_range = "22/22"
  priority = 100
  security_group_id = "${alicloud_security_group.default.0.id}"
  cidr_ip = "0.0.0.0/0"
  description = "abcd"
}
`

const testAccSecurityGroupRule_description = `

variable "name" {
  default = "tf-testAccSecurityGroupRule_description"
}

data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}
resource "alicloud_security_group" "default" {
  count = 2
  vpc_id = data.alicloud_vpcs.default.ids.0
  name = "${var.name}"
}

resource "alicloud_security_group_rule" "default" {
  type = "ingress"
  ip_protocol = "tcp"
  nic_type = "intranet"
  policy = "drop"
  port_range = "22/22"
  priority = 100
  security_group_id = "${alicloud_security_group.default.0.id}"
  cidr_ip = "0.0.0.0/0"
  description = "description"
}
`

const testAccSecurityGroupRule_all = `

variable "name" {
  default = "tf-testAccSecurityGroupRule_description"
}

data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}

resource "alicloud_security_group" "default" {
  count = 2
  vpc_id = data.alicloud_vpcs.default.ids.0
  name = "${var.name}"
}

resource "alicloud_security_group_rule" "default" {
  type = "ingress"
  ip_protocol = "tcp"
  nic_type = "intranet"
  policy = "drop"
  port_range = "22/22"
  priority = 100
  security_group_id = "${alicloud_security_group.default.0.id}"
  cidr_ip = "0.0.0.0/0"
  description = "abcd"
}
`

const testAccSecurityGroupRuleMulti = `
variable "name" {
  default = "tf-testAccSecurityGroupRuleBasic"
}

variable "cidr_ip_list" {
  type = "list"
  default = ["50.255.255.255/32", "75.250.250.250/32", "45.20.250.240/32"]
}

data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}

resource "alicloud_security_group" "default" {
  name = "${var.name}"
  description = "Security group for rules"
  vpc_id = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_security_group_rule" "default" {
  count = "${length(compact(var.cidr_ip_list))}"
  security_group_id = "${alicloud_security_group.default.id}"
  type = "ingress"
  policy = "drop"
  port_range = "22/22"
  ip_protocol = "tcp"
  nic_type = "intranet"
  priority = 100
  cidr_ip = "${element(var.cidr_ip_list, count.index)}"
}
`

const testAccSecurityGroupIngressRuleIpv6 = `
variable "name" {
  default = "tf-testAccSecurityGroupRuleBasic"
}

data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}
resource "alicloud_security_group" "default" {
  name = "${var.name}"
  description = "Security group for rules"
  vpc_id = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_security_group_rule" "default" {
  security_group_id = "${alicloud_security_group.default.id}"
  type = "ingress"
  policy = "drop"
  port_range = "22/22"
  ip_protocol = "tcp"
  nic_type = "intranet"
  priority = 100
  ipv6_cidr_ip = "2408:4004:cc:400::/56"
}
`

var testAccCheckSecurityGroupRuleBasicMap = map[string]string{
	"type":                     "ingress",
	"ip_protocol":              "tcp",
	"nic_type":                 "intranet",
	"policy":                   "drop",
	"port_range":               "22/22",
	"priority":                 "100",
	"security_group_id":        CHECKSET,
	"source_security_group_id": CHECKSET,
	"cidr_ip":                  "",
}

var testAccCheckSecurityGroupIngressRuleIpv6Map = map[string]string{
	"type":              "ingress",
	"ip_protocol":       "tcp",
	"nic_type":          "intranet",
	"policy":            "drop",
	"port_range":        "22/22",
	"priority":          "100",
	"security_group_id": CHECKSET,
	"cidr_ip":           "",
}

var testAccCheckSecurityGroupRulePrefixList = map[string]string{
	"type":              "ingress",
	"ip_protocol":       "tcp",
	"nic_type":          "intranet",
	"policy":            "accept",
	"port_range":        "22/22",
	"priority":          "100",
	"security_group_id": CHECKSET,
}

func testAccCheckSecurityGroupRuleDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	ecsService := EcsService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_security_group_rule" {
			continue
		}
		_, err := ecsService.DescribeSecurityGroupRule(rs.Primary.ID)

		// Verify the error is what we want
		if err != nil && !NotFoundError(err) {
			return WrapError(err)
		}
	}

	return nil
}

const testAccSecurityGroupEgressRule = `
variable "name" {
  default = "tf-testAccSecurityGroupRuleBasic"
}

data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}

resource "alicloud_security_group" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
  name = "${var.name}"
}

resource "alicloud_security_group_rule" "default" {
  type = "egress"
  ip_protocol = "tcp"
  nic_type = "intranet"
  policy = "accept"
  port_range = "443/443"
  priority = "1"
  security_group_id = "${alicloud_security_group.default.id}"
  cidr_ip = "182.254.11.243/32"
  description = "SHDRP-7513"
}
`

const testAccSecurityGroupEgressRuleIpv6 = `
variable "name" {
  default = "tf-testAccSecurityGroupRuleBasic"
}

data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}

resource "alicloud_security_group" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
  name = "${var.name}"
}

resource "alicloud_security_group_rule" "default" {
  type = "egress"
  ip_protocol = "tcp"
  nic_type = "intranet"
  policy = "accept"
  port_range = "443/443"
  priority = "1"
  security_group_id = "${alicloud_security_group.default.id}"
  ipv6_cidr_ip = "2408:4004:cc:400::/56"
  description = "SHDRP-7513"
}
`

const testAccSecurityGroupEgressRule_description = `
variable "name" {
  default = "tf-testAccSecurityGroupRuleBasic"
}

data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}

resource "alicloud_security_group" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
  name = "${var.name}"
}

resource "alicloud_security_group_rule" "default" {
  type = "egress"
  ip_protocol = "tcp"
  nic_type = "intranet"
  policy = "accept"
  port_range = "443/443"
  priority = "1"
  security_group_id = "${alicloud_security_group.default.id}"
  cidr_ip = "182.254.11.243/32"
  description = "SHDRP-7512"
}
`
