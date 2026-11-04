package rds

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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
			// Closing the proxy instances is asynchronous and has been observed to
			// take longer than 20 minutes for some proxies.
			Delete: schema.DefaultTimeout(45 * time.Minute),
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
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Required: true,
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
			"persistent_connection_status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"Enabled", "Disabled"}, false),
			},
			"db_proxy_nodes": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 2,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cpu_cores": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 16),
						},
						"node_counts": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 2),
						},
						"zone_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
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
	if v, ok := d.GetOk("db_proxy_nodes"); ok {
		dbProxyNodes, err := expandDBProxyNodes(v.(*schema.Set))
		if err != nil {
			return WrapError(err)
		}
		request["DBProxyNodes"] = dbProxyNodes
	}
	if err := retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *retry.RetryError {
		response, err := client.RpcPost("Rds", "2014-08-15", action, nil, request, false)
		if err != nil {
			return dbProxyModifyRetryError(err, false)
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
		if NotFoundError(proxySslError) {
			d.SetId("")
			return nil
		}
		return WrapError(proxySslError)
	}
	if proxySsl != nil {
		d.Set("ssl_expired_time", proxySsl["SslExpiredTime"])
		if v, ok := proxySsl["DbProxySslEnabled"]; ok {
			d.Set("db_proxy_ssl_enabled", v)
		}
	}
	if v, ok := proxy["DBProxyPersistentConnectionStatus"]; ok {
		d.Set("persistent_connection_status", v)
	}
	// The node layout is refreshed only when it is managed by the configuration,
	// otherwise the API layout would be added to resources that do not manage it.
	if _, ok := d.GetOk("db_proxy_nodes"); ok {
		nodes, err := describeDBProxyNodes(client, d.Id())
		if err != nil {
			return WrapError(err)
		}
		if len(nodes) > 0 {
			if flattenNodes := flattenDBProxyNodes(nodes); flattenNodes != nil {
				if err := d.Set("db_proxy_nodes", flattenNodes); err != nil {
					return WrapError(err)
				}
			}
		}
	}
	return nil
}

func resourceAlicloudRdsDBProxyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rdsService := RdsService{client}
	// The retries and waits of one update share a single budget: the timeout is
	// not restarted for every sequential step, so a deferred or slow proxy cannot
	// hold the update for a multiple of the configured timeout.
	updateBudget := newDBProxyUpdateBudget(time.Now, d.Timeout(schema.TimeoutUpdate))
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
			if err := retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *retry.RetryError {
				response, err := client.RpcPost("Rds", "2014-08-15", action, nil, request, false)
				if err != nil {
					if NeedRetry(err) {
						return retry.RetryableError(err)
					}
					return retry.NonRetryableError(err)
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
	_, dbProxyNodesConfigured := d.GetOk("db_proxy_nodes")
	// Changing the vSwitch or the VPC migrates the proxy to the target network
	// through MigrateAZ instead of replacing the resource.
	networkChanged := !d.IsNewResource() && d.HasChanges("vswitch_id", "vpc_id")
	if !d.IsNewResource() && (d.HasChange("db_proxy_instance_num") || d.HasChange("db_proxy_instance_type") || (d.HasChange("db_proxy_nodes") && dbProxyNodesConfigured) || networkChanged) {
		action := "ModifyDBProxyInstance"
		request := map[string]interface{}{
			"RegionId":            client.RegionId,
			"DBInstanceId":        d.Id(),
			"DBProxyEndpointId":   endpointInfo["DBProxyEndpointId"],
			"DBProxyInstanceNum":  d.Get("db_proxy_instance_num"),
			"DBProxyInstanceType": d.Get("db_proxy_instance_type"),
			"SourceIp":            client.SourceIp,
		}
		if networkChanged {
			migrateAz, err := buildDBProxyMigrateAz(endpointInfo["DBProxyEndpointId"], d.Get("vswitch_id"), d.Get("vpc_id"))
			if err != nil {
				return WrapError(err)
			}
			request["MigrateAZ"] = migrateAz
		}
		// EffectiveTime defaults to MaintainTime on the server, which defers the change to the
		// maintenance window and leaves the proxy in a changing state, so the configured value is
		// sent as-is and Immediate is sent when the attribute is absent from the configuration.
		effectiveTime, effectiveSpecificTime, err := dbProxyEffectiveTime(d)
		if err != nil {
			return WrapError(err)
		}
		request["EffectiveTime"] = effectiveTime
		if effectiveTime == "SpecificTime" {
			request["EffectiveSpecificTime"] = effectiveSpecificTime
		}
		targetZone := ""
		var currentLayout []map[string]interface{}
		if networkChanged {
			var err error
			targetZone, err = dbProxyTargetZone(client, d.Get("vswitch_id").(string))
			if err != nil {
				return WrapError(err)
			}
		}
		// The current layout is needed on every network change and on every
		// specification update that does not send the layout itself: an in-place
		// migration moves every node into the target zone, and the API requires
		// DBProxyNodes when a multi-AZ proxy is modified.
		if networkChanged || !dbProxyNodesConfigured {
			var err error
			currentLayout, err = currentDBProxyNodeLayout(client, d.Id())
			if err != nil {
				return WrapError(err)
			}
		}
		if err := dbProxyMultiAZGuard(networkChanged, dbProxyNodesConfigured, currentLayout); err != nil {
			return WrapError(err)
		}
		if v, ok := d.GetOk("db_proxy_nodes"); ok {
			if networkChanged {
				for _, item := range v.(*schema.Set).List() {
					node, _ := item.(map[string]interface{})
					zoneId, _ := node["zone_id"].(string)
					if zoneId != "" && zoneId != targetZone {
						return WrapError(fmt.Errorf("db_proxy_nodes.zone_id (%s) must be the availability zone of the target vswitch_id (%s), the proxy nodes move together with the proxy", zoneId, targetZone))
					}
				}
			}
			dbProxyNodes, err := expandDBProxyNodes(v.(*schema.Set))
			if err != nil {
				return WrapError(err)
			}
			request["DBProxyNodes"] = dbProxyNodes
		} else if networkChanged && len(currentLayout) == 1 {
			// MigrateAZ alone is accepted but does not move the proxy: the API needs
			// the target node layout. The current node size is resent in the target
			// availability zone.
			dbProxyNodes, err := dbProxyNodesOfZone(currentLayout, targetZone)
			if err != nil {
				return WrapError(err)
			}
			if dbProxyNodes != "" {
				request["DBProxyNodes"] = dbProxyNodes
			}
		}
		steps := []func(timeout time.Duration) error{
			func(timeout time.Duration) error {
				if err := retry.Retry(client.GetRetryTimeout(timeout), func() *retry.RetryError {
					response, err := client.RpcPost("Rds", "2014-08-15", action, nil, request, false)
					if err != nil {
						if IsExpectedErrors(err, OperationDeniedDBStatus) || IsExpectedErrors(err, []string{"TaskExists", "MaxscaleInternalError", "APICallingFailed"}) || NeedRetry(err) {
							return retry.RetryableError(err)
						}
						return retry.NonRetryableError(err)
					}
					addDebug(action, response, request)
					return nil
				}); err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
				}
				return nil
			},
			// waiting state changes to running
			func(timeout time.Duration) error {
				stateConf := BuildStateConf([]string{}, []string{"Running"}, timeout, 3*time.Minute, rdsService.RdsDBProxyStateRefreshFunc(d.Id(), []string{"Deleting"}))
				if _, err := stateConf.WaitForStateContext(context.Background()); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}
				return nil
			},
		}
		// A finished migration task is not reflected in DescribeDBProxy at once, so the
		// target network is polled before the state is refreshed from the API. A deferred
		// effective time is not polled because the migration only starts at that time.
		if networkChanged && effectiveTime == "Immediate" {
			vswitchId := d.Get("vswitch_id").(string)
			vpcId := d.Get("vpc_id").(string)
			steps = append(steps, func(timeout time.Duration) error {
				return WrapError(waitForDBProxyNetwork(client, d.Id(), vswitchId, vpcId, timeout))
			})
		}
		if err := runDBProxyUpdateSteps(updateBudget, steps...); err != nil {
			return err
		}
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
		if err := retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *retry.RetryError {
			response, err := client.RpcPost("Rds", "2014-08-15", action, nil, request, false)
			if err != nil {
				if NeedRetry(err) {
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
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

	if d.HasChange("persistent_connection_status") {
		if v, ok := d.GetOk("persistent_connection_status"); ok {
			action := "ModifyDBProxy"
			request := map[string]interface{}{
				"RegionId":                   client.RegionId,
				"DBInstanceId":               d.Id(),
				"ConfigDBProxyService":       "Modify",
				"PersistentConnectionStatus": v,
				"SourceIp":                   client.SourceIp,
			}
			steps := []func(timeout time.Duration) error{
				func(timeout time.Duration) error {
					if err := retry.Retry(client.GetRetryTimeout(timeout), func() *retry.RetryError {
						response, err := client.RpcPost("Rds", "2014-08-15", action, nil, request, false)
						if err != nil {
							return dbProxyModifyConfigRetryError(err)
						}
						addDebug(action, response, request)
						return nil
					}); err != nil {
						return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
					}
					return nil
				},
				func(timeout time.Duration) error {
					stateConf := BuildStateConf([]string{}, []string{"Running"}, timeout, 3*time.Minute, rdsService.RdsDBProxyStateRefreshFunc(d.Id(), []string{"Deleting"}))
					if _, err := stateConf.WaitForStateContext(context.Background()); err != nil {
						return WrapErrorf(err, IdMsg, d.Id())
					}
					return nil
				},
			}
			if err := runDBProxyUpdateSteps(updateBudget, steps...); err != nil {
				return err
			}
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
		if err := retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *retry.RetryError {
			response, err := client.RpcPost("Rds", "2014-08-15", action, nil, request, false)
			if err != nil {
				if NeedRetry(err) {
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
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
		if err := retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *retry.RetryError {
			response, err := client.RpcPost("Rds", "2014-08-15", action, nil, request, false)
			if err != nil {
				if NeedRetry(err) {
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
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
	if err := retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *retry.RetryError {
		response, err := client.RpcPost("Rds", "2014-08-15", action, nil, request, false)
		if err != nil {
			return dbProxyModifyRetryError(err, true)
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

func expandDBProxyNodes(nodes *schema.Set) (string, error) {
	items := nodes.List()
	dbProxyNodes := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		node, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		zoneId, _ := node["zone_id"].(string)
		cpuCores, _ := node["cpu_cores"].(int)
		nodeCounts, _ := node["node_counts"].(int)
		// Every sub-attribute has to be set explicitly: a default would be written
		// back from the API into the set element hash and never converge.
		if zoneId == "" || cpuCores == 0 || nodeCounts == 0 {
			return "", fmt.Errorf("db_proxy_nodes requires zone_id, cpu_cores and node_counts, got zone_id=%q cpu_cores=%d node_counts=%d", zoneId, cpuCores, nodeCounts)
		}
		dbProxyNodes = append(dbProxyNodes, map[string]interface{}{
			"cpuCores":   strconv.Itoa(cpuCores),
			"nodeCounts": strconv.Itoa(nodeCounts),
			"zoneId":     zoneId,
		})
	}
	result, err := json.Marshal(dbProxyNodes)
	if err != nil {
		return "", err
	}
	return string(result), nil
}

// dbProxyEffectiveTime returns the EffectiveTime that the request carries.
// The attribute is Optional+Computed, so the raw configuration decides whether a
// value was requested; Immediate is used when the attribute is absent.
func dbProxyEffectiveTime(d *schema.ResourceData) (string, interface{}, error) {
	requested := ""
	if raw := d.GetRawConfig(); raw != cty.NilVal && !raw.IsNull() {
		if v := raw.GetAttr("effective_time"); v.IsKnown() && !v.IsNull() {
			requested = v.AsString()
		}
	} else if v, ok := d.GetOk("effective_time"); ok {
		requested = v.(string)
	}
	switch requested {
	case "":
		return "Immediate", nil, nil
	case "SpecificTime":
		// The API requires the point in time together with SpecificTime, so a
		// missing value is reported here instead of as a server-side parameter error.
		specific, _ := d.Get("effective_specific_time").(string)
		if specific == "" {
			return "", nil, fmt.Errorf("effective_specific_time must be set when effective_time is SpecificTime")
		}
		return requested, specific, nil
	default:
		return requested, nil, nil
	}
}

// dbProxyTargetZone returns the availability zone of the vSwitch that the
// proxy is migrated to.
func dbProxyTargetZone(client *connectivity.AliyunClient, vswitchId string) (string, error) {
	vpcService := VpcService{client}
	vswitch, err := vpcService.DescribeVswitch(vswitchId)
	if err != nil {
		return "", err
	}
	zone, _ := vswitch["ZoneId"].(string)
	if zone == "" {
		return "", fmt.Errorf("cannot determine the availability zone of vswitch %q", vswitchId)
	}
	return zone, nil
}

// currentDBProxyNodeLayout returns the node layout that the proxy currently
// reports, aggregated per availability zone.
func currentDBProxyNodeLayout(client *connectivity.AliyunClient, id string) ([]map[string]interface{}, error) {
	nodes, err := describeDBProxyNodes(client, id)
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, nil
	}
	return flattenDBProxyNodes(nodes), nil
}

// dbProxyNodesOfZone rebuilds the DBProxyNodes request value for the target
// availability zone from the current node layout.
func dbProxyNodesOfZone(layout []map[string]interface{}, zoneId string) (string, error) {
	items := make([]map[string]interface{}, 0, len(layout))
	for _, node := range layout {
		cpuCores, ok := node["cpu_cores"].(int)
		nodeCounts, _ := node["node_counts"].(int)
		if !ok || cpuCores <= 0 || nodeCounts <= 0 || zoneId == "" {
			return "", nil
		}
		items = append(items, map[string]interface{}{
			"cpuCores":   strconv.Itoa(cpuCores),
			"nodeCounts": strconv.Itoa(nodeCounts),
			"zoneId":     zoneId,
		})
	}
	result, err := json.Marshal(items)
	if err != nil {
		return "", err
	}
	return string(result), nil
}

// dbProxyMultiAZGuard reports the error that stops a proxy change before the
// request is sent when the proxy currently spans several availability zones. An
// in-place network change cannot move a multi-AZ proxy, and the API requires
// DBProxyNodes when any other change is applied to such a proxy.
func dbProxyMultiAZGuard(networkChanged, dbProxyNodesConfigured bool, layout []map[string]interface{}) error {
	if len(layout) <= 1 {
		return nil
	}
	zones := make([]string, 0, len(layout))
	for _, node := range layout {
		zones = append(zones, fmt.Sprint(node["zone_id"]))
	}
	zoneList := strings.Join(zones, ", ")
	if networkChanged {
		return fmt.Errorf("the proxy is deployed across availability zones (%s); changing vswitch_id or vpc_id in place is not supported for a multi-AZ proxy, recreate the proxy instead", zoneList)
	}
	if !dbProxyNodesConfigured {
		return fmt.Errorf("the proxy is deployed across availability zones (%s); db_proxy_nodes must be configured when a multi-AZ proxy is modified", zoneList)
	}
	return nil
}

// dbProxyModifyServiceTransientErrors are the ModifyDBProxy state conflicts
// that both Startup and Shutdown can observe while another proxy task is
// running. MaxscaleInternalError is HTTP 400 and MaxscaleCreating is HTTP 404,
// so they are matched before Shutdown can read the not-found status as a proxy
// that is already deleted.
var dbProxyModifyServiceTransientErrors = []string{"MaxscaleInternalError", "MaxscaleCreating"}

// dbProxyModifyRetryError classifies a ModifyDBProxy Startup or Shutdown error.
// A nil return means the error is treated as success, which only the Shutdown
// path does for a missing proxy.
func dbProxyModifyRetryError(err error, shutdown bool) *retry.RetryError {
	if IsExpectedErrors(err, dbProxyModifyServiceTransientErrors) ||
		(shutdown && IsExpectedErrors(err, OperationDeniedDBStatus)) ||
		NeedRetry(err) {
		return retry.RetryableError(err)
	}
	if shutdown && NotFoundError(err) {
		return nil
	}
	return retry.NonRetryableError(err)
}

// dbProxyModifyConfigTransientErrors are the ModifyDBProxy state conflicts that
// the ConfigDBProxyService=Modify write can observe while another proxy task is
// running. MaxscaleCreating is HTTP 404 but only reports that the associated
// proxy is being created or deleted, so it has to be retried instead of failing
// the persistent connection update.
var dbProxyModifyConfigTransientErrors = []string{"TaskExists", "MaxscaleInternalError", "MaxscaleCreating"}

// dbProxyModifyConfigRetryError classifies a ModifyDBProxy
// ConfigDBProxyService=Modify error. A missing proxy stays terminal here: only
// the Shutdown path treats it as an already deleted proxy.
func dbProxyModifyConfigRetryError(err error) *retry.RetryError {
	if IsExpectedErrors(err, OperationDeniedDBStatus) ||
		IsExpectedErrors(err, dbProxyModifyConfigTransientErrors) ||
		NeedRetry(err) {
		return retry.RetryableError(err)
	}
	return retry.NonRetryableError(err)
}

// dbProxyUpdateBudget holds the single timeout that the sequential retries and
// waits of one update share. The SDK does not deduct elapsed time from
// ResourceData.Timeout, so passing the full timeout to every step would let a
// slow or deferred proxy run for a multiple of it.
type dbProxyUpdateBudget struct {
	deadline time.Time
	total    time.Duration
	now      func() time.Time
}

func newDBProxyUpdateBudget(now func() time.Time, total time.Duration) *dbProxyUpdateBudget {
	return &dbProxyUpdateBudget{deadline: now().Add(total), total: total, now: now}
}

// Remaining returns the time left in the shared budget.
func (b *dbProxyUpdateBudget) Remaining() (time.Duration, error) {
	remaining := b.deadline.Sub(b.now())
	if remaining <= 0 {
		return 0, fmt.Errorf("the %s update timeout was exhausted before the next proxy change could be applied", b.total)
	}
	return remaining, nil
}

// runDBProxyUpdateSteps runs the sequential retry and wait steps of one update.
// Every step receives the time left in the shared budget, and no step is started
// once that budget is exhausted.
func runDBProxyUpdateSteps(budget *dbProxyUpdateBudget, steps ...func(timeout time.Duration) error) error {
	for _, step := range steps {
		timeout, err := budget.Remaining()
		if err != nil {
			return WrapError(err)
		}
		if err := step(timeout); err != nil {
			return err
		}
	}
	return nil
}

// waitForDBProxyNetwork waits until the proxy reports the target network,
// because MigrateAZ returns before the migration is reflected by DescribeDBProxy.
func waitForDBProxyNetwork(client *connectivity.AliyunClient, id, vswitchId, vpcId string, timeout time.Duration) error {
	rdsService := RdsService{client}
	stateConf := &retry.StateChangeConf{
		Pending: []string{"Pending"},
		Target:  []string{"Converged"},
		Refresh: func() (interface{}, string, error) {
			object, err := rdsService.DescribeDBProxy(id)
			if err != nil {
				if NotFoundError(err) {
					return nil, "", nil
				}
				return nil, "", WrapError(err)
			}
			if fmt.Sprint(object["DBProxyVswitchId"]) == vswitchId && fmt.Sprint(object["DBProxyVpcId"]) == vpcId {
				return object, "Converged", nil
			}
			return object, "Pending", nil
		},
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(context.Background())
	return err
}

func buildDBProxyMigrateAz(endpointId, vswitchId, vpcId interface{}) (string, error) {
	endpoint, _ := endpointId.(string)
	vswitch, _ := vswitchId.(string)
	vpc, _ := vpcId.(string)
	if endpoint == "" || vswitch == "" || vpc == "" {
		return "", fmt.Errorf("migrating the proxy requires a proxy endpoint, vswitch_id and vpc_id, got endpoint=%q vswitch_id=%q vpc_id=%q", endpoint, vswitch, vpc)
	}
	migrateAz, err := json.Marshal([]map[string]interface{}{
		{
			"dbProxyEndpointId": endpoint,
			"destVSwitchId":     vswitch,
			"destVpcId":         vpc,
		},
	})
	if err != nil {
		return "", err
	}
	return string(migrateAz), nil
}

func describeDBProxyNodes(client *connectivity.AliyunClient, id string) ([]interface{}, error) {
	action := "DescribeDBProxy"
	request := map[string]interface{}{
		"RegionId":     client.RegionId,
		"DBInstanceId": id,
		"SourceIp":     client.SourceIp,
	}
	response, err := client.RpcPost("Rds", "2014-08-15", action, nil, request, true)
	if err != nil {
		// The node layout is auxiliary data of the proxy that is already checked by the caller,
		// so a vanished proxy only skips the refresh instead of failing the whole read.
		if IsExpectedErrors(err, []string{"InvalidDBInstanceId.NotFound"}) {
			return nil, nil
		}
		return nil, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	addDebug(action, response, request)
	proxyNodes, ok := response["DBProxyNodes"].(map[string]interface{})
	if !ok {
		return nil, nil
	}
	nodes, _ := proxyNodes["DBProxyNodes"].([]interface{})
	return nodes, nil
}

// DescribeDBProxy returns one entry per proxy node, while the configuration is
// expressed per availability zone, so the nodes are grouped by zone.
func flattenDBProxyNodes(apiNodes []interface{}) []map[string]interface{} {
	zoneIds := make([]string, 0)
	aggregated := make(map[string]map[string]interface{}, 0)
	for _, item := range apiNodes {
		apiNode, ok := item.(map[string]interface{})
		if !ok {
			return nil
		}
		zoneId, _ := apiNode["zoneId"].(string)
		if zoneId == "" {
			return nil
		}
		node, ok := aggregated[zoneId]
		if !ok {
			cpuCores, ok := dbProxyNodeCpuCores(apiNode["cpuCores"])
			if !ok {
				return nil
			}
			node = map[string]interface{}{
				"zone_id":     zoneId,
				"cpu_cores":   cpuCores,
				"node_counts": 0,
			}
			aggregated[zoneId] = node
			zoneIds = append(zoneIds, zoneId)
		}
		node["node_counts"] = node["node_counts"].(int) + 1
	}
	result := make([]map[string]interface{}, 0, len(zoneIds))
	for _, zoneId := range zoneIds {
		result = append(result, aggregated[zoneId])
	}
	return result
}

func dbProxyNodeCpuCores(source interface{}) (int, bool) {
	switch v := source.(type) {
	case string:
		if cores, err := strconv.Atoi(v); err == nil {
			return cores, true
		}
	case float64:
		return int(v), true
	}
	return 0, false
}
