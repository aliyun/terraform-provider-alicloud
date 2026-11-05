package ens

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/validators"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceEnsLoadBalancerHttpListener() *schema.Resource {
	return &schema.Resource{
		Create: resourceEnsLoadBalancerHttpListenerCreate,
		Read:   resourceEnsLoadBalancerHttpListenerRead,
		Update: resourceEnsLoadBalancerHttpListenerUpdate,
		Delete: resourceEnsLoadBalancerHttpListenerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(15 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(15 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"load_balancer_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"listener_port": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validators.IntBetween(1, 65535),
			},
			"backend_server_port": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validators.IntBetween(1, 65535),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"scheduler": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "wrr",
				ValidateFunc: validators.StringInSlice([]string{"wrr", "wlc", "rr", "sch", "qch", "iqch"}, false),
			},
			"health_check": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "off",
				ValidateFunc: validators.StringInSlice([]string{"on", "off"}, false),
			},
			"health_check_domain": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: ensLoadBalancerHttpListenerHealthCheckOffDiffSuppressFunc,
			},
			"health_check_uri": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: ensLoadBalancerHttpListenerHealthCheckOffDiffSuppressFunc,
			},
			"healthy_threshold": {
				Type:             schema.TypeInt,
				Optional:         true,
				Default:          3,
				ValidateFunc:     validators.IntBetween(2, 10),
				DiffSuppressFunc: ensLoadBalancerHttpListenerHealthCheckOffDiffSuppressFunc,
			},
			"unhealthy_threshold": {
				Type:             schema.TypeInt,
				Optional:         true,
				Default:          3,
				ValidateFunc:     validators.IntBetween(2, 10),
				DiffSuppressFunc: ensLoadBalancerHttpListenerHealthCheckOffDiffSuppressFunc,
			},
			"health_check_timeout": {
				Type:             schema.TypeInt,
				Optional:         true,
				Default:          5,
				ValidateFunc:     validators.IntBetween(1, 300),
				DiffSuppressFunc: ensLoadBalancerHttpListenerHealthCheckOffDiffSuppressFunc,
			},
			"health_check_connect_port": {
				Type:             schema.TypeInt,
				Optional:         true,
				Computed:         true,
				ValidateFunc:     validators.IntBetween(1, 65535),
				DiffSuppressFunc: ensLoadBalancerHttpListenerHealthCheckOffDiffSuppressFunc,
			},
			"health_check_interval": {
				Type:             schema.TypeInt,
				Optional:         true,
				Default:          2,
				ValidateFunc:     validators.IntBetween(1, 50),
				DiffSuppressFunc: ensLoadBalancerHttpListenerHealthCheckOffDiffSuppressFunc,
			},
			"health_check_http_code": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "http_2xx",
				ValidateFunc:     validators.StringInSlice([]string{"http_2xx", "http_3xx", "http_4xx", "http_5xx"}, false),
				DiffSuppressFunc: ensLoadBalancerHttpListenerHealthCheckOffDiffSuppressFunc,
			},
			"health_check_method": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "head",
				ValidateFunc:     validators.StringInSlice([]string{"head", "get"}, false),
				DiffSuppressFunc: ensLoadBalancerHttpListenerHealthCheckOffDiffSuppressFunc,
			},
			"idle_timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      15,
				ValidateFunc: validators.IntBetween(1, 60),
			},
			"request_timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      60,
				ValidateFunc: validators.IntBetween(1, 180),
			},
			"x_forwarded_for": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "off",
				ValidateFunc: validators.StringInSlice([]string{"on", "off"}, false),
			},
			"listener_forward": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "off",
				ValidateFunc: validators.StringInSlice([]string{"on", "off"}, false),
			},
			"forward_port": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validators.StringInSlice([]string{"running", "stopped"}, false),
			},
		},
	}
}

func resourceEnsLoadBalancerHttpListenerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "CreateLoadBalancerHTTPListener"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	request["LoadBalancerId"] = d.Get("load_balancer_id")
	request["ListenerPort"] = d.Get("listener_port")
	if v, ok := d.GetOkExists("backend_server_port"); ok {
		request["BackendServerPort"] = v
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	request["Scheduler"] = d.Get("scheduler")
	request["HealthCheck"] = d.Get("health_check")
	if d.Get("health_check").(string) == "on" {
		if v, ok := d.GetOk("health_check_domain"); ok {
			request["HealthCheckDomain"] = v
		}
		if v, ok := d.GetOk("health_check_uri"); ok {
			request["HealthCheckURI"] = v
		}
		request["HealthyThreshold"] = d.Get("healthy_threshold")
		request["UnhealthyThreshold"] = d.Get("unhealthy_threshold")
		request["HealthCheckTimeout"] = d.Get("health_check_timeout")
		if v, ok := d.GetOkExists("health_check_connect_port"); ok {
			request["HealthCheckConnectPort"] = v
		}
		request["HealthCheckInterval"] = d.Get("health_check_interval")
		request["HealthCheckHttpCode"] = d.Get("health_check_http_code")
		request["HealthCheckMethod"] = d.Get("health_check_method")
	}
	request["IdleTimeout"] = d.Get("idle_timeout")
	request["RequestTimeout"] = d.Get("request_timeout")
	request["XForwardedFor"] = d.Get("x_forwarded_for")
	request["ListenerForward"] = d.Get("listener_forward")
	if v, ok := d.GetOkExists("forward_port"); ok {
		request["ForwardPort"] = v
	}

	wait := helper.IncrementalWait(10*time.Second, 15*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("Ens", "2017-11-10", action, query, request, true)
		if err != nil {
			// A same-port ForceNew destroy (backend_server_port is ForceNew
			// while listener_port is not) may leave the listener port still
			// releasing on the ENS backend; absorb the consistency window.
			if errs.IsExpectedErrors(err, []string{"ListenerAlreadyExists"}) {
				wait()
				return retry.RetryableError(err)
			}
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_ens_load_balancer_http_listener", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", d.Get("load_balancer_id"), d.Get("listener_port")))

	ensServiceV2 := NewService(client)
	stateConf := helper.BuildStateConf([]string{}, []string{"stopped"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, ensServiceV2.EnsLoadBalancerHttpListenerStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	// Start the listener when the desired status is running.
	if d.Get("status").(string) == "running" {
		if err := ensServiceV2.StartEnsLoadBalancerHttpListener(d.Id()); err != nil {
			return errs.WrapError(err)
		}
		startStateConf := helper.BuildStateConf([]string{}, []string{"running"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, ensServiceV2.EnsLoadBalancerHttpListenerStateRefreshFunc(d.Id(), "Status", []string{}))
		if _, err := startStateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
		}
	}

	return resourceEnsLoadBalancerHttpListenerRead(d, meta)
}

func resourceEnsLoadBalancerHttpListenerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ensServiceV2 := NewService(client)

	objectRaw, err := ensServiceV2.DescribeEnsLoadBalancerHttpListener(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ens_load_balancer_http_listener DescribeEnsLoadBalancerHttpListener Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	parts := strings.Split(d.Id(), ":")
	if len(parts) == 2 {
		d.Set("load_balancer_id", parts[0])
		listenerPort, _ := strconv.Atoi(parts[1])
		d.Set("listener_port", listenerPort)
	}

	vVal, err := convert.FormatInt(objectRaw["BackendServerPort"])
	if err != nil {
		return err
	}
	if v := vVal; v > 0 {
		d.Set("backend_server_port", v)
	}
	d.Set("description", objectRaw["Description"])
	d.Set("scheduler", objectRaw["Scheduler"])
	d.Set("health_check", objectRaw["HealthCheck"])
	// When health check is off, the ENS Describe API omits the health check
	// parameters (HealthCheckConnectPort is still returned but is meaningless),
	// so they must not be set into state. DiffSuppressFunc keeps them out of
	// the diff; Read keeps them out of state, so they stay absent and do not
	// surface as unexpected attributes after refresh.
	if objectRaw["HealthCheck"] == "on" {
		d.Set("health_check_domain", objectRaw["HealthCheckDomain"])
		d.Set("health_check_uri", objectRaw["HealthCheckURI"])
		healthy_thresholdVal, err := convert.FormatInt(objectRaw["HealthyThreshold"])
		if err != nil {
			return err
		}
		d.Set("healthy_threshold", healthy_thresholdVal)
		unhealthy_thresholdVal, err := convert.FormatInt(objectRaw["UnhealthyThreshold"])
		if err != nil {
			return err
		}
		d.Set("unhealthy_threshold", unhealthy_thresholdVal)
		health_check_timeoutVal, err := convert.FormatInt(objectRaw["HealthCheckTimeout"])
		if err != nil {
			return err
		}
		d.Set("health_check_timeout", health_check_timeoutVal)
		healthCheckConnectPort, err := convert.FormatInt(objectRaw["HealthCheckConnectPort"])
		if err != nil {
			return err
		}
		if v := healthCheckConnectPort; v > 0 {
			d.Set("health_check_connect_port", v)
		}
		health_check_intervalVal, err := convert.FormatInt(objectRaw["HealthCheckInterval"])
		if err != nil {
			return err
		}
		d.Set("health_check_interval", health_check_intervalVal)
		d.Set("health_check_http_code", objectRaw["HealthCheckHttpCode"])
		d.Set("health_check_method", objectRaw["HealthCheckMethod"])
	}
	idle_timeoutVal, err := convert.FormatInt(objectRaw["IdleTimeout"])
	if err != nil {
		return err
	}
	d.Set("idle_timeout", idle_timeoutVal)
	request_timeoutVal, err := convert.FormatInt(objectRaw["RequestTimeout"])
	if err != nil {
		return err
	}
	d.Set("request_timeout", request_timeoutVal)
	d.Set("x_forwarded_for", objectRaw["XForwardedFor"])
	d.Set("listener_forward", objectRaw["ListenerForward"])
	forwardPort, err := convert.FormatInt(objectRaw["ForwardPort"])
	if err != nil {
		return err
	}
	if v := forwardPort; v > 0 {
		d.Set("forward_port", v)
	}
	d.Set("status", objectRaw["Status"])

	return nil
}

func resourceEnsLoadBalancerHttpListenerUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ensServiceV2 := NewService(client)

	attributeChanged := d.HasChanges([]string{
		"description", "scheduler", "health_check", "health_check_domain",
		"health_check_uri", "healthy_threshold", "unhealthy_threshold",
		"health_check_timeout", "health_check_connect_port", "health_check_interval",
		"health_check_http_code", "health_check_method", "idle_timeout",
		"request_timeout", "x_forwarded_for",
	}...)
	if attributeChanged {
		action := "SetLoadBalancerHTTPListenerAttribute"
		var request map[string]interface{}
		var response map[string]interface{}
		query := make(map[string]interface{})
		var err error
		request = make(map[string]interface{})

		parts := strings.Split(d.Id(), ":")
		request["LoadBalancerId"] = parts[0]
		request["ListenerPort"] = parts[1]

		request["Description"] = d.Get("description")
		request["Scheduler"] = d.Get("scheduler")
		request["HealthCheck"] = d.Get("health_check")
		if d.Get("health_check").(string) == "on" {
			if v, ok := d.GetOk("health_check_domain"); ok {
				request["HealthCheckDomain"] = v
			}
			if v, ok := d.GetOk("health_check_uri"); ok {
				request["HealthCheckURI"] = v
			}
			request["HealthyThreshold"] = d.Get("healthy_threshold")
			request["UnhealthyThreshold"] = d.Get("unhealthy_threshold")
			request["HealthCheckTimeout"] = d.Get("health_check_timeout")
			if v, ok := d.GetOkExists("health_check_connect_port"); ok {
				request["HealthCheckConnectPort"] = v
			}
			request["HealthCheckInterval"] = d.Get("health_check_interval")
			request["HealthCheckHttpCode"] = d.Get("health_check_http_code")
			request["HealthCheckMethod"] = d.Get("health_check_method")
		}
		request["IdleTimeout"] = d.Get("idle_timeout")
		request["RequestTimeout"] = d.Get("request_timeout")
		request["XForwardedFor"] = d.Get("x_forwarded_for")

		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("Ens", "2017-11-10", action, query, request, true)
			if err != nil {
				// The Set API requires the listener to be in the stopped state.
				// When the listener is running, stop it and retry the modification.
				if errs.IsExpectedErrors(err, []string{"IncorrectListenerStatus", "IncorrectInstanceStatus"}) {
					_ = ensServiceV2.StopEnsLoadBalancerHttpListener(d.Id())
					wait()
					return retry.RetryableError(err)
				}
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

	// Converge the listener status to the desired value after an explicit
	// status change or an attribute change (the Set API may have stopped the listener).
	if d.HasChange("status") || attributeChanged {
		targetStatus := d.Get("status").(string)
		if targetStatus == "running" {
			if err := ensServiceV2.StartEnsLoadBalancerHttpListener(d.Id()); err != nil {
				return errs.WrapError(err)
			}
			startStateConf := helper.BuildStateConf([]string{}, []string{"running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, ensServiceV2.EnsLoadBalancerHttpListenerStateRefreshFunc(d.Id(), "Status", []string{}))
			if _, err := startStateConf.WaitForState(); err != nil {
				return errs.WrapErrorf(err, errs.IdMsg, d.Id())
			}
		} else {
			if err := ensServiceV2.StopEnsLoadBalancerHttpListener(d.Id()); err != nil {
				return errs.WrapError(err)
			}
			stopStateConf := helper.BuildStateConf([]string{}, []string{"stopped"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, ensServiceV2.EnsLoadBalancerHttpListenerStateRefreshFunc(d.Id(), "Status", []string{}))
			if _, err := stopStateConf.WaitForState(); err != nil {
				return errs.WrapErrorf(err, errs.IdMsg, d.Id())
			}
		}
	}

	return resourceEnsLoadBalancerHttpListenerRead(d, meta)
}

func resourceEnsLoadBalancerHttpListenerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ensServiceV2 := NewService(client)

	action := "DeleteLoadBalancerHTTPListener"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	parts := strings.Split(d.Id(), ":")
	request["LoadBalancerId"] = parts[0]
	request["ListenerPort"] = parts[1]

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("Ens", "2017-11-10", action, query, request, true)
		if err != nil {
			// A running listener must be stopped before it can be deleted.
			if errs.IsExpectedErrors(err, []string{"IncorrectListenerStatus", "IncorrectInstanceStatus"}) {
				_ = ensServiceV2.StopEnsLoadBalancerHttpListener(d.Id())
				wait()
				return retry.RetryableError(err)
			}
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
		if errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	// Wait until the listener is confirmed gone before returning, so a
	// same-port ForceNew recreate does not collide with a not-yet-released
	// listener on the ENS backend. Poll DescribeLoadBalancerHTTPListenerAttribute
	// explicitly until it reports the listener as not found; the previous
	// empty-target StateConf could return without an actual Describe poll,
	// leaving the port apparently busy for the immediate recreate.
	goneWait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		_, e := ensServiceV2.DescribeEnsLoadBalancerHttpListener(d.Id())
		if e != nil {
			if errs.NotFoundError(e) {
				return nil
			}
			if helper.NeedRetry(e) {
				goneWait()
				return retry.RetryableError(e)
			}
			return retry.NonRetryableError(e)
		}
		goneWait()
		return retry.RetryableError(errs.WrapError(errs.Error("LoadBalancerHttpListener %s still exists after delete", d.Id())))
	})
	if err != nil && !errs.NotFoundError(err) {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return nil
}

// ensLoadBalancerHttpListenerHealthCheckOffDiffSuppressFunc suppresses the
// ENS HTTP listener health check parameter fields when health_check is "off".
// The ENS Describe API omits HealthCheckDomain/URI/HealthyThreshold/
// UnhealthyThreshold/HealthCheckTimeout/HealthCheckInterval/HealthCheckHttpCode/
// HealthCheckMethod while health check is off (HealthCheckConnectPort is still
// returned), so these parameters must neither be planned into the diff nor set
// into state; otherwise they persist as empty/zero values and surface as
// unexpected attributes after refresh.
func ensLoadBalancerHttpListenerHealthCheckOffDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	return d.Get("health_check").(string) == "off"
}
