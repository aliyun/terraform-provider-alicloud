package sms

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceSmsShortUrl() *schema.Resource {
	return &schema.Resource{
		Create: resourceSmsShortUrlCreate,
		Read:   resourceSmsShortUrlRead,
		Delete: resourceSmsShortUrlDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(2 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"effective_days": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntInSlice([]int{30, 60, 90}),
			},
			"short_url_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"source_url": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceSmsShortUrlCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	var err error
	action := "AddShortUrl"
	request := make(map[string]interface{})
	request["EffectiveDays"] = d.Get("effective_days")
	request["ShortUrlName"] = d.Get("short_url_name")
	request["SourceUrl"] = d.Get("source_url")
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("Dysmsapi", "2017-05-25", action, nil, request, false)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"undefined"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_sms_short_url", action, errs.AlibabaCloudSdkGoERROR)
	}
	responseData := response["Data"].(map[string]interface{})
	d.SetId(fmt.Sprint(responseData["ShortUrl"]))

	return resourceSmsShortUrlRead(d, meta)
}
func resourceSmsShortUrlRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	dysmsapiService := NewService(client)
	object, err := dysmsapiService.DescribeSmsShortUrl(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_sms_short_url dysmsapiService.DescribeSmsShortUrl Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	n, err := helper.GetDaysBetween2Date("2006-01-02 15:04:05", object["CreateDate"].(string), object["ExpireDate"].(string))
	if err != nil {
		return errs.WrapError(err)
	}
	d.Set("effective_days", n)
	d.Set("short_url_name", object["ShortUrlName"])
	d.Set("source_url", object["SourceUrl"])
	d.Set("status", object["ShortUrlStatus"])
	return nil
}
func resourceSmsShortUrlDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	var err error
	action := "DeleteShortUrl"
	request := make(map[string]interface{})
	request["SourceUrl"] = d.Id()
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("Dysmsapi", "2017-05-25", action, nil, request, false)
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

	return nil
}
