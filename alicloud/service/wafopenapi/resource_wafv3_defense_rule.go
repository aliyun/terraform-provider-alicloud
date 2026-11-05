// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package wafopenapi

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/validators"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceWafv3DefenseRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceWafv3DefenseRuleCreate,
		Read:   resourceWafv3DefenseRuleRead,
		Update: resourceWafv3DefenseRuleUpdate,
		Delete: resourceWafv3DefenseRuleDelete,
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
										ValidateFunc: validators.StringInSlice([]string{"cookie.acw_tc", "header", "queryarg", "cookie", "remote_addr"}, false),
									},
									"gray_rate": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validators.IntBetween(0, 100),
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
							ValidateFunc: validators.StringInSlice([]string{"block", "monitor", "js", "captcha", "captcha_strict", "filter"}, false),
						},
						"mode": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validators.IntInSlice([]int{0, 1}),
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
										ValidateFunc: validators.IntBetween(-12, 12),
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
										ValidateFunc: validators.StringInSlice([]string{"permanent", "period", "cycle"}, false),
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
										ValidateFunc: validators.StringInSlice([]string{"not-contain", "contain", "none", "ne", "eq", "lt", "gt", "len-lt", "len-eq", "len-gt", "not-match", "match-one", "all-not-match", "all-not-contain", "contain-one", "not-regex", "regex", "all-not-regex", "regex-one", "prefix-match", "suffix-match", "empty", "exists", "inl", "in-list", "not-in-list"}, false),
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
													ValidateFunc: validators.IntBetween(0, 100),
												},
												"count": {
													Type:         schema.TypeInt,
													Optional:     true,
													ValidateFunc: validators.IntBetween(2, 50000),
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
										ValidateFunc: validators.StringInSlice([]string{"cookie.acw_tc", "header", "queryarg", "cookie", "remote_addr"}, false),
									},
									"ttl": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validators.IntBetween(60, 86400),
									},
									"interval": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validators.IntBetween(0, 1800),
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
							Computed:     true,
							ValidateFunc: validators.StringInSlice([]string{"service", "rule"}, false),
						},
						"throttle_threhold": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validators.IntBetween(0, 5000000),
						},
						"throttle_type": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validators.StringInSlice([]string{"qps", "ratio"}, false),
						},
						"auto_update": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"cc_status": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validators.IntInSlice([]int{0, 1}),
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
										ValidateFunc: validators.IntBetween(0, 20),
									},
									"decode_type": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validators.StringInSlice([]string{"plain", "basic", "jwt"}, false),
									},
									"key": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validators.StringInSlice([]string{"Query-Arg", "Cookie-Exact", "Post-Arg", "Header"}, false),
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
							ValidateFunc: validators.IntInSlice([]int{0, 1}),
						},
						"waf_base_config": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"rule_type": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validators.StringInSlice([]string{"system", "custom"}, false),
									},
									"rule_detail": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"rule_id": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"rule_action": {
													Type:         schema.TypeString,
													Optional:     true,
													ValidateFunc: validators.StringInSlice([]string{"block", "monitor"}, false),
												},
												"rule_status": {
													Type:         schema.TypeInt,
													Optional:     true,
													ValidateFunc: validators.IntInSlice([]int{0, 1}),
												},
											},
										},
									},
									"rule_batch_operation_config": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validators.StringInSlice([]string{"default", "all_on", "all_off", "all_block", "all_monitor"}, false),
									},
								},
							},
						},
						"protocol": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validators.StringInSlice([]string{"https", "http"}, false),
						},
						"codec_list": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"defense_origin": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validators.StringInSlice([]string{"custom"}, false),
			},
			"defense_scene": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validators.StringInSlice([]string{"ip_blacklist", "custom_acl", "whitelist", "region_block", "cc", "tamperproof", "dlp", "spike_throttle", "account_identifier", "waf_base", "waf_codec"}, false),
			},
			"defense_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validators.StringInSlice([]string{"template", "resource", "global"}, false),
			},
			"gmt_modified": {
				Type:     schema.TypeString,
				Computed: true,
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
				ValidateFunc: validators.IntInSlice([]int{0, 1}),
			},
			"template_id": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceWafv3DefenseRuleCreate(d *schema.ResourceData, meta interface{}) error {

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

	rulesDataList := make(map[string]interface{})

	if v, ok := d.GetOk("config"); ok {
		mode1, _ := jsonpath.Get("$[0].mode", v)
		if mode1 != nil && mode1 != "" {
			rulesDataList["mode"] = mode1
		}
	}

	if v := d.Get("config"); !helper.IsNil(v) {
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
			rulesDataList["ratelimit"] = ratelimit
		}
	}

	if v, ok := d.GetOk("config"); ok {
		grayStatus1, _ := jsonpath.Get("$[0].gray_status", v)
		if grayStatus1 != nil && grayStatus1 != "" {
			rulesDataList["grayStatus"] = grayStatus1
		}
	}

	if v := d.Get("config"); !helper.IsNil(v) {
		timeConfig := make(map[string]interface{})
		localData, err := jsonpath.Get("$[0].time_config[0].time_periods", v)
		if err != nil {
			localData = make([]interface{}, 0)
		}
		localMaps := make([]interface{}, 0)
		for _, dataLoop := range convert.ConvertToInterfaceArray(localData) {
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

		timeZone1, _ := jsonpath.Get("$[0].time_config[0].time_zone", d.Get("config"))
		if timeZone1 != nil && timeZone1 != "" {
			timeConfig["timeZone"] = timeZone1
		}
		localData1, err := jsonpath.Get("$[0].time_config[0].week_time_periods", v)
		if err != nil {
			localData1 = make([]interface{}, 0)
		}
		localMaps1 := make([]interface{}, 0)
		for _, dataLoop1 := range convert.ConvertToInterfaceArray(localData1) {
			dataLoop1Tmp := make(map[string]interface{})
			if dataLoop1 != nil {
				dataLoop1Tmp = dataLoop1.(map[string]interface{})
			}
			dataLoop1Map := make(map[string]interface{})
			localMaps2 := make([]interface{}, 0)
			localData2 := dataLoop1Tmp["day_periods"]
			for _, dataLoop2 := range convert.ConvertToInterfaceArray(localData2) {
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

		timeScope1, _ := jsonpath.Get("$[0].time_config[0].time_scope", d.Get("config"))
		if timeScope1 != nil && timeScope1 != "" {
			timeConfig["timeScope"] = timeScope1
		}

		if len(timeConfig) > 0 {
			rulesDataList["timeConfig"] = timeConfig
		}
	}

	if v := d.Get("config"); !helper.IsNil(v) {
		localData3, err := jsonpath.Get("$[0].waf_base_config", v)
		if err != nil {
			localData3 = make([]interface{}, 0)
		}
		localMaps3 := make([]interface{}, 0)
		for _, dataLoop3 := range convert.ConvertToInterfaceArray(localData3) {
			dataLoop3Tmp := make(map[string]interface{})
			if dataLoop3 != nil {
				dataLoop3Tmp = dataLoop3.(map[string]interface{})
			}
			dataLoop3Map := make(map[string]interface{})
			dataLoop3Map["ruleBatchOperationConfig"] = dataLoop3Tmp["rule_batch_operation_config"]
			localMaps4 := make([]interface{}, 0)
			localData4 := dataLoop3Tmp["rule_detail"]
			for _, dataLoop4 := range convert.ConvertToInterfaceArray(localData4) {
				dataLoop4Tmp := dataLoop4.(map[string]interface{})
				dataLoop4Map := make(map[string]interface{})
				dataLoop4Map["ruleId"] = dataLoop4Tmp["rule_id"]
				dataLoop4Map["ruleStatus"] = dataLoop4Tmp["rule_status"]
				dataLoop4Map["ruleAction"] = dataLoop4Tmp["rule_action"]
				localMaps4 = append(localMaps4, dataLoop4Map)
			}
			dataLoop3Map["ruleDetail"] = localMaps4
			dataLoop3Map["ruleType"] = dataLoop3Tmp["rule_type"]
			localMaps3 = append(localMaps3, dataLoop3Map)
		}
		rulesDataList["config"] = localMaps3

	}

	if v := d.Get("config"); !helper.IsNil(v) {
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
			rulesDataList["grayConfig"] = grayConfig
		}
	}

	if v, ok := d.GetOk("config"); ok {
		ccEffect, _ := jsonpath.Get("$[0].cc_effect", v)
		if ccEffect != nil && ccEffect != "" {
			rulesDataList["effect"] = ccEffect
		}
	}

	if v, ok := d.GetOk("config"); ok {
		throttleType, _ := jsonpath.Get("$[0].throttle_type", v)
		if throttleType != nil && throttleType != "" {
			rulesDataList["type"] = throttleType
		}
	}

	if v, ok := d.GetOk("config"); ok {
		throttleThrehold, _ := jsonpath.Get("$[0].throttle_threhold", v)
		if throttleThrehold != nil && throttleThrehold != "" && throttleThrehold.(int) > 0 {
			rulesDataList["threshold"] = throttleThrehold
		}
	}

	if v, ok := d.GetOk("rule_name"); ok {
		rulesDataList["name"] = v
	}

	if v, ok := d.GetOk("config"); ok {
		url1, _ := jsonpath.Get("$[0].url", v)
		if url1 != nil && url1 != "" {
			rulesDataList["url"] = url1
		}
	}

	if v, ok := d.GetOk("config"); ok {
		protocol1, _ := jsonpath.Get("$[0].protocol", v)
		if protocol1 != nil && protocol1 != "" {
			rulesDataList["protocol"] = protocol1
		}
	}

	if v, ok := d.GetOk("config"); ok {
		ua1, _ := jsonpath.Get("$[0].ua", v)
		if ua1 != nil && ua1 != "" {
			rulesDataList["ua"] = ua1
		}
	}

	if v, ok := d.GetOk("config"); ok {
		autoUpdate1, _ := jsonpath.Get("$[0].auto_update", v)
		if autoUpdate1 != nil && autoUpdate1 != "" {
			rulesDataList["autoUpdate"] = autoUpdate1
		}
	}

	if v := d.Get("config"); !helper.IsNil(v) {
		localData5, err := jsonpath.Get("$[0].account_identifiers", v)
		if err != nil {
			localData5 = make([]interface{}, 0)
		}
		localMaps5 := make([]interface{}, 0)
		for _, dataLoop5 := range convert.ConvertToInterfaceArray(localData5) {
			dataLoop5Tmp := make(map[string]interface{})
			if dataLoop5 != nil {
				dataLoop5Tmp = dataLoop5.(map[string]interface{})
			}
			dataLoop5Map := make(map[string]interface{})
			dataLoop5Map["decodeType"] = dataLoop5Tmp["decode_type"]
			dataLoop5Map["position"] = dataLoop5Tmp["position"]
			dataLoop5Map["priority"] = dataLoop5Tmp["priority"]
			dataLoop5Map["subKey"] = dataLoop5Tmp["sub_key"]
			dataLoop5Map["key"] = dataLoop5Tmp["key"]
			localMaps5 = append(localMaps5, dataLoop5Map)
		}
		rulesDataList["accountIdentifiers"] = localMaps5

	}

	if v, ok := d.GetOk("config"); ok {
		ruleAction2, _ := jsonpath.Get("$[0].rule_action", v)
		if ruleAction2 != nil && ruleAction2 != "" {
			rulesDataList["action"] = ruleAction2
		}
	}

	if v, ok := d.GetOk("defense_origin"); ok {
		rulesDataList["origin"] = v
	}

	if v, ok := d.GetOk("config"); ok {
		ccStatus1, _ := jsonpath.Get("$[0].cc_status", v)
		if ccStatus1 != nil && ccStatus1 != "" {
			rulesDataList["ccStatus"] = ccStatus1
		}
	}

	if v, ok := d.GetOk("config"); ok {
		codecList1, _ := jsonpath.Get("$[0].codec_list", v)
		if codecList1 != nil && codecList1 != "" {
			rulesDataList["codecList"] = convert.ConvertToInterfaceArray(codecList1)
		}
	}

	if v := d.Get("config"); !helper.IsNil(v) {
		localData6, err := jsonpath.Get("$[0].conditions", v)
		if err != nil {
			localData6 = make([]interface{}, 0)
		}
		localMaps6 := make([]interface{}, 0)
		for _, dataLoop6 := range convert.ConvertToInterfaceArray(localData6) {
			dataLoop6Tmp := make(map[string]interface{})
			if dataLoop6 != nil {
				dataLoop6Tmp = dataLoop6.(map[string]interface{})
			}
			dataLoop6Map := make(map[string]interface{})
			dataLoop6Map["values"] = dataLoop6Tmp["values"]
			dataLoop6Map["opValue"] = dataLoop6Tmp["op_value"]
			dataLoop6Map["subKey"] = dataLoop6Tmp["sub_key"]
			dataLoop6Map["key"] = dataLoop6Tmp["key"]
			localMaps6 = append(localMaps6, dataLoop6Map)
		}
		rulesDataList["conditions"] = localMaps6

	}

	if v, ok := d.GetOkExists("rule_status"); ok {
		rulesDataList["status"] = v
	}

	if v, ok := d.GetOk("config"); ok {
		bypassRegularTypes, _ := jsonpath.Get("$[0].bypass_regular_types", v)
		if bypassRegularTypes != nil && bypassRegularTypes != "" {
			rulesDataList["regularTypes"] = convert.ConvertToInterfaceArray(bypassRegularTypes)
		}
	}

	if v, ok := d.GetOk("config"); ok {
		bypassTags, _ := jsonpath.Get("$[0].bypass_tags", v)
		if bypassTags != nil && bypassTags != "" {
			rulesDataList["tags"] = convert.ConvertToInterfaceArray(bypassTags)
		}
	}

	if v, ok := d.GetOk("config"); ok {
		abroadRegions, _ := jsonpath.Get("$[0].abroad_regions", v)
		if abroadRegions != nil && abroadRegions != "" {
			rulesDataList["abroadRegionList"] = abroadRegions
		}
	}

	if v, ok := d.GetOk("config"); ok {
		remoteAddr1, _ := jsonpath.Get("$[0].remote_addr", v)
		if remoteAddr1 != nil && remoteAddr1 != "" {
			rulesDataList["remoteAddr"] = convert.ConvertToInterfaceArray(remoteAddr1)
		}
	}

	if v, ok := d.GetOk("config"); ok {
		bypassRegularRules, _ := jsonpath.Get("$[0].bypass_regular_rules", v)
		if bypassRegularRules != nil && bypassRegularRules != "" {
			rulesDataList["regularRules"] = convert.ConvertToInterfaceArray(bypassRegularRules)
		}
	}

	if v, ok := d.GetOk("config"); ok {
		cnRegions, _ := jsonpath.Get("$[0].cn_regions", v)
		if cnRegions != nil && cnRegions != "" {
			rulesDataList["cnRegionList"] = cnRegions
		}
	}

	RulesMap := make([]map[string]interface{}, 0)
	RulesMap = append(RulesMap, rulesDataList)
	request["Rules"], _ = convert.ConvertListMapToJsonString(RulesMap)

	if v, ok := d.GetOk("resource"); ok {
		request["Resource"] = v
	}
	if v, ok := d.GetOkExists("template_id"); ok {
		request["TemplateId"] = v
	}
	request["DefenseScene"] = d.Get("defense_scene")
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("waf-openapi", "2021-10-01", action, query, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_wafv3_defense_rule", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v:%v", request["InstanceId"], request["DefenseType"], response["RuleIds"]))

	return resourceWafv3DefenseRuleRead(d, meta)
}

func resourceWafv3DefenseRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	wafv3ServiceV2 := Wafv3ServiceV2{client}

	objectRaw, err := wafv3ServiceV2.DescribeWafv3DefenseRule(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_wafv3_defense_rule DescribeWafv3DefenseRule Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("defense_origin", objectRaw["DefenseOrigin"])
	d.Set("defense_scene", objectRaw["DefenseScene"])
	d.Set("gmt_modified", objectRaw["GmtModified"])
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
		configRaw = convert.ConvertJsonStringToObject(objectRaw["Config"])
	}
	if len(configRaw) > 0 {
		configMap["abroad_regions"] = configRaw["abroadRegionList"]
		actualAutoUpdate, err := convert.FormatBool(configRaw["autoUpdate"])
		if err != nil {
			return err
		}
		configMap["auto_update"] = actualAutoUpdate
		configMap["cc_effect"] = configRaw["effect"]
		actualCcStatus, err := convert.FormatInt(configRaw["ccStatus"])
		if err != nil {
			return err
		}
		configMap["cc_status"] = actualCcStatus
		configMap["cn_regions"] = configRaw["cnRegionList"]
		actualGrayStatus, err := convert.FormatInt(configRaw["grayStatus"])
		if err != nil {
			return err
		}
		configMap["gray_status"] = actualGrayStatus
		actualMode, err := convert.FormatInt(configRaw["mode"])
		if err != nil {
			return err
		}
		configMap["mode"] = actualMode
		configMap["protocol"] = configRaw["protocol"]
		configMap["rule_action"] = configRaw["action"]
		throttleThrehold, err := convert.FormatInt(configRaw["threshold"])
		if err != nil {
			return err
		}
		configMap["throttle_threhold"] = throttleThrehold
		configMap["throttle_type"] = configRaw["type"]
		configMap["ua"] = configRaw["ua"]
		configMap["url"] = configRaw["url"]

		accountIdentifiersRaw := configRaw["accountIdentifiers"]
		accountIdentifiersMaps := make([]map[string]interface{}, 0)
		if accountIdentifiersRaw != nil {
			for _, accountIdentifiersChildRaw := range convert.ConvertToInterfaceArray(accountIdentifiersRaw) {
				accountIdentifiersMap := make(map[string]interface{})
				accountIdentifiersChildRaw := accountIdentifiersChildRaw.(map[string]interface{})
				accountIdentifiersMap["decode_type"] = accountIdentifiersChildRaw["decodeType"]
				accountIdentifiersMap["key"] = accountIdentifiersChildRaw["key"]
				accountIdentifiersMap["position"] = accountIdentifiersChildRaw["position"]
				accountIdentifiersChildRawPriority, err := convert.FormatInt(accountIdentifiersChildRaw["priority"])
				if err != nil {
					return err
				}
				accountIdentifiersMap["priority"] = accountIdentifiersChildRawPriority
				accountIdentifiersMap["sub_key"] = accountIdentifiersChildRaw["subKey"]

				accountIdentifiersMaps = append(accountIdentifiersMaps, accountIdentifiersMap)
			}
		}
		configMap["account_identifiers"] = accountIdentifiersMaps
		regularRulesRaw := make([]interface{}, 0)
		if configRaw["regularRules"] != nil {
			regularRulesRaw = convert.ConvertToInterfaceArray(configRaw["regularRules"])
		}

		configMap["bypass_regular_rules"] = regularRulesRaw
		regularTypesRaw := make([]interface{}, 0)
		if configRaw["regularTypes"] != nil {
			regularTypesRaw = convert.ConvertToInterfaceArray(configRaw["regularTypes"])
		}

		configMap["bypass_regular_types"] = regularTypesRaw
		tagsRaw := make([]interface{}, 0)
		if configRaw["tags"] != nil {
			tagsRaw = convert.ConvertToInterfaceArray(configRaw["tags"])
		}

		configMap["bypass_tags"] = tagsRaw
		codecListRaw := make([]interface{}, 0)
		if configRaw["codecList"] != nil {
			codecListRaw = convert.ConvertToInterfaceArray(configRaw["codecList"])
		}

		configMap["codec_list"] = codecListRaw
		conditionsRaw := configRaw["conditions"]
		conditionsMaps := make([]map[string]interface{}, 0)
		if conditionsRaw != nil {
			for _, conditionsChildRaw := range convert.ConvertToInterfaceArray(conditionsRaw) {
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
			grayConfigRawGrayRate, err := convert.FormatInt(grayConfigRaw["grayRate"])
			if err != nil {
				return err
			}
			grayConfigMap["gray_rate"] = grayConfigRawGrayRate
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
			ratelimitRawInterval, err := convert.FormatInt(ratelimitRaw["interval"])
			if err != nil {
				return err
			}
			rateLimitMap["interval"] = ratelimitRawInterval
			rateLimitMap["sub_key"] = ratelimitRaw["subKey"]
			rateLimitMap["target"] = ratelimitRaw["target"]
			ratelimitRawThreshold, err := convert.FormatInt(ratelimitRaw["threshold"])
			if err != nil {
				return err
			}
			rateLimitMap["threshold"] = ratelimitRawThreshold
			ratelimitRawTtl, err := convert.FormatInt(ratelimitRaw["ttl"])
			if err != nil {
				return err
			}
			rateLimitMap["ttl"] = ratelimitRawTtl

			statusMaps := make([]map[string]interface{}, 0)
			statusMap := make(map[string]interface{})
			statusRaw := make(map[string]interface{})
			if ratelimitRaw["status"] != nil {
				statusRaw = ratelimitRaw["status"].(map[string]interface{})
			}
			if len(statusRaw) > 0 {
				actualCode, err := convert.FormatInt(statusRaw["code"])
				if err != nil {
					return err
				}
				statusMap["code"] = actualCode
				actualCount, err := convert.FormatInt(statusRaw["count"])
				if err != nil {
					return err
				}
				statusMap["count"] = actualCount
				actualRatio, err := convert.FormatInt(statusRaw["ratio"])
				if err != nil {
					return err
				}
				statusMap["ratio"] = actualRatio

				statusMaps = append(statusMaps, statusMap)
			}
			rateLimitMap["status"] = statusMaps
			rateLimitMaps = append(rateLimitMaps, rateLimitMap)
		}
		configMap["rate_limit"] = rateLimitMaps
		remoteAddrRaw := make([]interface{}, 0)
		if configRaw["remoteAddr"] != nil {
			remoteAddrRaw = convert.ConvertToInterfaceArray(configRaw["remoteAddr"])
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
			timeConfigRawTimeZone, err := convert.FormatInt(timeConfigRaw["timeZone"])
			if err != nil {
				return err
			}
			timeConfigMap["time_zone"] = timeConfigRawTimeZone

			timePeriodsRaw := timeConfigRaw["timePeriods"]
			timePeriodsMaps := make([]map[string]interface{}, 0)
			if timePeriodsRaw != nil {
				for _, timePeriodsChildRaw := range convert.ConvertToInterfaceArray(timePeriodsRaw) {
					timePeriodsMap := make(map[string]interface{})
					timePeriodsChildRaw := timePeriodsChildRaw.(map[string]interface{})
					timePeriodsChildRawEnd, err := convert.FormatInt(timePeriodsChildRaw["end"])
					if err != nil {
						return err
					}
					timePeriodsMap["end"] = timePeriodsChildRawEnd
					timePeriodsChildRawStart, err := convert.FormatInt(timePeriodsChildRaw["start"])
					if err != nil {
						return err
					}
					timePeriodsMap["start"] = timePeriodsChildRawStart

					timePeriodsMaps = append(timePeriodsMaps, timePeriodsMap)
				}
			}
			timeConfigMap["time_periods"] = timePeriodsMaps
			weekTimePeriodsRaw := timeConfigRaw["weekTimePeriods"]
			weekTimePeriodsMaps := make([]map[string]interface{}, 0)
			if weekTimePeriodsRaw != nil {
				for _, weekTimePeriodsChildRaw := range convert.ConvertToInterfaceArray(weekTimePeriodsRaw) {
					weekTimePeriodsMap := make(map[string]interface{})
					weekTimePeriodsChildRaw := weekTimePeriodsChildRaw.(map[string]interface{})
					weekTimePeriodsMap["day"] = weekTimePeriodsChildRaw["day"]

					dayPeriodsRaw := weekTimePeriodsChildRaw["dayPeriods"]
					dayPeriodsMaps := make([]map[string]interface{}, 0)
					if dayPeriodsRaw != nil {
						for _, dayPeriodsChildRaw := range convert.ConvertToInterfaceArray(dayPeriodsRaw) {
							dayPeriodsMap := make(map[string]interface{})
							dayPeriodsChildRaw := dayPeriodsChildRaw.(map[string]interface{})
							dayPeriodsChildRawEnd, err := convert.FormatInt(dayPeriodsChildRaw["end"])
							if err != nil {
								return err
							}
							dayPeriodsMap["end"] = dayPeriodsChildRawEnd
							dayPeriodsChildRawStart, err := convert.FormatInt(dayPeriodsChildRaw["start"])
							if err != nil {
								return err
							}
							dayPeriodsMap["start"] = dayPeriodsChildRawStart

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
		config1Raw := configRaw["config"]
		wafBaseConfigMaps := make([]map[string]interface{}, 0)
		if config1Raw != nil {
			for _, configChildRaw := range convert.ConvertToInterfaceArray(config1Raw) {
				wafBaseConfigMap := make(map[string]interface{})
				configChildRaw := configChildRaw.(map[string]interface{})
				wafBaseConfigMap["rule_batch_operation_config"] = configChildRaw["ruleBatchOperationConfig"]
				wafBaseConfigMap["rule_type"] = configChildRaw["ruleType"]

				ruleDetailRaw := configChildRaw["ruleDetail"]
				ruleDetailMaps := make([]map[string]interface{}, 0)
				if ruleDetailRaw != nil {
					for _, ruleDetailChildRaw := range convert.ConvertToInterfaceArray(ruleDetailRaw) {
						ruleDetailMap := make(map[string]interface{})
						ruleDetailChildRaw := ruleDetailChildRaw.(map[string]interface{})
						ruleDetailMap["rule_action"] = ruleDetailChildRaw["ruleAction"]
						if v := ruleDetailChildRaw["ruleId"]; v != nil {
							ruleDetailMap["rule_id"] = fmt.Sprint(v)
						}
						ruleDetailChildRawRuleStatus, err := convert.FormatInt(ruleDetailChildRaw["ruleStatus"])
						if err != nil {
							return err
						}
						ruleDetailMap["rule_status"] = ruleDetailChildRawRuleStatus

						ruleDetailMaps = append(ruleDetailMaps, ruleDetailMap)
					}
				}
				wafBaseConfigMap["rule_detail"] = ruleDetailMaps
				wafBaseConfigMaps = append(wafBaseConfigMaps, wafBaseConfigMap)
			}
		}
		configMap["waf_base_config"] = wafBaseConfigMaps
		configMaps = append(configMaps, configMap)
	}
	if err := d.Set("config", configMaps); err != nil {
		return err
	}

	parts := strings.Split(d.Id(), ":")
	d.Set("instance_id", parts[0])

	return nil
}

func resourceWafv3DefenseRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	parts := strings.Split(d.Id(), ":")
	action := "ModifyDefenseRule"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["InstanceId"] = parts[0]
	request["DefenseType"] = parts[1]
	request["RegionId"] = client.RegionId
	rulesDataList := make(map[string]interface{})

	if d.HasChange("config") {
		update = true
	}
	mode1, _ := jsonpath.Get("$[0].mode", d.Get("config"))
	if mode1 != nil && mode1 != "" {
		rulesDataList["mode"] = mode1
	}

	if d.HasChange("config") {
		update = true
	}
	if v := d.Get("config"); v != nil {
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
			rulesDataList["ratelimit"] = ratelimit
		}
	}

	if d.HasChange("config") {
		update = true
	}
	grayStatus1, _ := jsonpath.Get("$[0].gray_status", d.Get("config"))
	if grayStatus1 != nil && grayStatus1 != "" {
		rulesDataList["grayStatus"] = grayStatus1
	}

	if d.HasChange("config") {
		update = true
	}
	if v := d.Get("config"); v != nil {
		timeConfig := make(map[string]interface{})
		localData, err := jsonpath.Get("$[0].time_config[0].time_periods", v)
		if err != nil {
			localData = make([]interface{}, 0)
		}
		localMaps := make([]interface{}, 0)
		for _, dataLoop := range convert.ConvertToInterfaceArray(localData) {
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

		timeZone1, _ := jsonpath.Get("$[0].time_config[0].time_zone", d.Get("config"))
		if timeZone1 != nil && timeZone1 != "" {
			timeConfig["timeZone"] = timeZone1
		}
		localData1, err := jsonpath.Get("$[0].time_config[0].week_time_periods", v)
		if err != nil {
			localData1 = make([]interface{}, 0)
		}
		localMaps1 := make([]interface{}, 0)
		for _, dataLoop1 := range convert.ConvertToInterfaceArray(localData1) {
			dataLoop1Tmp := make(map[string]interface{})
			if dataLoop1 != nil {
				dataLoop1Tmp = dataLoop1.(map[string]interface{})
			}
			dataLoop1Map := make(map[string]interface{})
			localMaps2 := make([]interface{}, 0)
			localData2 := dataLoop1Tmp["day_periods"]
			for _, dataLoop2 := range convert.ConvertToInterfaceArray(localData2) {
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

		timeScope1, _ := jsonpath.Get("$[0].time_config[0].time_scope", d.Get("config"))
		if timeScope1 != nil && timeScope1 != "" {
			timeConfig["timeScope"] = timeScope1
		}

		if len(timeConfig) > 0 {
			rulesDataList["timeConfig"] = timeConfig
		}
	}

	if d.HasChange("config") {
		update = true
	}
	if v := d.Get("config"); v != nil {
		localData3, err := jsonpath.Get("$[0].waf_base_config", v)
		if err != nil {
			localData3 = make([]interface{}, 0)
		}
		localMaps3 := make([]interface{}, 0)
		for _, dataLoop3 := range convert.ConvertToInterfaceArray(localData3) {
			dataLoop3Tmp := make(map[string]interface{})
			if dataLoop3 != nil {
				dataLoop3Tmp = dataLoop3.(map[string]interface{})
			}
			dataLoop3Map := make(map[string]interface{})
			dataLoop3Map["ruleBatchOperationConfig"] = dataLoop3Tmp["rule_batch_operation_config"]
			localMaps4 := make([]interface{}, 0)
			localData4 := dataLoop3Tmp["rule_detail"]
			for _, dataLoop4 := range convert.ConvertToInterfaceArray(localData4) {
				dataLoop4Tmp := dataLoop4.(map[string]interface{})
				dataLoop4Map := make(map[string]interface{})
				dataLoop4Map["ruleId"] = dataLoop4Tmp["rule_id"]
				dataLoop4Map["ruleStatus"] = dataLoop4Tmp["rule_status"]
				dataLoop4Map["ruleAction"] = dataLoop4Tmp["rule_action"]
				localMaps4 = append(localMaps4, dataLoop4Map)
			}
			dataLoop3Map["ruleDetail"] = localMaps4
			dataLoop3Map["ruleType"] = dataLoop3Tmp["rule_type"]
			localMaps3 = append(localMaps3, dataLoop3Map)
		}
		rulesDataList["config"] = localMaps3

	}

	if d.HasChange("config") {
		update = true
	}
	if v := d.Get("config"); v != nil {
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
			rulesDataList["grayConfig"] = grayConfig
		}
	}

	if d.HasChange("config") {
		update = true
	}
	ccEffect, _ := jsonpath.Get("$[0].cc_effect", d.Get("config"))
	if ccEffect != nil && ccEffect != "" {
		rulesDataList["effect"] = ccEffect
	}

	if d.HasChange("config") {
		update = true
	}
	throttleType, _ := jsonpath.Get("$[0].throttle_type", d.Get("config"))
	if throttleType != nil && throttleType != "" {
		rulesDataList["type"] = throttleType
	}

	if d.HasChange("config") {
		update = true
	}
	throttleThrehold, _ := jsonpath.Get("$[0].throttle_threhold", d.Get("config"))
	if throttleThrehold != nil && throttleThrehold != "" && throttleThrehold.(int) > 0 {
		rulesDataList["threshold"] = throttleThrehold
	}

	if d.HasChange("rule_name") {
		update = true
	}
	if v, ok := d.GetOk("rule_name"); ok {
		rulesDataList["name"] = v
	}

	if d.HasChange("config") {
		update = true
	}
	url1, _ := jsonpath.Get("$[0].url", d.Get("config"))
	if url1 != nil && url1 != "" {
		rulesDataList["url"] = url1
	}

	if d.HasChange("config") {
		update = true
	}
	protocol1, _ := jsonpath.Get("$[0].protocol", d.Get("config"))
	if protocol1 != nil && protocol1 != "" {
		rulesDataList["protocol"] = protocol1
	}

	if d.HasChange("config") {
		update = true
	}
	ua1, _ := jsonpath.Get("$[0].ua", d.Get("config"))
	if ua1 != nil && ua1 != "" {
		rulesDataList["ua"] = ua1
	}

	if d.HasChange("config") {
		update = true
	}
	autoUpdate1, _ := jsonpath.Get("$[0].auto_update", d.Get("config"))
	if autoUpdate1 != nil && autoUpdate1 != "" {
		rulesDataList["autoUpdate"] = autoUpdate1
	}

	if d.HasChange("config") {
		update = true
	}
	if v := d.Get("config"); v != nil {
		localData5, err := jsonpath.Get("$[0].account_identifiers", v)
		if err != nil {
			localData5 = make([]interface{}, 0)
		}
		localMaps5 := make([]interface{}, 0)
		for _, dataLoop5 := range convert.ConvertToInterfaceArray(localData5) {
			dataLoop5Tmp := make(map[string]interface{})
			if dataLoop5 != nil {
				dataLoop5Tmp = dataLoop5.(map[string]interface{})
			}
			dataLoop5Map := make(map[string]interface{})
			dataLoop5Map["decodeType"] = dataLoop5Tmp["decode_type"]
			dataLoop5Map["position"] = dataLoop5Tmp["position"]
			dataLoop5Map["priority"] = dataLoop5Tmp["priority"]
			dataLoop5Map["subKey"] = dataLoop5Tmp["sub_key"]
			dataLoop5Map["key"] = dataLoop5Tmp["key"]
			localMaps5 = append(localMaps5, dataLoop5Map)
		}
		rulesDataList["accountIdentifiers"] = localMaps5

	}

	if d.HasChange("config") {
		update = true
	}
	ruleAction2, _ := jsonpath.Get("$[0].rule_action", d.Get("config"))
	if ruleAction2 != nil && ruleAction2 != "" {
		rulesDataList["action"] = ruleAction2
	}

	if d.HasChange("config") {
		update = true
	}
	ccStatus1, _ := jsonpath.Get("$[0].cc_status", d.Get("config"))
	if ccStatus1 != nil && ccStatus1 != "" {
		rulesDataList["ccStatus"] = ccStatus1
	}

	if d.HasChange("config") {
		update = true
	}
	codecList1, _ := jsonpath.Get("$[0].codec_list", d.Get("config"))
	if codecList1 != nil && codecList1 != "" {
		rulesDataList["codecList"] = convert.ConvertToInterfaceArray(codecList1)
	}

	if d.HasChange("config") {
		update = true
	}
	if v := d.Get("config"); v != nil {
		localData6, err := jsonpath.Get("$[0].conditions", v)
		if err != nil {
			localData6 = make([]interface{}, 0)
		}
		localMaps6 := make([]interface{}, 0)
		for _, dataLoop6 := range convert.ConvertToInterfaceArray(localData6) {
			dataLoop6Tmp := make(map[string]interface{})
			if dataLoop6 != nil {
				dataLoop6Tmp = dataLoop6.(map[string]interface{})
			}
			dataLoop6Map := make(map[string]interface{})
			dataLoop6Map["values"] = dataLoop6Tmp["values"]
			dataLoop6Map["opValue"] = dataLoop6Tmp["op_value"]
			dataLoop6Map["subKey"] = dataLoop6Tmp["sub_key"]
			dataLoop6Map["key"] = dataLoop6Tmp["key"]
			localMaps6 = append(localMaps6, dataLoop6Map)
		}
		rulesDataList["conditions"] = localMaps6

	}

	if d.HasChange("config") {
		update = true
	}
	bypassRegularTypes, _ := jsonpath.Get("$[0].bypass_regular_types", d.Get("config"))
	if bypassRegularTypes != nil && bypassRegularTypes != "" {
		rulesDataList["regularTypes"] = convert.ConvertToInterfaceArray(bypassRegularTypes)
	}

	rulesDataList["id"] = parts[2]
	if d.HasChange("rule_id") {
		update = true
	}
	if v, ok := d.GetOkExists("rule_id"); ok {
		rulesDataList["id"] = v
	}

	if d.HasChange("config") {
		update = true
	}
	bypassTags, _ := jsonpath.Get("$[0].bypass_tags", d.Get("config"))
	if bypassTags != nil && bypassTags != "" {
		rulesDataList["tags"] = convert.ConvertToInterfaceArray(bypassTags)
	}

	if d.HasChange("config") {
		update = true
	}
	abroadRegions, _ := jsonpath.Get("$[0].abroad_regions", d.Get("config"))
	if abroadRegions != nil && abroadRegions != "" {
		rulesDataList["abroadRegionList"] = abroadRegions
	}

	if d.HasChange("config") {
		update = true
	}
	remoteAddr1, _ := jsonpath.Get("$[0].remote_addr", d.Get("config"))
	if remoteAddr1 != nil && remoteAddr1 != "" {
		rulesDataList["remoteAddr"] = convert.ConvertToInterfaceArray(remoteAddr1)
	}

	if d.HasChange("config") {
		update = true
	}
	bypassRegularRules, _ := jsonpath.Get("$[0].bypass_regular_rules", d.Get("config"))
	if bypassRegularRules != nil && bypassRegularRules != "" {
		rulesDataList["regularRules"] = convert.ConvertToInterfaceArray(bypassRegularRules)
	}

	if d.HasChange("config") {
		update = true
	}
	cnRegions, _ := jsonpath.Get("$[0].cn_regions", d.Get("config"))
	if cnRegions != nil && cnRegions != "" {
		rulesDataList["cnRegionList"] = cnRegions
	}

	RulesMap := make([]map[string]interface{}, 0)
	RulesMap = append(RulesMap, rulesDataList)
	request["Rules"], _ = convert.ConvertListMapToJsonString(RulesMap)

	request["InstanceId"] = parts[0]
	request["DefenseType"] = parts[1]

	if update {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("waf-openapi", "2021-10-01", action, query, request, true)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
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
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("waf-openapi", "2021-10-01", action, query, request, true)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
	}

	return resourceWafv3DefenseRuleRead(d, meta)
}

func resourceWafv3DefenseRuleDelete(d *schema.ResourceData, meta interface{}) error {

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
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("waf-openapi", "2021-10-01", action, query, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		if errs.IsExpectedErrors(err, []string{"Defense.Control.DefenseRuleNotExist"}) || errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	return nil
}
