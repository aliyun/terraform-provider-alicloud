// Package alicloud. This file is hand-written from the Cms 2024-03-30 OpenAPI definition.
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

func resourceAliCloudCmsOncallSchedule() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudCmsOncallScheduleCreate,
		Read:   resourceAliCloudCmsOncallScheduleRead,
		Update: resourceAliCloudCmsOncallScheduleUpdate,
		Delete: resourceAliCloudCmsOncallScheduleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"oncall_schedule_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"oncall_schedule_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"rotations": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"active_days": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeInt},
						},
						"contacts": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"rotation_end_time": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"rotation_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"rotation_start_time": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"shift_length": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"shift_recurrence_frequency": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"start_date": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"time_zone": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"shift_robot_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"source": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"substitudes": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceAliCloudCmsOncallScheduleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "/oncallschedule"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	var err error
	request = make(map[string]interface{})

	request["oncallScheduleName"] = d.Get("oncall_schedule_name")
	if v, ok := d.GetOk("rotations"); ok {
		request["rotations"] = expandOncallScheduleRotations(v)
	}
	if v, ok := d.GetOk("shift_robot_id"); ok {
		request["shiftRobotId"] = v
	}
	if v, ok := d.GetOk("source"); ok {
		request["source"] = v
	}
	if v, ok := d.GetOk("substitudes"); ok {
		request["substitudes"] = expandOncallScheduleSubstitudes(v)
	}
	body := request

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RoaPost("Cms", "2024-03-30", action, query, nil, body, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cms_oncall_schedule", action, AlibabaCloudSdkGoERROR)
	}

	id, ok := response["oncallScheduleId"].(string)
	if !ok || id == "" {
		// Some Create responses do not echo the generated id; fall back to
		// listing by name to resolve the freshly created schedule.
		name := d.Get("oncall_schedule_name").(string)
		listAction := "/oncallSchedules"
		listQuery := make(map[string]*string)
		listQuery["oncallScheduleName"] = StringPointer(name)
		var listResponse map[string]interface{}
		wait2 := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
			listResponse, err = client.RoaGet("Cms", "2024-03-30", listAction, listQuery, nil, nil)
			if err != nil {
				if NeedRetry(err) {
					wait2()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(listAction, listResponse, listQuery)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "alicloud_cms_oncall_schedule", listAction, AlibabaCloudSdkGoERROR)
		}
		items, _ := jsonpath.Get("$.oncallSchedules[*]", listResponse)
		if list, ok := items.([]interface{}); ok && len(list) > 0 {
			if m, ok := list[0].(map[string]interface{}); ok {
				id = fmt.Sprint(m["oncallScheduleId"])
			}
		}
	}
	if id == "" {
		return WrapError(Error("create Cms OncallSchedule failed: oncallScheduleId not found in response"))
	}

	d.SetId(id)

	return resourceAliCloudCmsOncallScheduleRead(d, meta)
}

func resourceAliCloudCmsOncallScheduleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cmsServiceV2 := CmsServiceV2{client}

	objectRaw, err := cmsServiceV2.DescribeCmsOncallSchedule(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cms_oncall_schedule DescribeCmsOncallSchedule Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("oncall_schedule_id", objectRaw["oncallScheduleId"])
	d.Set("oncall_schedule_name", objectRaw["oncallScheduleName"])
	d.Set("shift_robot_id", objectRaw["shiftRobotId"])
	d.Set("source", objectRaw["source"])
	if err := d.Set("rotations", flattenOncallScheduleRotations(objectRaw["rotations"])); err != nil {
		return WrapError(err)
	}
	if err := d.Set("substitudes", flattenOncallScheduleSubstitudes(objectRaw["substitudes"])); err != nil {
		return WrapError(err)
	}

	return nil
}

func resourceAliCloudCmsOncallScheduleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := fmt.Sprintf("/oncallSchedule/%s", d.Id())
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	var err error
	request = make(map[string]interface{})

	if d.HasChange("oncall_schedule_name") {
		request["oncallScheduleName"] = d.Get("oncall_schedule_name")
	}
	if d.HasChange("rotations") {
		request["rotations"] = expandOncallScheduleRotations(d.Get("rotations"))
	}
	if d.HasChange("shift_robot_id") {
		request["shiftRobotId"] = d.Get("shift_robot_id")
	}
	if d.HasChange("source") {
		request["source"] = d.Get("source")
	}
	if d.HasChange("substitudes") {
		request["substitudes"] = expandOncallScheduleSubstitudes(d.Get("substitudes"))
	}
	body := request

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
		response, err = client.RoaPatch("Cms", "2024-03-30", action, query, nil, body, true)
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
			d.SetId("")
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return resourceAliCloudCmsOncallScheduleRead(d, meta)
}

func resourceAliCloudCmsOncallScheduleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	// DeleteOncallSchedules is a plural, query-style operation:
	// DELETE /oncallSchedule?oncallScheduleIds=["<id>"]
	action := "/oncallSchedule"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	var err error
	request = make(map[string]interface{})
	query["oncallScheduleIds"] = StringPointer(fmt.Sprintf(`["%s"]`, d.Id()))

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RoaDelete("Cms", "2024-03-30", action, query, nil, nil, true)
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
		if IsExpectedErrors(err, []string{"ResourceNotFound", "data_not_exist"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
