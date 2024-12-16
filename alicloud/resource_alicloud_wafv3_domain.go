package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudWafv3Domain() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudWafv3DomainCreate,
		Read:   resourceAlicloudWafv3DomainRead,
		Update: resourceAlicloudWafv3DomainUpdate,
		Delete: resourceAlicloudWafv3DomainDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"access_type": {
				Optional:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"share"}, false),
			},
			"domain": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"instance_id": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"listen": {
				Required: true,
				MaxItems: 1,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cert_id": {
							Optional: true,
							Type:     schema.TypeString,
						},
						"cipher_suite": {
							Optional:     true,
							Type:         schema.TypeInt,
							ValidateFunc: validation.IntInSlice([]int{1, 2, 99}),
						},
						"custom_ciphers": {
							Optional: true,
							Type:     schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"enable_tlsv3": {
							Optional: true,
							Type:     schema.TypeBool,
						},
						"exclusive_ip": {
							Optional: true,
							Type:     schema.TypeBool,
						},
						"focus_https": {
							Optional: true,
							Type:     schema.TypeBool,
						},
						"http2_enabled": {
							Optional: true,
							Type:     schema.TypeBool,
						},
						"http_ports": {
							Optional: true,
							Type:     schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
						},
						"https_ports": {
							Optional: true,
							Type:     schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
						},
						"ipv6_enabled": {
							Optional: true,
							Type:     schema.TypeBool,
						},
						"protection_resource": {
							Optional: true,
							Type:     schema.TypeString,
						},
						"tls_version": {
							Optional:     true,
							Type:         schema.TypeString,
							ValidateFunc: validation.StringInSlice([]string{"tlsv1", "tlsv1.1", "tlsv1.2"}, false),
						},
						"xff_header_mode": {
							Optional:     true,
							Type:         schema.TypeInt,
							ValidateFunc: validation.IntInSlice([]int{0, 1, 2}),
						},
						"xff_headers": {
							Optional: true,
							Type:     schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"redirect": {
				Required: true,
				MaxItems: 1,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"backends": {
							Optional: true,
							Type:     schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"connect_timeout": {
							Optional: true,
							Type:     schema.TypeInt,
						},
						"focus_http_backend": {
							Optional: true,
							Type:     schema.TypeBool,
						},
						"keepalive": {
							Optional: true,
							Type:     schema.TypeBool,
						},
						"keepalive_requests": {
							Optional: true,
							Type:     schema.TypeInt,
						},
						"keepalive_timeout": {
							Optional: true,
							Type:     schema.TypeInt,
						},
						"loadbalance": {
							Required:     true,
							Type:         schema.TypeString,
							ValidateFunc: validation.StringInSlice([]string{"iphash", "roundRobin", "leastTime"}, false),
						},
						"read_timeout": {
							Optional: true,
							Type:     schema.TypeInt,
						},
						"request_headers": {
							Optional: true,
							Type:     schema.TypeSet,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Optional: true,
										Type:     schema.TypeString,
									},
									"value": {
										Optional: true,
										Type:     schema.TypeString,
									},
								},
							},
						},
						"retry": {
							Optional: true,
							Type:     schema.TypeBool,
						},
						"sni_enabled": {
							Optional: true,
							Type:     schema.TypeBool,
						},
						"sni_host": {
							Optional: true,
							Type:     schema.TypeString,
						},
						"write_timeout": {
							Optional: true,
							Type:     schema.TypeInt,
						},
					},
				},
			},
			"resource_manager_resource_group_id": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"status": {
				Computed: true,
				Type:     schema.TypeString,
			},
		},
	}
}

func resourceAlicloudWafv3DomainCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	wafOpenapiService := WafOpenapiService{client}

	request := map[string]interface{}{
		"RegionId": client.RegionId,
	}
	var err error

	if v, ok := d.GetOk("access_type"); ok {
		request["AccessType"] = v
	}
	if v, ok := d.GetOk("domain"); ok {
		request["Domain"] = v
	}
	if v, ok := d.GetOk("instance_id"); ok {
		request["InstanceId"] = v
	}

	listenMap := make(map[string]interface{}, 0)
	for _, v := range d.Get("listen").([]interface{}) {
		listenObj := v.(map[string]interface{})
		listenMap["CertId"] = listenObj["cert_id"]
		listenMap["CipherSuite"] = listenObj["cipher_suite"]
		listenMap["CustomCiphers"] = listenObj["custom_ciphers"].([]interface{})
		listenMap["EnableTLSv3"] = listenObj["enable_tlsv3"]
		listenMap["ExclusiveIp"] = listenObj["exclusive_ip"]
		listenMap["FocusHttps"] = listenObj["focus_https"]
		listenMap["Http2Enabled"] = listenObj["http2_enabled"]
		listenMap["HttpPorts"] = listenObj["http_ports"].([]interface{})
		listenMap["HttpsPorts"] = listenObj["https_ports"].([]interface{})
		listenMap["IPv6Enabled"] = listenObj["ipv6_enabled"]
		listenMap["ProtectionResource"] = listenObj["protection_resource"]
		listenMap["TLSVersion"] = listenObj["tls_version"]
		listenMap["XffHeaderMode"] = listenObj["xff_header_mode"]
		listenMap["XffHeaders"] = listenObj["xff_headers"]
	}
	request["Listen"], _ = convertMaptoJsonString(listenMap)

	redirectMap := make(map[string]interface{}, 0)
	for _, v := range d.Get("redirect").([]interface{}) {
		redirectObj := v.(map[string]interface{})

		redirectMap["Backends"] = redirectObj["backends"].([]interface{})
		redirectMap["ConnectTimeout"] = redirectObj["connect_timeout"]
		redirectMap["FocusHttpBackend"] = redirectObj["focus_http_backend"]
		redirectMap["Keepalive"] = redirectObj["keepalive"]
		redirectMap["KeepaliveRequests"] = redirectObj["keepalive_requests"]
		redirectMap["KeepaliveTimeout"] = redirectObj["keepalive_timeout"]
		redirectMap["Retry"] = redirectObj["retry"]
		redirectMap["SniEnabled"] = redirectObj["sni_enabled"]
		redirectMap["SniHost"] = redirectObj["sni_host"]
		redirectMap["Loadbalance"] = redirectObj["loadbalance"]
		redirectMap["WriteTimeout"] = redirectObj["write_timeout"]
		redirectMap["ReadTimeout"] = redirectObj["read_timeout"]

		requestHeaderMap := make([]map[string]interface{}, 0)
		if v, ok := redirectObj["request_headers"]; ok {
			for _, requestHeader := range v.(*schema.Set).List() {
				requestHeaderObj := requestHeader.(map[string]interface{})
				requestHeaderMap = append(requestHeaderMap, map[string]interface{}{
					"Key":   requestHeaderObj["key"],
					"Value": requestHeaderObj["value"],
				})
			}
			redirectMap["RequestHeaders"] = requestHeaderMap
		}
	}
	request["Redirect"], _ = convertMaptoJsonString(redirectMap)

	var response map[string]interface{}
	action := "CreateDomain"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = client.RpcPost("waf-openapi", "2021-10-01", action, nil, request, false)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_wafv3_domain", action, AlibabaCloudSdkGoERROR)
	}
	domainValue, err := jsonpath.Get("$.DomainInfo.Domain", response)
	if err != nil || domainValue == nil {
		return WrapErrorf(err, IdMsg, "alicloud_wafv3_domain")
	}

	d.SetId(fmt.Sprint(request["InstanceId"], ":", domainValue))

	stateConf := BuildStateConf([]string{}, []string{"1"}, d.Timeout(schema.TimeoutDelete), 5*time.Second, wafOpenapiService.Wafv3DomainStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudWafv3DomainRead(d, meta)
}

func resourceAlicloudWafv3DomainRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	wafOpenapiService := WafOpenapiService{client}

	object, err := wafOpenapiService.DescribeWafv3Domain(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_wafv3_domain wafOpenapiService.DescribeWafv3Domain Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	d.Set("instance_id", parts[0])
	d.Set("domain", parts[1])
	listen19Maps := make([]map[string]interface{}, 0)
	listen19Map := make(map[string]interface{})
	listen19Raw := object["Listen"].(map[string]interface{})
	listen19Map["cert_id"] = listen19Raw["CertId"]
	listen19Map["cipher_suite"] = listen19Raw["CipherSuite"]
	if v, ok := listen19Raw["CustomCiphers"]; ok {
		listen19Map["custom_ciphers"] = v.([]interface{})
	}
	listen19Map["enable_tlsv3"] = listen19Raw["EnableTLSv3"]
	listen19Map["exclusive_ip"] = listen19Raw["ExclusiveIp"]
	listen19Map["focus_https"] = listen19Raw["FocusHttps"]
	listen19Map["http2_enabled"] = listen19Raw["Http2Enabled"]
	if v, ok := listen19Raw["HttpPorts"]; ok {
		listen19Map["http_ports"] = v.([]interface{})
	}
	if v, ok := listen19Raw["HttpsPorts"]; ok {
		listen19Map["https_ports"] = v.([]interface{})
	}
	listen19Map["ipv6_enabled"] = listen19Raw["IPv6Enabled"]
	listen19Map["protection_resource"] = listen19Raw["ProtectionResource"]
	listen19Map["tls_version"] = listen19Raw["TLSVersion"]
	listen19Map["xff_header_mode"] = listen19Raw["XffHeaderMode"]
	if v, ok := listen19Raw["XffHeaders"]; ok {
		listen19Map["xff_headers"] = v.([]interface{})
	}
	listen19Maps = append(listen19Maps, listen19Map)
	d.Set("listen", listen19Maps)
	redirect81Maps := make([]map[string]interface{}, 0)
	redirect81Map := make(map[string]interface{})
	redirect81Raw := object["Redirect"].(map[string]interface{})
	if v, ok := redirect81Raw["AllBackends"]; ok {
		redirect81Map["backends"] = v.([]interface{})
	}

	redirect81Map["connect_timeout"] = redirect81Raw["ConnectTimeout"]
	redirect81Map["focus_http_backend"] = redirect81Raw["FocusHttpBackend"]
	redirect81Map["keepalive"] = redirect81Raw["Keepalive"]
	redirect81Map["keepalive_requests"] = redirect81Raw["KeepaliveRequests"]
	redirect81Map["keepalive_timeout"] = redirect81Raw["KeepaliveTimeout"]
	redirect81Map["loadbalance"] = redirect81Raw["Loadbalance"]
	redirect81Map["read_timeout"] = redirect81Raw["ReadTimeout"]
	requestHeaders81Maps := make([]map[string]interface{}, 0)
	if v, ok := redirect81Raw["RequestHeaders"]; ok && v != nil {
		for _, value1 := range v.([]interface{}) {
			requestHeaders81 := value1.(map[string]interface{})
			requestHeaders81Map := make(map[string]interface{})
			requestHeaders81Map["key"] = requestHeaders81["Key"]
			requestHeaders81Map["value"] = requestHeaders81["Value"]
			requestHeaders81Maps = append(requestHeaders81Maps, requestHeaders81Map)
		}
	}
	redirect81Map["request_headers"] = requestHeaders81Maps
	redirect81Map["retry"] = redirect81Raw["Retry"]
	redirect81Map["sni_enabled"] = redirect81Raw["SniEnabled"]
	redirect81Map["sni_host"] = redirect81Raw["SniHost"]
	redirect81Map["write_timeout"] = redirect81Raw["WriteTimeout"]
	redirect81Maps = append(redirect81Maps, redirect81Map)
	d.Set("redirect", redirect81Maps)
	resourceManagerResourceGroupId52 := object["ResourceManagerResourceGroupId"]
	d.Set("resource_manager_resource_group_id", resourceManagerResourceGroupId52)
	d.Set("status", fmt.Sprint(object["Status"]))

	return nil
}

func resourceAlicloudWafv3DomainUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	wafOpenapiService := WafOpenapiService{client}
	var response map[string]interface{}
	var err error
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"InstanceId": parts[0],
		"Domain":     parts[1],
		"RegionId":   client.RegionId,
	}

	if v, ok := d.GetOk("access_type"); ok {
		request["AccessType"] = v
	}

	redirectMap := make(map[string]interface{}, 0)
	for _, v := range d.Get("redirect").([]interface{}) {
		redirectObj := v.(map[string]interface{})

		redirectMap["Backends"] = redirectObj["backends"].([]interface{})
		redirectMap["ConnectTimeout"] = redirectObj["connect_timeout"]
		redirectMap["FocusHttpBackend"] = redirectObj["focus_http_backend"]
		redirectMap["Keepalive"] = redirectObj["keepalive"]
		redirectMap["KeepaliveRequests"] = redirectObj["keepalive_requests"]
		redirectMap["KeepaliveTimeout"] = redirectObj["keepalive_timeout"]
		redirectMap["Retry"] = redirectObj["retry"]
		redirectMap["SniEnabled"] = redirectObj["sni_enabled"]
		redirectMap["SniHost"] = redirectObj["sni_host"]
		redirectMap["Loadbalance"] = redirectObj["loadbalance"]
		redirectMap["WriteTimeout"] = redirectObj["write_timeout"]
		redirectMap["ReadTimeout"] = redirectObj["read_timeout"]

		requestHeaderMap := make([]map[string]interface{}, 0)
		if v, ok := redirectObj["request_headers"]; ok {
			for _, requestHeader := range v.(*schema.Set).List() {
				requestHeaderObj := requestHeader.(map[string]interface{})
				requestHeaderMap = append(requestHeaderMap, map[string]interface{}{
					"Key":   requestHeaderObj["key"],
					"Value": requestHeaderObj["value"],
				})
			}
			redirectMap["RequestHeaders"] = requestHeaderMap
		}
	}
	request["Redirect"], _ = convertMaptoJsonString(redirectMap)

	listenMap := make(map[string]interface{}, 0)
	for _, v := range d.Get("listen").([]interface{}) {
		listenObj := v.(map[string]interface{})
		listenMap["CertId"] = listenObj["cert_id"]
		listenMap["CipherSuite"] = listenObj["cipher_suite"]
		listenMap["CustomCiphers"] = listenObj["custom_ciphers"].([]interface{})
		listenMap["EnableTLSv3"] = listenObj["enable_tlsv3"]
		listenMap["ExclusiveIp"] = listenObj["exclusive_ip"]
		listenMap["FocusHttps"] = listenObj["focus_https"]
		listenMap["Http2Enabled"] = listenObj["http2_enabled"]
		listenMap["HttpPorts"] = listenObj["http_ports"].([]interface{})
		listenMap["HttpsPorts"] = listenObj["https_ports"].([]interface{})
		listenMap["IPv6Enabled"] = listenObj["ipv6_enabled"]
		listenMap["ProtectionResource"] = listenObj["protection_resource"]
		listenMap["TLSVersion"] = listenObj["tls_version"]
		listenMap["XffHeaderMode"] = listenObj["xff_header_mode"]
		listenMap["XffHeaders"] = listenObj["xff_headers"]
	}
	request["Listen"], _ = convertMaptoJsonString(listenMap)

	action := "ModifyDomain"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
		response, err = client.RpcPost("waf-openapi", "2021-10-01", action, nil, request, false)
		if err != nil {
			if NeedRetry(err) || IsExpectedErrors(err, []string{"Waf.Pullin.ResourceProcessing"}) {
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
	stateConf := BuildStateConf([]string{}, []string{"1"}, d.Timeout(schema.TimeoutDelete), 5*time.Second, wafOpenapiService.Wafv3DomainStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudWafv3DomainRead(d, meta)
}

func resourceAlicloudWafv3DomainDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	wafOpenapiService := WafOpenapiService{client}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	var response map[string]interface{}
	request := map[string]interface{}{
		"InstanceId": parts[0],
		"Domain":     parts[1],
		"RegionId":   client.RegionId,
	}

	action := "DeleteDomain"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = client.RpcPost("waf-openapi", "2021-10-01", action, nil, request, false)
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
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, wafOpenapiService.Wafv3DomainStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
