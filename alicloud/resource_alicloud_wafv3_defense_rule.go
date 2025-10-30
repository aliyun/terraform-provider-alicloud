// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudWafv3DefenseRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudWafv3DefenseRuleCreate,
		Read:   resourceAliCloudWafv3DefenseRuleRead,
		Update: resourceAliCloudWafv3DefenseRuleUpdate,
		Delete: resourceAliCloudWafv3DefenseRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"config": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bypass_tags": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"gray_config": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"gray_target": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: StringInSlice([]string{"cookie.acw_tc", "header", "queryarg", "cookie", "remote_addr"}, false),
									},
									"gray_rate": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: IntBetween(0, 100),
									},
									"gray_sub_key": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"rule_action": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: StringInSlice([]string{"block", "monitor", "js", "captcha", "captcha_strict", "filter"}, false),
						},
						"mode": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: IntInSlice([]int{0, 1}),
						},
						"bypass_regular_rules": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"cn_regions": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"time_config": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"time_zone": {
										Type:         schema.TypeInt,
										Optional:     true,
										Computed:     true,
										ValidateFunc: IntBetween(-12, 12),
									},
									"time_periods": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"start": {
													Type:     schema.TypeInt,
													Optional: true,
												},
												"end": {
													Type:     schema.TypeInt,
													Optional: true,
												},
											},
										},
									},
									"time_scope": {
										Type:         schema.TypeString,
										Optional:     true,
										Computed:     true,
										ValidateFunc: StringInSlice([]string{"permanent", "period", "cycle"}, false),
									},
									"week_time_periods": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"day_periods": {
													Type:     schema.TypeSet,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"start": {
																Type:     schema.TypeInt,
																Optional: true,
															},
															"end": {
																Type:     schema.TypeInt,
																Optional: true,
															},
														},
													},
												},
												"day": {
													Type:     schema.TypeString,
													Optional: true,
												},
											},
										},
									},
								},
							},
						},
						"bypass_regular_types": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"ua": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"remote_addr": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"conditions": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"op_value": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: StringInSlice([]string{"not-contain", "contain", "none", "ne", "eq", "lt", "gt", "len-lt", "len-eq", "len-gt", "not-match", "match-one", "all-not-match", "all-not-contain", "contain-one", "not-regex", "regex", "all-not-regex", "regex-one", "prefix-match", "suffix-match", "empty", "exists", "inl"}, false),
									},
									"values": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"key": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"sub_key": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"url": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"abroad_regions": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"rate_limit": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"status": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"ratio": {
													Type:         schema.TypeInt,
													Optional:     true,
													ValidateFunc: IntBetween(0, 100),
												},
												"count": {
													Type:         schema.TypeInt,
													Optional:     true,
													ValidateFunc: IntBetween(2, 50000),
												},
												"code": {
													Type:     schema.TypeInt,
													Optional: true,
												},
											},
										},
									},
									"target": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: StringInSlice([]string{"cookie.acw_tc", "header", "queryarg", "cookie", "remote_addr"}, false),
									},
									"ttl": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: IntBetween(60, 86400),
									},
									"interval": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: IntBetween(0, 1800),
									},
									"threshold": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"sub_key": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"cc_effect": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: StringInSlice([]string{"service", "rule"}, false),
						},
						"throttle_threhold": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: IntBetween(0, 5000000),
						},
						"throttle_type": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: StringInSlice([]string{"qps", "ratio"}, false),
						},
						"cc_status": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: IntInSlice([]int{0, 1}),
						},
						"account_identifiers": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"position": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"priority": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: IntBetween(0, 20),
									},
									"decode_type": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: StringInSlice([]string{"plain", "basic", "jwt"}, false),
									},
									"key": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: StringInSlice([]string{"Query-Arg", "Cookie-Exact", "Post-Arg", "Header"}, false),
									},
									"sub_key": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"gray_status": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: IntInSlice([]int{0, 1}),
						},
						"protocol": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: StringInSlice([]string{"https", "http"}, false),
						},
					},
				},
			},
			"defense_origin": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"custom"}, false),
			},
			"defense_scene": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"ip_blacklist", "custom_acl", "whitelist", "region_block", "cc", "tamperproof", "dlp", "spike_throttle", "account_identifier"}, false),
			},
			"defense_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"template", "resource", "global"}, false),
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"resource": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"rule_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"rule_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"rule_status": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: IntInSlice([]int{0, 1}),
			},
			"template_id": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudWafv3DefenseRuleCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateDefenseRule"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		request["InstanceId"] = v
	}
	if v, ok := d.GetOk("defense_type"); ok {
		request["DefenseType"] = v
	}
	request["RegionId"] = client.RegionId

	dataList := make(map[string]interface{})

	if v, ok := d.GetOk("config"); ok {
		mode1, _ := jsonpath.Get("$[0].mode", v)
		if mode1 != nil && mode1 != "" {
			dataList["mode"] = mode1
		}
	}

	if v := d.Get("config"); !IsNil(v) {
		ratelimit := make(map[string]interface{})
		interval1, _ := jsonpath.Get("$[0].rate_limit[0].interval", d.Get("config"))
		if interval1 != nil && interval1 != "" && interval1.(int) > 0 {
			ratelimit["interval"] = interval1
		}
		status := make(map[string]interface{})
		count1, _ := jsonpath.Get("$[0].rate_limit[0].status[0].count", d.Get("config"))
		if count1 != nil && count1 != "" && count1.(int) > 0 {
			status["count"] = count1
		}
		code1, _ := jsonpath.Get("$[0].rate_limit[0].status[0].code", d.Get("config"))
		if code1 != nil && code1 != "" {
			status["code"] = code1
		}
		ratio1, _ := jsonpath.Get("$[0].rate_limit[0].status[0].ratio", d.Get("config"))
		if ratio1 != nil && ratio1 != "" && ratio1.(int) > 0 {
			status["ratio"] = ratio1
		}

		if len(status) > 0 {
			ratelimit["status"] = status
		}
		target1, _ := jsonpath.Get("$[0].rate_limit[0].target", d.Get("config"))
		if target1 != nil && target1 != "" {
			ratelimit["target"] = target1
		}
		threshold1, _ := jsonpath.Get("$[0].rate_limit[0].threshold", d.Get("config"))
		if threshold1 != nil && threshold1 != "" {
			ratelimit["threshold"] = threshold1
		}
		ttl1, _ := jsonpath.Get("$[0].rate_limit[0].ttl", d.Get("config"))
		if ttl1 != nil && ttl1 != "" && ttl1.(int) > 0 {
			ratelimit["ttl"] = ttl1
		}
		subKey1, _ := jsonpath.Get("$[0].rate_limit[0].sub_key", d.Get("config"))
		if subKey1 != nil && subKey1 != "" {
			ratelimit["subKey"] = subKey1
		}

		if len(ratelimit) > 0 {
			dataList["ratelimit"] = ratelimit
		}
	}

	if v, ok := d.GetOk("config"); ok {
		grayStatus1, _ := jsonpath.Get("$[0].gray_status", v)
		if grayStatus1 != nil && grayStatus1 != "" {
			dataList["grayStatus"] = grayStatus1
		}
	}

	if v := d.Get("config"); !IsNil(v) {
		timeConfig := make(map[string]interface{})
		if v, ok := d.GetOk("config"); ok {
			localData, err := jsonpath.Get("$[0].time_config[0].time_periods", v)
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
				dataLoopMap["start"] = dataLoopTmp["start"]
				dataLoopMap["end"] = dataLoopTmp["end"]
				localMaps = append(localMaps, dataLoopMap)
			}
			timeConfig["timePeriods"] = localMaps
		}

		timeZone1, _ := jsonpath.Get("$[0].time_config[0].time_zone", d.Get("config"))
		if timeZone1 != nil && timeZone1 != "" {
			timeConfig["timeZone"] = timeZone1
		}
		if v, ok := d.GetOk("config"); ok {
			localData1, err := jsonpath.Get("$[0].time_config[0].week_time_periods", v)
			if err != nil {
				localData1 = make([]interface{}, 0)
			}
			localMaps1 := make([]interface{}, 0)
			for _, dataLoop1 := range convertToInterfaceArray(localData1) {
				dataLoop1Tmp := make(map[string]interface{})
				if dataLoop1 != nil {
					dataLoop1Tmp = dataLoop1.(map[string]interface{})
				}
				dataLoop1Map := make(map[string]interface{})
				localMaps2 := make([]interface{}, 0)
				localData2 := dataLoop1Tmp["day_periods"]
				for _, dataLoop2 := range convertToInterfaceArray(localData2) {
					dataLoop2Tmp := dataLoop2.(map[string]interface{})
					dataLoop2Map := make(map[string]interface{})
					dataLoop2Map["start"] = dataLoop2Tmp["start"]
					dataLoop2Map["end"] = dataLoop2Tmp["end"]
					localMaps2 = append(localMaps2, dataLoop2Map)
				}
				dataLoop1Map["dayPeriods"] = localMaps2
				dataLoop1Map["day"] = dataLoop1Tmp["day"]
				localMaps1 = append(localMaps1, dataLoop1Map)
			}
			timeConfig["weekTimePeriods"] = localMaps1
		}

		timeScope1, _ := jsonpath.Get("$[0].time_config[0].time_scope", d.Get("config"))
		if timeScope1 != nil && timeScope1 != "" {
			timeConfig["timeScope"] = timeScope1
		}

		if len(timeConfig) > 0 {
			dataList["timeConfig"] = timeConfig
		}
	}

	if v := d.Get("config"); !IsNil(v) {
		grayConfig := make(map[string]interface{})
		grayTarget1, _ := jsonpath.Get("$[0].gray_config[0].gray_target", d.Get("config"))
		if grayTarget1 != nil && grayTarget1 != "" {
			grayConfig["grayTarget"] = grayTarget1
		}
		graySubKey1, _ := jsonpath.Get("$[0].gray_config[0].gray_sub_key", d.Get("config"))
		if graySubKey1 != nil && graySubKey1 != "" {
			grayConfig["graySubKey"] = graySubKey1
		}
		grayRate1, _ := jsonpath.Get("$[0].gray_config[0].gray_rate", d.Get("config"))
		if grayRate1 != nil && grayRate1 != "" && grayRate1.(int) > 0 {
			grayConfig["grayRate"] = grayRate1
		}

		if len(grayConfig) > 0 {
			dataList["grayConfig"] = grayConfig
		}
	}

	if v, ok := d.GetOk("config"); ok {
		ccEffect, _ := jsonpath.Get("$[0].cc_effect", v)
		if ccEffect != nil && ccEffect != "" {
			dataList["effect"] = ccEffect
		}
	}

	if v, ok := d.GetOk("config"); ok {
		throttleType, _ := jsonpath.Get("$[0].throttle_type", v)
		if throttleType != nil && throttleType != "" {
			dataList["type"] = throttleType
		}
	}

	if v, ok := d.GetOk("config"); ok {
		throttleThrehold, _ := jsonpath.Get("$[0].throttle_threhold", v)
		if throttleThrehold != nil && throttleThrehold != "" && throttleThrehold.(int) > 0 {
			dataList["threshold"] = throttleThrehold
		}
	}

	if v, ok := d.GetOk("rule_name"); ok {
		dataList["name"] = v
	}

	if v, ok := d.GetOk("config"); ok {
		url1, _ := jsonpath.Get("$[0].url", v)
		if url1 != nil && url1 != "" {
			dataList["url"] = url1
		}
	}

	if v, ok := d.GetOk("config"); ok {
		protocol1, _ := jsonpath.Get("$[0].protocol", v)
		if protocol1 != nil && protocol1 != "" {
			dataList["protocol"] = protocol1
		}
	}

	if v, ok := d.GetOk("config"); ok {
		ua1, _ := jsonpath.Get("$[0].ua", v)
		if ua1 != nil && ua1 != "" {
			dataList["ua"] = ua1
		}
	}

	if v := d.Get("config"); !IsNil(v) {
		if v, ok := d.GetOk("config"); ok {
			localData3, err := jsonpath.Get("$[0].account_identifiers", v)
			if err != nil {
				localData3 = make([]interface{}, 0)
			}
			localMaps3 := make([]interface{}, 0)
			for _, dataLoop3 := range convertToInterfaceArray(localData3) {
				dataLoop3Tmp := make(map[string]interface{})
				if dataLoop3 != nil {
					dataLoop3Tmp = dataLoop3.(map[string]interface{})
				}
				dataLoop3Map := make(map[string]interface{})
				dataLoop3Map["decodeType"] = dataLoop3Tmp["decode_type"]
				dataLoop3Map["position"] = dataLoop3Tmp["position"]
				dataLoop3Map["priority"] = dataLoop3Tmp["priority"]
				dataLoop3Map["subKey"] = dataLoop3Tmp["sub_key"]
				dataLoop3Map["key"] = dataLoop3Tmp["key"]
				localMaps3 = append(localMaps3, dataLoop3Map)
			}
			dataList["accountIdentifiers"] = localMaps3
		}

	}

	if v, ok := d.GetOk("config"); ok {
		ruleAction, _ := jsonpath.Get("$[0].rule_action", v)
		if ruleAction != nil && ruleAction != "" {
			dataList["action"] = ruleAction
		}
	}

	if v, ok := d.GetOk("defense_origin"); ok {
		dataList["origin"] = v
	}

	if v, ok := d.GetOk("config"); ok {
		ccStatus1, _ := jsonpath.Get("$[0].cc_status", v)
		if ccStatus1 != nil && ccStatus1 != "" {
			dataList["ccStatus"] = ccStatus1
		}
	}

	if v := d.Get("config"); !IsNil(v) {
		if v, ok := d.GetOk("config"); ok {
			localData4, err := jsonpath.Get("$[0].conditions", v)
			if err != nil {
				localData4 = make([]interface{}, 0)
			}
			localMaps4 := make([]interface{}, 0)
			for _, dataLoop4 := range convertToInterfaceArray(localData4) {
				dataLoop4Tmp := make(map[string]interface{})
				if dataLoop4 != nil {
					dataLoop4Tmp = dataLoop4.(map[string]interface{})
				}
				dataLoop4Map := make(map[string]interface{})
				dataLoop4Map["values"] = dataLoop4Tmp["values"]
				dataLoop4Map["opValue"] = dataLoop4Tmp["op_value"]
				dataLoop4Map["subKey"] = dataLoop4Tmp["sub_key"]
				dataLoop4Map["key"] = dataLoop4Tmp["key"]
				localMaps4 = append(localMaps4, dataLoop4Map)
			}
			dataList["conditions"] = localMaps4
		}

	}

	if v, ok := d.GetOkExists("rule_status"); ok {
		dataList["status"] = v
	}

	if v, ok := d.GetOk("config"); ok {
		bypassRegularTypes, _ := jsonpath.Get("$[0].bypass_regular_types", v)
		if bypassRegularTypes != nil && bypassRegularTypes != "" {
			dataList["regularTypes"] = convertToInterfaceArray(bypassRegularTypes)
		}
	}

	if v, ok := d.GetOk("config"); ok {
		bypassTags, _ := jsonpath.Get("$[0].bypass_tags", v)
		if bypassTags != nil && bypassTags != "" {
			dataList["tags"] = convertToInterfaceArray(bypassTags)
		}
	}

	if v, ok := d.GetOk("config"); ok {
		abroadRegions, _ := jsonpath.Get("$[0].abroad_regions", v)
		if abroadRegions != nil && abroadRegions != "" {
			dataList["abroadRegionList"] = abroadRegions
		}
	}

	if v, ok := d.GetOk("config"); ok {
		remoteAddr1, _ := jsonpath.Get("$[0].remote_addr", v)
		if remoteAddr1 != nil && remoteAddr1 != "" {
			dataList["remoteAddr"] = convertToInterfaceArray(remoteAddr1)
		}
	}

	if v, ok := d.GetOk("config"); ok {
		bypassRegularRules, _ := jsonpath.Get("$[0].bypass_regular_rules", v)
		if bypassRegularRules != nil && bypassRegularRules != "" {
			dataList["regularRules"] = convertToInterfaceArray(bypassRegularRules)
		}
	}

	if v, ok := d.GetOk("config"); ok {
		cnRegions, _ := jsonpath.Get("$[0].cn_regions", v)
		if cnRegions != nil && cnRegions != "" {
			dataList["cnRegionList"] = cnRegions
		}
	}

	RulesMap := make([]map[string]interface{}, 0)
	RulesMap = append(RulesMap, dataList)
	request["Rules"], _ = convertListMapToJsonString(RulesMap)

	if v, ok := d.GetOk("resource"); ok {
		request["Resource"] = v
	}
	if v, ok := d.GetOkExists("template_id"); ok {
		request["TemplateId"] = v
	}
	request["DefenseScene"] = d.Get("defense_scene")
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("waf-openapi", "2021-10-01", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_wafv3_defense_rule", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v:%v", request["InstanceId"], request["DefenseType"], response["RuleIds"]))

	return resourceAliCloudWafv3DefenseRuleRead(d, meta)
}

func resourceAliCloudWafv3DefenseRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	wafv3ServiceV2 := Wafv3ServiceV2{client}

	objectRaw, err := wafv3ServiceV2.DescribeWafv3DefenseRule(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_wafv3_defense_rule DescribeWafv3DefenseRule Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("defense_origin", objectRaw["DefenseOrigin"])
	d.Set("defense_scene", objectRaw["DefenseScene"])
	d.Set("resource", objectRaw["Resource"])
	d.Set("rule_name", objectRaw["RuleName"])
	d.Set("rule_status", objectRaw["Status"])
	d.Set("template_id", objectRaw["TemplateId"])
	d.Set("defense_type", objectRaw["DefenseType"])
	d.Set("rule_id", objectRaw["RuleId"])

	configMaps := make([]map[string]interface{}, 0)
	configMap := make(map[string]interface{})
	configRaw := make(map[string]interface{})
	if objectRaw["Config"] != nil {
		configRaw = convertJsonStringToObject(objectRaw["Config"])
	}
	if len(configRaw) > 0 {
		configMap["abroad_regions"] = configRaw["abroadRegionList"]
		configMap["cc_effect"] = configRaw["effect"]
		configMap["cc_status"] = formatInt(configRaw["ccStatus"])
		configMap["cn_regions"] = configRaw["cnRegionList"]
		configMap["gray_status"] = formatInt(configRaw["grayStatus"])
		configMap["mode"] = formatInt(configRaw["mode"])
		configMap["protocol"] = configRaw["protocol"]
		configMap["rule_action"] = configRaw["action"]
		configMap["throttle_threhold"] = formatInt(configRaw["threshold"])
		configMap["throttle_type"] = configRaw["type"]
		configMap["ua"] = configRaw["ua"]
		configMap["url"] = configRaw["url"]

		accountIdentifiersRaw := configRaw["accountIdentifiers"]
		accountIdentifiersMaps := make([]map[string]interface{}, 0)
		if accountIdentifiersRaw != nil {
			for _, accountIdentifiersChildRaw := range convertToInterfaceArray(accountIdentifiersRaw) {
				accountIdentifiersMap := make(map[string]interface{})
				accountIdentifiersChildRaw := accountIdentifiersChildRaw.(map[string]interface{})
				accountIdentifiersMap["decode_type"] = accountIdentifiersChildRaw["decodeType"]
				accountIdentifiersMap["key"] = accountIdentifiersChildRaw["key"]
				accountIdentifiersMap["position"] = accountIdentifiersChildRaw["position"]
				accountIdentifiersMap["priority"] = formatInt(accountIdentifiersChildRaw["priority"])
				accountIdentifiersMap["sub_key"] = accountIdentifiersChildRaw["subKey"]

				accountIdentifiersMaps = append(accountIdentifiersMaps, accountIdentifiersMap)
			}
		}
		configMap["account_identifiers"] = accountIdentifiersMaps
		regularRulesRaw := make([]interface{}, 0)
		if configRaw["regularRules"] != nil {
			regularRulesRaw = convertToInterfaceArray(configRaw["regularRules"])
		}

		configMap["bypass_regular_rules"] = regularRulesRaw
		regularTypesRaw := make([]interface{}, 0)
		if configRaw["regularTypes"] != nil {
			regularTypesRaw = convertToInterfaceArray(configRaw["regularTypes"])
		}

		configMap["bypass_regular_types"] = regularTypesRaw
		tagsRaw := make([]interface{}, 0)
		if configRaw["tags"] != nil {
			tagsRaw = convertToInterfaceArray(configRaw["tags"])
		}

		configMap["bypass_tags"] = tagsRaw
		conditionsRaw := configRaw["conditions"]
		conditionsMaps := make([]map[string]interface{}, 0)
		if conditionsRaw != nil {
			for _, conditionsChildRaw := range convertToInterfaceArray(conditionsRaw) {
				conditionsMap := make(map[string]interface{})
				conditionsChildRaw := conditionsChildRaw.(map[string]interface{})
				conditionsMap["key"] = conditionsChildRaw["key"]
				conditionsMap["op_value"] = conditionsChildRaw["opValue"]
				conditionsMap["sub_key"] = conditionsChildRaw["subKey"]
				conditionsMap["values"] = conditionsChildRaw["values"]

				conditionsMaps = append(conditionsMaps, conditionsMap)
			}
		}
		configMap["conditions"] = conditionsMaps
		grayConfigMaps := make([]map[string]interface{}, 0)
		grayConfigMap := make(map[string]interface{})
		grayConfigRaw := make(map[string]interface{})
		if configRaw["grayConfig"] != nil {
			grayConfigRaw = configRaw["grayConfig"].(map[string]interface{})
		}
		if len(grayConfigRaw) > 0 {
			grayConfigMap["gray_rate"] = formatInt(grayConfigRaw["grayRate"])
			grayConfigMap["gray_sub_key"] = grayConfigRaw["graySubKey"]
			grayConfigMap["gray_target"] = grayConfigRaw["grayTarget"]

			grayConfigMaps = append(grayConfigMaps, grayConfigMap)
		}
		configMap["gray_config"] = grayConfigMaps
		rateLimitMaps := make([]map[string]interface{}, 0)
		rateLimitMap := make(map[string]interface{})
		ratelimitRaw := make(map[string]interface{})
		if configRaw["ratelimit"] != nil {
			ratelimitRaw = configRaw["ratelimit"].(map[string]interface{})
		}
		if len(ratelimitRaw) > 0 {
			rateLimitMap["interval"] = formatInt(ratelimitRaw["interval"])
			rateLimitMap["sub_key"] = ratelimitRaw["subKey"]
			rateLimitMap["target"] = ratelimitRaw["target"]
			rateLimitMap["threshold"] = formatInt(ratelimitRaw["threshold"])
			rateLimitMap["ttl"] = formatInt(ratelimitRaw["ttl"])

			statusMaps := make([]map[string]interface{}, 0)
			statusMap := make(map[string]interface{})
			statusRaw := make(map[string]interface{})
			if ratelimitRaw["status"] != nil {
				statusRaw = ratelimitRaw["status"].(map[string]interface{})
			}
			if len(statusRaw) > 0 {
				statusMap["code"] = formatInt(statusRaw["code"])
				statusMap["count"] = formatInt(statusRaw["count"])
				statusMap["ratio"] = formatInt(statusRaw["ratio"])

				statusMaps = append(statusMaps, statusMap)
			}
			rateLimitMap["status"] = statusMaps
			rateLimitMaps = append(rateLimitMaps, rateLimitMap)
		}
		configMap["rate_limit"] = rateLimitMaps
		remoteAddrRaw := make([]interface{}, 0)
		if configRaw["remoteAddr"] != nil {
			remoteAddrRaw = convertToInterfaceArray(configRaw["remoteAddr"])
		}

		configMap["remote_addr"] = remoteAddrRaw
		timeConfigMaps := make([]map[string]interface{}, 0)
		timeConfigMap := make(map[string]interface{})
		timeConfigRaw := make(map[string]interface{})
		if configRaw["timeConfig"] != nil {
			timeConfigRaw = configRaw["timeConfig"].(map[string]interface{})
		}
		if len(timeConfigRaw) > 0 {
			timeConfigMap["time_scope"] = timeConfigRaw["timeScope"]
			timeConfigMap["time_zone"] = formatInt(timeConfigRaw["timeZone"])

			timePeriodsRaw := timeConfigRaw["timePeriods"]
			timePeriodsMaps := make([]map[string]interface{}, 0)
			if timePeriodsRaw != nil {
				for _, timePeriodsChildRaw := range convertToInterfaceArray(timePeriodsRaw) {
					timePeriodsMap := make(map[string]interface{})
					timePeriodsChildRaw := timePeriodsChildRaw.(map[string]interface{})
					timePeriodsMap["end"] = formatInt(timePeriodsChildRaw["end"])
					timePeriodsMap["start"] = formatInt(timePeriodsChildRaw["start"])

					timePeriodsMaps = append(timePeriodsMaps, timePeriodsMap)
				}
			}
			timeConfigMap["time_periods"] = timePeriodsMaps
			weekTimePeriodsRaw := timeConfigRaw["weekTimePeriods"]
			weekTimePeriodsMaps := make([]map[string]interface{}, 0)
			if weekTimePeriodsRaw != nil {
				for _, weekTimePeriodsChildRaw := range convertToInterfaceArray(weekTimePeriodsRaw) {
					weekTimePeriodsMap := make(map[string]interface{})
					weekTimePeriodsChildRaw := weekTimePeriodsChildRaw.(map[string]interface{})
					weekTimePeriodsMap["day"] = weekTimePeriodsChildRaw["day"]

					dayPeriodsRaw := weekTimePeriodsChildRaw["dayPeriods"]
					dayPeriodsMaps := make([]map[string]interface{}, 0)
					if dayPeriodsRaw != nil {
						for _, dayPeriodsChildRaw := range convertToInterfaceArray(dayPeriodsRaw) {
							dayPeriodsMap := make(map[string]interface{})
							dayPeriodsChildRaw := dayPeriodsChildRaw.(map[string]interface{})
							dayPeriodsMap["end"] = formatInt(dayPeriodsChildRaw["end"])
							dayPeriodsMap["start"] = formatInt(dayPeriodsChildRaw["start"])

							dayPeriodsMaps = append(dayPeriodsMaps, dayPeriodsMap)
						}
					}
					weekTimePeriodsMap["day_periods"] = dayPeriodsMaps
					weekTimePeriodsMaps = append(weekTimePeriodsMaps, weekTimePeriodsMap)
				}
			}
			timeConfigMap["week_time_periods"] = weekTimePeriodsMaps
			timeConfigMaps = append(timeConfigMaps, timeConfigMap)
		}
		configMap["time_config"] = timeConfigMaps
		configMaps = append(configMaps, configMap)
	}
	if err := d.Set("config", configMaps); err != nil {
		return err
	}

	parts := strings.Split(d.Id(), ":")
	d.Set("instance_id", parts[0])

	return nil
}

func resourceAliCloudWafv3DefenseRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	d.Partial(true)

	var err error
	parts := strings.Split(d.Id(), ":")
	action := "ModifyDefenseRule"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["InstanceId"] = parts[0]
	request["DefenseType"] = parts[1]
	request["RegionId"] = client.RegionId
	dataList := make(map[string]interface{})

	if d.HasChange("config") {
		update = true
	}
	mode1, _ := jsonpath.Get("$[0].mode", d.Get("config"))
	if mode1 != nil && (d.HasChange("config.0.mode") || mode1 != "") {
		dataList["mode"] = mode1
	}

	if d.HasChange("config") {
		update = true
	}
	if v := d.Get("config"); v != nil {
		ratelimit := make(map[string]interface{})
		interval1, _ := jsonpath.Get("$[0].rate_limit[0].interval", d.Get("config"))
		if interval1 != nil && (d.HasChange("config.0.rate_limit.0.interval") || interval1 != "") && interval1.(int) > 0 {
			ratelimit["interval"] = interval1
		}
		status := make(map[string]interface{})
		count1, _ := jsonpath.Get("$[0].rate_limit[0].status[0].count", d.Get("config"))
		if count1 != nil && (d.HasChange("config.0.rate_limit.0.status.0.count") || count1 != "") && count1.(int) > 0 {
			status["count"] = count1
		}
		code1, _ := jsonpath.Get("$[0].rate_limit[0].status[0].code", d.Get("config"))
		if code1 != nil && (d.HasChange("config.0.rate_limit.0.status.0.code") || code1 != "") {
			status["code"] = code1
		}
		ratio1, _ := jsonpath.Get("$[0].rate_limit[0].status[0].ratio", d.Get("config"))
		if ratio1 != nil && (d.HasChange("config.0.rate_limit.0.status.0.ratio") || ratio1 != "") && ratio1.(int) > 0 {
			status["ratio"] = ratio1
		}

		if len(status) > 0 {
			ratelimit["status"] = status
		}
		target1, _ := jsonpath.Get("$[0].rate_limit[0].target", d.Get("config"))
		if target1 != nil && (d.HasChange("config.0.rate_limit.0.target") || target1 != "") {
			ratelimit["target"] = target1
		}
		threshold1, _ := jsonpath.Get("$[0].rate_limit[0].threshold", d.Get("config"))
		if threshold1 != nil && (d.HasChange("config.0.rate_limit.0.threshold") || threshold1 != "") {
			ratelimit["threshold"] = threshold1
		}
		ttl1, _ := jsonpath.Get("$[0].rate_limit[0].ttl", d.Get("config"))
		if ttl1 != nil && (d.HasChange("config.0.rate_limit.0.ttl") || ttl1 != "") && ttl1.(int) > 0 {
			ratelimit["ttl"] = ttl1
		}
		subKey1, _ := jsonpath.Get("$[0].rate_limit[0].sub_key", d.Get("config"))
		if subKey1 != nil && (d.HasChange("config.0.rate_limit.0.sub_key") || subKey1 != "") {
			ratelimit["subKey"] = subKey1
		}

		if len(ratelimit) > 0 {
			dataList["ratelimit"] = ratelimit
		}
	}

	if d.HasChange("config") {
		update = true
	}
	grayStatus1, _ := jsonpath.Get("$[0].gray_status", d.Get("config"))
	if grayStatus1 != nil && (d.HasChange("config.0.gray_status") || grayStatus1 != "") {
		dataList["grayStatus"] = grayStatus1
	}

	if d.HasChange("config") {
		update = true
	}
	if v := d.Get("config"); v != nil {
		timeConfig := make(map[string]interface{})
		if v, ok := d.GetOk("config"); ok {
			localData, err := jsonpath.Get("$[0].time_config[0].time_periods", v)
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
				dataLoopMap["start"] = dataLoopTmp["start"]
				dataLoopMap["end"] = dataLoopTmp["end"]
				localMaps = append(localMaps, dataLoopMap)
			}
			timeConfig["timePeriods"] = localMaps
		}

		timeZone1, _ := jsonpath.Get("$[0].time_config[0].time_zone", d.Get("config"))
		if timeZone1 != nil && (d.HasChange("config.0.time_config.0.time_zone") || timeZone1 != "") {
			timeConfig["timeZone"] = timeZone1
		}
		if v, ok := d.GetOk("config"); ok {
			localData1, err := jsonpath.Get("$[0].time_config[0].week_time_periods", v)
			if err != nil {
				localData1 = make([]interface{}, 0)
			}
			localMaps1 := make([]interface{}, 0)
			for _, dataLoop1 := range convertToInterfaceArray(localData1) {
				dataLoop1Tmp := make(map[string]interface{})
				if dataLoop1 != nil {
					dataLoop1Tmp = dataLoop1.(map[string]interface{})
				}
				dataLoop1Map := make(map[string]interface{})
				localMaps2 := make([]interface{}, 0)
				localData2 := dataLoop1Tmp["day_periods"]
				for _, dataLoop2 := range convertToInterfaceArray(localData2) {
					dataLoop2Tmp := dataLoop2.(map[string]interface{})
					dataLoop2Map := make(map[string]interface{})
					dataLoop2Map["start"] = dataLoop2Tmp["start"]
					dataLoop2Map["end"] = dataLoop2Tmp["end"]
					localMaps2 = append(localMaps2, dataLoop2Map)
				}
				dataLoop1Map["dayPeriods"] = localMaps2
				dataLoop1Map["day"] = dataLoop1Tmp["day"]
				localMaps1 = append(localMaps1, dataLoop1Map)
			}
			timeConfig["weekTimePeriods"] = localMaps1
		}

		timeScope1, _ := jsonpath.Get("$[0].time_config[0].time_scope", d.Get("config"))
		if timeScope1 != nil && (d.HasChange("config.0.time_config.0.time_scope") || timeScope1 != "") {
			timeConfig["timeScope"] = timeScope1
		}

		if len(timeConfig) > 0 {
			dataList["timeConfig"] = timeConfig
		}
	}

	if d.HasChange("config") {
		update = true
	}
	if v := d.Get("config"); v != nil {
		grayConfig := make(map[string]interface{})
		grayTarget1, _ := jsonpath.Get("$[0].gray_config[0].gray_target", d.Get("config"))
		if grayTarget1 != nil && (d.HasChange("config.0.gray_config.0.gray_target") || grayTarget1 != "") {
			grayConfig["grayTarget"] = grayTarget1
		}
		graySubKey1, _ := jsonpath.Get("$[0].gray_config[0].gray_sub_key", d.Get("config"))
		if graySubKey1 != nil && (d.HasChange("config.0.gray_config.0.gray_sub_key") || graySubKey1 != "") {
			grayConfig["graySubKey"] = graySubKey1
		}
		grayRate1, _ := jsonpath.Get("$[0].gray_config[0].gray_rate", d.Get("config"))
		if grayRate1 != nil && (d.HasChange("config.0.gray_config.0.gray_rate") || grayRate1 != "") && grayRate1.(int) > 0 {
			grayConfig["grayRate"] = grayRate1
		}

		if len(grayConfig) > 0 {
			dataList["grayConfig"] = grayConfig
		}
	}

	if d.HasChange("config") {
		update = true
	}
	ccEffect, _ := jsonpath.Get("$[0].cc_effect", d.Get("config"))
	if ccEffect != nil && (d.HasChange("config.0.cc_effect") || ccEffect != "") {
		dataList["effect"] = ccEffect
	}

	if d.HasChange("config") {
		update = true
	}
	throttleType, _ := jsonpath.Get("$[0].throttle_type", d.Get("config"))
	if throttleType != nil && (d.HasChange("config.0.throttle_type") || throttleType != "") {
		dataList["type"] = throttleType
	}

	if d.HasChange("config") {
		update = true
	}
	throttleThrehold, _ := jsonpath.Get("$[0].throttle_threhold", d.Get("config"))
	if throttleThrehold != nil && (d.HasChange("config.0.throttle_threhold") || throttleThrehold != "") && throttleThrehold.(int) > 0 {
		dataList["threshold"] = throttleThrehold
	}

	if d.HasChange("rule_name") {
		update = true
	}
	if v, ok := d.GetOk("rule_name"); ok {
		dataList["name"] = v
	}

	if d.HasChange("config") {
		update = true
	}
	url1, _ := jsonpath.Get("$[0].url", d.Get("config"))
	if url1 != nil && (d.HasChange("config.0.url") || url1 != "") {
		dataList["url"] = url1
	}

	if d.HasChange("config") {
		update = true
	}
	protocol1, _ := jsonpath.Get("$[0].protocol", d.Get("config"))
	if protocol1 != nil && (d.HasChange("config.0.protocol") || protocol1 != "") {
		dataList["protocol"] = protocol1
	}

	if d.HasChange("config") {
		update = true
	}
	ua1, _ := jsonpath.Get("$[0].ua", d.Get("config"))
	if ua1 != nil && (d.HasChange("config.0.ua") || ua1 != "") {
		dataList["ua"] = ua1
	}

	if d.HasChange("config") {
		update = true
	}
	if v := d.Get("config"); v != nil {
		if v, ok := d.GetOk("config"); ok {
			localData3, err := jsonpath.Get("$[0].account_identifiers", v)
			if err != nil {
				localData3 = make([]interface{}, 0)
			}
			localMaps3 := make([]interface{}, 0)
			for _, dataLoop3 := range convertToInterfaceArray(localData3) {
				dataLoop3Tmp := make(map[string]interface{})
				if dataLoop3 != nil {
					dataLoop3Tmp = dataLoop3.(map[string]interface{})
				}
				dataLoop3Map := make(map[string]interface{})
				dataLoop3Map["decodeType"] = dataLoop3Tmp["decode_type"]
				dataLoop3Map["position"] = dataLoop3Tmp["position"]
				dataLoop3Map["priority"] = dataLoop3Tmp["priority"]
				dataLoop3Map["subKey"] = dataLoop3Tmp["sub_key"]
				dataLoop3Map["key"] = dataLoop3Tmp["key"]
				localMaps3 = append(localMaps3, dataLoop3Map)
			}
			dataList["accountIdentifiers"] = localMaps3
		}

	}

	if d.HasChange("config") {
		update = true
	}
	ruleAction, _ := jsonpath.Get("$[0].rule_action", d.Get("config"))
	if ruleAction != nil && (d.HasChange("config.0.rule_action") || ruleAction != "") {
		dataList["action"] = ruleAction
	}

	if d.HasChange("config") {
		update = true
	}
	ccStatus1, _ := jsonpath.Get("$[0].cc_status", d.Get("config"))
	if ccStatus1 != nil && (d.HasChange("config.0.cc_status") || ccStatus1 != "") {
		dataList["ccStatus"] = ccStatus1
	}

	if d.HasChange("config") {
		update = true
	}
	if v := d.Get("config"); v != nil {
		if v, ok := d.GetOk("config"); ok {
			localData4, err := jsonpath.Get("$[0].conditions", v)
			if err != nil {
				localData4 = make([]interface{}, 0)
			}
			localMaps4 := make([]interface{}, 0)
			for _, dataLoop4 := range convertToInterfaceArray(localData4) {
				dataLoop4Tmp := make(map[string]interface{})
				if dataLoop4 != nil {
					dataLoop4Tmp = dataLoop4.(map[string]interface{})
				}
				dataLoop4Map := make(map[string]interface{})
				dataLoop4Map["values"] = dataLoop4Tmp["values"]
				dataLoop4Map["opValue"] = dataLoop4Tmp["op_value"]
				dataLoop4Map["subKey"] = dataLoop4Tmp["sub_key"]
				dataLoop4Map["key"] = dataLoop4Tmp["key"]
				localMaps4 = append(localMaps4, dataLoop4Map)
			}
			dataList["conditions"] = localMaps4
		}

	}

	if d.HasChange("config") {
		update = true
	}
	bypassRegularTypes, _ := jsonpath.Get("$[0].bypass_regular_types", d.Get("config"))
	if bypassRegularTypes != nil && (d.HasChange("config.0.bypass_regular_types") || bypassRegularTypes != "") {
		dataList["regularTypes"] = convertToInterfaceArray(bypassRegularTypes)
	}

	dataList["id"] = parts[2]
	if d.HasChange("rule_id") {
		update = true
	}
	if v, ok := d.GetOk("rule_id"); ok {
		dataList["id"] = v
	}

	if d.HasChange("config") {
		update = true
	}
	bypassTags, _ := jsonpath.Get("$[0].bypass_tags", d.Get("config"))
	if bypassTags != nil && (d.HasChange("config.0.bypass_tags") || bypassTags != "") {
		dataList["tags"] = convertToInterfaceArray(bypassTags)
	}

	if d.HasChange("config") {
		update = true
	}
	abroadRegions, _ := jsonpath.Get("$[0].abroad_regions", d.Get("config"))
	if abroadRegions != nil && (d.HasChange("config.0.abroad_regions") || abroadRegions != "") {
		dataList["abroadRegionList"] = abroadRegions
	}

	if d.HasChange("config") {
		update = true
	}
	remoteAddr1, _ := jsonpath.Get("$[0].remote_addr", d.Get("config"))
	if remoteAddr1 != nil && (d.HasChange("config.0.remote_addr") || remoteAddr1 != "") {
		dataList["remoteAddr"] = convertToInterfaceArray(remoteAddr1)
	}

	if d.HasChange("config") {
		update = true
	}
	bypassRegularRules, _ := jsonpath.Get("$[0].bypass_regular_rules", d.Get("config"))
	if bypassRegularRules != nil && (d.HasChange("config.0.bypass_regular_rules") || bypassRegularRules != "") {
		dataList["regularRules"] = convertToInterfaceArray(bypassRegularRules)
	}

	if d.HasChange("config") {
		update = true
	}
	cnRegions, _ := jsonpath.Get("$[0].cn_regions", d.Get("config"))
	if cnRegions != nil && (d.HasChange("config.0.cn_regions") || cnRegions != "") {
		dataList["cnRegionList"] = cnRegions
	}

	RulesMap := make([]map[string]interface{}, 0)
	RulesMap = append(RulesMap, dataList)
	request["Rules"], _ = convertListMapToJsonString(RulesMap)

	request["InstanceId"] = parts[0]
	request["DefenseType"] = parts[1]

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("waf-openapi", "2021-10-01", action, query, request, true)
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
	update = false
	parts = strings.Split(d.Id(), ":")
	action = "ModifyDefenseRuleStatus"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["InstanceId"] = parts[0]
	request["RuleId"] = parts[2]
	request["DefenseType"] = parts[1]
	request["RegionId"] = client.RegionId
	if d.HasChange("rule_status") {
		update = true
	}
	request["RuleStatus"] = d.Get("rule_status")
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("waf-openapi", "2021-10-01", action, query, request, true)
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

	d.Partial(false)
	return resourceAliCloudWafv3DefenseRuleRead(d, meta)
}

func resourceAliCloudWafv3DefenseRuleDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeleteDefenseRule"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["InstanceId"] = parts[0]
	request["RuleIds"] = parts[2]
	request["DefenseType"] = parts[1]
	request["RegionId"] = client.RegionId

	if v, ok := d.GetOkExists("template_id"); ok {
		request["TemplateId"] = v
	}
	if v, ok := d.GetOk("resource"); ok {
		request["Resource"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("waf-openapi", "2021-10-01", action, query, request, true)
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
		if IsExpectedErrors(err, []string{"Defense.Control.DefenseRuleNotExist"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
