package alicloud

import (
	"fmt"
	"regexp"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
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
			"message_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"ConfigurationItemChangeNotification", "ScheduledNotification", "ConfigurationSnapshotDeliveryCompleted"}, false),
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
						"compliance": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"compliance_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"count": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
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
						"event_source": {
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
						"scope_compliance_resource_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"scope_compliance_resource_types": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"source_detail_message_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_identifier": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_maximum_execution_frequency": {
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

	action := "ListConfigRules"
	request := make(map[string]interface{})
	if v, ok := d.GetOk("config_rule_state"); ok {
		request["ConfigRuleState"] = v
	}
	if v, ok := d.GetOk("member_id"); ok {
		request["MemberId"] = v
	}
	if v, ok := d.GetOk("message_type"); ok {
		request["MessageType"] = v
	}
	if v, ok := d.GetOkExists("multi_account"); ok {
		request["MultiAccount"] = v
	}
	if v, ok := d.GetOk("risk_level"); ok {
		request["RiskLevel"] = v
	}
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1
	var objects []map[string]interface{}
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
	var response map[string]interface{}
	conn, err := client.NewConfigClient()
	if err != nil {
		return WrapError(err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("GET"), StringPointer("2019-01-08"), StringPointer("AK"), request, nil, &runtime)
			if err != nil {
				if IsExpectedErrors(err, []string{"Throttling.User"}) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_config_rules", action, AlibabaCloudSdkGoERROR)
		}

		resp, err := jsonpath.Get("$.ConfigRules.ConfigRuleList", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.ConfigRules.ConfigRuleList", response)
		}
		for _, v := range resp.([]interface{}) {
			item := v.(map[string]interface{})
			if ruleNameRegex != nil {
				if !ruleNameRegex.MatchString(item["ConfigRuleName"].(string)) {
					continue
				}
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["ConfigRuleId"])]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}
		if len(resp.([]interface{})) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	ids := make([]string, 0)
	names := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"account_id":        formatInt(object["AccountId"]),
			"config_rule_arn":   object["ConfigRuleArn"],
			"id":                fmt.Sprint(object["ConfigRuleId"]),
			"config_rule_id":    fmt.Sprint(object["ConfigRuleId"]),
			"config_rule_state": object["ConfigRuleState"],
			"description":       object["Description"],
			"risk_level":        formatInt(object["RiskLevel"]),
			"rule_name":         object["ConfigRuleName"],
			"source_identifier": object["SourceIdentifier"],
			"source_owner":      object["SourceOwner"],
		}

		complianceSli := make([]map[string]interface{}, 0)
		if len(object["Compliance"].(map[string]interface{})) > 0 {
			compliance := object["Compliance"]
			complianceMap := make(map[string]interface{})
			complianceMap["compliance_type"] = compliance.(map[string]interface{})["ComplianceType"]
			complianceMap["count"] = compliance.(map[string]interface{})["Count"]
			complianceSli = append(complianceSli, complianceMap)
		}
		mapping["compliance"] = complianceSli

		ids = append(ids, fmt.Sprint(object["ConfigRuleId"]))
		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			names = append(names, object["ConfigRuleName"].(string))
			s = append(s, mapping)
			continue
		}

		configService := ConfigService{client}
		id := fmt.Sprint(object["ConfigRuleId"])
		getResp, err := configService.DescribeConfigRule(id)
		if err != nil {
			return WrapError(err)
		}
		mapping["create_timestamp"] = getResp["CreateTimestamp"]
		mapping["input_parameters"] = getResp["InputParameters"]
		mapping["modified_timestamp"] = getResp["ModifiedTimestamp"]
		mapping["scope_compliance_resource_id"] = getResp["Scope"].(map[string]interface{})["ComplianceResourceId"]
		mapping["scope_compliance_resource_types"] = getResp["Scope"].(map[string]interface{})["ComplianceResourceTypes"]
		mapping["source_maximum_execution_frequency"] = getResp["MaximumExecutionFrequency"]
		if v := getResp["Source"].(map[string]interface{})["SourceDetails"].([]interface{}); len(v) > 0 {
			mapping["event_source"] = v[0].(map[string]interface{})["EventSource"]
			mapping["source_detail_message_type"] = v[0].(map[string]interface{})["MessageType"]
		}
		names = append(names, object["ConfigRuleName"].(string))
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
