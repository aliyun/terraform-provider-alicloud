package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAliCloudWafv3Domains() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudWafv3DomainsRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"domain": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"backend": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"enable_details": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"page_number": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"page_size": {
				Optional: true,
				Type:     schema.TypeInt,
				Default:  10,
			},
			"domains": {
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"domain": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_manager_resource_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"listen": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cert_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"cipher_suite": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"enable_tlsv3": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"exclusive_ip": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"focus_https": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"http2_enabled": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"ipv6_enabled": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"protection_resource": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"tls_version": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"xff_header_mode": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"http_ports": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeInt},
									},
									"https_ports": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeInt},
									},
									"custom_ciphers": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"xff_headers": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"redirect": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"loadbalance": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"focus_http_backend": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"keepalive": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"keepalive_requests": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"retry": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"sni_enabled": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"sni_host": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"connect_timeout": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"keepalive_timeout": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"read_timeout": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"write_timeout": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"backends": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"request_headers": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"key": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"value": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceAliCloudWafv3DomainsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeDomains"
	request := make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["InstanceId"] = d.Get("instance_id")
	setPagingRequest(d, request, PageSizeLarge)

	if v, ok := d.GetOk("domain"); ok {
		request["Domain"] = v
	}

	if v, ok := d.GetOk("backend"); ok {
		request["Backend"] = v
	}

	var objects []map[string]interface{}
	wafOpenapiService := WafOpenapiService{client}
	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}

	var response map[string]interface{}
	var err error

	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("waf-openapi", "2021-10-01", action, nil, request, true)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_wafv3_domains", action, AlibabaCloudSdkGoERROR)
		}

		resp, err := jsonpath.Get("$.Domains", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Domains", response)
		}

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(request["InstanceId"], ":", item["Domain"])]; !ok {
					continue
				}
			}

			objects = append(objects, item)
		}

		if len(result) < request["PageSize"].(int) {
			break
		}

		request["PageNumber"] = request["PageNumber"].(int) + 1
	}

	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"id":                                 fmt.Sprintf("%v:%v", request["InstanceId"], object["Domain"]),
			"domain":                             object["Domain"],
			"resource_manager_resource_group_id": object["ResourceManagerResourceGroupId"],
			"status":                             fmt.Sprint(object["Status"]),
		}

		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			s = append(s, mapping)
			continue
		}

		object, err = wafOpenapiService.DescribeWafv3Domain(fmt.Sprint(mapping["id"]))
		if err != nil {
			return WrapError(err)
		}

		if listen, ok := object["Listen"]; ok {
			listenMaps := make([]map[string]interface{}, 0)
			listenArg := listen.(map[string]interface{})
			listenMap := map[string]interface{}{}

			if certId, ok := listenArg["CertId"]; ok {
				listenMap["cert_id"] = certId
			}

			if cipherSuite, ok := listenArg["CipherSuite"]; ok {
				listenMap["cipher_suite"] = cipherSuite
			}

			if enableTLSv3, ok := listenArg["EnableTLSv3"]; ok {
				listenMap["enable_tlsv3"] = enableTLSv3
			}

			if exclusiveIp, ok := listenArg["ExclusiveIp"]; ok {
				listenMap["exclusive_ip"] = exclusiveIp
			}

			if focusHttps, ok := listenArg["FocusHttps"]; ok {
				listenMap["focus_https"] = focusHttps
			}

			if http2Enabled, ok := listenArg["Http2Enabled"]; ok {
				listenMap["http2_enabled"] = http2Enabled
			}

			if iPv6Enabled, ok := listenArg["IPv6Enabled"]; ok {
				listenMap["ipv6_enabled"] = iPv6Enabled
			}

			if protectionResource, ok := listenArg["ProtectionResource"]; ok {
				listenMap["protection_resource"] = protectionResource
			}

			if tlsVersion, ok := listenArg["TLSVersion"]; ok {
				listenMap["tls_version"] = tlsVersion
			}

			if xffHeaderMode, ok := listenArg["XffHeaderMode"]; ok {
				listenMap["xff_header_mode"] = xffHeaderMode
			}

			if httpPorts, ok := listenArg["HttpPorts"]; ok {
				listenMap["http_ports"] = httpPorts
			}

			if httpsPorts, ok := listenArg["HttpsPorts"]; ok {
				listenMap["https_ports"] = httpsPorts
			}

			if customCiphers, ok := listenArg["CustomCiphers"]; ok {
				listenMap["custom_ciphers"] = customCiphers
			}

			if xffHeaders, ok := listenArg["XffHeaders"]; ok {
				listenMap["xff_headers"] = xffHeaders
			}

			listenMaps = append(listenMaps, listenMap)
			mapping["listen"] = listenMaps
		}

		if redirect, ok := object["Redirect"]; ok {
			redirectMaps := make([]map[string]interface{}, 0)
			redirectArg := redirect.(map[string]interface{})
			redirectMap := map[string]interface{}{}

			if loadBalance, ok := redirectArg["Loadbalance"]; ok {
				redirectMap["loadbalance"] = loadBalance
			}

			if focusHttpBackend, ok := redirectArg["FocusHttpBackend"]; ok {
				redirectMap["focus_http_backend"] = focusHttpBackend
			}

			if keepalive, ok := redirectArg["Keepalive"]; ok {
				redirectMap["keepalive"] = keepalive
			}

			if keepaliveRequests, ok := redirectArg["KeepaliveRequests"]; ok {
				redirectMap["keepalive_requests"] = keepaliveRequests
			}

			if retry, ok := redirectArg["Retry"]; ok {
				redirectMap["retry"] = retry
			}

			if sniEnabled, ok := redirectArg["SniEnabled"]; ok {
				redirectMap["sni_enabled"] = sniEnabled
			}

			if sniHost, ok := redirectArg["SniHost"]; ok {
				redirectMap["sni_host"] = sniHost
			}

			if connectTimeout, ok := redirectArg["ConnectTimeout"]; ok {
				redirectMap["connect_timeout"] = connectTimeout
			}

			if keepaliveTimeout, ok := redirectArg["KeepaliveTimeout"]; ok {
				redirectMap["keepalive_timeout"] = keepaliveTimeout
			}

			if readTimeout, ok := redirectArg["ReadTimeout"]; ok {
				redirectMap["read_timeout"] = readTimeout
			}

			if writeTimeout, ok := redirectArg["WriteTimeout"]; ok {
				redirectMap["write_timeout"] = writeTimeout
			}

			if allBackends, ok := redirectArg["AllBackends"]; ok {
				redirectMap["backends"] = allBackends
			}

			if requestHeadersList, ok := redirectArg["RequestHeaders"]; ok {
				requestHeadersMaps := make([]map[string]interface{}, 0)
				for _, requestHeaders := range requestHeadersList.([]interface{}) {
					requestHeadersArg := requestHeaders.(map[string]interface{})
					requestHeadersMap := map[string]interface{}{}

					if key, ok := requestHeadersArg["Key"]; ok {
						requestHeadersMap["key"] = key
					}

					if value, ok := requestHeadersArg["Value"]; ok {
						requestHeadersMap["value"] = value
					}

					requestHeadersMaps = append(requestHeadersMaps, requestHeadersMap)
				}

				redirectMap["request_headers"] = requestHeadersMaps
			}

			redirectMaps = append(redirectMaps, redirectMap)
			mapping["redirect"] = redirectMaps
		}

		ids = append(ids, fmt.Sprint(mapping["id"]))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("domains", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
