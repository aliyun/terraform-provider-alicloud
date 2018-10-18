package alicloud

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
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
				Default:      WRRScheduler,
			},
			"server_group_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"acl_status": &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: validateAllowedStringValue([]string{string(OnFlag), string(OffFlag)}),
				Optional:     true,
				Default:      OffFlag,
			},
			"acl_type": &schema.Schema{
				Type:             schema.TypeString,
				ValidateFunc:     validateAllowedStringValue([]string{string(AclTypeBlack), string(AclTypeWhite)}),
				Optional:         true,
				DiffSuppressFunc: slbAclDiffSuppressFunc,
			},
			"acl_id": &schema.Schema{
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: slbAclDiffSuppressFunc,
			},
			//http & https
			"sticky_session": &schema.Schema{
				Type:             schema.TypeString,
				ValidateFunc:     validateAllowedStringValue([]string{string(OnFlag), string(OffFlag)}),
				Optional:         true,
				Default:          OffFlag,
				DiffSuppressFunc: httpHttpsDiffSuppressFunc,
			},
			//http & https
			"sticky_session_type": &schema.Schema{
				Type: schema.TypeString,
				ValidateFunc: validateAllowedStringValue([]string{
					string(InsertStickySessionType),
					string(ServerStickySessionType)}),
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
				Type:             schema.TypeString,
				ValidateFunc:     validateAllowedStringValue([]string{string(OnFlag), string(OffFlag)}),
				Optional:         true,
				Default:          OnFlag,
				DiffSuppressFunc: httpHttpsDiffSuppressFunc,
			},
			//tcp
			"health_check_type": &schema.Schema{
				Type: schema.TypeString,
				ValidateFunc: validateAllowedStringValue([]string{
					string(TCPHealthCheckType),
					string(HTTPHealthCheckType)}),
				Optional:         true,
				Default:          TCPHealthCheckType,
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
					string(HTTP_2XX), string(HTTP_3XX), string(HTTP_4XX), string(HTTP_5XX)}, ","),
				Optional:         true,
				Default:          HTTP_2XX,
				DiffSuppressFunc: httpHttpsTcpDiffSuppressFunc,
			},
			//https
			"ssl_certificate_id": &schema.Schema{
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: sslCertificateIdDiffSuppressFunc,
			},

			//http, https
			"gzip": &schema.Schema{
				Type:             schema.TypeBool,
				Optional:         true,
				Default:          true,
				DiffSuppressFunc: httpHttpsDiffSuppressFunc,
			},
			"x_forwarded_for": &schema.Schema{
				Type:             schema.TypeList,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: httpHttpsDiffSuppressFunc,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						// At present, retrive client ip can not be modified, and it default to true.
						"retrive_client_ip": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"retrive_slb_ip": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"retrive_slb_id": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"retrive_slb_proto": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
				MaxItems: 1,
			},
		},
	}
}

func resourceAliyunSlbListenerCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*AliyunClient)

	protocol := d.Get("protocol").(string)
	lb_id := d.Get("load_balancer_id").(string)
	frontend := d.Get("frontend_port").(int)

	req := buildListenerCommonArgs(d, meta)
	req.ApiName = fmt.Sprintf("CreateLoadBalancer%sListener", strings.ToUpper(protocol))

	if Protocol(protocol) == Http || Protocol(protocol) == Https {
		reqHttp, err := buildHttpListenerArgs(d, req)
		if err != nil {
			return err
		}
		req = reqHttp
		if Protocol(protocol) == Https {
			ssl_id, ok := d.GetOk("ssl_certificate_id")
			if !ok || ssl_id.(string) == "" {
				return fmt.Errorf("'ssl_certificate_id': required field is not set when the protocol is 'https'.")
			}
			req.QueryParams["ServerCertificateId"] = ssl_id.(string)
		}
	}

	if _, err := client.slbconn.ProcessCommonRequest(req); err != nil {
		if IsExceptedErrors(err, []string{ListenerAlreadyExists}) {
			return fmt.Errorf("The listener with the frontend port %d already exists. Please define a new 'alicloud_slb_listener' resource and "+
				"use ID '%s:%d' to import it or modify its frontend port and then try again.", frontend, lb_id, frontend)
		}
		return fmt.Errorf("Create %s Listener got an error: %#v", protocol, err)
	}

	d.SetId(lb_id + ":" + strconv.Itoa(frontend))

	if err := client.WaitForListener(lb_id, frontend, Protocol(protocol), Stopped, DefaultTimeout); err != nil {
		return fmt.Errorf("WaitForListener %s got error: %#v", Stopped, err)
	}

	reqStart := slb.CreateStartLoadBalancerListenerRequest()
	reqStart.LoadBalancerId = lb_id
	reqStart.ListenerPort = requests.NewInteger(frontend)
	if _, err := client.slbconn.StartLoadBalancerListener(reqStart); err != nil {
		return err
	}

	if err := client.WaitForListener(lb_id, frontend, Protocol(protocol), Running, DefaultTimeout); err != nil {
		return fmt.Errorf("WaitForListener %s got error: %#v", Running, err)
	}

	return resourceAliyunSlbListenerUpdate(d, meta)
}

func resourceAliyunSlbListenerRead(d *schema.ResourceData, meta interface{}) error {
	lb_id, protocol, port, err := parseListenerId(d, meta)
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Get slb listener got an error: %#v", err)
	}

	d.Set("protocol", protocol)
	d.Set("load_balancer_id", lb_id)

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		listener, err := meta.(*AliyunClient).DescribeLoadBalancerListenerAttribute(lb_id, port, Protocol(protocol))
		return readListenerAttribute(d, protocol, listener, err)
	})
}

func resourceAliyunSlbListenerUpdate(d *schema.ResourceData, meta interface{}) error {

	slbconn := meta.(*AliyunClient).slbconn
	protocol := Protocol(d.Get("protocol").(string))

	d.Partial(true)

	commonArgs := buildListenerCommonArgs(d, meta)
	commonArgs.ApiName = fmt.Sprintf("SetLoadBalancer%sListenerAttribute", strings.ToUpper(string(protocol)))

	update := false
	if d.HasChange("scheduler") {
		commonArgs.QueryParams["Scheduler"] = d.Get("scheduler").(string)
		d.SetPartial("scheduler")
		update = true
	}

	if d.HasChange("server_group_id") {
		commonArgs.QueryParams["VServerGroupId"] = d.Get("server_group_id").(string)
		d.SetPartial("server_group_id")
		update = true
	}

	if d.HasChange("acl_status") {
		commonArgs.QueryParams["AclStatus"] = d.Get("acl_status").(string)
		d.SetPartial("acl_status")
		update = true
	}

	if d.HasChange("acl_type") {
		commonArgs.QueryParams["AclType"] = d.Get("acl_type").(string)
		d.SetPartial("acl_type")
		update = true
	}

	if d.HasChange("acl_id") {
		commonArgs.QueryParams["AclId"] = d.Get("acl_id").(string)
		d.SetPartial("acl_id")
		update = true
	}

	httpArgs, err := buildHttpListenerArgs(d, commonArgs)
	if (protocol == Https || protocol == Http) && err != nil {
		return err
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
	if d.HasChange("health_check") {
		d.SetPartial("health_check")
		update = true
	}
	if d.HasChange("gzip") || d.HasChange("x_forwarded_for") {
		update = true
		d.SetPartial("gzip")
		if d.Get("gzip").(bool) {
			httpArgs.QueryParams["Gzip"] = string(OnFlag)
		} else {
			httpArgs.QueryParams["Gzip"] = string(OffFlag)
		}

		d.SetPartial("x_forwarded_for")
		if len(d.Get("x_forwarded_for").([]interface{})) > 0 {
			xff := d.Get("x_forwarded_for").([]interface{})[0].(map[string]interface{})
			if xff["retrive_slb_ip"].(bool) {
				httpArgs.QueryParams["XForwardedFor_SLBIP"] = string(OnFlag)
			} else {
				httpArgs.QueryParams["XForwardedFor_SLBIP"] = string(OffFlag)
			}
			if xff["retrive_slb_id"].(bool) {
				httpArgs.QueryParams["XForwardedFor_SLBID"] = string(OnFlag)
			} else {
				httpArgs.QueryParams["XForwardedFor_SLBID"] = string(OffFlag)
			}
			if xff["retrive_slb_proto"].(bool) {
				httpArgs.QueryParams["XForwardedFor_proto"] = string(OnFlag)
			} else {
				httpArgs.QueryParams["XForwardedFor_proto"] = string(OffFlag)
			}
		}
	}

	// http https tcp udp and health_check=on
	if d.HasChange("unhealthy_threshold") {
		commonArgs.QueryParams["UnhealthyThreshold"] = string(requests.NewInteger(d.Get("unhealthy_threshold").(int)))
		d.SetPartial("unhealthy_threshold")
		update = true
		//}
	}
	if d.HasChange("healthy_threshold") {
		commonArgs.QueryParams["HealthyThreshold"] = string(requests.NewInteger(d.Get("healthy_threshold").(int)))
		d.SetPartial("healthy_threshold")
		update = true
	}
	if d.HasChange("health_check_timeout") {
		commonArgs.QueryParams["HealthCheckConnectTimeout"] = string(requests.NewInteger(d.Get("health_check_timeout").(int)))
		d.SetPartial("health_check_timeout")
		update = true
	}
	if d.HasChange("health_check_interval") {
		commonArgs.QueryParams["HealthCheckInterval"] = string(requests.NewInteger(d.Get("health_check_interval").(int)))
		d.SetPartial("health_check_interval")
		update = true
	}
	if d.HasChange("health_check_connect_port") {
		if port, ok := d.GetOk("health_check_connect_port"); ok {
			httpArgs.QueryParams["HealthCheckConnectPort"] = string(requests.NewInteger(port.(int)))
			commonArgs.QueryParams["HealthCheckConnectPort"] = string(requests.NewInteger(port.(int)))
			d.SetPartial("health_check_connect_port")
			update = true
		}
	}

	// tcp and udp
	if d.HasChange("persistence_timeout") {
		commonArgs.QueryParams["PersistenceTimeout"] = string(requests.NewInteger(d.Get("persistence_timeout").(int)))
		d.SetPartial("persistence_timeout")
		update = true
	}

	tcpArgs := commonArgs
	udpArgs := commonArgs

	// http https tcp
	if d.HasChange("health_check_domain") {
		if domain, ok := d.GetOk("health_check_domain"); ok {
			httpArgs.QueryParams["HealthCheckDomain"] = domain.(string)
			tcpArgs.QueryParams["HealthCheckDomain"] = domain.(string)
			d.SetPartial("health_check_domain")
			update = true
		}
	}
	if d.HasChange("health_check_uri") {
		tcpArgs.QueryParams["HealthCheckURI"] = d.Get("health_check_uri").(string)
		d.SetPartial("health_check_uri")
		update = true
	}
	if d.HasChange("health_check_http_code") {
		tcpArgs.QueryParams["HealthCheckHttpCode"] = d.Get("health_check_http_code").(string)
		d.SetPartial("health_check_http_code")
		update = true
	}

	// tcp
	if d.HasChange("health_check_type") {
		tcpArgs.QueryParams["HealthCheckType"] = d.Get("health_check_type").(string)
		d.SetPartial("health_check_type")
		update = true
	}

	// https
	httpsArgs := httpArgs
	if protocol == Https {
		ssl_id, ok := d.GetOk("ssl_certificate_id")
		if !ok && ssl_id == "" {
			return fmt.Errorf("'ssl_certificate_id': required field is not set when the protocol is 'https'.")
		}

		httpsArgs.QueryParams["ServerCertificateId"] = ssl_id.(string)
		if d.HasChange("ssl_certificate_id") {
			d.SetPartial("ssl_certificate_id")
			update = true
		}
	}

	if update {
		var request *requests.CommonRequest
		switch protocol {
		case Https:
			request = httpsArgs
		case Tcp:
			request = tcpArgs
		case Udp:
			request = udpArgs
		default:
			request = httpArgs
		}
		if _, err = slbconn.ProcessCommonRequest(request); err != nil {
			return fmt.Errorf("%s got an error: %#v", request.ApiName, err)
		}
	}

	d.Partial(false)

	return resourceAliyunSlbListenerRead(d, meta)
}

func resourceAliyunSlbListenerDelete(d *schema.ResourceData, meta interface{}) error {
	slbconn := meta.(*AliyunClient).slbconn
	lb_id, protocol, port, err := parseListenerId(d, meta)
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Get slb listener got an error: %#v", err)
	}

	req := slb.CreateDeleteLoadBalancerListenerRequest()
	req.LoadBalancerId = lb_id
	req.ListenerPort = requests.NewInteger(port)
	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err := slbconn.DeleteLoadBalancerListener(req)

		if err != nil {
			if IsExceptedErrors(err, SlbIsBusy) {
				return resource.RetryableError(fmt.Errorf("Delete load balancer listener timeout and got an error: %#v.", err))
			}
			return resource.NonRetryableError(err)
		}

		listener, err := meta.(*AliyunClient).DescribeLoadBalancerListenerAttribute(lb_id, port, Protocol(protocol))
		return ensureListenerAbsent(d, protocol, listener, err)
	})
}

func buildListenerCommonArgs(d *schema.ResourceData, meta interface{}) *requests.CommonRequest {
	req := meta.(*AliyunClient).BuildSlbCommonRequest()
	req.QueryParams["LoadBalancerId"] = d.Get("load_balancer_id").(string)
	req.QueryParams["ListenerPort"] = string(requests.NewInteger(d.Get("frontend_port").(int)))
	req.QueryParams["BackendServerPort"] = string(requests.NewInteger(d.Get("backend_port").(int)))
	req.QueryParams["Bandwidth"] = string(requests.NewInteger(d.Get("bandwidth").(int)))
	if groupId, ok := d.GetOk("server_group_id"); ok && groupId.(string) != "" {
		req.QueryParams["VServerGroupId"] = groupId.(string)
	}
	// acl status
	if aclStatus, ok := d.GetOk("acl_status"); ok && aclStatus.(string) != "" {
		req.QueryParams["AclStatus"] = aclStatus.(string)
	}
	// acl type
	if aclType, ok := d.GetOk("acl_type"); ok && aclType.(string) != "" {
		req.QueryParams["AclType"] = aclType.(string)
	}
	// acl id
	if aclId, ok := d.GetOk("acl_id"); ok && aclId.(string) != "" {
		req.QueryParams["AclId"] = aclId.(string)
	}

	return req

}
func buildHttpListenerArgs(d *schema.ResourceData, req *requests.CommonRequest) (*requests.CommonRequest, error) {
	stickySession := d.Get("sticky_session").(string)
	healthCheck := d.Get("health_check").(string)
	req.QueryParams["StickySession"] = stickySession
	req.QueryParams["HealthCheck"] = healthCheck

	if stickySession == string(OnFlag) {
		sessionType, ok := d.GetOk("sticky_session_type")
		if !ok || sessionType.(string) == "" {
			return req, fmt.Errorf("'sticky_session_type': required field is not set when the StickySession is %s.", OnFlag)
		} else {
			req.QueryParams["StickySessionType"] = sessionType.(string)

		}
		if sessionType.(string) == string(InsertStickySessionType) {
			if timeout, ok := d.GetOk("cookie_timeout"); !ok || timeout == 0 {
				return req, fmt.Errorf("'cookie_timeout': required field is not set when the StickySession is %s "+
					"and StickySessionType is %s.", OnFlag, InsertStickySessionType)
			} else {
				req.QueryParams["CookieTimeout"] = string(requests.NewInteger(timeout.(int)))
			}
		} else {
			if cookie, ok := d.GetOk("cookie"); !ok || cookie.(string) == "" {
				return req, fmt.Errorf("'cookie': required field is not set when the StickySession is %s "+
					"and StickySessionType is %s.", OnFlag, ServerStickySessionType)
			} else {
				req.QueryParams["Cookie"] = cookie.(string)
			}
		}
	}
	if healthCheck == string(OnFlag) {
		req.QueryParams["HealthCheckURI"] = d.Get("health_check_uri").(string)
		if port, ok := d.GetOk("health_check_connect_port"); !ok || port.(int) == 0 {
			return req, fmt.Errorf("'health_check_connect_port': required field is not set when the HealthCheck is %s.", OnFlag)
		} else {
			req.QueryParams["HealthCheckConnectPort"] = string(requests.NewInteger(port.(int)))
		}
		req.QueryParams["HealthyThreshold"] = string(requests.NewInteger(d.Get("healthy_threshold").(int)))
		req.QueryParams["UnhealthyThreshold"] = string(requests.NewInteger(d.Get("unhealthy_threshold").(int)))
		req.QueryParams["HealthCheckTimeout"] = string(requests.NewInteger(d.Get("health_check_timeout").(int)))
		req.QueryParams["HealthCheckInterval"] = string(requests.NewInteger(d.Get("health_check_interval").(int)))
		req.QueryParams["HealthCheckHttpCode"] = d.Get("health_check_http_code").(string)
	}
	return req, nil
}

func parseListenerId(d *schema.ResourceData, meta interface{}) (string, string, int, error) {
	parts := strings.Split(d.Id(), ":")
	port, err := strconv.Atoi(parts[1])
	if err != nil {
		return "", "", 0, fmt.Errorf("Parsing SlbListener's id got an error: %#v", err)
	}
	loadBalancer, err := meta.(*AliyunClient).DescribeLoadBalancerAttribute(parts[0])
	if err != nil {
		return "", "", 0, err
	}
	for _, portAndProtocol := range loadBalancer.ListenerPortsAndProtocol.ListenerPortAndProtocol {
		if portAndProtocol.ListenerPort == port {
			return loadBalancer.LoadBalancerId, portAndProtocol.ListenerProtocol, port, nil
		}
	}
	return "", "", 0, nil
}

func readListenerAttribute(d *schema.ResourceData, protocol string, listen map[string]interface{}, err error) *resource.RetryError {
	if err != nil {
		if IsExceptedError(err, ListenerNotFound) {
			d.SetId("")
			return nil
		}
		if IsExceptedErrors(err, SlbIsBusy) {
			return resource.RetryableError(fmt.Errorf("DescribeLoadBalancer%sListenerAttribute timeout and got an error: %#v", strings.ToUpper(protocol), err))
		}
		return resource.NonRetryableError(fmt.Errorf("DescribeLoadBalancer%sListenerAttribute got an error: %#v", strings.ToUpper(protocol), err))
	}

	if port, ok := listen["ListenerPort"]; ok && port.(float64) > 0 {
		readListener(d, listen)
	} else {
		d.SetId("")
	}
	return nil
}

func readListener(d *schema.ResourceData, listener map[string]interface{}) {
	if val, ok := listener["BackendServerPort"]; ok {
		d.Set("backend_port", val.(float64))
	}
	if val, ok := listener["ListenerPort"]; ok {
		d.Set("frontend_port", val.(float64))
	}
	if val, ok := listener["Bandwidth"]; ok {
		d.Set("bandwidth", val.(float64))
	}
	if val, ok := listener["Scheduler"]; ok {
		d.Set("scheduler", val.(string))
	}
	if val, ok := listener["VServerGroupId"]; ok {
		d.Set("server_group_id", val.(string))
	}
	if val, ok := listener["AclStatus"]; ok {
		d.Set("acl_status", val.(string))
	}
	if val, ok := listener["AclType"]; ok {
		d.Set("acl_type", val.(string))
	}
	if val, ok := listener["AclId"]; ok {
		d.Set("acl_id", val.(string))
	}
	if val, ok := listener["HealthCheck"]; ok {
		d.Set("health_check", val.(string))
	}
	if val, ok := listener["StickySession"]; ok {
		d.Set("sticky_session", val.(string))
	}
	if val, ok := listener["StickySessionType"]; ok {
		d.Set("sticky_session_type", val.(string))
	}
	if val, ok := listener["CookieTimeout"]; ok {
		d.Set("cookie_timeout", val.(float64))
	}
	if val, ok := listener["Cookie"]; ok {
		d.Set("cookie", val.(string))
	}
	if val, ok := listener["PersistenceTimeout"]; ok {
		d.Set("persistence_timeout", val.(float64))
	}
	if val, ok := listener["HealthCheckType"]; ok {
		d.Set("health_check_type", val.(string))
	}
	if val, ok := listener["HealthCheckDomain"]; ok {
		d.Set("health_check_domain", val.(string))
	}
	if val, ok := listener["HealthCheckConnectPort"]; ok {
		d.Set("health_check_connect_port", val.(float64))
	}
	if val, ok := listener["HealthCheckURI"]; ok {
		d.Set("health_check_uri", val.(string))
	}
	if val, ok := listener["HealthyThreshold"]; ok {
		d.Set("healthy_threshold", val.(float64))
	}
	if val, ok := listener["UnhealthyThreshold"]; ok {
		d.Set("unhealthy_threshold", val.(float64))
	}
	if val, ok := listener["HealthCheckTimeout"]; ok {
		d.Set("health_check_timeout", val.(float64))
	}
	if val, ok := listener["HealthCheckConnectTimeout"]; ok {
		d.Set("health_check_timeout", val.(float64))
	}
	if val, ok := listener["HealthCheckInterval"]; ok {
		d.Set("health_check_interval", val.(float64))
	}
	if val, ok := listener["HealthCheckHttpCode"]; ok {
		d.Set("health_check_http_code", val.(string))
	}
	if val, ok := listener["ServerCertificateId"]; ok {
		d.Set("ssl_certificate_id", val.(string))
	}
	if val, ok := listener["Gzip"]; ok {
		d.Set("gzip", val.(string) == string(OnFlag))
	}
	xff := make(map[string]interface{})
	if val, ok := listener["XForwardedFor"]; ok {
		xff["retrive_client_ip"] = val.(string) == string(OnFlag)
	}
	if val, ok := listener["XForwardedFor_SLBIP"]; ok {
		xff["retrive_slb_ip"] = val.(string) == string(OnFlag)
	}
	if val, ok := listener["XForwardedFor_SLBID"]; ok {
		xff["retrive_slb_id"] = val.(string) == string(OnFlag)
	}
	if val, ok := listener["XForwardedFor_proto"]; ok {
		xff["retrive_slb_proto"] = val.(string) == string(OnFlag)
	}

	if len(xff) > 0 {
		d.Set("x_forwarded_for", []map[string]interface{}{xff})
	}

	return
}

func ensureListenerAbsent(d *schema.ResourceData, protocol string, listener map[string]interface{}, err error) *resource.RetryError {

	if err != nil {
		if IsExceptedError(err, ListenerNotFound) {
			d.SetId("")
			return nil
		}
		if IsExceptedErrors(err, SlbIsBusy) {
			return resource.RetryableError(fmt.Errorf("While deleting listener, DescribeLoadBalancer%sListenerAttribute timeout and got an error: %#v", protocol, err))
		}
		return resource.NonRetryableError(fmt.Errorf("While deleting listener, DescribeLoadBalancer%sListenerAttribute got an error: %#v", protocol, err))
	}
	if port, ok := listener["ListenerPort"]; ok && port.(float64) > 0 {
		return resource.RetryableError(fmt.Errorf("Delete load balancer listener timeout and got an error: %#v.", err))
	}
	d.SetId("")
	return nil
}
