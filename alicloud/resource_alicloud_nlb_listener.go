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

func resourceAliCloudNlbListener() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudNlbListenerCreate,
		Read:   resourceAliCloudNlbListenerRead,
		Update: resourceAliCloudNlbListenerUpdate,
		Delete: resourceAliCloudNlbListenerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"alpn_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"alpn_policy": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ca_certificate_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"ca_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"certificate_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"cps": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: IntBetween(0, 1000000),
			},
			"end_port": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: IntBetween(0, 65535),
			},
			"idle_timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: IntBetween(0, 900),
			},
			"listener_description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringLenBetween(2, 256),
			},
			"listener_port": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: IntBetween(0, 65535),
			},
			"listener_protocol": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"TCP", "UDP", "TCPSSL"}, false),
			},
			"load_balancer_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"mss": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: IntBetween(0, 1500),
			},
			"proxy_protocol_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"sec_sensor_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"security_policy_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"server_group_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"start_port": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: IntBetween(0, 65535),
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"Stopped", "Running"}, false),
			},
			"tags": tagsSchema(),
		},
	}
}

func resourceAliCloudNlbListenerCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateListener"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	conn, err := client.NewNlbClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	request["ListenerProtocol"] = d.Get("listener_protocol")
	request["ListenerPort"] = d.Get("listener_port")
	if v, ok := d.GetOk("listener_description"); ok {
		request["ListenerDescription"] = v
	}
	request["LoadBalancerId"] = d.Get("load_balancer_id")
	request["ServerGroupId"] = d.Get("server_group_id")
	if v, ok := d.GetOk("idle_timeout"); ok {
		request["IdleTimeout"] = v
	}
	if v, ok := d.GetOk("security_policy_id"); ok {
		request["SecurityPolicyId"] = v
	}
	if v, ok := d.GetOk("certificate_ids"); ok {
		certificateIdsMaps := v.([]interface{})
		request["CertificateIds"] = certificateIdsMaps
	}

	if v, ok := d.GetOk("ca_certificate_ids"); ok {
		caCertificateIdsMaps := v.([]interface{})
		request["CaCertificateIds"] = caCertificateIdsMaps
	}

	if v, ok := d.GetOk("alpn_policy"); ok {
		request["AlpnPolicy"] = v
	}
	if v, ok := d.GetOkExists("proxy_protocol_enabled"); ok {
		request["ProxyProtocolEnabled"] = v
	}
	if v, ok := d.GetOkExists("sec_sensor_enabled"); ok {
		request["SecSensorEnabled"] = v
	}
	if v, ok := d.GetOkExists("alpn_enabled"); ok {
		request["AlpnEnabled"] = v
	}
	if v, ok := d.GetOkExists("ca_enabled"); ok {
		request["CaEnabled"] = v
	}
	if v, ok := d.GetOk("start_port"); ok {
		request["StartPort"] = v
	}
	if v, ok := d.GetOk("end_port"); ok {
		request["EndPort"] = v
	}
	if v, ok := d.GetOk("cps"); ok {
		request["Cps"] = v
	}
	if v, ok := d.GetOk("mss"); ok {
		request["Mss"] = v
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2022-04-30"), StringPointer("AK"), query, request, &runtime)
		request["ClientToken"] = buildClientToken(action)

		if err != nil {
			if IsExpectedErrors(err, []string{"SystemBusy"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_nlb_listener", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["ListenerId"]))

	nlbServiceV2 := NlbServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, nlbServiceV2.NlbListenerStateRefreshFunc(d.Id(), "ListenerStatus", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudNlbListenerUpdate(d, meta)
}

func resourceAliCloudNlbListenerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	nlbServiceV2 := NlbServiceV2{client}

	objectRaw, err := nlbServiceV2.DescribeNlbListener(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_nlb_listener DescribeNlbListener Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("alpn_enabled", objectRaw["AlpnEnabled"])
	d.Set("alpn_policy", objectRaw["AlpnPolicy"])
	d.Set("ca_enabled", objectRaw["CaEnabled"])
	d.Set("cps", objectRaw["Cps"])
	d.Set("end_port", formatInt(objectRaw["EndPort"]))
	d.Set("idle_timeout", objectRaw["IdleTimeout"])
	d.Set("listener_description", objectRaw["ListenerDescription"])
	d.Set("listener_port", objectRaw["ListenerPort"])
	d.Set("listener_protocol", objectRaw["ListenerProtocol"])
	d.Set("load_balancer_id", objectRaw["LoadBalancerId"])
	d.Set("mss", objectRaw["Mss"])
	d.Set("proxy_protocol_enabled", objectRaw["ProxyProtocolEnabled"])
	d.Set("sec_sensor_enabled", objectRaw["SecSensorEnabled"])
	d.Set("security_policy_id", objectRaw["SecurityPolicyId"])
	d.Set("server_group_id", objectRaw["ServerGroupId"])
	d.Set("start_port", formatInt(objectRaw["StartPort"]))
	d.Set("status", objectRaw["ListenerStatus"])

	caCertificateIds1Raw := make([]interface{}, 0)
	if objectRaw["CaCertificateIds"] != nil {
		caCertificateIds1Raw = objectRaw["CaCertificateIds"].([]interface{})
	}

	d.Set("ca_certificate_ids", caCertificateIds1Raw)
	certificateIds1Raw := make([]interface{}, 0)
	if objectRaw["CertificateIds"] != nil {
		certificateIds1Raw = objectRaw["CertificateIds"].([]interface{})
	}

	d.Set("certificate_ids", certificateIds1Raw)
	tagsMaps := objectRaw["Tags"]
	d.Set("tags", tagsToMap(tagsMaps))

	return nil
}

func resourceAliCloudNlbListenerUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	action := "UpdateListenerAttribute"
	conn, err := client.NewNlbClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ListenerId"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)
	if !d.IsNewResource() && d.HasChange("listener_description") {
		update = true
		request["ListenerDescription"] = d.Get("listener_description")
	}

	if !d.IsNewResource() && d.HasChange("server_group_id") {
		update = true
	}
	request["ServerGroupId"] = d.Get("server_group_id")
	if !d.IsNewResource() && d.HasChange("security_policy_id") {
		update = true
		request["SecurityPolicyId"] = d.Get("security_policy_id")
	}

	if !d.IsNewResource() && d.HasChange("certificate_ids") {
		update = true
		if v, ok := d.GetOk("certificate_ids"); ok {
			certificateIdsMaps := v.([]interface{})
			request["CertificateIds"] = certificateIdsMaps
		}
	}

	if !d.IsNewResource() && d.HasChange("ca_certificate_ids") {
		update = true
		if v, ok := d.GetOk("ca_certificate_ids"); ok {
			caCertificateIdsMaps := v.([]interface{})
			request["CaCertificateIds"] = caCertificateIdsMaps
		}
	}

	if !d.IsNewResource() && d.HasChange("idle_timeout") {
		update = true
		request["IdleTimeout"] = d.Get("idle_timeout")
	}

	if !d.IsNewResource() && d.HasChange("alpn_policy") {
		update = true
		request["AlpnPolicy"] = d.Get("alpn_policy")
	}

	if !d.IsNewResource() && d.HasChange("ca_enabled") {
		update = true
		request["CaEnabled"] = d.Get("ca_enabled")
	}

	if !d.IsNewResource() && d.HasChange("proxy_protocol_enabled") {
		update = true
		request["ProxyProtocolEnabled"] = d.Get("proxy_protocol_enabled")
	}

	if !d.IsNewResource() && d.HasChange("sec_sensor_enabled") {
		update = true
		request["SecSensorEnabled"] = d.Get("sec_sensor_enabled")
	}

	if !d.IsNewResource() && d.HasChange("alpn_enabled") {
		update = true
		request["AlpnEnabled"] = d.Get("alpn_enabled")
	}

	if !d.IsNewResource() && d.HasChange("cps") {
		update = true
		request["Cps"] = d.Get("cps")
	}

	if !d.IsNewResource() && d.HasChange("mss") {
		update = true
		request["Mss"] = d.Get("mss")
	}

	if update {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2022-04-30"), StringPointer("AK"), query, request, &runtime)
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
		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 15*time.Second, nlbServiceV2.NlbListenerStateRefreshFunc(d.Id(), "ListenerStatus", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	if d.HasChange("status") {
		client := meta.(*connectivity.AliyunClient)
		nlbServiceV2 := NlbServiceV2{client}
		object, err := nlbServiceV2.DescribeNlbListener(d.Id())
		if err != nil {
			return WrapError(err)
		}

		target := d.Get("status").(string)
		if object["ListenerStatus"].(string) != target {
			if target == "Running" {
				action = "StartListener"
				conn, err = client.NewNlbClient()
				if err != nil {
					return WrapError(err)
				}
				request = make(map[string]interface{})
				query = make(map[string]interface{})
				request["ListenerId"] = d.Id()
				request["RegionId"] = client.RegionId
				request["ClientToken"] = buildClientToken(action)
				runtime := util.RuntimeOptions{}
				runtime.SetAutoretry(true)
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2022-04-30"), StringPointer("AK"), query, request, &runtime)
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
				stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, nlbServiceV2.NlbListenerStateRefreshFunc(d.Id(), "ListenerStatus", []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}

			}
			if target == "Stopped" {
				action = "StopListener"
				conn, err = client.NewNlbClient()
				if err != nil {
					return WrapError(err)
				}
				request = make(map[string]interface{})
				query = make(map[string]interface{})
				request["ListenerId"] = d.Id()
				request["RegionId"] = client.RegionId
				request["ClientToken"] = buildClientToken(action)
				runtime := util.RuntimeOptions{}
				runtime.SetAutoretry(true)
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2022-04-30"), StringPointer("AK"), query, request, &runtime)
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
				stateConf := BuildStateConf([]string{}, []string{"Stopped"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, nlbServiceV2.NlbListenerStateRefreshFunc(d.Id(), "ListenerStatus", []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}

			}
		}
	}

	if d.HasChange("tags") {
		nlbServiceV2 := NlbServiceV2{client}
		if err := nlbServiceV2.SetResourceTags(d, "listener"); err != nil {
			return WrapError(err)
		}
		d.SetPartial("tags")
	}
	return resourceAliCloudNlbListenerRead(d, meta)
}

func resourceAliCloudNlbListenerDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteListener"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	conn, err := client.NewNlbClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["ListenerId"] = d.Id()
	request["RegionId"] = client.RegionId

	request["ClientToken"] = buildClientToken(action)

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2022-04-30"), StringPointer("AK"), query, request, &runtime)
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
	stateConf := BuildStateConf([]string{}, []string{"Succeeded"}, d.Timeout(schema.TimeoutDelete), 5*time.Second, nlbServiceV2.DescribeAsyncNlbListenerStateRefreshFunc(d, response, "$.Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
