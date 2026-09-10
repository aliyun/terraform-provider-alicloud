// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudRocketmqTopic() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudRocketmqTopicCreate,
		Read:   resourceAliCloudRocketmqTopicRead,
		Update: resourceAliCloudRocketmqTopicUpdate,
		Delete: resourceAliCloudRocketmqTopicDelete,
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
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"max_send_tps": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"message_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"NORMAL", "FIFO", "DELAY", "TRANSACTION"}, false),
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"remark": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringMatch(regexp.MustCompile("^[\u4E00-\u9FA5A-Za-z0-9_]+$"), "Custom remarks"),
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"topic_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringMatch(regexp.MustCompile("^[%a-zA-Z0-9_-]+$"), "Topic name and identification"),
			},
		},
	}
}

func resourceAliCloudRocketmqTopicCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	instanceId := d.Get("instance_id")
	topicName := d.Get("topic_name")
	action := fmt.Sprintf("/instances/%s/topics/%s", instanceId, topicName)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	if v, ok := d.GetOk("message_type"); ok {
		request["messageType"] = v
	}
	if v, ok := d.GetOk("remark"); ok {
		request["remark"] = v
	}
	if v, ok := d.GetOkExists("max_send_tps"); ok {
		request["maxSendTps"] = v
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_rocketmq_topic", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", instanceId, topicName))

	rocketmqServiceV2 := RocketmqServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"RUNNING"}, d.Timeout(schema.TimeoutCreate), 50*time.Second, rocketmqServiceV2.RocketmqTopicStateRefreshFunc(d.Id(), "status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudRocketmqTopicRead(d, meta)
}

func resourceAliCloudRocketmqTopicRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rocketmqServiceV2 := RocketmqServiceV2{client}

	objectRaw, err := rocketmqServiceV2.DescribeRocketmqTopic(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_rocketmq_topic DescribeRocketmqTopic Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("create_time", objectRaw["createTime"])
	d.Set("max_send_tps", objectRaw["maxSendTps"])
	d.Set("message_type", objectRaw["messageType"])
	d.Set("region_id", objectRaw["regionId"])
	d.Set("remark", objectRaw["remark"])
	d.Set("status", objectRaw["status"])
	d.Set("instance_id", objectRaw["instanceId"])
	d.Set("topic_name", objectRaw["topicName"])

	return nil
}

func resourceAliCloudRocketmqTopicUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var body map[string]interface{}
	update := false

	var err error
	parts := strings.Split(d.Id(), ":")
	instanceId := parts[0]
	topicName := parts[1]
	action := fmt.Sprintf("/instances/%s/topics/%s", instanceId, topicName)
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})

	if d.HasChange("remark") {
		update = true
	}
	if v, ok := d.GetOk("remark"); ok {
		request["remark"] = v
	}
	if d.HasChange("max_send_tps") {
		update = true
	}
	if v, ok := d.GetOkExists("max_send_tps"); ok && v.(int) != 0 {
		request["maxSendTps"] = v
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
		stateConf := BuildStateConf([]string{}, []string{"RUNNING"}, d.Timeout(schema.TimeoutUpdate), 30*time.Second, rocketmqServiceV2.RocketmqTopicStateRefreshFunc(d.Id(), "status", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	return resourceAliCloudRocketmqTopicRead(d, meta)
}

func resourceAliCloudRocketmqTopicDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	instanceId := parts[0]
	topicName := parts[1]
	action := fmt.Sprintf("/instances/%s/topics/%s", instanceId, topicName)
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
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 10*time.Second, rocketmqServiceV2.RocketmqTopicStateRefreshFunc(d.Id(), "topicName", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}
