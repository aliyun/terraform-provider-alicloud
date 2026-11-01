package polardb

import (
	"fmt"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceGatewayCostRule() *schema.Resource {
	return &schema.Resource{
		Create:   resourceGatewayCostRuleCreate,
		Read:     resourceGatewayCostRuleRead,
		Update:   resourceGatewayCostRuleUpdate,
		Delete:   resourceGatewayCostRuleDelete,
		Importer: &schema.ResourceImporter{State: schema.ImportStatePassthrough},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"cost_rule_id": {
				Type: schema.TypeString, Computed: true,
				Description: "The ID of the cost rule.",
			},
			"gateway_id": {
				Type: schema.TypeString, Required: true, ForceNew: true,
				Description: "The ID of the PolarDB AI gateway.",
			},
			"model_name": {
				Type: schema.TypeString, Required: true, ForceNew: true,
				Description: "The name of the model to which the cost rule applies.",
			},
			"model_service_id": {
				Type: schema.TypeString, Required: true, ForceNew: true,
				Description: "The ID of the model service to which the cost rule applies.",
			},
			"input_cost_points_per_million": {
				Type: schema.TypeString, Optional: true, Default: "0",
				StateFunc:   normalizePolarDBGatewayAIDecimal,
				Description: "The input cost in points per million tokens.",
			},
			"output_cost_points_per_million": {
				Type: schema.TypeString, Optional: true, Default: "0",
				StateFunc:   normalizePolarDBGatewayAIDecimal,
				Description: "The output cost in points per million tokens.",
			},
			"cache_cost_points_per_million": {
				Type: schema.TypeString, Optional: true, Default: "0",
				StateFunc:   normalizePolarDBGatewayAIDecimal,
				Description: "The cached-token cost in points per million tokens.",
			},
			"create_time": {
				Type: schema.TypeString, Computed: true,
				Description: "The time when the cost rule was created.",
			},
			"modify_time": {
				Type: schema.TypeString, Computed: true,
				Description: "The time when the cost rule was last modified.",
			},
		},
	}
}

func resourceGatewayCostRuleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "CreateCostRule"
	request := polarDBGatewayCostRuleRequest(d, client.RegionId)
	response, err := callPolarDBGatewayAI(client, action, request, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_polardb_gateway_cost_rule", action, errs.AlibabaCloudSdkGoERROR)
	}
	costRuleID := fmt.Sprint(response["CostRuleId"])
	if costRuleID == "" || costRuleID == "<nil>" {
		return errs.WrapError(fmt.Errorf("CreateCostRule returned empty CostRuleId"))
	}
	d.SetId(polarDBGatewayAIChildID(d.Get("gateway_id").(string), costRuleID))
	return resourceGatewayCostRuleRead(d, meta)
}

func resourceGatewayCostRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gatewayID, costRuleID, err := ParsePolarDBGatewayAIChildID(d.Id())
	if err != nil {
		return errs.WrapError(err)
	}
	object, err := FindPolarDBGatewayAIItem(client, "DescribeCostRules", gatewayID, "", "CostRuleId", costRuleID)
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), "DescribeCostRules", errs.AlibabaCloudSdkGoERROR)
	}
	if object == nil {
		d.SetId("")
		return nil
	}
	if err = d.Set("gateway_id", gatewayID); err != nil {
		return errs.WrapError(err)
	}
	if err = d.Set("cost_rule_id", costRuleID); err != nil {
		return errs.WrapError(err)
	}
	if err = d.Set("model_name", object["Model"]); err != nil {
		return errs.WrapError(err)
	}
	if err = d.Set("model_service_id", object["ModelServiceId"]); err != nil {
		return errs.WrapError(err)
	}
	if err = d.Set("input_cost_points_per_million", normalizePolarDBGatewayAIDecimal(object["InputCostPointsPerMillion"])); err != nil {
		return errs.WrapError(err)
	}
	if err = d.Set("output_cost_points_per_million", normalizePolarDBGatewayAIDecimal(object["OutputCostPointsPerMillion"])); err != nil {
		return errs.WrapError(err)
	}
	if err = d.Set("cache_cost_points_per_million", normalizePolarDBGatewayAIDecimal(object["CacheCostPointsPerMillion"])); err != nil {
		return errs.WrapError(err)
	}
	if err = d.Set("create_time", object["GmtCreated"]); err != nil {
		return errs.WrapError(err)
	}
	if err = d.Set("modify_time", object["GmtModified"]); err != nil {
		return errs.WrapError(err)
	}
	if _, ok := object["Model"]; !ok {
		if value, exists := object["ModelName"]; exists && value != nil {
			if err = d.Set("model_name", value); err != nil {
				return errs.WrapError(err)
			}
		}
	}
	return nil
}

func resourceGatewayCostRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gatewayID, costRuleID, err := ParsePolarDBGatewayAIChildID(d.Id())
	if err != nil {
		return errs.WrapError(err)
	}
	request := polarDBGatewayCostRuleRequest(d, client.RegionId)
	request["GwClusterId"] = gatewayID
	request["CostRuleId"] = costRuleID
	response, err := callPolarDBGatewayAI(client, "ModifyCostRule", request, d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), "ModifyCostRule", errs.AlibabaCloudSdkGoERROR)
	}
	// ModifyCostRule replaces the backend rule and returns a new ID. Track the
	// replacement so the following Read and eventual Delete address the live
	// rule instead of treating the resource as disappeared.
	if newCostRuleID := fmt.Sprint(response["CostRuleId"]); newCostRuleID != "" && newCostRuleID != "<nil>" {
		d.SetId(polarDBGatewayAIChildID(gatewayID, newCostRuleID))
	}
	return resourceGatewayCostRuleRead(d, meta)
}

func resourceGatewayCostRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gatewayID, costRuleID, err := ParsePolarDBGatewayAIChildID(d.Id())
	if err != nil {
		return errs.WrapError(err)
	}
	request := map[string]interface{}{"RegionId": client.RegionId, "GwClusterId": gatewayID, "CostRuleId": costRuleID}
	if _, err = callPolarDBGatewayAI(client, "DeleteCostRule", request, d.Timeout(schema.TimeoutDelete)); err != nil {
		if errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), "DeleteCostRule", errs.AlibabaCloudSdkGoERROR)
	}
	d.SetId("")
	return nil
}

func polarDBGatewayCostRuleRequest(d *schema.ResourceData, regionID string) map[string]interface{} {
	return map[string]interface{}{
		"RegionId": regionID, "GwClusterId": d.Get("gateway_id"), "ModelName": d.Get("model_name"),
		"ModelServiceId":             d.Get("model_service_id"),
		"InputCostPointsPerMillion":  d.Get("input_cost_points_per_million"),
		"OutputCostPointsPerMillion": d.Get("output_cost_points_per_million"),
		"CacheCostPointsPerMillion":  d.Get("cache_cost_points_per_million"),
	}
}
