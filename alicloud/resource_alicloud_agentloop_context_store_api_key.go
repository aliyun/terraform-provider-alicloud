// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudAgentloopContextStoreApiKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudAgentloopContextStoreApiKeyCreate,
		Read:   resourceAliCloudAgentloopContextStoreApiKeyRead,
		Delete: resourceAliCloudAgentloopContextStoreApiKeyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"agent_space": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"api_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"context_store_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudAgentloopContextStoreApiKeyCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	agentSpace := d.Get("agent_space")
	contextStoreName := d.Get("context_store_name")
	action := fmt.Sprintf("/agentspace/%s/contextstore/%s/apikey", agentSpace, contextStoreName)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("name"); ok {
		request["name"] = v
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
	// The create response carries the full apiKey (returned only once); mask it
	// before debug logging so the secret never lands in DEBUG=terraform logs.
	debugResponse := response
	if response != nil && response["apiKey"] != nil {
		debugResponse = make(map[string]interface{}, len(response))
		for k, v := range response {
			debugResponse[k] = v
		}
		debugResponse["apiKey"] = "***"
	}
	addDebug(action, debugResponse, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_agentloop_context_store_api_key", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v:%v", agentSpace, contextStoreName, request["name"]))
	// The complete apiKey is only returned at creation time. The Get API only returns a
	// desensitized prefix (e.g. sk-abcd****), so it is persisted into state here and must
	// NOT be refreshed in Read.
	if response["apiKey"] != nil {
		d.Set("api_key", response["apiKey"])
	}

	return resourceAliCloudAgentloopContextStoreApiKeyRead(d, meta)
}

func resourceAliCloudAgentloopContextStoreApiKeyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	agentloopServiceV2 := AgentloopServiceV2{client}

	objectRaw, err := agentloopServiceV2.DescribeAgentloopContextStoreApiKey(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_agentloop_context_store_api_key DescribeAgentloopContextStoreApiKey Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("create_time", objectRaw["createTime"])
	d.Set("region_id", objectRaw["regionId"])
	d.Set("agent_space", objectRaw["agentSpace"])
	d.Set("context_store_name", objectRaw["contextStoreName"])
	d.Set("name", objectRaw["name"])

	return nil
}

func resourceAliCloudAgentloopContextStoreApiKeyDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	agentSpace := parts[0]
	contextStoreName := parts[1]
	name := parts[2]
	action := fmt.Sprintf("/agentspace/%s/contextstore/%s/apikey/%s", agentSpace, contextStoreName, name)
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
		if IsExpectedErrors(err, []string{"NotFound", "ContextStoreNotExist", "AgentSpaceNotExist"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
