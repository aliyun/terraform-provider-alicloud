// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
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
		err = retry.Retry(d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
			response, err = client.RpcPost("waf-openapi", "2021-10-01", action, query, request, true)

			if err != nil {
				if NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
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
				v502, err502 := formatBool(configRaw["autoUpdate"])
				if err502 != nil {
					return err502
				}
				configMap["auto_update"] = v502
				configMap["cc_effect"] = configRaw["effect"]
				v504, err504 := formatInt(configRaw["ccStatus"])
				if err504 != nil {
					return err504
				}
				configMap["cc_status"] = v504
				configMap["cn_regions"] = configRaw["cnRegionList"]
				v506, err506 := formatInt(configRaw["grayStatus"])
				if err506 != nil {
					return err506
				}
				configMap["gray_status"] = v506
				v507, err507 := formatInt(configRaw["mode"])
				if err507 != nil {
					return err507
				}
				configMap["mode"] = v507
				configMap["protocol"] = configRaw["protocol"]
				configMap["rule_action"] = configRaw["action"]
				v510, err510 := formatInt(configRaw["threshold"])
				if err510 != nil {
					return err510
				}
				configMap["throttle_threhold"] = v510
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
						v524, err524 := formatInt(accountIdentifiersChildRaw["priority"])
						if err524 != nil {
							return err524
						}
						accountIdentifiersMap["priority"] = v524
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
					v577, err577 := formatInt(grayConfigRaw["grayRate"])
					if err577 != nil {
						return err577
					}
					grayConfigMap["gray_rate"] = v577
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
					v591, err591 := formatInt(ratelimitRaw["interval"])
					if err591 != nil {
						return err591
					}
					rateLimitMap["interval"] = v591
					rateLimitMap["sub_key"] = ratelimitRaw["subKey"]
					rateLimitMap["target"] = ratelimitRaw["target"]
					v594, err594 := formatInt(ratelimitRaw["threshold"])
					if err594 != nil {
						return err594
					}
					rateLimitMap["threshold"] = v594
					v595, err595 := formatInt(ratelimitRaw["ttl"])
					if err595 != nil {
						return err595
					}
					rateLimitMap["ttl"] = v595

					statusMaps := make([]map[string]interface{}, 0)
					statusMap := make(map[string]interface{})
					statusRaw := make(map[string]interface{})
					if ratelimitRaw["status"] != nil {
						statusRaw = ratelimitRaw["status"].(map[string]interface{})
					}
					if len(statusRaw) > 0 {
						v604, err604 := formatInt(statusRaw["code"])
						if err604 != nil {
							return err604
						}
						statusMap["code"] = v604
						v605, err605 := formatInt(statusRaw["count"])
						if err605 != nil {
							return err605
						}
						statusMap["count"] = v605
						v606, err606 := formatInt(statusRaw["ratio"])
						if err606 != nil {
							return err606
						}
						statusMap["ratio"] = v606

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
					v628, err628 := formatInt(timeConfigRaw["timeZone"])
					if err628 != nil {
						return err628
					}
					timeConfigMap["time_zone"] = v628

					timePeriodsRaw := timeConfigRaw["timePeriods"]
					timePeriodsMaps := make([]map[string]interface{}, 0)
					if timePeriodsRaw != nil {
						for _, timePeriodsChildRaw := range convertToInterfaceArray(timePeriodsRaw) {
							timePeriodsMap := make(map[string]interface{})
							timePeriodsChildRaw := timePeriodsChildRaw.(map[string]interface{})
							v636, err636 := formatInt(timePeriodsChildRaw["end"])
							if err636 != nil {
								return err636
							}
							timePeriodsMap["end"] = v636
							v637, err637 := formatInt(timePeriodsChildRaw["start"])
							if err637 != nil {
								return err637
							}
							timePeriodsMap["start"] = v637

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
									v657, err657 := formatInt(dayPeriodsChildRaw["end"])
									if err657 != nil {
										return err657
									}
									dayPeriodsMap["end"] = v657
									v658, err658 := formatInt(dayPeriodsChildRaw["start"])
									if err658 != nil {
										return err658
									}
									dayPeriodsMap["start"] = v658

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
								v690, err690 := formatInt(ruleDetailChildRaw["ruleStatus"])
								if err690 != nil {
									return err690
								}
								ruleDetailMap["rule_status"] = v690

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
			v767, err767 := formatBool(configRaw["autoUpdate"])
			if err767 != nil {
				return nil, err767
			}
			configMap["auto_update"] = v767
			configMap["cc_effect"] = configRaw["effect"]
			v769, err769 := formatInt(configRaw["ccStatus"])
			if err769 != nil {
				return nil, err769
			}
			configMap["cc_status"] = v769
			configMap["cn_regions"] = configRaw["cnRegionList"]
			v771, err771 := formatInt(configRaw["grayStatus"])
			if err771 != nil {
				return nil, err771
			}
			configMap["gray_status"] = v771
			v772, err772 := formatInt(configRaw["mode"])
			if err772 != nil {
				return nil, err772
			}
			configMap["mode"] = v772
			configMap["protocol"] = configRaw["protocol"]
			configMap["rule_action"] = configRaw["action"]
			v775, err775 := formatInt(configRaw["threshold"])
			if err775 != nil {
				return nil, err775
			}
			configMap["throttle_threhold"] = v775
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
					v789, err789 := formatInt(accountIdentifiersChildRaw["priority"])
					if err789 != nil {
						return nil, err789
					}
					accountIdentifiersMap["priority"] = v789
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
				v842, err842 := formatInt(grayConfigRaw["grayRate"])
				if err842 != nil {
					return nil, err842
				}
				grayConfigMap["gray_rate"] = v842
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
				v856, err856 := formatInt(ratelimitRaw["interval"])
				if err856 != nil {
					return nil, err856
				}
				rateLimitMap["interval"] = v856
				rateLimitMap["sub_key"] = ratelimitRaw["subKey"]
				rateLimitMap["target"] = ratelimitRaw["target"]
				v859, err859 := formatInt(ratelimitRaw["threshold"])
				if err859 != nil {
					return nil, err859
				}
				rateLimitMap["threshold"] = v859
				v860, err860 := formatInt(ratelimitRaw["ttl"])
				if err860 != nil {
					return nil, err860
				}
				rateLimitMap["ttl"] = v860

				statusMaps := make([]map[string]interface{}, 0)
				statusMap := make(map[string]interface{})
				statusRaw := make(map[string]interface{})
				if ratelimitRaw["status"] != nil {
					statusRaw = ratelimitRaw["status"].(map[string]interface{})
				}
				if len(statusRaw) > 0 {
					v869, err869 := formatInt(statusRaw["code"])
					if err869 != nil {
						return nil, err869
					}
					statusMap["code"] = v869
					v870, err870 := formatInt(statusRaw["count"])
					if err870 != nil {
						return nil, err870
					}
					statusMap["count"] = v870
					v871, err871 := formatInt(statusRaw["ratio"])
					if err871 != nil {
						return nil, err871
					}
					statusMap["ratio"] = v871

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
				v893, err893 := formatInt(timeConfigRaw["timeZone"])
				if err893 != nil {
					return nil, err893
				}
				timeConfigMap["time_zone"] = v893

				timePeriodsRaw := timeConfigRaw["timePeriods"]
				timePeriodsMaps := make([]map[string]interface{}, 0)
				if timePeriodsRaw != nil {
					for _, timePeriodsChildRaw := range convertToInterfaceArray(timePeriodsRaw) {
						timePeriodsMap := make(map[string]interface{})
						timePeriodsChildRaw := timePeriodsChildRaw.(map[string]interface{})
						v901, err901 := formatInt(timePeriodsChildRaw["end"])
						if err901 != nil {
							return nil, err901
						}
						timePeriodsMap["end"] = v901
						v902, err902 := formatInt(timePeriodsChildRaw["start"])
						if err902 != nil {
							return nil, err902
						}
						timePeriodsMap["start"] = v902

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
								v922, err922 := formatInt(dayPeriodsChildRaw["end"])
								if err922 != nil {
									return nil, err922
								}
								dayPeriodsMap["end"] = v922
								v923, err923 := formatInt(dayPeriodsChildRaw["start"])
								if err923 != nil {
									return nil, err923
								}
								dayPeriodsMap["start"] = v923

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
							v955, err955 := formatInt(ruleDetailChildRaw["ruleStatus"])
							if err955 != nil {
								return nil, err955
							}
							ruleDetailMap["rule_status"] = v955

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
