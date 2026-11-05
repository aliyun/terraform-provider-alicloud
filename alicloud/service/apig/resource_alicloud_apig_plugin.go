// Package apig. This file is generated automatically. Please do not modify it manually, thank you!
package apig

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourcePlugin() *schema.Resource {
	return &schema.Resource{
		Create: resourcePluginCreate,
		Read:   resourcePluginRead,
		Delete: resourcePluginDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"gateway_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"plugin_class_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"gateway_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"plugin_class_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourcePluginCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := fmt.Sprintf("/v1/plugins/")
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	request["pluginClassId"] = d.Get("plugin_class_id")
	if v, ok := d.GetOk("gateway_id"); ok {
		localData, _ := jsonpath.Get("$", v)
		gatewayIdsMapsArray := convert.ConvertToInterfaceArray(localData)

		request["gatewayIds"] = gatewayIdsMapsArray
	}

	body = request
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RoaPost("APIG", "2024-03-27", action, query, nil, body, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_apig_plugin", action, errs.AlibabaCloudSdkGoERROR)
	}

	id, _ := jsonpath.Get("$.data.installPluginResults[0].pluginId", response)
	d.SetId(fmt.Sprint(id))

	return resourcePluginRead(d, meta)
}

func resourcePluginRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	apigServiceV2 := NewService(client)

	objectRaw, err := apigServiceV2.DescribeApigPlugin(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_apig_plugin DescribeApigPlugin Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	if gatewayInfo, ok := objectRaw["gatewayInfo"]; ok {
		if gatewayInfoMap, ok := gatewayInfo.(map[string]interface{}); ok {
			d.Set("gateway_name", gatewayInfoMap["name"])
			d.Set("gateway_id", gatewayInfoMap["gatewayId"])
		}
	}

	if pluginClassInfo, ok := objectRaw["pluginClassInfo"]; ok {
		if pluginClassInfoMap, ok := pluginClassInfo.(map[string]interface{}); ok {
			d.Set("plugin_class_name", pluginClassInfoMap["name"])
			d.Set("plugin_class_id", pluginClassInfoMap["pluginClassId"])
		}
	}

	return nil
}

func resourcePluginDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	pluginId := d.Id()
	action := fmt.Sprintf("/v1/plugins/%s", pluginId)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	var err error
	request = make(map[string]interface{})

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RoaDelete("APIG", "2024-03-27", action, query, nil, nil, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		if errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	return nil
}
