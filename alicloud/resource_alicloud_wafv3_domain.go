// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudWafv3Domain() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudWafv3DomainCreate,
		Read:   resourceAliCloudWafv3DomainRead,
		Update: resourceAliCloudWafv3DomainUpdate,
		Delete: resourceAliCloudWafv3DomainDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"access_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"domain": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"domain_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"listen": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protection_resource": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: StringInSlice([]string{"share", "gslb"}, false),
						},
						"https_ports": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeInt},
						},
						"custom_ciphers": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"tls_version": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: StringInSlice([]string{"tlsv1", "tlsv1.1", "tlsv1.2"}, false),
						},
						"http2_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"cert_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"cipher_suite": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: IntInSlice([]int{0, 1, 2, 99}),
						},
						"enable_tlsv3": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"ipv6_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"focus_https": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"sm2_access_only": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"xff_headers": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"xff_header_mode": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: IntInSlice([]int{0, 1, 2}),
						},
						"sm2_cert_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"exclusive_ip": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"sm2_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"http_ports": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeInt},
						},
					},
				},
			},
			"redirect": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"connect_timeout": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: IntBetween(0, 3600),
						},
						"keepalive": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"sni_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"backup_backends": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"read_timeout": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: IntBetween(0, 3600),
						},
						"keepalive_timeout": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: IntBetween(0, 60),
						},
						"sni_host": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"backends": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"focus_http_backend": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"write_timeout": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: IntBetween(0, 3600),
						},
						"retry": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"request_headers": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"value": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"key": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"keepalive_requests": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: IntBetween(60, 1000),
						},
						"xff_proto": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"loadbalance": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: StringInSlice([]string{"iphash", "roundRobin", "leastTime"}, false),
						},
					},
				},
			},
			"resource_manager_resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"tags": tagsSchema(),
		},
	}
}

func resourceAliCloudWafv3DomainCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateDomain"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("domain"); ok {
		request["Domain"] = v
	}
	if v, ok := d.GetOk("instance_id"); ok {
		request["InstanceId"] = v
	}
	request["RegionId"] = client.RegionId

	objectDataLocalMap := make(map[string]interface{})

	if v := d.Get("listen"); v != nil {
		cipherSuite1, _ := jsonpath.Get("$[0].cipher_suite", v)
		if cipherSuite1 != nil && cipherSuite1 != "" {
			objectDataLocalMap["CipherSuite"] = cipherSuite1
		}
		httpPorts1, _ := jsonpath.Get("$[0].http_ports", v)
		if httpPorts1 != nil && httpPorts1 != "" {
			objectDataLocalMap["HttpPorts"] = httpPorts1
		}
		sm2Enabled, _ := jsonpath.Get("$[0].sm2_enabled", v)
		if sm2Enabled != nil && sm2Enabled != "" {
			objectDataLocalMap["SM2Enabled"] = sm2Enabled
		}
		focusHttps1, _ := jsonpath.Get("$[0].focus_https", v)
		if focusHttps1 != nil && focusHttps1 != "" {
			objectDataLocalMap["FocusHttps"] = focusHttps1
		}
		certId1, _ := jsonpath.Get("$[0].cert_id", v)
		if certId1 != nil && certId1 != "" {
			objectDataLocalMap["CertId"] = certId1
		}
		tLSVersion1, _ := jsonpath.Get("$[0].tls_version", v)
		if tLSVersion1 != nil && tLSVersion1 != "" {
			objectDataLocalMap["TLSVersion"] = tLSVersion1
		}
		customCiphers1, _ := jsonpath.Get("$[0].custom_ciphers", v)
		if customCiphers1 != nil && customCiphers1 != "" {
			objectDataLocalMap["CustomCiphers"] = customCiphers1
		}
		protectionResource1, _ := jsonpath.Get("$[0].protection_resource", v)
		if protectionResource1 != nil && protectionResource1 != "" {
			objectDataLocalMap["ProtectionResource"] = protectionResource1
		}
		sm2CertId, _ := jsonpath.Get("$[0].sm2_cert_id", v)
		if sm2CertId != nil && sm2CertId != "" {
			objectDataLocalMap["SM2CertId"] = sm2CertId
		}
		enableTLSv31, _ := jsonpath.Get("$[0].enable_tlsv3", v)
		if enableTLSv31 != nil && enableTLSv31 != "" {
			objectDataLocalMap["EnableTLSv3"] = enableTLSv31
		}
		exclusiveIp1, _ := jsonpath.Get("$[0].exclusive_ip", v)
		if exclusiveIp1 != nil && exclusiveIp1 != "" {
			objectDataLocalMap["ExclusiveIp"] = exclusiveIp1
		}
		xffHeaderMode1, _ := jsonpath.Get("$[0].xff_header_mode", v)
		if xffHeaderMode1 != nil && xffHeaderMode1 != "" {
			objectDataLocalMap["XffHeaderMode"] = xffHeaderMode1
		}
		http2Enabled1, _ := jsonpath.Get("$[0].http2_enabled", v)
		if http2Enabled1 != nil && http2Enabled1 != "" {
			objectDataLocalMap["Http2Enabled"] = http2Enabled1
		}
		xffHeaders1, _ := jsonpath.Get("$[0].xff_headers", v)
		if xffHeaders1 != nil && xffHeaders1 != "" {
			objectDataLocalMap["XffHeaders"] = xffHeaders1
		}
		iPv6Enabled1, _ := jsonpath.Get("$[0].ipv6_enabled", v)
		if iPv6Enabled1 != nil && iPv6Enabled1 != "" {
			objectDataLocalMap["IPv6Enabled"] = iPv6Enabled1
		}
		sm2AccessOnly, _ := jsonpath.Get("$[0].sm2_access_only", v)
		if sm2AccessOnly != nil && sm2AccessOnly != "" {
			objectDataLocalMap["SM2AccessOnly"] = sm2AccessOnly
		}
		httpsPorts1, _ := jsonpath.Get("$[0].https_ports", v)
		if httpsPorts1 != nil && httpsPorts1 != "" {
			objectDataLocalMap["HttpsPorts"] = httpsPorts1
		}

		objectDataLocalMapJson, err := json.Marshal(objectDataLocalMap)
		if err != nil {
			return WrapError(err)
		}
		request["Listen"] = string(objectDataLocalMapJson)
	}

	objectDataLocalMap1 := make(map[string]interface{})

	if v := d.Get("redirect"); v != nil {
		if v, ok := d.GetOk("redirect"); ok {
			localData, err := jsonpath.Get("$[0].request_headers", v)
			if err != nil {
				localData = make([]interface{}, 0)
			}
			localMaps := make([]interface{}, 0)
			for _, dataLoop := range localData.([]interface{}) {
				dataLoopTmp := make(map[string]interface{})
				if dataLoop != nil {
					dataLoopTmp = dataLoop.(map[string]interface{})
				}
				dataLoopMap := make(map[string]interface{})
				dataLoopMap["Key"] = dataLoopTmp["key"]
				dataLoopMap["Value"] = dataLoopTmp["value"]
				localMaps = append(localMaps, dataLoopMap)
			}
			objectDataLocalMap1["RequestHeaders"] = localMaps
		}

		sniHost1, _ := jsonpath.Get("$[0].sni_host", v)
		if sniHost1 != nil && sniHost1 != "" {
			objectDataLocalMap1["SniHost"] = sniHost1
		}
		xffProto1, _ := jsonpath.Get("$[0].xff_proto", v)
		if xffProto1 != nil && xffProto1 != "" {
			objectDataLocalMap1["XffProto"] = xffProto1
		}
		retry1, _ := jsonpath.Get("$[0].retry", v)
		if retry1 != nil && retry1 != "" {
			objectDataLocalMap1["Retry"] = retry1
		}
		keepalive1, _ := jsonpath.Get("$[0].keepalive", v)
		if keepalive1 != nil && keepalive1 != "" {
			objectDataLocalMap1["Keepalive"] = keepalive1
		}
		focusHttpBackend1, _ := jsonpath.Get("$[0].focus_http_backend", v)
		if focusHttpBackend1 != nil && focusHttpBackend1 != "" {
			objectDataLocalMap1["FocusHttpBackend"] = focusHttpBackend1
		}
		loadbalance1, _ := jsonpath.Get("$[0].loadbalance", v)
		if loadbalance1 != nil && loadbalance1 != "" {
			objectDataLocalMap1["Loadbalance"] = loadbalance1
		}
		sniEnabled1, _ := jsonpath.Get("$[0].sni_enabled", v)
		if sniEnabled1 != nil && sniEnabled1 != "" {
			objectDataLocalMap1["SniEnabled"] = sniEnabled1
		}
		keepaliveRequests1, _ := jsonpath.Get("$[0].keepalive_requests", v)
		if keepaliveRequests1 != nil && keepaliveRequests1 != "" && keepaliveRequests1.(int) > 0 {
			objectDataLocalMap1["KeepaliveRequests"] = keepaliveRequests1
		}
		connectTimeout1, _ := jsonpath.Get("$[0].connect_timeout", v)
		if connectTimeout1 != nil && connectTimeout1 != "" && connectTimeout1.(int) > 0 {
			objectDataLocalMap1["ConnectTimeout"] = connectTimeout1
		}
		writeTimeout1, _ := jsonpath.Get("$[0].write_timeout", v)
		if writeTimeout1 != nil && writeTimeout1 != "" && writeTimeout1.(int) > 0 {
			objectDataLocalMap1["WriteTimeout"] = writeTimeout1
		}
		backends1, _ := jsonpath.Get("$[0].backends", v)
		if backends1 != nil && backends1 != "" {
			objectDataLocalMap1["Backends"] = backends1
		}
		keepaliveTimeout1, _ := jsonpath.Get("$[0].keepalive_timeout", v)
		if keepaliveTimeout1 != nil && keepaliveTimeout1 != "" && keepaliveTimeout1.(int) > 0 {
			objectDataLocalMap1["KeepaliveTimeout"] = keepaliveTimeout1
		}
		backupBackends1, _ := jsonpath.Get("$[0].backup_backends", v)
		if backupBackends1 != nil && backupBackends1 != "" {
			objectDataLocalMap1["BackupBackends"] = backupBackends1
		}
		readTimeout1, _ := jsonpath.Get("$[0].read_timeout", v)
		if readTimeout1 != nil && readTimeout1 != "" && readTimeout1.(int) > 0 {
			objectDataLocalMap1["ReadTimeout"] = readTimeout1
		}

		objectDataLocalMap1Json, err := json.Marshal(objectDataLocalMap1)
		if err != nil {
			return WrapError(err)
		}
		request["Redirect"] = string(objectDataLocalMap1Json)
	}

	if v, ok := d.GetOk("tags"); ok {
		tagsMap := ConvertTags(v.(map[string]interface{}))
		request = expandTagsToMap(request, tagsMap)
	}

	if v, ok := d.GetOk("resource_manager_resource_group_id"); ok {
		request["ResourceManagerResourceGroupId"] = v
	}
	if v, ok := d.GetOk("access_type"); ok {
		request["AccessType"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("waf-openapi", "2021-10-01", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_wafv3_domain", action, AlibabaCloudSdkGoERROR)
	}

	DomainInfoDomainVar, _ := jsonpath.Get("$.DomainInfo.Domain", response)
	d.SetId(fmt.Sprintf("%v:%v", request["InstanceId"], DomainInfoDomainVar))

	wafv3ServiceV2 := Wafv3ServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"1"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, wafv3ServiceV2.Wafv3DomainStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudWafv3DomainUpdate(d, meta)
}

func resourceAliCloudWafv3DomainRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	wafv3ServiceV2 := Wafv3ServiceV2{client}

	objectRaw, err := wafv3ServiceV2.DescribeWafv3Domain(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_wafv3_domain DescribeWafv3Domain Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("domain_id", objectRaw["DomainId"])
	d.Set("resource_manager_resource_group_id", objectRaw["ResourceManagerResourceGroupId"])
	d.Set("status", objectRaw["Status"])
	d.Set("domain", objectRaw["Domain"])

	listenMaps := make([]map[string]interface{}, 0)
	listenMap := make(map[string]interface{})
	listenRaw := make(map[string]interface{})
	if objectRaw["Listen"] != nil {
		listenRaw = objectRaw["Listen"].(map[string]interface{})
	}
	if len(listenRaw) > 0 {
		listenMap["cert_id"] = listenRaw["CertId"]
		listenMap["cipher_suite"] = listenRaw["CipherSuite"]
		listenMap["enable_tlsv3"] = listenRaw["EnableTLSv3"]
		listenMap["exclusive_ip"] = listenRaw["ExclusiveIp"]
		listenMap["focus_https"] = listenRaw["FocusHttps"]
		listenMap["http2_enabled"] = listenRaw["Http2Enabled"]
		listenMap["ipv6_enabled"] = listenRaw["IPv6Enabled"]
		listenMap["protection_resource"] = listenRaw["ProtectionResource"]
		listenMap["sm2_access_only"] = listenRaw["SM2AccessOnly"]
		listenMap["sm2_cert_id"] = listenRaw["SM2CertId"]
		listenMap["sm2_enabled"] = listenRaw["SM2Enabled"]
		listenMap["tls_version"] = listenRaw["TLSVersion"]
		listenMap["xff_header_mode"] = listenRaw["XffHeaderMode"]

		customCiphersRaw := make([]interface{}, 0)
		if listenRaw["CustomCiphers"] != nil {
			customCiphersRaw = listenRaw["CustomCiphers"].([]interface{})
		}

		listenMap["custom_ciphers"] = customCiphersRaw
		httpPortsRaw := make([]interface{}, 0)
		if listenRaw["HttpPorts"] != nil {
			httpPortsRaw = listenRaw["HttpPorts"].([]interface{})
		}

		listenMap["http_ports"] = httpPortsRaw
		httpsPortsRaw := make([]interface{}, 0)
		if listenRaw["HttpsPorts"] != nil {
			httpsPortsRaw = listenRaw["HttpsPorts"].([]interface{})
		}

		listenMap["https_ports"] = httpsPortsRaw
		xffHeadersRaw := make([]interface{}, 0)
		if listenRaw["XffHeaders"] != nil {
			xffHeadersRaw = listenRaw["XffHeaders"].([]interface{})
		}

		listenMap["xff_headers"] = xffHeadersRaw
		listenMaps = append(listenMaps, listenMap)
	}
	if err := d.Set("listen", listenMaps); err != nil {
		return err
	}
	redirectMaps := make([]map[string]interface{}, 0)
	redirectMap := make(map[string]interface{})
	redirectRaw := make(map[string]interface{})
	if objectRaw["Redirect"] != nil {
		redirectRaw = objectRaw["Redirect"].(map[string]interface{})
	}
	if len(redirectRaw) > 0 {
		redirectMap["connect_timeout"] = redirectRaw["ConnectTimeout"]
		redirectMap["focus_http_backend"] = redirectRaw["FocusHttpBackend"]
		redirectMap["keepalive"] = redirectRaw["Keepalive"]
		redirectMap["keepalive_requests"] = redirectRaw["KeepaliveRequests"]
		redirectMap["keepalive_timeout"] = redirectRaw["KeepaliveTimeout"]
		redirectMap["loadbalance"] = redirectRaw["Loadbalance"]
		redirectMap["read_timeout"] = redirectRaw["ReadTimeout"]
		redirectMap["retry"] = redirectRaw["Retry"]
		redirectMap["sni_enabled"] = redirectRaw["SniEnabled"]
		redirectMap["sni_host"] = redirectRaw["SniHost"]
		redirectMap["write_timeout"] = redirectRaw["WriteTimeout"]
		redirectMap["xff_proto"] = redirectRaw["XffProto"]

		backendsRaw := redirectRaw["AllBackends"]
		if backendsRaw != nil {
			redirectMap["backends"] = backendsRaw
		}

		requestHeadersRaw := redirectRaw["RequestHeaders"]
		requestHeadersMaps := make([]map[string]interface{}, 0)
		if requestHeadersRaw != nil {
			for _, requestHeadersChildRaw := range requestHeadersRaw.([]interface{}) {
				requestHeadersMap := make(map[string]interface{})
				requestHeadersChildRaw := requestHeadersChildRaw.(map[string]interface{})
				requestHeadersMap["key"] = requestHeadersChildRaw["Key"]
				requestHeadersMap["value"] = requestHeadersChildRaw["Value"]

				requestHeadersMaps = append(requestHeadersMaps, requestHeadersMap)
			}
		}
		redirectMap["request_headers"] = requestHeadersMaps
		redirectMaps = append(redirectMaps, redirectMap)
	}
	if err := d.Set("redirect", redirectMaps); err != nil {
		return err
	}

	objectRaw, err = wafv3ServiceV2.DescribeDomainListTagResources(d.Id())
	if err != nil && !NotFoundError(err) {
		return WrapError(err)
	}

	tagsMaps := objectRaw["TagResources"]
	d.Set("tags", tagsToMap(tagsMaps))

	parts := strings.Split(d.Id(), ":")
	d.Set("instance_id", parts[0])

	return nil
}

func resourceAliCloudWafv3DomainUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	parts := strings.Split(d.Id(), ":")
	action := "ModifyDomain"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["Domain"] = parts[1]
	request["InstanceId"] = parts[0]
	request["RegionId"] = client.RegionId
	if !d.IsNewResource() && d.HasChange("listen") {
		update = true
	}
	objectDataLocalMap := make(map[string]interface{})

	if v := d.Get("listen"); v != nil {
		enableTLSv31, _ := jsonpath.Get("$[0].enable_tlsv3", v)
		if enableTLSv31 != nil && (d.HasChange("listen.0.enable_tlsv3") || enableTLSv31 != "") {
			objectDataLocalMap["EnableTLSv3"] = enableTLSv31
		}
		cipherSuite1, _ := jsonpath.Get("$[0].cipher_suite", v)
		if cipherSuite1 != nil && (d.HasChange("listen.0.cipher_suite") || cipherSuite1 != "") {
			objectDataLocalMap["CipherSuite"] = cipherSuite1
		}
		exclusiveIp1, _ := jsonpath.Get("$[0].exclusive_ip", v)
		if exclusiveIp1 != nil && (d.HasChange("listen.0.exclusive_ip") || exclusiveIp1 != "") {
			objectDataLocalMap["ExclusiveIp"] = exclusiveIp1
		}
		xffHeaderMode1, _ := jsonpath.Get("$[0].xff_header_mode", v)
		if xffHeaderMode1 != nil && (d.HasChange("listen.0.xff_header_mode") || xffHeaderMode1 != "") {
			objectDataLocalMap["XffHeaderMode"] = xffHeaderMode1
		}
		http2Enabled1, _ := jsonpath.Get("$[0].http2_enabled", v)
		if http2Enabled1 != nil && (d.HasChange("listen.0.http2_enabled") || http2Enabled1 != "") {
			objectDataLocalMap["Http2Enabled"] = http2Enabled1
		}
		httpPorts1, _ := jsonpath.Get("$[0].http_ports", d.Get("listen"))
		if httpPorts1 != nil && (d.HasChange("listen.0.http_ports") || httpPorts1 != "") {
			objectDataLocalMap["HttpPorts"] = httpPorts1
		}
		xffHeaders1, _ := jsonpath.Get("$[0].xff_headers", d.Get("listen"))
		if xffHeaders1 != nil && (d.HasChange("listen.0.xff_headers") || xffHeaders1 != "") {
			objectDataLocalMap["XffHeaders"] = xffHeaders1
		}
		sm2Enabled, _ := jsonpath.Get("$[0].sm2_enabled", v)
		if sm2Enabled != nil && (d.HasChange("listen.0.sm2_enabled") || sm2Enabled != "") {
			objectDataLocalMap["SM2Enabled"] = sm2Enabled
		}
		iPv6Enabled1, _ := jsonpath.Get("$[0].ipv6_enabled", v)
		if iPv6Enabled1 != nil && (d.HasChange("listen.0.ipv6_enabled") || iPv6Enabled1 != "") {
			objectDataLocalMap["IPv6Enabled"] = iPv6Enabled1
		}
		sm2AccessOnly, _ := jsonpath.Get("$[0].sm2_access_only", v)
		if sm2AccessOnly != nil && (d.HasChange("listen.0.sm2_access_only") || sm2AccessOnly != "") {
			objectDataLocalMap["SM2AccessOnly"] = sm2AccessOnly
		}
		focusHttps1, _ := jsonpath.Get("$[0].focus_https", v)
		if focusHttps1 != nil && (d.HasChange("listen.0.focus_https") || focusHttps1 != "") {
			objectDataLocalMap["FocusHttps"] = focusHttps1
		}
		certId1, _ := jsonpath.Get("$[0].cert_id", v)
		if certId1 != nil && (d.HasChange("listen.0.cert_id") || certId1 != "") {
			objectDataLocalMap["CertId"] = certId1
		}
		tLSVersion1, _ := jsonpath.Get("$[0].tls_version", v)
		if tLSVersion1 != nil && (d.HasChange("listen.0.tls_version") || tLSVersion1 != "") {
			objectDataLocalMap["TLSVersion"] = tLSVersion1
		}
		customCiphers1, _ := jsonpath.Get("$[0].custom_ciphers", d.Get("listen"))
		if customCiphers1 != nil && (d.HasChange("listen.0.custom_ciphers") || customCiphers1 != "") {
			objectDataLocalMap["CustomCiphers"] = customCiphers1
		}
		protectionResource1, _ := jsonpath.Get("$[0].protection_resource", v)
		if protectionResource1 != nil && (d.HasChange("listen.0.protection_resource") || protectionResource1 != "") {
			objectDataLocalMap["ProtectionResource"] = protectionResource1
		}
		httpsPorts1, _ := jsonpath.Get("$[0].https_ports", d.Get("listen"))
		if httpsPorts1 != nil && (d.HasChange("listen.0.https_ports") || httpsPorts1 != "") {
			objectDataLocalMap["HttpsPorts"] = httpsPorts1
		}
		sm2CertId, _ := jsonpath.Get("$[0].sm2_cert_id", v)
		if sm2CertId != nil && (d.HasChange("listen.0.sm2_cert_id") || sm2CertId != "") {
			objectDataLocalMap["SM2CertId"] = sm2CertId
		}

		objectDataLocalMapJson, err := json.Marshal(objectDataLocalMap)
		if err != nil {
			return WrapError(err)
		}
		request["Listen"] = string(objectDataLocalMapJson)
	}

	if !d.IsNewResource() && d.HasChange("redirect") {
		update = true
	}
	objectDataLocalMap1 := make(map[string]interface{})

	if v := d.Get("redirect"); v != nil {
		sniEnabled1, _ := jsonpath.Get("$[0].sni_enabled", v)
		if sniEnabled1 != nil && (d.HasChange("redirect.0.sni_enabled") || sniEnabled1 != "") {
			objectDataLocalMap1["SniEnabled"] = sniEnabled1
		}
		if v, ok := d.GetOk("redirect"); ok {
			localData, err := jsonpath.Get("$[0].request_headers", v)
			if err != nil {
				localData = make([]interface{}, 0)
			}
			localMaps := make([]interface{}, 0)
			for _, dataLoop := range localData.([]interface{}) {
				dataLoopTmp := make(map[string]interface{})
				if dataLoop != nil {
					dataLoopTmp = dataLoop.(map[string]interface{})
				}
				dataLoopMap := make(map[string]interface{})
				dataLoopMap["Key"] = dataLoopTmp["key"]
				dataLoopMap["Value"] = dataLoopTmp["value"]
				localMaps = append(localMaps, dataLoopMap)
			}
			objectDataLocalMap1["RequestHeaders"] = localMaps
		}

		sniHost1, _ := jsonpath.Get("$[0].sni_host", v)
		if sniHost1 != nil && (d.HasChange("redirect.0.sni_host") || sniHost1 != "") {
			objectDataLocalMap1["SniHost"] = sniHost1
		}
		xffProto1, _ := jsonpath.Get("$[0].xff_proto", v)
		if xffProto1 != nil && (d.HasChange("redirect.0.xff_proto") || xffProto1 != "") {
			objectDataLocalMap1["XffProto"] = xffProto1
		}
		keepaliveRequests1, _ := jsonpath.Get("$[0].keepalive_requests", v)
		if keepaliveRequests1 != nil && (d.HasChange("redirect.0.keepalive_requests") || keepaliveRequests1 != "") && keepaliveRequests1.(int) > 0 {
			objectDataLocalMap1["KeepaliveRequests"] = keepaliveRequests1
		}
		retry1, _ := jsonpath.Get("$[0].retry", v)
		if retry1 != nil && (d.HasChange("redirect.0.retry") || retry1 != "") {
			objectDataLocalMap1["Retry"] = retry1
		}
		connectTimeout1, _ := jsonpath.Get("$[0].connect_timeout", v)
		if connectTimeout1 != nil && (d.HasChange("redirect.0.connect_timeout") || connectTimeout1 != "") && connectTimeout1.(int) > 0 {
			objectDataLocalMap1["ConnectTimeout"] = connectTimeout1
		}
		keepalive1, _ := jsonpath.Get("$[0].keepalive", v)
		if keepalive1 != nil && (d.HasChange("redirect.0.keepalive") || keepalive1 != "") {
			objectDataLocalMap1["Keepalive"] = keepalive1
		}
		writeTimeout1, _ := jsonpath.Get("$[0].write_timeout", v)
		if writeTimeout1 != nil && (d.HasChange("redirect.0.write_timeout") || writeTimeout1 != "") && writeTimeout1.(int) > 0 {
			objectDataLocalMap1["WriteTimeout"] = writeTimeout1
		}
		focusHttpBackend1, _ := jsonpath.Get("$[0].focus_http_backend", v)
		if focusHttpBackend1 != nil && (d.HasChange("redirect.0.focus_http_backend") || focusHttpBackend1 != "") {
			objectDataLocalMap1["FocusHttpBackend"] = focusHttpBackend1
		}
		loadbalance1, _ := jsonpath.Get("$[0].loadbalance", v)
		if loadbalance1 != nil && (d.HasChange("redirect.0.loadbalance") || loadbalance1 != "") {
			objectDataLocalMap1["Loadbalance"] = loadbalance1
		}
		backends1, _ := jsonpath.Get("$[0].backends", d.Get("redirect"))
		if backends1 != nil && (d.HasChange("redirect.0.backends") || backends1 != "") {
			objectDataLocalMap1["Backends"] = backends1
		}
		keepaliveTimeout1, _ := jsonpath.Get("$[0].keepalive_timeout", v)
		if keepaliveTimeout1 != nil && (d.HasChange("redirect.0.keepalive_timeout") || keepaliveTimeout1 != "") && keepaliveTimeout1.(int) > 0 {
			objectDataLocalMap1["KeepaliveTimeout"] = keepaliveTimeout1
		}
		backupBackends1, _ := jsonpath.Get("$[0].backup_backends", d.Get("redirect"))
		if backupBackends1 != nil && (d.HasChange("redirect.0.backup_backends") || backupBackends1 != "") {
			objectDataLocalMap1["BackupBackends"] = backupBackends1
		}
		readTimeout1, _ := jsonpath.Get("$[0].read_timeout", v)
		if readTimeout1 != nil && (d.HasChange("redirect.0.read_timeout") || readTimeout1 != "") && readTimeout1.(int) > 0 {
			objectDataLocalMap1["ReadTimeout"] = readTimeout1
		}

		objectDataLocalMap1Json, err := json.Marshal(objectDataLocalMap1)
		if err != nil {
			return WrapError(err)
		}
		request["Redirect"] = string(objectDataLocalMap1Json)
	}

	if d.HasChange("domain_id") {
		update = true
		request["DomainId"] = d.Get("domain_id")
	}

	if v, ok := d.GetOk("access_type"); ok {
		request["AccessType"] = v
	}
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("waf-openapi", "2021-10-01", action, query, request, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"Waf.Pullin.ResourceProcessing"}) || NeedRetry(err) {
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
		wafv3ServiceV2 := Wafv3ServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"1"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, wafv3ServiceV2.Wafv3DomainStateRefreshFunc(d.Id(), "Status", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	if d.HasChange("tags") {
		wafv3ServiceV2 := Wafv3ServiceV2{client}
		if err := wafv3ServiceV2.SetResourceTags(d, "ALIYUN::WAF::DEFENSERESOURCE"); err != nil {
			return WrapError(err)
		}
	}
	return resourceAliCloudWafv3DomainRead(d, meta)
}

func resourceAliCloudWafv3DomainDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeleteDomain"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["InstanceId"] = parts[0]
	request["Domain"] = parts[1]
	request["RegionId"] = client.RegionId

	if v, ok := d.GetOk("domain_id"); ok {
		request["DomainId"] = v
	}
	if v, ok := d.GetOk("access_type"); ok {
		request["AccessType"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("waf-openapi", "2021-10-01", action, query, request, true)

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

	wafv3ServiceV2 := Wafv3ServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, wafv3ServiceV2.Wafv3DomainStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}
