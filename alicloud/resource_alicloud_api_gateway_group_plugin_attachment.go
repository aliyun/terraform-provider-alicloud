package alicloud

import (
	"fmt"
	"github.com/PaesslerAG/jsonpath"
	"log"
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
	var response map[string]interface{}

	action := "DescribePluginsByGroup"
	request := make(map[string]interface{})
	var err error
	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return WrapErrorf(fmt.Errorf("id is not a expe"), DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	request["GroupId"] = parts[0]
	request["StageName"] = parts[2]
	request["PageSize"] = 100
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

	// 获取数据
	v, err := jsonpath.Get("$.Plugins.PluginAttribute[*]", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, d.Id(), "$.Plugins.PluginAttribute[*]", response)
	}

	// 检查数据是否存在
	if len(v.([]interface{})) == 0 {
		if !d.IsNewResource() {
			log.Printf("[DEBUG] Resource not found, removing from state")
			d.SetId("")
			return nil
		}
		return WrapErrorf(NotFoundErr("GroupPluginAttachment", d.Id()), NotFoundMsg, response)
	}

	// 处理查询到的数据
	plugins := v.([]interface{})

	// 找到匹配的插件（根据 plugin_id）
	pluginId := parts[1]

	var foundPlugin map[string]interface{}

	for _, plugin := range plugins {
		pluginMap := plugin.(map[string]interface{})
		if pluginMap["PluginId"].(string) == pluginId {
			foundPlugin = pluginMap
			break
		}
	}

	if foundPlugin == nil {
		if !d.IsNewResource() {
			log.Printf("[DEBUG] Plugin attachment not found, removing from state")
			d.SetId("")
			return nil
		}
		return WrapErrorf(NotFoundErr("GroupPluginAttachment", d.Id()), NotFoundMsg, response)
	}

	d.Set("group_id", parts[0])   // 从 ID 中获取
	d.Set("stage_name", parts[2]) // 从 ID 中获取 (假设 parts[2] 是 stage_name)
	d.Set("plugin_id", parts[1])  // 从 ID 中获取 (假设 parts[1] 是 plugin_id)

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
