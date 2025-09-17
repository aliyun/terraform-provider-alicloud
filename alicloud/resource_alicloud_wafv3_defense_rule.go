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
		ratelimit := make(map[string]interface{})
		threshold1, _ := jsonpath.Get("$[0].rate_limit[0].threshold", v)
		if threshold1 != nil && threshold1 != "" {
			ratelimit["threshold"] = threshold1
		}
		status := make(map[string]interface{})
		code1, _ := jsonpath.Get("$[0].rate_limit[0].status[0].code", v)
		if code1 != nil && code1 != "" {
			status["code"] = code1
		}
		count1, _ := jsonpath.Get("$[0].rate_limit[0].status[0].count", v)
		if count1 != nil && count1 != "" && count1.(int) > 0 {
			status["count"] = count1
		}
		ratio1, _ := jsonpath.Get("$[0].rate_limit[0].status[0].ratio", v)
		if ratio1 != nil && ratio1 != "" && ratio1.(int) > 0 {
			status["ratio"] = ratio1
		}

		statusRaw, _ := jsonpath.Get("$[0].rate_limit[0].status", v)
		if statusRaw != nil {
			ratelimit["status"] = status
		}
		interval1, _ := jsonpath.Get("$[0].rate_limit[0].interval", v)
		if interval1 != nil && interval1 != "" && interval1.(int) > 0 {
			ratelimit["interval"] = interval1
		}
		ttl1, _ := jsonpath.Get("$[0].rate_limit[0].ttl", v)
		if ttl1 != nil && ttl1 != "" && ttl1.(int) > 0 {
			ratelimit["ttl"] = ttl1
		}
		target1, _ := jsonpath.Get("$[0].rate_limit[0].target", v)
		if target1 != nil && target1 != "" {
			ratelimit["target"] = target1
		}
		subKey1, _ := jsonpath.Get("$[0].rate_limit[0].sub_key", v)
		if subKey1 != nil && subKey1 != "" {
			ratelimit["subKey"] = subKey1
		}

		if len(ratelimit) > 0 {
			dataList["ratelimit"] = ratelimit
		}
	}

	if v, ok := d.GetOk("config"); ok {
		mode1, _ := jsonpath.Get("$[0].mode", v)
		if mode1 != nil && mode1 != "" {
			dataList["mode"] = mode1
		}
	}

	if v := d.Get("config"); !IsNil(v) {
		if v, ok := d.GetOk("config"); ok {
			localData, err := jsonpath.Get("$[0].conditions", v)
			if err != nil {
				localData = make([]interface{}, 0)
			}
			localMaps := make([]interface{}, 0)
			for _, dataLoop := range localData.(*schema.Set).List() {
				dataLoopTmp := make(map[string]interface{})
				if dataLoop != nil {
					dataLoopTmp = dataLoop.(map[string]interface{})
				}
				dataLoopMap := make(map[string]interface{})
				dataLoopMap["values"] = dataLoopTmp["values"]
				dataLoopMap["opValue"] = dataLoopTmp["op_value"]
				dataLoopMap["subKey"] = dataLoopTmp["sub_key"]
				dataLoopMap["key"] = dataLoopTmp["key"]
				localMaps = append(localMaps, dataLoopMap)
			}
			dataList["conditions"] = localMaps
		}

	}

	if v, ok := d.GetOk("config"); ok {
		ccEffect, _ := jsonpath.Get("$[0].cc_effect", v)
		if ccEffect != nil && ccEffect != "" {
			dataList["effect"] = ccEffect
		}
	}

	if v, ok := d.GetOkExists("rule_status"); ok {
		dataList["status"] = v
	}

	if v := d.Get("config"); !IsNil(v) {
		if v, ok := d.GetOk("config"); ok {
			localData1, err := jsonpath.Get("$[0].account_identifiers", v)
			if err != nil {
				localData1 = make([]interface{}, 0)
			}
			localMaps1 := make([]interface{}, 0)
			for _, dataLoop1 := range localData1.(*schema.Set).List() {
				dataLoop1Tmp := make(map[string]interface{})
				if dataLoop1 != nil {
					dataLoop1Tmp = dataLoop1.(map[string]interface{})
				}
				dataLoop1Map := make(map[string]interface{})
				dataLoop1Map["position"] = dataLoop1Tmp["position"]
				dataLoop1Map["priority"] = dataLoop1Tmp["priority"]
				dataLoop1Map["subKey"] = dataLoop1Tmp["sub_key"]
				dataLoop1Map["key"] = dataLoop1Tmp["key"]
				dataLoop1Map["decodeType"] = dataLoop1Tmp["decode_type"]
				localMaps1 = append(localMaps1, dataLoop1Map)
			}
			dataList["accountIdentifiers"] = localMaps1
		}

	}

	if v, ok := d.GetOk("config"); ok {
		bypassRegularTypes, _ := jsonpath.Get("$[0].bypass_regular_types", v)
		if bypassRegularTypes != nil && bypassRegularTypes != "" {
			dataList["regularTypes"] = bypassRegularTypes.(*schema.Set).List()
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

	if v, ok := d.GetOk("config"); ok {
		bypassTags, _ := jsonpath.Get("$[0].bypass_tags", v)
		if bypassTags != nil && bypassTags != "" {
			dataList["tags"] = bypassTags.(*schema.Set).List()
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
			dataList["remoteAddr"] = remoteAddr1.(*schema.Set).List()
		}
	}

	if v, ok := d.GetOk("config"); ok {
		ruleAction, _ := jsonpath.Get("$[0].rule_action", v)
		if ruleAction != nil && ruleAction != "" {
			dataList["action"] = ruleAction
		}
	}

	if v, ok := d.GetOk("config"); ok {
		bypassRegularRules, _ := jsonpath.Get("$[0].bypass_regular_rules", v)
		if bypassRegularRules != nil && bypassRegularRules != "" {
			dataList["regularRules"] = bypassRegularRules.(*schema.Set).List()
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

	if v, ok := d.GetOkExists("template_id"); ok {
		request["TemplateId"] = v
	}
	if v, ok := d.GetOk("resource"); ok {
		request["Resource"] = v
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

	return resourceAliCloudWafv3DefenseRuleUpdate(d, meta)
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
			for _, accountIdentifiersChildRaw := range accountIdentifiersRaw.([]interface{}) {
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
			regularRulesRaw = configRaw["regularRules"].([]interface{})
		}

		configMap["bypass_regular_rules"] = regularRulesRaw
		regularTypesRaw := make([]interface{}, 0)
		if configRaw["regularTypes"] != nil {
			regularTypesRaw = configRaw["regularTypes"].([]interface{})
		}

		configMap["bypass_regular_types"] = regularTypesRaw
		tagsRaw := make([]interface{}, 0)
		if configRaw["tags"] != nil {
			tagsRaw = configRaw["tags"].([]interface{})
		}

		configMap["bypass_tags"] = tagsRaw
		conditionsRaw := configRaw["conditions"]
		conditionsMaps := make([]map[string]interface{}, 0)
		if conditionsRaw != nil {
			for _, conditionsChildRaw := range conditionsRaw.([]interface{}) {
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
			remoteAddrRaw = configRaw["remoteAddr"].([]interface{})
		}

		configMap["remote_addr"] = remoteAddrRaw
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
	v := d.Get("config")
	ccStatus1, _ := jsonpath.Get("$[0].cc_status", v)
	if ccStatus1 != nil && (d.HasChange("config.0.cc_status") || ccStatus1 != "") {
		dataList["ccStatus"] = ccStatus1
	}

	if d.HasChange("config") {
		update = true
	}
	if v := d.Get("config"); v != nil {
		ratelimit := make(map[string]interface{})
		threshold1, _ := jsonpath.Get("$[0].rate_limit[0].threshold", v)
		if threshold1 != nil && (d.HasChange("config.0.rate_limit.0.threshold") || threshold1 != "") {
			ratelimit["threshold"] = threshold1
		}
		status := make(map[string]interface{})
		code1, _ := jsonpath.Get("$[0].rate_limit[0].status[0].code", v)
		if code1 != nil && (d.HasChange("config.0.rate_limit.0.status.0.code") || code1 != "") {
			status["code"] = code1
		}
		count1, _ := jsonpath.Get("$[0].rate_limit[0].status[0].count", v)
		if count1 != nil && (d.HasChange("config.0.rate_limit.0.status.0.count") || count1 != "") && count1.(int) > 0 {
			status["count"] = count1
		}
		ratio1, _ := jsonpath.Get("$[0].rate_limit[0].status[0].ratio", v)
		if ratio1 != nil && (d.HasChange("config.0.rate_limit.0.status.0.ratio") || ratio1 != "") && ratio1.(int) > 0 {
			status["ratio"] = ratio1
		}
		statusJsonNode, _ := jsonpath.Get("$[0].rate_limit[0].status", v)
		if statusJsonNode != nil {
			ratelimit["status"] = status
		}
		interval1, _ := jsonpath.Get("$[0].rate_limit[0].interval", v)
		if interval1 != nil && (d.HasChange("config.0.rate_limit.0.interval") || interval1 != "") && interval1.(int) > 0 {
			ratelimit["interval"] = interval1
		}
		ttl1, _ := jsonpath.Get("$[0].rate_limit[0].ttl", v)
		if ttl1 != nil && (d.HasChange("config.0.rate_limit.0.ttl") || ttl1 != "") && ttl1.(int) > 0 {
			ratelimit["ttl"] = ttl1
		}
		target1, _ := jsonpath.Get("$[0].rate_limit[0].target", v)
		if target1 != nil && (d.HasChange("config.0.rate_limit.0.target") || target1 != "") {
			ratelimit["target"] = target1
		}
		subKey1, _ := jsonpath.Get("$[0].rate_limit[0].sub_key", v)
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
	mode1, _ := jsonpath.Get("$[0].mode", v)
	if mode1 != nil && (d.HasChange("config.0.mode") || mode1 != "") {
		dataList["mode"] = mode1
	}

	if d.HasChange("config") {
		update = true
	}
	if v := d.Get("config"); v != nil {
		if v, ok := d.GetOk("config"); ok {
			localData, err := jsonpath.Get("$[0].conditions", v)
			if err != nil {
				localData = make([]interface{}, 0)
			}
			localMaps := make([]interface{}, 0)
			for _, dataLoop := range localData.(*schema.Set).List() {
				dataLoopTmp := make(map[string]interface{})
				if dataLoop != nil {
					dataLoopTmp = dataLoop.(map[string]interface{})
				}
				dataLoopMap := make(map[string]interface{})
				dataLoopMap["values"] = dataLoopTmp["values"]
				dataLoopMap["opValue"] = dataLoopTmp["op_value"]
				dataLoopMap["subKey"] = dataLoopTmp["sub_key"]
				dataLoopMap["key"] = dataLoopTmp["key"]
				localMaps = append(localMaps, dataLoopMap)
			}
			dataList["conditions"] = localMaps
		}

	}

	if d.HasChange("config") {
		update = true
	}
	ccEffect, _ := jsonpath.Get("$[0].cc_effect", v)
	if ccEffect != nil && (d.HasChange("config.0.cc_effect") || ccEffect != "") {
		dataList["effect"] = ccEffect
	}

	if d.HasChange("config") {
		update = true
	}
	if v := d.Get("config"); v != nil {
		if v, ok := d.GetOk("config"); ok {
			localData1, err := jsonpath.Get("$[0].account_identifiers", v)
			if err != nil {
				localData1 = make([]interface{}, 0)
			}
			localMaps1 := make([]interface{}, 0)
			for _, dataLoop1 := range localData1.(*schema.Set).List() {
				dataLoop1Tmp := make(map[string]interface{})
				if dataLoop1 != nil {
					dataLoop1Tmp = dataLoop1.(map[string]interface{})
				}
				dataLoop1Map := make(map[string]interface{})
				dataLoop1Map["position"] = dataLoop1Tmp["position"]
				dataLoop1Map["priority"] = dataLoop1Tmp["priority"]
				dataLoop1Map["subKey"] = dataLoop1Tmp["sub_key"]
				dataLoop1Map["key"] = dataLoop1Tmp["key"]
				dataLoop1Map["decodeType"] = dataLoop1Tmp["decode_type"]
				localMaps1 = append(localMaps1, dataLoop1Map)
			}
			dataList["accountIdentifiers"] = localMaps1
		}

	}

	if d.HasChange("config") {
		update = true
	}
	bypassRegularTypes, _ := jsonpath.Get("$[0].bypass_regular_types", d.Get("config"))
	if bypassRegularTypes != nil && (d.HasChange("config.0.bypass_regular_types") || bypassRegularTypes != "") {
		dataList["regularTypes"] = bypassRegularTypes.(*schema.Set).List()
	}

	if d.HasChange("config") {
		update = true
	}
	throttleType, _ := jsonpath.Get("$[0].throttle_type", v)
	if throttleType != nil && (d.HasChange("config.0.throttle_type") || throttleType != "") {
		dataList["type"] = throttleType
	}

	if d.HasChange("config") {
		update = true
	}
	throttleThrehold, _ := jsonpath.Get("$[0].throttle_threhold", v)
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
	url1, _ := jsonpath.Get("$[0].url", v)
	if url1 != nil && (d.HasChange("config.0.url") || url1 != "") {
		dataList["url"] = url1
	}

	if d.HasChange("config") {
		update = true
	}
	protocol1, _ := jsonpath.Get("$[0].protocol", v)
	if protocol1 != nil && (d.HasChange("config.0.protocol") || protocol1 != "") {
		dataList["protocol"] = protocol1
	}

	if d.HasChange("config") {
		update = true
	}
	ua1, _ := jsonpath.Get("$[0].ua", v)
	if ua1 != nil && (d.HasChange("config.0.ua") || ua1 != "") {
		dataList["ua"] = ua1
	}

	dataList["id"] = parts[2]

	if d.HasChange("config") {
		update = true
	}
	bypassTags, _ := jsonpath.Get("$[0].bypass_tags", d.Get("config"))
	if bypassTags != nil && (d.HasChange("config.0.bypass_tags") || bypassTags != "") {
		dataList["tags"] = bypassTags.(*schema.Set).List()
	}

	if d.HasChange("config") {
		update = true
	}
	abroadRegions, _ := jsonpath.Get("$[0].abroad_regions", v)
	if abroadRegions != nil && (d.HasChange("config.0.abroad_regions") || abroadRegions != "") {
		dataList["abroadRegionList"] = abroadRegions
	}

	if d.HasChange("config") {
		update = true
	}
	remoteAddr1, _ := jsonpath.Get("$[0].remote_addr", d.Get("config"))
	if remoteAddr1 != nil && (d.HasChange("config.0.remote_addr") || remoteAddr1 != "") {
		dataList["remoteAddr"] = remoteAddr1.(*schema.Set).List()
	}

	if d.HasChange("config") {
		update = true
	}
	bypassRegularRules, _ := jsonpath.Get("$[0].bypass_regular_rules", d.Get("config"))
	if bypassRegularRules != nil && (d.HasChange("config.0.bypass_regular_rules") || bypassRegularRules != "") {
		dataList["regularRules"] = bypassRegularRules.(*schema.Set).List()
	}

	if d.HasChange("config") {
		update = true
	}
	cnRegions, _ := jsonpath.Get("$[0].cn_regions", v)
	if cnRegions != nil && (d.HasChange("config.0.cn_regions") || cnRegions != "") {
		dataList["cnRegionList"] = cnRegions
	}

	if d.HasChange("config") {
		update = true
	}
	ruleAction, _ := jsonpath.Get("$[0].rule_action", v)
	if ruleAction != nil && (d.HasChange("config.0.rule_action") || ruleAction != "") {
		dataList["action"] = ruleAction
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
	if !d.IsNewResource() && d.HasChange("rule_status") {
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
