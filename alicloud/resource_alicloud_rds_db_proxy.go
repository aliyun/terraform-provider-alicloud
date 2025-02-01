package alicloud

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudRdsDBProxy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudRdsDBProxyCreate,
		Read:   resourceAlicloudRdsDBProxyRead,
		Update: resourceAlicloudRdsDBProxyUpdate,
		Delete: resourceAlicloudRdsDBProxyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"db_proxy_instance_num": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(1, 60),
			},
			"instance_network_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"VPC"}, false),
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"db_proxy_connection_prefix": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"read_only_instance_distribution_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"Standard", "Custom"}, false),
			},
			"db_proxy_instance_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"common", "exclusive"}, false),
			},
			"read_only_instance_weight": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"weight": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"read_only_instance_max_delay_time": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"db_proxy_connection_string": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"db_proxy_connect_string_port": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"db_proxy_endpoint_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"net_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"db_proxy_endpoint_aliases": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"db_proxy_endpoint_read_write_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"db_proxy_features": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"effective_time": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"Immediate", "MaintainTime", "SpecificTime"}, false),
			},
			"effective_specific_time": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"db_proxy_ssl_enabled": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"Open", "Close", "Update"}, false),
			},
			"upgrade_time": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"Immediate", "MaintainTime", "SpecifyTime"}, false),
			},
			"switch_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ssl_expired_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudRdsDBProxyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rdsService := RdsService{client}
	action := "ModifyDBProxy"
	request := map[string]interface{}{
		"RegionId":             client.RegionId,
		"DBInstanceId":         Trim(d.Get("instance_id").(string)),
		"ConfigDBProxyService": "Startup",
		"DBProxyInstanceNum":   d.Get("db_proxy_instance_num"),
		"SourceIp":             client.SourceIp,
	}
	v, ok := d.GetOk("instance_network_type")
	if ok && v.(string) != "" {
		request["InstanceNetworkType"] = v
	}
	vpcId, ok := d.GetOk("vpc_id")
	if ok && vpcId.(string) != "" {
		request["VPCId"] = vpcId
	}
	vSwithId, ok := d.GetOk("vswitch_id")
	if ok && vpcId.(string) != "" {
		request["VSwitchId"] = vSwithId
	}
	dBProxyInstanceType, ok := d.GetOk("db_proxy_instance_type")
	if ok && vpcId.(string) != "" {
		request["DBProxyInstanceType"] = dBProxyInstanceType
	}
	resourceGroupId, ok := d.GetOk("resource_group_id")
	if ok && vpcId.(string) != "" {
		request["ResourceGroupId"] = resourceGroupId
	}
	if err := resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err := client.RpcPost("Rds", "2014-08-15", action, nil, request, false)
		if err != nil {
			if NeedRetry(err) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	d.SetId(request["DBInstanceId"].(string))
	stateConf := BuildStateConf([]string{"Creating"}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 3*time.Minute, rdsService.RdsDBProxyStateRefreshFunc(request["DBInstanceId"].(string), []string{"Deleting"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return resourceAlicloudRdsDBProxyUpdate(d, meta)
}

func resourceAlicloudRdsDBProxyRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	rdsService := RdsService{client}
	proxy, proxyErr := rdsService.DescribeDBProxy(d.Id())
	if proxyErr != nil {
		if NotFoundError(proxyErr) {
			d.SetId("")
			return nil
		}
		return WrapError(proxyErr)
	}
	if _, ok := proxy["DBProxyInstanceStatus"]; !ok {
		d.SetId("")
		return nil
	}

	endpointInfo, endpointError := rdsService.DescribeRdsProxyEndpoint(d.Id())
	if endpointError != nil {
		if NotFoundError(endpointError) {
			d.SetId("")
			return nil
		}
		return WrapError(endpointError)
	}
	d.Set("instance_id", d.Id())
	d.Set("db_proxy_instance_type", convertDBProxyInstanceTypeResponse(proxy["DBProxyInstanceType"]))
	d.Set("db_proxy_instance_num", proxy["DBProxyInstanceNum"])
	d.Set("vpc_id", proxy["DBProxyVpcId"])
	d.Set("vswitch_id", proxy["DBProxyVswitchId"])
	d.Set("instance_network_type", "VPC")
	d.Set("db_proxy_endpoint_id", endpointInfo["DBProxyEndpointId"])
	d.Set("db_proxy_connection_string", endpointInfo["DBProxyConnectString"])
	d.Set("db_proxy_connection_prefix", strings.Split(endpointInfo["DBProxyConnectString"].(string), ".")[0])
	d.Set("read_only_instance_distribution_type", endpointInfo["ReadOnlyInstanceDistributionType"])
	d.Set("db_proxy_endpoint_aliases", endpointInfo["DbProxyEndpointAliases"])
	d.Set("db_proxy_endpoint_read_write_mode", endpointInfo["DbProxyEndpointReadWriteMode"])
	if dbProxyConnectStringPort, err := strconv.Atoi(endpointInfo["DBProxyConnectStringPort"].(string)); err == nil {
		d.Set("db_proxy_connect_string_port", dbProxyConnectStringPort)

	}
	if d.Get("db_proxy_endpoint_read_write_mode") == "ReadWrite" {
		if mdt, err := strconv.Atoi(endpointInfo["ReadOnlyInstanceMaxDelayTime"].(string)); err == nil {
			d.Set("read_only_instance_max_delay_time", mdt)
		}
	} else {
		d.Set("read_only_instance_max_delay_time", nil)
	}
	var weight []map[string]interface{}
	rawData := []byte(endpointInfo["ReadOnlyInstanceWeight"].(string))
	parseErr := json.Unmarshal(rawData, &weight)
	if parseErr != nil {
		return WrapError(parseErr)
	}
	readOnlyInstanceWeight := make([]map[string]interface{}, 0)
	for _, val := range weight {
		v := strconv.FormatFloat(val["Weight"].(float64), 'G', -1, 64)
		readOnlyInstanceWeight = append(readOnlyInstanceWeight, map[string]interface{}{
			"instance_id": val["DBInstanceId"],
			"weight":      v,
		})
	}
	if err := d.Set("read_only_instance_weight", readOnlyInstanceWeight); err != nil {
		return WrapError(err)
	}
	proxySsl, proxySslError := rdsService.GetDbProxyInstanceSsl(d.Id())
	if proxySslError != nil {
		if NotFoundError(endpointError) {
			d.SetId("")
			return nil
		}
		return WrapError(endpointError)
	}
	if proxySsl != nil {
		d.Set("ssl_expired_time", proxySsl["SslExpiredTime"])
	}
	return nil
}

func resourceAlicloudRdsDBProxyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rdsService := RdsService{client}
	proxy, proxyErr := rdsService.DescribeDBProxy(d.Id())
	if proxyErr != nil {
		if NotFoundError(proxyErr) {
			d.SetId("")
			return nil
		}
		return WrapError(proxyErr)
	}
	if _, ok := proxy["DBProxyInstanceStatus"]; !ok {
		return nil
	}
	endpointInfo, endpointError := rdsService.DescribeRdsProxyEndpoint(d.Id())
	if endpointError != nil {
		return WrapError(endpointError)
	}
	if d.HasChanges("db_proxy_connection_prefix", "db_proxy_connect_string_port") && d.Get("instance_network_type") != "" {
		action := "ModifyDBProxyEndpointAddress"
		request := map[string]interface{}{
			"RegionId":                    client.RegionId,
			"DBInstanceId":                d.Id(),
			"DBProxyEndpointId":           endpointInfo["DBProxyEndpointId"],
			"DBProxyConnectStringNetType": d.Get("instance_network_type"),
			"SourceIp":                    client.SourceIp,
		}
		portAddressUpdate := false
		if v, ok := d.GetOk("db_proxy_connection_prefix"); ok {
			request["DBProxyNewConnectString"] = v
			portAddressUpdate = true
		}
		if v, ok := d.GetOk("db_proxy_connect_string_port"); ok {
			request["DBProxyNewConnectStringPort"] = v
			portAddressUpdate = true
		}
		if portAddressUpdate {
			if err := resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
				response, err := client.RpcPost("Rds", "2014-08-15", action, nil, request, false)
				if err != nil {
					if NeedRetry(err) {
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				addDebug(action, response, request)
				return nil
			}); err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}
			// waiting state changes to running
			stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 3*time.Minute, rdsService.RdsDBProxyStateRefreshFunc(d.Id(), []string{"Deleting"}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}
		}
	}
	endpointInfo, endpointError = rdsService.DescribeRdsProxyEndpoint(d.Id())
	if endpointError != nil {
		return WrapError(endpointError)
	}
	if !d.IsNewResource() && (d.HasChange("db_proxy_instance_num") || d.HasChange("db_proxy_instance_type")) {
		action := "ModifyDBProxyInstance"
		request := map[string]interface{}{
			"RegionId":            client.RegionId,
			"DBInstanceId":        d.Id(),
			"DBProxyEndpointId":   endpointInfo["DBProxyEndpointId"],
			"DBProxyInstanceNum":  d.Get("db_proxy_instance_num"),
			"DBProxyInstanceType": d.Get("db_proxy_instance_type"),
			"SourceIp":            client.SourceIp,
		}
		if d.HasChange("effective_time") {
			request["EffectiveTime"] = d.Get("effective_time")
		}
		if d.HasChange("effective_specific_time") {
			request["EffectiveSpecificTime"] = d.Get("effective_specific_time")
		}
		if err := resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err := client.RpcPost("Rds", "2014-08-15", action, nil, request, false)
			if err != nil {
				if NeedRetry(err) {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		}); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		// waiting state changes to running
		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 3*time.Minute, rdsService.RdsDBProxyStateRefreshFunc(d.Id(), []string{"Deleting"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("effective_time")
		d.SetPartial("effective_specific_time")
	}
	if d.HasChanges("db_proxy_endpoint_read_write_mode", "read_only_instance_max_delay_time", "read_only_instance_distribution_type", "db_proxy_features", "read_only_instance_weight") {
		action := "ModifyDBProxyEndpoint"
		request := map[string]interface{}{
			"RegionId":          client.RegionId,
			"DBInstanceId":      d.Id(),
			"DBProxyEndpointId": endpointInfo["DBProxyEndpointId"],
			"DbEndpointAliases": endpointInfo["DbProxyEndpointAliases"],
			"SourceIp":          client.SourceIp,
		}
		if v, ok := d.GetOk("db_proxy_endpoint_read_write_mode"); ok {
			request["DbEndpointReadWriteMode"] = v
		}
		if v, ok := d.GetOk("read_only_instance_max_delay_time"); ok {
			request["ReadOnlyInstanceMaxDelayTime"] = v
		}
		if v, ok := d.GetOk("read_only_instance_distribution_type"); ok {
			request["ReadOnlyInstanceDistributionType"] = v
		}
		if v, ok := d.GetOk("db_proxy_features"); ok {
			request["ConfigDBProxyFeatures"] = v
		}

		if v, ok := d.GetOk("read_only_instance_weight"); ok {
			list := v.(*schema.Set).List()
			weightMap := map[string]interface{}{}
			for _, v := range list {
				v := v.(map[string]interface{})
				weightMap[v["instance_id"].(string)] = v["weight"]
			}
			weightStr, _ := convertMaptoJsonString(weightMap)
			request["ReadOnlyInstanceWeight"] = weightStr
			//request["ReadOnlyInstanceWeight"] = v
		}
		if err := resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err := client.RpcPost("Rds", "2014-08-15", action, nil, request, false)
			if err != nil {
				if NeedRetry(err) {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		}); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		// waiting state changes to running
		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 3*time.Minute, rdsService.RdsDBProxyStateRefreshFunc(d.Id(), []string{"Deleting"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	if d.HasChange("db_proxy_ssl_enabled") {
		action := "ModifyDbProxyInstanceSsl"
		request := map[string]interface{}{
			"RegionId":          client.RegionId,
			"DBInstanceId":      d.Id(),
			"DBProxyEndpointId": endpointInfo["DBProxyEndpointId"],
			"SourceIp":          client.SourceIp,
		}
		proxySsl := d.Get("db_proxy_ssl_enabled").(string)
		if proxySsl == "Close" {
			request["DbProxySslEnabled"] = 0
		}
		if proxySsl == "Open" {
			request["DbProxySslEnabled"] = 1
		}
		if proxySsl == "Update" {
			request["DbProxySslEnabled"] = 2
		}
		ProxyEndpoint, ProxyEndpointErr := rdsService.DescribeRdsProxyEndpoint(d.Id())
		if ProxyEndpointErr != nil {
			if NotFoundError(ProxyEndpointErr) {
				d.SetId("")
				return nil
			}
			return WrapError(proxyErr)
		}
		request["DbProxyConnectString"] = ProxyEndpoint["DBProxyConnectString"]
		if err := resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err := client.RpcPost("Rds", "2014-08-15", action, nil, request, false)
			if err != nil {
				if NeedRetry(err) {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		}); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		// waiting state changes to running
		if err := rdsService.WaitForDBInstance(request["DBInstanceId"].(string), Running, 60*60); err != nil {
			return WrapError(err)
		}
		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 3*time.Minute, rdsService.RdsDBProxyStateRefreshFunc(d.Id(), []string{"Deleting"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	if d.HasChange("upgrade_time") {
		action := "UpgradeDBProxyInstanceKernelVersion"
		request := map[string]interface{}{
			"RegionId":     client.RegionId,
			"DBInstanceId": d.Id(),
			"UpgradeTime":  d.Get("upgrade_time"),
			"SourceIp":     client.SourceIp,
		}
		if d.Get("upgrade_time") == "SpecificTime" && d.HasChange("switch_time") {
			request["SwitchTime"] = d.Get("switch_time")
		}
		if err := resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err := client.RpcPost("Rds", "2014-08-15", action, nil, request, false)
			if err != nil {
				if NeedRetry(err) {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		}); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		if err := rdsService.WaitForDBInstance(request["DBInstanceId"].(string), Running, 60*60); err != nil {
			return WrapError(err)
		}
		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 3*time.Minute, rdsService.RdsDBProxyStateRefreshFunc(d.Id(), []string{"Deleting"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}

	}
	return resourceAlicloudRdsDBProxyRead(d, meta)
}

func resourceAlicloudRdsDBProxyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rdsService := RdsService{client}
	if d.Id() == "" {
		return nil
	}
	_, proxyErr := rdsService.DescribeDBProxy(d.Id())
	if proxyErr != nil {
		if NotFoundError(proxyErr) {
			return nil
		}
		return WrapError(proxyErr)
	}
	action := "ModifyDBProxy"
	request := map[string]interface{}{
		"RegionId":             client.RegionId,
		"DBInstanceId":         d.Id(),
		"ConfigDBProxyService": "Shutdown",
		"SourceIp":             client.SourceIp,
	}
	if err := resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err := client.RpcPost("Rds", "2014-08-15", action, nil, request, false)
		if err != nil {
			if IsExpectedErrors(err, OperationDeniedDBStatus) || NeedRetry(err) {
				return resource.RetryableError(err)
			}
			if NotFoundError(err) {
				return nil
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 3*time.Minute, rdsService.RdsDBProxyStateRefreshFunc(d.Id(), []string{"Deleted"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}

func convertDBProxyInstanceTypeResponse(source interface{}) interface{} {
	switch source {
	case "2":
		return "exclusive"
	case "3":
		return "common"
	}
	return ""
}
