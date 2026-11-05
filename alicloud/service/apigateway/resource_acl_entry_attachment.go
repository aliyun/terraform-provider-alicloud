package apigateway

import (
	"fmt"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceAclEntryAttachment() *schema.Resource {
	return &schema.Resource{
		Create: ResourceAclEntryAttachmentCreate,
		Read:   ResourceAclEntryAttachmentRead,
		Delete: ResourceAclEntryAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"acl_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"entry": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"comment": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func ResourceAclEntryAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "AddAccessControlListEntry"
	var response map[string]interface{}
	var err error
	request := map[string]interface{}{
		"AclId": d.Get("acl_id"),
	}

	aclMaps := make([]map[string]interface{}, 1)
	aclEntry := make(map[string]interface{}, 2)
	aclEntry["entry"] = d.Get("entry")
	if v, ok := d.GetOk("comment"); ok {
		aclEntry["comment"] = v
	}
	aclMaps = append(aclMaps, aclEntry)
	aclEntriesJSON, err := convert.ConvertListMapToJsonString(aclMaps)

	if err != nil {
		return errs.WrapError(err)
	}
	request["AclEntrys"] = aclEntriesJSON
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("CloudAPI", "2016-07-14", action, nil, request, true)
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_api_gateway_acl_entry_attachment", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["AclId"], ":", aclEntry["entry"]))
	return ResourceAclEntryAttachmentRead(d, meta)
}

func ResourceAclEntryAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	apiGatewayServiceV2 := ApiGatewayServiceV2{client}

	objectRaw, err := apiGatewayServiceV2.DescribeApiGatewayAclEntryAttachmentAttribute(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}
	parts, err := helper.ParseResourceId(d.Id(), 2)
	if err != nil {
		return errs.WrapError(err)
	}
	err = d.Set("acl_id", parts[0])
	if err != nil {
		return errs.WrapError(err)
	}
	err = d.Set("entry", objectRaw["AclEntryIp"])
	if err != nil {
		return errs.WrapError(err)
	}
	err = d.Set("comment", objectRaw["AclEntryComment"])
	if err != nil {
		return errs.WrapError(err)
	}
	return nil
}

func ResourceAclEntryAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "RemoveAccessControlListEntry"
	var response map[string]interface{}
	var err error
	parts, err := helper.ParseResourceId(d.Id(), 2)
	request := map[string]interface{}{
		"AclId": parts[0],
	}

	aclMaps := make([]map[string]interface{}, 1)
	aclMaps = append(aclMaps, map[string]interface{}{
		"entry": parts[1],
	})
	aclEntriesJSON, err := convert.ConvertListMapToJsonString(aclMaps)
	if err != nil {
		return errs.WrapError(err)
	}
	request["AclEntrys"] = aclEntriesJSON

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("CloudAPI", "2016-07-14", action, nil, request, true)
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
