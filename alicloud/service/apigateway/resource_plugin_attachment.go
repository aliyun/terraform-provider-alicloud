package apigateway

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourcePluginAttachment() *schema.Resource {
	return &schema.Resource{
		Create: ResourcePluginAttachmentCreate,
		Read:   ResourcePluginAttachmentRead,
		Delete: ResourcePluginAttachmentDelete,
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

func ResourcePluginAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "AttachPlugin"
	request := make(map[string]interface{})
	var err error
	request["PluginId"] = d.Get("plugin_id")
	request["GroupId"] = d.Get("group_id")
	request["ApiId"] = d.Get("api_id")
	request["StageName"] = d.Get("stage_name")
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *retry.RetryError {
		response, err = client.RpcPost("CloudAPI", "2016-07-14", action, nil, request, false)
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_api_gateway_plugin_attachment", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["GroupId"], ":", request["ApiId"], ":", request["PluginId"], ":", request["StageName"]))

	return ResourcePluginAttachmentRead(d, meta)
}

func ResourcePluginAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudApiService := CloudApiService{client}

	_, err := cloudApiService.DescribeApiGatewayPluginAttachment(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_api_gateway_plugin_attachment DescribeApiGatewayPluginAttachment Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	parts, err := helper.ParseResourceId(d.Id(), 4)
	if err != nil {
		return errs.WrapError(err)
	}
	d.Set("group_id", parts[0])
	d.Set("api_id", parts[1])
	d.Set("plugin_id", parts[2])
	d.Set("stage_name", parts[3])

	return nil
}

func ResourcePluginAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := helper.ParseResourceId(d.Id(), 4)
	if err != nil {
		return errs.WrapError(err)
	}
	var response map[string]interface{}
	request := make(map[string]interface{})
	action := "DetachPlugin"
	request["RegionId"] = client.RegionId
	request["GroupId"] = parts[0]
	request["ApiId"] = parts[1]
	request["PluginId"] = parts[2]
	request["StageName"] = parts[3]

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *retry.RetryError {
		response, err = client.RpcPost("CloudAPI", "2016-07-14", action, nil, request, false)
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
