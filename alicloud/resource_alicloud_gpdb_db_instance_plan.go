package alicloud

import (
	"fmt"
	"log"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudGpdbDbInstancePlan() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudGpdbDbInstancePlanCreate,
		Read:   resourceAlicloudGpdbDbInstancePlanRead,
		Update: resourceAlicloudGpdbDbInstancePlanUpdate,
		Delete: resourceAlicloudGpdbDbInstancePlanDelete,
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
				Required: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"plan_config": {
				Required: true,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"pause": {
							Optional: true,
							MaxItems: 1,
							Type:     schema.TypeList,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"execute_time": {
										Optional: true,
										Computed: true,
										Type:     schema.TypeString,
									},
									"plan_cron_time": {
										Optional: true,
										Type:     schema.TypeString,
									},
								},
							},
						},
						"resume": {
							Optional: true,
							MaxItems: 1,
							Type:     schema.TypeList,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"execute_time": {
										Optional: true,
										Computed: true,
										Type:     schema.TypeString,
									},
									"plan_cron_time": {
										Optional: true,
										Type:     schema.TypeString,
									},
								},
							},
						},
						"scale_in": {
							Optional: true,
							MaxItems: 1,
							Type:     schema.TypeList,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"execute_time": {
										Optional: true,
										Computed: true,
										Type:     schema.TypeString,
									},
									"plan_cron_time": {
										Optional: true,
										Type:     schema.TypeString,
									},
									"segment_node_num": {
										Optional: true,
										Type:     schema.TypeString,
									},
								},
							},
						},
						"scale_out": {
							Optional: true,
							MaxItems: 1,
							Type:     schema.TypeList,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"execute_time": {
										Optional: true,
										Computed: true,
										Type:     schema.TypeString,
									},
									"plan_cron_time": {
										Optional: true,
										Type:     schema.TypeString,
									},
									"segment_node_num": {
										Optional: true,
										Type:     schema.TypeString,
									},
								},
							},
						},
					},
				},
			},
			"plan_desc": {
				Optional: true,
				Computed: true,
				Type:     schema.TypeString,
			},
			"plan_end_date": {
				Optional: true,
				Computed: true,
				Type:     schema.TypeString,
			},
			"db_instance_plan_name": {
				Required: true,
				Type:     schema.TypeString,
			},
			"plan_schedule_type": {
				Required:     true,
				ForceNew:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"Postpone", "Regular"}, false),
			},
			"plan_start_date": {
				Optional: true,
				Computed: true,
				Type:     schema.TypeString,
			},
			"plan_type": {
				Required:     true,
				ForceNew:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"PauseResume", "Resize"}, false),
			},
			"status": {
				Optional:     true,
				Computed:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"active", "cancel"}, false),
			},
			"plan_id": {
				Computed: true,
				Type:     schema.TypeString,
			},
		},
	}
}

func resourceAlicloudGpdbDbInstancePlanCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateDBInstancePlan"
	request := make(map[string]interface{})
	conn, err := client.NewGpdbClient()
	if err != nil {
		return WrapError(err)
	}
	request["PlanName"] = d.Get("db_instance_plan_name")
	if v, ok := d.GetOk("plan_desc"); ok {
		request["PlanDesc"] = v
	}
	if v, ok := d.GetOk("plan_end_date"); ok {
		request["PlanEndDate"] = v
	}
	request["DBInstanceId"] = d.Get("db_instance_id")
	request["PlanScheduleType"] = d.Get("plan_schedule_type")
	if v, ok := d.GetOk("plan_start_date"); ok {
		request["PlanStartDate"] = v
	}

	planConfig := make(map[string]interface{}, 0)
	if len(d.Get("plan_config").([]interface{})) > 0 {
		planConfigMap := d.Get("plan_config").([]interface{})[0].(map[string]interface{})

		if v, ok := planConfigMap["pause"]; ok && len(v.([]interface{})) > 0 {
			pause := make(map[string]interface{}, 0)
			for _, item := range v.([]interface{}) {
				pauseMap := item.(map[string]interface{})
				pause["planCronTime"] = pauseMap["plan_cron_time"]
				pause["executeTime"] = pauseMap["execute_time"]
			}
			planConfig["pause"] = pause
		}

		if v, ok := planConfigMap["resume"]; ok && len(v.([]interface{})) > 0 {
			resume := make(map[string]interface{}, 0)
			for _, item := range v.([]interface{}) {
				resumeMap := item.(map[string]interface{})
				resume["planCronTime"] = resumeMap["plan_cron_time"]
				resume["executeTime"] = resumeMap["execute_time"]
			}
			planConfig["resume"] = resume
		}

		if v, ok := planConfigMap["scale_in"]; ok && len(v.([]interface{})) > 0 {
			scaleIn := make(map[string]interface{}, 0)
			for _, item := range v.([]interface{}) {
				resumeMap := item.(map[string]interface{})
				scaleIn["segmentNodeNum"] = resumeMap["segment_node_num"]
				scaleIn["planCronTime"] = resumeMap["plan_cron_time"]
				scaleIn["executeTime"] = resumeMap["execute_time"]
			}
			planConfig["scaleIn"] = scaleIn
		}

		if v, ok := planConfigMap["scale_out"]; ok && len(v.([]interface{})) > 0 {
			scaleOut := make(map[string]interface{}, 0)
			for _, item := range v.([]interface{}) {
				resumeMap := item.(map[string]interface{})
				scaleOut["segmentNodeNum"] = resumeMap["segment_node_num"]
				scaleOut["planCronTime"] = resumeMap["plan_cron_time"]
				scaleOut["executeTime"] = resumeMap["execute_time"]
			}
			planConfig["scaleOut"] = scaleOut
		}
	}
	request["PlanConfig"], err = convertArrayObjectToJsonString(planConfig)
	if err != nil {
		return WrapError(err)
	}

	request["PlanType"] = d.Get("plan_type")
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-05-03"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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

	d.SetId(fmt.Sprint(request["DBInstanceId"], ":", response["PlanId"]))

	return resourceAlicloudGpdbDbInstancePlanUpdate(d, meta)
}
func resourceAlicloudGpdbDbInstancePlanRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gpdbService := GpdbService{client}
	object, err := gpdbService.DescribeGpdbDbInstancePlan(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_gpdb_db_instance_plan gpdbService.DescribeGpdbDbInstancePlan Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	d.Set("db_instance_id", parts[0])
	d.Set("plan_id", parts[1])

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
				scaleOutSli = append(scaleOutSli, scaleOutMap)
			}
		}
		planConfigMap["scale_out"] = scaleOutSli
		planConfigSli = append(planConfigSli, planConfigMap)
	}
	d.Set("plan_config", planConfigSli)
	d.Set("plan_end_date", object["PlanEndDate"])
	d.Set("db_instance_plan_name", object["PlanName"])
	d.Set("plan_start_date", object["PlanStartDate"])
	d.Set("status", object["PlanStatus"])
	d.Set("plan_desc", object["PlanDesc"])
	d.Set("plan_schedule_type", object["PlanScheduleType"])
	d.Set("plan_type", object["PlanType"])
	return nil
}
func resourceAlicloudGpdbDbInstancePlanUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gpdbService := GpdbService{client}
	conn, err := client.NewGpdbClient()
	if err != nil {
		return WrapError(err)
	}
	var response map[string]interface{}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	d.Partial(true)

	if d.HasChange("status") {
		object, err := gpdbService.DescribeGpdbDbInstancePlan(d.Id())
		if err != nil {
			return WrapError(err)
		}
		target := d.Get("status").(string)

		if object["PlanStatus"].(string) != target {
			action := "SetDBInstancePlanStatus"
			request := map[string]interface{}{
				"PlanId":       parts[1],
				"DBInstanceId": parts[0],
				"PlanStatus":   target,
			}
			if target == "cancel" {
				request["PlanStatus"] = "disable"
			}
			if target == "active" {
				request["PlanStatus"] = "enable"
			}
			wait := incrementalWait(3*time.Second, 3*time.Second)
			err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
				response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-05-03"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
			d.SetPartial("status")
		}

	}

	update := false
	updateDBInstancePlanReq := map[string]interface{}{
		"PlanId":       parts[1],
		"DBInstanceId": parts[0],
	}
	if !d.IsNewResource() && d.HasChange("plan_config") {
		update = true
		planConfig := make(map[string]interface{}, 0)
		if len(d.Get("plan_config").([]interface{})) > 0 {
			planConfigMap := d.Get("plan_config").([]interface{})[0].(map[string]interface{})

			if v, ok := planConfigMap["pause"]; ok && len(v.([]interface{})) > 0 {
				pause := make(map[string]interface{}, 0)
				for _, item := range v.([]interface{}) {
					pauseMap := item.(map[string]interface{})
					pause["planCronTime"] = pauseMap["plan_cron_time"]
					pause["executeTime"] = pauseMap["execute_time"]
				}
				planConfig["pause"] = pause
			}

			if v, ok := planConfigMap["resume"]; ok && len(v.([]interface{})) > 0 {
				resume := make(map[string]interface{}, 0)
				for _, item := range v.([]interface{}) {
					resumeMap := item.(map[string]interface{})
					resume["planCronTime"] = resumeMap["plan_cron_time"]
					resume["executeTime"] = resumeMap["execute_time"]
				}
				planConfig["resume"] = resume
			}

			if v, ok := planConfigMap["scale_in"]; ok && len(v.([]interface{})) > 0 {
				scaleIn := make(map[string]interface{}, 0)
				for _, item := range v.([]interface{}) {
					resumeMap := item.(map[string]interface{})
					scaleIn["segmentNodeNum"] = resumeMap["segment_node_num"]
					scaleIn["planCronTime"] = resumeMap["plan_cron_time"]
					scaleIn["executeTime"] = resumeMap["execute_time"]
				}
				planConfig["scaleIn"] = scaleIn
			}

			if v, ok := planConfigMap["scale_out"]; ok && len(v.([]interface{})) > 0 {
				scaleOut := make(map[string]interface{}, 0)
				for _, item := range v.([]interface{}) {
					resumeMap := item.(map[string]interface{})
					scaleOut["segmentNodeNum"] = resumeMap["segment_node_num"]
					scaleOut["planCronTime"] = resumeMap["plan_cron_time"]
					scaleOut["executeTime"] = resumeMap["execute_time"]
				}
				planConfig["scaleOut"] = scaleOut
			}
		}
		updateDBInstancePlanReq["PlanConfig"], err = convertArrayObjectToJsonString(planConfig)
		if err != nil {
			return WrapError(err)
		}
	}
	if v, ok := d.GetOk("plan_desc"); ok {
		updateDBInstancePlanReq["PlanDesc"] = v
	}
	if !d.IsNewResource() && d.HasChange("plan_desc") {
		update = true
	}
	if !d.IsNewResource() && d.HasChange("plan_end_date") {
		update = true
		if v, ok := d.GetOk("plan_end_date"); ok {
			updateDBInstancePlanReq["PlanEndDate"] = v
		}
	}
	if !d.IsNewResource() && d.HasChange("db_instance_plan_name") {
		update = true
		updateDBInstancePlanReq["PlanName"] = d.Get("db_instance_plan_name")
	}
	if !d.IsNewResource() && d.HasChange("plan_start_date") {
		update = true
		if v, ok := d.GetOk("plan_start_date"); ok {
			updateDBInstancePlanReq["PlanStartDate"] = v
		}
	}
	if update {
		action := "UpdateDBInstancePlan"
		conn, err := client.NewGpdbClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-05-03"), StringPointer("AK"), nil, updateDBInstancePlanReq, &util.RuntimeOptions{})
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
		d.SetPartial("db_instance_plan_name")
		d.SetPartial("plan_desc")
		d.SetPartial("plan_end_date")
		d.SetPartial("plan_start_date")
	}
	d.Partial(false)
	return resourceAlicloudGpdbDbInstancePlanRead(d, meta)
}
func resourceAlicloudGpdbDbInstancePlanDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	action := "DeleteDBInstancePlan"
	var response map[string]interface{}
	conn, err := client.NewGpdbClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"PlanId":       parts[1],
		"DBInstanceId": parts[0],
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-05-03"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
	return nil
}
