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
			"dlq_policy": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"max_receive_count": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"dead_letter_target_queue": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
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
			"tags": tagsSchema(),
			"visibility_timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: IntBetween(0, 43200),
			},
		},
	}
}

func resourceAliCloudMessageServiceQueueCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateQueue"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["QueueName"] = d.Get("queue_name")

	if v, ok := d.GetOkExists("message_retention_period"); ok {
		request["MessageRetentionPeriod"] = v
	}
	if v, ok := d.GetOkExists("polling_wait_seconds"); ok {
		request["PollingWaitSeconds"] = v
	}
	if v, ok := d.GetOkExists("visibility_timeout"); ok {
		request["VisibilityTimeout"] = v
	}
	if v, ok := d.GetOkExists("delay_seconds"); ok {
		request["DelaySeconds"] = v
	}
	if v, ok := d.GetOkExists("maximum_message_size"); ok {
		request["MaximumMessageSize"] = v
	}
	if v, ok := d.GetOkExists("logging_enabled"); ok {
		request["EnableLogging"] = v
	}
	if v, ok := d.GetOk("tags"); ok {
		tagsMap := ConvertTags(v.(map[string]interface{}))
		request = expandTagsToMap(request, tagsMap)
	}

	objectDataLocalMap := make(map[string]interface{})

	if v := d.Get("dlq_policy"); !IsNil(v) {
		enabled1, _ := jsonpath.Get("$[0].enabled", v)
		if enabled1 != nil && enabled1 != "" {
			objectDataLocalMap["Enabled"] = enabled1
		}

		deadLetterTargetQueue1, _ := jsonpath.Get("$[0].dead_letter_target_queue", v)
		if deadLetterTargetQueue1 != nil && deadLetterTargetQueue1 != "" {
			objectDataLocalMap["DeadLetterTargetQueue"] = deadLetterTargetQueue1
		}
		maxReceiveCount1, _ := jsonpath.Get("$[0].max_receive_count", v)
		if maxReceiveCount1 != nil && maxReceiveCount1 != "" {
			objectDataLocalMap["MaxReceiveCount"] = maxReceiveCount1
		}

		objectDataLocalMapJson, err := json.Marshal(objectDataLocalMap)
		if err != nil {
			return WrapError(err)
		}
		request["DlqPolicy"] = string(objectDataLocalMapJson)
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Mns-open", "2022-01-19", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_message_service_queue", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["QueueName"]))

	return resourceAliCloudMessageServiceQueueRead(d, meta)
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

	dlqPolicyMaps := make([]map[string]interface{}, 0)
	dlqPolicyMap := make(map[string]interface{})
	dlqPolicyRaw := make(map[string]interface{})
	if objectRaw["DlqPolicy"] != nil {
		dlqPolicyRaw = objectRaw["DlqPolicy"].(map[string]interface{})
	}
	if len(dlqPolicyRaw) > 0 {
		dlqPolicyMap["dead_letter_target_queue"] = dlqPolicyRaw["DeadLetterTargetQueue"]
		dlqPolicyMap["enabled"] = dlqPolicyRaw["Enabled"]
		dlqPolicyMap["max_receive_count"] = formatInt(dlqPolicyRaw["MaxReceiveCount"])

		dlqPolicyMaps = append(dlqPolicyMaps, dlqPolicyMap)
	}
	if err := d.Set("dlq_policy", dlqPolicyMaps); err != nil {
		return err
	}
	tagsMaps := objectRaw["Tags"]
	d.Set("tags", tagsToMap(tagsMaps))

	return nil
}

func resourceAliCloudMessageServiceQueueUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	d.Partial(true)

	action := "SetQueueAttributes"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["QueueName"] = d.Id()
	if !d.IsNewResource() && d.HasChange("message_retention_period") {
		update = true

		if v, ok := d.GetOk("message_retention_period"); ok {
			request["MessageRetentionPeriod"] = v
		}
	}

	if !d.IsNewResource() && d.HasChange("polling_wait_seconds") {
		update = true

		if v, ok := d.GetOkExists("polling_wait_seconds"); ok {
			request["PollingWaitSeconds"] = v
		}
	}

	if !d.IsNewResource() && d.HasChange("visibility_timeout") {
		update = true

		if v, ok := d.GetOk("visibility_timeout"); ok {
			request["VisibilityTimeout"] = v
		}
	}

	if !d.IsNewResource() && d.HasChange("delay_seconds") {
		update = true

		if v, ok := d.GetOkExists("delay_seconds"); ok {
			request["DelaySeconds"] = v
		}
	}

	if !d.IsNewResource() && d.HasChange("maximum_message_size") {
		update = true

		if v, ok := d.GetOkExists("maximum_message_size"); ok {
			request["MaximumMessageSize"] = v
		}
	}

	if !d.IsNewResource() && d.HasChange("logging_enabled") {
		update = true
	}
	if v, ok := d.GetOkExists("logging_enabled"); ok {
		request["EnableLogging"] = v
	}

	if !d.IsNewResource() && d.HasChange("dlq_policy") {
		update = true
		objectDataLocalMap := make(map[string]interface{})

		if v := d.Get("dlq_policy"); v != nil {
			enabled1, _ := jsonpath.Get("$[0].enabled", v)
			if enabled1 != nil && (d.HasChange("dlq_policy.0.enabled") || enabled1 != "") {
				objectDataLocalMap["Enabled"] = enabled1

				if objectDataLocalMap["Enabled"].(bool) {
					deadLetterTargetQueue1, _ := jsonpath.Get("$[0].dead_letter_target_queue", v)
					if deadLetterTargetQueue1 != nil && (d.HasChange("dlq_policy.0.dead_letter_target_queue") || deadLetterTargetQueue1 != "") {
						objectDataLocalMap["DeadLetterTargetQueue"] = deadLetterTargetQueue1
					}
					maxReceiveCount1, _ := jsonpath.Get("$[0].max_receive_count", v)
					if maxReceiveCount1 != nil && (d.HasChange("dlq_policy.0.max_receive_count") || maxReceiveCount1 != "") {
						objectDataLocalMap["MaxReceiveCount"] = maxReceiveCount1
					}
				}
			}

			objectDataLocalMapJson, err := json.Marshal(objectDataLocalMap)
			if err != nil {
				return WrapError(err)
			}
			request["DlqPolicy"] = string(objectDataLocalMapJson)
		}
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Mns-open", "2022-01-19", action, query, request, true)
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

	if !d.IsNewResource() && d.HasChange("tags") {
		messageServiceServiceV2 := MessageServiceServiceV2{client}
		if err := messageServiceServiceV2.SetResourceTags(d, "queue"); err != nil {
			return WrapError(err)
		}
	}
	d.Partial(false)
	return resourceAliCloudMessageServiceQueueRead(d, meta)
}

func resourceAliCloudMessageServiceQueueDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteQueue"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["QueueName"] = d.Id()

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Mns-open", "2022-01-19", action, query, request, true)

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
