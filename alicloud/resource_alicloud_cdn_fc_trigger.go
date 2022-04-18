package alicloud

import (
	"fmt"
	"log"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudCdnFcTrigger() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCdnFcTriggerCreate,
		Read:   resourceAlicloudCdnFcTriggerRead,
		Update: resourceAlicloudCdnFcTriggerUpdate,
		Delete: resourceAlicloudCdnFcTriggerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"event_meta_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"event_meta_version": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"function_arn": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"notes": {
				Type:     schema.TypeString,
				Required: true,
			},
			"role_arn": {
				Type:     schema.TypeString,
				Required: true,
			},
			"source_arn": {
				Type:     schema.TypeString,
				Required: true,
			},
			"trigger_arn": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudCdnFcTriggerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "AddFCTrigger"
	request := make(map[string]interface{})
	conn, err := client.NewCdnClient()
	if err != nil {
		return WrapError(err)
	}
	request["EventMetaName"] = d.Get("event_meta_name")
	request["EventMetaVersion"] = d.Get("event_meta_version")
	if v, ok := d.GetOk("function_arn"); ok {
		request["FunctionARN"] = v
	}
	request["Notes"] = d.Get("notes")
	request["RoleARN"] = d.Get("role_arn")
	request["SourceARN"] = d.Get("source_arn")
	request["TriggerARN"] = d.Get("trigger_arn")
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-05-10"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cdn_fc_trigger", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["TriggerARN"]))

	return resourceAlicloudCdnFcTriggerRead(d, meta)
}
func resourceAlicloudCdnFcTriggerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cdnService := CdnService{client}
	object, err := cdnService.DescribeCdnFcTrigger(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cdn_fc_trigger cdnService.DescribeCdnFcTrigger Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("trigger_arn", d.Id())
	d.Set("event_meta_name", object["EventMetaName"])
	d.Set("event_meta_version", object["EventMetaVersion"])
	d.Set("notes", object["Notes"])
	d.Set("role_arn", object["RoleARN"])
	d.Set("source_arn", object["SourceArn"])
	return nil
}
func resourceAlicloudCdnFcTriggerUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	conn, err := client.NewCdnClient()
	if err != nil {
		return WrapError(err)
	}
	var response map[string]interface{}
	update := false
	request := map[string]interface{}{
		"TriggerARN": d.Id(),
	}
	if d.HasChange("notes") {
		update = true
		request["Notes"] = d.Get("notes")
	}
	if d.HasChange("role_arn") {
		update = true
		request["RoleARN"] = d.Get("role_arn")
	}
	if d.HasChange("source_arn") {
		update = true
		request["SourceARN"] = d.Get("source_arn")
	}
	if update {
		if v, ok := d.GetOk("function_arn"); ok {
			request["FunctionARN"] = v
		}
		action := "UpdateFCTrigger"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-05-10"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
	return resourceAlicloudCdnFcTriggerRead(d, meta)
}
func resourceAlicloudCdnFcTriggerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteFCTrigger"
	var response map[string]interface{}
	conn, err := client.NewCdnClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"TriggerARN": d.Id(),
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-05-10"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
