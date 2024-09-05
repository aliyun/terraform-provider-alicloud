// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
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

func resourceAliCloudFcv3ProvisionConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudFcv3ProvisionConfigCreate,
		Read:   resourceAliCloudFcv3ProvisionConfigRead,
		Update: resourceAliCloudFcv3ProvisionConfigUpdate,
		Delete: resourceAliCloudFcv3ProvisionConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"always_allocate_cpu": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"always_allocate_gpu": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"function_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"qualifier": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"scheduled_actions": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"schedule_expression": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"target": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"time_zone": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"end_time": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"start_time": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"target": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: IntBetween(0, 10000),
			},
			"target_tracking_policies": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"metric_target": {
							Type:     schema.TypeFloat,
							Optional: true,
						},
						"time_zone": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"end_time": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"metric_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"start_time": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"min_capacity": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"max_capacity": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceAliCloudFcv3ProvisionConfigCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	functionName := d.Get("function_name")
	action := fmt.Sprintf("/2023-03-30/functions/%s/provision-config", functionName)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	conn, err := client.NewFcv2Client()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	if v, ok := d.GetOk("function_name"); ok {
		request["functionName"] = v
	}

	if v, ok := d.GetOk("qualifier"); ok {
		query["qualifier"] = StringPointer(v.(string))
	}

	if v, ok := d.GetOk("target_tracking_policies"); ok {
		targetTrackingPoliciesMaps := make([]interface{}, 0)
		for _, dataLoop := range v.([]interface{}) {
			dataLoopTmp := dataLoop.(map[string]interface{})
			dataLoopMap := make(map[string]interface{})
			dataLoopMap["endTime"] = dataLoopTmp["end_time"]
			dataLoopMap["name"] = dataLoopTmp["name"]
			dataLoopMap["startTime"] = dataLoopTmp["start_time"]
			dataLoopMap["maxCapacity"] = dataLoopTmp["max_capacity"]
			dataLoopMap["metricTarget"] = dataLoopTmp["metric_target"]
			dataLoopMap["metricType"] = dataLoopTmp["metric_type"]
			dataLoopMap["minCapacity"] = dataLoopTmp["min_capacity"]
			dataLoopMap["timeZone"] = dataLoopTmp["time_zone"]
			targetTrackingPoliciesMaps = append(targetTrackingPoliciesMaps, dataLoopMap)
		}
		request["targetTrackingPolicies"] = targetTrackingPoliciesMaps
	}

	if v, ok := d.GetOk("scheduled_actions"); ok {
		scheduledActionsMaps := make([]interface{}, 0)
		for _, dataLoop1 := range v.([]interface{}) {
			dataLoop1Tmp := dataLoop1.(map[string]interface{})
			dataLoop1Map := make(map[string]interface{})
			dataLoop1Map["scheduleExpression"] = dataLoop1Tmp["schedule_expression"]
			dataLoop1Map["target"] = dataLoop1Tmp["target"]
			dataLoop1Map["endTime"] = dataLoop1Tmp["end_time"]
			dataLoop1Map["name"] = dataLoop1Tmp["name"]
			dataLoop1Map["startTime"] = dataLoop1Tmp["start_time"]
			dataLoop1Map["timeZone"] = dataLoop1Tmp["time_zone"]
			scheduledActionsMaps = append(scheduledActionsMaps, dataLoop1Map)
		}
		request["scheduledActions"] = scheduledActionsMaps
	}

	if v, ok := d.GetOk("target"); ok {
		request["target"] = v
	}
	if v, ok := d.GetOk("always_allocate_cpu"); ok {
		request["alwaysAllocateCPU"] = v
	}
	if v, ok := d.GetOk("always_allocate_gpu"); ok {
		request["alwaysAllocateGPU"] = v
	}
	body = request
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer("2023-03-30"), nil, StringPointer("PUT"), StringPointer("AK"), StringPointer(action), query, nil, body, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_fcv3_provision_config", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(functionName))

	return resourceAliCloudFcv3ProvisionConfigRead(d, meta)
}

func resourceAliCloudFcv3ProvisionConfigRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	fcv3ServiceV2 := Fcv3ServiceV2{client}

	objectRaw, err := fcv3ServiceV2.DescribeFcv3ProvisionConfig(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_fcv3_provision_config DescribeFcv3ProvisionConfig Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	if objectRaw["alwaysAllocateCPU"] != nil {
		d.Set("always_allocate_cpu", objectRaw["alwaysAllocateCPU"])
	}
	if objectRaw["alwaysAllocateGPU"] != nil {
		d.Set("always_allocate_gpu", objectRaw["alwaysAllocateGPU"])
	}
	if objectRaw["target"] != nil {
		d.Set("target", objectRaw["target"])
	}

	scheduledActions1Raw := objectRaw["scheduledActions"]
	scheduledActionsMaps := make([]map[string]interface{}, 0)
	if scheduledActions1Raw != nil {
		for _, scheduledActionsChild1Raw := range scheduledActions1Raw.([]interface{}) {
			scheduledActionsMap := make(map[string]interface{})
			scheduledActionsChild1Raw := scheduledActionsChild1Raw.(map[string]interface{})
			scheduledActionsMap["end_time"] = scheduledActionsChild1Raw["endTime"]
			scheduledActionsMap["name"] = scheduledActionsChild1Raw["name"]
			scheduledActionsMap["schedule_expression"] = scheduledActionsChild1Raw["scheduleExpression"]
			scheduledActionsMap["start_time"] = scheduledActionsChild1Raw["startTime"]
			scheduledActionsMap["target"] = scheduledActionsChild1Raw["target"]
			scheduledActionsMap["time_zone"] = scheduledActionsChild1Raw["timeZone"]

			scheduledActionsMaps = append(scheduledActionsMaps, scheduledActionsMap)
		}
	}
	if objectRaw["scheduledActions"] != nil {
		if err := d.Set("scheduled_actions", scheduledActionsMaps); err != nil {
			return err
		}
	}
	targetTrackingPolicies1Raw := objectRaw["targetTrackingPolicies"]
	targetTrackingPoliciesMaps := make([]map[string]interface{}, 0)
	if targetTrackingPolicies1Raw != nil {
		for _, targetTrackingPoliciesChild1Raw := range targetTrackingPolicies1Raw.([]interface{}) {
			targetTrackingPoliciesMap := make(map[string]interface{})
			targetTrackingPoliciesChild1Raw := targetTrackingPoliciesChild1Raw.(map[string]interface{})
			targetTrackingPoliciesMap["end_time"] = targetTrackingPoliciesChild1Raw["endTime"]
			targetTrackingPoliciesMap["max_capacity"] = targetTrackingPoliciesChild1Raw["maxCapacity"]
			targetTrackingPoliciesMap["metric_target"] = targetTrackingPoliciesChild1Raw["metricTarget"]
			targetTrackingPoliciesMap["metric_type"] = targetTrackingPoliciesChild1Raw["metricType"]
			targetTrackingPoliciesMap["min_capacity"] = targetTrackingPoliciesChild1Raw["minCapacity"]
			targetTrackingPoliciesMap["name"] = targetTrackingPoliciesChild1Raw["name"]
			targetTrackingPoliciesMap["start_time"] = targetTrackingPoliciesChild1Raw["startTime"]
			targetTrackingPoliciesMap["time_zone"] = targetTrackingPoliciesChild1Raw["timeZone"]

			targetTrackingPoliciesMaps = append(targetTrackingPoliciesMaps, targetTrackingPoliciesMap)
		}
	}
	if objectRaw["targetTrackingPolicies"] != nil {
		if err := d.Set("target_tracking_policies", targetTrackingPoliciesMaps); err != nil {
			return err
		}
	}

	d.Set("function_name", d.Id())

	return nil
}

func resourceAliCloudFcv3ProvisionConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var body map[string]interface{}
	update := false
	functionName := d.Id()
	action := fmt.Sprintf("/2023-03-30/functions/%s/provision-config", functionName)
	conn, err := client.NewFcv2Client()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})
	request["functionName"] = d.Id()

	if v, ok := d.GetOk("qualifier"); ok {
		query["qualifier"] = StringPointer(v.(string))
	}

	if d.HasChange("target_tracking_policies") {
		update = true
	}
	if v, ok := d.GetOk("target_tracking_policies"); ok || d.HasChange("target_tracking_policies") {
		targetTrackingPoliciesMaps := make([]interface{}, 0)
		for _, dataLoop := range v.([]interface{}) {
			dataLoopTmp := dataLoop.(map[string]interface{})
			dataLoopMap := make(map[string]interface{})
			dataLoopMap["endTime"] = dataLoopTmp["end_time"]
			dataLoopMap["name"] = dataLoopTmp["name"]
			dataLoopMap["startTime"] = dataLoopTmp["start_time"]
			dataLoopMap["maxCapacity"] = dataLoopTmp["max_capacity"]
			dataLoopMap["metricTarget"] = dataLoopTmp["metric_target"]
			dataLoopMap["metricType"] = dataLoopTmp["metric_type"]
			dataLoopMap["minCapacity"] = dataLoopTmp["min_capacity"]
			dataLoopMap["timeZone"] = dataLoopTmp["time_zone"]
			targetTrackingPoliciesMaps = append(targetTrackingPoliciesMaps, dataLoopMap)
		}
		request["targetTrackingPolicies"] = targetTrackingPoliciesMaps
	}

	if d.HasChange("scheduled_actions") {
		update = true
	}
	if v, ok := d.GetOk("scheduled_actions"); ok || d.HasChange("scheduled_actions") {
		scheduledActionsMaps := make([]interface{}, 0)
		for _, dataLoop1 := range v.([]interface{}) {
			dataLoop1Tmp := dataLoop1.(map[string]interface{})
			dataLoop1Map := make(map[string]interface{})
			dataLoop1Map["scheduleExpression"] = dataLoop1Tmp["schedule_expression"]
			dataLoop1Map["target"] = dataLoop1Tmp["target"]
			dataLoop1Map["endTime"] = dataLoop1Tmp["end_time"]
			dataLoop1Map["name"] = dataLoop1Tmp["name"]
			dataLoop1Map["startTime"] = dataLoop1Tmp["start_time"]
			dataLoop1Map["timeZone"] = dataLoop1Tmp["time_zone"]
			scheduledActionsMaps = append(scheduledActionsMaps, dataLoop1Map)
		}
		request["scheduledActions"] = scheduledActionsMaps
	}

	if d.HasChange("target") {
		update = true
	}
	if v, ok := d.GetOk("target"); ok || d.HasChange("target") {
		request["target"] = v
	}
	if d.HasChange("always_allocate_cpu") {
		update = true
	}
	if v, ok := d.GetOk("always_allocate_cpu"); ok || d.HasChange("always_allocate_cpu") {
		request["alwaysAllocateCPU"] = v
	}
	if d.HasChange("always_allocate_gpu") {
		update = true
	}
	if v, ok := d.GetOk("always_allocate_gpu"); ok || d.HasChange("always_allocate_gpu") {
		request["alwaysAllocateGPU"] = v
	}
	body = request
	if update {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer("2023-03-30"), nil, StringPointer("PUT"), StringPointer("AK"), StringPointer(action), query, nil, body, &runtime)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}

	return resourceAliCloudFcv3ProvisionConfigRead(d, meta)
}

func resourceAliCloudFcv3ProvisionConfigDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	functionName := d.Id()
	action := fmt.Sprintf("/2023-03-30/functions/%s/provision-config", functionName)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	conn, err := client.NewFcv2Client()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["functionName"] = d.Id()

	if v, ok := d.GetOk("qualifier"); ok {
		query["qualifier"] = StringPointer(v.(string))
	}

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer("2023-03-30"), nil, StringPointer("DELETE"), StringPointer("AK"), StringPointer(action), query, nil, nil, &runtime)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
