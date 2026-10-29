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

func TestAccAliCloudCloudFirewallControlPolicy_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_firewall_control_policy.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudCloudFirewallControlPolicyMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cloudfirewall.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCloudFirewallControlPolicy")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scloudfirewallcontrolpolicy%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudCloudFirewallControlPolicyBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"direction":        "in",
					"application_name": "ANY",
					"description":      name,
					"acl_action":       "accept",
					"source":           "127.0.0.1/32",
					"source_type":      "net",
					"destination":      "127.0.0.2/32",
					"destination_type": "net",
					"proto":            "TCP",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"direction":        "in",
						"application_name": "ANY",
						"description":      name,
						"acl_action":       "accept",
						"source":           "127.0.0.1/32",
						"source_type":      "net",
						"destination":      "127.0.0.2/32",
						"destination_type": "net",
						"proto":            "TCP",
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
				Config: testAccConfig(map[string]interface{}{
					"source": "${alicloud_cloud_firewall_address_book.ip_update.group_name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source": acctest.CHECKSET,
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
					"destination":      "${alicloud_cloud_firewall_address_book.ip.group_name}",
					"destination_type": "group",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"destination":      acctest.CHECKSET,
						"destination_type": "group",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"destination": "${alicloud_cloud_firewall_address_book.ip_update.group_name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"destination": acctest.CHECKSET,
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
					"proto": "ANY",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"proto": "ANY",
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
					"dest_port":       "20/22",
					"dest_port_group": acctest.REMOVEKEY,
					"dest_port_type":  "port",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dest_port":       "20/22",
						"dest_port_group": acctest.REMOVEKEY,
						"dest_port_type":  "port",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"repeat_type": "None",
					"start_time":  "1716998400",
					"end_time":    "1717083000",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"repeat_type": "None",
						"start_time":  acctest.CHECKSET,
						"end_time":    acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"repeat_type":       "Daily",
					"repeat_start_time": "08:00",
					"repeat_end_time":   "10:00",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"repeat_type":       "Daily",
						"repeat_start_time": "08:00",
						"repeat_end_time":   "10:00",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"repeat_type": "Weekly",
					"repeat_days": []string{"1", "2", "3"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"repeat_type":   "Weekly",
						"repeat_days.#": "3",
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
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"source_ip", "lang"},
			},
		},
	})
}

func TestAccAliCloudCloudFirewallControlPolicy_basic0_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_firewall_control_policy.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudCloudFirewallControlPolicyMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cloudfirewall.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCloudFirewallControlPolicy")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scloudfirewallcontrolpolicy%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudCloudFirewallControlPolicyBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"direction":         "in",
					"application_name":  "ANY",
					"description":       name,
					"acl_action":        "accept",
					"source":            "::1/128",
					"source_type":       "net",
					"destination":       "::2/128",
					"destination_type":  "net",
					"proto":             "TCP",
					"dest_port":         "20/22",
					"dest_port_type":    "port",
					"ip_version":        "6",
					"repeat_type":       "Weekly",
					"start_time":        "1716998400",
					"end_time":          "1717083000",
					"repeat_start_time": "08:00",
					"repeat_end_time":   "10:00",
					"repeat_days":       []string{"1", "2", "3"},
					"release":           "false",
					"source_ip":         "127.0.0.1",
					"lang":              "zh",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"direction":         "in",
						"application_name":  "ANY",
						"description":       name,
						"acl_action":        "accept",
						"source":            "::1/128",
						"source_type":       "net",
						"destination":       "::2/128",
						"destination_type":  "net",
						"proto":             "TCP",
						"dest_port":         "20/22",
						"dest_port_type":    "port",
						"ip_version":        "6",
						"repeat_type":       "Weekly",
						"start_time":        acctest.CHECKSET,
						"end_time":          acctest.CHECKSET,
						"repeat_start_time": "08:00",
						"repeat_end_time":   "10:00",
						"repeat_days.#":     "3",
						"release":           "false",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"source_ip", "lang"},
			},
		},
	})
}

func TestAccAliCloudCloudFirewallControlPolicy_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_firewall_control_policy.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudCloudFirewallControlPolicyMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cloudfirewall.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCloudFirewallControlPolicy")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scloudfirewallcontrolpolicy%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudCloudFirewallControlPolicyBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"direction":        "out",
					"application_name": "ANY",
					"description":      name,
					"acl_action":       "accept",
					"source":           "${alicloud_cloud_firewall_address_book.ip.group_name}",
					"source_type":      "group",
					"destination":      "${alicloud_cloud_firewall_address_book.ip.group_name}",
					"destination_type": "group",
					"proto":            "TCP",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"direction":        "out",
						"application_name": "ANY",
						"description":      name,
						"acl_action":       "accept",
						"source":           acctest.CHECKSET,
						"source_type":      "group",
						"destination":      acctest.CHECKSET,
						"destination_type": "group",
						"proto":            "TCP",
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
					"source": "${alicloud_cloud_firewall_address_book.ip_update.group_name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"source":      "127.0.0.2/32",
					"source_type": "net",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source":      "127.0.0.2/32",
						"source_type": "net",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"source": "127.0.0.1/32",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source": "127.0.0.1/32",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"destination": "${alicloud_cloud_firewall_address_book.ip_update.group_name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"destination": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"destination":      "127.0.0.3/32",
					"destination_type": "net",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"destination":      "127.0.0.3/32",
						"destination_type": "net",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"destination": "127.0.0.2/32",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"destination": "127.0.0.2/32",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"destination":      "${alicloud_cloud_firewall_address_book.domain.group_name}",
					"destination_type": "group",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"destination":      acctest.CHECKSET,
						"destination_type": "group",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"dest_port_group": "${alicloud_cloud_firewall_address_book.port.group_name}",
					"dest_port_type":  "group",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dest_port_group": acctest.CHECKSET,
						"dest_port_type":  "group",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"dest_port":       "20/22",
					"dest_port_group": acctest.REMOVEKEY,
					"dest_port_type":  "port",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dest_port":       "20/22",
						"dest_port_group": acctest.REMOVEKEY,
						"dest_port_type":  "port",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"dest_port":       acctest.REMOVEKEY,
					"dest_port_group": "${alicloud_cloud_firewall_address_book.port_update.group_name}",
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
					"repeat_type": "None",
					"start_time":  "1716998400",
					"end_time":    "1717083000",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"repeat_type": "None",
						"start_time":  acctest.CHECKSET,
						"end_time":    acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"repeat_type":       "Daily",
					"repeat_start_time": "08:00",
					"repeat_end_time":   "10:00",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"repeat_type":       "Daily",
						"repeat_start_time": "08:00",
						"repeat_end_time":   "10:00",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"repeat_type": "Weekly",
					"repeat_days": []string{"1", "2", "3"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"repeat_type":   "Weekly",
						"repeat_days.#": "3",
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
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"source_ip", "lang"},
			},
		},
	})
}

func TestAccAliCloudCloudFirewallControlPolicy_basic1_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_firewall_control_policy.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudCloudFirewallControlPolicyMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cloudfirewall.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCloudFirewallControlPolicy")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scloudfirewallcontrolpolicy%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudCloudFirewallControlPolicyBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"direction":           "out",
					"application_name":    "ANY",
					"description":         name,
					"acl_action":          "accept",
					"source":              "${alicloud_cloud_firewall_address_book.ip.group_name}",
					"source_type":         "group",
					"destination":         "${alicloud_cloud_firewall_address_book.domain.group_name}",
					"destination_type":    "group",
					"proto":               "TCP",
					"dest_port_group":     "${alicloud_cloud_firewall_address_book.port.group_name}",
					"dest_port_type":      "group",
					"ip_version":          "4",
					"domain_resolve_type": "DNS",
					"repeat_type":         "Weekly",
					"start_time":          "1716998400",
					"end_time":            "1717083000",
					"repeat_start_time":   "08:00",
					"repeat_end_time":     "10:00",
					"repeat_days":         []string{"1", "2", "3"},
					"release":             "false",
					"source_ip":           "127.0.0.1",
					"lang":                "en",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"direction":           "out",
						"application_name":    "ANY",
						"description":         name,
						"acl_action":          "accept",
						"source":              acctest.CHECKSET,
						"source_type":         "group",
						"destination":         acctest.CHECKSET,
						"destination_type":    "group",
						"proto":               "TCP",
						"dest_port_group":     acctest.CHECKSET,
						"dest_port_type":      "group",
						"ip_version":          "4",
						"domain_resolve_type": "DNS",
						"repeat_type":         "Weekly",
						"start_time":          acctest.CHECKSET,
						"end_time":            acctest.CHECKSET,
						"repeat_start_time":   "08:00",
						"repeat_end_time":     "10:00",
						"repeat_days.#":       "3",
						"release":             "false",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"source_ip", "lang"},
			},
		},
	})
}

func TestAccAliCloudCloudFirewallControlPolicy_basic2(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_firewall_control_policy.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudCloudFirewallControlPolicyMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cloudfirewall.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCloudFirewallControlPolicy")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scloudfirewallcontrolpolicy%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudCloudFirewallControlPolicyBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"direction":             "in",
					"application_name_list": []string{"ANY"},
					"description":           name,
					"acl_action":            "accept",
					"source":                "127.0.0.1/32",
					"source_type":           "net",
					"destination":           "127.0.0.2/32",
					"destination_type":      "net",
					"proto":                 "TCP",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"direction":               "in",
						"application_name_list.#": "1",
						"description":             name,
						"acl_action":              "accept",
						"source":                  "127.0.0.1/32",
						"source_type":             "net",
						"destination":             "127.0.0.2/32",
						"destination_type":        "net",
						"proto":                   "TCP",
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
				Config: testAccConfig(map[string]interface{}{
					"source": "${alicloud_cloud_firewall_address_book.ip_update.group_name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source": acctest.CHECKSET,
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
					"destination":      "${alicloud_cloud_firewall_address_book.ip.group_name}",
					"destination_type": "group",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"destination":      acctest.CHECKSET,
						"destination_type": "group",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"destination": "${alicloud_cloud_firewall_address_book.ip_update.group_name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"destination": acctest.CHECKSET,
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
					"proto": "ANY",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"proto": "ANY",
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
					"dest_port":       "20/22",
					"dest_port_group": acctest.REMOVEKEY,
					"dest_port_type":  "port",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dest_port":       "20/22",
						"dest_port_group": acctest.REMOVEKEY,
						"dest_port_type":  "port",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"repeat_type": "None",
					"start_time":  "1716998400",
					"end_time":    "1717083000",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"repeat_type": "None",
						"start_time":  acctest.CHECKSET,
						"end_time":    acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"repeat_type":       "Daily",
					"repeat_start_time": "08:00",
					"repeat_end_time":   "10:00",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"repeat_type":       "Daily",
						"repeat_start_time": "08:00",
						"repeat_end_time":   "10:00",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"repeat_type": "Weekly",
					"repeat_days": []string{"1", "2", "3"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"repeat_type":   "Weekly",
						"repeat_days.#": "3",
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
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"source_ip", "lang"},
			},
		},
	})
}

func TestAccAliCloudCloudFirewallControlPolicy_basic2_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_firewall_control_policy.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudCloudFirewallControlPolicyMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cloudfirewall.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCloudFirewallControlPolicy")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scloudfirewallcontrolpolicy%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudCloudFirewallControlPolicyBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"direction":             "in",
					"application_name_list": []string{"ANY"},
					"description":           name,
					"acl_action":            "accept",
					"source":                "::1/128",
					"source_type":           "net",
					"destination":           "::2/128",
					"destination_type":      "net",
					"proto":                 "TCP",
					"dest_port":             "20/22",
					"dest_port_type":        "port",
					"ip_version":            "6",
					"repeat_type":           "Weekly",
					"start_time":            "1716998400",
					"end_time":              "1717083000",
					"repeat_start_time":     "08:00",
					"repeat_end_time":       "10:00",
					"repeat_days":           []string{"1", "2", "3"},
					"release":               "false",
					"source_ip":             "127.0.0.1",
					"lang":                  "zh",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"direction":               "in",
						"application_name_list.#": "1",
						"description":             name,
						"acl_action":              "accept",
						"source":                  "::1/128",
						"source_type":             "net",
						"destination":             "::2/128",
						"destination_type":        "net",
						"proto":                   "TCP",
						"dest_port":               "20/22",
						"dest_port_type":          "port",
						"ip_version":              "6",
						"repeat_type":             "Weekly",
						"start_time":              acctest.CHECKSET,
						"end_time":                acctest.CHECKSET,
						"repeat_start_time":       "08:00",
						"repeat_end_time":         "10:00",
						"repeat_days.#":           "3",
						"release":                 "false",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"source_ip", "lang"},
			},
		},
	})
}

func TestAccAliCloudCloudFirewallControlPolicy_basic3(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_firewall_control_policy.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudCloudFirewallControlPolicyMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cloudfirewall.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCloudFirewallControlPolicy")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scloudfirewallcontrolpolicy%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudCloudFirewallControlPolicyBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"direction":             "out",
					"application_name_list": []string{"ANY"},
					"description":           name,
					"acl_action":            "accept",
					"source":                "${alicloud_cloud_firewall_address_book.ip.group_name}",
					"source_type":           "group",
					"destination":           "${alicloud_cloud_firewall_address_book.ip.group_name}",
					"destination_type":      "group",
					"proto":                 "TCP",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"direction":               "out",
						"application_name_list.#": "1",
						"description":             name,
						"acl_action":              "accept",
						"source":                  acctest.CHECKSET,
						"source_type":             "group",
						"destination":             acctest.CHECKSET,
						"destination_type":        "group",
						"proto":                   "TCP",
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
					"source": "${alicloud_cloud_firewall_address_book.ip_update.group_name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"source":      "127.0.0.2/32",
					"source_type": "net",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source":      "127.0.0.2/32",
						"source_type": "net",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"source": "127.0.0.1/32",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source": "127.0.0.1/32",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"destination": "${alicloud_cloud_firewall_address_book.ip_update.group_name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"destination": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"destination":      "127.0.0.3/32",
					"destination_type": "net",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"destination":      "127.0.0.3/32",
						"destination_type": "net",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"destination": "127.0.0.2/32",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"destination": "127.0.0.2/32",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"destination":      "${alicloud_cloud_firewall_address_book.domain.group_name}",
					"destination_type": "group",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"destination":      acctest.CHECKSET,
						"destination_type": "group",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"dest_port_group": "${alicloud_cloud_firewall_address_book.port.group_name}",
					"dest_port_type":  "group",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dest_port_group": acctest.CHECKSET,
						"dest_port_type":  "group",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"dest_port":       "20/22",
					"dest_port_group": acctest.REMOVEKEY,
					"dest_port_type":  "port",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dest_port":       "20/22",
						"dest_port_group": acctest.REMOVEKEY,
						"dest_port_type":  "port",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"dest_port":       acctest.REMOVEKEY,
					"dest_port_group": "${alicloud_cloud_firewall_address_book.port_update.group_name}",
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
					"repeat_type": "None",
					"start_time":  "1716998400",
					"end_time":    "1717083000",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"repeat_type": "None",
						"start_time":  acctest.CHECKSET,
						"end_time":    acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"repeat_type":       "Daily",
					"repeat_start_time": "08:00",
					"repeat_end_time":   "10:00",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"repeat_type":       "Daily",
						"repeat_start_time": "08:00",
						"repeat_end_time":   "10:00",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"repeat_type": "Weekly",
					"repeat_days": []string{"1", "2", "3"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"repeat_type":   "Weekly",
						"repeat_days.#": "3",
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
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"source_ip", "lang"},
			},
		},
	})
}

func TestAccAliCloudCloudFirewallControlPolicy_basic3_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_firewall_control_policy.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudCloudFirewallControlPolicyMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cloudfirewall.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCloudFirewallControlPolicy")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scloudfirewallcontrolpolicy%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudCloudFirewallControlPolicyBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"direction":             "out",
					"application_name_list": []string{"ANY"},
					"description":           name,
					"acl_action":            "accept",
					"source":                "${alicloud_cloud_firewall_address_book.ip.group_name}",
					"source_type":           "group",
					"destination":           "${alicloud_cloud_firewall_address_book.domain.group_name}",
					"destination_type":      "group",
					"proto":                 "TCP",
					"dest_port_group":       "${alicloud_cloud_firewall_address_book.port.group_name}",
					"dest_port_type":        "group",
					"ip_version":            "4",
					"domain_resolve_type":   "DNS",
					"repeat_type":           "Weekly",
					"start_time":            "1716998400",
					"end_time":              "1717083000",
					"repeat_start_time":     "08:00",
					"repeat_end_time":       "10:00",
					"repeat_days":           []string{"1", "2", "3"},
					"release":               "false",
					"source_ip":             "127.0.0.1",
					"lang":                  "en",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"direction":               "out",
						"application_name_list.#": "1",
						"description":             name,
						"acl_action":              "accept",
						"source":                  acctest.CHECKSET,
						"source_type":             "group",
						"destination":             acctest.CHECKSET,
						"destination_type":        "group",
						"proto":                   "TCP",
						"dest_port_group":         acctest.CHECKSET,
						"dest_port_type":          "group",
						"ip_version":              "4",
						"domain_resolve_type":     "DNS",
						"repeat_type":             "Weekly",
						"start_time":              acctest.CHECKSET,
						"end_time":                acctest.CHECKSET,
						"repeat_start_time":       "08:00",
						"repeat_end_time":         "10:00",
						"repeat_days.#":           "3",
						"release":                 "false",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"source_ip", "lang"},
			},
		},
	})
}

var AliCloudCloudFirewallControlPolicyMap0 = map[string]string{
	"dest_port":      acctest.CHECKSET,
	"dest_port_type": acctest.CHECKSET,
	"ip_version":     acctest.CHECKSET,
	"repeat_type":    acctest.CHECKSET,
	"release":        acctest.CHECKSET,
	"acl_uuid":       acctest.CHECKSET,
	"create_time":    acctest.CHECKSET,
}

func AliCloudCloudFirewallControlPolicyBasicDependence0(name string) string {
	return fmt.Sprintf(` 
	variable "name" {
  		default = "%s"
	}

	resource "alicloud_cloud_firewall_address_book" "ip" {
  		group_name   = var.name
  		group_type   = "ip"
  		description  = var.name
  		address_list = ["10.21.0.0/16", "10.168.0.0/16"]
	}

	resource "alicloud_cloud_firewall_address_book" "ip_update" {
  		group_name   = "${var.name}-update"
  		group_type   = "ip"
  		description  = "${var.name}-update"
  		address_list = ["10.21.0.0/16", "10.22.0.0/16", "10.168.0.0/16"]
	}

	resource "alicloud_cloud_firewall_address_book" "port" {
  		group_name   = var.name
  		group_type   = "port"
  		description  = var.name
  		address_list = ["1/1", "22/22", "88/88"]
	}

	resource "alicloud_cloud_firewall_address_book" "port_update" {
  		group_name   = "${var.name}-update"
  		group_type   = "port"
  		description  = "${var.name}-update"
  		address_list = ["22/22", "88/88"]
	}

	resource "alicloud_cloud_firewall_address_book" "domain" {
  		group_name   = "${var.name}-domain"
  		group_type   = "domain"
  		description  = "${var.name}-domain"
  		address_list = ["alibaba.com", "aliyun.com", "alicloud.com"]
	}
`, name)
}

// TestAccAliCloudCloudFirewallControlPolicyApplicationNameListOrder verifies that
// reordering application_name_list (a plain TypeList of TypeString with no
// DiffSuppressFunc) converges to an empty post-apply plan.
func TestAccAliCloudCloudFirewallControlPolicyApplicationNameListOrder(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_firewall_control_policy.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudCloudFirewallControlPolicyMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cloudfirewall.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCloudFirewallControlPolicy")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	name := fmt.Sprintf("tf-testacc%scloudfirewallcontrolpolicy%d", acctest.DefaultRegionToTest, tfacctest.RandIntRange(10000, 99999))
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudCloudFirewallControlPolicyBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"direction":             "in",
					"description":           name,
					"acl_action":            "accept",
					"source":                "127.0.0.1/32",
					"source_type":           "net",
					"destination":           "127.0.0.2/32",
					"destination_type":      "net",
					"proto":                 "TCP",
					"application_name_list": []string{"HTTP", "HTTPS"},
				}),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"application_name_list": []string{"HTTPS", "HTTP"},
				}),
			},
		},
	})
}

// TestAccAliCloudCloudFirewallControlPolicyRepeatDaysOrder verifies that
// reordering repeat_days (a plain TypeList of TypeInt with no DiffSuppressFunc)
// converges to an empty post-apply plan.
func TestAccAliCloudCloudFirewallControlPolicyRepeatDaysOrder(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_firewall_control_policy.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudCloudFirewallControlPolicyMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cloudfirewall.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCloudFirewallControlPolicy")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	name := fmt.Sprintf("tf-testacc%scloudfirewallcontrolpolicy%d", acctest.DefaultRegionToTest, tfacctest.RandIntRange(10000, 99999))
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudCloudFirewallControlPolicyBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"direction":        "in",
					"description":      name,
					"acl_action":       "accept",
					"source":           "127.0.0.1/32",
					"source_type":      "net",
					"destination":      "127.0.0.2/32",
					"destination_type": "net",
					"proto":            "TCP",
					"repeat_type":      "Weekly",
					"repeat_days":      []string{"1", "2"},
				}),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"repeat_days": []string{"2", "1"},
				}),
			},
		},
	})
}
