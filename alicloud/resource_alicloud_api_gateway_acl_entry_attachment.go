package alicloud

import (
	"fmt"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudApiGatewayAclEntryAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudApiGatewayAclEntryAttachmentCreate,
		Read:   resourceAliCloudApiGatewayAclEntryAttachmentRead,
		Delete: resourceAliCloudApiGatewayAclEntryAttachmentDelete,
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

func resourceAliCloudApiGatewayAclEntryAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
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
	aclEntriesJSON, err := convertListMapToJsonString(aclMaps)

	if err != nil {
		return WrapError(err)
	}
	request["AclEntrys"] = aclEntriesJSON
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("CloudAPI", "2016-07-14", action, nil, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_api_gateway_acl_entry_attachment", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["AclId"], ":", aclEntry["entry"]))
	return resourceAliCloudApiGatewayAclEntryAttachmentRead(d, meta)
}

func resourceAliCloudApiGatewayAclEntryAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	apiGatewayServiceV2 := ApiGatewayServiceV2{client}

	objectRaw, err := apiGatewayServiceV2.DescribeApiGatewayAclEntryAttachmentAttribute(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	err = d.Set("acl_id", parts[0])
	if err != nil {
		return WrapError(err)
	}
	err = d.Set("entry", objectRaw["AclEntryIp"])
	if err != nil {
		return WrapError(err)
	}
	err = d.Set("comment", objectRaw["AclEntryComment"])
	if err != nil {
		return WrapError(err)
	}
	return nil
}

func resourceAliCloudApiGatewayAclEntryAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "RemoveAccessControlListEntry"
	var response map[string]interface{}
	var err error
	parts, err := ParseResourceId(d.Id(), 2)
	request := map[string]interface{}{
		"AclId": parts[0],
	}

	aclMaps := make([]map[string]interface{}, 1)
	aclMaps = append(aclMaps, map[string]interface{}{
		"entry": parts[1],
	})
	aclEntriesJSON, err := convertListMapToJsonString(aclMaps)
	if err != nil {
		return WrapError(err)
	}
	request["AclEntrys"] = aclEntriesJSON

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("CloudAPI", "2016-07-14", action, nil, request, true)
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
