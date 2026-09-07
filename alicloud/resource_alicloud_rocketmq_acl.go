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

func resourceAliCloudRocketmqAcl() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudRocketmqAclCreate,
		Read:   resourceAliCloudRocketmqAclRead,
		Update: resourceAliCloudRocketmqAclUpdate,
		Delete: resourceAliCloudRocketmqAclDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"actions": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"decision": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: StringInSlice([]string{"Deny", "Allow"}, false),
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ip_whitelists": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"resource_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"resource_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"Group", "Topic"}, false),
			},
			"username": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudRocketmqAclCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	instanceId := d.Get("instance_id")
	username := d.Get("username")
	action := fmt.Sprintf("/instances/%s/acl/account/%s", instanceId, username)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["resourceType"] = d.Get("resource_type")
	request["resourceName"] = d.Get("resource_name")
	request["decision"] = d.Get("decision")
	if v, ok := d.GetOk("ip_whitelists"); ok {
		ipWhitelistsMapsArray := v.([]interface{})
		request["ipWhitelists"] = ipWhitelistsMapsArray
	}

	if v, ok := d.GetOk("actions"); ok {
		actionsMapsArray := v.([]interface{})
		request["actions"] = actionsMapsArray
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_rocketmq_acl", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v:%v:%v", instanceId, username, request["resourceType"], request["resourceName"]))

	return resourceAliCloudRocketmqAclRead(d, meta)
}

func resourceAliCloudRocketmqAclRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rocketmqServiceV2 := RocketmqServiceV2{client}

	objectRaw, err := rocketmqServiceV2.DescribeRocketmqAcl(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_rocketmq_acl DescribeRocketmqAcl Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("decision", objectRaw["decision"])
	d.Set("instance_id", objectRaw["instanceId"])
	d.Set("resource_name", objectRaw["resourceName"])
	d.Set("resource_type", objectRaw["resourceType"])
	d.Set("username", objectRaw["username"])

	actionsRaw := make([]interface{}, 0)
	if objectRaw["actions"] != nil {
		actionsRaw = objectRaw["actions"].([]interface{})
	}

	d.Set("actions", actionsRaw)
	ipWhitelistsRaw := make([]interface{}, 0)
	if objectRaw["ipWhitelists"] != nil {
		ipWhitelistsRaw = objectRaw["ipWhitelists"].([]interface{})
	}

	d.Set("ip_whitelists", ipWhitelistsRaw)

	return nil
}

func resourceAliCloudRocketmqAclUpdate(d *schema.ResourceData, meta interface{}) error {
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
	action := fmt.Sprintf("/instances/%s/acl/account/%s", instanceId, username)
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})
	request["resourceType"] = parts[2]
	request["resourceName"] = parts[3]

	if d.HasChange("ip_whitelists") {
		update = true
	}
	if v, ok := d.GetOk("ip_whitelists"); ok {
		ipWhitelistsMapsArray := v.([]interface{})
		request["ipWhitelists"] = ipWhitelistsMapsArray
	}

	if d.HasChange("actions") {
		update = true
	}
	if v, ok := d.GetOk("actions"); ok {
		actionsMapsArray := v.([]interface{})
		request["actions"] = actionsMapsArray
	}

	if d.HasChange("decision") {
		update = true
	}
	request["decision"] = d.Get("decision")
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

	return resourceAliCloudRocketmqAclRead(d, meta)
}

func resourceAliCloudRocketmqAclDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	instanceId := parts[0]
	username := parts[1]
	action := fmt.Sprintf("/instances/%s/acl/account/%s", instanceId, username)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	var err error
	request = make(map[string]interface{})
	query["resourceType"] = StringPointer(parts[2])
	query["resourceName"] = StringPointer(parts[3])

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
