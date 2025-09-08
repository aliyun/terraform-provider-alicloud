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
			"dlq_policy": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
			"endpoint": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"sts_role_arn": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"filter_tag": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"notify_content_format": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"notify_strategy": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"push_type": {
				Type:     schema.TypeString,
				Required: true,
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
	query := make(map[string]interface{})
	var err error
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
	if v, ok := d.GetOk("sts_role_arn"); ok {
		request["StsRoleArn"] = v
	}
	request["PushType"] = d.Get("push_type")
	objectDataLocalMap := make(map[string]interface{})

	if v := d.Get("dlq_policy"); !IsNil(v) {
		deadLetterTargetQueue1, _ := jsonpath.Get("$[0].dead_letter_target_queue", v)
		if deadLetterTargetQueue1 != nil && deadLetterTargetQueue1 != "" {
			objectDataLocalMap["DeadLetterTargetQueue"] = deadLetterTargetQueue1
		}
		enabled1, _ := jsonpath.Get("$[0].enabled", v)
		if enabled1 != nil && enabled1 != "" {
			objectDataLocalMap["Enabled"] = enabled1
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

	dlqPolicyMaps := make([]map[string]interface{}, 0)
	dlqPolicyMap := make(map[string]interface{})
	dlqPolicyRaw := make(map[string]interface{})
	if objectRaw["DlqPolicy"] != nil {
		dlqPolicyRaw = objectRaw["DlqPolicy"].(map[string]interface{})
	}
	if len(dlqPolicyRaw) > 0 {
		dlqPolicyMap["dead_letter_target_queue"] = dlqPolicyRaw["DeadLetterTargetQueue"]
		dlqPolicyMap["enabled"] = dlqPolicyRaw["Enabled"]

		dlqPolicyMaps = append(dlqPolicyMaps, dlqPolicyMap)
	}
	if err := d.Set("dlq_policy", dlqPolicyMaps); err != nil {
		return err
	}

	return nil
}

func resourceAliCloudMessageServiceSubscriptionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	action := "SetSubscriptionAttributes"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["TopicName"] = parts[0]
	request["SubscriptionName"] = parts[1]

	if d.HasChange("notify_strategy") {
		update = true
	}
	if v, ok := d.GetOk("notify_strategy"); ok {
		request["NotifyStrategy"] = v
	}

	if d.HasChange("dlq_policy") {
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

	return resourceAliCloudMessageServiceSubscriptionRead(d, meta)
}

func resourceAliCloudMessageServiceSubscriptionDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	action := "Unsubscribe"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	request = make(map[string]interface{})
	request["SubscriptionName"] = parts[1]
	request["TopicName"] = parts[0]

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
