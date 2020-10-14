package alicloud

import (
	"regexp"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/config"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudConfigRules() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudConfigRulesRead,
		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
				ForceNew:     true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"config_rule_state": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"ACTIVE", "DELETING", "DELETING_RESULTS", "EVALUATING", "INACTIVE"}, false),
			},
			"member_id": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"multi_account": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
			},
			"risk_level": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntInSlice([]int{1, 2, 3}),
			},
			"names": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"rules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"account_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"config_rule_arn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"config_rule_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"config_rule_state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_timestamp": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"input_parameters": {
							Type:     schema.TypeMap,
							Computed: true,
						},
						"modified_timestamp": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"risk_level": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"rule_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_details": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"event_source": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"maximum_execution_frequency": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"message_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"source_identifier": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_owner": {
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

func dataSourceAlicloudConfigRulesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := config.CreateListConfigRulesRequest()
	if v, ok := d.GetOk("config_rule_state"); ok {
		request.ConfigRuleState = v.(string)
	}
	if v, ok := d.GetOk("member_id"); ok {
		request.MemberId = requests.NewInteger(v.(int))
	}
	if v, ok := d.GetOkExists("multi_account"); ok {
		request.MultiAccount = requests.NewBoolean(v.(bool))
	}
	if v, ok := d.GetOk("risk_level"); ok {
		request.RiskLevel = requests.NewInteger(v.(int))
	}
	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)
	var objects []config.ConfigRule
	var ruleNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		ruleNameRegex = r
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
	var response *config.ListConfigRulesResponse
	for {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
			raw, err := client.WithConfigClient(func(configClient *config.Client) (interface{}, error) {
				return configClient.ListConfigRules(request)
			})
			if err != nil {
				if IsExpectedErrors(err, []string{"Throttling.User"}) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(request.GetActionName(), raw)
			response, _ = raw.(*config.ListConfigRulesResponse)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_config_rules", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}

		for _, item := range response.ConfigRules.ConfigRuleList {
			if ruleNameRegex != nil {
				if !ruleNameRegex.MatchString(item.ConfigRuleName) {
					continue
				}
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[item.ConfigRuleId]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}
		if len(response.ConfigRules.ConfigRuleList) < PageSizeLarge {
			break
		}

		page, err := getNextpageNumber(request.PageNumber)
		if err != nil {
			return WrapError(err)
		}
		request.PageNumber = page
	}
	ids := make([]string, 0)
	names := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"account_id":        object.AccountId,
			"config_rule_arn":   object.ConfigRuleArn,
			"id":                object.ConfigRuleId,
			"config_rule_id":    object.ConfigRuleId,
			"config_rule_state": object.ConfigRuleState,
			"description":       object.Description,
			"risk_level":        object.RiskLevel,
			"rule_name":         object.ConfigRuleName,
			"source_identifier": object.SourceIdentifier,
			"source_owner":      object.SourceOwner,
		}
		ids = append(ids, object.ConfigRuleId)
		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			names = append(names, object.ConfigRuleName)
			s = append(s, mapping)
			continue
		}

		request := config.CreateDescribeConfigRuleRequest()
		request.RegionId = client.RegionId
		request.ConfigRuleId = object.ConfigRuleId
		raw, err := client.WithConfigClient(func(configClient *config.Client) (interface{}, error) {
			return configClient.DescribeConfigRule(request)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_config_rules", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		responseGet, _ := raw.(*config.DescribeConfigRuleResponse)
		mapping["create_timestamp"] = responseGet.ConfigRule.CreateTimestamp
		mapping["input_parameters"] = responseGet.ConfigRule.InputParameters
		mapping["modified_timestamp"] = responseGet.ConfigRule.ModifiedTimestamp

		sourceDetails := make([]map[string]interface{}, len(responseGet.ConfigRule.ManagedRule.SourceDetails))
		for i, v := range responseGet.ConfigRule.ManagedRule.SourceDetails {
			mapping1 := map[string]interface{}{
				"event_source":                v.EventSource,
				"maximum_execution_frequency": v.MaximumExecutionFrequency,
				"message_type":                v.MessageType,
			}
			sourceDetails[i] = mapping1
		}
		mapping["source_details"] = sourceDetails
		names = append(names, object.ConfigRuleName)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("rules", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
