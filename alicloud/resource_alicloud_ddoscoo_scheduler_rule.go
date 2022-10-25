package alicloud

import (
	"fmt"
	"log"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudDdoscooSchedulerRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudDdoscooSchedulerRuleCreate,
		Read:   resourceAlicloudDdoscooSchedulerRuleRead,
		Update: resourceAlicloudDdoscooSchedulerRuleUpdate,
		Delete: resourceAlicloudDdoscooSchedulerRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"cname": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"param": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"rule_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"rule_type": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntInSlice([]int{2, 3, 6}),
			},
			"rules": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"type": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"A", "CNAME"}, false),
						},
						"value": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"value_type": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntInSlice([]int{1, 2, 3, 6}),
						},
						"priority": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceAlicloudDdoscooSchedulerRuleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateSchedulerRule"
	request := make(map[string]interface{})
	conn, err := client.NewDdoscooClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOk("param"); ok {
		request["Param"] = v
	}

	if v, ok := d.GetOk("rules"); ok {
		ruleMaps := make([]map[string]interface{}, 0)
		for _, rule := range v.(*schema.Set).List() {
			ruleMap := map[string]interface{}{}
			ruleArg := rule.(map[string]interface{})
			ruleMap["Priority"] = ruleArg["priority"]
			ruleMap["RegionId"] = ruleArg["region_id"]
			ruleMap["Type"] = ruleArg["type"]
			ruleMap["Value"] = ruleArg["value"]
			ruleMap["ValueType"] = ruleArg["value_type"]
			ruleMaps = append(ruleMaps, ruleMap)
		}

		rules, err := convertListMapToJsonString(ruleMaps)
		if err != nil {
			return WrapError(err)
		}
		request["Rules"] = rules
	}

	request["RuleName"] = d.Get("rule_name")
	request["RuleType"] = d.Get("rule_type")
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-01-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ddoscoo_scheduler_rule", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["RuleName"]))

	return resourceAlicloudDdoscooSchedulerRuleRead(d, meta)
}
func resourceAlicloudDdoscooSchedulerRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ddoscooService := DdoscooService{client}
	object, err := ddoscooService.DescribeDdoscooSchedulerRule(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ddoscoo_scheduler_rule ddoscooService.DescribeDdoscooSchedulerRule Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("rule_name", d.Id())
	d.Set("cname", object["Cname"])
	d.Set("rule_type", formatInt(object["RuleType"]))
	ruleMaps := make([]map[string]interface{}, 0)
	if v, ok := object["Rules"].([]interface{}); ok {
		for _, rule := range v {
			ruleMap := map[string]interface{}{}
			ruleArg := rule.(map[string]interface{})
			ruleMap["priority"] = formatInt(ruleArg["Priority"])
			ruleMap["region_id"] = ruleArg["RegionId"]
			ruleMap["status"] = formatInt(ruleArg["Status"])
			ruleMap["type"] = ruleArg["Type"]
			ruleMap["value"] = ruleArg["Value"]
			ruleMap["value_type"] = formatInt(ruleArg["ValueType"])
			ruleMaps = append(ruleMaps, ruleMap)
		}
	}
	d.Set("rules", ruleMaps)
	return nil
}
func resourceAlicloudDdoscooSchedulerRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	conn, err := client.NewDdoscooClient()
	if err != nil {
		return WrapError(err)
	}
	var response map[string]interface{}
	request := map[string]interface{}{}
	request["RuleName"] = d.Id()
	update := false
	if d.HasChange("rule_type") {
		update = true
	}
	request["RuleType"] = d.Get("rule_type")
	if d.HasChange("rules") {
		update = true
	}
	if v, ok := d.GetOk("rules"); ok {
		ruleMaps := make([]map[string]interface{}, 0)
		for _, rule := range v.(*schema.Set).List() {
			ruleMap := map[string]interface{}{}
			ruleArg := rule.(map[string]interface{})
			ruleMap["Priority"] = ruleArg["priority"]
			ruleMap["RegionId"] = ruleArg["region_id"]
			ruleMap["Type"] = ruleArg["type"]
			ruleMap["Value"] = ruleArg["value"]
			ruleMap["ValueType"] = ruleArg["value_type"]
			ruleMaps = append(ruleMaps, ruleMap)
		}

		rules, err := convertListMapToJsonString(ruleMaps)
		if err != nil {
			return WrapError(err)
		}
		request["Rules"] = rules
	}

	if d.HasChange("param") {
		update = true
		if v, ok := d.GetOk("param"); ok {
			request["Param"] = v
		}
	}
	if d.HasChange("resource_group_id") {
		update = true
		if v, ok := d.GetOk("resource_group_id"); ok {
			request["ResourceGroupId"] = v
		}
	}
	if update {
		action := "ModifySchedulerRule"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-01-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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

	return resourceAlicloudDdoscooSchedulerRuleRead(d, meta)
}
func resourceAlicloudDdoscooSchedulerRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	conn, err := client.NewDdoscooClient()
	if err != nil {
		return WrapError(err)
	}
	action := "DeleteSchedulerRule"
	var response map[string]interface{}
	request := map[string]interface{}{}
	request["RuleName"] = d.Id()
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-01-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
