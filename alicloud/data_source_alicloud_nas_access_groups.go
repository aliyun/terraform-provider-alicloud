package alicloud

import (
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/nas"
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
			"type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},

			// groups values
			"groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"rule_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mount_target_count": {
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
	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		if r, err := regexp.Compile(Trim(v.(string))); err == nil {
			nameRegex = r
		}
	}
	invoker := NewInvoker()
	for {
		var raw interface{}
		if err := invoker.Run(func() error {
			rsp, err := client.WithNasClient(func(nasClient *nas.Client) (interface{}, error) {
				return nasClient.DescribeAccessGroups(args)
			})
			raw = rsp
			return err
		}); err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_nas_access_groups", args.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		resp, _ := raw.(*nas.DescribeAccessGroupsResponse)
		if resp == nil || len(resp.AccessGroups.AccessGroup) < 1 {
			break
		}
		for _, ag := range resp.AccessGroups.AccessGroup {
			if v, ok := d.GetOk("type"); ok && ag.AccessGroupType != Trim(v.(string)) {
				continue
			}

			if v, ok := d.GetOk("description"); ok && ag.Description != v.(string) {
				continue
			}
			if nameRegex != nil {
				if !nameRegex.MatchString(ag.AccessGroupName) {
					continue
				}
			}
			allAgs = append(allAgs, ag)
		}
		if len(resp.AccessGroups.AccessGroup) < PageSizeLarge {
			break
		}

		if page, err := getNextpageNumber(args.PageNumber); err != nil {
			return WrapError(err)
		} else {
			args.PageNumber = page
		}
	}
	return AccessGroupsDecriptionAttributes(d, allAgs, meta)
}

func AccessGroupsDecriptionAttributes(d *schema.ResourceData, nasSetTypes []nas.AccessGroup, meta interface{}) error {
	var ids []string
	var s []map[string]interface{}
	for _, ag := range nasSetTypes {
		mapping := map[string]interface{}{
			"id":                 ag.AccessGroupName,
			"type":               ag.AccessGroupType,
			"description":        ag.Description,
			"mount_target_count": ag.MountTargetCount,
			"rule_count":         ag.RuleCount,
		}
		ids = append(ids, ag.AccessGroupName)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("groups", s); err != nil {
		return WrapError(err)
	}
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
