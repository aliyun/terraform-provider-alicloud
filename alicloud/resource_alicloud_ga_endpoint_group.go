package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudGaEndpointGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudGaEndpointGroupCreate,
		Read:   resourceAliCloudGaEndpointGroupRead,
		Update: resourceAliCloudGaEndpointGroupUpdate,
		Delete: resourceAliCloudGaEndpointGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(15 * time.Minute),
			Update: schema.DefaultTimeout(3 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"accelerator_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"listener_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"endpoint_group_region": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"endpoint_group_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"default", "virtual"}, false),
			},
			"endpoint_request_protocol": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"HTTP", "HTTPS"}, false),
			},
			"endpoint_protocol_version": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"HTTP1.1", "HTTP2"}, false),
			},
			"health_check_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"health_check_path": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"health_check_port": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"health_check_protocol": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"TCP", "HTTP", "HTTPS", "tcp", "http", "https"}, false),
			},
			"health_check_interval_seconds": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"threshold_count": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"traffic_percentage": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": tagsSchema(),
			"endpoint_configurations": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"endpoint": {
							Type:     schema.TypeString,
							Required: true,
						},
						"type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: StringInSlice([]string{"Domain", "Ip", "IpTarget", "PublicIp", "ECS", "SLB", "ALB", "NLB", "ENI", "OSS"}, false),
						},
						"weight": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: IntBetween(0, 255),
						},
						"sub_address": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"enable_proxy_protocol": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"enable_clientip_preservation": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"vswitch_ids": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"port_overrides": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"endpoint_port": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"listener_port": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
			"endpoint_group_ip_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudGaEndpointGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := GaService{client}
	var response map[string]interface{}
	action := "CreateEndpointGroup"
	request := make(map[string]interface{})
	var err error

	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken("CreateEndpointGroup")
	request["AcceleratorId"] = d.Get("accelerator_id")
	request["ListenerId"] = d.Get("listener_id")
	request["EndpointGroupRegion"] = d.Get("endpoint_group_region")

	if v, ok := d.GetOk("endpoint_group_type"); ok {
		request["EndpointGroupType"] = v
	}

	if v, ok := d.GetOk("endpoint_request_protocol"); ok {
		request["EndpointRequestProtocol"] = v
	}

	if v, ok := d.GetOk("endpoint_protocol_version"); ok {
		request["EndpointProtocolVersion"] = v
	}

	if v, ok := d.GetOkExists("health_check_enabled"); ok {
		request["HealthCheckEnabled"] = v
	}

	if v, ok := d.GetOk("health_check_path"); ok {
		request["HealthCheckPath"] = v
	}

	if v, ok := d.GetOkExists("health_check_port"); ok {
		request["HealthCheckPort"] = v
	}

	if v, ok := d.GetOk("health_check_protocol"); ok {
		request["HealthCheckProtocol"] = v
	}

	if v, ok := d.GetOkExists("health_check_interval_seconds"); ok {
		request["HealthCheckIntervalSeconds"] = v
	}

	if v, ok := d.GetOkExists("threshold_count"); ok {
		request["ThresholdCount"] = v
	}

	if v, ok := d.GetOkExists("traffic_percentage"); ok {
		request["TrafficPercentage"] = v
	}

	if v, ok := d.GetOk("name"); ok {
		request["Name"] = v
	}

	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}

	endpointConfigurations := d.Get("endpoint_configurations")
	endpointConfigurationsMaps := make([]map[string]interface{}, 0)
	for _, endpointConfigurationsList := range endpointConfigurations.([]interface{}) {
		endpointConfigurationsMap := map[string]interface{}{}
		endpointConfigurationsArg := endpointConfigurationsList.(map[string]interface{})

		endpointConfigurationsMap["Endpoint"] = endpointConfigurationsArg["endpoint"]
		endpointConfigurationsMap["Type"] = endpointConfigurationsArg["type"]
		endpointConfigurationsMap["Weight"] = endpointConfigurationsArg["weight"]

		if subAddress, ok := endpointConfigurationsArg["sub_address"]; ok {
			endpointConfigurationsMap["SubAddress"] = subAddress
		}

		if enableProxyProtocol, ok := endpointConfigurationsArg["enable_proxy_protocol"]; ok {
			endpointConfigurationsMap["EnableProxyProtocol"] = enableProxyProtocol
		}

		if enableClientIPPreservation, ok := endpointConfigurationsArg["enable_clientip_preservation"]; ok {
			endpointConfigurationsMap["EnableClientIPPreservation"] = enableClientIPPreservation
		}

		if vpcId, ok := endpointConfigurationsArg["vpc_id"]; ok {
			endpointConfigurationsMap["VpcId"] = vpcId
		}

		if vSwitchIds, ok := endpointConfigurationsArg["vswitch_ids"]; ok {
			endpointConfigurationsMap["VSwitchIds"] = vSwitchIds
		}

		endpointConfigurationsMaps = append(endpointConfigurationsMaps, endpointConfigurationsMap)
	}

	request["EndpointConfigurations"] = endpointConfigurationsMaps

	if v, ok := d.GetOk("port_overrides"); ok {
		portOverridesMaps := make([]map[string]interface{}, 0)
		for _, portOverrides := range v.([]interface{}) {
			portOverridesMap := map[string]interface{}{}
			portOverridesArg := portOverrides.(map[string]interface{})

			if endpointPort, ok := portOverridesArg["endpoint_port"]; ok {
				portOverridesMap["EndpointPort"] = endpointPort
			}

			if listenerPort, ok := portOverridesArg["listener_port"]; ok {
				portOverridesMap["ListenerPort"] = listenerPort
			}

			portOverridesMaps = append(portOverridesMaps, portOverridesMap)
		}

		request["PortOverrides"] = portOverridesMaps
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"GA_NOT_STEADY", "StateError.Accelerator", "StateError.EndPointGroup", "NotActive.Listener"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ga_endpoint_group", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["EndpointGroupId"]))

	stateConf := BuildStateConf([]string{}, []string{"active"}, d.Timeout(schema.TimeoutCreate), 30*time.Second, gaService.GaEndpointGroupStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudGaEndpointGroupUpdate(d, meta)
}

func resourceAliCloudGaEndpointGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := GaService{client}

	object, err := gaService.DescribeGaEndpointGroup(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ga_endpoint_group gaService.DescribeGaEndpointGroup Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("accelerator_id", object["AcceleratorId"])
	d.Set("listener_id", object["ListenerId"])
	d.Set("endpoint_group_region", object["EndpointGroupRegion"])
	d.Set("endpoint_group_type", object["EndpointGroupType"])
	d.Set("endpoint_request_protocol", object["EndpointRequestProtocol"])
	d.Set("endpoint_protocol_version", object["EndpointProtocolVersion"])
	d.Set("health_check_enabled", object["HealthCheckEnabled"])
	d.Set("health_check_path", object["HealthCheckPath"])
	d.Set("health_check_port", formatInt(object["HealthCheckPort"]))
	d.Set("health_check_protocol", object["HealthCheckProtocol"])
	d.Set("health_check_interval_seconds", formatInt(object["HealthCheckIntervalSeconds"]))
	d.Set("threshold_count", formatInt(object["ThresholdCount"]))
	d.Set("traffic_percentage", formatInt(object["TrafficPercentage"]))
	d.Set("name", object["Name"])
	d.Set("description", object["Description"])

	if endpointConfigurations, ok := object["EndpointConfigurations"]; ok {
		endpointConfigurationsMaps := make([]map[string]interface{}, 0)
		for _, endpointConfigurationsList := range endpointConfigurations.([]interface{}) {
			endpointConfigurationsArg := endpointConfigurationsList.(map[string]interface{})
			endpointConfigurationsMap := make(map[string]interface{})

			if endpoint, ok := endpointConfigurationsArg["Endpoint"]; ok {
				endpointConfigurationsMap["endpoint"] = endpoint
			}

			if endpointType, ok := endpointConfigurationsArg["Type"]; ok {
				endpointConfigurationsMap["type"] = endpointType
			}

			if weight, ok := endpointConfigurationsArg["Weight"]; ok {
				endpointConfigurationsMap["weight"] = weight
			}

			if subAddress, ok := endpointConfigurationsArg["SubAddress"]; ok {
				endpointConfigurationsMap["sub_address"] = subAddress
			}

			if enableProxyProtocol, ok := endpointConfigurationsArg["EnableProxyProtocol"]; ok {
				endpointConfigurationsMap["enable_proxy_protocol"] = enableProxyProtocol
			}

			if enableClientIPPreservation, ok := endpointConfigurationsArg["EnableClientIPPreservation"]; ok {
				endpointConfigurationsMap["enable_clientip_preservation"] = enableClientIPPreservation
			}

			if vpcId, ok := endpointConfigurationsArg["VpcId"]; ok {
				endpointConfigurationsMap["vpc_id"] = vpcId
			}

			if vSwitchIds, ok := endpointConfigurationsArg["VSwitchIds"]; ok {
				endpointConfigurationsMap["vswitch_ids"] = vSwitchIds
			}

			endpointConfigurationsMaps = append(endpointConfigurationsMaps, endpointConfigurationsMap)
		}

		if err := d.Set("endpoint_configurations", endpointConfigurationsMaps); err != nil {
			return WrapError(err)
		}
	}

	if portOverrides, ok := object["PortOverrides"]; ok {
		portOverridesMaps := make([]map[string]interface{}, 0)
		for _, portOverridesList := range portOverrides.([]interface{}) {
			portOverridesArg := portOverridesList.(map[string]interface{})
			portOverridesMap := make(map[string]interface{})

			if endpointPort, ok := portOverridesArg["EndpointPort"]; ok {
				portOverridesMap["endpoint_port"] = endpointPort
			}

			if listenerPort, ok := portOverridesArg["ListenerPort"]; ok {
				portOverridesMap["listener_port"] = listenerPort
			}

			portOverridesMaps = append(portOverridesMaps, portOverridesMap)
		}

		if err := d.Set("port_overrides", portOverridesMaps); err != nil {
			return WrapError(err)
		}
	}

	d.Set("endpoint_group_ip_list", object["EndpointGroupIpList"])
	d.Set("status", object["State"])

	listTagResourcesObject, err := gaService.ListTagResources(d.Id(), "endpointgroup")
	if err != nil {
		return WrapError(err)
	}

	d.Set("tags", tagsToMap(listTagResourcesObject))

	return nil
}

func resourceAliCloudGaEndpointGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := GaService{client}
	var response map[string]interface{}
	d.Partial(true)

	if d.HasChange("tags") {
		if err := gaService.SetResourceTags(d, "endpointgroup"); err != nil {
			return WrapError(err)
		}

		d.SetPartial("tags")
	}

	update := false
	request := map[string]interface{}{
		"RegionId":            client.RegionId,
		"ClientToken":         buildClientToken("UpdateEndpointGroup"),
		"EndpointGroupRegion": d.Get("endpoint_group_region"),
		"EndpointGroupId":     d.Id(),
	}

	if !d.IsNewResource() && d.HasChange("endpoint_request_protocol") {
		update = true
	}
	if v, ok := d.GetOk("endpoint_request_protocol"); ok {
		request["EndpointRequestProtocol"] = v
	}

	if !d.IsNewResource() && d.HasChange("endpoint_protocol_version") {
		update = true
	}
	if v, ok := d.GetOk("endpoint_protocol_version"); ok {
		request["EndpointProtocolVersion"] = v
	}

	if !d.IsNewResource() && d.HasChange("health_check_enabled") {
		update = true
	}
	if v, ok := d.GetOkExists("health_check_enabled"); ok {
		request["HealthCheckEnabled"] = v
	}

	if !d.IsNewResource() && d.HasChange("health_check_path") {
		update = true
	}
	if v, ok := d.GetOk("health_check_path"); ok {
		request["HealthCheckPath"] = v
	}

	if !d.IsNewResource() && d.HasChange("health_check_port") {
		update = true

		if v, ok := d.GetOkExists("health_check_port"); ok {
			request["HealthCheckPort"] = v
		}
	}

	if !d.IsNewResource() && d.HasChange("health_check_protocol") {
		update = true
	}
	if v, ok := d.GetOk("health_check_protocol"); ok {
		request["HealthCheckProtocol"] = v
	}

	if !d.IsNewResource() && d.HasChange("health_check_interval_seconds") {
		update = true

		if v, ok := d.GetOkExists("health_check_interval_seconds"); ok {
			request["HealthCheckIntervalSeconds"] = v
		}
	}

	if !d.IsNewResource() && d.HasChange("threshold_count") {
		update = true
	}
	if v, ok := d.GetOkExists("threshold_count"); ok {
		request["ThresholdCount"] = v
	}

	if !d.IsNewResource() && d.HasChange("traffic_percentage") {
		update = true

		if v, ok := d.GetOkExists("traffic_percentage"); ok {
			request["TrafficPercentage"] = v
		}
	}

	if !d.IsNewResource() && d.HasChange("name") {
		update = true
	}
	if v, ok := d.GetOk("name"); ok {
		request["Name"] = v
	}

	if !d.IsNewResource() && d.HasChange("description") {
		update = true
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}

	if !d.IsNewResource() && d.HasChange("endpoint_configurations") {
		update = true
	}
	endpointConfigurations := d.Get("endpoint_configurations")
	endpointConfigurationsMaps := make([]map[string]interface{}, 0)
	for _, endpointConfigurationsList := range endpointConfigurations.([]interface{}) {
		endpointConfigurationsMap := map[string]interface{}{}
		endpointConfigurationsArg := endpointConfigurationsList.(map[string]interface{})

		endpointConfigurationsMap["Endpoint"] = endpointConfigurationsArg["endpoint"]
		endpointConfigurationsMap["Type"] = endpointConfigurationsArg["type"]
		endpointConfigurationsMap["Weight"] = endpointConfigurationsArg["weight"]

		if subAddress, ok := endpointConfigurationsArg["sub_address"]; ok {
			endpointConfigurationsMap["SubAddress"] = subAddress
		}

		if enableProxyProtocol, ok := endpointConfigurationsArg["enable_proxy_protocol"]; ok {
			endpointConfigurationsMap["EnableProxyProtocol"] = enableProxyProtocol
		}

		if enableClientIPPreservation, ok := endpointConfigurationsArg["enable_clientip_preservation"]; ok {
			endpointConfigurationsMap["EnableClientIPPreservation"] = enableClientIPPreservation
		}

		if vpcId, ok := endpointConfigurationsArg["vpc_id"]; ok {
			endpointConfigurationsMap["VpcId"] = vpcId
		}

		if vSwitchIds, ok := endpointConfigurationsArg["vswitch_ids"]; ok {
			endpointConfigurationsMap["VSwitchIds"] = vSwitchIds
		}

		endpointConfigurationsMaps = append(endpointConfigurationsMaps, endpointConfigurationsMap)
	}

	request["EndpointConfigurations"] = endpointConfigurationsMaps

	if !d.IsNewResource() && d.HasChange("port_overrides") {
		update = true
	}
	if v, ok := d.GetOk("port_overrides"); ok {
		portOverridesMaps := make([]map[string]interface{}, 0)
		for _, portOverrides := range v.([]interface{}) {
			portOverridesMap := map[string]interface{}{}
			portOverridesArg := portOverrides.(map[string]interface{})

			if endpointPort, ok := portOverridesArg["endpoint_port"]; ok {
				portOverridesMap["EndpointPort"] = endpointPort
			}

			if listenerPort, ok := portOverridesArg["listener_port"]; ok {
				portOverridesMap["ListenerPort"] = listenerPort
			}

			portOverridesMaps = append(portOverridesMaps, portOverridesMap)
		}

		request["PortOverrides"] = portOverridesMaps
	}

	if update {
		action := "UpdateEndpointGroup"
		var err error

		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"StateError.Accelerator", "StateError.EndPointGroup", "NotActive.Listener"}) || NeedRetry(err) {
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

		stateConf := BuildStateConf([]string{}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 30*time.Second, gaService.GaEndpointGroupStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}

		d.SetPartial("endpoint_request_protocol")
		d.SetPartial("endpoint_protocol_version")
		d.SetPartial("health_check_enabled")
		d.SetPartial("health_check_path")
		d.SetPartial("health_check_port")
		d.SetPartial("health_check_protocol")
		d.SetPartial("health_check_interval_seconds")
		d.SetPartial("threshold_count")
		d.SetPartial("traffic_percentage")
		d.SetPartial("name")
		d.SetPartial("description")
		d.SetPartial("endpoint_configurations")
		d.SetPartial("port_overrides")
	}

	d.Partial(false)

	return resourceAliCloudGaEndpointGroupRead(d, meta)
}

func resourceAliCloudGaEndpointGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := GaService{client}
	action := "DeleteEndpointGroup"
	var response map[string]interface{}

	var err error

	request := map[string]interface{}{
		"ClientToken":     buildClientToken("DeleteEndpointGroup"),
		"AcceleratorId":   d.Get("accelerator_id"),
		"EndpointGroupId": d.Id(),
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"StateError.Accelerator", "StateError.EndPointGroup", "NotActive.Listener"}) || NeedRetry(err) {
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

	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, gaService.GaEndpointGroupStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}
