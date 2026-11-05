// Package ens. This file is generated automatically. Please do not modify it manually, thank you!
package ens

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/validators"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceEnsLoadBalancerUdpListener() *schema.Resource {
	return &schema.Resource{
		Create: resourceEnsLoadBalancerUdpListenerCreate,
		Read:   resourceEnsLoadBalancerUdpListenerRead,
		Update: resourceEnsLoadBalancerUdpListenerUpdate,
		Delete: resourceEnsLoadBalancerUdpListenerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"backend_server_port": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validators.IntBetween(0, 65535),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"eip_transmit": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"established_timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validators.IntBetween(10, 900),
			},
			"health_check_connect_port": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validators.IntBetween(0, 65535),
			},
			"health_check_connect_timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validators.IntBetween(0, 300),
			},
			"health_check_exp": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"health_check_interval": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validators.IntBetween(0, 50),
			},
			"health_check_req": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"healthy_threshold": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validators.IntBetween(2, 10),
			},
			"listener_port": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"load_balancer_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"scheduler": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validators.StringInSlice([]string{"wrr", "wlc", "rr", "sch", "qch", "iqch", "tch"}, false),
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"unhealthy_threshold": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validators.IntBetween(2, 10),
			},
		},
	}
}

func resourceEnsLoadBalancerUdpListenerCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateLoadBalancerUDPListener"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("load_balancer_id"); ok {
		request["LoadBalancerId"] = v
	}
	if v, ok := d.GetOk("listener_port"); ok {
		request["ListenerPort"] = v
	}

	if v, ok := d.GetOkExists("health_check_connect_port"); ok {
		request["HealthCheckConnectPort"] = v
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	if v, ok := d.GetOkExists("established_timeout"); ok && v.(int) > 0 {
		request["EstablishedTimeout"] = v
	}
	if v, ok := d.GetOkExists("health_check_connect_timeout"); ok && v.(int) > 0 {
		request["HealthCheckConnectTimeout"] = v
	}
	if v, ok := d.GetOk("scheduler"); ok {
		request["Scheduler"] = v
	}
	if v, ok := d.GetOk("health_check_req"); ok {
		request["HealthCheckReq"] = v
	}
	if v, ok := d.GetOk("health_check_exp"); ok {
		request["HealthCheckExp"] = v
	}
	if v, ok := d.GetOk("eip_transmit"); ok {
		request["EipTransmit"] = v
	}
	if v, ok := d.GetOkExists("health_check_interval"); ok && v.(int) > 0 {
		request["HealthCheckInterval"] = v
	}
	if v, ok := d.GetOkExists("backend_server_port"); ok {
		request["BackendServerPort"] = v
	}
	if v, ok := d.GetOkExists("unhealthy_threshold"); ok && v.(int) > 0 {
		request["UnhealthyThreshold"] = v
	}
	if v, ok := d.GetOkExists("healthy_threshold"); ok && v.(int) > 0 {
		request["HealthyThreshold"] = v
	}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("Ens", "2017-11-10", action, query, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_ens_load_balancer_udp_listener", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", request["LoadBalancerId"], request["ListenerPort"]))

	ensServiceV2 := NewService(client)
	stateConf := helper.BuildStateConf([]string{}, []string{"Stopped"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, ensServiceV2.EnsLoadBalancerUdpListenerStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return resourceEnsLoadBalancerUdpListenerUpdate(d, meta)
}

func resourceEnsLoadBalancerUdpListenerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ensServiceV2 := NewService(client)

	objectRaw, err := ensServiceV2.DescribeEnsLoadBalancerUdpListener(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ens_load_balancer_udp_listener DescribeEnsLoadBalancerUdpListener Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("backend_server_port", objectRaw["BackendServerPort"])
	d.Set("description", objectRaw["Description"])
	d.Set("eip_transmit", objectRaw["EipTransmit"])
	d.Set("established_timeout", objectRaw["EstablishedTimeout"])
	d.Set("health_check_connect_port", objectRaw["HealthCheckConnectPort"])
	d.Set("health_check_connect_timeout", objectRaw["HealthCheckConnectTimeout"])
	d.Set("health_check_exp", objectRaw["HealthCheckExp"])
	d.Set("health_check_interval", objectRaw["HealthCheckInterval"])
	d.Set("health_check_req", objectRaw["HealthCheckReq"])
	d.Set("healthy_threshold", objectRaw["HealthyThreshold"])
	d.Set("scheduler", objectRaw["Scheduler"])
	d.Set("status", objectRaw["Status"])
	d.Set("unhealthy_threshold", objectRaw["UnhealthyThreshold"])
	d.Set("listener_port", objectRaw["ListenerPort"])

	parts := strings.Split(d.Id(), ":")
	d.Set("load_balancer_id", parts[0])

	return nil
}

func resourceEnsLoadBalancerUdpListenerUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	ensServiceV2 := NewService(client)
	objectRaw, _ := ensServiceV2.DescribeEnsLoadBalancerUdpListener(d.Id())

	if d.HasChange("status") {
		var err error
		target := d.Get("status").(string)

		currentStatus, err := jsonpath.Get("Status", objectRaw)
		if err != nil {
			return errs.WrapErrorf(err, errs.FailedGetAttributeMsg, d.Id(), "Status", objectRaw)
		}
		if fmt.Sprint(currentStatus) != target {
			if target == "Running" {
				parts := strings.Split(d.Id(), ":")
				action := "StartLoadBalancerListener"
				request = make(map[string]interface{})
				query = make(map[string]interface{})
				request["LoadBalancerId"] = parts[0]
				request["ListenerPort"] = parts[1]

				request["ListenerProtocol"] = "udp"
				wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
				err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
					response, err = client.RpcPost("Ens", "2017-11-10", action, query, request, true)
					if err != nil {
						if helper.NeedRetry(err) {
							wait()
							return retry.RetryableError(err)
						}
						return retry.NonRetryableError(err)
					}
					return nil
				})
				logs.AddDebug(action, response, request)
				if err != nil {
					return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
				}
				ensServiceV2 := NewService(client)
				stateConf := helper.BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, ensServiceV2.EnsLoadBalancerUdpListenerStateRefreshFunc(d.Id(), "Status", []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return errs.WrapErrorf(err, errs.IdMsg, d.Id())
				}

			}
			if target == "Stopped" {
				parts := strings.Split(d.Id(), ":")
				action := "StopLoadBalancerListener"
				request = make(map[string]interface{})
				query = make(map[string]interface{})
				request["LoadBalancerId"] = parts[0]
				request["ListenerPort"] = parts[1]

				request["ListenerProtocol"] = "udp"
				wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
				err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
					response, err = client.RpcPost("Ens", "2017-11-10", action, query, request, true)
					if err != nil {
						if helper.NeedRetry(err) {
							wait()
							return retry.RetryableError(err)
						}
						return retry.NonRetryableError(err)
					}
					return nil
				})
				logs.AddDebug(action, response, request)
				if err != nil {
					return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
				}
				ensServiceV2 := NewService(client)
				stateConf := helper.BuildStateConf([]string{}, []string{"Stopped"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, ensServiceV2.EnsLoadBalancerUdpListenerStateRefreshFunc(d.Id(), "Status", []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return errs.WrapErrorf(err, errs.IdMsg, d.Id())
				}

			}
		}
	}

	var err error
	parts := strings.Split(d.Id(), ":")
	action := "SetLoadBalancerUDPListenerAttribute"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["LoadBalancerId"] = parts[0]
	request["ListenerPort"] = parts[1]

	if !d.IsNewResource() && d.HasChange("health_check_connect_port") {
		update = true
	}
	if v, ok := d.GetOkExists("health_check_connect_port"); ok || d.HasChange("health_check_connect_port") {
		request["HealthCheckConnectPort"] = v
	}
	if !d.IsNewResource() && d.HasChange("description") {
		update = true
	}
	if v, ok := d.GetOk("description"); ok || d.HasChange("description") {
		request["Description"] = v
	}
	if !d.IsNewResource() && d.HasChange("established_timeout") {
		update = true
	}
	if v, ok := d.GetOkExists("established_timeout"); (ok || d.HasChange("established_timeout")) && v.(int) > 0 {
		request["EstablishedTimeout"] = v
	}
	if !d.IsNewResource() && d.HasChange("health_check_connect_timeout") {
		update = true
	}
	if v, ok := d.GetOkExists("health_check_connect_timeout"); (ok || d.HasChange("health_check_connect_timeout")) && v.(int) > 0 {
		request["HealthCheckConnectTimeout"] = v
	}
	if !d.IsNewResource() && d.HasChange("scheduler") {
		update = true
	}
	if v, ok := d.GetOk("scheduler"); ok || d.HasChange("scheduler") {
		request["Scheduler"] = v
	}
	if !d.IsNewResource() && d.HasChange("health_check_req") {
		update = true
	}
	if v, ok := d.GetOk("health_check_req"); ok || d.HasChange("health_check_req") {
		request["HealthCheckReq"] = v
	}
	if !d.IsNewResource() && d.HasChange("health_check_exp") {
		update = true
	}
	if v, ok := d.GetOk("health_check_exp"); ok || d.HasChange("health_check_exp") {
		request["HealthCheckExp"] = v
	}
	if !d.IsNewResource() && d.HasChange("eip_transmit") {
		update = true
	}
	if v, ok := d.GetOk("eip_transmit"); ok || d.HasChange("eip_transmit") {
		request["EipTransmit"] = v
	}
	if !d.IsNewResource() && d.HasChange("health_check_interval") {
		update = true
	}
	if v, ok := d.GetOkExists("health_check_interval"); (ok || d.HasChange("health_check_interval")) && v.(int) > 0 {
		request["HealthCheckInterval"] = v
	}
	if !d.IsNewResource() && d.HasChange("unhealthy_threshold") {
		update = true
	}
	if v, ok := d.GetOkExists("unhealthy_threshold"); (ok || d.HasChange("unhealthy_threshold")) && v.(int) > 0 {
		request["UnhealthyThreshold"] = v
	}
	if !d.IsNewResource() && d.HasChange("healthy_threshold") {
		update = true
	}
	if v, ok := d.GetOkExists("healthy_threshold"); (ok || d.HasChange("healthy_threshold")) && v.(int) > 0 {
		request["HealthyThreshold"] = v
	}
	if update {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("Ens", "2017-11-10", action, query, request, true)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
	}

	return resourceEnsLoadBalancerUdpListenerRead(d, meta)
}

func resourceEnsLoadBalancerUdpListenerDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeleteLoadBalancerListener"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["LoadBalancerId"] = parts[0]
	request["ListenerPort"] = parts[1]

	request["ListenerProtocol"] = "udp"
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("Ens", "2017-11-10", action, query, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		if errs.IsExpectedErrors(err, []string{"ListenerNotFound", "LoadBalancerNotFound"}) || errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	ensServiceV2 := NewService(client)
	stateConf := helper.BuildStateConf([]string{}, []string{""}, d.Timeout(schema.TimeoutDelete), 5*time.Second, ensServiceV2.EnsLoadBalancerUdpListenerStateRefreshFunc(d.Id(), "ListenerPort", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return nil
}
