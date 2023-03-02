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

func dataSourceAlicloudWafv3Domains() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudWafv3DomainsRead,
		Schema: map[string]*schema.Schema{
			"backend": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"domain": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"instance_id": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"ids": {
				Optional: true,
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"enable_details": {
				Optional: true,
				Type:     schema.TypeBool,
			},
			"output_file": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"page_number": {
				Optional: true,
				Type:     schema.TypeInt,
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
							Computed: true,
							Type:     schema.TypeString,
						},
						"domain": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"status": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"listen": {
							Computed: true,
							Type:     schema.TypeList,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cert_id": {
										Computed: true,
										Type:     schema.TypeString,
									},
									"cipher_suite": {
										Computed: true,
										Type:     schema.TypeInt,
									},
									"custom_ciphers": {
										Computed: true,
										Type:     schema.TypeList,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"enable_tlsv3": {
										Computed: true,
										Type:     schema.TypeBool,
									},
									"exclusive_ip": {
										Computed: true,
										Type:     schema.TypeBool,
									},
									"focus_https": {
										Computed: true,
										Type:     schema.TypeBool,
									},
									"http2_enabled": {
										Computed: true,
										Type:     schema.TypeBool,
									},
									"http_ports": {
										Computed: true,
										Type:     schema.TypeList,
										Elem: &schema.Schema{
											Type: schema.TypeInt,
										},
									},
									"https_ports": {
										Computed: true,
										Type:     schema.TypeList,
										Elem: &schema.Schema{
											Type: schema.TypeInt,
										},
									},
									"ipv6_enabled": {
										Computed: true,
										Type:     schema.TypeBool,
									},
									"protection_resource": {
										Computed: true,
										Type:     schema.TypeString,
									},
									"tls_version": {
										Computed: true,
										Type:     schema.TypeString,
									},
									"xff_header_mode": {
										Computed: true,
										Type:     schema.TypeInt,
									},
									"xff_headers": {
										Computed: true,
										Type:     schema.TypeList,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
						"redirect": {
							Computed: true,
							Type:     schema.TypeList,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"backends": {
										Computed: true,
										Type:     schema.TypeList,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"connect_timeout": {
										Computed: true,
										Type:     schema.TypeInt,
									},
									"focus_http_backend": {
										Computed: true,
										Type:     schema.TypeBool,
									},
									"keepalive": {
										Computed: true,
										Type:     schema.TypeBool,
									},
									"keepalive_requests": {
										Computed: true,
										Type:     schema.TypeInt,
									},
									"keepalive_timeout": {
										Computed: true,
										Type:     schema.TypeInt,
									},
									"loadbalance": {
										Computed: true,
										Type:     schema.TypeString,
									},
									"read_timeout": {
										Computed: true,
										Type:     schema.TypeInt,
									},
									"request_headers": {
										Computed: true,
										Type:     schema.TypeList,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"key": {
													Computed: true,
													Type:     schema.TypeString,
												},
												"value": {
													Computed: true,
													Type:     schema.TypeString,
												},
											},
										},
									},
									"retry": {
										Computed: true,
										Type:     schema.TypeBool,
									},
									"sni_enabled": {
										Computed: true,
										Type:     schema.TypeBool,
									},
									"sni_host": {
										Computed: true,
										Type:     schema.TypeString,
									},
									"write_timeout": {
										Computed: true,
										Type:     schema.TypeInt,
									},
								},
							},
						},
						"resource_manager_resource_group_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudWafv3DomainsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := map[string]interface{}{
		"RegionId": client.RegionId,
	}

	if v, ok := d.GetOk("backend"); ok {
		request["Backend"] = v
	}
	if v, ok := d.GetOk("domain"); ok {
		request["Domain"] = v
	}
	request["InstanceId"] = d.Get("instance_id")
	setPagingRequest(d, request, PageSizeLarge)

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}

	conn, err := client.NewWafClient()
	if err != nil {
		return WrapError(err)
	}
	var objects []interface{}
	var response map[string]interface{}

	for {
		action := "DescribeDomains"
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			resp, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2021-10-01"), StringPointer("AK"), nil, request, &runtime)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			response = resp
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_wafv3_domains", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.Domains", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Domains", response)
		}
		result, _ := resp.([]interface{})
		if isPagingRequest(d) {
			objects = result
			break
		}
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
	wafOpenapiService := WafOpenapiService{client}
	for _, v := range objects {
		object := v.(map[string]interface{})

		mapping := map[string]interface{}{
			"id": fmt.Sprint(request["InstanceId"], ":", object["Domain"]),
		}

		resource85 := object["Domain"]
		mapping["domain"] = resource85
		mapping["status"] = fmt.Sprint(object["Status"])

		resourceManagerResourceGroupId10 := object["ResourceManagerResourceGroupId"]
		mapping["resource_manager_resource_group_id"] = resourceManagerResourceGroupId10

		ids = append(ids, fmt.Sprint(request["InstanceId"], ":", object["Domain"]))

		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			s = append(s, mapping)
			continue
		}
		id := fmt.Sprint(request["InstanceId"], ":", object["Domain"])
		object, err = wafOpenapiService.DescribeWafv3Domain(id)
		if err != nil {
			return WrapError(err)
		}
		mapping["domain"] = object["Domain"]
		listen87Maps := make([]map[string]interface{}, 0)
		listen87Map := make(map[string]interface{})
		listen87Raw := object["Listen"].(map[string]interface{})
		listen87Map["cert_id"] = listen87Raw["CertId"]
		listen87Map["cipher_suite"] = listen87Raw["CipherSuite"]
		listen87Map["custom_ciphers"] = listen87Raw["CustomCiphers"].([]interface{})
		listen87Map["enable_tlsv3"] = listen87Raw["EnableTLSv3"]
		listen87Map["exclusive_ip"] = listen87Raw["ExclusiveIp"]
		listen87Map["focus_https"] = listen87Raw["FocusHttps"]
		listen87Map["http2_enabled"] = listen87Raw["Http2Enabled"]
		listen87Map["http_ports"] = listen87Raw["HttpsPorts"].([]interface{})
		listen87Map["https_ports"] = listen87Raw["HttpPorts"].([]interface{})
		listen87Map["ipv6_enabled"] = listen87Raw["IPv6Enabled"]
		listen87Map["protection_resource"] = listen87Raw["ProtectionResource"]
		listen87Map["tls_version"] = listen87Raw["TLSVersion"]
		listen87Map["xff_header_mode"] = listen87Raw["XffHeaderMode"]
		listen87Map["xff_headers"] = listen87Raw["XffHeaders"].([]interface{})
		listen87Maps = append(listen87Maps, listen87Map)
		mapping["listen"] = listen87Maps
		redirect49Maps := make([]map[string]interface{}, 0)
		redirect49Map := make(map[string]interface{})
		redirect49Raw := object["Redirect"].(map[string]interface{})
		redirect49Map["backends"] = redirect49Raw["AllBackends"].([]interface{})
		redirect49Map["connect_timeout"] = redirect49Raw["ConnectTimeout"]
		redirect49Map["focus_http_backend"] = redirect49Raw["FocusHttpBackend"]
		redirect49Map["keepalive"] = redirect49Raw["Keepalive"]
		redirect49Map["keepalive_requests"] = redirect49Raw["KeepaliveRequests"]
		redirect49Map["keepalive_timeout"] = redirect49Raw["KeepaliveTimeout"]
		redirect49Map["loadbalance"] = redirect49Raw["Loadbalance"]
		redirect49Map["read_timeout"] = redirect49Raw["ReadTimeout"]
		requestHeaders49Maps := make([]map[string]interface{}, 0)
		requestHeaders49Raw := redirect49Raw["RequestHeaders"]
		for _, value1 := range requestHeaders49Raw.([]interface{}) {
			requestHeaders49 := value1.(map[string]interface{})
			requestHeaders49Map := make(map[string]interface{})
			requestHeaders49Map["key"] = requestHeaders49["Key"]
			requestHeaders49Map["value"] = requestHeaders49["Value"]
			requestHeaders49Maps = append(requestHeaders49Maps, requestHeaders49Map)
		}
		redirect49Map["request_headers"] = requestHeaders49Maps
		redirect49Map["retry"] = redirect49Raw["Retry"]
		redirect49Map["sni_enabled"] = redirect49Raw["SniEnabled"]
		redirect49Map["sni_host"] = redirect49Raw["SniHost"]
		redirect49Map["write_timeout"] = redirect49Raw["WriteTimeout"]
		redirect49Maps = append(redirect49Maps, redirect49Map)
		mapping["redirect"] = redirect49Maps
		mapping["resource_manager_resource_group_id"] = object["ResourceManagerResourceGroupId"]
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
