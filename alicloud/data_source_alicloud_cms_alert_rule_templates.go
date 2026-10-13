package alicloud

import (
	"fmt"
	"regexp"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudCmsAlertRuleTemplates() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudCmsAlertRuleTemplatesRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
			},
			"biz_source": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"CI", "CMS_ENT"}, false),
			},
			"enable_details": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"names": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"templates": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"alert_rule_template_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"biz_source": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"display_name_cn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"display_name_en": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description_cn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description_en": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"level": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"interval": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"labels": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"message_cn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"message_en": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"condition": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"alert_count": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"no_data_append_value": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"nodata_alert_level": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"case_list": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"condition": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"count_condition": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"express": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"first_param": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"level": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"oper": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"second_param": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"type": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"compare_list": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"aggregate": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"oper": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"threshold": {
													Type:     schema.TypeFloat,
													Computed: true,
												},
												"yoy_time_unit": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"yoy_time_value": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"value_level_list": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"level": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"value": {
																Type:     schema.TypeFloat,
																Computed: true,
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
						"datasource": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"instance_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"namespace": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"ds_list": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"project": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"region_id": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"store": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
								},
							},
						},
						"query": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"duration": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"expr": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"group_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"group_field_list": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"first_join": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"type": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"conditions": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"first_field": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"oper": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"second_field": {
																Type:     schema.TypeString,
																Computed: true,
															},
														},
													},
												},
											},
										},
									},
									"second_join": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"type": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"conditions": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"first_field": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"oper": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"second_field": {
																Type:     schema.TypeString,
																Computed: true,
															},
														},
													},
												},
											},
										},
									},
									"queries": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"duration": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"end": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"expr": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"start": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"time_unit": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"window": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudCmsAlertRuleTemplatesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "/alertRuleTemplates"
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

	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		nameRegex = r
	}

	query := make(map[string]*string)
	query["RegionId"] = StringPointer(client.RegionId)
	query["bizSource"] = StringPointer(d.Get("biz_source").(string))
	includeDetails := "false"
	if v, ok := d.GetOk("enable_details"); ok && v.(bool) {
		includeDetails = "true"
	}
	query["includeDetails"] = StringPointer(includeDetails)
	if v, ok := d.GetOk("ids"); ok && len(v.([]interface{})) == 1 {
		if vv, ok := v.([]interface{})[0].(string); ok && vv != "" {
			query["alertRuleTemplateId"] = StringPointer(vv)
		}
	}
	query["maxResults"] = StringPointer(fmt.Sprintf("%d", PageSizeLarge))

	var response map[string]interface{}
	var err error
	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RoaGet("Cms", "2024-03-30", action, query, nil, nil)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, query)
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_cms_alert_rule_templates", action, AlibabaCloudSdkGoERROR)
		}
		if fmt.Sprint(response["success"]) == "false" {
			return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
		}
		resp, gerr := jsonpath.Get("$.data[*]", response)
		if gerr != nil {
			return WrapErrorf(gerr, FailedGetAttributeMsg, action, "$.data[*]", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item, ok := v.(map[string]interface{})
			if !ok {
				continue
			}
			displayName := fmt.Sprint(item["displayNameEn"])
			if nameRegex != nil && !nameRegex.MatchString(displayName) {
				continue
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["alertRuleTemplateId"])]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}
		if nextToken, ok := response["nextToken"].(string); ok && nextToken != "" {
			query["nextToken"] = StringPointer(nextToken)
			continue
		}
		break
	}

	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		templateId := fmt.Sprint(object["alertRuleTemplateId"])
		mapping := map[string]interface{}{
			"id":                     templateId,
			"alert_rule_template_id": templateId,
			"biz_source":             object["bizSource"],
			"display_name_cn":        object["displayNameCn"],
			"display_name_en":        object["displayNameEn"],
			"description_cn":         object["descriptionCn"],
			"description_en":         object["descriptionEn"],
			"level":                  object["level"],
			"interval":               object["interval"],
			"labels":                 object["labels"],
			"message_cn":             object["messageCn"],
			"message_en":             object["messageEn"],
		}

		conditionMaps := make([]map[string]interface{}, 0)
		if conditionRaw, ok := object["condition"].(map[string]interface{}); ok && len(conditionRaw) > 0 {
			conditionMap := map[string]interface{}{
				"alert_count":          conditionRaw["alertCount"],
				"no_data_append_value": conditionRaw["noDataAppendValue"],
				"nodata_alert_level":   conditionRaw["nodataAlertLevel"],
				"type":                 conditionRaw["type"],
			}
			caseListMaps := make([]map[string]interface{}, 0)
			if caseListRaw := conditionRaw["caseList"]; caseListRaw != nil {
				for _, caseItemRaw := range convertToInterfaceArray(caseListRaw) {
					caseItem, ok := caseItemRaw.(map[string]interface{})
					if !ok {
						continue
					}
					caseListMaps = append(caseListMaps, map[string]interface{}{
						"condition":       caseItem["condition"],
						"count_condition": caseItem["countCondition"],
						"express":         caseItem["express"],
						"first_param":     caseItem["firstParam"],
						"level":           caseItem["level"],
						"oper":            caseItem["oper"],
						"second_param":    caseItem["secondParam"],
						"type":            caseItem["type"],
					})
				}
			}
			conditionMap["case_list"] = caseListMaps
			compareListMaps := make([]map[string]interface{}, 0)
			if compareListRaw := conditionRaw["compareList"]; compareListRaw != nil {
				for _, compareItemRaw := range convertToInterfaceArray(compareListRaw) {
					compareItem, ok := compareItemRaw.(map[string]interface{})
					if !ok {
						continue
					}
					compareMap := map[string]interface{}{
						"aggregate":      compareItem["aggregate"],
						"oper":           compareItem["oper"],
						"threshold":      compareItem["value"],
						"yoy_time_unit":  compareItem["yoyTimeUnit"],
						"yoy_time_value": compareItem["yoyTimeValue"],
					}
					valueLevelListMaps := make([]map[string]interface{}, 0)
					if valueLevelListRaw := compareItem["valueLevelList"]; valueLevelListRaw != nil {
						for _, vllItemRaw := range convertToInterfaceArray(valueLevelListRaw) {
							vllItem, ok := vllItemRaw.(map[string]interface{})
							if !ok {
								continue
							}
							valueLevelListMaps = append(valueLevelListMaps, map[string]interface{}{
								"level": vllItem["level"],
								"value": vllItem["value"],
							})
						}
					}
					compareMap["value_level_list"] = valueLevelListMaps
					compareListMaps = append(compareListMaps, compareMap)
				}
			}
			conditionMap["compare_list"] = compareListMaps
			conditionMaps = append(conditionMaps, conditionMap)
		}
		mapping["condition"] = conditionMaps

		datasourceMaps := make([]map[string]interface{}, 0)
		if dsRaw, ok := object["datasource"].(map[string]interface{}); ok && len(dsRaw) > 0 {
			dsMap := map[string]interface{}{
				"instance_id": dsRaw["instanceId"],
				"namespace":   dsRaw["namespace"],
				"type":        dsRaw["type"],
			}
			dsListMaps := make([]map[string]interface{}, 0)
			if dsListRaw := dsRaw["dsList"]; dsListRaw != nil {
				for _, dsItemRaw := range convertToInterfaceArray(dsListRaw) {
					dsItem, ok := dsItemRaw.(map[string]interface{})
					if !ok {
						continue
					}
					dsListMaps = append(dsListMaps, map[string]interface{}{
						"project":   dsItem["project"],
						"region_id": dsItem["regionId"],
						"store":     dsItem["store"],
					})
				}
			}
			dsMap["ds_list"] = dsListMaps
			datasourceMaps = append(datasourceMaps, dsMap)
		}
		mapping["datasource"] = datasourceMaps

		queryMaps := make([]map[string]interface{}, 0)
		if qRaw, ok := object["query"].(map[string]interface{}); ok && len(qRaw) > 0 {
			qMap := map[string]interface{}{
				"duration":   qRaw["duration"],
				"expr":       qRaw["expr"],
				"group_type": qRaw["groupType"],
				"type":       qRaw["type"],
			}
			groupFieldListRaw := make([]interface{}, 0)
			if gfl := qRaw["groupFieldList"]; gfl != nil {
				groupFieldListRaw = convertToInterfaceArray(gfl)
			}
			qMap["group_field_list"] = groupFieldListRaw
			qMap["first_join"] = buildJoinMaps(qRaw["firstJoin"])
			qMap["second_join"] = buildJoinMaps(qRaw["secondJoin"])
			queriesMaps := make([]map[string]interface{}, 0)
			if queriesRaw := qRaw["queries"]; queriesRaw != nil {
				for _, qItemRaw := range convertToInterfaceArray(queriesRaw) {
					qItem, ok := qItemRaw.(map[string]interface{})
					if !ok {
						continue
					}
					queriesMaps = append(queriesMaps, map[string]interface{}{
						"duration":  qItem["duration"],
						"end":       qItem["end"],
						"expr":      qItem["expr"],
						"start":     qItem["start"],
						"time_unit": qItem["timeUnit"],
						"window":    qItem["window"],
					})
				}
			}
			qMap["queries"] = queriesMaps
			queryMaps = append(queryMaps, qMap)
		}
		mapping["query"] = queryMaps

		ids = append(ids, templateId)
		names = append(names, object["displayNameEn"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}
	if err := d.Set("templates", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}

// buildJoinMaps flattens a firstJoin/secondJoin object into a single-element
// list of maps, mirroring the ListAlertRuleTemplates output structure.
func buildJoinMaps(raw interface{}) []map[string]interface{} {
	joinMaps := make([]map[string]interface{}, 0)
	joinRaw, ok := raw.(map[string]interface{})
	if !ok || len(joinRaw) == 0 {
		return joinMaps
	}
	joinMap := map[string]interface{}{
		"type": joinRaw["type"],
	}
	conditionsMaps := make([]map[string]interface{}, 0)
	if conditionsRaw := joinRaw["conditions"]; conditionsRaw != nil {
		for _, condItemRaw := range convertToInterfaceArray(conditionsRaw) {
			condItem, ok := condItemRaw.(map[string]interface{})
			if !ok {
				continue
			}
			conditionsMaps = append(conditionsMaps, map[string]interface{}{
				"first_field":  condItem["firstField"],
				"oper":         condItem["oper"],
				"second_field": condItem["secondField"],
			})
		}
	}
	joinMap["conditions"] = conditionsMaps
	joinMaps = append(joinMaps, joinMap)
	return joinMaps
}
