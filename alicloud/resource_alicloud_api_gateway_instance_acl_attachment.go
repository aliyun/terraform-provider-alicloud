package alicloud

import (
	"fmt"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudApiGatewayInstanceAclAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudApiGatewayInstanceAclAttachmentCreate,
		Read:   resourceAliCloudApiGatewayInstanceAclAttachmentRead,
		Delete: resourceAliCloudApiGatewayInstanceAclAttachmentDelete,
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
			"acl_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"acl_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"white", "black"}, true),
			},
		},
	}
}

func resourceAliCloudApiGatewayInstanceAclAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "EnableInstanceAccessControl"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["AclId"] = d.Get("acl_id")
	request["InstanceId"] = d.Get("instance_id")
	request["AclType"] = d.Get("acl_type")

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("CloudAPI", "2016-07-14", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_api_gateway_instance_acl_attachment", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%s:%s:%s", d.Get("instance_id").(string), d.Get("acl_id").(string), d.Get("acl_type").(string)))
	return resourceAliCloudApiGatewayInstanceAclAttachmentRead(d, meta)
}

func resourceAliCloudApiGatewayInstanceAclAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	apiGatewayServiceV2 := ApiGatewayServiceV2{client}

	objectRaw, err := apiGatewayServiceV2.DescribeApiGatewayInstanceAclAttachmentAttribute(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	err = d.Set("instance_id", objectRaw["InstanceId"].(string))
	if err != nil {
		return WrapError(err)
	}
	err = d.Set("acl_id", objectRaw["AclId"].(string))
	if err != nil {
		return WrapError(err)
	}
	err = d.Set("acl_type", objectRaw["AclType"].(string))
	if err != nil {
		return WrapError(err)
	}
	return nil
}

func resourceAliCloudApiGatewayInstanceAclAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DisableInstanceAccessControl"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["InstanceId"] = d.Get("instance_id")
	request["AclId"] = d.Get("acl_id")

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("CloudAPI", "2016-07-14", action, query, request, true)
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
	return nil
}
