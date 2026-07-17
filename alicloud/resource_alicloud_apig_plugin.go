// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAliCloudApigPlugin() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudApigPluginCreate,
		Read:   resourceAliCloudApigPluginRead,
		Delete: resourceAliCloudApigPluginDelete,
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
			"gateway_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"plugin_class_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"plugin_class_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudApigPluginCreate(d *schema.ResourceData, meta interface{}) error {

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
		gatewayIdsMapsArray := convertToInterfaceArray(localData)

		request["gatewayIds"] = gatewayIdsMapsArray
	}

	body = request
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RoaPost("APIG", "2024-03-27", action, query, nil, body, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_apig_plugin", action, AlibabaCloudSdkGoERROR)
	}

	id, _ := jsonpath.Get("$.data.installPluginResults[0].pluginId", response)
	d.SetId(fmt.Sprint(id))

	return resourceAliCloudApigPluginRead(d, meta)
}

func resourceAliCloudApigPluginRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	apigServiceV2 := ApigServiceV2{client}

	objectRaw, err := apigServiceV2.DescribeApigPlugin(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_apig_plugin DescribeApigPlugin Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	gatewayInfoRawObj, _ := jsonpath.Get("$.gatewayInfo", objectRaw)
	gatewayInfoRaw := make(map[string]interface{})
	if gatewayInfoRawObj != nil {
		gatewayInfoRaw = gatewayInfoRawObj.(map[string]interface{})
	}
	d.Set("gateway_name", gatewayInfoRaw["name"])
	d.Set("gateway_id", gatewayInfoRaw["gatewayId"])

	pluginClassInfoRawObj, _ := jsonpath.Get("$.pluginClassInfo", objectRaw)
	pluginClassInfoRaw := make(map[string]interface{})
	if pluginClassInfoRawObj != nil {
		pluginClassInfoRaw = pluginClassInfoRawObj.(map[string]interface{})
	}
	d.Set("plugin_class_name", pluginClassInfoRaw["name"])
	d.Set("plugin_class_id", pluginClassInfoRaw["pluginClassId"])

	return nil
}

func resourceAliCloudApigPluginDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	pluginId := d.Id()
	action := fmt.Sprintf("/v1/plugins/%s", pluginId)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	var err error
	request = make(map[string]interface{})

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RoaDelete("APIG", "2024-03-27", action, query, nil, nil, true)
		if err != nil {
			if NeedRetry(err) {
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
