package alicloud

import (
	"fmt"
	"testing"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test DdosBgp Policy. >>> Resource test cases, automatically generated.
// Case l3策略类型测试_2 7021
func TestAccAliCloudDdosBgpPolicy_basic7021(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ddos_bgp_policy.default"
	ra := resourceAttrInit(resourceId, AlicloudDdosBgpPolicyMap7021)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DdosBgpServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDdosBgpPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testacc_bgp%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudDdosBgpPolicyBasicDependence7021)
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
					"type":        "l3",
					"policy_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"type":        "l3",
						"policy_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"content": []map[string]interface{}{
						{
							"black_ip_list_expire_at": fmt.Sprint(time.Now().Add(10 * time.Minute).Unix()),
							"intelligence_level":      "weak",
							"whiten_gfbr_nets":        "false",
							"region_block_country_list": []string{
								"${var.region_block_country_list}"},
							"region_block_province_list": []string{
								"${var.region_block_province_list}"},
							"source_limit": []map[string]interface{}{
								{
									"pps":     "32",
									"bps":     "1024",
									"syn_pps": "1",
									"syn_bps": "1024",
								},
							},
							"reflect_block_udp_port_list": []string{
								"${var.reflect_block_udp_port_list}", "${var.reflect_block_udp_port_list-update}", "${var.reflect_block_udp_port_list2}"},
							"enable_intelligence": "false",
							"enable_drop_icmp":    "false",
							"port_rule_list": []map[string]interface{}{
								{
									"port_rule_id":   "1111",
									"protocol":       "tcp",
									"src_port_start": "0",
									"src_port_end":   "65535",
									"dst_port_start": "0",
									"dst_port_end":   "65531",
									"seq_no":         "2",
									"match_action":   "drop",
								},
								{
									"port_rule_id":   "2222",
									"protocol":       "tcp",
									"src_port_start": "2",
									"src_port_end":   "3",
									"dst_port_start": "4",
									"dst_port_end":   "5",
									"match_action":   "drop",
									"seq_no":         "3",
								},
								{
									"port_rule_id":   "3333",
									"protocol":       "tcp",
									"src_port_start": "4",
									"src_port_end":   "5",
									"dst_port_start": "5",
									"dst_port_end":   "6",
									"match_action":   "drop",
									"seq_no":         "3",
								},
							},
							"finger_print_rule_list": []map[string]interface{}{
								{
									"protocol":             "tcp",
									"src_port_start":       "0",
									"src_port_end":         "65535",
									"dst_port_start":       "0",
									"dst_port_end":         "65535",
									"min_pkt_len":          "11",
									"max_pkt_len":          "22",
									"offset":               "3",
									"payload_bytes":        "11111111",
									"match_action":         "accept",
									"rate_value":           "1",
									"seq_no":               "1",
									"finger_print_rule_id": "1",
								},
								{
									"protocol":             "tcp",
									"src_port_start":       "0",
									"src_port_end":         "65535",
									"dst_port_start":       "0",
									"dst_port_end":         "65535",
									"min_pkt_len":          "11",
									"max_pkt_len":          "22",
									"offset":               "3",
									"payload_bytes":        "11111111",
									"match_action":         "accept",
									"rate_value":           "1",
									"seq_no":               "1",
									"finger_print_rule_id": "2",
								},
								{
									"protocol":             "tcp",
									"src_port_start":       "0",
									"src_port_end":         "65535",
									"dst_port_start":       "0",
									"dst_port_end":         "65535",
									"min_pkt_len":          "11",
									"max_pkt_len":          "22",
									"offset":               "3",
									"payload_bytes":        "11111111",
									"match_action":         "accept",
									"rate_value":           "1",
									"seq_no":               "1",
									"finger_print_rule_id": "3",
								},
							},
							"source_block_list": []map[string]interface{}{
								{
									"type":                 "4",
									"block_expire_seconds": "60",
									"every_seconds":        "720",
									"exceed_limit_times":   "5",
								},
								{
									"type":                 "5",
									"block_expire_seconds": "80",
									"every_seconds":        "80",
									"exceed_limit_times":   "22",
								},
								{
									"type":                 "6",
									"block_expire_seconds": "80",
									"every_seconds":        "66",
									"exceed_limit_times":   "77",
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
				Config: testAccConfig(map[string]interface{}{
					"content": []map[string]interface{}{
						{
							"black_ip_list_expire_at": fmt.Sprint(time.Now().Add(10 * time.Minute).Unix()),
							"whiten_gfbr_nets":        "true",
							"enable_drop_icmp":        "true",
							"region_block_country_list": []string{
								"${var.region_block_country_list-update}", "${var.region_block_country_list}", "${var.region_block_country_list2}"},
							"region_block_province_list": []string{
								"${var.region_block_province_list-_update}", "${var.region_block_province_list}", "${var.region_block_province_list-_update1}"},
							"intelligence_level": "default",
							"port_rule_list": []map[string]interface{}{
								{
									"protocol":       "udp",
									"src_port_start": "3",
									"src_port_end":   "65533",
									"dst_port_start": "1",
									"dst_port_end":   "65535",
									"match_action":   "drop",
									"seq_no":         "4",
									"port_rule_id":   "2222",
								},
							},
							"finger_print_rule_list": []map[string]interface{}{
								{
									"finger_print_rule_id": "111",
									"protocol":             "udp",
									"src_port_start":       "2",
									"src_port_end":         "65534",
									"dst_port_start":       "1",
									"dst_port_end":         "65534",
									"min_pkt_len":          "12",
									"max_pkt_len":          "21",
									"offset":               "4",
									"payload_bytes":        "112221",
									"match_action":         "ip_rate",
									"rate_value":           "12",
									"seq_no":               "2",
								},
							},
							"source_limit": []map[string]interface{}{
								{
									"pps":     "33",
									"bps":     "1066",
									"syn_pps": "22",
									"syn_bps": "1056",
								},
							},
							"source_block_list": []map[string]interface{}{
								{
									"type":                 "5",
									"block_expire_seconds": "70",
									"every_seconds":        "80",
									"exceed_limit_times":   "6",
								},
							},
							"enable_intelligence": "true",
							"reflect_block_udp_port_list": []string{
								"${var.reflect_block_udp_port_list-update}", "${var.reflect_block_udp_port_list}"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"content": []map[string]interface{}{
						{
							"intelligence_level": "default",
							"whiten_gfbr_nets":   "true",
							"enable_drop_icmp":   "true",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"type": "l3",
					"content": []map[string]interface{}{
						{
							"black_ip_list_expire_at": fmt.Sprint(time.Now().Add(10 * time.Minute).Unix()),
							"intelligence_level":      "weak",
							"whiten_gfbr_nets":        "false",
							"region_block_country_list": []string{
								"${var.region_block_country_list}"},
							"region_block_province_list": []string{
								"${var.region_block_province_list}"},
							"source_limit": []map[string]interface{}{
								{
									"pps":     "32",
									"bps":     "1024",
									"syn_pps": "1",
									"syn_bps": "1024",
								},
							},
							"reflect_block_udp_port_list": []string{
								"${var.reflect_block_udp_port_list}", "${var.reflect_block_udp_port_list-update}", "${var.reflect_block_udp_port_list2}"},
							"enable_intelligence": "false",
							"enable_drop_icmp":    "false",
							"port_rule_list": []map[string]interface{}{
								{
									"port_rule_id":   "1111",
									"protocol":       "tcp",
									"src_port_start": "0",
									"src_port_end":   "65535",
									"dst_port_start": "0",
									"dst_port_end":   "65531",
									"seq_no":         "2",
									"match_action":   "drop",
								},
								{
									"port_rule_id":   "2222",
									"protocol":       "tcp",
									"src_port_start": "2",
									"src_port_end":   "3",
									"dst_port_start": "4",
									"dst_port_end":   "5",
									"match_action":   "drop",
									"seq_no":         "3",
								},
								{
									"port_rule_id":   "3333",
									"protocol":       "tcp",
									"src_port_start": "4",
									"src_port_end":   "5",
									"dst_port_start": "5",
									"dst_port_end":   "6",
									"match_action":   "drop",
									"seq_no":         "3",
								},
							},
							"finger_print_rule_list": []map[string]interface{}{
								{
									"protocol":             "tcp",
									"src_port_start":       "0",
									"src_port_end":         "65535",
									"dst_port_start":       "0",
									"dst_port_end":         "65535",
									"min_pkt_len":          "11",
									"max_pkt_len":          "22",
									"offset":               "3",
									"payload_bytes":        "11111111",
									"match_action":         "accept",
									"rate_value":           "1",
									"seq_no":               "1",
									"finger_print_rule_id": "1",
								},
								{
									"protocol":             "tcp",
									"src_port_start":       "0",
									"src_port_end":         "65535",
									"dst_port_start":       "0",
									"dst_port_end":         "65535",
									"min_pkt_len":          "11",
									"max_pkt_len":          "22",
									"offset":               "3",
									"payload_bytes":        "11111111",
									"match_action":         "accept",
									"rate_value":           "1",
									"seq_no":               "1",
									"finger_print_rule_id": "2",
								},
								{
									"protocol":             "tcp",
									"src_port_start":       "0",
									"src_port_end":         "65535",
									"dst_port_start":       "0",
									"dst_port_end":         "65535",
									"min_pkt_len":          "11",
									"max_pkt_len":          "22",
									"offset":               "3",
									"payload_bytes":        "11111111",
									"match_action":         "accept",
									"rate_value":           "1",
									"seq_no":               "1",
									"finger_print_rule_id": "3",
								},
							},
							"source_block_list": []map[string]interface{}{
								{
									"type":                 "4",
									"block_expire_seconds": "60",
									"every_seconds":        "720",
									"exceed_limit_times":   "5",
								},
								{
									"type":                 "5",
									"block_expire_seconds": "80",
									"every_seconds":        "80",
									"exceed_limit_times":   "22",
								},
								{
									"type":                 "6",
									"block_expire_seconds": "80",
									"every_seconds":        "66",
									"exceed_limit_times":   "77",
								},
							},
						},
					},
					"policy_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"type":        "l3",
						"policy_name": name + "_update",
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

var AlicloudDdosBgpPolicyMap7021 = map[string]string{}

func AlicloudDdosBgpPolicyBasicDependence7021(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "reflect_block_udp_port_list-update" {
  default = "9090"
}

variable "region_block_country_list-update" {
  default = "3"
}

variable "ip_list" {
  default = "1.1.1.1"
}

variable "policy_name" {
  default = "policy_test_l3"
}

variable "reflect_block_udp_port_list2" {
  default = "7070"
}

variable "region_block_country_list" {
  default = "2"
}

variable "region_block_province_list" {
  default = "11"
}

variable "reflect_block_udp_port_list" {
  default = "8888"
}

variable "region_block_province_list-_update" {
  default = "65"
}

variable "region_block_country_list2" {
  default = "4"
}

variable "region_block_province_list-_update1" {
  default = "61"
}


`, name)
}

// Case 更新l4类型智能开关_2 7022
func TestAccAliCloudDdosBgpPolicy_basic7022(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ddos_bgp_policy.default"
	ra := resourceAttrInit(resourceId, AlicloudDdosBgpPolicyMap7022)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DdosBgpServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDdosBgpPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testacc_bgp%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudDdosBgpPolicyBasicDependence7022)
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
					"type":        "l4",
					"policy_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"type":        "l4",
						"policy_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"content": []map[string]interface{}{
						{
							"enable_defense": "false",
							"layer4_rule_list": []map[string]interface{}{
								{
									"name":     "11",
									"priority": "10",
									"method":   "hex",
									"match":    "1",
									"action":   "1",
									"limited":  "0",
									"condition_list": []map[string]interface{}{
										{
											"arg":      "3C",
											"position": "1",
											"depth":    "2",
										},
									},
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
				Config: testAccConfig(map[string]interface{}{
					"content": []map[string]interface{}{
						{
							"enable_defense": "true",
							"layer4_rule_list": []map[string]interface{}{
								{
									"name":     "11",
									"priority": "34",
									"method":   "hex",
									"match":    "0",
									"action":   "2",
									"limited":  "4",
									"condition_list": []map[string]interface{}{
										{
											"arg":      "9C",
											"position": "2",
											"depth":    "3",
										},
									},
								},
								{
									"name":     "333",
									"priority": "34",
									"method":   "hex",
									"match":    "0",
									"action":   "1",
									"limited":  "4",
									"condition_list": []map[string]interface{}{
										{
											"arg":      "9C",
											"position": "2",
											"depth":    "3",
										},
									},
								},
								{
									"name":     "55",
									"priority": "34",
									"method":   "char",
									"match":    "0",
									"action":   "2",
									"limited":  "4",
									"condition_list": []map[string]interface{}{
										{
											"arg":      "3333",
											"position": "2",
											"depth":    "3",
										},
									},
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
				Config: testAccConfig(map[string]interface{}{
					"content": []map[string]interface{}{
						{
							"enable_defense": "true",
							"layer4_rule_list": []map[string]interface{}{
								{
									"name":     "11",
									"priority": "34",
									"method":   "hex",
									"match":    "0",
									"action":   "1",
									"limited":  "4",
									"condition_list": []map[string]interface{}{
										{
											"arg":      "9C",
											"position": "2",
											"depth":    "3",
										},
									},
								},
								{
									"name":     "333",
									"priority": "34",
									"method":   "hex",
									"match":    "0",
									"action":   "2",
									"limited":  "4",
									"condition_list": []map[string]interface{}{
										{
											"arg":      "9C",
											"position": "2",
											"depth":    "3",
										},
									},
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
				Config: testAccConfig(map[string]interface{}{
					"type":        "l4",
					"policy_name": name + "_update",
					"content": []map[string]interface{}{
						{
							"enable_defense": "false",
							"layer4_rule_list": []map[string]interface{}{
								{
									"name":     "11",
									"priority": "10",
									"method":   "hex",
									"match":    "1",
									"action":   "1",
									"limited":  "0",
									"condition_list": []map[string]interface{}{
										{
											"arg":      "3C",
											"position": "1",
											"depth":    "2",
										},
									},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"type":        "l4",
						"policy_name": name + "_update",
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

var AlicloudDdosBgpPolicyMap7022 = map[string]string{}

func AlicloudDdosBgpPolicyBasicDependence7022(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "policy_name" {
  default = "test_l4_policy"
}


`, name)
}

// Case l3策略类型测试 6912
func TestAccAliCloudDdosBgpPolicy_basic6912(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ddos_bgp_policy.default"
	ra := resourceAttrInit(resourceId, AlicloudDdosBgpPolicyMap6912)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DdosBgpServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDdosBgpPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testacc_bgp%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudDdosBgpPolicyBasicDependence6912)
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
					"type":        "l3",
					"policy_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"type":        "l3",
						"policy_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"content": []map[string]interface{}{
						{
							"black_ip_list_expire_at": fmt.Sprint(time.Now().Add(10 * time.Minute).Unix()),
							"intelligence_level":      "weak",
							"whiten_gfbr_nets":        "false",
							"region_block_country_list": []string{
								"${var.region_block_country_list}"},
							"region_block_province_list": []string{
								"${var.region_block_province_list}"},
							"source_limit": []map[string]interface{}{
								{
									"pps":     "32",
									"bps":     "1024",
									"syn_pps": "1",
									"syn_bps": "1024",
								},
							},
							"reflect_block_udp_port_list": []string{
								"${var.reflect_block_udp_port_list}", "${var.reflect_block_udp_port_list-update}", "${var.reflect_block_udp_port_list2}"},
							"enable_intelligence": "false",
							"enable_drop_icmp":    "false",
							"port_rule_list": []map[string]interface{}{
								{
									"port_rule_id":   "1111",
									"protocol":       "tcp",
									"src_port_start": "0",
									"src_port_end":   "65535",
									"dst_port_start": "0",
									"dst_port_end":   "65531",
									"seq_no":         "2",
									"match_action":   "drop",
								},
								{
									"port_rule_id":   "2222",
									"protocol":       "tcp",
									"src_port_start": "2",
									"src_port_end":   "3",
									"dst_port_start": "4",
									"dst_port_end":   "5",
									"match_action":   "drop",
									"seq_no":         "3",
								},
								{
									"port_rule_id":   "3333",
									"protocol":       "tcp",
									"src_port_start": "4",
									"src_port_end":   "5",
									"dst_port_start": "5",
									"dst_port_end":   "6",
									"match_action":   "drop",
									"seq_no":         "3",
								},
							},
							"finger_print_rule_list": []map[string]interface{}{
								{
									"protocol":             "tcp",
									"src_port_start":       "0",
									"src_port_end":         "65535",
									"dst_port_start":       "0",
									"dst_port_end":         "65535",
									"min_pkt_len":          "11",
									"max_pkt_len":          "22",
									"offset":               "3",
									"payload_bytes":        "11111111",
									"match_action":         "accept",
									"rate_value":           "1",
									"seq_no":               "1",
									"finger_print_rule_id": "1",
								},
								{
									"protocol":             "tcp",
									"src_port_start":       "0",
									"src_port_end":         "65535",
									"dst_port_start":       "0",
									"dst_port_end":         "65535",
									"min_pkt_len":          "11",
									"max_pkt_len":          "22",
									"offset":               "3",
									"payload_bytes":        "11111111",
									"match_action":         "accept",
									"rate_value":           "1",
									"seq_no":               "1",
									"finger_print_rule_id": "2",
								},
								{
									"protocol":             "tcp",
									"src_port_start":       "0",
									"src_port_end":         "65535",
									"dst_port_start":       "0",
									"dst_port_end":         "65535",
									"min_pkt_len":          "11",
									"max_pkt_len":          "22",
									"offset":               "3",
									"payload_bytes":        "11111111",
									"match_action":         "accept",
									"rate_value":           "1",
									"seq_no":               "1",
									"finger_print_rule_id": "3",
								},
							},
							"source_block_list": []map[string]interface{}{
								{
									"type":                 "4",
									"block_expire_seconds": "60",
									"every_seconds":        "720",
									"exceed_limit_times":   "5",
								},
								{
									"type":                 "5",
									"block_expire_seconds": "80",
									"every_seconds":        "80",
									"exceed_limit_times":   "22",
								},
								{
									"type":                 "6",
									"block_expire_seconds": "80",
									"every_seconds":        "66",
									"exceed_limit_times":   "77",
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
				Config: testAccConfig(map[string]interface{}{
					"content": []map[string]interface{}{
						{
							"black_ip_list_expire_at": fmt.Sprint(time.Now().Add(10 * time.Minute).Unix()),
							"whiten_gfbr_nets":        "true",
							"enable_drop_icmp":        "true",
							"region_block_country_list": []string{
								"${var.region_block_country_list-update}", "${var.region_block_country_list}", "${var.region_block_country_list2}"},
							"region_block_province_list": []string{
								"${var.region_block_province_list-_update}", "${var.region_block_province_list}", "${var.region_block_province_list-_update1}"},
							"intelligence_level": "default",
							"port_rule_list": []map[string]interface{}{
								{
									"protocol":       "udp",
									"src_port_start": "3",
									"src_port_end":   "65533",
									"dst_port_start": "1",
									"dst_port_end":   "65535",
									"match_action":   "drop",
									"seq_no":         "4",
									"port_rule_id":   "2222",
								},
							},
							"finger_print_rule_list": []map[string]interface{}{
								{
									"finger_print_rule_id": "111",
									"protocol":             "udp",
									"src_port_start":       "2",
									"src_port_end":         "65534",
									"dst_port_start":       "1",
									"dst_port_end":         "65534",
									"min_pkt_len":          "12",
									"max_pkt_len":          "21",
									"offset":               "4",
									"payload_bytes":        "112221",
									"match_action":         "ip_rate",
									"rate_value":           "12",
									"seq_no":               "2",
								},
							},
							"source_limit": []map[string]interface{}{
								{
									"pps":     "33",
									"bps":     "1066",
									"syn_pps": "22",
									"syn_bps": "1056",
								},
							},
							"source_block_list": []map[string]interface{}{
								{
									"type":                 "5",
									"block_expire_seconds": "70",
									"every_seconds":        "80",
									"exceed_limit_times":   "6",
								},
							},
							"enable_intelligence": "true",
							"reflect_block_udp_port_list": []string{
								"${var.reflect_block_udp_port_list-update}", "${var.reflect_block_udp_port_list}"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"content": []map[string]interface{}{
						{
							"intelligence_level": "default",
							"whiten_gfbr_nets":   "true",
							"enable_drop_icmp":   "true",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"type": "l3",
					"content": []map[string]interface{}{
						{
							"black_ip_list_expire_at": fmt.Sprint(time.Now().Add(10 * time.Minute).Unix()),
							"intelligence_level":      "weak",
							"whiten_gfbr_nets":        "false",
							"region_block_country_list": []string{
								"${var.region_block_country_list}"},
							"region_block_province_list": []string{
								"${var.region_block_province_list}"},
							"source_limit": []map[string]interface{}{
								{
									"pps":     "32",
									"bps":     "1024",
									"syn_pps": "1",
									"syn_bps": "1024",
								},
							},
							"reflect_block_udp_port_list": []string{
								"${var.reflect_block_udp_port_list}", "${var.reflect_block_udp_port_list-update}", "${var.reflect_block_udp_port_list2}"},
							"enable_intelligence": "false",
							"enable_drop_icmp":    "false",
							"port_rule_list": []map[string]interface{}{
								{
									"port_rule_id":   "1111",
									"protocol":       "tcp",
									"src_port_start": "0",
									"src_port_end":   "65535",
									"dst_port_start": "0",
									"dst_port_end":   "65531",
									"seq_no":         "2",
									"match_action":   "drop",
								},
								{
									"port_rule_id":   "2222",
									"protocol":       "tcp",
									"src_port_start": "2",
									"src_port_end":   "3",
									"dst_port_start": "4",
									"dst_port_end":   "5",
									"match_action":   "drop",
									"seq_no":         "3",
								},
								{
									"port_rule_id":   "3333",
									"protocol":       "tcp",
									"src_port_start": "4",
									"src_port_end":   "5",
									"dst_port_start": "5",
									"dst_port_end":   "6",
									"match_action":   "drop",
									"seq_no":         "3",
								},
							},
							"finger_print_rule_list": []map[string]interface{}{
								{
									"protocol":             "tcp",
									"src_port_start":       "0",
									"src_port_end":         "65535",
									"dst_port_start":       "0",
									"dst_port_end":         "65535",
									"min_pkt_len":          "11",
									"max_pkt_len":          "22",
									"offset":               "3",
									"payload_bytes":        "11111111",
									"match_action":         "accept",
									"rate_value":           "1",
									"seq_no":               "1",
									"finger_print_rule_id": "1",
								},
								{
									"protocol":             "tcp",
									"src_port_start":       "0",
									"src_port_end":         "65535",
									"dst_port_start":       "0",
									"dst_port_end":         "65535",
									"min_pkt_len":          "11",
									"max_pkt_len":          "22",
									"offset":               "3",
									"payload_bytes":        "11111111",
									"match_action":         "accept",
									"rate_value":           "1",
									"seq_no":               "1",
									"finger_print_rule_id": "2",
								},
								{
									"protocol":             "tcp",
									"src_port_start":       "0",
									"src_port_end":         "65535",
									"dst_port_start":       "0",
									"dst_port_end":         "65535",
									"min_pkt_len":          "11",
									"max_pkt_len":          "22",
									"offset":               "3",
									"payload_bytes":        "11111111",
									"match_action":         "accept",
									"rate_value":           "1",
									"seq_no":               "1",
									"finger_print_rule_id": "3",
								},
							},
							"source_block_list": []map[string]interface{}{
								{
									"type":                 "4",
									"block_expire_seconds": "60",
									"every_seconds":        "720",
									"exceed_limit_times":   "5",
								},
								{
									"type":                 "5",
									"block_expire_seconds": "80",
									"every_seconds":        "80",
									"exceed_limit_times":   "22",
								},
								{
									"type":                 "6",
									"block_expire_seconds": "80",
									"every_seconds":        "66",
									"exceed_limit_times":   "77",
								},
							},
						},
					},
					"policy_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"type":        "l3",
						"policy_name": name + "_update",
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

var AlicloudDdosBgpPolicyMap6912 = map[string]string{}

func AlicloudDdosBgpPolicyBasicDependence6912(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "reflect_block_udp_port_list-update" {
  default = "9090"
}

variable "region_block_country_list-update" {
  default = "3"
}

variable "ip_list" {
  default = "1.1.1.1"
}

variable "policy_name" {
  default = "policy_test_l3"
}

variable "reflect_block_udp_port_list2" {
  default = "7070"
}

variable "region_block_country_list" {
  default = "2"
}

variable "region_block_province_list" {
  default = "11"
}

variable "reflect_block_udp_port_list" {
  default = "8888"
}

variable "region_block_province_list-_update" {
  default = "65"
}

variable "region_block_country_list2" {
  default = "4"
}

variable "region_block_province_list-_update1" {
  default = "61"
}


`, name)
}

// Case 更新l4类型智能开关 6856
func TestAccAliCloudDdosBgpPolicy_basic6856(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ddos_bgp_policy.default"
	ra := resourceAttrInit(resourceId, AlicloudDdosBgpPolicyMap6856)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DdosBgpServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDdosBgpPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testacc_bgp%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudDdosBgpPolicyBasicDependence6856)
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
					"type":        "l4",
					"policy_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"type":        "l4",
						"policy_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"content": []map[string]interface{}{
						{
							"enable_defense": "false",
							"layer4_rule_list": []map[string]interface{}{
								{
									"name":     "11",
									"priority": "10",
									"method":   "hex",
									"match":    "1",
									"action":   "1",
									"limited":  "0",
									"condition_list": []map[string]interface{}{
										{
											"arg":      "3C",
											"position": "1",
											"depth":    "2",
										},
									},
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
				Config: testAccConfig(map[string]interface{}{
					"content": []map[string]interface{}{
						{
							"enable_defense": "true",
							"layer4_rule_list": []map[string]interface{}{
								{
									"name":     "11",
									"priority": "34",
									"method":   "hex",
									"match":    "0",
									"action":   "2",
									"limited":  "4",
									"condition_list": []map[string]interface{}{
										{
											"arg":      "9C",
											"position": "2",
											"depth":    "3",
										},
									},
								},
								{
									"name":     "333",
									"priority": "34",
									"method":   "hex",
									"match":    "0",
									"action":   "1",
									"limited":  "4",
									"condition_list": []map[string]interface{}{
										{
											"arg":      "9C",
											"position": "2",
											"depth":    "3",
										},
									},
								},
								{
									"name":     "55",
									"priority": "34",
									"method":   "char",
									"match":    "0",
									"action":   "2",
									"limited":  "4",
									"condition_list": []map[string]interface{}{
										{
											"arg":      "3333",
											"position": "2",
											"depth":    "3",
										},
									},
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
				Config: testAccConfig(map[string]interface{}{
					"content": []map[string]interface{}{
						{
							"enable_defense": "true",
							"layer4_rule_list": []map[string]interface{}{
								{
									"name":     "11",
									"priority": "34",
									"method":   "hex",
									"match":    "0",
									"action":   "1",
									"limited":  "4",
									"condition_list": []map[string]interface{}{
										{
											"arg":      "9C",
											"position": "2",
											"depth":    "3",
										},
									},
								},
								{
									"name":     "333",
									"priority": "34",
									"method":   "hex",
									"match":    "0",
									"action":   "2",
									"limited":  "4",
									"condition_list": []map[string]interface{}{
										{
											"arg":      "9C",
											"position": "2",
											"depth":    "3",
										},
									},
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
				Config: testAccConfig(map[string]interface{}{
					"type":        "l4",
					"policy_name": name + "_update",
					"content": []map[string]interface{}{
						{
							"enable_defense": "false",
							"layer4_rule_list": []map[string]interface{}{
								{
									"name":     "11",
									"priority": "10",
									"method":   "hex",
									"match":    "1",
									"action":   "1",
									"limited":  "0",
									"condition_list": []map[string]interface{}{
										{
											"arg":      "3C",
											"position": "1",
											"depth":    "2",
										},
									},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"type":        "l4",
						"policy_name": name + "_update",
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

var AlicloudDdosBgpPolicyMap6856 = map[string]string{}

func AlicloudDdosBgpPolicyBasicDependence6856(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "policy_name" {
  default = "test_l4_policy"
}


`, name)
}

// Case l3策略类型测试_2 7021  twin
func TestAccAliCloudDdosBgpPolicy_basic7021_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ddos_bgp_policy.default"
	ra := resourceAttrInit(resourceId, AlicloudDdosBgpPolicyMap7021)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DdosBgpServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDdosBgpPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testacc_bgp%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudDdosBgpPolicyBasicDependence7021)
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
					"type": "l3",
					"content": []map[string]interface{}{
						{
							"black_ip_list_expire_at": fmt.Sprint(time.Now().Add(10 * time.Minute).Unix()),
							"intelligence_level":      "weak",
							"whiten_gfbr_nets":        "false",
							"region_block_country_list": []string{
								"${var.region_block_country_list}"},
							"region_block_province_list": []string{
								"${var.region_block_province_list}"},
							"source_limit": []map[string]interface{}{
								{
									"pps":     "32",
									"bps":     "1024",
									"syn_pps": "1",
									"syn_bps": "1024",
								},
							},
							"reflect_block_udp_port_list": []string{
								"${var.reflect_block_udp_port_list}", "${var.reflect_block_udp_port_list-update}", "${var.reflect_block_udp_port_list2}"},
							"enable_intelligence": "false",
							"enable_drop_icmp":    "false",
							"port_rule_list": []map[string]interface{}{
								{
									"port_rule_id":   "1111",
									"protocol":       "tcp",
									"src_port_start": "0",
									"src_port_end":   "65535",
									"dst_port_start": "0",
									"dst_port_end":   "65531",
									"seq_no":         "2",
									"match_action":   "drop",
								},
								{
									"port_rule_id":   "2222",
									"protocol":       "tcp",
									"src_port_start": "2",
									"src_port_end":   "3",
									"dst_port_start": "4",
									"dst_port_end":   "5",
									"match_action":   "drop",
									"seq_no":         "3",
								},
								{
									"port_rule_id":   "3333",
									"protocol":       "tcp",
									"src_port_start": "4",
									"src_port_end":   "5",
									"dst_port_start": "5",
									"dst_port_end":   "6",
									"match_action":   "drop",
									"seq_no":         "3",
								},
							},
							"finger_print_rule_list": []map[string]interface{}{
								{
									"protocol":             "tcp",
									"src_port_start":       "0",
									"src_port_end":         "65535",
									"dst_port_start":       "0",
									"dst_port_end":         "65535",
									"min_pkt_len":          "11",
									"max_pkt_len":          "22",
									"offset":               "3",
									"payload_bytes":        "11111111",
									"match_action":         "accept",
									"rate_value":           "1",
									"seq_no":               "1",
									"finger_print_rule_id": "1",
								},
								{
									"protocol":             "tcp",
									"src_port_start":       "0",
									"src_port_end":         "65535",
									"dst_port_start":       "0",
									"dst_port_end":         "65535",
									"min_pkt_len":          "11",
									"max_pkt_len":          "22",
									"offset":               "3",
									"payload_bytes":        "11111111",
									"match_action":         "accept",
									"rate_value":           "1",
									"seq_no":               "1",
									"finger_print_rule_id": "2",
								},
								{
									"protocol":             "tcp",
									"src_port_start":       "0",
									"src_port_end":         "65535",
									"dst_port_start":       "0",
									"dst_port_end":         "65535",
									"min_pkt_len":          "11",
									"max_pkt_len":          "22",
									"offset":               "3",
									"payload_bytes":        "11111111",
									"match_action":         "accept",
									"rate_value":           "1",
									"seq_no":               "1",
									"finger_print_rule_id": "3",
								},
							},
							"source_block_list": []map[string]interface{}{
								{
									"type":                 "4",
									"block_expire_seconds": "60",
									"every_seconds":        "720",
									"exceed_limit_times":   "5",
								},
								{
									"type":                 "5",
									"block_expire_seconds": "80",
									"every_seconds":        "80",
									"exceed_limit_times":   "22",
								},
								{
									"type":                 "6",
									"block_expire_seconds": "80",
									"every_seconds":        "66",
									"exceed_limit_times":   "77",
								},
							},
						},
					},
					"policy_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"type":        "l3",
						"policy_name": name,
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

// Case 更新l4类型智能开关_2 7022  twin
func TestAccAliCloudDdosBgpPolicy_basic7022_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ddos_bgp_policy.default"
	ra := resourceAttrInit(resourceId, AlicloudDdosBgpPolicyMap7022)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DdosBgpServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDdosBgpPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testacc_bgp%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudDdosBgpPolicyBasicDependence7022)
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
					"type":        "l4",
					"policy_name": name,
					"content": []map[string]interface{}{
						{
							"enable_defense": "false",
							"layer4_rule_list": []map[string]interface{}{
								{
									"name":     "11",
									"priority": "10",
									"method":   "hex",
									"match":    "1",
									"action":   "1",
									"limited":  "0",
									"condition_list": []map[string]interface{}{
										{
											"arg":      "3C",
											"position": "1",
											"depth":    "2",
										},
									},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"type":        "l4",
						"policy_name": name,
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

// Case l3策略类型测试 6912  twin
func TestAccAliCloudDdosBgpPolicy_basic6912_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ddos_bgp_policy.default"
	ra := resourceAttrInit(resourceId, AlicloudDdosBgpPolicyMap6912)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DdosBgpServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDdosBgpPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testacc_bgp%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudDdosBgpPolicyBasicDependence6912)
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
					"type": "l3",
					"content": []map[string]interface{}{
						{
							"black_ip_list_expire_at": fmt.Sprint(time.Now().Add(10 * time.Minute).Unix()),
							"intelligence_level":      "weak",
							"whiten_gfbr_nets":        "false",
							"region_block_country_list": []string{
								"${var.region_block_country_list}"},
							"region_block_province_list": []string{
								"${var.region_block_province_list}"},
							"source_limit": []map[string]interface{}{
								{
									"pps":     "32",
									"bps":     "1024",
									"syn_pps": "1",
									"syn_bps": "1024",
								},
							},
							"reflect_block_udp_port_list": []string{
								"${var.reflect_block_udp_port_list}", "${var.reflect_block_udp_port_list-update}", "${var.reflect_block_udp_port_list2}"},
							"enable_intelligence": "false",
							"enable_drop_icmp":    "false",
							"port_rule_list": []map[string]interface{}{
								{
									"port_rule_id":   "1111",
									"protocol":       "tcp",
									"src_port_start": "0",
									"src_port_end":   "65535",
									"dst_port_start": "0",
									"dst_port_end":   "65531",
									"seq_no":         "2",
									"match_action":   "drop",
								},
								{
									"port_rule_id":   "2222",
									"protocol":       "tcp",
									"src_port_start": "2",
									"src_port_end":   "3",
									"dst_port_start": "4",
									"dst_port_end":   "5",
									"match_action":   "drop",
									"seq_no":         "3",
								},
								{
									"port_rule_id":   "3333",
									"protocol":       "tcp",
									"src_port_start": "4",
									"src_port_end":   "5",
									"dst_port_start": "5",
									"dst_port_end":   "6",
									"match_action":   "drop",
									"seq_no":         "3",
								},
							},
							"finger_print_rule_list": []map[string]interface{}{
								{
									"protocol":             "tcp",
									"src_port_start":       "0",
									"src_port_end":         "65535",
									"dst_port_start":       "0",
									"dst_port_end":         "65535",
									"min_pkt_len":          "11",
									"max_pkt_len":          "22",
									"offset":               "3",
									"payload_bytes":        "11111111",
									"match_action":         "accept",
									"rate_value":           "1",
									"seq_no":               "1",
									"finger_print_rule_id": "1",
								},
								{
									"protocol":             "tcp",
									"src_port_start":       "0",
									"src_port_end":         "65535",
									"dst_port_start":       "0",
									"dst_port_end":         "65535",
									"min_pkt_len":          "11",
									"max_pkt_len":          "22",
									"offset":               "3",
									"payload_bytes":        "11111111",
									"match_action":         "accept",
									"rate_value":           "1",
									"seq_no":               "1",
									"finger_print_rule_id": "2",
								},
								{
									"protocol":             "tcp",
									"src_port_start":       "0",
									"src_port_end":         "65535",
									"dst_port_start":       "0",
									"dst_port_end":         "65535",
									"min_pkt_len":          "11",
									"max_pkt_len":          "22",
									"offset":               "3",
									"payload_bytes":        "11111111",
									"match_action":         "accept",
									"rate_value":           "1",
									"seq_no":               "1",
									"finger_print_rule_id": "3",
								},
							},
							"source_block_list": []map[string]interface{}{
								{
									"type":                 "4",
									"block_expire_seconds": "60",
									"every_seconds":        "720",
									"exceed_limit_times":   "5",
								},
								{
									"type":                 "5",
									"block_expire_seconds": "80",
									"every_seconds":        "80",
									"exceed_limit_times":   "22",
								},
								{
									"type":                 "6",
									"block_expire_seconds": "80",
									"every_seconds":        "66",
									"exceed_limit_times":   "77",
								},
							},
						},
					},
					"policy_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"type":        "l3",
						"policy_name": name,
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

// Case 更新l4类型智能开关 6856  twin
func TestAccAliCloudDdosBgpPolicy_basic6856_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ddos_bgp_policy.default"
	ra := resourceAttrInit(resourceId, AlicloudDdosBgpPolicyMap6856)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DdosBgpServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDdosBgpPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testacc_bgp%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudDdosBgpPolicyBasicDependence6856)
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
					"type":        "l4",
					"policy_name": name,
					"content": []map[string]interface{}{
						{
							"enable_defense": "false",
							"layer4_rule_list": []map[string]interface{}{
								{
									"name":     "11",
									"priority": "10",
									"method":   "hex",
									"match":    "1",
									"action":   "1",
									"limited":  "0",
									"condition_list": []map[string]interface{}{
										{
											"arg":      "3C",
											"position": "1",
											"depth":    "2",
										},
									},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"type":        "l4",
						"policy_name": name,
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

// Case l3策略类型测试_2 7021  raw
func TestAccAliCloudDdosBgpPolicy_basic7021_raw(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ddos_bgp_policy.default"
	ra := resourceAttrInit(resourceId, AlicloudDdosBgpPolicyMap7021)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DdosBgpServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDdosBgpPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testacc_bgp%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudDdosBgpPolicyBasicDependence7021)
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
					"type": "l3",
					"content": []map[string]interface{}{
						{
							"black_ip_list_expire_at": fmt.Sprint(time.Now().Add(10 * time.Minute).Unix()),
							"intelligence_level":      "weak",
							"whiten_gfbr_nets":        "false",
							"region_block_country_list": []string{
								"${var.region_block_country_list}"},
							"region_block_province_list": []string{
								"${var.region_block_province_list}"},
							"source_limit": []map[string]interface{}{
								{
									"pps":     "32",
									"bps":     "1024",
									"syn_pps": "1",
									"syn_bps": "1024",
								},
							},
							"reflect_block_udp_port_list": []string{
								"${var.reflect_block_udp_port_list}", "${var.reflect_block_udp_port_list-update}", "${var.reflect_block_udp_port_list2}"},
							"enable_intelligence": "false",
							"enable_drop_icmp":    "false",
							"port_rule_list": []map[string]interface{}{
								{
									"port_rule_id":   "1111",
									"protocol":       "tcp",
									"src_port_start": "0",
									"src_port_end":   "65535",
									"dst_port_start": "0",
									"dst_port_end":   "65531",
									"seq_no":         "2",
									"match_action":   "drop",
								},
								{
									"port_rule_id":   "2222",
									"protocol":       "tcp",
									"src_port_start": "2",
									"src_port_end":   "3",
									"dst_port_start": "4",
									"dst_port_end":   "5",
									"match_action":   "drop",
									"seq_no":         "3",
								},
								{
									"port_rule_id":   "3333",
									"protocol":       "tcp",
									"src_port_start": "4",
									"src_port_end":   "5",
									"dst_port_start": "5",
									"dst_port_end":   "6",
									"match_action":   "drop",
									"seq_no":         "3",
								},
							},
							"finger_print_rule_list": []map[string]interface{}{
								{
									"protocol":             "tcp",
									"src_port_start":       "0",
									"src_port_end":         "65535",
									"dst_port_start":       "0",
									"dst_port_end":         "65535",
									"min_pkt_len":          "11",
									"max_pkt_len":          "22",
									"offset":               "3",
									"payload_bytes":        "11111111",
									"match_action":         "accept",
									"rate_value":           "1",
									"seq_no":               "1",
									"finger_print_rule_id": "1",
								},
								{
									"protocol":             "tcp",
									"src_port_start":       "0",
									"src_port_end":         "65535",
									"dst_port_start":       "0",
									"dst_port_end":         "65535",
									"min_pkt_len":          "11",
									"max_pkt_len":          "22",
									"offset":               "3",
									"payload_bytes":        "11111111",
									"match_action":         "accept",
									"rate_value":           "1",
									"seq_no":               "1",
									"finger_print_rule_id": "2",
								},
								{
									"protocol":             "tcp",
									"src_port_start":       "0",
									"src_port_end":         "65535",
									"dst_port_start":       "0",
									"dst_port_end":         "65535",
									"min_pkt_len":          "11",
									"max_pkt_len":          "22",
									"offset":               "3",
									"payload_bytes":        "11111111",
									"match_action":         "accept",
									"rate_value":           "1",
									"seq_no":               "1",
									"finger_print_rule_id": "3",
								},
							},
							"source_block_list": []map[string]interface{}{
								{
									"type":                 "4",
									"block_expire_seconds": "60",
									"every_seconds":        "720",
									"exceed_limit_times":   "5",
								},
								{
									"type":                 "5",
									"block_expire_seconds": "80",
									"every_seconds":        "80",
									"exceed_limit_times":   "22",
								},
								{
									"type":                 "6",
									"block_expire_seconds": "80",
									"every_seconds":        "66",
									"exceed_limit_times":   "77",
								},
							},
						},
					},
					"policy_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"type":        "l3",
						"policy_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"content": []map[string]interface{}{
						{
							"black_ip_list_expire_at": fmt.Sprint(time.Now().Add(10 * time.Minute).Unix()),
							"whiten_gfbr_nets":        "true",
							"enable_drop_icmp":        "true",
							"region_block_country_list": []string{
								"${var.region_block_country_list-update}", "${var.region_block_country_list}", "${var.region_block_country_list2}"},
							"region_block_province_list": []string{
								"${var.region_block_province_list-_update}", "${var.region_block_province_list}", "${var.region_block_province_list-_update1}"},
							"intelligence_level": "default",
							"port_rule_list": []map[string]interface{}{
								{
									"protocol":       "udp",
									"src_port_start": "3",
									"src_port_end":   "65533",
									"dst_port_start": "1",
									"dst_port_end":   "65535",
									"match_action":   "drop",
									"seq_no":         "4",
									"port_rule_id":   "2222",
								},
							},
							"finger_print_rule_list": []map[string]interface{}{
								{
									"finger_print_rule_id": "111",
									"protocol":             "udp",
									"src_port_start":       "2",
									"src_port_end":         "65534",
									"dst_port_start":       "1",
									"dst_port_end":         "65534",
									"min_pkt_len":          "12",
									"max_pkt_len":          "21",
									"offset":               "4",
									"payload_bytes":        "112221",
									"match_action":         "ip_rate",
									"rate_value":           "12",
									"seq_no":               "2",
								},
							},
							"source_limit": []map[string]interface{}{
								{
									"pps":     "33",
									"bps":     "1066",
									"syn_pps": "22",
									"syn_bps": "1056",
								},
							},
							"source_block_list": []map[string]interface{}{
								{
									"type":                 "5",
									"block_expire_seconds": "70",
									"every_seconds":        "80",
									"exceed_limit_times":   "6",
								},
							},
							"enable_intelligence": "true",
							"reflect_block_udp_port_list": []string{
								"${var.reflect_block_udp_port_list-update}", "${var.reflect_block_udp_port_list}"},
						},
					},
					"policy_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"content": []map[string]interface{}{
						{
							"intelligence_level": "default",
							"whiten_gfbr_nets":   "true",
							"enable_drop_icmp":   "true",
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

// Case 更新l4类型智能开关_2 7022  raw
func TestAccAliCloudDdosBgpPolicy_basic7022_raw(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ddos_bgp_policy.default"
	ra := resourceAttrInit(resourceId, AlicloudDdosBgpPolicyMap7022)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DdosBgpServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDdosBgpPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testacc_bgp%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudDdosBgpPolicyBasicDependence7022)
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
					"type":        "l4",
					"policy_name": name,
					"content": []map[string]interface{}{
						{
							"enable_defense": "false",
							"layer4_rule_list": []map[string]interface{}{
								{
									"name":     "11",
									"priority": "10",
									"method":   "hex",
									"match":    "1",
									"action":   "1",
									"limited":  "0",
									"condition_list": []map[string]interface{}{
										{
											"arg":      "3C",
											"position": "1",
											"depth":    "2",
										},
									},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"type":        "l4",
						"policy_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"content": []map[string]interface{}{
						{
							"enable_defense": "true",
							"layer4_rule_list": []map[string]interface{}{
								{
									"name":     "11",
									"priority": "34",
									"method":   "hex",
									"match":    "0",
									"action":   "2",
									"limited":  "4",
									"condition_list": []map[string]interface{}{
										{
											"arg":      "9C",
											"position": "2",
											"depth":    "3",
										},
									},
								},
								{
									"name":     "333",
									"priority": "34",
									"method":   "hex",
									"match":    "0",
									"action":   "1",
									"limited":  "4",
									"condition_list": []map[string]interface{}{
										{
											"arg":      "9C",
											"position": "2",
											"depth":    "3",
										},
									},
								},
								{
									"name":     "55",
									"priority": "34",
									"method":   "char",
									"match":    "0",
									"action":   "2",
									"limited":  "4",
									"condition_list": []map[string]interface{}{
										{
											"arg":      "3333",
											"position": "2",
											"depth":    "3",
										},
									},
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
				Config: testAccConfig(map[string]interface{}{
					"content": []map[string]interface{}{
						{
							"enable_defense": "true",
							"layer4_rule_list": []map[string]interface{}{
								{
									"name":     "11",
									"priority": "34",
									"method":   "hex",
									"match":    "0",
									"action":   "1",
									"limited":  "4",
									"condition_list": []map[string]interface{}{
										{
											"arg":      "9C",
											"position": "2",
											"depth":    "3",
										},
									},
								},
								{
									"name":     "333",
									"priority": "34",
									"method":   "hex",
									"match":    "0",
									"action":   "2",
									"limited":  "4",
									"condition_list": []map[string]interface{}{
										{
											"arg":      "9C",
											"position": "2",
											"depth":    "3",
										},
									},
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

// Case l3策略类型测试 6912  raw
func TestAccAliCloudDdosBgpPolicy_basic6912_raw(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ddos_bgp_policy.default"
	ra := resourceAttrInit(resourceId, AlicloudDdosBgpPolicyMap6912)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DdosBgpServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDdosBgpPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testacc_bgp%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudDdosBgpPolicyBasicDependence6912)
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
					"type": "l3",
					"content": []map[string]interface{}{
						{
							"black_ip_list_expire_at": fmt.Sprint(time.Now().Add(10 * time.Minute).Unix()),
							"intelligence_level":      "weak",
							"whiten_gfbr_nets":        "false",
							"region_block_country_list": []string{
								"${var.region_block_country_list}"},
							"region_block_province_list": []string{
								"${var.region_block_province_list}"},
							"source_limit": []map[string]interface{}{
								{
									"pps":     "32",
									"bps":     "1024",
									"syn_pps": "1",
									"syn_bps": "1024",
								},
							},
							"reflect_block_udp_port_list": []string{
								"${var.reflect_block_udp_port_list}", "${var.reflect_block_udp_port_list-update}", "${var.reflect_block_udp_port_list2}"},
							"enable_intelligence": "false",
							"enable_drop_icmp":    "false",
							"port_rule_list": []map[string]interface{}{
								{
									"port_rule_id":   "1111",
									"protocol":       "tcp",
									"src_port_start": "0",
									"src_port_end":   "65535",
									"dst_port_start": "0",
									"dst_port_end":   "65531",
									"seq_no":         "2",
									"match_action":   "drop",
								},
								{
									"port_rule_id":   "2222",
									"protocol":       "tcp",
									"src_port_start": "2",
									"src_port_end":   "3",
									"dst_port_start": "4",
									"dst_port_end":   "5",
									"match_action":   "drop",
									"seq_no":         "3",
								},
								{
									"port_rule_id":   "3333",
									"protocol":       "tcp",
									"src_port_start": "4",
									"src_port_end":   "5",
									"dst_port_start": "5",
									"dst_port_end":   "6",
									"match_action":   "drop",
									"seq_no":         "3",
								},
							},
							"finger_print_rule_list": []map[string]interface{}{
								{
									"protocol":             "tcp",
									"src_port_start":       "0",
									"src_port_end":         "65535",
									"dst_port_start":       "0",
									"dst_port_end":         "65535",
									"min_pkt_len":          "11",
									"max_pkt_len":          "22",
									"offset":               "3",
									"payload_bytes":        "11111111",
									"match_action":         "accept",
									"rate_value":           "1",
									"seq_no":               "1",
									"finger_print_rule_id": "1",
								},
								{
									"protocol":             "tcp",
									"src_port_start":       "0",
									"src_port_end":         "65535",
									"dst_port_start":       "0",
									"dst_port_end":         "65535",
									"min_pkt_len":          "11",
									"max_pkt_len":          "22",
									"offset":               "3",
									"payload_bytes":        "11111111",
									"match_action":         "accept",
									"rate_value":           "1",
									"seq_no":               "1",
									"finger_print_rule_id": "2",
								},
								{
									"protocol":             "tcp",
									"src_port_start":       "0",
									"src_port_end":         "65535",
									"dst_port_start":       "0",
									"dst_port_end":         "65535",
									"min_pkt_len":          "11",
									"max_pkt_len":          "22",
									"offset":               "3",
									"payload_bytes":        "11111111",
									"match_action":         "accept",
									"rate_value":           "1",
									"seq_no":               "1",
									"finger_print_rule_id": "3",
								},
							},
							"source_block_list": []map[string]interface{}{
								{
									"type":                 "4",
									"block_expire_seconds": "60",
									"every_seconds":        "720",
									"exceed_limit_times":   "5",
								},
								{
									"type":                 "5",
									"block_expire_seconds": "80",
									"every_seconds":        "80",
									"exceed_limit_times":   "22",
								},
								{
									"type":                 "6",
									"block_expire_seconds": "80",
									"every_seconds":        "66",
									"exceed_limit_times":   "77",
								},
							},
						},
					},
					"policy_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"type":        "l3",
						"policy_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"content": []map[string]interface{}{
						{
							"black_ip_list_expire_at": fmt.Sprint(time.Now().Add(10 * time.Minute).Unix()),
							"whiten_gfbr_nets":        "true",
							"enable_drop_icmp":        "true",
							"region_block_country_list": []string{
								"${var.region_block_country_list-update}", "${var.region_block_country_list}", "${var.region_block_country_list2}"},
							"region_block_province_list": []string{
								"${var.region_block_province_list-_update}", "${var.region_block_province_list}", "${var.region_block_province_list-_update1}"},
							"intelligence_level": "default",
							"port_rule_list": []map[string]interface{}{
								{
									"protocol":       "udp",
									"src_port_start": "3",
									"src_port_end":   "65533",
									"dst_port_start": "1",
									"dst_port_end":   "65535",
									"match_action":   "drop",
									"seq_no":         "4",
									"port_rule_id":   "2222",
								},
							},
							"finger_print_rule_list": []map[string]interface{}{
								{
									"finger_print_rule_id": "111",
									"protocol":             "udp",
									"src_port_start":       "2",
									"src_port_end":         "65534",
									"dst_port_start":       "1",
									"dst_port_end":         "65534",
									"min_pkt_len":          "12",
									"max_pkt_len":          "21",
									"offset":               "4",
									"payload_bytes":        "112221",
									"match_action":         "ip_rate",
									"rate_value":           "12",
									"seq_no":               "2",
								},
							},
							"source_limit": []map[string]interface{}{
								{
									"pps":     "33",
									"bps":     "1066",
									"syn_pps": "22",
									"syn_bps": "1056",
								},
							},
							"source_block_list": []map[string]interface{}{
								{
									"type":                 "5",
									"block_expire_seconds": "70",
									"every_seconds":        "80",
									"exceed_limit_times":   "6",
								},
							},
							"enable_intelligence": "true",
							"reflect_block_udp_port_list": []string{
								"${var.reflect_block_udp_port_list-update}", "${var.reflect_block_udp_port_list}"},
						},
					},
					"policy_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"content": []map[string]interface{}{
						{
							"intelligence_level": "default",
							"whiten_gfbr_nets":   "true",
							"enable_drop_icmp":   "true",
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

// Case 更新l4类型智能开关 6856  raw
func TestAccAliCloudDdosBgpPolicy_basic6856_raw(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ddos_bgp_policy.default"
	ra := resourceAttrInit(resourceId, AlicloudDdosBgpPolicyMap6856)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DdosBgpServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDdosBgpPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testacc_bgp%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudDdosBgpPolicyBasicDependence6856)
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
					"type":        "l4",
					"policy_name": name,
					"content": []map[string]interface{}{
						{
							"enable_defense": "false",
							"layer4_rule_list": []map[string]interface{}{
								{
									"name":     "11",
									"priority": "10",
									"method":   "hex",
									"match":    "1",
									"action":   "1",
									"limited":  "0",
									"condition_list": []map[string]interface{}{
										{
											"arg":      "3C",
											"position": "1",
											"depth":    "2",
										},
									},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"type":        "l4",
						"policy_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"content": []map[string]interface{}{
						{
							"enable_defense": "true",
							"layer4_rule_list": []map[string]interface{}{
								{
									"name":     "11",
									"priority": "34",
									"method":   "hex",
									"match":    "0",
									"action":   "2",
									"limited":  "4",
									"condition_list": []map[string]interface{}{
										{
											"arg":      "9C",
											"position": "2",
											"depth":    "3",
										},
									},
								},
								{
									"name":     "333",
									"priority": "34",
									"method":   "hex",
									"match":    "0",
									"action":   "1",
									"limited":  "4",
									"condition_list": []map[string]interface{}{
										{
											"arg":      "9C",
											"position": "2",
											"depth":    "3",
										},
									},
								},
								{
									"name":     "55",
									"priority": "34",
									"method":   "char",
									"match":    "0",
									"action":   "2",
									"limited":  "4",
									"condition_list": []map[string]interface{}{
										{
											"arg":      "3333",
											"position": "2",
											"depth":    "3",
										},
									},
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
				Config: testAccConfig(map[string]interface{}{
					"content": []map[string]interface{}{
						{
							"enable_defense": "true",
							"layer4_rule_list": []map[string]interface{}{
								{
									"name":     "11",
									"priority": "34",
									"method":   "hex",
									"match":    "0",
									"action":   "1",
									"limited":  "4",
									"condition_list": []map[string]interface{}{
										{
											"arg":      "9C",
											"position": "2",
											"depth":    "3",
										},
									},
								},
								{
									"name":     "333",
									"priority": "34",
									"method":   "hex",
									"match":    "0",
									"action":   "2",
									"limited":  "4",
									"condition_list": []map[string]interface{}{
										{
											"arg":      "9C",
											"position": "2",
											"depth":    "3",
										},
									},
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

// Test DdosBgp Policy. <<< Resource test cases, automatically generated.
