// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudEsaWaitingRoomEvent() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudEsaWaitingRoomEventCreate,
		Read:   resourceAliCloudEsaWaitingRoomEventRead,
		Update: resourceAliCloudEsaWaitingRoomEventUpdate,
		Delete: resourceAliCloudEsaWaitingRoomEventDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"custom_page_html": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"disable_session_renewal_enable": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"end_time": {
				Type:     schema.TypeString,
				Required: true,
			},
			"json_response_enable": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"language": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"new_users_per_minute": {
				Type:     schema.TypeString,
				Required: true,
			},
			"pre_queue_enable": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"pre_queue_start_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"queuing_method": {
				Type:     schema.TypeString,
				Required: true,
			},
			"queuing_status_code": {
				Type:     schema.TypeString,
				Required: true,
			},
			"random_pre_queue_enable": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"session_duration": {
				Type:     schema.TypeString,
				Required: true,
			},
			"site_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"start_time": {
				Type:     schema.TypeString,
				Required: true,
			},
			"status": {
				Type:     schema.TypeString,
				Required: true,
			},
			"total_active_users": {
				Type:     schema.TypeString,
				Required: true,
			},
			"waiting_room_event_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"waiting_room_event_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"waiting_room_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"waiting_room_type": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceAliCloudEsaWaitingRoomEventCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateWaitingRoomEvent"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("site_id"); ok {
		request["SiteId"] = v
	}
	if v, ok := d.GetOk("waiting_room_id"); ok {
		request["WaitingRoomId"] = v
	}
	request["RegionId"] = client.RegionId

	request["QueuingMethod"] = d.Get("queuing_method")
	request["WaitingRoomType"] = d.Get("waiting_room_type")
	if v, ok := d.GetOk("custom_page_html"); ok {
		request["CustomPageHtml"] = v
	}
	if v, ok := d.GetOk("pre_queue_start_time"); ok {
		request["PreQueueStartTime"] = v
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	request["StartTime"] = d.Get("start_time")
	request["Enable"] = d.Get("status")
	request["QueuingStatusCode"] = d.Get("queuing_status_code")
	if v, ok := d.GetOk("language"); ok {
		request["Language"] = v
	}
	if v, ok := d.GetOk("pre_queue_enable"); ok {
		request["PreQueueEnable"] = v
	}
	request["SessionDuration"] = d.Get("session_duration")
	if v, ok := d.GetOk("disable_session_renewal_enable"); ok {
		request["DisableSessionRenewalEnable"] = v
	}
	if v, ok := d.GetOk("random_pre_queue_enable"); ok {
		request["RandomPreQueueEnable"] = v
	}
	if v, ok := d.GetOk("json_response_enable"); ok {
		request["JsonResponseEnable"] = v
	}
	request["TotalActiveUsers"] = d.Get("total_active_users")
	request["NewUsersPerMinute"] = d.Get("new_users_per_minute")
	request["Name"] = d.Get("waiting_room_event_name")
	request["EndTime"] = d.Get("end_time")
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("ESA", "2024-09-10", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_esa_waiting_room_event", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v:%v", request["SiteId"], request["WaitingRoomId"], response["WaitingRoomEventId"]))

	return resourceAliCloudEsaWaitingRoomEventRead(d, meta)
}

func resourceAliCloudEsaWaitingRoomEventRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	esaServiceV2 := EsaServiceV2{client}

	objectRaw, err := esaServiceV2.DescribeEsaWaitingRoomEvent(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_esa_waiting_room_event DescribeEsaWaitingRoomEvent Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("custom_page_html", objectRaw["CustomPageHtml"])
	d.Set("description", objectRaw["Description"])
	d.Set("disable_session_renewal_enable", objectRaw["DisableSessionRenewalEnable"])
	d.Set("end_time", objectRaw["EndTime"])
	d.Set("json_response_enable", objectRaw["JsonResponseEnable"])
	d.Set("language", objectRaw["Language"])
	d.Set("new_users_per_minute", objectRaw["NewUsersPerMinute"])
	d.Set("pre_queue_enable", objectRaw["PreQueueEnable"])
	d.Set("pre_queue_start_time", objectRaw["PreQueueStartTime"])
	d.Set("queuing_method", objectRaw["QueuingMethod"])
	d.Set("queuing_status_code", objectRaw["QueuingStatusCode"])
	d.Set("random_pre_queue_enable", objectRaw["RandomPreQueueEnable"])
	d.Set("session_duration", objectRaw["SessionDuration"])
	d.Set("start_time", objectRaw["StartTime"])
	d.Set("status", objectRaw["Enable"])
	d.Set("total_active_users", objectRaw["TotalActiveUsers"])
	d.Set("waiting_room_event_name", objectRaw["Name"])
	d.Set("waiting_room_type", objectRaw["WaitingRoomType"])
	d.Set("waiting_room_event_id", objectRaw["WaitingRoomEventId"])
	d.Set("waiting_room_id", objectRaw["WaitingRoomId"])

	parts := strings.Split(d.Id(), ":")
	d.Set("site_id", formatInt(parts[0]))

	return nil
}

func resourceAliCloudEsaWaitingRoomEventUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	parts := strings.Split(d.Id(), ":")
	action := "UpdateWaitingRoomEvent"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["SiteId"] = parts[0]
	request["WaitingRoomEventId"] = parts[2]
	request["RegionId"] = client.RegionId
	if d.HasChange("queuing_method") {
		update = true
	}
	request["QueuingMethod"] = d.Get("queuing_method")
	if d.HasChange("waiting_room_type") {
		update = true
	}
	request["WaitingRoomType"] = d.Get("waiting_room_type")
	if d.HasChange("custom_page_html") {
		update = true
		request["CustomPageHtml"] = d.Get("custom_page_html")
	}

	if d.HasChange("pre_queue_start_time") {
		update = true
		request["PreQueueStartTime"] = d.Get("pre_queue_start_time")
	}

	if d.HasChange("description") {
		update = true
		request["Description"] = d.Get("description")
	}

	if d.HasChange("start_time") {
		update = true
	}
	request["StartTime"] = d.Get("start_time")
	if d.HasChange("status") {
		update = true
	}
	request["Enable"] = d.Get("status")
	if d.HasChange("queuing_status_code") {
		update = true
	}
	request["QueuingStatusCode"] = d.Get("queuing_status_code")
	if d.HasChange("language") {
		update = true
		request["Language"] = d.Get("language")
	}

	if d.HasChange("pre_queue_enable") {
		update = true
		request["PreQueueEnable"] = d.Get("pre_queue_enable")
	}

	if d.HasChange("session_duration") {
		update = true
	}
	request["SessionDuration"] = d.Get("session_duration")
	if d.HasChange("disable_session_renewal_enable") {
		update = true
		request["DisableSessionRenewalEnable"] = d.Get("disable_session_renewal_enable")
	}

	if d.HasChange("random_pre_queue_enable") {
		update = true
		request["RandomPreQueueEnable"] = d.Get("random_pre_queue_enable")
	}

	if d.HasChange("json_response_enable") {
		update = true
		request["JsonResponseEnable"] = d.Get("json_response_enable")
	}

	if d.HasChange("total_active_users") {
		update = true
	}
	request["TotalActiveUsers"] = d.Get("total_active_users")
	if d.HasChange("new_users_per_minute") {
		update = true
	}
	request["NewUsersPerMinute"] = d.Get("new_users_per_minute")
	if d.HasChange("waiting_room_event_name") {
		update = true
	}
	request["Name"] = d.Get("waiting_room_event_name")
	if d.HasChange("end_time") {
		update = true
	}
	request["EndTime"] = d.Get("end_time")
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("ESA", "2024-09-10", action, query, request, true)
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

	return resourceAliCloudEsaWaitingRoomEventRead(d, meta)
}

func resourceAliCloudEsaWaitingRoomEventDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeleteWaitingRoomEvent"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["SiteId"] = parts[0]
	request["WaitingRoomEventId"] = parts[2]
	request["RegionId"] = client.RegionId

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("ESA", "2024-09-10", action, query, request, true)

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
