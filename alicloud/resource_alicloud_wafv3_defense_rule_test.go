// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Case DefenseRule__longitude_header_test 12001
func TestAccAliCloudWafv3DefenseRule_longitude12001(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_wafv3_defense_rule.this"
	ra := resourceAttrInit(resourceId, AlicloudWafv3DefenseRuleMapLongitude12001)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &Wafv3ServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeWafv3DefenseRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccwafv3%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudWafv3DefenseRuleBasicDependenceLongitude12001)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":    "${data.alicloud_wafv3_instances.default.ids.0}",
					"defense_origin": "custom",
					"defense_scene":  "custom_acl",
					"defense_type":   "template",
					"rule_name":      "longitude",
					"rule_status":    "0",
					"template_id":    "${alicloud_wafv3_defense_template.defaultfIoHt5-hf.defense_template_id}",
					"config": []map[string]interface{}{
						{
							"mode":        "0",
							"cc_effect":   "service",
							"cc_status":   "0",
							"rule_action": "monitor",
							"conditions": []map[string]interface{}{
								{
									"key":      "Header",
									"op_value": "none",
									"sub_key":  "alicdn-viewer-longitude",
								},
							},
							"rate_limit": []map[string]interface{}{
								{
									"interval":  "0",
									"threshold": "0",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":    CHECKSET,
						"defense_origin": "custom",
						"defense_scene":  "custom_acl",
						"defense_type":   "template",
						"rule_name":      "longitude",
						"rule_status":    "0",
						"template_id":    CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rule_status": "1",
					"rule_name":   "longitude-updated",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule_status": "1",
						"rule_name":   "longitude-updated",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rule_status": "0",
					"rule_name":   "longitude-final",
					"config": []map[string]interface{}{
						{
							"cc_effect":   "service",
							"cc_status":   "0",
							"rule_action": "js",
							"conditions": []map[string]interface{}{
								{
									"key":      "Header",
									"op_value": "none",
									"sub_key":  "alicdn-viewer-longitude",
								},
							},
							"rate_limit": []map[string]interface{}{
								{
									"interval":  "0",
									"threshold": "0",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule_status": "0",
						"rule_name":   "longitude-final",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudWafv3DefenseRuleMapLongitude12001 = map[string]string{
	"rule_id": CHECKSET,
}

func AlicloudWafv3DefenseRuleBasicDependenceLongitude12001(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "region_id" {
  default = "cn-hangzhou"
}

data "alicloud_wafv3_instances" "default" {
}

resource "alicloud_wafv3_defense_template" "defaultfIoHt5-hf" {
  instance_id           = data.alicloud_wafv3_instances.default.ids.0
  template_origin       = "custom"
  defense_template_name = "1754448878"
  defense_scene         = "custom_acl"
  template_type         = "user_custom"
  status                = "1"
  description           = "testCreate"
}

`, name)
}

// Test Wafv3 DefenseRule. >>> Resource test cases, automatically generated.
// Case  DefenseRule-20250715_resource 11029
func TestAccAliCloudWafv3DefenseRule_basic11029(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_wafv3_defense_rule.default"
	ra := resourceAttrInit(resourceId, AlicloudWafv3DefenseRuleMap11029)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &Wafv3ServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeWafv3DefenseRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccwafv3%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudWafv3DefenseRuleBasicDependence11029)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":    "${data.alicloud_wafv3_instances.default.ids.0}",
					"defense_type":   "resource",
					"defense_scene":  "account_identifier",
					"rule_status":    "1",
					"resource":       "${alicloud_wafv3_domain.defaultICMRhk.domain_id}",
					"defense_origin": "custom",
					"config": []map[string]interface{}{
						{
							"account_identifiers": []map[string]interface{}{
								{
									"position":    "jwt",
									"priority":    "2",
									"decode_type": "jwt",
									"key":         "Query-Arg",
									"sub_key":     "adb",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":    CHECKSET,
						"defense_type":   "resource",
						"defense_scene":  "account_identifier",
						"rule_status":    "1",
						"resource":       CHECKSET,
						"defense_origin": "custom",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"config": []map[string]interface{}{
						{
							"account_identifiers": []map[string]interface{}{
								{
									"priority":    "18",
									"decode_type": "basic",
									"key":         "Header",
									"sub_key":     "asdsd",
								},
								{
									"priority":    "12",
									"decode_type": "plain",
									"key":         "Post-Arg",
									"sub_key":     "22222",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudWafv3DefenseRuleMap11029 = map[string]string{
	"rule_id": CHECKSET,
}

func AlicloudWafv3DefenseRuleBasicDependence11029(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "region_id" {
  default = "cn-hangzhou"
}

data "alicloud_wafv3_instances" "default" {
}

resource "alicloud_wafv3_domain" "defaultICMRhk" {
  redirect {
    loadbalance = "iphash"
    backends    = ["39.98.217.197"]
    connect_timeout = 5
    read_timeout    = 120
    write_timeout   = 120
  }
  domain      = "testfromtftest01.wafqax.top"
  access_type = "share"
  instance_id = data.alicloud_wafv3_instances.default.ids.0
  listen {
    http_ports = ["80"]
  }
}


`, name)
}

// Case  DefenseRule__20250715__自定义规则 11017
func TestAccAliCloudWafv3DefenseRule_basic11017(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_wafv3_defense_rule.default"
	ra := resourceAttrInit(resourceId, AlicloudWafv3DefenseRuleMap11017)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &Wafv3ServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeWafv3DefenseRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccwafv3%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudWafv3DefenseRuleBasicDependence11017)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"defense_origin": "custom",
					"instance_id":    "${data.alicloud_wafv3_instances.default.ids.0}",
					"config": []map[string]interface{}{
						{
							"rule_action": "block",
							"conditions": []map[string]interface{}{
								{
									"op_value": "contain",
									"values":   "abc",
									"key":      "URL",
								},
								{
									"op_value": "contain",
									"values":   "abc",
									"key":      "URLPath",
								},
								{
									"op_value": "contain",
									"values":   "1.1.1.2",
									"key":      "IP",
								},
							},
							"cc_status": "1",
							"cc_effect": "service",
							"rate_limit": []map[string]interface{}{
								{
									"target":    "remote_addr",
									"interval":  "16",
									"threshold": "204",
									"ttl":       "68",
									"status": []map[string]interface{}{
										{
											"code":  "414",
											"count": "333",
										},
									},
									"sub_key": "testky1",
								},
							},
						},
					},
					"defense_scene": "custom_acl",
					"rule_status":   "1",
					"defense_type":  "template",
					"template_id":   "${alicloud_wafv3_defense_template.defaultfIoHt5.defense_template_id}",
					"rule_name":     "custom_acl-create",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"defense_origin": "custom",
						"instance_id":    CHECKSET,
						"defense_scene":  "custom_acl",
						"rule_status":    "1",
						"defense_type":   "template",
						"template_id":    CHECKSET,
						"rule_name":      "custom_acl-create",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"config": []map[string]interface{}{
						{
							"rule_action": "monitor",
							"conditions": []map[string]interface{}{
								{
									"op_value": "contain",
									"values":   "abcd",
									"key":      "Header",
									"sub_key":  "testkey",
								},
							},
							"cc_status": "1",
							"cc_effect": "rule",
							"rate_limit": []map[string]interface{}{
								{
									"target":    "header",
									"interval":  "6",
									"threshold": "3",
									"ttl":       "61",
									"status": []map[string]interface{}{
										{
											"code":  "404",
											"ratio": "34",
										},
									},
									"sub_key": "abc",
								},
							},
						},
					},
					"rule_name": "testtt",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule_name": "testtt",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"config": []map[string]interface{}{
						{
							"rule_action": "js",
							"cc_status":   "0",
							"conditions": []map[string]interface{}{
								{
									"op_value": "contain",
									"values":   "123",
									"key":      "Header",
									"sub_key":  "test",
								},
							},
						},
					},
					"rule_status": "0",
					"rule_name":   "custom_acl_update1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule_status": "0",
						"rule_name":   "custom_acl_update1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudWafv3DefenseRuleMap11017 = map[string]string{
	"rule_id": CHECKSET,
}

func AlicloudWafv3DefenseRuleBasicDependence11017(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "region_id" {
  default = "cn-hangzhou"
}

data "alicloud_wafv3_instances" "default" {
}

resource "alicloud_wafv3_defense_template" "defaultfIoHt5" {
  instance_id           = data.alicloud_wafv3_instances.default.ids.0
  template_origin       = "custom"
  defense_template_name = "1758078439"
  defense_scene         = "custom_acl"
  template_type         = "user_custom"
  status                = "1"
  description           = "testCreate"
}


`, name)
}

// Case  DefenseRule__20250723—— 白名单_2 11097
func TestAccAliCloudWafv3DefenseRule_basic11097(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_wafv3_defense_rule.default"
	ra := resourceAttrInit(resourceId, AlicloudWafv3DefenseRuleMap11097)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &Wafv3ServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeWafv3DefenseRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccwafv3%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudWafv3DefenseRuleBasicDependence11097)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":    "${data.alicloud_wafv3_instances.default.ids.0}",
					"defense_type":   "template",
					"defense_scene":  "whitelist",
					"rule_status":    "1",
					"template_id":    "${alicloud_wafv3_defense_template.defaultZmPPmw.defense_template_id}",
					"rule_name":      "tf-whitelist",
					"defense_origin": "custom",
					"config": []map[string]interface{}{
						{
							"conditions": []map[string]interface{}{
								{
									"op_value": "contain",
									"values":   "abc",
									"key":      "URL",
									"sub_key":  "request_url",
								},
								{
									"op_value": "contain",
									"values":   "test",
									"key":      "URLPath",
									"sub_key":  "reqeust-url",
								},
								{
									"op_value": "eq",
									"values":   "2.2.2.2",
									"key":      "IP",
									"sub_key":  "requset-ip",
								},
							},
							"bypass_tags": []string{
								"customrule", "blacklist", "antiscan", "regular_rule"},
							"bypass_regular_types": []string{},
							"bypass_regular_rules": []string{
								"130068", "900928", "900814"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":    CHECKSET,
						"defense_type":   "template",
						"defense_scene":  "whitelist",
						"rule_status":    "1",
						"template_id":    CHECKSET,
						"rule_name":      "tf-whitelist",
						"defense_origin": "custom",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rule_status": "0",
					"rule_name":   "test-whietest-update",
					"config": []map[string]interface{}{
						{
							"conditions": []map[string]interface{}{
								{
									"op_value": "lt",
									"values":   "10",
									"key":      "Content-Length",
									"sub_key":  "requeset-content-length",
								},
							},
							"bypass_tags": []string{
								"cc", "regular_rule"},
							"bypass_regular_rules": []string{
								"130068", "900928"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule_status": "0",
						"rule_name":   "test-whietest-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rule_status": "1",
					"rule_name":   "whitelist-update-2",
					"config": []map[string]interface{}{
						{
							"conditions": []map[string]interface{}{
								{
									"op_value": "contain",
									"values":   "abc",
									"key":      "URL",
									"sub_key":  "testurl",
								},
							},
							"bypass_tags": []string{
								"regular_type"},
							"bypass_regular_types": []string{
								"sqli", "xss", "code_exec"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule_status": "1",
						"rule_name":   "whitelist-update-2",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudWafv3DefenseRuleMap11097 = map[string]string{
	"rule_id": CHECKSET,
}

func AlicloudWafv3DefenseRuleBasicDependence11097(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "region_id" {
  default = "cn-hangzhou"
}

data "alicloud_wafv3_instances" "default" {
}

resource "alicloud_wafv3_defense_template" "defaultZmPPmw" {
  instance_id           = data.alicloud_wafv3_instances.default.ids.0
  template_origin       = "custom"
  defense_template_name = "1758078467"
  defense_scene         = "whitelist"
  template_type         = "user_custom"
  status                = "1"
  description           = "testCreate"
}


`, name)
}

// Case  DefenseRule__20250722——IP黑名单 11070
func TestAccAliCloudWafv3DefenseRule_basic11070(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_wafv3_defense_rule.default"
	ra := resourceAttrInit(resourceId, AlicloudWafv3DefenseRuleMap11070)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &Wafv3ServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeWafv3DefenseRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccwafv3%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudWafv3DefenseRuleBasicDependence11070)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":    "${data.alicloud_wafv3_instances.default.ids.0}",
					"defense_type":   "template",
					"defense_scene":  "ip_blacklist",
					"rule_status":    "1",
					"template_id":    "${alicloud_wafv3_defense_template.defaultZmPPmw-ipblack.defense_template_id}",
					"rule_name":      "tf-test-ip-blacklist",
					"defense_origin": "custom",
					"config": []map[string]interface{}{
						{
							"rule_action": "block",
							"remote_addr": []string{
								"1.1.1.1", "2.2.2.2", "3.3.3.3"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":    CHECKSET,
						"defense_type":   "template",
						"defense_scene":  "ip_blacklist",
						"rule_status":    "1",
						"template_id":    CHECKSET,
						"rule_name":      "tf-test-ip-blacklist",
						"defense_origin": "custom",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rule_status": "0",
					"rule_name":   "tf-test-ip-blacklist-update",
					"config": []map[string]interface{}{
						{
							"rule_action": "monitor",
							"remote_addr": []string{
								"2.2.2.2", "3.3.3.3"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule_status": "0",
						"rule_name":   "tf-test-ip-blacklist-update",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudWafv3DefenseRuleMap11070 = map[string]string{
	"rule_id": CHECKSET,
}

func AlicloudWafv3DefenseRuleBasicDependence11070(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "region_id" {
  default = "cn-hangzhou"
}

data "alicloud_wafv3_instances" "default" {
}

resource "alicloud_wafv3_defense_template" "defaultZmPPmw-ipblack" {
  instance_id           = data.alicloud_wafv3_instances.default.ids.0
  template_origin       = "custom"
  defense_template_name = "1758078494"
  defense_scene         = "ip_blacklist"
  template_type         = "user_custom"
  status                = "1"
  description           = "testCreate"
}


`, name)
}

// Case  DefenseRule__20250723——信息泄露防护规则 11081
func TestAccAliCloudWafv3DefenseRule_basic11081(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_wafv3_defense_rule.default"
	ra := resourceAttrInit(resourceId, AlicloudWafv3DefenseRuleMap11081)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &Wafv3ServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeWafv3DefenseRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccwafv3%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudWafv3DefenseRuleBasicDependence11081)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":    "${data.alicloud_wafv3_instances.default.ids.0}",
					"defense_type":   "template",
					"defense_scene":  "dlp",
					"rule_status":    "1",
					"template_id":    "${alicloud_wafv3_defense_template.defaultZmPPmw-dlp.defense_template_id}",
					"defense_origin": "custom",
					"config": []map[string]interface{}{
						{
							"rule_action": "monitor",
							"conditions": []map[string]interface{}{
								{
									"op_value": "contain",
									"values":   "abc",
									"key":      "URL",
									"sub_key":  "request-uri",
								},
								{
									"op_value": "contain",
									"values":   "phone",
									"key":      "SensitiveInfo",
									"sub_key":  "test",
								},
							},
						},
					},
					"rule_name": "dlp-create-name",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":    CHECKSET,
						"defense_type":   "template",
						"defense_scene":  "dlp",
						"rule_status":    "1",
						"template_id":    CHECKSET,
						"defense_origin": "custom",
						"rule_name":      "dlp-create-name",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rule_status": "0",
					"config": []map[string]interface{}{
						{
							"rule_action": "block",
							"conditions": []map[string]interface{}{
								{
									"op_value": "contain",
									"values":   "ccd",
									"key":      "URL",
									"sub_key":  "rul-aa",
								},
								{
									"op_value": "contain",
									"values":   "401",
									"key":      "HttpCode",
									"sub_key":  "test222",
								},
							},
						},
					},
					"rule_name": "dlp-update-name",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule_status": "0",
						"rule_name":   "dlp-update-name",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudWafv3DefenseRuleMap11081 = map[string]string{
	"rule_id": CHECKSET,
}

func AlicloudWafv3DefenseRuleBasicDependence11081(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "region_id" {
  default = "cn-hangzhou"
}

data "alicloud_wafv3_instances" "default" {
}

resource "alicloud_wafv3_defense_template" "defaultZmPPmw-dlp" {
  instance_id           = data.alicloud_wafv3_instances.default.ids.0
  template_origin       = "custom"
  defense_template_name = "1758078522"
  defense_scene         = "dlp"
  template_type         = "user_custom"
  status                = "1"
  description           = "testCreate"
}


`, name)
}

// Case  DefenseRule__20250723——网页防篡改 11079
func TestAccAliCloudWafv3DefenseRule_basic11079(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_wafv3_defense_rule.default"
	ra := resourceAttrInit(resourceId, AlicloudWafv3DefenseRuleMap11079)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &Wafv3ServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeWafv3DefenseRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccwafv3%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudWafv3DefenseRuleBasicDependence11079)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":    "${data.alicloud_wafv3_instances.default.ids.0}",
					"defense_type":   "template",
					"defense_scene":  "tamperproof",
					"rule_status":    "1",
					"template_id":    "${alicloud_wafv3_defense_template.defaultZmPPmw-fcg.defense_template_id}",
					"defense_origin": "custom",
					"config": []map[string]interface{}{
						{
							"url":      "/abc",
							"ua":       "app",
							"protocol": "https",
						},
					},
					"rule_name": "tamperproof-create-name",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":    CHECKSET,
						"defense_type":   "template",
						"defense_scene":  "tamperproof",
						"rule_status":    "1",
						"template_id":    CHECKSET,
						"defense_origin": "custom",
						"rule_name":      "tamperproof-create-name",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rule_status": "0",
					"config": []map[string]interface{}{
						{
							"url":      "/abcd",
							"ua":       "app2",
							"protocol": "http",
						},
					},
					"rule_name": "tamperproof-update-name",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule_status": "0",
						"rule_name":   "tamperproof-update-name",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudWafv3DefenseRuleMap11079 = map[string]string{
	"rule_id": CHECKSET,
}

func AlicloudWafv3DefenseRuleBasicDependence11079(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "region_id" {
  default = "cn-hangzhou"
}

data "alicloud_wafv3_instances" "default" {
}

resource "alicloud_wafv3_defense_template" "defaultZmPPmw-fcg" {
  instance_id           = data.alicloud_wafv3_instances.default.ids.0
  template_origin       = "custom"
  defense_template_name = "1758078550"
  defense_scene         = "tamperproof"
  template_type         = "user_custom"
  status                = "1"
  description           = "testCreate"
}


`, name)
}

// Case  DefenseRule__20250723—— 白名单 11076
func TestAccAliCloudWafv3DefenseRule_basic11076(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_wafv3_defense_rule.default"
	ra := resourceAttrInit(resourceId, AlicloudWafv3DefenseRuleMap11076)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &Wafv3ServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeWafv3DefenseRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccwafv3%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudWafv3DefenseRuleBasicDependence11076)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":    "${data.alicloud_wafv3_instances.default.ids.0}",
					"defense_type":   "template",
					"defense_scene":  "whitelist",
					"rule_status":    "1",
					"template_id":    "${alicloud_wafv3_defense_template.defaultBQg9ZY.defense_template_id}",
					"rule_name":      "tf-whitelist",
					"defense_origin": "custom",
					"config": []map[string]interface{}{
						{
							"conditions": []map[string]interface{}{
								{
									"op_value": "contain",
									"values":   "abc",
									"key":      "URL",
									"sub_key":  "request_url",
								},
								{
									"op_value": "contain",
									"values":   "test",
									"key":      "URLPath",
									"sub_key":  "reqeust-url",
								},
								{
									"op_value": "eq",
									"values":   "2.2.2.2",
									"key":      "IP",
									"sub_key":  "requset-ip",
								},
							},
							"bypass_tags": []string{
								"customrule", "blacklist", "antiscan", "regular_type"},
							"bypass_regular_types": []string{
								"xss", "sql", "code_exec"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":    CHECKSET,
						"defense_type":   "template",
						"defense_scene":  "whitelist",
						"rule_status":    "1",
						"template_id":    CHECKSET,
						"rule_name":      "tf-whitelist",
						"defense_origin": "custom",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rule_status": "0",
					"rule_name":   "test-whietest-update",
					"config": []map[string]interface{}{
						{
							"conditions": []map[string]interface{}{
								{
									"op_value": "lt",
									"values":   "10",
									"key":      "Content-Length",
									"sub_key":  "requeset-content-length",
								},
							},
							"bypass_tags": []string{
								"cc", "regular_type"},
							"bypass_regular_types": []string{
								"xss"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule_status": "0",
						"rule_name":   "test-whietest-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rule_status": "1",
					"rule_name":   "whitelist-update-2",
					"config": []map[string]interface{}{
						{
							"conditions": []map[string]interface{}{
								{
									"op_value": "contain",
									"values":   "abc",
									"key":      "URL",
									"sub_key":  "testurl",
								},
							},
							"bypass_tags": []string{
								"regular_type"},
							"bypass_regular_types": []string{
								"sqli", "xss", "code_exec"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule_status": "1",
						"rule_name":   "whitelist-update-2",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudWafv3DefenseRuleMap11076 = map[string]string{
	"rule_id": CHECKSET,
}

func AlicloudWafv3DefenseRuleBasicDependence11076(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "region_id" {
  default = "cn-hangzhou"
}

data "alicloud_wafv3_instances" "default" {
}

resource "alicloud_wafv3_defense_template" "defaultBQg9ZY" {
  instance_id           = data.alicloud_wafv3_instances.default.ids.0
  template_origin       = "custom"
  defense_template_name = "1758078580"
  defense_scene         = "whitelist"
  template_type         = "user_custom"
  status                = "1"
  description           = "testCreate"
}


`, name)
}

// Case  DefenseRule__20250723_洪峰限流 11084
func TestAccAliCloudWafv3DefenseRule_basic11084(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_wafv3_defense_rule.default"
	ra := resourceAttrInit(resourceId, AlicloudWafv3DefenseRuleMap11084)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &Wafv3ServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeWafv3DefenseRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccwafv3%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudWafv3DefenseRuleBasicDependence11084)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":    "${data.alicloud_wafv3_instances.default.ids.0}",
					"defense_type":   "template",
					"defense_scene":  "spike_throttle",
					"rule_status":    "1",
					"template_id":    "${alicloud_wafv3_defense_template.defaultfIoHt5-hf.defense_template_id}",
					"rule_name":      "spike_throttle_create",
					"defense_origin": "custom",
					"config": []map[string]interface{}{
						{
							"rule_action": "block",
							"conditions": []map[string]interface{}{
								{
									"op_value": "contain",
									"values":   "abc",
									"key":      "URL",
									"sub_key":  "request_uri",
								},
								{
									"op_value": "contain",
									"values":   "abcde",
									"key":      "URLPath",
									"sub_key":  "request_uri1",
								},
								{
									"op_value": "eq",
									"values":   "1.1.1.1",
									"key":      "IP",
									"sub_key":  "TEST",
								},
							},
							"cn_regions":        "110000,120000,130000",
							"abroad_regions":    "AD,AE,AF",
							"throttle_type":     "qps",
							"throttle_threhold": "500",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":    CHECKSET,
						"defense_type":   "template",
						"defense_scene":  "spike_throttle",
						"rule_status":    "1",
						"template_id":    CHECKSET,
						"rule_name":      "spike_throttle_create",
						"defense_origin": "custom",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rule_status": "0",
					"rule_name":   "spike_throttle-update",
					"config": []map[string]interface{}{
						{
							"rule_action": "monitor",
							"conditions": []map[string]interface{}{
								{
									"op_value": "contain",
									"values":   "abcd",
									"key":      "URLPath",
								},
							},
							"cn_regions":        "110000,120000",
							"abroad_regions":    "AD,AE",
							"throttle_type":     "ratio",
							"throttle_threhold": "41",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule_status": "0",
						"rule_name":   "spike_throttle-update",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudWafv3DefenseRuleMap11084 = map[string]string{
	"rule_id": CHECKSET,
}

func AlicloudWafv3DefenseRuleBasicDependence11084(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "region_id" {
  default = "cn-hangzhou"
}

data "alicloud_wafv3_instances" "default" {
}

resource "alicloud_wafv3_defense_template" "defaultfIoHt5-hf" {
  instance_id           = data.alicloud_wafv3_instances.default.ids.0
  template_origin       = "custom"
  defense_template_name = "1758078609"
  defense_scene         = "spike_throttle"
  template_type         = "user_custom"
  status                = "1"
  description           = "testCreate"
}


`, name)
}

// Test Wafv3 DefenseRule. <<< Resource test cases, automatically generated.
