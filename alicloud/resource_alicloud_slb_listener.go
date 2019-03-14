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
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
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
			"load_balancer_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"frontend_port": {
				Type:         schema.TypeInt,
				ValidateFunc: validateInstancePort,
				Required:     true,
				ForceNew:     true,
			},
			"lb_port": {
				Type:       schema.TypeInt,
				Optional:   true,
				Deprecated: "Field 'lb_port' has been deprecated, and using 'frontend_port' to replace.",
			},

			"backend_port": {
				Type:         schema.TypeInt,
				ValidateFunc: validateInstancePort,
				Required:     true,
				ForceNew:     true,
			},

			"instance_port": {
				Type:       schema.TypeInt,
				Optional:   true,
				Deprecated: "Field 'instance_port' has been deprecated, and using 'backend_port' to replace.",
			},

			"lb_protocol": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "Field 'lb_protocol' has been deprecated, and using 'protocol' to replace.",
			},

			"protocol": {
				Type:         schema.TypeString,
				ValidateFunc: validateInstanceProtocol,
				Required:     true,
				ForceNew:     true,
			},

			"bandwidth": {
				Type:         schema.TypeInt,
				ValidateFunc: validateSlbListenerBandwidth,
				Required:     true,
			},
			"scheduler": {
				Type:         schema.TypeString,
				ValidateFunc: validateSlbListenerScheduler,
				Optional:     true,
				Default:      WRRScheduler,
			},
			"server_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"acl_status": {
				Type:         schema.TypeString,
				ValidateFunc: validateAllowedStringValue([]string{string(OnFlag), string(OffFlag)}),
				Optional:     true,
				Default:      OffFlag,
			},
			"acl_type": {
				Type:             schema.TypeString,
				ValidateFunc:     validateAllowedStringValue([]string{string(AclTypeBlack), string(AclTypeWhite)}),
				Optional:         true,
				DiffSuppressFunc: slbAclDiffSuppressFunc,
			},
			"acl_id": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: slbAclDiffSuppressFunc,
			},
			//http & https
			"sticky_session": {
				Type:             schema.TypeString,
				ValidateFunc:     validateAllowedStringValue([]string{string(OnFlag), string(OffFlag)}),
				Optional:         true,
				Default:          OffFlag,
				DiffSuppressFunc: httpHttpsDiffSuppressFunc,
			},
			//http & https
			"sticky_session_type": {
				Type: schema.TypeString,
				ValidateFunc: validateAllowedStringValue([]string{
					string(InsertStickySessionType),
					string(ServerStickySessionType)}),
				Optional:         true,
				DiffSuppressFunc: stickySessionTypeDiffSuppressFunc,
			},
			//http & https
			"cookie_timeout": {
				Type:             schema.TypeInt,
				ValidateFunc:     validateSlbListenerCookieTimeout,
				Optional:         true,
				DiffSuppressFunc: cookieTimeoutDiffSuppressFunc,
			},
			//http & https
			"cookie": {
				Type:             schema.TypeString,
				ValidateFunc:     validateSlbListenerCookie,
				Optional:         true,
				DiffSuppressFunc: cookieDiffSuppressFunc,
			},
			//tcp & udp
			"persistence_timeout": {
				Type:             schema.TypeInt,
				ValidateFunc:     validateSlbListenerPersistenceTimeout,
				Optional:         true,
				Default:          0,
				DiffSuppressFunc: tcpUdpDiffSuppressFunc,
			},
			//http & https
			"health_check": {
				Type:             schema.TypeString,
				ValidateFunc:     validateAllowedStringValue([]string{string(OnFlag), string(OffFlag)}),
				Optional:         true,
				Default:          OnFlag,
				DiffSuppressFunc: httpHttpsDiffSuppressFunc,
			},
			//tcp
			"health_check_type": {
				Type: schema.TypeString,
				ValidateFunc: validateAllowedStringValue([]string{
					string(TCPHealthCheckType),
					string(HTTPHealthCheckType)}),
				Optional:         true,
				Default:          TCPHealthCheckType,
				DiffSuppressFunc: healthCheckTypeDiffSuppressFunc,
			},
			//http & https & tcp
			"health_check_domain": {
				Type:             schema.TypeString,
				ValidateFunc:     validateSlbListenerHealthCheckDomain,
				Optional:         true,
				DiffSuppressFunc: httpHttpsTcpDiffSuppressFunc,
			},
			//http & https & tcp
			"health_check_uri": {
				Type:             schema.TypeString,
				ValidateFunc:     validateSlbListenerHealthCheckUri,
				Optional:         true,
				Default:          "/",
				DiffSuppressFunc: httpHttpsTcpDiffSuppressFunc,
			},
			"health_check_connect_port": {
				Type:             schema.TypeInt,
				ValidateFunc:     validateSlbListenerHealthCheckConnectPort,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: healthCheckDiffSuppressFunc,
			},
			"healthy_threshold": {
				Type:             schema.TypeInt,
				ValidateFunc:     validateIntegerInRange(1, 10),
				Optional:         true,
				Default:          3,
				DiffSuppressFunc: healthCheckDiffSuppressFunc,
			},
			"unhealthy_threshold": {
				Type:             schema.TypeInt,
				ValidateFunc:     validateIntegerInRange(1, 10),
				Optional:         true,
				Default:          3,
				DiffSuppressFunc: healthCheckDiffSuppressFunc,
			},

			"health_check_timeout": {
				Type:             schema.TypeInt,
				ValidateFunc:     validateIntegerInRange(1, 300),
				Optional:         true,
				Default:          5,
				DiffSuppressFunc: healthCheckDiffSuppressFunc,
			},
			"health_check_interval": {
				Type:             schema.TypeInt,
				ValidateFunc:     validateIntegerInRange(1, 50),
				Optional:         true,
				Default:          2,
				DiffSuppressFunc: healthCheckDiffSuppressFunc,
			},
			//http & https & tcp
			"health_check_http_code": {
				Type: schema.TypeString,
				ValidateFunc: validateAllowedSplitStringValue([]string{
					string(HTTP_2XX), string(HTTP_3XX), string(HTTP_4XX), string(HTTP_5XX)}, ","),
				Optional:         true,
				Default:          HTTP_2XX,
				DiffSuppressFunc: httpHttpsTcpDiffSuppressFunc,
			},
			//https
			"ssl_certificate_id": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: sslCertificateIdDiffSuppressFunc,
			},

			//http, https
			"gzip": {
				Type:             schema.TypeBool,
				Optional:         true,
				Default:          true,
				DiffSuppressFunc: httpHttpsDiffSuppressFunc,
			},
			"x_forwarded_for": {
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
			//tcp
			"established_timeout": {
				Type:             schema.TypeInt,
				ValidateFunc:     validateIntegerInRange(10, 900),
				Optional:         true,
				Default:          900,
				DiffSuppressFunc: establishedTimeoutDiffSuppressFunc,
			},

			//http & https
			"idle_timeout": {
				Type:             schema.TypeInt,
				ValidateFunc:     validateIntegerInRange(1, 60),
				Optional:         true,
				Default:          15,
				DiffSuppressFunc: httpHttpsDiffSuppressFunc,
			},

			//http & https
			"request_timeout": {
				Type:             schema.TypeInt,
				ValidateFunc:     validateIntegerInRange(1, 180),
				Optional:         true,
				Default:          60,
				DiffSuppressFunc: httpHttpsDiffSuppressFunc,
			},

			//https
			"enable_http2": {
				Type:             schema.TypeString,
				ValidateFunc:     validateAllowedStringValue([]string{string(OnFlag), string(OffFlag)}),
				Optional:         true,
				Default:          OnFlag,
				DiffSuppressFunc: httpsDiffSuppressFunc,
			},

			//https
			"tls_cipher_policy": {
				Type:    schema.TypeString,
				Default: string(TlsCipherPolicy_1_0),
				ValidateFunc: validateAllowedStringValue([]string{string(TlsCipherPolicy_1_0),
					string(TlsCipherPolicy_1_1), string(TlsCipherPolicy_1_2), string(TlsCipherPolicy_1_2_STRICT)}),
				Optional:         true,
				DiffSuppressFunc: httpsDiffSuppressFunc,
			},
		},
	}
}

func resourceAliyunSlbListenerCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	slbService := SlbService{client}

	protocol := d.Get("protocol").(string)
	lb_id := d.Get("load_balancer_id").(string)
	frontend := d.Get("frontend_port").(int)

	req, err := buildListenerCommonArgs(d, meta)
	if err != nil {
		return WrapError(err)
	}
	req.ApiName = fmt.Sprintf("CreateLoadBalancer%sListener", strings.ToUpper(protocol))

	if Protocol(protocol) == Http || Protocol(protocol) == Https {
		reqHttp, err := buildHttpListenerArgs(d, req)
		if err != nil {
			return WrapError(err)
		}
		req = reqHttp
		if Protocol(protocol) == Https {
			ssl_id, ok := d.GetOk("ssl_certificate_id")
			if !ok || ssl_id.(string) == "" {
				return WrapError(Error(`'ssl_certificate_id': required field is not set when the protocol is 'https'.`))
			}
			req.QueryParams["ServerCertificateId"] = ssl_id.(string)
		}
	}
	raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
		return slbClient.ProcessCommonRequest(req)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "slb_listener", req.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(req.GetActionName(), raw)
	d.SetId(lb_id + ":" + strconv.Itoa(frontend))

	if err := slbService.WaitForListener(lb_id, frontend, Protocol(protocol), Stopped, DefaultTimeout); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), req.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	reqStart := slb.CreateStartLoadBalancerListenerRequest()
	reqStart.LoadBalancerId = lb_id
	reqStart.ListenerPort = requests.NewInteger(frontend)
	raw, err = client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
		return slbClient.StartLoadBalancerListener(reqStart)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "slb_listener", reqStart.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(reqStart.GetActionName(), raw)
	if err = slbService.WaitForListener(lb_id, frontend, Protocol(protocol), Running, DefaultTimeout); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "slb_listener", reqStart.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	return resourceAliyunSlbListenerUpdate(d, meta)
}

func resourceAliyunSlbListenerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	slbService := SlbService{client}

	lb_id, protocol, port, err := parseListenerId(d, meta)
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("protocol", protocol)
	d.Set("load_balancer_id", lb_id)

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		listener, err := slbService.DescribeLoadBalancerListenerAttribute(lb_id, port, Protocol(protocol))
		return readListenerAttribute(d, protocol, listener, err)
	})
}

func resourceAliyunSlbListenerUpdate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	protocol := Protocol(d.Get("protocol").(string))

	d.Partial(true)

	commonArgs, err := buildListenerCommonArgs(d, meta)
	if err != nil {
		return WrapError(err)
	}
	commonArgs.ApiName = fmt.Sprintf("SetLoadBalancer%sListenerAttribute", strings.ToUpper(string(protocol)))

	update := false
	if d.HasChange("scheduler") {
		commonArgs.QueryParams["Scheduler"] = d.Get("scheduler").(string)
		d.SetPartial("scheduler")
		update = true
	}

	if d.HasChange("server_group_id") {
		serverGroupId := d.Get("server_group_id").(string)
		if serverGroupId != "" {
			commonArgs.QueryParams["VServerGroup"] = string(OnFlag)
			commonArgs.QueryParams["VServerGroupId"] = d.Get("server_group_id").(string)
		} else {
			commonArgs.QueryParams["VServerGroup"] = string(OffFlag)
		}
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
		return WrapError(err)
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

	// http https
	if d.HasChange("idle_timeout") {
		d.SetPartial("idle_timeout")
		update = true
	}

	// http https
	if d.HasChange("request_timeout") {
		d.SetPartial("request_timeout")
		update = true
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

	// tcp
	if d.HasChange("established_timeout") {
		tcpArgs.QueryParams["EstablishedTimeout"] = string(requests.NewInteger(d.Get("established_timeout").(int)))
		d.SetPartial("established_timeout")
		update = true
	}

	// https
	httpsArgs := httpArgs
	if protocol == Https {
		ssl_id, ok := d.GetOk("ssl_certificate_id")
		if !ok && ssl_id == "" {
			return WrapError(Error("'ssl_certificate_id': required field is not set when the protocol is 'https'."))
		}

		httpsArgs.QueryParams["ServerCertificateId"] = ssl_id.(string)
		if d.HasChange("ssl_certificate_id") {
			d.SetPartial("ssl_certificate_id")
			update = true
		}

		if d.HasChange("enable_http2") {
			httpsArgs.QueryParams["EnableHttp2"] = d.Get("enable_http2").(string)
			d.SetPartial("enable_http2")
			update = true
		}

		if d.HasChange("tls_cipher_policy") {
			// spec changes check, can not be updated when load balancer instance is "Shared-Performance".
			slbService := SlbService{client}
			loadBalancer, _ := slbService.DescribeLoadBalancerAttribute(d.Get("load_balancer_id").(string))
			spec := loadBalancer.LoadBalancerSpec
			if spec == "" {
				if !d.IsNewResource() || string(TlsCipherPolicy_1_0) != d.Get("tls_cipher_policy").(string) {
					return WrapError(Error("Currently the param \"tls_cipher_policy\" can not be updated when load balancer instance is \"Shared-Performance\"."))
				}
			} else {
				httpsArgs.QueryParams["TLSCipherPolicy"] = d.Get("tls_cipher_policy").(string)
				d.SetPartial("tls_cipher_policy")
				update = true
			}
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
		raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
			return slbClient.ProcessCommonRequest(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw)
	}

	d.Partial(false)

	return resourceAliyunSlbListenerRead(d, meta)
}

func resourceAliyunSlbListenerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	slbService := SlbService{client}

	lb_id, protocol, port, err := parseListenerId(d, meta)
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	req := slb.CreateDeleteLoadBalancerListenerRequest()
	req.LoadBalancerId = lb_id
	req.ListenerPort = requests.NewInteger(port)
	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
			return slbClient.DeleteLoadBalancerListener(req)
		})

		if err != nil {
			if IsExceptedErrors(err, SlbIsBusy) {
				return resource.RetryableError(WrapErrorf(err, DefaultErrorMsg, d.Id(), req.GetActionName(), AlibabaCloudSdkGoERROR))
			}
			return resource.NonRetryableError(WrapErrorf(err, DefaultErrorMsg, d.Id(), req.GetActionName(), AlibabaCloudSdkGoERROR))
		}

		listener, err := slbService.DescribeLoadBalancerListenerAttribute(lb_id, port, Protocol(protocol))
		return ensureListenerAbsent(d, protocol, listener, err)
	})
}

func buildListenerCommonArgs(d *schema.ResourceData, meta interface{}) (*requests.CommonRequest, error) {
	client := meta.(*connectivity.AliyunClient)
	slbService := SlbService{client}

	req, err := slbService.BuildSlbCommonRequest()
	if err != nil {
		return req, WrapError(err)
	}
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

	return req, nil

}
func buildHttpListenerArgs(d *schema.ResourceData, req *requests.CommonRequest) (*requests.CommonRequest, error) {
	stickySession := d.Get("sticky_session").(string)
	healthCheck := d.Get("health_check").(string)
	req.QueryParams["StickySession"] = stickySession
	req.QueryParams["HealthCheck"] = healthCheck

	if stickySession == string(OnFlag) {
		sessionType, ok := d.GetOk("sticky_session_type")
		if !ok || sessionType.(string) == "" {
			return req, WrapError(fmt.Errorf("'sticky_session_type': required field is not set when the StickySession is %s.", OnFlag))
		} else {
			req.QueryParams["StickySessionType"] = sessionType.(string)

		}
		if sessionType.(string) == string(InsertStickySessionType) {
			if timeout, ok := d.GetOk("cookie_timeout"); !ok || timeout == 0 {
				return req, WrapError(fmt.Errorf("'cookie_timeout': required field is not set when the StickySession is %s "+
					"and StickySessionType is %s.", OnFlag, InsertStickySessionType))
			} else {
				req.QueryParams["CookieTimeout"] = string(requests.NewInteger(timeout.(int)))
			}
		} else {
			if cookie, ok := d.GetOk("cookie"); !ok || cookie.(string) == "" {
				return req, WrapError(fmt.Errorf("'cookie': required field is not set when the StickySession is %s "+
					"and StickySessionType is %s.", OnFlag, ServerStickySessionType))
			} else {
				req.QueryParams["Cookie"] = cookie.(string)
			}
		}
	}
	if healthCheck == string(OnFlag) {
		req.QueryParams["HealthCheckURI"] = d.Get("health_check_uri").(string)
		if port, ok := d.GetOk("health_check_connect_port"); !ok || port.(int) == 0 {
			return req, WrapError(fmt.Errorf("'health_check_connect_port': required field is not set when the HealthCheck is %s.", OnFlag))
		} else {
			req.QueryParams["HealthCheckConnectPort"] = string(requests.NewInteger(port.(int)))
		}
		req.QueryParams["HealthyThreshold"] = string(requests.NewInteger(d.Get("healthy_threshold").(int)))
		req.QueryParams["UnhealthyThreshold"] = string(requests.NewInteger(d.Get("unhealthy_threshold").(int)))
		req.QueryParams["HealthCheckTimeout"] = string(requests.NewInteger(d.Get("health_check_timeout").(int)))
		req.QueryParams["HealthCheckInterval"] = string(requests.NewInteger(d.Get("health_check_interval").(int)))
		req.QueryParams["HealthCheckHttpCode"] = d.Get("health_check_http_code").(string)

		req.QueryParams["IdleTimeout"] = string(requests.NewInteger(d.Get("idle_timeout").(int)))
		req.QueryParams["RequestTimeout"] = string(requests.NewInteger(d.Get("request_timeout").(int)))
	}
	return req, nil
}

func parseListenerId(d *schema.ResourceData, meta interface{}) (string, string, int, error) {
	client := meta.(*connectivity.AliyunClient)
	slbService := SlbService{client}

	parts := strings.Split(d.Id(), ":")
	port, err := strconv.Atoi(parts[1])
	if err != nil {
		return "", "", 0, WrapError(err)
	}
	loadBalancer, err := slbService.DescribeLoadBalancerAttribute(parts[0])
	if err != nil {
		return "", "", 0, WrapError(err)
	}
	for _, portAndProtocol := range loadBalancer.ListenerPortsAndProtocol.ListenerPortAndProtocol {
		if portAndProtocol.ListenerPort == port {
			return loadBalancer.LoadBalancerId, portAndProtocol.ListenerProtocol, port, nil
		}
	}
	return "", "", 0, GetNotFoundErrorFromString(GetNotFoundMessage("Listener", d.Id()))
}

func readListenerAttribute(d *schema.ResourceData, protocol string, listen map[string]interface{}, err error) *resource.RetryError {
	if err != nil {
		if IsExceptedError(err, ListenerNotFound) {
			d.SetId("")
			return nil
		}
		if IsExceptedErrors(err, SlbIsBusy) {
			return resource.RetryableError(WrapError(err))
		}
		return resource.NonRetryableError(WrapError(err))
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
	if val, ok := listener["EstablishedTimeout"]; ok {
		d.Set("established_timeout", val.(float64))
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

	if val, ok := listener["EnableHttp2"]; ok {
		d.Set("enable_http2", val.(string))
	}

	if val, ok := listener["TLSCipherPolicy"]; ok {
		d.Set("tls_cipher_policy", val.(string))
	}

	if val, ok := listener["IdleTimeout"]; ok {
		d.Set("idle_timeout", val.(float64))
	}

	if val, ok := listener["RequestTimeout"]; ok {
		d.Set("request_timeout", val.(float64))
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
			return resource.RetryableError(WrapError(err))
		}
		return resource.NonRetryableError(WrapError(err))
	}
	if port, ok := listener["ListenerPort"]; ok && port.(float64) > 0 {
		return resource.RetryableError(WrapError(err))
	}
	d.SetId("")
	return nil
}
