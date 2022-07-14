package alicloud

import (
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"log"
	"time"
)

func resourceAlicloudMNSQueue() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudMNSQueueCreate,
		Read:   resourceAlicloudMNSQueueRead,
		Update: resourceAlicloudMNSQueueUpdate,
		Delete: resourceAlicloudMNSQueueDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(3, 256),
			},
			"delay_seconds": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      0,
				ValidateFunc: validation.IntBetween(0, 604800),
			},
			"maximum_message_size": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      65536,
				ValidateFunc: validation.IntBetween(1024, 65536),
			},
			"message_retention_period": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      345600,
				ValidateFunc: validation.IntBetween(60, 604800),
			},
			"visibility_timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      30,
				ValidateFunc: validation.IntBetween(1, 43200),
			},
			"polling_wait_seconds": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      0,
				ValidateFunc: validation.IntBetween(0, 1800),
			},
			"logging_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceAlicloudMNSQueueCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateQueue"
	request := make(map[string]interface{})
	conn, err := client.NewMnsClient()
	if err != nil {
		return WrapError(err)
	}

	if v, ok := d.GetOk("name"); ok {
		request["QueueName"] = v
	}
	if v, ok := d.GetOkExists("delay_seconds"); ok {
		request["DelaySeconds"] = v
	}
	if v, ok := d.GetOk("maximum_message_size"); ok {
		request["MaximumMessageSize"] = v
	}
	if v, ok := d.GetOk("message_retention_period"); ok {
		request["MessageRetentionPeriod"] = v
	}
	if v, ok := d.GetOkExists("polling_wait_seconds"); ok {
		request["PollingWaitSeconds"] = v
	}
	if v, ok := d.GetOk("visibility_timeout"); ok {
		request["VisibilityTimeout"] = v
	}
	if v, ok := d.GetOkExists("logging_enabled"); ok {
		request["EnableLogging"] = v
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		resp, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2022-01-19"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		response = resp
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_mns_queue", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(d.Get("name").(string))

	return resourceAlicloudMNSQueueRead(d, meta)
}

func resourceAlicloudMNSQueueRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	mnsService := MnsService{client}

	object, err := mnsService.DescribeMessageServiceQueue(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_message_service_queue mnsOpenService.DescribeMessageServiceQueue Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("name", object["QueueName"])
	d.Set("delay_seconds", object["DelaySeconds"])
	d.Set("maximum_message_size", object["MaximumMessageSize"])
	d.Set("message_retention_period", object["MessageRetentionPeriod"])
	d.Set("visibility_timeout", object["VisibilityTimeout"])
	d.Set("polling_wait_seconds", object["PollingWaitSeconds"])
	d.Set("logging_enabled", object["LoggingEnabled"])

	return nil
}

func resourceAlicloudMNSQueueUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	conn, err := client.NewMnsClient()
	if err != nil {
		return WrapError(err)
	}

	update := false
	request := map[string]interface{}{
		"QueueName": d.Id(),
	}

	if !d.IsNewResource() && d.HasChange("delay_seconds") {
		update = true
		if v, ok := d.GetOkExists("delay_seconds"); ok {
			request["DelaySeconds"] = v
		}
	}

	if !d.IsNewResource() && d.HasChange("maximum_message_size") {
		update = true
		if v, ok := d.GetOk("maximum_message_size"); ok {
			request["MaximumMessageSize"] = v
		}
	}

	if !d.IsNewResource() && d.HasChange("message_retention_period") {
		update = true
		if v, ok := d.GetOk("message_retention_period"); ok {
			request["MessageRetentionPeriod"] = v
		}
	}

	if !d.IsNewResource() && d.HasChange("visibility_timeout") {
		update = true
		if v, ok := d.GetOk("visibility_timeout"); ok {
			request["VisibilityTimeout"] = v
		}
	}

	if !d.IsNewResource() && d.HasChange("polling_wait_seconds") {
		update = true
		if v, ok := d.GetOkExists("polling_wait_seconds"); ok {
			request["PollingWaitSeconds"] = v
		}
	}

	if !d.IsNewResource() && d.HasChange("logging_enabled") {
		update = true
		if v, ok := d.GetOkExists("logging_enabled"); ok {
			request["EnableLogging"] = v
		}
	}

	if update {
		action := "SetQueueAttributes"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			resp, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2022-01-19"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, resp, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}

	return resourceAlicloudMNSQueueRead(d, meta)
}

func resourceAlicloudMNSQueueDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteQueue"
	conn, err := client.NewMnsClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"QueueName": d.Id(),
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		resp, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2022-01-19"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, resp, request)
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
