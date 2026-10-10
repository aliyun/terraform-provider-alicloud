package alicloud

import (
	"fmt"
	"strings"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudPolarDBGatewayBudgetPolicies() *schema.Resource {
	return &schema.Resource{Read: dataSourceAlicloudPolarDBGatewayBudgetPoliciesRead, Schema: map[string]*schema.Schema{
		"gateway_id":              {Type: schema.TypeString, Required: true},
		"ids":                     {Type: schema.TypeList, Optional: true, Elem: &schema.Schema{Type: schema.TypeString}},
		"status":                  {Type: schema.TypeString, Optional: true, ValidateFunc: StringInSlice([]string{"Enabled", "Disabled", "Disenabled"}, false)},
		"budget_dimension_type":   {Type: schema.TypeString, Optional: true, ValidateFunc: StringInSlice([]string{"ConsumerGroup", "Consumer"}, false)},
		"budget_dimension_ref_id": {Type: schema.TypeString, Optional: true},
		"scope_ref_name":          {Type: schema.TypeString, Optional: true},
		"policies": {Type: schema.TypeList, Computed: true, Elem: &schema.Resource{Schema: map[string]*schema.Schema{
			"id": {Type: schema.TypeString, Computed: true}, "budget_type": {Type: schema.TypeString, Computed: true},
			"budget_dimension_type": {Type: schema.TypeString, Computed: true}, "budget_dimension_ref_id": {Type: schema.TypeString, Computed: true},
			"reset_day_of_month": {Type: schema.TypeInt, Computed: true}, "budget_points": {Type: schema.TypeString, Computed: true},
			"alert_threshold_pct": {Type: schema.TypeInt, Computed: true}, "status": {Type: schema.TypeString, Computed: true},
			"used_points": {Type: schema.TypeInt, Computed: true}, "alert_triggered": {Type: schema.TypeBool, Computed: true},
			"exceeded": {Type: schema.TypeString, Computed: true}, "create_time": {Type: schema.TypeString, Computed: true},
			"modify_time": {Type: schema.TypeString, Computed: true},
		}}},
	}}
}

func dataSourceAlicloudPolarDBGatewayBudgetPoliciesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := map[string]interface{}{"GwClusterId": d.Get("gateway_id")}
	copyPolarDBGatewayAIFilters(d, request, map[string]string{
		"status": "Status", "budget_dimension_type": "BudgetDimensionType",
		"budget_dimension_ref_id": "BudgetDimensionRefId", "scope_ref_name": "ScopeRefName",
	})
	if ids := stringListFromSchema(d.Get("ids")); len(ids) == 1 {
		request["BudgetPolicyId"] = ids[0]
	}
	objects, err := listPolarDBGatewayAI(client, "DescribeBudgetPolicies", request)
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_polardb_gateway_budget_policies", "DescribeBudgetPolicies", AlibabaCloudSdkGoERROR)
	}
	requestedIDs := make(map[string]struct{})
	for _, id := range stringListFromSchema(d.Get("ids")) {
		requestedIDs[id] = struct{}{}
	}
	items := make([]map[string]interface{}, 0, len(objects))
	ids := make([]string, 0, len(objects))
	for _, object := range objects {
		id := strings.TrimSpace(fmt.Sprint(object["BudgetPolicyId"]))
		if len(requestedIDs) > 0 {
			if _, ok := requestedIDs[id]; !ok {
				continue
			}
		}
		ids = append(ids, id)
		items = append(items, map[string]interface{}{
			"id": id, "budget_type": object["BudgetType"], "budget_dimension_type": object["BudgetDimensionType"],
			"budget_dimension_ref_id": object["BudgetDimensionRefId"], "reset_day_of_month": polarDBGatewayAIInt(object["ResetDayOfMonth"]),
			"budget_points": object["BudgetPoints"], "alert_threshold_pct": polarDBGatewayAIInt(object["AlertThresholdPct"]),
			"status": object["Status"], "used_points": polarDBGatewayAIInt(object["UsedPoints"]), "alert_triggered": object["AlertTriggered"],
			"exceeded": object["Exceeded"], "create_time": object["GmtCreated"], "modify_time": object["GmtModified"],
		})
	}
	if err = d.Set("policies", items); err != nil {
		return WrapError(err)
	}
	d.SetId(dataResourceIdHash(ids))
	return nil
}
