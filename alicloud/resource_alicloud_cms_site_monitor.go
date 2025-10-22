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

func resourceAliCloudCloudMonitorServiceSiteMonitor() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudCloudMonitorServiceSiteMonitorCreate,
		Read:   resourceAliCloudCloudMonitorServiceSiteMonitorRead,
		Update: resourceAliCloudCloudMonitorServiceSiteMonitorUpdate,
		Delete: resourceAliCloudCloudMonitorServiceSiteMonitorDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"address": {
				Type:     schema.TypeString,
				Required: true,
			},
			"agent_group": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"custom_schedule": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"time_zone": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"start_hour": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"days": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeInt},
						},
						"end_hour": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
			"interval": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"1", "5", "15", "30", "60"}, false),
			},
			"isp_cities": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"isp": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"city": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"option_json": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"response_content": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"expect_value": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"is_base_encode": {
							Type:     schema.TypeBool,
							Optional: true,
							ForceNew: true,
							Computed: true,
						},
						"ping_num": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"match_rule": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: IntInSlice([]int{0, 1}),
						},
						"failure_rate": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"request_content": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"attempts": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"request_format": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: StringInSlice([]string{"hex", "text"}, false),
						},
						"password": {
							Type:      schema.TypeString,
							Optional:  true,
							Sensitive: true,
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								if old != "" && new != "" && old != new {
									return true
								}
								return false
							},
						},
						"diagnosis_ping": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"response_format": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: StringInSlice([]string{"hex", "text"}, false),
						},
						"cookie": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"ping_port": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"user_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"dns_match_rule": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: StringInSlice([]string{"IN_DNS", "DNS_IN", "EQUAL", "ANY"}, false),
						},
						"timeout": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"assertions": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"operator": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"target": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"type": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"dns_server": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"diagnosis_mtr": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"header": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"min_tls_version": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"ping_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"dns_type": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: StringInSlice([]string{"A", "CNAME", "NS", "MX", "TXT"}, false),
						},
						"dns_hijack_whitelist": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"http_method": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: StringInSlice([]string{"get", "post"}, false),
						},
					},
				},
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"task_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"task_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"options_json": {
				Type:       schema.TypeString,
				Optional:   true,
				Computed:   true,
				Deprecated: "Field `options_json` has been deprecated from provider version 1.262.0. New field `option_json` instead",
			},
			"task_state": {
				Type:       schema.TypeString,
				Computed:   true,
				Deprecated: "Field `task_state` has been deprecated from provider version 1.262.0. New field `status` instead.",
			},
			"create_time": {
				Type:       schema.TypeString,
				Computed:   true,
				Deprecated: "Field `create_time` has been deprecated from provider version 1.262.0.",
			},
			"update_time": {
				Type:       schema.TypeString,
				Computed:   true,
				Deprecated: "Field `update_time` has been deprecated from provider version 1.262.0.",
			},
			"alert_ids": {
				Type:       schema.TypeList,
				Optional:   true,
				Elem:       &schema.Schema{Type: schema.TypeString},
				Deprecated: "Field `alert_ids` has been deprecated from provider version 1.262.0.",
			},
		},
	}
}

func resourceAliCloudCloudMonitorServiceSiteMonitorCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateSiteMonitor"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	request["TaskType"] = d.Get("task_type")
	dataList := make(map[string]interface{})

	if v := d.Get("option_json"); !IsNil(v) {
		responseFormat, _ := jsonpath.Get("$[0].response_format", v)
		if responseFormat != nil && responseFormat != "" {
			dataList["response_format"] = responseFormat
		}
		dnsMatchRule, _ := jsonpath.Get("$[0].dns_match_rule", v)
		if dnsMatchRule != nil && dnsMatchRule != "" {
			dataList["dns_match_rule"] = dnsMatchRule
		}
		timeout, _ := jsonpath.Get("$[0].timeout", v)
		if timeout != nil && timeout != "" {
			dataList["time_out"] = timeout
		}
		if v, ok := d.GetOk("option_json"); ok {
			localData, err := jsonpath.Get("$[0].assertions", v)
			if err != nil {
				localData = make([]interface{}, 0)
			}
			localMaps := make([]interface{}, 0)
			for _, dataLoop := range convertToInterfaceArray(localData) {
				dataLoopTmp := make(map[string]interface{})
				if dataLoop != nil {
					dataLoopTmp = dataLoop.(map[string]interface{})
				}
				dataLoopMap := make(map[string]interface{})
				dataLoopMap["target"] = dataLoopTmp["target"]
				dataLoopMap["operator"] = dataLoopTmp["operator"]
				dataLoopMap["type"] = dataLoopTmp["type"]
				localMaps = append(localMaps, dataLoopMap)
			}
			dataList["assertions"] = localMaps
		}

		attempts1, _ := jsonpath.Get("$[0].attempts", v)
		if attempts1 != nil && attempts1 != "" {
			dataList["attempts"] = attempts1
		}
		pingType, _ := jsonpath.Get("$[0].ping_type", v)
		if pingType != nil && pingType != "" {
			dataList["ping_type"] = pingType
		}
		expectValue, _ := jsonpath.Get("$[0].expect_value", v)
		if expectValue != nil && expectValue != "" {
			dataList["expect_value"] = expectValue
		}
		httpMethod, _ := jsonpath.Get("$[0].http_method", v)
		if httpMethod != nil && httpMethod != "" {
			dataList["http_method"] = httpMethod
		}
		cookie1, _ := jsonpath.Get("$[0].cookie", v)
		if cookie1 != nil && cookie1 != "" {
			dataList["cookie"] = cookie1
		}
		responseContent, _ := jsonpath.Get("$[0].response_content", v)
		if responseContent != nil && responseContent != "" {
			dataList["response_content"] = responseContent
		}
		userName, _ := jsonpath.Get("$[0].user_name", v)
		if userName != nil && userName != "" {
			dataList["username"] = userName
		}
		pingNum, _ := jsonpath.Get("$[0].ping_num", v)
		if pingNum != nil && pingNum != "" {
			dataList["ping_num"] = pingNum
		}
		pingPort, _ := jsonpath.Get("$[0].ping_port", v)
		if pingPort != nil && pingPort != "" {
			dataList["ping_port"] = pingPort
		}
		header1, _ := jsonpath.Get("$[0].header", v)
		if header1 != nil && header1 != "" {
			dataList["header"] = header1
		}
		dnsHijackWhitelist, _ := jsonpath.Get("$[0].dns_hijack_whitelist", v)
		if dnsHijackWhitelist != nil && dnsHijackWhitelist != "" {
			dataList["dns_hijack_whitelist"] = dnsHijackWhitelist
		}
		diagnosisPing, _ := jsonpath.Get("$[0].diagnosis_ping", v)
		if diagnosisPing != nil && diagnosisPing != "" {
			dataList["diagnosis_ping"] = diagnosisPing
		}
		diagnosisMtr, _ := jsonpath.Get("$[0].diagnosis_mtr", v)
		if diagnosisMtr != nil && diagnosisMtr != "" {
			dataList["diagnosis_mtr"] = diagnosisMtr
		}
		dnsServer, _ := jsonpath.Get("$[0].dns_server", v)
		if dnsServer != nil && dnsServer != "" {
			dataList["dns_server"] = dnsServer
		}
		failureRate, _ := jsonpath.Get("$[0].failure_rate", v)
		if failureRate != nil && failureRate != "" {
			dataList["failure_rate"] = failureRate
		}
		password1, _ := jsonpath.Get("$[0].password", v)
		if password1 != nil && password1 != "" {
			dataList["password"] = password1
		}
		matchRule, _ := jsonpath.Get("$[0].match_rule", v)
		if matchRule != nil && matchRule != "" {
			dataList["match_rule"] = matchRule
		}
		requestContent, _ := jsonpath.Get("$[0].request_content", v)
		if requestContent != nil && requestContent != "" {
			dataList["request_content"] = requestContent
		}
		isBaseEncode, _ := jsonpath.Get("$[0].is_base_encode", v)
		if isBaseEncode != nil && isBaseEncode != "" {
			dataList["isBase64Encode"] = isBaseEncode
		}
		requestFormat, _ := jsonpath.Get("$[0].request_format", v)
		if requestFormat != nil && requestFormat != "" {
			dataList["request_format"] = requestFormat
		}
		dnsType, _ := jsonpath.Get("$[0].dns_type", v)
		if dnsType != nil && dnsType != "" {
			dataList["dns_type"] = dnsType
		}
		port1, _ := jsonpath.Get("$[0].port", v)
		if port1 != nil && port1 != "" {
			dataList["port"] = port1
		}
		minTlsVersion, _ := jsonpath.Get("$[0].min_tls_version", v)
		if minTlsVersion != nil && minTlsVersion != "" {
			dataList["min_tls_version"] = minTlsVersion
		}

		request["OptionsJson"] = convertMapToJsonStringIgnoreError(dataList)
	} else if v, ok := d.GetOk("options_json"); ok {
		request["OptionsJson"] = v
	}

	dataList1 := make(map[string]interface{})

	if v := d.Get("custom_schedule"); !IsNil(v) {
		days1, _ := jsonpath.Get("$[0].days", v)
		if days1 != nil && days1 != "" {
			dataList1["days"] = days1
		}
		startHour, _ := jsonpath.Get("$[0].start_hour", v)
		if startHour != nil && startHour != "" {
			dataList1["start_hour"] = startHour
		}
		endHour, _ := jsonpath.Get("$[0].end_hour", v)
		if endHour != nil && endHour != "" {
			dataList1["end_hour"] = endHour
		}
		timeZone, _ := jsonpath.Get("$[0].time_zone", v)
		if timeZone != nil && timeZone != "" {
			dataList1["time_zone"] = timeZone
		}

		customScheduleJson, err := convertMaptoJsonString(dataList1)
		if err != nil {
			return WrapError(err)
		}

		request["CustomSchedule"] = customScheduleJson
	}

	if v, ok := d.GetOk("interval"); ok {
		request["Interval"] = v
	}

	if v, ok := d.GetOk("isp_cities"); ok {
		ispCitiesMapsArray := make([]interface{}, 0)
		for _, dataLoop1 := range convertToInterfaceArray(v) {
			dataLoop1Tmp := dataLoop1.(map[string]interface{})
			dataLoop1Map := make(map[string]interface{})
			dataLoop1Map["Type"] = dataLoop1Tmp["type"]
			dataLoop1Map["Isp"] = dataLoop1Tmp["isp"]
			dataLoop1Map["City"] = dataLoop1Tmp["city"]
			ispCitiesMapsArray = append(ispCitiesMapsArray, dataLoop1Map)
		}

		ispCitiesJson, err := convertInterfaceToJsonString(ispCitiesMapsArray)
		if err != nil {
			return WrapError(err)
		}

		request["IspCities"] = ispCitiesJson
	}

	request["TaskName"] = d.Get("task_name")
	request["Address"] = d.Get("address")
	if v, ok := d.GetOk("agent_group"); ok {
		request["AgentGroup"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Cms", "2019-01-01", action, query, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cms_site_monitor", action, AlibabaCloudSdkGoERROR)
	}

	id, _ := jsonpath.Get("$.CreateResultList.CreateResultList[0].TaskId", response)
	d.SetId(fmt.Sprint(id))

	return resourceAliCloudCloudMonitorServiceSiteMonitorUpdate(d, meta)
}

func resourceAliCloudCloudMonitorServiceSiteMonitorRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudMonitorServiceServiceV2 := CloudMonitorServiceServiceV2{client}

	objectRaw, err := cloudMonitorServiceServiceV2.DescribeCloudMonitorServiceSiteMonitor(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cms_site_monitor DescribeCloudMonitorServiceSiteMonitor Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("address", objectRaw["Address"])
	d.Set("agent_group", objectRaw["AgentGroup"])
	d.Set("interval", objectRaw["Interval"])
	d.Set("status", fmt.Sprint(objectRaw["TaskState"]))
	d.Set("task_name", objectRaw["TaskName"])
	d.Set("task_type", objectRaw["TaskType"])
	d.Set("task_state", fmt.Sprint(objectRaw["TaskState"]))

	customScheduleMaps := make([]map[string]interface{}, 0)
	customScheduleMap := make(map[string]interface{})
	customScheduleRaw := make(map[string]interface{})
	if objectRaw["CustomSchedule"] != nil {
		customScheduleRaw = objectRaw["CustomSchedule"].(map[string]interface{})
	}
	if len(customScheduleRaw) > 0 {
		customScheduleMap["end_hour"] = customScheduleRaw["end_hour"]
		customScheduleMap["start_hour"] = customScheduleRaw["start_hour"]
		customScheduleMap["time_zone"] = customScheduleRaw["time_zone"]

		daysRaw, _ := jsonpath.Get("$.CustomSchedule.days.days", objectRaw)
		customScheduleMap["days"] = daysRaw
		customScheduleMaps = append(customScheduleMaps, customScheduleMap)
	}
	if err := d.Set("custom_schedule", customScheduleMaps); err != nil {
		return err
	}

	ispCityRaw, _ := jsonpath.Get("$.IspCities.IspCity", objectRaw)
	ispCitiesMaps := make([]map[string]interface{}, 0)
	if ispCityRaw != nil {
		for _, ispCityChildRaw := range convertToInterfaceArray(ispCityRaw) {
			ispCitiesMap := make(map[string]interface{})
			ispCityChildRaw := ispCityChildRaw.(map[string]interface{})
			ispCitiesMap["city"] = ispCityChildRaw["City"]
			ispCitiesMap["isp"] = ispCityChildRaw["Isp"]
			ispCitiesMap["type"] = ispCityChildRaw["Type"]

			ispCitiesMaps = append(ispCitiesMaps, ispCitiesMap)
		}
	}
	if err := d.Set("isp_cities", ispCitiesMaps); err != nil {
		return err
	}

	optionJsonMaps := make([]map[string]interface{}, 0)
	optionJsonMap := make(map[string]interface{})
	optionJsonRaw := make(map[string]interface{})
	if objectRaw["OptionJson"] != nil {
		optionJsonRaw = objectRaw["OptionJson"].(map[string]interface{})

		optionJsonJson, err := convertInterfaceToJsonString(objectRaw["OptionJson"])
		if err != nil {
			return WrapError(err)
		}

		d.Set("options_json", optionJsonJson)
	}
	if len(optionJsonRaw) > 0 {
		optionJsonMap["attempts"] = optionJsonRaw["attempts"]
		optionJsonMap["cookie"] = optionJsonRaw["cookie"]
		optionJsonMap["diagnosis_mtr"] = optionJsonRaw["diagnosis_mtr"]
		optionJsonMap["diagnosis_ping"] = optionJsonRaw["diagnosis_ping"]
		optionJsonMap["dns_hijack_whitelist"] = optionJsonRaw["dns_hijack_whitelist"]
		optionJsonMap["dns_match_rule"] = optionJsonRaw["dns_match_rule"]
		optionJsonMap["dns_server"] = optionJsonRaw["dns_server"]
		optionJsonMap["dns_type"] = optionJsonRaw["dns_type"]
		optionJsonMap["expect_value"] = optionJsonRaw["expect_value"]
		optionJsonMap["failure_rate"] = optionJsonRaw["failure_rate"]
		optionJsonMap["header"] = optionJsonRaw["header"]
		optionJsonMap["http_method"] = optionJsonRaw["http_method"]
		optionJsonMap["is_base_encode"] = optionJsonRaw["isBase64Encode"]
		optionJsonMap["match_rule"] = optionJsonRaw["match_rule"]
		optionJsonMap["min_tls_version"] = optionJsonRaw["min_tls_version"]
		optionJsonMap["password"] = optionJsonRaw["password"]
		optionJsonMap["ping_num"] = optionJsonRaw["ping_num"]
		optionJsonMap["ping_port"] = optionJsonRaw["ping_port"]
		optionJsonMap["ping_type"] = optionJsonRaw["ping_type"]
		optionJsonMap["port"] = optionJsonRaw["port"]
		optionJsonMap["request_content"] = optionJsonRaw["request_content"]
		optionJsonMap["request_format"] = optionJsonRaw["request_format"]
		optionJsonMap["response_content"] = optionJsonRaw["response_content"]
		optionJsonMap["response_format"] = optionJsonRaw["response_format"]
		optionJsonMap["timeout"] = optionJsonRaw["time_out"]
		optionJsonMap["user_name"] = optionJsonRaw["username"]

		assertionsRaw, _ := jsonpath.Get("$.OptionJson.assertions.assertions", objectRaw)
		assertionsMaps := make([]map[string]interface{}, 0)
		if assertionsRaw != nil {
			for _, assertionsChildRaw := range convertToInterfaceArray(assertionsRaw) {
				assertionsMap := make(map[string]interface{})
				assertionsChildRaw := assertionsChildRaw.(map[string]interface{})
				assertionsMap["operator"] = assertionsChildRaw["operator"]
				assertionsMap["target"] = assertionsChildRaw["target"]
				assertionsMap["type"] = assertionsChildRaw["type"]

				assertionsMaps = append(assertionsMaps, assertionsMap)
			}
		}
		optionJsonMap["assertions"] = assertionsMaps
		optionJsonMaps = append(optionJsonMaps, optionJsonMap)
	}
	if err := d.Set("option_json", optionJsonMaps); err != nil {
		return err
	}

	return nil
}

func resourceAliCloudCloudMonitorServiceSiteMonitorUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	cloudMonitorServiceServiceV2 := CloudMonitorServiceServiceV2{client}
	objectRaw, _ := cloudMonitorServiceServiceV2.DescribeCloudMonitorServiceSiteMonitor(d.Id())

	if d.HasChange("status") {
		var err error
		target := d.Get("status").(string)
		if fmt.Sprint(objectRaw["TaskState"]) != target {
			if target == "1" {
				action := "EnableSiteMonitors"
				request = make(map[string]interface{})
				query = make(map[string]interface{})
				request["TaskIds"] = d.Id()

				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = client.RpcPost("Cms", "2019-01-01", action, query, request, true)
					if err != nil {
						if NeedRetry(err) {
							wait()
							return resource.RetryableError(err)
						}
						return resource.NonRetryableError(err)
					}
					return nil
				})
				addDebug(action, response, request)
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
				}

			}
			if target == "2" {
				action := "DisableSiteMonitors"
				request = make(map[string]interface{})
				query = make(map[string]interface{})
				request["TaskIds"] = d.Id()

				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = client.RpcPost("Cms", "2019-01-01", action, query, request, true)
					if err != nil {
						if NeedRetry(err) {
							wait()
							return resource.RetryableError(err)
						}
						return resource.NonRetryableError(err)
					}
					return nil
				})
				addDebug(action, response, request)
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
				}

			}
		}
	}

	var err error
	action := "ModifySiteMonitor"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["TaskId"] = d.Id()

	if !d.IsNewResource() && d.HasChange("option_json") {
		update = true
		dataList := make(map[string]interface{})

		if v := d.Get("option_json"); v != nil {
			responseFormat, _ := jsonpath.Get("$[0].response_format", v)
			if responseFormat != nil && (d.HasChange("option_json.0.response_format") || responseFormat != "") {
				dataList["response_format"] = responseFormat
			}
			dnsMatchRule, _ := jsonpath.Get("$[0].dns_match_rule", v)
			if dnsMatchRule != nil && (d.HasChange("option_json.0.dns_match_rule") || dnsMatchRule != "") {
				dataList["dns_match_rule"] = dnsMatchRule
			}
			if v, ok := d.GetOk("option_json"); ok {
				localData, err := jsonpath.Get("$[0].assertions", v)
				if err != nil {
					localData = make([]interface{}, 0)
				}
				localMaps := make([]interface{}, 0)
				for _, dataLoop := range convertToInterfaceArray(localData) {
					dataLoopTmp := make(map[string]interface{})
					if dataLoop != nil {
						dataLoopTmp = dataLoop.(map[string]interface{})
					}
					dataLoopMap := make(map[string]interface{})
					dataLoopMap["target"] = dataLoopTmp["target"]
					dataLoopMap["operator"] = dataLoopTmp["operator"]
					dataLoopMap["type"] = dataLoopTmp["type"]
					localMaps = append(localMaps, dataLoopMap)
				}
				dataList["assertions"] = localMaps
			}

			attempts1, _ := jsonpath.Get("$[0].attempts", v)
			if attempts1 != nil && (d.HasChange("option_json.0.attempts") || attempts1 != "") {
				dataList["attempts"] = attempts1
			}
			pingType, _ := jsonpath.Get("$[0].ping_type", v)
			if pingType != nil && (d.HasChange("option_json.0.ping_type") || pingType != "") {
				dataList["ping_type"] = pingType
			}
			expectValue, _ := jsonpath.Get("$[0].expect_value", v)
			if expectValue != nil && (d.HasChange("option_json.0.expect_value") || expectValue != "") {
				dataList["expect_value"] = expectValue
			}
			httpMethod, _ := jsonpath.Get("$[0].http_method", v)
			if httpMethod != nil && (d.HasChange("option_json.0.http_method") || httpMethod != "") {
				dataList["http_method"] = httpMethod
			}
			cookie1, _ := jsonpath.Get("$[0].cookie", v)
			if cookie1 != nil && (d.HasChange("option_json.0.cookie") || cookie1 != "") {
				dataList["cookie"] = cookie1
			}
			responseContent, _ := jsonpath.Get("$[0].response_content", v)
			if responseContent != nil && (d.HasChange("option_json.0.response_content") || responseContent != "") {
				dataList["response_content"] = responseContent
			}
			userName, _ := jsonpath.Get("$[0].user_name", v)
			if userName != nil && (d.HasChange("option_json.0.user_name") || userName != "") {
				dataList["username"] = userName
			}
			pingNum, _ := jsonpath.Get("$[0].ping_num", v)
			if pingNum != nil && (d.HasChange("option_json.0.ping_num") || pingNum != "") {
				dataList["ping_num"] = pingNum
			}
			pingPort, _ := jsonpath.Get("$[0].ping_port", v)
			if pingPort != nil && (d.HasChange("option_json.0.ping_port") || pingPort != "") {
				dataList["ping_port"] = pingPort
			}
			header1, _ := jsonpath.Get("$[0].header", v)
			if header1 != nil && (d.HasChange("option_json.0.header") || header1 != "") {
				dataList["header"] = header1
			}
			dnsHijackWhitelist, _ := jsonpath.Get("$[0].dns_hijack_whitelist", v)
			if dnsHijackWhitelist != nil && (d.HasChange("option_json.0.dns_hijack_whitelist") || dnsHijackWhitelist != "") {
				dataList["dns_hijack_whitelist"] = dnsHijackWhitelist
			}
			diagnosisPing, _ := jsonpath.Get("$[0].diagnosis_ping", v)
			if diagnosisPing != nil && (d.HasChange("option_json.0.diagnosis_ping") || diagnosisPing != "") {
				dataList["diagnosis_ping"] = diagnosisPing
			}
			diagnosisMtr, _ := jsonpath.Get("$[0].diagnosis_mtr", v)
			if diagnosisMtr != nil && (d.HasChange("option_json.0.diagnosis_mtr") || diagnosisMtr != "") {
				dataList["diagnosis_mtr"] = diagnosisMtr
			}
			dnsServer, _ := jsonpath.Get("$[0].dns_server", v)
			if dnsServer != nil && (d.HasChange("option_json.0.dns_server") || dnsServer != "") {
				dataList["dns_server"] = dnsServer
			}
			failureRate, _ := jsonpath.Get("$[0].failure_rate", v)
			if failureRate != nil && (d.HasChange("option_json.0.failure_rate") || failureRate != "") {
				dataList["failure_rate"] = failureRate
			}
			password1, _ := jsonpath.Get("$[0].password", v)
			if password1 != nil && (d.HasChange("option_json.0.password") || password1 != "") {
				dataList["password"] = password1
			}
			matchRule, _ := jsonpath.Get("$[0].match_rule", v)
			if matchRule != nil && (d.HasChange("option_json.0.match_rule") || matchRule != "") {
				dataList["match_rule"] = matchRule
			}
			requestContent, _ := jsonpath.Get("$[0].request_content", v)
			if requestContent != nil && (d.HasChange("option_json.0.request_content") || requestContent != "") {
				dataList["request_content"] = requestContent
			}
			requestFormat, _ := jsonpath.Get("$[0].request_format", v)
			if requestFormat != nil && (d.HasChange("option_json.0.request_format") || requestFormat != "") {
				dataList["request_format"] = requestFormat
			}
			dnsType, _ := jsonpath.Get("$[0].dns_type", v)
			if dnsType != nil && (d.HasChange("option_json.0.dns_type") || dnsType != "") {
				dataList["dns_type"] = dnsType
			}
			port1, _ := jsonpath.Get("$[0].port", v)
			if port1 != nil && (d.HasChange("option_json.0.port") || port1 != "") {
				dataList["port"] = port1
			}
			minTlsVersion, _ := jsonpath.Get("$[0].min_tls_version", v)
			if minTlsVersion != nil && (d.HasChange("option_json.0.min_tls_version") || minTlsVersion != "") {
				dataList["min_tls_version"] = minTlsVersion
			}

			request["OptionsJson"] = convertMapToJsonStringIgnoreError(dataList)
		}
	}

	if !d.IsNewResource() && d.HasChange("options_json") {
		update = true

		if v, ok := d.GetOk("options_json"); ok {
			request["OptionsJson"] = v
		}
	}

	if !d.IsNewResource() && d.HasChange("custom_schedule") {
		update = true
		dataList1 := make(map[string]interface{})

		if v := d.Get("custom_schedule"); v != nil {
			days1, _ := jsonpath.Get("$[0].days", v)
			if days1 != nil && (d.HasChange("custom_schedule.0.days") || days1 != "") {
				dataList1["days"] = days1
			}
			startHour, _ := jsonpath.Get("$[0].start_hour", v)
			if startHour != nil && (d.HasChange("custom_schedule.0.start_hour") || startHour != "") {
				dataList1["start_hour"] = startHour
			}
			endHour, _ := jsonpath.Get("$[0].end_hour", v)
			if endHour != nil && (d.HasChange("custom_schedule.0.end_hour") || endHour != "") {
				dataList1["end_hour"] = endHour
			}
			timeZone, _ := jsonpath.Get("$[0].time_zone", v)
			if timeZone != nil && (d.HasChange("custom_schedule.0.time_zone") || timeZone != "") {
				dataList1["time_zone"] = timeZone
			}

			customScheduleJson, err := convertMaptoJsonString(dataList1)
			if err != nil {
				return WrapError(err)
			}

			request["CustomSchedule"] = customScheduleJson
		}
	}

	if !d.IsNewResource() && d.HasChange("interval") {
		update = true
		request["Interval"] = d.Get("interval")
	}

	if !d.IsNewResource() && d.HasChange("isp_cities") {
		update = true
		if v, ok := d.GetOk("isp_cities"); ok || d.HasChange("isp_cities") {
			ispCitiesMapsArray := make([]interface{}, 0)
			for _, dataLoop1 := range convertToInterfaceArray(v) {
				dataLoop1Tmp := dataLoop1.(map[string]interface{})
				dataLoop1Map := make(map[string]interface{})
				dataLoop1Map["Type"] = dataLoop1Tmp["type"]
				dataLoop1Map["isp"] = dataLoop1Tmp["isp"]
				dataLoop1Map["city"] = dataLoop1Tmp["city"]
				ispCitiesMapsArray = append(ispCitiesMapsArray, dataLoop1Map)
			}

			ispCitiesJson, err := convertInterfaceToJsonString(ispCitiesMapsArray)
			if err != nil {
				return WrapError(err)
			}

			request["IspCities"] = ispCitiesJson
		}
	}

	if !d.IsNewResource() && d.HasChange("task_name") {
		update = true
	}
	request["TaskName"] = d.Get("task_name")

	if !d.IsNewResource() && d.HasChange("address") {
		update = true
	}
	request["Address"] = d.Get("address")
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Cms", "2019-01-01", action, query, request, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}

	return resourceAliCloudCloudMonitorServiceSiteMonitorRead(d, meta)
}

func resourceAliCloudCloudMonitorServiceSiteMonitorDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteSiteMonitors"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["TaskIds"] = d.Id()

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Cms", "2019-01-01", action, query, request, true)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		if IsExpectedErrors(err, []string{"ResourceNotFound"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
