package alicloud

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudNlbServerGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudNlbServerGroupCreate,
		Read:   resourceAlicloudNlbServerGroupRead,
		Update: resourceAlicloudNlbServerGroupUpdate,
		Delete: resourceAlicloudNlbServerGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"address_ip_version": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Ipv4", "DualStack"}, false),
			},
			"connection_drain": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"connection_drain_timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(10, 900),
			},
			"health_check": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"health_check_connect_port": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.IntBetween(0, 65535),
						},
						"health_check_connect_timeout": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.IntBetween(1, 300),
						},
						"health_check_domain": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"health_check_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"health_check_interval": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: intBetween(5, 40),
						},
						"health_check_type": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringInSlice([]string{"TCP", "HTTP"}, false),
						},
						"health_check_url": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringLenBetween(1, 80),
						},
						"healthy_threshold": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: intBetween(2, 10),
						},
						"unhealthy_threshold": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: intBetween(2, 10),
						},
						"http_check_method": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringInSlice([]string{"HEAD", "GET"}, false),
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
			"protocol": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"TCP", "UDP", "TCPSSL"}, false),
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"scheduler": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"Wrr", "Rr", "Sch", "Tch", "Qch"}, false),
			},
			"server_group_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9-._]{1,127}$`), "The name must be 2 to 128 characters in length, and can contain letters, digits, periods (.), underscores (_), and hyphens (-). The name must start with a letter."),
			},
			"server_group_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Instance", "Ip"}, false),
			},
			"preserve_client_ip_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
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

func resourceAlicloudNlbServerGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateServerGroup"
	request := make(map[string]interface{})
	conn, err := client.NewNlbClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("address_ip_version"); ok {
		request["AddressIPVersion"] = v
	}
	if v, ok := d.GetOkExists("connection_drain"); ok {
		request["ConnectionDrainEnabled"] = v
	}
	if v, ok := d.GetOk("connection_drain_timeout"); ok {
		request["ConnectionDrainTimeout"] = v
	}
	if v, ok := d.GetOk("protocol"); ok {
		request["Protocol"] = v
	}
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOk("scheduler"); ok {
		request["Scheduler"] = v
	}
	request["ServerGroupName"] = d.Get("server_group_name")
	if v, ok := d.GetOk("server_group_type"); ok {
		request["ServerGroupType"] = v
	}
	if v, ok := d.GetOkExists("preserve_client_ip_enabled"); ok {
		request["PreserveClientIpEnabled"] = v
	}
	if v, ok := d.GetOk("health_check"); ok {
		healthCheckConfig := make(map[string]interface{})
		for _, healthCheckConfigArgs := range v.([]interface{}) {
			healthCheckConfigArg := healthCheckConfigArgs.(map[string]interface{})
			healthCheckConfig["HealthCheckEnabled"] = healthCheckConfigArg["health_check_enabled"]
			if healthCheckConfig["HealthCheckEnabled"] == true {
				healthCheckConfig["HealthCheckConnectPort"] = healthCheckConfigArg["health_check_connect_port"]
				healthCheckConfig["HealthCheckType"] = healthCheckConfigArg["health_check_type"]
				healthCheckConfig["HealthyThreshold"] = healthCheckConfigArg["healthy_threshold"]
				healthCheckConfig["UnhealthyThreshold"] = healthCheckConfigArg["unhealthy_threshold"]
				healthCheckConfig["HealthCheckConnectTimeout"] = healthCheckConfigArg["health_check_connect_timeout"]
				healthCheckConfig["HealthCheckInterval"] = healthCheckConfigArg["health_check_interval"]
				if v, ok := healthCheckConfigArg["health_check_domain"]; ok && fmt.Sprint(v) != "" {
					healthCheckConfig["HealthCheckDomain"] = v
				}
				if v, ok := healthCheckConfigArg["health_check_url"]; ok && fmt.Sprint(v) != "" {
					healthCheckConfig["HealthCheckUrl"] = v
				}
				if v, ok := healthCheckConfigArg["http_check_method"]; ok && fmt.Sprint(v) != "" {
					healthCheckConfig["HttpCheckMethod"] = v
				}
				healthCheckConfig["HealthCheckHttpCode"] = healthCheckConfigArg["health_check_http_code"]
			}
		}
		request["HealthCheckConfig"] = healthCheckConfig
	}
	request["VpcId"] = d.Get("vpc_id")
	request["ClientToken"] = buildClientToken("CreateServerGroup")
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
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
	if fmt.Sprint(response["Success"]) == "false" {
		return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}

	d.SetId(fmt.Sprint(response["ServerGroupId"]))

	nlbService := NlbService{client}
	stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, nlbService.NlbServerGroupStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return resourceAlicloudNlbServerGroupUpdate(d, meta)
}
func resourceAlicloudNlbServerGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	nlbService := NlbService{client}
	object, err := nlbService.DescribeNlbServerGroup(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_nlb_server_group nlbService.DescribeNlbServerGroup Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("address_ip_version", object["AddressIPVersion"])
	d.Set("connection_drain", object["ConnectionDrainEnabled"])
	if v, ok := object["ConnectionDrainTimeout"]; ok && fmt.Sprint(v) != "0" {
		d.Set("connection_drain_timeout", formatInt(v))
	}

	healthCheckSli := make([]map[string]interface{}, 0)
	if v, ok := object["HealthCheck"]; ok {
		if len(v.(map[string]interface{})) > 0 {
			healthCheck := v.(map[string]interface{})
			healthCheckMap := make(map[string]interface{})
			healthCheckMap["health_check_connect_port"] = healthCheck["HealthCheckConnectPort"]
			healthCheckMap["health_check_connect_timeout"] = healthCheck["HealthCheckConnectTimeout"]
			healthCheckMap["health_check_domain"] = healthCheck["HealthCheckDomain"]
			healthCheckMap["health_check_enabled"] = healthCheck["HealthCheckEnabled"]
			healthCheckMap["health_check_interval"] = healthCheck["HealthCheckInterval"]
			healthCheckMap["health_check_type"] = healthCheck["HealthCheckType"]
			healthCheckMap["health_check_url"] = healthCheck["HealthCheckUrl"]
			healthCheckMap["healthy_threshold"] = healthCheck["HealthyThreshold"]
			healthCheckMap["unhealthy_threshold"] = healthCheck["UnhealthyThreshold"]
			healthCheckMap["http_check_method"] = healthCheck["HttpCheckMethod"]
			healthCheckMap["health_check_http_code"] = healthCheck["HealthCheckHttpCode"]
			healthCheckSli = append(healthCheckSli, healthCheckMap)
		}
	}
	d.Set("health_check", healthCheckSli)
	d.Set("protocol", object["Protocol"])
	d.Set("scheduler", object["Scheduler"])
	d.Set("server_group_name", object["ServerGroupName"])
	d.Set("server_group_type", object["ServerGroupType"])
	d.Set("status", object["ServerGroupStatus"])
	d.Set("vpc_id", object["VpcId"])
	d.Set("resource_group_id", object["ResourceGroupId"])
	d.Set("preserve_client_ip_enabled", object["PreserveClientIpEnabled"])
	if v, ok := object["Tags"]; ok && len(v.([]interface{})) > 0 {
		d.Set("tags", tagsToMap(v))
	}

	return nil
}
func resourceAlicloudNlbServerGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	nlbService := NlbService{client}
	var response map[string]interface{}
	conn, err := client.NewNlbClient()
	if err != nil {
		return WrapError(err)
	}
	d.Partial(true)

	if d.HasChange("tags") {
		if err := nlbService.SetResourceTags(d, "servergroup"); err != nil {
			return WrapError(err)
		}
		d.SetPartial("tags")
	}
	update := false
	request := map[string]interface{}{
		"ServerGroupId": d.Id(),
	}
	if !d.IsNewResource() && d.HasChange("connection_drain") {
		update = true
		if v, ok := d.GetOkExists("connection_drain"); ok {
			request["ConnectionDrainEnabled"] = v
		}
	}
	if !d.IsNewResource() && d.HasChange("connection_drain_timeout") {
		update = true
		if v, ok := d.GetOk("connection_drain_timeout"); ok {
			request["ConnectionDrainTimeout"] = v
		}
	}
	request["RegionId"] = client.RegionId
	if !d.IsNewResource() && d.HasChange("scheduler") {
		update = true
		if v, ok := d.GetOk("scheduler"); ok {
			request["Scheduler"] = v
		}
	}

	if !d.IsNewResource() && d.HasChange("health_check") {
		update = true
		if v, ok := d.GetOk("health_check"); ok {
			healthCheckConfig := make(map[string]interface{})
			for _, healthCheckConfigArgs := range v.([]interface{}) {
				healthCheckConfigArg := healthCheckConfigArgs.(map[string]interface{})
				healthCheckConfig["HealthCheckEnabled"] = healthCheckConfigArg["health_check_enabled"]
				if healthCheckConfig["HealthCheckEnabled"] == true {
					healthCheckConfig["HealthCheckConnectPort"] = healthCheckConfigArg["health_check_connect_port"]
					healthCheckConfig["HealthCheckType"] = healthCheckConfigArg["health_check_type"]
					healthCheckConfig["HealthyThreshold"] = healthCheckConfigArg["healthy_threshold"]
					healthCheckConfig["UnhealthyThreshold"] = healthCheckConfigArg["unhealthy_threshold"]
					healthCheckConfig["HealthCheckConnectTimeout"] = healthCheckConfigArg["health_check_connect_timeout"]
					healthCheckConfig["HealthCheckInterval"] = healthCheckConfigArg["health_check_interval"]
					if v, ok := healthCheckConfigArg["health_check_domain"]; ok && fmt.Sprint(v) != "" {
						healthCheckConfig["HealthCheckDomain"] = v
					}
					if v, ok := healthCheckConfigArg["health_check_url"]; ok && fmt.Sprint(v) != "" {
						healthCheckConfig["HealthCheckUrl"] = v
					}
					healthCheckConfig["HttpCheckMethod"] = healthCheckConfigArg["http_check_method"]
					healthCheckConfig["HealthCheckHttpCode"] = healthCheckConfigArg["health_check_http_code"]
				}
			}
			request["HealthCheckConfig"] = healthCheckConfig
		}
	}
	if update {
		action := "UpdateServerGroupAttribute"
		request["ClientToken"] = buildClientToken("UpdateServerGroupAttribute")
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
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
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		if fmt.Sprint(response["Success"]) == "false" {
			return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
		}
		stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, nlbService.NlbServerGroupStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("connection_drain")
		d.SetPartial("connection_drain_timeout")
		d.SetPartial("scheduler")
	}
	update = false
	addServersToServerGroupReq := map[string]interface{}{
		"ServerGroupId": d.Id(),
	}
	if !d.IsNewResource() && d.HasChange("server_group_name") {
		update = true
		addServersToServerGroupReq["ServerGroupName"] = d.Get("server_group_name")
	}
	if update {
		action := "UpdateServerGroupAttribute"
		addServersToServerGroupReq["ClientToken"] = buildClientToken("UpdateServerGroupAttribute")
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2022-04-30"), StringPointer("AK"), nil, addServersToServerGroupReq, &runtime)
			if err != nil {
				if NeedRetry(err) {
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
		if fmt.Sprint(response["Success"]) == "false" {
			return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
		}
		d.SetPartial("server_group_name")
	}
	d.Partial(false)
	return resourceAlicloudNlbServerGroupRead(d, meta)
}
func resourceAlicloudNlbServerGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteServerGroup"
	var response map[string]interface{}
	conn, err := client.NewNlbClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"ServerGroupId": d.Id(),
	}

	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken("DeleteServerGroup")
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2022-04-30"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) || IsExpectedErrors(err, []string{"SystemBusy"}) {
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
