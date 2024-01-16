package alicloud

import (
	"fmt"
	"log"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
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
			Create: schema.DefaultTimeout(6 * time.Minute),
			Delete: schema.DefaultTimeout(6 * time.Minute),
			Update: schema.DefaultTimeout(6 * time.Minute),
		},
		CustomizeDiff: resourceAlbServerGroupCustomizeDiff,
		Schema: map[string]*schema.Schema{
			"server_group_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"server_group_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"Instance", "Ip", "Fc"}, false),
			},
			"protocol": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"HTTP", "HTTPS", "gRPC"}, false),
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"scheduler": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"Sch", "Wlc", "Wrr"}, false),
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"sticky_session_config": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"sticky_session_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"sticky_session_type": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: StringInSlice([]string{"Server", "Insert"}, false),
						},
						"cookie": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"cookie_timeout": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.Any(validation.IntInSlice([]int{0}), validation.IntBetween(1, 86400)),
						},
					},
				},
			},
			"health_check_config": {
				Type:     schema.TypeSet,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"health_check_enabled": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"health_check_connect_port": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: IntBetween(0, 65535),
						},
						"health_check_host": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"health_check_http_version": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: StringInSlice([]string{"HTTP1.0", "HTTP1.1"}, false),
						},
						"health_check_interval": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: IntBetween(1, 50),
						},
						"health_check_method": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: StringInSlice([]string{"GET", "POST", "HEAD"}, false),
						},
						"health_check_path": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"health_check_protocol": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"health_check_timeout": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: IntBetween(1, 300),
						},
						"healthy_threshold": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: IntBetween(2, 10),
						},
						"unhealthy_threshold": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: IntBetween(2, 10),
						},
						"health_check_codes": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"servers": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"server_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"server_type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: StringInSlice([]string{"Ecs", "Eni", "Eci", "Ip", "Fc"}, false),
						},
						"server_ip": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"port": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: IntBetween(0, 65535),
						},
						"remote_ip_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"weight": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      100,
							ValidateFunc: IntBetween(0, 100),
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"tags": tagsSchema(),
			"dry_run": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudAlbServerGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	albService := AlbService{client}
	var response map[string]interface{}
	action := "CreateServerGroup"
	request := make(map[string]interface{})
	conn, err := client.NewAlbClient()
	if err != nil {
		return WrapError(err)
	}

	request["ClientToken"] = buildClientToken("CreateServerGroup")
	request["ServerGroupName"] = d.Get("server_group_name")

	if v, ok := d.GetOk("server_group_type"); ok {
		request["ServerGroupType"] = v
	}

	if v, ok := d.GetOk("protocol"); ok {
		request["Protocol"] = v
	}

	if v, ok := d.GetOk("vpc_id"); ok {
		request["VpcId"] = v
	}

	if v, ok := d.GetOk("scheduler"); ok {
		request["Scheduler"] = v
	}

	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}

	if v, ok := d.GetOk("sticky_session_config"); ok {
		stickySessionConfigMap := make(map[string]interface{})
		for _, stickySessionConfigList := range v.(*schema.Set).List() {
			stickySessionConfigArg := stickySessionConfigList.(map[string]interface{})

			if stickySessionEnabled, ok := stickySessionConfigArg["sticky_session_enabled"]; ok {
				stickySessionConfigMap["StickySessionEnabled"] = stickySessionEnabled
			}

			if stickySessionConfigMap["StickySessionEnabled"] == true {
				if stickySessionType, ok := stickySessionConfigArg["sticky_session_type"]; ok {
					stickySessionConfigMap["StickySessionType"] = stickySessionType
				}

				if stickySessionConfigMap["StickySessionType"] == "Server" {
					if cookie, ok := stickySessionConfigArg["cookie"]; ok {
						stickySessionConfigMap["Cookie"] = cookie
					}
				}

				if stickySessionConfigMap["StickySessionType"] == "Insert" {
					if cookieTimeout, ok := stickySessionConfigArg["cookie_timeout"]; ok {
						stickySessionConfigMap["CookieTimeout"] = cookieTimeout
					}
				}
			}
		}

		request["StickySessionConfig"] = stickySessionConfigMap
	}

	healthCheckConfig := d.Get("health_check_config")
	healthCheckConfigMap := make(map[string]interface{})
	for _, healthCheckConfigList := range healthCheckConfig.(*schema.Set).List() {
		healthCheckConfigArg := healthCheckConfigList.(map[string]interface{})

		healthCheckConfigMap["HealthCheckEnabled"] = healthCheckConfigArg["health_check_enabled"]

		if healthCheckConfigMap["HealthCheckEnabled"] == true {
			if healthCheckConnectPort, ok := healthCheckConfigArg["health_check_connect_port"]; ok {
				healthCheckConfigMap["HealthCheckConnectPort"] = healthCheckConnectPort
			}

			if healthCheckHost, ok := healthCheckConfigArg["health_check_host"]; ok {
				healthCheckConfigMap["HealthCheckHost"] = healthCheckHost
			}

			if healthCheckHttpVersion, ok := healthCheckConfigArg["health_check_http_version"]; ok {
				healthCheckConfigMap["HealthCheckHttpVersion"] = healthCheckHttpVersion
			}

			if healthCheckInterval, ok := healthCheckConfigArg["health_check_interval"]; ok {
				healthCheckConfigMap["HealthCheckInterval"] = healthCheckInterval
			}

			if healthCheckMethod, ok := healthCheckConfigArg["health_check_method"]; ok {
				healthCheckConfigMap["HealthCheckMethod"] = healthCheckMethod
			}

			if healthCheckPath, ok := healthCheckConfigArg["health_check_path"]; ok {
				healthCheckConfigMap["HealthCheckPath"] = healthCheckPath
			}

			if healthCheckProtocol, ok := healthCheckConfigArg["health_check_protocol"]; ok {
				healthCheckConfigMap["HealthCheckProtocol"] = healthCheckProtocol
			}

			if healthCheckTimeout, ok := healthCheckConfigArg["health_check_timeout"]; ok {
				healthCheckConfigMap["HealthCheckTimeout"] = healthCheckTimeout
			}

			if healthyThreshold, ok := healthCheckConfigArg["healthy_threshold"]; ok {
				healthCheckConfigMap["HealthyThreshold"] = healthyThreshold
			}

			if unhealthyThreshold, ok := healthCheckConfigArg["unhealthy_threshold"]; ok {
				healthCheckConfigMap["UnhealthyThreshold"] = unhealthyThreshold
			}

			if healthCheckCodes, ok := healthCheckConfigArg["health_check_codes"]; ok {
				healthCheckConfigMap["HealthCheckCodes"] = healthCheckCodes
			}
		}
	}

	request["HealthCheckConfig"] = healthCheckConfigMap

	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-06-16"), StringPointer("AK"), nil, request, &runtime)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_alb_server_group", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["ServerGroupId"]))

	stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, albService.AlbServerGroupStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudAlbServerGroupUpdate(d, meta)
}

func resourceAliCloudAlbServerGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	albService := AlbService{client}

	object, err := albService.DescribeAlbServerGroup(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_alb_server_group albService.DescribeAlbServerGroup Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("server_group_name", object["ServerGroupName"])
	d.Set("server_group_type", object["ServerGroupType"])
	d.Set("protocol", object["Protocol"])
	d.Set("vpc_id", object["VpcId"])
	d.Set("scheduler", object["Scheduler"])
	d.Set("resource_group_id", object["ResourceGroupId"])

	if stickySessionConfig, ok := object["StickySessionConfig"]; ok {
		stickySessionConfigMaps := make([]map[string]interface{}, 0)
		stickySessionConfigArg := stickySessionConfig.(map[string]interface{})
		stickySessionConfigMap := make(map[string]interface{})

		if stickySessionEnabled, ok := stickySessionConfigArg["StickySessionEnabled"]; ok {
			stickySessionConfigMap["sticky_session_enabled"] = stickySessionEnabled
		}

		if stickySessionType, ok := stickySessionConfigArg["StickySessionType"]; ok {
			stickySessionConfigMap["sticky_session_type"] = stickySessionType
		}

		if cookie, ok := stickySessionConfigArg["Cookie"]; ok {
			stickySessionConfigMap["cookie"] = cookie
		}

		if cookieTimeout, ok := stickySessionConfigArg["CookieTimeout"]; ok {
			stickySessionConfigMap["cookie_timeout"] = cookieTimeout
		}

		stickySessionConfigMaps = append(stickySessionConfigMaps, stickySessionConfigMap)

		d.Set("sticky_session_config", stickySessionConfigMaps)
	}

	if healthCheckConfig, ok := object["HealthCheckConfig"]; ok {
		healthCheckConfigMaps := make([]map[string]interface{}, 0)
		healthCheckConfigArg := healthCheckConfig.(map[string]interface{})
		healthCheckConfigMap := make(map[string]interface{})

		if healthCheckEnabled, ok := healthCheckConfigArg["HealthCheckEnabled"]; ok {
			healthCheckConfigMap["health_check_enabled"] = healthCheckEnabled
		}

		if healthCheckConnectPort, ok := healthCheckConfigArg["HealthCheckConnectPort"]; ok {
			healthCheckConfigMap["health_check_connect_port"] = formatInt(healthCheckConnectPort)
		}

		if healthCheckHost, ok := healthCheckConfigArg["HealthCheckHost"]; ok {
			healthCheckConfigMap["health_check_host"] = healthCheckHost
		}

		if healthCheckHttpVersion, ok := healthCheckConfigArg["HealthCheckHttpVersion"]; ok {
			healthCheckConfigMap["health_check_http_version"] = healthCheckHttpVersion
		}

		if healthCheckInterval, ok := healthCheckConfigArg["HealthCheckInterval"]; ok {
			healthCheckConfigMap["health_check_interval"] = formatInt(healthCheckInterval)
		}

		if healthCheckMethod, ok := healthCheckConfigArg["HealthCheckMethod"]; ok {
			healthCheckConfigMap["health_check_method"] = healthCheckMethod
		}

		if healthCheckPath, ok := healthCheckConfigArg["HealthCheckPath"]; ok {
			healthCheckConfigMap["health_check_path"] = healthCheckPath
		}

		if healthCheckProtocol, ok := healthCheckConfigArg["HealthCheckProtocol"]; ok {
			healthCheckConfigMap["health_check_protocol"] = healthCheckProtocol
		}

		if healthCheckTimeout, ok := healthCheckConfigArg["HealthCheckTimeout"]; ok {
			healthCheckConfigMap["health_check_timeout"] = formatInt(healthCheckTimeout)
		}

		if healthyThreshold, ok := healthCheckConfigArg["HealthyThreshold"]; ok {
			healthCheckConfigMap["healthy_threshold"] = formatInt(healthyThreshold)
		}

		if unhealthyThreshold, ok := healthCheckConfigArg["UnhealthyThreshold"]; ok {
			healthCheckConfigMap["unhealthy_threshold"] = formatInt(unhealthyThreshold)
		}

		if healthCheckCodes, ok := healthCheckConfigArg["HealthCheckCodes"]; ok {
			healthCheckConfigMap["health_check_codes"] = healthCheckCodes
		}

		healthCheckConfigMaps = append(healthCheckConfigMaps, healthCheckConfigMap)

		d.Set("health_check_config", healthCheckConfigMaps)
	}

	serversList, err := albService.ListServerGroupServers(d.Id())
	if err != nil {
		return WrapError(err)
	}

	if len(serversList) > 0 {
		serversMaps := make([]map[string]interface{}, 0)
		for _, servers := range serversList {
			serversArg := servers.(map[string]interface{})
			serversMap := map[string]interface{}{}

			serversMap["server_id"] = serversArg["ServerId"]
			serversMap["server_type"] = serversArg["ServerType"]

			if serverIp, ok := serversArg["ServerIp"]; ok {
				serversMap["server_ip"] = serverIp
			}

			if port, ok := serversArg["Port"]; ok {
				serversMap["port"] = port
			}

			if remoteIpEnabled, ok := serversArg["RemoteIpEnabled"]; ok {
				serversMap["remote_ip_enabled"] = remoteIpEnabled
			}

			if weight, ok := serversArg["Weight"]; ok {
				serversMap["weight"] = weight
			}

			if description, ok := serversArg["Description"]; ok {
				serversMap["description"] = description
			}

			if status, ok := serversArg["Status"]; ok {
				serversMap["status"] = status
			}

			serversMaps = append(serversMaps, serversMap)
		}

		d.Set("servers", serversMaps)
	}

	d.Set("status", object["ServerGroupStatus"])

	listTagResourcesObject, err := albService.ListTagResources(d.Id(), "servergroup")
	if err != nil {
		return WrapError(err)
	}

	d.Set("tags", tagsToMap(listTagResourcesObject))

	return nil
}

func resourceAliCloudAlbServerGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	albService := AlbService{client}
	var response map[string]interface{}
	d.Partial(true)

	if d.HasChange("tags") {
		if err := albService.SetResourceTags(d, "servergroup"); err != nil {
			return WrapError(err)
		}

		d.SetPartial("tags")
	}

	update := false
	updateServerGroupAttributeReq := map[string]interface{}{
		"ClientToken":   buildClientToken("UpdateServerGroupAttribute"),
		"ServerGroupId": d.Id(),
	}

	if !d.IsNewResource() && d.HasChange("server_group_name") {
		update = true
	}
	if v, ok := d.GetOk("server_group_name"); ok {
		updateServerGroupAttributeReq["ServerGroupName"] = v
	}

	if !d.IsNewResource() && d.HasChange("scheduler") {
		update = true
	}
	if v, ok := d.GetOk("scheduler"); ok {
		updateServerGroupAttributeReq["Scheduler"] = v
	}

	if !d.IsNewResource() && d.HasChange("sticky_session_config") {
		update = true
	}
	if v, ok := d.GetOk("sticky_session_config"); ok {
		stickySessionConfigMap := make(map[string]interface{})
		for _, stickySessionConfigList := range v.(*schema.Set).List() {
			stickySessionConfigArg := stickySessionConfigList.(map[string]interface{})

			if stickySessionEnabled, ok := stickySessionConfigArg["sticky_session_enabled"]; ok {
				stickySessionConfigMap["StickySessionEnabled"] = stickySessionEnabled
			}

			if stickySessionConfigMap["StickySessionEnabled"] == true {
				if stickySessionType, ok := stickySessionConfigArg["sticky_session_type"]; ok {
					stickySessionConfigMap["StickySessionType"] = stickySessionType
				}

				if stickySessionConfigMap["StickySessionType"] == "Server" {
					if cookie, ok := stickySessionConfigArg["cookie"]; ok {
						stickySessionConfigMap["Cookie"] = cookie
					}
				}

				if stickySessionConfigMap["StickySessionType"] == "Insert" {
					if cookieTimeout, ok := stickySessionConfigArg["cookie_timeout"]; ok {
						stickySessionConfigMap["CookieTimeout"] = cookieTimeout
					}
				}
			}
		}

		updateServerGroupAttributeReq["StickySessionConfig"] = stickySessionConfigMap
	}

	if !d.IsNewResource() && d.HasChange("health_check_config") {
		update = true
	}
	healthCheckConfig := d.Get("health_check_config")
	healthCheckConfigMap := make(map[string]interface{})
	for _, healthCheckConfigList := range healthCheckConfig.(*schema.Set).List() {
		healthCheckConfigArg := healthCheckConfigList.(map[string]interface{})

		healthCheckConfigMap["HealthCheckEnabled"] = healthCheckConfigArg["health_check_enabled"]

		if healthCheckConfigMap["HealthCheckEnabled"] == true {
			if healthCheckConnectPort, ok := healthCheckConfigArg["health_check_connect_port"]; ok {
				healthCheckConfigMap["HealthCheckConnectPort"] = healthCheckConnectPort
			}

			if healthCheckHost, ok := healthCheckConfigArg["health_check_host"]; ok {
				healthCheckConfigMap["HealthCheckHost"] = healthCheckHost
			}

			if healthCheckHttpVersion, ok := healthCheckConfigArg["health_check_http_version"]; ok {
				healthCheckConfigMap["HealthCheckHttpVersion"] = healthCheckHttpVersion
			}

			if healthCheckInterval, ok := healthCheckConfigArg["health_check_interval"]; ok {
				healthCheckConfigMap["HealthCheckInterval"] = healthCheckInterval
			}

			if healthCheckMethod, ok := healthCheckConfigArg["health_check_method"]; ok {
				healthCheckConfigMap["HealthCheckMethod"] = healthCheckMethod
			}

			if healthCheckPath, ok := healthCheckConfigArg["health_check_path"]; ok {
				healthCheckConfigMap["HealthCheckPath"] = healthCheckPath
			}

			if healthCheckProtocol, ok := healthCheckConfigArg["health_check_protocol"]; ok {
				healthCheckConfigMap["HealthCheckProtocol"] = healthCheckProtocol
			}

			if healthCheckTimeout, ok := healthCheckConfigArg["health_check_timeout"]; ok {
				healthCheckConfigMap["HealthCheckTimeout"] = healthCheckTimeout
			}

			if healthyThreshold, ok := healthCheckConfigArg["healthy_threshold"]; ok {
				healthCheckConfigMap["HealthyThreshold"] = healthyThreshold
			}

			if unhealthyThreshold, ok := healthCheckConfigArg["unhealthy_threshold"]; ok {
				healthCheckConfigMap["UnhealthyThreshold"] = unhealthyThreshold
			}

			if healthCheckCodes, ok := healthCheckConfigArg["health_check_codes"]; ok {
				healthCheckConfigMap["HealthCheckCodes"] = healthCheckCodes
			}
		}
	}

	updateServerGroupAttributeReq["HealthCheckConfig"] = healthCheckConfigMap

	if v, ok := d.GetOkExists("dry_run"); ok {
		updateServerGroupAttributeReq["DryRun"] = v
	}

	if update {
		action := "UpdateServerGroupAttribute"
		conn, err := client.NewAlbClient()
		if err != nil {
			return WrapError(err)
		}

		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-06-16"), StringPointer("AK"), nil, updateServerGroupAttributeReq, &runtime)
			if err != nil {
				if IsExpectedErrors(err, []string{"IncorrectStatus.ServerGroup", "ResourceNotFound.ServerGroup"}) || NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, updateServerGroupAttributeReq)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, albService.AlbServerGroupStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}

		d.SetPartial("server_group_name")
		d.SetPartial("scheduler")
		d.SetPartial("sticky_session_config")
		d.SetPartial("health_check_config")
	}

	update = false
	moveResourceGroup := map[string]interface{}{
		"ResourceType": "servergroup",
		"ResourceId":   d.Id(),
	}

	if !d.IsNewResource() && d.HasChange("resource_group_id") {
		update = true
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		moveResourceGroup["NewResourceGroupId"] = v
	}

	if update {
		action := "MoveResourceGroup"
		conn, err := client.NewAlbClient()
		if err != nil {
			return WrapError(err)
		}

		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-06-16"), StringPointer("AK"), nil, moveResourceGroup, &runtime)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, moveResourceGroup)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		d.SetPartial("resource_group_id")
	}

	update = false
	if d.HasChange("servers") {
		removed, added := d.GetChange("servers")
		removeServersFromServerGroupReq := map[string]interface{}{
			"ClientToken":   buildClientToken("RemoveServersFromServerGroup"),
			"ServerGroupId": d.Id(),
		}

		removeServersMaps := make([]map[string]interface{}, 0)
		for _, servers := range removed.(*schema.Set).List() {
			update = true
			serversMap := map[string]interface{}{}
			serversArg := servers.(map[string]interface{})

			serversMap["ServerId"] = serversArg["server_id"]
			serversMap["ServerType"] = serversArg["server_type"]

			if serverIp, ok := serversArg["server_ip"]; ok {
				serversMap["ServerIp"] = serverIp
			}

			if port, ok := serversArg["port"]; ok {
				serversMap["Port"] = port
			}

			removeServersMaps = append(removeServersMaps, serversMap)
		}

		removeServersFromServerGroupReq["Servers"] = removeServersMaps

		if update {
			action := "RemoveServersFromServerGroup"
			conn, err := client.NewAlbClient()
			if err != nil {
				return WrapError(err)
			}

			runtime := util.RuntimeOptions{}
			runtime.SetAutoretry(true)
			wait := incrementalWait(3*time.Second, 3*time.Second)
			err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
				response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-06-16"), StringPointer("AK"), nil, removeServersFromServerGroupReq, &runtime)
				if err != nil {
					if IsExpectedErrors(err, []string{"IncorrectStatus.ServerGroup"}) || NeedRetry(err) {
						wait()
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				return nil
			})
			addDebug(action, response, removeServersFromServerGroupReq)

			stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, albService.AlbServerGroupStateRefreshFunc(d.Id(), []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}
		}

		update = false
		addServersToServerGroupReq := map[string]interface{}{
			"ClientToken":   buildClientToken("AddServersToServerGroup"),
			"ServerGroupId": d.Id(),
		}

		addServersMaps := make([]map[string]interface{}, 0)
		for _, servers := range added.(*schema.Set).List() {
			update = true
			serversArg := servers.(map[string]interface{})
			serversMap := map[string]interface{}{}

			serversMap["ServerId"] = serversArg["server_id"]
			serversMap["ServerType"] = serversArg["server_type"]

			if serverIp, ok := serversArg["server_ip"]; ok {
				serversMap["ServerIp"] = serverIp
			}

			if port, ok := serversArg["port"]; ok {
				serversMap["Port"] = port
			}

			if remoteIpEnabled, ok := serversArg["remote_ip_enabled"]; ok {
				serversMap["RemoteIpEnabled"] = remoteIpEnabled
			}

			if weight, ok := serversArg["weight"]; ok {
				serversMap["Weight"] = weight
			}

			if description, ok := serversArg["description"]; ok && fmt.Sprint(description) != "" {
				serversMap["Description"] = description
			}

			addServersMaps = append(addServersMaps, serversMap)
		}

		addServersToServerGroupReq["Servers"] = addServersMaps

		if update {
			action := "AddServersToServerGroup"
			conn, err := client.NewAlbClient()
			if err != nil {
				return WrapError(err)
			}

			runtime := util.RuntimeOptions{}
			runtime.SetAutoretry(true)
			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
				response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-06-16"), StringPointer("AK"), nil, addServersToServerGroupReq, &runtime)
				if err != nil {
					if IsExpectedErrors(err, []string{"IncorrectStatus.ServerGroup"}) || NeedRetry(err) {
						wait()
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				return nil
			})
			addDebug(action, response, addServersToServerGroupReq)

			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}

			stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, albService.AlbServerGroupStateRefreshFunc(d.Id(), []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}
		}

		d.SetPartial("servers")
	}

	d.Partial(false)

	return resourceAliCloudAlbServerGroupRead(d, meta)
}

func resourceAliCloudAlbServerGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	albService := AlbService{client}
	action := "DeleteServerGroup"
	var response map[string]interface{}

	conn, err := client.NewAlbClient()
	if err != nil {
		return WrapError(err)
	}

	request := map[string]interface{}{
		"ClientToken":   buildClientToken("DeleteServerGroup"),
		"ServerGroupId": d.Id(),
	}

	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-06-16"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if IsExpectedErrors(err, []string{"IncorrectStatus.ServerGroup", "ResourceInUse.ServerGroup"}) || NeedRetry(err) {
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

	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, albService.AlbServerGroupStateRefreshFunc(d.Id(), []string{}))
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
