// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
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

func resourceAliCloudCdnRealTimeLogDelivery() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudCdnRealTimeLogDeliveryCreate,
		Read:   resourceAliCloudCdnRealTimeLogDeliveryRead,
		Update: resourceAliCloudCdnRealTimeLogDeliveryUpdate,
		Delete: resourceAliCloudCdnRealTimeLogDeliveryDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"domain": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"logstore": {
				Type:     schema.TypeString,
				Required: true,
			},
			"project": {
				Type:     schema.TypeString,
				Required: true,
			},
			"sls_region": {
				Type:     schema.TypeString,
				Required: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudCdnRealTimeLogDeliveryCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateRealTimeLogDelivery"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	conn, err := client.NewCdnClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	if v, ok := d.GetOk("domain"); ok {
		query["Domain"] = v
	}

	if v, ok := d.GetOk("project"); ok {
		query["Project"] = StringPointer(v.(string))
	}

	if v, ok := d.GetOk("logstore"); ok {
		query["Logstore"] = StringPointer(v.(string))
	}

	if v, ok := d.GetOk("sls_region"); ok {
		query["Region"] = StringPointer(v.(string))
	}

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("GET"), StringPointer("2018-05-10"), StringPointer("AK"), query, request, &runtime)
		if err != nil {
			if IsExpectedErrors(err, []string{"InternalError"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cdn_real_time_log_delivery", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(query["Domain"]))

	return resourceAliCloudCdnRealTimeLogDeliveryUpdate(d, meta)
}

func resourceAliCloudCdnRealTimeLogDeliveryRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cdnServiceV2 := CdnServiceV2{client}

	objectRaw, err := cdnServiceV2.DescribeCdnRealTimeLogDelivery(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cdn_real_time_log_delivery DescribeCdnRealTimeLogDelivery Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	if objectRaw["Logstore"] != nil {
		d.Set("logstore", objectRaw["Logstore"])
	}
	if objectRaw["Project"] != nil {
		d.Set("project", objectRaw["Project"])
	}
	if objectRaw["Region"] != nil {
		d.Set("sls_region", objectRaw["Region"])
	}
	if objectRaw["Status"] != nil {
		d.Set("status", objectRaw["Status"])
	}

	d.Set("domain", d.Id())

	return nil
}

func resourceAliCloudCdnRealTimeLogDeliveryUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	if d.HasChange("status") {
		cdnServiceV2 := CdnServiceV2{client}
		object, err := cdnServiceV2.DescribeCdnRealTimeLogDelivery(d.Id())
		if err != nil {
			return WrapError(err)
		}

		target := d.Get("status").(string)
		if object["Status"].(string) != target {
			if target == "online" {
				action := "EnableRealtimeLogDelivery"
				conn, err := client.NewCdnClient()
				if err != nil {
					return WrapError(err)
				}
				request = make(map[string]interface{})
				query = make(map[string]interface{})
				query["Domain"] = d.Id()

				runtime := util.RuntimeOptions{}
				runtime.SetAutoretry(true)
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("GET"), StringPointer("2018-05-10"), StringPointer("AK"), query, request, &runtime)
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
				cdnServiceV2 := CdnServiceV2{client}
				stateConf := BuildStateConf([]string{}, []string{"online"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, cdnServiceV2.CdnRealTimeLogDeliveryStateRefreshFunc(d.Id(), "Status", []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}

			}
			if target == "offline" {
				action := "DisableRealtimeLogDelivery"
				conn, err := client.NewCdnClient()
				if err != nil {
					return WrapError(err)
				}
				request = make(map[string]interface{})
				query = make(map[string]interface{})
				query["Domain"] = d.Id()

				runtime := util.RuntimeOptions{}
				runtime.SetAutoretry(true)
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("GET"), StringPointer("2018-05-10"), StringPointer("AK"), query, request, &runtime)
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
				cdnServiceV2 := CdnServiceV2{client}
				stateConf := BuildStateConf([]string{}, []string{"offline"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, cdnServiceV2.CdnRealTimeLogDeliveryStateRefreshFunc(d.Id(), "Status", []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}

			}
		}
	}

	action := "ModifyRealtimeLogDelivery"
	conn, err := client.NewCdnClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["Domain"] = d.Id()

	if d.HasChange("project") {
		update = true
	}
	if v, ok := d.GetOk("project"); ok {
		query["Project"] = StringPointer(v.(string))
	}

	if d.HasChange("logstore") {
		update = true
	}
	if v, ok := d.GetOk("logstore"); ok {
		query["Logstore"] = StringPointer(v.(string))
	}

	if d.HasChange("sls_region") {
		update = true
	}
	if v, ok := d.GetOk("sls_region"); ok {
		query["Region"] = StringPointer(v.(string))
	}

	if update {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("GET"), StringPointer("2018-05-10"), StringPointer("AK"), query, request, &runtime)
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

	return resourceAliCloudCdnRealTimeLogDeliveryRead(d, meta)
}

func resourceAliCloudCdnRealTimeLogDeliveryDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteRealtimeLogDelivery"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	conn, err := client.NewCdnClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	query["Domain"] = d.Id()

	if v, ok := d.GetOk("project"); ok {
		query["Project"] = StringPointer(v.(string))
	}

	if v, ok := d.GetOk("logstore"); ok {
		query["Logstore"] = StringPointer(v.(string))
	}

	if v, ok := d.GetOk("sls_region"); ok {
		query["Region"] = StringPointer(v.(string))
	}

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("GET"), StringPointer("2018-05-10"), StringPointer("AK"), query, request, &runtime)

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
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
