// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// expandAgentloopEvaluatorConfig converts string values that carry JSON
// documents (e.g. config.variables, which the API requires to be a real
// array) into their parsed form before sending them to the backend.
func expandAgentloopEvaluatorConfig(raw interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	rawMap, ok := raw.(map[string]interface{})
	if !ok {
		return result
	}
	for k, v := range rawMap {
		s, ok := v.(string)
		if !ok {
			result[k] = v
			continue
		}
		trimmed := strings.TrimSpace(s)
		if strings.HasPrefix(trimmed, "[") || strings.HasPrefix(trimmed, "{") {
			var parsed interface{}
			if err := json.Unmarshal([]byte(trimmed), &parsed); err == nil {
				result[k] = parsed
				continue
			}
		}
		result[k] = s
	}
	return result
}

// flattenAgentloopEvaluatorConfig converts non-string values returned by the
// API (e.g. arrays/objects) back into JSON strings so they fit the TypeMap
// schema attribute.
func flattenAgentloopEvaluatorConfig(raw interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	rawMap, ok := raw.(map[string]interface{})
	if !ok {
		return result
	}
	for k, v := range rawMap {
		if s, ok := v.(string); ok {
			result[k] = s
			continue
		}
		if bs, err := json.Marshal(v); err == nil {
			result[k] = string(bs)
		}
	}
	return result
}

func resourceAliCloudAgentloopEvaluator() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudAgentloopEvaluatorCreate,
		Read:   resourceAliCloudAgentloopEvaluatorRead,
		Update: resourceAliCloudAgentloopEvaluatorUpdate,
		Delete: resourceAliCloudAgentloopEvaluatorDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"agent_space": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"annotations": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"config": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"created_at": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"current_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"metric_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"properties": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"updated_at": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"version": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"version_description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudAgentloopEvaluatorCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	agentSpace := d.Get("agent_space")
	action := fmt.Sprintf("/api/v1/evaluators/%s", agentSpace)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("name"); ok {
		request["name"] = v
	}

	if v, ok := d.GetOk("annotations"); ok {
		annotationsMapsArray := convertToInterfaceArray(v)

		request["annotations"] = annotationsMapsArray
	}

	if v, ok := d.GetOk("config"); ok {
		request["config"] = expandAgentloopEvaluatorConfig(v)
	}
	if v, ok := d.GetOk("description"); ok {
		request["description"] = v
	}
	request["metricName"] = d.Get("metric_name")
	if v, ok := d.GetOk("version_description"); ok {
		request["versionDescription"] = v
	}
	request["type"] = d.Get("type")
	if v, ok := d.GetOk("display_name"); ok {
		request["displayName"] = v
	}
	if v, ok := d.GetOk("version"); ok {
		request["version"] = v
	}
	if v, ok := d.GetOk("properties"); ok {
		request["properties"] = v
	}
	body = request
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RoaPost("AgentLoop", "2026-05-20", action, query, nil, body, true)
		if err != nil {
			if NeedRetry(err) || IsExpectedErrors(err, []string{"Conflict."}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_agentloop_evaluator", action, AlibabaCloudSdkGoERROR)
	}

	if response["name"] != nil {
		d.SetId(fmt.Sprintf("%v:%v", agentSpace, response["name"]))
	} else {
		d.SetId(fmt.Sprintf("%v:%v", agentSpace, request["name"]))
	}

	return resourceAliCloudAgentloopEvaluatorRead(d, meta)
}

func resourceAliCloudAgentloopEvaluatorRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	agentloopServiceV2 := AgentloopServiceV2{client}

	objectRaw, err := agentloopServiceV2.DescribeAgentloopEvaluator(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_agentloop_evaluator DescribeAgentloopEvaluator Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	evaluatorRawObj, _ := jsonpath.Get("$.evaluator", objectRaw)
	evaluatorRaw := make(map[string]interface{})
	if evaluatorRawMap, ok := evaluatorRawObj.(map[string]interface{}); ok {
		evaluatorRaw = evaluatorRawMap
	}
	d.Set("config", flattenAgentloopEvaluatorConfig(evaluatorRaw["config"]))
	d.Set("created_at", evaluatorRaw["createdAt"])
	d.Set("current_version", evaluatorRaw["currentVersion"])
	d.Set("description", evaluatorRaw["description"])
	d.Set("display_name", evaluatorRaw["displayName"])
	d.Set("metric_name", evaluatorRaw["metricName"])
	d.Set("properties", evaluatorRaw["properties"])
	d.Set("type", evaluatorRaw["type"])
	d.Set("updated_at", evaluatorRaw["updatedAt"])
	d.Set("agent_space", evaluatorRaw["agentSpace"])
	d.Set("name", evaluatorRaw["name"])

	annotationsRaw, _ := jsonpath.Get("$.evaluator.annotations", objectRaw)
	d.Set("annotations", annotationsRaw)

	return nil
}

func resourceAliCloudAgentloopEvaluatorUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var body map[string]interface{}
	update := false
	d.Partial(true)

	var err error
	parts := strings.Split(d.Id(), ":")
	agentSpace := parts[0]
	name := parts[1]
	action := fmt.Sprintf("/api/v1/evaluators/%s/%s", agentSpace, name)
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})

	if d.HasChange("annotations") {
		update = true
	}
	if v, ok := d.GetOk("annotations"); ok || d.HasChange("annotations") {
		annotationsMapsArray := convertToInterfaceArray(v)

		request["annotations"] = annotationsMapsArray
	}

	if d.HasChange("config") {
		update = true
	}
	if v, ok := d.GetOk("config"); ok || d.HasChange("config") {
		request["config"] = expandAgentloopEvaluatorConfig(v)
	}
	if d.HasChange("description") {
		update = true
	}
	if v, ok := d.GetOk("description"); ok || d.HasChange("description") {
		request["description"] = v
	}
	if d.HasChange("display_name") {
		update = true
	}
	if v, ok := d.GetOk("display_name"); ok || d.HasChange("display_name") {
		request["displayName"] = v
	}
	// The backend treats a PUT carrying "version" as "create that version" and
	// rejects it with InvalidOperation.EvaluatorVersionAlreadyExists when the
	// version already exists. Only submit version when it actually changes
	// (i.e. a new version is being created); attribute-only updates must omit
	// it, otherwise any update that keeps the configured version fails.
	if d.HasChange("version") {
		update = true
		if v, ok := d.GetOk("version"); ok {
			request["version"] = v
		}
	}
	if d.HasChange("properties") {
		update = true
	}
	if v, ok := d.GetOk("properties"); ok || d.HasChange("properties") {
		request["properties"] = v
	}
	body = request
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RoaPut("AgentLoop", "2026-05-20", action, query, nil, body, true)
			if err != nil {
				if NeedRetry(err) || IsExpectedErrors(err, []string{"Conflict."}) {
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
	return resourceAliCloudAgentloopEvaluatorRead(d, meta)
}

func resourceAliCloudAgentloopEvaluatorDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	agentSpace := parts[0]
	name := parts[1]
	action := fmt.Sprintf("/api/v1/evaluators/%s/%s", agentSpace, name)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	var err error
	request = make(map[string]interface{})

	// The unversioned DeleteEvaluator removes the whole evaluator including all
	// of its versions (probe-verified). A version-scoped delete only removes that
	// single version and refuses to remove the last remaining one
	// (InvalidOperation.EvaluatorLastVersionDeleteForbidden). Since the Terraform
	// resource models the evaluator as a whole, the version query is NOT set here.

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RoaDelete("AgentLoop", "2026-05-20", action, query, nil, nil, true)
		if err != nil {
			if NeedRetry(err) || IsExpectedErrors(err, []string{"Conflict."}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		// Only genuine NotFound errors are swallowed here. A refusal such as
		// InvalidOperation.EvaluatorLastVersionDeleteForbidden means the backend
		// rejected the deletion; returning nil would drop the state while the
		// remote evaluator still exists.
		if IsExpectedErrors(err, []string{"EvaluatorNotFound", "EvaluatorVersionNotFound"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
