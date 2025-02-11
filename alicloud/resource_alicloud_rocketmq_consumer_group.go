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
						"retry_policy": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: StringInSlice([]string{"DefaultRetryPolicy", "FixedRetryPolicy", "BackoffRetryPolicy"}, false),
						},
						"max_retry_times": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: IntBetween(1, 1000),
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
		nodeNative, _ := jsonpath.Get("$[0].max_retry_times", v)
		if nodeNative != "" {
			objectDataLocalMap["maxRetryTimes"] = nodeNative
		}
		nodeNative1, _ := jsonpath.Get("$[0].retry_policy", v)
		if nodeNative1 != "" {
			objectDataLocalMap["retryPolicy"] = nodeNative1
		}
	}
	request["consumeRetryPolicy"] = objectDataLocalMap

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
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_rocketmq_consumer_group", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", instanceId, consumerGroupId))

	rocketmqServiceV2 := RocketmqServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"RUNNING"}, d.Timeout(schema.TimeoutCreate), 30*time.Second, rocketmqServiceV2.RocketmqConsumerGroupStateRefreshFunc(d.Id(), "status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudRocketmqConsumerGroupUpdate(d, meta)
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
	d.Set("remark", objectRaw["remark"])
	d.Set("status", objectRaw["status"])
	d.Set("consumer_group_id", objectRaw["consumerGroupId"])
	d.Set("instance_id", objectRaw["instanceId"])
	consumeRetryPolicyMaps := make([]map[string]interface{}, 0)
	consumeRetryPolicyMap := make(map[string]interface{})
	consumeRetryPolicy1Raw := make(map[string]interface{})
	if objectRaw["consumeRetryPolicy"] != nil {
		consumeRetryPolicy1Raw = objectRaw["consumeRetryPolicy"].(map[string]interface{})
	}
	if len(consumeRetryPolicy1Raw) > 0 {
		consumeRetryPolicyMap["max_retry_times"] = consumeRetryPolicy1Raw["maxRetryTimes"]
		consumeRetryPolicyMap["retry_policy"] = consumeRetryPolicy1Raw["retryPolicy"]
		consumeRetryPolicyMaps = append(consumeRetryPolicyMaps, consumeRetryPolicyMap)
	}
	d.Set("consume_retry_policy", consumeRetryPolicyMaps)

	return nil
}

func resourceAliCloudRocketmqConsumerGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var body map[string]interface{}
	update := false
	parts := strings.Split(d.Id(), ":")
	instanceId := parts[0]
	consumerGroupId := parts[1]
	action := fmt.Sprintf("/instances/%s/consumerGroups/%s", instanceId, consumerGroupId)
	var err error
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})
	if !d.IsNewResource() && d.HasChange("remark") {
		update = true
	}
	request["remark"] = d.Get("remark")

	if !d.IsNewResource() && d.HasChange("delivery_order_type") {
		update = true
	}
	request["deliveryOrderType"] = d.Get("delivery_order_type")

	if !d.IsNewResource() && d.HasChange("consume_retry_policy") {
		update = true
	}
	objectDataLocalMap := make(map[string]interface{})
	if v := d.Get("consume_retry_policy"); !IsNil(v) {
		nodeNative, _ := jsonpath.Get("$[0].max_retry_times", v)
		if nodeNative != "" {
			objectDataLocalMap["maxRetryTimes"] = nodeNative
		}
		nodeNative1, _ := jsonpath.Get("$[0].retry_policy", v)
		if nodeNative1 != "" {
			objectDataLocalMap["retryPolicy"] = nodeNative1
		}
	}
	request["consumeRetryPolicy"] = objectDataLocalMap

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
			addDebug(action, response, request)
			return nil
		})
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
	body := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	body = request
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RoaDelete("RocketMQ", "2022-08-01", action, query, nil, body, true)

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

	rocketmqServiceV2 := RocketmqServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 30*time.Second, rocketmqServiceV2.RocketmqConsumerGroupStateRefreshFunc(d.Id(), "consumerGroupId", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
