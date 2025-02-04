// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudArmsSyntheticTask() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudArmsSyntheticTaskCreate,
		Read:   resourceAliCloudArmsSyntheticTaskRead,
		Update: resourceAliCloudArmsSyntheticTaskUpdate,
		Delete: resourceAliCloudArmsSyntheticTaskDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"available_assertions": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"operator": {
							Type:     schema.TypeString,
							Required: true,
						},
						"type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"target": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"expect": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"common_setting": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip_type": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"is_open_trace": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"trace_client_type": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"xtrace_region": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"custom_host": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"hosts": {
										Type:     schema.TypeList,
										Required: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"ip_type": {
													Type:     schema.TypeInt,
													Required: true,
												},
												"domain": {
													Type:     schema.TypeString,
													Required: true,
												},
												"ips": {
													Type:     schema.TypeList,
													Required: true,
													Elem:     &schema.Schema{Type: schema.TypeString},
												},
											},
										},
									},
									"select_type": {
										Type:     schema.TypeInt,
										Required: true,
									},
								},
							},
						},
						"monitor_samples": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
			"custom_period": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"start_hour": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"end_hour": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
			"frequency": {
				Type:     schema.TypeString,
				Required: true,
			},
			"monitor_category": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"monitor_conf": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"net_dns": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"query_method": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"timeout": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"ns_server": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"dns_server_ip_type": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"target_url": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
						"file_download": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"connection_timeout": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"custom_header_content": {
										Type:     schema.TypeMap,
										Optional: true,
									},
									"transmission_size": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"ignore_certificate_out_of_date_error": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"ignore_certificate_using_error": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"white_list": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"verify_way": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"target_url": {
										Type:     schema.TypeString,
										Required: true,
									},
									"ignore_certificate_status_error": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"ignore_certificate_canceled_error": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"monitor_timeout": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"validate_keywords": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"ignore_certificate_untrustworthy_error": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"redirection": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"ignore_certificate_auth_error": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"quick_protocol": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"download_kernel": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"ignore_invalid_host_error": {
										Type:     schema.TypeInt,
										Optional: true,
									},
								},
							},
						},
						"website": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"disable_compression": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"custom_header_content": {
										Type:     schema.TypeMap,
										Optional: true,
									},
									"automatic_scrolling": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"slow_element_threshold": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"verify_string_whitelist": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"ignore_certificate_error": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"target_url": {
										Type:     schema.TypeString,
										Required: true,
									},
									"flow_hijack_jump_times": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"disable_cache": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"monitor_timeout": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"wait_completion_time": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"verify_string_blacklist": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"flow_hijack_logo": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"element_blacklist": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"redirection": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"custom_header": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"dns_hijack_whitelist": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"page_tamper": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"filter_invalid_ip": {
										Type:     schema.TypeInt,
										Optional: true,
									},
								},
							},
						},
						"stream": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"stream_type": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"custom_header_content": {
										Type:     schema.TypeMap,
										Optional: true,
									},
									"stream_address_type": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"stream_monitor_timeout": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"player_type": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"white_list": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"target_url": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"net_tcp": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"tracert_enable": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"timeout": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"tracert_timeout": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"target_url": {
										Type:     schema.TypeString,
										Required: true,
									},
									"tracert_num_max": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"connect_times": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"interval": {
										Type:     schema.TypeInt,
										Optional: true,
									},
								},
							},
						},
						"net_icmp": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"tracert_enable": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"package_num": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"package_size": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"split_package": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"timeout": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"tracert_timeout": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"tracert_num_max": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"target_url": {
										Type:     schema.TypeString,
										Required: true,
									},
									"interval": {
										Type:     schema.TypeInt,
										Optional: true,
									},
								},
							},
						},
						"api_http": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"connect_timeout": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"timeout": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"target_url": {
										Type:     schema.TypeString,
										Required: true,
									},
									"method": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"request_headers": {
										Type:     schema.TypeMap,
										Optional: true,
									},
									"request_body": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"type": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"content": {
													Type:     schema.TypeString,
													Optional: true,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"monitors": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"city_code": {
							Type:     schema.TypeString,
							Required: true,
						},
						"client_type": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"operator_code": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"synthetic_task_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"tags": tagsSchema(),
			"task_type": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudArmsSyntheticTaskCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateTimingSyntheticTask"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId

	request["TaskType"] = d.Get("task_type")
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOk("tags"); ok {
		tagsMap := ConvertTags(v.(map[string]interface{}))
		request["Tags"] = tagsMap
	}

	if v, ok := d.GetOk("monitors"); ok {
		monitorsMaps := make([]map[string]interface{}, 0)
		for _, dataLoop1 := range v.([]interface{}) {
			dataLoop1Tmp := dataLoop1.(map[string]interface{})
			dataLoop1Map := make(map[string]interface{})
			dataLoop1Map["OperatorCode"] = dataLoop1Tmp["operator_code"]
			dataLoop1Map["CityCode"] = dataLoop1Tmp["city_code"]
			dataLoop1Map["ClientType"] = dataLoop1Tmp["client_type"]
			monitorsMaps = append(monitorsMaps, dataLoop1Map)
		}
		request["Monitors"], _ = convertListMapToJsonString(monitorsMaps)
	}

	request["MonitorCategory"] = d.Get("monitor_category")
	request["Name"] = d.Get("synthetic_task_name")
	request["Frequency"] = d.Get("frequency")
	objectDataLocalMap := make(map[string]interface{})
	if v := d.Get("monitor_conf"); !IsNil(v) {
		netICMP_map, _ := jsonpath.Get("$[0].net_icmp[0]", v)
		if !IsNil(netICMP_map) {
			netICMP := make(map[string]interface{})
			nodeNative5, _ := jsonpath.Get("$[0].net_icmp[0].package_num", d.Get("monitor_conf"))
			if nodeNative5 != "" {
				netICMP["PackageNum"] = nodeNative5
			}
			nodeNative6, _ := jsonpath.Get("$[0].net_icmp[0].package_size", d.Get("monitor_conf"))
			if nodeNative6 != "" {
				netICMP["PackageSize"] = nodeNative6
			}
			nodeNative7, _ := jsonpath.Get("$[0].net_icmp[0].split_package", d.Get("monitor_conf"))
			if nodeNative7 != "" {
				netICMP["SplitPackage"] = nodeNative7
			}
			nodeNative8, _ := jsonpath.Get("$[0].net_icmp[0].interval", d.Get("monitor_conf"))
			if nodeNative8 != "" {
				netICMP["Interval"] = nodeNative8
			}
			nodeNative9, _ := jsonpath.Get("$[0].net_icmp[0].tracert_enable", d.Get("monitor_conf"))
			if nodeNative9 != "" {
				netICMP["TracertEnable"] = nodeNative9
			}
			nodeNative10, _ := jsonpath.Get("$[0].net_icmp[0].tracert_num_max", d.Get("monitor_conf"))
			if nodeNative10 != "" {
				netICMP["TracertNumMax"] = nodeNative10
			}
			nodeNative11, _ := jsonpath.Get("$[0].net_icmp[0].tracert_timeout", d.Get("monitor_conf"))
			if nodeNative11 != "" {
				netICMP["TracertTimeout"] = nodeNative11
			}
			nodeNative12, _ := jsonpath.Get("$[0].net_icmp[0].timeout", d.Get("monitor_conf"))
			if nodeNative12 != "" {
				netICMP["Timeout"] = nodeNative12
			}
			nodeNative13, _ := jsonpath.Get("$[0].net_icmp[0].target_url", d.Get("monitor_conf"))
			if nodeNative13 != "" {
				netICMP["TargetUrl"] = nodeNative13
			}
			objectDataLocalMap["NetICMP"] = netICMP
		}
		netTCP_map, _ := jsonpath.Get("$[0].net_tcp[0]", v)
		if !IsNil(netTCP_map) {
			netTCP := make(map[string]interface{})
			nodeNative14, _ := jsonpath.Get("$[0].net_tcp[0].connect_times", d.Get("monitor_conf"))
			if nodeNative14 != "" {
				netTCP["ConnectTimes"] = nodeNative14
			}
			nodeNative15, _ := jsonpath.Get("$[0].net_tcp[0].interval", d.Get("monitor_conf"))
			if nodeNative15 != "" {
				netTCP["Interval"] = nodeNative15
			}
			nodeNative16, _ := jsonpath.Get("$[0].net_tcp[0].tracert_enable", d.Get("monitor_conf"))
			if nodeNative16 != "" {
				netTCP["TracertEnable"] = nodeNative16
			}
			nodeNative17, _ := jsonpath.Get("$[0].net_tcp[0].tracert_num_max", d.Get("monitor_conf"))
			if nodeNative17 != "" {
				netTCP["TracertNumMax"] = nodeNative17
			}
			nodeNative18, _ := jsonpath.Get("$[0].net_tcp[0].tracert_timeout", d.Get("monitor_conf"))
			if nodeNative18 != "" {
				netTCP["TracertTimeout"] = nodeNative18
			}
			nodeNative19, _ := jsonpath.Get("$[0].net_tcp[0].timeout", d.Get("monitor_conf"))
			if nodeNative19 != "" {
				netTCP["Timeout"] = nodeNative19
			}
			nodeNative20, _ := jsonpath.Get("$[0].net_tcp[0].target_url", d.Get("monitor_conf"))
			if nodeNative20 != "" {
				netTCP["TargetUrl"] = nodeNative20
			}
			objectDataLocalMap["NetTCP"] = netTCP
		}
		netDNS_map, _ := jsonpath.Get("$[0].net_dns[0]", v)
		if !IsNil(netDNS_map) {
			netDNS := make(map[string]interface{})
			nodeNative21, _ := jsonpath.Get("$[0].net_dns[0].dns_server_ip_type", d.Get("monitor_conf"))
			if nodeNative21 != "" {
				netDNS["DnsServerIpType"] = nodeNative21
			}
			nodeNative22, _ := jsonpath.Get("$[0].net_dns[0].ns_server", d.Get("monitor_conf"))
			if nodeNative22 != "" {
				netDNS["NsServer"] = nodeNative22
			}
			nodeNative23, _ := jsonpath.Get("$[0].net_dns[0].query_method", d.Get("monitor_conf"))
			if nodeNative23 != "" {
				netDNS["QueryMethod"] = nodeNative23
			}
			nodeNative24, _ := jsonpath.Get("$[0].net_dns[0].timeout", d.Get("monitor_conf"))
			if nodeNative24 != "" {
				netDNS["Timeout"] = nodeNative24
			}
			nodeNative25, _ := jsonpath.Get("$[0].net_dns[0].target_url", d.Get("monitor_conf"))
			if nodeNative25 != "" {
				netDNS["TargetUrl"] = nodeNative25
			}
			objectDataLocalMap["NetDNS"] = netDNS
		}
		apiHTTP_map, _ := jsonpath.Get("$[0].api_http[0]", v)
		if !IsNil(apiHTTP_map) {
			apiHTTP := make(map[string]interface{})
			nodeNative26, _ := jsonpath.Get("$[0].api_http[0].target_url", d.Get("monitor_conf"))
			if nodeNative26 != "" {
				apiHTTP["TargetUrl"] = nodeNative26
			}
			nodeNative27, _ := jsonpath.Get("$[0].api_http[0].method", d.Get("monitor_conf"))
			if nodeNative27 != "" {
				apiHTTP["Method"] = nodeNative27
			}
			requestBody_map, _ := jsonpath.Get("$[0].api_http[0].request_body[0]", v)
			if !IsNil(requestBody_map) {
				requestBody := make(map[string]interface{})
				nodeNative28, _ := jsonpath.Get("$[0].api_http[0].request_body[0].content", d.Get("monitor_conf"))
				if nodeNative28 != "" {
					requestBody["Content"] = nodeNative28
				}
				nodeNative29, _ := jsonpath.Get("$[0].api_http[0].request_body[0].type", d.Get("monitor_conf"))
				if nodeNative29 != "" {
					requestBody["Type"] = nodeNative29
				}
				apiHTTP["RequestBody"] = requestBody
			}
			nodeNative30, _ := jsonpath.Get("$[0].api_http[0].connect_timeout", d.Get("monitor_conf"))
			if nodeNative30 != "" {
				apiHTTP["ConnectTimeout"] = nodeNative30
			}
			nodeNative31, _ := jsonpath.Get("$[0].api_http[0].request_headers", d.Get("monitor_conf"))
			if nodeNative31 != "" {
				apiHTTP["RequestHeaders"] = nodeNative31
			}
			nodeNative32, _ := jsonpath.Get("$[0].api_http[0].timeout", d.Get("monitor_conf"))
			if nodeNative32 != "" {
				apiHTTP["Timeout"] = nodeNative32
			}
			objectDataLocalMap["ApiHTTP"] = apiHTTP
		}
		website_map, _ := jsonpath.Get("$[0].website[0]", v)
		if !IsNil(website_map) {
			website := make(map[string]interface{})
			nodeNative33, _ := jsonpath.Get("$[0].website[0].automatic_scrolling", d.Get("monitor_conf"))
			if nodeNative33 != "" {
				website["AutomaticScrolling"] = nodeNative33
			}
			nodeNative34, _ := jsonpath.Get("$[0].website[0].custom_header", d.Get("monitor_conf"))
			if nodeNative34 != "" {
				website["CustomHeader"] = nodeNative34
			}
			nodeNative35, _ := jsonpath.Get("$[0].website[0].disable_cache", d.Get("monitor_conf"))
			if nodeNative35 != "" {
				website["DisableCache"] = nodeNative35
			}
			nodeNative36, _ := jsonpath.Get("$[0].website[0].disable_compression", d.Get("monitor_conf"))
			if nodeNative36 != "" {
				website["DisableCompression"] = nodeNative36
			}
			nodeNative37, _ := jsonpath.Get("$[0].website[0].filter_invalid_ip", d.Get("monitor_conf"))
			if nodeNative37 != "" {
				website["FilterInvalidIP"] = nodeNative37
			}
			nodeNative38, _ := jsonpath.Get("$[0].website[0].ignore_certificate_error", d.Get("monitor_conf"))
			if nodeNative38 != "" {
				website["IgnoreCertificateError"] = nodeNative38
			}
			nodeNative39, _ := jsonpath.Get("$[0].website[0].slow_element_threshold", d.Get("monitor_conf"))
			if nodeNative39 != "" {
				website["SlowElementThreshold"] = nodeNative39
			}
			nodeNative40, _ := jsonpath.Get("$[0].website[0].wait_completion_time", d.Get("monitor_conf"))
			if nodeNative40 != "" {
				website["WaitCompletionTime"] = nodeNative40
			}
			nodeNative41, _ := jsonpath.Get("$[0].website[0].verify_string_blacklist", d.Get("monitor_conf"))
			if nodeNative41 != "" {
				website["VerifyStringBlacklist"] = nodeNative41
			}
			nodeNative42, _ := jsonpath.Get("$[0].website[0].verify_string_whitelist", d.Get("monitor_conf"))
			if nodeNative42 != "" {
				website["VerifyStringWhitelist"] = nodeNative42
			}
			nodeNative43, _ := jsonpath.Get("$[0].website[0].element_blacklist", d.Get("monitor_conf"))
			if nodeNative43 != "" {
				website["ElementBlacklist"] = nodeNative43
			}
			nodeNative44, _ := jsonpath.Get("$[0].website[0].dns_hijack_whitelist", d.Get("monitor_conf"))
			if nodeNative44 != "" {
				website["DNSHijackWhitelist"] = nodeNative44
			}
			nodeNative45, _ := jsonpath.Get("$[0].website[0].page_tamper", d.Get("monitor_conf"))
			if nodeNative45 != "" {
				website["PageTamper"] = nodeNative45
			}
			nodeNative46, _ := jsonpath.Get("$[0].website[0].flow_hijack_jump_times", d.Get("monitor_conf"))
			if nodeNative46 != "" {
				website["FlowHijackJumpTimes"] = nodeNative46
			}
			nodeNative47, _ := jsonpath.Get("$[0].website[0].flow_hijack_logo", d.Get("monitor_conf"))
			if nodeNative47 != "" {
				website["FlowHijackLogo"] = nodeNative47
			}
			nodeNative48, _ := jsonpath.Get("$[0].website[0].target_url", d.Get("monitor_conf"))
			if nodeNative48 != "" {
				website["TargetUrl"] = nodeNative48
			}
			nodeNative49, _ := jsonpath.Get("$[0].website[0].custom_header_content", d.Get("monitor_conf"))
			if nodeNative49 != "" {
				website["CustomHeaderContent"] = nodeNative49
			}
			nodeNative50, _ := jsonpath.Get("$[0].website[0].monitor_timeout", d.Get("monitor_conf"))
			if nodeNative50 != "" {
				website["MonitorTimeout"] = nodeNative50
			}
			nodeNative51, _ := jsonpath.Get("$[0].website[0].redirection", d.Get("monitor_conf"))
			if nodeNative51 != "" {
				website["Redirection"] = nodeNative51
			}
			objectDataLocalMap["Website"] = website
		}
		fileDownload_map, _ := jsonpath.Get("$[0].file_download[0]", v)
		if !IsNil(fileDownload_map) {
			fileDownload := make(map[string]interface{})
			nodeNative52, _ := jsonpath.Get("$[0].file_download[0].download_kernel", d.Get("monitor_conf"))
			if nodeNative52 != "" {
				fileDownload["DownloadKernel"] = nodeNative52
			}
			nodeNative53, _ := jsonpath.Get("$[0].file_download[0].quick_protocol", d.Get("monitor_conf"))
			if nodeNative53 != "" {
				fileDownload["QuickProtocol"] = nodeNative53
			}
			nodeNative54, _ := jsonpath.Get("$[0].file_download[0].connection_timeout", d.Get("monitor_conf"))
			if nodeNative54 != "" {
				fileDownload["ConnectionTimeout"] = nodeNative54
			}
			nodeNative55, _ := jsonpath.Get("$[0].file_download[0].ignore_certificate_status_error", d.Get("monitor_conf"))
			if nodeNative55 != "" {
				fileDownload["IgnoreCertificateStatusError"] = nodeNative55
			}
			nodeNative56, _ := jsonpath.Get("$[0].file_download[0].ignore_certificate_untrustworthy_error", d.Get("monitor_conf"))
			if nodeNative56 != "" {
				fileDownload["IgnoreCertificateUntrustworthyError"] = nodeNative56
			}
			nodeNative57, _ := jsonpath.Get("$[0].file_download[0].ignore_invalid_host_error", d.Get("monitor_conf"))
			if nodeNative57 != "" {
				fileDownload["IgnoreInvalidHostError"] = nodeNative57
			}
			nodeNative58, _ := jsonpath.Get("$[0].file_download[0].transmission_size", d.Get("monitor_conf"))
			if nodeNative58 != "" {
				fileDownload["TransmissionSize"] = nodeNative58
			}
			nodeNative59, _ := jsonpath.Get("$[0].file_download[0].ignore_certificate_canceled_error", d.Get("monitor_conf"))
			if nodeNative59 != "" {
				fileDownload["IgnoreCertificateCanceledError"] = nodeNative59
			}
			nodeNative60, _ := jsonpath.Get("$[0].file_download[0].ignore_certificate_auth_error", d.Get("monitor_conf"))
			if nodeNative60 != "" {
				fileDownload["IgnoreCertificateAuthError"] = nodeNative60
			}
			nodeNative61, _ := jsonpath.Get("$[0].file_download[0].ignore_certificate_out_of_date_error", d.Get("monitor_conf"))
			if nodeNative61 != "" {
				fileDownload["IgnoreCertificateOutOfDateError"] = nodeNative61
			}
			nodeNative62, _ := jsonpath.Get("$[0].file_download[0].ignore_certificate_using_error", d.Get("monitor_conf"))
			if nodeNative62 != "" {
				fileDownload["IgnoreCertificateUsingError"] = nodeNative62
			}
			nodeNative63, _ := jsonpath.Get("$[0].file_download[0].verify_way", d.Get("monitor_conf"))
			if nodeNative63 != "" {
				fileDownload["VerifyWay"] = nodeNative63
			}
			nodeNative64, _ := jsonpath.Get("$[0].file_download[0].validate_keywords", d.Get("monitor_conf"))
			if nodeNative64 != "" {
				fileDownload["ValidateKeywords"] = nodeNative64
			}
			nodeNative65, _ := jsonpath.Get("$[0].file_download[0].target_url", d.Get("monitor_conf"))
			if nodeNative65 != "" {
				fileDownload["TargetUrl"] = nodeNative65
			}
			nodeNative66, _ := jsonpath.Get("$[0].file_download[0].monitor_timeout", d.Get("monitor_conf"))
			if nodeNative66 != "" {
				fileDownload["MonitorTimeout"] = nodeNative66
			}
			nodeNative67, _ := jsonpath.Get("$[0].file_download[0].custom_header_content", d.Get("monitor_conf"))
			if nodeNative67 != "" {
				fileDownload["CustomHeaderContent"] = nodeNative67
			}
			nodeNative68, _ := jsonpath.Get("$[0].file_download[0].redirection", d.Get("monitor_conf"))
			if nodeNative68 != "" {
				fileDownload["Redirection"] = nodeNative68
			}
			nodeNative69, _ := jsonpath.Get("$[0].file_download[0].white_list", d.Get("monitor_conf"))
			if nodeNative69 != "" {
				fileDownload["WhiteList"] = nodeNative69
			}
			objectDataLocalMap["FileDownload"] = fileDownload
		}
		stream_map, _ := jsonpath.Get("$[0].stream[0]", v)
		if !IsNil(stream_map) {
			stream := make(map[string]interface{})
			nodeNative70, _ := jsonpath.Get("$[0].stream[0].stream_type", d.Get("monitor_conf"))
			if nodeNative70 != "" {
				stream["StreamType"] = nodeNative70
			}
			nodeNative71, _ := jsonpath.Get("$[0].stream[0].stream_monitor_timeout", d.Get("monitor_conf"))
			if nodeNative71 != "" {
				stream["StreamMonitorTimeout"] = nodeNative71
			}
			nodeNative72, _ := jsonpath.Get("$[0].stream[0].stream_address_type", d.Get("monitor_conf"))
			if nodeNative72 != "" {
				stream["StreamAddressType"] = nodeNative72
			}
			nodeNative73, _ := jsonpath.Get("$[0].stream[0].player_type", d.Get("monitor_conf"))
			if nodeNative73 != "" {
				stream["PlayerType"] = nodeNative73
			}
			nodeNative74, _ := jsonpath.Get("$[0].stream[0].white_list", d.Get("monitor_conf"))
			if nodeNative74 != "" {
				stream["WhiteList"] = nodeNative74
			}
			nodeNative75, _ := jsonpath.Get("$[0].stream[0].custom_header_content", d.Get("monitor_conf"))
			if nodeNative75 != "" {
				stream["CustomHeaderContent"] = nodeNative75
			}
			nodeNative76, _ := jsonpath.Get("$[0].stream[0].target_url", d.Get("monitor_conf"))
			if nodeNative76 != "" {
				stream["TargetUrl"] = nodeNative76
			}
			objectDataLocalMap["Stream"] = stream
		}
		request["MonitorConf"] = convertMapToJsonStringIgnoreError(objectDataLocalMap)
	}

	objectDataLocalMap1 := make(map[string]interface{})
	if v := d.Get("custom_period"); !IsNil(v) {
		nodeNative77, _ := jsonpath.Get("$[0].end_hour", d.Get("custom_period"))
		if nodeNative77 != "" {
			objectDataLocalMap1["EndHour"] = nodeNative77
		}
		nodeNative78, _ := jsonpath.Get("$[0].start_hour", d.Get("custom_period"))
		if nodeNative78 != "" {
			objectDataLocalMap1["StartHour"] = nodeNative78
		}
		request["CustomPeriod"] = convertMapToJsonStringIgnoreError(objectDataLocalMap1)
	}

	objectDataLocalMap2 := make(map[string]interface{})
	if v := d.Get("common_setting"); !IsNil(v) {
		customHost_map, _ := jsonpath.Get("$[0].custom_host[0]", v)
		if !IsNil(customHost_map) {
			customHost := make(map[string]interface{})
			if v, ok := d.GetOk("common_setting"); ok {
				localData2, err := jsonpath.Get("$[0].custom_host[0].hosts", v)
				if err != nil {
					return WrapError(err)
				}
				localMaps := make([]map[string]interface{}, 0)
				for _, dataLoop2 := range localData2.([]interface{}) {
					dataLoop2Tmp := dataLoop2.(map[string]interface{})
					dataLoop2Map := make(map[string]interface{})
					dataLoop2Map["Domain"] = dataLoop2Tmp["domain"]
					dataLoop2Map["Ips"] = dataLoop2Tmp["ips"]
					dataLoop2Map["IpType"] = dataLoop2Tmp["ip_type"]
					localMaps = append(localMaps, dataLoop2Map)
				}
				customHost["Hosts"] = localMaps
			}
			nodeNative82, _ := jsonpath.Get("$[0].custom_host[0].select_type", d.Get("common_setting"))
			if nodeNative82 != "" {
				customHost["SelectType"] = nodeNative82
			}
			objectDataLocalMap2["CustomHost"] = customHost
		}
		nodeNative83, _ := jsonpath.Get("$[0].monitor_samples", d.Get("common_setting"))
		if nodeNative83 != "" {
			objectDataLocalMap2["MonitorSamples"] = nodeNative83
		}
		nodeNative84, _ := jsonpath.Get("$[0].is_open_trace", d.Get("common_setting"))
		if nodeNative84 != "" {
			objectDataLocalMap2["IsOpenTrace"] = nodeNative84
		}
		nodeNative85, _ := jsonpath.Get("$[0].trace_client_type", d.Get("common_setting"))
		if nodeNative85 != "" {
			objectDataLocalMap2["TraceClientType"] = nodeNative85
		}
		nodeNative86, _ := jsonpath.Get("$[0].xtrace_region", d.Get("common_setting"))
		if nodeNative86 != "" {
			objectDataLocalMap2["XtraceRegion"] = nodeNative86
		}
		nodeNative87, _ := jsonpath.Get("$[0].ip_type", d.Get("common_setting"))
		if nodeNative87 != "" {
			objectDataLocalMap2["IpType"] = nodeNative87
		}
		request["CommonSetting"] = convertMapToJsonStringIgnoreError(objectDataLocalMap2)
	}

	if v, ok := d.GetOk("available_assertions"); ok {
		availableAssertionsMaps := make([]map[string]interface{}, 0)
		for _, dataLoop3 := range v.([]interface{}) {
			dataLoop3Tmp := dataLoop3.(map[string]interface{})
			dataLoop3Map := make(map[string]interface{})
			dataLoop3Map["Type"] = dataLoop3Tmp["type"]
			dataLoop3Map["Target"] = dataLoop3Tmp["target"]
			dataLoop3Map["Operator"] = dataLoop3Tmp["operator"]
			dataLoop3Map["Expect"] = dataLoop3Tmp["expect"]
			availableAssertionsMaps = append(availableAssertionsMaps, dataLoop3Map)
		}
		request["AvailableAssertions"], _ = convertListMapToJsonString(availableAssertionsMaps)
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("ARMS", "2019-08-08", action, query, request, true)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_arms_synthetic_task", action, AlibabaCloudSdkGoERROR)
	}

	id, _ := jsonpath.Get("$.Data.TaskId", response)
	d.SetId(fmt.Sprint(id))

	return resourceAliCloudArmsSyntheticTaskUpdate(d, meta)
}

func resourceAliCloudArmsSyntheticTaskRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	armsServiceV2 := ArmsServiceV2{client}

	objectRaw, err := armsServiceV2.DescribeArmsSyntheticTask(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_arms_synthetic_task DescribeArmsSyntheticTask Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("frequency", objectRaw["Frequency"])
	d.Set("monitor_category", objectRaw["MonitorCategory"])
	d.Set("resource_group_id", objectRaw["ResourceGroupId"])
	d.Set("status", objectRaw["Status"])
	d.Set("synthetic_task_name", objectRaw["Name"])
	d.Set("task_type", objectRaw["TaskType"])

	availableAssertions1Raw := objectRaw["AvailableAssertions"]
	availableAssertionsMaps := make([]map[string]interface{}, 0)
	if availableAssertions1Raw != nil {
		for _, availableAssertionsChild1Raw := range availableAssertions1Raw.([]interface{}) {
			availableAssertionsMap := make(map[string]interface{})
			availableAssertionsChild1Raw := availableAssertionsChild1Raw.(map[string]interface{})
			availableAssertionsMap["expect"] = availableAssertionsChild1Raw["Expect"]
			availableAssertionsMap["operator"] = availableAssertionsChild1Raw["Operator"]
			availableAssertionsMap["target"] = availableAssertionsChild1Raw["Target"]
			availableAssertionsMap["type"] = availableAssertionsChild1Raw["Type"]

			availableAssertionsMaps = append(availableAssertionsMaps, availableAssertionsMap)
		}
	}
	d.Set("available_assertions", availableAssertionsMaps)
	commonSettingMaps := make([]map[string]interface{}, 0)
	commonSettingMap := make(map[string]interface{})
	commonSetting1Raw := make(map[string]interface{})
	if objectRaw["CommonSetting"] != nil {
		commonSetting1Raw = objectRaw["CommonSetting"].(map[string]interface{})
	}
	if len(commonSetting1Raw) > 0 {
		commonSettingMap["custom_host"] = commonSetting1Raw["CustomHost"]
		commonSettingMap["ip_type"] = commonSetting1Raw["IpType"]
		commonSettingMap["is_open_trace"] = commonSetting1Raw["IsOpenTrace"]
		commonSettingMap["monitor_samples"] = commonSetting1Raw["MonitorSamples"]
		commonSettingMap["trace_client_type"] = commonSetting1Raw["TraceClientType"]
		commonSettingMap["xtrace_region"] = commonSetting1Raw["XtraceRegion"]

		customHostMaps := make([]map[string]interface{}, 0)
		customHostMap := make(map[string]interface{})
		customHost3Raw := make(map[string]interface{})
		if commonSetting1Raw["CustomHost"] != nil {
			customHost3Raw = commonSetting1Raw["CustomHost"].(map[string]interface{})
		}
		if len(customHost3Raw) > 0 {
			customHostMap["hosts"] = customHost3Raw["Hosts"]
			customHostMap["select_type"] = customHost3Raw["SelectType"]

			hosts3Raw := customHost3Raw["Hosts"]
			hostsMaps := make([]map[string]interface{}, 0)
			if hosts3Raw != nil {
				for _, hostsChild1Raw := range hosts3Raw.([]interface{}) {
					hostsMap := make(map[string]interface{})
					hostsChild1Raw := hostsChild1Raw.(map[string]interface{})
					hostsMap["domain"] = hostsChild1Raw["Domain"]
					hostsMap["ip_type"] = hostsChild1Raw["IpType"]

					ips1Raw := make([]interface{}, 0)
					if hostsChild1Raw["Ips"] != nil {
						ips1Raw = hostsChild1Raw["Ips"].([]interface{})
					}

					hostsMap["ips"] = ips1Raw
					hostsMaps = append(hostsMaps, hostsMap)
				}
			}
			customHostMap["hosts"] = hostsMaps
			customHostMaps = append(customHostMaps, customHostMap)
		}
		commonSettingMap["custom_host"] = customHostMaps
		commonSettingMaps = append(commonSettingMaps, commonSettingMap)
	}
	d.Set("common_setting", commonSettingMaps)
	customPeriodMaps := make([]map[string]interface{}, 0)
	customPeriodMap := make(map[string]interface{})
	customPeriod1Raw := make(map[string]interface{})
	if objectRaw["CustomPeriod"] != nil {
		customPeriod1Raw = objectRaw["CustomPeriod"].(map[string]interface{})
	}
	if len(customPeriod1Raw) > 0 {
		customPeriodMap["end_hour"] = customPeriod1Raw["EndHour"]
		customPeriodMap["start_hour"] = customPeriod1Raw["StartHour"]

		customPeriodMaps = append(customPeriodMaps, customPeriodMap)
	}
	d.Set("custom_period", customPeriodMaps)
	monitorConfMaps := make([]map[string]interface{}, 0)
	monitorConfMap := make(map[string]interface{})
	monitorConf1Raw := make(map[string]interface{})
	if objectRaw["MonitorConf"] != nil {
		monitorConf1Raw = objectRaw["MonitorConf"].(map[string]interface{})
	}
	if len(monitorConf1Raw) > 0 {

		apiHTTPMaps := make([]map[string]interface{}, 0)
		apiHTTPMap := make(map[string]interface{})
		apiHTTP1Raw := make(map[string]interface{})
		if monitorConf1Raw["ApiHTTP"] != nil {
			apiHTTP1Raw = monitorConf1Raw["ApiHTTP"].(map[string]interface{})
		}
		if len(apiHTTP1Raw) > 0 {
			apiHTTPMap["connect_timeout"] = apiHTTP1Raw["ConnectTimeout"]
			apiHTTPMap["method"] = apiHTTP1Raw["Method"]
			apiHTTPMap["request_headers"] = apiHTTP1Raw["RequestHeaders"]
			apiHTTPMap["target_url"] = apiHTTP1Raw["TargetUrl"]
			apiHTTPMap["timeout"] = apiHTTP1Raw["Timeout"]

			requestBodyMaps := make([]map[string]interface{}, 0)
			requestBodyMap := make(map[string]interface{})
			requestBody1Raw := make(map[string]interface{})
			if apiHTTP1Raw["RequestBody"] != nil {
				requestBody1Raw = apiHTTP1Raw["RequestBody"].(map[string]interface{})
			}
			if len(requestBody1Raw) > 0 {
				requestBodyMap["content"] = requestBody1Raw["Content"]
				requestBodyMap["type"] = requestBody1Raw["Type"]

				requestBodyMaps = append(requestBodyMaps, requestBodyMap)
			}
			apiHTTPMap["request_body"] = requestBodyMaps
			apiHTTPMaps = append(apiHTTPMaps, apiHTTPMap)
		}
		monitorConfMap["api_http"] = apiHTTPMaps
		fileDownloadMaps := make([]map[string]interface{}, 0)
		fileDownloadMap := make(map[string]interface{})
		fileDownload1Raw := make(map[string]interface{})
		if monitorConf1Raw["FileDownload"] != nil {
			fileDownload1Raw = monitorConf1Raw["FileDownload"].(map[string]interface{})
		}
		if len(fileDownload1Raw) > 0 {
			fileDownloadMap["connection_timeout"] = fileDownload1Raw["ConnectionTimeout"]
			fileDownloadMap["custom_header_content"] = fileDownload1Raw["CustomHeaderContent"]
			fileDownloadMap["download_kernel"] = fileDownload1Raw["DownloadKernel"]
			fileDownloadMap["ignore_certificate_auth_error"] = fileDownload1Raw["IgnoreCertificateAuthError"]
			fileDownloadMap["ignore_certificate_canceled_error"] = fileDownload1Raw["IgnoreCertificateCanceledError"]
			fileDownloadMap["ignore_certificate_out_of_date_error"] = fileDownload1Raw["IgnoreCertificateOutOfDateError"]
			fileDownloadMap["ignore_certificate_status_error"] = fileDownload1Raw["IgnoreCertificateStatusError"]
			fileDownloadMap["ignore_certificate_untrustworthy_error"] = fileDownload1Raw["IgnoreCertificateUntrustworthyError"]
			fileDownloadMap["ignore_certificate_using_error"] = fileDownload1Raw["IgnoreCertificateUsingError"]
			fileDownloadMap["ignore_invalid_host_error"] = fileDownload1Raw["IgnoreInvalidHostError"]
			fileDownloadMap["monitor_timeout"] = fileDownload1Raw["MonitorTimeout"]
			fileDownloadMap["quick_protocol"] = fileDownload1Raw["QuickProtocol"]
			fileDownloadMap["redirection"] = fileDownload1Raw["Redirection"]
			fileDownloadMap["target_url"] = fileDownload1Raw["TargetUrl"]
			fileDownloadMap["transmission_size"] = fileDownload1Raw["TransmissionSize"]
			fileDownloadMap["validate_keywords"] = fileDownload1Raw["ValidateKeywords"]
			fileDownloadMap["verify_way"] = fileDownload1Raw["VerifyWay"]
			fileDownloadMap["white_list"] = fileDownload1Raw["WhiteList"]

			fileDownloadMaps = append(fileDownloadMaps, fileDownloadMap)
		}
		monitorConfMap["file_download"] = fileDownloadMaps
		netDnsMaps := make([]map[string]interface{}, 0)
		netDnsMap := make(map[string]interface{})
		netDNS1Raw := make(map[string]interface{})
		if monitorConf1Raw["NetDNS"] != nil {
			netDNS1Raw = monitorConf1Raw["NetDNS"].(map[string]interface{})
		}
		if len(netDNS1Raw) > 0 {
			netDnsMap["dns_server_ip_type"] = netDNS1Raw["DnsServerIpType"]
			netDnsMap["ns_server"] = netDNS1Raw["NsServer"]
			netDnsMap["query_method"] = netDNS1Raw["QueryMethod"]
			netDnsMap["target_url"] = netDNS1Raw["TargetUrl"]
			netDnsMap["timeout"] = netDNS1Raw["Timeout"]

			netDnsMaps = append(netDnsMaps, netDnsMap)
		}
		monitorConfMap["net_dns"] = netDnsMaps
		netIcmpMaps := make([]map[string]interface{}, 0)
		netIcmpMap := make(map[string]interface{})
		netICMP1Raw := make(map[string]interface{})
		if monitorConf1Raw["NetICMP"] != nil {
			netICMP1Raw = monitorConf1Raw["NetICMP"].(map[string]interface{})
		}
		if len(netICMP1Raw) > 0 {
			netIcmpMap["interval"] = netICMP1Raw["Interval"]
			netIcmpMap["package_num"] = netICMP1Raw["PackageNum"]
			netIcmpMap["package_size"] = netICMP1Raw["PackageSize"]
			netIcmpMap["split_package"] = netICMP1Raw["SplitPackage"]
			netIcmpMap["target_url"] = netICMP1Raw["TargetUrl"]
			netIcmpMap["timeout"] = netICMP1Raw["Timeout"]
			netIcmpMap["tracert_enable"] = netICMP1Raw["TracertEnable"]
			netIcmpMap["tracert_num_max"] = netICMP1Raw["TracertNumMax"]
			netIcmpMap["tracert_timeout"] = netICMP1Raw["TracertTimeout"]

			netIcmpMaps = append(netIcmpMaps, netIcmpMap)
		}
		monitorConfMap["net_icmp"] = netIcmpMaps
		netTCPMaps := make([]map[string]interface{}, 0)
		netTCPMap := make(map[string]interface{})
		netTCP1Raw := make(map[string]interface{})
		if monitorConf1Raw["NetTCP"] != nil {
			netTCP1Raw = monitorConf1Raw["NetTCP"].(map[string]interface{})
		}
		if len(netTCP1Raw) > 0 {
			netTCPMap["connect_times"] = netTCP1Raw["ConnectTimes"]
			netTCPMap["interval"] = netTCP1Raw["Interval"]
			netTCPMap["target_url"] = netTCP1Raw["TargetUrl"]
			netTCPMap["timeout"] = netTCP1Raw["Timeout"]
			netTCPMap["tracert_enable"] = netTCP1Raw["TracertEnable"]
			netTCPMap["tracert_num_max"] = netTCP1Raw["TracertNumMax"]
			netTCPMap["tracert_timeout"] = netTCP1Raw["TracertTimeout"]

			netTCPMaps = append(netTCPMaps, netTCPMap)
		}
		monitorConfMap["net_tcp"] = netTCPMaps
		streamMaps := make([]map[string]interface{}, 0)
		streamMap := make(map[string]interface{})
		stream1Raw := make(map[string]interface{})
		if monitorConf1Raw["Stream"] != nil {
			stream1Raw = monitorConf1Raw["Stream"].(map[string]interface{})
		}
		if len(stream1Raw) > 0 {
			streamMap["custom_header_content"] = stream1Raw["CustomHeaderContent"]
			streamMap["player_type"] = stream1Raw["PlayerType"]
			streamMap["stream_address_type"] = stream1Raw["StreamAddressType"]
			streamMap["stream_monitor_timeout"] = stream1Raw["StreamMonitorTimeout"]
			streamMap["stream_type"] = stream1Raw["StreamType"]
			streamMap["target_url"] = stream1Raw["TargetUrl"]
			streamMap["white_list"] = stream1Raw["WhiteList"]

			streamMaps = append(streamMaps, streamMap)
		}
		monitorConfMap["stream"] = streamMaps
		websiteMaps := make([]map[string]interface{}, 0)
		websiteMap := make(map[string]interface{})
		website1Raw := make(map[string]interface{})
		if monitorConf1Raw["Website"] != nil {
			website1Raw = monitorConf1Raw["Website"].(map[string]interface{})
		}
		if len(website1Raw) > 0 {
			websiteMap["automatic_scrolling"] = website1Raw["AutomaticScrolling"]
			websiteMap["custom_header"] = website1Raw["CustomHeader"]
			websiteMap["custom_header_content"] = website1Raw["CustomHeaderContent"]
			websiteMap["disable_cache"] = website1Raw["DisableCache"]
			websiteMap["disable_compression"] = website1Raw["DisableCompression"]
			websiteMap["dns_hijack_whitelist"] = website1Raw["DNSHijackWhitelist"]
			websiteMap["element_blacklist"] = website1Raw["ElementBlacklist"]
			websiteMap["filter_invalid_ip"] = website1Raw["FilterInvalidIP"]
			websiteMap["flow_hijack_jump_times"] = website1Raw["FlowHijackJumpTimes"]
			websiteMap["flow_hijack_logo"] = website1Raw["FlowHijackLogo"]
			websiteMap["ignore_certificate_error"] = website1Raw["IgnoreCertificateError"]
			websiteMap["monitor_timeout"] = website1Raw["MonitorTimeout"]
			websiteMap["page_tamper"] = website1Raw["PageTamper"]
			websiteMap["redirection"] = website1Raw["Redirection"]
			websiteMap["slow_element_threshold"] = website1Raw["SlowElementThreshold"]
			websiteMap["target_url"] = website1Raw["TargetUrl"]
			websiteMap["verify_string_blacklist"] = website1Raw["VerifyStringBlacklist"]
			websiteMap["verify_string_whitelist"] = website1Raw["VerifyStringWhitelist"]
			websiteMap["wait_completion_time"] = website1Raw["WaitCompletionTime"]

			websiteMaps = append(websiteMaps, websiteMap)
		}
		monitorConfMap["website"] = websiteMaps
		monitorConfMaps = append(monitorConfMaps, monitorConfMap)
	}
	d.Set("monitor_conf", monitorConfMaps)
	monitors1Raw := objectRaw["Monitors"]
	monitorsMaps := make([]map[string]interface{}, 0)
	if monitors1Raw != nil {
		for _, monitorsChild1Raw := range monitors1Raw.([]interface{}) {
			monitorsMap := make(map[string]interface{})
			monitorsChild1Raw := monitorsChild1Raw.(map[string]interface{})
			monitorsMap["city_code"] = monitorsChild1Raw["CityCode"]
			monitorsMap["client_type"] = monitorsChild1Raw["ClientType"]
			monitorsMap["operator_code"] = monitorsChild1Raw["OperatorCode"]

			monitorsMaps = append(monitorsMaps, monitorsMap)
		}
	}
	d.Set("monitors", monitorsMaps)
	tagsMaps := objectRaw["Tags"]
	d.Set("tags", tagsToMap(tagsMaps))

	return nil
}

func resourceAliCloudArmsSyntheticTaskUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	d.Partial(true)
	action := "UpdateTimingSyntheticTask"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["TaskId"] = d.Id()
	request["RegionId"] = client.RegionId
	if !d.IsNewResource() && d.HasChange("monitors") {
		update = true
	}
	if v, ok := d.GetOk("monitors"); ok {
		monitorsMaps := make([]map[string]interface{}, 0)
		for _, dataLoop := range v.([]interface{}) {
			dataLoopTmp := dataLoop.(map[string]interface{})
			dataLoopMap := make(map[string]interface{})
			dataLoopMap["CityCode"] = dataLoopTmp["city_code"]
			dataLoopMap["OperatorCode"] = dataLoopTmp["operator_code"]
			dataLoopMap["ClientType"] = dataLoopTmp["client_type"]
			monitorsMaps = append(monitorsMaps, dataLoopMap)
		}
		request["Monitors"], _ = convertListMapToJsonString(monitorsMaps)
	}

	if !d.IsNewResource() && d.HasChange("synthetic_task_name") {
		update = true
	}
	request["Name"] = d.Get("synthetic_task_name")
	if !d.IsNewResource() && d.HasChange("frequency") {
		update = true
	}
	request["Frequency"] = d.Get("frequency")
	if !d.IsNewResource() && d.HasChange("monitor_conf") {
		update = true
	}
	objectDataLocalMap := make(map[string]interface{})
	if v := d.Get("monitor_conf"); !IsNil(v) {
		stream_map, _ := jsonpath.Get("$[0].stream[0]", v)
		if !IsNil(stream_map) {
			stream := make(map[string]interface{})
			nodeNative3, _ := jsonpath.Get("$[0].stream[0].target_url", v)
			if nodeNative3 != "" {
				stream["TargetUrl"] = nodeNative3
			}
			nodeNative4, _ := jsonpath.Get("$[0].stream[0].white_list", v)
			if nodeNative4 != "" {
				stream["WhiteList"] = nodeNative4
			}
			nodeNative5, _ := jsonpath.Get("$[0].stream[0].stream_type", v)
			if nodeNative5 != "" {
				stream["StreamType"] = nodeNative5
			}
			nodeNative6, _ := jsonpath.Get("$[0].stream[0].stream_monitor_timeout", v)
			if nodeNative6 != "" {
				stream["StreamMonitorTimeout"] = nodeNative6
			}
			nodeNative7, _ := jsonpath.Get("$[0].stream[0].stream_address_type", v)
			if nodeNative7 != "" {
				stream["StreamAddressType"] = nodeNative7
			}
			nodeNative8, _ := jsonpath.Get("$[0].stream[0].player_type", v)
			if nodeNative8 != "" {
				stream["PlayerType"] = nodeNative8
			}
			nodeNative9, _ := jsonpath.Get("$[0].stream[0].custom_header_content", v)
			if nodeNative9 != "" {
				stream["CustomHeaderContent"] = nodeNative9
			}
			objectDataLocalMap["Stream"] = stream
		}
		netICMP_map, _ := jsonpath.Get("$[0].net_icmp[0]", v)
		if !IsNil(netICMP_map) {
			netICMP := make(map[string]interface{})
			nodeNative10, _ := jsonpath.Get("$[0].net_icmp[0].package_num", v)
			if nodeNative10 != "" {
				netICMP["PackageNum"] = nodeNative10
			}
			nodeNative11, _ := jsonpath.Get("$[0].net_icmp[0].package_size", v)
			if nodeNative11 != "" {
				netICMP["PackageSize"] = nodeNative11
			}
			nodeNative12, _ := jsonpath.Get("$[0].net_icmp[0].split_package", v)
			if nodeNative12 != "" {
				netICMP["SplitPackage"] = nodeNative12
			}
			nodeNative13, _ := jsonpath.Get("$[0].net_icmp[0].interval", v)
			if nodeNative13 != "" {
				netICMP["Interval"] = nodeNative13
			}
			nodeNative14, _ := jsonpath.Get("$[0].net_icmp[0].timeout", v)
			if nodeNative14 != "" {
				netICMP["Timeout"] = nodeNative14
			}
			nodeNative15, _ := jsonpath.Get("$[0].net_icmp[0].tracert_enable", v)
			if nodeNative15 != "" {
				netICMP["TracertEnable"] = nodeNative15
			}
			nodeNative16, _ := jsonpath.Get("$[0].net_icmp[0].tracert_num_max", v)
			if nodeNative16 != "" {
				netICMP["TracertNumMax"] = nodeNative16
			}
			nodeNative17, _ := jsonpath.Get("$[0].net_icmp[0].tracert_timeout", v)
			if nodeNative17 != "" {
				netICMP["TracertTimeout"] = nodeNative17
			}
			nodeNative18, _ := jsonpath.Get("$[0].net_icmp[0].target_url", v)
			if nodeNative18 != "" {
				netICMP["TargetUrl"] = nodeNative18
			}
			objectDataLocalMap["NetICMP"] = netICMP
		}
		apiHTTP_map, _ := jsonpath.Get("$[0].api_http[0]", v)
		if !IsNil(apiHTTP_map) {
			apiHTTP := make(map[string]interface{})
			nodeNative19, _ := jsonpath.Get("$[0].api_http[0].timeout", v)
			if nodeNative19 != "" {
				apiHTTP["Timeout"] = nodeNative19
			}
			nodeNative20, _ := jsonpath.Get("$[0].api_http[0].target_url", v)
			if nodeNative20 != "" {
				apiHTTP["TargetUrl"] = nodeNative20
			}
			nodeNative21, _ := jsonpath.Get("$[0].api_http[0].method", v)
			if nodeNative21 != "" {
				apiHTTP["Method"] = nodeNative21
			}
			requestBody_map, _ := jsonpath.Get("$[0].api_http[0].request_body[0]", v)
			if !IsNil(requestBody_map) {
				requestBody := make(map[string]interface{})
				nodeNative22, _ := jsonpath.Get("$[0].api_http[0].request_body[0].content", v)
				if nodeNative22 != "" {
					requestBody["Content"] = nodeNative22
				}
				nodeNative23, _ := jsonpath.Get("$[0].api_http[0].request_body[0].type", v)
				if nodeNative23 != "" {
					requestBody["Type"] = nodeNative23
				}
				apiHTTP["RequestBody"] = requestBody
			}
			nodeNative24, _ := jsonpath.Get("$[0].api_http[0].connect_timeout", v)
			if nodeNative24 != "" {
				apiHTTP["ConnectTimeout"] = nodeNative24
			}
			nodeNative25, _ := jsonpath.Get("$[0].api_http[0].request_headers", v)
			if nodeNative25 != "" {
				apiHTTP["RequestHeaders"] = nodeNative25
			}
			objectDataLocalMap["ApiHTTP"] = apiHTTP
		}
		fileDownload_map, _ := jsonpath.Get("$[0].file_download[0]", v)
		if !IsNil(fileDownload_map) {
			fileDownload := make(map[string]interface{})
			nodeNative26, _ := jsonpath.Get("$[0].file_download[0].target_url", v)
			if nodeNative26 != "" {
				fileDownload["TargetUrl"] = nodeNative26
			}
			nodeNative27, _ := jsonpath.Get("$[0].file_download[0].monitor_timeout", v)
			if nodeNative27 != "" {
				fileDownload["MonitorTimeout"] = nodeNative27
			}
			nodeNative28, _ := jsonpath.Get("$[0].file_download[0].redirection", v)
			if nodeNative28 != "" {
				fileDownload["Redirection"] = nodeNative28
			}
			nodeNative29, _ := jsonpath.Get("$[0].file_download[0].download_kernel", v)
			if nodeNative29 != "" {
				fileDownload["DownloadKernel"] = nodeNative29
			}
			nodeNative30, _ := jsonpath.Get("$[0].file_download[0].quick_protocol", v)
			if nodeNative30 != "" {
				fileDownload["QuickProtocol"] = nodeNative30
			}
			nodeNative31, _ := jsonpath.Get("$[0].file_download[0].connection_timeout", v)
			if nodeNative31 != "" {
				fileDownload["ConnectionTimeout"] = nodeNative31
			}
			nodeNative32, _ := jsonpath.Get("$[0].file_download[0].ignore_certificate_status_error", v)
			if nodeNative32 != "" {
				fileDownload["IgnoreCertificateStatusError"] = nodeNative32
			}
			nodeNative33, _ := jsonpath.Get("$[0].file_download[0].ignore_certificate_untrustworthy_error", v)
			if nodeNative33 != "" {
				fileDownload["IgnoreCertificateUntrustworthyError"] = nodeNative33
			}
			nodeNative34, _ := jsonpath.Get("$[0].file_download[0].ignore_invalid_host_error", v)
			if nodeNative34 != "" {
				fileDownload["IgnoreInvalidHostError"] = nodeNative34
			}
			nodeNative35, _ := jsonpath.Get("$[0].file_download[0].transmission_size", v)
			if nodeNative35 != "" {
				fileDownload["TransmissionSize"] = nodeNative35
			}
			nodeNative36, _ := jsonpath.Get("$[0].file_download[0].ignore_certificate_canceled_error", v)
			if nodeNative36 != "" {
				fileDownload["IgnoreCertificateCanceledError"] = nodeNative36
			}
			nodeNative37, _ := jsonpath.Get("$[0].file_download[0].ignore_certificate_auth_error", v)
			if nodeNative37 != "" {
				fileDownload["IgnoreCertificateAuthError"] = nodeNative37
			}
			nodeNative38, _ := jsonpath.Get("$[0].file_download[0].ignore_certificate_out_of_date_error", v)
			if nodeNative38 != "" {
				fileDownload["IgnoreCertificateOutOfDateError"] = nodeNative38
			}
			nodeNative39, _ := jsonpath.Get("$[0].file_download[0].ignore_certificate_using_error", v)
			if nodeNative39 != "" {
				fileDownload["IgnoreCertificateUsingError"] = nodeNative39
			}
			nodeNative40, _ := jsonpath.Get("$[0].file_download[0].verify_way", v)
			if nodeNative40 != "" {
				fileDownload["VerifyWay"] = nodeNative40
			}
			nodeNative41, _ := jsonpath.Get("$[0].file_download[0].validate_keywords", v)
			if nodeNative41 != "" {
				fileDownload["ValidateKeywords"] = nodeNative41
			}
			nodeNative42, _ := jsonpath.Get("$[0].file_download[0].white_list", v)
			if nodeNative42 != "" {
				fileDownload["WhiteList"] = nodeNative42
			}
			nodeNative43, _ := jsonpath.Get("$[0].file_download[0].custom_header_content", v)
			if nodeNative43 != "" {
				fileDownload["CustomHeaderContent"] = nodeNative43
			}
			objectDataLocalMap["FileDownload"] = fileDownload
		}
		netTCP_map, _ := jsonpath.Get("$[0].net_tcp[0]", v)
		if !IsNil(netTCP_map) {
			netTCP := make(map[string]interface{})
			nodeNative44, _ := jsonpath.Get("$[0].net_tcp[0].connect_times", v)
			if nodeNative44 != "" {
				netTCP["ConnectTimes"] = nodeNative44
			}
			nodeNative45, _ := jsonpath.Get("$[0].net_tcp[0].target_url", v)
			if nodeNative45 != "" {
				netTCP["TargetUrl"] = nodeNative45
			}
			nodeNative46, _ := jsonpath.Get("$[0].net_tcp[0].interval", v)
			if nodeNative46 != "" {
				netTCP["Interval"] = nodeNative46
			}
			nodeNative47, _ := jsonpath.Get("$[0].net_tcp[0].timeout", v)
			if nodeNative47 != "" {
				netTCP["Timeout"] = nodeNative47
			}
			nodeNative48, _ := jsonpath.Get("$[0].net_tcp[0].tracert_enable", v)
			if nodeNative48 != "" {
				netTCP["TracertEnable"] = nodeNative48
			}
			nodeNative49, _ := jsonpath.Get("$[0].net_tcp[0].tracert_num_max", v)
			if nodeNative49 != "" {
				netTCP["TracertNumMax"] = nodeNative49
			}
			nodeNative50, _ := jsonpath.Get("$[0].net_tcp[0].tracert_timeout", v)
			if nodeNative50 != "" {
				netTCP["TracertTimeout"] = nodeNative50
			}
			objectDataLocalMap["NetTCP"] = netTCP
		}
		netDNS_map, _ := jsonpath.Get("$[0].net_dns[0]", v)
		if !IsNil(netDNS_map) {
			netDNS := make(map[string]interface{})
			nodeNative51, _ := jsonpath.Get("$[0].net_dns[0].timeout", v)
			if nodeNative51 != "" {
				netDNS["Timeout"] = nodeNative51
			}
			nodeNative52, _ := jsonpath.Get("$[0].net_dns[0].dns_server_ip_type", v)
			if nodeNative52 != "" {
				netDNS["DnsServerIpType"] = nodeNative52
			}
			nodeNative53, _ := jsonpath.Get("$[0].net_dns[0].ns_server", v)
			if nodeNative53 != "" {
				netDNS["NsServer"] = nodeNative53
			}
			nodeNative54, _ := jsonpath.Get("$[0].net_dns[0].query_method", v)
			if nodeNative54 != "" {
				netDNS["QueryMethod"] = nodeNative54
			}
			nodeNative55, _ := jsonpath.Get("$[0].net_dns[0].target_url", v)
			if nodeNative55 != "" {
				netDNS["TargetUrl"] = nodeNative55
			}
			objectDataLocalMap["NetDNS"] = netDNS
		}
		website_map, _ := jsonpath.Get("$[0].website[0]", v)
		if !IsNil(website_map) {
			website := make(map[string]interface{})
			nodeNative56, _ := jsonpath.Get("$[0].website[0].target_url", v)
			if nodeNative56 != "" {
				website["TargetUrl"] = nodeNative56
			}
			nodeNative57, _ := jsonpath.Get("$[0].website[0].automatic_scrolling", v)
			if nodeNative57 != "" {
				website["AutomaticScrolling"] = nodeNative57
			}
			nodeNative58, _ := jsonpath.Get("$[0].website[0].custom_header", v)
			if nodeNative58 != "" {
				website["CustomHeader"] = nodeNative58
			}
			nodeNative59, _ := jsonpath.Get("$[0].website[0].disable_cache", v)
			if nodeNative59 != "" {
				website["DisableCache"] = nodeNative59
			}
			nodeNative60, _ := jsonpath.Get("$[0].website[0].disable_compression", v)
			if nodeNative60 != "" {
				website["DisableCompression"] = nodeNative60
			}
			nodeNative61, _ := jsonpath.Get("$[0].website[0].filter_invalid_ip", v)
			if nodeNative61 != "" {
				website["FilterInvalidIP"] = nodeNative61
			}
			nodeNative62, _ := jsonpath.Get("$[0].website[0].ignore_certificate_error", v)
			if nodeNative62 != "" {
				website["IgnoreCertificateError"] = nodeNative62
			}
			nodeNative63, _ := jsonpath.Get("$[0].website[0].slow_element_threshold", v)
			if nodeNative63 != "" {
				website["SlowElementThreshold"] = nodeNative63
			}
			nodeNative64, _ := jsonpath.Get("$[0].website[0].wait_completion_time", v)
			if nodeNative64 != "" {
				website["WaitCompletionTime"] = nodeNative64
			}
			nodeNative65, _ := jsonpath.Get("$[0].website[0].verify_string_blacklist", v)
			if nodeNative65 != "" {
				website["VerifyStringBlacklist"] = nodeNative65
			}
			nodeNative66, _ := jsonpath.Get("$[0].website[0].verify_string_whitelist", v)
			if nodeNative66 != "" {
				website["VerifyStringWhitelist"] = nodeNative66
			}
			nodeNative67, _ := jsonpath.Get("$[0].website[0].element_blacklist", v)
			if nodeNative67 != "" {
				website["ElementBlacklist"] = nodeNative67
			}
			nodeNative68, _ := jsonpath.Get("$[0].website[0].dns_hijack_whitelist", v)
			if nodeNative68 != "" {
				website["DNSHijackWhitelist"] = nodeNative68
			}
			nodeNative69, _ := jsonpath.Get("$[0].website[0].page_tamper", v)
			if nodeNative69 != "" {
				website["PageTamper"] = nodeNative69
			}
			nodeNative70, _ := jsonpath.Get("$[0].website[0].flow_hijack_jump_times", v)
			if nodeNative70 != "" {
				website["FlowHijackJumpTimes"] = nodeNative70
			}
			nodeNative71, _ := jsonpath.Get("$[0].website[0].flow_hijack_logo", v)
			if nodeNative71 != "" {
				website["FlowHijackLogo"] = nodeNative71
			}
			nodeNative72, _ := jsonpath.Get("$[0].website[0].monitor_timeout", v)
			if nodeNative72 != "" {
				website["MonitorTimeout"] = nodeNative72
			}
			nodeNative73, _ := jsonpath.Get("$[0].website[0].redirection", v)
			if nodeNative73 != "" {
				website["Redirection"] = nodeNative73
			}
			nodeNative74, _ := jsonpath.Get("$[0].website[0].custom_header_content", v)
			if nodeNative74 != "" {
				website["CustomHeaderContent"] = nodeNative74
			}
			objectDataLocalMap["Website"] = website
		}
		request["MonitorConf"] = convertMapToJsonStringIgnoreError(objectDataLocalMap)
	}

	if !d.IsNewResource() && d.HasChange("available_assertions") {
		update = true
	}
	if v, ok := d.GetOk("available_assertions"); ok {
		availableAssertionsMaps := make([]map[string]interface{}, 0)
		for _, dataLoop1 := range v.([]interface{}) {
			dataLoop1Tmp := dataLoop1.(map[string]interface{})
			dataLoop1Map := make(map[string]interface{})
			dataLoop1Map["Type"] = dataLoop1Tmp["type"]
			dataLoop1Map["Target"] = dataLoop1Tmp["target"]
			dataLoop1Map["Operator"] = dataLoop1Tmp["operator"]
			dataLoop1Map["Expect"] = dataLoop1Tmp["expect"]
			availableAssertionsMaps = append(availableAssertionsMaps, dataLoop1Map)
		}
		request["AvailableAssertions"], _ = convertListMapToJsonString(availableAssertionsMaps)
	}

	if !d.IsNewResource() && d.HasChange("custom_period") {
		update = true
	}
	objectDataLocalMap1 := make(map[string]interface{})
	if v := d.Get("custom_period"); !IsNil(v) {
		nodeNative79, _ := jsonpath.Get("$[0].end_hour", v)
		if nodeNative79 != "" {
			objectDataLocalMap1["EndHour"] = nodeNative79
		}
		nodeNative80, _ := jsonpath.Get("$[0].start_hour", v)
		if nodeNative80 != "" {
			objectDataLocalMap1["StartHour"] = nodeNative80
		}
		request["CustomPeriod"] = convertMapToJsonStringIgnoreError(objectDataLocalMap1)
	}

	if !d.IsNewResource() && d.HasChange("common_setting") {
		update = true
	}
	objectDataLocalMap2 := make(map[string]interface{})
	if v := d.Get("common_setting"); !IsNil(v) {
		nodeNative81, _ := jsonpath.Get("$[0].ip_type", v)
		if nodeNative81 != "" {
			objectDataLocalMap2["IpType"] = nodeNative81
		}
		customHost_map, _ := jsonpath.Get("$[0].custom_host[0]", v)
		if !IsNil(customHost_map) {
			customHost := make(map[string]interface{})
			if v, ok := d.GetOk("common_setting"); ok {
				localData2, err := jsonpath.Get("$[0].custom_host[0].hosts", v)
				if err != nil {
					return WrapError(err)
				}
				localMaps := make([]map[string]interface{}, 0)
				for _, dataLoop2 := range localData2.([]interface{}) {
					dataLoop2Tmp := dataLoop2.(map[string]interface{})
					dataLoop2Map := make(map[string]interface{})
					dataLoop2Map["Domain"] = dataLoop2Tmp["domain"]
					dataLoop2Map["Ips"] = dataLoop2Tmp["ips"]
					dataLoop2Map["IpType"] = dataLoop2Tmp["ip_type"]
					localMaps = append(localMaps, dataLoop2Map)
				}
				customHost["Hosts"] = localMaps
			}
			nodeNative86, _ := jsonpath.Get("$[0].custom_host[0].select_type", v)
			if nodeNative86 != "" {
				customHost["SelectType"] = nodeNative86
			}
			objectDataLocalMap2["CustomHost"] = customHost
		}
		nodeNative88, _ := jsonpath.Get("$[0].monitor_samples", v)
		if nodeNative88 != "" {
			objectDataLocalMap2["MonitorSamples"] = nodeNative88
		}
		nodeNative89, _ := jsonpath.Get("$[0].is_open_trace", v)
		if nodeNative89 != "" {
			objectDataLocalMap2["IsOpenTrace"] = nodeNative89
		}
		nodeNative90, _ := jsonpath.Get("$[0].trace_client_type", v)
		if nodeNative90 != "" {
			objectDataLocalMap2["TraceClientType"] = nodeNative90
		}
		nodeNative91, _ := jsonpath.Get("$[0].xtrace_region", v)
		if nodeNative91 != "" {
			objectDataLocalMap2["XtraceRegion"] = nodeNative91
		}
		request["CommonSetting"] = convertMapToJsonStringIgnoreError(objectDataLocalMap2)
	}

	if _, ok := d.GetOk("resource_group_id"); ok && !d.IsNewResource() && d.HasChange("resource_group_id") {
		update = true
		request["ResourceGroupId"] = d.Get("resource_group_id")
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("ARMS", "2019-08-08", action, query, request, true)

			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("synthetic_task_name")
		d.SetPartial("frequency")
		d.SetPartial("resource_group_id")
	}
	update = false
	action = "ChangeResourceGroup"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["ResourceId"] = d.Id()
	request["RegionId"] = client.RegionId
	if _, ok := d.GetOk("resource_group_id"); ok && !d.IsNewResource() && d.HasChange("resource_group_id") {
		update = true
		request["NewResourceGroupId"] = d.Get("resource_group_id")
	}

	request["ResourceType"] = "SYNTHETICTASK"
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("ARMS", "2019-08-08", action, query, request, true)

			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("resource_group_id")
	}

	if d.HasChange("status") {
		client := meta.(*connectivity.AliyunClient)
		armsServiceV2 := ArmsServiceV2{client}
		object, err := armsServiceV2.DescribeArmsSyntheticTask(d.Id())
		if err != nil {
			return WrapError(err)
		}

		target := d.Get("status").(string)
		if object["Status"].(string) != target {
			if target == "RUNNING" {
				action = "StartTimingSyntheticTask"
				request = make(map[string]interface{})
				query = make(map[string]interface{})
				query["TaskIds"] = "[\"" + d.Id() + "\"]"
				request["RegionId"] = client.RegionId
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = client.RpcPost("ARMS", "2019-08-08", action, query, request, true)

					if err != nil {
						if NeedRetry(err) {
							wait()
							return resource.RetryableError(err)
						}
						return resource.NonRetryableError(err)
					}
					addDebug(action, response, request)
					return nil
				})
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
				}

			}
			if target == "STOP" {
				action = "StopTimingSyntheticTask"
				request = make(map[string]interface{})
				query = make(map[string]interface{})
				query["TaskIds"] = "[\"" + d.Id() + "\"]"
				request["RegionId"] = client.RegionId
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = client.RpcPost("ARMS", "2019-08-08", action, query, request, true)

					if err != nil {
						if NeedRetry(err) {
							wait()
							return resource.RetryableError(err)
						}
						return resource.NonRetryableError(err)
					}
					addDebug(action, response, request)
					return nil
				})
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
				}

			}
		}
	}

	if d.HasChange("tags") {
		armsServiceV2 := ArmsServiceV2{client}
		if err := armsServiceV2.SetResourceTags(d, "SYNTHETICTASK"); err != nil {
			return WrapError(err)
		}
		d.SetPartial("tags")
	}
	d.Partial(false)
	return resourceAliCloudArmsSyntheticTaskRead(d, meta)
}

func resourceAliCloudArmsSyntheticTaskDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteTimingSyntheticTask"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["TaskId"] = d.Id()
	request["RegionId"] = client.RegionId

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("ARMS", "2019-08-08", action, query, request, true)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
