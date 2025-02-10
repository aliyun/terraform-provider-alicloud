// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"bytes"
	"fmt"
	"hash/crc32"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudGwlbServerGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudGwlbServerGroupCreate,
		Read:   resourceAliCloudGwlbServerGroupRead,
		Update: resourceAliCloudGwlbServerGroupUpdate,
		Delete: resourceAliCloudGwlbServerGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
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
			"dry_run": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"health_check_config": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
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
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"unhealthy_threshold": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: IntBetween(2, 10),
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
						"health_check_connect_timeout": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"health_check_protocol": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: StringInSlice([]string{"TCP", "HTTP"}, false),
						},
						"health_check_domain": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"health_check_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"health_check_http_code": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"protocol": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"GENEVE"}, false),
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
				ValidateFunc: StringInSlice([]string{"2TCH", "3TCH", "5TCH"}, false),
			},
			"server_group_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"server_group_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"Ip", "Instance"}, false),
			},
			"servers": {
				Set: func(v interface{}) int {
					var buf bytes.Buffer
					server := v.(map[string]interface{})
					if v, ok := server["server_id"]; ok {
						buf.WriteString(fmt.Sprintf("%s-", v.(string)))
					}
					if v, ok := server["server_ip"]; ok && server["server_type"] == "Eni" {
						buf.WriteString(fmt.Sprintf("%s-", v.(string)))
					}
					if v, ok := server["server_port"]; ok {
						buf.WriteString(fmt.Sprintf("%d", v.(int)))
					}
					return int(crc32.ChecksumIEEE([]byte(buf.String())))
				},
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"server_type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: StringInSlice([]string{"Ip", "Eci", "Eni", "Ecs"}, false),
						},
						"server_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"server_ip": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"server_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsSchema(),
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudGwlbServerGroupCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateServerGroup"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOk("health_check_config"); ok {
		jsonPathResult, err := jsonpath.Get("$[0].health_check_connect_timeout", v)
		if err == nil && jsonPathResult != "" {
			request["HealthCheckConfig.HealthCheckConnectTimeout"] = jsonPathResult
		}
	}
	if v, ok := d.GetOk("health_check_config"); ok {
		jsonPathResult1, err := jsonpath.Get("$[0].health_check_connect_port", v)
		if err == nil && jsonPathResult1 != "" {
			request["HealthCheckConfig.HealthCheckConnectPort"] = jsonPathResult1
		}
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOk("tags"); ok {
		tagsMap := ConvertTags(v.(map[string]interface{}))
		request["Tags"] = tagsMap
	}

	if v, ok := d.GetOk("protocol"); ok {
		request["Protocol"] = v
	}
	if v, ok := d.GetOk("health_check_config"); ok {
		jsonPathResult4, err := jsonpath.Get("$[0].health_check_path", v)
		if err == nil && jsonPathResult4 != "" {
			request["HealthCheckConfig.HealthCheckPath"] = jsonPathResult4
		}
	}
	if v, ok := d.GetOk("scheduler"); ok {
		request["Scheduler"] = v
	}
	if v, ok := d.GetOk("health_check_config"); ok {
		jsonPathResult6, err := jsonpath.Get("$[0].health_check_enabled", v)
		if err == nil && jsonPathResult6 != "" {
			request["HealthCheckConfig.HealthCheckEnabled"] = jsonPathResult6
		}
	}
	objectDataLocalMap := make(map[string]interface{})

	if v := d.Get("health_check_config"); !IsNil(v) {
		healthCheckConnectTimeout1, _ := jsonpath.Get("$[0].health_check_connect_timeout", d.Get("health_check_config"))
		if healthCheckConnectTimeout1 != nil && healthCheckConnectTimeout1 != "" {
			objectDataLocalMap["HealthCheckConnectTimeout"] = healthCheckConnectTimeout1
		}
		healthCheckConnectPort1, _ := jsonpath.Get("$[0].health_check_connect_port", d.Get("health_check_config"))
		if healthCheckConnectPort1 != nil && healthCheckConnectPort1 != "" {
			objectDataLocalMap["HealthCheckConnectPort"] = healthCheckConnectPort1
		}
		healthCheckPath1, _ := jsonpath.Get("$[0].health_check_path", d.Get("health_check_config"))
		if healthCheckPath1 != nil && healthCheckPath1 != "" {
			objectDataLocalMap["HealthCheckPath"] = healthCheckPath1
		}
		healthCheckEnabled1, _ := jsonpath.Get("$[0].health_check_enabled", d.Get("health_check_config"))
		if healthCheckEnabled1 != nil && healthCheckEnabled1 != "" {
			objectDataLocalMap["HealthCheckEnabled"] = healthCheckEnabled1
		}
		healthCheckHttpCode1, _ := jsonpath.Get("$[0].health_check_http_code", v)
		if healthCheckHttpCode1 != nil && healthCheckHttpCode1 != "" {
			objectDataLocalMap["HealthCheckHttpCode"] = healthCheckHttpCode1.(*schema.Set).List()
		}
		unhealthyThreshold1, _ := jsonpath.Get("$[0].unhealthy_threshold", d.Get("health_check_config"))
		if unhealthyThreshold1 != nil && unhealthyThreshold1 != "" && unhealthyThreshold1.(int) > 0 {
			objectDataLocalMap["UnhealthyThreshold"] = unhealthyThreshold1
		}
		healthCheckInterval1, _ := jsonpath.Get("$[0].health_check_interval", d.Get("health_check_config"))
		if healthCheckInterval1 != nil && healthCheckInterval1 != "" && healthCheckInterval1.(int) > 0 {
			objectDataLocalMap["HealthCheckInterval"] = healthCheckInterval1
		}
		healthCheckProtocol1, _ := jsonpath.Get("$[0].health_check_protocol", d.Get("health_check_config"))
		if healthCheckProtocol1 != nil && healthCheckProtocol1 != "" {
			objectDataLocalMap["HealthCheckProtocol"] = healthCheckProtocol1
		}
		healthCheckDomain1, _ := jsonpath.Get("$[0].health_check_domain", d.Get("health_check_config"))
		if healthCheckDomain1 != nil && healthCheckDomain1 != "" {
			objectDataLocalMap["HealthCheckDomain"] = healthCheckDomain1
		}
		healthyThreshold1, _ := jsonpath.Get("$[0].healthy_threshold", d.Get("health_check_config"))
		if healthyThreshold1 != nil && healthyThreshold1 != "" && healthyThreshold1.(int) > 0 {
			objectDataLocalMap["HealthyThreshold"] = healthyThreshold1
		}

		request["HealthCheckConfig"] = objectDataLocalMap
	}

	if v, ok := d.GetOk("server_group_name"); ok {
		request["ServerGroupName"] = v
	}
	if v, ok := d.GetOk("health_check_config"); ok {
		jsonPathResult8, err := jsonpath.Get("$[0].unhealthy_threshold", v)
		if err == nil && jsonPathResult8 != "" {
			request["HealthCheckConfig.UnhealthyThreshold"] = jsonPathResult8
		}
	}
	request["VpcId"] = d.Get("vpc_id")
	if v, ok := d.GetOk("server_group_type"); ok {
		request["ServerGroupType"] = v
	}
	if v, ok := d.GetOk("health_check_config"); ok {
		jsonPathResult11, err := jsonpath.Get("$[0].health_check_interval", v)
		if err == nil && jsonPathResult11 != "" {
			request["HealthCheckConfig.HealthCheckInterval"] = jsonPathResult11
		}
	}
	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	if v, ok := d.GetOk("health_check_config"); ok {
		jsonPathResult13, err := jsonpath.Get("$[0].health_check_protocol", v)
		if err == nil && jsonPathResult13 != "" {
			request["HealthCheckConfig.HealthCheckProtocol"] = jsonPathResult13
		}
	}
	if v, ok := d.GetOk("connection_drain_config"); ok {
		jsonPathResult14, err := jsonpath.Get("$[0].connection_drain_enabled", v)
		if err == nil && jsonPathResult14 != "" {
			request["ConnectionDrainConfig.ConnectionDrainEnabled"] = jsonPathResult14
		}
	}
	if v, ok := d.GetOk("connection_drain_config"); ok {
		jsonPathResult15, err := jsonpath.Get("$[0].connection_drain_timeout", v)
		if err == nil && jsonPathResult15 != "" {
			request["ConnectionDrainConfig.ConnectionDrainTimeout"] = jsonPathResult15
		}
	}
	if v, ok := d.GetOk("health_check_config"); ok {
		jsonPathResult16, err := jsonpath.Get("$[0].health_check_domain", v)
		if err == nil && jsonPathResult16 != "" {
			request["HealthCheckConfig.HealthCheckDomain"] = jsonPathResult16
		}
	}
	if v, ok := d.GetOk("health_check_config"); ok {
		jsonPathResult17, err := jsonpath.Get("$[0].healthy_threshold", v)
		if err == nil && jsonPathResult17 != "" {
			request["HealthCheckConfig.HealthyThreshold"] = jsonPathResult17
		}
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Gwlb", "2024-04-15", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_gwlb_server_group", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["ServerGroupId"]))

	gwlbServiceV2 := GwlbServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, gwlbServiceV2.GwlbServerGroupStateRefreshFunc(d.Id(), "ServerGroupStatus", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudGwlbServerGroupUpdate(d, meta)
}

func resourceAliCloudGwlbServerGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gwlbServiceV2 := GwlbServiceV2{client}

	objectRaw, err := gwlbServiceV2.DescribeGwlbServerGroup(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_gwlb_server_group DescribeGwlbServerGroup Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	if objectRaw["CreateTime"] != nil {
		d.Set("create_time", objectRaw["CreateTime"])
	}
	if objectRaw["Protocol"] != nil {
		d.Set("protocol", objectRaw["Protocol"])
	}
	if objectRaw["ResourceGroupId"] != nil {
		d.Set("resource_group_id", objectRaw["ResourceGroupId"])
	}
	if objectRaw["Scheduler"] != nil {
		d.Set("scheduler", objectRaw["Scheduler"])
	}
	if objectRaw["ServerGroupName"] != nil {
		d.Set("server_group_name", objectRaw["ServerGroupName"])
	}
	if objectRaw["ServerGroupType"] != nil {
		d.Set("server_group_type", objectRaw["ServerGroupType"])
	}
	if objectRaw["ServerGroupStatus"] != nil {
		d.Set("status", objectRaw["ServerGroupStatus"])
	}
	if objectRaw["VpcId"] != nil {
		d.Set("vpc_id", objectRaw["VpcId"])
	}

	connectionDrainConfigMaps := make([]map[string]interface{}, 0)
	connectionDrainConfigMap := make(map[string]interface{})
	connectionDrainConfig1Raw := make(map[string]interface{})
	if objectRaw["ConnectionDrainConfig"] != nil {
		connectionDrainConfig1Raw = objectRaw["ConnectionDrainConfig"].(map[string]interface{})
	}
	if len(connectionDrainConfig1Raw) > 0 {
		connectionDrainConfigMap["connection_drain_enabled"] = connectionDrainConfig1Raw["ConnectionDrainEnabled"]
		connectionDrainConfigMap["connection_drain_timeout"] = connectionDrainConfig1Raw["ConnectionDrainTimeout"]

		connectionDrainConfigMaps = append(connectionDrainConfigMaps, connectionDrainConfigMap)
	}
	if objectRaw["ConnectionDrainConfig"] != nil {
		if err := d.Set("connection_drain_config", connectionDrainConfigMaps); err != nil {
			return err
		}
	}
	healthCheckConfigMaps := make([]map[string]interface{}, 0)
	healthCheckConfigMap := make(map[string]interface{})
	healthCheckConfig1Raw := make(map[string]interface{})
	if objectRaw["HealthCheckConfig"] != nil {
		healthCheckConfig1Raw = objectRaw["HealthCheckConfig"].(map[string]interface{})
	}
	if len(healthCheckConfig1Raw) > 0 {
		healthCheckConfigMap["health_check_connect_port"] = healthCheckConfig1Raw["HealthCheckConnectPort"]
		healthCheckConfigMap["health_check_connect_timeout"] = healthCheckConfig1Raw["HealthCheckConnectTimeout"]
		healthCheckConfigMap["health_check_domain"] = healthCheckConfig1Raw["HealthCheckDomain"]
		healthCheckConfigMap["health_check_enabled"] = healthCheckConfig1Raw["HealthCheckEnabled"]
		healthCheckConfigMap["health_check_interval"] = healthCheckConfig1Raw["HealthCheckInterval"]
		healthCheckConfigMap["health_check_path"] = healthCheckConfig1Raw["HealthCheckPath"]
		healthCheckConfigMap["health_check_protocol"] = healthCheckConfig1Raw["HealthCheckProtocol"]
		healthCheckConfigMap["healthy_threshold"] = healthCheckConfig1Raw["HealthyThreshold"]
		healthCheckConfigMap["unhealthy_threshold"] = healthCheckConfig1Raw["UnhealthyThreshold"]

		healthCheckHttpCode1Raw := make([]interface{}, 0)
		if healthCheckConfig1Raw["HealthCheckHttpCode"] != nil {
			healthCheckHttpCode1Raw = healthCheckConfig1Raw["HealthCheckHttpCode"].([]interface{})
		}

		healthCheckConfigMap["health_check_http_code"] = healthCheckHttpCode1Raw
		healthCheckConfigMaps = append(healthCheckConfigMaps, healthCheckConfigMap)
	}
	if objectRaw["HealthCheckConfig"] != nil {
		if err := d.Set("health_check_config", healthCheckConfigMaps); err != nil {
			return err
		}
	}
	tagsMaps := objectRaw["Tags"]
	d.Set("tags", tagsToMap(tagsMaps))

	objectRaw, err = gwlbServiceV2.DescribeListServerGroupServers(d.Id())
	if err != nil {
		return WrapError(err)
	}

	servers1Raw, _ := jsonpath.Get("$.Servers", objectRaw)

	serversMaps := make([]map[string]interface{}, 0)
	if servers1Raw != nil {
		for _, serversChild1Raw := range servers1Raw.([]interface{}) {
			serversMap := make(map[string]interface{})
			serversChild1Raw := serversChild1Raw.(map[string]interface{})
			serversMap["port"] = 6081
			serversMap["server_group_id"] = serversChild1Raw["ServerGroupId"]
			serversMap["server_id"] = serversChild1Raw["ServerId"]
			serversMap["server_ip"] = serversChild1Raw["ServerIp"]
			serversMap["server_type"] = serversChild1Raw["ServerType"]
			serversMap["status"] = serversChild1Raw["Status"]

			serversMaps = append(serversMaps, serversMap)
		}
	}
	if objectRaw["Servers"] != nil {
		if err := d.Set("servers", serversMaps); err != nil {
			return err
		}
	}

	return nil
}

func resourceAliCloudGwlbServerGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	d.Partial(true)
	action := "UpdateServerGroupAttribute"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ServerGroupId"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)
	if !d.IsNewResource() && d.HasChange("health_check_config.0.health_check_connect_timeout") {
		update = true
		jsonPathResult, err := jsonpath.Get("$[0].health_check_connect_timeout", d.Get("health_check_config"))
		if err == nil {
			request["HealthCheckConfig.HealthCheckConnectTimeout"] = jsonPathResult
		}
	}

	if !d.IsNewResource() && d.HasChange("health_check_config.0.health_check_connect_port") {
		update = true
		jsonPathResult1, err := jsonpath.Get("$[0].health_check_connect_port", d.Get("health_check_config"))
		if err == nil {
			request["HealthCheckConfig.HealthCheckConnectPort"] = jsonPathResult1
		}
	}

	if !d.IsNewResource() && d.HasChange("health_check_config.0.health_check_path") {
		update = true
		jsonPathResult2, err := jsonpath.Get("$[0].health_check_path", d.Get("health_check_config"))
		if err == nil {
			request["HealthCheckConfig.HealthCheckPath"] = jsonPathResult2
		}
	}

	if !d.IsNewResource() && d.HasChange("scheduler") {
		update = true
		request["Scheduler"] = d.Get("scheduler")
	}

	if !d.IsNewResource() && d.HasChange("health_check_config.0.health_check_enabled") {
		update = true
		jsonPathResult4, err := jsonpath.Get("$[0].health_check_enabled", d.Get("health_check_config"))
		if err == nil {
			request["HealthCheckConfig.HealthCheckEnabled"] = jsonPathResult4
		}
	}

	if !d.IsNewResource() && d.HasChange("health_check_config") {
		update = true
		objectDataLocalMap := make(map[string]interface{})

		if v := d.Get("health_check_config"); !IsNil(v) {
			healthCheckConnectTimeout1, _ := jsonpath.Get("$[0].health_check_connect_timeout", v)
			if healthCheckConnectTimeout1 != nil && (d.HasChange("health_check_config.0.health_check_connect_timeout") || healthCheckConnectTimeout1 != "") {
				objectDataLocalMap["HealthCheckConnectTimeout"] = healthCheckConnectTimeout1
			}
			healthCheckConnectPort1, _ := jsonpath.Get("$[0].health_check_connect_port", v)
			if healthCheckConnectPort1 != nil && (d.HasChange("health_check_config.0.health_check_connect_port") || healthCheckConnectPort1 != "") {
				objectDataLocalMap["HealthCheckConnectPort"] = healthCheckConnectPort1
			}
			healthCheckPath1, _ := jsonpath.Get("$[0].health_check_path", v)
			if healthCheckPath1 != nil && (d.HasChange("health_check_config.0.health_check_path") || healthCheckPath1 != "") {
				objectDataLocalMap["HealthCheckPath"] = healthCheckPath1
			}
			healthCheckEnabled1, _ := jsonpath.Get("$[0].health_check_enabled", v)
			if healthCheckEnabled1 != nil && (d.HasChange("health_check_config.0.health_check_enabled") || healthCheckEnabled1 != "") {
				objectDataLocalMap["HealthCheckEnabled"] = healthCheckEnabled1
			}
			healthCheckHttpCode1, _ := jsonpath.Get("$[0].health_check_http_code", d.Get("health_check_config"))
			if healthCheckHttpCode1 != nil && (d.HasChange("health_check_config.0.health_check_http_code") || healthCheckHttpCode1 != "") {
				objectDataLocalMap["HealthCheckHttpCode"] = healthCheckHttpCode1.(*schema.Set).List()
			}
			unhealthyThreshold1, _ := jsonpath.Get("$[0].unhealthy_threshold", v)
			if unhealthyThreshold1 != nil && (d.HasChange("health_check_config.0.unhealthy_threshold") || unhealthyThreshold1 != "") && unhealthyThreshold1.(int) > 0 {
				objectDataLocalMap["UnhealthyThreshold"] = unhealthyThreshold1
			}
			healthCheckInterval1, _ := jsonpath.Get("$[0].health_check_interval", v)
			if healthCheckInterval1 != nil && (d.HasChange("health_check_config.0.health_check_interval") || healthCheckInterval1 != "") && healthCheckInterval1.(int) > 0 {
				objectDataLocalMap["HealthCheckInterval"] = healthCheckInterval1
			}
			healthCheckProtocol1, _ := jsonpath.Get("$[0].health_check_protocol", v)
			if healthCheckProtocol1 != nil && (d.HasChange("health_check_config.0.health_check_protocol") || healthCheckProtocol1 != "") {
				objectDataLocalMap["HealthCheckProtocol"] = healthCheckProtocol1
			}
			healthCheckDomain1, _ := jsonpath.Get("$[0].health_check_domain", v)
			if healthCheckDomain1 != nil && (d.HasChange("health_check_config.0.health_check_domain") || healthCheckDomain1 != "") {
				objectDataLocalMap["HealthCheckDomain"] = healthCheckDomain1
			}
			healthyThreshold1, _ := jsonpath.Get("$[0].healthy_threshold", v)
			if healthyThreshold1 != nil && (d.HasChange("health_check_config.0.healthy_threshold") || healthyThreshold1 != "") && healthyThreshold1.(int) > 0 {
				objectDataLocalMap["HealthyThreshold"] = healthyThreshold1
			}

			request["HealthCheckConfig"] = objectDataLocalMap
		}
	}

	if !d.IsNewResource() && d.HasChange("server_group_name") {
		update = true
		request["ServerGroupName"] = d.Get("server_group_name")
	}

	if !d.IsNewResource() && d.HasChange("health_check_config.0.unhealthy_threshold") {
		update = true
		jsonPathResult6, err := jsonpath.Get("$[0].unhealthy_threshold", d.Get("health_check_config"))
		if err == nil && jsonPathResult6.(int) > 0 {
			request["HealthCheckConfig.UnhealthyThreshold"] = jsonPathResult6
		}
	}

	if !d.IsNewResource() && d.HasChange("health_check_config.0.health_check_interval") {
		update = true
		jsonPathResult7, err := jsonpath.Get("$[0].health_check_interval", d.Get("health_check_config"))
		if err == nil && jsonPathResult7.(int) > 0 {
			request["HealthCheckConfig.HealthCheckInterval"] = jsonPathResult7
		}
	}

	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	if !d.IsNewResource() && d.HasChange("health_check_config.0.health_check_protocol") {
		update = true
		jsonPathResult9, err := jsonpath.Get("$[0].health_check_protocol", d.Get("health_check_config"))
		if err == nil {
			request["HealthCheckConfig.HealthCheckProtocol"] = jsonPathResult9
		}
	}

	if !d.IsNewResource() && d.HasChange("connection_drain_config.0.connection_drain_enabled") {
		update = true
		jsonPathResult10, err := jsonpath.Get("$[0].connection_drain_enabled", d.Get("connection_drain_config"))
		if err == nil {
			request["ConnectionDrainConfig.ConnectionDrainEnabled"] = jsonPathResult10
		}
	}

	if !d.IsNewResource() && d.HasChange("connection_drain_config.0.connection_drain_timeout") {
		update = true
		jsonPathResult11, err := jsonpath.Get("$[0].connection_drain_timeout", d.Get("connection_drain_config"))
		if err == nil {
			request["ConnectionDrainConfig.ConnectionDrainTimeout"] = jsonPathResult11
		}
	}

	if !d.IsNewResource() && d.HasChange("health_check_config.0.health_check_domain") {
		update = true
		jsonPathResult12, err := jsonpath.Get("$[0].health_check_domain", d.Get("health_check_config"))
		if err == nil {
			request["HealthCheckConfig.HealthCheckDomain"] = jsonPathResult12
		}
	}

	if !d.IsNewResource() && d.HasChange("health_check_config.0.healthy_threshold") {
		update = true
		jsonPathResult13, err := jsonpath.Get("$[0].healthy_threshold", d.Get("health_check_config"))
		if err == nil && jsonPathResult13.(int) > 0 {
			request["HealthCheckConfig.HealthyThreshold"] = jsonPathResult13
		}
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Gwlb", "2024-04-15", action, query, request, true)
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
		gwlbServiceV2 := GwlbServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, gwlbServiceV2.GwlbServerGroupStateRefreshFunc(d.Id(), "ServerGroupStatus", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	update = false
	action = "MoveResourceGroup"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ResourceId"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)
	if _, ok := d.GetOk("resource_group_id"); ok && !d.IsNewResource() && d.HasChange("resource_group_id") {
		update = true
	}
	request["NewResourceGroupId"] = d.Get("resource_group_id")
	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	request["ResourceType"] = "servergroup"
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Gwlb", "2024-04-15", action, query, request, true)
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
			request["RegionId"] = client.RegionId
			request["ClientToken"] = buildClientToken(action)
			localData := removed.List()
			serversMapsArray := make([]interface{}, 0)
			for _, dataLoop := range localData {
				dataLoopTmp := dataLoop.(map[string]interface{})
				dataLoopMap := make(map[string]interface{})
				dataLoopMap["Port"] = 6081
				dataLoopMap["ServerId"] = dataLoopTmp["server_id"]
				dataLoopMap["ServerIp"] = dataLoopTmp["server_ip"]
				dataLoopMap["ServerType"] = dataLoopTmp["server_type"]
				serversMapsArray = append(serversMapsArray, dataLoopMap)
			}
			request["Servers"] = serversMapsArray

			if v, ok := d.GetOkExists("dry_run"); ok {
				request["DryRun"] = v
			}
			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = client.RpcPost("Gwlb", "2024-04-15", action, query, request, true)
				if err != nil {
					if IsExpectedErrors(err, []string{"IncorrectStatus.ServerGroup"}) || NeedRetry(err) {
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
			gwlbServiceV2 := GwlbServiceV2{client}
			stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutUpdate), 20*time.Second, gwlbServiceV2.GwlbServerGroupStateRefreshFunc(d.Id(), "ServerGroupStatus", []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}

		}

		if added.Len() > 0 {
			action := "AddServersToServerGroup"
			request = make(map[string]interface{})
			query = make(map[string]interface{})
			request["ServerGroupId"] = d.Id()
			request["RegionId"] = client.RegionId
			request["ClientToken"] = buildClientToken(action)
			localData := added.List()
			serversMapsArray := make([]interface{}, 0)
			for _, dataLoop := range localData {
				dataLoopTmp := dataLoop.(map[string]interface{})
				dataLoopMap := make(map[string]interface{})
				dataLoopMap["Port"] = 6081
				dataLoopMap["ServerId"] = dataLoopTmp["server_id"]
				dataLoopMap["ServerIp"] = dataLoopTmp["server_ip"]
				dataLoopMap["ServerType"] = dataLoopTmp["server_type"]
				serversMapsArray = append(serversMapsArray, dataLoopMap)
			}
			request["Servers"] = serversMapsArray

			if v, ok := d.GetOkExists("dry_run"); ok {
				request["DryRun"] = v
			}
			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = client.RpcPost("Gwlb", "2024-04-15", action, query, request, true)
				if err != nil {
					if IsExpectedErrors(err, []string{"IncorrectStatus.ServerGroup"}) || NeedRetry(err) {
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
			gwlbServiceV2 := GwlbServiceV2{client}
			stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutUpdate), 20*time.Second, gwlbServiceV2.GwlbServerGroupStateRefreshFunc(d.Id(), "ServerGroupStatus", []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}

		}

	}
	if d.HasChange("tags") {
		gwlbServiceV2 := GwlbServiceV2{client}
		if err := gwlbServiceV2.SetResourceTags(d, "servergroup"); err != nil {
			return WrapError(err)
		}
	}
	d.Partial(false)
	return resourceAliCloudGwlbServerGroupRead(d, meta)
}

func resourceAliCloudGwlbServerGroupDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteServerGroup"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["ServerGroupId"] = d.Id()
	request["RegionId"] = client.RegionId

	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Gwlb", "2024-04-15", action, query, request, true)
		request["ClientToken"] = buildClientToken(action)

		if err != nil {
			if IsExpectedErrors(err, []string{"IncorrectStatus.ServerGroup"}) || NeedRetry(err) {
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

	gwlbServiceV2 := GwlbServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"[]"}, d.Timeout(schema.TimeoutDelete), 5*time.Second, gwlbServiceV2.DescribeAsyncGwlbServerGroupStateRefreshFunc(d, response, "$.Servers[*]", []string{}))
	if jobDetail, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id(), jobDetail)
	}

	return nil
}
