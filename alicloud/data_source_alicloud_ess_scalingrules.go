package alicloud

import (
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func dataSourceAlicloudEssScalingRules() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudEssScalingRulesRead,
		Schema: map[string]*schema.Schema{
			"scaling_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateNameRegex,
				ForceNew:     true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				ForceNew: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"rules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"scaling_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cooldown": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"adjustment_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"adjustment_value": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"min_adjustment_magnitude": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"scaling_rule_ari": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudEssScalingRulesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	args := ess.CreateDescribeScalingRulesRequest()
	args.PageSize = requests.NewInteger(PageSizeLarge)
	args.PageNumber = requests.NewInteger(1)

	if scalingGroupId, ok := d.GetOk("scaling_group_id"); ok && scalingGroupId.(string) != "" {
		args.ScalingGroupId = scalingGroupId.(string)
	}

	if ruleType, ok := d.GetOk("type"); ok && ruleType.(string) != "" {
		args.ScalingRuleType = ruleType.(string)
	}

	var allScalingRules []ess.ScalingRule

	for {
		ram, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
			return essClient.DescribeScalingRules(args)
		})
		if err != nil {
			return err
		}

		resp, _ := ram.(*ess.DescribeScalingRulesResponse)
		if resp == nil || len(resp.ScalingRules.ScalingRule) < 1 {
			break
		}

		allScalingRules = append(allScalingRules, resp.ScalingRules.ScalingRule...)

		if len(resp.ScalingRules.ScalingRule) < PageSizeLarge {
			break
		}

		if page, err := getNextpageNumber(args.PageNumber); err != nil {
			return err
		} else {
			args.PageNumber = page
		}
	}

	var filteredScalingRulesTemp = make([]ess.ScalingRule, 0)

	nameRegex, okNameRegex := d.GetOk("name_regex")
	idsMap := make(map[string]string)
	ids, okIds := d.GetOk("ids")
	if okIds {
		for _, i := range ids.([]interface{}) {
			idsMap[i.(string)] = i.(string)
		}
	}

	if okNameRegex || okIds {
		for _, rule := range allScalingRules {
			if okNameRegex && nameRegex != "" {
				var r = regexp.MustCompile(nameRegex.(string))
				if r != nil && !r.MatchString(rule.ScalingRuleName) {
					continue
				}
			}
			if okIds && len(idsMap) > 0 {
				if _, ok := idsMap[rule.ScalingRuleId]; !ok {
					continue
				}
			}
			filteredScalingRulesTemp = append(filteredScalingRulesTemp, rule)
		}
	} else {
		filteredScalingRulesTemp = allScalingRules
	}
	return scalingRulesDescriptionAttribute(d, filteredScalingRulesTemp, meta)
}

func scalingRulesDescriptionAttribute(d *schema.ResourceData, scalingRules []ess.ScalingRule, meta interface{}) error {
	var ids []string
	var names []string
	var s = make([]map[string]interface{}, 0)
	for _, scalingRule := range scalingRules {
		mapping := map[string]interface{}{
			"id":                       scalingRule.ScalingRuleId,
			"scaling_group_id":         scalingRule.ScalingGroupId,
			"name":                     scalingRule.ScalingRuleName,
			"type":                     scalingRule.ScalingRuleType,
			"cooldown":                 scalingRule.Cooldown,
			"adjustment_type":          scalingRule.AdjustmentType,
			"adjustment_value":         scalingRule.AdjustmentValue,
			"min_adjustment_magnitude": scalingRule.MinAdjustmentMagnitude,
			"scaling_rule_ari":         scalingRule.ScalingRuleAri,
		}
		ids = append(ids, scalingRule.ScalingRuleId)
		names = append(names, scalingRule.ScalingRuleName)
		s = append(s, mapping)
	}
	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("rules", s); err != nil {
		return err
	}

	if err := d.Set("ids", ids); err != nil {
		return err
	}

	if err := d.Set("names", names); err != nil {
		return err
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
