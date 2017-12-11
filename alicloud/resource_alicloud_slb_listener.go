package alicloud

import (
	"fmt"
	"github.com/denverdino/aliyungo/slb"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func resourceAliyunSlbListener() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunSlbListenerCreate,
		Read:   resourceAliyunSlbListenerRead,
		Update: resourceAliyunSlbListenerUpdate,
		Delete: resourceAliyunSlbListenerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"load_balancer_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"frontend_port": &schema.Schema{
				Type:         schema.TypeInt,
				ValidateFunc: validateInstancePort,
				Required:     true,
				ForceNew:     true,
			},
			"lb_port": &schema.Schema{
				Type:       schema.TypeInt,
				Optional:   true,
				Deprecated: "Field 'lb_port' has been deprecated, and using 'frontend_port' to replace.",
			},

			"backend_port": &schema.Schema{
				Type:         schema.TypeInt,
				ValidateFunc: validateInstancePort,
				Required:     true,
				ForceNew:     true,
			},

			"instance_port": &schema.Schema{
				Type:       schema.TypeInt,
				Optional:   true,
				Deprecated: "Field 'instance_port' has been deprecated, and using 'backend_port' to replace.",
			},

			"lb_protocol": &schema.Schema{
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "Field 'lb_protocol' has been deprecated, and using 'protocol' to replace.",
			},

			"protocol": &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: validateInstanceProtocol,
				Required:     true,
				ForceNew:     true,
			},

			"bandwidth": &schema.Schema{
				Type:         schema.TypeInt,
				ValidateFunc: validateSlbListenerBandwidth,
				Required:     true,
			},
			"scheduler": &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: validateSlbListenerScheduler,
				Optional:     true,
				Default:      slb.WRRScheduler,
			},
			//http & https
			"sticky_session": &schema.Schema{
				Type: schema.TypeString,
				ValidateFunc: validateAllowedStringValue([]string{
					string(slb.OnFlag),
					string(slb.OffFlag)}),
				Optional:         true,
				Default:          slb.OffFlag,
				DiffSuppressFunc: httpHttpsDiffSuppressFunc,
			},
			//http & https
			"sticky_session_type": &schema.Schema{
				Type: schema.TypeString,
				ValidateFunc: validateAllowedStringValue([]string{
					string(slb.InsertStickySessionType),
					string(slb.ServerStickySessionType)}),
				Optional:         true,
				DiffSuppressFunc: stickySessionTypeDiffSuppressFunc,
			},
			//http & https
			"cookie_timeout": &schema.Schema{
				Type:             schema.TypeInt,
				ValidateFunc:     validateSlbListenerCookieTimeout,
				Optional:         true,
				DiffSuppressFunc: cookieTimeoutDiffSuppressFunc,
			},
			//http & https
			"cookie": &schema.Schema{
				Type:             schema.TypeString,
				ValidateFunc:     validateSlbListenerCookie,
				Optional:         true,
				DiffSuppressFunc: cookieDiffSuppressFunc,
			},
			//tcp & udp
			"persistence_timeout": &schema.Schema{
				Type:             schema.TypeInt,
				ValidateFunc:     validateSlbListenerPersistenceTimeout,
				Optional:         true,
				Default:          0,
				DiffSuppressFunc: tcpUdpDiffSuppressFunc,
			},
			//http & https
			"health_check": &schema.Schema{
				Type: schema.TypeString,
				ValidateFunc: validateAllowedStringValue([]string{
					string(slb.OnFlag),
					string(slb.OffFlag)}),
				Optional:         true,
				Default:          slb.OnFlag,
				DiffSuppressFunc: httpHttpsDiffSuppressFunc,
			},
			//tcp
			"health_check_type": &schema.Schema{
				Type: schema.TypeString,
				ValidateFunc: validateAllowedStringValue([]string{
					string(slb.TCPHealthCheckType),
					string(slb.HTTPHealthCheckType)}),
				Optional:         true,
				Default:          slb.TCPHealthCheckType,
				DiffSuppressFunc: healthCheckTypeDiffSuppressFunc,
			},
			//http & https & tcp
			"health_check_domain": &schema.Schema{
				Type:             schema.TypeString,
				ValidateFunc:     validateSlbListenerHealthCheckDomain,
				Optional:         true,
				DiffSuppressFunc: httpHttpsTcpDiffSuppressFunc,
			},
			//http & https & tcp
			"health_check_uri": &schema.Schema{
				Type:             schema.TypeString,
				ValidateFunc:     validateSlbListenerHealthCheckUri,
				Optional:         true,
				Default:          "/",
				DiffSuppressFunc: httpHttpsTcpDiffSuppressFunc,
			},
			"health_check_connect_port": &schema.Schema{
				Type:             schema.TypeInt,
				ValidateFunc:     validateSlbListenerHealthCheckConnectPort,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: healthCheckDiffSuppressFunc,
			},
			"healthy_threshold": &schema.Schema{
				Type:             schema.TypeInt,
				ValidateFunc:     validateIntegerInRange(1, 10),
				Optional:         true,
				Default:          3,
				DiffSuppressFunc: healthCheckDiffSuppressFunc,
			},
			"unhealthy_threshold": &schema.Schema{
				Type:             schema.TypeInt,
				ValidateFunc:     validateIntegerInRange(1, 10),
				Optional:         true,
				Default:          3,
				DiffSuppressFunc: healthCheckDiffSuppressFunc,
			},

			"health_check_timeout": &schema.Schema{
				Type:             schema.TypeInt,
				ValidateFunc:     validateIntegerInRange(1, 300),
				Optional:         true,
				Default:          5,
				DiffSuppressFunc: healthCheckDiffSuppressFunc,
			},
			"health_check_interval": &schema.Schema{
				Type:             schema.TypeInt,
				ValidateFunc:     validateIntegerInRange(1, 50),
				Optional:         true,
				Default:          2,
				DiffSuppressFunc: healthCheckDiffSuppressFunc,
			},
			//http & https & tcp
			"health_check_http_code": &schema.Schema{
				Type: schema.TypeString,
				ValidateFunc: validateAllowedSplitStringValue([]string{
					string(slb.HTTP_2XX),
					string(slb.HTTP_3XX),
					string(slb.HTTP_4XX),
					string(slb.HTTP_5XX)}, ","),
				Optional:         true,
				Default:          slb.HTTP_2XX,
				DiffSuppressFunc: httpHttpsTcpDiffSuppressFunc,
			},
			//https
			"ssl_certificate_id": &schema.Schema{
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: sslCertificateIdDiffSuppressFunc,
			},
		},
	}
}

func resourceAliyunSlbListenerCreate(d *schema.ResourceData, meta interface{}) error {

	slbconn := meta.(*AliyunClient).slbconn

	protocol := d.Get("protocol").(string)
	lb_id := d.Get("load_balancer_id").(string)
	frontend := d.Get("frontend_port").(int)
	var err error

	switch Protocol(protocol) {
	case Https:
		ssl_id, ok := d.GetOk("ssl_certificate_id")
		if !ok || ssl_id == "" {
			return fmt.Errorf("'ssl_certificate_id': required field is not set when the protocol is 'https'.")
		}
		httpType, buildErr := buildHttpListenerType(d)
		if buildErr != nil {
			return buildErr
		}
		args := slb.CreateLoadBalancerHTTPSListenerArgs(slb.HTTPSListenerType{
			HTTPListenerType:    httpType,
			ServerCertificateId: ssl_id.(string),
		})
		err = slbconn.CreateLoadBalancerHTTPSListener(&args)
	case Tcp:
		args := buildTcpListenerArgs(d)
		err = slbconn.CreateLoadBalancerTCPListener(&args)
	case Udp:
		args := buildUdpListenerArgs(d)
		err = slbconn.CreateLoadBalancerUDPListener(&args)
	default:
		httpType, buildErr := buildHttpListenerType(d)
		if buildErr != nil {
			return buildErr
		}
		args := slb.CreateLoadBalancerHTTPListenerArgs(httpType)
		err = slbconn.CreateLoadBalancerHTTPListener(&args)
	}

	if err != nil {
		if IsExceptedError(err, ListenerAlreadyExists) {
			return fmt.Errorf("The listener with the frontend port %d already exists. Please define a new 'alicloud_slb_listener' resource and "+
				"use ID '%s:%d' to import it or modify its frontend port and then try again.", frontend, lb_id, frontend)
		}
		return fmt.Errorf("Create %s Listener got an error: %#v", protocol, err)
	}

	d.SetId(lb_id + ":" + strconv.Itoa(frontend))

	if err := slbconn.WaitForListenerAsyn(lb_id, frontend, slb.ListenerType(protocol), slb.Stopped, defaultTimeout); err != nil {
		return fmt.Errorf("WaitForListener %s got error: %#v", slb.Stopped, err)
	}

	if err := slbconn.StartLoadBalancerListener(lb_id, frontend); err != nil {
		return err
	}

	if err := slbconn.WaitForListenerAsyn(lb_id, frontend, slb.ListenerType(protocol), slb.Running, defaultTimeout); err != nil {
		return fmt.Errorf("WaitForListener %s got error: %#v", slb.Running, err)
	}

	return resourceAliyunSlbListenerUpdate(d, meta)
}

func resourceAliyunSlbListenerRead(d *schema.ResourceData, meta interface{}) error {
	slbconn := meta.(*AliyunClient).slbconn
	lb_id, protocol, port, err := parseListenerId(d, meta)
	if err != nil {
		return fmt.Errorf("Get slb listener got an error: %#v", err)
	}

	if protocol == "" {
		d.SetId("")
		return nil
	}
	d.Set("protocol", protocol)
	d.Set("load_balancer_id", lb_id)

	switch Protocol(protocol) {
	case Https:
		https_ls, err := slbconn.DescribeLoadBalancerHTTPSListenerAttribute(lb_id, port)
		return readListenerAttribute(d, protocol, https_ls, err)
	case Tcp:
		tcp_ls, err := slbconn.DescribeLoadBalancerTCPListenerAttribute(lb_id, port)
		return readListenerAttribute(d, protocol, tcp_ls, err)
	case Udp:
		udp_ls, err := slbconn.DescribeLoadBalancerUDPListenerAttribute(lb_id, port)
		return readListenerAttribute(d, protocol, udp_ls, err)
	default:
		http_ls, err := slbconn.DescribeLoadBalancerHTTPListenerAttribute(lb_id, port)
		return readListenerAttribute(d, protocol, http_ls, err)
	}
}

func resourceAliyunSlbListenerUpdate(d *schema.ResourceData, meta interface{}) error {

	slbconn := meta.(*AliyunClient).slbconn
	protocol := Protocol(d.Get("protocol").(string))

	d.Partial(true)

	httpType, err := buildHttpListenerType(d)
	if (protocol == Https || protocol == Http) && err != nil {
		return err
	}
	tcpArgs := slb.SetLoadBalancerTCPListenerAttributeArgs(buildTcpListenerArgs(d))
	udpArgs := slb.SetLoadBalancerUDPListenerAttributeArgs(buildUdpListenerArgs(d))
	httpsArgs := slb.SetLoadBalancerHTTPSListenerAttributeArgs(slb.CreateLoadBalancerHTTPSListenerArgs(slb.HTTPSListenerType{}))

	update := false
	if d.HasChange("scheduler") {
		scheduler := slb.SchedulerType(d.Get("scheduler").(string))
		httpType.Scheduler = scheduler
		tcpArgs.Scheduler = scheduler
		udpArgs.Scheduler = scheduler
		d.SetPartial("scheduler")
		update = true
	}
	// http https
	if d.HasChange("sticky_session") {
		d.SetPartial("sticky_session")
		update = true
	}
	if d.HasChange("sticky_session_type") {
		d.SetPartial("sticky_session_type")
		update = true
	}
	if d.HasChange("cookie_timeout") {
		d.SetPartial("cookie_timeout")
		update = true
	}
	if d.HasChange("cookie") {
		d.SetPartial("cookie")
		update = true
	}

	// http https
	if d.HasChange("health_check") {
		d.SetPartial("health_check")
		update = true
	}

	// http https tcp
	if d.HasChange("health_check_domain") {
		if domain, ok := d.GetOk("health_check_domain"); ok {
			httpType.HealthCheckDomain = domain.(string)
			tcpArgs.HealthCheckDomain = domain.(string)
			d.SetPartial("health_check_domain")
			update = true
		}
	}
	if d.HasChange("health_check_uri") {
		tcpArgs.HealthCheckURI = d.Get("health_check_uri").(string)
		d.SetPartial("health_check_uri")
		update = true
	}
	if d.HasChange("health_check_http_code") {
		tcpArgs.HealthCheckHttpCode = slb.HealthCheckHttpCodeType(d.Get("health_check_http_code").(string))
		d.SetPartial("health_check_http_code")
		update = true
	}

	// http https tcp udp and health_check=on
	if d.HasChange("unhealthy_threshold") {
		tcpArgs.UnhealthyThreshold = d.Get("unhealthy_threshold").(int)
		udpArgs.UnhealthyThreshold = d.Get("unhealthy_threshold").(int)
		d.SetPartial("unhealthy_threshold")
		update = true
		//}
	}
	if d.HasChange("healthy_threshold") {
		tcpArgs.HealthyThreshold = d.Get("healthy_threshold").(int)
		udpArgs.HealthyThreshold = d.Get("healthy_threshold").(int)
		d.SetPartial("healthy_threshold")
		update = true
	}
	if d.HasChange("health_check_timeout") {
		tcpArgs.HealthCheckConnectTimeout = d.Get("health_check_timeout").(int)
		udpArgs.HealthCheckConnectTimeout = d.Get("health_check_timeout").(int)
		d.SetPartial("health_check_timeout")
		update = true
	}
	if d.HasChange("health_check_interval") {
		tcpArgs.HealthCheckInterval = d.Get("health_check_interval").(int)
		udpArgs.HealthCheckInterval = d.Get("health_check_interval").(int)
		d.SetPartial("health_check_interval")
		update = true
	}
	if d.HasChange("health_check_connect_port") {
		if port, ok := d.GetOk("health_check_connect_port"); ok {
			httpType.HealthCheckConnectPort = port.(int)
			tcpArgs.HealthCheckConnectPort = port.(int)
			udpArgs.HealthCheckConnectPort = port.(int)
			d.SetPartial("health_check_connect_port")
			update = true
		}
	}

	// tcp and udp
	if d.HasChange("persistence_timeout") {
		tcpArgs.PersistenceTimeout = d.Get("persistence_timeout").(int)
		udpArgs.PersistenceTimeout = d.Get("persistence_timeout").(int)
		d.SetPartial("persistence_timeout")
		update = true
	}

	// tcp
	if d.HasChange("health_check_type") {
		tcpArgs.HealthCheckType = slb.HealthCheckType(d.Get("health_check_type").(string))
		d.SetPartial("health_check_type")
		update = true
	}

	// https
	if protocol == Https {
		ssl_id, ok := d.GetOk("ssl_certificate_id")
		if !ok && ssl_id == "" {
			return fmt.Errorf("'ssl_certificate_id': required field is not set when the protocol is 'https'.")
		}

		httpsArgs.ServerCertificateId = ssl_id.(string)
		if d.HasChange("ssl_certificate_id") {
			d.SetPartial("ssl_certificate_id")
			update = true
		}
	}

	if update {
		switch protocol {
		case Https:
			httpsArgs.HTTPListenerType = httpType
			if err := slbconn.SetLoadBalancerHTTPSListenerAttribute(&httpsArgs); err != nil {
				return fmt.Errorf("SetHTTPSListenerAttribute got an error: %#v", err)
			}
		case Tcp:
			if err := slbconn.SetLoadBalancerTCPListenerAttribute(&tcpArgs); err != nil {
				return fmt.Errorf("SetTCPListenerAttribute got an error: %#v", err)
			}
		case Udp:
			if err := slbconn.SetLoadBalancerUDPListenerAttribute(&udpArgs); err != nil {
				return fmt.Errorf("SetTCPListenerAttribute got an error: %#v", err)
			}
		default:
			httpArgs := slb.SetLoadBalancerHTTPListenerAttributeArgs(slb.CreateLoadBalancerHTTPListenerArgs(httpType))
			if err := slbconn.SetLoadBalancerHTTPListenerAttribute(&httpArgs); err != nil {
				return fmt.Errorf("SetHTTPListenerAttribute got an error: %#v", err)
			}
		}
	}

	d.Partial(false)

	return resourceAliyunSlbListenerRead(d, meta)
}

func resourceAliyunSlbListenerDelete(d *schema.ResourceData, meta interface{}) error {
	slbconn := meta.(*AliyunClient).slbconn
	lb_id, protocol, port, err := parseListenerId(d, meta)
	if err != nil {
		return fmt.Errorf("Get slb listener got an error: %#v", err)
	}

	if protocol == "" {
		d.SetId("")
		return nil
	}
	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		err := slbconn.DeleteLoadBalancerListener(lb_id, port)

		if err != nil {
			return resource.NonRetryableError(err)
		}

		switch Protocol(protocol) {
		case Https:
			https_ls, err := slbconn.DescribeLoadBalancerHTTPSListenerAttribute(lb_id, port)
			return ensureListenerAbsent(d, protocol, https_ls, err)
		case Tcp:
			tcp_ls, err := slbconn.DescribeLoadBalancerTCPListenerAttribute(lb_id, port)
			return ensureListenerAbsent(d, protocol, tcp_ls, err)
		case Udp:
			udp_ls, err := slbconn.DescribeLoadBalancerUDPListenerAttribute(lb_id, port)
			return ensureListenerAbsent(d, protocol, udp_ls, err)
		default:
			http_ls, err := slbconn.DescribeLoadBalancerHTTPListenerAttribute(lb_id, port)
			return ensureListenerAbsent(d, protocol, http_ls, err)
		}
	})
}

func buildHttpListenerType(d *schema.ResourceData) (slb.HTTPListenerType, error) {

	httpType := slb.HTTPListenerType{
		LoadBalancerId:    d.Get("load_balancer_id").(string),
		ListenerPort:      d.Get("frontend_port").(int),
		BackendServerPort: d.Get("backend_port").(int),
		Bandwidth:         d.Get("bandwidth").(int),
		StickySession:     slb.FlagType(d.Get("sticky_session").(string)),
		HealthCheck:       slb.FlagType(d.Get("health_check").(string)),
	}
	if httpType.StickySession == slb.OnFlag {
		if sessionType, ok := d.GetOk("sticky_session_type"); !ok || sessionType.(string) == "" {
			return httpType, fmt.Errorf("'sticky_session_type': required field is not set when the StickySession is %s.", slb.OnFlag)
		} else {
			httpType.StickySessionType = slb.StickySessionType(sessionType.(string))

		}
		if httpType.StickySessionType == slb.InsertStickySessionType {
			if timeout, ok := d.GetOk("cookie_timeout"); !ok || timeout == 0 {
				return httpType, fmt.Errorf("'cookie_timeout': required field is not set when the StickySession is %s "+
					"and StickySessionType is %s.", slb.OnFlag, slb.InsertStickySessionType)
			} else {
				httpType.CookieTimeout = timeout.(int)
			}
		} else {
			if cookie, ok := d.GetOk("cookie"); !ok || cookie.(string) == "" {
				return httpType, fmt.Errorf("'cookie': required field is not set when the StickySession is %s "+
					"and StickySessionType is %s.", slb.OnFlag, slb.ServerStickySessionType)
			} else {
				httpType.Cookie = cookie.(string)
			}
		}
	}
	if httpType.HealthCheck == slb.OnFlag {
		httpType.HealthCheckURI = d.Get("health_check_uri").(string)
		if port, ok := d.GetOk("health_check_connect_port"); !ok || port.(int) == 0 {
			return httpType, fmt.Errorf("'health_check_connect_port': required field is not set when the HealthCheck is %s.", slb.OnFlag)
		} else {
			httpType.HealthCheckConnectPort = port.(int)
		}
		httpType.HealthyThreshold = d.Get("healthy_threshold").(int)
		httpType.UnhealthyThreshold = d.Get("unhealthy_threshold").(int)
		httpType.HealthCheckTimeout = d.Get("health_check_timeout").(int)
		httpType.HealthCheckInterval = d.Get("health_check_interval").(int)
		httpType.HealthCheckHttpCode = slb.HealthCheckHttpCodeType(d.Get("health_check_http_code").(string))
	}
	return httpType, nil
}

func buildTcpListenerArgs(d *schema.ResourceData) slb.CreateLoadBalancerTCPListenerArgs {

	return slb.CreateLoadBalancerTCPListenerArgs(slb.TCPListenerType{
		LoadBalancerId:    d.Get("load_balancer_id").(string),
		ListenerPort:      d.Get("frontend_port").(int),
		BackendServerPort: d.Get("backend_port").(int),
		Bandwidth:         d.Get("bandwidth").(int),
	})
}
func buildUdpListenerArgs(d *schema.ResourceData) slb.CreateLoadBalancerUDPListenerArgs {

	return slb.CreateLoadBalancerUDPListenerArgs(slb.UDPListenerType{
		LoadBalancerId:    d.Get("load_balancer_id").(string),
		ListenerPort:      d.Get("frontend_port").(int),
		BackendServerPort: d.Get("backend_port").(int),
		Bandwidth:         d.Get("bandwidth").(int),
	})
}

func parseListenerId(d *schema.ResourceData, meta interface{}) (string, string, int, error) {
	slbconn := meta.(*AliyunClient).slbconn
	parts := strings.Split(d.Id(), ":")
	port, err := strconv.Atoi(parts[1])
	if err != nil {
		return "", "", 0, fmt.Errorf("Parsing SlbListener's id got an error: %#v", err)
	}
	loadBalancer, err := slbconn.DescribeLoadBalancerAttribute(parts[0])
	if err != nil {
		if IsExceptedError(err, LoadBalancerNotFound) {
			return "", "", 0, nil
		}
		return "", "", 0, fmt.Errorf("DescribeLoadBalancerAttribute got an error: %#v", parts[0])
	}
	for _, portAndProtocol := range loadBalancer.ListenerPortsAndProtocol.ListenerPortAndProtocol {
		if portAndProtocol.ListenerPort == port {
			return loadBalancer.LoadBalancerId, portAndProtocol.ListenerProtocol, port, nil
		}
	}
	return "", "", 0, nil
}

func readListenerAttribute(d *schema.ResourceData, protocol string, listen interface{}, err error) error {
	v := reflect.ValueOf(listen).Elem()

	if err != nil {
		if IsExceptedError(err, ListenerNotFound) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("DescribeLoadBalancer%sListenerAttribute got an error: %#v", strings.ToUpper(protocol), err)
	}
	if port := v.FieldByName("ListenerPort"); port.IsValid() && port.Interface().(int) > 0 {
		readListener(d, listen)
	} else {
		d.SetId("")
	}
	return nil
}

func readListener(d *schema.ResourceData, listen interface{}) {
	v := reflect.ValueOf(listen).Elem()

	if val := v.FieldByName("BackendServerPort"); val.IsValid() {
		d.Set("backend_port", val.Interface().(int))
	}
	if val := v.FieldByName("ListenerPort"); val.IsValid() {
		d.Set("frontend_port", val.Interface().(int))
	}
	if val := v.FieldByName("Bandwidth"); val.IsValid() {
		d.Set("bandwidth", val.Interface().(int))
	}
	if val := v.FieldByName("Scheduler"); val.IsValid() {
		d.Set("scheduler", string(val.Interface().(slb.SchedulerType)))
	}
	if val := v.FieldByName("HealthCheck"); val.IsValid() {
		d.Set("health_check", string(val.Interface().(slb.FlagType)))
	}
	if val := v.FieldByName("StickySession"); val.IsValid() {
		d.Set("sticky_session", string(val.Interface().(slb.FlagType)))
	}
	if val := v.FieldByName("StickySessionType"); val.IsValid() {
		d.Set("sticky_session_type", string(val.Interface().(slb.StickySessionType)))
	}
	if val := v.FieldByName("CookieTimeout"); val.IsValid() {
		d.Set("cookie_timeout", val.Interface().(int))
	}
	if val := v.FieldByName("Cookie"); val.IsValid() {
		d.Set("cookie", val.Interface().(string))
	}
	if val := v.FieldByName("PersistenceTimeout"); val.IsValid() {
		d.Set("persistence_timeout", val.Interface().(int))
	}
	if val := v.FieldByName("HealthCheckType"); val.IsValid() {
		d.Set("health_check_type", string(val.Interface().(slb.HealthCheckType)))
	}
	if val := v.FieldByName("HealthCheckDomain"); val.IsValid() {
		d.Set("health_check_domain", val.Interface().(string))
	}
	if val := v.FieldByName("HealthCheckConnectPort"); val.IsValid() {
		d.Set("health_check_connect_port", val.Interface().(int))
	}
	if val := v.FieldByName("HealthCheckURI"); val.IsValid() {
		d.Set("health_check_uri", val.Interface().(string))
	}
	if val := v.FieldByName("HealthyThreshold"); val.IsValid() {
		d.Set("healthy_threshold", val.Interface().(int))
	}
	if val := v.FieldByName("UnhealthyThreshold"); val.IsValid() {
		d.Set("unhealthy_threshold", val.Interface().(int))
	}
	if val := v.FieldByName("HealthCheckTimeout"); val.IsValid() {
		d.Set("health_check_timeout", val.Interface().(int))
	}
	if val := v.FieldByName("HealthCheckConnectTimeout"); val.IsValid() {
		d.Set("health_check_timeout", val.Interface().(int))
	}
	if val := v.FieldByName("HealthCheckInterval"); val.IsValid() {
		d.Set("health_check_interval", val.Interface().(int))
	}
	if val := v.FieldByName("HealthCheckHttpCode"); val.IsValid() {
		d.Set("health_check_http_code", string(val.Interface().(slb.HealthCheckHttpCodeType)))
	}
	if val := v.FieldByName("ServerCertificateId"); val.IsValid() {
		d.Set("ssl_certificate_id", val.Interface().(string))
	}

	return
}

func ensureListenerAbsent(d *schema.ResourceData, protocol string, listen interface{}, err error) *resource.RetryError {
	v := reflect.ValueOf(listen).Elem()

	if err != nil {
		if IsExceptedError(err, ListenerNotFound) {
			d.SetId("")
			return nil
		}
		return resource.NonRetryableError(fmt.Errorf("While deleting listener, DescribeLoadBalancer%sListenerAttribute got an error: %#v", protocol, err))
	}
	if port := v.FieldByName("ListenerPort"); port.IsValid() && port.Interface().(int) > 0 {
		return resource.RetryableError(fmt.Errorf("Listener in use - trying again while it deleted."))
	}
	d.SetId("")
	return nil
}
