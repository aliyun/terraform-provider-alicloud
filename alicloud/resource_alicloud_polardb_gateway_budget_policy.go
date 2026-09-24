package alicloud

import (
	"fmt"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudPolarDBGatewayBudgetPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudPolarDBGatewayBudgetPolicyCreate,
		Read:   resourceAlicloudPolarDBGatewayBudgetPolicyRead,
		Update: resourceAlicloudPolarDBGatewayBudgetPolicyUpdate,
		Delete: resourceAlicloudPolarDBGatewayBudgetPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		CustomizeDiff: validatePolarDBGatewayBudgetPolicyDiff,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"budget_policy_id": {
				Type: schema.TypeString, Computed: true,
				Description: "The ID of the budget policy.",
			},
			"gateway_id": {
				Type: schema.TypeString, Required: true, ForceNew: true,
				Description: "The ID of the PolarDB AI gateway.",
			},
			"budget_type": {
				Type: schema.TypeString, Required: true, ForceNew: true,
				ValidateFunc: StringInSlice([]string{"GlobalTotal", "ConsumerTotal", "ConsumerGroupTotal"}, false),
				Description:  "The budget type.",
			},
			"budget_dimension_ref_id": {
				Type: schema.TypeString, Optional: true, ForceNew: true,
				Description: "The consumer or consumer group ID to which the budget applies.",
			},
			"reset_day_of_month": {
				Type: schema.TypeInt, Required: true,
				ValidateFunc: IntBetween(1, 28),
				Description:  "The day of the month on which the budget is reset.",
			},
			"budget_points": {
				Type: schema.TypeString, Required: true,
				StateFunc:   normalizePolarDBGatewayAIDecimal,
				Description: "The budget amount in points.",
			},
			"alert_threshold_pct": {
				Type: schema.TypeInt, Optional: true,
				ValidateFunc: IntBetween(0, 100),
				Description:  "The budget usage percentage that triggers an alert.",
			},
			"status": {
				Type: schema.TypeString, Computed: true,
				Description: "The status of the budget policy.",
			},
			"budget_dimension_type": {
				Type: schema.TypeString, Computed: true,
				Description: "The dimension of the budget policy.",
			},
			"used_points": {
				Type: schema.TypeInt, Computed: true,
				Description: "The number of used points.",
			},
			"alert_triggered": {
				Type: schema.TypeBool, Computed: true,
				Description: "Whether the alert threshold was triggered.",
			},
			"exceeded": {
				Type: schema.TypeString, Computed: true,
				Description: "Whether the budget was exceeded.",
			},
			"create_time": {
				Type: schema.TypeString, Computed: true,
				Description: "The time when the budget policy was created.",
			},
			"modify_time": {
				Type: schema.TypeString, Computed: true,
				Description: "The time when the budget policy was last modified.",
			},
		},
	}
}

func validatePolarDBGatewayBudgetPolicyDiff(diff *schema.ResourceDiff, _ interface{}) error {
	budgetType := diff.Get("budget_type").(string)
	// References to resources created in the same plan are unknown during
	// CustomizeDiff. Defer validation until their concrete value is available.
	if !diff.NewValueKnown("budget_dimension_ref_id") {
		return nil
	}
	_, hasDimensionRefID := diff.GetOk("budget_dimension_ref_id")
	if budgetType == "GlobalTotal" && hasDimensionRefID {
		return fmt.Errorf("budget_dimension_ref_id must not be set when budget_type is GlobalTotal")
	}
	if budgetType != "GlobalTotal" && !hasDimensionRefID {
		return fmt.Errorf("budget_dimension_ref_id must be set when budget_type is %s", budgetType)
	}
	return nil
}

func resourceAlicloudPolarDBGatewayBudgetPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "CreateBudgetPolicy"
	request := polarDBGatewayBudgetPolicyRequest(d, client.RegionId)
	response, err := callPolarDBGatewayAI(client, action, request, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_polardb_gateway_budget_policy", action, AlibabaCloudSdkGoERROR)
	}
	budgetPolicyID := fmt.Sprint(response["BudgetPolicyId"])
	if budgetPolicyID == "" || budgetPolicyID == "<nil>" {
		return WrapError(fmt.Errorf("CreateBudgetPolicy returned empty BudgetPolicyId"))
	}
	d.SetId(polarDBGatewayAIChildID(d.Get("gateway_id").(string), budgetPolicyID))
	return resourceAlicloudPolarDBGatewayBudgetPolicyRead(d, meta)
}

func resourceAlicloudPolarDBGatewayBudgetPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gatewayID, budgetPolicyID, err := parsePolarDBGatewayAIChildID(d.Id())
	if err != nil {
		return WrapError(err)
	}
	object, err := findPolarDBGatewayAIItem(client, "DescribeBudgetPolicies", gatewayID, "BudgetPolicyId", "BudgetPolicyId", budgetPolicyID)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "DescribeBudgetPolicies", AlibabaCloudSdkGoERROR)
	}
	if object == nil {
		d.SetId("")
		return nil
	}
	values := map[string]interface{}{
		"gateway_id": gatewayID, "budget_policy_id": budgetPolicyID, "budget_type": object["BudgetType"],
		"budget_dimension_type": object["BudgetDimensionType"], "reset_day_of_month": polarDBGatewayAIInt(object["ResetDayOfMonth"]),
		"budget_points":       object["BudgetPoints"],
		"alert_threshold_pct": polarDBGatewayAIInt(object["AlertThresholdPct"]), "status": object["Status"], "used_points": polarDBGatewayAIInt(object["UsedPoints"]),
		"alert_triggered": object["AlertTriggered"], "exceeded": object["Exceeded"],
		"create_time": object["GmtCreated"], "modify_time": object["GmtModified"],
	}
	// The service reports the synthetic value "Global" for a GlobalTotal
	// policy. It is not a writable dimension reference and setting it would
	// force replacement on every plan because GlobalTotal must omit the field.
	if fmt.Sprint(object["BudgetType"]) != "GlobalTotal" {
		values["budget_dimension_ref_id"] = object["BudgetDimensionRefId"]
	}
	for key, value := range values {
		if err = d.Set(key, value); err != nil {
			return WrapError(err)
		}
	}
	return nil
}

func resourceAlicloudPolarDBGatewayBudgetPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gatewayID, budgetPolicyID, err := parsePolarDBGatewayAIChildID(d.Id())
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"RegionId": client.RegionId, "GwClusterId": gatewayID, "BudgetPolicyId": budgetPolicyID,
		"ResetDayOfMonth": d.Get("reset_day_of_month"), "BudgetPoints": d.Get("budget_points"),
	}
	if value, ok := d.GetOkExists("alert_threshold_pct"); ok {
		request["AlertThresholdPct"] = value
	}
	if _, err = callPolarDBGatewayAI(client, "ModifyBudgetPolicy", request, d.Timeout(schema.TimeoutUpdate)); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "ModifyBudgetPolicy", AlibabaCloudSdkGoERROR)
	}
	return resourceAlicloudPolarDBGatewayBudgetPolicyRead(d, meta)
}

func resourceAlicloudPolarDBGatewayBudgetPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gatewayID, budgetPolicyID, err := parsePolarDBGatewayAIChildID(d.Id())
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"RegionId": client.RegionId, "GwClusterId": gatewayID, "BudgetPolicyId": budgetPolicyID,
	}
	if _, err = callPolarDBGatewayAI(client, "DeleteBudgetPolicy", request, d.Timeout(schema.TimeoutDelete)); err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "DeleteBudgetPolicy", AlibabaCloudSdkGoERROR)
	}
	d.SetId("")
	return nil
}

func polarDBGatewayBudgetPolicyRequest(d *schema.ResourceData, regionID string) map[string]interface{} {
	request := map[string]interface{}{
		"RegionId": regionID, "GwClusterId": d.Get("gateway_id"), "BudgetType": d.Get("budget_type"),
		"ResetDayOfMonth": d.Get("reset_day_of_month"), "BudgetPoints": d.Get("budget_points"),
	}
	if value, ok := d.GetOk("budget_dimension_ref_id"); ok {
		request["BudgetDimensionRefId"] = value
	}
	if value, ok := d.GetOkExists("alert_threshold_pct"); ok {
		request["AlertThresholdPct"] = value
	}
	return request
}
