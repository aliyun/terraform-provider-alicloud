package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudDcdnWafRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudDcdnWafRuleCreate,
		Read:   resourceAlicloudDcdnWafRuleRead,
		Update: resourceAlicloudDcdnWafRuleUpdate,
		Delete: resourceAlicloudDcdnWafRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"cc_status": {
				Optional:     true,
				Computed:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"on", "off"}, false),
			},
			"cn_region_list": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"conditions": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Required: true,
							Type:     schema.TypeString,
						},
						"op_value": {
							Required: true,
							Type:     schema.TypeString,
						},
						"sub_key": {
							Optional: true,
							Type:     schema.TypeString,
						},
						"values": {
							Optional: true,
							Type:     schema.TypeString,
						},
					},
				},
			},
			"defense_scene": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"effect": {
				Optional:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"rule", "service"}, false),
			},
			"gmt_modified": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"other_region_list": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"policy_id": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"rate_limit": {
				Optional: true,
				MaxItems: 1,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"interval": {
							Optional:     true,
							Type:         schema.TypeInt,
							ValidateFunc: validation.IntBetween(5, 1800),
						},
						"status": {
							Optional: true,
							MaxItems: 1,
							Type:     schema.TypeList,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"code": {
										Optional: true,
										Type:     schema.TypeString,
									},
									"count": {
										Optional:     true,
										Type:         schema.TypeInt,
										ValidateFunc: validation.IntBetween(2, 50000),
									},
									"ratio": {
										Optional:     true,
										Type:         schema.TypeInt,
										ValidateFunc: validation.IntBetween(1, 100),
									},
								},
							},
						},
						"sub_key": {
							Optional: true,
							Type:     schema.TypeString,
						},
						"target": {
							Optional:     true,
							Type:         schema.TypeString,
							ValidateFunc: validation.StringInSlice([]string{"IP", "Header", "Query String Parameter", "Cookie Name", "Session"}, false),
						},
						"threshold": {
							Optional:     true,
							Type:         schema.TypeInt,
							ValidateFunc: validation.IntBetween(2, 500000),
						},
						"ttl": {
							Optional:     true,
							Type:         schema.TypeInt,
							ValidateFunc: validation.IntBetween(60, 86400),
						},
					},
				},
			},
			"regular_rules": {
				Optional: true,
				Type:     schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"regular_types": {
				Optional: true,
				Type:     schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"remote_addr": {
				Optional: true,
				Type:     schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"rule_name": {
				Required: true,
				Type:     schema.TypeString,
			},
			"scenes": {
				Optional: true,
				Type:     schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"status": {
				Optional:     true,
				Computed:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"on", "off"}, false),
			},
			"waf_group_ids": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"action": {
				Optional:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"deny", "js", "monitor"}, false),
			},
		},
	}
}

func resourceAlicloudDcdnWafRuleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := make(map[string]interface{})
	var err error
	ruleConfigMap := make(map[string]interface{}, 0)
	if v, ok := d.GetOk("cc_status"); ok {
		ruleConfigMap["ccStatus"] = v
	}
	if v, ok := d.GetOk("action"); ok {
		ruleConfigMap["action"] = v
	}
	if v, ok := d.GetOk("conditions"); ok {
		conditionsMaps := make([]map[string]interface{}, 0)
		for _, value0 := range v.(*schema.Set).List() {
			conditions := value0.(map[string]interface{})
			conditionsMap := make(map[string]interface{})
			conditionsMap["key"] = conditions["key"]
			conditionsMap["opValue"] = conditions["op_value"]
			if v, ok := conditions["sub_key"]; ok {
				conditionsMap["subKey"] = v
			}
			if v, ok := conditions["values"]; ok {
				conditionsMap["values"] = v
			}
			conditionsMaps = append(conditionsMaps, conditionsMap)
		}
		ruleConfigMap["conditions"] = conditionsMaps
	}
	if v, ok := d.GetOk("rate_limit"); ok {
		rateLimit := v.([]interface{})[0].(map[string]interface{})
		rateLimitMap := make(map[string]interface{})
		rateLimitMap["target"] = rateLimit["target"]
		rateLimitMap["interval"] = rateLimit["interval"]
		rateLimitMap["subKey"] = rateLimit["sub_key"]
		rateLimitMap["threshold"] = formatInt(rateLimit["threshold"])
		rateLimitMap["ttl"] = formatInt(rateLimit["ttl"])

		if v, ok := rateLimit["status"]; ok && len(v.([]interface{})) > 0 {
			rateLimitStatusMap := make(map[string]interface{})
			rateLimitStatus := v.([]interface{})[0].(map[string]interface{})
			rateLimitStatusMap["code"] = rateLimitStatus["code"]
			rateLimitStatusMap["ratio"] = formatInt(rateLimitStatus["ratio"])
			rateLimitStatusMap["count"] = formatInt(rateLimitStatus["count"])
			rateLimitMap["status"] = rateLimitStatusMap
		}

		ruleConfigMap["ratelimit"] = rateLimitMap
	}
	if v, ok := d.GetOk("effect"); ok {
		ruleConfigMap["effect"] = v
	}
	if v, ok := d.GetOk("other_region_list"); ok {
		ruleConfigMap["otherRegionList"] = v
	}
	if v, ok := d.GetOk("policy_id"); ok {
		request["PolicyId"] = formatInt(v)
	}
	if v, ok := d.GetOk("rule_name"); ok {
		ruleConfigMap["name"] = v
	}
	if v, ok := d.GetOk("status"); ok {
		ruleConfigMap["status"] = v
	}
	if v, ok := d.GetOk("waf_group_ids"); ok {
		ruleConfigMap["wafGroupIds"] = v
	}
	if v, ok := d.GetOk("scenes"); ok {
		ruleConfigMap["tags"] = v.([]interface{})
	}
	if v, ok := d.GetOk("regular_rules"); ok {
		ruleConfigMap["regularRules"] = v.([]interface{})
	}
	if v, ok := d.GetOk("regular_types"); ok {
		ruleConfigMap["regularTypes"] = v.([]interface{})
	}
	if v, ok := d.GetOk("remote_addr"); ok {
		ruleConfigMap["remoteAddr"] = v.([]interface{})
	}
	if v, ok := d.GetOk("cn_region_list"); ok {
		ruleConfigMap["cnRegionList"] = v
	}
	if v, err := convertListMapToJsonString([]map[string]interface{}{ruleConfigMap}); err != nil {
		return WrapError(err)
	} else {
		request["RuleConfigs"] = v
	}
	var response map[string]interface{}
	action := "BatchCreateDcdnWafRules"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = client.RpcPost("dcdn", "2018-01-15", action, request, nil, false)
		if err != nil {
			if NeedRetry(err) || IsExpectedErrors(err, []string{"InternalError"}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_dcdn_waf_rule", action, AlibabaCloudSdkGoERROR)
	}

	if v, err := jsonpath.Get("$.RuleIds.RuleId[0]", response); err != nil || v == nil {
		return WrapErrorf(err, IdMsg, "alicloud_dcdn_waf_rule")
	} else {
		d.SetId(fmt.Sprint(v))
	}

	return resourceAlicloudDcdnWafRuleRead(d, meta)
}

func resourceAlicloudDcdnWafRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	dcdnService := DcdnService{client}

	object, err := dcdnService.DescribeDcdnWafRule(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_dcdn_waf_rule dcdnService.DescribeDcdnWafRule Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("defense_scene", object["DefenseScene"])
	d.Set("gmt_modified", object["GmtModified"])
	d.Set("policy_id", object["PolicyId"])
	d.Set("rule_name", object["RuleName"])
	d.Set("status", object["RuleStatus"])

	ruleConfig, err := jsonpath.Get("$.RuleConfig", object)
	if err != nil {
		return WrapError(err)
	}
	ruleConfigObj, err := convertJsonStringToMap(ruleConfig.(string))
	if err != nil {
		return WrapError(err)
	}

	d.Set("action", ruleConfigObj["action"])
	d.Set("cc_status", ruleConfigObj["ccStatus"])
	d.Set("cn_region_list", ruleConfigObj["cnRegionList"])
	if conditions52Raw, ok := ruleConfigObj["conditions"]; ok {
		conditions52Maps := make([]map[string]interface{}, 0)
		for _, value0 := range conditions52Raw.([]interface{}) {
			conditions52 := value0.(map[string]interface{})
			conditions52Map := make(map[string]interface{})
			conditions52Map["key"] = conditions52["key"]
			conditions52Map["op_value"] = conditions52["opValue"]
			conditions52Map["sub_key"] = conditions52["subKey"]
			conditions52Map["values"] = conditions52["values"]
			conditions52Maps = append(conditions52Maps, conditions52Map)
		}
		d.Set("conditions", conditions52Maps)
	}

	if v, ok := ruleConfigObj["rateLimit"]; ok {
		rateLimitMap := make(map[string]interface{}, 0)
		rateLimitObj := v.(map[string]interface{})
		rateLimitMap["ttl"] = rateLimitObj["ttl"]
		rateLimitMap["sub_key"] = rateLimitObj["subKey"]
		rateLimitMap["target"] = rateLimitObj["target"]
		rateLimitMap["interval"] = rateLimitObj["interval"]
		rateLimitMap["threshold"] = rateLimitObj["threshold"]
		if v, ok := rateLimitObj["status"]; ok {
			statusMap := make(map[string]interface{}, 0)
			statusObj := v.(map[string]interface{})
			statusMap["code"] = statusObj["code"]
			statusMap["ratio"] = statusObj["ratio"]
			statusMap["count"] = statusObj["count"]
			rateLimitMap["status"] = []interface{}{statusMap}
		}
		d.Set("rate_limit", []interface{}{rateLimitMap})
	}

	d.Set("effect", ruleConfigObj["effect"])
	d.Set("other_region_list", ruleConfigObj["otherRegionList"])
	d.Set("regular_rules", ruleConfigObj["regularRules"])
	d.Set("regular_types", ruleConfigObj["regularTypes"])
	d.Set("remote_addr", ruleConfigObj["remoteAddr"])
	d.Set("scenes", ruleConfigObj["tags"])
	d.Set("waf_group_ids", ruleConfigObj["wafGroupIds"])

	return nil
}

func resourceAlicloudDcdnWafRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	var err error
	update := false
	request := map[string]interface{}{
		"RuleId": d.Id(),
	}
	if d.HasChange("cc_status") {
		update = true
	}
	if d.HasChange("cn_region_list") {
		update = true
	}
	if d.HasChange("conditions") {
		update = true
	}
	if d.HasChange("effect") {
		update = true
	}
	if d.HasChange("other_region_list") {
		update = true
	}
	if d.HasChange("regular_rules") {
		update = true
	}
	if d.HasChange("regular_types") {
		update = true
	}
	if d.HasChange("remote_addr") {
		update = true
	}
	if d.HasChange("rule_name") {
		update = true
		if v, ok := d.GetOk("rule_name"); ok {
			request["RuleName"] = v
		}
	}
	if d.HasChange("scenes") {
		update = true
	}
	if d.HasChange("status") {
		update = true
	}
	if d.HasChange("waf_group_ids") {
		update = true
	}
	if d.HasChange("rate_limit") {
		update = true
	}
	if d.HasChange("action") {
		update = true
	}

	ruleConfigMap := make(map[string]interface{}, 0)
	if v, ok := d.GetOk("cc_status"); ok {
		ruleConfigMap["ccStatus"] = v
	}
	if v, ok := d.GetOk("action"); ok {
		ruleConfigMap["action"] = v
	}
	if v, ok := d.GetOk("conditions"); ok {
		conditionsMaps := make([]map[string]interface{}, 0)
		for _, value0 := range v.(*schema.Set).List() {
			conditions := value0.(map[string]interface{})
			conditionsMap := make(map[string]interface{})
			conditionsMap["key"] = conditions["key"]
			conditionsMap["opValue"] = conditions["op_value"]
			if v, ok := conditions["sub_key"]; ok {
				conditionsMap["subKey"] = v
			}
			if v, ok := conditions["values"]; ok {
				conditionsMap["values"] = v
			}
			conditionsMaps = append(conditionsMaps, conditionsMap)
		}
		ruleConfigMap["conditions"] = conditionsMaps
	}
	if v, ok := d.GetOk("rate_limit"); ok {
		rateLimit := v.([]interface{})[0].(map[string]interface{})
		rateLimitMap := make(map[string]interface{})
		rateLimitMap["target"] = rateLimit["target"]
		rateLimitMap["interval"] = rateLimit["interval"]
		rateLimitMap["subKey"] = rateLimit["sub_key"]
		rateLimitMap["threshold"] = rateLimit["threshold"]
		rateLimitMap["ttl"] = rateLimit["ttl"]

		if v, ok := rateLimit["status"]; ok && len(v.([]interface{})) > 0 {
			rateLimitStatusMap := make(map[string]interface{})
			rateLimitStatus := v.([]interface{})[0].(map[string]interface{})
			rateLimitStatusMap["code"] = rateLimitStatus["code"]
			rateLimitStatusMap["ratio"] = rateLimitStatus["ratio"]
			rateLimitStatusMap["count"] = rateLimitStatus["count"]
			rateLimitMap["status"] = rateLimitStatusMap
		}

		ruleConfigMap["ratelimit"] = rateLimitMap
	}
	if v, ok := d.GetOk("effect"); ok {
		ruleConfigMap["effect"] = v
	}
	if v, ok := d.GetOk("other_region_list"); ok {
		ruleConfigMap["otherRegionList"] = v
	}
	if v, ok := d.GetOk("status"); ok {
		request["RuleStatus"] = v
	}
	if v, ok := d.GetOk("waf_group_ids"); ok {
		ruleConfigMap["wafGroupIds"] = v
	}
	if v, ok := d.GetOk("scenes"); ok {
		ruleConfigMap["tags"] = v.([]interface{})
	}
	if v, ok := d.GetOk("regular_rules"); ok {
		ruleConfigMap["regularRules"] = v.([]interface{})
	}
	if v, ok := d.GetOk("regular_types"); ok {
		ruleConfigMap["regularTypes"] = v.([]interface{})
	}
	if v, ok := d.GetOk("remote_addr"); ok {
		ruleConfigMap["remoteAddr"] = v.([]interface{})
	}
	if v, ok := d.GetOk("cn_region_list"); ok {
		ruleConfigMap["cnRegionList"] = v
	}

	if v, err := convertMaptoJsonString(ruleConfigMap); err != nil {
		return WrapError(err)
	} else {
		request["RuleConfig"] = v
	}
	if update {
		action := "ModifyDcdnWafRule"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err := client.RpcPost("dcdn", "2018-01-15", action, request, nil, false)
			if err != nil {
				if NeedRetry(err) || IsExpectedErrors(err, []string{"InternalError"}) {
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

	return resourceAlicloudDcdnWafRuleRead(d, meta)
}

func resourceAlicloudDcdnWafRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var err error

	request := map[string]interface{}{
		"RuleIds": d.Id(),
	}

	action := "BatchDeleteDcdnWafRules"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err := client.RpcPost("dcdn", "2018-01-15", action, request, nil, false)
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
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}
