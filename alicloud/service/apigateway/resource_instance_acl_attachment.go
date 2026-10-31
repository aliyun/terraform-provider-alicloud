package apigateway

import (
	"fmt"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/validators"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceInstanceAclAttachment() *schema.Resource {
	return &schema.Resource{
		Create: ResourceInstanceAclAttachmentCreate,
		Read:   ResourceInstanceAclAttachmentRead,
		Delete: ResourceInstanceAclAttachmentDelete,
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
				ValidateFunc: validators.StringInSlice([]string{"white", "black"}, true),
			},
		},
	}
}

func ResourceInstanceAclAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
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

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("CloudAPI", "2016-07-14", action, query, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_api_gateway_instance_acl_attachment", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%s:%s:%s", d.Get("instance_id").(string), d.Get("acl_id").(string), d.Get("acl_type").(string)))
	return ResourceInstanceAclAttachmentRead(d, meta)
}

func ResourceInstanceAclAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	apiGatewayServiceV2 := ApiGatewayServiceV2{client}

	objectRaw, err := apiGatewayServiceV2.DescribeApiGatewayInstanceAclAttachmentAttribute(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	err = d.Set("instance_id", objectRaw["InstanceId"].(string))
	if err != nil {
		return errs.WrapError(err)
	}
	err = d.Set("acl_id", objectRaw["AclId"].(string))
	if err != nil {
		return errs.WrapError(err)
	}
	err = d.Set("acl_type", objectRaw["AclType"].(string))
	if err != nil {
		return errs.WrapError(err)
	}
	return nil
}

func ResourceInstanceAclAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DisableInstanceAccessControl"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["InstanceId"] = d.Get("instance_id")
	request["AclId"] = d.Get("acl_id")

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("CloudAPI", "2016-07-14", action, query, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}
	return nil
}
