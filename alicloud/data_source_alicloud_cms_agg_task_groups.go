package alicloud

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAliCloudCmsAggTaskGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudCmsAggTaskGroupsRead,
		Schema: map[string]*schema.Schema{
			"source_prometheus_id": {
				Type:     schema.TypeString,
				Required: true,
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
			"names": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
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
						"agg_task_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"agg_task_group_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_prometheus_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"target_prometheus_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cron_expr": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"delay": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"from_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"interval": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"max_retries": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"max_run_time_in_seconds": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"schedule_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"schedule_time_expr": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"to_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"update_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tags": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
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

func dataSourceAliCloudCmsAggTaskGroupsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	var objects []map[string]interface{}
	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		nameRegex = r
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

	sourcePrometheusId := d.Get("source_prometheus_id").(string)
	action := fmt.Sprintf("/prometheus-instances/%s/agg-task-groups", sourcePrometheusId)
	var response map[string]interface{}
	var err error
	query := make(map[string]*string)
	query["maxResults"] = StringPointer("100")

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
			return WrapErrorf(err, DefaultErrorMsg, "alicloud_cms_agg_task_groups", action, AlibabaCloudSdkGoERROR)
		}

		resp, _ := jsonpath.Get("$.aggTaskGroups[*]", response)
		result, _ := resp.([]interface{})
		for _, v := range result {
			item, ok := v.(map[string]interface{})
			if !ok {
				continue
			}
			if nameRegex != nil && !nameRegex.MatchString(fmt.Sprint(item["aggTaskGroupName"])) {
				continue
			}
			if len(idsMap) > 0 {
				id := fmt.Sprintf("%v:%v", sourcePrometheusId, item["aggTaskGroupId"])
				if _, ok := idsMap[fmt.Sprint(item["aggTaskGroupId"])]; !ok {
					if _, ok := idsMap[id]; !ok {
						continue
					}
				}
			}
			objects = append(objects, item)
		}

		nextToken, _ := response["nextToken"].(string)
		if nextToken == "" {
			break
		}
		query["nextToken"] = StringPointer(nextToken)
	}

	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, objectRaw := range objects {
		mapping := map[string]interface{}{}

		id := fmt.Sprintf("%v:%v", sourcePrometheusId, objectRaw["aggTaskGroupId"])
		mapping["id"] = id
		mapping["agg_task_group_id"] = objectRaw["aggTaskGroupId"]
		mapping["agg_task_group_name"] = objectRaw["aggTaskGroupName"]
		mapping["source_prometheus_id"] = objectRaw["sourcePrometheusId"]
		mapping["target_prometheus_id"] = objectRaw["targetPrometheusId"]
		mapping["cron_expr"] = objectRaw["cronExpr"]
		mapping["delay"] = objectRaw["delay"]
		mapping["description"] = objectRaw["description"]
		mapping["from_time"] = objectRaw["fromTime"]
		mapping["interval"] = objectRaw["interval"]
		mapping["max_retries"] = objectRaw["maxRetries"]
		mapping["max_run_time_in_seconds"] = objectRaw["maxRunTimeInSeconds"]
		mapping["region_id"] = objectRaw["regionId"]
		mapping["schedule_mode"] = objectRaw["scheduleMode"]
		mapping["schedule_time_expr"] = objectRaw["scheduleTimeExpr"]
		mapping["status"] = objectRaw["status"]
		mapping["to_time"] = objectRaw["toTime"]
		mapping["update_time"] = objectRaw["updateTime"]

		tagsMaps := make(map[string]interface{})
		if tagsRaw, ok := objectRaw["tags"]; ok && tagsRaw != nil {
			for _, tagRaw := range convertToInterfaceArray(tagsRaw) {
				tagMap, ok := tagRaw.(map[string]interface{})
				if !ok {
					continue
				}
				key := tagMap["key"]
				if key == nil {
					key = tagMap["Key"]
				}
				value := tagMap["value"]
				if value == nil {
					value = tagMap["Value"]
				}
				if key != nil && value != nil {
					keyStr := fmt.Sprint(key)
					if strings.HasPrefix(keyStr, "acs:") {
						continue
					}
					tagsMaps[keyStr] = fmt.Sprint(value)
				}
			}
		}
		mapping["tags"] = tagsMaps

		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, objectRaw["aggTaskGroupName"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
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
