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

func resourceAliCloudMessageServiceTopic() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudMessageServiceTopicCreate,
		Read:   resourceAliCloudMessageServiceTopicRead,
		Update: resourceAliCloudMessageServiceTopicUpdate,
		Delete: resourceAliCloudMessageServiceTopicDelete,
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
				Type:     schema.TypeString,
				Computed: true,
			},
			"enable_logging": {
				Type:          schema.TypeBool,
				Optional:      true,
				ConflictsWith: []string{"logging_enabled"},
				Computed:      true,
			},
			"max_message_size": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: IntBetween(1024, 65536),
			},
			"tags": tagsSchema(),
			"topic_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"logging_enabled": {
				Type:       schema.TypeBool,
				Optional:   true,
				Computed:   true,
				Deprecated: "Field `logging_enabled` has been deprecated from provider version 1.241.0. New field `enable_logging` instead.",
			},
		},
	}
}

func resourceAliCloudMessageServiceTopicCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateTopic"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	conn, err := client.NewMnsClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["TopicName"] = d.Get("topic_name")

	if v, ok := d.GetOkExists("enable_logging"); ok {
		request["EnableLogging"] = v
	} else if v, ok := d.GetOkExists("logging_enabled"); ok {
		request["EnableLogging"] = v
	}
	if v, ok := d.GetOkExists("max_message_size"); ok {
		request["MaxMessageSize"] = v
	}
	if v, ok := d.GetOk("tags"); ok {
		tagsMap := ConvertTags(v.(map[string]interface{}))
		request = expandTagsToMap(request, tagsMap)
	}

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2022-01-19"), StringPointer("AK"), query, request, &runtime)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_message_service_topic", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["TopicName"]))

	return resourceAliCloudMessageServiceTopicRead(d, meta)
}

func resourceAliCloudMessageServiceTopicRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	messageServiceServiceV2 := MessageServiceServiceV2{client}

	objectRaw, err := messageServiceServiceV2.DescribeMessageServiceTopic(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_message_service_topic DescribeMessageServiceTopic Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	if objectRaw["CreateTime"] != nil {
		d.Set("create_time", objectRaw["CreateTime"])
	}
	if objectRaw["LoggingEnabled"] != nil {
		d.Set("enable_logging", objectRaw["LoggingEnabled"])
		d.Set("logging_enabled", objectRaw["LoggingEnabled"])
	}
	if objectRaw["MaxMessageSize"] != nil {
		d.Set("max_message_size", objectRaw["MaxMessageSize"])
	}
	if objectRaw["TopicName"] != nil {
		d.Set("topic_name", objectRaw["TopicName"])
	}

	if objectRaw["Tags"] != nil {
		d.Set("tags", tagsToMap(objectRaw["Tags"]))
	}

	return nil
}

func resourceAliCloudMessageServiceTopicUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	action := "SetTopicAttributes"
	conn, err := client.NewMnsClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["TopicName"] = d.Id()

	if d.HasChange("enable_logging") {
		update = true

		if v, ok := d.GetOkExists("enable_logging"); ok {
			request["EnableLogging"] = v
		}
	}

	if d.HasChange("logging_enabled") {
		update = true

		if v, ok := d.GetOkExists("logging_enabled"); ok {
			request["EnableLogging"] = v
		}
	}

	if d.HasChange("max_message_size") {
		update = true

		if v, ok := d.GetOkExists("max_message_size"); ok {
			request["MaxMessageSize"] = v
		}
	}

	if update {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2022-01-19"), StringPointer("AK"), query, request, &runtime)
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

	if d.HasChange("tags") {
		messageServiceServiceV2 := MessageServiceServiceV2{client}
		if err := messageServiceServiceV2.SetResourceTags(d, "topic"); err != nil {
			return WrapError(err)
		}
	}
	return resourceAliCloudMessageServiceTopicRead(d, meta)
}

func resourceAliCloudMessageServiceTopicDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteTopic"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	conn, err := client.NewMnsClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["TopicName"] = d.Id()

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2022-01-19"), StringPointer("AK"), query, request, &runtime)

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
