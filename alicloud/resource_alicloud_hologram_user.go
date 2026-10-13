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

func resourceAliCloudHologramUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudHologramUserCreate,
		Read:   resourceAliCloudHologramUserRead,
		Delete: resourceAliCloudHologramUserDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"super_user": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"user_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudHologramUserCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	instanceId := d.Get("instance_id")
	action := fmt.Sprintf("/api/v1/instances/%s/createUser", instanceId)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("user_name"); ok {
		request["userName"] = v
	}

	if v, ok := d.GetOkExists("super_user"); ok {
		request["superUser"] = v
	}
	body = request
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RoaPost("Hologram", "2022-06-01", action, query, nil, body, true)
		if err != nil {
			if NeedRetry(err) || strings.Contains(err.Error(), "Failed to establish a connection to the execution service") {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_hologram_user", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", instanceId, request["userName"]))

	return resourceAliCloudHologramUserRead(d, meta)
}

func resourceAliCloudHologramUserRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	hologramServiceV2 := HologramServiceV2{client}

	objectRaw, err := hologramServiceV2.DescribeHologramUser(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_hologram_user DescribeHologramUser Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("super_user", objectRaw["superUser"])
	d.Set("user_name", objectRaw["userName"])

	parts := strings.Split(d.Id(), ":")
	d.Set("instance_id", parts[0])

	return nil
}

func resourceAliCloudHologramUserDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	instanceId := parts[0]
	action := fmt.Sprintf("/api/v1/instances/%s/dropUser", instanceId)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["userName"] = parts[1]

	body = request
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RoaPost("Hologram", "2022-06-01", action, query, nil, body, true)
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
