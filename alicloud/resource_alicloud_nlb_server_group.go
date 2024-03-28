// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
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

func resourceAliCloudNlbServerGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudNlbServerGroupCreate,
		Read:   resourceAliCloudNlbServerGroupRead,
		Update: resourceAliCloudNlbServerGroupUpdate,
		Delete: resourceAliCloudNlbServerGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"server_group_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"server_group_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"Instance", "Ip"}, false),
			},
			"protocol": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"TCP", "UDP", "TCPSSL"}, false),
			},
			"address_ip_version": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"DualStack", "Ipv4"}, false),
			},
			"scheduler": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"Wrr", "Rr", "Qch", "Tch"}, false),
			},
			"any_port_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"connection_drain_enabled": {
				Type:          schema.TypeBool,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"connection_drain"},
			},
			"connection_drain_timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: IntBetween(10, 900),
			},
			"preserve_client_ip_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"health_check": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"health_check_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"health_check_type": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: StringInSlice([]string{"TCP", "HTTP"}, false),
						},
						"http_check_method": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: StringInSlice([]string{"GET", "HEAD"}, false),
						},
						"health_check_domain": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"health_check_url": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"health_check_connect_port": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: IntBetween(0, 65535),
						},
						"health_check_interval": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: IntBetween(1, 50),
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
						"health_check_connect_timeout": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: IntBetween(1, 300),
						},
						"health_check_http_code": {
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
							ValidateFunc: StringInSlice([]string{"Ecs", "Eni", "Eci", "Ip"}, false),
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
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"connection_drain": {
				Type:       schema.TypeBool,
				Optional:   true,
				Computed:   true,
				Deprecated: "Field 'connection_drain' has been deprecated since provider version 1.214.0. New field 'connection_drain_enabled' instead.",
			},
		},
	}
}

func resourceAliCloudNlbServerGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	nlbServiceV2 := NlbServiceV2{client}
	var response map[string]interface{}
	action := "CreateServerGroup"
	request := make(map[string]interface{})
	conn, err := client.NewNlbClient()
	if err != nil {
		return WrapError(err)
	}

	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken("CreateServerGroup")
	request["ServerGroupName"] = d.Get("server_group_name")
	request["VpcId"] = d.Get("vpc_id")

	if v, ok := d.GetOk("server_group_type"); ok {
		request["ServerGroupType"] = v
	}

	if v, ok := d.GetOk("protocol"); ok {
		request["Protocol"] = v
	}

	if v, ok := d.GetOk("address_ip_version"); ok {
		request["AddressIPVersion"] = v
	}

	if v, ok := d.GetOk("scheduler"); ok {
		request["Scheduler"] = v
	}

	if v, ok := d.GetOkExists("any_port_enabled"); ok {
		request["AnyPortEnabled"] = v
	}

	if v, ok := d.GetOkExists("connection_drain_enabled"); ok {
		request["ConnectionDrainEnabled"] = v
	} else if v, ok := d.GetOkExists("connection_drain"); ok {
		request["ConnectionDrainEnabled"] = v
	}

	if v, ok := d.GetOk("connection_drain_timeout"); ok {
		request["ConnectionDrainTimeout"] = v
	}

	if v, ok := d.GetOkExists("preserve_client_ip_enabled"); ok {
		request["PreserveClientIpEnabled"] = v
	}

	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}

	if _, ok := d.GetOk("health_check"); ok {
		healthCheckMap := make(map[string]interface{})

		if healthCheckEnabled, ok := d.GetOkExists("health_check.0.health_check_enabled"); ok {
			healthCheckMap["HealthCheckEnabled"] = healthCheckEnabled
		}

		if healthCheckMap["HealthCheckEnabled"] == true {
			if healthCheckType, ok := d.GetOk("health_check.0.health_check_type"); ok {
				healthCheckMap["HealthCheckType"] = healthCheckType
			}

			if httpCheckMethod, ok := d.GetOk("health_check.0.http_check_method"); ok {
				healthCheckMap["HttpCheckMethod"] = httpCheckMethod
			}

			if healthCheckDomain, ok := d.GetOk("health_check.0.health_check_domain"); ok {
				healthCheckMap["HealthCheckDomain"] = healthCheckDomain
			}

			if healthCheckUrl, ok := d.GetOk("health_check.0.health_check_url"); ok {
				healthCheckMap["HealthCheckUrl"] = healthCheckUrl
			}

			if healthCheckConnectPort, ok := d.GetOkExists("health_check.0.health_check_connect_port"); ok {
				healthCheckMap["HealthCheckConnectPort"] = healthCheckConnectPort
			}

			if healthCheckInterval, ok := d.GetOkExists("health_check.0.health_check_interval"); ok {
				healthCheckMap["HealthCheckInterval"] = healthCheckInterval
			}

			if healthyThreshold, ok := d.GetOkExists("health_check.0.healthy_threshold"); ok {
				healthCheckMap["HealthyThreshold"] = healthyThreshold
			}

			if unhealthyThreshold, ok := d.GetOkExists("health_check.0.unhealthy_threshold"); ok {
				healthCheckMap["UnhealthyThreshold"] = unhealthyThreshold
			}

			if healthCheckConnectTimeout, ok := d.GetOkExists("health_check.0.health_check_connect_timeout"); ok {
				healthCheckMap["HealthCheckConnectTimeout"] = healthCheckConnectTimeout
			}

			if healthCheckHttpCode, ok := d.GetOk("health_check.0.health_check_http_code"); ok {
				healthCheckMap["HealthCheckHttpCode"] = healthCheckHttpCode
			}
		}

		request["HealthCheckConfig"] = healthCheckMap
	}

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2022-04-30"), StringPointer("AK"), nil, request, &runtime)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_nlb_server_group", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["ServerGroupId"]))

	stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, nlbServiceV2.NlbServerGroupStateRefreshFunc(d.Id(), "ServerGroupStatus", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudNlbServerGroupUpdate(d, meta)
}

func resourceAliCloudNlbServerGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	nlbServiceV2 := NlbServiceV2{client}

	object, err := nlbServiceV2.DescribeNlbServerGroup(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_nlb_server_group DescribeNlbServerGroup Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("server_group_name", object["ServerGroupName"])
	d.Set("vpc_id", object["VpcId"])
	d.Set("server_group_type", object["ServerGroupType"])
	d.Set("protocol", object["Protocol"])
	d.Set("address_ip_version", object["AddressIPVersion"])
	d.Set("scheduler", object["Scheduler"])
	d.Set("any_port_enabled", object["AnyPortEnabled"])
	d.Set("connection_drain_enabled", object["ConnectionDrainEnabled"])
	d.Set("connection_drain_timeout", object["ConnectionDrainTimeout"])
	d.Set("preserve_client_ip_enabled", object["PreserveClientIpEnabled"])
	d.Set("resource_group_id", object["ResourceGroupId"])

	if healthCheck, ok := object["HealthCheck"]; ok {
		healthCheckMaps := make([]map[string]interface{}, 0)
		healthCheckArg := healthCheck.(map[string]interface{})
		healthCheckMap := make(map[string]interface{})

		if healthCheckEnabled, ok := healthCheckArg["HealthCheckEnabled"]; ok {
			healthCheckMap["health_check_enabled"] = healthCheckEnabled
		}

		if healthCheckType, ok := healthCheckArg["HealthCheckType"]; ok {
			healthCheckMap["health_check_type"] = healthCheckType
		}

		if httpCheckMethod, ok := healthCheckArg["HttpCheckMethod"]; ok {
			healthCheckMap["http_check_method"] = httpCheckMethod
		}

		if healthCheckDomain, ok := healthCheckArg["HealthCheckDomain"]; ok {
			healthCheckMap["health_check_domain"] = healthCheckDomain
		}

		if healthCheckUrl, ok := healthCheckArg["HealthCheckUrl"]; ok {
			healthCheckMap["health_check_url"] = healthCheckUrl
		}

		if healthCheckConnectPort, ok := healthCheckArg["HealthCheckConnectPort"]; ok {
			healthCheckMap["health_check_connect_port"] = formatInt(healthCheckConnectPort)
		}

		if healthCheckInterval, ok := healthCheckArg["HealthCheckInterval"]; ok {
			healthCheckMap["health_check_interval"] = formatInt(healthCheckInterval)
		}

		if healthyThreshold, ok := healthCheckArg["HealthyThreshold"]; ok {
			healthCheckMap["healthy_threshold"] = formatInt(healthyThreshold)
		}

		if unhealthyThreshold, ok := healthCheckArg["UnhealthyThreshold"]; ok {
			healthCheckMap["unhealthy_threshold"] = formatInt(unhealthyThreshold)
		}

		if healthCheckConnectTimeout, ok := healthCheckArg["HealthCheckConnectTimeout"]; ok {
			healthCheckMap["health_check_connect_timeout"] = formatInt(healthCheckConnectTimeout)
		}

		if healthCheckHttpCode, ok := healthCheckArg["HealthCheckHttpCode"]; ok {
			healthCheckMap["health_check_http_code"] = healthCheckHttpCode
		}

		healthCheckMaps = append(healthCheckMaps, healthCheckMap)

		d.Set("health_check", healthCheckMaps)
	}

	tagsMaps := object["Tags"]
	d.Set("tags", tagsToMap(tagsMaps))

	d.Set("status", object["ServerGroupStatus"])
	d.Set("connection_drain", d.Get("connection_drain_enabled"))

	return nil
}

func resourceAliCloudNlbServerGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	nlbServiceV2 := NlbServiceV2{client}
	var response map[string]interface{}
	d.Partial(true)

	if d.HasChange("tags") {
		if err := nlbServiceV2.SetResourceTags(d, "servergroup"); err != nil {
			return WrapError(err)
		}

		d.SetPartial("tags")
	}

	update := false
	updateServerGroupAttributeReq := map[string]interface{}{
		"RegionId":      client.RegionId,
		"ClientToken":   buildClientToken("UpdateServerGroupAttribute"),
		"ServerGroupId": d.Id(),
	}

	if !d.IsNewResource() && d.HasChange("server_group_name") {
		update = true
	}
	updateServerGroupAttributeReq["ServerGroupName"] = d.Get("server_group_name")

	if !d.IsNewResource() && d.HasChange("scheduler") {
		update = true
	}
	if v, ok := d.GetOk("scheduler"); ok {
		updateServerGroupAttributeReq["Scheduler"] = v
	}

	if !d.IsNewResource() && d.HasChange("connection_drain_enabled") {
		update = true

		if v, ok := d.GetOkExists("connection_drain_enabled"); ok {
			updateServerGroupAttributeReq["ConnectionDrainEnabled"] = v
		}
	}

	if !d.IsNewResource() && d.HasChange("connection_drain_timeout") {
		update = true

		if v, ok := d.GetOk("connection_drain_timeout"); ok {
			updateServerGroupAttributeReq["ConnectionDrainTimeout"] = v
		}
	}

	if !d.IsNewResource() && d.HasChange("preserve_client_ip_enabled") {
		update = true

		if v, ok := d.GetOkExists("preserve_client_ip_enabled"); ok {
			updateServerGroupAttributeReq["PreserveClientIpEnabled"] = v
		}
	}

	if !d.IsNewResource() && d.HasChange("health_check") {
		update = true
	}
	if _, ok := d.GetOk("health_check"); ok {
		healthCheckMap := make(map[string]interface{})

		if healthCheckEnabled, ok := d.GetOkExists("health_check.0.health_check_enabled"); ok {
			healthCheckMap["HealthCheckEnabled"] = healthCheckEnabled
		}

		if healthCheckMap["HealthCheckEnabled"] == true {
			if healthCheckType, ok := d.GetOk("health_check.0.health_check_type"); ok {
				healthCheckMap["HealthCheckType"] = healthCheckType
			}

			if httpCheckMethod, ok := d.GetOk("health_check.0.http_check_method"); ok {
				healthCheckMap["HttpCheckMethod"] = httpCheckMethod
			}

			if healthCheckDomain, ok := d.GetOk("health_check.0.health_check_domain"); ok {
				healthCheckMap["HealthCheckDomain"] = healthCheckDomain
			}

			if healthCheckUrl, ok := d.GetOk("health_check.0.health_check_url"); ok {
				healthCheckMap["HealthCheckUrl"] = healthCheckUrl
			}

			if healthCheckConnectPort, ok := d.GetOkExists("health_check.0.health_check_connect_port"); ok {
				healthCheckMap["HealthCheckConnectPort"] = healthCheckConnectPort
			}

			if healthCheckInterval, ok := d.GetOkExists("health_check.0.health_check_interval"); ok {
				healthCheckMap["HealthCheckInterval"] = healthCheckInterval
			}

			if healthyThreshold, ok := d.GetOkExists("health_check.0.healthy_threshold"); ok {
				healthCheckMap["HealthyThreshold"] = healthyThreshold
			}

			if unhealthyThreshold, ok := d.GetOkExists("health_check.0.unhealthy_threshold"); ok {
				healthCheckMap["UnhealthyThreshold"] = unhealthyThreshold
			}

			if healthCheckConnectTimeout, ok := d.GetOkExists("health_check.0.health_check_connect_timeout"); ok {
				healthCheckMap["HealthCheckConnectTimeout"] = healthCheckConnectTimeout
			}

			if healthCheckHttpCode, ok := d.GetOk("health_check.0.health_check_http_code"); ok {
				healthCheckMap["HealthCheckHttpCode"] = healthCheckHttpCode
			}
		}

		updateServerGroupAttributeReq["HealthCheckConfig"] = healthCheckMap
	}

	if !d.IsNewResource() && d.HasChange("connection_drain") {
		update = true

		if v, ok := d.GetOkExists("connection_drain"); ok {
			updateServerGroupAttributeReq["ConnectionDrainEnabled"] = v
		}
	}

	if update {
		action := "UpdateServerGroupAttribute"
		conn, err := client.NewNlbClient()
		if err != nil {
			return WrapError(err)
		}

		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2022-04-30"), StringPointer("AK"), nil, updateServerGroupAttributeReq, &runtime)
			if err != nil {
				if NeedRetry(err) {
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

		stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, nlbServiceV2.NlbServerGroupStateRefreshFunc(d.Id(), "ServerGroupStatus", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}

		d.SetPartial("server_group_name")
		d.SetPartial("scheduler")
		d.SetPartial("connection_drain_enabled")
		d.SetPartial("connection_drain_timeout")
		d.SetPartial("preserve_client_ip_enabled")
		d.SetPartial("health_check")
		d.SetPartial("connection_drain")
	}

	update = false
	moveResourceGroup := map[string]interface{}{
		"RegionId":     client.RegionId,
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
		conn, err := client.NewNlbClient()
		if err != nil {
			return WrapError(err)
		}

		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2022-04-30"), StringPointer("AK"), nil, moveResourceGroup, &runtime)
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
			"RegionId":      client.RegionId,
			"ClientToken":   buildClientToken("RemoveServersFromServerGroup"),
			"ServerGroupId": d.Id(),
		}

		removedList := removed.(*schema.Set).List()
		removeServersMaps := make([]map[string]interface{}, 0)

		if len(removedList) <= 200 {
			for _, servers := range removedList {
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
						if IsExpectedErrors(err, []string{"IncorrectStatus.serverGroup"}) || NeedRetry(err) {
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
		} else {
			for i := 0; i < len(removedList); i = i + 200 {
				if len(removedList[i:]) < 200 {
					for _, servers := range removedList[i:] {
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
								if IsExpectedErrors(err, []string{"IncorrectStatus.serverGroup"}) || NeedRetry(err) {
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
				} else {
					for _, servers := range removedList[i : i+200] {
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
								if IsExpectedErrors(err, []string{"IncorrectStatus.serverGroup"}) || NeedRetry(err) {
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
				}
			}
		}

		update = false
		addServersToServerGroupReq := map[string]interface{}{
			"RegionId":      client.RegionId,
			"ClientToken":   buildClientToken("AddServersToServerGroup"),
			"ServerGroupId": d.Id(),
		}

		addedList := added.(*schema.Set).List()
		addServersMaps := make([]map[string]interface{}, 0)

		if len(addedList) <= 200 {
			for _, servers := range addedList {
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
						if IsExpectedErrors(err, []string{"IncorrectStatus.serverGroup"}) || NeedRetry(err) {
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
		} else {
			for i := 0; i < len(addedList); i = i + 200 {
				if len(addedList[i:]) < 200 {
					for _, servers := range addedList[i:] {
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
								if IsExpectedErrors(err, []string{"IncorrectStatus.serverGroup"}) || NeedRetry(err) {
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
				} else {
					for _, servers := range addedList[i : i+200] {
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
								if IsExpectedErrors(err, []string{"IncorrectStatus.serverGroup"}) || NeedRetry(err) {
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
				}
			}
		}

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
					if IsExpectedErrors(err, []string{"IncorrectStatus.serverGroup"}) || NeedRetry(err) {
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

	return resourceAliCloudNlbServerGroupRead(d, meta)
}

func resourceAliCloudNlbServerGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	nlbServiceV2 := NlbServiceV2{client}
	action := "DeleteServerGroup"
	var response map[string]interface{}

	conn, err := client.NewNlbClient()
	if err != nil {
		return WrapError(err)
	}

	request := map[string]interface{}{
		"RegionId":      client.RegionId,
		"ClientToken":   buildClientToken("DeleteServerGroup"),
		"ServerGroupId": d.Id(),
	}

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2022-04-30"), StringPointer("AK"), nil, request, &runtime)
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
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	stateConf := BuildStateConf([]string{}, []string{"Succeeded"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, nlbServiceV2.DescribeAsyncNlbServerGroupStateRefreshFunc(d, response, "$.Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}
