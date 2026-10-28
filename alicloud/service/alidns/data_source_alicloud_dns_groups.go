package alidns

import (
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	alidnssdk "github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlicloudDnsGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudDnsGroupsRead,

		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},

			// Computed values
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"group_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudDnsGroupsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := alidnssdk.CreateDescribeDomainGroupsRequest()

	var allGroups []alidnssdk.DomainGroup
	request.RegionId = client.RegionId
	request.PageSize = requests.NewInteger(helper.PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)
	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}
	for {
		raw, err := client.WithDnsClient(func(dnsClient *alidnssdk.Client) (interface{}, error) {
			return dnsClient.DescribeDomainGroups(request)
		})
		if err != nil {
			return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_dns_groups", request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		}
		logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*alidnssdk.DescribeDomainGroupsResponse)
		groups := response.DomainGroups.DomainGroup
		for _, domainGroup := range groups {
			if len(idsMap) > 0 {
				if _, ok := idsMap[domainGroup.GroupId]; !ok {
					continue
				}
			}
			allGroups = append(allGroups, domainGroup)
		}
		if len(groups) < helper.PageSizeLarge {
			break
		}
		if page, err := helper.GetNextpageNumber(request.PageNumber); err != nil {
			return errs.WrapError(err)
		} else {
			request.PageNumber = page
		}
	}

	var filteredGroups []alidnssdk.DomainGroup
	if v, ok := d.GetOk("name_regex"); ok && v.(string) != "" {
		r := regexp.MustCompile(v.(string))

		for _, group := range allGroups {
			if r.MatchString(group.GroupName) {
				filteredGroups = append(filteredGroups, group)
			}
		}
	} else {
		filteredGroups = allGroups[:]
	}

	return groupsDecriptionAttributes(d, filteredGroups, meta)
}

func groupsDecriptionAttributes(d *schema.ResourceData, groupTypes []alidnssdk.DomainGroup, meta interface{}) error {
	var ids []string
	var names []string
	var s []map[string]interface{}
	for _, group := range groupTypes {
		mapping := map[string]interface{}{
			"group_id":   group.GroupId,
			"group_name": group.GroupName,
		}
		ids = append(ids, group.GroupId)
		names = append(names, group.GroupName)
		s = append(s, mapping)
	}

	d.SetId(helper.DataResourceIdHash(ids))
	if err := d.Set("groups", s); err != nil {
		return errs.WrapError(err)
	}
	if err := d.Set("ids", ids); err != nil {
		return errs.WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return errs.WrapError(err)
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		helper.WriteToFile(output.(string), s)
	}
	return nil
}
