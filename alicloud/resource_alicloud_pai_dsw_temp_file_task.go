// Package alicloud. This file is hand-written for PaiDsw TempFileTask.
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

func resourceAliCloudPaiDswTempFileTask() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudPaiDswTempFileTaskCreate,
		Read:   resourceAliCloudPaiDswTempFileTaskRead,
		Update: resourceAliCloudPaiDswTempFileTaskUpdate,
		Delete: resourceAliCloudPaiDswTempFileTaskDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"gmt_expired_time": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"owner_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"temp_file_task_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"gmt_modified_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudPaiDswTempFileTaskCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "/api/v2/tempfiletasks"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["RegionId"] = StringPointer(client.RegionId)

	request["InstanceId"] = d.Get("instance_id")

	body = request
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RoaPost("pai-dsw", "2022-01-01", action, query, nil, body, true)
		if err != nil {
			if NeedRetry(err) || IsExpectedErrors(err, []string{"400"}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_pai_dsw_temp_file_task", action, AlibabaCloudSdkGoERROR)
	}

	id, _ := jsonpath.Get("$.TempFileTaskId", response)
	d.SetId(fmt.Sprint(id))

	// GmtExpiredTime is not a Create input; if the user set it, call Update to apply it.
	if v, ok := d.GetOk("gmt_expired_time"); ok && v.(string) != "" {
		if upErr := resourceAliCloudPaiDswTempFileTaskUpdate(d, meta); upErr != nil {
			return WrapErrorf(upErr, DefaultErrorMsg, "alicloud_pai_dsw_temp_file_task", "UpdateTempFileTask", AlibabaCloudSdkGoERROR)
		}
	}

	return resourceAliCloudPaiDswTempFileTaskRead(d, meta)
}

func resourceAliCloudPaiDswTempFileTaskRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	paiDswService := PaiDswService{client}

	objectRaw, err := paiDswService.DescribePaiDswTempFileTask(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_pai_dsw_temp_file_task DescribePaiDswTempFileTask Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	if objectRaw["Uuid"] != nil {
		d.Set("temp_file_task_id", objectRaw["Uuid"])
	} else {
		d.Set("temp_file_task_id", d.Id())
	}
	if objectRaw["UserId"] != nil {
		d.Set("user_id", objectRaw["UserId"])
	}
	if objectRaw["OwnerId"] != nil {
		d.Set("owner_id", objectRaw["OwnerId"])
	}
	if objectRaw["InstanceId"] != nil {
		d.Set("instance_id", objectRaw["InstanceId"])
	}
	if objectRaw["GmtCreateTime"] != nil {
		d.Set("create_time", objectRaw["GmtCreateTime"])
	}
	if objectRaw["GmtModifiedTime"] != nil {
		d.Set("gmt_modified_time", objectRaw["GmtModifiedTime"])
	}
	if objectRaw["GmtExpiredTime"] != nil {
		d.Set("gmt_expired_time", objectRaw["GmtExpiredTime"])
	}
	d.Set("region_id", client.RegionId)

	return nil
}

func resourceAliCloudPaiDswTempFileTaskUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error
	update := false

	tempFileTaskId := d.Id()
	action := fmt.Sprintf("/api/v2/tempfiletasks/%s", tempFileTaskId)
	request = make(map[string]interface{})
	query["RegionId"] = StringPointer(client.RegionId)

	// GmtExpiredTime is not a Create input; on Create (IsNewResource) send PUT
	// to apply the user-set expiry, on Update send PUT only when it changed.
	if v, ok := d.GetOk("gmt_expired_time"); ok {
		if d.IsNewResource() || d.HasChange("gmt_expired_time") {
			request["GmtExpiredTime"] = v.(string)
			update = true
		}
	}

	body = request
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RoaPut("pai-dsw", "2022-01-01", action, query, nil, body, true)
			if err != nil {
				if NeedRetry(err) || IsExpectedErrors(err, []string{"400"}) {
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

	return resourceAliCloudPaiDswTempFileTaskRead(d, meta)
}

func resourceAliCloudPaiDswTempFileTaskDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	tempFileTaskId := d.Id()
	action := fmt.Sprintf("/api/v2/tempfiletasks/%s", tempFileTaskId)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	var err error
	request = make(map[string]interface{})
	query["RegionId"] = StringPointer(client.RegionId)

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RoaDelete("pai-dsw", "2022-01-01", action, query, nil, nil, true)
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
		if IsExpectedErrors(err, []string{"ValidationError"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
