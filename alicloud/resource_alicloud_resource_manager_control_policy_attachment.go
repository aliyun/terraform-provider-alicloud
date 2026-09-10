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

func resourceAliCloudResourceManagerControlPolicyAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudResourceManagerControlPolicyAttachmentCreate,
		Read:   resourceAliCloudResourceManagerControlPolicyAttachmentRead,
		Delete: resourceAliCloudResourceManagerControlPolicyAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"policy_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"target_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudResourceManagerControlPolicyAttachmentCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "AttachControlPolicy"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("target_id"); ok {
		request["TargetId"] = v
	}
	if v, ok := d.GetOk("policy_id"); ok {
		request["PolicyId"] = v
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("ResourceManager", "2020-03-31", action, query, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"ConcurrentCallNotSupported"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_resource_manager_control_policy_attachment", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", request["PolicyId"], request["TargetId"]))

	resourceManagerServiceV2 := ResourceManagerServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"Custom"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, resourceManagerServiceV2.ResourceManagerControlPolicyAttachmentStateRefreshFunc(d.Id(), "$.PolicyType", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudResourceManagerControlPolicyAttachmentRead(d, meta)
}

func resourceAliCloudResourceManagerControlPolicyAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	resourceManagerServiceV2 := ResourceManagerServiceV2{client}

	objectRaw, err := resourceManagerServiceV2.DescribeResourceManagerControlPolicyAttachment(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_resource_manager_control_policy_attachment DescribeResourceManagerControlPolicyAttachment Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("policy_id", objectRaw["PolicyId"])

	parts := strings.Split(d.Id(), ":")
	d.Set("target_id", parts[1])

	return nil
}

func resourceAliCloudResourceManagerControlPolicyAttachmentDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DetachControlPolicy"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["TargetId"] = parts[1]
	request["PolicyId"] = parts[0]

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("ResourceManager", "2020-03-31", action, query, request, true)

		if err != nil {
			if IsExpectedErrors(err, []string{"ConcurrentCallNotSupported"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		if IsExpectedErrors(err, []string{"EntityNotExists.Target"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
