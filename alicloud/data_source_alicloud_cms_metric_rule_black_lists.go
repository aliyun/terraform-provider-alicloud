package alicloud

import (
	"encoding/json"
	"fmt"
	"regexp"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudCmsMetricRuleBlackLists() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudCmsMetricRuleBlackListsRead,
		Schema: map[string]*schema.Schema{
			"category": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"metric_rule_black_list_id": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"namespace": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"order": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeInt,
			},
			"ids": {
				Optional: true,
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"name_regex": {
				Optional:     true,
				ForceNew:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.ValidateRegexp,
			},
			"names": {
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"output_file": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"page_number": {
				Optional: true,
				Type:     schema.TypeInt,
			},
			"page_size": {
				Optional: true,
				Type:     schema.TypeInt,
				Default:  10,
			},
			"lists": {
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"category": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"create_time": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"effective_time": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"enable_end_time": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"enable_start_time": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"instances": {
							Computed: true,
							Type:     schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"is_enable": {
							Computed: true,
							Type:     schema.TypeBool,
						},
						"metric_rule_black_list_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"metric_rule_black_list_name": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"metrics": {
							Computed: true,
							Type:     schema.TypeList,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"metric_name": {
										Computed: true,
										Type:     schema.TypeString,
									},
									"resource": {
										Computed: true,
										Type:     schema.TypeString,
									},
								},
							},
						},
						"namespace": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"scope_type": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"scope_value": {
							Computed: true,
							Type:     schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudCmsMetricRuleBlackListsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := make(map[string]interface{})

	if v, ok := d.GetOk("category"); ok {
		request["Category"] = v
	}
	if v, ok := d.GetOk("metric_rule_black_list_id"); ok {
		request["Ids"] = v
	}
	if v, ok := d.GetOk("namespace"); ok {
		request["Namespace"] = v
	}
	if v, ok := d.GetOk("order"); ok {
		request["Order"] = v
	}
	if v, ok := d.GetOk("page_size"); ok && v.(int) > 0 {
		request["PageSize"] = v.(int)
	} else {
		request["PageSize"] = PageSizeSmall
	}

	if v, ok := d.GetOk("page_number"); ok && v.(int) > 0 {
		request["PageNumber"] = v.(int)
	} else {
		request["PageNumber"] = 1
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

	var metricRuleBlackListNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		metricRuleBlackListNameRegex = r
	}

	var err error
	var objects []interface{}
	var response map[string]interface{}

	for {
		action := "DescribeMetricRuleBlackList"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			resp, err := client.RpcPost("Cms", "2019-01-01", action, nil, request, false)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			response = resp
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_cms_metric_rule_black_lists", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.DescribeMetricRuleBlackList", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.DescribeMetricRuleBlackList", response)
		}
		result, _ := resp.([]interface{})
		if isPagingRequest(d) {
			objects = result
			break
		}
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["Id"])]; !ok {
					continue
				}
			}

			if metricRuleBlackListNameRegex != nil && !metricRuleBlackListNameRegex.MatchString(fmt.Sprint(item["Name"])) {
				continue
			}
			objects = append(objects, item)
		}
		if len(result) < request["PageSize"].(int) {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}

	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, v := range objects {
		object := v.(map[string]interface{})
		instances := object["Instances"].([]interface{})
		instancesSet := make([]interface{}, 0)
		for _, valueOfInstance := range instances {
			jsonMarshalResult, err := json.Marshal(valueOfInstance)
			if err != nil {
				return WrapError(err)
			}
			instancesSet = append(instancesSet, string(jsonMarshalResult))
		}
		mapping := map[string]interface{}{
			"id":                          fmt.Sprint(object["Id"]),
			"category":                    object["Category"],
			"create_time":                 object["CreateTime"],
			"effective_time":              object["EffectiveTime"],
			"enable_end_time":             object["EnableEndTime"],
			"enable_start_time":           object["EnableStartTime"],
			"instances":                   instancesSet,
			"is_enable":                   object["IsEnable"],
			"metric_rule_black_list_id":   object["Id"],
			"metric_rule_black_list_name": object["Name"],
			"namespace":                   object["Namespace"],
			"scope_type":                  object["ScopeType"],
		}

		scopeValue, _ := jsonpath.Get("$.ScopeValue", object)
		mapping["scope_value"] = scopeValue

		metricsMaps := make([]map[string]interface{}, 0)
		metricsRaw := object["Metrics"]
		for _, value0 := range metricsRaw.([]interface{}) {
			metrics := value0.(map[string]interface{})
			metricsMap := make(map[string]interface{})
			metricsMap["metric_name"] = metrics["MetricName"]
			metricsMap["resource"] = metrics["Resource"]
			metricsMaps = append(metricsMaps, metricsMap)
		}
		mapping["metrics"] = metricsMaps

		ids = append(ids, fmt.Sprint(object["Id"]))
		names = append(names, object["Name"])

		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("lists", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
