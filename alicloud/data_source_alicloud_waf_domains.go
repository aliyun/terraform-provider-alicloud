package alicloud

import (
	"regexp"

	waf_openapi "github.com/aliyun/alibaba-cloud-sdk-go/services/waf-openapi"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudWafDomains() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudWafDomainsRead,
		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
				ForceNew:     true,
			},
			"names": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"domains": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cluster_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cname": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"connection_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"domain_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"domain": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"http2_port": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"http_port": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"http_to_user_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"https_port": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"https_redirect": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_access_product": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"load_balancing": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"log_headers": {
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
						"read_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"resource_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_ips": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"version": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"write_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"enable_details": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func dataSourceAlicloudWafDomainsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := waf_openapi.CreateDescribeDomainNamesRequest()
	if v, ok := d.GetOk("instance_id"); ok {
		request.InstanceId = v.(string)
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request.ResourceGroupId = v.(string)
	}
	var domainNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		domainNameRegex = r
	}

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}
	var response *waf_openapi.DescribeDomainNamesResponse
	raw, err := client.WithWafOpenapiClient(func(waf_openapiClient *waf_openapi.Client) (interface{}, error) {
		return waf_openapiClient.DescribeDomainNames(request)
	})
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_waf_domains", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)
	response, _ = raw.(*waf_openapi.DescribeDomainNamesResponse)

	ids := make([]string, 0)
	names := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range response.DomainNames {
		if domainNameRegex != nil {
			if !domainNameRegex.MatchString(object) {
				continue
			}
		}
		if len(idsMap) > 0 {
			if _, ok := idsMap[object]; !ok {
				continue
			}
		}
		mapping := map[string]interface{}{
			"id":          object,
			"domain_name": object,
			"domain":      object,
		}
		ids = append(ids, object)
		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			names = append(names, object)
			s = append(s, mapping)
			continue
		}
		request := waf_openapi.CreateDescribeDomainRequest()
		request.RegionId = client.RegionId
		request.Domain = object
		request.InstanceId = d.Get("instance_id").(string)
		raw, err := client.WithWafOpenapiClient(func(waf_openapiClient *waf_openapi.Client) (interface{}, error) {
			return waf_openapiClient.DescribeDomain(request)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_waf_domains", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		responseGet, _ := raw.(*waf_openapi.DescribeDomainResponse)
		mapping["cluster_type"] = convertClusterTypeResponse(responseGet.Domain.ClusterType)
		mapping["cname"] = responseGet.Domain.Cname
		mapping["connection_time"] = responseGet.Domain.ConnectionTime
		mapping["http2_port"] = responseGet.Domain.Http2Port
		mapping["http_port"] = responseGet.Domain.HttpPort
		mapping["http_to_user_ip"] = convertHttpToUserIpResponse(responseGet.Domain.HttpToUserIp)
		mapping["https_port"] = responseGet.Domain.HttpsPort
		mapping["https_redirect"] = convertHttpsRedirectResponse(responseGet.Domain.HttpsRedirect)
		mapping["is_access_product"] = convertIsAccessProductResponse(responseGet.Domain.IsAccessProduct)
		mapping["load_balancing"] = convertLoadBalancingResponse(responseGet.Domain.LoadBalancing)
		logHeaders := make([]map[string]interface{}, len(responseGet.Domain.LogHeaders))
		for i, v := range responseGet.Domain.LogHeaders {
			logHeaders[i] = map[string]interface{}{
				"key":   v.K,
				"value": v.V,
			}
		}
		mapping["log_headers"] = logHeaders
		mapping["read_time"] = responseGet.Domain.ReadTime
		mapping["resource_group_id"] = responseGet.Domain.ResourceGroupId
		mapping["source_ips"] = responseGet.Domain.SourceIps
		mapping["version"] = responseGet.Domain.Version
		mapping["write_time"] = responseGet.Domain.WriteTime
		names = append(names, object)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
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
