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
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAliCloudCmsAlertWebhook() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudCmsAlertWebhookCreate,
		Read:   resourceAliCloudCmsAlertWebhookRead,
		Update: resourceAliCloudCmsAlertWebhookUpdate,
		Delete: resourceAliCloudCmsAlertWebhookDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"alert_webhook_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"url": {
				Type:     schema.TypeString,
				Required: true,
			},
			"content_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"JSON", "FORM"}, false),
			},
			"method": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"GET", "POST"}, false),
			},
			"lang": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"zh_CN", "en_US"}, false),
			},
			"headers": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"workspace": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudCmsAlertWebhookCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "/webhook"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	request["name"] = d.Get("alert_webhook_name")
	request["url"] = d.Get("url")
	if v, ok := d.GetOk("content_type"); ok {
		request["contentType"] = v
	}
	if v, ok := d.GetOk("method"); ok {
		request["method"] = v
	}
	if v, ok := d.GetOk("lang"); ok {
		request["lang"] = v
	}
	if v, ok := d.GetOk("headers"); ok {
		request["headers"] = v
	}
	body = request
	wait := incrementalWait(3*time.Second, 0*time.Second)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cms_alert_webhook", action, AlibabaCloudSdkGoERROR)
	}

	id, _ := jsonpath.Get("$.alertWebhookId", response)
	if id == nil {
		return WrapError(fmt.Errorf("%s failed, response does not contain alertWebhookId: %v", action, response))
	}
	d.SetId(fmt.Sprint(id))

	return resourceAliCloudCmsAlertWebhookRead(d, meta)
}

func resourceAliCloudCmsAlertWebhookRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cmsServiceV2 := CmsServiceV2{client}

	objectRaw, err := cmsServiceV2.DescribeCmsAlertWebhook(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cms_alert_webhook DescribeCmsAlertWebhook Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("alert_webhook_name", objectRaw["name"])
	d.Set("content_type", objectRaw["contentType"])
	d.Set("lang", objectRaw["lang"])
	d.Set("method", objectRaw["method"])
	d.Set("url", objectRaw["url"])
	d.Set("workspace", objectRaw["workspace"])
	d.Set("headers", objectRaw["headers"])

	return nil
}

func resourceAliCloudCmsAlertWebhookUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var body map[string]interface{}
	update := false

	var err error
	action := fmt.Sprintf("/webhook/%s", d.Id())
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})

	if d.HasChange("alert_webhook_name") {
		update = true
	}
	request["name"] = d.Get("alert_webhook_name")

	if d.HasChange("url") {
		update = true
	}
	request["url"] = d.Get("url")

	if d.HasChange("content_type") {
		update = true
	}
	if v := d.Get("content_type").(string); v != "" {
		request["contentType"] = v
	}

	if d.HasChange("method") {
		update = true
	}
	if v := d.Get("method").(string); v != "" {
		request["method"] = v
	}

	if d.HasChange("lang") {
		update = true
	}
	if v := d.Get("lang").(string); v != "" {
		request["lang"] = v
	}

	if d.HasChange("headers") {
		update = true
		request["headers"] = d.Get("headers")
	}

	body = request
	if update {
		wait := incrementalWait(3*time.Second, 0*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RoaPatch("Cms", "2024-03-30", action, query, nil, body, true)
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
	}

	return resourceAliCloudCmsAlertWebhookRead(d, meta)
}

func resourceAliCloudCmsAlertWebhookDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "/webhooks"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	var err error
	request = make(map[string]interface{})

	// DeleteAlertWebhooks is a batch API. The single identifier is wrapped
	// into the webhookIds query parameter (JSON array string).
	webhookIdsJson, jsonErr := json.Marshal([]string{d.Id()})
	if jsonErr != nil {
		return WrapError(jsonErr)
	}
	query["webhookIds"] = StringPointer(string(webhookIdsJson))

	wait := incrementalWait(3*time.Second, 0*time.Second)
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
		if IsExpectedErrors(err, []string{"404", "ResourceNotFound"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
