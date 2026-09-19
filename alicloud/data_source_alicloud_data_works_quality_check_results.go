package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudDataWorksQualityCheckResults() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudDataWorksQualityCheckResultsRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"project_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"start_date": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"end_date": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"entity_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"rule_id"},
			},
			"rule_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"entity_id"},
			},
			"page_size": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Default:  10,
			},
			"page_number": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Default:  1,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"results": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"qualityt_check_result_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"entity_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"rule_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"project_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"table_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"op": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"task_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"checker_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"checker_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"checker_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"begin_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"end_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"expect_value": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"upper_value": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"lower_value": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"warning_threshold": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"critical_threshold": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"trend": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"check_result": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"fixed_check": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_prediction": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"block_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"match_expression": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"actual_expression": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"where_condition": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"property": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"method_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"date_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"rule_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"template_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"template_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"comment": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"external_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"external_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudDataWorksQualityCheckResultsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	var action string
	request := make(map[string]interface{})
	request["ProjectName"] = d.Get("project_name")
	request["StartDate"] = d.Get("start_date")
	request["EndDate"] = d.Get("end_date")

	if v, ok := d.GetOk("entity_id"); ok && v.(string) != "" {
		action = "ListQualityResultsByEntity"
		request["EntityId"] = v
	} else if v, ok := d.GetOk("rule_id"); ok && v.(string) != "" {
		action = "ListQualityResultsByRule"
		request["RuleId"] = v
	} else {
		action = "ListQualityResultsByEntity"
	}

	pageSize := d.Get("page_size").(int)
	if pageSize <= 0 || pageSize > 20 {
		pageSize = 10
	}
	request["PageSize"] = pageSize
	request["PageNumber"] = d.Get("page_number").(int)

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}

	var objects []map[string]interface{}
	var response map[string]interface{}
	var err error
	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("dataworks-public", "2020-05-18", action, nil, request, true)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_data_works_quality_check_results", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.Data.RuleChecks", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Data.RuleChecks", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["Id"])]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}
		if len(result) < pageSize {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}

	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"id":                       fmt.Sprint(object["Id"]),
			"qualityt_check_result_id": fmt.Sprint(object["Id"]),
			"entity_id":                fmt.Sprint(object["EntityId"]),
			"rule_id":                  fmt.Sprint(object["RuleId"]),
			"project_name":             fmt.Sprint(object["ProjectName"]),
			"table_name":               fmt.Sprint(object["TableName"]),
			"op":                       fmt.Sprint(object["Op"]),
			"task_id":                  fmt.Sprint(object["TaskId"]),
			"checker_name":             fmt.Sprint(object["CheckerName"]),
			"checker_id":               fmt.Sprint(object["CheckerId"]),
			"checker_type":             fmt.Sprint(object["CheckerType"]),
			"begin_time":               fmt.Sprint(object["BeginTime"]),
			"end_time":                 fmt.Sprint(object["EndTime"]),
			"expect_value":             fmt.Sprint(object["ExpectValue"]),
			"upper_value":              fmt.Sprint(object["UpperValue"]),
			"lower_value":              fmt.Sprint(object["LowerValue"]),
			"warning_threshold":        fmt.Sprint(object["WarningThreshold"]),
			"critical_threshold":       fmt.Sprint(object["CriticalThreshold"]),
			"trend":                    fmt.Sprint(object["Trend"]),
			"check_result":             fmt.Sprint(object["CheckResult"]),
			"fixed_check":              fmt.Sprint(object["FixedCheck"]),
			"is_prediction":            fmt.Sprint(object["IsPrediction"]),
			"block_type":               fmt.Sprint(object["BlockType"]),
			"match_expression":         fmt.Sprint(object["MatchExpression"]),
			"actual_expression":        fmt.Sprint(object["ActualExpression"]),
			"where_condition":          fmt.Sprint(object["WhereCondition"]),
			"property":                 fmt.Sprint(object["Property"]),
			"method_name":              fmt.Sprint(object["MethodName"]),
			"date_type":                fmt.Sprint(object["DateType"]),
			"rule_name":                fmt.Sprint(object["RuleName"]),
			"template_name":            fmt.Sprint(object["TemplateName"]),
			"template_id":              fmt.Sprint(object["TemplateId"]),
			"comment":                  fmt.Sprint(object["Comment"]),
			"external_id":              fmt.Sprint(object["ExternalId"]),
			"external_type":            fmt.Sprint(object["ExternalType"]),
		}
		ids = append(ids, fmt.Sprint(mapping["id"]))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("results", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
