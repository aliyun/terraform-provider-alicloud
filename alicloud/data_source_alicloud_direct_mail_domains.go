package alicloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"regexp"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAliCloudDirectMailDomains() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudDirectMailDomainsRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.ValidateRegexp,
			},
			"key_word": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"0", "1", "2", "3", "4"}, false),
			},
			"enable_details": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"domains": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"domain_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"domain_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"domain_record": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"domain_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cname_auth_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cname_confirm_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cname_record": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"icp_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mx_auth_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mx_record": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"spf_auth_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"spf_record": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"default_domain": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"host_record": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dns_mx": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dns_txt": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dns_spf": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dns_dmarc": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dkim_auth_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dkim_rr": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dkim_public_key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dmarc_auth_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dmarc_record": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dmarc_host_record": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tl_domain_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tracef_record": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAliCloudDirectMailDomainsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "QueryDomainByParam"
	request := make(map[string]interface{})
	request["PageSize"] = PageSizeLarge
	request["PageNo"] = 1

	if v, ok := d.GetOk("key_word"); ok {
		request["KeyWord"] = v
	}

	if v, ok := d.GetOk("status"); ok {
		request["Status"] = v
	}

	var objects []map[string]interface{}
	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}

	var domainNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		domainNameRegex = r
	}

	var response map[string]interface{}
	var err error
	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("Dm", "2015-11-23", action, nil, request, true)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_direct_mail_domains", action, AlibabaCloudSdkGoERROR)
		}

		resp, err := jsonpath.Get("$.data.domain", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.data.domain", response)
		}

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["DomainId"])]; !ok {
					continue
				}
			}

			if domainNameRegex != nil && !domainNameRegex.MatchString(fmt.Sprint(item["DomainName"])) {
				continue
			}

			objects = append(objects, item)
		}

		if len(result) < request["PageSize"].(int) {
			break
		}

		request["PageNo"] = request["PageNo"].(int) + 1
	}

	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"id":                fmt.Sprint(object["DomainId"]),
			"domain_id":         fmt.Sprint(object["DomainId"]),
			"domain_name":       object["DomainName"],
			"domain_record":     object["DomainRecord"],
			"cname_auth_status": object["CnameAuthStatus"],
			"icp_status":        object["IcpStatus"],
			"mx_auth_status":    object["MxAuthStatus"],
			"spf_auth_status":   object["SpfAuthStatus"],
			"status":            object["DomainStatus"],
			"create_time":       object["CreateTime"],
		}

		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, object["DomainName"])

		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			s = append(s, mapping)
			continue
		}

		id := fmt.Sprint(object["DomainId"])
		dmService := DmService{client}

		directMailDomainDetail, err := dmService.DescribeDirectMailDomain(id)
		if err != nil {
			return WrapError(err)
		}

		mapping["domain_type"] = directMailDomainDetail["DomainType"]
		mapping["cname_confirm_status"] = directMailDomainDetail["CnameConfirmStatus"]
		mapping["cname_record"] = directMailDomainDetail["CnameRecord"]
		mapping["mx_record"] = directMailDomainDetail["MxRecord"]
		mapping["spf_record"] = directMailDomainDetail["SpfRecord"]
		mapping["default_domain"] = directMailDomainDetail["DefaultDomain"]
		mapping["host_record"] = directMailDomainDetail["HostRecord"]
		mapping["dns_mx"] = directMailDomainDetail["DnsMx"]
		mapping["dns_txt"] = directMailDomainDetail["DnsTxt"]
		mapping["dns_spf"] = directMailDomainDetail["DnsSpf"]
		mapping["dns_dmarc"] = directMailDomainDetail["DnsDmarc"]
		mapping["dkim_auth_status"] = directMailDomainDetail["DkimAuthStatus"]
		mapping["dkim_rr"] = directMailDomainDetail["DkimRR"]
		mapping["dkim_public_key"] = directMailDomainDetail["DkimPublicKey"]
		mapping["dmarc_auth_status"] = directMailDomainDetail["DmarcAuthStatus"]
		mapping["dmarc_record"] = directMailDomainDetail["DmarcRecord"]
		mapping["dmarc_host_record"] = directMailDomainDetail["DmarcHostRecord"]
		mapping["tl_domain_name"] = directMailDomainDetail["TlDomainName"]
		mapping["tracef_record"] = directMailDomainDetail["TracefRecord"]

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
