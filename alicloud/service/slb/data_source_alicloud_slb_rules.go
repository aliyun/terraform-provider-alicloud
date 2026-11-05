package slb

import (
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceRules() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceRulesRead,

		Schema: map[string]*schema.Schema{
			"load_balancer_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"frontend_port": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				MinItems: 1,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
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
			"slb_rules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"domain": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"server_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceRulesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := slb.CreateDescribeRulesRequest()
	request.RegionId = client.RegionId
	request.LoadBalancerId = d.Get("load_balancer_id").(string)
	request.ListenerPort = requests.NewInteger(d.Get("frontend_port").(int))

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[helper.Trim(vv.(string))] = helper.Trim(vv.(string))
		}
	}

	raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
		return slbClient.DescribeRules(request)
	})
	if err != nil {
		return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_slb_rules", request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*slb.DescribeRulesResponse)
	var filteredRulesTemp []slb.Rule
	nameRegex, ok := d.GetOk("name_regex")
	if (ok && nameRegex.(string) != "") || (len(idsMap) > 0) {
		var r *regexp.Regexp
		if nameRegex != "" {
			r, err = regexp.Compile(nameRegex.(string))
			if err != nil {
				return errs.WrapError(err)
			}
		}
		for _, rule := range response.Rules.Rule {
			if r != nil && !r.MatchString(rule.RuleName) {
				continue
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[rule.RuleId]; !ok {
					continue
				}
			}

			filteredRulesTemp = append(filteredRulesTemp, rule)
		}
	} else {
		filteredRulesTemp = response.Rules.Rule
	}

	return slbRulesDescriptionAttributes(d, filteredRulesTemp)
}

func slbRulesDescriptionAttributes(d *schema.ResourceData, rules []slb.Rule) error {
	var ids []string
	var names []string
	var s []map[string]interface{}

	for _, rule := range rules {
		mapping := map[string]interface{}{
			"id":              rule.RuleId,
			"name":            rule.RuleName,
			"domain":          rule.Domain,
			"url":             rule.Url,
			"server_group_id": rule.VServerGroupId,
		}

		ids = append(ids, rule.RuleId)
		names = append(names, rule.RuleName)
		s = append(s, mapping)
	}

	d.SetId(helper.DataResourceIdHash(ids))
	if err := d.Set("slb_rules", s); err != nil {
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
