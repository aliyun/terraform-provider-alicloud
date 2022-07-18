package alicloud

import (
	"fmt"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	ali_mns "github.com/aliyun/aliyun-mns-go-sdk"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudMNSSubscription() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudMNSSubscriptionCreate,
		Read:   resourceAlicloudMNSSubscriptionRead,
		Update: resourceAlicloudMNSSubscriptionUpdate,
		Delete: resourceAlicloudMNSSubscriptionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"topic_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(3, 256),
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(3, 256),
			},
			"endpoint": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"filter_tag": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(0, 16),
			},
			"notify_strategy": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  string(ali_mns.BACKOFF_RETRY),
				ValidateFunc: validation.StringInSlice([]string{
					string(ali_mns.BACKOFF_RETRY), string(ali_mns.EXPONENTIAL_DECAY_RETRY),
				}, false),
			},
			"notify_content_format": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  string(ali_mns.SIMPLIFIED),
				ForceNew: true,
			},
			"push_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"http", "queue", "mpush", "email"}, false),
			},
		},
	}
}

func resourceAlicloudMNSSubscriptionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := make(map[string]interface{})
	conn, err := client.NewMnsClient()
	if err != nil {
		return WrapError(err)
	}

	request["TopicName"] = d.Get("topic_name")
	request["SubscriptionName"] = d.Get("name")
	request["Endpoint"] = d.Get("endpoint")
	request["PushType"] = d.Get("push_type")
	if v, ok := d.GetOk("filter_tag"); ok {
		request["MessageTag"] = v
	}
	if v, ok := d.GetOk("notify_strategy"); ok {
		request["NotifyStrategy"] = v
	}
	if v, ok := d.GetOk("notify_content_format"); ok {
		request["NotifyContentFormat"] = v
	}

	var response map[string]interface{}
	action := "Subscribe"
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_mns_topic_subscription", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%s%s%s", request["TopicName"], COLON_SEPARATED, request["SubscriptionName"]))
	return resourceAlicloudMNSSubscriptionRead(d, meta)
}

func resourceAlicloudMNSSubscriptionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	mnsService := MnsService{client}

	object, err := mnsService.DescribeMessageServiceSubscription(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	d.Set("topic_name", parts[0])
	d.Set("name", parts[1])
	d.Set("endpoint", object["Endpoint"])
	d.Set("filter_tag", object["FilterTag"])
	d.Set("notify_strategy", object["NotifyStrategy"])
	d.Set("notify_content_format", object["NotifyContentFormat"])
	return nil
}

func resourceAlicloudMNSSubscriptionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	conn, err := client.NewMnsClient()
	if err != nil {
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	update := false
	request := map[string]interface{}{
		"TopicName":        parts[0],
		"SubscriptionName": parts[1],
	}

	if !d.IsNewResource() && d.HasChange("notify_strategy") {
		update = true
		if v, ok := d.GetOk("notify_strategy"); ok {
			request["NotifyStrategy"] = v
		}
	}

	if update {
		action := "SetSubscriptionAttributes"
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

	return resourceAlicloudMNSSubscriptionRead(d, meta)
}

func resourceAlicloudMNSSubscriptionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	conn, err := client.NewMnsClient()
	if err != nil {
		return WrapError(err)
	}

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}

	request := map[string]interface{}{
		"TopicName": parts[0], "SubscriptionName": parts[1],
	}

	action := "Unsubscribe"
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
