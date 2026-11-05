package mscopensubscription

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceSubscription() *schema.Resource {
	return &schema.Resource{
		Create: resourceSubscriptionCreate,
		Read:   resourceSubscriptionRead,
		Update: resourceSubscriptionUpdate,
		Delete: resourceSubscriptionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"channel": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"contact_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"email_status": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntInSlice([]int{-1, -2, 0, 1}),
			},
			"item_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"pmsg_status": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntInSlice([]int{-1, -2, 0, 1}),
			},
			"sms_status": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntInSlice([]int{-1, -2, 0, 1}),
			},
			"tts_status": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntInSlice([]int{-1, -2, 0, 1}),
			},
			"webhook_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"webhook_status": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntInSlice([]int{-1, -2, 0, 1}),
			},
		},
	}
}

func resourceSubscriptionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateSubscriptionItem"
	request := make(map[string]interface{})
	var err error
	request["ItemName"] = d.Get("item_name")
	request["Locale"] = "en"
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("MscOpenSubscription", "2021-07-13", action, nil, request, false)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_msc_sub_subscription", action, errs.AlibabaCloudSdkGoERROR)
	}
	responseSubscriptionItem := make(map[string]interface{}, 0)
	if v, ok := response["SubscriptionItem"].(map[string]interface{}); ok {
		responseSubscriptionItem = v
	}
	if len(responseSubscriptionItem) == 0 {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_msc_sub_subscription", action, fmt.Sprintf("The item name %s does not support subscription.", request["ItemName"]))
	}
	d.SetId(fmt.Sprint(responseSubscriptionItem["ItemId"]))

	return resourceSubscriptionUpdate(d, meta)
}
func resourceSubscriptionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	mscOpenSubscriptionService := NewService(client)
	object, err := mscOpenSubscriptionService.DescribeSubscription(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_msc_sub_subscription mscOpenSubscriptionService.DescribeSubscription Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}
	d.Set("channel", object["Channel"])

	if contactIds, ok := object["ContactIds"]; ok && contactIds != nil {
		contact_idsVal, err := convert.ConvertJsonStringToStringList(contactIds)
		if err != nil {
			return err
		}
		d.Set("contact_ids", contact_idsVal)
	}
	d.Set("description", object["Description"])
	if v, ok := object["EmailStatus"]; ok {
		email_statusVal, err := convert.FormatInt(v)
		if err != nil {
			return err
		}
		d.Set("email_status", email_statusVal)
	}
	d.Set("item_name", object["ItemName"])
	if v, ok := object["PmsgStatus"]; ok {
		pmsg_statusVal, err := convert.FormatInt(v)
		if err != nil {
			return err
		}
		d.Set("pmsg_status", pmsg_statusVal)
	}
	if v, ok := object["SmsStatus"]; ok {
		sms_statusVal, err := convert.FormatInt(v)
		if err != nil {
			return err
		}
		d.Set("sms_status", sms_statusVal)
	}
	if v, ok := object["TtsStatus"]; ok {
		tts_statusVal, err := convert.FormatInt(v)
		if err != nil {
			return err
		}
		d.Set("tts_status", tts_statusVal)
	}

	if webhookIds, ok := object["WebhookIds"]; ok && webhookIds != nil {
		webhook_idsVal, err := convert.ConvertJsonStringToStringList(webhookIds)
		if err != nil {
			return err
		}
		d.Set("webhook_ids", webhook_idsVal)
	}
	if v, ok := object["WebhookStatus"]; ok {
		webhook_statusVal, err := convert.FormatInt(v)
		if err != nil {
			return err
		}
		d.Set("webhook_status", webhook_statusVal)
	}
	return nil
}
func resourceSubscriptionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	var err error
	update := false
	request := map[string]interface{}{
		"ItemId": d.Id(),
	}
	if d.HasChange("email_status") {
		update = true
	}
	if v, ok := d.GetOkExists("email_status"); ok {
		request["EmailStatus"] = v
	}
	if d.HasChange("pmsg_status") {
		update = true
	}
	if v, ok := d.GetOkExists("pmsg_status"); ok {
		request["PmsgStatus"] = v
	}
	if d.HasChange("sms_status") {
		update = true
	}
	if v, ok := d.GetOkExists("sms_status"); ok {
		request["SmsStatus"] = v
	}
	if d.HasChange("tts_status") {
		update = true
	}
	if v, ok := d.GetOkExists("tts_status"); ok {
		request["TtsStatus"] = v
	}
	if d.HasChange("webhook_status") {
		update = true
	}
	if v, ok := d.GetOkExists("webhook_status"); ok {
		request["WebhookStatus"] = v
	}
	if d.HasChange("contact_ids") {
		update = true
		if v, ok := d.GetOk("contact_ids"); ok {
			request["ContactIds"] = convert.ConvertListToJsonString(v.(*schema.Set).List())
		}
	}
	request["Locale"] = "en"
	request["RegionId"] = client.RegionId
	if d.HasChange("webhook_ids") {
		update = true
		if v, ok := d.GetOk("webhook_ids"); ok {
			request["WebhookIds"] = convert.ConvertListToJsonString(v.([]interface{}))
		}
	}

	if update {
		action := "UpdateSubscriptionItem"
		request["ClientToken"] = helper.BuildClientToken("UpdateSubscriptionItem")
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
			return nil
		})
		logs.AddDebug(action, response, request)
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
	}

	return resourceSubscriptionRead(d, meta)
}
func resourceSubscriptionDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resourceSubscription. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}
