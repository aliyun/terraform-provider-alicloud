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

func resourceAliCloudAlbServerGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudAlbServerGroupCreate,
		Read:   resourceAliCloudAlbServerGroupRead,
		Update: resourceAliCloudAlbServerGroupUpdate,
		Delete: resourceAliCloudAlbServerGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		CustomizeDiff: resourceAlbServerGroupCustomizeDiff,
		Schema: map[string]*schema.Schema{
			"connection_drain_config": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"connection_drain_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"connection_drain_timeout": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cross_zone_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"dry_run": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"health_check_config": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"health_check_interval": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: IntBetween(0, 50),
						},
						"health_check_connect_port": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: IntBetween(0, 65535),
						},
						"health_check_codes": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"unhealthy_threshold": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: IntBetween(2, 10),
						},
						"health_check_method": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: StringInSlice([]string{"GET", "POST", "HEAD"}, false),
						},
						"health_check_host": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"health_check_path": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"healthy_threshold": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: IntBetween(2, 10),
						},
						"health_check_protocol": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: StringInSlice([]string{"HTTP", "HTTPS", "TCP", "GRPC"}, false),
						},
						"health_check_http_version": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: StringInSlice([]string{"HTTP1.0", "HTTP1.1"}, false),
						},
						"health_check_enabled": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"health_check_timeout": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: IntBetween(0, 300),
						},
					},
				},
			},
			"health_check_template_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ipv6_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"protocol": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"HTTP", "HTTPS", "GRPC"}, false),
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"scheduler": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"Wrr", "Wlc", "Sch", "Uch"}, false),
			},
			"server_group_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"server_group_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"servers": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"server_type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"server_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"port": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: IntBetween(0, 65535),
						},
						"server_ip": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"remote_ip_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"server_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"weight": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"service_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"slow_start_config": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"slow_start_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"slow_start_duration": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sticky_session_config": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cookie": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"cookie_timeout": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: IntBetween(0, 86400),
						},
						"sticky_session_type": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: StringInSlice([]string{"Insert", "Server"}, false),
						},
						"sticky_session_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
			"tags": tagsSchema(),
			"uch_config": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: StringInSlice([]string{"QueryString"}, false),
						},
						"value": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"upstream_keepalive_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudAlbServerGroupCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateServerGroup"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	request["ClientToken"] = buildClientToken(action)

	objectDataLocalMap := make(map[string]interface{})

	if v := d.Get("sticky_session_config"); !IsNil(v) {
		cookieTimeout1, _ := jsonpath.Get("$[0].cookie_timeout", v)
		if cookieTimeout1 != nil && cookieTimeout1 != "" && cookieTimeout1.(int) > 0 {
			objectDataLocalMap["CookieTimeout"] = cookieTimeout1
		}
		cookie1, _ := jsonpath.Get("$[0].cookie", v)
		if cookie1 != nil && cookie1 != "" {
			objectDataLocalMap["Cookie"] = cookie1
		}
		stickySessionEnabled1, _ := jsonpath.Get("$[0].sticky_session_enabled", v)
		if stickySessionEnabled1 != nil && stickySessionEnabled1 != "" {
			objectDataLocalMap["StickySessionEnabled"] = stickySessionEnabled1
		}
		stickySessionType1, _ := jsonpath.Get("$[0].sticky_session_type", v)
		if stickySessionType1 != nil && stickySessionType1 != "" {
			objectDataLocalMap["StickySessionType"] = stickySessionType1
		}

		request["StickySessionConfig"] = objectDataLocalMap
	}

	objectDataLocalMap1 := make(map[string]interface{})

	if v := d.Get("health_check_config"); v != nil {
		healthCheckCodes1, _ := jsonpath.Get("$[0].health_check_codes", v)
		if healthCheckCodes1 != nil && healthCheckCodes1 != "" {
			objectDataLocalMap1["HealthCheckCodes"] = healthCheckCodes1
		}
		healthCheckHost1, _ := jsonpath.Get("$[0].health_check_host", v)
		if healthCheckHost1 != nil && healthCheckHost1 != "" {
			objectDataLocalMap1["HealthCheckHost"] = healthCheckHost1
		}
		healthCheckPath1, _ := jsonpath.Get("$[0].health_check_path", v)
		if healthCheckPath1 != nil && healthCheckPath1 != "" {
			objectDataLocalMap1["HealthCheckPath"] = healthCheckPath1
		}
		healthCheckEnabled1, _ := jsonpath.Get("$[0].health_check_enabled", v)
		if healthCheckEnabled1 != nil && healthCheckEnabled1 != "" {
			objectDataLocalMap1["HealthCheckEnabled"] = healthCheckEnabled1
		}
		healthCheckTimeout1, _ := jsonpath.Get("$[0].health_check_timeout", v)
		if healthCheckTimeout1 != nil && healthCheckTimeout1 != "" && healthCheckTimeout1.(int) > 0 {
			objectDataLocalMap1["HealthCheckTimeout"] = healthCheckTimeout1
		}
		unhealthyThreshold1, _ := jsonpath.Get("$[0].unhealthy_threshold", v)
		if unhealthyThreshold1 != nil && unhealthyThreshold1 != "" && unhealthyThreshold1.(int) > 0 {
			objectDataLocalMap1["UnhealthyThreshold"] = unhealthyThreshold1
		}
		healthCheckInterval1, _ := jsonpath.Get("$[0].health_check_interval", v)
		if healthCheckInterval1 != nil && healthCheckInterval1 != "" && healthCheckInterval1.(int) > 0 {
			objectDataLocalMap1["HealthCheckInterval"] = healthCheckInterval1
		}
		healthCheckConnectPort1, _ := jsonpath.Get("$[0].health_check_connect_port", v)
		if healthCheckConnectPort1 != nil && healthCheckConnectPort1 != "" {
			objectDataLocalMap1["HealthCheckConnectPort"] = healthCheckConnectPort1
		}
		healthCheckHttpVersion1, _ := jsonpath.Get("$[0].health_check_http_version", v)
		if healthCheckHttpVersion1 != nil && healthCheckHttpVersion1 != "" {
			objectDataLocalMap1["HealthCheckHttpVersion"] = healthCheckHttpVersion1
		}
		healthCheckMethod1, _ := jsonpath.Get("$[0].health_check_method", v)
		if healthCheckMethod1 != nil && healthCheckMethod1 != "" {
			objectDataLocalMap1["HealthCheckMethod"] = healthCheckMethod1
		}
		healthCheckProtocol1, _ := jsonpath.Get("$[0].health_check_protocol", v)
		if healthCheckProtocol1 != nil && healthCheckProtocol1 != "" {
			objectDataLocalMap1["HealthCheckProtocol"] = healthCheckProtocol1
		}
		healthyThreshold1, _ := jsonpath.Get("$[0].healthy_threshold", v)
		if healthyThreshold1 != nil && healthyThreshold1 != "" && healthyThreshold1.(int) > 0 {
			objectDataLocalMap1["HealthyThreshold"] = healthyThreshold1
		}

		request["HealthCheckConfig"] = objectDataLocalMap1
	}

	if v, ok := d.GetOk("protocol"); ok {
		request["Protocol"] = v
	}
	objectDataLocalMap2 := make(map[string]interface{})

	if v := d.Get("slow_start_config"); !IsNil(v) {
		slowStartEnabled1, _ := jsonpath.Get("$[0].slow_start_enabled", v)
		if slowStartEnabled1 != nil && slowStartEnabled1 != "" {
			objectDataLocalMap2["SlowStartEnabled"] = slowStartEnabled1
		}
		slowStartDuration1, _ := jsonpath.Get("$[0].slow_start_duration", v)
		if slowStartDuration1 != nil && slowStartDuration1 != "" {
			objectDataLocalMap2["SlowStartDuration"] = slowStartDuration1
		}

		request["SlowStartConfig"] = objectDataLocalMap2
	}

	request["ServerGroupName"] = d.Get("server_group_name")
	if v, ok := d.GetOk("server_group_type"); ok {
		request["ServerGroupType"] = v
	}
	if v, ok := d.GetOk("tags"); ok {
		tagsMap := ConvertTags(v.(map[string]interface{}))
		request["Tags"] = tagsMap
	}

	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	objectDataLocalMap3 := make(map[string]interface{})

	if v := d.Get("connection_drain_config"); !IsNil(v) {
		connectionDrainEnabled1, _ := jsonpath.Get("$[0].connection_drain_enabled", v)
		if connectionDrainEnabled1 != nil && connectionDrainEnabled1 != "" {
			objectDataLocalMap3["ConnectionDrainEnabled"] = connectionDrainEnabled1
		}
		connectionDrainTimeout1, _ := jsonpath.Get("$[0].connection_drain_timeout", v)
		if connectionDrainTimeout1 != nil && connectionDrainTimeout1 != "" {
			objectDataLocalMap3["ConnectionDrainTimeout"] = connectionDrainTimeout1
		}

		request["ConnectionDrainConfig"] = objectDataLocalMap3
	}

	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	objectDataLocalMap4 := make(map[string]interface{})

	if v := d.Get("uch_config"); v != nil {
		value2, _ := jsonpath.Get("$[0].value", v)
		if value2 != nil && value2 != "" {
			objectDataLocalMap4["Value"] = value2
		}
		type1, _ := jsonpath.Get("$[0].type", v)
		if type1 != nil && type1 != "" {
			objectDataLocalMap4["Type"] = type1
		}

		request["UchConfig"] = objectDataLocalMap4
	}

	if v, ok := d.GetOk("scheduler"); ok {
		request["Scheduler"] = v
	}
	if v, ok := d.GetOkExists("upstream_keepalive_enabled"); ok {
		request["UpstreamKeepaliveEnabled"] = v
	}
	if v, ok := d.GetOkExists("ipv6_enabled"); ok {
		request["Ipv6Enabled"] = v
	}
	if v, ok := d.GetOk("vpc_id"); ok {
		request["VpcId"] = v
	}
	if v, ok := d.GetOk("service_name"); ok {
		request["ServiceName"] = v
	}
	if v, ok := d.GetOkExists("cross_zone_enabled"); ok {
		request["CrossZoneEnabled"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Alb", "2020-06-16", action, query, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"SystemBusy", "OperationFailed.ResourceGroupStatusCheckFail", "IdempotenceProcessing"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_alb_server_group", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["ServerGroupId"]))

	albServiceV2 := AlbServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutCreate), 0, albServiceV2.AlbServerGroupStateRefreshFunc(d.Id(), "ServerGroupStatus", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudAlbServerGroupUpdate(d, meta)
}

func resourceAliCloudAlbServerGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	albServiceV2 := AlbServiceV2{client}

	objectRaw, err := albServiceV2.DescribeAlbServerGroup(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_alb_server_group DescribeAlbServerGroup Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("create_time", objectRaw["CreateTime"])
	d.Set("cross_zone_enabled", objectRaw["CrossZoneEnabled"])
	d.Set("ipv6_enabled", objectRaw["Ipv6Enabled"])
	d.Set("protocol", objectRaw["Protocol"])
	d.Set("resource_group_id", objectRaw["ResourceGroupId"])
	d.Set("scheduler", objectRaw["Scheduler"])
	d.Set("server_group_name", objectRaw["ServerGroupName"])
	d.Set("server_group_type", objectRaw["ServerGroupType"])
	d.Set("status", objectRaw["ServerGroupStatus"])
	d.Set("upstream_keepalive_enabled", objectRaw["UpstreamKeepaliveEnabled"])
	d.Set("vpc_id", objectRaw["VpcId"])

	connectionDrainConfigMaps := make([]map[string]interface{}, 0)
	connectionDrainConfigMap := make(map[string]interface{})
	connectionDrainConfigRaw := make(map[string]interface{})
	if objectRaw["ConnectionDrainConfig"] != nil {
		connectionDrainConfigRaw = objectRaw["ConnectionDrainConfig"].(map[string]interface{})
	}
	if len(connectionDrainConfigRaw) > 0 {
		connectionDrainConfigMap["connection_drain_enabled"] = connectionDrainConfigRaw["ConnectionDrainEnabled"]
		connectionDrainConfigMap["connection_drain_timeout"] = connectionDrainConfigRaw["ConnectionDrainTimeout"]

		connectionDrainConfigMaps = append(connectionDrainConfigMaps, connectionDrainConfigMap)
	}
	if err := d.Set("connection_drain_config", connectionDrainConfigMaps); err != nil {
		return err
	}
	healthCheckConfigMaps := make([]map[string]interface{}, 0)
	healthCheckConfigMap := make(map[string]interface{})
	healthCheckConfigRaw := make(map[string]interface{})
	if objectRaw["HealthCheckConfig"] != nil {
		healthCheckConfigRaw = objectRaw["HealthCheckConfig"].(map[string]interface{})
	}
	if len(healthCheckConfigRaw) > 0 {
		healthCheckConfigMap["health_check_connect_port"] = healthCheckConfigRaw["HealthCheckConnectPort"]
		healthCheckConfigMap["health_check_enabled"] = healthCheckConfigRaw["HealthCheckEnabled"]
		healthCheckConfigMap["health_check_host"] = healthCheckConfigRaw["HealthCheckHost"]
		healthCheckConfigMap["health_check_http_version"] = healthCheckConfigRaw["HealthCheckHttpVersion"]
		healthCheckConfigMap["health_check_interval"] = healthCheckConfigRaw["HealthCheckInterval"]
		healthCheckConfigMap["health_check_method"] = healthCheckConfigRaw["HealthCheckMethod"]
		healthCheckConfigMap["health_check_path"] = healthCheckConfigRaw["HealthCheckPath"]
		healthCheckConfigMap["health_check_protocol"] = healthCheckConfigRaw["HealthCheckProtocol"]
		healthCheckConfigMap["health_check_timeout"] = healthCheckConfigRaw["HealthCheckTimeout"]
		healthCheckConfigMap["healthy_threshold"] = healthCheckConfigRaw["HealthyThreshold"]
		healthCheckConfigMap["unhealthy_threshold"] = healthCheckConfigRaw["UnhealthyThreshold"]

		healthCheckCodesRaw := make([]interface{}, 0)
		if healthCheckConfigRaw["HealthCheckCodes"] != nil {
			healthCheckCodesRaw = healthCheckConfigRaw["HealthCheckCodes"].([]interface{})
		}

		healthCheckConfigMap["health_check_codes"] = healthCheckCodesRaw
		healthCheckConfigMaps = append(healthCheckConfigMaps, healthCheckConfigMap)
	}
	if err := d.Set("health_check_config", healthCheckConfigMaps); err != nil {
		return err
	}
	slowStartConfigMaps := make([]map[string]interface{}, 0)
	slowStartConfigMap := make(map[string]interface{})
	slowStartConfigRaw := make(map[string]interface{})
	if objectRaw["SlowStartConfig"] != nil {
		slowStartConfigRaw = objectRaw["SlowStartConfig"].(map[string]interface{})
	}
	if len(slowStartConfigRaw) > 0 {
		slowStartConfigMap["slow_start_duration"] = slowStartConfigRaw["SlowStartDuration"]
		slowStartConfigMap["slow_start_enabled"] = slowStartConfigRaw["SlowStartEnabled"]

		slowStartConfigMaps = append(slowStartConfigMaps, slowStartConfigMap)
	}
	if err := d.Set("slow_start_config", slowStartConfigMaps); err != nil {
		return err
	}
	stickySessionConfigMaps := make([]map[string]interface{}, 0)
	stickySessionConfigMap := make(map[string]interface{})
	stickySessionConfigRaw := make(map[string]interface{})
	if objectRaw["StickySessionConfig"] != nil {
		stickySessionConfigRaw = objectRaw["StickySessionConfig"].(map[string]interface{})
	}
	if len(stickySessionConfigRaw) > 0 {
		stickySessionConfigMap["cookie"] = stickySessionConfigRaw["Cookie"]
		stickySessionConfigMap["cookie_timeout"] = stickySessionConfigRaw["CookieTimeout"]
		stickySessionConfigMap["sticky_session_enabled"] = stickySessionConfigRaw["StickySessionEnabled"]
		stickySessionConfigMap["sticky_session_type"] = stickySessionConfigRaw["StickySessionType"]

		stickySessionConfigMaps = append(stickySessionConfigMaps, stickySessionConfigMap)
	}
	if err := d.Set("sticky_session_config", stickySessionConfigMaps); err != nil {
		return err
	}
	tagsMaps := objectRaw["Tags"]
	d.Set("tags", tagsToMap(tagsMaps))
	uchConfigMaps := make([]map[string]interface{}, 0)
	uchConfigMap := make(map[string]interface{})
	uchConfigRaw := make(map[string]interface{})
	if objectRaw["UchConfig"] != nil {
		uchConfigRaw = objectRaw["UchConfig"].(map[string]interface{})
	}
	if len(uchConfigRaw) > 0 {
		uchConfigMap["type"] = uchConfigRaw["Type"]
		uchConfigMap["value"] = uchConfigRaw["Value"]

		uchConfigMaps = append(uchConfigMaps, uchConfigMap)
	}
	if err := d.Set("uch_config", uchConfigMaps); err != nil {
		return err
	}

	objectRaw, err = albServiceV2.DescribeServerGroupListServerGroupServers(d.Id())
	if err != nil && !NotFoundError(err) {
		return WrapError(err)
	}

	serversRaw, _ := jsonpath.Get("$.Servers", objectRaw)

	serversMaps := make([]map[string]interface{}, 0)
	if serversRaw != nil {
		for _, serversChildRaw := range serversRaw.([]interface{}) {
			serversMap := make(map[string]interface{})
			serversChildRaw := serversChildRaw.(map[string]interface{})
			serversMap["description"] = serversChildRaw["Description"]
			serversMap["port"] = serversChildRaw["Port"]
			serversMap["remote_ip_enabled"] = serversChildRaw["RemoteIpEnabled"]
			serversMap["server_group_id"] = serversChildRaw["ServerGroupId"]
			serversMap["server_id"] = serversChildRaw["ServerId"]
			serversMap["server_ip"] = serversChildRaw["ServerIp"]
			serversMap["server_type"] = serversChildRaw["ServerType"]
			serversMap["status"] = serversChildRaw["Status"]
			serversMap["weight"] = serversChildRaw["Weight"]

			serversMaps = append(serversMaps, serversMap)
		}
	}
	if err := d.Set("servers", serversMaps); err != nil {
		return err
	}

	return nil
}

func resourceAliCloudAlbServerGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	d.Partial(true)

	var err error
	action := "UpdateServerGroupAttribute"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ServerGroupId"] = d.Id()

	request["ClientToken"] = buildClientToken(action)
	if !d.IsNewResource() && d.HasChange("sticky_session_config") {
		update = true
		objectDataLocalMap := make(map[string]interface{})

		if v := d.Get("sticky_session_config"); v != nil {
			cookieTimeout1, _ := jsonpath.Get("$[0].cookie_timeout", v)
			if cookieTimeout1 != nil && (d.HasChange("sticky_session_config.0.cookie_timeout") || cookieTimeout1 != "") && cookieTimeout1.(int) > 0 {
				objectDataLocalMap["CookieTimeout"] = cookieTimeout1
			}
			cookie1, _ := jsonpath.Get("$[0].cookie", v)
			if cookie1 != nil && (d.HasChange("sticky_session_config.0.cookie") || cookie1 != "") {
				objectDataLocalMap["Cookie"] = cookie1
			}
			stickySessionEnabled1, _ := jsonpath.Get("$[0].sticky_session_enabled", v)
			if stickySessionEnabled1 != nil && (d.HasChange("sticky_session_config.0.sticky_session_enabled") || stickySessionEnabled1 != "") {
				objectDataLocalMap["StickySessionEnabled"] = stickySessionEnabled1
			}
			stickySessionType1, _ := jsonpath.Get("$[0].sticky_session_type", v)
			if stickySessionType1 != nil && (d.HasChange("sticky_session_config.0.sticky_session_type") || stickySessionType1 != "") {
				objectDataLocalMap["StickySessionType"] = stickySessionType1
			}

			request["StickySessionConfig"] = objectDataLocalMap
		}
	}

	if !d.IsNewResource() && d.HasChange("health_check_config") {
		update = true
	}
	objectDataLocalMap1 := make(map[string]interface{})

	if v := d.Get("health_check_config"); v != nil {
		healthCheckCodes1, _ := jsonpath.Get("$[0].health_check_codes", d.Get("health_check_config"))
		if healthCheckCodes1 != nil && (d.HasChange("health_check_config.0.health_check_codes") || healthCheckCodes1 != "") {
			objectDataLocalMap1["HealthCheckCodes"] = healthCheckCodes1
		}
		healthCheckHost1, _ := jsonpath.Get("$[0].health_check_host", v)
		if healthCheckHost1 != nil && (d.HasChange("health_check_config.0.health_check_host") || healthCheckHost1 != "") {
			objectDataLocalMap1["HealthCheckHost"] = healthCheckHost1
		}
		healthCheckPath1, _ := jsonpath.Get("$[0].health_check_path", v)
		if healthCheckPath1 != nil && (d.HasChange("health_check_config.0.health_check_path") || healthCheckPath1 != "") {
			objectDataLocalMap1["HealthCheckPath"] = healthCheckPath1
		}
		healthCheckEnabled1, _ := jsonpath.Get("$[0].health_check_enabled", v)
		if healthCheckEnabled1 != nil && (d.HasChange("health_check_config.0.health_check_enabled") || healthCheckEnabled1 != "") {
			objectDataLocalMap1["HealthCheckEnabled"] = healthCheckEnabled1
		}
		healthCheckTimeout1, _ := jsonpath.Get("$[0].health_check_timeout", v)
		if healthCheckTimeout1 != nil && (d.HasChange("health_check_config.0.health_check_timeout") || healthCheckTimeout1 != "") && healthCheckTimeout1.(int) > 0 {
			objectDataLocalMap1["HealthCheckTimeout"] = healthCheckTimeout1
		}
		unhealthyThreshold1, _ := jsonpath.Get("$[0].unhealthy_threshold", v)
		if unhealthyThreshold1 != nil && (d.HasChange("health_check_config.0.unhealthy_threshold") || unhealthyThreshold1 != "") && unhealthyThreshold1.(int) > 0 {
			objectDataLocalMap1["UnhealthyThreshold"] = unhealthyThreshold1
		}
		healthCheckInterval1, _ := jsonpath.Get("$[0].health_check_interval", v)
		if healthCheckInterval1 != nil && (d.HasChange("health_check_config.0.health_check_interval") || healthCheckInterval1 != "") && healthCheckInterval1.(int) > 0 {
			objectDataLocalMap1["HealthCheckInterval"] = healthCheckInterval1
		}
		healthCheckConnectPort1, _ := jsonpath.Get("$[0].health_check_connect_port", v)
		if healthCheckConnectPort1 != nil && (d.HasChange("health_check_config.0.health_check_connect_port") || healthCheckConnectPort1 != "") {
			objectDataLocalMap1["HealthCheckConnectPort"] = healthCheckConnectPort1
		}
		healthCheckHttpVersion1, _ := jsonpath.Get("$[0].health_check_http_version", v)
		if healthCheckHttpVersion1 != nil && (d.HasChange("health_check_config.0.health_check_http_version") || healthCheckHttpVersion1 != "") {
			objectDataLocalMap1["HealthCheckHttpVersion"] = healthCheckHttpVersion1
		}
		healthCheckMethod1, _ := jsonpath.Get("$[0].health_check_method", v)
		if healthCheckMethod1 != nil && (d.HasChange("health_check_config.0.health_check_method") || healthCheckMethod1 != "") {
			objectDataLocalMap1["HealthCheckMethod"] = healthCheckMethod1
		}
		healthCheckProtocol1, _ := jsonpath.Get("$[0].health_check_protocol", v)
		if healthCheckProtocol1 != nil && (d.HasChange("health_check_config.0.health_check_protocol") || healthCheckProtocol1 != "") {
			objectDataLocalMap1["HealthCheckProtocol"] = healthCheckProtocol1
		}
		healthyThreshold1, _ := jsonpath.Get("$[0].healthy_threshold", v)
		if healthyThreshold1 != nil && (d.HasChange("health_check_config.0.healthy_threshold") || healthyThreshold1 != "") && healthyThreshold1.(int) > 0 {
			objectDataLocalMap1["HealthyThreshold"] = healthyThreshold1
		}

		request["HealthCheckConfig"] = objectDataLocalMap1
	}

	if !d.IsNewResource() && d.HasChange("slow_start_config") {
		update = true
		objectDataLocalMap2 := make(map[string]interface{})

		if v := d.Get("slow_start_config"); v != nil {
			slowStartEnabled1, _ := jsonpath.Get("$[0].slow_start_enabled", v)
			if slowStartEnabled1 != nil && (d.HasChange("slow_start_config.0.slow_start_enabled") || slowStartEnabled1 != "") {
				objectDataLocalMap2["SlowStartEnabled"] = slowStartEnabled1
			}
			slowStartDuration1, _ := jsonpath.Get("$[0].slow_start_duration", v)
			if slowStartDuration1 != nil && (d.HasChange("slow_start_config.0.slow_start_duration") || slowStartDuration1 != "") {
				objectDataLocalMap2["SlowStartDuration"] = slowStartDuration1
			}

			request["SlowStartConfig"] = objectDataLocalMap2
		}
	}

	if !d.IsNewResource() && d.HasChange("server_group_name") {
		update = true
	}
	request["ServerGroupName"] = d.Get("server_group_name")
	if !d.IsNewResource() && d.HasChange("connection_drain_config") {
		update = true
		objectDataLocalMap3 := make(map[string]interface{})

		if v := d.Get("connection_drain_config"); v != nil {
			connectionDrainEnabled1, _ := jsonpath.Get("$[0].connection_drain_enabled", v)
			if connectionDrainEnabled1 != nil && (d.HasChange("connection_drain_config.0.connection_drain_enabled") || connectionDrainEnabled1 != "") {
				objectDataLocalMap3["ConnectionDrainEnabled"] = connectionDrainEnabled1
			}
			connectionDrainTimeout1, _ := jsonpath.Get("$[0].connection_drain_timeout", v)
			if connectionDrainTimeout1 != nil && (d.HasChange("connection_drain_config.0.connection_drain_timeout") || connectionDrainTimeout1 != "") {
				objectDataLocalMap3["ConnectionDrainTimeout"] = connectionDrainTimeout1
			}

			request["ConnectionDrainConfig"] = objectDataLocalMap3
		}
	}

	if !d.IsNewResource() && d.HasChange("uch_config") {
		update = true
		objectDataLocalMap4 := make(map[string]interface{})
		if v := d.Get("uch_config"); v != nil {
			value1, _ := jsonpath.Get("$[0].value", v)
			if value1 != nil && (d.HasChange("uch_config.0.value") || value1 != "") {
				objectDataLocalMap4["Value"] = value1
			}
			type1, _ := jsonpath.Get("$[0].type", v)
			if type1 != nil && (d.HasChange("uch_config.0.type") || type1 != "") {
				objectDataLocalMap4["Type"] = type1
			}

			request["UchConfig"] = objectDataLocalMap4
		}
	}

	if !d.IsNewResource() && d.HasChange("scheduler") {
		update = true
		request["Scheduler"] = d.Get("scheduler")
	}

	if !d.IsNewResource() && d.HasChange("upstream_keepalive_enabled") {
		update = true
		request["UpstreamKeepaliveEnabled"] = d.Get("upstream_keepalive_enabled")
	}

	if !d.IsNewResource() && d.HasChange("cross_zone_enabled") {
		update = true
		request["CrossZoneEnabled"] = d.Get("cross_zone_enabled")
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Alb", "2020-06-16", action, query, request, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"SystemBusy", "IdempotenceProcessing", "IncorrectStatus.ServerGroup"}) || NeedRetry(err) {
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
		albServiceV2 := AlbServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, albServiceV2.AlbServerGroupStateRefreshFunc(d.Id(), "ServerGroupStatus", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	update = false
	action = "MoveResourceGroup"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ResourceId"] = d.Id()

	if _, ok := d.GetOk("resource_group_id"); ok && !d.IsNewResource() && d.HasChange("resource_group_id") {
		update = true
	}
	request["NewResourceGroupId"] = d.Get("resource_group_id")
	request["ResourceType"] = "servergroup"
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Alb", "2020-06-16", action, query, request, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"NotExist.ResourceGroup"}) || NeedRetry(err) {
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
	action = "ApplyHealthCheckTemplateToServerGroup"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ServerGroupId"] = d.Id()

	request["ClientToken"] = buildClientToken(action)
	request["HealthCheckTemplateId"] = d.Get("health_check_template_id")
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Alb", "2020-06-16", action, query, request, true)
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
	}

	if d.HasChange("servers") {
		oldEntry, newEntry := d.GetChange("servers")
		oldEntrySet := oldEntry.(*schema.Set)
		newEntrySet := newEntry.(*schema.Set)
		removed := oldEntrySet.Difference(newEntrySet)
		added := newEntrySet.Difference(oldEntrySet)

		if removed.Len() > 0 {
			action := "RemoveServersFromServerGroup"
			request = make(map[string]interface{})
			query = make(map[string]interface{})
			request["ServerGroupId"] = d.Id()

			request["ClientToken"] = buildClientToken(action)
			localData := removed.List()
			serversMapsArray := make([]interface{}, 0)
			for _, dataLoop := range localData {
				dataLoopTmp := dataLoop.(map[string]interface{})
				dataLoopMap := make(map[string]interface{})
				if dataLoopTmp["port"].(int) > 0 {
					dataLoopMap["Port"] = dataLoopTmp["port"]
				}
				dataLoopMap["ServerId"] = dataLoopTmp["server_id"]
				dataLoopMap["ServerIp"] = dataLoopTmp["server_ip"]
				dataLoopMap["ServerType"] = dataLoopTmp["server_type"]
				serversMapsArray = append(serversMapsArray, dataLoopMap)
			}
			request["Servers"] = serversMapsArray

			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = client.RpcPost("Alb", "2020-06-16", action, query, request, true)
				if err != nil {
					if IsExpectedErrors(err, []string{"SystemBusy", "IdempotenceProcessing", "IncorrectStatus.ServerGroup"}) || NeedRetry(err) {
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
			albServiceV2 := AlbServiceV2{client}
			stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, albServiceV2.AlbServerGroupStateRefreshFunc(d.Id(), "ServerGroupStatus", []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}

		}

		if added.Len() > 0 {
			action := "AddServersToServerGroup"
			request = make(map[string]interface{})
			query = make(map[string]interface{})
			request["ServerGroupId"] = d.Id()

			request["ClientToken"] = buildClientToken(action)
			localData := added.List()
			serversMapsArray := make([]interface{}, 0)
			for _, dataLoop := range localData {
				dataLoopTmp := dataLoop.(map[string]interface{})
				dataLoopMap := make(map[string]interface{})
				if dataLoopTmp["port"].(int) > 0 {
					dataLoopMap["Port"] = dataLoopTmp["port"]
				}
				dataLoopMap["ServerId"] = dataLoopTmp["server_id"]
				dataLoopMap["Weight"] = dataLoopTmp["weight"]
				dataLoopMap["RemoteIpEnabled"] = dataLoopTmp["remote_ip_enabled"]
				if dataLoopTmp["description"] != "" {
					dataLoopMap["Description"] = dataLoopTmp["description"]
				}
				dataLoopMap["ServerIp"] = dataLoopTmp["server_ip"]
				dataLoopMap["ServerType"] = dataLoopTmp["server_type"]
				serversMapsArray = append(serversMapsArray, dataLoopMap)
			}
			request["Servers"] = serversMapsArray

			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = client.RpcPost("Alb", "2020-06-16", action, query, request, true)
				if err != nil {
					if IsExpectedErrors(err, []string{"SystemBusy", "IdempotenceProcessing", "IncorrectStatus.ServerGroup"}) || NeedRetry(err) {
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
			albServiceV2 := AlbServiceV2{client}
			stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, albServiceV2.AlbServerGroupStateRefreshFunc(d.Id(), "ServerGroupStatus", []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}

		}

	}
	if d.HasChange("tags") {
		albServiceV2 := AlbServiceV2{client}
		if err := albServiceV2.SetResourceTags(d, "servergroup"); err != nil {
			return WrapError(err)
		}
	}
	d.Partial(false)
	return resourceAliCloudAlbServerGroupRead(d, meta)
}

func resourceAliCloudAlbServerGroupDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteServerGroup"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["ServerGroupId"] = d.Id()

	request["ClientToken"] = buildClientToken(action)

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Alb", "2020-06-16", action, query, request, true)
		request["ClientToken"] = buildClientToken(action)

		if err != nil {
			if IsExpectedErrors(err, []string{"SystemBusy", "ResourceInUse.ServerGroup", "IncorrectStatus.ServerGroup", "IdempotenceProcessing"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		if IsExpectedErrors(err, []string{"ResourceNotFound.ServerGroup"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	albServiceV2 := AlbServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, albServiceV2.AlbServerGroupStateRefreshFunc(d.Id(), "ServerGroupStatus", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}

func resourceAlbServerGroupCustomizeDiff(diff *schema.ResourceDiff, v interface{}) error {
	groupType := diff.Get("server_group_type").(string)
	if groupType == "Fc" {
		// Fc load balancers do not support vpc_id, protocol
		if diff.Get("vpc_id") != "" {
			return fmt.Errorf("fc server group type do not support vpc_id")
		}

		if diff.Get("protocol") != "" {
			return fmt.Errorf("fc server group type do not support protocol")
		}
	}

	return nil
}
