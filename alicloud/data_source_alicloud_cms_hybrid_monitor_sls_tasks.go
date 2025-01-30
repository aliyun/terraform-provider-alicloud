package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudCmsHybridMonitorSlsTasks() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudCmsHybridMonitorSlsTasksRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"keyword": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"namespace": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"page_number": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"page_size": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  50,
			},
			"tasks": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"attach_labels": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
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
						"collect_interval": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"collect_target_endpoint": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"collect_target_path": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"collect_target_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"collect_timout": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"extra_info": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"hybrid_monitor_sls_task_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instances": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"log_file_path": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"log_process": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"log_sample": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"log_split": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"match_express": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"function": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"name": {
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
						"match_express_relation": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"namespace": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"network_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sls_process": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sls_process_config": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"express": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"alias": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"express": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"filter": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"filters": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"operator": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"sls_key_name": {
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
												"relation": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"group_by": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"alias": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"sls_key_name": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"statistics": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"alias": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"function": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"parameter_one": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"parameter_two": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"sls_key_name": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
								},
							},
						},
						"task_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"task_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"upload_region": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"yarm_config": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudCmsHybridMonitorSlsTasksRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeHybridMonitorTaskList"
	request := make(map[string]interface{})
	if v, ok := d.GetOk("keyword"); ok {
		request["Keyword"] = v
	}
	if v, ok := d.GetOk("namespace"); ok {
		request["Namespace"] = v
	}
	request["TaskType"] = "aliyun_sls"
	setPagingRequest(d, request, PageSizeLarge)
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
	var err error
	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("Cms", "2019-01-01", action, nil, request, false)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_cms_hybrid_monitor_sls_tasks", action, AlibabaCloudSdkGoERROR)
		}
		if fmt.Sprint(response["Success"]) == "false" {
			return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
		}
		resp, err := jsonpath.Get("$.TaskList", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.TaskList", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["TaskId"])]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}
		if len(result) < request["PageSize"].(int) {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"collect_interval":           formatInt(object["CollectInterval"]),
			"collect_target_endpoint":    object["CollectTargetEndpoint"],
			"collect_target_path":        object["CollectTargetPath"],
			"collect_target_type":        object["CollectTargetType"],
			"collect_timout":             formatInt(object["CollectTimout"]),
			"create_time":                object["CreateTime"],
			"description":                object["Description"],
			"extra_info":                 object["ExtraInfo"],
			"group_id":                   object["GroupId"],
			"id":                         fmt.Sprint(object["TaskId"]),
			"hybrid_monitor_sls_task_id": fmt.Sprint(object["TaskId"]),
			"instances":                  object["Instances"],
			"log_file_path":              object["LogFilePath"],
			"log_process":                object["LogProcess"],
			"log_sample":                 object["LogSample"],
			"log_split":                  object["LogSplit"],
			"match_express_relation":     object["MatchExpressRelation"],
			"namespace":                  object["Namespace"],
			"network_type":               object["NetworkType"],
			"sls_process":                object["SLSProcess"],
			"task_name":                  object["TaskName"],
			"task_type":                  object["TaskType"],
			"upload_region":              object["UploadRegion"],
			"yarm_config":                object["YARMConfig"],
		}

		attachLabels := make([]map[string]interface{}, 0)
		if attachLabelsList, ok := object["AttachLabels"].([]interface{}); ok {
			for _, v := range attachLabelsList {
				if m1, ok := v.(map[string]interface{}); ok {
					temp1 := map[string]interface{}{
						"name":  m1["Name"],
						"value": m1["Value"],
					}
					attachLabels = append(attachLabels, temp1)
				}
			}
		}
		mapping["attach_labels"] = attachLabels

		matchExpress := make([]map[string]interface{}, 0)
		if matchExpressList, ok := object["MatchExpress"].([]interface{}); ok {
			for _, v := range matchExpressList {
				if m1, ok := v.(map[string]interface{}); ok {
					temp1 := map[string]interface{}{
						"function": m1["Function"],
						"name":     m1["Name"],
						"value":    m1["Value"],
					}
					matchExpress = append(matchExpress, temp1)
				}
			}
		}
		mapping["match_express"] = matchExpress

		slsProcessConfigSli := make([]map[string]interface{}, 0)
		if len(object["SLSProcessConfig"].(map[string]interface{})) > 0 {
			slsProcessConfig := object["SLSProcessConfig"]
			slsProcessConfigMap := make(map[string]interface{})

			expressSli := make([]map[string]interface{}, 0)
			if len(slsProcessConfig.(map[string]interface{})["Express"].([]interface{})) > 0 {
				for _, express := range slsProcessConfig.(map[string]interface{})["Express"].([]interface{}) {
					expressMap := make(map[string]interface{})
					expressMap["alias"] = express.(map[string]interface{})["Alias"]
					expressMap["express"] = express.(map[string]interface{})["Express"]
					expressSli = append(expressSli, expressMap)
				}
			}
			slsProcessConfigMap["express"] = expressSli

			filterSli := make([]map[string]interface{}, 0)
			if len(slsProcessConfig.(map[string]interface{})["Filter"].(map[string]interface{})) > 0 {
				filter := slsProcessConfig.(map[string]interface{})["Filter"]
				filterMap := make(map[string]interface{})

				filtersSli := make([]map[string]interface{}, 0)
				if len(filter.(map[string]interface{})["Filters"].([]interface{})) > 0 {
					for _, filters := range filter.(map[string]interface{})["Filters"].([]interface{}) {
						filtersMap := make(map[string]interface{})
						filtersMap["operator"] = filters.(map[string]interface{})["Operator"]
						filtersMap["sls_key_name"] = filters.(map[string]interface{})["SLSKeyName"]
						filtersMap["value"] = filters.(map[string]interface{})["Value"]
						filtersSli = append(filtersSli, filtersMap)
					}
				}
				filterMap["filters"] = filtersSli
				filterMap["relation"] = filter.(map[string]interface{})["Relation"]
				filterSli = append(filterSli, filterMap)
			}
			slsProcessConfigMap["filter"] = filterSli

			groupBySli := make([]map[string]interface{}, 0)
			if len(slsProcessConfig.(map[string]interface{})["GroupBy"].([]interface{})) > 0 {
				for _, groupBy := range slsProcessConfig.(map[string]interface{})["GroupBy"].([]interface{}) {
					groupByMap := make(map[string]interface{})
					groupByMap["alias"] = groupBy.(map[string]interface{})["Alias"]
					groupByMap["sls_key_name"] = groupBy.(map[string]interface{})["SLSKeyName"]
					groupBySli = append(groupBySli, groupByMap)
				}
			}
			slsProcessConfigMap["group_by"] = groupBySli

			statisticsSli := make([]map[string]interface{}, 0)
			if len(slsProcessConfig.(map[string]interface{})["Statistics"].([]interface{})) > 0 {
				for _, statistics := range slsProcessConfig.(map[string]interface{})["Statistics"].([]interface{}) {
					statisticsMap := make(map[string]interface{})
					statisticsMap["alias"] = statistics.(map[string]interface{})["Alias"]
					statisticsMap["function"] = statistics.(map[string]interface{})["Function"]
					statisticsMap["parameter_one"] = statistics.(map[string]interface{})["Parameter1"]
					statisticsMap["parameter_two"] = statistics.(map[string]interface{})["Parameter2"]
					statisticsMap["sls_key_name"] = statistics.(map[string]interface{})["SLSKeyName"]
					statisticsSli = append(statisticsSli, statisticsMap)
				}
			}
			slsProcessConfigMap["statistics"] = statisticsSli
			slsProcessConfigSli = append(slsProcessConfigSli, slsProcessConfigMap)
		}
		mapping["sls_process_config"] = slsProcessConfigSli
		ids = append(ids, fmt.Sprint(mapping["id"]))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("tasks", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
