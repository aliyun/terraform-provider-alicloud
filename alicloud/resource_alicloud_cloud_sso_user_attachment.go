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

func resourceAliCloudCloudSSOUserAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudCloudSSOUserAttachmentCreate,
		Read:   resourceAliCloudCloudSSOUserAttachmentRead,
		Delete: resourceAliCloudCloudSSOUserAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"directory_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"user_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudCloudSSOUserAttachmentCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "AddUserToGroup"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["GroupId"] = d.Get("group_id")
	request["UserId"] = d.Get("user_id")
	request["DirectoryId"] = d.Get("directory_id")

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("cloudsso", "2021-05-15", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cloud_sso_user_attachment", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v:%v", request["DirectoryId"], request["GroupId"], request["UserId"]))

	return resourceAliCloudCloudSSOUserAttachmentRead(d, meta)
}

func resourceAliCloudCloudSSOUserAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudSSOServiceV2 := CloudSSOServiceV2{client}

	objectRaw, err := cloudSSOServiceV2.DescribeCloudSSOUserAttachment(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cloud_sso_user_attachment DescribeCloudSSOUserAttachment Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("group_id", objectRaw["GroupId"])
	d.Set("user_id", objectRaw["UserId"])

	parts := strings.Split(d.Id(), ":")
	d.Set("directory_id", parts[0])

	return nil
}

func resourceAliCloudCloudSSOUserAttachmentDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "RemoveUserFromGroup"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["GroupId"] = parts[1]
	request["UserId"] = parts[2]
	request["DirectoryId"] = parts[0]

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("cloudsso", "2021-05-15", action, query, request, true)

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
