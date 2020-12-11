package alicloud

import (
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dcdn"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudDcdnDomains() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudDcdnDomainsRead,
		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
				ForceNew:     true,
			},
			"change_end_time": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"change_start_time": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"check_domain_show": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
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
			"domain_search_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"security_token": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"check_failed", "checking", "configure_failed", "configuring", "offline", "online"}, false),
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
						"cert_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cname": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
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
						"gmt_modified": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ssl_protocol": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ssl_pub": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"scope": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sources": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"content": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"enabled": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"port": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"priority": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"weight": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"status": {
							Type:     schema.TypeString,
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

func dataSourceAlicloudDcdnDomainsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := dcdn.CreateDescribeDcdnUserDomainsRequest()
	if v, ok := d.GetOk("change_end_time"); ok {
		request.ChangeEndTime = v.(string)
	}
	if v, ok := d.GetOk("change_start_time"); ok {
		request.ChangeStartTime = v.(string)
	}
	if v, ok := d.GetOkExists("check_domain_show"); ok {
		request.CheckDomainShow = requests.NewBoolean(v.(bool))
	}
	if v, ok := d.GetOk("domain_search_type"); ok {
		request.DomainSearchType = v.(string)
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request.ResourceGroupId = v.(string)
	}
	if v, ok := d.GetOk("security_token"); ok {
		request.SecurityToken = v.(string)
	}
	if v, ok := d.GetOk("status"); ok {
		request.DomainStatus = v.(string)
	}
	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)
	var objects []dcdn.PageData
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
	var response *dcdn.DescribeDcdnUserDomainsResponse
	for {
		raw, err := client.WithDcdnClient(func(dcdnClient *dcdn.Client) (interface{}, error) {
			return dcdnClient.DescribeDcdnUserDomains(request)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_dcdn_domains", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw)
		response, _ = raw.(*dcdn.DescribeDcdnUserDomainsResponse)

		for _, item := range response.Domains.PageData {
			if domainNameRegex != nil {
				if !domainNameRegex.MatchString(item.DomainName) {
					continue
				}
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[item.DomainName]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}
		if len(response.Domains.PageData) < PageSizeLarge {
			break
		}

		page, err := getNextpageNumber(request.PageNumber)
		if err != nil {
			return WrapError(err)
		}
		request.PageNumber = page
	}
	ids := make([]string, len(objects))
	names := make([]string, len(objects))
	s := make([]map[string]interface{}, len(objects))
	for i, object := range objects {
		mapping := map[string]interface{}{
			"cname":             object.Cname,
			"id":                object.DomainName,
			"domain_name":       object.DomainName,
			"gmt_modified":      object.GmtModified,
			"resource_group_id": object.ResourceGroupId,
			"ssl_protocol":      object.SSLProtocol,
			"status":            object.DomainStatus,
		}

		var sourcesList []map[string]interface{}
		for _, v := range object.Sources.Source {
			sourcesList = append(sourcesList, map[string]interface{}{
				"content":  v.Content,
				"enabled":  v.Enabled,
				"port":     v.Port,
				"priority": v.Priority,
				"type":     v.Type,
				"weight":   v.Weight,
			})
		}
		mapping["sources"] = sourcesList
		ids[i] = object.DomainName
		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			names[i] = object.DomainName
			s[i] = mapping
			continue
		}

		request := dcdn.CreateDescribeDcdnDomainCertificateInfoRequest()
		request.RegionId = client.RegionId
		request.DomainName = object.DomainName
		raw, err := client.WithDcdnClient(func(dcdnClient *dcdn.Client) (interface{}, error) {
			return dcdnClient.DescribeDcdnDomainCertificateInfo(request)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_dcdn_domains", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		responseGet, _ := raw.(*dcdn.DescribeDcdnDomainCertificateInfoResponse)
		mapping["cert_name"] = responseGet.CertInfos.CertInfo[0].CertName
		mapping["ssl_pub"] = responseGet.CertInfos.CertInfo[0].SSLPub

		request1 := dcdn.CreateDescribeDcdnDomainDetailRequest()
		request1.RegionId = client.RegionId
		request1.DomainName = object.DomainName
		raw1, err := client.WithDcdnClient(func(dcdnClient *dcdn.Client) (interface{}, error) {
			return dcdnClient.DescribeDcdnDomainDetail(request1)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_dcdn_domains", request1.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request1.GetActionName(), raw1, request1.RpcRequest, request1)
		responseGet1, _ := raw1.(*dcdn.DescribeDcdnDomainDetailResponse)
		mapping["description"] = responseGet1.DomainDetail.Description
		mapping["scope"] = responseGet1.DomainDetail.Scope

		names[i] = object.DomainName
		s[i] = mapping
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
