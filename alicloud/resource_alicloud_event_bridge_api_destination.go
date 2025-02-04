package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudEventBridgeApiDestination() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudEventBridgeApiDestinationCreate,
		Read:   resourceAliCloudEventBridgeApiDestinationRead,
		Update: resourceAliCloudEventBridgeApiDestinationUpdate,
		Delete: resourceAliCloudEventBridgeApiDestinationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"connection_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"api_destination_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"http_api_parameters": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"endpoint": {
							Type:     schema.TypeString,
							Required: true,
						},
						"method": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: StringInSlice([]string{"GET", "POST", "HEAD", "DELETE", "PUT", "PATCH"}, false),
						},
					},
				},
			},
			"create_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudEventBridgeApiDestinationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateApiDestination"
	request := make(map[string]interface{})
	var err error

	request["ConnectionName"] = d.Get("connection_name")
	request["ApiDestinationName"] = d.Get("api_destination_name")

	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}

	httpApiParameters := d.Get("http_api_parameters")
	httpApiParametersMap := map[string]interface{}{}
	for _, httpApiParametersList := range httpApiParameters.([]interface{}) {
		httpApiParametersArg := httpApiParametersList.(map[string]interface{})

		httpApiParametersMap["Endpoint"] = httpApiParametersArg["endpoint"]
		httpApiParametersMap["Method"] = httpApiParametersArg["method"]
	}

	httpApiParametersJson, err := convertMaptoJsonString(httpApiParametersMap)
	if err != nil {
		return WrapError(err)
	}

	request["HttpApiParameters"] = httpApiParametersJson
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = client.RpcPost("eventbridge", "2020-04-01", action, nil, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_event_bridge_api_destination", action, AlibabaCloudSdkGoERROR)
	}

	if resp, err := jsonpath.Get("$.Date", response); err != nil || resp == nil {
		return WrapErrorf(err, IdMsg, "alicloud_event_bridge_api_destination")
	} else {
		apiDestinationName := resp.(map[string]interface{})["ApiDestinationName"]
		d.SetId(fmt.Sprint(apiDestinationName))
	}

	return resourceAliCloudEventBridgeApiDestinationRead(d, meta)
}

func resourceAliCloudEventBridgeApiDestinationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	eventBridgeServiceV2 := EventBridgeServiceV2{client}

	object, err := eventBridgeServiceV2.DescribeEventBridgeApiDestination(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("connection_name", object["ConnectionName"])
	d.Set("api_destination_name", object["ApiDestinationName"])
	d.Set("description", object["Description"])
	d.Set("create_time", object["GmtCreate"])

	if httpApiParameters, ok := object["HttpApiParameters"]; ok {
		httpApiParametersMaps := make([]map[string]interface{}, 0)
		httpApiParametersArg := httpApiParameters.(map[string]interface{})
		httpApiParametersMap := make(map[string]interface{})

		if endpoint, ok := httpApiParametersArg["Endpoint"]; ok {
			httpApiParametersMap["endpoint"] = endpoint
		}

		if method, ok := httpApiParametersArg["Method"]; ok {
			httpApiParametersMap["method"] = method
		}

		httpApiParametersMaps = append(httpApiParametersMaps, httpApiParametersMap)

		d.Set("http_api_parameters", httpApiParametersMaps)
	}

	return nil
}

func resourceAliCloudEventBridgeApiDestinationUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	update := false

	request := map[string]interface{}{
		"ApiDestinationName": d.Id(),
	}

	if d.HasChange("description") {
		update = true
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}

	if d.HasChange("http_api_parameters") {
		update = true
	}
	httpApiParameters := d.Get("http_api_parameters")
	httpApiParametersMap := map[string]interface{}{}
	for _, httpApiParametersList := range httpApiParameters.([]interface{}) {
		httpApiParametersArg := httpApiParametersList.(map[string]interface{})

		httpApiParametersMap["Endpoint"] = httpApiParametersArg["endpoint"]
		httpApiParametersMap["Method"] = httpApiParametersArg["method"]
	}

	httpApiParametersJson, err := convertMaptoJsonString(httpApiParametersMap)
	if err != nil {
		return WrapError(err)
	}

	request["HttpApiParameters"] = httpApiParametersJson

	if update {
		action := "UpdateApiDestination"
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.RpcPost("eventbridge", "2020-04-01", action, nil, request, true)
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
	}

	return resourceAliCloudEventBridgeApiDestinationRead(d, meta)
}

func resourceAliCloudEventBridgeApiDestinationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteApiDestination"
	var response map[string]interface{}
	var err error

	request := map[string]interface{}{
		"ApiDestinationName": d.Id(),
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = client.RpcPost("eventbridge", "2020-04-01", action, nil, request, true)
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

	return nil
}
