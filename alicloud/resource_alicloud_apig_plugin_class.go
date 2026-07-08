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

func resourceAliCloudApigPluginClass() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudApigPluginClassCreate,
		Read:   resourceAliCloudApigPluginClassRead,
		Delete: resourceAliCloudApigPluginClassDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"alias": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"document": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"execute_priority": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: IntAtLeast(1),
			},
			"execute_stage": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"plugin_class_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"supported_min_gateway_version": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"version": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"version_description": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"wasm_language": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"wasm_url": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudApigPluginClassCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := fmt.Sprintf("/v1/plugin-classes")
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	request["wasmUrl"] = d.Get("wasm_url")
	if v, ok := d.GetOk("alias"); ok {
		request["alias"] = v
	}
	if v, ok := d.GetOk("supported_min_gateway_version"); ok {
		request["supportedMinGatewayVersion"] = v
	}
	request["description"] = d.Get("description")
	request["wasmLanguage"] = d.Get("wasm_language")
	request["name"] = d.Get("plugin_class_name")
	request["executePriority"] = d.Get("execute_priority")
	request["versionDescription"] = d.Get("version_description")
	request["executeStage"] = d.Get("execute_stage")
	request["version"] = d.Get("version")
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_apig_plugin_class", action, AlibabaCloudSdkGoERROR)
	}

	id, _ := jsonpath.Get("$.data.pluginClassId", response)
	d.SetId(fmt.Sprint(id))

	return resourceAliCloudApigPluginClassRead(d, meta)
}

func resourceAliCloudApigPluginClassRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	apigServiceV2 := ApigServiceV2{client}

	objectRaw, err := apigServiceV2.DescribeApigPluginClass(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_apig_plugin_class DescribeApigPluginClass Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("alias", objectRaw["alias"])
	d.Set("description", objectRaw["description"])
	d.Set("document", objectRaw["document"])
	d.Set("plugin_class_name", objectRaw["name"])
	d.Set("status", objectRaw["publishStatus"])
	d.Set("supported_min_gateway_version", objectRaw["supportedMinGatewayVersion"])
	d.Set("type", objectRaw["type"])
	d.Set("version", objectRaw["version"])
	d.Set("wasm_language", objectRaw["wasmLanguage"])

	return nil
}

func resourceAliCloudApigPluginClassDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	pluginClassId := d.Id()
	action := fmt.Sprintf("/v1/plugin-classes/%s", pluginClassId)
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
