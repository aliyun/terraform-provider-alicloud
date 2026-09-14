// Package alicloud. Resource for ENS Load Balancer TCP Listener.
package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudEnsLoadBalancerTcpListener() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudEnsLoadBalancerTcpListenerCreate,
		Read:   resourceAliCloudEnsLoadBalancerTcpListenerRead,
		Update: resourceAliCloudEnsLoadBalancerTcpListenerUpdate,
		Delete: resourceAliCloudEnsLoadBalancerTcpListenerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"load_balancer_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"listener_port": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"backend_server_port": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"scheduler": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"wrr", "wlc", "rr", "sch", "qch", "iqch"}, false),
			},
			"persistence_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"established_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"healthy_threshold": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"unhealthy_threshold": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"health_check_connect_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"health_check_interval": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"health_check_connect_port": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"health_check_domain": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"health_check_http_code": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"http_2xx", "http_3xx", "http_4xx", "http_5xx", "http_2xx,http_3xx", "http_2xx,http_3xx,http_4xx", "http_2xx,http_3xx,http_4xx,http_5xx", "http_3xx,http_4xx", "http_3xx,http_4xx,http_5xx", "http_4xx,http_5xx"}, false),
			},
			"health_check_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"tcp", "http"}, false),
			},
			"health_check_uri": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"eip_transmit": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"on", "off"}, false),
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"Running", "Stopped"}, false),
			},
			"protocol": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudEnsLoadBalancerTcpListenerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "CreateLoadBalancerTCPListener"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	request["ClientToken"] = buildClientToken(action)
	request["LoadBalancerId"] = d.Get("load_balancer_id")
	request["ListenerPort"] = d.Get("listener_port")
	if v, ok := d.GetOk("backend_server_port"); ok {
		request["BackendServerPort"] = v
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	if v, ok := d.GetOk("scheduler"); ok {
		request["Scheduler"] = v
	}
	if v, ok := d.GetOk("persistence_timeout"); ok {
		request["PersistenceTimeout"] = v
	}
	if v, ok := d.GetOk("established_timeout"); ok {
		request["EstablishedTimeout"] = v
	}
	if v, ok := d.GetOk("healthy_threshold"); ok {
		request["HealthyThreshold"] = v
	}
	if v, ok := d.GetOk("unhealthy_threshold"); ok {
		request["UnhealthyThreshold"] = v
	}
	if v, ok := d.GetOk("health_check_connect_timeout"); ok {
		request["HealthCheckConnectTimeout"] = v
	}
	if v, ok := d.GetOk("health_check_interval"); ok {
		request["HealthCheckInterval"] = v
	}
	if v, ok := d.GetOk("health_check_connect_port"); ok {
		request["HealthCheckConnectPort"] = v
	}
	if v, ok := d.GetOk("health_check_domain"); ok {
		request["HealthCheckDomain"] = v
	}
	if v, ok := d.GetOk("health_check_http_code"); ok {
		request["HealthCheckHttpCode"] = v
	}
	if v, ok := d.GetOk("health_check_type"); ok {
		request["HealthCheckType"] = v
	}
	if v, ok := d.GetOk("health_check_uri"); ok {
		request["HealthCheckURI"] = v
	}
	if v, ok := d.GetOk("eip_transmit"); ok {
		request["EipTransmit"] = v
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Ens", "2017-11-10", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ens_load_balancer_tcp_listener", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["LoadBalancerId"], ":", request["ListenerPort"]))

	ensServiceV2 := EnsServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"Running", "Stopped"}, d.Timeout(schema.TimeoutCreate), 10*time.Second, ensServiceV2.EnsLoadBalancerTcpListenerStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	// If the user requested the listener to be Running, start it (default ENS
	// creates a listener in Stopped state and the API has no create-time
	// status flag, so an explicit Running must be applied post-create).
	if d.Get("status").(string) == "Running" {
		if err := ensLoadBalancerTcpListenerSetStatus(client, d, "Running", d.Timeout(schema.TimeoutCreate)); err != nil {
			return err
		}
	}

	return resourceAliCloudEnsLoadBalancerTcpListenerRead(d, meta)
}

func resourceAliCloudEnsLoadBalancerTcpListenerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ensServiceV2 := EnsServiceV2{client}

	objectRaw, err := ensServiceV2.DescribeEnsLoadBalancerTcpListener(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ens_load_balancer_tcp_listener DescribeEnsLoadBalancerTcpListener Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	// load_balancer_id is not returned by DescribeLoadBalancerTCPListenerAttribute;
	// parse it from the resource ID so both refresh and import populate state.
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	d.Set("load_balancer_id", parts[0])
	d.Set("listener_port", formatInt(objectRaw["ListenerPort"]))
	d.Set("backend_server_port", formatInt(objectRaw["BackendServerPort"]))
	d.Set("description", objectRaw["Description"])
	d.Set("scheduler", objectRaw["Scheduler"])
	d.Set("persistence_timeout", formatInt(objectRaw["PersistenceTimeout"]))
	d.Set("established_timeout", formatInt(objectRaw["EstablishedTimeout"]))
	d.Set("healthy_threshold", formatInt(objectRaw["HealthyThreshold"]))
	d.Set("unhealthy_threshold", formatInt(objectRaw["UnhealthyThreshold"]))
	d.Set("health_check_connect_timeout", formatInt(objectRaw["HealthCheckConnectTimeout"]))
	d.Set("health_check_interval", formatInt(objectRaw["HealthCheckInterval"]))
	d.Set("health_check_connect_port", formatInt(objectRaw["HealthCheckConnectPort"]))
	d.Set("health_check_type", objectRaw["HealthCheckType"])
	d.Set("eip_transmit", objectRaw["EipTransmit"])
	d.Set("status", objectRaw["Status"])
	// protocol is always "tcp" for this TCP listener resource; the API
	// does not return it in the Describe response.
	d.Set("protocol", "tcp")
	// health_check_domain/http_code/uri are not returned by the Describe
	// API for TCP listeners; preserve the configured value unless the
	// API provides one (e.g. when health_check_type is "http").
	if v, ok := objectRaw["HealthCheckDomain"]; ok {
		d.Set("health_check_domain", v)
	}
	if v, ok := objectRaw["HealthCheckHttpCode"]; ok {
		d.Set("health_check_http_code", v)
	}
	if v, ok := objectRaw["HealthCheckURI"]; ok {
		d.Set("health_check_uri", v)
	}
	return nil
}

func resourceAliCloudEnsLoadBalancerTcpListenerUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	loadBalancerId := parts[0]
	listenerPort := parts[1]

	// 1. Attribute changes go through SetLoadBalancerTCPListenerAttribute.
	update := false
	request := make(map[string]interface{})
	query := make(map[string]interface{})
	request["LoadBalancerId"] = loadBalancerId
	request["ListenerPort"] = listenerPort
	if !d.IsNewResource() && d.HasChange("description") {
		update = true
		request["Description"] = d.Get("description")
	}
	if !d.IsNewResource() && d.HasChange("scheduler") {
		update = true
		request["Scheduler"] = d.Get("scheduler")
	}
	if !d.IsNewResource() && d.HasChange("persistence_timeout") {
		update = true
		request["PersistenceTimeout"] = d.Get("persistence_timeout")
	}
	if !d.IsNewResource() && d.HasChange("established_timeout") {
		update = true
		request["EstablishedTimeout"] = d.Get("established_timeout")
	}
	if !d.IsNewResource() && d.HasChange("healthy_threshold") {
		update = true
		request["HealthyThreshold"] = d.Get("healthy_threshold")
	}
	if !d.IsNewResource() && d.HasChange("unhealthy_threshold") {
		update = true
		request["UnhealthyThreshold"] = d.Get("unhealthy_threshold")
	}
	if !d.IsNewResource() && d.HasChange("health_check_connect_timeout") {
		update = true
		request["HealthCheckConnectTimeout"] = d.Get("health_check_connect_timeout")
	}
	if !d.IsNewResource() && d.HasChange("health_check_interval") {
		update = true
		request["HealthCheckInterval"] = d.Get("health_check_interval")
	}
	if !d.IsNewResource() && d.HasChange("health_check_connect_port") {
		update = true
		request["HealthCheckConnectPort"] = d.Get("health_check_connect_port")
	}
	if !d.IsNewResource() && d.HasChange("health_check_domain") {
		update = true
		request["HealthCheckDomain"] = d.Get("health_check_domain")
	}
	if !d.IsNewResource() && d.HasChange("health_check_http_code") {
		update = true
		request["HealthCheckHttpCode"] = d.Get("health_check_http_code")
	}
	if !d.IsNewResource() && d.HasChange("health_check_type") {
		update = true
		request["HealthCheckType"] = d.Get("health_check_type")
	}
	if !d.IsNewResource() && d.HasChange("health_check_uri") {
		update = true
		request["HealthCheckURI"] = d.Get("health_check_uri")
	}
	if !d.IsNewResource() && d.HasChange("eip_transmit") {
		update = true
		request["EipTransmit"] = d.Get("eip_transmit")
	}
	if update {
		action := "SetLoadBalancerTCPListenerAttribute"
		wait := incrementalWait(3*time.Second, 5*time.Second)
		var response map[string]interface{}
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Ens", "2017-11-10", action, query, request, true)
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

	// 2. Status changes go through Start/StopLoadBalancerListener.
	if d.HasChange("status") {
		targetStatus := d.Get("status").(string)
		if err := ensLoadBalancerTcpListenerSetStatus(client, d, targetStatus, d.Timeout(schema.TimeoutUpdate)); err != nil {
			return err
		}
	}

	return resourceAliCloudEnsLoadBalancerTcpListenerRead(d, meta)
}

func resourceAliCloudEnsLoadBalancerTcpListenerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	loadBalancerId := parts[0]
	listenerPort := parts[1]

	action := "DeleteLoadBalancerListener"
	query := make(map[string]interface{})
	request := make(map[string]interface{})
	request["LoadBalancerId"] = loadBalancerId
	request["ListenerPort"] = listenerPort

	wait := incrementalWait(3*time.Second, 5*time.Second)
	var response map[string]interface{}
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Ens", "2017-11-10", action, query, request, true)
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
	return nil
}

// ensLoadBalancerTcpListenerSetStatus starts or stops an ENS Load Balancer TCP
// listener. ENS listener lifecycle (Start/Stop) is shared across TCP/HTTP/UDP
// listener types and keyed by LoadBalancerId + ListenerPort + ListenerProtocol.
// The protocol is derived from the listener's computed protocol value.
func ensLoadBalancerTcpListenerSetStatus(client *connectivity.AliyunClient, d *schema.ResourceData, targetStatus string, timeout time.Duration) error {
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	loadBalancerId := parts[0]
	listenerPort := parts[1]

	// Read the listener protocol to pass it to Start/StopLoadBalancerListener.
	ensServiceV2 := EnsServiceV2{client}
	object, err := ensServiceV2.DescribeEnsLoadBalancerTcpListener(d.Id())
	if err != nil {
		return WrapError(err)
	}
	protocol, _ := object["Protocol"].(string)
	if protocol == "" {
		protocol = "tcp"
	}

	var action string
	switch targetStatus {
	case "Running":
		action = "StartLoadBalancerListener"
	case "Stopped":
		action = "StopLoadBalancerListener"
	default:
		return nil
	}

	query := make(map[string]interface{})
	request := make(map[string]interface{})
	request["LoadBalancerId"] = loadBalancerId
	request["ListenerPort"] = listenerPort
	request["ListenerProtocol"] = protocol

	wait := incrementalWait(3*time.Second, 5*time.Second)
	var response map[string]interface{}
	err = resource.Retry(timeout, func() *resource.RetryError {
		response, err = client.RpcPost("Ens", "2017-11-10", action, query, request, true)
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

	target := []string{"Running"}
	if targetStatus == "Stopped" {
		target = []string{"Stopped"}
	}
	stateConf := BuildStateConf([]string{}, target, timeout, 10*time.Second, ensServiceV2.EnsLoadBalancerTcpListenerStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
