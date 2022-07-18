package alicloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"log"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudSmsShortUrl() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudSmsShortUrlCreate,
		Read:   resourceAlicloudSmsShortUrlRead,
		Delete: resourceAlicloudSmsShortUrlDelete,
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

func resourceAlicloudSmsShortUrlCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "AddShortUrl"
	request := make(map[string]interface{})
	conn, err := client.NewDysmsClient()
	if err != nil {
		return WrapError(err)
	}
	request["EffectiveDays"] = d.Get("effective_days")
	request["ShortUrlName"] = d.Get("short_url_name")
	request["SourceUrl"] = d.Get("source_url")
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-05-25"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if IsExpectedErrors(err, []string{"undefined"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_sms_short_url", action, AlibabaCloudSdkGoERROR)
	}
	responseData := response["Data"].(map[string]interface{})
	d.SetId(fmt.Sprint(responseData["ShortUrl"]))

	return resourceAlicloudSmsShortUrlRead(d, meta)
}
func resourceAlicloudSmsShortUrlRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	dysmsapiService := DysmsapiService{client}
	object, err := dysmsapiService.DescribeSmsShortUrl(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_sms_short_url dysmsapiService.DescribeSmsShortUrl Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	n, err := GetDaysBetween2Date("2006-01-02 15:04:05", object["CreateDate"].(string), object["ExpireDate"].(string))
	if err != nil {
		return WrapError(err)
	}
	d.Set("effective_days", n)
	d.Set("short_url_name", object["ShortUrlName"])
	d.Set("source_url", object["SourceUrl"])
	d.Set("status", object["ShortUrlStatus"])
	return nil
}
func resourceAlicloudSmsShortUrlDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "DeleteShortUrl"
	request := make(map[string]interface{})
	conn, err := client.NewDysmsClient()
	if err != nil {
		return WrapError(err)
	}
	request["SourceUrl"] = d.Id()
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-05-25"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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

	return nil
}
