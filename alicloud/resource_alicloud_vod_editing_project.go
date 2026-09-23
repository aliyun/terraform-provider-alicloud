package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudVodEditingProject() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudVodEditingProjectCreate,
		Read:   resourceAlicloudVodEditingProjectRead,
		Update: resourceAlicloudVodEditingProjectUpdate,
		Delete: resourceAlicloudVodEditingProjectDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"cover_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"division": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"editing_project_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"timeline": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"title": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceAlicloudVodEditingProjectCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "AddEditingProject"
	request := make(map[string]interface{})
	var err error
	if v, ok := d.GetOk("cover_url"); ok {
		request["CoverURL"] = v
	}
	if v, ok := d.GetOk("division"); ok {
		request["Division"] = v
	}
	if v, ok := d.GetOk("editing_project_name"); ok {
		request["Description"] = v
	}
	if v, ok := d.GetOk("timeline"); ok {
		request["Timeline"] = v
	}
	request["Title"] = d.Get("title")
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = client.RpcPost("vod", "2017-03-21", action, nil, request, false)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_vod_editing_project", action, AlibabaCloudSdkGoERROR)
	}
	responseProject := response["Project"].(map[string]interface{})
	d.SetId(fmt.Sprint(responseProject["ProjectId"]))

	return resourceAlicloudVodEditingProjectRead(d, meta)
}
func resourceAlicloudVodEditingProjectRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vodService := VodService{client}
	object, err := vodService.DescribeVodEditingProject(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_vod_editing_project vodService.DescribeVodEditingProject Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("editing_project_name", object["Description"])
	d.Set("status", object["Status"])
	d.Set("timeline", object["Timeline"])
	d.Set("title", object["Title"])
	d.Set("cover_url", object["CoverURL"])
	return nil
}
func resourceAlicloudVodEditingProjectUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	var err error
	update := false
	request := map[string]interface{}{
		"ProjectId": d.Id(),
	}
	if d.HasChange("editing_project_name") {
		update = true
		if v, ok := d.GetOk("editing_project_name"); ok {
			request["Description"] = v
		}
	}
	if d.HasChange("timeline") {
		update = true
		if v, ok := d.GetOk("timeline"); ok {
			request["Timeline"] = v
		}
	}
	if d.HasChange("title") {
		update = true
		request["Title"] = d.Get("title")
	}
	if d.HasChange("cover_url") {
		update = true
		if v, ok := d.GetOk("cover_url"); ok {
			request["CoverURL"] = v
		}
	}
	if update {
		action := "UpdateEditingProject"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.RpcPost("vod", "2017-03-21", action, nil, request, false)
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
	return resourceAlicloudVodEditingProjectRead(d, meta)
}
func resourceAlicloudVodEditingProjectDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteEditingProject"
	var response map[string]interface{}
	var err error
	request := map[string]interface{}{
		"ProjectIds": d.Id(),
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = client.RpcPost("vod", "2017-03-21", action, nil, request, false)
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
