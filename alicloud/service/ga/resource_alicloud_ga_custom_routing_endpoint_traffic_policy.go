package ga

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceAliCloudGaCustomRoutingEndpointTrafficPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudGaCustomRoutingEndpointTrafficPolicyCreate,
		Read:   resourceAliCloudGaCustomRoutingEndpointTrafficPolicyRead,
		Update: resourceAliCloudGaCustomRoutingEndpointTrafficPolicyUpdate,
		Delete: resourceAliCloudGaCustomRoutingEndpointTrafficPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"endpoint_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"address": {
				Type:     schema.TypeString,
				Required: true,
			},
			"port_ranges": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"from_port": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"to_port": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
			"accelerator_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"listener_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"endpoint_group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"custom_routing_endpoint_traffic_policy_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudGaCustomRoutingEndpointTrafficPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := Service{client}
	var response map[string]interface{}
	action := "CreateCustomRoutingEndpointTrafficPolicies"
	request := make(map[string]interface{})
	var err error

	request["RegionId"] = client.RegionId
	request["ClientToken"] = helper.BuildClientToken("CreateCustomRoutingEndpointTrafficPolicies")
	request["EndpointId"] = d.Get("endpoint_id")

	policyConfigurationsMaps := make([]map[string]interface{}, 0)
	policyConfigurationsMap := map[string]interface{}{}
	policyConfigurationsMap["Address"] = d.Get("address")

	if v, ok := d.GetOk("port_ranges"); ok {
		portRangesMaps := make([]map[string]interface{}, 0)
		for _, portRanges := range v.([]interface{}) {
			portRangesMap := map[string]interface{}{}
			portRangesArg := portRanges.(map[string]interface{})

			if fromPort, ok := portRangesArg["from_port"]; ok {
				portRangesMap["FromPort"] = fromPort
			}

			if toPort, ok := portRangesArg["to_port"]; ok {
				portRangesMap["ToPort"] = toPort
			}

			portRangesMaps = append(portRangesMaps, portRangesMap)
		}

		policyConfigurationsMap["PortRanges"] = portRangesMaps
	}

	policyConfigurationsMaps = append(policyConfigurationsMaps, policyConfigurationsMap)
	request["PolicyConfigurations"] = policyConfigurationsMaps

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *retry.RetryError {
		response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"NotActive.Listener", "StateError.Accelerator", "StateError.EndPointGroup"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_ga_custom_routing_endpoint_traffic_policy", action, errs.AlibabaCloudSdkGoERROR)
	}

	if resp, err := jsonpath.Get("$.PolicyIds", response); err != nil || resp == nil {
		return errs.WrapErrorf(err, errs.IdMsg, "alicloud_ga_custom_routing_endpoint_traffic_policy")
	} else {
		policyId := resp.([]interface{})[0]
		d.SetId(fmt.Sprintf("%v:%v", request["EndpointId"], policyId))
	}

	stateConf := helper.BuildStateConf([]string{}, []string{"active"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, gaService.GaCustomRoutingEndpointTrafficPolicyStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return resourceAliCloudGaCustomRoutingEndpointTrafficPolicyRead(d, meta)
}

func resourceAliCloudGaCustomRoutingEndpointTrafficPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := Service{client}

	object, err := gaService.DescribeGaCustomRoutingEndpointTrafficPolicy(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("endpoint_id", object["EndpointId"])
	d.Set("address", object["Address"])

	if portRangesList, ok := object["PortRanges"]; ok {
		portRangesMaps := make([]map[string]interface{}, 0)
		for _, portRanges := range portRangesList.([]interface{}) {
			portRangesArg := portRanges.(map[string]interface{})
			portRangesMap := map[string]interface{}{}

			if fromPort, ok := portRangesArg["FromPort"]; ok {
				portRangesMap["from_port"] = fromPort
			}

			if toPort, ok := portRangesArg["ToPort"]; ok {
				portRangesMap["to_port"] = toPort
			}

			portRangesMaps = append(portRangesMaps, portRangesMap)
		}

		d.Set("port_ranges", portRangesMaps)
	}

	d.Set("accelerator_id", object["AcceleratorId"])
	d.Set("listener_id", object["ListenerId"])
	d.Set("endpoint_group_id", object["EndpointGroupId"])
	d.Set("custom_routing_endpoint_traffic_policy_id", object["PolicyId"])
	d.Set("status", object["State"])

	return nil
}

func resourceAliCloudGaCustomRoutingEndpointTrafficPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := Service{client}
	var response map[string]interface{}
	update := false

	parts, err := helper.ParseResourceId(d.Id(), 2)
	if err != nil {
		return errs.WrapError(err)
	}

	request := map[string]interface{}{
		"RegionId":    client.RegionId,
		"ClientToken": helper.BuildClientToken("UpdateCustomRoutingEndpointTrafficPolicies"),
		"EndpointId":  parts[0],
	}

	policyConfigurationsMaps := make([]map[string]interface{}, 0)
	policyConfigurationsMap := map[string]interface{}{}
	policyConfigurationsMap["PolicyId"] = parts[1]

	if d.HasChange("address") {
		update = true
	}
	policyConfigurationsMap["Address"] = d.Get("address")

	if d.HasChange("port_ranges") {
		update = true
	}
	if v, ok := d.GetOk("port_ranges"); ok {
		portRangesMaps := make([]map[string]interface{}, 0)
		for _, portRanges := range v.([]interface{}) {
			portRangesMap := map[string]interface{}{}
			portRangesArg := portRanges.(map[string]interface{})

			if fromPort, ok := portRangesArg["from_port"]; ok {
				portRangesMap["FromPort"] = fromPort
			}

			if toPort, ok := portRangesArg["to_port"]; ok {
				portRangesMap["ToPort"] = toPort
			}

			portRangesMaps = append(portRangesMaps, portRangesMap)
		}

		policyConfigurationsMap["PortRanges"] = portRangesMaps
	}

	policyConfigurationsMaps = append(policyConfigurationsMaps, policyConfigurationsMap)
	request["PolicyConfigurations"] = policyConfigurationsMaps

	if update {
		action := "UpdateCustomRoutingEndpointTrafficPolicies"
		var err error

		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *retry.RetryError {
			response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
			if err != nil {
				if errs.IsExpectedErrors(err, []string{"NotActive.Listener", "StateError.Accelerator", "StateError.EndPointGroup"}) || helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)

		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}

		stateConf := helper.BuildStateConf([]string{}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, gaService.GaCustomRoutingEndpointTrafficPolicyStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
		}
	}

	return resourceAliCloudGaCustomRoutingEndpointTrafficPolicyRead(d, meta)
}

func resourceAliCloudGaCustomRoutingEndpointTrafficPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := Service{client}
	action := "DeleteCustomRoutingEndpointTrafficPolicies"
	var response map[string]interface{}

	var err error

	parts, err := helper.ParseResourceId(d.Id(), 2)
	if err != nil {
		return errs.WrapError(err)
	}

	request := map[string]interface{}{
		"RegionId":    client.RegionId,
		"ClientToken": helper.BuildClientToken("DeleteCustomRoutingEndpointTrafficPolicies"),
		"EndpointId":  parts[0],
		"PolicyIds":   []string{parts[1]},
	}

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *retry.RetryError {
		response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"NotActive.Listener", "StateError.Accelerator", "StateError.EndPointGroup"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		if errs.IsExpectedErrors(err, []string{"IllegalParameter.TrafficPolicyAndEndpointIdMismatch"}) || errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	stateConf := helper.BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, gaService.GaCustomRoutingEndpointTrafficPolicyStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return nil
}
