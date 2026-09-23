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

func resourceAliCloudGpdbStreamingDataSource() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudGpdbStreamingDataSourceCreate,
		Read:   resourceAliCloudGpdbStreamingDataSourceRead,
		Update: resourceAliCloudGpdbStreamingDataSourceUpdate,
		Delete: resourceAliCloudGpdbStreamingDataSourceDelete,
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
			"db_instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"data_source_config": {
				Type:     schema.TypeString,
				Required: true,
			},
			"data_source_description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"data_source_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"data_source_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"data_source_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"service_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudGpdbStreamingDataSourceCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateStreamingDataSource"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["DBInstanceId"] = d.Get("db_instance_id")
	request["RegionId"] = client.RegionId

	request["DataSourceName"] = d.Get("data_source_name")
	request["DataSourceConfig"] = d.Get("data_source_config")
	request["ServiceId"] = d.Get("service_id")
	request["DataSourceType"] = d.Get("data_source_type")
	if v, ok := d.GetOk("data_source_description"); ok {
		request["DataSourceDescription"] = v
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_gpdb_streaming_data_source", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", query["DBInstanceId"], response["DataSourceId"]))

	return resourceAliCloudGpdbStreamingDataSourceRead(d, meta)
}

func resourceAliCloudGpdbStreamingDataSourceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gpdbServiceV2 := GpdbServiceV2{client}

	objectRaw, err := gpdbServiceV2.DescribeGpdbStreamingDataSource(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_gpdb_streaming_data_source DescribeGpdbStreamingDataSource Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	if objectRaw["CreateTime"] != nil {
		d.Set("create_time", objectRaw["CreateTime"])
	}
	if objectRaw["DataSourceConfig"] != nil {
		d.Set("data_source_config", objectRaw["DataSourceConfig"])
	}
	if objectRaw["DataSourceDescription"] != nil {
		d.Set("data_source_description", objectRaw["DataSourceDescription"])
	}
	if objectRaw["DataSourceName"] != nil {
		d.Set("data_source_name", objectRaw["DataSourceName"])
	}
	if objectRaw["DataSourceType"] != nil {
		d.Set("data_source_type", objectRaw["DataSourceType"])
	}
	if objectRaw["ServiceId"] != nil {
		d.Set("service_id", objectRaw["ServiceId"])
	}
	if objectRaw["Status"] != nil {
		d.Set("status", objectRaw["Status"])
	}
	if objectRaw["DataSourceId"] != nil {
		d.Set("data_source_id", objectRaw["DataSourceId"])
	}

	parts := strings.Split(d.Id(), ":")
	d.Set("db_instance_id", parts[0])
	d.Set("data_source_id", parts[1])

	return nil
}

func resourceAliCloudGpdbStreamingDataSourceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	parts := strings.Split(d.Id(), ":")
	action := "ModifyStreamingDataSource"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["DataSourceId"] = parts[1]
	query["DBInstanceId"] = parts[0]
	request["RegionId"] = client.RegionId
	if d.HasChange("data_source_description") {
		update = true
		request["DataSourceDescription"] = d.Get("data_source_description")
	}

	if d.HasChange("data_source_config") {
		update = true
	}
	request["DataSourceConfig"] = d.Get("data_source_config")
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

	return resourceAliCloudGpdbStreamingDataSourceRead(d, meta)
}

func resourceAliCloudGpdbStreamingDataSourceDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeleteStreamingDataSource"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["DBInstanceId"] = parts[0]
	query["DataSourceId"] = parts[1]
	request["RegionId"] = client.RegionId

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
