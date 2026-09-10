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

func resourceAliCloudDdosCooWebCcRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudDdosCooWebCcRuleCreate,
		Read:   resourceAliCloudDdosCooWebCcRuleRead,
		Update: resourceAliCloudDdosCooWebCcRuleUpdate,
		Delete: resourceAliCloudDdosCooWebCcRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"domain": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"rule_detail": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"condition": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"match_method": {
										Type:     schema.TypeString,
										Required: true,
									},
									"field": {
										Type:     schema.TypeString,
										Required: true,
									},
									"header_name": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"content": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
						"rate_limit": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"target": {
										Type:     schema.TypeString,
										Required: true,
									},
									"ttl": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"sub_key": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"interval": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"threshold": {
										Type:     schema.TypeInt,
										Required: true,
									},
								},
							},
						},
						"action": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: StringInSlice([]string{"block", "challenge", "watch", "accept"}, false),
						},
						"statistics": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"field": {
										Type:     schema.TypeString,
										Required: true,
									},
									"header_name": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"mode": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
						"status_code": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"count_threshold": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"ratio_threshold": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"enabled": {
										Type:     schema.TypeBool,
										Required: true,
									},
									"code": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"use_ratio": {
										Type:     schema.TypeBool,
										Required: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceAliCloudDdosCooWebCcRuleCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "ConfigWebCCRuleV2"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("domain"); ok {
		request["Domain"] = v
	}
	request["RegionId"] = client.RegionId

	ruleListDataList := make(map[string]interface{})

	if v, ok := d.GetOk("rule_detail"); ok {
		action1, _ := jsonpath.Get("$[0].action", v)
		if action1 != nil && action1 != "" {
			ruleListDataList["action"] = action1
		}
	}

	if v := d.Get("rule_detail"); !IsNil(v) {
		ratelimit := make(map[string]interface{})
		threshold1, _ := jsonpath.Get("$[0].rate_limit[0].threshold", d.Get("rule_detail"))
		if threshold1 != nil && threshold1 != "" {
			ratelimit["threshold"] = threshold1
		}
		subKey, _ := jsonpath.Get("$[0].rate_limit[0].sub_key", d.Get("rule_detail"))
		if subKey != nil && subKey != "" {
			ratelimit["subkey"] = subKey
		}
		target1, _ := jsonpath.Get("$[0].rate_limit[0].target", d.Get("rule_detail"))
		if target1 != nil && target1 != "" {
			ratelimit["target"] = target1
		}
		interval1, _ := jsonpath.Get("$[0].rate_limit[0].interval", d.Get("rule_detail"))
		if interval1 != nil && interval1 != "" {
			ratelimit["interval"] = interval1
		}
		ttl1, _ := jsonpath.Get("$[0].rate_limit[0].ttl", d.Get("rule_detail"))
		if ttl1 != nil && ttl1 != "" {
			ratelimit["ttl"] = ttl1
		}

		if len(ratelimit) > 0 {
			ruleListDataList["ratelimit"] = ratelimit
		}
	}

	if v := d.Get("rule_detail"); !IsNil(v) {
		status_code := make(map[string]interface{})
		countThreshold, _ := jsonpath.Get("$[0].status_code[0].count_threshold", d.Get("rule_detail"))
		if countThreshold != nil && countThreshold != "" && !IsEmpty(countThreshold) {
			status_code["count_threshold"] = countThreshold
		}
		ratioThreshold, _ := jsonpath.Get("$[0].status_code[0].ratio_threshold", d.Get("rule_detail"))
		if ratioThreshold != nil && ratioThreshold != "" && !IsEmpty(ratioThreshold) {
			status_code["ratio_threshold"] = ratioThreshold
		}
		useRatio, _ := jsonpath.Get("$[0].status_code[0].use_ratio", d.Get("rule_detail"))
		if useRatio != nil && useRatio != "" {
			status_code["use_ratio"] = useRatio
		}
		enabled1, _ := jsonpath.Get("$[0].status_code[0].enabled", d.Get("rule_detail"))
		if enabled1 != nil && enabled1 != "" {
			status_code["enabled"] = enabled1
		}
		code1, _ := jsonpath.Get("$[0].status_code[0].code", d.Get("rule_detail"))
		if code1 != nil && code1 != "" {
			status_code["code"] = code1
		}

		if len(status_code) > 0 {
			ruleListDataList["status_code"] = status_code
		}
	}

	if v := d.Get("rule_detail"); !IsNil(v) {
		statistics := make(map[string]interface{})
		mode1, _ := jsonpath.Get("$[0].statistics[0].mode", d.Get("rule_detail"))
		if mode1 != nil && mode1 != "" {
			statistics["mode"] = mode1
		}
		headerName, _ := jsonpath.Get("$[0].statistics[0].header_name", d.Get("rule_detail"))
		if headerName != nil && headerName != "" {
			statistics["header_name"] = headerName
		}
		field1, _ := jsonpath.Get("$[0].statistics[0].field", d.Get("rule_detail"))
		if field1 != nil && field1 != "" {
			statistics["field"] = field1
		}

		if len(statistics) > 0 {
			ruleListDataList["statistics"] = statistics
		}
	}

	if v, ok := d.GetOk("name"); ok {
		ruleListDataList["name"] = v
	}

	if v := d.Get("rule_detail"); !IsNil(v) {
		localData, err := jsonpath.Get("$[0].condition", v)
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
			dataLoopMap["content"] = dataLoopTmp["content"]
			dataLoopMap["field"] = dataLoopTmp["field"]
			dataLoopMap["match_method"] = dataLoopTmp["match_method"]
			dataLoopMap["header_name"] = dataLoopTmp["header_name"]
			localMaps = append(localMaps, dataLoopMap)
		}
		ruleListDataList["condition"] = localMaps

	}

	RuleListMap := make([]interface{}, 0)
	RuleListMap = append(RuleListMap, ruleListDataList)
	request["RuleList"] = convertObjectToJsonString(RuleListMap)

	request["Expires"] = "0"

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("ddoscoo", "2020-01-01", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ddoscoo_web_cc_rule", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", request["Domain"], ruleListDataList["name"]))

	return resourceAliCloudDdosCooWebCcRuleRead(d, meta)
}

func resourceAliCloudDdosCooWebCcRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ddosCooServiceV2 := DdosCooServiceV2{client}

	objectRaw, err := ddosCooServiceV2.DescribeDdosCooWebCcRule(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ddoscoo_web_cc_rule DescribeDdosCooWebCcRule Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("name", objectRaw["Name"])

	parts := strings.Split(d.Id(), ":")
	d.Set("domain", parts[0])

	ruleDetailMaps := make([]map[string]interface{}, 0)
	ruleDetailMap := make(map[string]interface{})
	ruleDetailRaw := make(map[string]interface{})
	if objectRaw["RuleDetail"] != nil {
		ruleDetailRaw = objectRaw["RuleDetail"].(map[string]interface{})
	}
	if len(ruleDetailRaw) > 0 {
		ruleDetailMap["action"] = ruleDetailRaw["Action"]

		conditionChildRaw := ruleDetailRaw["Condition"]
		conditionMaps := make([]map[string]interface{}, 0)
		if conditionChildRaw != nil {
			for _, contentRaw := range convertToInterfaceArray(conditionChildRaw) {
				conditionMap := make(map[string]interface{})
				conditionChildRaw := contentRaw.(map[string]interface{})
				conditionMap["content"] = conditionChildRaw["Content"]
				conditionMap["field"] = conditionChildRaw["Field"]
				conditionMap["header_name"] = conditionChildRaw["HeaderName"]
				conditionMap["match_method"] = conditionChildRaw["MatchMethod"]

				conditionMaps = append(conditionMaps, conditionMap)
			}
		}
		ruleDetailMap["condition"] = conditionMaps
		rateLimitMaps := make([]map[string]interface{}, 0)
		rateLimitMap := make(map[string]interface{})
		rateLimitRaw := make(map[string]interface{})
		if ruleDetailRaw["RateLimit"] != nil {
			rateLimitRaw = ruleDetailRaw["RateLimit"].(map[string]interface{})
		}
		if len(rateLimitRaw) > 0 {
			rateLimitMap["interval"] = rateLimitRaw["Interval"]
			rateLimitMap["sub_key"] = rateLimitRaw["SubKey"]
			rateLimitMap["target"] = rateLimitRaw["Target"]
			rateLimitMap["threshold"] = rateLimitRaw["Threshold"]
			rateLimitMap["ttl"] = rateLimitRaw["Ttl"]

			rateLimitMaps = append(rateLimitMaps, rateLimitMap)
		}
		ruleDetailMap["rate_limit"] = rateLimitMaps
		statisticsMaps := make([]map[string]interface{}, 0)
		statisticsMap := make(map[string]interface{})
		statisticsRaw := make(map[string]interface{})
		if ruleDetailRaw["Statistics"] != nil {
			statisticsRaw = ruleDetailRaw["Statistics"].(map[string]interface{})
		}
		if len(statisticsRaw) > 0 {
			statisticsMap["field"] = statisticsRaw["Field"]
			statisticsMap["header_name"] = statisticsRaw["HeaderName"]
			statisticsMap["mode"] = statisticsRaw["Mode"]

			statisticsMaps = append(statisticsMaps, statisticsMap)
		}
		ruleDetailMap["statistics"] = statisticsMaps
		statusCodeMaps := make([]map[string]interface{}, 0)
		statusCodeMap := make(map[string]interface{})
		statusCodeRaw := make(map[string]interface{})
		if ruleDetailRaw["StatusCode"] != nil {
			statusCodeRaw = ruleDetailRaw["StatusCode"].(map[string]interface{})
		}
		if len(statusCodeRaw) > 0 {
			statusCodeMap["code"] = statusCodeRaw["Code"]
			statusCodeMap["count_threshold"] = statusCodeRaw["CountThreshold"]
			statusCodeMap["enabled"] = statusCodeRaw["Enabled"]
			statusCodeMap["ratio_threshold"] = statusCodeRaw["RatioThreshold"]
			statusCodeMap["use_ratio"] = statusCodeRaw["UseRatio"]

			statusCodeMaps = append(statusCodeMaps, statusCodeMap)
		}
		ruleDetailMap["status_code"] = statusCodeMaps
		ruleDetailMaps = append(ruleDetailMaps, ruleDetailMap)
	}
	if err := d.Set("rule_detail", ruleDetailMaps); err != nil {
		return err
	}

	return nil
}

func resourceAliCloudDdosCooWebCcRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	parts := strings.Split(d.Id(), ":")
	action := "ConfigWebCCRuleV2"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["Domain"] = parts[0]
	request["RegionId"] = client.RegionId
	ruleListDataList := make(map[string]interface{})

	if d.HasChange("rule_detail") {
		update = true
	}
	action1, _ := jsonpath.Get("$[0].action", d.Get("rule_detail"))
	if action1 != nil && action1 != "" {
		ruleListDataList["action"] = action1
	}

	ruleListDataList["name"] = parts[1]

	if v := d.Get("rule_detail"); !IsNil(v) {
		ratelimit := make(map[string]interface{})
		threshold1, _ := jsonpath.Get("$[0].rate_limit[0].threshold", d.Get("rule_detail"))
		if threshold1 != nil && threshold1 != "" {
			ratelimit["threshold"] = threshold1
		}
		subKey, _ := jsonpath.Get("$[0].rate_limit[0].sub_key", d.Get("rule_detail"))
		if subKey != nil && subKey != "" {
			ratelimit["subkey"] = subKey
		}
		target1, _ := jsonpath.Get("$[0].rate_limit[0].target", d.Get("rule_detail"))
		if target1 != nil && target1 != "" {
			ratelimit["target"] = target1
		}
		interval1, _ := jsonpath.Get("$[0].rate_limit[0].interval", d.Get("rule_detail"))
		if interval1 != nil && interval1 != "" {
			ratelimit["interval"] = interval1
		}
		ttl1, _ := jsonpath.Get("$[0].rate_limit[0].ttl", d.Get("rule_detail"))
		if ttl1 != nil && ttl1 != "" {
			ratelimit["ttl"] = ttl1
		}

		if len(ratelimit) > 0 {
			ruleListDataList["ratelimit"] = ratelimit
		}
	}

	if v := d.Get("rule_detail"); !IsNil(v) {
		status_code := make(map[string]interface{})
		countThreshold, _ := jsonpath.Get("$[0].status_code[0].count_threshold", d.Get("rule_detail"))
		if countThreshold != nil && countThreshold != "" && !IsEmpty(countThreshold) {
			status_code["count_threshold"] = countThreshold
		}
		ratioThreshold, _ := jsonpath.Get("$[0].status_code[0].ratio_threshold", d.Get("rule_detail"))
		if ratioThreshold != nil && ratioThreshold != "" && !IsEmpty(ratioThreshold) {
			status_code["ratio_threshold"] = ratioThreshold
		}
		useRatio, _ := jsonpath.Get("$[0].status_code[0].use_ratio", d.Get("rule_detail"))
		if useRatio != nil && useRatio != "" {
			status_code["use_ratio"] = useRatio
		}
		enabled1, _ := jsonpath.Get("$[0].status_code[0].enabled", d.Get("rule_detail"))
		if enabled1 != nil && enabled1 != "" {
			status_code["enabled"] = enabled1
		}
		code1, _ := jsonpath.Get("$[0].status_code[0].code", d.Get("rule_detail"))
		if code1 != nil && code1 != "" {
			status_code["code"] = code1
		}

		if len(status_code) > 0 {
			ruleListDataList["status_code"] = status_code
		}
	}

	if v := d.Get("rule_detail"); !IsNil(v) {
		statistics := make(map[string]interface{})
		mode1, _ := jsonpath.Get("$[0].statistics[0].mode", d.Get("rule_detail"))
		if mode1 != nil && mode1 != "" {
			statistics["mode"] = mode1
		}
		headerName, _ := jsonpath.Get("$[0].statistics[0].header_name", d.Get("rule_detail"))
		if headerName != nil && headerName != "" {
			statistics["header_name"] = headerName
		}
		field1, _ := jsonpath.Get("$[0].statistics[0].field", d.Get("rule_detail"))
		if field1 != nil && field1 != "" {
			statistics["field"] = field1
		}

		if len(statistics) > 0 {
			ruleListDataList["statistics"] = statistics
		}
	}

	if v := d.Get("rule_detail"); !IsNil(v) {
		localData, err := jsonpath.Get("$[0].condition", v)
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
			dataLoopMap["content"] = dataLoopTmp["content"]
			dataLoopMap["field"] = dataLoopTmp["field"]
			dataLoopMap["match_method"] = dataLoopTmp["match_method"]
			dataLoopMap["header_name"] = dataLoopTmp["header_name"]
			localMaps = append(localMaps, dataLoopMap)
		}
		ruleListDataList["condition"] = localMaps

	}

	RuleListMap := make([]interface{}, 0)
	RuleListMap = append(RuleListMap, ruleListDataList)
	request["RuleList"] = convertObjectToJsonString(RuleListMap)

	request["Expires"] = "0"

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("ddoscoo", "2020-01-01", action, query, request, true)
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

	return resourceAliCloudDdosCooWebCcRuleRead(d, meta)
}

func resourceAliCloudDdosCooWebCcRuleDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeleteWebCCRuleV2"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["RuleNames"] = convertListToJsonString([]interface{}{parts[1]})
	request["Domain"] = parts[0]
	request["RegionId"] = client.RegionId

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("ddoscoo", "2020-01-01", action, query, request, true)
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
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
