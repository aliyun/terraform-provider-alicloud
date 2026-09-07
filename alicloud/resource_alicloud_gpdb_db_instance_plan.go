package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudGpdbDbInstancePlan() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudGpdbDbInstancePlanCreate,
		Read:   resourceAliCloudGpdbDbInstancePlanRead,
		Update: resourceAliCloudGpdbDbInstancePlanUpdate,
		Delete: resourceAliCloudGpdbDbInstancePlanDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"db_instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"db_instance_plan_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"plan_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"PauseResume", "Resize", "ModifySpec"}, false),
			},
			"plan_schedule_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"Postpone", "Regular"}, false),
			},
			"plan_start_date": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"plan_end_date": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"plan_desc": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"active", "cancel"}, false),
			},
			"plan_config": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resume": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"execute_time": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"plan_cron_time": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"plan_task_status": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"pause": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"execute_time": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"plan_cron_time": {
										Type:     schema.TypeString,
										Optional: true,
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
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"segment_node_num": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"execute_time": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"plan_cron_time": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"plan_task_status": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"scale_out": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"segment_node_num": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"execute_time": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"plan_cron_time": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"plan_task_status": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"scale_up": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"instance_spec": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"execute_time": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"plan_cron_time": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"plan_task_status": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"scale_down": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"instance_spec": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"execute_time": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"plan_cron_time": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"plan_task_status": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"plan_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudGpdbDbInstancePlanCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gpdbService := GpdbService{client}
	var response map[string]interface{}
	action := "CreateDBInstancePlan"
	request := make(map[string]interface{})
	var err error

	request["DBInstanceId"] = d.Get("db_instance_id")
	request["PlanName"] = d.Get("db_instance_plan_name")
	request["PlanType"] = d.Get("plan_type")
	request["PlanScheduleType"] = d.Get("plan_schedule_type")

	if v, ok := d.GetOk("plan_start_date"); ok {
		request["PlanStartDate"] = v
	}

	if v, ok := d.GetOk("plan_end_date"); ok {
		request["PlanEndDate"] = v
	}

	if v, ok := d.GetOk("plan_desc"); ok {
		request["PlanDesc"] = v
	}

	planConfig := d.Get("plan_config")
	planConfigMap := map[string]interface{}{}
	for _, planConfigList := range planConfig.([]interface{}) {
		planConfigArg := planConfigList.(map[string]interface{})

		if resume, ok := planConfigArg["resume"]; ok && len(resume.([]interface{})) > 0 {
			resumeMap := map[string]interface{}{}
			for _, resumeList := range resume.([]interface{}) {
				resumeArg := resumeList.(map[string]interface{})

				if executeTime, ok := resumeArg["execute_time"]; ok {
					resumeMap["executeTime"] = executeTime
				}

				if planCronTime, ok := resumeArg["plan_cron_time"]; ok {
					resumeMap["planCronTime"] = planCronTime
				}
			}

			planConfigMap["resume"] = resumeMap
		}

		if pause, ok := planConfigArg["pause"]; ok && len(pause.([]interface{})) > 0 {
			pauseMap := map[string]interface{}{}
			for _, pauseList := range pause.([]interface{}) {
				pauseArg := pauseList.(map[string]interface{})

				if executeTime, ok := pauseArg["execute_time"]; ok {
					pauseMap["executeTime"] = executeTime
				}

				if planCronTime, ok := pauseArg["plan_cron_time"]; ok {
					pauseMap["planCronTime"] = planCronTime
				}
			}

			planConfigMap["pause"] = pauseMap
		}

		if scaleIn, ok := planConfigArg["scale_in"]; ok && len(scaleIn.([]interface{})) > 0 {
			scaleInMap := map[string]interface{}{}
			for _, scaleInList := range scaleIn.([]interface{}) {
				scaleInArg := scaleInList.(map[string]interface{})

				if segmentNodeNum, ok := scaleInArg["segment_node_num"]; ok {
					scaleInMap["segmentNodeNum"] = segmentNodeNum
				}

				if executeTime, ok := scaleInArg["execute_time"]; ok {
					scaleInMap["executeTime"] = executeTime
				}

				if planCronTime, ok := scaleInArg["plan_cron_time"]; ok {
					scaleInMap["planCronTime"] = planCronTime
				}
			}

			planConfigMap["scaleIn"] = scaleInMap
		}

		if scaleOut, ok := planConfigArg["scale_out"]; ok && len(scaleOut.([]interface{})) > 0 {
			scaleOutMap := map[string]interface{}{}
			for _, scaleOutList := range scaleOut.([]interface{}) {
				scaleOutArg := scaleOutList.(map[string]interface{})

				if segmentNodeNum, ok := scaleOutArg["segment_node_num"]; ok {
					scaleOutMap["segmentNodeNum"] = segmentNodeNum
				}

				if executeTime, ok := scaleOutArg["execute_time"]; ok {
					scaleOutMap["executeTime"] = executeTime
				}

				if planCronTime, ok := scaleOutArg["plan_cron_time"]; ok {
					scaleOutMap["planCronTime"] = planCronTime
				}
			}

			planConfigMap["scaleOut"] = scaleOutMap
		}

		if scaleUp, ok := planConfigArg["scale_up"]; ok && len(scaleUp.([]interface{})) > 0 {
			scaleUpMap := map[string]interface{}{}
			for _, scaleUpList := range scaleUp.([]interface{}) {
				scaleUpArg := scaleUpList.(map[string]interface{})

				if instanceSpec, ok := scaleUpArg["instance_spec"]; ok {
					scaleUpMap["instanceSpec"] = instanceSpec
				}

				if executeTime, ok := scaleUpArg["execute_time"]; ok {
					scaleUpMap["executeTime"] = executeTime
				}

				if planCronTime, ok := scaleUpArg["plan_cron_time"]; ok {
					scaleUpMap["planCronTime"] = planCronTime
				}
			}

			planConfigMap["scaleUp"] = scaleUpMap
		}

		if scaleDown, ok := planConfigArg["scale_down"]; ok && len(scaleDown.([]interface{})) > 0 {
			scaleDownMap := map[string]interface{}{}
			for _, scaleDownList := range scaleDown.([]interface{}) {
				scaleDownArg := scaleDownList.(map[string]interface{})

				if instanceSpec, ok := scaleDownArg["instance_spec"]; ok {
					scaleDownMap["instanceSpec"] = instanceSpec
				}

				if executeTime, ok := scaleDownArg["execute_time"]; ok {
					scaleDownMap["executeTime"] = executeTime
				}

				if planCronTime, ok := scaleDownArg["plan_cron_time"]; ok {
					scaleDownMap["planCronTime"] = planCronTime
				}
			}

			planConfigMap["scaleDown"] = scaleDownMap
		}
	}

	planConfigJson, err := convertMaptoJsonString(planConfigMap)
	if err != nil {
		return WrapError(err)
	}

	request["PlanConfig"] = planConfigJson

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_gpdb_db_instance_plan", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", request["DBInstanceId"], response["PlanId"]))

	stateConf := BuildStateConf([]string{}, []string{"active"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, gpdbService.GpdbDbInstancePlanStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudGpdbDbInstancePlanUpdate(d, meta)
}

func resourceAliCloudGpdbDbInstancePlanRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gpdbService := GpdbService{client}

	object, err := gpdbService.DescribeGpdbDbInstancePlan(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_gpdb_db_instance_plan gpdbService.DescribeGpdbDbInstancePlan Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("db_instance_id", object["DBInstanceId"])
	d.Set("plan_id", object["PlanId"])
	d.Set("db_instance_plan_name", object["PlanName"])
	d.Set("plan_type", object["PlanType"])
	d.Set("plan_schedule_type", object["PlanScheduleType"])
	d.Set("plan_start_date", object["PlanStartDate"])
	d.Set("plan_end_date", object["PlanEndDate"])
	d.Set("plan_desc", object["PlanDesc"])
	d.Set("status", object["PlanStatus"])

	if planConfig, ok := object["PlanConfig"].(string); ok && planConfig != "" {
		planConfigArg, err := convertJsonStringToMap(planConfig)
		if err != nil {
			return WrapError(err)
		}

		planConfigMaps := make([]map[string]interface{}, 0)
		planConfigMap := make(map[string]interface{})

		if resume, ok := planConfigArg["resume"]; ok {
			resumeMaps := make([]map[string]interface{}, 0)
			resumeArg := resume.(map[string]interface{})
			resumeMap := map[string]interface{}{}

			if executeTime, ok := resumeArg["executeTime"]; ok {
				resumeMap["execute_time"] = executeTime
			}

			if planCronTime, ok := resumeArg["planCronTime"]; ok {
				resumeMap["plan_cron_time"] = planCronTime
			}

			if planTaskStatus, ok := resumeArg["planTaskStatus"]; ok {
				resumeMap["plan_task_status"] = planTaskStatus
			}

			resumeMaps = append(resumeMaps, resumeMap)

			planConfigMap["resume"] = resumeMaps
		}

		if pause, ok := planConfigArg["pause"]; ok {
			pauseMaps := make([]map[string]interface{}, 0)
			pauseArg := pause.(map[string]interface{})
			pauseMap := map[string]interface{}{}

			if executeTime, ok := pauseArg["executeTime"]; ok {
				pauseMap["execute_time"] = executeTime
			}

			if planCronTime, ok := pauseArg["planCronTime"]; ok {
				pauseMap["plan_cron_time"] = planCronTime
			}

			if planTaskStatus, ok := pauseArg["planTaskStatus"]; ok {
				pauseArg["plan_task_status"] = planTaskStatus
			}

			pauseMaps = append(pauseMaps, pauseMap)

			planConfigMap["pause"] = pauseMaps
		}

		if scaleIn, ok := planConfigArg["scaleIn"]; ok {
			scaleInMaps := make([]map[string]interface{}, 0)
			scaleInArg := scaleIn.(map[string]interface{})
			scaleInMap := map[string]interface{}{}

			if segmentNodeNum, ok := scaleInArg["segmentNodeNum"]; ok {
				scaleInMap["segment_node_num"] = segmentNodeNum
			}

			if executeTime, ok := scaleInArg["executeTime"]; ok {
				scaleInMap["execute_time"] = executeTime
			}

			if planCronTime, ok := scaleInArg["planCronTime"]; ok {
				scaleInMap["plan_cron_time"] = planCronTime
			}

			if planTaskStatus, ok := scaleInArg["planTaskStatus"]; ok {
				scaleInArg["plan_task_status"] = planTaskStatus
			}

			scaleInMaps = append(scaleInMaps, scaleInMap)

			planConfigMap["scale_in"] = scaleInMaps
		}

		if scaleOut, ok := planConfigArg["scaleOut"]; ok {
			scaleOutMaps := make([]map[string]interface{}, 0)
			scaleOutArg := scaleOut.(map[string]interface{})
			scaleOutMap := map[string]interface{}{}

			if segmentNodeNum, ok := scaleOutArg["segmentNodeNum"]; ok {
				scaleOutMap["segment_node_num"] = segmentNodeNum
			}

			if executeTime, ok := scaleOutArg["executeTime"]; ok {
				scaleOutMap["execute_time"] = executeTime
			}

			if planCronTime, ok := scaleOutArg["planCronTime"]; ok {
				scaleOutMap["plan_cron_time"] = planCronTime
			}

			if planTaskStatus, ok := scaleOutArg["planTaskStatus"]; ok {
				scaleOutArg["plan_task_status"] = planTaskStatus
			}

			scaleOutMaps = append(scaleOutMaps, scaleOutMap)

			planConfigMap["scale_out"] = scaleOutMaps
		}

		if scaleUp, ok := planConfigArg["scaleUp"]; ok {
			scaleUpMaps := make([]map[string]interface{}, 0)
			scaleUpArg := scaleUp.(map[string]interface{})
			scaleUpMap := map[string]interface{}{}

			if instanceSpec, ok := scaleUpArg["instanceSpec"]; ok {
				scaleUpMap["instance_spec"] = instanceSpec
			}

			if executeTime, ok := scaleUpArg["executeTime"]; ok {
				scaleUpMap["execute_time"] = executeTime
			}

			if planCronTime, ok := scaleUpArg["planCronTime"]; ok {
				scaleUpMap["plan_cron_time"] = planCronTime
			}

			if planTaskStatus, ok := scaleUpArg["planTaskStatus"]; ok {
				scaleUpArg["plan_task_status"] = planTaskStatus
			}

			scaleUpMaps = append(scaleUpMaps, scaleUpMap)

			planConfigMap["scale_up"] = scaleUpMaps
		}

		if scaleDown, ok := planConfigArg["scaleDown"]; ok {
			scaleDownMaps := make([]map[string]interface{}, 0)
			scaleDownArg := scaleDown.(map[string]interface{})
			scaleDownMap := map[string]interface{}{}

			if instanceSpec, ok := scaleDownArg["instanceSpec"]; ok {
				scaleDownMap["instance_spec"] = instanceSpec
			}

			if executeTime, ok := scaleDownArg["executeTime"]; ok {
				scaleDownMap["execute_time"] = executeTime
			}

			if planCronTime, ok := scaleDownArg["planCronTime"]; ok {
				scaleDownMap["plan_cron_time"] = planCronTime
			}

			if planTaskStatus, ok := scaleDownArg["planTaskStatus"]; ok {
				scaleDownArg["plan_task_status"] = planTaskStatus
			}

			scaleDownMaps = append(scaleDownMaps, scaleDownMap)

			planConfigMap["scale_down"] = scaleDownMaps
		}

		planConfigMaps = append(planConfigMaps, planConfigMap)

		d.Set("plan_config", planConfigMaps)
	}

	return nil
}

func resourceAliCloudGpdbDbInstancePlanUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gpdbService := GpdbService{client}
	var response map[string]interface{}
	d.Partial(true)

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}

	update := false
	updateDBInstancePlanReq := map[string]interface{}{
		"DBInstanceId": parts[0],
		"PlanId":       parts[1],
	}

	if !d.IsNewResource() && d.HasChange("db_instance_plan_name") {
		update = true
	}
	updateDBInstancePlanReq["PlanName"] = d.Get("db_instance_plan_name")

	if !d.IsNewResource() && d.HasChange("plan_start_date") {
		update = true
	}
	if v, ok := d.GetOk("plan_start_date"); ok {
		updateDBInstancePlanReq["PlanStartDate"] = v
	}

	if !d.IsNewResource() && d.HasChange("plan_end_date") {
		update = true
	}
	if v, ok := d.GetOk("plan_end_date"); ok {
		updateDBInstancePlanReq["PlanEndDate"] = v
	}

	if !d.IsNewResource() && d.HasChange("plan_desc") {
		update = true
	}
	if v, ok := d.GetOk("plan_desc"); ok {
		updateDBInstancePlanReq["PlanDesc"] = v
	}

	if !d.IsNewResource() && d.HasChange("plan_config") {
		update = true
	}
	planConfig := d.Get("plan_config")
	planConfigMap := map[string]interface{}{}
	for _, planConfigList := range planConfig.([]interface{}) {
		planConfigArg := planConfigList.(map[string]interface{})

		if resume, ok := planConfigArg["resume"]; ok && len(resume.([]interface{})) > 0 {
			resumeMap := map[string]interface{}{}
			for _, resumeList := range resume.([]interface{}) {
				resumeArg := resumeList.(map[string]interface{})

				if executeTime, ok := resumeArg["execute_time"]; ok {
					resumeMap["executeTime"] = executeTime
				}

				if planCronTime, ok := resumeArg["plan_cron_time"]; ok {
					resumeMap["planCronTime"] = planCronTime
				}
			}

			planConfigMap["resume"] = resumeMap
		}

		if pause, ok := planConfigArg["pause"]; ok && len(pause.([]interface{})) > 0 {
			pauseMap := map[string]interface{}{}
			for _, pauseList := range pause.([]interface{}) {
				pauseArg := pauseList.(map[string]interface{})

				if executeTime, ok := pauseArg["execute_time"]; ok {
					pauseMap["executeTime"] = executeTime
				}

				if planCronTime, ok := pauseArg["plan_cron_time"]; ok {
					pauseMap["planCronTime"] = planCronTime
				}
			}

			planConfigMap["pause"] = pauseMap
		}

		if scaleIn, ok := planConfigArg["scale_in"]; ok && len(scaleIn.([]interface{})) > 0 {
			scaleInMap := map[string]interface{}{}
			for _, scaleInList := range scaleIn.([]interface{}) {
				scaleInArg := scaleInList.(map[string]interface{})

				if segmentNodeNum, ok := scaleInArg["segment_node_num"]; ok {
					scaleInMap["segmentNodeNum"] = segmentNodeNum
				}

				if executeTime, ok := scaleInArg["execute_time"]; ok {
					scaleInMap["executeTime"] = executeTime
				}

				if planCronTime, ok := scaleInArg["plan_cron_time"]; ok {
					scaleInMap["planCronTime"] = planCronTime
				}
			}

			planConfigMap["scaleIn"] = scaleInMap
		}

		if scaleOut, ok := planConfigArg["scale_out"]; ok && len(scaleOut.([]interface{})) > 0 {
			scaleOutMap := map[string]interface{}{}
			for _, scaleOutList := range scaleOut.([]interface{}) {
				scaleOutArg := scaleOutList.(map[string]interface{})

				if segmentNodeNum, ok := scaleOutArg["segment_node_num"]; ok {
					scaleOutMap["segmentNodeNum"] = segmentNodeNum
				}

				if executeTime, ok := scaleOutArg["execute_time"]; ok {
					scaleOutMap["executeTime"] = executeTime
				}

				if planCronTime, ok := scaleOutArg["plan_cron_time"]; ok {
					scaleOutMap["planCronTime"] = planCronTime
				}
			}

			planConfigMap["scaleOut"] = scaleOutMap
		}

		if scaleUp, ok := planConfigArg["scale_up"]; ok && len(scaleUp.([]interface{})) > 0 {
			scaleUpMap := map[string]interface{}{}
			for _, scaleUpList := range scaleUp.([]interface{}) {
				scaleUpArg := scaleUpList.(map[string]interface{})

				if instanceSpec, ok := scaleUpArg["instance_spec"]; ok {
					scaleUpMap["instanceSpec"] = instanceSpec
				}

				if executeTime, ok := scaleUpArg["execute_time"]; ok {
					scaleUpMap["executeTime"] = executeTime
				}

				if planCronTime, ok := scaleUpArg["plan_cron_time"]; ok {
					scaleUpMap["planCronTime"] = planCronTime
				}
			}

			planConfigMap["scaleUp"] = scaleUpMap
		}

		if scaleDown, ok := planConfigArg["scale_down"]; ok && len(scaleDown.([]interface{})) > 0 {
			scaleDownMap := map[string]interface{}{}
			for _, scaleDownList := range scaleDown.([]interface{}) {
				scaleDownArg := scaleDownList.(map[string]interface{})

				if instanceSpec, ok := scaleDownArg["instance_spec"]; ok {
					scaleDownMap["instanceSpec"] = instanceSpec
				}

				if executeTime, ok := scaleDownArg["execute_time"]; ok {
					scaleDownMap["executeTime"] = executeTime
				}

				if planCronTime, ok := scaleDownArg["plan_cron_time"]; ok {
					scaleDownMap["planCronTime"] = planCronTime
				}
			}

			planConfigMap["scaleDown"] = scaleDownMap
		}
	}

	planConfigJson, err := convertMaptoJsonString(planConfigMap)
	if err != nil {
		return WrapError(err)
	}

	updateDBInstancePlanReq["PlanConfig"] = planConfigJson

	if update {
		action := "UpdateDBInstancePlan"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.RpcPost("gpdb", "2016-05-03", action, nil, updateDBInstancePlanReq, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, updateDBInstancePlanReq)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		stateConf := BuildStateConf([]string{}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, gpdbService.GpdbDbInstancePlanStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}

		d.SetPartial("db_instance_plan_name")
		d.SetPartial("plan_start_date")
		d.SetPartial("plan_end_date")
		d.SetPartial("plan_desc")
		d.SetPartial("plan_config")
	}

	if d.HasChange("status") {
		object, err := gpdbService.DescribeGpdbDbInstancePlan(d.Id())
		if err != nil {
			return WrapError(err)
		}

		target := d.Get("status").(string)
		if object["PlanStatus"].(string) != target {
			action := "SetDBInstancePlanStatus"
			request := map[string]interface{}{
				"DBInstanceId": parts[0],
				"PlanId":       parts[1],
			}

			switch target {
			case "active":
				request["PlanStatus"] = "enable"
			case "cancel":
				request["PlanStatus"] = "disable"
			}
			wait := incrementalWait(3*time.Second, 3*time.Second)
			err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
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
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}

			stateConf := BuildStateConf([]string{}, []string{target}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, gpdbService.GpdbDbInstancePlanStateRefreshFunc(d.Id(), []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}

			d.SetPartial("status")
		}
	}

	d.Partial(false)

	return resourceAliCloudGpdbDbInstancePlanRead(d, meta)
}

func resourceAliCloudGpdbDbInstancePlanDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gpdbService := GpdbService{client}
	action := "DeleteDBInstancePlan"
	var response map[string]interface{}

	var err error

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}

	request := map[string]interface{}{
		"DBInstanceId": parts[0],
		"PlanId":       parts[1],
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
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
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, gpdbService.GpdbDbInstancePlanStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}
