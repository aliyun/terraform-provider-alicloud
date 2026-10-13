package alicloud

import (
	"fmt"
	"regexp"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudCmsBizTraces() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudCmsBizTracesRead,
		Schema: map[string]*schema.Schema{
			"workspace": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"name_regex": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"biz_traces": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"biz_trace_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"biz_trace_code": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"biz_trace_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"rule_config": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"advanced_config": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"workspace": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudCmsBizTracesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "/bizTraces"
	query := make(map[string]*string)
	if v, ok := d.GetOk("workspace"); ok {
		query["workspace"] = StringPointer(v.(string))
	}
	query["maxResults"] = StringPointer(fmt.Sprint(PageSizeLarge))

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
		wait := incrementalWait(3*time.Second, 5*time.Second)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_cms_biz_traces", action, AlibabaCloudSdkGoERROR)
		}

		resp, err := jsonpath.Get("$.items", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.items", response)
		}
		result, _ := resp.([]interface{})
		if result == nil {
			result = []interface{}{}
		}
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["bizTraceId"])]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}

		nextToken, _ := jsonpath.Get("$.nextToken", response)
		nextTokenStr := fmt.Sprint(nextToken)
		if nextTokenStr == "" || nextTokenStr == "<nil>" || len(result) == 0 {
			break
		}
		query["nextToken"] = StringPointer(nextTokenStr)
	}

	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	var nameRegexFilter *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		nameRegexFilter = r
	}
	for _, object := range objects {
		if nameRegexFilter != nil && !nameRegexFilter.MatchString(fmt.Sprint(object["bizTraceName"])) {
			continue
		}
		mapping := map[string]interface{}{
			"biz_trace_id":    fmt.Sprint(object["bizTraceId"]),
			"biz_trace_code":  object["bizTraceCode"],
			"biz_trace_name":  object["bizTraceName"],
			"rule_config":     object["ruleConfig"],
			"advanced_config": object["advancedConfig"],
			"workspace":       object["workspace"],
			"create_time":     object["createTime"],
			"region_id":       object["regionId"],
		}
		ids = append(ids, fmt.Sprint(mapping["biz_trace_id"]))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("biz_traces", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
