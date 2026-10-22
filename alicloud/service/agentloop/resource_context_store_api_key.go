package agentloop

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceContextStoreApiKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceContextStoreApiKeyCreate,
		Read:   resourceContextStoreApiKeyRead,
		Delete: resourceContextStoreApiKeyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: agentloopImportStatePassthrough(3),
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

func resourceContextStoreApiKeyCreate(d *schema.ResourceData, meta interface{}) error {

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
	logs.AddDebug(action, debugResponse, request)

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_agentloop_context_store_api_key", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v:%v", agentSpace, contextStoreName, request["name"]))
	// The complete apiKey is only returned at creation time. The Get API only returns a
	// desensitized prefix (e.g. sk-abcd****), so it is persisted into state here and must
	// NOT be refreshed in Read.
	if response["apiKey"] != nil {
		if err := d.Set("api_key", response["apiKey"]); err != nil {
			return err
		}
	}

	return resourceContextStoreApiKeyRead(d, meta)
}

func resourceContextStoreApiKeyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	agentloopService := NewService(client)

	objectRaw, err := agentloopService.DescribeContextStoreApiKey(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_agentloop_context_store_api_key DescribeContextStoreApiKey Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	if err := d.Set("create_time", objectRaw["createTime"]); err != nil {
		return err
	}
	if err := d.Set("region_id", objectRaw["regionId"]); err != nil {
		return err
	}
	if err := d.Set("agent_space", objectRaw["agentSpace"]); err != nil {
		return err
	}
	if err := d.Set("context_store_name", objectRaw["contextStoreName"]); err != nil {
		return err
	}
	if err := d.Set("name", objectRaw["name"]); err != nil {
		return err
	}

	return nil
}

func resourceContextStoreApiKeyDelete(d *schema.ResourceData, meta interface{}) error {

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

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.RetryContext(context.Background(), d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
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
		if errs.IsExpectedErrors(err, []string{"NotFound", "ContextStoreNotExist", "AgentSpaceNotExist"}) || errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	return nil
}
