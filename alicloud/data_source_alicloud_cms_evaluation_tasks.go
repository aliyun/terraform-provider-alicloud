// Package alicloud. This file is hand-written from the Cms 2024-03-30 CloudSpec EvaluationTask definition.
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
)

func dataSourceAliCloudCmsEvaluationTasks() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudCmsEvaluationTasksRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"task_name_regex": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"workspace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"task_mode": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"channel": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"evaluation_tasks": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"channel": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"data_filter": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"data_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"evaluators": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"run_strategies": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"task_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"task_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"task_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"workspace": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAliCloudCmsEvaluationTasksRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	var objects []map[string]interface{}

	var taskNameRegex *regexp.Regexp
	if v, ok := d.GetOk("task_name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		taskNameRegex = r
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

	workspace := d.Get("workspace").(string)
	query := make(map[string]*string)
	query["workspace"] = StringPointer(workspace)
	if v, ok := d.GetOk("task_mode"); ok {
		query["taskMode"] = StringPointer(v.(string))
	}
	if v, ok := d.GetOk("status"); ok {
		query["status"] = StringPointer(v.(string))
	}
	if v, ok := d.GetOk("channel"); ok {
		query["channel"] = StringPointer(v.(string))
	}
	query["maxResults"] = StringPointer(strconv.Itoa(PageSizeLarge))

	var response map[string]interface{}
	var err error
	for {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutRead), func() *resource.RetryError {
			response, err = client.RoaGet("Cms", "2024-03-30", "/api/v1/evaluation-tasks", query, nil, nil)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug("/api/v1/evaluation-tasks", response, query)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), "/api/v1/evaluation-tasks", AlibabaCloudSdkGoERROR)
		}

		resp, _ := jsonpath.Get("$.evaluationTasks[*]", response)
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			itemId := fmt.Sprintf("%v:%v", workspace, item["taskId"])
			if taskNameRegex != nil {
				if !taskNameRegex.MatchString(fmt.Sprint(item["taskName"])) {
					continue
				}
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[itemId]; !ok {
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

	ids := make([]string, 0, len(objects))
	s := make([]map[string]interface{}, 0, len(objects))
	for _, objectRaw := range objects {
		mapping := map[string]interface{}{}
		mapping["channel"] = objectRaw["channel"]
		mapping["create_time"] = objectRaw["createdAt"]
		mapping["data_filter"] = objectRaw["dataFilter"]
		mapping["data_type"] = objectRaw["dataType"]
		mapping["description"] = objectRaw["description"]
		mapping["evaluators"] = objectRaw["evaluators"]
		mapping["id"] = fmt.Sprintf("%v:%v", workspace, objectRaw["taskId"])
		mapping["run_strategies"] = objectRaw["runStrategies"]
		mapping["status"] = objectRaw["status"]
		mapping["task_id"] = objectRaw["taskId"]
		mapping["task_mode"] = objectRaw["taskMode"]
		mapping["task_name"] = objectRaw["taskName"]
		mapping["workspace"] = objectRaw["workspace"]

		ids = append(ids, mapping["id"].(string))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("evaluation_tasks", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		if err := writeToFile(output.(string), s); err != nil {
			return WrapError(err)
		}
	}

	return nil
}
