package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudCmsMetricRuleBlackList() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCmsMetricRuleBlackListCreate,
		Read:   resourceAlicloudCmsMetricRuleBlackListRead,
		Update: resourceAlicloudCmsMetricRuleBlackListUpdate,
		Delete: resourceAlicloudCmsMetricRuleBlackListDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"category": {
				Required: true,
				Type:     schema.TypeString,
			},
			"create_time": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"update_time": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"effective_time": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"enable_end_time": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"enable_start_time": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"instances": {
				Required: true,
				Type:     schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"is_enable": {
				Optional: true,
				Type:     schema.TypeBool,
				Computed: true,
			},
			"metric_rule_black_list_id": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"metric_rule_black_list_name": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"metrics": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"metric_name": {
							Required: true,
							Type:     schema.TypeString,
						},
						"resource": {
							Optional: true,
							Type:     schema.TypeString,
						},
					},
				},
			},
			"namespace": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"scope_type": {
				Optional: true,
				Type:     schema.TypeString,
				Computed: true,
			},
			"scope_value": {
				Optional: true,
				Type:     schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceAlicloudCmsMetricRuleBlackListCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := make(map[string]interface{})
	var err error

	if v, ok := d.GetOk("category"); ok {
		request["Category"] = v
	}
	if v, ok := d.GetOk("effective_time"); ok {
		request["EffectiveTime"] = v
	}
	if v, ok := d.GetOk("enable_end_time"); ok {
		request["EnableEndTime"] = v
	}
	if v, ok := d.GetOk("enable_start_time"); ok {
		request["EnableStartTime"] = v
	}
	if v, ok := d.GetOk("instances"); ok {
		request["Instances"] = v.([]interface{})
	}
	if v, ok := d.GetOk("metric_rule_black_list_name"); ok {
		request["Name"] = v
	}
	if v, ok := d.GetOk("metrics"); ok {
		metricsMaps := make([]map[string]interface{}, 0)
		for _, value0 := range v.(*schema.Set).List() {
			metrics := value0.(map[string]interface{})
			metricsMap := make(map[string]interface{})
			metricsMap["MetricName"] = metrics["metric_name"]
			metricsMap["Resource"] = metrics["resource"]
			metricsMaps = append(metricsMaps, metricsMap)
		}
		request["Metrics"] = metricsMaps
	}
	if v, ok := d.GetOk("namespace"); ok {
		request["Namespace"] = v
	}
	if v, ok := d.GetOk("scope_type"); ok {
		request["ScopeType"] = v
	}
	if v, ok := d.GetOk("scope_value"); ok {
		request["ScopeValue"] = convertListToCommaSeparate(v.([]interface{}))
	}

	var response map[string]interface{}
	action := "CreateMetricRuleBlackList"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		resp, err := client.RpcPost("Cms", "2019-01-01", action, nil, request, false)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		response = resp
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cms_metric_rule_black_list", action, AlibabaCloudSdkGoERROR)
	}

	if v, err := jsonpath.Get("$.Id", response); err != nil || v == nil {
		return WrapErrorf(err, IdMsg, "alicloud_cms_metric_rule_black_list")
	} else {
		d.SetId(fmt.Sprint(v))
	}

	return resourceAlicloudCmsMetricRuleBlackListUpdate(d, meta)
}

func resourceAlicloudCmsMetricRuleBlackListRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cmsService := CmsService{client}

	object, err := cmsService.DescribeCmsMetricRuleBlackList(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cms_metric_rule_black_list cmsService.DescribeCmsMetricRuleBlackList Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("metric_rule_black_list_id", object["Id"])
	d.Set("category", object["Category"])
	d.Set("effective_time", object["EffectiveTime"])
	d.Set("enable_end_time", object["EnableEndTime"])
	d.Set("enable_start_time", object["EnableStartTime"])
	d.Set("create_time", object["CreateTime"])
	d.Set("update_time", object["UpdateTime"])
	d.Set("enable_start_time", object["EnableStartTime"])
	instances, _ := jsonpath.Get("$.Instances", object)
	instancesSet := make([]interface{}, 0)
	for _, valueOfInstance := range instances.([]interface{}) {
		jsonMarshalResult, err := json.Marshal(valueOfInstance)
		if err != nil {
			return WrapError(err)
		}
		instancesSet = append(instancesSet, string(jsonMarshalResult))
	}
	if err := d.Set("instances", instancesSet); err != nil {
		return WrapError(err)
	}
	d.Set("is_enable", object["IsEnable"])
	d.Set("metric_rule_black_list_name", object["Name"])
	metricsMaps := make([]map[string]interface{}, 0)
	metricsRaw := object["Metrics"]
	for _, value0 := range metricsRaw.([]interface{}) {
		metrics := value0.(map[string]interface{})
		metricsMap := make(map[string]interface{})
		metricsMap["metric_name"] = metrics["MetricName"]
		metricsMap["resource"] = metrics["Resource"]
		metricsMaps = append(metricsMaps, metricsMap)
	}
	d.Set("metrics", metricsMaps)
	d.Set("namespace", object["Namespace"])
	d.Set("scope_type", object["ScopeType"])
	scopeValue, _ := jsonpath.Get("$.ScopeValue", object)
	d.Set("scope_value", scopeValue)

	return nil
}

func resourceAlicloudCmsMetricRuleBlackListUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	var err error
	d.Partial(true)
	update := false
	request := map[string]interface{}{
		"Id": d.Id(),
	}

	if d.HasChange("is_enable") {
		update = true
	}
	request["IsEnable"] = d.Get("is_enable")

	if update {
		action := "EnableMetricRuleBlackList"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			resp, err := client.RpcPost("Cms", "2019-01-01", action, nil, request, false)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, resp, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("is_enable")
	}

	update = false
	request = map[string]interface{}{
		"Id": d.Id(),
	}

	if !d.IsNewResource() && d.HasChange("category") {
		update = true
	}
	request["Category"] = d.Get("category")
	if !d.IsNewResource() && d.HasChange("effective_time") {
		update = true
		if v, ok := d.GetOk("effective_time"); ok {
			request["EffectiveTime"] = v
		}
	}
	if !d.IsNewResource() && d.HasChange("enable_end_time") {
		update = true
		if v, ok := d.GetOk("enable_end_time"); ok {
			request["EnableEndTime"] = v
		}
	}
	if !d.IsNewResource() && d.HasChange("enable_start_time") {
		update = true
		if v, ok := d.GetOk("enable_start_time"); ok {
			request["EnableStartTime"] = v
		}
	}
	if !d.IsNewResource() && d.HasChange("instances") {
		update = true
		if v, ok := d.GetOk("instances"); ok {
			request["Instances"] = v.([]interface{})
		}
	}
	if !d.IsNewResource() && d.HasChange("metric_rule_black_list_name") {
		update = true
		if v, ok := d.GetOk("metric_rule_black_list_name"); ok {
			request["Name"] = v
		}
	}

	if !d.IsNewResource() && d.HasChange("metrics") {
		update = true
		if v, ok := d.GetOk("metrics"); ok {
			metricsMaps := make([]map[string]interface{}, 0)
			for _, value0 := range v.(*schema.Set).List() {
				metrics := value0.(map[string]interface{})
				metricsMap := make(map[string]interface{})
				metricsMap["MetricName"] = metrics["metric_name"]
				metricsMap["Resource"] = metrics["resource"]
				metricsMaps = append(metricsMaps, metricsMap)
			}
			request["Metrics"] = metricsMaps
		}
	}
	if !d.IsNewResource() && d.HasChange("namespace") {
		update = true
		if v, ok := d.GetOk("namespace"); ok {
			request["Namespace"] = v
		}
	}
	if !d.IsNewResource() && d.HasChange("scope_type") {
		update = true
		if v, ok := d.GetOk("scope_type"); ok {
			request["ScopeType"] = v
		}
	}
	if !d.IsNewResource() && d.HasChange("scope_value") {
		update = true
		if v, ok := d.GetOk("scope_value"); ok {
			request["ScopeValue"] = convertListToCommaSeparate(v.([]interface{}))
		}
	}

	if update {
		action := "ModifyMetricRuleBlackList"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			resp, err := client.RpcPost("Cms", "2019-01-01", action, nil, request, false)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, resp, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("category")
		d.SetPartial("effective_time")
		d.SetPartial("enable_end_time")
		d.SetPartial("enable_start_time")
		d.SetPartial("instances")
		d.SetPartial("metric_rule_black_list_name")
		d.SetPartial("metrics")
		d.SetPartial("namespace")
		d.SetPartial("scope_type")
		d.SetPartial("scope_value")
	}

	d.Partial(false)
	return resourceAlicloudCmsMetricRuleBlackListRead(d, meta)
}

func resourceAlicloudCmsMetricRuleBlackListDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var err error

	request := map[string]interface{}{

		"Id": d.Id(),
	}

	action := "DeleteMetricRuleBlackList"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		resp, err := client.RpcPost("Cms", "2019-01-01", action, nil, request, false)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, resp, request)
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
