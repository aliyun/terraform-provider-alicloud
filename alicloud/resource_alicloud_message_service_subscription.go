// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
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
			"dm_attributes": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"subject": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"account_name": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"dysms_attributes": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"template_code": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"sign_name": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
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
			"last_modify_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"notify_content_format": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
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
			"sts_role_arn": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"subscription_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"tenant_rate_limit_policy": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"max_receives_per_second": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
			"topic_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"topic_owner": {
				Type:     schema.TypeString,
				Computed: true,
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
	if v, ok := d.GetOk("topic_name"); ok {
		request["TopicName"] = v
	}
	if v, ok := d.GetOk("subscription_name"); ok {
		request["SubscriptionName"] = v
	}
	request["RegionId"] = client.RegionId

	dlqPolicy := make(map[string]interface{})

	if v := d.Get("dlq_policy"); !IsNil(v) {
		deadLetterTargetQueue1, _ := jsonpath.Get("$[0].dead_letter_target_queue", v)
		if deadLetterTargetQueue1 != nil && deadLetterTargetQueue1 != "" {
			dlqPolicy["DeadLetterTargetQueue"] = deadLetterTargetQueue1
		}
		enabled1, _ := jsonpath.Get("$[0].enabled", v)
		if enabled1 != nil && enabled1 != "" {
			dlqPolicy["Enabled"] = enabled1
		}

		dlqPolicyJson, err := json.Marshal(dlqPolicy)
		if err != nil {
			return WrapError(err)
		}
		request["DlqPolicy"] = string(dlqPolicyJson)
	}

	if v, ok := d.GetOk("notify_strategy"); ok {
		request["NotifyStrategy"] = v
	}
	if v, ok := d.GetOk("sts_role_arn"); ok {
		request["StsRoleArn"] = v
	}
	tenantRateLimitPolicy := make(map[string]interface{})

	if v := d.Get("tenant_rate_limit_policy"); !IsNil(v) {
		enabled3, _ := jsonpath.Get("$[0].enabled", v)
		if enabled3 != nil && enabled3 != "" {
			tenantRateLimitPolicy["Enabled"] = enabled3
		}
		maxReceivesPerSecond1, _ := jsonpath.Get("$[0].max_receives_per_second", v)
		if maxReceivesPerSecond1 != nil && maxReceivesPerSecond1 != "" {
			tenantRateLimitPolicy["MaxReceivesPerSecond"] = maxReceivesPerSecond1
		}

		tenantRateLimitPolicyJson, err := json.Marshal(tenantRateLimitPolicy)
		if err != nil {
			return WrapError(err)
		}
		request["TenantRateLimitPolicy"] = string(tenantRateLimitPolicyJson)
	}

	request["PushType"] = d.Get("push_type")
	dmAttributes := make(map[string]interface{})

	if v := d.Get("dm_attributes"); !IsNil(v) {
		subject1, _ := jsonpath.Get("$[0].subject", v)
		if subject1 != nil && subject1 != "" {
			dmAttributes["Subject"] = subject1
		}
		accountName1, _ := jsonpath.Get("$[0].account_name", v)
		if accountName1 != nil && accountName1 != "" {
			dmAttributes["AccountName"] = accountName1
		}

		dmAttributesJson, err := json.Marshal(dmAttributes)
		if err != nil {
			return WrapError(err)
		}
		request["DmAttributes"] = string(dmAttributesJson)
	}

	if v, ok := d.GetOk("notify_content_format"); ok {
		request["NotifyContentFormat"] = v
	}
	if v, ok := d.GetOk("filter_tag"); ok {
		request["MessageTag"] = v
	}
	request["Endpoint"] = d.Get("endpoint")
	dysmsAttributes := make(map[string]interface{})

	if v := d.Get("dysms_attributes"); !IsNil(v) {
		templateCode1, _ := jsonpath.Get("$[0].template_code", v)
		if templateCode1 != nil && templateCode1 != "" {
			dysmsAttributes["TemplateCode"] = templateCode1
		}
		signName1, _ := jsonpath.Get("$[0].sign_name", v)
		if signName1 != nil && signName1 != "" {
			dysmsAttributes["SignName"] = signName1
		}

		dysmsAttributesJson, err := json.Marshal(dysmsAttributes)
		if err != nil {
			return WrapError(err)
		}
		request["DysmsAttributes"] = string(dysmsAttributesJson)
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
	d.Set("last_modify_time", objectRaw["LastModifyTime"])
	d.Set("notify_content_format", objectRaw["NotifyContentFormat"])
	d.Set("notify_strategy", objectRaw["NotifyStrategy"])
	d.Set("topic_owner", objectRaw["TopicOwner"])
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
	tenantRateLimitPolicyMaps := make([]map[string]interface{}, 0)
	tenantRateLimitPolicyMap := make(map[string]interface{})
	tenantRateLimitPolicyRaw := make(map[string]interface{})
	if objectRaw["TenantRateLimitPolicy"] != nil {
		tenantRateLimitPolicyRaw = objectRaw["TenantRateLimitPolicy"].(map[string]interface{})
	}
	if len(tenantRateLimitPolicyRaw) > 0 {
		tenantRateLimitPolicyMap["enabled"] = tenantRateLimitPolicyRaw["Enabled"]
		tenantRateLimitPolicyMap["max_receives_per_second"] = tenantRateLimitPolicyRaw["MaxReceivesPerSecond"]

		tenantRateLimitPolicyMaps = append(tenantRateLimitPolicyMaps, tenantRateLimitPolicyMap)
	}
	if err := d.Set("tenant_rate_limit_policy", tenantRateLimitPolicyMaps); err != nil {
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
	parts := strings.Split(d.Id(), ":")
	action := "SetSubscriptionAttributes"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["TopicName"] = parts[0]
	request["SubscriptionName"] = parts[1]
	request["RegionId"] = client.RegionId
	if d.HasChange("dlq_policy") {
		update = true
		dlqPolicy := make(map[string]interface{})

		if v := d.Get("dlq_policy"); v != nil {
			enabled1, _ := jsonpath.Get("$[0].enabled", v)
			if enabled1 != nil && (d.HasChange("dlq_policy.0.enabled") || enabled1 != "") {
				dlqPolicy["Enabled"] = enabled1

				if dlqPolicy["Enabled"].(bool) {
					deadLetterTargetQueue1, _ := jsonpath.Get("$[0].dead_letter_target_queue", v)
					if deadLetterTargetQueue1 != nil && (d.HasChange("dlq_policy.0.dead_letter_target_queue") || deadLetterTargetQueue1 != "") {
						dlqPolicy["DeadLetterTargetQueue"] = deadLetterTargetQueue1
					}
				}
			}

			dlqPolicyJson, err := json.Marshal(dlqPolicy)
			if err != nil {
				return WrapError(err)
			}
			request["DlqPolicy"] = string(dlqPolicyJson)
		}
	}

	if d.HasChange("notify_strategy") {
		update = true
	}
	// NotifyStrategy is required by SetSubscriptionAttributes even when only
	// other attributes change; omitting it fails with Missing:NotifyStrategy.
	if v, ok := d.GetOk("notify_strategy"); ok {
		request["NotifyStrategy"] = v
	}

	if d.HasChange("tenant_rate_limit_policy") {
		update = true
		tenantRateLimitPolicy := make(map[string]interface{})

		if v := d.Get("tenant_rate_limit_policy"); v != nil {
			enabled3, _ := jsonpath.Get("$[0].enabled", v)
			if enabled3 != nil && enabled3 != "" {
				tenantRateLimitPolicy["Enabled"] = enabled3
			}
			maxReceivesPerSecond1, _ := jsonpath.Get("$[0].max_receives_per_second", v)
			if maxReceivesPerSecond1 != nil && maxReceivesPerSecond1 != "" {
				tenantRateLimitPolicy["MaxReceivesPerSecond"] = maxReceivesPerSecond1
			}

			tenantRateLimitPolicyJson, err := json.Marshal(tenantRateLimitPolicy)
			if err != nil {
				return WrapError(err)
			}
			request["TenantRateLimitPolicy"] = string(tenantRateLimitPolicyJson)
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
	parts := strings.Split(d.Id(), ":")
	action := "Unsubscribe"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["TopicName"] = parts[0]
	request["SubscriptionName"] = parts[1]
	request["RegionId"] = client.RegionId

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
