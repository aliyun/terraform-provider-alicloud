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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceAliCloudGaCustomRoutingEndpoint() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudGaCustomRoutingEndpointCreate,
		Read:   resourceAliCloudGaCustomRoutingEndpointRead,
		Update: resourceAliCloudGaCustomRoutingEndpointUpdate,
		Delete: resourceAliCloudGaCustomRoutingEndpointDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"endpoint_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"endpoint": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"PrivateSubNet"}, false),
			},
			"traffic_to_endpoint_policy": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"DenyAll", "AllowAll", "AllowCustom"}, false),
			},
			"accelerator_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"listener_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"custom_routing_endpoint_id": {
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

func resourceAliCloudGaCustomRoutingEndpointCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := Service{client}
	var response map[string]interface{}
	action := "CreateCustomRoutingEndpoints"
	request := make(map[string]interface{})
	var err error

	request["RegionId"] = client.RegionId
	request["ClientToken"] = helper.BuildClientToken("CreateCustomRoutingEndpoints")
	request["EndpointGroupId"] = d.Get("endpoint_group_id")

	endpointConfigurationsMaps := make([]map[string]interface{}, 0)
	endpointConfigurationsMap := map[string]interface{}{}
	endpointConfigurationsMap["Endpoint"] = d.Get("endpoint")
	endpointConfigurationsMap["Type"] = d.Get("type")

	if v, ok := d.GetOk("traffic_to_endpoint_policy"); ok {
		endpointConfigurationsMap["TrafficToEndpointPolicy"] = v
	}

	endpointConfigurationsMaps = append(endpointConfigurationsMaps, endpointConfigurationsMap)
	request["EndpointConfigurations"] = endpointConfigurationsMaps

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *retry.RetryError {
		response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"NotActive.Listener", "StateError.EndPointGroup"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_ga_custom_routing_endpoint", action, errs.AlibabaCloudSdkGoERROR)
	}

	if resp, err := jsonpath.Get("$.EndpointIds", response); err != nil || resp == nil {
		return errs.WrapErrorf(err, errs.IdMsg, "alicloud_ga_custom_routing_endpoint")
	} else {
		endpointId := resp.([]interface{})[0]
		d.SetId(fmt.Sprintf("%v:%v", request["EndpointGroupId"], endpointId))
	}

	stateConf := helper.BuildStateConf([]string{}, []string{"active"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, gaService.GaCustomRoutingEndpointStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return resourceAliCloudGaCustomRoutingEndpointRead(d, meta)
}

func resourceAliCloudGaCustomRoutingEndpointRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := Service{client}

	object, err := gaService.DescribeGaCustomRoutingEndpoint(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("endpoint_group_id", object["EndpointGroupId"])
	d.Set("endpoint", object["Endpoint"])
	d.Set("type", object["Type"])
	d.Set("traffic_to_endpoint_policy", object["TrafficToEndpointPolicy"])
	d.Set("accelerator_id", object["AcceleratorId"])
	d.Set("listener_id", object["ListenerId"])
	d.Set("custom_routing_endpoint_id", object["EndpointId"])
	d.Set("status", object["State"])

	return nil
}

func resourceAliCloudGaCustomRoutingEndpointUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := Service{client}
	var response map[string]interface{}
	update := false

	parts, err := helper.ParseResourceId(d.Id(), 2)
	if err != nil {
		return errs.WrapError(err)
	}

	request := map[string]interface{}{
		"RegionId":        client.RegionId,
		"ClientToken":     helper.BuildClientToken("UpdateCustomRoutingEndpoints"),
		"EndpointGroupId": parts[0],
	}

	endpointConfigurationsMaps := make([]map[string]interface{}, 0)
	endpointConfigurationsMap := map[string]interface{}{}
	endpointConfigurationsMap["EndpointId"] = parts[1]

	if d.HasChange("traffic_to_endpoint_policy") {
		update = true
	}
	if v, ok := d.GetOk("traffic_to_endpoint_policy"); ok {
		endpointConfigurationsMap["TrafficToEndpointPolicy"] = v
	}

	endpointConfigurationsMaps = append(endpointConfigurationsMaps, endpointConfigurationsMap)
	request["EndpointConfigurations"] = endpointConfigurationsMaps

	if update {
		action := "UpdateCustomRoutingEndpoints"
		var err error

		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *retry.RetryError {
			response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
			if err != nil {
				if errs.IsExpectedErrors(err, []string{"NotActive.Listener", "StateError.EndPointGroup"}) || helper.NeedRetry(err) {
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

		stateConf := helper.BuildStateConf([]string{}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, gaService.GaCustomRoutingEndpointStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
		}
	}

	return resourceAliCloudGaCustomRoutingEndpointRead(d, meta)
}

func resourceAliCloudGaCustomRoutingEndpointDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := Service{client}
	action := "DeleteCustomRoutingEndpoints"
	var response map[string]interface{}

	var err error

	parts, err := helper.ParseResourceId(d.Id(), 2)
	if err != nil {
		return errs.WrapError(err)
	}

	request := map[string]interface{}{
		"RegionId":        client.RegionId,
		"ClientToken":     helper.BuildClientToken("DeleteCustomRoutingEndpoints"),
		"EndpointGroupId": parts[0],
		"EndpointIds":     []string{parts[1]},
	}

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *retry.RetryError {
		response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"NotActive.Listener", "StateError.EndPointGroup"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		if errs.IsExpectedErrors(err, []string{"NotExist.EndPointGroup", "NotExist.Listener"}) || errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	stateConf := helper.BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, gaService.GaCustomRoutingEndpointStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return nil
}
