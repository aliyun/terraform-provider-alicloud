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

func resourceAliCloudGpdbStreamingJob() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudGpdbStreamingJobCreate,
		Read:   resourceAliCloudGpdbStreamingJobRead,
		Update: resourceAliCloudGpdbStreamingJobUpdate,
		Delete: resourceAliCloudGpdbStreamingJobDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"account": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"consistency": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"db_instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"data_source_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"dest_columns": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"dest_database": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dest_schema": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dest_table": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"error_limit_count": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"fallback_offset": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"group_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"job_config": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"job_description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"job_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"job_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"match_columns": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"mode": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"password": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"src_columns": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"try_run": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"update_columns": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"write_mode": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAliCloudGpdbStreamingJobCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateStreamingJob"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["DBInstanceId"] = d.Get("db_instance_id")
	query["RegionId"] = client.RegionId

	request["DataSourceId"] = d.Get("data_source_id")
	if v, ok := d.GetOk("job_description"); ok {
		request["JobDescription"] = v
	}
	request["JobName"] = d.Get("job_name")
	if v, ok := d.GetOk("mode"); ok {
		request["Mode"] = v
	}
	if v, ok := d.GetOk("src_columns"); ok {
		srcColumnsMaps := v.([]interface{})
		srcColumnsMapsJson, err := json.Marshal(srcColumnsMaps)
		if err != nil {
			return WrapError(err)
		}
		request["SrcColumns"] = string(srcColumnsMapsJson)
	}

	if v, ok := d.GetOk("dest_columns"); ok {
		destColumnsMaps := v.([]interface{})
		destColumnsMapsJson, err := json.Marshal(destColumnsMaps)
		if err != nil {
			return WrapError(err)
		}
		request["DestColumns"] = string(destColumnsMapsJson)
	}

	if v, ok := d.GetOk("account"); ok {
		request["Account"] = v
	}
	if v, ok := d.GetOk("password"); ok {
		request["Password"] = v
	}
	if v, ok := d.GetOk("dest_database"); ok {
		request["DestDatabase"] = v
	}
	if v, ok := d.GetOk("dest_schema"); ok {
		request["DestSchema"] = v
	}
	if v, ok := d.GetOk("dest_table"); ok {
		request["DestTable"] = v
	}
	if v, ok := d.GetOk("write_mode"); ok {
		request["WriteMode"] = v
	}
	if v, ok := d.GetOk("job_config"); ok {
		request["JobConfig"] = v
	}
	if v, ok := d.GetOk("group_name"); ok {
		request["GroupName"] = v
	}
	if v, ok := d.GetOk("fallback_offset"); ok {
		request["FallbackOffset"] = v
	}
	if v, ok := d.GetOk("match_columns"); ok {
		matchColumnsMaps := v.([]interface{})
		matchColumnsMapsJson, err := json.Marshal(matchColumnsMaps)
		if err != nil {
			return WrapError(err)
		}
		request["MatchColumns"] = string(matchColumnsMapsJson)
	}

	if v, ok := d.GetOk("update_columns"); ok {
		updateColumnsMaps := v.([]interface{})
		updateColumnsMapsJson, err := json.Marshal(updateColumnsMaps)
		if err != nil {
			return WrapError(err)
		}
		request["UpdateColumns"] = string(updateColumnsMapsJson)
	}

	if v, ok := d.GetOk("consistency"); ok {
		request["Consistency"] = v
	}
	if v, ok := d.GetOk("error_limit_count"); ok {
		request["ErrorLimitCount"] = v
	}
	if v, ok := d.GetOkExists("try_run"); ok {
		request["TryRun"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("gpdb", "2016-05-03", action, query, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_gpdb_streaming_job", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", query["DBInstanceId"], response["JobId"]))

	return resourceAliCloudGpdbStreamingJobRead(d, meta)
}

func resourceAliCloudGpdbStreamingJobRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gpdbServiceV2 := GpdbServiceV2{client}

	objectRaw, err := gpdbServiceV2.DescribeGpdbStreamingJob(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_gpdb_streaming_job DescribeGpdbStreamingJob Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	if objectRaw["Account"] != nil {
		d.Set("account", objectRaw["Account"])
	}
	if objectRaw["Consistency"] != nil {
		d.Set("consistency", objectRaw["Consistency"])
	}
	if objectRaw["CreateTime"] != nil {
		d.Set("create_time", objectRaw["CreateTime"])
	}
	if objectRaw["DataSourceId"] != nil {
		d.Set("data_source_id", objectRaw["DataSourceId"])
	}
	if objectRaw["DestDatabase"] != nil {
		d.Set("dest_database", objectRaw["DestDatabase"])
	}
	if objectRaw["DestSchema"] != nil {
		d.Set("dest_schema", objectRaw["DestSchema"])
	}
	if objectRaw["DestTable"] != nil {
		d.Set("dest_table", objectRaw["DestTable"])
	}
	if objectRaw["ErrorLimitCount"] != nil {
		d.Set("error_limit_count", objectRaw["ErrorLimitCount"])
	}
	if objectRaw["FallbackOffset"] != nil {
		d.Set("fallback_offset", objectRaw["FallbackOffset"])
	}
	if objectRaw["GroupName"] != nil {
		d.Set("group_name", objectRaw["GroupName"])
	}
	if objectRaw["JobConfig"] != nil {
		d.Set("job_config", objectRaw["JobConfig"])
	}
	if objectRaw["JobDescription"] != nil {
		d.Set("job_description", objectRaw["JobDescription"])
	}
	if objectRaw["JobName"] != nil {
		d.Set("job_name", objectRaw["JobName"])
	}
	if objectRaw["Mode"] != nil {
		d.Set("mode", objectRaw["Mode"])
	}
	if objectRaw["Password"] != nil {
		d.Set("password", objectRaw["Password"])
	}
	if objectRaw["Status"] != nil {
		d.Set("status", objectRaw["Status"])
	}
	if objectRaw["WriteMode"] != nil {
		d.Set("write_mode", objectRaw["WriteMode"])
	}
	if objectRaw["JobId"] != nil {
		d.Set("job_id", objectRaw["JobId"])
	}

	destColumns1Raw := make([]interface{}, 0)
	if objectRaw["DestColumns"] != nil {
		destColumns1Raw = objectRaw["DestColumns"].([]interface{})
	}

	d.Set("dest_columns", destColumns1Raw)
	matchColumns1Raw := make([]interface{}, 0)
	if objectRaw["MatchColumns"] != nil {
		matchColumns1Raw = objectRaw["MatchColumns"].([]interface{})
	}

	d.Set("match_columns", matchColumns1Raw)
	srcColumns1Raw := make([]interface{}, 0)
	if objectRaw["SrcColumns"] != nil {
		srcColumns1Raw = objectRaw["SrcColumns"].([]interface{})
	}

	d.Set("src_columns", srcColumns1Raw)
	updateColumns1Raw := make([]interface{}, 0)
	if objectRaw["UpdateColumns"] != nil {
		updateColumns1Raw = objectRaw["UpdateColumns"].([]interface{})
	}

	d.Set("update_columns", updateColumns1Raw)

	parts := strings.Split(d.Id(), ":")
	d.Set("db_instance_id", parts[0])
	d.Set("job_id", parts[1])

	return nil
}

func resourceAliCloudGpdbStreamingJobUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	parts := strings.Split(d.Id(), ":")
	action := "ModifyStreamingJob"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["JobId"] = parts[1]
	query["DBInstanceId"] = parts[0]
	query["RegionId"] = client.RegionId
	if d.HasChange("job_description") {
		update = true
		request["JobDescription"] = d.Get("job_description")
	}

	if d.HasChange("account") {
		update = true
		request["Account"] = d.Get("account")
	}

	if d.HasChange("password") {
		update = true
		request["Password"] = d.Get("password")
	}

	if d.HasChange("write_mode") {
		update = true
		request["WriteMode"] = d.Get("write_mode")
	}

	if d.HasChange("job_config") {
		update = true
		request["JobConfig"] = d.Get("job_config")
	}

	if d.HasChange("group_name") {
		update = true
		request["GroupName"] = d.Get("group_name")
	}

	if d.HasChange("fallback_offset") {
		update = true
		request["FallbackOffset"] = d.Get("fallback_offset")
	}

	if d.HasChange("consistency") {
		update = true
		request["Consistency"] = d.Get("consistency")
	}

	if d.HasChange("error_limit_count") {
		update = true
		request["ErrorLimitCount"] = d.Get("error_limit_count")
	}

	if d.HasChange("src_columns") {
		update = true
		if v, ok := d.GetOk("src_columns"); ok || d.HasChange("src_columns") {
			srcColumnsMaps := v.([]interface{})
			srcColumnsMapsJson, err := json.Marshal(srcColumnsMaps)
			if err != nil {
				return WrapError(err)
			}
			request["SrcColumns"] = string(srcColumnsMapsJson)
		}
	}

	if d.HasChange("update_columns") {
		update = true
		if v, ok := d.GetOk("update_columns"); ok || d.HasChange("update_columns") {
			updateColumnsMaps := v.([]interface{})
			updateColumnsMapsJson, err := json.Marshal(updateColumnsMaps)
			if err != nil {
				return WrapError(err)
			}
			request["UpdateColumns"] = string(updateColumnsMapsJson)
		}
	}

	if d.HasChange("dest_columns") {
		update = true
		if v, ok := d.GetOk("dest_columns"); ok || d.HasChange("dest_columns") {
			destColumnsMaps := v.([]interface{})
			destColumnsMapsJson, err := json.Marshal(destColumnsMaps)
			if err != nil {
				return WrapError(err)
			}
			request["DestColumns"] = string(destColumnsMapsJson)
		}
	}

	if d.HasChange("match_columns") {
		update = true
		if v, ok := d.GetOk("match_columns"); ok || d.HasChange("match_columns") {
			matchColumnsMaps := v.([]interface{})
			matchColumnsMapsJson, err := json.Marshal(matchColumnsMaps)
			if err != nil {
				return WrapError(err)
			}
			request["MatchColumns"] = string(matchColumnsMapsJson)
		}
	}

	if v, ok := d.GetOkExists("try_run"); ok {
		request["TryRun"] = v
	}
	if d.HasChange("dest_schema") {
		update = true
		request["DestSchema"] = d.Get("dest_schema")
	}

	if d.HasChange("dest_database") {
		update = true
		request["DestDatabase"] = d.Get("dest_database")
	}

	if d.HasChange("dest_table") {
		update = true
		request["DestTable"] = d.Get("dest_table")
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("gpdb", "2016-05-03", action, query, request, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}

	return resourceAliCloudGpdbStreamingJobRead(d, meta)
}

func resourceAliCloudGpdbStreamingJobDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeleteStreamingJob"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["JobId"] = parts[1]
	query["DBInstanceId"] = parts[0]
	query["RegionId"] = client.RegionId

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("gpdb", "2016-05-03", action, query, request, true)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
