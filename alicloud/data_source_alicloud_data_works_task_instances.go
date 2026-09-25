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

func dataSourceAlicloudDataWorksTaskInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudDataWorksTaskInstancesRead,
		Schema: map[string]*schema.Schema{
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
			// --- ListTaskInstances request inputs (PascalCase API params) ---
			// Bizdate and ProjectId are required by the ListTaskInstances API.
			// Bizdate is a UNIX timestamp in milliseconds, e.g. 1743350400000.
			"bizdate": {
				Type:     schema.TypeString,
				Required: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			// project_env accepts Prod or Dev (first-letter capital) per the
			// ListTaskInstances API; PROD/DEV/prod/dev are rejected with
			// InvalidProjectEnv. ValidateFunc gives a clear provider-side error
			// instead of surfacing the API-side 400.
			"project_env": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"Prod", "Dev"}, false),
			},
			"owner": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"page_number": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
			},
			"page_size": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  PageSizeLarge,
			},
			"sort_by": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"task_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"task_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"task_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"trigger_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"workflow_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"workflow_instance_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// --- ListTaskInstances response outputs (36 attribute mappings) ---
			"instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"task_instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"task_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"trigger_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"bizdate": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"task_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"project_env": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"owner": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"workflow_instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"workflow_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"workflow_instance_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"runtime_resource_resource_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"task_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"workflow_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"timeout": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"rerun_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"run_number": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"baseline_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"priority": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"runtime_resource_image": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"runtime_resource_cu": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"runtime_process_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"runtime_gateway": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"trigger_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"started_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"finished_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tenant_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"modify_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_user": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"modify_user": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"data_source_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"period_number": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudDataWorksTaskInstancesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListTaskInstances"
	request := make(map[string]interface{})
	// Bizdate and ProjectId are required by the ListTaskInstances API and are therefore always set.
	request["Bizdate"] = d.Get("bizdate").(string)
	request["ProjectId"] = d.Get("project_id").(string)
	if v, ok := d.GetOk("project_env"); ok {
		request["ProjectEnv"] = v
	}
	if v, ok := d.GetOk("owner"); ok {
		request["Owner"] = v
	}
	if v, ok := d.GetOk("sort_by"); ok {
		request["SortBy"] = v
	}
	if v, ok := d.GetOk("status"); ok {
		request["Status"] = v
	}
	if v, ok := d.GetOk("task_id"); ok {
		request["TaskId"] = v
	}
	if v, ok := d.GetOk("task_name"); ok {
		request["TaskName"] = v
	}
	if v, ok := d.GetOk("task_type"); ok {
		request["TaskType"] = v
	}
	if v, ok := d.GetOk("trigger_type"); ok {
		request["TriggerType"] = v
	}
	if v, ok := d.GetOk("workflow_id"); ok {
		request["WorkflowId"] = v
	}
	if v, ok := d.GetOk("workflow_instance_id"); ok {
		request["WorkflowInstanceId"] = v
	}

	pageSize := PageSizeLarge
	if v, ok := d.GetOk("page_size"); ok && v.(int) > 0 {
		pageSize = v.(int)
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

	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		nameRegex = r
	}

	var objects []map[string]interface{}
	var response map[string]interface{}
	var err error
	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("dataworks-public", "2024-05-18", action, nil, request, true)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_data_works_task_instances", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.Data.TaskInstances", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Data.TaskInstances", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item, _ := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["TaskInstanceId"])]; !ok {
					continue
				}
			}
			if nameRegex != nil {
				taskName := fmt.Sprint(item["TaskName"])
				if !nameRegex.MatchString(taskName) {
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
			"id":                     fmt.Sprint(object["TaskInstanceId"]),
			"task_instance_id":       fmt.Sprint(object["TaskInstanceId"]),
			"task_type":              fmt.Sprint(object["TaskType"]),
			"trigger_type":           fmt.Sprint(object["TriggerType"]),
			"bizdate":                fmt.Sprint(object["Bizdate"]),
			"task_id":                fmt.Sprint(object["TaskId"]),
			"project_env":            fmt.Sprint(object["ProjectEnv"]),
			"owner":                  fmt.Sprint(object["Owner"]),
			"workflow_instance_id":   fmt.Sprint(object["WorkflowInstanceId"]),
			"project_id":             fmt.Sprint(object["ProjectId"]),
			"workflow_id":            fmt.Sprint(object["WorkflowId"]),
			"workflow_instance_type": fmt.Sprint(object["WorkflowInstanceType"]),
			"task_name":              fmt.Sprint(object["TaskName"]),
			"region_id":              fmt.Sprint(object["RegionId"]),
			"description":            fmt.Sprint(object["Description"]),
			"workflow_name":          fmt.Sprint(object["WorkflowName"]),
			"timeout":                fmt.Sprint(object["Timeout"]),
			"rerun_mode":             fmt.Sprint(object["RerunMode"]),
			"run_number":             fmt.Sprint(object["RunNumber"]),
			"baseline_id":            fmt.Sprint(object["BaselineId"]),
			"priority":               fmt.Sprint(object["Priority"]),
			"trigger_time":           fmt.Sprint(object["TriggerTime"]),
			"started_time":           fmt.Sprint(object["StartedTime"]),
			"finished_time":          fmt.Sprint(object["FinishedTime"]),
			"tenant_id":              fmt.Sprint(object["TenantId"]),
			"create_time":            fmt.Sprint(object["CreateTime"]),
			"modify_time":            fmt.Sprint(object["ModifyTime"]),
			"create_user":            fmt.Sprint(object["CreateUser"]),
			"modify_user":            fmt.Sprint(object["ModifyUser"]),
			"status":                 fmt.Sprint(object["Status"]),
			"period_number":          fmt.Sprint(object["PeriodNumber"]),
		}
		// nested fields (RuntimeResource.* / Runtime.* / DataSource.*) — safe extraction
		if rr, ok := object["RuntimeResource"].(map[string]interface{}); ok {
			mapping["runtime_resource_resource_group_id"] = fmt.Sprint(rr["ResourceGroupId"])
			mapping["runtime_resource_image"] = fmt.Sprint(rr["Image"])
			mapping["runtime_resource_cu"] = fmt.Sprint(rr["Cu"])
		}
		if rt, ok := object["Runtime"].(map[string]interface{}); ok {
			mapping["runtime_process_id"] = fmt.Sprint(rt["ProcessId"])
			mapping["runtime_gateway"] = fmt.Sprint(rt["Gateway"])
		}
		if ds, ok := object["DataSource"].(map[string]interface{}); ok {
			mapping["data_source_name"] = fmt.Sprint(ds["Name"])
		}
		ids = append(ids, fmt.Sprint(mapping["id"]))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if err := d.Set("instances", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
