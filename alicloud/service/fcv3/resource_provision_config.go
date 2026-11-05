// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package fcv3

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/validators"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceProvisionConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceProvisionConfigCreate,
		Read:   resourceProvisionConfigRead,
		Update: resourceProvisionConfigUpdate,
		Delete: resourceProvisionConfigDelete,
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
				Computed: true,
			},
			"always_allocate_gpu": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"current": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"current_error": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"function_arn": {
				Type:     schema.TypeString,
				Computed: true,
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
				Computed:     true,
				ValidateFunc: validators.IntBetween(0, 10000),
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

func resourceProvisionConfigCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	functionName := d.Get("function_name")
	action := fmt.Sprintf("/2023-03-30/functions/%s/provision-config", functionName)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("function_name"); ok {
		request["functionName"] = v
	}

	if v, ok := d.GetOk("qualifier"); ok {
		query["qualifier"] = helper.StringPointer(v.(string))
	}

	if v, ok := d.GetOk("target_tracking_policies"); ok {
		targetTrackingPoliciesMapsArray := make([]interface{}, 0)
		for _, dataLoop := range v.([]interface{}) {
			dataLoopTmp, ok := dataLoop.(map[string]interface{})
			if !ok {
				continue
			}
			dataLoopMap := make(map[string]interface{})
			dataLoopMap["endTime"] = dataLoopTmp["end_time"]
			dataLoopMap["name"] = dataLoopTmp["name"]
			dataLoopMap["startTime"] = dataLoopTmp["start_time"]
			dataLoopMap["maxCapacity"] = dataLoopTmp["max_capacity"]
			dataLoopMap["metricTarget"] = dataLoopTmp["metric_target"]
			dataLoopMap["metricType"] = dataLoopTmp["metric_type"]
			dataLoopMap["minCapacity"] = dataLoopTmp["min_capacity"]
			dataLoopMap["timeZone"] = dataLoopTmp["time_zone"]
			targetTrackingPoliciesMapsArray = append(targetTrackingPoliciesMapsArray, dataLoopMap)
		}
		request["targetTrackingPolicies"] = targetTrackingPoliciesMapsArray
	}

	if v, ok := d.GetOk("scheduled_actions"); ok {
		scheduledActionsMapsArray := make([]interface{}, 0)
		for _, dataLoop1 := range v.([]interface{}) {
			dataLoop1Tmp, ok := dataLoop1.(map[string]interface{})
			if !ok {
				continue
			}
			dataLoop1Map := make(map[string]interface{})
			dataLoop1Map["scheduleExpression"] = dataLoop1Tmp["schedule_expression"]
			dataLoop1Map["target"] = dataLoop1Tmp["target"]
			dataLoop1Map["endTime"] = dataLoop1Tmp["end_time"]
			dataLoop1Map["name"] = dataLoop1Tmp["name"]
			dataLoop1Map["startTime"] = dataLoop1Tmp["start_time"]
			dataLoop1Map["timeZone"] = dataLoop1Tmp["time_zone"]
			scheduledActionsMapsArray = append(scheduledActionsMapsArray, dataLoop1Map)
		}
		request["scheduledActions"] = scheduledActionsMapsArray
	}

	if v, ok := d.GetOkExists("target"); ok {
		request["target"] = v
		// PutProvisionConfig deprecates target in favour of defaultTarget, so
		// send the same value through both fields.
		request["defaultTarget"] = v
	}
	if v, ok := d.GetOkExists("always_allocate_cpu"); ok {
		request["alwaysAllocateCPU"] = v
	}
	if v, ok := d.GetOkExists("always_allocate_gpu"); ok {
		request["alwaysAllocateGPU"] = v
	}
	body = request
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RoaPut("FC", "2023-03-30", action, query, nil, body, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_fcv3_provision_config", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(functionName))

	return resourceProvisionConfigRead(d, meta)
}

func resourceProvisionConfigRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	service := NewService(client)

	objectRaw, err := service.DescribeProvisionConfig(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_fcv3_provision_config DescribeProvisionConfig Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	if objectRaw["alwaysAllocateCPU"] != nil {
		d.Set("always_allocate_cpu", objectRaw["alwaysAllocateCPU"])
	}
	if objectRaw["alwaysAllocateGPU"] != nil {
		d.Set("always_allocate_gpu", objectRaw["alwaysAllocateGPU"])
	}
	if objectRaw["current"] != nil {
		d.Set("current", objectRaw["current"])
	}
	if objectRaw["currentError"] != nil {
		d.Set("current_error", objectRaw["currentError"])
	}
	if objectRaw["functionArn"] != nil {
		d.Set("function_arn", objectRaw["functionArn"])
	}
	// The service reports target as the effective target, which a scheduled
	// action or a target tracking policy can override while it is active, and
	// defaultTarget as the configured default of the provision configuration.
	// Map target from defaultTarget: the policy-computed runtime value must
	// never become the stored desired configuration, because the replacement
	// PUT performed by Update sends the value from state and would turn a
	// temporary policy target into the persistent default. A response that
	// omits defaultTarget is only trusted when no scaling policy is present,
	// where the effective target equals the configured default.
	if objectRaw["defaultTarget"] != nil {
		d.Set("target", objectRaw["defaultTarget"])
	} else if objectRaw["target"] != nil && !provisionConfigHasScalingPolicy(objectRaw) {
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

// provisionConfigHasScalingPolicy reports whether the response carries any
// scheduled action or target tracking policy. Without one the effective target
// equals defaultTarget, so target is a safe fallback when defaultTarget is not
// returned.
func provisionConfigHasScalingPolicy(objectRaw map[string]interface{}) bool {
	for _, key := range []string{"scheduledActions", "targetTrackingPolicies"} {
		if list, ok := objectRaw[key].([]interface{}); ok && len(list) > 0 {
			return true
		}
	}
	return false
}

// provisionConfigUpdateTarget returns the value the replacement PUT sends as
// the configured default. An explicitly configured or changed target,
// including zero, is honored from the configuration/plan. Otherwise state
// normally holds the configured default mapped by Read, so it is reused. A
// zero without a configured value means no default was ever established in
// state, for example an import whose response omitted defaultTarget while a
// scaling policy was active; in that case the service is queried for its
// configured default instead of persisting zero.
func provisionConfigUpdateTarget(d *schema.ResourceData, client *connectivity.AliyunClient) (interface{}, error) {
	target := d.Get("target")
	if d.HasChange("target") || provisionConfigTargetConfigured(d) {
		return target, nil
	}
	if targetValue, ok := target.(int); !ok || targetValue != 0 {
		return target, nil
	}

	service := NewService(client)
	object, err := service.DescribeProvisionConfig(d.Id())
	if err != nil {
		return nil, errs.WrapError(err)
	}
	switch v := object["defaultTarget"].(type) {
	case float64:
		return int(v), nil
	case int:
		return v, nil
	case nil:
	default:
		return nil, fmt.Errorf("unexpected defaultTarget type %T for provision config %q", v, d.Id())
	}
	if !provisionConfigHasScalingPolicy(object) {
		return target, nil
	}
	return nil, fmt.Errorf("cannot determine the configured default target of provision config %q: the service response omits defaultTarget while a scaling policy is active", d.Id())
}

// provisionConfigTargetConfigured reports whether the configuration sets
// target explicitly, so an explicitly configured zero is never replaced by the
// value reported by the service. The raw configuration is unavailable during
// import and in some test fixtures; in that case the caller falls back to the
// state value.
func provisionConfigTargetConfigured(d *schema.ResourceData) bool {
	raw := d.GetRawConfig()
	if raw.IsNull() || !raw.Type().IsObjectType() || !raw.Type().HasAttribute("target") {
		return false
	}
	value := raw.GetAttr("target")
	return value.IsKnown() && !value.IsNull()
}

func resourceProvisionConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var body map[string]interface{}
	update := false
	functionName := d.Id()
	action := fmt.Sprintf("/2023-03-30/functions/%s/provision-config", functionName)
	var err error
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})
	request["functionName"] = d.Id()

	if v, ok := d.GetOk("qualifier"); ok {
		query["qualifier"] = helper.StringPointer(v.(string))
	}

	if d.HasChange("target_tracking_policies") {
		update = true
	}
	if v, ok := d.GetOk("target_tracking_policies"); ok {
		targetTrackingPoliciesMapsArray := make([]interface{}, 0)
		for _, dataLoop := range v.([]interface{}) {
			dataLoopTmp, ok := dataLoop.(map[string]interface{})
			if !ok {
				continue
			}
			dataLoopMap := make(map[string]interface{})
			dataLoopMap["endTime"] = dataLoopTmp["end_time"]
			dataLoopMap["name"] = dataLoopTmp["name"]
			dataLoopMap["startTime"] = dataLoopTmp["start_time"]
			dataLoopMap["maxCapacity"] = dataLoopTmp["max_capacity"]
			dataLoopMap["metricTarget"] = dataLoopTmp["metric_target"]
			dataLoopMap["metricType"] = dataLoopTmp["metric_type"]
			dataLoopMap["minCapacity"] = dataLoopTmp["min_capacity"]
			dataLoopMap["timeZone"] = dataLoopTmp["time_zone"]
			targetTrackingPoliciesMapsArray = append(targetTrackingPoliciesMapsArray, dataLoopMap)
		}
		request["targetTrackingPolicies"] = targetTrackingPoliciesMapsArray
	} else if d.HasChange("target_tracking_policies") {
		request["targetTrackingPolicies"] = []interface{}{}
	}

	if d.HasChange("scheduled_actions") {
		update = true
	}
	if v, ok := d.GetOk("scheduled_actions"); ok {
		scheduledActionsMapsArray := make([]interface{}, 0)
		for _, dataLoop1 := range v.([]interface{}) {
			dataLoop1Tmp, ok := dataLoop1.(map[string]interface{})
			if !ok {
				continue
			}
			dataLoop1Map := make(map[string]interface{})
			dataLoop1Map["scheduleExpression"] = dataLoop1Tmp["schedule_expression"]
			dataLoop1Map["target"] = dataLoop1Tmp["target"]
			dataLoop1Map["endTime"] = dataLoop1Tmp["end_time"]
			dataLoop1Map["name"] = dataLoop1Tmp["name"]
			dataLoop1Map["startTime"] = dataLoop1Tmp["start_time"]
			dataLoop1Map["timeZone"] = dataLoop1Tmp["time_zone"]
			scheduledActionsMapsArray = append(scheduledActionsMapsArray, dataLoop1Map)
		}
		request["scheduledActions"] = scheduledActionsMapsArray
	} else if d.HasChange("scheduled_actions") {
		request["scheduledActions"] = []interface{}{}
	}

	if d.HasChange("target") {
		update = true
	}

	if d.HasChange("always_allocate_cpu") {
		update = true
	}

	if d.HasChange("always_allocate_gpu") {
		update = true
	}

	// PutProvisionConfig replaces the stored provision configuration and
	// requires target, so every update carries target plus the current
	// always-allocate values; otherwise a single-field update could drop
	// them or fall back to the service defaults. target is Optional+Computed:
	// when it is omitted or comes from an import, state holds the configured
	// default mapped from defaultTarget by Read, so the replacement PUT keeps
	// that default instead of writing zero or the policy-computed runtime
	// target. An explicitly configured value, including zero, is still sent.
	if update {
		targetValue, targetErr := provisionConfigUpdateTarget(d, client)
		if targetErr != nil {
			return targetErr
		}
		request["target"] = targetValue
		// target is deprecated but still required by PutProvisionConfig; send
		// the same configured default through defaultTarget as well.
		request["defaultTarget"] = targetValue
		request["alwaysAllocateCPU"] = d.Get("always_allocate_cpu")
		request["alwaysAllocateGPU"] = d.Get("always_allocate_gpu")
	}

	body = request
	if update {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RoaPut("FC", "2023-03-30", action, query, nil, body, true)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
	}

	return resourceProvisionConfigRead(d, meta)
}

func resourceProvisionConfigDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	functionName := d.Id()
	action := fmt.Sprintf("/2023-03-30/functions/%s/provision-config", functionName)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	var err error
	request = make(map[string]interface{})
	request["functionName"] = d.Id()

	if v, ok := d.GetOk("qualifier"); ok {
		query["qualifier"] = helper.StringPointer(v.(string))
	}

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RoaDelete("FC", "2023-03-30", action, query, nil, nil, true)

		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		if errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	service := NewService(client)
	deleteWait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	return retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		object, err := service.DescribeProvisionConfig(functionName)
		if err != nil {
			if errs.NotFoundError(err) {
				return nil
			}
			return retry.NonRetryableError(err)
		}
		if object == nil {
			return nil
		}
		deleteWait()
		return retry.RetryableError(errs.WrapErrorf(errs.Error("ProvisionConfig [%s] deletion is in progress", functionName), "please wait a moment and try again"))
	})
}
