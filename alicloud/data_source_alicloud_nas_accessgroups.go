package alicloud

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/nas"
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func dataSourceAlicloudAccessGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudAccessGroupsRead,

		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateNameRegex,
				ForceNew:     true,
			},
			"accessgroup_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"accessgroup_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},


			// Computed values
			"accessgroups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"accessgroup_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"rule_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"accessgroup_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mounttarget_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},

					},
				},
			},
		},
	}
}

func dataSourceAlicloudAccessGroupsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	args := nas.CreateDescribeAccessGroupsRequest()
	args.RegionId = string(client.Region)
	args.PageSize = requests.NewInteger(PageSizeLarge)
	args.PageNumber = requests.NewInteger(1)

	var allAgs []nas.AccessGroup

	for {
		raw, err := client.WithNasClient(func(nasClient *nas.Client) (interface{}, error) {
			return nasClient.DescribeAccessGroups(args)
		})
		if err != nil {
			return err
		}
		resp, _ := raw.(*nas.DescribeAccessGroupsResponse)
		if resp == nil || len(resp.AccessGroups.AccessGroup) < 1 {
			break
		}

		allAgs = append(allAgs, resp.AccessGroups.AccessGroup...)

		if len(resp.AccessGroups.AccessGroup) < PageSizeLarge {
			break
		}

		if page, err := getNextpageNumber(args.PageNumber); err != nil {
			return err
		} else {
			args.PageNumber = page
		}
	}

	var filteredAgs []nas.AccessGroup
	var r *regexp.Regexp
	if nameRegex, ok := d.GetOk("name_regex"); ok && nameRegex.(string) != "" {
		r = regexp.MustCompile(nameRegex.(string))
	}

	for _, v := range allAgs {
		if r != nil && !r.MatchString(v.AccessGroupName) {
			continue
		}

		filteredAgs = append(filteredAgs, v)
	}

	return accessGroupsDecriptionAttributes(d, filteredAgs, meta)
}

func accessGroupsDecriptionAttributes(d *schema.ResourceData, nasSetTypes []nas.AccessGroup, meta interface{}) error {
	var ids []string
	var s []map[string]interface{}
	for _, ag := range nasSetTypes {
		mapping := map[string]interface{}{
			"accessgroup_name":      ag.AccessGroupName,
			"accessgroup_type":      ag.AccessGroupType,
			"description":         	 ag.Description,
			"mounttarget_count":     ag.MountTargetCount,
			"rule_count":    		 ag.RuleCount,
		}
		ids = append(ids, ag.AccessGroupName)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("accessgroups", s); err != nil {
		return err
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
