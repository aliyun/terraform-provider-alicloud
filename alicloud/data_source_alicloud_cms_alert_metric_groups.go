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

func dataSourceAliCloudCmsAlertMetricGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudCmsAlertMetricGroupsRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"datasource_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"include_details": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"alert_metric_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"datasource_types": {
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
						"order_index": {
							Type:     schema.TypeInt,
							Computed: true,
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
									"label_disabled": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"dim_disabled": {
										Type:     schema.TypeBool,
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
									"placeholder_cn": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"placeholder_en": {
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
		},
	}
}

func dataSourceAliCloudCmsAlertMetricGroupsRead(d *schema.ResourceData, meta interface{}) error {
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

	var response map[string]interface{}
	var query map[string]*string
	// ListAlertMetricGroups
	action := "/alertMetricGroups"
	var err error
	query = make(map[string]*string)
	if v, ok := d.GetOk("datasource_type"); ok {
		query["datasourceType"] = StringPointer(v.(string))
	}
	if v, ok := d.GetOkExists("include_details"); ok {
		query["includeDetails"] = StringPointer(fmt.Sprint(v.(bool)))
	}
	query["maxResults"] = StringPointer(strconv.Itoa(PageSizeLarge))

	for {
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
			addDebug(action, response, nil)
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
				if _, ok := idsMap[fmt.Sprint(item["alertMetricGroupId"])]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}

		if nextToken, ok := response["nextToken"].(string); ok && nextToken != "" {
			query["nextToken"] = StringPointer(nextToken)
		} else {
			break
		}
	}

	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, objectRaw := range objects {
		mapping := map[string]interface{}{}

		mapping["id"] = objectRaw["alertMetricGroupId"]
		mapping["alert_metric_group_id"] = objectRaw["alertMetricGroupId"]
		mapping["datasource_types"] = objectRaw["datasourceTypes"]
		mapping["display_name_cn"] = objectRaw["displayNameCn"]
		mapping["display_name_en"] = objectRaw["displayNameEn"]
		mapping["description_cn"] = objectRaw["descriptionCn"]
		mapping["description_en"] = objectRaw["descriptionEn"]
		mapping["order_index"] = objectRaw["orderIndex"]

		filtersMaps := make([]map[string]interface{}, 0)
		if objectRaw["filters"] != nil {
			for _, filtersChildRaw := range convertToInterfaceArray(objectRaw["filters"]) {
				filtersMap := make(map[string]interface{})
				filtersChildRaw := filtersChildRaw.(map[string]interface{})
				filtersMap["dim"] = filtersChildRaw["dim"]
				filtersMap["display_name_cn"] = filtersChildRaw["displayNameCn"]
				filtersMap["display_name_en"] = filtersChildRaw["displayNameEn"]
				filtersMap["hidden"] = filtersChildRaw["hidden"]
				filtersMap["opt"] = filtersChildRaw["opt"]
				filtersMap["label_disabled"] = filtersChildRaw["labelDisabled"]
				filtersMap["dim_disabled"] = filtersChildRaw["dimDisabled"]

				supportedOptsMaps := make([]map[string]interface{}, 0)
				if filtersChildRaw["supportedOpts"] != nil {
					for _, supportedOptsChildRaw := range convertToInterfaceArray(filtersChildRaw["supportedOpts"]) {
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
		if objectRaw["params"] != nil {
			for _, paramsChildRaw := range convertToInterfaceArray(objectRaw["params"]) {
				paramsMap := make(map[string]interface{})
				paramsChildRaw := paramsChildRaw.(map[string]interface{})
				paramsMap["max_width"] = paramsChildRaw["maxWidth"]
				paramsMap["min_width"] = paramsChildRaw["minWidth"]
				paramsMap["name"] = paramsChildRaw["name"]
				paramsMap["placeholder_cn"] = paramsChildRaw["placeholderCn"]
				paramsMap["placeholder_en"] = paramsChildRaw["placeholderEn"]
				paramsMap["type"] = paramsChildRaw["type"]
				paramsMap["value"] = paramsChildRaw["value"]

				valuesMaps := make([]map[string]interface{}, 0)
				if paramsChildRaw["values"] != nil {
					for _, valuesChildRaw := range convertToInterfaceArray(paramsChildRaw["values"]) {
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

	if err := d.Set("groups", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
