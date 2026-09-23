// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudAlikafkaScheduledScalingRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudAlikafkaScheduledScalingRuleCreate,
		Read:   resourceAliCloudAlikafkaScheduledScalingRuleRead,
		Update: resourceAliCloudAlikafkaScheduledScalingRuleUpdate,
		Delete: resourceAliCloudAlikafkaScheduledScalingRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"duration_minutes": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"enable": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"first_scheduled_time": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"repeat_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"reserved_pub_flow": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"reserved_sub_flow": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"rule_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"schedule_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"time_zone": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"weekly_types": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceAliCloudAlikafkaScheduledScalingRuleCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateScheduledScalingRule"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		request["InstanceId"] = v
	}
	if v, ok := d.GetOk("rule_name"); ok {
		request["RuleName"] = v
	}
	request["RegionId"] = client.RegionId

	request["ReservedPubFlow"] = d.Get("reserved_pub_flow")
	request["ReservedSubFlow"] = d.Get("reserved_sub_flow")
	request["ScheduleType"] = d.Get("schedule_type")
	request["DurationMinutes"] = d.Get("duration_minutes")
	request["TimeZone"] = d.Get("time_zone")
	request["FirstScheduledTime"] = d.Get("first_scheduled_time")
	if v, ok := d.GetOkExists("enable"); ok {
		request["Enable"] = v
	}
	if v, ok := d.GetOk("repeat_type"); ok {
		request["RepeatType"] = v
	}
	if v, ok := d.GetOk("weekly_types"); ok {
		weeklyTypesMapsArray := convertToInterfaceArray(v)

		weeklyTypesMapsJson, err := json.Marshal(weeklyTypesMapsArray)
		if err != nil {
			return WrapError(err)
		}
		request["WeeklyTypes"] = string(weeklyTypesMapsJson)
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("alikafka", "2019-09-16", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_alikafka_scheduled_scaling_rule", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", request["InstanceId"], request["RuleName"]))

	return resourceAliCloudAlikafkaScheduledScalingRuleRead(d, meta)
}

func resourceAliCloudAlikafkaScheduledScalingRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	alikafkaServiceV2 := AlikafkaServiceV2{client}

	objectRaw, err := alikafkaServiceV2.DescribeAlikafkaScheduledScalingRule(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_alikafka_scheduled_scaling_rule DescribeAlikafkaScheduledScalingRule Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("duration_minutes", objectRaw["DurationMinutes"])
	d.Set("enable", objectRaw["Enable"])
	d.Set("first_scheduled_time", objectRaw["FirstScheduledTime"])
	d.Set("repeat_type", objectRaw["RepeatType"])
	d.Set("reserved_pub_flow", objectRaw["ReservedPubFlow"])
	d.Set("reserved_sub_flow", objectRaw["ReservedSubFlow"])
	d.Set("schedule_type", objectRaw["ScheduleType"])
	d.Set("time_zone", objectRaw["TimeZone"])
	d.Set("rule_name", objectRaw["RuleName"])

	weeklyTypesRaw, _ := jsonpath.Get("$.WeeklyTypes.WeeklyTypes", objectRaw)
	d.Set("weekly_types", weeklyTypesRaw)

	parts := strings.Split(d.Id(), ":")
	d.Set("instance_id", parts[0])

	return nil
}

func resourceAliCloudAlikafkaScheduledScalingRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	parts := strings.Split(d.Id(), ":")
	action := "ModifyScheduledScalingRule"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["InstanceId"] = parts[0]
	request["RuleName"] = parts[1]
	request["RegionId"] = client.RegionId
	if d.HasChange("enable") {
		update = true

		if v, ok := d.GetOkExists("enable"); ok {
			request["Enable"] = v
		}
	}
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("alikafka", "2019-09-16", action, query, request, true)
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
	}

	return resourceAliCloudAlikafkaScheduledScalingRuleRead(d, meta)
}

func resourceAliCloudAlikafkaScheduledScalingRuleDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeleteScheduledScalingRule"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["InstanceId"] = parts[0]
	request["RuleName"] = parts[1]
	request["RegionId"] = client.RegionId

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("alikafka", "2019-09-16", action, query, request, true)
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

	return nil
}
