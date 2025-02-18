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
			"address_ip_version": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"DualStack", "Ipv4"}, false),
			},
			"any_port_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"connection_drain_enabled": {
				Type:          schema.TypeBool,
				Optional:      true,
				ConflictsWith: []string{"connection_drain"},
				Computed:      true,
			},
			"connection_drain_timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: IntBetween(10, 900),
			},
			"health_check": {
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
						"health_check_url": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
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
							ValidateFunc: IntBetween(0, 300),
						},
						"health_check_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"health_check_connect_port": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: IntBetween(0, 65535),
						},
						"health_check_req": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"healthy_threshold": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: IntBetween(2, 10),
						},
						"http_check_method": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: StringInSlice([]string{"GET", "HEAD"}, false),
						},
						"health_check_exp": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"health_check_domain": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"health_check_http_code": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"health_check_type": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: StringInSlice([]string{"TCP", "HTTP", "UDP"}, false),
						},
					},
				},
			},
			"preserve_client_ip_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"protocol": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"TCP", "UDP", "TCPSSL"}, false),
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
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
				ValidateFunc: StringInSlice([]string{"Wrr", "Rr", "Qch", "Tch", "Sch", "Wlc"}, false),
			},
			"server_group_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"server_group_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"Instance", "Ip"}, false),
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

	action := "CreateServerGroup"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOk("server_group_type"); ok {
		request["ServerGroupType"] = v
	}
	request["ServerGroupName"] = d.Get("server_group_name")
	if v, ok := d.GetOk("protocol"); ok {
		request["Protocol"] = v
	}
	request["VpcId"] = d.Get("vpc_id")
	if v, ok := d.GetOkExists("any_port_enabled"); ok {
		request["AnyPortEnabled"] = v
	}
	if v, ok := d.GetOkExists("connection_drain"); ok || d.HasChange("connection_drain") {
		request["ConnectionDrainEnabled"] = v
	}

	if v, ok := d.GetOkExists("connection_drain_enabled"); ok {
		request["ConnectionDrainEnabled"] = v
	}
	if v, ok := d.GetOkExists("connection_drain_timeout"); ok && v.(int) > 0 {
		request["ConnectionDrainTimeout"] = v
	}
	if v, ok := d.GetOk("scheduler"); ok {
		request["Scheduler"] = v
	}
	if v, ok := d.GetOkExists("preserve_client_ip_enabled"); ok {
		request["PreserveClientIpEnabled"] = v
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOk("address_ip_version"); ok {
		request["AddressIPVersion"] = v
	}
	objectDataLocalMap := make(map[string]interface{})

	if v := d.Get("health_check"); !IsNil(v) {
		healthCheckEnabled1, _ := jsonpath.Get("$[0].health_check_enabled", v)
		if healthCheckEnabled1 != nil && healthCheckEnabled1 != "" {
			objectDataLocalMap["HealthCheckEnabled"] = healthCheckEnabled1
		}
		healthCheckType1, _ := jsonpath.Get("$[0].health_check_type", v)
		if healthCheckType1 != nil && healthCheckType1 != "" {
			objectDataLocalMap["HealthCheckType"] = healthCheckType1
		}
		healthCheckConnectPort1, _ := jsonpath.Get("$[0].health_check_connect_port", v)
		if healthCheckConnectPort1 != nil && healthCheckConnectPort1 != "" {
			objectDataLocalMap["HealthCheckConnectPort"] = healthCheckConnectPort1
		}
		healthyThreshold1, _ := jsonpath.Get("$[0].healthy_threshold", v)
		if healthyThreshold1 != nil && healthyThreshold1 != "" && healthyThreshold1.(int) > 0 {
			objectDataLocalMap["HealthyThreshold"] = healthyThreshold1
		}
		unhealthyThreshold1, _ := jsonpath.Get("$[0].unhealthy_threshold", v)
		if unhealthyThreshold1 != nil && unhealthyThreshold1 != "" && unhealthyThreshold1.(int) > 0 {
			objectDataLocalMap["UnhealthyThreshold"] = unhealthyThreshold1
		}
		healthCheckConnectTimeout1, _ := jsonpath.Get("$[0].health_check_connect_timeout", v)
		if healthCheckConnectTimeout1 != nil && healthCheckConnectTimeout1 != "" && healthCheckConnectTimeout1.(int) > 0 {
			objectDataLocalMap["HealthCheckConnectTimeout"] = healthCheckConnectTimeout1
		}
		healthCheckInterval1, _ := jsonpath.Get("$[0].health_check_interval", v)
		if healthCheckInterval1 != nil && healthCheckInterval1 != "" && healthCheckInterval1.(int) > 0 {
			objectDataLocalMap["HealthCheckInterval"] = healthCheckInterval1
		}
		healthCheckDomain1, _ := jsonpath.Get("$[0].health_check_domain", v)
		if healthCheckDomain1 != nil && healthCheckDomain1 != "" {
			objectDataLocalMap["HealthCheckDomain"] = healthCheckDomain1
		}
		healthCheckUrl1, _ := jsonpath.Get("$[0].health_check_url", v)
		if healthCheckUrl1 != nil && healthCheckUrl1 != "" {
			objectDataLocalMap["HealthCheckUrl"] = healthCheckUrl1
		}
		httpCheckMethod1, _ := jsonpath.Get("$[0].http_check_method", v)
		if httpCheckMethod1 != nil && httpCheckMethod1 != "" {
			objectDataLocalMap["HttpCheckMethod"] = httpCheckMethod1
		}
		healthCheckHttpCode1, _ := jsonpath.Get("$[0].health_check_http_code", v)
		if healthCheckHttpCode1 != nil && healthCheckHttpCode1 != "" {
			objectDataLocalMap["HealthCheckHttpCode"] = healthCheckHttpCode1
		}
		healthCheckReq1, _ := jsonpath.Get("$[0].health_check_req", v)
		if healthCheckReq1 != nil && healthCheckReq1 != "" {
			objectDataLocalMap["HealthCheckReq"] = healthCheckReq1
		}
		healthCheckExp1, _ := jsonpath.Get("$[0].health_check_exp", v)
		if healthCheckExp1 != nil && healthCheckExp1 != "" {
			objectDataLocalMap["HealthCheckExp"] = healthCheckExp1
		}

		request["HealthCheckConfig"] = objectDataLocalMap
	}
	if v, ok := d.GetOk("tags"); ok {
		tagsMap := ConvertTags(v.(map[string]interface{}))
		request = expandTagsToMap(request, tagsMap)
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Nlb", "2022-04-30", action, query, request, true)
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

	nlbServiceV2 := NlbServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, nlbServiceV2.NlbServerGroupStateRefreshFunc(d.Id(), "ServerGroupStatus", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudNlbServerGroupRead(d, meta)
}

func resourceAliCloudNlbServerGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	nlbServiceV2 := NlbServiceV2{client}

	objectRaw, err := nlbServiceV2.DescribeNlbServerGroup(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_nlb_server_group DescribeNlbServerGroup Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("address_ip_version", objectRaw["AddressIPVersion"])
	d.Set("any_port_enabled", objectRaw["AnyPortEnabled"])
	d.Set("connection_drain_enabled", objectRaw["ConnectionDrainEnabled"])
	d.Set("connection_drain_timeout", objectRaw["ConnectionDrainTimeout"])
	d.Set("preserve_client_ip_enabled", objectRaw["PreserveClientIpEnabled"])
	d.Set("protocol", objectRaw["Protocol"])
	d.Set("region_id", objectRaw["RegionId"])
	d.Set("resource_group_id", objectRaw["ResourceGroupId"])
	d.Set("scheduler", objectRaw["Scheduler"])
	d.Set("server_group_name", objectRaw["ServerGroupName"])
	d.Set("server_group_type", objectRaw["ServerGroupType"])
	d.Set("status", objectRaw["ServerGroupStatus"])
	d.Set("vpc_id", objectRaw["VpcId"])

	healthCheckMaps := make([]map[string]interface{}, 0)
	healthCheckMap := make(map[string]interface{})
	healthCheckRaw := make(map[string]interface{})
	if objectRaw["HealthCheck"] != nil {
		healthCheckRaw = objectRaw["HealthCheck"].(map[string]interface{})
	}
	if len(healthCheckRaw) > 0 {
		healthCheckMap["health_check_connect_port"] = healthCheckRaw["HealthCheckConnectPort"]
		healthCheckMap["health_check_connect_timeout"] = healthCheckRaw["HealthCheckConnectTimeout"]
		healthCheckMap["health_check_domain"] = healthCheckRaw["HealthCheckDomain"]
		healthCheckMap["health_check_enabled"] = healthCheckRaw["HealthCheckEnabled"]
		healthCheckMap["health_check_exp"] = healthCheckRaw["HealthCheckExp"]
		healthCheckMap["health_check_interval"] = healthCheckRaw["HealthCheckInterval"]
		healthCheckMap["health_check_req"] = healthCheckRaw["HealthCheckReq"]
		healthCheckMap["health_check_type"] = healthCheckRaw["HealthCheckType"]
		healthCheckMap["health_check_url"] = healthCheckRaw["HealthCheckUrl"]
		healthCheckMap["healthy_threshold"] = healthCheckRaw["HealthyThreshold"]
		healthCheckMap["http_check_method"] = healthCheckRaw["HttpCheckMethod"]
		healthCheckMap["unhealthy_threshold"] = healthCheckRaw["UnhealthyThreshold"]

		healthCheckHttpCodeRaw := make([]interface{}, 0)
		if healthCheckRaw["HealthCheckHttpCode"] != nil {
			healthCheckHttpCodeRaw = healthCheckRaw["HealthCheckHttpCode"].([]interface{})
		}

		healthCheckMap["health_check_http_code"] = healthCheckHttpCodeRaw
		healthCheckMaps = append(healthCheckMaps, healthCheckMap)
	}
	if err := d.Set("health_check", healthCheckMaps); err != nil {
		return err
	}
	tagsMaps := objectRaw["Tags"]
	d.Set("tags", tagsToMap(tagsMaps))

	d.Set("connection_drain", d.Get("connection_drain_enabled"))
	return nil
}

func resourceAliCloudNlbServerGroupUpdate(d *schema.ResourceData, meta interface{}) error {
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
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)
	if d.HasChange("connection_drain") {
		update = true
		request["ConnectionDrainEnabled"] = d.Get("connection_drain")
	}

	if d.HasChange("connection_drain_enabled") {
		update = true
		request["ConnectionDrainEnabled"] = d.Get("connection_drain_enabled")
	}

	if d.HasChange("connection_drain_timeout") {
		update = true
		request["ConnectionDrainTimeout"] = d.Get("connection_drain_timeout")
	}

	if d.HasChange("scheduler") {
		update = true
		request["Scheduler"] = d.Get("scheduler")
	}

	if d.HasChange("preserve_client_ip_enabled") {
		update = true
		request["PreserveClientIpEnabled"] = d.Get("preserve_client_ip_enabled")
	}

	if d.HasChange("health_check") {
		update = true
		objectDataLocalMap := make(map[string]interface{})

		if v := d.Get("health_check"); v != nil {
			healthCheckEnabled1, _ := jsonpath.Get("$[0].health_check_enabled", v)
			if healthCheckEnabled1 != nil && (d.HasChange("health_check.0.health_check_enabled") || healthCheckEnabled1 != "") {
				objectDataLocalMap["HealthCheckEnabled"] = healthCheckEnabled1
			}
			healthCheckType1, _ := jsonpath.Get("$[0].health_check_type", v)
			if healthCheckType1 != nil && (d.HasChange("health_check.0.health_check_type") || healthCheckType1 != "") {
				objectDataLocalMap["HealthCheckType"] = healthCheckType1
			}
			healthCheckConnectPort1, _ := jsonpath.Get("$[0].health_check_connect_port", v)
			if healthCheckConnectPort1 != nil && (d.HasChange("health_check.0.health_check_connect_port") || healthCheckConnectPort1 != "") {
				objectDataLocalMap["HealthCheckConnectPort"] = healthCheckConnectPort1
			}
			healthyThreshold1, _ := jsonpath.Get("$[0].healthy_threshold", v)
			if healthyThreshold1 != nil && (d.HasChange("health_check.0.healthy_threshold") || healthyThreshold1 != "") && healthyThreshold1.(int) > 0 {
				objectDataLocalMap["HealthyThreshold"] = healthyThreshold1
			}
			unhealthyThreshold1, _ := jsonpath.Get("$[0].unhealthy_threshold", v)
			if unhealthyThreshold1 != nil && (d.HasChange("health_check.0.unhealthy_threshold") || unhealthyThreshold1 != "") && unhealthyThreshold1.(int) > 0 {
				objectDataLocalMap["UnhealthyThreshold"] = unhealthyThreshold1
			}
			healthCheckConnectTimeout1, _ := jsonpath.Get("$[0].health_check_connect_timeout", v)
			if healthCheckConnectTimeout1 != nil && (d.HasChange("health_check.0.health_check_connect_timeout") || healthCheckConnectTimeout1 != "") && healthCheckConnectTimeout1.(int) > 0 {
				objectDataLocalMap["HealthCheckConnectTimeout"] = healthCheckConnectTimeout1
			}
			healthCheckInterval1, _ := jsonpath.Get("$[0].health_check_interval", v)
			if healthCheckInterval1 != nil && (d.HasChange("health_check.0.health_check_interval") || healthCheckInterval1 != "") && healthCheckInterval1.(int) > 0 {
				objectDataLocalMap["HealthCheckInterval"] = healthCheckInterval1
			}
			healthCheckDomain1, _ := jsonpath.Get("$[0].health_check_domain", v)
			if healthCheckDomain1 != nil && (d.HasChange("health_check.0.health_check_domain") || healthCheckDomain1 != "") {
				objectDataLocalMap["HealthCheckDomain"] = healthCheckDomain1
			}
			healthCheckUrl1, _ := jsonpath.Get("$[0].health_check_url", v)
			if healthCheckUrl1 != nil && (d.HasChange("health_check.0.health_check_url") || healthCheckUrl1 != "") {
				objectDataLocalMap["HealthCheckUrl"] = healthCheckUrl1
			}
			httpCheckMethod1, _ := jsonpath.Get("$[0].http_check_method", v)
			if httpCheckMethod1 != nil && (d.HasChange("health_check.0.http_check_method") || httpCheckMethod1 != "") {
				objectDataLocalMap["HttpCheckMethod"] = httpCheckMethod1
			}
			healthCheckHttpCode1, _ := jsonpath.Get("$[0].health_check_http_code", d.Get("health_check"))
			if healthCheckHttpCode1 != nil && (d.HasChange("health_check.0.health_check_http_code") || healthCheckHttpCode1 != "") {
				objectDataLocalMap["HealthCheckHttpCode"] = healthCheckHttpCode1
			}
			healthCheckReq1, _ := jsonpath.Get("$[0].health_check_req", v)
			if healthCheckReq1 != nil && (d.HasChange("health_check.0.health_check_req") || healthCheckReq1 != "") {
				objectDataLocalMap["HealthCheckReq"] = healthCheckReq1
			}
			healthCheckExp1, _ := jsonpath.Get("$[0].health_check_exp", v)
			if healthCheckExp1 != nil && (d.HasChange("health_check.0.health_check_exp") || healthCheckExp1 != "") {
				objectDataLocalMap["HealthCheckExp"] = healthCheckExp1
			}

			request["HealthCheckConfig"] = objectDataLocalMap
		}
	}

	if d.HasChange("server_group_name") {
		update = true
	}
	request["ServerGroupName"] = d.Get("server_group_name")
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Nlb", "2022-04-30", action, query, request, true)
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
		nlbServiceV2 := NlbServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, nlbServiceV2.NlbServerGroupStateRefreshFunc(d.Id(), "ServerGroupStatus", []string{}))
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
	if _, ok := d.GetOk("resource_group_id"); ok && d.HasChange("resource_group_id") {
		update = true
	}
	request["NewResourceGroupId"] = d.Get("resource_group_id")
	request["ResourceType"] = "servergroup"
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Nlb", "2022-04-30", action, query, request, true)
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

	if d.HasChange("tags") {
		nlbServiceV2 := NlbServiceV2{client}
		if err := nlbServiceV2.SetResourceTags(d, "servergroup"); err != nil {
			return WrapError(err)
		}
	}
	d.Partial(false)
	return resourceAliCloudNlbServerGroupRead(d, meta)
}

func resourceAliCloudNlbServerGroupDelete(d *schema.ResourceData, meta interface{}) error {

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

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Nlb", "2022-04-30", action, query, request, true)
		request["ClientToken"] = buildClientToken(action)

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

	nlbServiceV2 := NlbServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"Succeeded"}, d.Timeout(schema.TimeoutDelete), 5*time.Second, nlbServiceV2.DescribeAsyncNlbServerGroupStateRefreshFunc(d, response, "$.Status", []string{}))
	if jobDetail, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id(), jobDetail)
	}

	return nil
}
