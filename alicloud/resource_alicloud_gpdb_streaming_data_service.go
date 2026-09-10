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

func resourceAliCloudGpdbStreamingDataService() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudGpdbStreamingDataServiceCreate,
		Read:   resourceAliCloudGpdbStreamingDataServiceRead,
		Update: resourceAliCloudGpdbStreamingDataServiceUpdate,
		Delete: resourceAliCloudGpdbStreamingDataServiceDelete,
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
			"service_description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"service_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"service_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"service_spec": {
				Type:     schema.TypeString,
				Required: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudGpdbStreamingDataServiceCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateStreamingDataService"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("db_instance_id"); ok {
		request["DBInstanceId"] = v
	}
	request["RegionId"] = client.RegionId

	if v, ok := d.GetOk("service_description"); ok {
		request["ServiceDescription"] = v
	}
	request["ServiceSpec"] = d.Get("service_spec")
	request["ServiceName"] = d.Get("service_name")
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
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_gpdb_streaming_data_service", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", request["DBInstanceId"], response["ServiceId"]))

	gpdbServiceV2 := GpdbServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"Running", "running"}, d.Timeout(schema.TimeoutCreate), 10*time.Second, gpdbServiceV2.GpdbStreamingDataServiceStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudGpdbStreamingDataServiceRead(d, meta)
}

func resourceAliCloudGpdbStreamingDataServiceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gpdbServiceV2 := GpdbServiceV2{client}

	objectRaw, err := gpdbServiceV2.DescribeGpdbStreamingDataService(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_gpdb_streaming_data_service DescribeGpdbStreamingDataService Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("create_time", objectRaw["CreateTime"])
	d.Set("service_description", objectRaw["ServiceDescription"])
	d.Set("service_name", objectRaw["ServiceName"])
	d.Set("service_spec", objectRaw["ServiceSpec"])
	d.Set("status", objectRaw["Status"])
	d.Set("service_id", objectRaw["ServiceId"])

	parts := strings.Split(d.Id(), ":")
	d.Set("db_instance_id", parts[0])

	return nil
}

func resourceAliCloudGpdbStreamingDataServiceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	parts := strings.Split(d.Id(), ":")
	action := "ModifyStreamingDataService"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ServiceId"] = parts[1]
	request["DBInstanceId"] = parts[0]
	request["RegionId"] = client.RegionId
	if d.HasChange("service_description") {
		update = true
		request["ServiceDescription"] = d.Get("service_description")
	}

	if d.HasChange("service_spec") {
		update = true
	}
	request["ServiceSpec"] = d.Get("service_spec")
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
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		gpdbServiceV2 := GpdbServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"Running", "running"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, gpdbServiceV2.GpdbStreamingDataServiceStateRefreshFunc(d.Id(), "Status", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	return resourceAliCloudGpdbStreamingDataServiceRead(d, meta)
}

func resourceAliCloudGpdbStreamingDataServiceDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeleteStreamingDataService"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["ServiceId"] = parts[1]
	request["DBInstanceId"] = parts[0]
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
