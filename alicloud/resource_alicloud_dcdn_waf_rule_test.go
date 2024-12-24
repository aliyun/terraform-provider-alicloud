package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Case 1
func TestAccAliCloudDcdnWafRule_basic2264(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_dcdn_waf_rule.default"
	ra := resourceAttrInit(resourceId, AlicloudDcdnWafRuleMap2264)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DcdnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDcdnWafRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testaccDcdnWafRule%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudDcdnWafRuleBasicDependence2264)
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
					"policy_id":     "${alicloud_dcdn_waf_policy.default.id}",
					"rule_name":     "${var.name}",
					"waf_group_ids": "1012",
					"conditions": []map[string]interface{}{
						{
							"key":      "URI",
							"op_value": "ne",
							"values":   "/login.php",
						},
						{
							"key":      "Header",
							"sub_key":  "a",
							"op_value": "eq",
							"values":   "b",
						},
					},
					"status":    "on",
					"cc_status": "on",
					"action":    "monitor",
					"effect":    "rule",
					"rate_limit": []map[string]interface{}{
						{
							"target":    "IP",
							"interval":  "5",
							"threshold": "5",
							"ttl":       "1800",
							"status": []map[string]interface{}{
								{
									"code":  "200",
									"ratio": "60",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy_id":                   CHECKSET,
						"rule_name":                   name,
						"waf_group_ids":               "1012",
						"status":                      "on",
						"cc_status":                   "on",
						"action":                      "monitor",
						"effect":                      "rule",
						"conditions.#":                "2",
						"rate_limit.#":                "1",
						"rate_limit.0.target":         "IP",
						"rate_limit.0.interval":       "5",
						"rate_limit.0.threshold":      "5",
						"rate_limit.0.ttl":            "1800",
						"rate_limit.0.status.#":       "1",
						"rate_limit.0.status.0.code":  "200",
						"rate_limit.0.status.0.ratio": "60",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rule_name": "${var.name}_update",
					"conditions": []map[string]interface{}{
						{
							"key":      "IP",
							"op_value": "ip-contain",
							"values":   "1.1.1.1",
						},
						{
							"key":      "Header",
							"sub_key":  "a",
							"op_value": "eq",
							"values":   "b",
						},
					},
					"status":    "off",
					"cc_status": "on",
					"action":    "deny",
					"effect":    "rule",
					"rate_limit": []map[string]interface{}{
						{
							"target":    "Session",
							"interval":  "10",
							"threshold": "10",
							"ttl":       "1200",
							"status": []map[string]interface{}{
								{
									"code":  "500",
									"ratio": "40",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule_name":                   name + "_update",
						"status":                      "off",
						"cc_status":                   "on",
						"action":                      "deny",
						"effect":                      "rule",
						"conditions.#":                "2",
						"rate_limit.#":                "1",
						"rate_limit.0.target":         "Session",
						"rate_limit.0.interval":       "10",
						"rate_limit.0.threshold":      "10",
						"rate_limit.0.ttl":            "1200",
						"rate_limit.0.status.#":       "1",
						"rate_limit.0.status.0.code":  "500",
						"rate_limit.0.status.0.ratio": "40",
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

var AlicloudDcdnWafRuleMap2264 = map[string]string{}

func AlicloudDcdnWafRuleBasicDependence2264(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_dcdn_waf_policy" "default" {
  defense_scene = "custom_acl"
  policy_name   = var.name
  policy_type   = "custom"
  status        = "on"
}

`, name)
}

// Case 3
func TestAccAliCloudDcdnWafRule_basic2624(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_dcdn_waf_rule.default"
	ra := resourceAttrInit(resourceId, AlicloudDcdnWafRuleMap2624)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DcdnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDcdnWafRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testaccDcdnWafRule%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudDcdnWafRuleBasicDependence2624)
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
					"action":      "deny",
					"policy_id":   "${alicloud_dcdn_waf_policy.default.id}",
					"remote_addr": []string{"1.1.1.1", "2.2.2.0/24", "::1", "abcd::abcd"},
					"rule_name":   "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"action":        "deny",
						"policy_id":     CHECKSET,
						"remote_addr.#": "4",
						"rule_name":     name,
					}),
				),
			}, {
				Config: testAccConfig(map[string]interface{}{
					"action":      "monitor",
					"remote_addr": []string{"1.1.1.1"},
					"rule_name":   "${var.name}_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"action":        "monitor",
						"remote_addr.#": "1",
						"rule_name":     name + "_update",
					}),
				),
			}, {
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudDcdnWafRuleMap2624 = map[string]string{}

func AlicloudDcdnWafRuleBasicDependence2624(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_dcdn_waf_policy" "default" {
  defense_scene = "ip_blacklist"
  policy_name   = var.name
  policy_type   = "custom"
  status        = "on"
}

`, name)
}

// Case 5
func TestAccAliCloudDcdnWafRule_basic2749(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_dcdn_waf_rule.default"
	ra := resourceAttrInit(resourceId, AlicloudDcdnWafRuleMap2749)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DcdnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDcdnWafRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testaccDcdnWafRule%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudDcdnWafRuleBasicDependence2749)
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
					"status": "on",
					"scenes": []string{"waf_group", "custom_acl", "ip_blacklist", "region_block", "bot", "anti_scan"},
					"conditions": []map[string]interface{}{
						{
							"key":      "Http-Method",
							"op_value": "match-one",
							"values":   "GET,POST,DELETE",
						},
					},
					"rule_name":     "${var.name}",
					"policy_id":     "${alicloud_dcdn_waf_policy.default.id}",
					"regular_types": []string{"sqli", "xss", "code_exec"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status":          "on",
						"scenes.#":        "6",
						"conditions.#":    "1",
						"rule_name":       name,
						"policy_id":       CHECKSET,
						"regular_types.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status":        "off",
					"scenes":        []string{"waf_group"},
					"regular_rules": []string{"100003", "100004"},
					"regular_types": REMOVEKEY,
					"conditions": []map[string]interface{}{
						{
							"key":      "Referer",
							"op_value": "none",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status":          "off",
						"scenes.#":        "1",
						"regular_rules.#": "2",
						"conditions.#":    "1",
						"regular_types.#": "0",
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

var AlicloudDcdnWafRuleMap2749 = map[string]string{}

func AlicloudDcdnWafRuleBasicDependence2749(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_dcdn_waf_policy" "default" {
  defense_scene = "whitelist"
  policy_name   = var.name
  policy_type   = "custom"
  status        = "on"
}

`, name)
}

// Case 6
func TestAccAliCloudDcdnWafRule_basic2753(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_dcdn_waf_rule.default"
	ra := resourceAttrInit(resourceId, AlicloudDcdnWafRuleMap2753)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DcdnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDcdnWafRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testaccDcdnWafRule%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudDcdnWafRuleBasicDependence2753)
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
					"status":            "on",
					"action":            "deny",
					"rule_name":         "${var.name}",
					"other_region_list": "JP,GB",
					"cn_region_list":    "110000,TW,MO",
					"policy_id":         "${alicloud_dcdn_waf_policy.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status":            "on",
						"rule_name":         name,
						"action":            "deny",
						"other_region_list": "JP,GB",
						"cn_region_list":    "110000,TW,MO",
						"policy_id":         CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"action":            "monitor",
					"other_region_list": "DE,GB",
					"cn_region_list":    "510000,430000",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"action":            "monitor",
						"other_region_list": "DE,GB",
						"cn_region_list":    "510000,430000",
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

var AlicloudDcdnWafRuleMap2753 = map[string]string{}

func AlicloudDcdnWafRuleBasicDependence2753(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_dcdn_waf_policy" "default" {
  defense_scene = "region_block"
  policy_name   = var.name
  policy_type   = "custom"
  status        = "on"
}


`, name)
}
