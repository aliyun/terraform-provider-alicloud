// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudMessageServiceSubscription() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudMessageServiceSubscriptionCreate,
		Read:   resourceAliCloudMessageServiceSubscriptionRead,
		Update: resourceAliCloudMessageServiceSubscriptionUpdate,
		Delete: resourceAliCloudMessageServiceSubscriptionDelete,
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
			"endpoint": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"filter_tag": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"notify_content_format": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"XML", "JSON", "SIMPLIFIED"}, false),
			},
			"notify_strategy": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"BACKOFF_RETRY", "EXPONENTIAL_DECAY_RETRY"}, false),
			},
			"push_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: StringInSlice([]string{"http", "queue", "mpush", "alisms", "email"}, false),
			},
			"subscription_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"topic_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudMessageServiceSubscriptionCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "Subscribe"
	var request map[string]interface{}
	var response map[string]interface{}
	conn, err := client.NewMnsClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["TopicName"] = d.Get("topic_name")
	request["SubscriptionName"] = d.Get("subscription_name")

	if v, ok := d.GetOk("notify_strategy"); ok {
		request["NotifyStrategy"] = v
	}
	if v, ok := d.GetOk("notify_content_format"); ok {
		request["NotifyContentFormat"] = v
	}
	if v, ok := d.GetOk("filter_tag"); ok {
		request["MessageTag"] = v
	}
	request["Endpoint"] = d.Get("endpoint")
	request["PushType"] = d.Get("push_type")
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_message_service_subscription", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", request["TopicName"], request["SubscriptionName"]))

	return resourceAliCloudMessageServiceSubscriptionRead(d, meta)
}

func resourceAliCloudMessageServiceSubscriptionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	messageServiceServiceV2 := MessageServiceServiceV2{client}

	objectRaw, err := messageServiceServiceV2.DescribeMessageServiceSubscription(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_message_service_subscription DescribeMessageServiceSubscription Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("create_time", objectRaw["CreateTime"])
	d.Set("endpoint", objectRaw["Endpoint"])
	d.Set("filter_tag", objectRaw["FilterTag"])
	d.Set("notify_content_format", objectRaw["NotifyContentFormat"])
	d.Set("notify_strategy", objectRaw["NotifyStrategy"])
	d.Set("subscription_name", objectRaw["SubscriptionName"])
	d.Set("topic_name", objectRaw["TopicName"])

	return nil
}

func resourceAliCloudMessageServiceSubscriptionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	update := false
	parts := strings.Split(d.Id(), ":")
	action := "SetSubscriptionAttributes"
	conn, err := client.NewMnsClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["TopicName"] = parts[0]
	request["SubscriptionName"] = parts[1]
	if d.HasChange("notify_strategy") {
		update = true
		request["NotifyStrategy"] = d.Get("notify_strategy")
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

	return resourceAliCloudMessageServiceSubscriptionRead(d, meta)
}

func resourceAliCloudMessageServiceSubscriptionDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "Unsubscribe"
	var request map[string]interface{}
	var response map[string]interface{}
	conn, err := client.NewMnsClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["SubscriptionName"] = parts[1]
	request["TopicName"] = parts[0]

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
