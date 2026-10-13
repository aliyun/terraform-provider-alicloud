package alicloud

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAliCloudEnsLoadBalancerHttpListener() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudEnsLoadBalancerHttpListenerCreate,
		Read:   resourceAliCloudEnsLoadBalancerHttpListenerRead,
		Update: resourceAliCloudEnsLoadBalancerHttpListenerUpdate,
		Delete: resourceAliCloudEnsLoadBalancerHttpListenerDelete,
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
				ValidateFunc: IntBetween(1, 65535),
			},
			"backend_server_port": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: IntBetween(1, 65535),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"scheduler": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "wrr",
				ValidateFunc: StringInSlice([]string{"wrr", "wlc", "rr", "sch", "qch", "iqch"}, false),
			},
			"health_check": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "off",
				ValidateFunc: StringInSlice([]string{"on", "off"}, false),
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
				ValidateFunc:     IntBetween(2, 10),
				DiffSuppressFunc: ensLoadBalancerHttpListenerHealthCheckOffDiffSuppressFunc,
			},
			"unhealthy_threshold": {
				Type:             schema.TypeInt,
				Optional:         true,
				Default:          3,
				ValidateFunc:     IntBetween(2, 10),
				DiffSuppressFunc: ensLoadBalancerHttpListenerHealthCheckOffDiffSuppressFunc,
			},
			"health_check_timeout": {
				Type:             schema.TypeInt,
				Optional:         true,
				Default:          5,
				ValidateFunc:     IntBetween(1, 300),
				DiffSuppressFunc: ensLoadBalancerHttpListenerHealthCheckOffDiffSuppressFunc,
			},
			"health_check_connect_port": {
				Type:             schema.TypeInt,
				Optional:         true,
				Computed:         true,
				ValidateFunc:     IntBetween(1, 65535),
				DiffSuppressFunc: ensLoadBalancerHttpListenerHealthCheckOffDiffSuppressFunc,
			},
			"health_check_interval": {
				Type:             schema.TypeInt,
				Optional:         true,
				Default:          2,
				ValidateFunc:     IntBetween(1, 50),
				DiffSuppressFunc: ensLoadBalancerHttpListenerHealthCheckOffDiffSuppressFunc,
			},
			"health_check_http_code": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "http_2xx",
				ValidateFunc:     StringInSlice([]string{"http_2xx", "http_3xx", "http_4xx", "http_5xx"}, false),
				DiffSuppressFunc: ensLoadBalancerHttpListenerHealthCheckOffDiffSuppressFunc,
			},
			"health_check_method": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "head",
				ValidateFunc:     StringInSlice([]string{"head", "get"}, false),
				DiffSuppressFunc: ensLoadBalancerHttpListenerHealthCheckOffDiffSuppressFunc,
			},
			"idle_timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      15,
				ValidateFunc: IntBetween(1, 60),
			},
			"request_timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      60,
				ValidateFunc: IntBetween(1, 180),
			},
			"x_forwarded_for": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "off",
				ValidateFunc: StringInSlice([]string{"on", "off"}, false),
			},
			"listener_forward": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "off",
				ValidateFunc: StringInSlice([]string{"on", "off"}, false),
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
				ValidateFunc: StringInSlice([]string{"running", "stopped"}, false),
			},
		},
	}
}

func resourceAliCloudEnsLoadBalancerHttpListenerCreate(d *schema.ResourceData, meta interface{}) error {
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

	wait := incrementalWait(10*time.Second, 15*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("Ens", "2017-11-10", action, query, request, true)
		if err != nil {
			// A same-port ForceNew destroy (backend_server_port is ForceNew
			// while listener_port is not) may leave the listener port still
			// releasing on the ENS backend; absorb the consistency window.
			if IsExpectedErrors(err, []string{"ListenerAlreadyExists"}) {
				wait()
				return retry.RetryableError(err)
			}
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ens_load_balancer_http_listener", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", d.Get("load_balancer_id"), d.Get("listener_port")))

	ensServiceV2 := EnsServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"stopped"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, ensServiceV2.EnsLoadBalancerHttpListenerStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	// Start the listener when the desired status is running.
	if d.Get("status").(string) == "running" {
		if err := ensServiceV2.StartEnsLoadBalancerHttpListener(d.Id()); err != nil {
			return WrapError(err)
		}
		startStateConf := BuildStateConf([]string{}, []string{"running"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, ensServiceV2.EnsLoadBalancerHttpListenerStateRefreshFunc(d.Id(), "Status", []string{}))
		if _, err := startStateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	return resourceAliCloudEnsLoadBalancerHttpListenerRead(d, meta)
}

func resourceAliCloudEnsLoadBalancerHttpListenerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ensServiceV2 := EnsServiceV2{client}

	objectRaw, err := ensServiceV2.DescribeEnsLoadBalancerHttpListener(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ens_load_balancer_http_listener DescribeEnsLoadBalancerHttpListener Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	parts := strings.Split(d.Id(), ":")
	if len(parts) == 2 {
		d.Set("load_balancer_id", parts[0])
		listenerPort, _ := strconv.Atoi(parts[1])
		d.Set("listener_port", listenerPort)
	}

	if v := formatInt(objectRaw["BackendServerPort"]); v > 0 {
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
		d.Set("healthy_threshold", formatInt(objectRaw["HealthyThreshold"]))
		d.Set("unhealthy_threshold", formatInt(objectRaw["UnhealthyThreshold"]))
		d.Set("health_check_timeout", formatInt(objectRaw["HealthCheckTimeout"]))
		if v := formatInt(objectRaw["HealthCheckConnectPort"]); v > 0 {
			d.Set("health_check_connect_port", v)
		}
		d.Set("health_check_interval", formatInt(objectRaw["HealthCheckInterval"]))
		d.Set("health_check_http_code", objectRaw["HealthCheckHttpCode"])
		d.Set("health_check_method", objectRaw["HealthCheckMethod"])
	}
	d.Set("idle_timeout", formatInt(objectRaw["IdleTimeout"]))
	d.Set("request_timeout", formatInt(objectRaw["RequestTimeout"]))
	d.Set("x_forwarded_for", objectRaw["XForwardedFor"])
	d.Set("listener_forward", objectRaw["ListenerForward"])
	if v := formatInt(objectRaw["ForwardPort"]); v > 0 {
		d.Set("forward_port", v)
	}
	d.Set("status", objectRaw["Status"])

	return nil
}

func resourceAliCloudEnsLoadBalancerHttpListenerUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ensServiceV2 := EnsServiceV2{client}

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

		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("Ens", "2017-11-10", action, query, request, true)
			if err != nil {
				// The Set API requires the listener to be in the stopped state.
				// When the listener is running, stop it and retry the modification.
				if IsExpectedErrors(err, []string{"IncorrectListenerStatus", "IncorrectInstanceStatus"}) {
					_ = ensServiceV2.StopEnsLoadBalancerHttpListener(d.Id())
					wait()
					return retry.RetryableError(err)
				}
				if NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}

	// Converge the listener status to the desired value after an explicit
	// status change or an attribute change (the Set API may have stopped the listener).
	if d.HasChange("status") || attributeChanged {
		targetStatus := d.Get("status").(string)
		if targetStatus == "running" {
			if err := ensServiceV2.StartEnsLoadBalancerHttpListener(d.Id()); err != nil {
				return WrapError(err)
			}
			startStateConf := BuildStateConf([]string{}, []string{"running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, ensServiceV2.EnsLoadBalancerHttpListenerStateRefreshFunc(d.Id(), "Status", []string{}))
			if _, err := startStateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}
		} else {
			if err := ensServiceV2.StopEnsLoadBalancerHttpListener(d.Id()); err != nil {
				return WrapError(err)
			}
			stopStateConf := BuildStateConf([]string{}, []string{"stopped"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, ensServiceV2.EnsLoadBalancerHttpListenerStateRefreshFunc(d.Id(), "Status", []string{}))
			if _, err := stopStateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}
		}
	}

	return resourceAliCloudEnsLoadBalancerHttpListenerRead(d, meta)
}

func resourceAliCloudEnsLoadBalancerHttpListenerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ensServiceV2 := EnsServiceV2{client}

	action := "DeleteLoadBalancerHTTPListener"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	parts := strings.Split(d.Id(), ":")
	request["LoadBalancerId"] = parts[0]
	request["ListenerPort"] = parts[1]

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("Ens", "2017-11-10", action, query, request, true)
		if err != nil {
			// A running listener must be stopped before it can be deleted.
			if IsExpectedErrors(err, []string{"IncorrectListenerStatus", "IncorrectInstanceStatus"}) {
				_ = ensServiceV2.StopEnsLoadBalancerHttpListener(d.Id())
				wait()
				return retry.RetryableError(err)
			}
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
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

	// Wait until the listener is confirmed gone before returning, so a
	// same-port ForceNew recreate does not collide with a not-yet-released
	// listener on the ENS backend. Poll DescribeLoadBalancerHTTPListenerAttribute
	// explicitly until it reports the listener as not found; the previous
	// empty-target StateConf could return without an actual Describe poll,
	// leaving the port apparently busy for the immediate recreate.
	goneWait := incrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		_, e := ensServiceV2.DescribeEnsLoadBalancerHttpListener(d.Id())
		if e != nil {
			if NotFoundError(e) {
				return nil
			}
			if NeedRetry(e) {
				goneWait()
				return retry.RetryableError(e)
			}
			return retry.NonRetryableError(e)
		}
		goneWait()
		return retry.RetryableError(WrapError(Error("LoadBalancerHttpListener %s still exists after delete", d.Id())))
	})
	if err != nil && !NotFoundError(err) {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}
