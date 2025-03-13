// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudRocketmqAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudRocketmqAccountCreate,
		Read:   resourceAliCloudRocketmqAccountRead,
		Update: resourceAliCloudRocketmqAccountUpdate,
		Delete: resourceAliCloudRocketmqAccountDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"account_status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"DISABLE", "ENABLE"}, false),
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"password": {
				Type:     schema.TypeString,
				Required: true,
			},
			"username": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudRocketmqAccountCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	instanceId := d.Get("instance_id")
	action := fmt.Sprintf("/instances/%s/accounts", instanceId)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["username"] = d.Get("username")
	request["password"] = d.Get("password")
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_rocketmq_account", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", instanceId, request["username"]))

	return resourceAliCloudRocketmqAccountUpdate(d, meta)
}

func resourceAliCloudRocketmqAccountRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rocketmqServiceV2 := RocketmqServiceV2{client}

	objectRaw, err := rocketmqServiceV2.DescribeRocketmqAccount(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_rocketmq_account DescribeRocketmqAccount Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("account_status", objectRaw["accountStatus"])
	d.Set("password", objectRaw["password"])
	d.Set("username", objectRaw["username"])

	parts := strings.Split(d.Id(), ":")
	d.Set("instance_id", parts[0])

	return nil
}

func resourceAliCloudRocketmqAccountUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var body map[string]interface{}
	update := false

	var err error
	parts := strings.Split(d.Id(), ":")
	instanceId := parts[0]
	username := parts[1]
	action := fmt.Sprintf("/instances/%s/accounts/%s", instanceId, username)
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})

	if d.HasChange("account_status") {
		update = true
	}
	if v, ok := d.GetOk("account_status"); ok {
		query["accountStatus"] = StringPointer(v.(string))
	}

	if !d.IsNewResource() && d.HasChange("password") {
		update = true
	}
	if v, ok := d.GetOk("password"); ok {
		query["password"] = StringPointer(v.(string))
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
	}

	return resourceAliCloudRocketmqAccountRead(d, meta)
}

func resourceAliCloudRocketmqAccountDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	instanceId := parts[0]
	username := parts[1]
	action := fmt.Sprintf("/instances/%s/accounts/%s", instanceId, username)
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

	return nil
}
