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

func resourceAliCloudGaCustomRoutingEndpointGroupDestination() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudGaCustomRoutingEndpointGroupDestinationCreate,
		Read:   resourceAliCloudGaCustomRoutingEndpointGroupDestinationRead,
		Update: resourceAliCloudGaCustomRoutingEndpointGroupDestinationUpdate,
		Delete: resourceAliCloudGaCustomRoutingEndpointGroupDestinationDelete,
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
			"protocols": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"from_port": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(1, 65499),
			},
			"to_port": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(1, 65499),
			},
			"accelerator_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"listener_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"custom_routing_endpoint_group_destination_id": {
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

func resourceAliCloudGaCustomRoutingEndpointGroupDestinationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := GaService{client}
	var response map[string]interface{}
	action := "CreateCustomRoutingEndpointGroupDestinations"
	request := make(map[string]interface{})
	var err error

	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken("CreateCustomRoutingEndpointGroupDestinations")
	request["EndpointGroupId"] = d.Get("endpoint_group_id")

	destinationConfigurationsMaps := make([]map[string]interface{}, 0)
	destinationConfigurationsMap := map[string]interface{}{}
	destinationConfigurationsMap["Protocols"] = d.Get("protocols")
	destinationConfigurationsMap["FromPort"] = d.Get("from_port")
	destinationConfigurationsMap["ToPort"] = d.Get("to_port")
	destinationConfigurationsMaps = append(destinationConfigurationsMaps, destinationConfigurationsMap)
	request["DestinationConfigurations"] = destinationConfigurationsMaps

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"NotActive.Listener", "StateError.Accelerator", "StateError.EndPointGroup"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ga_custom_routing_endpoint_group_destination", action, AlibabaCloudSdkGoERROR)
	}

	if resp, err := jsonpath.Get("$.DestinationIds", response); err != nil || resp == nil {
		return WrapErrorf(err, IdMsg, "alicloud_ga_custom_routing_endpoint_group_destination")
	} else {
		destinationId := resp.([]interface{})[0]
		d.SetId(fmt.Sprintf("%v:%v", request["EndpointGroupId"], destinationId))
	}

	stateConf := BuildStateConf([]string{}, []string{"active"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, gaService.GaCustomRoutingEndpointGroupDestinationStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudGaCustomRoutingEndpointGroupDestinationRead(d, meta)
}

func resourceAliCloudGaCustomRoutingEndpointGroupDestinationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := GaService{client}

	object, err := gaService.DescribeGaCustomRoutingEndpointGroupDestination(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("endpoint_group_id", object["EndpointGroupId"])
	d.Set("protocols", object["Protocols"])
	d.Set("from_port", object["FromPort"])
	d.Set("to_port", object["ToPort"])
	d.Set("accelerator_id", object["AcceleratorId"])
	d.Set("listener_id", object["ListenerId"])
	d.Set("custom_routing_endpoint_group_destination_id", object["DestinationId"])
	d.Set("status", object["State"])

	return nil
}

func resourceAliCloudGaCustomRoutingEndpointGroupDestinationUpdate(d *schema.ResourceData, meta interface{}) error {
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
		"ClientToken":     buildClientToken("UpdateCustomRoutingEndpointGroupDestinations"),
		"EndpointGroupId": parts[0],
	}

	destinationConfigurationsMaps := make([]map[string]interface{}, 0)
	destinationConfigurationsMap := map[string]interface{}{}
	destinationConfigurationsMap["DestinationId"] = parts[1]

	if d.HasChange("protocols") {
		update = true
	}
	destinationConfigurationsMap["Protocols"] = d.Get("protocols")

	if d.HasChange("from_port") {
		update = true
	}
	destinationConfigurationsMap["FromPort"] = d.Get("from_port")

	if d.HasChange("to_port") {
		update = true
	}
	destinationConfigurationsMap["ToPort"] = d.Get("to_port")

	destinationConfigurationsMaps = append(destinationConfigurationsMaps, destinationConfigurationsMap)
	request["DestinationConfigurations"] = destinationConfigurationsMaps

	if update {
		action := "UpdateCustomRoutingEndpointGroupDestinations"
		var err error

		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"NotActive.Listener", "StateError.Accelerator", "StateError.EndPointGroup"}) || NeedRetry(err) {
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

		stateConf := BuildStateConf([]string{}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, gaService.GaCustomRoutingEndpointGroupDestinationStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	return resourceAliCloudGaCustomRoutingEndpointGroupDestinationRead(d, meta)
}

func resourceAliCloudGaCustomRoutingEndpointGroupDestinationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := GaService{client}
	action := "DeleteCustomRoutingEndpointGroupDestinations"
	var response map[string]interface{}

	var err error

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}

	request := map[string]interface{}{
		"RegionId":        client.RegionId,
		"ClientToken":     buildClientToken("DeleteCustomRoutingEndpointGroupDestinations"),
		"EndpointGroupId": parts[0],
		"DestinationIds":  []string{parts[1]},
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"NotActive.Listener", "StateError.Accelerator", "StateError.EndPointGroup"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, gaService.GaCustomRoutingEndpointGroupDestinationStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}
