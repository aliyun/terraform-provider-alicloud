// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudRealtimeComputeSqlFile() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudRealtimeComputeSqlFileCreate,
		Read:   resourceAliCloudRealtimeComputeSqlFileRead,
		Update: resourceAliCloudRealtimeComputeSqlFileUpdate,
		Delete: resourceAliCloudRealtimeComputeSqlFileDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"batch_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"namespace": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"parent_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"session_cluster_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"sql_file_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sql_script": {
				Type:     schema.TypeString,
				Required: true,
			},
			"workspace": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudRealtimeComputeSqlFileCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	namespace := d.Get("namespace")
	action := fmt.Sprintf("/api/v2/namespaces/%s/sql-file", namespace)
	var request map[string]interface{}
	var response map[string]interface{}
	header := make(map[string]*string)
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	header["workspace"] = StringPointer(d.Get("workspace").(string))

	if v, ok := d.GetOk("batch_mode"); ok {
		request["batchMode"] = v
	}
	request["sqlScript"] = d.Get("sql_script")
	if v, ok := d.GetOk("description"); ok {
		request["description"] = v
	}
	if v, ok := d.GetOk("session_cluster_name"); ok {
		request["sessionClusterName"] = v
	}
	if v, ok := d.GetOk("parent_id"); ok {
		request["parentId"] = v
	} else {
		realtimeComputeServiceV2 := RealtimeComputeServiceV2{client}
		rootFolderId, err := realtimeComputeServiceV2.GetSqlFileRootFolderId(d.Get("workspace").(string), namespace.(string))
		if err != nil {
			return WrapError(err)
		}
		request["parentId"] = rootFolderId
	}
	request["name"] = d.Get("name")
	body = request
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RoaPost("ververica", "2022-07-18", action, query, header, body, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_realtime_compute_sql_file", action, AlibabaCloudSdkGoERROR)
	}

	dataworkspaceVar, _ := jsonpath.Get("$.data.workspace", response)
	datanamespaceVar, _ := jsonpath.Get("$.data.namespace", response)
	datasqlFileIdVar, _ := jsonpath.Get("$.data.sqlFileId", response)
	d.SetId(fmt.Sprintf("%v:%v:%v", dataworkspaceVar, datanamespaceVar, datasqlFileIdVar))

	return resourceAliCloudRealtimeComputeSqlFileRead(d, meta)
}

func resourceAliCloudRealtimeComputeSqlFileRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	realtimeComputeServiceV2 := RealtimeComputeServiceV2{client}

	objectRaw, err := realtimeComputeServiceV2.DescribeRealtimeComputeSqlFile(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_realtime_compute_sql_file DescribeRealtimeComputeSqlFile Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	if v, ok := objectRaw["batchMode"]; ok && v != nil {
		d.Set("batch_mode", fmt.Sprint(v))
	}
	d.Set("description", objectRaw["description"])
	d.Set("name", objectRaw["name"])
	d.Set("parent_id", objectRaw["parentId"])
	d.Set("session_cluster_name", objectRaw["sessionClusterName"])
	d.Set("sql_script", objectRaw["sqlScript"])
	d.Set("namespace", objectRaw["namespace"])
	d.Set("sql_file_id", objectRaw["sqlFileId"])
	d.Set("workspace", objectRaw["workspace"])

	return nil
}

func resourceAliCloudRealtimeComputeSqlFileUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var header map[string]*string
	var body map[string]interface{}
	update := false

	var err error
	parts := strings.Split(d.Id(), ":")
	namespace := parts[1]
	sqlFileId := parts[2]
	action := fmt.Sprintf("/api/v2/namespaces/%s/sql-file/%s", namespace, sqlFileId)
	request = make(map[string]interface{})
	query = make(map[string]*string)
	header = make(map[string]*string)
	body = make(map[string]interface{})
	header["workspace"] = StringPointer(parts[0])

	if d.HasChange("batch_mode") {
		update = true
	}
	if v, ok := d.GetOk("batch_mode"); ok || d.HasChange("batch_mode") {
		request["batchMode"] = v
	}
	if d.HasChange("sql_script") {
		update = true
	}
	request["sqlScript"] = d.Get("sql_script")
	if d.HasChange("description") {
		update = true
	}
	if v, ok := d.GetOk("description"); ok || d.HasChange("description") {
		request["description"] = v
	}
	if d.HasChange("session_cluster_name") {
		update = true
	}
	if v, ok := d.GetOk("session_cluster_name"); ok || d.HasChange("session_cluster_name") {
		request["sessionClusterName"] = v
	}
	if d.HasChange("parent_id") {
		update = true
	}
	if v, ok := d.GetOk("parent_id"); ok || d.HasChange("parent_id") {
		request["parentId"] = v
	}
	if d.HasChange("name") {
		update = true
	}
	request["name"] = d.Get("name")
	body = request
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RoaPatch("ververica", "2022-07-18", action, query, header, body, true)
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

	return resourceAliCloudRealtimeComputeSqlFileRead(d, meta)
}

func resourceAliCloudRealtimeComputeSqlFileDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	namespace := parts[1]
	sqlFileId := parts[2]
	action := fmt.Sprintf("/api/v2/namespaces/%s/sql-file/%s", namespace, sqlFileId)
	var request map[string]interface{}
	var response map[string]interface{}
	header := make(map[string]*string)
	query := make(map[string]*string)
	var err error
	request = make(map[string]interface{})
	header["workspace"] = StringPointer(parts[0])

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RoaDelete("ververica", "2022-07-18", action, query, header, nil, true)
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
		if IsExpectedErrors(err, []string{"990301"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
