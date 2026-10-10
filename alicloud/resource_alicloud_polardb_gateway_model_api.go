package alicloud

import (
	"fmt"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudPolarDBGatewayModelAPI() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudPolarDBGatewayModelAPICreate,
		Read:   resourceAlicloudPolarDBGatewayModelAPIRead,
		Update: resourceAlicloudPolarDBGatewayModelAPIUpdate,
		Delete: resourceAlicloudPolarDBGatewayModelAPIDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"model_api_id": {
				Type: schema.TypeString, Computed: true,
				Description: "The ID of the model API.",
			},
			"gateway_id": {
				Type: schema.TypeString, Required: true, ForceNew: true,
				Description: "The ID of the PolarDB AI gateway.",
			},
			"name": {
				Type: schema.TypeString, Required: true, ForceNew: true,
				Description: "The name of the model API.",
			},
			"model_category": {
				Type: schema.TypeString, Required: true,
				ValidateFunc: StringInSlice([]string{"text", "embedding", "rerank"}, false),
				Description:  "The model category. Valid values: `text`, `embedding`, `rerank`.",
			},
			"path_prefix": {
				Type: schema.TypeString, Required: true,
				Description: "The API path prefix.",
			},
			"protocol": {
				Type: schema.TypeString, Required: true,
				ValidateFunc: StringInSlice([]string{"openai", "anthropic", "bailian", "vllm"}, false),
				Description:  "The protocol of the model API.",
			},
			"record_input": {
				Type: schema.TypeString, Optional: true,
				Description: "Whether or how much request input is recorded for billing.",
			},
			"record_output": {
				Type: schema.TypeString, Optional: true,
				Description: "Whether or how much response output is recorded for billing.",
			},
			"route_rules": {
				Type: schema.TypeString, Required: true,
				ValidateFunc: validatePolarDBGatewayAIRouteRules,
				StateFunc:    func(v interface{}) string { return normalizePolarDBGatewayAIJSON(v) },
				Description:  "The routing rules as a JSON array string.",
			},
			"force_model": {
				Type: schema.TypeString, Optional: true, ForceNew: true,
				Description: "The model to which requests are forcibly routed.",
			},
			"invoke_endpoint": {
				Type: schema.TypeString, Computed: true,
				Description: "The invocation endpoint of the model API.",
			},
			"status": {
				Type: schema.TypeString, Computed: true,
				Description: "The status of the model API.",
			},
			"create_time": {
				Type: schema.TypeString, Computed: true,
				Description: "The time when the model API was created.",
			},
		},
	}
}

func resourceAlicloudPolarDBGatewayModelAPICreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "CreateModelApi"
	request := polarDBGatewayModelAPIRequest(d, client.RegionId, true)
	response, err := client.RpcPost(polarDBGatewayAIProduct, polarDBGatewayAIVersion, action, nil, request, false)
	debugPolarDBGatewayAI(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_polardb_gateway_model_api", action, AlibabaCloudSdkGoERROR)
	}
	modelAPIID := fmt.Sprint(response["ModelApiId"])
	if modelAPIID == "" || modelAPIID == "<nil>" {
		return WrapError(fmt.Errorf("CreateModelApi returned empty ModelApiId"))
	}
	d.SetId(polarDBGatewayAIChildID(d.Get("gateway_id").(string), modelAPIID))
	return resourceAlicloudPolarDBGatewayModelAPIRead(d, meta)
}

func resourceAlicloudPolarDBGatewayModelAPIRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gatewayID, modelAPIID, err := parsePolarDBGatewayAIChildID(d.Id())
	if err != nil {
		return WrapError(err)
	}
	object, err := findPolarDBGatewayAIItem(client, "DescribeModelApis", gatewayID, "ModelApiIds", "ModelApiId", modelAPIID)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "DescribeModelApis", AlibabaCloudSdkGoERROR)
	}
	if object == nil {
		d.SetId("")
		return nil
	}
	if err = d.Set("gateway_id", gatewayID); err != nil {
		return WrapError(err)
	}
	if err = d.Set("model_api_id", modelAPIID); err != nil {
		return WrapError(err)
	}
	if err = d.Set("name", object["Name"]); err != nil {
		return WrapError(err)
	}
	if err = d.Set("model_category", object["Category"]); err != nil {
		return WrapError(err)
	}
	if err = d.Set("path_prefix", object["PathPrefix"]); err != nil {
		return WrapError(err)
	}
	if err = d.Set("protocol", object["Protocol"]); err != nil {
		return WrapError(err)
	}
	// Some deployments omit these write-only values, or return an empty value,
	// after create/update. Keep the configured value in that case so refresh
	// does not introduce a perpetual diff.
	if value := fmt.Sprint(object["RecordInput"]); value != "" && value != "<nil>" {
		if err = d.Set("record_input", value); err != nil {
			return WrapError(err)
		}
	}
	if value := fmt.Sprint(object["RecordOutput"]); value != "" && value != "<nil>" {
		if err = d.Set("record_output", value); err != nil {
			return WrapError(err)
		}
	}
	if value, ok := object["ForceModel"]; ok && value != nil {
		if err = d.Set("force_model", value); err != nil {
			return WrapError(err)
		}
	}
	if err = d.Set("invoke_endpoint", object["InvokeEndpoint"]); err != nil {
		return WrapError(err)
	}
	if err = d.Set("status", object["Status"]); err != nil {
		return WrapError(err)
	}
	if err = d.Set("create_time", object["GmtCreated"]); err != nil {
		return WrapError(err)
	}
	if _, ok := object["Category"]; !ok {
		if value, exists := object["ModelCategory"]; exists && value != nil {
			if err = d.Set("model_category", value); err != nil {
				return WrapError(err)
			}
		}
	}
	// DescribeModelApis returns an expanded runtime representation of the route
	// rules and omits request-only fields. Preserve the configured request form;
	// only populate the field when importing an existing resource.
	if value, ok := object["RouteRules"]; ok && value != nil && d.Get("route_rules").(string) == "" {
		if err = d.Set("route_rules", normalizePolarDBGatewayAIJSON(value)); err != nil {
			return WrapError(err)
		}
	}
	return nil
}

func resourceAlicloudPolarDBGatewayModelAPIUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gatewayID, modelAPIID, err := parsePolarDBGatewayAIChildID(d.Id())
	if err != nil {
		return WrapError(err)
	}
	request := polarDBGatewayModelAPIRequest(d, client.RegionId, false)
	request["GwClusterId"] = gatewayID
	request["ModelApiId"] = modelAPIID
	if _, err = callPolarDBGatewayAI(client, "ModifyModelApi", request, d.Timeout(schema.TimeoutUpdate)); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "ModifyModelApi", AlibabaCloudSdkGoERROR)
	}
	return resourceAlicloudPolarDBGatewayModelAPIRead(d, meta)
}

func resourceAlicloudPolarDBGatewayModelAPIDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gatewayID, modelAPIID, err := parsePolarDBGatewayAIChildID(d.Id())
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{"RegionId": client.RegionId, "GwClusterId": gatewayID, "ModelApiId": modelAPIID}
	if _, err = callPolarDBGatewayAI(client, "DeleteModelApi", request, d.Timeout(schema.TimeoutDelete)); err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "DeleteModelApi", AlibabaCloudSdkGoERROR)
	}
	d.SetId("")
	return nil
}

func polarDBGatewayModelAPIRequest(d *schema.ResourceData, regionID string, includeCreateOnly bool) map[string]interface{} {
	request := map[string]interface{}{
		"RegionId": regionID, "GwClusterId": d.Get("gateway_id"), "ModelCategory": d.Get("model_category"),
		"PathPrefix": d.Get("path_prefix"), "Protocol": d.Get("protocol"), "RouteRules": d.Get("route_rules"),
	}
	if includeCreateOnly {
		request["Name"] = d.Get("name")
	}
	for key, requestKey := range map[string]string{"record_input": "RecordInput", "record_output": "RecordOutput"} {
		if value, ok := d.GetOk(key); ok {
			request[requestKey] = value
		}
	}
	if includeCreateOnly {
		if value, ok := d.GetOk("force_model"); ok {
			request["ForceModel"] = value
		}
	}
	return request
}
