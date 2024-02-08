package alicloud

import (
	"fmt"
	"log"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
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
				ValidateFunc: StringInSlice([]string{"PauseResume", "Resize"}, false),
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
				MaxItems: 1,
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
	conn, err := client.NewGpdbClient()
	if err != nil {
		return WrapError(err)
	}

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
	}

	planConfigJson, err := convertMaptoJsonString(planConfigMap)
	if err != nil {
		return WrapError(err)
	}

	request["PlanConfig"] = planConfigJson

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-05-03"), StringPointer("AK"), nil, request, &runtime)
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

			scaleOutMaps = append(scaleOutMaps, scaleOutMap)

			planConfigMap["scale_out"] = scaleOutMaps
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
	}

	planConfigJson, err := convertMaptoJsonString(planConfigMap)
	if err != nil {
		return WrapError(err)
	}

	updateDBInstancePlanReq["PlanConfig"] = planConfigJson

	if update {
		action := "UpdateDBInstancePlan"
		conn, err := client.NewGpdbClient()
		if err != nil {
			return WrapError(err)
		}

		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-05-03"), StringPointer("AK"), nil, updateDBInstancePlanReq, &runtime)
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
			conn, err := client.NewGpdbClient()
			if err != nil {
				return WrapError(err)
			}

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

			runtime := util.RuntimeOptions{}
			runtime.SetAutoretry(true)
			wait := incrementalWait(3*time.Second, 3*time.Second)
			err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
				response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-05-03"), StringPointer("AK"), nil, request, &runtime)
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

	conn, err := client.NewGpdbClient()
	if err != nil {
		return WrapError(err)
	}

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}

	request := map[string]interface{}{
		"DBInstanceId": parts[0],
		"PlanId":       parts[1],
	}

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-05-03"), StringPointer("AK"), nil, request, &runtime)
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
