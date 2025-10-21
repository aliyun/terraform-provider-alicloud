// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudCloudMonitorServiceSiteMonitor() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudCloudMonitorServiceSiteMonitorCreate,
		Read:   resourceAliCloudCloudMonitorServiceSiteMonitorRead,
		Update: resourceAliCloudCloudMonitorServiceSiteMonitorUpdate,
		Delete: resourceAliCloudCloudMonitorServiceSiteMonitorDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"address": {
				Type:     schema.TypeString,
				Required: true,
			},
			"agent_group": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"custom_schedule": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"time_zone": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"start_hour": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"days": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeInt},
						},
						"end_hour": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
			"interval": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"1", "5", "15", "30", "60"}, false),
			},
			"isp_cities": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"isp": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"city": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"task_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"task_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"options_json": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					equal, _ := compareJsonTemplateAreEquivalent(old, new)
					return equal
				},
			},
			"task_state": {
				Type:       schema.TypeString,
				Computed:   true,
				Deprecated: "Field `task_state` has been deprecated from provider version 1.262.0. New field `status` instead.",
			},
			"create_time": {
				Type:       schema.TypeString,
				Computed:   true,
				Deprecated: "Field `create_time` has been deprecated from provider version 1.262.0.",
			},
			"update_time": {
				Type:       schema.TypeString,
				Computed:   true,
				Deprecated: "Field `update_time` has been deprecated from provider version 1.262.0.",
			},
			"alert_ids": {
				Type:       schema.TypeList,
				Optional:   true,
				Elem:       &schema.Schema{Type: schema.TypeString},
				Deprecated: "Field `alert_ids` has been deprecated from provider version 1.262.0.",
			},
		},
	}
}

func resourceAliCloudCloudMonitorServiceSiteMonitorCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateSiteMonitor"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	request["TaskType"] = d.Get("task_type")

	if v, ok := d.GetOk("options_json"); ok {
		request["OptionsJson"] = v
	}

	dataList1 := make(map[string]interface{})

	if v := d.Get("custom_schedule"); !IsNil(v) {
		days1, _ := jsonpath.Get("$[0].days", v)
		if days1 != nil && days1 != "" {
			dataList1["days"] = days1
		}
		startHour, _ := jsonpath.Get("$[0].start_hour", v)
		if startHour != nil && startHour != "" {
			dataList1["start_hour"] = startHour
		}
		endHour, _ := jsonpath.Get("$[0].end_hour", v)
		if endHour != nil && endHour != "" {
			dataList1["end_hour"] = endHour
		}
		timeZone, _ := jsonpath.Get("$[0].time_zone", v)
		if timeZone != nil && timeZone != "" {
			dataList1["time_zone"] = timeZone
		}

		customScheduleJson, err := convertMaptoJsonString(dataList1)
		if err != nil {
			return WrapError(err)
		}

		request["CustomSchedule"] = customScheduleJson
	}

	if v, ok := d.GetOk("interval"); ok {
		request["Interval"] = v
	}

	if v, ok := d.GetOk("isp_cities"); ok {
		ispCitiesMapsArray := make([]interface{}, 0)
		for _, dataLoop1 := range convertToInterfaceArray(v) {
			dataLoop1Tmp := dataLoop1.(map[string]interface{})
			dataLoop1Map := make(map[string]interface{})
			dataLoop1Map["Type"] = dataLoop1Tmp["type"]
			dataLoop1Map["Isp"] = dataLoop1Tmp["isp"]
			dataLoop1Map["City"] = dataLoop1Tmp["city"]
			ispCitiesMapsArray = append(ispCitiesMapsArray, dataLoop1Map)
		}

		ispCitiesJson, err := convertInterfaceToJsonString(ispCitiesMapsArray)
		if err != nil {
			return WrapError(err)
		}

		request["IspCities"] = ispCitiesJson
	}

	request["TaskName"] = d.Get("task_name")
	request["Address"] = d.Get("address")
	if v, ok := d.GetOk("agent_group"); ok {
		request["AgentGroup"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Cms", "2019-01-01", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cms_site_monitor", action, AlibabaCloudSdkGoERROR)
	}

	id, _ := jsonpath.Get("$.CreateResultList.CreateResultList[0].TaskId", response)
	d.SetId(fmt.Sprint(id))

	return resourceAliCloudCloudMonitorServiceSiteMonitorUpdate(d, meta)
}

func resourceAliCloudCloudMonitorServiceSiteMonitorRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudMonitorServiceServiceV2 := CloudMonitorServiceServiceV2{client}

	objectRaw, err := cloudMonitorServiceServiceV2.DescribeCloudMonitorServiceSiteMonitor(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cms_site_monitor DescribeCloudMonitorServiceSiteMonitor Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("address", objectRaw["Address"])
	d.Set("agent_group", objectRaw["AgentGroup"])
	d.Set("interval", objectRaw["Interval"])
	d.Set("status", fmt.Sprint(objectRaw["TaskState"]))
	d.Set("task_name", objectRaw["TaskName"])
	d.Set("task_type", objectRaw["TaskType"])
	d.Set("task_state", fmt.Sprint(objectRaw["TaskState"]))

	customScheduleMaps := make([]map[string]interface{}, 0)
	customScheduleMap := make(map[string]interface{})
	customScheduleRaw := make(map[string]interface{})
	if objectRaw["CustomSchedule"] != nil {
		customScheduleRaw = objectRaw["CustomSchedule"].(map[string]interface{})
	}
	if len(customScheduleRaw) > 0 {
		customScheduleMap["end_hour"] = customScheduleRaw["end_hour"]
		customScheduleMap["start_hour"] = customScheduleRaw["start_hour"]
		customScheduleMap["time_zone"] = customScheduleRaw["time_zone"]

		daysRaw, _ := jsonpath.Get("$.CustomSchedule.days.days", objectRaw)
		customScheduleMap["days"] = daysRaw
		customScheduleMaps = append(customScheduleMaps, customScheduleMap)
	}
	if err := d.Set("custom_schedule", customScheduleMaps); err != nil {
		return err
	}

	ispCityRaw, _ := jsonpath.Get("$.IspCities.IspCity", objectRaw)
	ispCitiesMaps := make([]map[string]interface{}, 0)
	if ispCityRaw != nil {
		for _, ispCityChildRaw := range convertToInterfaceArray(ispCityRaw) {
			ispCitiesMap := make(map[string]interface{})
			ispCityChildRaw := ispCityChildRaw.(map[string]interface{})
			ispCitiesMap["city"] = ispCityChildRaw["City"]
			ispCitiesMap["isp"] = ispCityChildRaw["Isp"]
			ispCitiesMap["type"] = ispCityChildRaw["Type"]

			ispCitiesMaps = append(ispCitiesMaps, ispCitiesMap)
		}
	}
	if err := d.Set("isp_cities", ispCitiesMaps); err != nil {
		return err
	}

	if objectRaw["OptionJson"] != nil {
		optionJsonJson, err := convertInterfaceToJsonString(objectRaw["OptionJson"])
		if err != nil {
			return WrapError(err)
		}

		d.Set("options_json", optionJsonJson)
	}

	return nil
}

func resourceAliCloudCloudMonitorServiceSiteMonitorUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	cloudMonitorServiceServiceV2 := CloudMonitorServiceServiceV2{client}
	objectRaw, _ := cloudMonitorServiceServiceV2.DescribeCloudMonitorServiceSiteMonitor(d.Id())

	if d.HasChange("status") {
		var err error
		target := d.Get("status").(string)
		if fmt.Sprint(objectRaw["TaskState"]) != target {
			if target == "1" {
				action := "EnableSiteMonitors"
				request = make(map[string]interface{})
				query = make(map[string]interface{})
				request["TaskIds"] = d.Id()

				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = client.RpcPost("Cms", "2019-01-01", action, query, request, true)
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
			if target == "2" {
				action := "DisableSiteMonitors"
				request = make(map[string]interface{})
				query = make(map[string]interface{})
				request["TaskIds"] = d.Id()

				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = client.RpcPost("Cms", "2019-01-01", action, query, request, true)
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
		}
	}

	var err error
	action := "ModifySiteMonitor"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["TaskId"] = d.Id()

	if !d.IsNewResource() && d.HasChange("options_json") {
		update = true

		if v, ok := d.GetOk("options_json"); ok {
			request["OptionsJson"] = v
		}
	}

	if !d.IsNewResource() && d.HasChange("custom_schedule") {
		update = true
		dataList1 := make(map[string]interface{})

		if v := d.Get("custom_schedule"); v != nil {
			days1, _ := jsonpath.Get("$[0].days", v)
			if days1 != nil && (d.HasChange("custom_schedule.0.days") || days1 != "") {
				dataList1["days"] = days1
			}
			startHour, _ := jsonpath.Get("$[0].start_hour", v)
			if startHour != nil && (d.HasChange("custom_schedule.0.start_hour") || startHour != "") {
				dataList1["start_hour"] = startHour
			}
			endHour, _ := jsonpath.Get("$[0].end_hour", v)
			if endHour != nil && (d.HasChange("custom_schedule.0.end_hour") || endHour != "") {
				dataList1["end_hour"] = endHour
			}
			timeZone, _ := jsonpath.Get("$[0].time_zone", v)
			if timeZone != nil && (d.HasChange("custom_schedule.0.time_zone") || timeZone != "") {
				dataList1["time_zone"] = timeZone
			}

			customScheduleJson, err := convertMaptoJsonString(dataList1)
			if err != nil {
				return WrapError(err)
			}

			request["CustomSchedule"] = customScheduleJson
		}
	}

	if !d.IsNewResource() && d.HasChange("interval") {
		update = true
		request["Interval"] = d.Get("interval")
	}

	if !d.IsNewResource() && d.HasChange("isp_cities") {
		update = true
		if v, ok := d.GetOk("isp_cities"); ok || d.HasChange("isp_cities") {
			ispCitiesMapsArray := make([]interface{}, 0)
			for _, dataLoop1 := range convertToInterfaceArray(v) {
				dataLoop1Tmp := dataLoop1.(map[string]interface{})
				dataLoop1Map := make(map[string]interface{})
				dataLoop1Map["Type"] = dataLoop1Tmp["type"]
				dataLoop1Map["isp"] = dataLoop1Tmp["isp"]
				dataLoop1Map["city"] = dataLoop1Tmp["city"]
				ispCitiesMapsArray = append(ispCitiesMapsArray, dataLoop1Map)
			}

			ispCitiesJson, err := convertInterfaceToJsonString(ispCitiesMapsArray)
			if err != nil {
				return WrapError(err)
			}

			request["IspCities"] = ispCitiesJson
		}
	}

	if !d.IsNewResource() && d.HasChange("task_name") {
		update = true
	}
	request["TaskName"] = d.Get("task_name")

	if !d.IsNewResource() && d.HasChange("address") {
		update = true
	}
	request["Address"] = d.Get("address")
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Cms", "2019-01-01", action, query, request, true)
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

	return resourceAliCloudCloudMonitorServiceSiteMonitorRead(d, meta)
}

func resourceAliCloudCloudMonitorServiceSiteMonitorDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteSiteMonitors"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["TaskIds"] = d.Id()

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Cms", "2019-01-01", action, query, request, true)

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
		if IsExpectedErrors(err, []string{"ResourceNotFound"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
