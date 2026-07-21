// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudApiGatewayPlugin() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudApiGatewayPluginCreate,
		Read:   resourceAliCloudApiGatewayPluginRead,
		Update: resourceAliCloudApiGatewayPluginUpdate,
		Delete: resourceAliCloudApiGatewayPluginDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringMatch(regexp.MustCompile("(.*)"), "The description of the plug-in, which cannot exceed 200 characters."),
			},
			"plugin_data": {
				Type:     schema.TypeString,
				Required: true,
			},
			"plugin_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: StringMatch(regexp.MustCompile("^[\u4E00-\u9FA5A-Za-z0-9_]+$"), "The name of the plug-in that you want to create. It can contain uppercase English letters, lowercase English letters, Chinese characters, numbers, and underscores (_). It must be 4 to 50 characters in length and cannot start with an underscore (_)."),
			},
			"plugin_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"trafficControl", "ipControl", "backendSignature", "jwtAuth", "basicAuth", "cors", "caching", "routing", "accessControl", "errorMapping", "circuitBreaker", "remoteAuth", "logMask", "transformer"}, false),
			},
			"tags": tagsSchema(),
		},
	}
}

func resourceAliCloudApiGatewayPluginCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreatePlugin"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	request["PluginType"] = d.Get("plugin_type")
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	request["PluginName"] = d.Get("plugin_name")
	request["PluginData"] = d.Get("plugin_data")
	if v, ok := d.GetOk("tags"); ok {
		tagsMap := ConvertTags(v.(map[string]interface{}))
		request = expandTagsToMap(request, tagsMap)
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("CloudAPI", "2016-07-14", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_api_gateway_plugin", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["PluginId"]))

	return resourceAliCloudApiGatewayPluginRead(d, meta)
}

func resourceAliCloudApiGatewayPluginRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	apiGatewayServiceV2 := ApiGatewayServiceV2{client}

	objectRaw, err := apiGatewayServiceV2.DescribeApiGatewayPlugin(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_api_gateway_plugin DescribeApiGatewayPlugin Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("create_time", objectRaw["CreatedTime"])
	d.Set("description", objectRaw["Description"])
	d.Set("plugin_data", objectRaw["PluginData"])
	d.Set("plugin_name", objectRaw["PluginName"])
	d.Set("plugin_type", objectRaw["PluginType"])

	tagsMaps, _ := jsonpath.Get("$.Tags.TagInfo", objectRaw)
	d.Set("tags", tagsToMap(tagsMaps))

	return nil
}

func resourceAliCloudApiGatewayPluginUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	action := "ModifyPlugin"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["PluginId"] = d.Id()
	if d.HasChange("description") {
		update = true
		request["Description"] = d.Get("description")
	}

	if d.HasChange("plugin_name") {
		update = true
	}
	request["PluginName"] = d.Get("plugin_name")
	if d.HasChange("plugin_data") {
		update = true
	}
	request["PluginData"] = d.Get("plugin_data")
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("CloudAPI", "2016-07-14", action, query, request, true)
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
	}

	if d.HasChange("tags") {
		apiGatewayServiceV2 := ApiGatewayServiceV2{client}
		if err := apiGatewayServiceV2.SetResourceTags(d, "plugin"); err != nil {
			return WrapError(err)
		}
	}
	return resourceAliCloudApiGatewayPluginRead(d, meta)
}

func resourceAliCloudApiGatewayPluginDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeletePlugin"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["PluginId"] = d.Id()

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("CloudAPI", "2016-07-14", action, query, request, true)

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
		if IsExpectedErrors(err, []string{"500"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
