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

func resourceAliCloudEcsRamRoleAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudEcsRamRoleAttachmentCreate,
		Read:   resourceAliCloudEcsRamRoleAttachmentRead,
		Update: resourceAliCloudEcsRamRoleAttachmentUpdate,
		Delete: resourceAliCloudEcsRamRoleAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"policy": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ram_role_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudEcsRamRoleAttachmentCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "AttachInstanceRamRole"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["RamRoleName"] = d.Get("ram_role_name")
	request["InstanceIds"] = fmt.Sprintf("[\"%s\"]", d.Get("instance_id"))

	request["RegionId"] = client.RegionId

	if v, ok := d.GetOk("policy"); ok {
		request["Policy"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Ecs", "2014-05-26", action, query, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"unexpected end of JSON input"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ecs_ram_role_attachment", action, AlibabaCloudSdkGoERROR)
	}

	AttachInstanceRamRoleResultsAttachInstanceRamRoleResultInstanceIdVar, _ := jsonpath.Get("$.AttachInstanceRamRoleResults.AttachInstanceRamRoleResult[0].InstanceId", response)
	d.SetId(fmt.Sprintf("%v:%v", AttachInstanceRamRoleResultsAttachInstanceRamRoleResultInstanceIdVar, response["RamRoleName"]))

	return resourceAliCloudEcsRamRoleAttachmentRead(d, meta)
}

func resourceAliCloudEcsRamRoleAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsServiceV2 := EcsServiceV2{client}

	objectRaw, err := ecsServiceV2.DescribeEcsRamRoleAttachment(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ecs_ram_role_attachment DescribeEcsRamRoleAttachment Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("instance_id", objectRaw["InstanceId"])
	d.Set("ram_role_name", objectRaw["RamRoleName"])

	return nil
}

func resourceAliCloudEcsRamRoleAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Cannot update resource Alicloud Resource Ram Role Attachment.")
	return nil
}

func resourceAliCloudEcsRamRoleAttachmentDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DetachInstanceRamRole"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["RamRoleName"] = parts[1]
	request["InstanceIds"] = fmt.Sprintf("[\"%s\"]", parts[0])
	request["RegionId"] = client.RegionId

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Ecs", "2014-05-26", action, query, request, true)

		if err != nil {
			if IsExpectedErrors(err, []string{"unexpected end of JSON input"}) || NeedRetry(err) {
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
