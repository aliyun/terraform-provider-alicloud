package cloudfirewall_test

import (
	"fmt"
	"testing"

	acctest "github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	connectivity "github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	cloudfirewall "github.com/aliyun/terraform-provider-alicloud/alicloud/service/cloudfirewall"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	resource "github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudCloudFirewallVpcFirewallControlPolicy_basic_no_order(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_firewall_vpc_firewall_control_policy.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudCloudFirewallVpcFirewallControlPolicyMap11929)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cloudfirewall.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCloudFirewallVpcFirewallControlPolicy")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccloudfirewall%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudCloudFirewallVpcFirewallControlPolicyBasicDependence11929)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"acl_action":       "accept",
					"description":      name,
					"destination":      "127.0.0.2/32",
					"destination_type": "net",
					"proto":            "TCP",
					"source":           "127.0.0.1/32",
					"source_type":      "net",
					"vpc_firewall_id":  "${alicloud_cen_instance.default.id}",
					"application_name": "ANY",
					"dest_port":        "80/88",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"acl_action":       "accept",
						"description":      name,
						"destination":      "127.0.0.2/32",
						"destination_type": "net",
						"proto":            "TCP",
						"order":            "1", // if only one policy and order not set, the default order should be 1
						"source":           "127.0.0.1/32",
						"source_type":      "net",
						"vpc_firewall_id":  acctest.CHECKSET,
						"application_name": "ANY",
						"dest_port":        "80/88",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"lang"},
			},
		},
	})
}

// Test CloudFirewall VpcFirewallControlPolicy. >>> Resource test cases, automatically generated.
// Case VPC边界安全策略_ip修改为域名 11929
func TestAccAliCloudCloudFirewallVpcFirewallControlPolicy_basic11929(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_firewall_vpc_firewall_control_policy.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudCloudFirewallVpcFirewallControlPolicyMap11929)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cloudfirewall.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCloudFirewallVpcFirewallControlPolicy")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccloudfirewall%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudCloudFirewallVpcFirewallControlPolicyOrderDependence11929)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"acl_action":       "accept",
					"description":      name,
					"destination":      "127.0.0.2/32",
					"destination_type": "net",
					"proto":            "TCP",
					"order":            "1",
					"source":           "127.0.0.1/32",
					"source_type":      "net",
					"vpc_firewall_id":  "${alicloud_cen_instance.default.id}",
					"application_name": "ANY",
					"dest_port":        "80/88",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"acl_action":       "accept",
						"description":      name,
						"destination":      "127.0.0.2/32",
						"destination_type": "net",
						"proto":            "TCP",
						"order":            "1",
						"source":           "127.0.0.1/32",
						"source_type":      "net",
						"vpc_firewall_id":  acctest.CHECKSET,
						"application_name": "ANY",
						"dest_port":        "80/88",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"order": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{"order": "2"}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"order": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{"order": "1"}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"acl_action": "drop",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"acl_action": "drop",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"application_name": "HTTP",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"application_name": "HTTP",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"dest_port": "20/22",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dest_port": "20/22",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"dest_port":       acctest.REMOVEKEY,
					"dest_port_group": "${alicloud_cloud_firewall_address_book.port.group_name}",
					"dest_port_type":  "group",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dest_port":       acctest.REMOVEKEY,
						"dest_port_group": acctest.CHECKSET,
						"dest_port_type":  "group",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"destination": "127.0.0.3/32",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"destination": "127.0.0.3/32",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"destination":      "alibaba.com",
					"destination_type": "domain",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"destination":      acctest.CHECKSET,
						"destination_type": "domain",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"domain_resolve_type": "DNS",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"domain_resolve_type": "DNS",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"repeat_days": []string{
						"1", "2", "3"},
					"end_time":          "1766676600",
					"start_time":        "1765209600",
					"repeat_end_time":   "19:00",
					"repeat_type":       "Weekly",
					"repeat_start_time": "18:00",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"repeat_days.#":     "3",
						"end_time":          "1766676600",
						"start_time":        "1765209600",
						"repeat_end_time":   "19:00",
						"repeat_type":       "Weekly",
						"repeat_start_time": "18:00",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"application_name": "ANY",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"application_name": "ANY",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"proto": "UDP",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"proto": "UDP",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"release": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"release": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"source": "127.0.0.2/32",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source": "127.0.0.2/32",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"source":      "${alicloud_cloud_firewall_address_book.ip.group_name}",
					"source_type": "group",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source":      acctest.CHECKSET,
						"source_type": "group",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"lang"},
			},
		},
	})
}

func TestAccAliCloudCloudFirewallVpcFirewallControlPolicy_basic11929_twin0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_firewall_vpc_firewall_control_policy.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudCloudFirewallVpcFirewallControlPolicyMap11929)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cloudfirewall.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCloudFirewallVpcFirewallControlPolicy")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccloudfirewall%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudCloudFirewallVpcFirewallControlPolicyBasicDependence11929)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"acl_action":          "accept",
					"description":         name,
					"destination":         "alibaba.com",
					"destination_type":    "domain",
					"proto":               "TCP",
					"order":               "1",
					"source":              "127.0.0.1/32",
					"source_type":         "net",
					"vpc_firewall_id":     "${alicloud_cen_instance.default.id}",
					"application_name":    "ANY",
					"dest_port":           "80/88",
					"dest_port_type":      "port",
					"domain_resolve_type": "DNS",
					"repeat_days": []string{
						"1", "2", "3"},
					"end_time":          "1766676600",
					"start_time":        "1765209600",
					"repeat_end_time":   "19:00",
					"repeat_type":       "Weekly",
					"repeat_start_time": "18:00",
					"member_uid":        "${data.alicloud_account.current.id}",
					"release":           "false",
					"lang":              "zh",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"acl_action":          "accept",
						"description":         name,
						"destination":         "alibaba.com",
						"destination_type":    "domain",
						"proto":               "TCP",
						"order":               "1",
						"source":              "127.0.0.1/32",
						"source_type":         "net",
						"vpc_firewall_id":     acctest.CHECKSET,
						"application_name":    "ANY",
						"dest_port":           "80/88",
						"dest_port_type":      "port",
						"domain_resolve_type": "DNS",
						"repeat_days.#":       "3",
						"end_time":            "1766676600",
						"start_time":          "1765209600",
						"repeat_end_time":     "19:00",
						"repeat_type":         "Weekly",
						"repeat_start_time":   "18:00",
						"member_uid":          acctest.CHECKSET,
						"release":             "false",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"lang"},
			},
		},
	})
}

func TestAccAliCloudCloudFirewallVpcFirewallControlPolicy_basic11929_twin1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_firewall_vpc_firewall_control_policy.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudCloudFirewallVpcFirewallControlPolicyMap11929)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cloudfirewall.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCloudFirewallVpcFirewallControlPolicy")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccloudfirewall%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudCloudFirewallVpcFirewallControlPolicyBasicDependence11929)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"acl_action":          "accept",
					"description":         name,
					"destination":         "alibaba.com",
					"destination_type":    "domain",
					"proto":               "TCP",
					"order":               "1",
					"source":              "127.0.0.1/32",
					"source_type":         "net",
					"vpc_firewall_id":     "${alicloud_cen_instance.default.id}",
					"application_name":    "ANY",
					"dest_port_group":     "${alicloud_cloud_firewall_address_book.port.group_name}",
					"dest_port_type":      "group",
					"domain_resolve_type": "DNS",
					"repeat_days": []string{
						"1", "2", "3"},
					"end_time":          "1766676600",
					"start_time":        "1765209600",
					"repeat_end_time":   "19:00",
					"repeat_type":       "Weekly",
					"repeat_start_time": "18:00",
					"member_uid":        "${data.alicloud_account.current.id}",
					"release":           "false",
					"lang":              "zh",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"acl_action":          "accept",
						"description":         name,
						"destination":         "alibaba.com",
						"destination_type":    "domain",
						"proto":               "TCP",
						"order":               "1",
						"source":              "127.0.0.1/32",
						"source_type":         "net",
						"vpc_firewall_id":     acctest.CHECKSET,
						"application_name":    "ANY",
						"dest_port_group":     acctest.CHECKSET,
						"dest_port_type":      "group",
						"domain_resolve_type": "DNS",
						"repeat_days.#":       "3",
						"end_time":            "1766676600",
						"start_time":          "1765209600",
						"repeat_end_time":     "19:00",
						"repeat_type":         "Weekly",
						"repeat_start_time":   "18:00",
						"member_uid":          acctest.CHECKSET,
						"release":             "false",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"lang"},
			},
		},
	})
}

var AliCloudCloudFirewallVpcFirewallControlPolicyMap11929 = map[string]string{
	"acl_uuid":                  acctest.CHECKSET,
	"application_id":            acctest.CHECKSET,
	"create_time":               acctest.CHECKSET,
	"dest_port":                 acctest.CHECKSET,
	"dest_port_group_ports.#":   acctest.CHECKSET,
	"dest_port_type":            acctest.CHECKSET,
	"destination_group_cidrs.#": acctest.CHECKSET,
	"hit_times":                 acctest.CHECKSET,
	"member_uid":                acctest.CHECKSET,
	"release":                   acctest.CHECKSET,
	"repeat_type":               acctest.CHECKSET,
	"source_group_cidrs.#":      acctest.CHECKSET,
}

// A second policy makes moving the tested policy to position two valid.
func AliCloudCloudFirewallVpcFirewallControlPolicyOrderDependence11929(name string) string {
	return AliCloudCloudFirewallVpcFirewallControlPolicyBasicDependence11929(name) + `
resource "alicloud_cloud_firewall_vpc_firewall_control_policy" "sibling" {
  acl_action       = "accept"
  description      = "${var.name}-sibling"
  destination      = "127.0.0.3/32"
  destination_type = "net"
  proto            = "TCP"
  source           = "127.0.0.1/32"
  source_type      = "net"
  vpc_firewall_id  = alicloud_cen_instance.default.id
  application_name = "ANY"
  dest_port        = "80/88"
}
`
}

func AliCloudCloudFirewallVpcFirewallControlPolicyBasicDependence11929(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_account" "current" {
}

resource "alicloud_cen_instance" "default" {
  cen_instance_name = var.name
}

resource "alicloud_cloud_firewall_address_book" "port" {
  description  = "${var.name}-port"
  group_name   = "${var.name}-port"
  group_type   = "port"
  address_list = ["8081/8081"]
}

resource "alicloud_cloud_firewall_address_book" "ip" {
  description  = var.name
  group_name   = var.name
  group_type   = "ip"
  address_list = ["10.21.0.0/16", "10.168.0.0/16"]
}
`, name)
}

// Case VPC边界安全策略_IP_策略有效期(重复周期) 12022
func TestAccAliCloudCloudFirewallVpcFirewallControlPolicy_basic12022(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_firewall_vpc_firewall_control_policy.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudCloudFirewallVpcFirewallControlPolicyMap12022)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cloudfirewall.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCloudFirewallVpcFirewallControlPolicy")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccloudfirewall%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudCloudFirewallVpcFirewallControlPolicyBasicDependence12022)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"acl_action":            "accept",
					"description":           name,
					"destination":           "127.0.0.2/32",
					"destination_type":      "net",
					"proto":                 "TCP",
					"order":                 "1",
					"source":                "127.0.0.1/32",
					"source_type":           "net",
					"vpc_firewall_id":       "${alicloud_cen_instance.default.id}",
					"application_name_list": []string{"ANY"},
					"dest_port":             "80/88",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"acl_action":              "accept",
						"description":             name,
						"destination":             "127.0.0.2/32",
						"destination_type":        "net",
						"proto":                   "TCP",
						"order":                   "1",
						"source":                  "127.0.0.1/32",
						"source_type":             "net",
						"vpc_firewall_id":         acctest.CHECKSET,
						"application_name_list.#": "1",
						"dest_port":               "80/88",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"acl_action": "drop",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"acl_action": "drop",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"application_name_list": []string{"HTTP", "SMTP", "HTTPS"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"application_name_list.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"dest_port": "20/22",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dest_port": "20/22",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"dest_port":       acctest.REMOVEKEY,
					"dest_port_group": "${alicloud_cloud_firewall_address_book.port.group_name}",
					"dest_port_type":  "group",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dest_port":       acctest.REMOVEKEY,
						"dest_port_group": acctest.CHECKSET,
						"dest_port_type":  "group",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"destination": "127.0.0.3/32",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"destination": "127.0.0.3/32",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"destination":      "alibaba.com",
					"destination_type": "domain",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"destination":      acctest.CHECKSET,
						"destination_type": "domain",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"domain_resolve_type": "DNS",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"domain_resolve_type": "DNS",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"repeat_days": []string{
						"1", "2", "3"},
					"end_time":          "1766676600",
					"start_time":        "1765209600",
					"repeat_end_time":   "19:00",
					"repeat_type":       "Weekly",
					"repeat_start_time": "18:00",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"repeat_days.#":     "3",
						"end_time":          "1766676600",
						"start_time":        "1765209600",
						"repeat_end_time":   "19:00",
						"repeat_type":       "Weekly",
						"repeat_start_time": "18:00",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"application_name_list": []string{"ANY"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"application_name_list.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"proto": "UDP",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"proto": "UDP",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"release": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"release": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"source": "127.0.0.2/32",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source": "127.0.0.2/32",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"source":      "${alicloud_cloud_firewall_address_book.ip.group_name}",
					"source_type": "group",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source":      acctest.CHECKSET,
						"source_type": "group",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"lang"},
			},
		},
	})
}

func TestAccAliCloudCloudFirewallVpcFirewallControlPolicy_basic12022_twin0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_firewall_vpc_firewall_control_policy.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudCloudFirewallVpcFirewallControlPolicyMap12022)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cloudfirewall.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCloudFirewallVpcFirewallControlPolicy")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccloudfirewall%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudCloudFirewallVpcFirewallControlPolicyBasicDependence12022)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"acl_action":            "accept",
					"description":           name,
					"destination":           "alibaba.com",
					"destination_type":      "domain",
					"proto":                 "TCP",
					"order":                 "1",
					"source":                "127.0.0.1/32",
					"source_type":           "net",
					"vpc_firewall_id":       "${alicloud_cen_instance.default.id}",
					"application_name_list": []string{"ANY"},
					"dest_port":             "80/88",
					"dest_port_type":        "port",
					"domain_resolve_type":   "DNS",
					"repeat_days": []string{
						"1", "2", "3"},
					"end_time":          "1766676600",
					"start_time":        "1765209600",
					"repeat_end_time":   "19:00",
					"repeat_type":       "Weekly",
					"repeat_start_time": "18:00",
					"member_uid":        "${data.alicloud_account.current.id}",
					"release":           "false",
					"lang":              "zh",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"acl_action":              "accept",
						"description":             name,
						"destination":             "alibaba.com",
						"destination_type":        "domain",
						"proto":                   "TCP",
						"order":                   "1",
						"source":                  "127.0.0.1/32",
						"source_type":             "net",
						"vpc_firewall_id":         acctest.CHECKSET,
						"application_name_list.#": "1",
						"dest_port":               "80/88",
						"dest_port_type":          "port",
						"domain_resolve_type":     "DNS",
						"repeat_days.#":           "3",
						"end_time":                "1766676600",
						"start_time":              "1765209600",
						"repeat_end_time":         "19:00",
						"repeat_type":             "Weekly",
						"repeat_start_time":       "18:00",
						"member_uid":              acctest.CHECKSET,
						"release":                 "false",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"lang"},
			},
		},
	})
}

func TestAccAliCloudCloudFirewallVpcFirewallControlPolicy_basic12022_twin1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_firewall_vpc_firewall_control_policy.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudCloudFirewallVpcFirewallControlPolicyMap12022)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cloudfirewall.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCloudFirewallVpcFirewallControlPolicy")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccloudfirewall%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudCloudFirewallVpcFirewallControlPolicyBasicDependence12022)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"acl_action":            "accept",
					"description":           name,
					"destination":           "alibaba.com",
					"destination_type":      "domain",
					"proto":                 "TCP",
					"order":                 "1",
					"source":                "127.0.0.1/32",
					"source_type":           "net",
					"vpc_firewall_id":       "${alicloud_cen_instance.default.id}",
					"application_name_list": []string{"ANY"},
					"dest_port_group":       "${alicloud_cloud_firewall_address_book.port.group_name}",
					"dest_port_type":        "group",
					"domain_resolve_type":   "DNS",
					"repeat_days": []string{
						"1", "2", "3"},
					"end_time":          "1766676600",
					"start_time":        "1765209600",
					"repeat_end_time":   "19:00",
					"repeat_type":       "Weekly",
					"repeat_start_time": "18:00",
					"member_uid":        "${data.alicloud_account.current.id}",
					"release":           "false",
					"lang":              "zh",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"acl_action":              "accept",
						"description":             name,
						"destination":             "alibaba.com",
						"destination_type":        "domain",
						"proto":                   "TCP",
						"order":                   "1",
						"source":                  "127.0.0.1/32",
						"source_type":             "net",
						"vpc_firewall_id":         acctest.CHECKSET,
						"application_name_list.#": "1",
						"dest_port_group":         acctest.CHECKSET,
						"dest_port_type":          "group",
						"domain_resolve_type":     "DNS",
						"repeat_days.#":           "3",
						"end_time":                "1766676600",
						"start_time":              "1765209600",
						"repeat_end_time":         "19:00",
						"repeat_type":             "Weekly",
						"repeat_start_time":       "18:00",
						"member_uid":              acctest.CHECKSET,
						"release":                 "false",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"lang"},
			},
		},
	})
}

var AliCloudCloudFirewallVpcFirewallControlPolicyMap12022 = map[string]string{
	"acl_uuid":                  acctest.CHECKSET,
	"application_id":            acctest.CHECKSET,
	"create_time":               acctest.CHECKSET,
	"dest_port":                 acctest.CHECKSET,
	"dest_port_group_ports.#":   acctest.CHECKSET,
	"dest_port_type":            acctest.CHECKSET,
	"destination_group_cidrs.#": acctest.CHECKSET,
	"hit_times":                 acctest.CHECKSET,
	"member_uid":                acctest.CHECKSET,
	"release":                   acctest.CHECKSET,
	"repeat_type":               acctest.CHECKSET,
	"source_group_cidrs.#":      acctest.CHECKSET,
}

func AliCloudCloudFirewallVpcFirewallControlPolicyBasicDependence12022(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_account" "current" {
}

resource "alicloud_cen_instance" "default" {
  cen_instance_name = var.name
}

resource "alicloud_cloud_firewall_address_book" "port" {
  description  = "${var.name}-port"
  group_name   = "${var.name}-port"
  group_type   = "port"
  address_list = ["8081/8081"]
}

resource "alicloud_cloud_firewall_address_book" "ip" {
  description  = var.name
  group_name   = var.name
  group_type   = "ip"
  address_list = ["10.21.0.0/16", "10.168.0.0/16"]
}
`, name)
}

// TestAccAliCloudCloudFirewallVpcFirewallControlPolicyApplicationNameListOrder verifies that
// reordering application_name_list (a plain TypeList of TypeString with no
// DiffSuppressFunc) converges to an empty post-apply plan.
func TestAccAliCloudCloudFirewallVpcFirewallControlPolicyApplicationNameListOrder(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_firewall_vpc_firewall_control_policy.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudCloudFirewallVpcFirewallControlPolicyMap12022)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cloudfirewall.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCloudFirewallVpcFirewallControlPolicy")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	name := fmt.Sprintf("tf-testacc%scloudfirewallvpcfirewallcontrolpolicy%d", acctest.DefaultRegionToTest, tfacctest.RandIntRange(10000, 99999))
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudCloudFirewallVpcFirewallControlPolicyBasicDependence12022)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"acl_action":            "accept",
					"description":           name,
					"destination":           "127.0.0.2/32",
					"destination_type":      "net",
					"proto":                 "TCP",
					"order":                 "1",
					"source":                "127.0.0.1/32",
					"source_type":           "net",
					"vpc_firewall_id":       "${alicloud_cen_instance.default.id}",
					"application_name_list": []string{"ANY", "HTTP", "HTTPS"},
					"dest_port":             "80/88",
				}),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"acl_action":            "accept",
					"description":           name,
					"destination":           "127.0.0.2/32",
					"destination_type":      "net",
					"proto":                 "TCP",
					"order":                 "1",
					"source":                "127.0.0.1/32",
					"source_type":           "net",
					"vpc_firewall_id":       "${alicloud_cen_instance.default.id}",
					"application_name_list": []string{"HTTPS", "ANY", "HTTP"},
					"dest_port":             "80/88",
				}),
			},
		},
	})
}

// TestAccAliCloudCloudFirewallVpcFirewallControlPolicyRepeatDaysOrder verifies that
// reordering repeat_days (a plain TypeList of TypeInt with no DiffSuppressFunc)
// converges to an empty post-apply plan.
func TestAccAliCloudCloudFirewallVpcFirewallControlPolicyRepeatDaysOrder(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_firewall_vpc_firewall_control_policy.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudCloudFirewallVpcFirewallControlPolicyMap11929)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cloudfirewall.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCloudFirewallVpcFirewallControlPolicy")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	name := fmt.Sprintf("tf-testacc%scloudfirewallvpcfirewallcontrolpolicy%d", acctest.DefaultRegionToTest, tfacctest.RandIntRange(10000, 99999))
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudCloudFirewallVpcFirewallControlPolicyBasicDependence11929)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"acl_action":        "accept",
					"description":       name,
					"destination":       "127.0.0.2/32",
					"destination_type":  "net",
					"proto":             "TCP",
					"order":             "1",
					"source":            "127.0.0.1/32",
					"source_type":       "net",
					"vpc_firewall_id":   "${alicloud_cen_instance.default.id}",
					"application_name":  "ANY",
					"dest_port":         "80/88",
					"repeat_days":       []string{"1", "2", "3"},
					"end_time":          "1766676600",
					"start_time":        "1765209600",
					"repeat_end_time":   "19:00",
					"repeat_type":       "Weekly",
					"repeat_start_time": "18:00",
				}),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"acl_action":        "accept",
					"description":       name,
					"destination":       "127.0.0.2/32",
					"destination_type":  "net",
					"proto":             "TCP",
					"order":             "1",
					"source":            "127.0.0.1/32",
					"source_type":       "net",
					"vpc_firewall_id":   "${alicloud_cen_instance.default.id}",
					"application_name":  "ANY",
					"dest_port":         "80/88",
					"repeat_days":       []string{"3", "1", "2"},
					"end_time":          "1766676600",
					"start_time":        "1765209600",
					"repeat_end_time":   "19:00",
					"repeat_type":       "Weekly",
					"repeat_start_time": "18:00",
				}),
			},
		},
	})
}

// Test CloudFirewall VpcFirewallControlPolicy. <<< Resource test cases, automatically generated.
