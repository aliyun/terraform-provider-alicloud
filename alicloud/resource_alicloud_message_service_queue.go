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

func resourceAliCloudMessageServiceQueue() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudMessageServiceQueueCreate,
		Read:   resourceAliCloudMessageServiceQueueRead,
		Update: resourceAliCloudMessageServiceQueueUpdate,
		Delete: resourceAliCloudMessageServiceQueueDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"create_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"delay_seconds": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: IntBetween(0, 604800),
			},
			"logging_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"maximum_message_size": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: IntBetween(1024, 65536),
			},
			"message_retention_period": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: IntBetween(60, 604800),
			},
			"polling_wait_seconds": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: IntBetween(0, 30),
			},
			"queue_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"visibility_timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: IntBetween(1, 43200),
			},
		},
	}
}

func resourceAliCloudMessageServiceQueueCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateQueue"
	var request map[string]interface{}
	var response map[string]interface{}
	conn, err := client.NewMnsClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["QueueName"] = d.Get("queue_name")

	if v, ok := d.GetOk("message_retention_period"); ok {
		request["MessageRetentionPeriod"] = v
	}
	if v, ok := d.GetOk("polling_wait_seconds"); ok {
		request["PollingWaitSeconds"] = v
	}
	if v, ok := d.GetOk("visibility_timeout"); ok {
		request["VisibilityTimeout"] = v
	}
	if v, ok := d.GetOk("delay_seconds"); ok {
		request["DelaySeconds"] = v
	}
	if v, ok := d.GetOk("maximum_message_size"); ok {
		request["MaximumMessageSize"] = v
	}
	if v, ok := d.GetOkExists("logging_enabled"); ok {
		request["EnableLogging"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2022-01-19"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_message_service_queue", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["QueueName"]))

	return resourceAliCloudMessageServiceQueueUpdate(d, meta)
}

func resourceAliCloudMessageServiceQueueRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	messageServiceServiceV2 := MessageServiceServiceV2{client}

	objectRaw, err := messageServiceServiceV2.DescribeMessageServiceQueue(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_message_service_queue DescribeMessageServiceQueue Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("create_time", objectRaw["CreateTime"])
	d.Set("delay_seconds", objectRaw["DelaySeconds"])
	d.Set("logging_enabled", objectRaw["LoggingEnabled"])
	d.Set("maximum_message_size", objectRaw["MaximumMessageSize"])
	d.Set("message_retention_period", objectRaw["MessageRetentionPeriod"])
	d.Set("polling_wait_seconds", objectRaw["PollingWaitSeconds"])
	d.Set("visibility_timeout", objectRaw["VisibilityTimeout"])
	d.Set("queue_name", objectRaw["QueueName"])

	return nil
}

func resourceAliCloudMessageServiceQueueUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	update := false
	action := "SetQueueAttributes"
	conn, err := client.NewMnsClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["QueueName"] = d.Id()
	if !d.IsNewResource() && d.HasChange("message_retention_period") {
		update = true
		request["MessageRetentionPeriod"] = d.Get("message_retention_period")
	}

	if !d.IsNewResource() && d.HasChange("polling_wait_seconds") {
		update = true
		request["PollingWaitSeconds"] = d.Get("polling_wait_seconds")
	}

	if !d.IsNewResource() && d.HasChange("visibility_timeout") {
		update = true
		request["VisibilityTimeout"] = d.Get("visibility_timeout")
	}

	if !d.IsNewResource() && d.HasChange("delay_seconds") {
		update = true
		request["DelaySeconds"] = d.Get("delay_seconds")
	}

	if !d.IsNewResource() && d.HasChange("maximum_message_size") {
		update = true
		request["MaximumMessageSize"] = d.Get("maximum_message_size")
	}

	if !d.IsNewResource() && d.HasChange("logging_enabled") {
		update = true
		request["EnableLogging"] = d.Get("logging_enabled")
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2022-01-19"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})

			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}

	return resourceAliCloudMessageServiceQueueRead(d, meta)
}

func resourceAliCloudMessageServiceQueueDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteQueue"
	var request map[string]interface{}
	var response map[string]interface{}
	conn, err := client.NewMnsClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["QueueName"] = d.Id()

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2022-01-19"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
