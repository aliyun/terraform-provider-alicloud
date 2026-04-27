package alicloud

import (
	"fmt"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudApiGatewayGroupPluginAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudApiGatewayGroupPluginAttachmentCreate,
		Read:   resourceAlicloudApiGatewayGroupPluginAttachmentRead,
		Delete: resourceAlicloudApiGatewayGroupPluginAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"plugin_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"stage_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudApiGatewayGroupPluginAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "AttachGroupPlugin"
	request := make(map[string]interface{})
	var err error
	request["GroupId"] = d.Get("group_id")
	request["PluginId"] = d.Get("plugin_id")
	request["StageName"] = d.Get("stage_name")
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = client.RpcPost("CloudAPI", "2016-07-14", action, nil, request, false)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_api_gateway_group_plugin_attachment", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["GroupId"], ":", request["PluginId"], ":", request["StageName"]))

	return resourceAlicloudApiGatewayGroupPluginAttachmentRead(d, meta)
}

func resourceAlicloudApiGatewayGroupPluginAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	apiGatewayServiceV2 := ApiGatewayServiceV2{client}

	objectRaw, err := apiGatewayServiceV2.DescribeApiGatewayGroupPluginAttachment(d.Id())

	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	err = d.Set("group_id", objectRaw["GroupId"].(string))
	if err != nil {
		return WrapError(err)
	}
	err = d.Set("stage_name", objectRaw["StageName"].(string))
	if err != nil {
		return WrapError(err)
	}
	err = d.Set("plugin_id", objectRaw["PluginId"].(string))
	if err != nil {
		return WrapError(err)
	}
	return nil

}

func resourceAlicloudApiGatewayGroupPluginAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return WrapError(err)
	}
	var response map[string]interface{}
	request := make(map[string]interface{})
	action := "DetachGroupPlugin"
	request["RegionId"] = client.RegionId
	request["GroupId"] = parts[0]
	request["PluginId"] = parts[1]
	request["StageName"] = parts[2]

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = client.RpcPost("CloudAPI", "2016-07-14", action, nil, request, false)
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
