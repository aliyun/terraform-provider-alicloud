// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudPaiWorkspaceRun() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudPaiWorkspaceRunCreate,
		Read:   resourceAliCloudPaiWorkspaceRunRead,
		Update: resourceAliCloudPaiWorkspaceRunUpdate,
		Delete: resourceAliCloudPaiWorkspaceRunDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"experiment_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"run_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"source_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"source_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudPaiWorkspaceRunCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := fmt.Sprintf("/api/v1/runs")
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	if v, ok := d.GetOk("source_type"); ok {
		request["SourceType"] = v
	}
	if v, ok := d.GetOk("source_id"); ok {
		request["SourceId"] = v
	}
	request["ExperimentId"] = d.Get("experiment_id")
	if v, ok := d.GetOk("run_name"); ok {
		request["Name"] = v
	}
	body = request
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RoaPost("AIWorkSpace", "2021-02-04", action, query, nil, body, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_pai_workspace_run", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["RunId"]))

	return resourceAliCloudPaiWorkspaceRunRead(d, meta)
}

func resourceAliCloudPaiWorkspaceRunRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	paiWorkspaceServiceV2 := PaiWorkspaceServiceV2{client}

	objectRaw, err := paiWorkspaceServiceV2.DescribePaiWorkspaceRun(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_pai_workspace_run DescribePaiWorkspaceRun Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	if objectRaw["GmtCreateTime"] != nil {
		d.Set("create_time", objectRaw["GmtCreateTime"])
	}
	if objectRaw["ExperimentId"] != nil {
		d.Set("experiment_id", objectRaw["ExperimentId"])
	}
	if objectRaw["Name"] != nil {
		d.Set("run_name", objectRaw["Name"])
	}
	if objectRaw["SourceId"] != nil {
		d.Set("source_id", objectRaw["SourceId"])
	}
	if objectRaw["SourceType"] != nil {
		d.Set("source_type", objectRaw["SourceType"])
	}

	return nil
}

func resourceAliCloudPaiWorkspaceRunUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var body map[string]interface{}
	update := false

	RunId := d.Id()
	action := fmt.Sprintf("/api/v1/runs/%s", RunId)
	var err error
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})
	request["RunId"] = d.Id()

	if d.HasChange("run_name") {
		update = true
	}
	if v, ok := d.GetOk("run_name"); ok || d.HasChange("run_name") {
		request["Name"] = v
	}
	body = request
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RoaPut("AIWorkSpace", "2021-02-04", action, query, nil, body, true)
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

	return resourceAliCloudPaiWorkspaceRunRead(d, meta)
}

func resourceAliCloudPaiWorkspaceRunDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	RunId := d.Id()
	action := fmt.Sprintf("/api/v1/runs/%s", RunId)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	var err error
	request = make(map[string]interface{})
	request["RunId"] = d.Id()

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RoaDelete("AIWorkSpace", "2021-02-04", action, query, nil, nil, true)

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
