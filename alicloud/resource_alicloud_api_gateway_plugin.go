package alicloud

import (
	"fmt"
	"log"
	"regexp"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudApiGatewayPlugin() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudApiGatewayPluginCreate,
		Read:   resourceAlicloudApiGatewayPluginRead,
		Update: resourceAlicloudApiGatewayPluginUpdate,
		Delete: resourceAlicloudApiGatewayPluginDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 200),
			},
			"plugin_data": {
				Type:     schema.TypeString,
				Required: true,
			},
			"plugin_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[a-zA-Z][A-Za-z0-9_]{3,49}$`), "It can contain uppercase English letters, lowercase English letters, Chinese characters, numbers, and underscores (_). It must be 4 to 50 characters in length and cannot start with an underscore (_)"),
			},
			"plugin_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"backendSignature", "caching", "cors", "ipControl", "jwtAuth", "trafficControl"}, false),
			},
			"tags": tagsSchema(),
		},
	}
}

func resourceAlicloudApiGatewayPluginCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreatePlugin"
	request := make(map[string]interface{})
	conn, err := client.NewApigatewayClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	request["PluginData"] = d.Get("plugin_data")
	request["PluginName"] = d.Get("plugin_name")
	request["PluginType"] = d.Get("plugin_type")
	if v, ok := d.GetOk("tags"); ok {
		count := 1
		for key, value := range v.(map[string]interface{}) {
			request[fmt.Sprintf("Tag.%d.Key", count)] = key
			request[fmt.Sprintf("Tag.%d.Value", count)] = value
			count++
		}
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-07-14"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_api_gateway_plugin", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["PluginId"]))

	return resourceAlicloudApiGatewayPluginRead(d, meta)
}
func resourceAlicloudApiGatewayPluginRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudApiService := CloudApiService{client}
	object, err := cloudApiService.DescribeApiGatewayPlugin(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_api_gateway_plugin cloudApiService.DescribeApiGatewayPlugin Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("description", object["Description"])
	d.Set("plugin_data", object["PluginData"])
	d.Set("plugin_name", object["PluginName"])
	d.Set("plugin_type", object["PluginType"])
	if tag, ok := object["Tags"]; ok {
		if v, ok := tag.(map[string]interface{}); ok {
			d.Set("tags", tagsToMap(v["TagInfo"]))
		}
	}

	return nil
}
func resourceAlicloudApiGatewayPluginUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudApiService := CloudApiService{client}
	conn, err := client.NewApigatewayClient()
	if err != nil {
		return WrapError(err)
	}
	var response map[string]interface{}
	d.Partial(true)

	if err := cloudApiService.SetResourceTags(d, "plugin"); err != nil {
		return WrapError(err)
	}
	update := false
	request := map[string]interface{}{
		"PluginId": d.Id(),
	}
	if d.HasChange("description") {
		update = true
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	if d.HasChange("plugin_data") {
		update = true
	}
	request["PluginData"] = d.Get("plugin_data")
	if !d.IsNewResource() && d.HasChange("plugin_name") {
		update = true
	}
	request["PluginName"] = d.Get("plugin_name")
	if update {
		action := "ModifyPlugin"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-07-14"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("description")
		d.SetPartial("plugin_data")
		d.SetPartial("plugin_name")
	}
	d.Partial(false)
	return resourceAlicloudApiGatewayPluginRead(d, meta)
}
func resourceAlicloudApiGatewayPluginDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeletePlugin"
	var response map[string]interface{}
	conn, err := client.NewApigatewayClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"PluginId": d.Id(),
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-07-14"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		if IsExpectedErrors(err, []string{"500"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}
