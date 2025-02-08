package alicloud

import (
	"fmt"
	"log"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"time"
)

func resourceAlicloudApiGatewayPluginAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudApiGatewayPluginAttachmentCreate,
		Read:   resourceAlicloudApiGatewayPluginAttachmentRead,
		Delete: resourceAlicloudApiGatewayPluginAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"plugin_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"api_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"stage_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"PRE", "RELEASE", "TEST"}, false),
			},
		},
	}
}

func resourceAlicloudApiGatewayPluginAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "AttachPlugin"
	request := make(map[string]interface{})
	var err error
	request["PluginId"] = d.Get("plugin_id")
	request["GroupId"] = d.Get("group_id")
	request["ApiId"] = d.Get("api_id")
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_api_gateway_plugin_attachment", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["GroupId"], ":", request["ApiId"], ":", request["PluginId"], ":", request["StageName"]))

	return resourceAlicloudApiGatewayPluginAttachmentRead(d, meta)
}

func resourceAlicloudApiGatewayPluginAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudApiService := CloudApiService{client}

	_, err := cloudApiService.DescribeApiGatewayPluginAttachment(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_api_gateway_plugin_attachment DescribeApiGatewayPluginAttachment Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	parts, err := ParseResourceId(d.Id(), 4)
	if err != nil {
		return WrapError(err)
	}
	d.Set("group_id", parts[0])
	d.Set("api_id", parts[1])
	d.Set("plugin_id", parts[2])
	d.Set("stage_name", parts[3])

	return nil
}

func resourceAlicloudApiGatewayPluginAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := ParseResourceId(d.Id(), 4)
	if err != nil {
		return WrapError(err)
	}
	var response map[string]interface{}
	request := make(map[string]interface{})
	action := "DetachPlugin"
	request["RegionId"] = client.RegionId
	request["GroupId"] = parts[0]
	request["ApiId"] = parts[1]
	request["PluginId"] = parts[2]
	request["StageName"] = parts[3]

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
