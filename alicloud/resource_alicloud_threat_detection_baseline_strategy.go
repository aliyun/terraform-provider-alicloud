package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudThreatDetectionBaselineStrategy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudThreatDetectionBaselineStrategyCreate,
		Read:   resourceAlicloudThreatDetectionBaselineStrategyRead,
		Update: resourceAlicloudThreatDetectionBaselineStrategyUpdate,
		Delete: resourceAlicloudThreatDetectionBaselineStrategyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"baseline_strategy_id": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"baseline_strategy_name": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"custom_type": {
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"common", "custom"}, false),
				Type:         schema.TypeString,
			},
			"cycle_days": {
				Required:     true,
				ValidateFunc: validation.IntInSlice([]int{1, 3, 7, 30}),
				Type:         schema.TypeInt,
			},
			"cycle_start_time": {
				Optional:     true,
				Computed:     true,
				Type:         schema.TypeInt,
				ValidateFunc: validation.IntInSlice([]int{0, 6, 12, 18}),
			},
			"end_time": {
				Required: true,
				Type:     schema.TypeString,
			},
			"risk_sub_type_name": {
				Required: true,
				Type:     schema.TypeString,
			},
			"start_time": {
				Required: true,
				Type:     schema.TypeString,
			},
			"target_type": {
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"groupId", "uuid"}, false),
			},
		},
	}
}

func resourceAlicloudThreatDetectionBaselineStrategyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := make(map[string]interface{})
	var err error

	if v, ok := d.GetOk("baseline_strategy_name"); ok {
		request["Name"] = v
	}
	if v, ok := d.GetOk("custom_type"); ok {
		request["CustomType"] = v
	}
	if v, ok := d.GetOk("cycle_days"); ok {
		request["CycleDays"] = v
	}
	if v, ok := d.GetOk("cycle_start_time"); ok {
		request["CycleStartTime"] = v
	}
	if v, ok := d.GetOk("id"); ok {
		request["Id"] = v
	}
	if v, ok := d.GetOk("risk_sub_type_name"); ok {
		request["RiskSubTypeName"] = v
	}
	if v, ok := d.GetOk("target_type"); ok {
		request["TargetType"] = v
	}
	if v, ok := d.GetOk("start_time"); ok {
		request["StartTime"] = v
	}
	if v, ok := d.GetOk("end_time"); ok {
		request["EndTime"] = v
	}
	request["RiskCustomParams"] = "[]"

	var response map[string]interface{}
	action := "ModifyStrategy"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		resp, err := client.RpcPost("Sas", "2018-12-03", action, nil, request, false)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		response = resp
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_threat_detection_baseline_strategy", action, AlibabaCloudSdkGoERROR)
	}

	if v, err := jsonpath.Get("$.Result.StrategyId", response); err != nil || v == nil {
		return WrapErrorf(err, IdMsg, "alicloud_threat_detection_baseline_strategy")
	} else {
		d.SetId(fmt.Sprint(v))
	}

	return resourceAlicloudThreatDetectionBaselineStrategyRead(d, meta)
}

func resourceAlicloudThreatDetectionBaselineStrategyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	sasService := SasService{client}
	objectFromList, err := sasService.DescribeStrategy(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_threat_detection_baseline_strategy sasService.DescribeThreatDetectionBaselineStrategy Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	object, err := sasService.DescribeThreatDetectionBaselineStrategy(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_threat_detection_baseline_strategy sasService.DescribeThreatDetectionBaselineStrategy Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("baseline_strategy_name", object["Name"])
	d.Set("custom_type", objectFromList["CustomType"])
	d.Set("cycle_days", object["CycleDays"])
	d.Set("cycle_start_time", object["CycleStartTime"])
	d.Set("end_time", object["EndTime"])
	//d.Set("risk_sub_type_name", object["RiskSubTypeName"])
	d.Set("start_time", object["StartTime"])
	d.Set("target_type", object["TargetType"])

	return nil
}

func resourceAlicloudThreatDetectionBaselineStrategyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	var err error
	update := false
	request := map[string]interface{}{
		"Id": d.Id(),
	}
	request["RiskCustomParams"] = "[]"

	if !d.IsNewResource() && d.HasChange("baseline_strategy_name") {
		update = true
	}
	if v, ok := d.GetOk("baseline_strategy_name"); ok {
		request["Name"] = v
	}
	if !d.IsNewResource() && d.HasChange("custom_type") {
		update = true
	}
	if v, ok := d.GetOk("custom_type"); ok {
		request["CustomType"] = v
	}
	if !d.IsNewResource() && d.HasChange("cycle_days") {
		update = true
	}
	if v, ok := d.GetOk("cycle_days"); ok {
		request["CycleDays"] = v
	}
	if !d.IsNewResource() && d.HasChange("cycle_start_time") {
		update = true
	}
	if v, ok := d.GetOk("cycle_start_time"); ok {
		request["CycleStartTime"] = v
	}
	if !d.IsNewResource() && d.HasChange("end_time") {
		update = true
	}
	if v, ok := d.GetOk("end_time"); ok {
		request["EndTime"] = v
	}
	if !d.IsNewResource() && d.HasChange("risk_sub_type_name") {
		update = true
	}
	if v, ok := d.GetOk("risk_sub_type_name"); ok {
		request["RiskSubTypeName"] = v
	}
	if !d.IsNewResource() && d.HasChange("start_time") {
		update = true
	}
	if v, ok := d.GetOk("start_time"); ok {
		request["StartTime"] = v
	}
	if !d.IsNewResource() && d.HasChange("target_type") {
		update = true
	}
	if v, ok := d.GetOk("target_type"); ok {
		request["TargetType"] = v
	}

	if update {
		action := "ModifyStrategy"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			resp, err := client.RpcPost("Sas", "2018-12-03", action, nil, request, false)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, resp, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}
	return resourceAlicloudThreatDetectionBaselineStrategyRead(d, meta)
}

func resourceAlicloudThreatDetectionBaselineStrategyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var err error

	request := map[string]interface{}{
		"Id": d.Id(),
	}

	action := "DeleteStrategy"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		resp, err := client.RpcPost("Sas", "2018-12-03", action, nil, request, false)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, resp, request)
		return nil
	})
	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}
