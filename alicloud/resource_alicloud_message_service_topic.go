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
				Computed:      true,
				ConflictsWith: []string{"logging_enabled"},
			},
			"max_message_size": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: IntBetween(1024, 65536),
			},
			"topic_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"logging_enabled": {
				Type:       schema.TypeBool,
				Optional:   true,
				Computed:   true,
				Deprecated: "Field 'logging_enabled' has been deprecated since provider version 1.211.0. New field 'enable_logging' instead.",
			},
		},
	}
}

func resourceAliCloudMessageServiceTopicCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateTopic"
	var request map[string]interface{}
	var response map[string]interface{}
	conn, err := client.NewMnsClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["TopicName"] = d.Get("topic_name")

	if v, ok := d.GetOkExists("enable_logging"); ok {
		request["EnableLogging"] = v
	}
	if v, ok := d.GetOk("max_message_size"); ok {
		request["MaxMessageSize"] = v
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_message_service_topic", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["TopicName"]))

	return resourceAliCloudMessageServiceTopicUpdate(d, meta)
}

func resourceAliCloudMessageServiceTopicRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	messageServiceServiceV2 := MessageServiceServiceV2{client}

	objectRaw, err := messageServiceServiceV2.DescribeMessageServiceTopic(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_message_service_topic DescribeMessageServiceTopic Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("create_time", objectRaw["CreateTime"])
	d.Set("enable_logging", objectRaw["LoggingEnabled"])
	d.Set("max_message_size", objectRaw["MaxMessageSize"])
	d.Set("topic_name", objectRaw["TopicName"])

	d.Set("logging_enabled", d.Get("enable_logging"))
	return nil
}

func resourceAliCloudMessageServiceTopicUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	update := false
	action := "SetTopicAttributes"
	conn, err := client.NewMnsClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["TopicName"] = d.Id()
	if !d.IsNewResource() && d.HasChange("logging_enabled") {
		update = true
		request["EnableLogging"] = d.Get("logging_enabled")
	}

	if !d.IsNewResource() && d.HasChange("enable_logging") {
		update = true
		request["EnableLogging"] = d.Get("enable_logging")
	}

	if !d.IsNewResource() && d.HasChange("max_message_size") {
		update = true
		request["MaxMessageSize"] = d.Get("max_message_size")
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

	return resourceAliCloudMessageServiceTopicRead(d, meta)
}

func resourceAliCloudMessageServiceTopicDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteTopic"
	var request map[string]interface{}
	var response map[string]interface{}
	conn, err := client.NewMnsClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["TopicName"] = d.Id()

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
