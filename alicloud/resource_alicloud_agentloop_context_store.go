// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
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

func resourceAliCloudAgentloopContextStore() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudAgentloopContextStoreCreate,
		Read:   resourceAliCloudAgentloopContextStoreRead,
		Update: resourceAliCloudAgentloopContextStoreUpdate,
		Delete: resourceAliCloudAgentloopContextStoreDelete,
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
							ForceNew: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"project": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
									"logstore": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
									"start_time": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
								},
							},
						},
					},
				},
			},
			"context_store_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
		project1, _ := jsonpath.Get("$[0].source[0].project", v)
		if project1 != nil && project1 != "" {
			source["project"] = project1
		}
		logstore1, _ := jsonpath.Get("$[0].source[0].logstore", v)
		if logstore1 != nil && logstore1 != "" {
			source["logstore"] = logstore1
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

	d.Set("context_type", objectRaw["contextType"])
	d.Set("create_time", objectRaw["createTime"])
	d.Set("description", objectRaw["description"])
	d.Set("region_id", objectRaw["regionId"])
	d.Set("status", objectRaw["status"])
	d.Set("update_time", objectRaw["updateTime"])
	d.Set("agent_space", objectRaw["agentSpace"])
	d.Set("context_store_name", objectRaw["contextStoreName"])

	// config is intentionally NOT read back:
	// - serviceNames is immutable for now (documented in the API spec) and the
	//   Update API silently ignores it, so the configured value is always truth.
	// - metadataField user keys are consumed at create/update time and the Get
	//   API only returns internal bookkeeping keys (defaultRecallPolicyId,
	//   experienceMetadataVersion, experienceRecallStore, miningProgress).
	// - source is reduced to agentSpace/startTime in the Get response; logstore
	//   and project are never echoed.
	// All config attributes are ForceNew, so any change recreates the resource
	// and the configured values are re-applied. They are listed in
	// ImportStateVerifyIgnore of the acceptance test via the "config" prefix.

	return nil
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
