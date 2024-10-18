package alicloud

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/stretchr/testify/assert"

	"github.com/alibabacloud-go/tea-rpc/client"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudCloudFirewallControlPolicy_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_firewall_control_policy.default"
	ra := resourceAttrInit(resourceId, AliCloudCloudFirewallControlPolicyMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudfwService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudFirewallControlPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scloudfirewallcontrolpolicy%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCloudFirewallControlPolicyBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
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
						"source":      CHECKSET,
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
						"source": CHECKSET,
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
						"destination":      CHECKSET,
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
						"destination": CHECKSET,
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
					"dest_port":       REMOVEKEY,
					"dest_port_group": "${alicloud_cloud_firewall_address_book.port.group_name}",
					"dest_port_type":  "group",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dest_port":       REMOVEKEY,
						"dest_port_group": CHECKSET,
						"dest_port_type":  "group",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"dest_port":       "20/22",
					"dest_port_group": REMOVEKEY,
					"dest_port_type":  "port",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dest_port":       "20/22",
						"dest_port_group": REMOVEKEY,
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
						"start_time":  CHECKSET,
						"end_time":    CHECKSET,
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
	ra := resourceAttrInit(resourceId, AliCloudCloudFirewallControlPolicyMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudfwService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudFirewallControlPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scloudfirewallcontrolpolicy%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCloudFirewallControlPolicyBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
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
						"start_time":        CHECKSET,
						"end_time":          CHECKSET,
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
	ra := resourceAttrInit(resourceId, AliCloudCloudFirewallControlPolicyMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudfwService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudFirewallControlPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scloudfirewallcontrolpolicy%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCloudFirewallControlPolicyBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
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
						"source":           CHECKSET,
						"source_type":      "group",
						"destination":      CHECKSET,
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
						"source": CHECKSET,
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
						"destination": CHECKSET,
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
						"destination":      CHECKSET,
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
						"dest_port_group": CHECKSET,
						"dest_port_type":  "group",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"dest_port":       "20/22",
					"dest_port_group": REMOVEKEY,
					"dest_port_type":  "port",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dest_port":       "20/22",
						"dest_port_group": REMOVEKEY,
						"dest_port_type":  "port",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"dest_port":       REMOVEKEY,
					"dest_port_group": "${alicloud_cloud_firewall_address_book.port_update.group_name}",
					"dest_port_type":  "group",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dest_port":       REMOVEKEY,
						"dest_port_group": CHECKSET,
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
						"start_time":  CHECKSET,
						"end_time":    CHECKSET,
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
	ra := resourceAttrInit(resourceId, AliCloudCloudFirewallControlPolicyMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudfwService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudFirewallControlPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scloudfirewallcontrolpolicy%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCloudFirewallControlPolicyBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
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
						"source":              CHECKSET,
						"source_type":         "group",
						"destination":         CHECKSET,
						"destination_type":    "group",
						"proto":               "TCP",
						"dest_port_group":     CHECKSET,
						"dest_port_type":      "group",
						"ip_version":          "4",
						"domain_resolve_type": "DNS",
						"repeat_type":         "Weekly",
						"start_time":          CHECKSET,
						"end_time":            CHECKSET,
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
	ra := resourceAttrInit(resourceId, AliCloudCloudFirewallControlPolicyMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudfwService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudFirewallControlPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scloudfirewallcontrolpolicy%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCloudFirewallControlPolicyBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
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
						"source":      CHECKSET,
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
						"source": CHECKSET,
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
						"destination":      CHECKSET,
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
						"destination": CHECKSET,
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
					"dest_port":       REMOVEKEY,
					"dest_port_group": "${alicloud_cloud_firewall_address_book.port.group_name}",
					"dest_port_type":  "group",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dest_port":       REMOVEKEY,
						"dest_port_group": CHECKSET,
						"dest_port_type":  "group",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"dest_port":       "20/22",
					"dest_port_group": REMOVEKEY,
					"dest_port_type":  "port",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dest_port":       "20/22",
						"dest_port_group": REMOVEKEY,
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
						"start_time":  CHECKSET,
						"end_time":    CHECKSET,
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
	ra := resourceAttrInit(resourceId, AliCloudCloudFirewallControlPolicyMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudfwService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudFirewallControlPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scloudfirewallcontrolpolicy%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCloudFirewallControlPolicyBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
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
						"start_time":              CHECKSET,
						"end_time":                CHECKSET,
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
	ra := resourceAttrInit(resourceId, AliCloudCloudFirewallControlPolicyMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudfwService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudFirewallControlPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scloudfirewallcontrolpolicy%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCloudFirewallControlPolicyBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
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
						"source":                  CHECKSET,
						"source_type":             "group",
						"destination":             CHECKSET,
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
						"source": CHECKSET,
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
						"destination": CHECKSET,
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
						"destination":      CHECKSET,
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
						"dest_port_group": CHECKSET,
						"dest_port_type":  "group",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"dest_port":       "20/22",
					"dest_port_group": REMOVEKEY,
					"dest_port_type":  "port",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dest_port":       "20/22",
						"dest_port_group": REMOVEKEY,
						"dest_port_type":  "port",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"dest_port":       REMOVEKEY,
					"dest_port_group": "${alicloud_cloud_firewall_address_book.port_update.group_name}",
					"dest_port_type":  "group",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dest_port":       REMOVEKEY,
						"dest_port_group": CHECKSET,
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
						"start_time":  CHECKSET,
						"end_time":    CHECKSET,
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
	ra := resourceAttrInit(resourceId, AliCloudCloudFirewallControlPolicyMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudfwService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudFirewallControlPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scloudfirewallcontrolpolicy%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCloudFirewallControlPolicyBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
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
						"source":                  CHECKSET,
						"source_type":             "group",
						"destination":             CHECKSET,
						"destination_type":        "group",
						"proto":                   "TCP",
						"dest_port_group":         CHECKSET,
						"dest_port_type":          "group",
						"ip_version":              "4",
						"domain_resolve_type":     "DNS",
						"repeat_type":             "Weekly",
						"start_time":              CHECKSET,
						"end_time":                CHECKSET,
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
	"dest_port":      CHECKSET,
	"dest_port_type": CHECKSET,
	"ip_version":     CHECKSET,
	"repeat_type":    CHECKSET,
	"release":        CHECKSET,
	"acl_uuid":       CHECKSET,
	"create_time":    CHECKSET,
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

func TestUnitAliCloudCloudFirewallControlPolicy(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	d, _ := schema.InternalMap(p["alicloud_cloud_firewall_control_policy"].Schema).Data(nil, nil)
	dCreate, _ := schema.InternalMap(p["alicloud_cloud_firewall_control_policy"].Schema).Data(nil, nil)
	dCreate.MarkNewResource()
	for key, value := range map[string]interface{}{
		"application_name": "ANY",
		"acl_action":       "accept",
		"dest_port":        "0/0",
		"description":      "description",
		"destination_type": "net",
		"destination":      "100.1.1.0/24",
		"direction":        "direction",
		"proto":            "ANY",
		"source":           "1.2.3.0/24",
		"source_type":      "net",
		"dest_port_group":  "group",
		"dest_port_type":   "port",
		"ip_version":       "ip_version",
		"lang":             "lang",
		"source_ip":        "source_ip",
	} {
		err := dCreate.Set(key, value)
		assert.Nil(t, err)
		err = d.Set(key, value)
		assert.Nil(t, err)
	}
	region := os.Getenv("ALICLOUD_REGION")
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		t.Skipf("Skipping the test case with err: %s", err)
		t.Skipped()
	}
	rawClient = rawClient.(*connectivity.AliyunClient)
	ReadMockResponse := map[string]interface{}{
		"Policys": []interface{}{
			map[string]interface{}{
				"AclUuid":         "MockAclUuid",
				"Direction":       "direction",
				"AclAction":       "accept",
				"ApplicationName": "ANY",
				"Description":     "description",
				"DestPort":        "dest_port",
				"DestPortGroup":   "ANY",
				"DestPortType":    "group",
				"Destination":     "100.1.1.0/24",
				"DestinationType": "net",
				"Proto":           "proto",
				"Release":         "release",
				"Source":          "source",
				"SourceType":      "source_type",
			},
		},
	}

	responseMock := map[string]func(errorCode string) (map[string]interface{}, error){
		"RetryError": func(errorCode string) (map[string]interface{}, error) {
			return nil, &tea.SDKError{
				Code:       String(errorCode),
				Data:       String(errorCode),
				Message:    String(errorCode),
				StatusCode: tea.Int(400),
			}
		},
		"NotFoundError": func(errorCode string) (map[string]interface{}, error) {
			return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_cloud_firewall_control_policy", "MockAclUuid"))
		},
		"NoRetryError": func(errorCode string) (map[string]interface{}, error) {
			return nil, &tea.SDKError{
				Code:       String(errorCode),
				Data:       String(errorCode),
				Message:    String(errorCode),
				StatusCode: tea.Int(400),
			}
		},
		"CreateNormal": func(errorCode string) (map[string]interface{}, error) {
			result := ReadMockResponse
			result["AclUuid"] = "MockAclUuid"
			return result, nil
		},
		"UpdateNormal": func(errorCode string) (map[string]interface{}, error) {
			result := ReadMockResponse
			return result, nil
		},
		"DeleteNormal": func(errorCode string) (map[string]interface{}, error) {
			result := ReadMockResponse
			return result, nil
		},
		"ReadNormal": func(errorCode string) (map[string]interface{}, error) {
			result := ReadMockResponse
			return result, nil
		},
	}
	// Create
	t.Run("CreateClientAbnormal", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewCloudfwClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:       String("loadEndpoint error"),
				Data:       String("loadEndpoint error"),
				Message:    String("loadEndpoint error"),
				StatusCode: tea.Int(400),
			}
		})
		err := resourceAliCloudCloudFirewallControlPolicyCreate(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("CreateAbnormal", func(t *testing.T) {
		retryFlag := true
		noRetryFlag := true
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["CreateNormal"]("")
		})
		err := resourceAliCloudCloudFirewallControlPolicyCreate(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("CreateNormal", func(t *testing.T) {
		retryFlag := false
		noRetryFlag := false
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["CreateNormal"]("")
		})
		err := resourceAliCloudCloudFirewallControlPolicyCreate(dCreate, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})

	// Set ID for Update and Delete Method
	d.SetId(fmt.Sprint("MockAclUuid", ":", "direction"))
	// Update
	t.Run("UpdateClientAbnormal", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewCloudfwClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:       String("loadEndpoint error"),
				Data:       String("loadEndpoint error"),
				Message:    String("loadEndpoint error"),
				StatusCode: tea.Int(400),
			}
		})

		err := resourceAliCloudCloudFirewallControlPolicyUpdate(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})

	t.Run("UpdateModifyControlPolicyAbnormal", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		for _, key := range []string{"acl_action", "application_name", "application_name", "description", "destination", "destination_type", "proto", "source", "source_type", "dest_port", "dest_port_group", "dest_port_type", "lang", "release", "source_ip"} {
			switch p["alicloud_cloud_firewall_control_policy"].Schema[key].Type {
			case schema.TypeString:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: d.Get(key).(string), New: d.Get(key).(string) + "_update"})
			case schema.TypeBool:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: strconv.FormatBool(d.Get(key).(bool)), New: strconv.FormatBool(true)})
			case schema.TypeInt:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: strconv.Itoa(d.Get(key).(int)), New: strconv.Itoa(3)})
			case schema.TypeMap:
				diff.SetAttribute("tags.%", &terraform.ResourceAttrDiff{Old: "0", New: "2"})
				diff.SetAttribute("tags.For", &terraform.ResourceAttrDiff{Old: "", New: "Test"})
				diff.SetAttribute("tags.Created", &terraform.ResourceAttrDiff{Old: "", New: "TF"})
			}
		}
		resourceData1, _ := schema.InternalMap(p["alicloud_cloud_firewall_control_policy"].Schema).Data(nil, diff)
		resourceData1.SetId(d.Id())
		retryFlag := true
		noRetryFlag := true
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["UpdateNormal"]("")
		})
		err := resourceAliCloudCloudFirewallControlPolicyUpdate(resourceData1, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})

	t.Run("UpdateModifyControlPolicyNormal", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		for _, key := range []string{"acl_action", "application_name", "application_name", "description", "destination", "destination_type", "proto", "source", "source_type", "dest_port", "dest_port_group", "dest_port_type", "lang", "release", "source_ip"} {
			switch p["alicloud_cloud_firewall_control_policy"].Schema[key].Type {
			case schema.TypeString:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: d.Get(key).(string), New: d.Get(key).(string) + "_update"})
			case schema.TypeBool:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: strconv.FormatBool(d.Get(key).(bool)), New: strconv.FormatBool(true)})
			case schema.TypeInt:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: strconv.Itoa(d.Get(key).(int)), New: strconv.Itoa(3)})
			case schema.TypeMap:
				diff.SetAttribute("tags.%", &terraform.ResourceAttrDiff{Old: "0", New: "2"})
				diff.SetAttribute("tags.For", &terraform.ResourceAttrDiff{Old: "", New: "Test"})
				diff.SetAttribute("tags.Created", &terraform.ResourceAttrDiff{Old: "", New: "TF"})
			}
		}
		resourceData1, _ := schema.InternalMap(p["alicloud_cloud_firewall_control_policy"].Schema).Data(nil, diff)
		resourceData1.SetId(d.Id())
		retryFlag := false
		noRetryFlag := false
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["UpdateNormal"]("")
		})
		err := resourceAliCloudCloudFirewallControlPolicyUpdate(resourceData1, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})

	t.Run("UpdateMockParseResourceId", func(t *testing.T) {
		resourceData1, _ := schema.InternalMap(p["alicloud_cloud_firewall_control_policy"].Schema).Data(nil, nil)
		resourceData1.SetId("MockId")
		retryFlag := false
		noRetryFlag := false
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("RetryError")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["UpdateNormal"]("")
		})
		err := resourceAliCloudCloudFirewallControlPolicyUpdate(resourceData1, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})

	// Delete
	t.Run("DeleteClientAbnormal", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewCloudfwClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:       String("loadEndpoint error"),
				Data:       String("loadEndpoint error"),
				Message:    String("loadEndpoint error"),
				StatusCode: tea.Int(400),
			}
		})
		err := resourceAliCloudCloudFirewallControlPolicyDelete(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("DeleteMockAbnormal", func(t *testing.T) {
		retryFlag := true
		noRetryFlag := true
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["DeleteNormal"]("")
		})
		err := resourceAliCloudCloudFirewallControlPolicyDelete(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("DeleteMockNormal", func(t *testing.T) {
		retryFlag := false
		noRetryFlag := false
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["DeleteNormal"]("")
		})
		patcheDescribeCenInstance := gomonkey.ApplyMethod(reflect.TypeOf(&CbnService{}), "DescribeCenInstance", func(*CbnService, string) (map[string]interface{}, error) {
			return responseMock["NotFoundError"]("ResourceNotfound")
		})
		err := resourceAliCloudCloudFirewallControlPolicyDelete(d, rawClient)
		patches.Reset()
		patcheDescribeCenInstance.Reset()
		assert.Nil(t, err)
	})

	t.Run("DeleteMockParseResourceId", func(t *testing.T) {
		resourceData1, _ := schema.InternalMap(p["alicloud_cloud_firewall_control_policy"].Schema).Data(nil, nil)
		resourceData1.SetId("MockId")
		retryFlag := false
		noRetryFlag := false
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("RetryError")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["DeleteNormal"]("")
		})
		err := resourceAliCloudCloudFirewallControlPolicyDelete(resourceData1, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})

	//Read
	t.Run("ReadDescribeCloudFirewallControlPolicyNotFound", func(t *testing.T) {
		patcheDorequest := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			NotFoundFlag := true
			noRetryFlag := false
			if NotFoundFlag {
				return responseMock["NotFoundError"]("ResourceNotfound")
			} else if noRetryFlag {
				return responseMock["NoRetryError"]("NoRetryError")
			}
			return responseMock["ReadNormal"]("")
		})
		err := resourceAliCloudCloudFirewallControlPolicyRead(d, rawClient)
		patcheDorequest.Reset()
		assert.Nil(t, err)
	})

	t.Run("ReadDescribeCloudFirewallControlPolicyAbnormal", func(t *testing.T) {
		patcheDorequest := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			retryFlag := false
			noRetryFlag := true
			if retryFlag {
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["ReadNormal"]("")
		})
		err := resourceAliCloudCloudFirewallControlPolicyRead(d, rawClient)
		patcheDorequest.Reset()
		assert.NotNil(t, err)
	})
}
