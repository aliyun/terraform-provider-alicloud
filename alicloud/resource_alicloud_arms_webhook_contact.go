// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
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

func resourceAliCloudArmsWebhookContact() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudArmsWebhookContactCreate,
		Read:   resourceAliCloudArmsWebhookContactRead,
		Delete: resourceAliCloudArmsWebhookContactDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"webhook": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"biz_params": {
							Type:     schema.TypeMap,
							Optional: true,
							ForceNew: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"biz_headers": {
							Type:     schema.TypeMap,
							Optional: true,
							ForceNew: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"recover_body": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"method": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"body": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"url": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},
			"webhook_contact_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudArmsWebhookContactCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateOrUpdateWebhookContact"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	webhookBizHeadersJsonPath, err := jsonpath.Get("$[0].biz_headers", d.Get("webhook"))
	if err == nil {
		if headersJson, err := armsWebhookContactKvMapToJson(webhookBizHeadersJsonPath); err == nil {
			request["BizHeaders"] = headersJson
		}
	}

	webhookBizParamsJsonPath, err := jsonpath.Get("$[0].biz_params", d.Get("webhook"))
	if err == nil {
		if paramsJson, err := armsWebhookContactKvMapToJson(webhookBizParamsJsonPath); err == nil {
			request["BizParams"] = paramsJson
		}
	}

	webhookBodyJsonPath, err := jsonpath.Get("$[0].body", d.Get("webhook"))
	if err == nil {
		request["Body"] = webhookBodyJsonPath
	}

	webhookMethodJsonPath, err := jsonpath.Get("$[0].method", d.Get("webhook"))
	if err == nil {
		request["Method"] = webhookMethodJsonPath
	}

	webhookRecoverBodyJsonPath, err := jsonpath.Get("$[0].recover_body", d.Get("webhook"))
	if err == nil {
		request["RecoverBody"] = webhookRecoverBodyJsonPath
	}

	webhookUrlJsonPath, err := jsonpath.Get("$[0].url", d.Get("webhook"))
	if err == nil {
		request["Url"] = webhookUrlJsonPath
	}

	request["WebhookName"] = d.Get("webhook_contact_name")
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("ARMS", "2019-08-08", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_arms_webhook_contact", action, AlibabaCloudSdkGoERROR)
	}

	id, _ := jsonpath.Get("$.WebhookContact.WebhookId", response)
	d.SetId(fmt.Sprint(id))

	return resourceAliCloudArmsWebhookContactRead(d, meta)
}

func resourceAliCloudArmsWebhookContactRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	armsServiceV2 := ArmsServiceV2{client}

	objectRaw, err := armsServiceV2.DescribeArmsWebhookContact(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_arms_webhook_contact DescribeArmsWebhookContact Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("webhook_contact_name", objectRaw["WebhookName"])

	webhookMaps := make([]map[string]interface{}, 0)
	webhookMap := make(map[string]interface{})
	webhookRaw := make(map[string]interface{})
	if objectRaw["Webhook"] != nil {
		webhookRaw = objectRaw["Webhook"].(map[string]interface{})
	}
	if len(webhookRaw) > 0 {
		webhookMap["biz_headers"] = armsWebhookContactKvMapFromResponse(webhookRaw["BizHeaders"])
		webhookMap["biz_params"] = armsWebhookContactKvMapFromResponse(webhookRaw["BizParams"])
		webhookMap["body"] = webhookRaw["Body"]
		webhookMap["method"] = webhookRaw["Method"]
		webhookMap["recover_body"] = webhookRaw["RecoverBody"]
		webhookMap["url"] = webhookRaw["Url"]

		webhookMaps = append(webhookMaps, webhookMap)
	}
	if err := d.Set("webhook", webhookMaps); err != nil {
		return err
	}

	return nil
}

func resourceAliCloudArmsWebhookContactDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteWebhookContact"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["WebhookId"] = d.Id()

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("ARMS", "2019-08-08", action, query, request, true)
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
		if IsExpectedErrors(err, []string{"404"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}

// The published ARMS API models BizHeaders/BizParams as string params carrying
// a JSON array of single-entry objects (e.g. [{"Content-Type":"application/json"}]);
// the CloudSpec model declares them as plain objects, so the generated pass-through
// would send dotted form keys the API does not accept.
func armsWebhookContactKvMapToJson(raw interface{}) (string, error) {
	m, ok := raw.(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("unexpected kv map type: %T", raw)
	}
	kvList := make([]map[string]interface{}, 0, len(m))
	for k, v := range m {
		kvList = append(kvList, map[string]interface{}{k: v})
	}
	b, err := json.Marshal(kvList)
	return string(b), err
}

// DescribeWebhookContacts may return the fields as a JSON string or as an array
// of single-entry objects; both are normalized to a flat map for TypeMap fields.
func armsWebhookContactKvMapFromResponse(raw interface{}) map[string]interface{} {
	switch v := raw.(type) {
	case string:
		var parsed interface{}
		if err := json.Unmarshal([]byte(v), &parsed); err != nil {
			return map[string]interface{}{}
		}
		return armsWebhookContactKvMapFromResponse(parsed)
	case map[string]interface{}:
		return v
	case []interface{}:
		result := make(map[string]interface{})
		for _, item := range v {
			if m, ok := item.(map[string]interface{}); ok {
				for k, val := range m {
					result[k] = val
				}
			}
		}
		return result
	}
	return map[string]interface{}{}
}
