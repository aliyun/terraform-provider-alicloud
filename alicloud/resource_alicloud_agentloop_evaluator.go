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
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

// evaluatorConfigStructuredKeys lists the config keys whose values are
// structured (array/object) in the API contract: variables and parameters are
// arrays/objects, outputSchema and faasConfig are objects. Every other key
// (prompt, codeType, ...) carries a plain string and must be sent verbatim,
// even when its content happens to start with "{" or "[" (e.g. a prompt that
// embeds a JSON example); guessing from the value shape would corrupt those
// strings.
var evaluatorConfigStructuredKeys = map[string]bool{
	"variables":    true,
	"outputSchema": true,
	"faasConfig":   true,
	"parameters":   true,
}

// expandAgentloopEvaluatorConfig converts string values that carry JSON
// documents (e.g. config.variables, which the API requires to be a real
// array) into their parsed form before sending them to the backend. Only the
// keys known to be structured in the API contract are expanded.
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
		if evaluatorConfigStructuredKeys[k] && (strings.HasPrefix(trimmed, "[") || strings.HasPrefix(trimmed, "{")) {
			// Decode with UseNumber so integer literals beyond the float64
			// precision range reach the API unrounded.
			if parsed, err := agentloopDecodeJSON(trimmed); err == nil {
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
// schema attribute. When the serialized form is semantically equal to the text
// already held in the state (the user-supplied value), the state text is kept
// verbatim so formatting differences such as key order or whitespace do not
// turn into a perpetual diff, mirroring the Dataset and Pipeline handling.
func flattenAgentloopEvaluatorConfig(raw interface{}, configured map[string]interface{}) map[string]interface{} {
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
			flat := string(bs)
			if orig, ok := configured[k].(string); ok && agentloopJsonStringsEquivalent(orig, flat) {
				flat = orig
			}
			result[k] = flat
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
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"AGENT", "CODE"}, false),
			},
			"updated_at": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"version": {
				// The initial version number is required by CreateEvaluator.
				Type:     schema.TypeString,
				Required: true,
			},
			"version_description": {
				// Computed so that an imported Evaluator keeps the backend's current
				// version description without forcing a diff when the configuration
				// omits this optional attribute.
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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
	request["version"] = d.Get("version")
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
	configuredConfig := make(map[string]interface{})
	if current, ok := d.GetOk("config"); ok {
		configuredConfig = current.(map[string]interface{})
	}
	d.Set("config", flattenAgentloopEvaluatorConfig(evaluatorRaw["config"], configuredConfig))
	d.Set("created_at", evaluatorRaw["createdAt"])
	d.Set("current_version", evaluatorRaw["currentVersion"])
	// The Get API does not echo the configured version back; it reports the
	// current version together with every version's description. Restore both
	// attributes from the response so an imported Evaluator converges on the
	// first plan: without this the state would lack the required version, the
	// next apply would submit it as a new version, and the backend would
	// reject it with InvalidOperation.EvaluatorVersionAlreadyExists (probe
	// verified that a version-creating PUT also moves currentVersion).
	d.Set("version", evaluatorRaw["currentVersion"])
	if versionsRaw, ok := evaluatorRaw["versions"].([]interface{}); ok {
		currentVersion := fmt.Sprint(evaluatorRaw["currentVersion"])
		for _, item := range versionsRaw {
			versionMap, ok := item.(map[string]interface{})
			if !ok {
				continue
			}
			if fmt.Sprint(versionMap["version"]) == currentVersion {
				d.Set("version_description", versionMap["versionDescription"])
				break
			}
		}
	}
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
	// versionDescription is mutable per the UpdateEvaluator contract; it is
	// submitted together with a version change or on its own.
	if d.HasChange("version_description") {
		update = true
	}
	if v, ok := d.GetOk("version_description"); ok || d.HasChange("version_description") {
		request["versionDescription"] = v
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
