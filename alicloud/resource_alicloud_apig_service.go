// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudApigService() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudApigServiceCreate,
		Read:   resourceAliCloudApigServiceRead,
		Update: resourceAliCloudApigServiceUpdate,
		Delete: resourceAliCloudApigServiceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(6 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"addresses": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"create_timestamp": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"dns_servers": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"express_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"gateway_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"health_check_config": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"http_path": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"unhealthy_threshold": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"timeout": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"http_host": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"healthy_threshold": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"enable": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"expected_statuses": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"protocol": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"interval": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
			"health_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"healthy_panic_threshold": {
				Type:     schema.TypeFloat,
				Optional: true,
				Computed: true,
			},
			"namespace": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"outlier_detection_config": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"failure_percentage_threshold": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"enable": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"base_ejection_time": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"failure_percentage_minimum_hosts": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"interval": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
			"outlier_endpoints": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"ports": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"protocol": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"protocol": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"qualifier": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"runtime_detail_error_code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"runtime_detail_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"service_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"source_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"DNS", "VIP", "FC3"}, false),
			},
			"unhealthy_endpoints": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"update_timestamp": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudApigServiceCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := fmt.Sprintf("/v1/services")
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	serviceConfigsDataList := make(map[string]interface{})
	if v, ok := d.GetOk("namespace"); ok {
		serviceConfigsDataList["namespace"] = v
	}
	if v, ok := d.GetOk("dns_servers"); ok {
		dnsServers1, _ := jsonpath.Get("$", v)
		if dnsServers1 != nil && dnsServers1 != "" {
			serviceConfigsDataList["dnsServers"] = dnsServers1
		}
	}
	if v, ok := d.GetOk("qualifier"); ok {
		serviceConfigsDataList["qualifier"] = v
	}
	if v, ok := d.GetOk("addresses"); ok {
		addresses1, _ := jsonpath.Get("$", v)
		if addresses1 != nil && addresses1 != "" {
			serviceConfigsDataList["addresses"] = addresses1
		}
	}
	if v, ok := d.GetOk("express_type"); ok {
		serviceConfigsDataList["expressType"] = v
	}
	if v, ok := d.GetOk("service_name"); ok {
		serviceConfigsDataList["name"] = v
	}

	serviceConfigsMap := make([]interface{}, 0)
	serviceConfigsMap = append(serviceConfigsMap, serviceConfigsDataList)
	request["serviceConfigs"] = serviceConfigsMap

	if v, ok := d.GetOk("resource_group_id"); ok {
		request["resourceGroupId"] = v
	}
	if v, ok := d.GetOk("source_type"); ok {
		request["sourceType"] = v
	}
	if v, ok := d.GetOk("gateway_id"); ok {
		request["gatewayId"] = v
	}
	body = request
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RoaPost("APIG", "2024-03-27", action, query, nil, body, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_apig_service", action, AlibabaCloudSdkGoERROR)
	}

	id, _ := jsonpath.Get("$.data.serviceIds[0]", response)
	d.SetId(fmt.Sprint(id))

	return resourceAliCloudApigServiceUpdate(d, meta)
}

func resourceAliCloudApigServiceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	apigServiceV2 := ApigServiceV2{client}

	objectRaw, err := apigServiceV2.DescribeApigService(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_apig_service DescribeApigService Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("create_timestamp", objectRaw["createTimestamp"])
	d.Set("express_type", objectRaw["expressType"])
	d.Set("gateway_id", objectRaw["gatewayId"])
	d.Set("health_status", objectRaw["healthStatus"])
	d.Set("healthy_panic_threshold", objectRaw["healthyPanicThreshold"])
	d.Set("namespace", objectRaw["namespace"])
	d.Set("qualifier", objectRaw["qualifier"])
	d.Set("resource_group_id", objectRaw["resourceGroupId"])
	d.Set("runtime_detail_error_code", objectRaw["runtimeDetailErrorCode"])
	d.Set("runtime_detail_status", objectRaw["runtimeDetailStatus"])
	d.Set("service_name", objectRaw["name"])
	d.Set("source_type", objectRaw["sourceType"])
	d.Set("update_timestamp", objectRaw["updateTimestamp"])
	d.Set("protocol", objectRaw["protocol"])

	addressesRaw := make([]interface{}, 0)
	if objectRaw["addresses"] != nil {
		addressesRaw = convertToInterfaceArray(objectRaw["addresses"])
	}

	d.Set("addresses", addressesRaw)
	dnsServersRaw := make([]interface{}, 0)
	if objectRaw["dnsServers"] != nil {
		dnsServersRaw = convertToInterfaceArray(objectRaw["dnsServers"])
	}

	d.Set("dns_servers", dnsServersRaw)
	healthCheckConfigMaps := make([]map[string]interface{}, 0)
	healthCheckConfigMap := make(map[string]interface{})
	healthCheckRaw := make(map[string]interface{})
	if objectRaw["healthCheck"] != nil {
		healthCheckRaw = objectRaw["healthCheck"].(map[string]interface{})
	}
	if len(healthCheckRaw) > 0 {
		healthCheckConfigMap["enable"] = healthCheckRaw["enable"]
		healthCheckConfigMap["healthy_threshold"] = healthCheckRaw["healthyThreshold"]
		healthCheckConfigMap["http_host"] = healthCheckRaw["httpHost"]
		healthCheckConfigMap["http_path"] = healthCheckRaw["httpPath"]
		healthCheckConfigMap["interval"] = healthCheckRaw["interval"]
		healthCheckConfigMap["protocol"] = healthCheckRaw["protocol"]
		healthCheckConfigMap["timeout"] = healthCheckRaw["timeout"]
		healthCheckConfigMap["unhealthy_threshold"] = healthCheckRaw["unhealthyThreshold"]

		expectedStatusesRaw := make([]interface{}, 0)
		if healthCheckRaw["expectedStatuses"] != nil {
			expectedStatusesRaw = convertToInterfaceArray(healthCheckRaw["expectedStatuses"])
		}
		healthCheckConfigMap["expected_statuses"] = expectedStatusesRaw

		healthCheckConfigMaps = append(healthCheckConfigMaps, healthCheckConfigMap)
	}
	if err := d.Set("health_check_config", healthCheckConfigMaps); err != nil {
		return err
	}
	outlierDetectionConfigMaps := make([]map[string]interface{}, 0)
	outlierDetectionConfigMap := make(map[string]interface{})
	outlierDetectionRaw := make(map[string]interface{})
	if objectRaw["outlierDetection"] != nil {
		outlierDetectionRaw = objectRaw["outlierDetection"].(map[string]interface{})
	}
	if len(outlierDetectionRaw) > 0 {
		outlierDetectionConfigMap["base_ejection_time"] = outlierDetectionRaw["baseEjectionTime"]
		outlierDetectionConfigMap["enable"] = outlierDetectionRaw["enable"]
		outlierDetectionConfigMap["failure_percentage_minimum_hosts"] = outlierDetectionRaw["failurePercentageMinimumHosts"]
		outlierDetectionConfigMap["failure_percentage_threshold"] = outlierDetectionRaw["failurePercentageThreshold"]
		outlierDetectionConfigMap["interval"] = outlierDetectionRaw["interval"]

		outlierDetectionConfigMaps = append(outlierDetectionConfigMaps, outlierDetectionConfigMap)
	}
	if err := d.Set("outlier_detection_config", outlierDetectionConfigMaps); err != nil {
		return err
	}
	outlierEndpointsRaw := make([]interface{}, 0)
	if objectRaw["outlierEndpoints"] != nil {
		outlierEndpointsRaw = convertToInterfaceArray(objectRaw["outlierEndpoints"])
	}

	d.Set("outlier_endpoints", outlierEndpointsRaw)
	unhealthyEndpointsRaw := make([]interface{}, 0)
	if objectRaw["unhealthyEndpoints"] != nil {
		unhealthyEndpointsRaw = convertToInterfaceArray(objectRaw["unhealthyEndpoints"])
	}

	d.Set("unhealthy_endpoints", unhealthyEndpointsRaw)
	portsRaw := objectRaw["ports"]
	portsMaps := make([]map[string]interface{}, 0)
	if portsRaw != nil {
		for _, portsChildRaw := range convertToInterfaceArray(portsRaw) {
			portsMap := make(map[string]interface{})
			portsChildRaw := portsChildRaw.(map[string]interface{})
			portsMap["name"] = portsChildRaw["name"]
			portsMap["port"] = portsChildRaw["port"]
			portsMap["protocol"] = portsChildRaw["protocol"]

			portsMaps = append(portsMaps, portsMap)
		}
	}
	if err := d.Set("ports", portsMaps); err != nil {
		return err
	}

	return nil
}

func resourceAliCloudApigServiceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var body map[string]interface{}
	update := false

	var err error
	serviceId := d.Id()
	action := fmt.Sprintf("/v1/services/%s", serviceId)
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})

	if !d.IsNewResource() && d.HasChange("dns_servers") {
		update = true
	}
	if v, ok := d.GetOk("dns_servers"); ok || d.HasChange("dns_servers") {
		dnsServersMapsArray := convertToInterfaceArray(v)

		request["dnsServers"] = dnsServersMapsArray
	}

	if v, ok := d.GetOk("protocol"); ok {
		request["protocol"] = v
	}
	if d.HasChange("health_check_config") {
		update = true
	}
	healthCheckConfig := make(map[string]interface{})

	if v := d.Get("health_check_config"); !IsNil(v) || d.HasChange("health_check_config") {
		protocol1, _ := jsonpath.Get("$[0].protocol", v)
		if protocol1 != nil && protocol1 != "" {
			healthCheckConfig["protocol"] = protocol1
		}
		expectedStatuses1, _ := jsonpath.Get("$[0].expected_statuses", v)
		if expectedStatuses1 != nil && expectedStatuses1 != "" {
			healthCheckConfig["expectedStatuses"] = expectedStatuses1
		}
		unhealthyThreshold1, _ := jsonpath.Get("$[0].unhealthy_threshold", v)
		if unhealthyThreshold1 != nil && unhealthyThreshold1 != "" {
			healthCheckConfig["unhealthyThreshold"] = unhealthyThreshold1
		}
		interval1, _ := jsonpath.Get("$[0].interval", v)
		if interval1 != nil && interval1 != "" {
			healthCheckConfig["interval"] = interval1
		}
		timeout1, _ := jsonpath.Get("$[0].timeout", v)
		if timeout1 != nil && timeout1 != "" {
			healthCheckConfig["timeout"] = timeout1
		}
		enable1, _ := jsonpath.Get("$[0].enable", v)
		if enable1 != nil && enable1 != "" {
			healthCheckConfig["enable"] = enable1
		}
		healthyThreshold1, _ := jsonpath.Get("$[0].healthy_threshold", v)
		if healthyThreshold1 != nil && healthyThreshold1 != "" {
			healthCheckConfig["healthyThreshold"] = healthyThreshold1
		}
		httpPath1, _ := jsonpath.Get("$[0].http_path", v)
		if httpPath1 != nil && httpPath1 != "" {
			healthCheckConfig["httpPath"] = httpPath1
		}
		httpHost1, _ := jsonpath.Get("$[0].http_host", v)
		if httpHost1 != nil && httpHost1 != "" {
			healthCheckConfig["httpHost"] = httpHost1
		}

		request["healthCheckConfig"] = healthCheckConfig
	}

	if d.HasChange("outlier_detection_config") {
		update = true
	}
	outlierDetectionConfig := make(map[string]interface{})

	if v := d.Get("outlier_detection_config"); !IsNil(v) || d.HasChange("outlier_detection_config") {
		interval3, _ := jsonpath.Get("$[0].interval", v)
		if interval3 != nil && interval3 != "" {
			outlierDetectionConfig["interval"] = interval3
		}
		failurePercentageMinimumHosts1, _ := jsonpath.Get("$[0].failure_percentage_minimum_hosts", v)
		if failurePercentageMinimumHosts1 != nil && failurePercentageMinimumHosts1 != "" {
			outlierDetectionConfig["failurePercentageMinimumHosts"] = failurePercentageMinimumHosts1
		}
		enable3, _ := jsonpath.Get("$[0].enable", v)
		if enable3 != nil && enable3 != "" {
			outlierDetectionConfig["enable"] = enable3
		}
		baseEjectionTime1, _ := jsonpath.Get("$[0].base_ejection_time", v)
		if baseEjectionTime1 != nil && baseEjectionTime1 != "" {
			outlierDetectionConfig["baseEjectionTime"] = baseEjectionTime1
		}
		failurePercentageThreshold1, _ := jsonpath.Get("$[0].failure_percentage_threshold", v)
		if failurePercentageThreshold1 != nil && failurePercentageThreshold1 != "" {
			outlierDetectionConfig["failurePercentageThreshold"] = failurePercentageThreshold1
		}

		request["outlierDetectionConfig"] = outlierDetectionConfig
	}

	if !d.IsNewResource() && d.HasChange("addresses") {
		update = true
	}
	if v, ok := d.GetOk("addresses"); ok || d.HasChange("addresses") {
		addressesMapsArray := convertToInterfaceArray(v)

		request["addresses"] = addressesMapsArray
	}

	if d.HasChange("healthy_panic_threshold") {
		update = true
	}
	if v, ok := d.GetOk("healthy_panic_threshold"); ok || d.HasChange("healthy_panic_threshold") {
		request["healthyPanicThreshold"] = v
	}
	body = request
	if update {
		wait := incrementalWait(5*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RoaPut("APIG", "2024-03-27", action, query, nil, body, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"Conflict.LockFailed"}) || NeedRetry(err) {
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
	}
	update = false
	action = fmt.Sprintf("/move-resource-group")
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})
	query["ResourceId"] = StringPointer(d.Id())
	query["RegionId"] = StringPointer(client.RegionId)
	if !d.IsNewResource() && d.HasChange("resource_group_id") {
		update = true
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		query["ResourceGroupId"] = StringPointer(v.(string))
	}

	query["Service"] = StringPointer("APIG")
	query["ResourceType"] = StringPointer("Service")
	body = request
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RoaPost("APIG", "2024-03-27", action, query, nil, body, true)
			if err != nil {
				if NeedRetry(err) {
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
		apigServiceV2 := ApigServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{fmt.Sprint(d.Get("resource_group_id"))}, d.Timeout(schema.TimeoutUpdate), 35*time.Second, apigServiceV2.ApigServiceStateRefreshFunc(d.Id(), "resourceGroupId", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	return resourceAliCloudApigServiceRead(d, meta)
}

func resourceAliCloudApigServiceDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	serviceId := d.Id()
	action := fmt.Sprintf("/v1/services/%s", serviceId)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	var err error
	request = make(map[string]interface{})

	wait := incrementalWait(5*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RoaDelete("APIG", "2024-03-27", action, query, nil, nil, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"Conflict.LockFailed"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		if IsExpectedErrors(err, []string{"NotFound.ServiceNotFound", "NotFound", "404"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
