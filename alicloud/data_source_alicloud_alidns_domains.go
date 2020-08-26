package alicloud

import (
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudAlidnsDomains() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudAlidnsDomainsRead,
		Schema: map[string]*schema.Schema{
			"ali_domain": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"domain_name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
				ForceNew:     true,
			},
			"group_name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
				ForceNew:     true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"version_code": {
				Type:     schema.TypeString,
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
			"group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"key_word": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"lang": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"search_mode": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"starmark": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"tags": tagsSchema(),
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"domains": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ali_domain": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"dns_servers": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"domain_id": {
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
						"group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"group_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"in_black_hole": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"in_clean": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"line_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"min_ttl": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"puny_code": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"record_line_tree_json": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region_lines": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"remark": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"slave_dns": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"tags": {
							Type:     schema.TypeMap,
							Computed: true,
						},
						"version_code": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"version_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudAlidnsDomainsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := alidns.CreateDescribeDomainsRequest()
	if v, ok := d.GetOk("group_id"); ok {
		request.GroupId = v.(string)
	}
	if v, ok := d.GetOk("key_word"); ok {
		request.KeyWord = v.(string)
	}
	if v, ok := d.GetOk("lang"); ok {
		request.Lang = v.(string)
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request.ResourceGroupId = v.(string)
	}
	if v, ok := d.GetOk("search_mode"); ok {
		request.SearchMode = v.(string)
	}
	if v, ok := d.GetOkExists("starmark"); ok {
		request.Starmark = requests.NewBoolean(v.(bool))
	}
	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)
	var objects []alidns.Domain
	var domain_nameRegex *regexp.Regexp
	if v, ok := d.GetOk("domain_name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		domain_nameRegex = r
	}
	var group_nameRegex *regexp.Regexp
	if v, ok := d.GetOk("group_name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		group_nameRegex = r
	}

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			idsMap[vv.(string)] = vv.(string)
		}
	}
	tagsMap := make(map[string]interface{})
	if v, ok := d.GetOk("tags"); ok && len(v.(map[string]interface{})) > 0 {
		tagsMap = v.(map[string]interface{})
	}
	for {
		raw, err := client.WithDnsClient(func(alidnsClient *alidns.Client) (interface{}, error) {
			return alidnsClient.DescribeDomains(request)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_alidns_domains", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw)
		response, _ := raw.(*alidns.DescribeDomainsResponse)

		for _, item := range response.Domains.Domain {
			if v, ok := d.GetOk("ali_domain"); ok && item.AliDomain != v.(bool) {
				continue
			}
			if domain_nameRegex != nil {
				if !domain_nameRegex.MatchString(item.DomainName) {
					continue
				}
			}
			if group_nameRegex != nil {
				if !group_nameRegex.MatchString(item.GroupName) {
					continue
				}
			}
			if v, ok := d.GetOk("instance_id"); ok && v.(string) != "" && item.InstanceId != v.(string) {
				continue
			}
			if v, ok := d.GetOk("version_code"); ok && v.(string) != "" && item.VersionCode != v.(string) {
				continue
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[item.DomainName]; !ok {
					continue
				}
			}
			if len(tagsMap) > 0 {
				if len(item.Tags.Tag) != len(tagsMap) {
					continue
				}
				match := true
				for _, tag := range item.Tags.Tag {
					if v, ok := tagsMap[tag.Key]; !ok || v.(string) != tag.Value {
						match = false
						break
					}
				}
				if !match {
					continue
				}
			}
			objects = append(objects, item)
		}
		if len(response.Domains.Domain) < PageSizeLarge {
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
		request := alidns.CreateDescribeDomainInfoRequest()
		request.RegionId = client.RegionId
		request.DomainName = object.DomainName
		raw, err := client.WithDnsClient(func(alidnsClient *alidns.Client) (interface{}, error) {
			return alidnsClient.DescribeDomainInfo(request)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_alidns_domains", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*alidns.DescribeDomainInfoResponse)
		mapping := map[string]interface{}{
			"ali_domain":            object.AliDomain,
			"dns_servers":           response.DnsServers.DnsServer,
			"domain_id":             object.DomainId,
			"id":                    object.DomainName,
			"domain_name":           object.DomainName,
			"group_id":              object.GroupId,
			"group_name":            object.GroupName,
			"in_black_hole":         response.InBlackHole,
			"in_clean":              response.InClean,
			"instance_id":           object.InstanceId,
			"line_type":             response.LineType,
			"min_ttl":               response.MinTtl,
			"puny_code":             object.PunyCode,
			"record_line_tree_json": response.RecordLineTreeJson,
			"region_lines":          response.RegionLines,
			"remark":                object.Remark,
			"slave_dns":             response.SlaveDns,
			"version_code":          object.VersionCode,
			"version_name":          object.VersionName,
		}
		tags := make(map[string]string)
		for _, t := range object.Tags.Tag {
			tags[t.Key] = t.Value
		}
		mapping["tags"] = tags
		ids[i] = object.DomainName
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
