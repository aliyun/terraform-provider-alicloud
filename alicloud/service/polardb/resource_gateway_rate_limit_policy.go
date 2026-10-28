package polardb

import (
	"fmt"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/validators"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceGatewayRateLimitPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceGatewayRateLimitPolicyCreate,
		Read:   resourceGatewayRateLimitPolicyRead,
		Update: resourceGatewayRateLimitPolicyUpdate,
		Delete: resourceGatewayRateLimitPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"policy_id": {
				Type: schema.TypeString, Computed: true,
				Description: "The ID of the rate limit policy.",
			},
			"gateway_id": {
				Type: schema.TypeString, Required: true, ForceNew: true,
				Description: "The ID of the PolarDB AI gateway.",
			},
			"scope_type": {
				Type: schema.TypeString, Required: true, ForceNew: true,
				ValidateFunc: validators.StringInSlice([]string{"ConsumerGroup", "Consumer"}, false),
				Description:  "The rate limit dimension.",
			},
			"scope_ref_id": {
				Type: schema.TypeString, Required: true, ForceNew: true,
				Description: "The ID of the consumer group or consumer.",
			},
			"rate_limit_rpm": {
				Type: schema.TypeString, Required: true,
				Description: "The maximum number of requests per minute.",
			},
			"rate_limit_tpm": {
				Type: schema.TypeString, Required: true,
				Description: "The maximum number of tokens per minute.",
			},
			"policy_type": {
				Type: schema.TypeString, Computed: true,
				Description: "The policy type.",
			},
			"status": {
				Type: schema.TypeString, Computed: true,
				Description: "The status of the rate limit policy.",
			},
			"create_time": {
				Type: schema.TypeString, Computed: true,
				Description: "The time when the rate limit policy was created.",
			},
			"modify_time": {
				Type: schema.TypeString, Computed: true,
				Description: "The time when the rate limit policy was last modified.",
			},
		},
	}
}

func resourceGatewayRateLimitPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "CreateRateLimitPolicy"
	request := polarDBGatewayRateLimitPolicyRequest(d, client.RegionId)
	response, err := callPolarDBGatewayAI(client, action, request, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_polardb_gateway_rate_limit_policy", action, errs.AlibabaCloudSdkGoERROR)
	}
	policyID := fmt.Sprint(response["PolicyId"])
	if policyID == "" || policyID == "<nil>" {
		return errs.WrapError(fmt.Errorf("CreateRateLimitPolicy returned empty PolicyId"))
	}
	d.SetId(polarDBGatewayAIChildID(d.Get("gateway_id").(string), policyID))
	return resourceGatewayRateLimitPolicyRead(d, meta)
}

func resourceGatewayRateLimitPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gatewayID, policyID, err := ParsePolarDBGatewayAIChildID(d.Id())
	if err != nil {
		return errs.WrapError(err)
	}
	object, err := FindPolarDBGatewayAIItem(client, "DescribeRateLimitPolicy", gatewayID, "PolicyId", "PolicyId", policyID)
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), "DescribeRateLimitPolicy", errs.AlibabaCloudSdkGoERROR)
	}
	if object == nil {
		d.SetId("")
		return nil
	}
	if err = d.Set("gateway_id", gatewayID); err != nil {
		return errs.WrapError(err)
	}
	if err = d.Set("policy_id", policyID); err != nil {
		return errs.WrapError(err)
	}
	if err = d.Set("scope_type", object["ScopeType"]); err != nil {
		return errs.WrapError(err)
	}
	if err = d.Set("scope_ref_id", object["ScopeRefId"]); err != nil {
		return errs.WrapError(err)
	}
	if err = d.Set("rate_limit_rpm", object["RateLimitRpm"]); err != nil {
		return errs.WrapError(err)
	}
	if err = d.Set("rate_limit_tpm", object["RateLimitTpm"]); err != nil {
		return errs.WrapError(err)
	}
	if err = d.Set("policy_type", object["PolicyType"]); err != nil {
		return errs.WrapError(err)
	}
	if err = d.Set("status", object["Status"]); err != nil {
		return errs.WrapError(err)
	}
	if err = d.Set("create_time", object["GmtCreated"]); err != nil {
		return errs.WrapError(err)
	}
	if err = d.Set("modify_time", object["GmtModified"]); err != nil {
		return errs.WrapError(err)
	}
	return nil
}

func resourceGatewayRateLimitPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gatewayID, policyID, err := ParsePolarDBGatewayAIChildID(d.Id())
	if err != nil {
		return errs.WrapError(err)
	}
	request := map[string]interface{}{
		"RegionId": client.RegionId, "GwClusterId": gatewayID, "PolicyId": policyID,
		"RateLimitRpm": d.Get("rate_limit_rpm"), "RateLimitTpm": d.Get("rate_limit_tpm"),
	}
	if _, err = callPolarDBGatewayAI(client, "ModifyRateLimitPolicy", request, d.Timeout(schema.TimeoutUpdate)); err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), "ModifyRateLimitPolicy", errs.AlibabaCloudSdkGoERROR)
	}
	return resourceGatewayRateLimitPolicyRead(d, meta)
}

func resourceGatewayRateLimitPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gatewayID, policyID, err := ParsePolarDBGatewayAIChildID(d.Id())
	if err != nil {
		return errs.WrapError(err)
	}
	request := map[string]interface{}{
		"RegionId": client.RegionId, "GwClusterId": gatewayID, "PolicyId": policyID,
	}
	if _, err = callPolarDBGatewayAI(client, "DeleteRateLimitPolicy", request, d.Timeout(schema.TimeoutDelete)); err != nil {
		if errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), "DeleteRateLimitPolicy", errs.AlibabaCloudSdkGoERROR)
	}
	d.SetId("")
	return nil
}

func polarDBGatewayRateLimitPolicyRequest(d *schema.ResourceData, regionID string) map[string]interface{} {
	return map[string]interface{}{
		"RegionId": regionID, "GwClusterId": d.Get("gateway_id"), "ScopeType": d.Get("scope_type"),
		"ScopeRefId": d.Get("scope_ref_id"), "RateLimitRpm": d.Get("rate_limit_rpm"), "RateLimitTpm": d.Get("rate_limit_tpm"),
	}
}
