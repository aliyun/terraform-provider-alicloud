package alicloud

import (
	"fmt"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudMessageServiceSubscription() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudMessageServiceSubscriptionCreate,
		Read:   resourceAlicloudMessageServiceSubscriptionRead,
		Update: resourceAlicloudMessageServiceSubscriptionUpdate,
		Delete: resourceAlicloudMessageServiceSubscriptionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(3 * time.Minute),
			Update: schema.DefaultTimeout(3 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"topic_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"subscription_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"endpoint": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"push_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"http", "queue", "mpush", "alisms", "email"}, false),
			},
			"filter_tag": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(0, 16),
			},
			"notify_content_format": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"XML", "JSON", "SIMPLIFIED"}, false),
			},
			"notify_strategy": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"BACKOFF_RETRY", "EXPONENTIAL_DECAY_RETRY"}, false),
			},
		},
	}
}

func resourceAlicloudMessageServiceSubscriptionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "Subscribe"
	request := make(map[string]interface{})
	var err error

	request["TopicName"] = d.Get("topic_name")
	request["SubscriptionName"] = d.Get("subscription_name")
	request["Endpoint"] = d.Get("endpoint")
	request["PushType"] = d.Get("push_type")

	if v, ok := d.GetOk("filter_tag"); ok {
		request["MessageTag"] = v
	}

	if v, ok := d.GetOk("notify_content_format"); ok {
		request["NotifyContentFormat"] = v
	}

	if v, ok := d.GetOk("notify_strategy"); ok {
		request["NotifyStrategy"] = v
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = client.RpcPost("Mns-open", "2022-01-19", action, nil, request, false)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_message_service_subscription", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["TopicName"], ":", request["SubscriptionName"]))

	return resourceAlicloudMessageServiceSubscriptionRead(d, meta)
}

func resourceAlicloudMessageServiceSubscriptionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	mnsOpenService := MnsOpenService{client}

	object, err := mnsOpenService.DescribeMessageServiceSubscription(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("topic_name", object["TopicName"])
	d.Set("subscription_name", object["SubscriptionName"])
	d.Set("endpoint", object["Endpoint"])
	d.Set("filter_tag", object["FilterTag"])
	d.Set("notify_content_format", object["NotifyContentFormat"])
	d.Set("notify_strategy", object["NotifyStrategy"])

	return nil
}

func resourceAlicloudMessageServiceSubscriptionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	update := false
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"TopicName":        parts[0],
		"SubscriptionName": parts[1],
	}

	if d.HasChange("notify_strategy") {
		update = true
	}
	if v, ok := d.GetOk("notify_strategy"); ok {
		request["NotifyStrategy"] = v
	}

	if update {
		action := "SetSubscriptionAttributes"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.RpcPost("Mns-open", "2022-01-19", action, nil, request, false)
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

	return resourceAlicloudMessageServiceSubscriptionRead(d, meta)
}

func resourceAlicloudMessageServiceSubscriptionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "Unsubscribe"
	var response map[string]interface{}
	var err error

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"TopicName":        parts[0],
		"SubscriptionName": parts[1],
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = client.RpcPost("Mns-open", "2022-01-19", action, nil, request, false)
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
