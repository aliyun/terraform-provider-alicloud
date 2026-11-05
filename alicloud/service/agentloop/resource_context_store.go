package agentloop

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceContextStore() *schema.Resource {
	return &schema.Resource{
		Create: resourceContextStoreCreate,
		Read:   resourceContextStoreRead,
		Update: resourceContextStoreUpdate,
		Delete: resourceContextStoreDelete,
		Importer: &schema.ResourceImporter{
			StateContext: agentloopImportStatePassthrough(2),
		},
		CustomizeDiff: resourceContextStoreCustomizeDiff,
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
			"config": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"metadata_field": {
							Type:     schema.TypeMap,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"service_names": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						// UpdateContextStore rejects any config.source with
						// 400 InvalidParams "config.source is immutable after
						// store creation", so source stays ForceNew.
						"source": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							ForceNew: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									// The contract only defines agentSpace and startTime here;
									// project/logstore are not accepted (the data source is
									// located via the AgentSpace and serviceNames).
									"agent_space": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
										ForceNew: true,
									},
									"start_time": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
										ForceNew: true,
									},
								},
							},
						},
					},
				},
			},
			"context_store_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[a-z0-9_]+$`), "may only contain lowercase letters, digits, and underscores"),
			},
			// UpdateContextStore accepts contextType but ignores it (the
			// stored type is unchanged), so a type change needs a new store.
			"context_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceContextStoreCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	agentSpace := d.Get("agent_space")
	action := fmt.Sprintf("/agentspace/%s/contextstore", agentSpace)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("context_store_name"); ok {
		request["contextStoreName"] = v
	}

	config := make(map[string]interface{})

	if v := d.Get("config"); !helper.IsNil(v) {
		source := make(map[string]interface{})
		agentSpace1, _ := jsonpath.Get("$[0].source[0].agent_space", v)
		if agentSpace1 != nil && agentSpace1 != "" {
			source["agentSpace"] = agentSpace1
		}
		startTime1, _ := jsonpath.Get("$[0].source[0].start_time", v)
		if startTime1 != nil && startTime1 != "" {
			source["startTime"] = startTime1
		}

		if len(source) > 0 {
			config["source"] = source
		}
		metadataField1, _ := jsonpath.Get("$[0].metadata_field", v)
		if metadataField1 != nil && metadataField1 != "" {
			config["metadataField"] = metadataField1
		}
		serviceNames1, _ := jsonpath.Get("$[0].service_names", v)
		if serviceNames1 != nil {
			if nameList, ok := serviceNames1.([]interface{}); ok && len(nameList) > 0 {
				config["serviceNames"] = nameList
			}
		}

		if len(config) > 0 {
			request["config"] = config
		}
	}

	if v, ok := d.GetOk("description"); ok {
		request["description"] = v
	}
	request["contextType"] = d.Get("context_type")
	body = request
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.RetryContext(context.Background(), d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RoaPost("AgentLoop", "2026-05-20", action, query, nil, body, true)
		if err != nil {
			if helper.NeedRetry(err) || errs.IsExpectedErrors(err, []string{"Conflict."}) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_agentloop_context_store", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", agentSpace, request["contextStoreName"]))

	return resourceContextStoreRead(d, meta)
}

func resourceContextStoreRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	agentloopService := NewService(client)

	objectRaw, err := agentloopService.DescribeContextStore(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_agentloop_context_store DescribeContextStore Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	if err := d.Set("context_type", objectRaw["contextType"]); err != nil {
		return err
	}
	if err := d.Set("create_time", objectRaw["createTime"]); err != nil {
		return err
	}
	if err := d.Set("description", objectRaw["description"]); err != nil {
		return err
	}
	if err := d.Set("region_id", objectRaw["regionId"]); err != nil {
		return err
	}
	if err := d.Set("status", objectRaw["status"]); err != nil {
		return err
	}
	if err := d.Set("update_time", objectRaw["updateTime"]); err != nil {
		return err
	}
	if err := d.Set("agent_space", objectRaw["agentSpace"]); err != nil {
		return err
	}
	if err := d.Set("context_store_name", objectRaw["contextStoreName"]); err != nil {
		return err
	}

	if err := d.Set("config", flattenAgentloopContextStoreConfig(objectRaw["config"], d.Get("config"))); err != nil {
		return err
	}

	return nil
}

// flattenAgentloopContextStoreConfig converts the Get config into state.
// serviceNames and source are echoed verbatim by the API and are read back so
// drift and imports are detected. The Get API mixes the metadataField mapping
// with internal bookkeeping keys (defaultRecallPolicyId,
// experienceMetadataVersion, ...) that cannot be told apart from user keys, so
// the mapping is retained from prior state and cannot be recovered on import.
// Server-injected keys such as miningInterval are not managed by this resource.
func flattenAgentloopContextStoreConfig(raw interface{}, prior interface{}) []map[string]interface{} {
	apiConfig, _ := raw.(map[string]interface{})
	var priorConfig map[string]interface{}
	if l, ok := prior.([]interface{}); ok && len(l) > 0 {
		priorConfig, _ = l[0].(map[string]interface{})
	}

	result := make(map[string]interface{})
	if m, ok := priorConfig["metadata_field"].(map[string]interface{}); ok && len(m) > 0 {
		result["metadata_field"] = m
	}
	if names, ok := apiConfig["serviceNames"].([]interface{}); ok && len(names) > 0 {
		result["service_names"] = names
	}
	if src, ok := apiConfig["source"].(map[string]interface{}); ok && len(src) > 0 {
		var priorSource map[string]interface{}
		if l, ok := priorConfig["source"].([]interface{}); ok && len(l) > 0 {
			priorSource, _ = l[0].(map[string]interface{})
		}
		sourceMap := make(map[string]interface{})
		if v, ok := src["agentSpace"].(string); ok {
			sourceMap["agent_space"] = v
		}
		if v, ok := src["startTime"].(string); ok {
			sourceMap["start_time"] = v
			// Keep the configured spelling when it denotes the same instant
			// (e.g. "Z" versus "+00:00").
			if orig, ok := priorSource["start_time"].(string); ok && orig != v {
				origTime, err1 := time.Parse(time.RFC3339, orig)
				apiTime, err2 := time.Parse(time.RFC3339, v)
				if err1 == nil && err2 == nil && origTime.Equal(apiTime) {
					sourceMap["start_time"] = orig
				}
			}
		}
		if len(sourceMap) > 0 {
			result["source"] = []interface{}{sourceMap}
		}
	}

	if len(result) == 0 {
		return []map[string]interface{}{}
	}
	return []map[string]interface{}{result}
}

// expandAgentloopContextStoreUpdateConfig builds the UpdateContextStore config:
// the complete serviceNames list and metadataField mapping (an empty mapping
// clears the user keys), never source.
func expandAgentloopContextStoreUpdateConfig(v interface{}) map[string]interface{} {
	config := map[string]interface{}{
		"metadataField": map[string]interface{}{},
	}
	l, _ := v.([]interface{})
	if len(l) == 0 {
		return config
	}
	m, _ := l[0].(map[string]interface{})
	if fields, ok := m["metadata_field"].(map[string]interface{}); ok && len(fields) > 0 {
		config["metadataField"] = fields
	}
	if names, ok := m["service_names"].([]interface{}); ok && len(names) > 0 {
		config["serviceNames"] = names
	}
	return config
}

// resourceContextStoreCustomizeDiff forces replacement when service_names is
// cleared or the config block is removed: UpdateContextStore keeps the stored
// serviceNames when they are omitted and rejects an empty list, so neither can
// be applied in place.
func resourceContextStoreCustomizeDiff(_ context.Context, diff *schema.ResourceDiff, _ interface{}) error {
	if diff.Id() == "" || !diff.HasChange("config.0.service_names") || !diff.NewValueKnown("config.0.service_names") {
		return nil
	}
	oldValue, newValue := diff.GetChange("config.0.service_names")
	oldNames, _ := oldValue.([]interface{})
	newNames, _ := newValue.([]interface{})
	if len(oldNames) > 0 && len(newNames) == 0 {
		return diff.ForceNew("config.0.service_names")
	}
	return nil
}

func resourceContextStoreUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var body map[string]interface{}
	update := false

	var err error
	parts := strings.Split(d.Id(), ":")
	agentSpace := parts[0]
	contextStoreName := parts[1]
	action := fmt.Sprintf("/agentspace/%s/contextstore/%s", agentSpace, contextStoreName)
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})

	if d.HasChange("description") {
		update = true
	}
	if v, ok := d.GetOk("description"); ok || d.HasChange("description") {
		request["description"] = v
	}
	// A passed config replaces serviceNames and metadataField in place; source
	// must be omitted because the API rejects it after creation.
	if d.HasChanges("config.0.service_names", "config.0.metadata_field") {
		update = true
		request["config"] = expandAgentloopContextStoreUpdateConfig(d.Get("config"))
	}
	body = request
	if update {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.RetryContext(context.Background(), d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RoaPut("AgentLoop", "2026-05-20", action, query, nil, body, true)
			if err != nil {
				if helper.NeedRetry(err) || errs.IsExpectedErrors(err, []string{"Conflict."}) {
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

	return resourceContextStoreRead(d, meta)
}

func resourceContextStoreDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	agentSpace := parts[0]
	contextStoreName := parts[1]
	action := fmt.Sprintf("/agentspace/%s/contextstore/%s", agentSpace, contextStoreName)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	var err error
	request = make(map[string]interface{})

	ctx, cancel := context.WithTimeout(context.Background(), d.Timeout(schema.TimeoutDelete))
	defer cancel()
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RoaDelete("AgentLoop", "2026-05-20", action, query, nil, nil, true)
		if err != nil {
			if helper.NeedRetry(err) || errs.IsExpectedErrors(err, []string{"Conflict."}) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		if errs.IsExpectedErrors(err, []string{"ContextStoreNotExist"}) || errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	// Delete 返回成功时对象仍可能处于 Deleting，同名替换须等待实际不存在。
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RoaGet("AgentLoop", "2026-05-20", action, query, nil, nil)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"ContextStoreNotExist"}) || errs.NotFoundError(err) {
				return nil
			}
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		// 任意成功回读均表示仍存在，不依赖可能扩展的后端状态枚举。
		wait()
		return retry.RetryableError(fmt.Errorf("context store %s still exists", d.Id()))
	})
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}
	return nil
}
