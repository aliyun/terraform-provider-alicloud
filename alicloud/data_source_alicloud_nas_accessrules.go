package alicloud

import (
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/nas"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func dataSourceAlicloudAccessRules() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudAccessRulesRead,

		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateNameRegex,
				ForceNew:     true,
			},
			"sourcecidr_ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"accessgroup_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},

			// Computed values
			"accessrules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"sourcecidr_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"priority": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"accessrule_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"user_access": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"rw_access": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudAccessRulesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	args := nas.CreateDescribeAccessRulesRequest()
	args.AccessGroupName = d.Get("accessgroup_name").(string)
	args.RegionId = string(client.Region)
	args.PageSize = requests.NewInteger(PageSizeLarge)
	args.PageNumber = requests.NewInteger(1)

	var allAgs []nas.AccessRule

	for {
		raw, err := client.WithNasClient(func(nasClient *nas.Client) (interface{}, error) {
			return nasClient.DescribeAccessRules(args)
		})
		if err != nil {
			return err
		}
		resp, _ := raw.(*nas.DescribeAccessRulesResponse)
		if resp == nil || len(resp.AccessRules.AccessRule) < 1 {
			break
		}

		allAgs = append(allAgs, resp.AccessRules.AccessRule...)

		if len(resp.AccessRules.AccessRule) < PageSizeLarge {
			break
		}

		if page, err := getNextpageNumber(args.PageNumber); err != nil {
			return err
		} else {
			args.PageNumber = page
		}
	}

	var filteredArs []nas.AccessRule
	var r *regexp.Regexp
	if nameRegex, ok := d.GetOk("name_regex"); ok && nameRegex.(string) != "" {
		r = regexp.MustCompile(nameRegex.(string))
	}

	for _, v := range allAgs {
		if r != nil && !r.MatchString(v.AccessRuleId) {
			continue
		}

		filteredArs = append(filteredArs, v)
	}

	return accessRulesDecriptionAttributes(d, filteredArs, meta)
}

func accessRulesDecriptionAttributes(d *schema.ResourceData, nasSetTypes []nas.AccessRule, meta interface{}) error {
	var ids []string
	var s []map[string]interface{}

	for _, ag := range nasSetTypes {
		mapping := map[string]interface{}{
			"sourcecidr_ip": ag.SourceCidrIp,
			"priority":      ag.Priority,
			"accessrule_id": ag.AccessRuleId,
			"user_access":   ag.UserAccess,
			"rw_access":     ag.RWAccess,
		}
		ids = append(ids, ag.AccessRuleId)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("accessrules", s); err != nil {
		return err
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
