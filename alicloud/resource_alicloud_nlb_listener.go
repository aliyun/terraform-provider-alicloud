// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
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
				Type:     schema.TypeString,
				Optional: true,
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
			"proxy_protocol_config": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"proxy_protocol_config_private_link_eps_id_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"proxy_protocol_config_private_link_ep_id_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"proxy_protocol_config_vpc_id_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
			"proxy_protocol_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"region_id": {
				Type:     schema.TypeString,
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
	var err error
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
	if v, ok := d.GetOkExists("idle_timeout"); ok && v.(int) > 0 {
		request["IdleTimeout"] = v
	}
	if v, ok := d.GetOk("security_policy_id"); ok {
		request["SecurityPolicyId"] = v
	}
	if v, ok := d.GetOk("certificate_ids"); ok {
		certificateIdsMapsArray := v.([]interface{})
		request["CertificateIds"] = certificateIdsMapsArray
	}

	if v, ok := d.GetOk("ca_certificate_ids"); ok {
		caCertificateIdsMapsArray := v.([]interface{})
		request["CaCertificateIds"] = caCertificateIdsMapsArray
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
	if v, ok := d.GetOkExists("start_port"); ok {
		request["StartPort"] = v
	}
	if v, ok := d.GetOkExists("end_port"); ok {
		request["EndPort"] = v
	}
	if v, ok := d.GetOkExists("cps"); ok {
		request["Cps"] = v
	}
	if v, ok := d.GetOkExists("mss"); ok {
		request["Mss"] = v
	}
	if v, ok := d.GetOk("tags"); ok {
		tagsMap := ConvertTags(v.(map[string]interface{}))
		request = expandTagsToMap(request, tagsMap)
	}

	objectDataLocalMap := make(map[string]interface{})

	if v := d.Get("proxy_protocol_config"); !IsNil(v) {
		proxyProtocolConfigVpcIdEnabled, _ := jsonpath.Get("$[0].proxy_protocol_config_vpc_id_enabled", v)
		if proxyProtocolConfigVpcIdEnabled != nil && proxyProtocolConfigVpcIdEnabled != "" {
			objectDataLocalMap["Ppv2VpcIdEnabled"] = proxyProtocolConfigVpcIdEnabled
		}
		proxyProtocolConfigPrivateLinkEpIdEnabled, _ := jsonpath.Get("$[0].proxy_protocol_config_private_link_ep_id_enabled", v)
		if proxyProtocolConfigPrivateLinkEpIdEnabled != nil && proxyProtocolConfigPrivateLinkEpIdEnabled != "" {
			objectDataLocalMap["Ppv2PrivateLinkEpIdEnabled"] = proxyProtocolConfigPrivateLinkEpIdEnabled
		}
		proxyProtocolConfigPrivateLinkEpsIdEnabled, _ := jsonpath.Get("$[0].proxy_protocol_config_private_link_eps_id_enabled", v)
		if proxyProtocolConfigPrivateLinkEpsIdEnabled != nil && proxyProtocolConfigPrivateLinkEpsIdEnabled != "" {
			objectDataLocalMap["Ppv2PrivateLinkEpsIdEnabled"] = proxyProtocolConfigPrivateLinkEpsIdEnabled
		}

		objectDataLocalMapJson, err := json.Marshal(objectDataLocalMap)
		if err != nil {
			return WrapError(err)
		}
		request["ProxyProtocolV2Config"] = string(objectDataLocalMapJson)
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Nlb", "2022-04-30", action, query, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"SystemBusy"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

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

	if objectRaw["AlpnEnabled"] != nil {
		d.Set("alpn_enabled", objectRaw["AlpnEnabled"])
	}
	if objectRaw["AlpnPolicy"] != nil {
		d.Set("alpn_policy", objectRaw["AlpnPolicy"])
	}
	if objectRaw["CaEnabled"] != nil {
		d.Set("ca_enabled", objectRaw["CaEnabled"])
	}
	if objectRaw["Cps"] != nil {
		d.Set("cps", objectRaw["Cps"])
	}
	if objectRaw["EndPort"] != nil {
		d.Set("end_port", formatInt(objectRaw["EndPort"]))
	}
	if objectRaw["IdleTimeout"] != nil {
		d.Set("idle_timeout", objectRaw["IdleTimeout"])
	}
	if objectRaw["ListenerDescription"] != nil {
		d.Set("listener_description", objectRaw["ListenerDescription"])
	}
	if objectRaw["ListenerPort"] != nil {
		d.Set("listener_port", objectRaw["ListenerPort"])
	}
	if objectRaw["ListenerProtocol"] != nil {
		d.Set("listener_protocol", objectRaw["ListenerProtocol"])
	}
	if objectRaw["LoadBalancerId"] != nil {
		d.Set("load_balancer_id", objectRaw["LoadBalancerId"])
	}
	if objectRaw["Mss"] != nil {
		d.Set("mss", objectRaw["Mss"])
	}
	if objectRaw["ProxyProtocolEnabled"] != nil {
		d.Set("proxy_protocol_enabled", objectRaw["ProxyProtocolEnabled"])
	}
	if objectRaw["RegionId"] != nil {
		d.Set("region_id", objectRaw["RegionId"])
	}
	if objectRaw["SecSensorEnabled"] != nil {
		d.Set("sec_sensor_enabled", objectRaw["SecSensorEnabled"])
	}
	if objectRaw["SecurityPolicyId"] != nil {
		d.Set("security_policy_id", objectRaw["SecurityPolicyId"])
	}
	if objectRaw["ServerGroupId"] != nil {
		d.Set("server_group_id", objectRaw["ServerGroupId"])
	}
	if objectRaw["StartPort"] != nil {
		d.Set("start_port", formatInt(objectRaw["StartPort"]))
	}
	if objectRaw["ListenerStatus"] != nil {
		d.Set("status", objectRaw["ListenerStatus"])
	}

	caCertificateIds2Raw := make([]interface{}, 0)
	if objectRaw["CaCertificateIds"] != nil {
		caCertificateIds2Raw = objectRaw["CaCertificateIds"].([]interface{})
	}

	d.Set("ca_certificate_ids", caCertificateIds2Raw)
	certificateIds2Raw := make([]interface{}, 0)
	if objectRaw["CertificateIds"] != nil {
		certificateIds2Raw = objectRaw["CertificateIds"].([]interface{})
	}

	d.Set("certificate_ids", certificateIds2Raw)
	proxyProtocolConfigMaps := make([]map[string]interface{}, 0)
	proxyProtocolConfigMap := make(map[string]interface{})
	proxyProtocolV2Config2Raw := make(map[string]interface{})
	if objectRaw["ProxyProtocolV2Config"] != nil {
		proxyProtocolV2Config2Raw = objectRaw["ProxyProtocolV2Config"].(map[string]interface{})
	}
	if len(proxyProtocolV2Config2Raw) > 0 {
		proxyProtocolConfigMap["proxy_protocol_config_private_link_ep_id_enabled"] = proxyProtocolV2Config2Raw["Ppv2PrivateLinkEpIdEnabled"]
		proxyProtocolConfigMap["proxy_protocol_config_private_link_eps_id_enabled"] = proxyProtocolV2Config2Raw["Ppv2PrivateLinkEpsIdEnabled"]
		proxyProtocolConfigMap["proxy_protocol_config_vpc_id_enabled"] = proxyProtocolV2Config2Raw["Ppv2VpcIdEnabled"]

		proxyProtocolConfigMaps = append(proxyProtocolConfigMaps, proxyProtocolConfigMap)
	}
	if objectRaw["ProxyProtocolV2Config"] != nil {
		if err := d.Set("proxy_protocol_config", proxyProtocolConfigMaps); err != nil {
			return err
		}
	}
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

	if d.HasChange("status") {
		nlbServiceV2 := NlbServiceV2{client}
		object, err := nlbServiceV2.DescribeNlbListener(d.Id())
		if err != nil {
			return WrapError(err)
		}

		target := d.Get("status").(string)
		if object["ListenerStatus"].(string) != target {
			if target == "Running" {
				action := "StartListener"
				request = make(map[string]interface{})
				query = make(map[string]interface{})
				request["ListenerId"] = d.Id()
				request["RegionId"] = client.RegionId
				request["ClientToken"] = buildClientToken(action)
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
				stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, nlbServiceV2.NlbListenerStateRefreshFunc(d.Id(), "ListenerStatus", []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}

			}
			if target == "Stopped" {
				action := "StopListener"
				request = make(map[string]interface{})
				query = make(map[string]interface{})
				request["ListenerId"] = d.Id()
				request["RegionId"] = client.RegionId
				request["ClientToken"] = buildClientToken(action)
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
				stateConf := BuildStateConf([]string{}, []string{"Stopped"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, nlbServiceV2.NlbListenerStateRefreshFunc(d.Id(), "ListenerStatus", []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}

			}
		}
	}

	action := "UpdateListenerAttribute"
	var err error
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
		if v, ok := d.GetOk("certificate_ids"); ok || d.HasChange("certificate_ids") {
			certificateIdsMapsArray := v.([]interface{})
			request["CertificateIds"] = certificateIdsMapsArray
		}
	}

	if !d.IsNewResource() && d.HasChange("ca_certificate_ids") {
		update = true
		if v, ok := d.GetOk("ca_certificate_ids"); ok || d.HasChange("ca_certificate_ids") {
			caCertificateIdsMapsArray := v.([]interface{})
			request["CaCertificateIds"] = caCertificateIdsMapsArray
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

	if !d.IsNewResource() && d.HasChange("proxy_protocol_config") {
		update = true
		objectDataLocalMap := make(map[string]interface{})

		if v := d.Get("proxy_protocol_config"); v != nil {
			proxyProtocolConfigVpcIdEnabled, _ := jsonpath.Get("$[0].proxy_protocol_config_vpc_id_enabled", v)
			if proxyProtocolConfigVpcIdEnabled != nil && (d.HasChange("proxy_protocol_config.0.proxy_protocol_config_vpc_id_enabled") || proxyProtocolConfigVpcIdEnabled != "") {
				objectDataLocalMap["Ppv2VpcIdEnabled"] = proxyProtocolConfigVpcIdEnabled
			}
			proxyProtocolConfigPrivateLinkEpIdEnabled, _ := jsonpath.Get("$[0].proxy_protocol_config_private_link_ep_id_enabled", v)
			if proxyProtocolConfigPrivateLinkEpIdEnabled != nil && (d.HasChange("proxy_protocol_config.0.proxy_protocol_config_private_link_ep_id_enabled") || proxyProtocolConfigPrivateLinkEpIdEnabled != "") {
				objectDataLocalMap["Ppv2PrivateLinkEpIdEnabled"] = proxyProtocolConfigPrivateLinkEpIdEnabled
			}
			proxyProtocolConfigPrivateLinkEpsIdEnabled, _ := jsonpath.Get("$[0].proxy_protocol_config_private_link_eps_id_enabled", v)
			if proxyProtocolConfigPrivateLinkEpsIdEnabled != nil && (d.HasChange("proxy_protocol_config.0.proxy_protocol_config_private_link_eps_id_enabled") || proxyProtocolConfigPrivateLinkEpsIdEnabled != "") {
				objectDataLocalMap["Ppv2PrivateLinkEpsIdEnabled"] = proxyProtocolConfigPrivateLinkEpsIdEnabled
			}

			objectDataLocalMapJson, err := json.Marshal(objectDataLocalMap)
			if err != nil {
				return WrapError(err)
			}
			request["ProxyProtocolV2Config"] = string(objectDataLocalMapJson)
		}
	}

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
		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 15*time.Second, nlbServiceV2.NlbListenerStateRefreshFunc(d.Id(), "ListenerStatus", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	if d.HasChange("tags") {
		nlbServiceV2 := NlbServiceV2{client}
		if err := nlbServiceV2.SetResourceTags(d, "listener"); err != nil {
			return WrapError(err)
		}
	}
	return resourceAliCloudNlbListenerRead(d, meta)
}

func resourceAliCloudNlbListenerDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteListener"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["ListenerId"] = d.Id()
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
	stateConf := BuildStateConf([]string{}, []string{"Succeeded"}, d.Timeout(schema.TimeoutDelete), 30*time.Second, nlbServiceV2.DescribeAsyncNlbListenerStateRefreshFunc(d, response, "$.Status", []string{}))
	if jobDetail, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id(), jobDetail)
	}

	return nil
}
