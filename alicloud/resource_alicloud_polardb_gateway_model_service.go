package alicloud

import (
	"fmt"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudPolarDBGatewayModelService() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudPolarDBGatewayModelServiceCreate,
		Read:   resourceAlicloudPolarDBGatewayModelServiceRead,
		Update: resourceAlicloudPolarDBGatewayModelServiceUpdate,
		Delete: resourceAlicloudPolarDBGatewayModelServiceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"model_service_id": {
				Type: schema.TypeString, Computed: true,
				Description: "The ID of the model service.",
			},
			"gateway_id": {
				Type: schema.TypeString, Required: true, ForceNew: true,
				Description: "The ID of the PolarDB AI gateway.",
			},
			"name": {
				Type: schema.TypeString, Required: true, ForceNew: true,
				Description: "The name of the model service.",
			},
			"model_category": {
				Type: schema.TypeString, Required: true,
				ValidateFunc: StringInSlice([]string{"text", "embedding", "rerank"}, false),
				Description:  "The model category. Valid values: `text`, `embedding`, `rerank`.",
			},
			"protocol": {
				Type: schema.TypeString, Required: true,
				ValidateFunc: StringInSlice([]string{"openai", "anthropic", "bailian", "vllm"}, false),
				Description:  "The protocol used by the upstream model service.",
			},
			"base_url": {
				Type: schema.TypeString, Required: true,
				ValidateFunc: validatePolarDBGatewayAIURL,
				Description:  "The HTTP or HTTPS URL of the upstream model service.",
			},
			"api_key": {
				Type: schema.TypeString, Required: true, Sensitive: true,
				Description: "The API key of the upstream model service. The value is stored in Terraform state.",
			},
			"vendor": {
				Type: schema.TypeString, Required: true, ForceNew: true,
				Description: "The vendor of the model service.",
			},
			"input_cost_points_per_million": {
				Type: schema.TypeString, Optional: true,
				StateFunc:   normalizePolarDBGatewayAIDecimal,
				Description: "The input cost in points per million tokens.",
			},
			"output_cost_points_per_million": {
				Type: schema.TypeString, Optional: true,
				StateFunc:   normalizePolarDBGatewayAIDecimal,
				Description: "The output cost in points per million tokens.",
			},
			"request_cost_points": {
				Type: schema.TypeString, Optional: true,
				StateFunc:   normalizePolarDBGatewayAIDecimal,
				Description: "The cost in points per request.",
			},
			"status": {
				Type: schema.TypeString, Computed: true,
				Description: "The status of the model service.",
			},
			"create_time": {
				Type: schema.TypeString, Computed: true,
				Description: "The time when the model service was created.",
			},
		},
	}
}

func resourceAlicloudPolarDBGatewayModelServiceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "CreateModelService"
	request := polarDBGatewayModelServiceRequest(d, client.RegionId)
	response, err := client.RpcPost(polarDBGatewayAIProduct, polarDBGatewayAIVersion, action, nil, request, false)
	debugPolarDBGatewayAI(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_polardb_gateway_model_service", action, AlibabaCloudSdkGoERROR)
	}
	modelServiceID := fmt.Sprint(response["ModelServiceId"])
	if modelServiceID == "" || modelServiceID == "<nil>" {
		return WrapError(fmt.Errorf("CreateModelService returned empty ModelServiceId"))
	}
	d.SetId(polarDBGatewayAIChildID(d.Get("gateway_id").(string), modelServiceID))
	return resourceAlicloudPolarDBGatewayModelServiceRead(d, meta)
}

func resourceAlicloudPolarDBGatewayModelServiceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gatewayID, modelServiceID, err := parsePolarDBGatewayAIChildID(d.Id())
	if err != nil {
		return WrapError(err)
	}
	object, err := findPolarDBGatewayAIItem(client, "DescribeModelServices", gatewayID, "ModelServiceIds", "ModelServiceId", modelServiceID)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "DescribeModelServices", AlibabaCloudSdkGoERROR)
	}
	if object == nil {
		d.SetId("")
		return nil
	}
	if err = d.Set("gateway_id", gatewayID); err != nil {
		return WrapError(err)
	}
	if err = d.Set("model_service_id", modelServiceID); err != nil {
		return WrapError(err)
	}
	if err = d.Set("name", object["Name"]); err != nil {
		return WrapError(err)
	}
	// DescribeModelServices does not return all writable fields on every
	// deployment. Keep the configured value when a field is absent instead of
	// replacing it with an empty value and creating a perpetual diff.
	if value, ok := object["ModelCategory"]; ok {
		if err = d.Set("model_category", value); err != nil {
			return WrapError(err)
		}
	}
	if err = d.Set("protocol", object["Protocol"]); err != nil {
		return WrapError(err)
	}
	if err = d.Set("base_url", object["BaseUrl"]); err != nil {
		return WrapError(err)
	}
	if err = d.Set("vendor", object["Vendor"]); err != nil {
		return WrapError(err)
	}
	if value, ok := object["InputCostPointsPerMillion"]; ok {
		if err = d.Set("input_cost_points_per_million", value); err != nil {
			return WrapError(err)
		}
	}
	if value, ok := object["OutputCostPointsPerMillion"]; ok {
		if err = d.Set("output_cost_points_per_million", value); err != nil {
			return WrapError(err)
		}
	}
	if value, ok := object["RequestCostPoints"]; ok {
		if err = d.Set("request_cost_points", value); err != nil {
			return WrapError(err)
		}
	}
	if err = d.Set("status", object["Status"]); err != nil {
		return WrapError(err)
	}
	if err = d.Set("create_time", object["GmtCreated"]); err != nil {
		return WrapError(err)
	}
	// ApiKey is deliberately not read back: the API may return a masked value,
	// and replacing the configured secret would cause a perpetual diff.
	return nil
}

func resourceAlicloudPolarDBGatewayModelServiceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gatewayID, modelServiceID, err := parsePolarDBGatewayAIChildID(d.Id())
	if err != nil {
		return WrapError(err)
	}
	request := polarDBGatewayModelServiceRequest(d, client.RegionId)
	request["GwClusterId"] = gatewayID
	request["ModelServiceId"] = modelServiceID
	delete(request, "Vendor")
	if _, err = callPolarDBGatewayAI(client, "ModifyModelService", request, d.Timeout(schema.TimeoutUpdate)); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "ModifyModelService", AlibabaCloudSdkGoERROR)
	}
	return resourceAlicloudPolarDBGatewayModelServiceRead(d, meta)
}

func resourceAlicloudPolarDBGatewayModelServiceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gatewayID, _, err := parsePolarDBGatewayAIChildID(d.Id())
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"RegionId": client.RegionId, "GwClusterId": gatewayID, "ModelName": d.Get("name").(string),
	}
	if _, err = callPolarDBGatewayAI(client, "DeleteModelService", request, d.Timeout(schema.TimeoutDelete)); err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "DeleteModelService", AlibabaCloudSdkGoERROR)
	}
	d.SetId("")
	return nil
}

func polarDBGatewayModelServiceRequest(d *schema.ResourceData, regionID string) map[string]interface{} {
	request := map[string]interface{}{
		"RegionId": regionID, "GwClusterId": d.Get("gateway_id"), "Name": d.Get("name"),
		"ModelCategory": d.Get("model_category"), "Protocol": d.Get("protocol"), "BaseUrl": d.Get("base_url"),
		"ApiKey": d.Get("api_key"),
	}
	for key, requestKey := range map[string]string{
		"vendor": "Vendor", "input_cost_points_per_million": "InputCostPointsPerMillion",
		"output_cost_points_per_million": "OutputCostPointsPerMillion", "request_cost_points": "RequestCostPoints",
	} {
		if value, ok := d.GetOk(key); ok {
			request[requestKey] = value
		}
	}
	return request
}
