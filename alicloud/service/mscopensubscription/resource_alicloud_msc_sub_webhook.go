package mscopensubscription

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceWebhook() *schema.Resource {
	return &schema.Resource{
		Create: resourceWebhookCreate,
		Read:   resourceWebhookRead,
		Update: resourceWebhookUpdate,
		Delete: resourceWebhookDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"server_url": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^https://oapi.dingtalk.com/robot/send\?access_token=[a-zA-Z0-9]+$`), "The serverUrl of the Webhook."),
			},
			"webhook_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[a-zA-Z]{2,12}$`), "The name must be 2 to 12 characters in length, and can contain uppercase and lowercase letters."),
			},
		},
	}
}

func resourceWebhookCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := make(map[string]interface{})
	request["Locale"] = "en"
	request["ServerUrl"] = d.Get("server_url")
	request["WebhookName"] = d.Get("webhook_name")
	request["ClientToken"] = helper.BuildClientToken("CreateWebhook")
	var response map[string]interface{}
	var err error
	action := "CreateWebhook"
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("MscOpenSubscription", "2021-07-13", action, nil, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_msc_sub_webhook", action, errs.AlibabaCloudSdkGoERROR)
	}
	if fmt.Sprint(response["Code"]) != "200" {
		return errs.WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}

	d.SetId(fmt.Sprint(response["WebhookId"]))

	return resourceWebhookRead(d, meta)
}
func resourceWebhookRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	mscOpenSubscriptionService := NewService(client)
	object, err := mscOpenSubscriptionService.DescribeWebhook(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_msc_sub_webhook mscOpenSubscriptionService.DescribeWebhook Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}
	d.Set("server_url", object["ServerUrl"])
	d.Set("webhook_name", object["WebhookName"])
	return nil
}
func resourceWebhookUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	var err error
	update := false
	request := map[string]interface{}{
		"WebhookId": d.Id(),
	}
	if d.HasChange("server_url") {
		update = true
	}
	request["ServerUrl"] = d.Get("server_url")
	if d.HasChange("webhook_name") {
		update = true
	}
	request["WebhookName"] = d.Get("webhook_name")
	if update {
		request["Locale"] = "en"
		action := "UpdateWebhook"
		request["ClientToken"] = helper.BuildClientToken("UpdateWebhook")
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("MscOpenSubscription", "2021-07-13", action, nil, request, true)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			logs.AddDebug(action, response, request)
			return nil
		})
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
	}
	return resourceWebhookRead(d, meta)
}
func resourceWebhookDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	var err error
	request := map[string]interface{}{
		"WebhookId": d.Id(),
		"Locale":    "en",
	}

	action := "DeleteWebhook"
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("MscOpenSubscription", "2021-07-13", action, nil, request, false)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"ResourceNotFound"}) || errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}
	return nil
}
