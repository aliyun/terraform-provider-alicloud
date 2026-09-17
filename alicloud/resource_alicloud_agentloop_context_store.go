// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAliCloudAgentloopContextStore() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudAgentloopContextStoreCreate,
		Read:   resourceAliCloudAgentloopContextStoreRead,
		Update: resourceAliCloudAgentloopContextStoreUpdate,
		Delete: resourceAliCloudAgentloopContextStoreDelete,
		Importer: &schema.ResourceImporter{
			State: agentloopImportStatePassthrough(2),
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
			"config": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"metadata_field": {
							Type:     schema.TypeMap,
							Optional: true,
							ForceNew: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"service_names": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
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

func resourceAliCloudAgentloopContextStoreCreate(d *schema.ResourceData, meta interface{}) error {

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

	if v := d.Get("config"); !IsNil(v) {
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_agentloop_context_store", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", agentSpace, request["contextStoreName"]))

	return resourceAliCloudAgentloopContextStoreRead(d, meta)
}

func resourceAliCloudAgentloopContextStoreRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	agentloopServiceV2 := AgentloopServiceV2{client}

	objectRaw, err := agentloopServiceV2.DescribeAgentloopContextStore(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_agentloop_context_store DescribeAgentloopContextStore Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
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
// drift and imports are detected. metadataField is write-only: the Get API
// replaces the user mapping with internal bookkeeping keys
// (defaultRecallPolicyId, experienceMetadataVersion, ...), so the mapping is
// retained from prior state and cannot be recovered on import. Server-injected
// keys such as miningInterval are not managed by this resource.
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

func resourceAliCloudAgentloopContextStoreUpdate(d *schema.ResourceData, meta interface{}) error {
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
	contextStoreName := parts[1]
	action := fmt.Sprintf("/agentspace/%s/contextstore/%s", agentSpace, contextStoreName)
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})

	// NOTE: probe-verified backend behaviour — UpdateContextStore returns 200 for
	// every field but only description actually takes effect; contextType,
	// serviceNames, source and metadataField are silently ignored. Those
	// attributes are therefore marked ForceNew and must never be sent here.
	if d.HasChange("description") {
		update = true
	}
	if v, ok := d.GetOk("description"); ok || d.HasChange("description") {
		request["description"] = v
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
	return resourceAliCloudAgentloopContextStoreRead(d, meta)
}

func resourceAliCloudAgentloopContextStoreDelete(d *schema.ResourceData, meta interface{}) error {

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
		if IsExpectedErrors(err, []string{"ContextStoreNotExist"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
