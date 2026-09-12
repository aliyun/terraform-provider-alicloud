// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"strconv"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAliCloudCmsAlertMetrics() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudCmsAlertMetricsRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"group": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"include_details": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"metrics": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"alert_metric_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"group": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"alert_level": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"alert_message_cn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"alert_message_en": {
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
						"display_statement_cn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"display_statement_en": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"duration": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"unit_cn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"unit_en": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ext_info": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"expr_template": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"tpl": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"filters": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"dim": {
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
									"hidden": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"opt": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"supported_opts": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"display_name_cn": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"display_name_en": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"value": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
								},
							},
						},
						"params": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"max_width": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"min_width": {
										Type:     schema.TypeInt,
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
									"value": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"values": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"label_cn": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"label_en": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"value": {
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
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceAliCloudCmsAlertMetricsRead(d *schema.ResourceData, meta interface{}) error {
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
	var query map[string]*string
	// ListAlertMetrics
	action := fmt.Sprintf("/alertMetrics")
	var err error
	request = make(map[string]interface{})
	query = make(map[string]*string)
	query["RegionId"] = StringPointer(client.RegionId)
	if v := d.Get("group"); !IsNil(v) {
		query["group"] = StringPointer(v.(string))
	}
	if v := d.Get("include_details"); !IsNil(v) {
		query["includeDetails"] = StringPointer(strconv.FormatBool(v.(bool)))
	}
	// ListAlertMetrics returns the full set when maxResults is omitted, and
	// following its nextToken cursor currently fails with HTTP 500 on the
	// second page. Fetch the full set in a single request without pagination.
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutRead), func() *resource.RetryError {
		response, err = client.RoaGet("Cms", "2024-03-30", action, query, nil, nil)

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

	resp, _ := jsonpath.Get("$.data[*]", response)

	result, _ := resp.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		if len(idsMap) > 0 {
			if _, ok := idsMap[fmt.Sprint(item["alertMetricId"])]; !ok {
				continue
			}
		}
		objects = append(objects, item)
	}

	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, objectRaw := range objects {
		mapping := map[string]interface{}{}

		mapping["id"] = objectRaw["alertMetricId"]

		mapping["alert_level"] = objectRaw["alertLevel"]
		mapping["alert_message_cn"] = objectRaw["alertMessageCn"]
		mapping["alert_message_en"] = objectRaw["alertMessageEn"]
		mapping["alert_metric_id"] = objectRaw["alertMetricId"]
		mapping["display_name_cn"] = objectRaw["displayNameCn"]
		mapping["display_name_en"] = objectRaw["displayNameEn"]
		mapping["display_statement_cn"] = objectRaw["displayStatementCn"]
		mapping["display_statement_en"] = objectRaw["displayStatementEn"]
		mapping["duration"] = objectRaw["duration"]
		mapping["ext_info"] = objectRaw["extInfo"]
		mapping["group"] = objectRaw["group"]
		mapping["unit_cn"] = objectRaw["unitCn"]
		mapping["unit_en"] = objectRaw["unitEn"]

		exprTemplateMaps := make([]map[string]interface{}, 0)
		exprTemplateMap := make(map[string]interface{})

		exprTemplateRaw := make(map[string]interface{})
		if objectRaw["exprTemplate"] != nil {
			exprTemplateRaw = objectRaw["exprTemplate"].(map[string]interface{})
		}
		if len(exprTemplateRaw) > 0 {
			exprTemplateMap["tpl"] = exprTemplateRaw["tpl"]
			exprTemplateMap["type"] = exprTemplateRaw["type"]

			exprTemplateMaps = append(exprTemplateMaps, exprTemplateMap)
		}
		mapping["expr_template"] = exprTemplateMaps
		filtersMaps := make([]map[string]interface{}, 0)
		filtersRaw := objectRaw["filters"]
		if filtersRaw != nil {
			for _, filtersChildRaw := range convertToInterfaceArray(filtersRaw) {
				filtersMap := make(map[string]interface{})

				filtersChildRaw := filtersChildRaw.(map[string]interface{})
				filtersMap["dim"] = filtersChildRaw["dim"]
				filtersMap["display_name_cn"] = filtersChildRaw["displayNameCn"]
				filtersMap["display_name_en"] = filtersChildRaw["displayNameEn"]
				filtersMap["hidden"] = filtersChildRaw["hidden"]
				filtersMap["opt"] = filtersChildRaw["opt"]

				supportedOptsRaw := filtersChildRaw["supportedOpts"]
				supportedOptsMaps := make([]map[string]interface{}, 0)
				if supportedOptsRaw != nil {
					for _, supportedOptsChildRaw := range convertToInterfaceArray(supportedOptsRaw) {
						supportedOptsMap := make(map[string]interface{})

						supportedOptsChildRaw := supportedOptsChildRaw.(map[string]interface{})
						supportedOptsMap["display_name_cn"] = supportedOptsChildRaw["displayNameCn"]
						supportedOptsMap["display_name_en"] = supportedOptsChildRaw["displayNameEn"]
						supportedOptsMap["value"] = supportedOptsChildRaw["value"]

						supportedOptsMaps = append(supportedOptsMaps, supportedOptsMap)
					}
				}
				filtersMap["supported_opts"] = supportedOptsMaps
				filtersMaps = append(filtersMaps, filtersMap)
			}
		}
		mapping["filters"] = filtersMaps
		paramsMaps := make([]map[string]interface{}, 0)
		paramsRaw := objectRaw["params"]
		if paramsRaw != nil {
			for _, paramsChildRaw := range convertToInterfaceArray(paramsRaw) {
				paramsMap := make(map[string]interface{})

				paramsChildRaw := paramsChildRaw.(map[string]interface{})
				paramsMap["max_width"] = paramsChildRaw["maxWidth"]
				paramsMap["min_width"] = paramsChildRaw["minWidth"]
				paramsMap["name"] = paramsChildRaw["name"]
				paramsMap["type"] = paramsChildRaw["type"]
				paramsMap["value"] = paramsChildRaw["value"]

				valuesRaw := paramsChildRaw["values"]
				valuesMaps := make([]map[string]interface{}, 0)
				if valuesRaw != nil {
					for _, valuesChildRaw := range convertToInterfaceArray(valuesRaw) {
						valuesMap := make(map[string]interface{})

						valuesChildRaw := valuesChildRaw.(map[string]interface{})
						valuesMap["label_cn"] = valuesChildRaw["labelCn"]
						valuesMap["label_en"] = valuesChildRaw["labelEn"]
						valuesMap["value"] = valuesChildRaw["value"]

						valuesMaps = append(valuesMaps, valuesMap)
					}
				}
				paramsMap["values"] = valuesMaps
				paramsMaps = append(paramsMaps, paramsMap)
			}
		}
		mapping["params"] = paramsMaps

		ids = append(ids, fmt.Sprint(mapping["id"]))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("metrics", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
