// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAliCloudWafv3DefenseRules() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudWafv3DefenseRuleRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"defense_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"query": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"rule_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"rules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"config": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"bypass_tags": {
										Type:     schema.TypeSet,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"gray_config": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"gray_target": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"gray_rate": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"gray_sub_key": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"rule_action": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"mode": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"bypass_regular_rules": {
										Type:     schema.TypeSet,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"cn_regions": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"time_config": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"time_zone": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"time_periods": {
													Type:     schema.TypeSet,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"start": {
																Type:     schema.TypeInt,
																Computed: true,
															},
															"end": {
																Type:     schema.TypeInt,
																Computed: true,
															},
														},
													},
												},
												"time_scope": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"week_time_periods": {
													Type:     schema.TypeSet,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"day_periods": {
																Type:     schema.TypeSet,
																Computed: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"start": {
																			Type:     schema.TypeInt,
																			Computed: true,
																		},
																		"end": {
																			Type:     schema.TypeInt,
																			Computed: true,
																		},
																	},
																},
															},
															"day": {
																Type:     schema.TypeString,
																Computed: true,
															},
														},
													},
												},
											},
										},
									},
									"bypass_regular_types": {
										Type:     schema.TypeSet,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"ua": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"remote_addr": {
										Type:     schema.TypeSet,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"conditions": {
										Type:     schema.TypeSet,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"op_value": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"values": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"key": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"sub_key": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"url": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"abroad_regions": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"rate_limit": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"status": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"ratio": {
																Type:     schema.TypeInt,
																Computed: true,
															},
															"count": {
																Type:     schema.TypeInt,
																Computed: true,
															},
															"code": {
																Type:     schema.TypeInt,
																Computed: true,
															},
														},
													},
												},
												"target": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"ttl": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"interval": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"threshold": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"sub_key": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"cc_effect": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"throttle_threhold": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"throttle_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"auto_update": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"cc_status": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"account_identifiers": {
										Type:     schema.TypeSet,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"position": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"priority": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"decode_type": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"key": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"sub_key": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"gray_status": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"waf_base_config": {
										Type:     schema.TypeSet,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"rule_type": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"rule_detail": {
													Type:     schema.TypeSet,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"rule_id": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"rule_action": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"rule_status": {
																Type:     schema.TypeInt,
																Computed: true,
															},
														},
													},
												},
												"rule_batch_operation_config": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"protocol": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"codec_list": {
										Type:     schema.TypeSet,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"defense_origin": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"defense_scene": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"defense_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"gmt_modified": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"rule_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"rule_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"rule_status": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"template_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_details": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func dataSourceAliCloudWafv3DefenseRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	var objects []map[string]interface{}

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}

	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	action := "DescribeDefenseRules"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["InstanceId"] = d.Get("instance_id")
	if v, ok := d.GetOk("defense_type"); ok {
		request["DefenseType"] = v
	}
	if v, ok := d.GetOk("query"); ok {
		request["Query"] = v
	}
	if v, ok := d.GetOk("rule_type"); ok {
		request["RuleType"] = v
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1
	for {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutRead), func() *resource.RetryError {
			response, err = client.RpcPost("waf-openapi", "2021-10-01", action, query, request, true)

			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		resp, _ := jsonpath.Get("$.Rules[*]", response)

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(request["InstanceId"], ":", item["DefenseType"], ":", item["RuleId"])]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}

		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}

	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, objectRaw := range objects {
		mapping := map[string]interface{}{}

		mapping["id"] = fmt.Sprint(request["InstanceId"], ":", objectRaw["DefenseType"], ":", objectRaw["RuleId"])

		mapping["defense_origin"] = objectRaw["DefenseOrigin"]
		mapping["defense_scene"] = objectRaw["DefenseScene"]
		mapping["gmt_modified"] = objectRaw["GmtModified"]
		mapping["resource"] = objectRaw["Resource"]
		mapping["rule_name"] = objectRaw["RuleName"]
		mapping["rule_status"] = objectRaw["Status"]
		mapping["template_id"] = objectRaw["TemplateId"]
		mapping["defense_type"] = objectRaw["DefenseType"]
		mapping["rule_id"] = objectRaw["RuleId"]

		configMaps := make([]map[string]interface{}, 0)
		configMap := make(map[string]interface{})
		configRaw := make(map[string]interface{})
		if objectRaw["Config"] != nil {
			configRaw = convertJsonStringToObject(objectRaw["Config"])
			if len(configRaw) > 0 {
				configMap["abroad_regions"] = configRaw["abroadRegionList"]
				configMap["auto_update"] = formatBool(configRaw["autoUpdate"])
				configMap["cc_effect"] = configRaw["effect"]
				configMap["cc_status"] = formatInt(configRaw["ccStatus"])
				configMap["cn_regions"] = configRaw["cnRegionList"]
				configMap["gray_status"] = formatInt(configRaw["grayStatus"])
				configMap["mode"] = formatInt(configRaw["mode"])
				configMap["protocol"] = configRaw["protocol"]
				configMap["rule_action"] = configRaw["action"]
				configMap["throttle_threhold"] = formatInt(configRaw["threshold"])
				configMap["throttle_type"] = configRaw["type"]
				configMap["ua"] = configRaw["ua"]
				configMap["url"] = configRaw["url"]

				accountIdentifiersRaw := configRaw["accountIdentifiers"]
				accountIdentifiersMaps := make([]map[string]interface{}, 0)
				if accountIdentifiersRaw != nil {
					for _, accountIdentifiersChildRaw := range convertToInterfaceArray(accountIdentifiersRaw) {
						accountIdentifiersMap := make(map[string]interface{})
						accountIdentifiersChildRaw := accountIdentifiersChildRaw.(map[string]interface{})
						accountIdentifiersMap["decode_type"] = accountIdentifiersChildRaw["decodeType"]
						accountIdentifiersMap["key"] = accountIdentifiersChildRaw["key"]
						accountIdentifiersMap["position"] = accountIdentifiersChildRaw["position"]
						accountIdentifiersMap["priority"] = formatInt(accountIdentifiersChildRaw["priority"])
						accountIdentifiersMap["sub_key"] = accountIdentifiersChildRaw["subKey"]

						accountIdentifiersMaps = append(accountIdentifiersMaps, accountIdentifiersMap)
					}
				}
				configMap["account_identifiers"] = accountIdentifiersMaps
				regularRulesRaw := make([]interface{}, 0)
				if configRaw["regularRules"] != nil {
					regularRulesRaw = convertToInterfaceArray(configRaw["regularRules"])
				}

				configMap["bypass_regular_rules"] = regularRulesRaw
				regularTypesRaw := make([]interface{}, 0)
				if configRaw["regularTypes"] != nil {
					regularTypesRaw = convertToInterfaceArray(configRaw["regularTypes"])
				}

				configMap["bypass_regular_types"] = regularTypesRaw
				tagsRaw := make([]interface{}, 0)
				if configRaw["tags"] != nil {
					tagsRaw = convertToInterfaceArray(configRaw["tags"])
				}

				configMap["bypass_tags"] = tagsRaw
				codecListRaw := make([]interface{}, 0)
				if configRaw["codecList"] != nil {
					codecListRaw = convertToInterfaceArray(configRaw["codecList"])
				}

				configMap["codec_list"] = codecListRaw
				conditionsRaw := configRaw["conditions"]
				conditionsMaps := make([]map[string]interface{}, 0)
				if conditionsRaw != nil {
					for _, conditionsChildRaw := range convertToInterfaceArray(conditionsRaw) {
						conditionsMap := make(map[string]interface{})
						conditionsChildRaw := conditionsChildRaw.(map[string]interface{})
						conditionsMap["key"] = conditionsChildRaw["key"]
						conditionsMap["op_value"] = conditionsChildRaw["opValue"]
						conditionsMap["sub_key"] = conditionsChildRaw["subKey"]
						conditionsMap["values"] = conditionsChildRaw["values"]

						conditionsMaps = append(conditionsMaps, conditionsMap)
					}
				}
				configMap["conditions"] = conditionsMaps
				grayConfigMaps := make([]map[string]interface{}, 0)
				grayConfigMap := make(map[string]interface{})
				grayConfigRaw := make(map[string]interface{})
				if configRaw["grayConfig"] != nil {
					grayConfigRaw = configRaw["grayConfig"].(map[string]interface{})
				}
				if len(grayConfigRaw) > 0 {
					grayConfigMap["gray_rate"] = formatInt(grayConfigRaw["grayRate"])
					grayConfigMap["gray_sub_key"] = grayConfigRaw["graySubKey"]
					grayConfigMap["gray_target"] = grayConfigRaw["grayTarget"]

					grayConfigMaps = append(grayConfigMaps, grayConfigMap)
				}
				configMap["gray_config"] = grayConfigMaps
				rateLimitMaps := make([]map[string]interface{}, 0)
				rateLimitMap := make(map[string]interface{})
				ratelimitRaw := make(map[string]interface{})
				if configRaw["ratelimit"] != nil {
					ratelimitRaw = configRaw["ratelimit"].(map[string]interface{})
				}
				if len(ratelimitRaw) > 0 {
					rateLimitMap["interval"] = formatInt(ratelimitRaw["interval"])
					rateLimitMap["sub_key"] = ratelimitRaw["subKey"]
					rateLimitMap["target"] = ratelimitRaw["target"]
					rateLimitMap["threshold"] = formatInt(ratelimitRaw["threshold"])
					rateLimitMap["ttl"] = formatInt(ratelimitRaw["ttl"])

					statusMaps := make([]map[string]interface{}, 0)
					statusMap := make(map[string]interface{})
					statusRaw := make(map[string]interface{})
					if ratelimitRaw["status"] != nil {
						statusRaw = ratelimitRaw["status"].(map[string]interface{})
					}
					if len(statusRaw) > 0 {
						statusMap["code"] = formatInt(statusRaw["code"])
						statusMap["count"] = formatInt(statusRaw["count"])
						statusMap["ratio"] = formatInt(statusRaw["ratio"])

						statusMaps = append(statusMaps, statusMap)
					}
					rateLimitMap["status"] = statusMaps
					rateLimitMaps = append(rateLimitMaps, rateLimitMap)
				}
				configMap["rate_limit"] = rateLimitMaps
				remoteAddrRaw := make([]interface{}, 0)
				if configRaw["remoteAddr"] != nil {
					remoteAddrRaw = convertToInterfaceArray(configRaw["remoteAddr"])
				}

				configMap["remote_addr"] = remoteAddrRaw
				timeConfigMaps := make([]map[string]interface{}, 0)
				timeConfigMap := make(map[string]interface{})
				timeConfigRaw := make(map[string]interface{})
				if configRaw["timeConfig"] != nil {
					timeConfigRaw = configRaw["timeConfig"].(map[string]interface{})
				}
				if len(timeConfigRaw) > 0 {
					timeConfigMap["time_scope"] = timeConfigRaw["timeScope"]
					timeConfigMap["time_zone"] = formatInt(timeConfigRaw["timeZone"])

					timePeriodsRaw := timeConfigRaw["timePeriods"]
					timePeriodsMaps := make([]map[string]interface{}, 0)
					if timePeriodsRaw != nil {
						for _, timePeriodsChildRaw := range convertToInterfaceArray(timePeriodsRaw) {
							timePeriodsMap := make(map[string]interface{})
							timePeriodsChildRaw := timePeriodsChildRaw.(map[string]interface{})
							timePeriodsMap["end"] = formatInt(timePeriodsChildRaw["end"])
							timePeriodsMap["start"] = formatInt(timePeriodsChildRaw["start"])

							timePeriodsMaps = append(timePeriodsMaps, timePeriodsMap)
						}
					}
					timeConfigMap["time_periods"] = timePeriodsMaps
					weekTimePeriodsRaw := timeConfigRaw["weekTimePeriods"]
					weekTimePeriodsMaps := make([]map[string]interface{}, 0)
					if weekTimePeriodsRaw != nil {
						for _, weekTimePeriodsChildRaw := range convertToInterfaceArray(weekTimePeriodsRaw) {
							weekTimePeriodsMap := make(map[string]interface{})
							weekTimePeriodsChildRaw := weekTimePeriodsChildRaw.(map[string]interface{})
							weekTimePeriodsMap["day"] = weekTimePeriodsChildRaw["day"]

							dayPeriodsRaw := weekTimePeriodsChildRaw["dayPeriods"]
							dayPeriodsMaps := make([]map[string]interface{}, 0)
							if dayPeriodsRaw != nil {
								for _, dayPeriodsChildRaw := range convertToInterfaceArray(dayPeriodsRaw) {
									dayPeriodsMap := make(map[string]interface{})
									dayPeriodsChildRaw := dayPeriodsChildRaw.(map[string]interface{})
									dayPeriodsMap["end"] = formatInt(dayPeriodsChildRaw["end"])
									dayPeriodsMap["start"] = formatInt(dayPeriodsChildRaw["start"])

									dayPeriodsMaps = append(dayPeriodsMaps, dayPeriodsMap)
								}
							}
							weekTimePeriodsMap["day_periods"] = dayPeriodsMaps
							weekTimePeriodsMaps = append(weekTimePeriodsMaps, weekTimePeriodsMap)
						}
					}
					timeConfigMap["week_time_periods"] = weekTimePeriodsMaps
					timeConfigMaps = append(timeConfigMaps, timeConfigMap)
				}
				configMap["time_config"] = timeConfigMaps
				config1Raw := configRaw["config"]
				wafBaseConfigMaps := make([]map[string]interface{}, 0)
				if config1Raw != nil {
					for _, configChildRaw := range convertToInterfaceArray(config1Raw) {
						wafBaseConfigMap := make(map[string]interface{})
						configChildRaw := configChildRaw.(map[string]interface{})
						wafBaseConfigMap["rule_batch_operation_config"] = configChildRaw["ruleBatchOperationConfig"]
						wafBaseConfigMap["rule_type"] = configChildRaw["ruleType"]

						ruleDetailRaw := configChildRaw["ruleDetail"]
						ruleDetailMaps := make([]map[string]interface{}, 0)
						if ruleDetailRaw != nil {
							for _, ruleDetailChildRaw := range convertToInterfaceArray(ruleDetailRaw) {
								ruleDetailMap := make(map[string]interface{})
								ruleDetailChildRaw := ruleDetailChildRaw.(map[string]interface{})
								ruleDetailMap["rule_action"] = ruleDetailChildRaw["ruleAction"]
								if v := ruleDetailChildRaw["ruleId"]; v != nil {
									ruleDetailMap["rule_id"] = fmt.Sprint(v)
								}
								ruleDetailMap["rule_status"] = formatInt(ruleDetailChildRaw["ruleStatus"])

								ruleDetailMaps = append(ruleDetailMaps, ruleDetailMap)
							}
						}
						wafBaseConfigMap["rule_detail"] = ruleDetailMaps
						wafBaseConfigMaps = append(wafBaseConfigMaps, wafBaseConfigMap)
					}
				}
				configMap["waf_base_config"] = wafBaseConfigMaps
				configMaps = append(configMaps, configMap)
			}
		}
		mapping["config"] = configMaps

		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			ids = append(ids, fmt.Sprint(mapping["id"]))
			s = append(s, mapping)
			continue
		}

		id := fmt.Sprint(request["InstanceId"], ":", objectRaw["DefenseType"], ":", objectRaw["RuleId"])
		mapping, err = dataSourceAliCloudWafv3DefenseRuleReadDescription(d, id, mapping, meta)
		if err != nil {
			return WrapError(err)
		}

		ids = append(ids, fmt.Sprint(mapping["id"]))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
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

func dataSourceAliCloudWafv3DefenseRuleReadDescription(d *schema.ResourceData, id string, object map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	client := meta.(*connectivity.AliyunClient)

	wafv3ServiceV2 := Wafv3ServiceV2{client}
	getResp, err := wafv3ServiceV2.DescribeWafv3DefenseRule(id)
	if err != nil {
		return nil, WrapError(err)
	}

	// Merge additional fields from Get API response to mapping
	// Reuse the response mapping template from Resource's read function
	mapping := object
	objectRaw := getResp

	mapping["defense_origin"] = objectRaw["DefenseOrigin"]
	mapping["defense_scene"] = objectRaw["DefenseScene"]
	mapping["gmt_modified"] = objectRaw["GmtModified"]
	mapping["resource"] = objectRaw["Resource"]
	mapping["rule_name"] = objectRaw["RuleName"]
	mapping["rule_status"] = objectRaw["Status"]
	mapping["template_id"] = objectRaw["TemplateId"]
	mapping["defense_type"] = objectRaw["DefenseType"]
	mapping["rule_id"] = objectRaw["RuleId"]

	configMaps := make([]map[string]interface{}, 0)
	configMap := make(map[string]interface{})
	configRaw := make(map[string]interface{})
	if objectRaw["Config"] != nil {
		configRaw = convertJsonStringToObject(objectRaw["Config"])
		if len(configRaw) > 0 {
			configMap["abroad_regions"] = configRaw["abroadRegionList"]
			configMap["auto_update"] = formatBool(configRaw["autoUpdate"])
			configMap["cc_effect"] = configRaw["effect"]
			configMap["cc_status"] = formatInt(configRaw["ccStatus"])
			configMap["cn_regions"] = configRaw["cnRegionList"]
			configMap["gray_status"] = formatInt(configRaw["grayStatus"])
			configMap["mode"] = formatInt(configRaw["mode"])
			configMap["protocol"] = configRaw["protocol"]
			configMap["rule_action"] = configRaw["action"]
			configMap["throttle_threhold"] = formatInt(configRaw["threshold"])
			configMap["throttle_type"] = configRaw["type"]
			configMap["ua"] = configRaw["ua"]
			configMap["url"] = configRaw["url"]

			accountIdentifiersRaw := configRaw["accountIdentifiers"]
			accountIdentifiersMaps := make([]map[string]interface{}, 0)
			if accountIdentifiersRaw != nil {
				for _, accountIdentifiersChildRaw := range convertToInterfaceArray(accountIdentifiersRaw) {
					accountIdentifiersMap := make(map[string]interface{})
					accountIdentifiersChildRaw := accountIdentifiersChildRaw.(map[string]interface{})
					accountIdentifiersMap["decode_type"] = accountIdentifiersChildRaw["decodeType"]
					accountIdentifiersMap["key"] = accountIdentifiersChildRaw["key"]
					accountIdentifiersMap["position"] = accountIdentifiersChildRaw["position"]
					accountIdentifiersMap["priority"] = formatInt(accountIdentifiersChildRaw["priority"])
					accountIdentifiersMap["sub_key"] = accountIdentifiersChildRaw["subKey"]

					accountIdentifiersMaps = append(accountIdentifiersMaps, accountIdentifiersMap)
				}
			}
			configMap["account_identifiers"] = accountIdentifiersMaps
			regularRulesRaw := make([]interface{}, 0)
			if configRaw["regularRules"] != nil {
				regularRulesRaw = convertToInterfaceArray(configRaw["regularRules"])
			}

			configMap["bypass_regular_rules"] = regularRulesRaw
			regularTypesRaw := make([]interface{}, 0)
			if configRaw["regularTypes"] != nil {
				regularTypesRaw = convertToInterfaceArray(configRaw["regularTypes"])
			}

			configMap["bypass_regular_types"] = regularTypesRaw
			tagsRaw := make([]interface{}, 0)
			if configRaw["tags"] != nil {
				tagsRaw = convertToInterfaceArray(configRaw["tags"])
			}

			configMap["bypass_tags"] = tagsRaw
			codecListRaw := make([]interface{}, 0)
			if configRaw["codecList"] != nil {
				codecListRaw = convertToInterfaceArray(configRaw["codecList"])
			}

			configMap["codec_list"] = codecListRaw
			conditionsRaw := configRaw["conditions"]
			conditionsMaps := make([]map[string]interface{}, 0)
			if conditionsRaw != nil {
				for _, conditionsChildRaw := range convertToInterfaceArray(conditionsRaw) {
					conditionsMap := make(map[string]interface{})
					conditionsChildRaw := conditionsChildRaw.(map[string]interface{})
					conditionsMap["key"] = conditionsChildRaw["key"]
					conditionsMap["op_value"] = conditionsChildRaw["opValue"]
					conditionsMap["sub_key"] = conditionsChildRaw["subKey"]
					conditionsMap["values"] = conditionsChildRaw["values"]

					conditionsMaps = append(conditionsMaps, conditionsMap)
				}
			}
			configMap["conditions"] = conditionsMaps
			grayConfigMaps := make([]map[string]interface{}, 0)
			grayConfigMap := make(map[string]interface{})
			grayConfigRaw := make(map[string]interface{})
			if configRaw["grayConfig"] != nil {
				grayConfigRaw = configRaw["grayConfig"].(map[string]interface{})
			}
			if len(grayConfigRaw) > 0 {
				grayConfigMap["gray_rate"] = formatInt(grayConfigRaw["grayRate"])
				grayConfigMap["gray_sub_key"] = grayConfigRaw["graySubKey"]
				grayConfigMap["gray_target"] = grayConfigRaw["grayTarget"]

				grayConfigMaps = append(grayConfigMaps, grayConfigMap)
			}
			configMap["gray_config"] = grayConfigMaps
			rateLimitMaps := make([]map[string]interface{}, 0)
			rateLimitMap := make(map[string]interface{})
			ratelimitRaw := make(map[string]interface{})
			if configRaw["ratelimit"] != nil {
				ratelimitRaw = configRaw["ratelimit"].(map[string]interface{})
			}
			if len(ratelimitRaw) > 0 {
				rateLimitMap["interval"] = formatInt(ratelimitRaw["interval"])
				rateLimitMap["sub_key"] = ratelimitRaw["subKey"]
				rateLimitMap["target"] = ratelimitRaw["target"]
				rateLimitMap["threshold"] = formatInt(ratelimitRaw["threshold"])
				rateLimitMap["ttl"] = formatInt(ratelimitRaw["ttl"])

				statusMaps := make([]map[string]interface{}, 0)
				statusMap := make(map[string]interface{})
				statusRaw := make(map[string]interface{})
				if ratelimitRaw["status"] != nil {
					statusRaw = ratelimitRaw["status"].(map[string]interface{})
				}
				if len(statusRaw) > 0 {
					statusMap["code"] = formatInt(statusRaw["code"])
					statusMap["count"] = formatInt(statusRaw["count"])
					statusMap["ratio"] = formatInt(statusRaw["ratio"])

					statusMaps = append(statusMaps, statusMap)
				}
				rateLimitMap["status"] = statusMaps
				rateLimitMaps = append(rateLimitMaps, rateLimitMap)
			}
			configMap["rate_limit"] = rateLimitMaps
			remoteAddrRaw := make([]interface{}, 0)
			if configRaw["remoteAddr"] != nil {
				remoteAddrRaw = convertToInterfaceArray(configRaw["remoteAddr"])
			}

			configMap["remote_addr"] = remoteAddrRaw
			timeConfigMaps := make([]map[string]interface{}, 0)
			timeConfigMap := make(map[string]interface{})
			timeConfigRaw := make(map[string]interface{})
			if configRaw["timeConfig"] != nil {
				timeConfigRaw = configRaw["timeConfig"].(map[string]interface{})
			}
			if len(timeConfigRaw) > 0 {
				timeConfigMap["time_scope"] = timeConfigRaw["timeScope"]
				timeConfigMap["time_zone"] = formatInt(timeConfigRaw["timeZone"])

				timePeriodsRaw := timeConfigRaw["timePeriods"]
				timePeriodsMaps := make([]map[string]interface{}, 0)
				if timePeriodsRaw != nil {
					for _, timePeriodsChildRaw := range convertToInterfaceArray(timePeriodsRaw) {
						timePeriodsMap := make(map[string]interface{})
						timePeriodsChildRaw := timePeriodsChildRaw.(map[string]interface{})
						timePeriodsMap["end"] = formatInt(timePeriodsChildRaw["end"])
						timePeriodsMap["start"] = formatInt(timePeriodsChildRaw["start"])

						timePeriodsMaps = append(timePeriodsMaps, timePeriodsMap)
					}
				}
				timeConfigMap["time_periods"] = timePeriodsMaps
				weekTimePeriodsRaw := timeConfigRaw["weekTimePeriods"]
				weekTimePeriodsMaps := make([]map[string]interface{}, 0)
				if weekTimePeriodsRaw != nil {
					for _, weekTimePeriodsChildRaw := range convertToInterfaceArray(weekTimePeriodsRaw) {
						weekTimePeriodsMap := make(map[string]interface{})
						weekTimePeriodsChildRaw := weekTimePeriodsChildRaw.(map[string]interface{})
						weekTimePeriodsMap["day"] = weekTimePeriodsChildRaw["day"]

						dayPeriodsRaw := weekTimePeriodsChildRaw["dayPeriods"]
						dayPeriodsMaps := make([]map[string]interface{}, 0)
						if dayPeriodsRaw != nil {
							for _, dayPeriodsChildRaw := range convertToInterfaceArray(dayPeriodsRaw) {
								dayPeriodsMap := make(map[string]interface{})
								dayPeriodsChildRaw := dayPeriodsChildRaw.(map[string]interface{})
								dayPeriodsMap["end"] = formatInt(dayPeriodsChildRaw["end"])
								dayPeriodsMap["start"] = formatInt(dayPeriodsChildRaw["start"])

								dayPeriodsMaps = append(dayPeriodsMaps, dayPeriodsMap)
							}
						}
						weekTimePeriodsMap["day_periods"] = dayPeriodsMaps
						weekTimePeriodsMaps = append(weekTimePeriodsMaps, weekTimePeriodsMap)
					}
				}
				timeConfigMap["week_time_periods"] = weekTimePeriodsMaps
				timeConfigMaps = append(timeConfigMaps, timeConfigMap)
			}
			configMap["time_config"] = timeConfigMaps
			config1Raw := configRaw["config"]
			wafBaseConfigMaps := make([]map[string]interface{}, 0)
			if config1Raw != nil {
				for _, configChildRaw := range convertToInterfaceArray(config1Raw) {
					wafBaseConfigMap := make(map[string]interface{})
					configChildRaw := configChildRaw.(map[string]interface{})
					wafBaseConfigMap["rule_batch_operation_config"] = configChildRaw["ruleBatchOperationConfig"]
					wafBaseConfigMap["rule_type"] = configChildRaw["ruleType"]

					ruleDetailRaw := configChildRaw["ruleDetail"]
					ruleDetailMaps := make([]map[string]interface{}, 0)
					if ruleDetailRaw != nil {
						for _, ruleDetailChildRaw := range convertToInterfaceArray(ruleDetailRaw) {
							ruleDetailMap := make(map[string]interface{})
							ruleDetailChildRaw := ruleDetailChildRaw.(map[string]interface{})
							ruleDetailMap["rule_action"] = ruleDetailChildRaw["ruleAction"]
							if v := ruleDetailChildRaw["ruleId"]; v != nil {
								ruleDetailMap["rule_id"] = fmt.Sprint(v)
							}
							ruleDetailMap["rule_status"] = formatInt(ruleDetailChildRaw["ruleStatus"])

							ruleDetailMaps = append(ruleDetailMaps, ruleDetailMap)
						}
					}
					wafBaseConfigMap["rule_detail"] = ruleDetailMaps
					wafBaseConfigMaps = append(wafBaseConfigMaps, wafBaseConfigMap)
				}
			}
			configMap["waf_base_config"] = wafBaseConfigMaps
			configMaps = append(configMaps, configMap)
		}
	}
	mapping["config"] = configMaps

	return mapping, nil
}
