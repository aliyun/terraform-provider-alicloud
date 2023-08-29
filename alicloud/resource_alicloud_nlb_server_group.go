// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/PaesslerAG/jsonpath"
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
				Computed:      true,
				ConflictsWith: []string{"connection_drain"},
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
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"health_check_interval": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: IntBetween(1, 5000),
						},
						"health_check_connect_port": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: IntBetween(0, 65535),
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
						"healthy_threshold": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: IntBetween(2, 10),
						},
						"http_check_method": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: StringInSlice([]string{"GET", "HEAD"}, false),
						},
						"health_check_connect_timeout": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: IntBetween(1, 300),
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
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"health_check_type": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: StringInSlice([]string{"TCP", "HTTP"}, false),
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
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"scheduler": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"Wrr", "Rr", "Sch", "Tch", "Qch"}, false),
			},
			"server_group_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: StringMatch(regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9-._]{1,127}$`), "The name must be 2 to 128 characters in length, and can contain letters, digits, periods (.), underscores (_), and hyphens (-). The name must start with a letter."),
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
				Deprecated: "Field 'connection_drain' has been deprecated since provider version 1.210.0. New field 'connection_drain_enabled' instead.",
			},
		},
	}
}

func resourceAliCloudNlbServerGroupCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateServerGroup"
	var request map[string]interface{}
	var response map[string]interface{}
	conn, err := client.NewNlbClient()
	if err != nil {
		return WrapError(err)
	}
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
	if v, ok := d.GetOkExists("connection_drain"); ok {
		request["ConnectionDrainEnabled"] = v
	}
	if v, ok := d.GetOkExists("connection_drain_enabled"); ok {
		request["ConnectionDrainEnabled"] = v
	}
	if v, ok := d.GetOk("connection_drain_timeout"); ok {
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
		nodeNative, _ := jsonpath.Get("$[0].health_check_enabled", v)
		if nodeNative != "" {
			objectDataLocalMap["HealthCheckEnabled"] = nodeNative
		}
		if nodeNative == true {
			nodeNative1, _ := jsonpath.Get("$[0].health_check_type", v)
			if nodeNative1 != "" {
				objectDataLocalMap["HealthCheckType"] = nodeNative1
			}
			nodeNative2, _ := jsonpath.Get("$[0].health_check_connect_port", v)
			if nodeNative2 != "" {
				objectDataLocalMap["HealthCheckConnectPort"] = nodeNative2
			}
			nodeNative3, _ := jsonpath.Get("$[0].healthy_threshold", v)
			if nodeNative3 != "" {
				objectDataLocalMap["HealthyThreshold"] = nodeNative3
			}
			nodeNative4, _ := jsonpath.Get("$[0].unhealthy_threshold", v)
			if nodeNative4 != "" {
				objectDataLocalMap["UnhealthyThreshold"] = nodeNative4
			}
			nodeNative5, _ := jsonpath.Get("$[0].health_check_connect_timeout", v)
			if nodeNative5 != "" {
				objectDataLocalMap["HealthCheckConnectTimeout"] = nodeNative5
			}
			nodeNative6, _ := jsonpath.Get("$[0].health_check_interval", v)
			if nodeNative6 != "" {
				objectDataLocalMap["HealthCheckInterval"] = nodeNative6
			}
			nodeNative7, _ := jsonpath.Get("$[0].health_check_domain", v)
			if nodeNative7 != "" {
				objectDataLocalMap["HealthCheckDomain"] = nodeNative7
			}
			nodeNative8, _ := jsonpath.Get("$[0].health_check_url", v)
			if nodeNative8 != "" {
				objectDataLocalMap["HealthCheckUrl"] = nodeNative8
			}
			nodeNative9, _ := jsonpath.Get("$[0].http_check_method", v)
			if nodeNative9 != "" {
				objectDataLocalMap["HttpCheckMethod"] = nodeNative9
			}
			nodeNative10, _ := jsonpath.Get("$[0].health_check_http_code", v)
			if nodeNative10 != "" {
				objectDataLocalMap["HealthCheckHttpCode"] = nodeNative10
			}
		}
	}
	request["HealthCheckConfig"] = objectDataLocalMap

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2022-04-30"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		request["ClientToken"] = buildClientToken(action)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_nlb_server_group", action, AlibabaCloudSdkGoERROR)
	}
	if fmt.Sprint(response["Success"]) == "false" {
		return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}

	d.SetId(fmt.Sprint(response["ServerGroupId"]))

	nlbServiceV2 := NlbServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, nlbServiceV2.NlbServerGroupStateRefreshFunc(d.Id(), "ServerGroupStatus", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudNlbServerGroupUpdate(d, meta)
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
	d.Set("resource_group_id", objectRaw["ResourceGroupId"])
	d.Set("scheduler", objectRaw["Scheduler"])
	d.Set("server_group_name", objectRaw["ServerGroupName"])
	d.Set("server_group_type", objectRaw["ServerGroupType"])
	d.Set("status", objectRaw["ServerGroupStatus"])
	d.Set("vpc_id", objectRaw["VpcId"])
	healthCheckMaps := make([]map[string]interface{}, 0)
	healthCheckMap := make(map[string]interface{})
	healthCheck1Raw := make(map[string]interface{})
	if objectRaw["HealthCheck"] != nil {
		healthCheck1Raw = objectRaw["HealthCheck"].(map[string]interface{})
	}
	if len(healthCheck1Raw) > 0 {
		healthCheckMap["health_check_connect_port"] = healthCheck1Raw["HealthCheckConnectPort"]
		healthCheckMap["health_check_connect_timeout"] = healthCheck1Raw["HealthCheckConnectTimeout"]
		healthCheckMap["health_check_domain"] = healthCheck1Raw["HealthCheckDomain"]
		healthCheckMap["health_check_enabled"] = healthCheck1Raw["HealthCheckEnabled"]
		healthCheckMap["health_check_interval"] = healthCheck1Raw["HealthCheckInterval"]
		healthCheckMap["health_check_type"] = healthCheck1Raw["HealthCheckType"]
		healthCheckMap["health_check_url"] = healthCheck1Raw["HealthCheckUrl"]
		healthCheckMap["healthy_threshold"] = healthCheck1Raw["HealthyThreshold"]
		healthCheckMap["http_check_method"] = healthCheck1Raw["HttpCheckMethod"]
		healthCheckMap["unhealthy_threshold"] = healthCheck1Raw["UnhealthyThreshold"]
		healthCheckHttpCode1Raw := make([]interface{}, 0)
		if healthCheck1Raw["HealthCheckHttpCode"] != nil {
			healthCheckHttpCode1Raw = healthCheck1Raw["HealthCheckHttpCode"].([]interface{})
		}

		healthCheckMap["health_check_http_code"] = healthCheckHttpCode1Raw
		healthCheckMaps = append(healthCheckMaps, healthCheckMap)
	}
	d.Set("health_check", healthCheckMaps)
	tagsMaps := objectRaw["Tags"]
	d.Set("tags", tagsToMap(tagsMaps))

	d.Set("connection_drain", d.Get("connection_drain_enabled"))
	return nil
}

func resourceAliCloudNlbServerGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	update := false
	d.Partial(true)
	action := "UpdateServerGroupAttribute"
	conn, err := client.NewNlbClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["ServerGroupId"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)
	if !d.IsNewResource() && d.HasChange("connection_drain") {
		update = true
		request["ConnectionDrainEnabled"] = d.Get("connection_drain")
	}

	if !d.IsNewResource() && d.HasChange("connection_drain_enabled") {
		update = true
		request["ConnectionDrainEnabled"] = d.Get("connection_drain_enabled")
	}

	if !d.IsNewResource() && d.HasChange("connection_drain_timeout") {
		update = true
		request["ConnectionDrainTimeout"] = d.Get("connection_drain_timeout")
	}

	if !d.IsNewResource() && d.HasChange("scheduler") {
		update = true
		request["Scheduler"] = d.Get("scheduler")
	}

	if !d.IsNewResource() && d.HasChange("preserve_client_ip_enabled") {
		update = true
		request["PreserveClientIpEnabled"] = d.Get("preserve_client_ip_enabled")
	}

	if !d.IsNewResource() && d.HasChange("health_check") {
		update = true
		objectDataLocalMap := make(map[string]interface{})
		if v := d.Get("health_check"); !IsNil(v) {
			nodeNative, _ := jsonpath.Get("$[0].health_check_enabled", v)
			if nodeNative != "" {
				objectDataLocalMap["HealthCheckEnabled"] = nodeNative
			}
			nodeNative1, _ := jsonpath.Get("$[0].health_check_type", v)
			if nodeNative1 != "" {
				objectDataLocalMap["HealthCheckType"] = nodeNative1
			}
			nodeNative2, _ := jsonpath.Get("$[0].health_check_connect_port", v)
			if nodeNative2 != "" {
				objectDataLocalMap["HealthCheckConnectPort"] = nodeNative2
			}
			nodeNative3, _ := jsonpath.Get("$[0].healthy_threshold", v)
			if nodeNative3 != "" {
				objectDataLocalMap["HealthyThreshold"] = nodeNative3
			}
			nodeNative4, _ := jsonpath.Get("$[0].unhealthy_threshold", v)
			if nodeNative4 != "" {
				objectDataLocalMap["UnhealthyThreshold"] = nodeNative4
			}
			nodeNative5, _ := jsonpath.Get("$[0].health_check_connect_timeout", v)
			if nodeNative5 != "" {
				objectDataLocalMap["HealthCheckConnectTimeout"] = nodeNative5
			}
			nodeNative6, _ := jsonpath.Get("$[0].health_check_interval", v)
			if nodeNative6 != "" {
				objectDataLocalMap["HealthCheckInterval"] = nodeNative6
			}
			nodeNative7, _ := jsonpath.Get("$[0].health_check_domain", v)
			if nodeNative7 != "" {
				objectDataLocalMap["HealthCheckDomain"] = nodeNative7
			}
			nodeNative8, _ := jsonpath.Get("$[0].health_check_url", v)
			if nodeNative8 != "" {
				objectDataLocalMap["HealthCheckUrl"] = nodeNative8
			}
			nodeNative9, _ := jsonpath.Get("$[0].http_check_method", v)
			if nodeNative9 != "" {
				objectDataLocalMap["HttpCheckMethod"] = nodeNative9
			}
			nodeNative10, _ := jsonpath.Get("$[0].health_check_http_code", v)
			if nodeNative10 != "" {
				objectDataLocalMap["HealthCheckHttpCode"] = nodeNative10
			}
		}
		request["HealthCheckConfig"] = objectDataLocalMap
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2022-04-30"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			request["ClientToken"] = buildClientToken(action)

			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		nlbServiceV2 := NlbServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, nlbServiceV2.NlbServerGroupStateRefreshFunc(d.Id(), "ServerGroupStatus", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("connection_drain_enabled")
		d.SetPartial("connection_drain_timeout")
		d.SetPartial("scheduler")
		d.SetPartial("preserve_client_ip_enabled")
		d.SetPartial("health_check_enabled")
		d.SetPartial("health_check_type")
		d.SetPartial("health_check_connect_port")
		d.SetPartial("healthy_threshold")
		d.SetPartial("unhealthy_threshold")
		d.SetPartial("health_check_connect_timeout")
		d.SetPartial("health_check_interval")
		d.SetPartial("health_check_domain")
		d.SetPartial("health_check_url")
		d.SetPartial("http_check_method")
	}
	update = false
	action = "UpdateServerGroupAttribute"
	conn, err = client.NewNlbClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["ServerGroupId"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)
	if !d.IsNewResource() && d.HasChange("server_group_name") {
		update = true
	}
	request["ServerGroupName"] = d.Get("server_group_name")
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2022-04-30"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			request["ClientToken"] = buildClientToken(action)

			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		nlbServiceV2 := NlbServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, nlbServiceV2.NlbServerGroupStateRefreshFunc(d.Id(), "ServerGroupStatus", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("server_group_name")
	}
	update = false
	action = "MoveResourceGroup"
	conn, err = client.NewNlbClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["ResourceId"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)
	if !d.IsNewResource() && d.HasChange("resource_group_id") {
		update = true
		request["NewResourceGroupId"] = d.Get("resource_group_id")
	}

	request["ResourceType"] = "servergroup"
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2022-04-30"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			request["ClientToken"] = buildClientToken(action)

			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("resource_group_id")
	}

	if d.HasChange("tags") {
		nlbServiceV2 := NlbServiceV2{client}
		if err := nlbServiceV2.SetResourceTags(d, "servergroup"); err != nil {
			return WrapError(err)
		}
		d.SetPartial("tags")
	}
	d.Partial(false)
	return resourceAliCloudNlbServerGroupRead(d, meta)
}

func resourceAliCloudNlbServerGroupDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteServerGroup"
	var request map[string]interface{}
	var response map[string]interface{}
	conn, err := client.NewNlbClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["ServerGroupId"] = d.Id()
	request["RegionId"] = client.RegionId

	request["ClientToken"] = buildClientToken(action)

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2022-04-30"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		request["ClientToken"] = buildClientToken(action)

		if err != nil {
			if NeedRetry(err) || IsExpectedErrors(err, []string{"SystemBusy"}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	if fmt.Sprint(response["Success"]) == "false" {
		return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}
	nlbService := NlbService{client}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, nlbService.NlbServerGroupStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
