// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudEsaWaitingRoom() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudEsaWaitingRoomCreate,
		Read:   resourceAliCloudEsaWaitingRoomRead,
		Update: resourceAliCloudEsaWaitingRoomUpdate,
		Delete: resourceAliCloudEsaWaitingRoomDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"cookie_name": {
				Type:     schema.TypeString,
				Required: true,
			},
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
			"host_name_and_path": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"path": {
							Type:     schema.TypeString,
							Required: true,
						},
						"subdomain": {
							Type:     schema.TypeString,
							Required: true,
						},
						"domain": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
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
			"queue_all_enable": {
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
			"session_duration": {
				Type:     schema.TypeString,
				Required: true,
			},
			"site_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Required: true,
			},
			"total_active_users": {
				Type:     schema.TypeString,
				Required: true,
			},
			"waiting_room_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"waiting_room_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"waiting_room_type": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceAliCloudEsaWaitingRoomCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateWaitingRoom"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("site_id"); ok {
		request["SiteId"] = v
	}

	request["QueuingMethod"] = d.Get("queuing_method")
	request["WaitingRoomType"] = d.Get("waiting_room_type")
	if v, ok := d.GetOk("custom_page_html"); ok {
		request["CustomPageHtml"] = v
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	request["Enable"] = d.Get("status")
	if v, ok := d.GetOk("host_name_and_path"); ok {
		hostNameAndPathMapsArray := make([]interface{}, 0)
		for _, dataLoop := range convertToInterfaceArray(v) {
			dataLoopTmp := dataLoop.(map[string]interface{})
			dataLoopMap := make(map[string]interface{})
			dataLoopMap["Domain"] = dataLoopTmp["domain"]
			dataLoopMap["Subdomain"] = dataLoopTmp["subdomain"]
			dataLoopMap["Path"] = dataLoopTmp["path"]
			hostNameAndPathMapsArray = append(hostNameAndPathMapsArray, dataLoopMap)
		}
		hostNameAndPathMapsJson, err := json.Marshal(hostNameAndPathMapsArray)
		if err != nil {
			return WrapError(err)
		}
		request["HostNameAndPath"] = string(hostNameAndPathMapsJson)
	}

	request["QueuingStatusCode"] = d.Get("queuing_status_code")
	if v, ok := d.GetOk("language"); ok {
		request["Language"] = v
	}
	request["CookieName"] = d.Get("cookie_name")
	request["SessionDuration"] = d.Get("session_duration")
	if v, ok := d.GetOk("queue_all_enable"); ok {
		request["QueueAllEnable"] = v
	}
	if v, ok := d.GetOk("disable_session_renewal_enable"); ok {
		request["DisableSessionRenewalEnable"] = v
	}
	if v, ok := d.GetOk("json_response_enable"); ok {
		request["JsonResponseEnable"] = v
	}
	request["TotalActiveUsers"] = d.Get("total_active_users")
	request["NewUsersPerMinute"] = d.Get("new_users_per_minute")
	request["Name"] = d.Get("waiting_room_name")
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_esa_waiting_room", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", request["SiteId"], response["WaitingRoomId"]))

	return resourceAliCloudEsaWaitingRoomRead(d, meta)
}

func resourceAliCloudEsaWaitingRoomRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	esaServiceV2 := EsaServiceV2{client}

	objectRaw, err := esaServiceV2.DescribeEsaWaitingRoom(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_esa_waiting_room DescribeEsaWaitingRoom Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("cookie_name", objectRaw["CookieName"])
	d.Set("custom_page_html", objectRaw["CustomPageHtml"])
	d.Set("description", objectRaw["Description"])
	d.Set("disable_session_renewal_enable", objectRaw["DisableSessionRenewalEnable"])
	d.Set("json_response_enable", objectRaw["JsonResponseEnable"])
	d.Set("language", objectRaw["Language"])
	d.Set("new_users_per_minute", objectRaw["NewUsersPerMinute"])
	d.Set("queue_all_enable", objectRaw["QueueAllEnable"])
	d.Set("queuing_method", objectRaw["QueuingMethod"])
	d.Set("queuing_status_code", objectRaw["QueuingStatusCode"])
	d.Set("session_duration", objectRaw["SessionDuration"])
	d.Set("status", objectRaw["Enable"])
	d.Set("total_active_users", objectRaw["TotalActiveUsers"])
	d.Set("waiting_room_name", objectRaw["Name"])
	d.Set("waiting_room_type", objectRaw["WaitingRoomType"])
	d.Set("waiting_room_id", objectRaw["WaitingRoomId"])

	hostNameAndPathRaw := objectRaw["HostNameAndPath"]
	hostNameAndPathMaps := make([]map[string]interface{}, 0)
	if hostNameAndPathRaw != nil {
		for _, hostNameAndPathChildRaw := range convertToInterfaceArray(hostNameAndPathRaw) {
			hostNameAndPathMap := make(map[string]interface{})
			hostNameAndPathChildRaw := hostNameAndPathChildRaw.(map[string]interface{})
			hostNameAndPathMap["domain"] = hostNameAndPathChildRaw["Domain"]
			hostNameAndPathMap["path"] = hostNameAndPathChildRaw["Path"]
			hostNameAndPathMap["subdomain"] = hostNameAndPathChildRaw["Subdomain"]

			hostNameAndPathMaps = append(hostNameAndPathMaps, hostNameAndPathMap)
		}
	}
	if err := d.Set("host_name_and_path", hostNameAndPathMaps); err != nil {
		return err
	}

	parts := strings.Split(d.Id(), ":")
	d.Set("site_id", fmt.Sprintf(parts[0]))

	return nil
}

func resourceAliCloudEsaWaitingRoomUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	parts := strings.Split(d.Id(), ":")
	action := "UpdateWaitingRoom"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["SiteId"] = parts[0]
	request["WaitingRoomId"] = parts[1]

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

	if d.HasChange("description") {
		update = true
		request["Description"] = d.Get("description")
	}

	if d.HasChange("status") {
		update = true
	}
	request["Enable"] = d.Get("status")
	if d.HasChange("host_name_and_path") {
		update = true
	}
	if v, ok := d.GetOk("host_name_and_path"); ok || d.HasChange("host_name_and_path") {
		hostNameAndPathMapsArray := make([]interface{}, 0)
		for _, dataLoop := range convertToInterfaceArray(v) {
			dataLoopTmp := dataLoop.(map[string]interface{})
			dataLoopMap := make(map[string]interface{})
			dataLoopMap["Domain"] = dataLoopTmp["domain"]
			dataLoopMap["Subdomain"] = dataLoopTmp["subdomain"]
			dataLoopMap["Path"] = dataLoopTmp["path"]
			hostNameAndPathMapsArray = append(hostNameAndPathMapsArray, dataLoopMap)
		}
		hostNameAndPathMapsJson, err := json.Marshal(hostNameAndPathMapsArray)
		if err != nil {
			return WrapError(err)
		}
		request["HostNameAndPath"] = string(hostNameAndPathMapsJson)
	}

	if d.HasChange("queuing_status_code") {
		update = true
	}
	request["QueuingStatusCode"] = d.Get("queuing_status_code")
	if d.HasChange("language") {
		update = true
		request["Language"] = d.Get("language")
	}

	if d.HasChange("cookie_name") {
		update = true
	}
	request["CookieName"] = d.Get("cookie_name")
	if d.HasChange("session_duration") {
		update = true
	}
	request["SessionDuration"] = d.Get("session_duration")
	if d.HasChange("queue_all_enable") {
		update = true
		request["QueueAllEnable"] = d.Get("queue_all_enable")
	}

	if d.HasChange("disable_session_renewal_enable") {
		update = true
		request["DisableSessionRenewalEnable"] = d.Get("disable_session_renewal_enable")
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
	if d.HasChange("waiting_room_name") {
		update = true
	}
	request["Name"] = d.Get("waiting_room_name")
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

	return resourceAliCloudEsaWaitingRoomRead(d, meta)
}

func resourceAliCloudEsaWaitingRoomDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeleteWaitingRoom"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["SiteId"] = parts[0]
	request["WaitingRoomId"] = parts[1]

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
