package alidns

import (
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	alidnssdk "github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/tags"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlicloudAlidnsDomains() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudAlidnsDomainsRead,
		Schema: map[string]*schema.Schema{
			"ali_domain": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"domain_name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
			},
			"group_name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"version_code": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"names": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"key_word": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"lang": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"search_mode": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"starmark": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"tags": tags.TagsSchema(),
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
						"available_ttls": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeInt},
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
						"record_lines": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"father_code": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"line_code": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"line_display_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"line_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"region_lines": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"remark": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_group_id": {
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
							Elem:     &schema.Schema{Type: schema.TypeString},
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
			"enable_details": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func dataSourceAlicloudAlidnsDomainsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := alidnssdk.CreateDescribeDomainsRequest()
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
	request.PageSize = requests.NewInteger(helper.PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)
	var objects []alidnssdk.DomainInDescribeDomains
	var domainNameRegex *regexp.Regexp
	if v, ok := d.GetOk("domain_name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return errs.WrapError(err)
		}
		domainNameRegex = r
	}
	var groupNameRegex *regexp.Regexp
	if v, ok := d.GetOk("group_name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return errs.WrapError(err)
		}
		groupNameRegex = r
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
	tagsMap := make(map[string]interface{})
	if v, ok := d.GetOk("tags"); ok && len(v.(map[string]interface{})) > 0 {
		tagsMap = v.(map[string]interface{})
	}
	var response *alidnssdk.DescribeDomainsResponse
	for {
		raw, err := client.WithAlidnsClient(func(alidnsClient *alidnssdk.Client) (interface{}, error) {
			return alidnsClient.DescribeDomains(request)
		})
		if err != nil {
			return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_alidns_domains", request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		}
		logs.AddDebug(request.GetActionName(), raw)
		response, _ = raw.(*alidnssdk.DescribeDomainsResponse)

		for _, item := range response.Domains.Domain {
			if v, ok := d.GetOk("ali_domain"); ok && item.AliDomain != v.(bool) {
				continue
			}
			if domainNameRegex != nil {
				if !domainNameRegex.MatchString(item.DomainName) {
					continue
				}
			}
			if groupNameRegex != nil {
				if !groupNameRegex.MatchString(item.GroupName) {
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
		if len(response.Domains.Domain) < helper.PageSizeLarge {
			break
		}

		page, err := helper.GetNextpageNumber(request.PageNumber)
		if err != nil {
			return errs.WrapError(err)
		}
		request.PageNumber = page
	}
	ids := make([]string, 0)
	names := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"ali_domain":        object.AliDomain,
			"domain_id":         object.DomainId,
			"id":                object.DomainName,
			"domain_name":       object.DomainName,
			"group_id":          object.GroupId,
			"group_name":        object.GroupName,
			"instance_id":       object.InstanceId,
			"puny_code":         object.PunyCode,
			"remark":            object.Remark,
			"resource_group_id": object.ResourceGroupId,
			"version_code":      object.VersionCode,
			"version_name":      object.VersionName,
		}
		ids = append(ids, object.DomainName)
		tags := make(map[string]string)
		for _, t := range object.Tags.Tag {
			tags[t.Key] = t.Value
		}
		mapping["tags"] = tags
		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			names = append(names, object.DomainName)
			s = append(s, mapping)
			continue
		}

		request := alidnssdk.CreateDescribeDomainInfoRequest()
		request.RegionId = client.RegionId
		request.DomainName = object.DomainName
		raw, err := client.WithAlidnsClient(func(alidnsClient *alidnssdk.Client) (interface{}, error) {
			return alidnsClient.DescribeDomainInfo(request)
		})
		if err != nil {
			return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_alidns_domains", request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		}
		logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
		responseGet, _ := raw.(*alidnssdk.DescribeDomainInfoResponse)
		mapping["available_ttls"] = responseGet.AvailableTtls.AvailableTtl
		mapping["dns_servers"] = responseGet.DnsServers.DnsServer
		mapping["in_black_hole"] = responseGet.InBlackHole
		mapping["in_clean"] = responseGet.InClean
		mapping["line_type"] = responseGet.LineType
		mapping["min_ttl"] = responseGet.MinTtl
		mapping["record_line_tree_json"] = responseGet.RecordLineTreeJson

		recordLines := make([]map[string]interface{}, len(responseGet.RecordLines.RecordLine))
		for i, v := range responseGet.RecordLines.RecordLine {
			mapping1 := map[string]interface{}{
				"father_code":       v.FatherCode,
				"line_code":         v.LineCode,
				"line_display_name": v.LineDisplayName,
				"line_name":         v.LineName,
			}
			recordLines[i] = mapping1
		}
		mapping["record_lines"] = recordLines
		mapping["region_lines"] = responseGet.RegionLines
		mapping["slave_dns"] = responseGet.SlaveDns
		names = append(names, object.DomainName)
		s = append(s, mapping)
	}

	d.SetId(helper.DataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return errs.WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return errs.WrapError(err)
	}

	if err := d.Set("domains", s); err != nil {
		return errs.WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		helper.WriteToFile(output.(string), s)
	}

	return nil
}
