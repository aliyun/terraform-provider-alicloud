package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Arms SyntheticTask. >>> Resource test cases, automatically generated.
// Case 5711
func TestAccAliCloudArmsSyntheticTask_basic5711(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_arms_synthetic_task.default"
	ra := resourceAttrInit(resourceId, AlicloudArmsSyntheticTaskMap5711)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ArmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeArmsSyntheticTask")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sarmssynthetictask%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudArmsSyntheticTaskBasicDependence5711)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.ARMSSyncTaskSupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"monitor_conf": []map[string]interface{}{
						{
							"net_icmp": []map[string]interface{}{
								{
									"interval":        "200",
									"package_num":     "4",
									"package_size":    "64",
									"timeout":         "1000",
									"tracert_num_max": "20",
									"tracert_timeout": "1000",
									"target_url":      "www.aliyun.com",
									"split_package":   "true",
								},
							},
							"net_tcp": []map[string]interface{}{
								{
									"target_url":      "www.aliyun.com",
									"connect_times":   "1",
									"interval":        "200",
									"timeout":         "1000",
									"tracert_enable":  "true",
									"tracert_num_max": "20",
									"tracert_timeout": "1000",
								},
							},
							"net_dns": []map[string]interface{}{
								{
									"target_url":         "www.aliyun.com",
									"dns_server_ip_type": "0",
									"ns_server":          "114.114.114.114",
									"query_method":       "0",
									"timeout":            "1000",
								},
							},
							"api_http": []map[string]interface{}{
								{
									"target_url": "https://www.aliyun.com",
									"method":     "GET",
									"request_headers": map[string]interface{}{
										"key1": "value1",
									},
									"request_body": []map[string]interface{}{
										{
											"content": "test",
											"type":    "text/plain",
										},
									},
									"connect_timeout": "5000",
									"timeout":         "10000",
								},
							},
							"website": []map[string]interface{}{
								{
									"target_url":               "https://www.aliyun.com",
									"automatic_scrolling":      "0",
									"custom_header":            "0",
									"disable_cache":            "0",
									"disable_compression":      "0",
									"ignore_certificate_error": "1",
									"monitor_timeout":          "10000",
									"redirection":              "1",
									"slow_element_threshold":   "5000",
									"wait_completion_time":     "5000",
									"verify_string_blacklist":  "Error",
									"verify_string_whitelist":  "Success",
									"element_blacklist":        "a.gif",
									"dns_hijack_whitelist":     "www.aliyun.com:203.0.3.55|203.3.44.67",
									"page_tamper":              "www.aliyun.com:|/cc/bb/a.gif|/vv/bb/cc.jpg",
									"flow_hijack_jump_times":   "1",
									"flow_hijack_logo":         "senyuan",
									"custom_header_content": map[string]interface{}{
										"key1": "value1",
									},
									"filter_invalid_ip": "1",
								},
							},
							"file_download": []map[string]interface{}{
								{
									"target_url":                             "http://www.aliyun.com",
									"download_kernel":                        "1",
									"quick_protocol":                         "1",
									"connection_timeout":                     "5000",
									"monitor_timeout":                        "10000",
									"ignore_certificate_status_error":        "0",
									"ignore_certificate_untrustworthy_error": "1",
									"ignore_invalid_host_error":              "1",
									"redirection":                            "1",
									"transmission_size":                      "2048",
									"ignore_certificate_canceled_error":      "1",
									"ignore_certificate_auth_error":          "1",
									"ignore_certificate_out_of_date_error":   "1",
									"ignore_certificate_using_error":         "1",
									"verify_way":                             "1",
									"validate_keywords":                      "senyuan",
									"white_list":                             "www.aliyun.com:203.0.3.55|203.3.44.67",
									"custom_header_content": map[string]interface{}{
										"key1": "value1",
									},
								},
							},
							"stream": []map[string]interface{}{
								{
									"target_url":             "https://acd-assets.alicdn.com:443/2021productweek/week1_ys.mp4",
									"stream_type":            "0",
									"stream_monitor_timeout": "60",
									"stream_address_type":    "1",
									"player_type":            "12",
									"custom_header_content": map[string]interface{}{
										"key1": "value1",
									},
									"white_list": "www.aliyun.com:203.0.3.55|203.3.44.67",
								},
							},
						},
					},
					"task_type": "1",
					"monitors": []map[string]interface{}{
						{
							"city_code":     "1200101",
							"operator_code": "246",
							"client_type":   "4",
						},
					},
					"frequency":           "1h",
					"monitor_category":    "1",
					"synthetic_task_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"task_type":           "1",
						"monitors.#":          "1",
						"frequency":           "1h",
						"monitor_category":    "1",
						"synthetic_task_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"common_setting": []map[string]interface{}{
						{
							"ip_type":           "0",
							"trace_client_type": "0",
							"xtrace_region":     "cn-hangzhou",
							"custom_host": []map[string]interface{}{
								{
									"select_type": "0",
									"hosts": []map[string]interface{}{
										{
											"domain": "www.aliyun.com",
											"ips": []string{
												"114.114.114.114"},
											"ip_type": "0",
										},
										{
											"domain": "www.tencet.com",
											"ips": []string{
												"153.3.238.102"},
											"ip_type": "1",
										},
										{
											"domain": "www.baidu.com",
											"ips": []string{
												"153.3.238.110", "180.101.50.242", "180.101.50.188"},
											"ip_type": "0",
										},
									},
								},
							},
							"monitor_samples": "0",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"custom_period": []map[string]interface{}{
						{
							"end_hour":   "20",
							"start_hour": "2",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"available_assertions": []map[string]interface{}{
						{
							"type":     "IcmpPackLoss",
							"operator": "lt",
							"expect":   "100",
							"target":   "testKey",
						},
						{
							"type":     "IcmpPackAvgLatency",
							"operator": "lte",
							"expect":   "1000",
						},
						{
							"type":     "IcmpPackMaxLatency",
							"operator": "lte",
							"expect":   "10000",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"available_assertions.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "RUNNING",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "RUNNING",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"monitor_conf": []map[string]interface{}{
						{
							"net_icmp": []map[string]interface{}{
								{
									"interval":        "200",
									"package_num":     "4",
									"package_size":    "64",
									"timeout":         "1000",
									"tracert_num_max": "20",
									"tracert_timeout": "1000",
									"target_url":      "www.aliyun.com",
									"split_package":   "true",
								},
							},
							"net_tcp": []map[string]interface{}{
								{
									"target_url":      "www.aliyun.com",
									"connect_times":   "1",
									"interval":        "200",
									"timeout":         "1000",
									"tracert_enable":  "true",
									"tracert_num_max": "20",
									"tracert_timeout": "1000",
								},
							},
							"net_dns": []map[string]interface{}{
								{
									"target_url":         "www.aliyun.com",
									"dns_server_ip_type": "0",
									"ns_server":          "114.114.114.114",
									"query_method":       "0",
									"timeout":            "1000",
								},
							},
							"api_http": []map[string]interface{}{
								{
									"target_url": "https://www.aliyun.com",
									"method":     "GET",
									"request_headers": map[string]interface{}{
										"key1": "value1",
									},
									"request_body": []map[string]interface{}{
										{
											"content": "test",
											"type":    "text/plain",
										},
									},
									"connect_timeout": "5000",
									"timeout":         "10000",
								},
							},
							"website": []map[string]interface{}{
								{
									"target_url":               "https://www.aliyun.com",
									"automatic_scrolling":      "0",
									"custom_header":            "0",
									"disable_cache":            "0",
									"disable_compression":      "0",
									"ignore_certificate_error": "1",
									"monitor_timeout":          "10000",
									"redirection":              "1",
									"slow_element_threshold":   "5000",
									"wait_completion_time":     "5000",
									"verify_string_blacklist":  "Error",
									"verify_string_whitelist":  "Success",
									"element_blacklist":        "a.gif",
									"dns_hijack_whitelist":     "www.aliyun.com:203.0.3.55|203.3.44.67",
									"page_tamper":              "www.aliyun.com:|/cc/bb/a.gif|/vv/bb/cc.jpg",
									"flow_hijack_jump_times":   "1",
									"flow_hijack_logo":         "senyuan",
									"custom_header_content": map[string]interface{}{
										"key1": "value1",
									},
									"filter_invalid_ip": "1",
								},
							},
							"file_download": []map[string]interface{}{
								{
									"target_url":                             "http://www.aliyun.com",
									"download_kernel":                        "1",
									"quick_protocol":                         "1",
									"connection_timeout":                     "5000",
									"monitor_timeout":                        "10000",
									"ignore_certificate_status_error":        "0",
									"ignore_certificate_untrustworthy_error": "1",
									"ignore_invalid_host_error":              "1",
									"redirection":                            "1",
									"transmission_size":                      "2048",
									"ignore_certificate_canceled_error":      "1",
									"ignore_certificate_auth_error":          "1",
									"ignore_certificate_out_of_date_error":   "1",
									"ignore_certificate_using_error":         "1",
									"verify_way":                             "1",
									"validate_keywords":                      "senyuan",
									"white_list":                             "www.aliyun.com:203.0.3.55|203.3.44.67",
									"custom_header_content": map[string]interface{}{
										"key1": "value1",
									},
								},
							},
							"stream": []map[string]interface{}{
								{
									"target_url":             "https://acd-assets.alicdn.com:443/2021productweek/week1_ys.mp4",
									"stream_type":            "0",
									"stream_monitor_timeout": "60",
									"stream_address_type":    "1",
									"player_type":            "12",
									"custom_header_content": map[string]interface{}{
										"key1": "value1",
									},
									"white_list": "www.aliyun.com:203.0.3.55|203.3.44.67",
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
					"monitors": []map[string]interface{}{
						{
							"city_code":     "1200101",
							"operator_code": "246",
							"client_type":   "4",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"monitors.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"frequency": "1h",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"frequency": "1h",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"synthetic_task_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"synthetic_task_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"monitor_conf": []map[string]interface{}{
						{
							"net_icmp": []map[string]interface{}{
								{
									"interval":        "300",
									"package_num":     "2",
									"package_size":    "128",
									"timeout":         "2000",
									"tracert_enable":  "true",
									"tracert_num_max": "10",
									"tracert_timeout": "4000",
									"target_url":      "www.baidu.com",
									"split_package":   "false",
								},
							},
							"net_tcp": []map[string]interface{}{
								{
									"target_url":      "www.baidu.com",
									"connect_times":   "6",
									"interval":        "300",
									"timeout":         "3000",
									"tracert_num_max": "2",
									"tracert_timeout": "1050",
									"tracert_enable":  "false",
								},
							},
							"net_dns": []map[string]interface{}{
								{
									"target_url":         "www.baidu.com",
									"dns_server_ip_type": "1",
									"ns_server":          "61.128.114.167",
									"query_method":       "1",
									"timeout":            "5050",
								},
							},
							"api_http": []map[string]interface{}{
								{
									"target_url": "https://www.baidu.com",
									"method":     "POST",
									"request_headers": map[string]interface{}{
										"key2": "value2",
									},
									"request_body": []map[string]interface{}{
										{
											"content": "test2",
											"type":    "text/html",
										},
									},
									"connect_timeout": "6000",
									"timeout":         "10050",
								},
							},
							"website": []map[string]interface{}{
								{
									"target_url":               "http://www.baidu.com",
									"automatic_scrolling":      "1",
									"custom_header":            "1",
									"disable_cache":            "1",
									"disable_compression":      "1",
									"ignore_certificate_error": "0",
									"monitor_timeout":          "20000",
									"redirection":              "0",
									"slow_element_threshold":   "5005",
									"wait_completion_time":     "5005",
									"verify_string_blacklist":  "Failed",
									"verify_string_whitelist":  "Senyuan",
									"element_blacklist":        "a.jpg",
									"dns_hijack_whitelist":     "www.aliyun.com:203.0.3.55",
									"page_tamper":              "www.aliyun.com:|/cc/bb/a.gif",
									"flow_hijack_jump_times":   "10",
									"flow_hijack_logo":         "senyuan1",
									"custom_header_content": map[string]interface{}{
										"key2": "value2",
									},
									"filter_invalid_ip": "0",
								},
							},
							"file_download": []map[string]interface{}{
								{
									"target_url":                             "https://www.baidu.com",
									"download_kernel":                        "0",
									"quick_protocol":                         "2",
									"connection_timeout":                     "6090",
									"monitor_timeout":                        "1050",
									"ignore_certificate_status_error":        "1",
									"ignore_certificate_untrustworthy_error": "0",
									"ignore_invalid_host_error":              "0",
									"redirection":                            "0",
									"transmission_size":                      "128",
									"ignore_certificate_canceled_error":      "0",
									"ignore_certificate_auth_error":          "0",
									"ignore_certificate_out_of_date_error":   "0",
									"ignore_certificate_using_error":         "0",
									"verify_way":                             "0",
									"validate_keywords":                      "senyuan1",
									"white_list":                             "www.aliyun.com:203.0.3.55",
									"custom_header_content": map[string]interface{}{
										"key2": "value2",
									},
								},
							},
							"stream": []map[string]interface{}{
								{
									"target_url":             "https://acd-assets.alicdn.com:443/2021productweek/week1_s.mp4",
									"stream_type":            "1",
									"stream_monitor_timeout": "10",
									"stream_address_type":    "0",
									"player_type":            "2",
									"custom_header_content": map[string]interface{}{
										"key2": "value2",
									},
									"white_list": "www.aliyun.com:203.0.3.55",
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
					"monitors": []map[string]interface{}{
						{
							"city_code":     "1102001",
							"operator_code": "18",
							"client_type":   "1",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"monitors.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"frequency": "12h",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"frequency": "12h",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"synthetic_task_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"synthetic_task_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"common_setting": []map[string]interface{}{
						{
							"ip_type":           "1",
							"is_open_trace":     "true",
							"trace_client_type": "1",
							"xtrace_region":     "cn-beijing",
							"custom_host": []map[string]interface{}{
								{
									"select_type": "1",
									"hosts": []map[string]interface{}{
										{
											"domain": "www.a.baidu.com",
											"ips": []string{
												"153.3.238.102"},
											"ip_type": "1",
										},
										{
											"domain": "www.shifen.com",
											"ips": []string{
												"153.3.238.110"},
											"ip_type": "1",
										},
									},
								},
							},
							"monitor_samples": "1",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"custom_period": []map[string]interface{}{
						{
							"end_hour":   "12",
							"start_hour": "11",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"available_assertions": []map[string]interface{}{
						{
							"type":     "IcmpPackLoss",
							"operator": "neq",
							"expect":   "200",
							"target":   "test",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"available_assertions.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "STOP",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "STOP",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"monitor_conf": []map[string]interface{}{
						{
							"net_icmp": []map[string]interface{}{
								{
									"interval":        "200",
									"package_num":     "36",
									"package_size":    "512",
									"timeout":         "1000",
									"tracert_num_max": "1",
									"tracert_timeout": "1200",
									"target_url":      "www.aliyun.com",
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
					"monitors": []map[string]interface{}{
						{
							"city_code":     "1101702",
							"operator_code": "2",
							"client_type":   "1",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"monitors.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"frequency": "1h",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"frequency": "1h",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"synthetic_task_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"synthetic_task_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "RUNNING",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "RUNNING",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"monitor_conf": []map[string]interface{}{
						{
							"net_icmp": []map[string]interface{}{
								{
									"interval":        "200",
									"package_num":     "4",
									"package_size":    "64",
									"timeout":         "1000",
									"tracert_num_max": "20",
									"tracert_timeout": "1000",
									"target_url":      "www.aliyun.com",
									"split_package":   "true",
								},
							},
							"net_tcp": []map[string]interface{}{
								{
									"target_url":      "www.aliyun.com",
									"connect_times":   "1",
									"interval":        "200",
									"timeout":         "1000",
									"tracert_enable":  "true",
									"tracert_num_max": "20",
									"tracert_timeout": "1000",
								},
							},
							"net_dns": []map[string]interface{}{
								{
									"target_url":         "www.aliyun.com",
									"dns_server_ip_type": "0",
									"ns_server":          "114.114.114.114",
									"query_method":       "0",
									"timeout":            "1000",
								},
							},
							"api_http": []map[string]interface{}{
								{
									"target_url": "https://www.aliyun.com",
									"method":     "GET",
									"request_headers": map[string]interface{}{
										"key1": "value1",
									},
									"request_body": []map[string]interface{}{
										{
											"content": "test",
											"type":    "text/plain",
										},
									},
									"connect_timeout": "5000",
									"timeout":         "10000",
								},
							},
							"website": []map[string]interface{}{
								{
									"target_url":               "https://www.aliyun.com",
									"automatic_scrolling":      "0",
									"custom_header":            "0",
									"disable_cache":            "0",
									"disable_compression":      "0",
									"ignore_certificate_error": "1",
									"monitor_timeout":          "10000",
									"redirection":              "1",
									"slow_element_threshold":   "5000",
									"wait_completion_time":     "5000",
									"verify_string_blacklist":  "Error",
									"verify_string_whitelist":  "Success",
									"element_blacklist":        "a.gif",
									"dns_hijack_whitelist":     "www.aliyun.com:203.0.3.55|203.3.44.67",
									"page_tamper":              "www.aliyun.com:|/cc/bb/a.gif|/vv/bb/cc.jpg",
									"flow_hijack_jump_times":   "1",
									"flow_hijack_logo":         "senyuan",
									"custom_header_content": map[string]interface{}{
										"key1": "value1",
									},
									"filter_invalid_ip": "1",
								},
							},
							"file_download": []map[string]interface{}{
								{
									"target_url":                             "http://www.aliyun.com",
									"download_kernel":                        "1",
									"quick_protocol":                         "1",
									"connection_timeout":                     "5000",
									"monitor_timeout":                        "10000",
									"ignore_certificate_status_error":        "0",
									"ignore_certificate_untrustworthy_error": "1",
									"ignore_invalid_host_error":              "1",
									"redirection":                            "1",
									"transmission_size":                      "2048",
									"ignore_certificate_canceled_error":      "1",
									"ignore_certificate_auth_error":          "1",
									"ignore_certificate_out_of_date_error":   "1",
									"ignore_certificate_using_error":         "1",
									"verify_way":                             "1",
									"validate_keywords":                      "senyuan",
									"white_list":                             "www.aliyun.com:203.0.3.55|203.3.44.67",
									"custom_header_content": map[string]interface{}{
										"key1": "value1",
									},
								},
							},
							"stream": []map[string]interface{}{
								{
									"target_url":             "https://acd-assets.alicdn.com:443/2021productweek/week1_ys.mp4",
									"stream_type":            "0",
									"stream_monitor_timeout": "60",
									"stream_address_type":    "1",
									"player_type":            "12",
									"custom_header_content": map[string]interface{}{
										"key1": "value1",
									},
									"white_list": "www.aliyun.com:203.0.3.55|203.3.44.67",
								},
							},
						},
					},
					"task_type": "1",
					"monitors": []map[string]interface{}{
						{
							"city_code":     "1200101",
							"operator_code": "246",
							"client_type":   "4",
						},
					},
					"frequency":           "1h",
					"monitor_category":    "1",
					"synthetic_task_name": name + "_update",
					"common_setting": []map[string]interface{}{
						{
							"ip_type":           "0",
							"trace_client_type": "0",
							"xtrace_region":     "cn-hangzhou",
							"custom_host": []map[string]interface{}{
								{
									"select_type": "0",
									"hosts": []map[string]interface{}{
										{
											"domain": "www.aliyun.com",
											"ips": []string{
												"114.114.114.114"},
											"ip_type": "0",
										},
										{
											"domain": "www.tencet.com",
											"ips": []string{
												"153.3.238.102"},
											"ip_type": "1",
										},
										{
											"domain": "www.baidu.com",
											"ips": []string{
												"153.3.238.110", "180.101.50.242", "180.101.50.188"},
											"ip_type": "0",
										},
									},
								},
							},
							"monitor_samples": "0",
						},
					},
					"custom_period": []map[string]interface{}{
						{
							"end_hour":   "20",
							"start_hour": "2",
						},
					},
					"available_assertions": []map[string]interface{}{
						{
							"type":     "IcmpPackLoss",
							"operator": "lt",
							"expect":   "100",
							"target":   "testKey",
						},
						{
							"type":     "IcmpPackAvgLatency",
							"operator": "lte",
							"expect":   "1000",
						},
						{
							"type":     "IcmpPackMaxLatency",
							"operator": "lte",
							"expect":   "10000",
						},
					},
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"status":            "RUNNING",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"task_type":              "1",
						"monitors.#":             "1",
						"frequency":              "1h",
						"monitor_category":       "1",
						"synthetic_task_name":    name + "_update",
						"available_assertions.#": "3",
						"resource_group_id":      CHECKSET,
						"status":                 "RUNNING",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
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

var AlicloudArmsSyntheticTaskMap5711 = map[string]string{
	"status": CHECKSET,
}

func AlicloudArmsSyntheticTaskBasicDependence5711(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {}


`, name)
}

// Case 5711  twin
func TestAccAliCloudArmsSyntheticTask_basic5711_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_arms_synthetic_task.default"
	ra := resourceAttrInit(resourceId, AlicloudArmsSyntheticTaskMap5711)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ArmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeArmsSyntheticTask")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sarmssynthetictask%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudArmsSyntheticTaskBasicDependence5711)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.ARMSSyncTaskSupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"monitor_conf": []map[string]interface{}{
						{
							"net_icmp": []map[string]interface{}{
								{
									"interval":        "200",
									"package_num":     "36",
									"package_size":    "512",
									"timeout":         "1000",
									"tracert_enable":  "true",
									"tracert_num_max": "1",
									"tracert_timeout": "1200",
									"target_url":      "www.aliyun.com",
								},
							},
							"net_tcp": []map[string]interface{}{
								{
									"target_url":      "www.baidu.com",
									"connect_times":   "6",
									"interval":        "300",
									"timeout":         "3000",
									"tracert_num_max": "2",
									"tracert_timeout": "1050",
								},
							},
							"net_dns": []map[string]interface{}{
								{
									"target_url":         "www.baidu.com",
									"dns_server_ip_type": "1",
									"ns_server":          "61.128.114.167",
									"query_method":       "1",
									"timeout":            "5050",
								},
							},
							"api_http": []map[string]interface{}{
								{
									"target_url": "https://www.baidu.com",
									"method":     "POST",
									"request_headers": map[string]interface{}{
										"key1": "value1",
									},
									"request_body": []map[string]interface{}{
										{
											"content": "test2",
											"type":    "text/html",
										},
									},
									"connect_timeout": "6000",
									"timeout":         "10050",
								},
							},
							"website": []map[string]interface{}{
								{
									"target_url":               "http://www.baidu.com",
									"automatic_scrolling":      "1",
									"custom_header":            "1",
									"disable_cache":            "1",
									"disable_compression":      "1",
									"ignore_certificate_error": "0",
									"monitor_timeout":          "20000",
									"redirection":              "0",
									"slow_element_threshold":   "5005",
									"wait_completion_time":     "5005",
									"verify_string_blacklist":  "Failed",
									"verify_string_whitelist":  "Senyuan",
									"element_blacklist":        "a.jpg",
									"dns_hijack_whitelist":     "www.aliyun.com:203.0.3.55",
									"page_tamper":              "www.aliyun.com:|/cc/bb/a.gif",
									"flow_hijack_jump_times":   "10",
									"flow_hijack_logo":         "senyuan1",
									"custom_header_content": map[string]interface{}{
										"key1": "value1",
									},
									"filter_invalid_ip": "0",
								},
							},
							"file_download": []map[string]interface{}{
								{
									"target_url":                             "https://www.baidu.com",
									"download_kernel":                        "0",
									"quick_protocol":                         "2",
									"connection_timeout":                     "6090",
									"monitor_timeout":                        "1050",
									"ignore_certificate_status_error":        "1",
									"ignore_certificate_untrustworthy_error": "0",
									"ignore_invalid_host_error":              "0",
									"redirection":                            "0",
									"transmission_size":                      "128",
									"ignore_certificate_canceled_error":      "0",
									"ignore_certificate_auth_error":          "0",
									"ignore_certificate_out_of_date_error":   "0",
									"ignore_certificate_using_error":         "0",
									"verify_way":                             "0",
									"validate_keywords":                      "senyuan1",
									"white_list":                             "www.aliyun.com:203.0.3.55",
									"custom_header_content": map[string]interface{}{
										"key1": "value1",
									},
								},
							},
							"stream": []map[string]interface{}{
								{
									"target_url":             "https://acd-assets.alicdn.com:443/2021productweek/week1_s.mp4",
									"stream_type":            "1",
									"stream_monitor_timeout": "10",
									"stream_address_type":    "0",
									"player_type":            "2",
									"custom_header_content": map[string]interface{}{
										"key1": "value1",
									},
									"white_list": "www.aliyun.com:203.0.3.55",
								},
							},
						},
					},
					"task_type": "1",
					"monitors": []map[string]interface{}{
						{
							"city_code":     "1200101",
							"operator_code": "246",
							"client_type":   "4",
						},
					},
					"frequency":           "1h",
					"monitor_category":    "1",
					"synthetic_task_name": name,
					"common_setting": []map[string]interface{}{
						{
							"ip_type":           "1",
							"is_open_trace":     "true",
							"trace_client_type": "1",
							"xtrace_region":     "cn-beijing",
							"custom_host": []map[string]interface{}{
								{
									"select_type": "1",
									"hosts": []map[string]interface{}{
										{
											"domain": "www.a.baidu.com",
											"ips": []string{
												"153.3.238.102"},
											"ip_type": "0",
										},
										{
											"domain": "www.shifen.com",
											"ips": []string{
												"153.3.238.110", "114.114.114.114", "127.0.0.1"},
											"ip_type": "1",
										},
										{
											"domain": "www.baidu.com",
											"ips": []string{
												"153.3.238.110", "180.101.50.242", "180.101.50.188"},
											"ip_type": "0",
										},
									},
								},
							},
							"monitor_samples": "1",
						},
					},
					"custom_period": []map[string]interface{}{
						{
							"end_hour":   "12",
							"start_hour": "11",
						},
					},
					"available_assertions": []map[string]interface{}{
						{
							"type":     "IcmpPackLoss",
							"operator": "neq",
							"expect":   "200",
							"target":   "test",
						},
						{
							"type":     "IcmpPackAvgLatency",
							"operator": "lte",
							"expect":   "1000",
						},
						{
							"type":     "IcmpPackMaxLatency",
							"operator": "lte",
							"expect":   "10000",
						},
					},
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"status":            "RUNNING",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"task_type":              "1",
						"monitors.#":             "1",
						"frequency":              "1h",
						"monitor_category":       "1",
						"synthetic_task_name":    name,
						"available_assertions.#": "3",
						"resource_group_id":      CHECKSET,
						"status":                 "RUNNING",
						"tags.%":                 "2",
						"tags.Created":           "TF",
						"tags.For":               "Test",
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

// Test Arms SyntheticTask. <<< Resource test cases, automatically generated.
