// Package alicloud. This file is hand-written from the Cms 2024-03-30 CloudSpec EvaluationTask definition.
package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAliCloudCmsEvaluationTask() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudCmsEvaluationTaskCreate,
		Read:   resourceAliCloudCmsEvaluationTaskRead,
		Update: resourceAliCloudCmsEvaluationTaskUpdate,
		Delete: resourceAliCloudCmsEvaluationTaskDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"channel": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"config": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsJSON,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					equal, _ := compareJsonTemplateAreEquivalent(old, new)
					return equal
				},
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"data_filter": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"data_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"evaluators": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"config": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsJSON,
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								equal, _ := compareJsonTemplateAreEquivalent(old, new)
								return equal
							},
						},
						"filters": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsJSON,
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								equal, _ := compareJsonTemplateAreEquivalent(old, new)
								return equal
							},
						},
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"result_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"result_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"variable_mapping": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsJSON,
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								equal, _ := compareJsonTemplateAreEquivalent(old, new)
								return equal
							},
						},
					},
				},
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"run_strategies": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Running",
					"Pendding",
					"Completed",
					"Failed",
				}, false),
			},
			"tags": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsJSON,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					equal, _ := compareJsonTemplateAreEquivalent(old, new)
					return equal
				},
			},
			"task_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"task_mode": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"task_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"workspace": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudCmsEvaluationTaskCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	workspace := d.Get("workspace").(string)
	action := fmt.Sprintf("/api/v1/evaluation-task/%s", workspace)
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error

	body["taskName"] = d.Get("task_name")
	if v, ok := d.GetOk("task_mode"); ok {
		body["taskMode"] = v
	}
	if v, ok := d.GetOk("data_type"); ok {
		body["dataType"] = v
	}
	if v, ok := d.GetOk("data_filter"); ok {
		body["dataFilter"] = v
	}
	if v, ok := d.GetOk("channel"); ok {
		body["channel"] = v
	}
	if v, ok := d.GetOk("description"); ok {
		body["description"] = v
	}
	if v, ok := d.GetOk("run_strategies"); ok {
		body["runStrategies"] = v
	}
	if v, ok := d.GetOk("status"); ok {
		body["status"] = v
	}
	if v, ok := d.GetOk("config"); ok {
		configMap := map[string]interface{}{}
		if err := json.Unmarshal([]byte(v.(string)), &configMap); err != nil {
			return WrapError(err)
		}
		body["config"] = configMap
	}
	if v, ok := d.GetOk("tags"); ok {
		tagsMap := map[string]interface{}{}
		if err := json.Unmarshal([]byte(v.(string)), &tagsMap); err != nil {
			return WrapError(err)
		}
		body["tags"] = tagsMap
	}
	if evaluators, ok := d.GetOk("evaluators"); ok {
		body["evaluators"] = expandCmsEvaluationTaskEvaluators(evaluators.([]interface{}))
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
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
	addDebug(action, response, body)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cms_evaluation_task", action, AlibabaCloudSdkGoERROR)
	}

	taskId, _ := response["taskId"].(string)
	if taskId == "" {
		return WrapError(fmt.Errorf("failed to create Cms EvaluationTask: taskId is empty in response"))
	}

	d.SetId(fmt.Sprintf("%s:%s", workspace, taskId))

	return resourceAliCloudCmsEvaluationTaskRead(d, meta)
}

func resourceAliCloudCmsEvaluationTaskRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cmsServiceV2 := CmsServiceV2{client}

	objectRaw, err := cmsServiceV2.DescribeCmsEvaluationTask(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cms_evaluation_task DescribeCmsEvaluationTask Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("channel", objectRaw["channel"])
	d.Set("create_time", objectRaw["createdAt"])
	d.Set("data_filter", objectRaw["dataFilter"])
	d.Set("data_type", objectRaw["dataType"])
	d.Set("description", objectRaw["description"])
	d.Set("region_id", objectRaw["regionId"])
	d.Set("run_strategies", objectRaw["runStrategies"])
	d.Set("status", objectRaw["status"])
	d.Set("task_id", objectRaw["taskId"])
	d.Set("task_mode", objectRaw["taskMode"])
	d.Set("task_name", objectRaw["taskName"])
	d.Set("workspace", objectRaw["workspace"])
	if v, ok := objectRaw["config"]; ok && v != nil {
		if configBytes, err := json.Marshal(v); err == nil {
			d.Set("config", string(configBytes))
		} else {
			return WrapError(err)
		}
	}
	if v, ok := objectRaw["tags"]; ok && v != nil {
		if tagsBytes, err := json.Marshal(v); err == nil {
			d.Set("tags", string(tagsBytes))
		} else {
			return WrapError(err)
		}
	}
	if v, ok := objectRaw["evaluators"]; ok && v != nil {
		d.Set("evaluators", flattenCmsEvaluationTaskEvaluators(v))
	}

	return nil
}

func resourceAliCloudCmsEvaluationTaskUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	parts := strings.Split(d.Id(), ":")
	if len(parts) != 2 {
		return WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", d.Id(), 2, len(parts)))
	}
	workspace := parts[0]
	taskId := parts[1]
	action := fmt.Sprintf("/api/v1/evaluation-task/%s/%s", workspace, taskId)
	query := make(map[string]*string)
	body := make(map[string]interface{})

	body["taskName"] = d.Get("task_name")
	if v, ok := d.GetOk("task_mode"); ok {
		body["taskMode"] = v
	}
	if v, ok := d.GetOk("data_type"); ok {
		body["dataType"] = v
	}
	if v, ok := d.GetOk("data_filter"); ok || d.HasChange("data_filter") {
		body["dataFilter"] = v
	}
	if v, ok := d.GetOk("channel"); ok || d.HasChange("channel") {
		body["channel"] = v
	}
	if v, ok := d.GetOk("description"); ok || d.HasChange("description") {
		body["description"] = v
	}
	if v, ok := d.GetOk("run_strategies"); ok || d.HasChange("run_strategies") {
		body["runStrategies"] = v
	}
	if v, ok := d.GetOk("status"); ok || d.HasChange("status") {
		body["status"] = v
	}
	if v, ok := d.GetOk("config"); ok || d.HasChange("config") {
		configMap := map[string]interface{}{}
		if v != nil && v.(string) != "" {
			if err := json.Unmarshal([]byte(v.(string)), &configMap); err != nil {
				return WrapError(err)
			}
		}
		body["config"] = configMap
	}
	if v, ok := d.GetOk("tags"); ok || d.HasChange("tags") {
		tagsMap := map[string]interface{}{}
		if v != nil && v.(string) != "" {
			if err := json.Unmarshal([]byte(v.(string)), &tagsMap); err != nil {
				return WrapError(err)
			}
		}
		body["tags"] = tagsMap
	}
	if v, ok := d.GetOk("evaluators"); ok || d.HasChange("evaluators") {
		body["evaluators"] = expandCmsEvaluationTaskEvaluators(v.([]interface{}))
	}

	var response map[string]interface{}
	var err error
	wait := incrementalWait(3*time.Second, 5*time.Second)
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
	addDebug(action, response, body)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return resourceAliCloudCmsEvaluationTaskRead(d, meta)
}

func resourceAliCloudCmsEvaluationTaskDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	parts := strings.Split(d.Id(), ":")
	if len(parts) != 2 {
		return WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", d.Id(), 2, len(parts)))
	}
	workspace := parts[0]
	taskId := parts[1]
	action := fmt.Sprintf("/api/v1/evaluation-task/%s/%s", workspace, taskId)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	var err error
	request = make(map[string]interface{})

	wait := incrementalWait(3*time.Second, 5*time.Second)
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
		if IsExpectedErrors(err, []string{"EvaluationTaskNotFound"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}

func expandCmsEvaluationTaskEvaluators(evaluators []interface{}) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(evaluators))
	for _, e := range evaluators {
		if e == nil {
			continue
		}
		raw := e.(map[string]interface{})
		item := map[string]interface{}{
			"name":       raw["name"],
			"resultName": raw["result_name"],
		}
		if v, ok := raw["result_type"]; ok && v.(string) != "" {
			item["resultType"] = v
		}
		if v, ok := raw["config"]; ok && v.(string) != "" {
			configMap := map[string]interface{}{}
			if err := json.Unmarshal([]byte(v.(string)), &configMap); err == nil {
				item["config"] = configMap
			}
		}
		if v, ok := raw["filters"]; ok && v.(string) != "" {
			filtersMap := map[string]interface{}{}
			if err := json.Unmarshal([]byte(v.(string)), &filtersMap); err == nil {
				item["filters"] = filtersMap
			}
		}
		if v, ok := raw["variable_mapping"]; ok && v.(string) != "" {
			variableMappingMap := map[string]interface{}{}
			if err := json.Unmarshal([]byte(v.(string)), &variableMappingMap); err == nil {
				item["variableMapping"] = variableMappingMap
			}
		}
		result = append(result, item)
	}
	return result
}

func flattenCmsEvaluationTaskEvaluators(raw interface{}) []map[string]interface{} {
	result := make([]map[string]interface{}, 0)
	if raw == nil {
		return result
	}
	resp, ok := raw.([]interface{})
	if !ok {
		return result
	}
	for _, v := range resp {
		item := v.(map[string]interface{})
		mapping := map[string]interface{}{
			"name":        item["name"],
			"result_name": item["resultName"],
		}
		if v, ok := item["resultType"]; ok && v != nil {
			mapping["result_type"] = v
		}
		if v, ok := item["config"]; ok && v != nil {
			if configBytes, err := json.Marshal(v); err == nil {
				mapping["config"] = string(configBytes)
			}
		}
		if v, ok := item["filters"]; ok && v != nil {
			if filtersBytes, err := json.Marshal(v); err == nil {
				mapping["filters"] = string(filtersBytes)
			}
		}
		if v, ok := item["variableMapping"]; ok && v != nil {
			if vmBytes, err := json.Marshal(v); err == nil {
				mapping["variable_mapping"] = string(vmBytes)
			}
		}
		result = append(result, mapping)
	}
	return result
}
