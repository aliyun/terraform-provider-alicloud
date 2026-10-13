package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAliCloudCmsAlertEventIntegrationPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudCmsAlertEventIntegrationPolicyCreate,
		Read:   resourceAliCloudCmsAlertEventIntegrationPolicyRead,
		Update: resourceAliCloudCmsAlertEventIntegrationPolicyUpdate,
		Delete: resourceAliCloudCmsAlertEventIntegrationPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"alert_event_integration_policy_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"workspace": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"SYS_EVENT", "CMS_ALERT", "K8S_EVENT", "CUSTOM",
					"PROMETHEUS", "GRAFANA", "ZABBIX", "SKYWALKING",
					"OPEN_FALCON", "NAGIOS",
				}, false),
			},
			"integration_setting": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"enable": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"filter_setting": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"conditions": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"field": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"value": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"op": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"expression": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"relation": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"transformer_setting": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"filter_setting": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"conditions": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"field": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"value": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"op": {
													Type:     schema.TypeString,
													Optional: true,
												},
											},
										},
									},
									"expression": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"relation": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
								},
							},
						},
						"label_key": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"mapping": {
							Type:     schema.TypeMap,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Optional: true,
							Computed: true,
						},
						"reg_exp": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"source": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"target": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"value": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"variable": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudCmsAlertEventIntegrationPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "/alertEventIntegrationPolicies"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error

	request = make(map[string]interface{})
	query["workspace"] = StringPointer(d.Get("workspace").(string))

	body["alertEventIntegrationPolicyName"] = d.Get("alert_event_integration_policy_name")
	if v, ok := d.GetOk("description"); ok {
		body["description"] = v
	}
	if v, ok := d.GetOk("type"); ok {
		body["type"] = v
	}
	if v, ok := d.GetOk("integration_setting"); ok {
		body["integrationSetting"] = v
	}
	if v, ok := d.GetOk("filter_setting"); ok {
		if fs := expandCmsAlertEventIntegrationPolicyFilterSetting(v.([]interface{})); fs != nil {
			body["filterSetting"] = fs
		}
	}
	if v, ok := d.GetOk("transformer_setting"); ok {
		if ts := expandCmsAlertEventIntegrationPolicyTransformerSetting(v.([]interface{})); ts != nil {
			body["transformerSetting"] = ts
		}
	}
	request = body
	wait := incrementalWait(3*time.Second, 1*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RoaPost("Cms", "2024-03-30", action, query, nil, body, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cms_alert_event_integration_policy", action, AlibabaCloudSdkGoERROR)
	}

	id, _ := jsonpath.Get("$.data", response)
	if id == nil || fmt.Sprint(id) == "" {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cms_alert_event_integration_policy", action, AlibabaCloudSdkGoERROR)
	}
	d.SetId(fmt.Sprint(id))

	// If enable is explicitly set to false, disable the policy after creation.
	if v, ok := d.GetOkExists("enable"); ok && !v.(bool) {
		if err := resourceAliCloudCmsAlertEventIntegrationPolicyToggle(d, meta, false); err != nil {
			return WrapError(err)
		}
	}

	return resourceAliCloudCmsAlertEventIntegrationPolicyRead(d, meta)
}

func resourceAliCloudCmsAlertEventIntegrationPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cmsServiceV2 := CmsServiceV2{client}

	objectRaw, err := cmsServiceV2.DescribeCmsAlertEventIntegrationPolicy(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cms_alert_event_integration_policy DescribeCmsAlertEventIntegrationPolicy Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	policyObj, _ := jsonpath.Get("$.data", objectRaw)
	policy := make(map[string]interface{})
	if policyObj != nil {
		policy = policyObj.(map[string]interface{})
	}

	d.Set("alert_event_integration_policy_name", policy["alertEventIntegrationPolicyName"])
	d.Set("description", policy["description"])
	d.Set("type", policy["type"])
	d.Set("integration_setting", policy["integrationSetting"])
	d.Set("workspace", policy["workspace"])
	d.Set("region_id", policy["regionId"])
	d.Set("create_time", policy["createTime"])
	d.Set("update_time", policy["updateTime"])
	d.Set("user_id", policy["userId"])
	d.Set("enable", policy["enabled"])

	if fsList := flattenCmsAlertEventIntegrationPolicyFilterSetting(policy["filterSetting"]); len(fsList) > 0 {
		if err := d.Set("filter_setting", fsList); err != nil {
			return WrapError(err)
		}
	} else {
		d.Set("filter_setting", []interface{}{})
	}

	if tsList := flattenCmsAlertEventIntegrationPolicyTransformerSetting(policy["transformerSetting"]); len(tsList) > 0 {
		if err := d.Set("transformer_setting", tsList); err != nil {
			return WrapError(err)
		}
	} else {
		d.Set("transformer_setting", []interface{}{})
	}

	return nil
}

func resourceAliCloudCmsAlertEventIntegrationPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	policyId := d.Id()
	action := fmt.Sprintf("/alertEventIntegrationPolicies/%s", policyId)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error

	if v, ok := d.GetOk("workspace"); ok {
		query["workspace"] = StringPointer(v.(string))
	}

	body["alertEventIntegrationPolicyName"] = d.Get("alert_event_integration_policy_name")
	body["description"] = d.Get("description")
	body["type"] = d.Get("type")
	body["integrationSetting"] = d.Get("integration_setting")
	if v, ok := d.GetOk("filter_setting"); ok {
		if fs := expandCmsAlertEventIntegrationPolicyFilterSetting(v.([]interface{})); fs != nil {
			body["filterSetting"] = fs
		}
	}
	if v, ok := d.GetOk("transformer_setting"); ok {
		if ts := expandCmsAlertEventIntegrationPolicyTransformerSetting(v.([]interface{})); ts != nil {
			body["transformerSetting"] = ts
		}
	}
	request = body
	wait := incrementalWait(3*time.Second, 1*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
		response, err = client.RoaPut("Cms", "2024-03-30", action, query, nil, body, true)
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

	// Toggle enable/disable if the enable field changed.
	if d.HasChange("enable") {
		enableVal := d.Get("enable").(bool)
		if err := resourceAliCloudCmsAlertEventIntegrationPolicyToggle(d, meta, enableVal); err != nil {
			return WrapError(err)
		}
	}

	return resourceAliCloudCmsAlertEventIntegrationPolicyRead(d, meta)
}

func resourceAliCloudCmsAlertEventIntegrationPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	policyId := d.Id()
	action := fmt.Sprintf("/alertEventIntegrationPolicies/%s", policyId)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	var err error
	request = make(map[string]interface{})

	if v, ok := d.GetOk("workspace"); ok {
		query["workspace"] = StringPointer(v.(string))
	}

	wait := incrementalWait(3*time.Second, 1*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RoaDelete("Cms", "2024-03-30", action, query, nil, nil, true)
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
		if IsExpectedErrors(err, []string{"ResourceNotFound", "NotFound", "404"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}

func resourceAliCloudCmsAlertEventIntegrationPolicyToggle(d *schema.ResourceData, meta interface{}, enable bool) error {
	client := meta.(*connectivity.AliyunClient)
	policyId := d.Id()
	suffix := "enable"
	if !enable {
		suffix = "disable"
	}
	action := fmt.Sprintf("/alertEventIntegrationPolicies/%s/%s", policyId, suffix)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	var err error
	request = make(map[string]interface{})

	if v, ok := d.GetOk("workspace"); ok {
		query["workspace"] = StringPointer(v.(string))
	}

	wait := incrementalWait(3*time.Second, 1*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
		response, err = client.RoaPut("Cms", "2024-03-30", action, query, nil, nil, true)
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
	return nil
}

func expandCmsAlertEventIntegrationPolicyFilterSetting(list []interface{}) map[string]interface{} {
	if len(list) == 0 || list[0] == nil {
		return nil
	}
	m, ok := list[0].(map[string]interface{})
	if !ok {
		return nil
	}
	result := make(map[string]interface{})
	if conditions, ok := m["conditions"].([]interface{}); ok && len(conditions) > 0 {
		condList := make([]interface{}, 0)
		for _, c := range conditions {
			cm, ok := c.(map[string]interface{})
			if !ok {
				continue
			}
			cond := make(map[string]interface{})
			if v, ok := cm["field"]; ok && v.(string) != "" {
				cond["field"] = v
			}
			if v, ok := cm["value"]; ok && v.(string) != "" {
				cond["value"] = v
			}
			if v, ok := cm["op"]; ok && v.(string) != "" {
				cond["op"] = v
			}
			condList = append(condList, cond)
		}
		if len(condList) > 0 {
			result["conditions"] = condList
		}
	}
	if v, ok := m["expression"]; ok && v.(string) != "" {
		result["expression"] = v
	}
	if v, ok := m["relation"]; ok && v.(string) != "" {
		result["relation"] = v
	}
	return result
}

func expandCmsAlertEventIntegrationPolicyTransformerSetting(list []interface{}) []interface{} {
	if len(list) == 0 {
		return nil
	}
	result := make([]interface{}, 0)
	for _, item := range list {
		m, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		ts := make(map[string]interface{})
		if v, ok := m["label_key"]; ok && v.(string) != "" {
			ts["labelKey"] = v
		}
		if v, ok := m["source"]; ok && v.(string) != "" {
			ts["source"] = v
		}
		if v, ok := m["target"]; ok && v.(string) != "" {
			ts["target"] = v
		}
		if v, ok := m["type"]; ok && v.(string) != "" {
			ts["type"] = v
		}
		if v, ok := m["value"]; ok && v.(string) != "" {
			ts["value"] = v
		}
		if v, ok := m["variable"]; ok && v.(string) != "" {
			ts["variable"] = v
		}
		if v, ok := m["reg_exp"]; ok && v.(string) != "" {
			ts["regExp"] = v
		}
		if v, ok := m["mapping"]; ok && len(v.(map[string]interface{})) > 0 {
			ts["mapping"] = v
		}
		if fsList, ok := m["filter_setting"].([]interface{}); ok {
			if fs := expandCmsAlertEventIntegrationPolicyFilterSetting(fsList); fs != nil {
				ts["filterSetting"] = fs
			}
		}
		result = append(result, ts)
	}
	return result
}

func flattenCmsAlertEventIntegrationPolicyFilterSetting(raw interface{}) []map[string]interface{} {
	result := make([]map[string]interface{}, 0)
	if raw == nil {
		return result
	}
	m, ok := raw.(map[string]interface{})
	if !ok {
		return result
	}
	fsResult := make(map[string]interface{})
	conditions := make([]map[string]interface{}, 0)
	if conds, ok := m["conditions"]; ok {
		for _, c := range convertToInterfaceArray(conds) {
			cm, ok := c.(map[string]interface{})
			if !ok {
				continue
			}
			conditions = append(conditions, map[string]interface{}{
				"field": cm["field"],
				"value": cm["value"],
				"op":    cm["op"],
			})
		}
	}
	fsResult["conditions"] = conditions
	fsResult["expression"] = m["expression"]
	fsResult["relation"] = m["relation"]
	result = append(result, fsResult)
	return result
}

func flattenCmsAlertEventIntegrationPolicyTransformerSetting(raw interface{}) []map[string]interface{} {
	result := make([]map[string]interface{}, 0)
	if raw == nil {
		return result
	}
	arr := convertToInterfaceArray(raw)
	for _, item := range arr {
		m, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		ts := make(map[string]interface{})
		ts["label_key"] = m["labelKey"]
		ts["source"] = m["source"]
		ts["target"] = m["target"]
		ts["type"] = m["type"]
		ts["value"] = m["value"]
		ts["variable"] = m["variable"]
		ts["reg_exp"] = m["regExp"]
		if mapping, ok := m["mapping"]; ok {
			ts["mapping"] = mapping
		}
		if fsList := flattenCmsAlertEventIntegrationPolicyFilterSetting(m["filterSetting"]); len(fsList) > 0 {
			ts["filter_setting"] = fsList
		} else {
			ts["filter_setting"] = []map[string]interface{}{}
		}
		result = append(result, ts)
	}
	return result
}
