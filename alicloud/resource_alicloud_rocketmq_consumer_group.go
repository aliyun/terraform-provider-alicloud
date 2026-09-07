// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudRocketmqConsumerGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudRocketmqConsumerGroupCreate,
		Read:   resourceAliCloudRocketmqConsumerGroupRead,
		Update: resourceAliCloudRocketmqConsumerGroupUpdate,
		Delete: resourceAliCloudRocketmqConsumerGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"consume_retry_policy": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"dead_letter_target_topic": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"retry_policy": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: StringInSlice([]string{"DefaultRetryPolicy", "FixedRetryPolicy", "BackoffRetryPolicy"}, false),
						},
						"max_retry_times": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: IntBetween(0, 1000),
						},
					},
				},
			},
			"consumer_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"delivery_order_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"Concurrently", "Orderly"}, false),
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"max_receive_tps": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"remark": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudRocketmqConsumerGroupCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	instanceId := d.Get("instance_id")
	consumerGroupId := d.Get("consumer_group_id")
	action := fmt.Sprintf("/instances/%s/consumerGroups/%s", instanceId, consumerGroupId)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	if v, ok := d.GetOk("remark"); ok {
		request["remark"] = v
	}
	if v, ok := d.GetOk("delivery_order_type"); ok {
		request["deliveryOrderType"] = v
	}
	objectDataLocalMap := make(map[string]interface{})

	if v := d.Get("consume_retry_policy"); !IsNil(v) {
		maxRetryTimes1, _ := jsonpath.Get("$[0].max_retry_times", v)
		if maxRetryTimes1 != nil && maxRetryTimes1 != "" {
			objectDataLocalMap["maxRetryTimes"] = maxRetryTimes1
		}
		retryPolicy1, _ := jsonpath.Get("$[0].retry_policy", v)
		if retryPolicy1 != nil && retryPolicy1 != "" {
			objectDataLocalMap["retryPolicy"] = retryPolicy1
		}
		deadLetterTargetTopic1, _ := jsonpath.Get("$[0].dead_letter_target_topic", v)
		if deadLetterTargetTopic1 != nil && deadLetterTargetTopic1 != "" {
			objectDataLocalMap["deadLetterTargetTopic"] = deadLetterTargetTopic1
		}

		request["consumeRetryPolicy"] = objectDataLocalMap
	}

	if v, ok := d.GetOkExists("max_receive_tps"); ok {
		request["maxReceiveTps"] = v
	}
	body = request
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RoaPost("RocketMQ", "2022-08-01", action, query, nil, body, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_rocketmq_consumer_group", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", instanceId, consumerGroupId))

	rocketmqServiceV2 := RocketmqServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"RUNNING"}, d.Timeout(schema.TimeoutCreate), 30*time.Second, rocketmqServiceV2.RocketmqConsumerGroupStateRefreshFunc(d.Id(), "status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudRocketmqConsumerGroupRead(d, meta)
}

func resourceAliCloudRocketmqConsumerGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rocketmqServiceV2 := RocketmqServiceV2{client}

	objectRaw, err := rocketmqServiceV2.DescribeRocketmqConsumerGroup(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_rocketmq_consumer_group DescribeRocketmqConsumerGroup Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("create_time", objectRaw["createTime"])
	d.Set("delivery_order_type", objectRaw["deliveryOrderType"])
	d.Set("max_receive_tps", objectRaw["maxReceiveTps"])
	d.Set("region_id", objectRaw["regionId"])
	d.Set("remark", objectRaw["remark"])
	d.Set("status", objectRaw["status"])
	d.Set("consumer_group_id", objectRaw["consumerGroupId"])
	d.Set("instance_id", objectRaw["instanceId"])

	consumeRetryPolicyMaps := make([]map[string]interface{}, 0)
	consumeRetryPolicyMap := make(map[string]interface{})
	consumeRetryPolicyRaw := make(map[string]interface{})
	if objectRaw["consumeRetryPolicy"] != nil {
		consumeRetryPolicyRaw = objectRaw["consumeRetryPolicy"].(map[string]interface{})
	}
	if len(consumeRetryPolicyRaw) > 0 {
		consumeRetryPolicyMap["dead_letter_target_topic"] = consumeRetryPolicyRaw["deadLetterTargetTopic"]
		consumeRetryPolicyMap["max_retry_times"] = consumeRetryPolicyRaw["maxRetryTimes"]
		consumeRetryPolicyMap["retry_policy"] = consumeRetryPolicyRaw["retryPolicy"]

		consumeRetryPolicyMaps = append(consumeRetryPolicyMaps, consumeRetryPolicyMap)
	}
	if err := d.Set("consume_retry_policy", consumeRetryPolicyMaps); err != nil {
		return err
	}

	return nil
}

func resourceAliCloudRocketmqConsumerGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var body map[string]interface{}
	update := false

	var err error
	parts := strings.Split(d.Id(), ":")
	instanceId := parts[0]
	consumerGroupId := parts[1]
	action := fmt.Sprintf("/instances/%s/consumerGroups/%s", instanceId, consumerGroupId)
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})

	if d.HasChange("remark") {
		update = true
	}
	if v, ok := d.GetOk("remark"); ok {
		request["remark"] = v
	}

	if d.HasChange("delivery_order_type") {
		update = true
	}
	if v, ok := d.GetOk("delivery_order_type"); ok {
		request["deliveryOrderType"] = v
	}

	if d.HasChange("consume_retry_policy") {
		update = true
	}
	objectDataLocalMap := make(map[string]interface{})

	if v := d.Get("consume_retry_policy"); !IsNil(v) {
		maxRetryTimes1, _ := jsonpath.Get("$[0].max_retry_times", v)
		if maxRetryTimes1 != nil && maxRetryTimes1 != "" {
			objectDataLocalMap["maxRetryTimes"] = maxRetryTimes1
		}
		retryPolicy1, _ := jsonpath.Get("$[0].retry_policy", v)
		if retryPolicy1 != nil && retryPolicy1 != "" {
			objectDataLocalMap["retryPolicy"] = retryPolicy1
		}
		deadLetterTargetTopic1, _ := jsonpath.Get("$[0].dead_letter_target_topic", v)
		if deadLetterTargetTopic1 != nil && deadLetterTargetTopic1 != "" {
			objectDataLocalMap["deadLetterTargetTopic"] = deadLetterTargetTopic1
		}

		request["consumeRetryPolicy"] = objectDataLocalMap
	}

	if d.HasChange("max_receive_tps") {
		update = true
	}
	if v, ok := d.GetOkExists("max_receive_tps"); ok && v.(int) != 0 {
		request["maxReceiveTps"] = v
	}
	body = request
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RoaPatch("RocketMQ", "2022-08-01", action, query, nil, body, true)
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
		rocketmqServiceV2 := RocketmqServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"RUNNING"}, d.Timeout(schema.TimeoutUpdate), 30*time.Second, rocketmqServiceV2.RocketmqConsumerGroupStateRefreshFunc(d.Id(), "status", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	return resourceAliCloudRocketmqConsumerGroupRead(d, meta)
}

func resourceAliCloudRocketmqConsumerGroupDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	instanceId := parts[0]
	consumerGroupId := parts[1]
	action := fmt.Sprintf("/instances/%s/consumerGroups/%s", instanceId, consumerGroupId)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	var err error
	request = make(map[string]interface{})

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RoaDelete("RocketMQ", "2022-08-01", action, query, nil, nil, true)

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

	rocketmqServiceV2 := RocketmqServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 30*time.Second, rocketmqServiceV2.RocketmqConsumerGroupStateRefreshFunc(d.Id(), "consumerGroupId", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}
