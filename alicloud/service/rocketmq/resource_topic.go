// Package rocketmq. This file is generated automatically. Please do not modify it manually, thank you!
package rocketmq

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/validators"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceTopic() *schema.Resource {
	return &schema.Resource{
		Create: resourceTopicCreate,
		Read:   resourceTopicRead,
		Update: resourceTopicUpdate,
		Delete: resourceTopicDelete,
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
				ValidateFunc: validators.StringInSlice([]string{"NORMAL", "FIFO", "DELAY", "TRANSACTION"}, false),
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"remark": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validators.StringMatch(regexp.MustCompile("^[\u4E00-\u9FA5A-Za-z0-9_]+$"), "Custom remarks"),
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"topic_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validators.StringMatch(regexp.MustCompile("^[%a-zA-Z0-9_-]+$"), "Topic name and identification"),
			},
		},
	}
}

func resourceTopicCreate(d *schema.ResourceData, meta interface{}) error {

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
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RoaPost("RocketMQ", "2022-08-01", action, query, nil, body, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_rocketmq_topic", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", instanceId, topicName))

	rocketmqService := NewService(client)
	stateConf := helper.BuildStateConf([]string{}, []string{"RUNNING"}, d.Timeout(schema.TimeoutCreate), 50*time.Second, rocketmqService.TopicStateRefreshFunc(d.Id(), "status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return resourceTopicRead(d, meta)
}

func resourceTopicRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rocketmqService := NewService(client)

	objectRaw, err := rocketmqService.DescribeTopic(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_rocketmq_topic DescribeTopic Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
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

func resourceTopicUpdate(d *schema.ResourceData, meta interface{}) error {
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
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RoaPatch("RocketMQ", "2022-08-01", action, query, nil, body, true)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
		rocketmqService := NewService(client)
		stateConf := helper.BuildStateConf([]string{}, []string{"RUNNING"}, d.Timeout(schema.TimeoutUpdate), 30*time.Second, rocketmqService.TopicStateRefreshFunc(d.Id(), "status", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
		}
	}

	return resourceTopicRead(d, meta)
}

func resourceTopicDelete(d *schema.ResourceData, meta interface{}) error {

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

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RoaDelete("RocketMQ", "2022-08-01", action, query, nil, nil, true)

		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		if errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	rocketmqService := NewService(client)
	stateConf := helper.BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 10*time.Second, rocketmqService.TopicStateRefreshFunc(d.Id(), "topicName", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return nil
}
