package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAliCloudECSSecurityGroupRuleBasic(t *testing.T) {
	var v ecs.Permission
	resourceId := "alicloud_security_group_rule.test"
	name := acctest.RandString(4)
	serviceFunc := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rac := resourceAttrCheckInit(
		resourceCheckInit(resourceId, &v, serviceFunc),
		resourceAttrInit(resourceId, testAccCheckSecurityGroupRuleBasicMap))
	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSecurityGroupRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: hclSecurityGroupRuleBasic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "abc",
					}),
				),
			},
			{
				Config: hclSecurityGroupRuleCidrIp(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source_security_group_id": REMOVEKEY,
						"cidr_ip":                  "0.0.0.0/0",
						"description":              "abcd",
					}),
				),
			},
			{
				Config: hclSecurityGroupRuleDescription(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "description",
					}),
				),
			},
			{
				Config: hclSecurityGroupRuleAll(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "abcd",
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

func TestAccAliCloudECSSecurityGroupRuleEgress(t *testing.T) {
	var v ecs.Permission
	resourceId := "alicloud_security_group_rule.test"
	name := acctest.RandString(4)
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
				Config: hclSecurityGroupEgressRule(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "SHDRP-7513",
					}),
				),
			},
			{
				Config: hclSecurityGroupEgressRuleDescription(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "SHDRP-7512",
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

func TestAccAliCloudECSSecurityGroupRuleMulti(t *testing.T) {
	var v ecs.Permission
	resourceId := "alicloud_security_group_rule.test.2"
	name := acctest.RandString(4)
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
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSecurityGroupRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: hclSecurityGroupRuleMulti(name),
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

func TestAccAliCloudECSSecurityGroupRulePrefixList(t *testing.T) {
	var v ecs.Permission
	resourceId := "alicloud_security_group_rule.test"
	name := acctest.RandString(4)
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
				Config: hclSecurityGroupRulePrefix(name),
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

func TestAccAliCloudECSSecurityGroupRuleEgressIpv6(t *testing.T) {
	var v ecs.Permission
	resourceId := "alicloud_security_group_rule.test"
	name := acctest.RandString(4)
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
				Config: hclSecurityGroupEgressRuleIpv6(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "SHDRP-7513",
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

func TestAccAliCloudECSSecurityGroupRuleIngressIpv6(t *testing.T) {
	var v ecs.Permission
	resourceId := "alicloud_security_group_rule.test"
	name := acctest.RandString(4)
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
				Config: hclSecurityGroupIngressRuleIpv6(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ipv6_cidr_ip": "2408:4004:cc:400::/56",
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

func TestAccAliCloudECSSecurityGroupRuleEgressOtherIpv6(t *testing.T) {
	var v ecs.Permission
	resourceId := "alicloud_security_group_rule.test"
	name := acctest.RandString(4)
	ra := resourceAttrInit(resourceId, map[string]string{
		"type":         "egress",
		"policy":       "accept",
		"description":  "SHDRP-7513",
		"port_range":   "443/443",
		"priority":     "1",
		"ipv6_cidr_ip": "2001:0db8:3c4d:0015:0000:0000:1a2f:1a2b/0",
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
				Config: hclSecurityGroupEgressRuleOtherIpv6(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ipv6_cidr_ip": "2001:db8:3c4d:15::1a2f:1a2b/0",
						"description":  "SHDRP-7513",
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

func TestAccAliCloudECSSecurityGroupRuleIngressOtherIpv6(t *testing.T) {
	var v ecs.Permission
	resourceId := "alicloud_security_group_rule.test"
	name := acctest.RandString(4)
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
				Config: hclSecurityGroupIngressRuleOtherIpv6(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ipv6_cidr_ip": "2001:db8:3c4d:15::1a2f:1a2b/0",
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

func TestAccAliCloudECSSecurityGroupRuleEgressICMPv6(t *testing.T) {
	var v ecs.Permission
	resourceId := "alicloud_security_group_rule.test"
	name := acctest.RandString(4)
	ra := resourceAttrInit(resourceId, map[string]string{
		"type":         "ingress",
		"policy":       "accept",
		"description":  "SHDRP-7513",
		"port_range":   "-1/-1",
		"priority":     "1",
		"ipv6_cidr_ip": "::/0",
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
				Config: hclSecurityGroupEgressRuleICMPv6(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "SHDRP-7513",
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

var testAccCheckSecurityGroupRuleBasicMap = map[string]string{
	"type":                     "ingress",
	"ip_protocol":              "tcp",
	"nic_type":                 "intranet",
	"policy":                   "drop",
	"port_range":               "22/22",
	"priority":                 "100",
	"security_group_id":        CHECKSET,
	"source_security_group_id": CHECKSET,
	"security_group_rule_id":   CHECKSET,
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

func hclSecurityGroupRuleBasic(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "tf-testAccSGRBase%s"
}

data "alicloud_vpcs" "test" {
  name_regex = "^default-NODELETING$"
}

resource "alicloud_security_group" "test" {
  count               = 2
  vpc_id              = data.alicloud_vpcs.test.ids.0
  security_group_name = var.name
}

resource "alicloud_security_group_rule" "test" {
  type                     = "ingress"
  ip_protocol              = "tcp"
  nic_type                 = "intranet"
  policy                   = "drop"
  port_range               = "22/22"
  priority                 = 100
  security_group_id        = alicloud_security_group.test.0.id
  source_security_group_id = alicloud_security_group.test.1.id
  description              = "abc"
}
`, name)
}

func hclSecurityGroupRulePrefix(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "tf-testAccSGRPrefix%s"
}

data "alicloud_vpcs" "test" {
  name_regex = "^default-NODELETING$"
}

resource "alicloud_security_group" "test" {
  vpc_id              = data.alicloud_vpcs.test.vpcs.0.id
  security_group_name = var.name
}

resource "alicloud_ecs_prefix_list" "test" {
  address_family   = "IPv4"
  max_entries      = 2
  prefix_list_name = "tftest"
  description      = "description"
  entry {
    cidr        = "192.168.0.0/24"
    description = "description"
  }
}

resource "alicloud_security_group_rule" "test" {
  type              = "ingress"
  ip_protocol       = "tcp"
  prefix_list_id    = alicloud_ecs_prefix_list.test.id
  nic_type          = "intranet"
  policy            = "accept"
  port_range        = "22/22"
  priority          = 100
  security_group_id = alicloud_security_group.test.id
  description       = "abc"
}
`, name)
}

func hclSecurityGroupRuleCidrIp(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "tf-testAccSGR-CIRDIP%s"
}

data "alicloud_vpcs" "test" {
  name_regex = "^default-NODELETING$"
}

resource "alicloud_security_group" "test" {
  count               = 2
  vpc_id              = data.alicloud_vpcs.test.ids.0
  security_group_name = var.name
}

resource "alicloud_security_group_rule" "test" {
  type              = "ingress"
  ip_protocol       = "tcp"
  nic_type          = "intranet"
  policy            = "drop"
  port_range        = "22/22"
  priority          = 100
  security_group_id = alicloud_security_group.test.0.id
  cidr_ip           = "0.0.0.0/0"
  description       = "abcd"
}
`, name)
}

func hclSecurityGroupRuleDescription(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "tf-testAccSGR-desc%s"
}

data "alicloud_vpcs" "test" {
  name_regex = "^default-NODELETING$"
}

resource "alicloud_security_group" "test" {
  count               = 2
  vpc_id              = data.alicloud_vpcs.test.ids.0
  security_group_name = var.name
}

resource "alicloud_security_group_rule" "test" {
  type              = "ingress"
  ip_protocol       = "tcp"
  nic_type          = "intranet"
  policy            = "drop"
  port_range        = "22/22"
  priority          = 100
  security_group_id = alicloud_security_group.test.0.id
  cidr_ip           = "0.0.0.0/0"
  description       = "description"
}
`, name)
}

func hclSecurityGroupRuleAll(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "tf-testAccSGR_all%s"
}

data "alicloud_vpcs" "test" {
  name_regex = "^default-NODELETING$"
}

resource "alicloud_security_group" "test" {
  count               = 2
  vpc_id              = data.alicloud_vpcs.test.ids.0
  security_group_name = var.name
}

resource "alicloud_security_group_rule" "test" {
  type              = "ingress"
  ip_protocol       = "tcp"
  nic_type          = "intranet"
  policy            = "drop"
  port_range        = "22/22"
  priority          = 100
  security_group_id = alicloud_security_group.test.0.id
  cidr_ip           = "0.0.0.0/0"
  description       = "abcd"
}
`, name)
}

func hclSecurityGroupRuleMulti(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "tf-testAccSGRMulti%s"
}

variable "cidr_ip_list" {
  type    = "list"
  default = ["50.255.255.255/32", "75.250.250.250/32", "45.20.250.240/32"]
}

data "alicloud_vpcs" "test" {
  name_regex = "^default-NODELETING$"
}

resource "alicloud_security_group" "test" {
  security_group_name = var.name
  description         = "Security group for rules"
  vpc_id              = data.alicloud_vpcs.test.ids.0
}

resource "alicloud_security_group_rule" "test" {
  count             = length(compact(var.cidr_ip_list))
  security_group_id = alicloud_security_group.test.id
  type              = "ingress"
  policy            = "drop"
  port_range        = "22/22"
  ip_protocol       = "tcp"
  nic_type          = "intranet"
  priority          = 100
  cidr_ip           = element(var.cidr_ip_list, count.index)
}
`, name)
}

func hclSecurityGroupIngressRuleIpv6(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "tf-testAccSGRIngressIpv6%s"
}

data "alicloud_vpcs" "test" {
  name_regex = "^default-NODELETING$"
}

resource "alicloud_security_group" "test" {
  security_group_name = var.name
  description         = "Security group for rules"
  vpc_id              = data.alicloud_vpcs.test.ids.0
}

resource "alicloud_security_group_rule" "test" {
  security_group_id = alicloud_security_group.test.id
  type              = "ingress"
  policy            = "drop"
  port_range        = "22/22"
  ip_protocol       = "tcp"
  nic_type          = "intranet"
  priority          = 100
  ipv6_cidr_ip      = "2408:4004:cc:400::/56"
}
`, name)
}

func hclSecurityGroupIngressRuleOtherIpv6(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "tf-testAccSGRIngressOtherIpv6%s"
}

data "alicloud_vpcs" "test" {
  name_regex = "^default-NODELETING$"
}

resource "alicloud_security_group" "test" {
  security_group_name = var.name
  description         = "Security group for rules"
  vpc_id              = data.alicloud_vpcs.test.ids.0
}

resource "alicloud_security_group_rule" "test" {
  security_group_id = alicloud_security_group.test.id
  type              = "ingress"
  policy            = "drop"
  port_range        = "22/22"
  ip_protocol       = "tcp"
  nic_type          = "intranet"
  priority          = 100
  ipv6_cidr_ip      = "2001:0db8:3c4d:0015:0000:0000:1a2f:1a2b/0"
}
`, name)
}

func hclSecurityGroupEgressRule(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "tf-testAccSGREgress%s"
}

data "alicloud_vpcs" "test" {
  name_regex = "^default-NODELETING$"
}

resource "alicloud_security_group" "test" {
  vpc_id              = data.alicloud_vpcs.test.ids.0
  security_group_name = var.name
}

resource "alicloud_security_group_rule" "test" {
  type              = "egress"
  ip_protocol       = "tcp"
  nic_type          = "intranet"
  policy            = "accept"
  port_range        = "443/443"
  priority          = "1"
  security_group_id = alicloud_security_group.test.id
  cidr_ip           = "182.254.11.243/32"
  description       = "SHDRP-7513"
}
`, name)
}

func hclSecurityGroupEgressRuleIpv6(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "tf-testAccSGREgressIpv6%s"
}

data "alicloud_vpcs" "test" {
  name_regex = "^default-NODELETING$"
}

resource "alicloud_security_group" "test" {
  vpc_id              = data.alicloud_vpcs.test.ids.0
  security_group_name = var.name
}

resource "alicloud_security_group_rule" "test" {
  type              = "egress"
  ip_protocol       = "tcp"
  nic_type          = "intranet"
  policy            = "accept"
  port_range        = "443/443"
  priority          = "1"
  security_group_id = alicloud_security_group.test.id
  ipv6_cidr_ip      = "2408:4004:cc:400::/56"
  description       = "SHDRP-7513"
}
`, name)
}

func hclSecurityGroupEgressRuleOtherIpv6(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "tf-testAccSGREgressOtherIpv6%s"
}

data "alicloud_vpcs" "test" {
  name_regex = "^default-NODELETING$"
}

resource "alicloud_security_group" "test" {
  vpc_id              = data.alicloud_vpcs.test.ids.0
  security_group_name = var.name
}

resource "alicloud_security_group_rule" "test" {
  type              = "egress"
  ip_protocol       = "tcp"
  nic_type          = "intranet"
  policy            = "accept"
  port_range        = "443/443"
  priority          = "1"
  security_group_id = alicloud_security_group.test.id
  ipv6_cidr_ip      = "2001:0db8:3c4d:0015:0000:0000:1a2f:1a2b/0"
  description       = "SHDRP-7513"
}
`, name)
}

func hclSecurityGroupEgressRuleDescription(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "tf-testAccSGREgressDesc%s"
}

data "alicloud_vpcs" "test" {
  name_regex = "^default-NODELETING$"
}

resource "alicloud_security_group" "test" {
  vpc_id              = data.alicloud_vpcs.test.ids.0
  security_group_name = var.name
  inner_access_policy = "Accept"
}

resource "alicloud_security_group_rule" "test" {
  type              = "egress"
  ip_protocol       = "tcp"
  nic_type          = "intranet"
  policy            = "accept"
  port_range        = "443/443"
  priority          = "1"
  security_group_id = alicloud_security_group.test.id
  cidr_ip           = "182.254.11.243/32"
  description       = "SHDRP-7512"
}
`, name)
}

func hclSecurityGroupEgressRuleICMPv6(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "tf-testAccSGREgressIpv6%s"
}

data "alicloud_vpcs" "test" {
  name_regex = "^default-NODELETING$"
}

resource "alicloud_security_group" "test" {
  vpc_id              = data.alicloud_vpcs.test.ids.0
  security_group_name = var.name
}

resource "alicloud_security_group_rule" "test" {
  type              = "ingress"
  ip_protocol       = "icmpv6"
  nic_type          = "intranet"
  policy            = "accept"
  port_range        = "-1/-1"
  priority          = "1"
  security_group_id = alicloud_security_group.test.id
  ipv6_cidr_ip      = "::/0"
  description       = "SHDRP-7513"
}
`, name)
}
