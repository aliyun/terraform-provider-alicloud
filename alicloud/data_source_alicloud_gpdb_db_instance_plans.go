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

func dataSourceAlicloudGpdbDbInstancePlans() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudGpdbDbInstancePlansRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
				ForceNew:     true,
			},
			"names": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"db_instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"plan_schedule_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Postpone", "Regular"}, false),
			},
			"plan_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"PauseResume", "Resize"}, false),
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"active", "cancel", "deleted", "finished"}, false),
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"plans": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"db_instance_plan_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"plan_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"plan_config": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"pause": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"execute_time": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"plan_cron_time": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"plan_task_status": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"resume": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"execute_time": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"plan_cron_time": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"plan_task_status": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"scale_in": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"execute_time": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"plan_cron_time": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"plan_task_status": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"segment_node_num": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"scale_out": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"execute_time": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"plan_cron_time": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"plan_task_status": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"segment_node_num": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
								},
							},
						},
						"plan_end_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"plan_start_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"plan_schedule_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"plan_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"plan_desc": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudGpdbDbInstancePlansRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeDBInstancePlans"
	request := make(map[string]interface{})
	request["DBInstanceId"] = d.Get("db_instance_id")

	if v, ok := d.GetOk("plan_schedule_type"); ok {
		request["PlanScheduleType"] = v
	}
	if v, ok := d.GetOk("plan_type"); ok {
		request["PlanType"] = v
	}
	var objects []map[string]interface{}
	var dbInstancePlanNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		dbInstancePlanNameRegex = r
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
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("gpdb", "2016-05-03", action, nil, request, true)
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
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_gpdb_db_instance_plans", action, AlibabaCloudSdkGoERROR)
	}
	resp, err := jsonpath.Get("$.Items.PlanList", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Items.PlanList", response)
	}
	result, _ := resp.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		if dbInstancePlanNameRegex != nil && !dbInstancePlanNameRegex.MatchString(fmt.Sprint(item["PlanName"])) {
			continue
		}
		if len(idsMap) > 0 {
			if _, ok := idsMap[fmt.Sprint(request["DBInstanceId"], ":", item["PlanId"])]; !ok {
				continue
			}
		}
		if statusOk && status.(string) != "" && status.(string) != item["PlanStatus"].(string) {
			continue
		}
		objects = append(objects, item)
	}
	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"plan_end_date":         object["PlanEndDate"],
			"id":                    fmt.Sprint(request["DBInstanceId"], ":", object["PlanId"]),
			"plan_id":               object["PlanId"],
			"db_instance_plan_name": object["PlanName"],
			"plan_start_date":       object["PlanStartDate"],
			"status":                object["PlanStatus"],
			"plan_schedule_type":    object["PlanScheduleType"],
			"plan_desc":             object["PlanDesc"],
			"plan_type":             object["PlanType"],
		}

		planConfigSli := make([]map[string]interface{}, 0)
		planConfig, err := convertJsonStringToMap(object["PlanConfig"].(string))
		if err != nil {
			return WrapError(err)
		}
		if len(planConfig) > 0 {
			planConfigMap := make(map[string]interface{})

			pauseSli := make([]map[string]interface{}, 0)
			if pause, ok := planConfig["pause"]; ok {
				if len(pause.(map[string]interface{})) > 0 {
					pauseMap := make(map[string]interface{})
					pauseMap["plan_cron_time"] = pause.(map[string]interface{})["planCronTime"]
					pauseMap["execute_time"] = pause.(map[string]interface{})["executeTime"]
					pauseMap["plan_task_status"] = pause.(map[string]interface{})["planTaskStatus"]
					pauseSli = append(pauseSli, pauseMap)
				}
			}
			planConfigMap["pause"] = pauseSli

			resumeSli := make([]map[string]interface{}, 0)
			if resume, ok := planConfig["resume"]; ok {
				if len(resume.(map[string]interface{})) > 0 {
					resumeMap := make(map[string]interface{})
					resumeMap["plan_cron_time"] = resume.(map[string]interface{})["planCronTime"]
					resumeMap["execute_time"] = resume.(map[string]interface{})["executeTime"]
					resumeMap["plan_task_status"] = resume.(map[string]interface{})["planTaskStatus"]
					resumeSli = append(resumeSli, resumeMap)
				}
			}
			planConfigMap["resume"] = resumeSli

			scaleInSli := make([]map[string]interface{}, 0)
			if scaleIn, ok := planConfig["scaleIn"]; ok {
				if len(scaleIn.(map[string]interface{})) > 0 {
					scaleInMap := make(map[string]interface{})
					scaleInMap["execute_time"] = scaleIn.(map[string]interface{})["executeTime"]
					scaleInMap["plan_cron_time"] = scaleIn.(map[string]interface{})["planCronTime"]
					scaleInMap["segment_node_num"] = scaleIn.(map[string]interface{})["segmentNodeNum"]
					scaleInMap["plan_task_status"] = scaleIn.(map[string]interface{})["planTaskStatus"]
					scaleInSli = append(scaleInSli, scaleInMap)
				}
			}
			planConfigMap["scale_in"] = scaleInSli

			scaleOutSli := make([]map[string]interface{}, 0)
			if scaleOut, ok := planConfig["scaleOut"]; ok {
				if len(scaleOut.(map[string]interface{})) > 0 {
					scaleOutMap := make(map[string]interface{})
					scaleOutMap["execute_time"] = scaleOut.(map[string]interface{})["executeTime"]
					scaleOutMap["plan_cron_time"] = scaleOut.(map[string]interface{})["planCronTime"]
					scaleOutMap["segment_node_num"] = scaleOut.(map[string]interface{})["segmentNodeNum"]
					scaleOutMap["plan_task_status"] = scaleOut.(map[string]interface{})["planTaskStatus"]

					scaleOutSli = append(scaleOutSli, scaleOutMap)
				}
			}
			planConfigMap["scale_out"] = scaleOutSli
			planConfigSli = append(planConfigSli, planConfigMap)
		}
		mapping["plan_config"] = planConfigSli
		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, object["PlanName"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("plans", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
