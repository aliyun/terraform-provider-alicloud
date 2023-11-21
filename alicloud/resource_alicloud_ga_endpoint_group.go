package alicloud

import (
	"fmt"
	"log"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
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
			Update: schema.DefaultTimeout(2 * time.Minute),
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
			"health_check_interval_seconds": {
				Type:     schema.TypeInt,
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
				ValidateFunc: StringInSlice([]string{"http", "https", "tcp"}, false),
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
							ValidateFunc: StringInSlice([]string{"Domain", "Ip", "PublicIp", "ECS", "SLB"}, false),
						},
						"weight": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: IntBetween(0, 255),
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
	conn, err := client.NewGaplusClient()
	if err != nil {
		return WrapError(err)
	}

	request["RegionId"] = client.RegionId
	request["AcceleratorId"] = d.Get("accelerator_id")
	request["ListenerId"] = d.Get("listener_id")
	request["EndpointGroupRegion"] = d.Get("endpoint_group_region")

	if v, ok := d.GetOk("endpoint_group_type"); ok {
		request["EndpointGroupType"] = v
	}

	if v, ok := d.GetOk("endpoint_request_protocol"); ok {
		request["EndpointRequestProtocol"] = v
	}

	if v, ok := d.GetOk("health_check_interval_seconds"); ok {
		request["HealthCheckIntervalSeconds"] = v
	}

	if v, ok := d.GetOk("health_check_path"); ok {
		request["HealthCheckPath"] = v
	}

	if v, ok := d.GetOk("health_check_port"); ok {
		request["HealthCheckPort"] = v
	}

	if v, ok := d.GetOk("health_check_protocol"); ok {
		request["HealthCheckProtocol"] = v
	}

	if v, ok := d.GetOk("threshold_count"); ok {
		request["ThresholdCount"] = v
	}

	if v, ok := d.GetOk("traffic_percentage"); ok {
		request["TrafficPercentage"] = v
	}

	if v, ok := d.GetOk("name"); ok {
		request["Name"] = v
	}

	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}

	EndpointConfigurations := make([]map[string]interface{}, len(d.Get("endpoint_configurations").([]interface{})))
	for i, EndpointConfigurationsValue := range d.Get("endpoint_configurations").([]interface{}) {
		EndpointConfigurationsMap := EndpointConfigurationsValue.(map[string]interface{})
		EndpointConfigurations[i] = make(map[string]interface{})
		EndpointConfigurations[i]["Endpoint"] = EndpointConfigurationsMap["endpoint"]
		EndpointConfigurations[i]["Type"] = EndpointConfigurationsMap["type"]
		EndpointConfigurations[i]["Weight"] = EndpointConfigurationsMap["weight"]
		EndpointConfigurations[i]["EnableProxyProtocol"] = EndpointConfigurationsMap["enable_proxy_protocol"]
		EndpointConfigurations[i]["EnableClientIPPreservation"] = EndpointConfigurationsMap["enable_clientip_preservation"]
	}
	request["EndpointConfigurations"] = EndpointConfigurations

	if v, ok := d.GetOk("port_overrides"); ok {
		PortOverrides := make([]map[string]interface{}, len(v.([]interface{})))
		for i, PortOverridesValue := range v.([]interface{}) {
			PortOverridesMap := PortOverridesValue.(map[string]interface{})
			PortOverrides[i] = make(map[string]interface{})
			PortOverrides[i]["EndpointPort"] = PortOverridesMap["endpoint_port"]
			PortOverrides[i]["ListenerPort"] = PortOverridesMap["listener_port"]
		}
		request["PortOverrides"] = PortOverrides
	}

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		request["ClientToken"] = buildClientToken("CreateEndpointGroup")
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-11-20"), StringPointer("AK"), nil, request, &runtime)
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
	d.Set("description", object["Description"])

	endpointConfigurations := make([]map[string]interface{}, 0)
	if endpointConfigurationsList, ok := object["EndpointConfigurations"].([]interface{}); ok {
		for _, v := range endpointConfigurationsList {
			if m1, ok := v.(map[string]interface{}); ok {
				temp1 := map[string]interface{}{
					"endpoint":                     m1["Endpoint"],
					"type":                         m1["Type"],
					"weight":                       m1["Weight"],
					"enable_proxy_protocol":        m1["EnableProxyProtocol"],
					"enable_clientip_preservation": m1["EnableClientIPPreservation"],
				}
				endpointConfigurations = append(endpointConfigurations, temp1)

			}
		}
	}

	if err := d.Set("endpoint_configurations", endpointConfigurations); err != nil {
		return WrapError(err)
	}
	d.Set("endpoint_group_region", object["EndpointGroupRegion"])
	d.Set("health_check_interval_seconds", formatInt(object["HealthCheckIntervalSeconds"]))
	d.Set("health_check_path", object["HealthCheckPath"])
	d.Set("health_check_port", formatInt(object["HealthCheckPort"]))
	d.Set("health_check_protocol", object["HealthCheckProtocol"])
	d.Set("listener_id", object["ListenerId"])
	d.Set("name", object["Name"])
	d.Set("accelerator_id", object["AcceleratorId"])
	d.Set("endpoint_group_type", object["EndpointGroupType"])
	d.Set("endpoint_request_protocol", object["EndpointRequestProtocol"])
	portOverrides := make([]map[string]interface{}, 0)
	if portOverridesList, ok := object["PortOverrides"].([]interface{}); ok {
		for _, v := range portOverridesList {
			if m1, ok := v.(map[string]interface{}); ok {
				temp1 := map[string]interface{}{
					"endpoint_port": m1["EndpointPort"],
					"listener_port": m1["ListenerPort"],
				}
				portOverrides = append(portOverrides, temp1)

			}
		}
	}
	if err := d.Set("port_overrides", portOverrides); err != nil {
		return WrapError(err)
	}

	d.Set("endpoint_group_ip_list", object["EndpointGroupIpList"])
	d.Set("status", object["State"])
	d.Set("threshold_count", formatInt(object["ThresholdCount"]))
	d.Set("traffic_percentage", formatInt(object["TrafficPercentage"]))

	listTagResourcesObject, err := gaService.ListTagResources(d.Id(), "endpointgroup")
	if err != nil {
		return WrapError(err)
	}

	d.Set("tags", tagsToMap(listTagResourcesObject))

	return nil
}

func resourceAliCloudGaEndpointGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	conn, err := client.NewGaplusClient()
	if err != nil {
		return WrapError(err)
	}
	gaService := GaService{client}
	var response map[string]interface{}
	update := false
	request := map[string]interface{}{
		"EndpointGroupId": d.Id(),
	}

	if d.HasChange("tags") {
		if err := gaService.SetResourceTags(d, "endpointgroup"); err != nil {
			return WrapError(err)
		}
		d.SetPartial("tags")
	}

	request["RegionId"] = client.RegionId
	request["EndpointGroupRegion"] = d.Get("endpoint_group_region")

	if !d.IsNewResource() && d.HasChange("health_check_interval_seconds") {
		update = true
		request["HealthCheckIntervalSeconds"] = d.Get("health_check_interval_seconds")
	}

	if !d.IsNewResource() && d.HasChange("health_check_path") {
		update = true
		request["HealthCheckPath"] = d.Get("health_check_path")
	}

	if !d.IsNewResource() && d.HasChange("health_check_port") {
		update = true
		request["HealthCheckPort"] = d.Get("health_check_port")
	}

	if !d.IsNewResource() && d.HasChange("health_check_protocol") {
		update = true
		request["HealthCheckProtocol"] = d.Get("health_check_protocol")
	}

	if !d.IsNewResource() && d.HasChange("name") {
		update = true
		request["Name"] = d.Get("name")
	}

	if !d.IsNewResource() && d.HasChange("description") {
		update = true
		request["Description"] = d.Get("description")
	}

	if !d.IsNewResource() && d.HasChange("threshold_count") {
		update = true
	}
	if v, ok := d.GetOk("threshold_count"); ok {
		request["ThresholdCount"] = v
	}

	if !d.IsNewResource() && d.HasChange("traffic_percentage") {
		update = true

		if v, ok := d.GetOkExists("traffic_percentage"); ok {
			request["TrafficPercentage"] = v
		}
	}

	if !d.IsNewResource() && d.HasChange("endpoint_configurations") {
		update = true
	}
	EndpointConfigurations := make([]map[string]interface{}, len(d.Get("endpoint_configurations").([]interface{})))
	for i, EndpointConfigurationsValue := range d.Get("endpoint_configurations").([]interface{}) {
		EndpointConfigurationsMap := EndpointConfigurationsValue.(map[string]interface{})
		EndpointConfigurations[i] = make(map[string]interface{})
		EndpointConfigurations[i]["Endpoint"] = EndpointConfigurationsMap["endpoint"]
		EndpointConfigurations[i]["Type"] = EndpointConfigurationsMap["type"]
		EndpointConfigurations[i]["Weight"] = EndpointConfigurationsMap["weight"]
		EndpointConfigurations[i]["EnableProxyProtocol"] = EndpointConfigurationsMap["enable_proxy_protocol"]
		EndpointConfigurations[i]["EnableClientIPPreservation"] = EndpointConfigurationsMap["enable_clientip_preservation"]
	}
	request["EndpointConfigurations"] = EndpointConfigurations

	if !d.IsNewResource() && d.HasChange("port_overrides") {
		update = true
	}
	PortOverrides := make([]map[string]interface{}, len(d.Get("port_overrides").([]interface{})))
	for i, PortOverridesValue := range d.Get("port_overrides").([]interface{}) {
		PortOverridesMap := PortOverridesValue.(map[string]interface{})
		PortOverrides[i] = make(map[string]interface{})
		PortOverrides[i]["EndpointPort"] = PortOverridesMap["endpoint_port"]
		PortOverrides[i]["ListenerPort"] = PortOverridesMap["listener_port"]
	}
	request["PortOverrides"] = PortOverrides

	if update {
		if v, ok := d.GetOk("endpoint_request_protocol"); ok {
			request["EndpointRequestProtocol"] = v
		}

		action := "UpdateEndpointGroup"
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			request["ClientToken"] = buildClientToken("UpdateEndpointGroup")
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-11-20"), StringPointer("AK"), nil, request, &runtime)
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

		d.SetPartial("health_check_interval_seconds")
		d.SetPartial("health_check_path")
		d.SetPartial("health_check_port")
		d.SetPartial("health_check_protocol")
		d.SetPartial("name")
		d.SetPartial("description")
		d.SetPartial("threshold_count")
		d.SetPartial("traffic_percentage")
		d.SetPartial("endpoint_configurations")
		d.SetPartial("port_overrides")
	}

	return resourceAliCloudGaEndpointGroupRead(d, meta)
}

func resourceAliCloudGaEndpointGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := GaService{client}
	action := "DeleteEndpointGroup"
	var response map[string]interface{}
	conn, err := client.NewGaplusClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"EndpointGroupId": d.Id(),
	}

	request["AcceleratorId"] = d.Get("accelerator_id")
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		request["ClientToken"] = buildClientToken("DeleteEndpointGroup")
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-11-20"), StringPointer("AK"), nil, request, &runtime)
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
