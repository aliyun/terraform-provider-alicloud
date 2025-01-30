package alicloud

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudCmsEventRules() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudCmsEventRulesRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"name_prefix": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
				ForceNew:     true,
			},
			"page_number": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"page_size": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  PageSizeLarge,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"DISABLED", "ENABLED"}, false),
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
			"rules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"event_rule_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"event_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"silence_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"event_pattern": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"event_type_list": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"level_list": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"name_list": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"product": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"sql_filter": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"keyword_filter": {
										Type:     schema.TypeSet,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"key_words": {
													Type:     schema.TypeList,
													Computed: true,
													Elem:     &schema.Schema{Type: schema.TypeString},
												},
												"relation": {
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

func dataSourceAlicloudCmsEventRulesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeEventRuleList"
	request := make(map[string]interface{})
	setPagingRequest(d, request, PageSizeLarge)
	if v, ok := d.GetOk("name_prefix"); ok {
		request["NamePrefix"] = v
	}
	var objects []map[string]interface{}
	var eventRuleNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		eventRuleNameRegex = r
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
	status, statusOk := d.GetOk("status")
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_cms_event_rules", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.EventRules.EventRule", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.EventRules.EventRule", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if eventRuleNameRegex != nil && !eventRuleNameRegex.MatchString(fmt.Sprint(item["Name"])) {
				continue
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["Name"])]; !ok {
					continue
				}
			}
			if statusOk && status.(string) != "" && status.(string) != item["State"].(string) {
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
	for _, object := range objects {
		silenceTime, _ := strconv.Atoi(fmt.Sprint(object["SilenceTime"]))
		mapping := map[string]interface{}{
			"id":              fmt.Sprint(object["Name"]),
			"event_rule_name": fmt.Sprint(object["Name"]),
			"description":     object["Description"],
			"event_type":      object["EventType"],
			"group_id":        object["GroupId"],
			"silence_time":    silenceTime,
			"status":          object["State"],
		}
		if v, ok := object["EventPattern"]; ok {
			eventPatternMaps := make([]map[string]interface{}, 0)
			eventPattern := v.(map[string]interface{})
			eventPatternMap := make(map[string]interface{})
			eventPatternMap["event_type_list"] = eventPattern["EventTypeList"].(map[string]interface{})["EventTypeList"]
			eventPatternMap["level_list"] = eventPattern["LevelList"].(map[string]interface{})["LevelList"]
			eventPatternMap["name_list"] = eventPattern["NameList"].(map[string]interface{})["NameList"]
			eventPatternMap["product"] = eventPattern["Product"]
			eventPatternMap["sql_filter"] = eventPattern["SQLFilter"]
			if keywordFilter, ok := object["KeywordFilter"]; ok {
				keywordFilterMaps := make([]map[string]interface{}, 0)
				keywordFilterObject := keywordFilter.(map[string]interface{})
				keywordFilterMap := make(map[string]interface{})
				keywordFilterMap["key_words"] = keywordFilterObject["Keywords"]
				keywordFilterMap["relation"] = keywordFilterObject["Relation"]
				keywordFilterMaps = append(keywordFilterMaps, keywordFilterMap)
				eventPatternMap["keyword_filter"] = keywordFilterMaps
			}
			eventPatternMaps = append(eventPatternMaps, eventPatternMap)
			mapping["event_pattern"] = eventPatternMaps
		}
		ids = append(ids, fmt.Sprint(mapping["id"]))
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

	if err := d.Set("rules", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
