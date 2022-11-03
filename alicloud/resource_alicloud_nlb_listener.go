package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudNlbListener() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudNlbListenerCreate,
		Read:   resourceAlicloudNlbListenerRead,
		Update: resourceAlicloudNlbListenerUpdate,
		Delete: resourceAlicloudNlbListenerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create:  schema.DefaultTimeout(1 * time.Minute),
			Update:  schema.DefaultTimeout(1 * time.Minute),
			Default: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"alpn_enabled": {
				Computed: true,
				Optional: true,
				Type:     schema.TypeBool,
			},
			"alpn_policy": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"ca_certificate_ids": {
				Optional: true,
				Type:     schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"ca_enabled": {
				Optional: true,
				Computed: true,
				Type:     schema.TypeBool,
			},
			"certificate_ids": {
				Optional: true,
				Type:     schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"cps": {
				Optional:     true,
				Type:         schema.TypeInt,
				ValidateFunc: validation.IntBetween(0, 1000000),
			},
			"end_port": {
				Optional:     true,
				ForceNew:     true,
				Type:         schema.TypeInt,
				ValidateFunc: validation.IntBetween(0, 65535),
			},
			"idle_timeout": {
				Optional:     true,
				Computed:     true,
				Type:         schema.TypeInt,
				ValidateFunc: validation.IntBetween(0, 900),
			},
			"listener_description": {
				Optional:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.StringLenBetween(2, 256),
			},
			"listener_port": {
				Required:     true,
				ForceNew:     true,
				Type:         schema.TypeInt,
				ValidateFunc: validation.IntBetween(0, 65535),
			},
			"listener_protocol": {
				Required:     true,
				ForceNew:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"TCP", "UDP", "TCPSSL"}, false),
			},
			"load_balancer_id": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"mss": {
				Optional:     true,
				Type:         schema.TypeInt,
				ValidateFunc: validation.IntBetween(0, 1500),
			},
			"proxy_protocol_enabled": {
				Optional: true,
				Computed: true,
				Type:     schema.TypeBool,
			},
			"sec_sensor_enabled": {
				Computed: true,
				Optional: true,
				Type:     schema.TypeBool,
			},
			"security_policy_id": {
				Computed:     true,
				Optional:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"tls_cipher_policy_1_0", "tls_cipher_policy_1_1", "tls_cipher_policy_1_2", "tls_cipher_policy_1_2_strict", "tls_cipher_policy_1_2_strict_with_1_3"}, false),
			},
			"server_group_id": {
				Required: true,
				Type:     schema.TypeString,
			},
			"start_port": {
				Optional:     true,
				ForceNew:     true,
				Type:         schema.TypeInt,
				ValidateFunc: validation.IntBetween(0, 65535),
			},
			"status": {
				Optional:     true,
				Computed:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"Stopped", "Running"}, false),
			},
		},
	}
}

func resourceAlicloudNlbListenerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	nlbService := NlbService{client}
	request := map[string]interface{}{
		"RegionId": client.RegionId,
	}
	conn, err := client.NewNlbClient()
	if err != nil {
		return WrapError(err)
	}

	if v, ok := d.GetOk("listener_description"); ok {
		request["ListenerDescription"] = v
	}

	request["ListenerPort"] = d.Get("listener_port")

	request["ListenerProtocol"] = d.Get("listener_protocol")

	request["LoadBalancerId"] = d.Get("load_balancer_id")

	request["ServerGroupId"] = d.Get("server_group_id")

	if v, ok := d.GetOkExists("idle_timeout"); ok {
		request["IdleTimeout"] = v
	}
	if v, ok := d.GetOkExists("cps"); ok {
		request["Cps"] = v
	}
	if v, ok := d.GetOkExists("proxy_protocol_enabled"); ok {
		request["ProxyProtocolEnabled"] = v
	}
	if v, ok := d.GetOkExists("mss"); ok {
		request["Mss"] = v
	}
	if v, ok := d.GetOkExists("sec_sensor_enabled"); ok {
		request["SecSensorEnabled"] = v
	}
	if v, ok := d.GetOkExists("ca_enabled"); ok {
		request["CaEnabled"] = v
	}
	if v, ok := d.GetOkExists("end_port"); ok {
		request["EndPort"] = v
	}
	if v, ok := d.GetOkExists("start_port"); ok {
		request["StartPort"] = v
	}
	if v, ok := d.GetOk("alpn_policy"); ok {
		request["AlpnPolicy"] = v
	}
	if v, ok := d.GetOkExists("alpn_enabled"); ok {
		request["AlpnEnabled"] = v
	}
	if v, ok := d.GetOk("ca_certificate_ids"); ok {
		request["CaCertificateIds"] = v.([]interface{})
	}
	if v, ok := d.GetOk("certificate_ids"); ok {
		request["CertificateIds"] = v.([]interface{})
	}
	if v, ok := d.GetOk("security_policy_id"); ok {
		request["SecurityPolicyId"] = v
	}

	request["ClientToken"] = buildClientToken("CreateListener")
	var response map[string]interface{}
	action := "CreateListener"
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		resp, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2022-04-30"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		response = resp
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_nlb_listener", action, AlibabaCloudSdkGoERROR)
	}

	if v, err := jsonpath.Get("$.ListenerId", response); err != nil || v == nil {
		return WrapErrorf(err, IdMsg, "alicloud_nlb_listener")
	} else {
		d.SetId(fmt.Sprint(v))
	}
	stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, nlbService.NlbListenerStateRefreshFunc(d, []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return resourceAlicloudNlbListenerUpdate(d, meta)
}

func resourceAlicloudNlbListenerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	nlbService := NlbService{client}

	object, err := nlbService.DescribeNlbListener(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_nlb_listener nlbService.DescribeNlbListener Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("alpn_enabled", object["AlpnEnabled"])
	d.Set("alpn_policy", object["AlpnPolicy"])
	caCertificateIds, _ := jsonpath.Get("$.CaCertificateIds", object)
	d.Set("ca_certificate_ids", caCertificateIds)
	d.Set("ca_enabled", object["CaEnabled"])
	certificateIds, _ := jsonpath.Get("$.CertificateIds", object)
	d.Set("certificate_ids", certificateIds)
	d.Set("cps", object["Cps"])
	d.Set("end_port", formatInt(object["EndPort"]))
	d.Set("idle_timeout", object["IdleTimeout"])
	d.Set("listener_description", object["ListenerDescription"])
	d.Set("listener_port", object["ListenerPort"])
	d.Set("listener_protocol", object["ListenerProtocol"])
	d.Set("load_balancer_id", object["LoadBalancerId"])
	d.Set("mss", object["Mss"])
	d.Set("proxy_protocol_enabled", object["ProxyProtocolEnabled"])
	d.Set("sec_sensor_enabled", object["SecSensorEnabled"])
	d.Set("security_policy_id", object["SecurityPolicyId"])
	d.Set("server_group_id", object["ServerGroupId"])
	d.Set("start_port", formatInt(object["StartPort"]))
	d.Set("status", object["ListenerStatus"])

	return nil
}

func resourceAlicloudNlbListenerUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	nlbService := NlbService{client}
	d.Partial(true)
	conn, err := client.NewNlbClient()
	if err != nil {
		return WrapError(err)
	}
	update := false
	request := map[string]interface{}{
		"ListenerId": d.Id(),
		"RegionId":   client.RegionId,
	}

	if !d.IsNewResource() && d.HasChange("alpn_enabled") {
		update = true
		if v, ok := d.GetOkExists("alpn_enabled"); ok {
			request["AlpnEnabled"] = v
		}
	}
	if !d.IsNewResource() && d.HasChange("alpn_policy") {
		update = true
		if v, ok := d.GetOk("alpn_policy"); ok {
			request["AlpnPolicy"] = v
		}
	}
	if !d.IsNewResource() && d.HasChange("ca_certificate_ids") {
		update = true
		if v, ok := d.GetOk("ca_certificate_ids"); ok {
			request["CaCertificateIds"] = v.([]interface{})
		}
	}
	if !d.IsNewResource() && d.HasChange("ca_enabled") {
		update = true
		if v, ok := d.GetOkExists("ca_enabled"); ok {
			request["CaEnabled"] = v
		}
	}
	if !d.IsNewResource() && d.HasChange("certificate_ids") {
		update = true
		if v, ok := d.GetOk("certificate_ids"); ok {
			request["CertificateIds"] = v.([]interface{})
		}
	}
	if !d.IsNewResource() && d.HasChange("cps") {
		update = true
		if v, ok := d.GetOk("cps"); ok {
			request["Cps"] = v
		}
	}
	if !d.IsNewResource() && d.HasChange("idle_timeout") {
		update = true
		if v, ok := d.GetOk("idle_timeout"); ok {
			request["IdleTimeout"] = v
		}
	}
	if !d.IsNewResource() && d.HasChange("listener_description") {
		update = true
		if v, ok := d.GetOk("listener_description"); ok {
			request["ListenerDescription"] = v
		}
	}
	if !d.IsNewResource() && d.HasChange("mss") {
		update = true
		if v, ok := d.GetOkExists("mss"); ok {
			request["Mss"] = v
		}
	}
	if !d.IsNewResource() && d.HasChange("proxy_protocol_enabled") {
		update = true
		if v, ok := d.GetOkExists("proxy_protocol_enabled"); ok {
			request["ProxyProtocolEnabled"] = v
		}
	}
	if !d.IsNewResource() && d.HasChange("sec_sensor_enabled") {
		update = true
		if v, ok := d.GetOkExists("sec_sensor_enabled"); ok {
			request["SecSensorEnabled"] = v
		}
	}
	if !d.IsNewResource() && d.HasChange("security_policy_id") {
		update = true
		if v, ok := d.GetOk("security_policy_id"); ok {
			request["SecurityPolicyId"] = v
		}
	}
	if !d.IsNewResource() && d.HasChange("server_group_id") {
		update = true
		if v, ok := d.GetOk("server_group_id"); ok {
			request["ServerGroupId"] = v
		}
	}

	if update {
		action := "UpdateListenerAttribute"
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			resp, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2022-04-30"), StringPointer("AK"), nil, request, &runtime)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, resp, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, nlbService.NlbListenerStateRefreshFunc(d, []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("alpn_enabled")
		d.SetPartial("alpn_policy")
		d.SetPartial("ca_certificate_ids")
		d.SetPartial("ca_enabled")
		d.SetPartial("certificate_ids")
		d.SetPartial("cps")
		d.SetPartial("idle_timeout")
		d.SetPartial("listener_description")
		d.SetPartial("mss")
		d.SetPartial("proxy_protocol_enabled")
		d.SetPartial("sec_sensor_enabled")
		d.SetPartial("security_policy_id")
		d.SetPartial("server_group_id")
	}

	if d.HasChange("status") {
		object, err := nlbService.DescribeNlbListener(d.Id())
		if err != nil {
			WrapError(err)
		}
		target := fmt.Sprint(d.Get("status"))
		if fmt.Sprint(object["Status"]) != target {
			if target == "Running" {
				request := map[string]interface{}{
					"ListenerId": d.Id(),
					"RegionId":   client.RegionId,
				}

				action := "StartListener"
				runtime := util.RuntimeOptions{}
				runtime.SetAutoretry(true)
				wait := incrementalWait(3*time.Second, 3*time.Second)
				err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
					resp, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2022-04-30"), StringPointer("AK"), nil, request, &runtime)
					if err != nil {
						if NeedRetry(err) {
							wait()
							return resource.RetryableError(err)
						}
						return resource.NonRetryableError(err)
					}
					addDebug(action, resp, request)
					return nil
				})
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
				}
				stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, nlbService.NlbListenerStateRefreshFunc(d, []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}

			}
			if target == "Stopped" {
				request := map[string]interface{}{
					"ListenerId": d.Id(),
					"RegionId":   client.RegionId,
				}

				action := "StopListener"
				runtime := util.RuntimeOptions{}
				runtime.SetAutoretry(true)
				wait := incrementalWait(3*time.Second, 3*time.Second)
				err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
					resp, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2022-04-30"), StringPointer("AK"), nil, request, &runtime)
					if err != nil {
						if NeedRetry(err) {
							wait()
							return resource.RetryableError(err)
						}
						return resource.NonRetryableError(err)
					}
					addDebug(action, resp, request)
					return nil
				})
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
				}
				stateConf := BuildStateConf([]string{}, []string{"Stopped"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, nlbService.NlbListenerStateRefreshFunc(d, []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}

			}
			d.SetPartial("status")
		}
	}
	d.Partial(false)
	return resourceAlicloudNlbListenerRead(d, meta)
}

func resourceAlicloudNlbListenerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	conn, err := client.NewNlbClient()
	if err != nil {
		return WrapError(err)
	}

	request := map[string]interface{}{
		"ListenerId": d.Id(),
		"RegionId":   client.RegionId,
	}

	request["ClientToken"] = buildClientToken("DeleteListener")
	action := "DeleteListener"
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		resp, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2022-04-30"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, resp, request)
		return nil
	})
	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}
