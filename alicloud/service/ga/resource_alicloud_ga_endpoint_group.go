package ga

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/tags"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/validators"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceAliCloudGaEndpointGroup() *schema.Resource {
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
				ValidateFunc: validators.StringInSlice([]string{"default", "virtual"}, false),
			},
			"endpoint_request_protocol": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validators.StringInSlice([]string{"HTTP", "HTTPS"}, false),
			},
			"endpoint_protocol_version": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validators.StringInSlice([]string{"HTTP1.1", "HTTP2"}, false),
			},
			"endpoint_ip_version": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validators.StringInSlice([]string{"IPv4", "IPv6", "ProtocolAffinity"}, false),
			},
			"health_check_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"health_check_host": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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
				ValidateFunc: validators.StringInSlice([]string{"TCP", "HTTP", "HTTPS", "tcp", "http", "https"}, false),
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
			"enable": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"access_log_record_customized_header_list": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 100,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"tags": tags.TagsSchema(),
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
							ValidateFunc: validators.StringInSlice([]string{"Domain", "Ip", "IpTarget", "PublicIp", "ECS", "SLB", "ALB", "NLB", "ENI", "OSS"}, false),
						},
						"weight": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validators.IntBetween(0, 255),
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
						"probe_port": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validators.IntBetween(1, 65535),
						},
						"probe_protocol": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validators.StringInSlice([]string{"tcp", "icmp"}, false),
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
	gaService := Service{client}
	var response map[string]interface{}
	action := "CreateEndpointGroup"
	request := make(map[string]interface{})
	var err error

	request["RegionId"] = client.RegionId
	request["ClientToken"] = helper.BuildClientToken("CreateEndpointGroup")
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

	if v, ok := d.GetOk("endpoint_ip_version"); ok {
		request["EndpointIpVersion"] = v
	}

	if v, ok := d.GetOkExists("health_check_enabled"); ok {
		request["HealthCheckEnabled"] = v
	}

	if v, ok := d.GetOk("health_check_path"); ok {
		request["HealthCheckPath"] = v
	}

	if v, ok := d.GetOk("health_check_host"); ok {
		request["HealthCheckHost"] = v
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

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *retry.RetryError {
		response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"GA_NOT_STEADY", "StateError.Accelerator", "StateError.EndPointGroup", "NotActive.Listener"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_ga_endpoint_group", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["EndpointGroupId"]))

	stateConf := helper.BuildStateConf([]string{}, []string{"active"}, d.Timeout(schema.TimeoutCreate), 30*time.Second, gaService.GaEndpointGroupStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	if err := resourceAliCloudGaEndpointGroupUpdateProbe(d, client); err != nil {
		return errs.WrapError(err)
	}

	if err := resourceAliCloudGaEndpointGroupUpdateAccessLog(d, client); err != nil {
		return errs.WrapError(err)
	}

	return resourceAliCloudGaEndpointGroupUpdate(d, meta)
}

func resourceAliCloudGaEndpointGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := Service{client}

	object, err := gaService.DescribeGaEndpointGroup(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ga_endpoint_group gaService.DescribeGaEndpointGroup Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("accelerator_id", object["AcceleratorId"])
	d.Set("listener_id", object["ListenerId"])
	d.Set("endpoint_group_region", object["EndpointGroupRegion"])
	d.Set("endpoint_group_type", object["EndpointGroupType"])
	d.Set("endpoint_request_protocol", object["EndpointRequestProtocol"])
	d.Set("endpoint_protocol_version", object["EndpointProtocolVersion"])
	d.Set("endpoint_ip_version", object["EndpointIpVersion"])
	d.Set("health_check_enabled", object["HealthCheckEnabled"])
	d.Set("health_check_path", object["HealthCheckPath"])
	d.Set("health_check_host", object["HealthCheckHost"])
	health_check_portVal, err := convert.FormatInt(object["HealthCheckPort"])
	if err != nil {
		return err
	}
	d.Set("health_check_port", health_check_portVal)
	d.Set("health_check_protocol", object["HealthCheckProtocol"])
	health_check_interval_secondsVal, err := convert.FormatInt(object["HealthCheckIntervalSeconds"])
	if err != nil {
		return err
	}
	d.Set("health_check_interval_seconds", health_check_interval_secondsVal)
	threshold_countVal, err := convert.FormatInt(object["ThresholdCount"])
	if err != nil {
		return err
	}
	d.Set("threshold_count", threshold_countVal)
	traffic_percentageVal, err := convert.FormatInt(object["TrafficPercentage"])
	if err != nil {
		return err
	}
	d.Set("traffic_percentage", traffic_percentageVal)
	d.Set("name", object["Name"])
	d.Set("description", object["Description"])
	d.Set("access_log_record_customized_header_list", object["AccessLogRecordCustomizedHeaderList"])

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

			if probePort, ok := endpointConfigurationsArg["ProbePort"]; ok {
				probePortVal, err := convert.FormatInt(probePort)
				if err != nil {
					return err
				}
				endpointConfigurationsMap["probe_port"] = probePortVal
			}

			if probeProtocol, ok := endpointConfigurationsArg["ProbeProtocol"]; ok {
				endpointConfigurationsMap["probe_protocol"] = probeProtocol
			}

			endpointConfigurationsMaps = append(endpointConfigurationsMaps, endpointConfigurationsMap)
		}

		if err := d.Set("endpoint_configurations", endpointConfigurationsMaps); err != nil {
			return errs.WrapError(err)
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
			return errs.WrapError(err)
		}
	}

	d.Set("endpoint_group_ip_list", object["EndpointGroupIpList"])
	d.Set("status", object["State"])

	listTagResourcesObject, err := gaService.ListTagResources(d.Id(), "endpointgroup")
	if err != nil {
		return errs.WrapError(err)
	}

	d.Set("tags", tags.TagsToMap(listTagResourcesObject))

	return nil
}

func resourceAliCloudGaEndpointGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := Service{client}
	var response map[string]interface{}
	d.Partial(true)

	if d.HasChange("tags") {
		if err := gaService.SetResourceTags(d, "endpointgroup"); err != nil {
			return errs.WrapError(err)
		}

	}

	update := false
	request := map[string]interface{}{
		"RegionId":            client.RegionId,
		"ClientToken":         helper.BuildClientToken("UpdateEndpointGroup"),
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

	if !d.IsNewResource() && d.HasChange("endpoint_ip_version") {
		update = true
	}
	if v, ok := d.GetOk("endpoint_ip_version"); ok {
		request["EndpointIpVersion"] = v
	}

	if !d.IsNewResource() && d.HasChange("health_check_host") {
		update = true
	}
	if v, ok := d.GetOk("health_check_host"); ok {
		request["HealthCheckHost"] = v
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

		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *retry.RetryError {
			response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
			if err != nil {
				if errs.IsExpectedErrors(err, []string{"StateError.Accelerator", "StateError.EndPointGroup", "NotActive.Listener"}) || helper.NeedRetry(err) {
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

		stateConf := helper.BuildStateConf([]string{}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 30*time.Second, gaService.GaEndpointGroupStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
		}

	}

	if !d.IsNewResource() {
		if err := resourceAliCloudGaEndpointGroupUpdateProbe(d, client); err != nil {
			return errs.WrapError(err)
		}

		if err := resourceAliCloudGaEndpointGroupUpdateAccessLog(d, client); err != nil {
			return errs.WrapError(err)
		}
	}

	d.Partial(false)

	return resourceAliCloudGaEndpointGroupRead(d, meta)
}

func resourceAliCloudGaEndpointGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := Service{client}
	action := "DeleteEndpointGroup"
	var response map[string]interface{}

	var err error

	request := map[string]interface{}{
		"ClientToken":     helper.BuildClientToken("DeleteEndpointGroup"),
		"AcceleratorId":   d.Get("accelerator_id"),
		"EndpointGroupId": d.Id(),
	}

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *retry.RetryError {
		response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"StateError.Accelerator", "StateError.EndPointGroup", "NotActive.Listener"}) || helper.NeedRetry(err) {
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

	stateConf := helper.BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, gaService.GaEndpointGroupStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return nil
}

func resourceAliCloudGaEndpointGroupUpdateProbe(d *schema.ResourceData, client *connectivity.AliyunClient) error {
	if !d.HasChange("enable") && !d.HasChange("endpoint_configurations") {
		return nil
	}

	enabled := strconv.FormatBool(d.Get("enable").(bool))
	gaService := Service{client}
	endpointConfigurations, _ := d.Get("endpoint_configurations").([]interface{})
	for _, item := range endpointConfigurations {
		endpointConfiguration, ok := item.(map[string]interface{})
		if !ok || endpointConfiguration == nil {
			continue
		}

		probePort, _ := endpointConfiguration["probe_port"].(int)
		probeProtocol, _ := endpointConfiguration["probe_protocol"].(string)
		if !d.HasChange("enable") && probePort == 0 && probeProtocol == "" {
			continue
		}

		endpointType, _ := endpointConfiguration["type"].(string)
		switch endpointType {
		case "Ip", "Domain", "PublicIp":
		default:
			return fmt.Errorf("alicloud_ga_endpoint_group: delay detection is supported only for endpoints of type Ip, Domain or PublicIp, but the endpoint %v uses type %s", endpointConfiguration["endpoint"], endpointType)
		}

		action := "ConfigEndpointProbe"
		request := map[string]interface{}{
			"RegionId":        client.RegionId,
			"ClientToken":     helper.BuildClientToken(action),
			"EndpointGroupId": d.Id(),
			"EndpointType":    endpointConfiguration["type"],
			"Endpoint":        endpointConfiguration["endpoint"],
			"Enable":          enabled,
		}
		if probePort != 0 {
			request["ProbePort"] = strconv.Itoa(probePort)
		}
		if probeProtocol != "" {
			request["ProbeProtocol"] = probeProtocol
		}

		var response map[string]interface{}
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err := retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *retry.RetryError {
			var err error
			response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
			if err != nil {
				if errs.IsExpectedErrors(err, []string{"StateError.Accelerator", "StateError.EndPointGroup", "NotActive.Listener"}) || helper.NeedRetry(err) {
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

		stateConf := helper.BuildStateConf([]string{}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 30*time.Second, gaService.GaEndpointGroupStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
		}
	}

	return nil
}

func resourceAliCloudGaEndpointGroupUpdateAccessLog(d *schema.ResourceData, client *connectivity.AliyunClient) error {
	if !d.HasChange("access_log_record_customized_header_list") {
		return nil
	}

	headerList := convert.ExpandStringList(d.Get("access_log_record_customized_header_list").(*schema.Set).List())
	gaService := Service{client}
	object, err := gaService.DescribeGaEndpointGroup(d.Id())
	if err != nil {
		return err
	}

	slsProjectName, _ := object["SlsProjectName"].(string)
	slsLogStoreName, _ := object["SlsLogStoreName"].(string)
	if slsProjectName == "" || slsLogStoreName == "" {
		if len(headerList) == 0 {
			return nil
		}
		return fmt.Errorf("alicloud_ga_endpoint_group: access_log_record_customized_header_list requires the access log of endpoint group %s to be enabled; attach it with alicloud_ga_access_log first", d.Id())
	}

	action := "UpdateLogStoreConfig"
	request := map[string]interface{}{
		"RegionId":        client.RegionId,
		"ClientToken":     helper.BuildClientToken(action),
		"EndpointGroupId": d.Id(),
		"AcceleratorId":   object["AcceleratorId"],
		"ListenerId":      object["ListenerId"],
		"SlsProjectName":  slsProjectName,
		"SlsLogStoreName": slsLogStoreName,
		"AccessLogRecordCustomizedHeadersEnabled": len(headerList) > 0,
		"AccessLogRecordCustomizedHeaderList":     headerList,
	}

	var response map[string]interface{}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *retry.RetryError {
		var err error
		response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"StateError.Accelerator", "StateError.EndPointGroup", "NotActive.Listener"}) || helper.NeedRetry(err) {
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

	return nil
}
