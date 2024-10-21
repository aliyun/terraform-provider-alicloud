package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudGaCustomRoutingEndpoint() *schema.Resource {
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
	gaService := GaService{client}
	var response map[string]interface{}
	action := "CreateCustomRoutingEndpoints"
	request := make(map[string]interface{})
	var err error

	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken("CreateCustomRoutingEndpoints")
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

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"NotActive.Listener", "StateError.EndPointGroup"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ga_custom_routing_endpoint", action, AlibabaCloudSdkGoERROR)
	}

	if resp, err := jsonpath.Get("$.EndpointIds", response); err != nil || resp == nil {
		return WrapErrorf(err, IdMsg, "alicloud_ga_custom_routing_endpoint")
	} else {
		endpointId := resp.([]interface{})[0]
		d.SetId(fmt.Sprintf("%v:%v", request["EndpointGroupId"], endpointId))
	}

	stateConf := BuildStateConf([]string{}, []string{"active"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, gaService.GaCustomRoutingEndpointStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudGaCustomRoutingEndpointRead(d, meta)
}

func resourceAliCloudGaCustomRoutingEndpointRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := GaService{client}

	object, err := gaService.DescribeGaCustomRoutingEndpoint(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
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
	gaService := GaService{client}
	var response map[string]interface{}
	update := false

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}

	request := map[string]interface{}{
		"RegionId":        client.RegionId,
		"ClientToken":     buildClientToken("UpdateCustomRoutingEndpoints"),
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

		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"NotActive.Listener", "StateError.EndPointGroup"}) || NeedRetry(err) {
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

		stateConf := BuildStateConf([]string{}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, gaService.GaCustomRoutingEndpointStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	return resourceAliCloudGaCustomRoutingEndpointRead(d, meta)
}

func resourceAliCloudGaCustomRoutingEndpointDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := GaService{client}
	action := "DeleteCustomRoutingEndpoints"
	var response map[string]interface{}

	var err error

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}

	request := map[string]interface{}{
		"RegionId":        client.RegionId,
		"ClientToken":     buildClientToken("DeleteCustomRoutingEndpoints"),
		"EndpointGroupId": parts[0],
		"EndpointIds":     []string{parts[1]},
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"NotActive.Listener", "StateError.EndPointGroup"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		if IsExpectedErrors(err, []string{"NotExist.EndPointGroup", "NotExist.Listener"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, gaService.GaCustomRoutingEndpointStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}
