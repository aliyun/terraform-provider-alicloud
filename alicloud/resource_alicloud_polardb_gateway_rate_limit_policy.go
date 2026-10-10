package alicloud

import (
	"fmt"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudPolarDBGatewayRateLimitPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudPolarDBGatewayRateLimitPolicyCreate,
		Read:   resourceAlicloudPolarDBGatewayRateLimitPolicyRead,
		Update: resourceAlicloudPolarDBGatewayRateLimitPolicyUpdate,
		Delete: resourceAlicloudPolarDBGatewayRateLimitPolicyDelete,
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
				ValidateFunc: StringInSlice([]string{"ConsumerGroup", "Consumer"}, false),
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

func resourceAlicloudPolarDBGatewayRateLimitPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "CreateRateLimitPolicy"
	request := polarDBGatewayRateLimitPolicyRequest(d, client.RegionId)
	response, err := client.RpcPost(polarDBGatewayAIProduct, polarDBGatewayAIVersion, action, nil, request, false)
	debugPolarDBGatewayAI(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_polardb_gateway_rate_limit_policy", action, AlibabaCloudSdkGoERROR)
	}
	policyID := fmt.Sprint(response["PolicyId"])
	if policyID == "" || policyID == "<nil>" {
		return WrapError(fmt.Errorf("CreateRateLimitPolicy returned empty PolicyId"))
	}
	d.SetId(polarDBGatewayAIChildID(d.Get("gateway_id").(string), policyID))
	return resourceAlicloudPolarDBGatewayRateLimitPolicyRead(d, meta)
}

func resourceAlicloudPolarDBGatewayRateLimitPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gatewayID, policyID, err := parsePolarDBGatewayAIChildID(d.Id())
	if err != nil {
		return WrapError(err)
	}
	object, err := findPolarDBGatewayAIItem(client, "DescribeRateLimitPolicy", gatewayID, "PolicyId", "PolicyId", policyID)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "DescribeRateLimitPolicy", AlibabaCloudSdkGoERROR)
	}
	if object == nil {
		d.SetId("")
		return nil
	}
	values := map[string]interface{}{
		"gateway_id": gatewayID, "policy_id": policyID, "scope_type": object["ScopeType"],
		"scope_ref_id": object["ScopeRefId"], "rate_limit_rpm": object["RateLimitRpm"],
		"rate_limit_tpm": object["RateLimitTpm"], "policy_type": object["PolicyType"],
		"status": object["Status"], "create_time": object["GmtCreated"], "modify_time": object["GmtModified"],
	}
	for key, value := range values {
		if err = d.Set(key, value); err != nil {
			return WrapError(err)
		}
	}
	return nil
}

func resourceAlicloudPolarDBGatewayRateLimitPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gatewayID, policyID, err := parsePolarDBGatewayAIChildID(d.Id())
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"RegionId": client.RegionId, "GwClusterId": gatewayID, "PolicyId": policyID,
		"RateLimitRpm": d.Get("rate_limit_rpm"), "RateLimitTpm": d.Get("rate_limit_tpm"),
	}
	if _, err = callPolarDBGatewayAI(client, "ModifyRateLimitPolicy", request, d.Timeout(schema.TimeoutUpdate)); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "ModifyRateLimitPolicy", AlibabaCloudSdkGoERROR)
	}
	return resourceAlicloudPolarDBGatewayRateLimitPolicyRead(d, meta)
}

func resourceAlicloudPolarDBGatewayRateLimitPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gatewayID, policyID, err := parsePolarDBGatewayAIChildID(d.Id())
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"RegionId": client.RegionId, "GwClusterId": gatewayID, "PolicyId": policyID,
	}
	if _, err = callPolarDBGatewayAI(client, "DeleteRateLimitPolicy", request, d.Timeout(schema.TimeoutDelete)); err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "DeleteRateLimitPolicy", AlibabaCloudSdkGoERROR)
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
